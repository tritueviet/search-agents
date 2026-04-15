// Package mcp implements the Model Context Protocol server.
package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/tritueviet/search-agents/internal/engine"
	"github.com/tritueviet/search-agents/internal/extractor"
	"github.com/tritueviet/search-agents/pkg/searchagents"
)

// Server represents the MCP server.
type Server struct {
	client  *searchagents.SearchAgents
	extract *extractor.Extractor
	reader  *bufio.Reader
	writer  io.Writer
}

// NewServer creates a new MCP server.
func NewServer(client *searchagents.SearchAgents, ext *extractor.Extractor) *Server {
	return &Server{
		client:  client,
		extract: ext,
		reader:  bufio.NewReader(os.Stdin),
		writer:  os.Stdout,
	}
}

// Request represents an MCP request.
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// Response represents an MCP response.
type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

// Error represents an MCP error.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Tool represents an MCP tool.
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// RunStdioAsync starts the MCP server over stdio.
func (s *Server) RunStdioAsync() error {
	fmt.Fprintln(os.Stderr, "MCP server started on stdio")

	for {
		line, err := s.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("failed to read input: %w", err)
		}

		var req Request
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			s.sendError(nil, -32700, "Parse error")
			continue
		}

		s.handleRequest(req)
	}
}

// handleRequest routes incoming requests.
func (s *Server) handleRequest(req Request) {
	switch req.Method {
	case "initialize":
		s.handleInitialize(req.ID)
	case "tools/list":
		s.handleToolsList(req.ID)
	case "tools/call":
		s.handleToolsCall(req.ID, req.Params)
	default:
		s.sendError(req.ID, -32601, fmt.Sprintf("Method not found: %s", req.Method))
	}
}

// handleInitialize responds to initialize request.
func (s *Server) handleInitialize(id interface{}) {
	result := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{},
		},
		"serverInfo": map[string]interface{}{
			"name":    "search-agents",
			"version": "0.1.0",
		},
	}
	s.sendResponse(id, result)
}

// handleToolsList returns available tools.
func (s *Server) handleToolsList(id interface{}) {
	tools := []Tool{
		{
			Name:        "search_text",
			Description: "Perform a web text search",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "Search query",
					},
					"max_results": map[string]interface{}{
						"type":        "integer",
						"description": "Maximum number of results",
						"default":     10,
					},
				},
				"required": []string{"query"},
			},
		},
		{
			Name:        "extract_content",
			Description: "Extract content from a URL",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"url": map[string]interface{}{
						"type":        "string",
						"description": "URL to extract content from",
					},
					"format": map[string]interface{}{
						"type":        "string",
						"description": "Output format (text_markdown, text_plain, text)",
						"default":     "text_markdown",
					},
				},
				"required": []string{"url"},
			},
		},
	}
	s.sendResponse(id, map[string]interface{}{"tools": tools})
}

// handleToolsCall executes a tool.
func (s *Server) handleToolsCall(id interface{}, params json.RawMessage) {
	var callParams struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	}
	if err := json.Unmarshal(params, &callParams); err != nil {
		s.sendError(id, -32602, "Invalid arguments")
		return
	}

	switch callParams.Name {
	case "search_text":
		s.callSearchText(id, callParams.Arguments)
	case "extract_content":
		s.callExtractContent(id, callParams.Arguments)
	default:
		s.sendError(id, -32601, fmt.Sprintf("Unknown tool: %s", callParams.Name))
	}
}

// callSearchText performs text search.
func (s *Server) callSearchText(id interface{}, args map[string]interface{}) {
	query, ok := args["query"].(string)
	if !ok || query == "" {
		s.sendError(id, -32602, "query is required")
		return
	}

	maxResults := 10
	if mr, ok := args["max_results"].(float64); ok {
		maxResults = int(mr)
	}

	opts := engine.SearchOptions{
		Region:     s.getStringArg(args, "region", "us-en"),
		SafeSearch: s.getStringArg(args, "safesearch", "moderate"),
		TimeLimit:  s.getStringArg(args, "timelimit", ""),
		Extra: map[string]string{
			"max_results": fmt.Sprintf("%d", maxResults),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	results, err := s.client.Text(ctx, query, opts)
	if err != nil {
		s.sendError(id, -32603, fmt.Sprintf("Search failed: %v", err))
		return
	}

	content, _ := json.Marshal(results)
	s.sendResponse(id, map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": string(content),
			},
		},
	})
}

// callExtractContent performs content extraction.
func (s *Server) callExtractContent(id interface{}, args map[string]interface{}) {
	url, ok := args["url"].(string)
	if !ok || url == "" {
		s.sendError(id, -32602, "url is required")
		return
	}

	format := s.getStringArg(args, "format", "text_markdown")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := s.extract.Extract(ctx, url, format)
	if err != nil {
		s.sendError(id, -32603, fmt.Sprintf("Extract failed: %v", err))
		return
	}

	content, _ := json.Marshal(result)
	s.sendResponse(id, map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": string(content),
			},
		},
	})
}

// getStringArg safely extracts a string argument.
func (s *Server) getStringArg(args map[string]interface{}, key, defaultVal string) string {
	if val, ok := args[key].(string); ok {
		return val
	}
	return defaultVal
}

// sendResponse sends a JSON-RPC response.
func (s *Server) sendResponse(id interface{}, result interface{}) {
	resp := Response{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
	s.sendJSON(resp)
}

// sendError sends a JSON-RPC error response.
func (s *Server) sendError(id interface{}, code int, message string) {
	resp := Response{
		JSONRPC: "2.0",
		ID:      id,
		Error: &Error{
			Code:    code,
			Message: message,
		},
	}
	s.sendJSON(resp)
}

// sendJSON marshals and sends a JSON object.
func (s *Server) sendJSON(v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal response: %v\n", err)
		return
	}
	fmt.Fprintf(s.writer, "%s\n", data)
}
