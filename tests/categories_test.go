package tests

import (
	"context"
	"testing"
	"time"

	"github.com/tritueviet/search-agents/internal/engine"
	"github.com/tritueviet/search-agents/pkg/searchagents"
)

// TestImagesSearch tests Bing Images engine
func TestImagesSearch(t *testing.T) {
	client, err := searchagents.New(searchagents.Options{
		Timeout: 15,
		Verify:  true,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	opts := engine.DefaultSearchOptions()
	opts.Extra = map[string]string{"max_results": "3"}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	results, err := client.Images(ctx, "butterfly", opts)
	if err != nil {
		t.Skipf("Images search failed: %v", err)
		return
	}

	if len(results) == 0 {
		t.Error("Expected images but got none")
		return
	}

	t.Logf("✓ Got %d images", len(results))
}

// TestNewsSearch tests news engines
func TestNewsSearch(t *testing.T) {
	client, err := searchagents.New(searchagents.Options{
		Timeout: 15,
		Verify:  true,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	opts := engine.DefaultSearchOptions()
	opts.Extra = map[string]string{"max_results": "3"}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	results, err := client.News(ctx, "AI technology", opts)
	if err != nil {
		t.Skipf("News search failed: %v", err)
		return
	}

	if len(results) == 0 {
		t.Error("Expected news but got none")
		return
	}

	t.Logf("✓ Got %d news articles", len(results))
}

// TestAllCategoriesSummary provides summary of all categories
func TestAllCategoriesSummary(t *testing.T) {
	client, err := searchagents.New(searchagents.Options{
		Timeout: 15,
		Verify:  true,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	query := "golang"
	opts := engine.DefaultSearchOptions()
	opts.Extra = map[string]string{"max_results": "2"}

	t.Log("\n============================================================")
	t.Log("CATEGORY TEST SUMMARY")
	t.Log("============================================================")

	// Test Text
	ctx1, cancel1 := context.WithTimeout(context.Background(), 15*time.Second)
	results1, err1 := client.Text(ctx1, query, opts)
	cancel1()
	if err1 != nil {
		t.Logf("❌ Text: %v", err1)
	} else {
		t.Logf("✅ Text: %d results", len(results1))
	}

	// Test Images
	ctx2, cancel2 := context.WithTimeout(context.Background(), 15*time.Second)
	results2, err2 := client.Images(ctx2, "butterfly", opts)
	cancel2()
	if err2 != nil {
		t.Logf("❌ Images: %v", err2)
	} else {
		t.Logf("✅ Images: %d results", len(results2))
	}

	// Test News
	ctx4, cancel4 := context.WithTimeout(context.Background(), 15*time.Second)
	results4, err4 := client.News(ctx4, "AI", opts)
	cancel4()
	if err4 != nil {
		t.Logf("❌ News: %v", err4)
	} else {
		t.Logf("✅ News: %d results", len(results4))
	}

	t.Log("============================================================")
}
