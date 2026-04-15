// Package openlibrary implements Open Library books search.
package openlibrary

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/tritueviet/search-agents/internal/engine"
	"github.com/tritueviet/search-agents/internal/httpclient"
	"github.com/tritueviet/search-agents/internal/utils"
)

// OpenLibrary implements Open Library books search engine.
type OpenLibrary struct {
	client   *httpclient.Client
	name     string
	category engine.Category
	provider string
	priority float64
}

const (
	openLibraryURL = "https://openlibrary.org/search.json"
)

// New creates a new Open Library engine.
func New(client *httpclient.Client) engine.SearchEngine {
	return &OpenLibrary{
		client:   client,
		name:     "openlibrary",
		category: engine.CategoryBooks,
		provider: "openlibrary",
		priority: 0.9,
	}
}

func (o *OpenLibrary) Name() string         { return o.name }
func (o *OpenLibrary) Category() engine.Category { return o.category }
func (o *OpenLibrary) Provider() string     { return o.provider }
func (o *OpenLibrary) Priority() float64    { return o.priority }

// Search performs Open Library books search.
func (o *OpenLibrary) Search(ctx context.Context, query string, opts engine.SearchOptions) ([]map[string]string, error) {
	payload := o.buildPayload(query, opts)
	apiURL := fmt.Sprintf("%s?%s", openLibraryURL, payload.Encode())

	resp, err := o.client.Do(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("search failed with status: %d", resp.StatusCode)
	}

	return o.extractResults(resp.Text)
}

func (o *OpenLibrary) buildPayload(query string, opts engine.SearchOptions) url.Values {
	payload := url.Values{}
	payload.Set("q", query)
	payload.Set("limit", "10")

	if opts.Page > 1 {
		payload.Set("offset", fmt.Sprintf("%d", (opts.Page-1)*10))
	}

	return payload
}

func (o *OpenLibrary) extractResults(jsonText string) ([]map[string]string, error) {
	var result struct {
		NumFound int `json:"numFound"`
		Docs     []struct {
			Title         string   `json:"title"`
			AuthorName    []string `json:"author_name"`
			Publisher     []string `json:"publisher"`
			PublishDate   []string `json:"publish_date"`
			ISBN          []string `json:"isbn"`
			Language      []string `json:"language"`
			Subject       []string `json:"subject"`
			Key           string   `json:"key"`
			FirstSentence []string `json:"first_sentence"`
		} `json:"docs"`
	}

	if err := json.Unmarshal([]byte(jsonText), &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if result.NumFound == 0 {
		return nil, fmt.Errorf("no books found")
	}

	books := make([]map[string]string, 0, len(result.Docs))
	for _, doc := range result.Docs {
		author := ""
		if len(doc.AuthorName) > 0 {
			author = doc.AuthorName[0]
		}

		publisher := ""
		if len(doc.Publisher) > 0 {
			publisher = doc.Publisher[0]
		}

		publishDate := ""
		if len(doc.PublishDate) > 0 {
			publishDate = doc.PublishDate[0]
		}

		description := ""
		if len(doc.FirstSentence) > 0 {
			description = doc.FirstSentence[0]
		}

		subjects := ""
		if len(doc.Subject) > 0 {
			max := 3
			if len(doc.Subject) < max {
				max = len(doc.Subject)
			}
			subjects = doc.Subject[0]
			for i := 1; i < max; i++ {
				subjects += ", " + doc.Subject[i]
			}
		}

		books = append(books, map[string]string{
			"title":       utils.NormalizeText(doc.Title),
			"author":      utils.NormalizeText(author),
			"publisher":   utils.NormalizeText(publisher),
			"info":        fmt.Sprintf("%s. %s. %s", publishDate, subjects, ""),
			"url":         fmt.Sprintf("https://openlibrary.org%s", doc.Key),
			"thumbnail":   fmt.Sprintf("https://covers.openlibrary.org/b/olid/%s-M.jpg", doc.Key),
			"description": utils.NormalizeText(description),
		})
	}

	if len(books) == 0 {
		return nil, fmt.Errorf("no books extracted")
	}

	return books, nil
}
