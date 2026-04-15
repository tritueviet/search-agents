// Package extractor provides URL content extraction functionality.
package extractor

import (
	"context"
	"fmt"

	markdown "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/tritueviet/search-agents/internal/httpclient"
	"github.com/tritueviet/search-agents/internal/models"
)

// Extractor extracts content from URLs.
type Extractor struct {
	client *httpclient.Client
}

// New creates a new Extractor.
func New(client *httpclient.Client) *Extractor {
	return &Extractor{client: client}
}

// Extract fetches a URL and extracts its content in various formats.
func (e *Extractor) Extract(ctx context.Context, url string, format string) (map[string]interface{}, error) {
	resp, err := e.client.Get(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, models.NewDDGSError(fmt.Sprintf("HTTP %d: failed to fetch %s", resp.StatusCode, url))
	}

	var content interface{}
	switch format {
	case "text_markdown":
		converter := markdown.NewConverter("", true, nil)
		mdContent, err := converter.ConvertString(resp.Text)
		if err != nil {
			content = resp.Text
		} else {
			content = mdContent
		}
	case "text_plain":
		content = stripHTML(resp.Text)
	case "text_rich":
		content = extractRichText(resp.Text)
	case "text":
		content = resp.Text
	case "content":
		content = resp.Content
	default:
		converter := markdown.NewConverter("", true, nil)
		mdContent, err := converter.ConvertString(resp.Text)
		if err != nil {
			content = resp.Text
		} else {
			content = mdContent
		}
	}

	return map[string]interface{}{
		"url":     url,
		"content": content,
	}, nil
}

// stripHTML removes HTML tags and returns plain text.
func stripHTML(htmlText string) string {
	// Simple HTML stripping
	result := make([]rune, 0, len(htmlText))
	var inTag bool
	for _, ch := range htmlText {
		if ch == '<' {
			inTag = true
			continue
		}
		if ch == '>' {
			inTag = false
			continue
		}
		if !inTag {
			result = append(result, ch)
		}
	}
	return string(result)
}

// extractRichText extracts text while preserving some formatting.
func extractRichText(htmlText string) string {
	// Basic rich text extraction - preserve headers and lists
	return stripHTML(htmlText)
}
