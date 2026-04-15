// Package bing_news implements Bing news search.
package bing_news

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

// BingNews implements Bing news search engine.
type BingNews struct {
	client   *httpclient.Client
	name     string
	category engine.Category
	provider string
	priority float64
}

const (
	bingNewsURL = "https://www.bing.com/news/search"
)

// New creates a new Bing news engine.
func New(client *httpclient.Client) engine.SearchEngine {
	return &BingNews{
		client:   client,
		name:     "bing_news",
		category: engine.CategoryNews,
		provider: "bing",
		priority: 0.85,
	}
}

func (b *BingNews) Name() string         { return b.name }
func (b *BingNews) Category() engine.Category { return b.category }
func (b *BingNews) Provider() string     { return b.provider }
func (b *BingNews) Priority() float64    { return b.priority }

// Search performs Bing news search.
func (b *BingNews) Search(ctx context.Context, query string, opts engine.SearchOptions) ([]map[string]string, error) {
	payload := b.buildPayload(query, opts)
	apiURL := fmt.Sprintf("%s?%s", bingNewsURL, payload.Encode())

	resp, err := b.client.Do(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("search failed with status: %d", resp.StatusCode)
	}

	return b.extractResults(resp.Text)
}

func (b *BingNews) buildPayload(query string, opts engine.SearchOptions) url.Values {
	payload := url.Values{}
	payload.Set("q", query)
	payload.Set("format", "json")

	if opts.TimeLimit != "" {
		hours := map[string]string{
			"d": "24",
			"w": "168",
			"m": "720",
		}
		if h, ok := hours[opts.TimeLimit]; ok {
			payload.Set("filter", fmt.Sprintf("publishedTime=Past%s", h))
		}
	}

	return payload
}

func (b *BingNews) extractResults(htmlText string) ([]map[string]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlText))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var results []map[string]string

	doc.Find("div.news-card").Each(func(i int, s *goquery.Selection) {
		title := s.Find("h2, a.title").First().Text()
		url, _ := s.Find("a.title").First().Attr("href")
		body := s.Find("p, .snippet").Text()
		image := ""
		s.Find("img").Each(func(j int, img *goquery.Selection) {
			src, _ := img.Attr("src")
			if src != "" {
				image = src
			}
		})
		source := s.Find(".source, .provider").Text()
		date := s.Find(".time, .datetime").Text()

		if title != "" {
			results = append(results, map[string]string{
				"title":  utils.NormalizeText(title),
				"url":    utils.NormalizeURL(url),
				"body":   utils.NormalizeText(body),
				"image":  utils.NormalizeURL(image),
				"source": utils.NormalizeText(source),
				"date":   date,
			})
		}
	})

	if len(results) == 0 {
		return nil, fmt.Errorf("no news extracted")
	}

	return results[:min(len(results), 10)], nil
}
