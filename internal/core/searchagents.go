// Package core provides the SearchAgents orchestrator for metasearch.
package core

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/tritueviet/search-agents/internal/engine"
	"github.com/tritueviet/search-agents/internal/httpclient"
	"github.com/tritueviet/search-agents/internal/models"
	"github.com/tritueviet/search-agents/internal/register"
)

// SearchAgents is the main orchestrator for metasearch operations.
type SearchAgents struct {
	proxy   string
	timeout int
	verify  bool
	registry *engine.Registry
	client  *httpclient.Client
	engines map[engine.Category][]engine.SearchEngine
}

// Options contains configuration for SearchAgents.
type Options struct {
	Proxy   string
	Timeout int
	Verify  bool
}

// New creates a new SearchAgents instance.
func New(opts Options) (*SearchAgents, error) {
	clientOpts := httpclient.Options{
		Proxy:   opts.Proxy,
		Verify:  opts.Verify,
	}
	if opts.Timeout > 0 {
		clientOpts.Timeout = 0 // Will be set per-request
	}

	client, err := httpclient.NewClient(clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP client: %w", err)
	}

	sa := &SearchAgents{
		proxy:    opts.Proxy,
		timeout:  opts.Timeout,
		verify:   opts.Verify,
		registry: engine.NewRegistry(),
		client:   client,
		engines:  make(map[engine.Category][]engine.SearchEngine),
	}

	// Register engines
	sa.registerEngines()

	return sa, nil
}

// registerEngines registers all available search engines.
func (sa *SearchAgents) registerEngines() {
	register.DefaultEngines(sa.client, sa.registry)
}

// Text performs a text search.
func (sa *SearchAgents) Text(ctx context.Context, query string, opts engine.SearchOptions) ([]map[string]string, error) {
	return sa.search(ctx, "text", query, opts)
}

// Images performs an image search.
func (sa *SearchAgents) Images(ctx context.Context, query string, opts engine.SearchOptions) ([]map[string]string, error) {
	return sa.search(ctx, "images", query, opts)
}

// Videos performs a video search.
func (sa *SearchAgents) Videos(ctx context.Context, query string, opts engine.SearchOptions) ([]map[string]string, error) {
	return sa.search(ctx, "videos", query, opts)
}

// News performs a news search.
func (sa *SearchAgents) News(ctx context.Context, query string, opts engine.SearchOptions) ([]map[string]string, error) {
	return sa.search(ctx, "news", query, opts)
}

// Books performs a books search.
func (sa *SearchAgents) Books(ctx context.Context, query string, opts engine.SearchOptions) ([]map[string]string, error) {
	return sa.search(ctx, "books", query, opts)
}

// search performs a search across engines in the given category.
func (sa *SearchAgents) search(ctx context.Context, category string, query string, opts engine.SearchOptions) ([]map[string]string, error) {
	if query == "" {
		return nil, models.NewDDGSError("query is mandatory")
	}

	// Get engines for category
	cat := engine.Category(category)
	engines := sa.registry.GetEngines(cat, []string{"auto"}, sa.client)
	if len(engines) == 0 {
		return nil, models.NewDDGSError(fmt.Sprintf("no engines found for category: %s", category))
	}

	// Sort engines by priority
	sort.Slice(engines, func(i, j int) bool {
		return engines[i].Priority() > engines[j].Priority()
	})

	// Perform parallel search
	return sa.parallelSearch(ctx, engines, query, opts)
}

// parallelSearch performs searches across multiple engines in parallel.
func (sa *SearchAgents) parallelSearch(ctx context.Context, engines []engine.SearchEngine, query string, opts engine.SearchOptions) ([]map[string]string, error) {
	maxResults := 10
	if opts.Extra["max_results"] != "" {
		fmt.Sscanf(opts.Extra["max_results"], "%d", &maxResults)
	}

	// Calculate max workers
	lenUniqueProviders := len(getUniqueProviders(engines))
	maxWorkers := int(math.Min(float64(lenUniqueProviders), math.Ceil(float64(maxResults)/10)+1))
	if maxWorkers < 1 {
		maxWorkers = 1
	}

	// Use goroutines for parallel search
	type searchResult struct {
		results []map[string]string
		err     error
		engine  string
	}

	resultChan := make(chan searchResult, len(engines))
	var wg sync.WaitGroup

	// Limit concurrent searches
	sem := make(chan struct{}, maxWorkers)

	for _, eng := range engines {
		wg.Add(1)
		go func(e engine.SearchEngine) {
			defer wg.Done()

			sem <- struct{}{}        // Acquire
			defer func() { <-sem }() // Release

			timeout := 10 * time.Second
			if sa.timeout > 0 {
				timeout = time.Duration(sa.timeout) * time.Second
			}
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			results, err := e.Search(ctx, query, opts)
			resultChan <- searchResult{
				results: results,
				err:     err,
				engine:  e.Name(),
			}
		}(eng)
	}

	// Close channel when all goroutines complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Aggregate results
	aggregator := NewResultsAggregator(map[string]bool{
		"href":  true,
		"url":   true,
		"image": true,
	})

	var errors []string
	for result := range resultChan {
		if result.err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", result.engine, result.err))
			continue
		}

		if len(result.results) > 0 {
			for _, item := range result.results {
				aggregator.Append(item)
			}
		}

		// Check if we have enough results
		if len(aggregator.Extract()) >= maxResults {
			break
		}
	}

	results := aggregator.Extract()
	if len(results) > 0 {
		if len(results) > maxResults {
			results = results[:maxResults]
		}
		return results, nil
	}

	// No results - provide detailed error
	errMsg := "no results found"
	if len(errors) > 0 {
		errMsg = fmt.Sprintf("no results found. Engine errors:\n  - %s", strings.Join(errors, "\n  - "))
	}
	return nil, models.NewDDGSError(errMsg)
}

// getUniqueProviders returns unique provider names from engines.
func getUniqueProviders(engines []engine.SearchEngine) map[string]bool {
	providers := make(map[string]bool)
	for _, e := range engines {
		providers[e.Provider()] = true
	}
	return providers
}
