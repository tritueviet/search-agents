// Package mojeek implements Mojeek search engine.
package mojeek

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

// Mojeek implements Mojeek search engine.
type Mojeek struct {
	client   *httpclient.Client
	name     string
	category engine.Category
	provider string
	priority float64
}

const (
	mojeekURL = "https://www.mojeek.com/search"
)

// New creates a new Mojeek engine.
func New(client *httpclient.Client) engine.SearchEngine {
	return &Mojeek{
		client:   client,
		name:     "mojeek",
		category: engine.CategoryText,
		provider: "mojeek",
		priority: 0.65,
	}
}

func (m *Mojeek) Name() string         { return m.name }
func (m *Mojeek) Category() engine.Category { return m.category }
func (m *Mojeek) Provider() string     { return m.provider }
func (m *Mojeek) Priority() float64    { return m.priority }

// Search performs Mojeek search.
func (m *Mojeek) Search(ctx context.Context, query string, opts engine.SearchOptions) ([]map[string]string, error) {
	payload := url.Values{}
	payload.Set("q", query)
	payload.Set("t", "10")

	if opts.Page > 1 {
		payload.Set("s", fmt.Sprintf("%d", (opts.Page-1)*10))
	}

	apiURL := fmt.Sprintf("%s?%s", mojeekURL, payload.Encode())

	resp, err := m.client.Do(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("search failed with status: %d", resp.StatusCode)
	}

	return m.extractResults(resp.Text)
}

func (m *Mojeek) extractResults(htmlText string) ([]map[string]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlText))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var results []map[string]string

	doc.Find("div.result, ul.results li").Each(func(i int, s *goquery.Selection) {
		title := s.Find("h2 a, a.title").First().Text()
		href, _ := s.Find("h2 a, a.title").First().Attr("href")
		body := s.Find("p.s, div.result-snippet").Text()

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
