// Package duckduckgo_videos implements DuckDuckGo videos search.
package duckduckgo_videos

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/tritueviet/search-agents/internal/engine"
	"github.com/tritueviet/search-agents/internal/httpclient"
	"github.com/tritueviet/search-agents/internal/utils"
)

// DuckduckgoVideos implements DuckDuckGo video search engine.
type DuckduckgoVideos struct {
	client   *httpclient.Client
	name     string
	category engine.Category
	provider string
	priority float64
}

const (
	ddgVideosURL = "https://duckduckgo.com/v.js"
	ddgMainURL   = "https://duckduckgo.com"
)

// New creates a new DuckDuckGo videos engine.
func New(client *httpclient.Client) engine.SearchEngine {
	return &DuckduckgoVideos{
		client:   client,
		name:     "duckduckgo_videos",
		category: engine.CategoryVideos,
		provider: "bing",
		priority: 0.8,
	}
}

func (d *DuckduckgoVideos) Name() string         { return d.name }
func (d *DuckduckgoVideos) Category() engine.Category { return d.category }
func (d *DuckduckgoVideos) Provider() string     { return d.provider }
func (d *DuckduckgoVideos) Priority() float64    { return d.priority }

// Search performs DuckDuckGo video search.
func (d *DuckduckgoVideos) Search(ctx context.Context, query string, opts engine.SearchOptions) ([]map[string]string, error) {
	vqd, err := d.getVQD(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get VQD: %w", err)
	}

	payload := d.buildPayload(query, opts, vqd)
	apiURL := fmt.Sprintf("%s?%s", ddgVideosURL, payload.Encode())

	resp, err := d.client.Do(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}

	if resp.StatusCode == 403 {
		return nil, fmt.Errorf("DuckDuckGo Videos API has been blocked (403 Forbidden). " +
			"Alternatives:\n" +
			"  1. Use YouTube Data API (requires API key)\n" +
			"  2. Use Google Custom Search API with video filter\n" +
			"  3. Use SerpAPI or similar third-party services")
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("search failed with status: %d", resp.StatusCode)
	}

	results, err := d.extractResults(resp.Text)
	if err != nil {
		return nil, fmt.Errorf("failed to extract results: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("DuckDuckGo Videos API returned empty results. " +
			"The API may be deprecated. Consider using:\n" +
			"  1. YouTube Data API\n" +
			"  2. SerpAPI Video Search")
	}

	return results, nil
}

// getVQD gets VQD token from DuckDuckGo.
func (d *DuckduckgoVideos) getVQD(ctx context.Context, query string) (string, error) {
	resp, err := d.client.Do(ctx, "GET", ddgMainURL+"?q="+url.QueryEscape(query), nil)
	if err != nil {
		return "", err
	}

	vqd := extractVQD(resp.Text)
	if vqd == "" {
		return "", fmt.Errorf("failed to extract VQD from response")
	}

	return vqd, nil
}

// extractVQD extracts VQD token from HTML with improved patterns.
func extractVQD(htmlText string) string {
	// Method 1: vqd="..."
	re1 := regexp.MustCompile(`vqd="([^"]+)"`)
	if match := re1.FindStringSubmatch(htmlText); len(match) > 1 {
		return match[1]
	}

	// Method 2: vqd='...'
	re2 := regexp.MustCompile(`vqd='([^']+)'`)
	if match := re2.FindStringSubmatch(htmlText); len(match) > 1 {
		return match[1]
	}

	// Method 3: In script tags
	re3 := regexp.MustCompile(`"vqd":"([^"]+)"`)
	if match := re3.FindStringSubmatch(htmlText); len(match) > 1 {
		return match[1]
	}

	return ""
}

func (d *DuckduckgoVideos) buildPayload(query string, opts engine.SearchOptions, vqd string) url.Values {
	payload := url.Values{}
	payload.Set("l", opts.Region)
	payload.Set("o", "json")
	payload.Set("q", query)
	payload.Set("vqd", vqd)

	safesearchMap := map[string]string{
		"on":       "1",
		"moderate": "-1",
		"off":      "-2",
	}
	payload.Set("p", safesearchMap[strings.ToLower(opts.SafeSearch)])

	filters := []string{}
	if opts.TimeLimit != "" {
		filters = append(filters, fmt.Sprintf("publishedAfter:%s", opts.TimeLimit))
	}
	payload.Set("f", strings.Join(filters, ","))

	if opts.Page > 1 {
		payload.Set("s", fmt.Sprintf("%d", (opts.Page-1)*60))
	}

	return payload
}

func (d *DuckduckgoVideos) extractResults(jsonText string) ([]map[string]string, error) {
	var result struct {
		Results []struct {
			Title       string            `json:"title"`
			Content     string            `json:"content"`
			Description string            `json:"description"`
			Duration    string            `json:"duration"`
			EmbedHTML   string            `json:"embed_html"`
			EmbedURL    string            `json:"embed_url"`
			ImageToken  string            `json:"image_token"`
			Images      map[string]string `json:"images"`
			Provider    string            `json:"provider"`
			Published   string            `json:"published"`
			Publisher   string            `json:"publisher"`
			Statistics  map[string]int    `json:"statistics"`
			Uploader    string            `json:"uploader"`
		} `json:"results"`
	}

	if err := json.Unmarshal([]byte(jsonText), &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	videos := make([]map[string]string, 0, len(result.Results))
	for _, r := range result.Results {
		viewCount := ""
		if r.Statistics != nil {
			viewCount = fmt.Sprintf("%d", r.Statistics["views"])
		}

		videos = append(videos, map[string]string{
			"title":       utils.NormalizeText(r.Title),
			"content":     utils.NormalizeURL(r.Content),
			"description": utils.NormalizeText(r.Description),
			"duration":    r.Duration,
			"embed_html":  r.EmbedHTML,
			"embed_url":   utils.NormalizeURL(r.EmbedURL),
			"image_token": r.ImageToken,
			"provider":    r.Provider,
			"published":   r.Published,
			"publisher":   utils.NormalizeText(r.Publisher),
			"uploader":    utils.NormalizeText(r.Uploader),
			"view_count":  viewCount,
		})
	}

	return videos, nil
}
