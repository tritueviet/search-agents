// Package engine provides search engine registry.
package engine

import (
	"github.com/tritueviet/search-agents/internal/httpclient"
)

// Registry holds all registered search engines.
type Registry struct {
	engines map[Category]map[string]func(client *httpclient.Client) SearchEngine
}

// NewRegistry creates a new engine registry.
func NewRegistry() *Registry {
	return &Registry{
		engines: make(map[Category]map[string]func(client *httpclient.Client) SearchEngine),
	}
}

// Register registers a search engine factory function.
func (r *Registry) Register(category Category, name string, factory func(client *httpclient.Client) SearchEngine) {
	if r.engines[category] == nil {
		r.engines[category] = make(map[string]func(client *httpclient.Client) SearchEngine)
	}
	r.engines[category][name] = factory
}

// GetEngines returns a list of engines for a category and backend selection.
func (r *Registry) GetEngines(category Category, backends []string, client *httpclient.Client) []SearchEngine {
	categoryEngines, exists := r.engines[category]
	if !exists {
		return nil
	}

	var result []SearchEngine

	// If "auto" or "all" is specified, return all engines
	if len(backends) == 0 || contains(backends, "auto") || contains(backends, "all") {
		for name, factory := range categoryEngines {
			_ = name
			result = append(result, factory(client))
		}
		return result
	}

	// Return specific engines
	for _, backend := range backends {
		if factory, exists := categoryEngines[backend]; exists {
			result = append(result, factory(client))
		}
	}

	return result
}

// Contains checks if a slice contains a string.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
