// Package bing_images implements Bing images search.
package bing_images

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tritueviet/search-agents/internal/engine"
	"github.com/tritueviet/search-agents/internal/httpclient"
	"github.com/tritueviet/search-agents/internal/utils"
)

// BingImages implements Bing images search engine.
type BingImages struct {
	client   *httpclient.Client
	name     string
	category engine.Category
	provider string
	priority float64
}

const (
	bingImagesURL = "https://www.bing.com/images/async"
)

// New creates a new Bing images engine.
func New(client *httpclient.Client) engine.SearchEngine {
	return &BingImages{
		client:   client,
		name:     "bing_images",
		category: engine.CategoryImages,
		provider: "bing",
		priority: 0.9,
	}
}

func (b *BingImages) Name() string         { return b.name }
func (b *BingImages) Category() engine.Category { return b.category }
func (b *BingImages) Provider() string     { return b.provider }
func (b *BingImages) Priority() float64    { return b.priority }

// Search performs Bing images search.
func (b *BingImages) Search(ctx context.Context, query string, opts engine.SearchOptions) ([]map[string]string, error) {
	payload := b.buildPayload(query, opts)
	apiURL := fmt.Sprintf("%s?%s", bingImagesURL, payload.Encode())

	resp, err := b.client.Do(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("search failed with status: %d", resp.StatusCode)
	}

	return b.extractResults(resp.Text)
}

func (b *BingImages) buildPayload(query string, opts engine.SearchOptions) url.Values {
	payload := url.Values{}
	payload.Set("q", query)
	payload.Set("async", "1")
	payload.Set("first", "1")
	payload.Set("count", "35")

	if opts.TimeLimit != "" {
		minutes := map[string]string{
			"d": "1440",
			"w": "10080",
			"m": "44640",
			"y": "525600",
		}
		if mins, ok := minutes[opts.TimeLimit]; ok {
			payload.Set("qft", fmt.Sprintf("filterui:age-lt%s", mins))
		}
	}

	return payload
}

func (b *BingImages) extractResults(htmlText string) ([]map[string]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlText))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var results []map[string]string

	doc.Find("div.imgpt").Each(func(i int, s *goquery.Selection) {
		iusc := s.Find("a.iusc")
		if iusc.Length() == 0 {
			return
		}

		metadata, exists := iusc.Attr("m")
		if !exists {
			return
		}

		var meta struct {
			Title    string `json:"t"`
			ImageURL string `json:"murl"`
			ThumbURL string `json:"turl"`
			PageURL  string `json:"purl"`
		}
		if err := json.Unmarshal([]byte(metadata), &meta); err != nil {
			return
		}

		width := ""
		height := ""
		s.Find("span.nowrap").Each(func(j int, span *goquery.Selection) {
			dim := span.Text()
			parts := strings.SplitN(strings.ReplaceAll(dim, "×", "x"), "x", 2)
			if len(parts) == 2 {
				width = strings.TrimSpace(parts[0])
				height = strings.TrimSpace(parts[1])
			}
		})

		source := ""
		s.Find("div.lnkw a").Each(func(j int, a *goquery.Selection) {
			source = a.Text()
		})

		if meta.Title != "" || meta.ImageURL != "" {
			results = append(results, map[string]string{
				"title":     utils.NormalizeText(meta.Title),
				"image":     utils.NormalizeURL(meta.ImageURL),
				"thumbnail": utils.NormalizeURL(meta.ThumbURL),
				"url":       utils.NormalizeURL(meta.PageURL),
				"width":     width,
				"height":    height,
				"source":    utils.NormalizeText(source),
			})
		}
	})

	if len(results) == 0 {
		return nil, fmt.Errorf("no images extracted")
	}

	return results[:min(len(results), 10)], nil
}
