// Package yahoo_news implements Yahoo news search.
package yahoo_news

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

// YahooNews implements Yahoo news search engine.
type YahooNews struct {
	client   *httpclient.Client
	name     string
	category engine.Category
	provider string
	priority float64
}

const (
	yahooNewsURL = "https://news.search.yahoo.com/search"
)

// New creates a new Yahoo news engine.
func New(client *httpclient.Client) engine.SearchEngine {
	return &YahooNews{
		client:   client,
		name:     "yahoo_news",
		category: engine.CategoryNews,
		provider: "yahoo",
		priority: 0.75,
	}
}

func (y *YahooNews) Name() string         { return y.name }
func (y *YahooNews) Category() engine.Category { return y.category }
func (y *YahooNews) Provider() string     { return y.provider }
func (y *YahooNews) Priority() float64    { return y.priority }

// Search performs Yahoo news search.
func (y *YahooNews) Search(ctx context.Context, query string, opts engine.SearchOptions) ([]map[string]string, error) {
	payload := y.buildPayload(query, opts)
	apiURL := fmt.Sprintf("%s?%s", yahooNewsURL, payload.Encode())

	resp, err := y.client.Do(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("search failed with status: %d", resp.StatusCode)
	}

	return y.extractResults(resp.Text)
}

func (y *YahooNews) buildPayload(query string, opts engine.SearchOptions) url.Values {
	payload := url.Values{}
	payload.Set("p", query)

	if opts.TimeLimit != "" {
		hours := map[string]string{
			"d": "24",
			"w": "168",
			"m": "720",
		}
		if h, ok := hours[opts.TimeLimit]; ok {
			payload.Set("tf", fmt.Sprintf("Past%s", h))
		}
	}

	return payload
}

func (y *YahooNews) extractResults(htmlText string) ([]map[string]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlText))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var results []map[string]string

	doc.Find("div.NewsArticle").Each(func(i int, s *goquery.Selection) {
		title := s.Find("h3, a.title").First().Text()
		url, _ := s.Find("a.title").First().Attr("href")
		body := s.Find("p, .excerpt").Text()
		image := ""
		s.Find("img").Each(func(j int, img *goquery.Selection) {
			src, _ := img.Attr("src")
			if src != "" {
				image = src
			}
		})
		source := s.Find(".provider, .source").Text()
		date := s.Find(".publish-time, .time").Text()

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
