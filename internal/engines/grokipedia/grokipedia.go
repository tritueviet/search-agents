// Package grokipedia implements Grokipedia search engine.
package grokipedia

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

// Grokipedia implements Grokipedia search engine.
type Grokipedia struct {
	client   *httpclient.Client
	name     string
	category engine.Category
	provider string
	priority float64
}

const (
	grokipediaURL = "https://www.grokipedia.com/search"
)

// New creates a new Grokipedia engine.
func New(client *httpclient.Client) engine.SearchEngine {
	return &Grokipedia{
		client:   client,
		name:     "grokipedia",
		category: engine.CategoryText,
		provider: "grokipedia",
		priority: 0.7,
	}
}

func (g *Grokipedia) Name() string         { return g.name }
func (g *Grokipedia) Category() engine.Category { return g.category }
func (g *Grokipedia) Provider() string     { return g.provider }
func (g *Grokipedia) Priority() float64    { return g.priority }

// Search performs Grokipedia search.
func (g *Grokipedia) Search(ctx context.Context, query string, opts engine.SearchOptions) ([]map[string]string, error) {
	params := url.Values{}
	params.Set("q", query)

	apiURL := fmt.Sprintf("%s?%s", grokipediaURL, params.Encode())

	resp, err := g.client.Do(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("search failed with status: %d", resp.StatusCode)
	}

	return g.extractResults(resp.Text)
}

func (g *Grokipedia) extractResults(htmlText string) ([]map[string]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlText))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var results []map[string]string

	doc.Find("div.result, article, div.search-result").Each(func(i int, s *goquery.Selection) {
		title := s.Find("h2, h3, a.title").First().Text()
		href, _ := s.Find("a").First().Attr("href")
		body := s.Find("p, .snippet, .description").Text()

		if title != "" {
			results = append(results, map[string]string{
				"title": utils.NormalizeText(title),
				"href":  utils.NormalizeURL(href),
				"body":  utils.NormalizeText(body),
			})
		}
	})

	if len(results) == 0 {
		return nil, fmt.Errorf("no results extracted")
	}

	return results[:min(len(results), 10)], nil
}
