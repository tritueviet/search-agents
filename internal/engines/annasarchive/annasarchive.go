// Package annasarchive implements Anna's Archive books search.
package annasarchive

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tritueviet/search-agents/internal/engine"
	"github.com/tritueviet/search-agents/internal/httpclient"
	"github.com/tritueviet/search-agents/internal/utils"
)

// AnnasArchive implements Anna's Archive books search engine.
type AnnasArchive struct {
	client   *httpclient.Client
	name     string
	category engine.Category
	provider string
	priority float64
}

const (
	annasArchiveURL = "https://annas-archive.li/search"
	userAgent       = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
)

// New creates a new Anna's Archive engine.
func New(client *httpclient.Client) engine.SearchEngine {
	eng := &AnnasArchive{
		client:   client,
		name:     "annasarchive",
		category: engine.CategoryBooks,
		provider: "annas_archive",
		priority: 0.8,
	}

	// Set browser-like headers
	eng.client.SetHeader("User-Agent", userAgent)
	eng.client.SetHeader("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	eng.client.SetHeader("Accept-Language", "en-US,en;q=0.5")

	return eng
}

func (a *AnnasArchive) Name() string         { return a.name }
func (a *AnnasArchive) Category() engine.Category { return a.category }
func (a *AnnasArchive) Provider() string     { return a.provider }
func (a *AnnasArchive) Priority() float64    { return a.priority }

// Search performs Anna's Archive books search.
func (a *AnnasArchive) Search(ctx context.Context, query string, opts engine.SearchOptions) ([]map[string]string, error) {
	payload := a.buildPayload(query, opts)
	apiURL := fmt.Sprintf("%s?%s", annasArchiveURL, payload.Encode())

	resp, err := a.client.Do(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("search failed with status: %d", resp.StatusCode)
	}

	return a.extractResults(resp.Text)
}

func (a *AnnasArchive) buildPayload(query string, opts engine.SearchOptions) url.Values {
	payload := url.Values{}
	payload.Set("q", query)

	if opts.Page > 1 {
		payload.Set("page", fmt.Sprintf("%d", opts.Page))
	}

	return payload
}

func (a *AnnasArchive) extractResults(htmlText string) ([]map[string]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlText))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var books []map[string]string

	// Try multiple selector strategies for Anna's Archive
	selectors := []string{
		"a[href*='/md5/']",
		"a[href*='/book/']",
		"main a[href]",
		"div a[href]",
	}

	var items *goquery.Selection
	for _, selector := range selectors {
		items = doc.Find(selector)
		if items.Length() > 0 {
			break
		}
	}

	if items.Length() == 0 {
		return nil, fmt.Errorf("no book links found with any selector")
	}

	items.Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Text())
		href, _ := s.Attr("href")

		// Skip empty or very short titles
		if len(title) < 10 {
			return
		}

		// Skip navigation/menu links
		if strings.Contains(href, "/faq") || 
		   strings.Contains(href, "/donate") ||
		   strings.Contains(href, "/login") ||
		   strings.Contains(href, "/register") {
			return
		}

		// Normalize URL
		if href != "" && !strings.HasPrefix(href, "http") {
			href = "https://annas-archive.li" + href
		}

		books = append(books, map[string]string{
			"title":     utils.NormalizeText(title),
			"author":    "",
			"publisher": "",
			"info":      "",
			"url":       utils.NormalizeURL(href),
			"thumbnail": "",
		})
	})

	// Deduplicate by URL
	seen := make(map[string]bool)
	var unique []map[string]string
	for _, book := range books {
		if !seen[book["url"]] {
			seen[book["url"]] = true
			unique = append(unique, book)
		}
	}

	if len(unique) == 0 {
		return nil, fmt.Errorf("no valid books found")
	}

	// Return max 10 results
	max := 10
	if len(unique) < max {
		max = len(unique)
	}

	return unique[:max], nil
}
