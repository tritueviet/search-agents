// Package engine provides search engine interfaces.
package engine

import (
	"context"
)

// Category represents a search category.
type Category string

const (
	CategoryText   Category = "text"
	CategoryImages Category = "images"
	CategoryVideos Category = "videos"
	CategoryNews   Category = "news"
	CategoryBooks  Category = "books"
)

// SearchEngine is the base interface for all search engines.
type SearchEngine interface {
	// Name returns the unique identifier for this engine.
	Name() string

	// Category returns the search category this engine supports.
	Category() Category

	// Provider returns the data source provider name.
	Provider() string

	// Priority returns the priority of this engine (higher = more preferred).
	Priority() float64

	// Search performs a search and returns results.
	Search(ctx context.Context, query string, opts SearchOptions) ([]map[string]string, error)
}

// SearchOptions contains options for a search request.
type SearchOptions struct {
	Region     string
	SafeSearch string
	TimeLimit  string
	Page       int
	Extra      map[string]string
}

// DefaultSearchOptions returns default search options.
func DefaultSearchOptions() SearchOptions {
	return SearchOptions{
		Region:     "us-en",
		SafeSearch: "moderate",
		Page:       1,
		Extra:      make(map[string]string),
	}
}
