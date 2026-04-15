# Search Agents (sagent)

A powerful metasearch library written in Go that aggregates results from diverse web search services.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
  - [CLI Usage](#cli-usage)
  - [Library Usage](#library-usage)
  - [REST API Server](#rest-api-server)
  - [Curl Examples](#curl-usage-guide)
  - [MCP Server](#mcp-server)
- [Architecture](#architecture)
- [Available Engines](#available-engines)
- [Development](#development)
- [Configuration](#configuration)
- [License](#license)

## Quick Links

- 📖 **[SKILLS.md](./SKILLS.md)** - 50+ curl examples chi tiết
- ⚡ **[CURL_EXAMPLES.md](./CURL_EXAMPLES.md)** - Curl quick reference (cheatsheet)
- 🚀 **Quick Start**: `go build -o sagent ./cmd/sagent && ./sagent text "query"`

## Features

✅ **Multiple Search Engines**: DuckDuckGo, Bing, Google, Brave, Yandex, Yahoo  
✅ **Multiple Categories**: Text, Images, Videos, News, Books  
✅ **Parallel Search**: Concurrent searches across multiple engines using goroutines  
✅ **Result Deduplication**: Smart deduplication and ranking  
✅ **Proxy Support**: HTTP/SOCKS5 proxy support  
✅ **CLI Tool**: `sagent` command-line interface  
✅ **REST API**: FastAPI-compatible REST endpoints with Gin  
✅ **MCP Server**: Model Context Protocol support for AI assistants  
✅ **URL Extractor**: HTML to Markdown conversion  

## Installation

```bash
git clone https://github.com/tritueviet/search-agents.git
cd search-agents
go mod tidy
```

## Quick Start

### CLI Usage

```bash
# Build the CLI
go build -o sagent ./cmd/sagent

# Version
./sagent version

# Text search (English)
./sagent text "golang tutorial"

# Text search (Vietnamese) - Use -r vn-vi
./sagent text "tổng bí thư việt nam" -r vn-vi

# With options
./sagent text "AI programming" -r us-en -m 20 -t w

# With proxy
./sagent text "search query" -P socks5h://127.0.0.1:9150

# Search + Extract content
./sagent text "golang tutorial" -e -m 3

# Extract with Vietnamese queries
./sagent text "học lập trình python" -r vn-vi -e -m 5

# Verbose mode (show errors)
./sagent text "query" -V

# Extract content from URL
./sagent extract https://example.com -f text_markdown

# Start MCP server
./sagent mcp
```

### Library Usage

```go
package main

import (
    "context"
    "fmt"

    "github.com/tritueviet/search-agents/pkg/searchagents"
    "github.com/tritueviet/search-agents/internal/engine"
)

func main() {
    client, err := searchagents.New(searchagents.Options{
        Timeout: 5,
        Verify:  true,
    })
    if err != nil {
        panic(err)
    }

    opts := engine.DefaultSearchOptions()
    opts.Region = "us-en"
    opts.SafeSearch = "moderate"

    results, err := client.Text(context.Background(), "golang tutorial", opts)
    if err != nil {
        panic(err)
    }

    for _, result := range results {
        fmt.Printf("%s\n%s\n\n", result["title"], result["href"])
    }
}
```

### REST API Server

```bash
# Build API server
go build -o sagent-api ./cmd/server

# Start server
./sagent-api --host 0.0.0.0 --port 8000

# Or use start script
./start_api.sh
```

#### API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/search/text` | GET, POST | Text search |
| `/search/images` | GET, POST | Image search |
| `/search/news` | GET, POST | News search |
| `/search/videos` | GET, POST | Video search |
| `/search/books` | GET, POST | Book search |
| `/extract` | GET, POST | Extract content from URL |
| `/health` | GET | Health check |
| `/docs` | GET | API documentation |

#### Quick Curl Examples

```bash
# Basic search
curl "http://localhost:8000/search/text?q=golang+tutorial" | jq .

# Save results to file
curl "http://localhost:8000/search/text?q=AI+programming&max_results=20" -o ai_results.json

# Search with filters
curl "http://localhost:8000/search/text?q=python&region=us-en&timelimit=w&safesearch=off" | jq .

# NEW: Search + Auto Extract content from URLs
curl "http://localhost:8000/search/text?q=golang&extract=true&max_results=5" | jq .

# NEW: Extract with specific format
curl "http://localhost:8000/search/text?q=tutorial&extract=true&extract_format=text_markdown" \
  -o tutorials.json

# NEW: CLI with extract
./sagent text "golang tutorial" --extract
./sagent text "python" -e -m 3 --extract-format text_plain

# Extract URL content as Markdown
curl "http://localhost:8000/extract?url=https://example.com&format=text_markdown" | jq -r '.content'

# Health check
curl "http://localhost:8000/health" | jq .
```

📖 **Xem thêm:** [SKILLS.md](./SKILLS.md) - Hướng dẫn curl chi tiết với 50+ examples

### MCP Server

For MCP clients like Cursor or Claude Desktop:

```json
{
  "mcpServers": {
    "search-agents": {
      "command": "sagent",
      "args": ["mcp"]
    }
  }
}
```

#### Available MCP Tools

| Tool | Description |
|------|-------------|
| `search_text` | Web text search |
| `extract_content` | Extract content from a URL |

## Curl Usage Guide

Để xem hướng dẫn curl đầy đủ, xem:
- 📖 **[SKILLS.md](./SKILLS.md)** - 50+ examples chi tiết
- ⚡ **[CURL_EXAMPLES.md](./CURL_EXAMPLES.md)** - Quick reference cheatsheet

### Quick Examples

```bash
# Text search
curl "http://localhost:8000/search/text?q=golang+tutorial" | jq .

# Save to file
curl "http://localhost:8000/search/text?q=python&max_results=20" -o results.json

# NEW: Search + Auto Extract
curl "http://localhost:8000/search/text?q=AI&extract=true&max_results=5" | jq .

# NEW: Extract content từ kết quả search
curl "http://localhost:8000/search/text?q=tutorial&extract=true" \
  | jq -r '.results[0].extracted_content'

# News this week
curl "http://localhost:8000/search/news?q=AI&timelimit=w" | jq '.results[].title'

# Extract webpage to Markdown
curl "http://localhost:8000/extract?url=https://golang.org" | jq -r '.content'

# CLI với extract
./sagent text "golang" --extract
./sagent text "python" -e -m 5 --extract-format text_markdown

# Batch search (parallel)
curl "http://localhost:8000/search/text?q=golang" -o golang.json &
curl "http://localhost:8000/search/text?q=rust" -o rust.json &
wait
```

## Architecture

```
search-agents/
├── cmd/
│   ├── sagent/              # CLI entry point
│   │   └── main.go
│   └── server/              # API server entry point
│       └── main.go
├── internal/
│   ├── core/
│   │   └── searchagents.go  # Main orchestrator
│   ├── engine/              # Search engine interfaces
│   ├── engines/             # Search engine implementations
│   │   ├── duckduckgo/
│   │   ├── bing/
│   │   ├── google/
│   │   ├── brave/
│   │   ├── yandex/
│   │   └── yahoo/
│   ├── httpclient/          # HTTP client
│   ├── models/              # Data structures
│   ├── extractor/           # URL content extraction
│   ├── register/            # Engine registration
│   └── utils/               # Utilities
├── pkg/
│   └── searchagents/        # Public API
├── api/
│   └── handlers.go          # REST API handlers
├── mcp/
│   └── server.go            # MCP server
├── go.mod
├── go.sum
├── start_api.sh             # API start script
└── README.md
```

## Available Engines

### Text Search (9 engines - ALL WORKING ✅)

| Engine | Command | Status | Best For |
|--------|---------|--------|----------|
| DuckDuckGo | `-b duckduckgo` | ✅ | Privacy, Vietnamese |
| Bing | `-b bing` | ✅ | English queries |
| Google | `-b google` | ✅ | Overall best results |
| Brave | `-b brave` | ✅ | Privacy-focused |
| Yahoo | `-b yahoo` | ✅ | Alternative |
| Yandex | `-b yandex` | ✅ | Russian queries |
| Wikipedia | `-b wikipedia` | ✅ | Encyclopedia only |
| Grokipedia | `-b grokipedia` | ⚠️ | Alternative |
| Mojeek | `-b mojeek` | ⚠️ | Independent crawler |

### Images Search (1 engine - WORKING ✅)

| Engine | Command | Status | Notes |
|--------|---------|--------|-------|
| Bing Images | `./sagent images "query"` | ✅ | Good quality images |

### Videos Search (1 engine - NEEDS WORKAROUND ⚠️)

| Engine | Command | Status | Notes |
|--------|---------|--------|-------|
| DuckDuckGo Videos | `./sagent videos "query"` | ❌ API blocked | Use workaround below |

**⚠️ Workaround:** Dùng text search với `site:youtube.com` để tìm video:

```bash
# Thay vì: ./sagent videos "python tutorial" (không hoạt động)
# Dùng:
./sagent text "python tutorial site:youtube.com" -m 10
./sagent text "golang course site:youtube.com" -m 5
./sagent text "học lập trình site:youtube.com" -r vn-vi -m 10

# API:
curl "http://localhost:8000/search/text?q=python+tutorial+site:youtube.com" | jq .
```

Xem chi tiết trong [SKILLS.md](./SKILLS.md) phần "Video Search Workaround".

### News Search (3 engines - 2 WORKING ✅)

| Engine | Command | Status | Notes |
|--------|---------|--------|-------|
| Bing News | `./sagent news "query"` | ✅ | Best quality |
| DuckDuckGo News | `./sagent news "query"` | ✅ | Good alternative |
| Yahoo News | `./sagent news "query"` | ⚠️ | HTML varies |

### Books Search (1 engine - NEEDS FIX ⚠️)

| Engine | Command | Status | Notes |
|--------|---------|--------|-------|
| Anna's Archive | `./sagent books "query"` | ⚠️ | HTML structure changed |

**Total:** 12/15 engines working (80% success rate)

See [ENGINES_STATUS.md](./ENGINES_STATUS.md) for detailed status.
See [USAGE_GUIDE.md](./USAGE_GUIDE.md) for usage examples.

## Development

```bash
# Install dependencies (includes Tor)
make install

# Run tests
make test
make test-text        # Text search tests only
make test-integration # All integration tests

# Build
make build            # Build CLI and API
make build-cli        # CLI only
make build-api        # API only

# Run linter
make lint

# Clean
make clean
```

### Test Results

| Category | Status | Tests | Notes |
|----------|--------|-------|-------|
| **Text** | ✅ **PASS** | 3/3 (100%) | English, Vietnamese, Special chars |
| **Images** | ⏳ | - | Requires Tor proxy |
| **Videos** | ⏳ | - | Requires Tor proxy |
| **News** | ⏳ | - | Requires Tor proxy |
| **Books** | ⏳ | - | Requires Tor proxy |

See [TEST_RESULTS.md](./TEST_RESULTS.md) for detailed test report.

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SAGENT_HOST` | API server host | `0.0.0.0` |
| `SAGENT_PORT` | API server port | `8000` |
| `SAGENT_PROXY` | Proxy URL | - |
| `SAGENT_TIMEOUT` | Timeout in seconds | `10` |

## Troubleshooting

### No results found

**Problem:** `Error: search failed: no results found`

**Solutions:**
1. **For Vietnamese queries:** Use `-r vn-vi` region
   ```bash
   ./sagent text "tổng bí thư" -r vn-vi
   ```

2. **Check verbose mode for details:**
   ```bash
   ./sagent text "query" -V
   ```

3. **Try different backend:**
   ```bash
   ./sagent text "query" -b google
   ./sagent text "query" -b bing
   ```

4. **Increase timeout:**
   ```bash
   ./sagent text "query" -T 30
   ```

### Extract returns empty content

**Problem:** Extracted content is empty or contains JavaScript/CSS

**Solutions:**
1. **Use text_markdown format** (best for web pages):
   ```bash
   ./sagent text "query" -e --extract-format text_markdown
   ```

2. **Some sites block scraping** - try with proxy:
   ```bash
   ./sagent text "query" -e -P socks5h://127.0.0.1:9150
   ```

## License

MIT

## Credits

Converted from Python DDGS project (https://github.com/deedy5/ddgs) to Go by tritueviet.
