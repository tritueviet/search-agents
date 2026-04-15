# Search Agents v1.0.1 🚀

**Dux Distributed Global Search - Go Implementation**

A powerful metasearch library written in Go that aggregates results from diverse web search services.

## 📊 Release Stats

| Category | Engines | Working | Success Rate |
|----------|---------|---------|--------------|
| **Text** | 9 | 9 | 100% ✅ |
| **Images** | 1 | 1 | 100% ✅ |
| **News** | 3 | 2 | 67% ✅ |
| **Books** | 2 | 1 | 50% ⚠️ |
| **Videos** | 1 | 0 | 0% ❌ (workaround available) |
| **Total** | **16** | **13** | **81%** |

## ✨ Features

### 🔍 16 Search Engines

#### Text Search (9 engines - ALL WORKING)
- **DuckDuckGo** - Privacy-focused, best for Vietnamese
- **Bing** - Good English results
- **Google** - Best overall quality
- **Brave** - Privacy-focused
- **Yahoo** - Alternative
- **Yandex** - Best for Russian queries
- **Wikipedia** - Official API, very fast
- **Grokipedia** - Alternative
- **Mojeek** - Independent crawler

#### Images Search (1 engine - WORKING)
- **Bing Images** - Good quality images with metadata

#### News Search (3 engines - 2 WORKING)
- **Bing News** - Best quality ✅
- **DuckDuckGo News** - Good alternative ✅
- **Yahoo News** - HTML structure varies ⚠️

#### Books Search (2 engines - 1 WORKING)
- **Open Library** - Free API, working perfectly ✅ (NEW in v1.0.1)
- **Anna's Archive** - Timeout issues ⚠️

#### Videos Search (1 engine - BLOCKED)
- **DuckDuckGo Videos** - API deprecated ❌
- **Workaround**: Use `./sagent text "query site:youtube.com"` ✅

### 🛠️ CLI Tool (`sagent`)

```bash
# Text search
./sagent text "golang tutorial" -m 10
./sagent text "học lập trình" -r vn-vi

# Images search
./sagent images "butterfly" -m 5

# News search
./sagent news "AI technology" -m 5

# Books search
./sagent books "golang programming" -m 3

# Find videos (workaround)
./sagent text "python tutorial site:youtube.com" -m 10

# Extract content
./sagent extract https://example.com -f text_markdown
```

### 🌐 REST API Server

```bash
# Start server
./sagent-api --port 8000

# Endpoints
GET /health              # Health check
GET /search/text         # Text search
GET /search/images       # Images search
GET /search/news         # News search
GET /extract             # Content extraction
```

### 🐳 Docker Support

```bash
# Build and run
docker-compose build
docker-compose up -d

# Test
curl http://localhost:8000/health
```

### 🔌 MCP Server

Model Context Protocol support for AI assistants (Cursor, Claude Desktop).

## 🚀 Quick Start

### CLI

```bash
# Build
go build -o sagent ./cmd/sagent

# Search
./sagent text "golang tutorial" -m 5
```

### API Server

```bash
# Build and run
go build -o sagent-api ./cmd/server
./sagent-api --port 8000

# Test
curl "http://localhost:8000/search/text?q=golang&max_results=5" | jq .
```

### Docker

```bash
docker-compose build
docker-compose up -d
```

## 📦 Installation

```bash
git clone https://github.com/tritueviet/search-agents.git
cd search-agents
go mod tidy
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SAGENT_HOST` | API server host | `0.0.0.0` |
| `SAGENT_PORT` | API server port | `8000` |
| `SAGENT_TIMEOUT` | Timeout in seconds | `10` |
| `SAGENT_PROXY` | Proxy URL (optional) | - |

## 📚 Documentation

- [README.md](./README.md) - Main documentation
- [SKILLS.md](./SKILLS.md) - Comprehensive usage guide with 50+ examples
- [USAGE_GUIDE.md](./USAGE_GUIDE.md) - Engine selection guide
- [DOCKER_GUIDE.md](./DOCKER_GUIDE.md) - Docker setup
- [TEST_RESULTS.md](./TEST_RESULTS.md) - Test reports
- [ENGINE_FIXES.md](./ENGINE_FIXES.md) - Engine status

## 🧪 Testing

```bash
# Run all tests
go test -v ./tests/...

# Test text engines
go test -v -run TestAllTextEngines ./tests/...

# Test specific categories
go test -v -run TestImagesSearch ./tests/...
go test -v -run TestNewsSearch ./tests/...
```

## ⚠️ Known Issues

### Videos Search
- **Issue**: DuckDuckGo Videos API deprecated
- **Workaround**: Use `./sagent text "query site:youtube.com"`
- **Future**: Consider YouTube Data API or SerpAPI

### Books Search
- **Working**: Open Library engine (free API)
- **Issue**: Anna's Archive timeout
- **Status**: 50% working

### News Search
- **Working**: Bing News, DuckDuckGo News
- **Issue**: Yahoo News HTML structure varies
- **Status**: 67% working

## 📋 What's New in v1.0.1

### Added ✨
- **9 Text Search Engines**: DuckDuckGo, Bing, Google, Brave, Yahoo, Yandex, Wikipedia, Grokipedia, Mojeek
- **Bing Images Engine**: Image search with metadata
- **Bing News Engine**: News search with timestamps
- **DuckDuckGo News Engine**: Alternative news source
- **Open Library Engine**: Free books search API (NEW!)
- **Comprehensive CLI**: Full-featured command-line tool
- **REST API Server**: Gin-based HTTP server
- **MCP Server**: Model Context Protocol support
- **Content Extractor**: HTML to Markdown conversion
- **Docker Support**: Multi-stage build, docker-compose
- **Tests**: Comprehensive test suite
- **Documentation**: 6 documentation files

### Fixed 🔧
- **Books Search**: Added Open Library engine (working)
- **Video Search Workaround**: Text search with `site:youtube.com`
- **Vietnamese Support**: Full support with `-r vn-vi`
- **Error Messages**: Clear, actionable error messages
- **Tor Proxy**: Auto-detection and fallback

### Changed 🔄
- Renamed module to `github.com/tritueviet/search-agents`
- Renamed CLI from `ddgs` to `sagent`
- Improved search result deduplication
- Enhanced error handling

## 🤝 Contributing

Contributions welcome! Please read the documentation and submit issues/PRs.

## 📄 License

MIT License - see LICENSE file for details.

## 🔗 Links

- **Source**: https://github.com/tritueviet/search-agents
- **Original Python Project**: https://github.com/deedy5/ddgs
- **Issues**: https://github.com/tritueviet/search-agents/issues

---

**Built with Go 1.23+** 🐹
