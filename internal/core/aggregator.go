// Package core provides the results aggregator for deduplication.
package core

import (
	"sort"
)

// ResultsAggregator aggregates and deduplicates search results.
type ResultsAggregator struct {
	cacheFields map[string]bool
	counter     map[string]int
	cache       map[string]map[string]string
}

// NewResultsAggregator creates a new results aggregator.
func NewResultsAggregator(cacheFields map[string]bool) *ResultsAggregator {
	return &ResultsAggregator{
		cacheFields: cacheFields,
		counter:     make(map[string]int),
		cache:       make(map[string]map[string]string),
	}
}

// Append adds a result to the aggregator.
func (r *ResultsAggregator) Append(item map[string]string) {
	key := r.getKey(item)
	if key == "" {
		return
	}

	// Store if not exists or if this has more content
	if existing, exists := r.cache[key]; !exists {
		r.cache[key] = item
	} else {
		// Keep the one with more content
		if len(item["body"]) > len(existing["body"]) {
			r.cache[key] = item
		}
	}

	r.counter[key]++
}

// Extract returns deduplicated results sorted by frequency.
func (r *ResultsAggregator) Extract() []map[string]string {
	type freqItem struct {
		key   string
		count int
	}

	var items []freqItem
	for key, count := range r.counter {
		items = append(items, freqItem{key, count})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].count > items[j].count
	})

	var results []map[string]string
	for _, item := range items {
		if result, exists := r.cache[item.key]; exists {
			results = append(results, result)
		}
	}

	return results
}

// getKey extracts a key from the item using cacheFields.
func (r *ResultsAggregator) getKey(item map[string]string) string {
	for field := range r.cacheFields {
		if val, exists := item[field]; exists && val != "" {
			return val
		}
	}
	return ""
}
