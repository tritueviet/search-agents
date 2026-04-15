// Package api provides REST API handlers.
package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tritueviet/search-agents/internal/engine"
	"github.com/tritueviet/search-agents/internal/extractor"
	"github.com/tritueviet/search-agents/pkg/searchagents"
)

// Server represents the API server.
type Server struct {
	router  *gin.Engine
	client  *searchagents.SearchAgents
	extract *extractor.Extractor
}

// NewServer creates a new API server.
func NewServer(client *searchagents.SearchAgents, ext *extractor.Extractor) *Server {
	router := gin.Default()

	s := &Server{
		router:  router,
		client:  client,
		extract: ext,
	}

	s.setupRoutes()

	return s
}

// setupRoutes configures API routes.
func (s *Server) setupRoutes() {
	// Health check
	s.router.GET("/health", s.healthCheck)

	// Search endpoints
	search := s.router.Group("/search")
	{
		search.GET("/text", s.textSearch)
		search.POST("/text", s.textSearch)
		search.GET("/images", s.imagesSearch)
		search.POST("/images", s.imagesSearch)
		search.GET("/news", s.newsSearch)
		search.POST("/news", s.newsSearch)
		search.GET("/videos", s.videosSearch)
		search.POST("/videos", s.videosSearch)
		search.GET("/books", s.booksSearch)
		search.POST("/books", s.booksSearch)
	}

	// Extract endpoint
	s.router.GET("/extract", s.extractContent)
	s.router.POST("/extract", s.extractContent)

	// Swagger documentation (placeholder)
	s.router.GET("/docs", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "API Documentation",
			"endpoints": map[string][]string{
				"GET/POST": {"/search/text", "/search/images", "/search/news", "/search/videos", "/search/books", "/extract"},
				"GET":      {"/health", "/docs"},
			},
		})
	})
}

// healthCheck returns server health status.
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "search-agents",
		"time":    time.Now().UTC().Format(time.RFC3339),
	})
}

// textSearch handles text search requests.
func (s *Server) textSearch(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	opts := s.parseSearchOptions(c)
	doExtract := c.DefaultQuery("extract", "false") == "true"
	extractFormat := c.DefaultQuery("extract_format", "text_markdown")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	results, err := s.client.Text(ctx, query, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Extract content from result URLs if requested
	if doExtract && len(results) > 0 && s.extract != nil {
		extractedResults := make([]map[string]interface{}, 0, len(results))
		for i, result := range results {
			href := result["href"]
			if href == "" {
				continue
			}

			extracted, err := s.extract.Extract(ctx, href, extractFormat)
			if err != nil {
				// Skip failed extractions
				continue
			}

			extractedResult := make(map[string]interface{})
			for k, v := range result {
				extractedResult[k] = v
			}
			extractedResult["index"] = i + 1
			extractedResult["extracted_content"] = extracted["content"]
			extractedResult["extract_format"] = extractFormat
			extractedResults = append(extractedResults, extractedResult)
		}

		c.JSON(http.StatusOK, gin.H{
			"query":          query,
			"extract":        true,
			"extract_format": extractFormat,
			"results":        extractedResults,
			"count":          len(extractedResults),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"query":   query,
		"results": results,
		"count":   len(results),
	})
}

// imagesSearch handles image search requests.
func (s *Server) imagesSearch(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "images search not yet implemented"})
}

// newsSearch handles news search requests.
func (s *Server) newsSearch(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "news search not yet implemented"})
}

// videosSearch handles video search requests.
func (s *Server) videosSearch(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "videos search not yet implemented"})
}

// booksSearch handles book search requests.
func (s *Server) booksSearch(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "books search not yet implemented"})
}

// extractContent handles content extraction requests.
func (s *Server) extractContent(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'url' is required"})
		return
	}

	format := c.DefaultQuery("format", "text_markdown")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	result, err := s.extract.Extract(ctx, url, format)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// parseSearchOptions extracts search options from request.
func (s *Server) parseSearchOptions(c *gin.Context) engine.SearchOptions {
	return engine.SearchOptions{
		Region:     c.DefaultQuery("region", "us-en"),
		SafeSearch: c.DefaultQuery("safesearch", "moderate"),
		TimeLimit:  c.Query("timelimit"),
		Page:       s.parseIntQuery(c, "page", 1),
		Extra: map[string]string{
			"max_results": c.DefaultQuery("max_results", "10"),
		},
	}
}

// parseIntQuery parses an integer query parameter with a default value.
func (s *Server) parseIntQuery(c *gin.Context, key string, defaultVal int) int {
	val := c.Query(key)
	if val == "" {
		return defaultVal
	}
	var result int
	if _, err := fmt.Sscanf(val, "%d", &result); err != nil {
		return defaultVal
	}
	return result
}

// Run starts the API server.
func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
