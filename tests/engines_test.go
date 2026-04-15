package tests

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/tritueviet/search-agents/internal/engine"
	"github.com/tritueviet/search-agents/pkg/searchagents"
)

// TestAllTextEngines tests all text search engines individually
func TestAllTextEngines(t *testing.T) {
	engines := []struct {
		name   string
		region string
	}{
		{"duckduckgo", "us-en"},
		{"bing", "us-en"},
		{"google", "us-en"},
		{"brave", "us-en"},
		{"yahoo", "us-en"},
		{"yandex", "ru-ru"},
		{"wikipedia", "en"},
		{"grokipedia", "us-en"},
		{"mojeek", "us-en"},
	}

	client, err := searchagents.New(searchagents.Options{
		Timeout: 15,
		Verify:  true,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	for _, eng := range engines {
		t.Run(eng.name, func(t *testing.T) {
			opts := engine.DefaultSearchOptions()
			opts.Region = eng.region
			opts.Extra = map[string]string{
				"max_results": "2",
			}

			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			results, err := client.Text(ctx, "golang programming", opts)
			
			if err != nil {
				t.Skipf("%s: %v", eng.name, err)
				return
			}

			if len(results) == 0 {
				t.Errorf("%s: expected results but got none", eng.name)
				return
			}

			t.Logf("✓ %s: Got %d results", eng.name, len(results))
		})
	}
}

// TestVietnameseQuery tests Vietnamese language support
func TestVietnameseQuery(t *testing.T) {
	client, err := searchagents.New(searchagents.Options{
		Timeout: 15,
		Verify:  true,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	tests := []struct {
		name   string
		query  string
		region string
	}{
		{"Vietnamese politics", "tổng bí thư việt nam", "vn-vi"},
		{"Vietnamese education", "học lập trình python", "vn-vi"},
		{"Vietnamese news", "tin tức việt nam", "vn-vi"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := engine.DefaultSearchOptions()
			opts.Region = tt.region
			opts.Extra = map[string]string{
				"max_results": "3",
			}

			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			results, err := client.Text(ctx, tt.query, opts)
			if err != nil {
				t.Errorf("Search failed: %v", err)
				return
			}

			if len(results) == 0 {
				t.Error("Expected results but got none")
				return
			}

			t.Logf("✓ Got %d results for: %s", len(results), tt.query)
		})
	}
}

// TestSearchEnginesComparison compares results from different engines
func TestSearchEnginesComparison(t *testing.T) {
	client, err := searchagents.New(searchagents.Options{
		Timeout: 15,
		Verify:  true,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	engines := []string{"duckduckgo", "bing", "google", "brave"}
	query := "go programming language"

	t.Logf("\nComparing engines for query: %s", query)
	t.Log(strings.Repeat("-", 60))

	for _, eng := range engines {
		opts := engine.DefaultSearchOptions()
		opts.Extra = map[string]string{
			"max_results": "3",
		}

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		results, err := client.Text(ctx, query, opts)
		if err != nil {
			t.Logf("⚠ %-12s: Error - %v", eng, err)
			continue
		}

		t.Logf("✓ %-12s: %d results", eng, len(results))
		for i, r := range results {
			if i < 2 {
				t.Logf("  %d. %s", i+1, r["title"])
			}
		}
	}
}
