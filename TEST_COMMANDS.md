# Search Agents - Test Commands

## Test Tất Cả Categories

### 1. Text Search ✅ Works

```bash
# English
./sagent text "golang tutorial" -m 3

# Vietnamese
./sagent text "tổng bí thư việt nam" -r vn-vi -m 3

# Text + Extract
./sagent text "golang" -m 2 -e --extract-format text_markdown

# Verbose mode
./sagent text "query" -V
```

### 2. Images Search ⚠️ Need API fix

```bash
# DuckDuckGo images returns 403
./sagent images "butterfly" -m 3
./sagent images "nature landscape" -m 5 -o images.json
```

**Status**: DuckDuckGo image API (`/i.js`) bị chặn. Cần implement alternative:
- Bing Images API
- Google Images scraper
- Alternative image search providers

### 3. Videos Search ⚠️ Need API fix

```bash
# DuckDuckGo videos returns 403  
./sagent videos "programming tutorial" -m 3
./sagent videos "python basics" -m 5 -o videos.json
```

**Status**: DuckDuckGo video API (`/v.js`) bị chặn. Cần alternative.

### 4. News Search ⚠️ Need API fix

```bash
./sagent news "artificial intelligence" -m 3
./sagent news "technology" -t d -m 5 -o news.json
```

**Status**: DuckDuckGo news API (`/news.js`) bị chặn.

### 5. Books Search ⚠️ HTML structure may differ

```bash
./sagent books "golang programming" -m 3
./sagent books "python cookbook" -m 5 -o books.json
```

**Status**: Anna's Archive HTML structure cần được kiểm tra lại.

## Test API Server

```bash
# Start server
./sagent-api --port 8000

# Test text search
curl "http://localhost:8000/search/text?q=golang&max_results=3" | jq .

# Test with extract
curl "http://localhost:8000/search/text?q=golang&extract=true&max_results=2" | jq .

# Test images (will fail - API issue)
curl "http://localhost:8000/search/images?q=butterfly" | jq .

# Test health
curl "http://localhost:8000/health" | jq .
```

## Test Library (Go code)

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/tritueviet/search-agents/pkg/searchagents"
    "github.com/tritueviet/search-agents/internal/engine"
)

func main() {
    client, _ := searchagents.New(searchagents.Options{
        Timeout: 10,
        Verify:  true,
    })
    
    opts := engine.DefaultSearchOptions()
    opts.Region = "us-en"
    
    // Text search - WORKS
    results, err := client.Text(context.Background(), "golang", opts)
    fmt.Printf("Text: %d results, err=%v\n", len(results), err)
    
    // Images - WILL FAIL (403)
    results, err = client.Images(context.Background(), "butterfly", opts)
    fmt.Printf("Images: %d results, err=%v\n", len(results), err)
    
    // Videos - WILL FAIL (403)
    results, err = client.Videos(context.Background(), "tutorial", opts)
    fmt.Printf("Videos: %d results, err=%v\n", len(results), err)
    
    // News - WILL FAIL (403)
    results, err = client.News(context.Background(), "AI", opts)
    fmt.Printf("News: %d results, err=%v\n", len(results), err)
    
    // Books - MAY WORK
    results, err = client.Books(context.Background(), "golang", opts)
    fmt.Printf("Books: %d results, err=%v\n", len(results), err)
}
```

## Test MCP Server

```bash
# Start MCP server
./sagent mcp

# Test with MCP client (e.g., Cursor)
# Configuration:
{
  "mcpServers": {
    "search-agents": {
      "command": "sagent",
      "args": ["mcp"]
    }
  }
}
```

## Current Status

| Category | CLI | API | Library | Notes |
|----------|-----|-----|---------|-------|
| **Text** | ✅ | ✅ | ✅ | Works perfectly without Tor |
| **Images** | ✅* | ✅* | ✅* | **Requires Tor proxy** |
| **Videos** | ✅* | ✅* | ✅* | **Requires Tor proxy** |
| **News** | ✅* | ✅* | ✅* | **Requires Tor proxy** |
| **Books** | ✅* | ✅* | ✅* | **Requires Tor proxy** |

\* Requires Tor to be installed and running (`sudo systemctl start tor`)

## Setup Required for Full Functionality

```bash
# 1. Install Tor
sudo apt install tor
sudo systemctl start tor

# 2. Test all categories
./sagent images "butterfly" -m 3
./sagent videos "tutorial" -m 3
./sagent news "AI" -m 3
./sagent books "golang" -m 3

# See TOR_SETUP.md for detailed instructions
```

## Next Steps to Fix 403 Errors

1. **Use alternative APIs:**
   - Bing Images/Video API (requires API key)
   - Google Custom Search API (requires API key)
   - SerpAPI, Serper, etc. (third-party services)

2. **Add browser headers:**
   - Implement proper User-Agent rotation
   - Add browser fingerprinting
   - Use headless browser (Playwright, etc.)

3. **Add proxy support:**
   - Route requests through proxies
   - Rotate IPs to avoid blocking

## Working Examples (Text Search)

```bash
# All of these work:
./sagent text "golang tutorial" -m 5
./sagent text "tổng bí thư" -r vn-vi -m 3
./sagent text "python programming" -e -m 3
./sagent text "AI news" -t w -m 10 -o results.json

# API server:
curl "http://localhost:8000/search/text?q=golang&max_results=5" | jq .
curl "http://localhost:8000/search/text?q=tổng+bí+thư&region=vn-vi" | jq .
```
