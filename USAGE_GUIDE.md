# Search Agents - CLI Usage Guide

## Quick Start

```bash
# Build
make build-cli

# Basic search
./sagent text "golang tutorial"

# Search with specific engine
./sagent text "golang" -b duckduckgo
./sagent text "golang" -b google
./sagent text "golang" -b bing
```

## Available Search Engines

### Text Search Engines (9 engines)

| Engine | Command | Region | Priority | Notes |
|--------|---------|--------|----------|-------|
| **DuckDuckGo** | `-b duckduckgo` | us-en, vn-vi, etc. | 1.0 | ✅ Best privacy, no tracking |
| **Bing** | `-b bing` | us-en, uk-en, etc. | 0.9 | ✅ Good for English queries |
| **Google** | `-b google` | us-en, etc. | 0.85 | ✅ Best overall results |
| **Brave** | `-b brave` | us-en, etc. | 0.8 | ✅ Privacy-focused |
| **Yahoo** | `-b yahoo` | us-en, etc. | 0.7 | ⚠️ Uses Bing results |
| **Yandex** | `-b yandex` | ru-ru, etc. | 0.75 | ✅ Best for Russian queries |
| **Wikipedia** | `-b wikipedia` | en, vi, etc. | 0.85 | ✅ Encyclopedia only |
| **Grokipedia** | `-b grokipedia` | us-en | 0.7 | ⚠️ May redirect |
| **Mojeek** | `-b mojeek` | us-en, uk-en | 0.65 | ✅ Independent crawler |

## CLI Commands

### 1. Text Search

```bash
# Basic search (auto-selects engines)
./sagent text "golang tutorial"

# Use specific engine
./sagent text "golang" -b duckduckgo
./sagent text "python" -b google
./sagent text "programming" -b brave

# Vietnamese queries
./sagent text "học lập trình" -r vn-vi
./sagent text "tổng bí thư" -r vn-vi -b duckduckgo

# Multiple options
./sagent text "AI programming" \
  -r us-en \
  -s moderate \
  -t w \
  -m 10 \
  -b google

# Save results
./sagent text "golang" -m 5 -o results.json

# Verbose mode (show errors)
./sagent text "query" -V
```

**Flags:**
- `-r, --region` - Region code (us-en, vn-vi, ru-ru, uk-en)
- `-s, --safesearch` - SafeSearch level (on, moderate, off)
- `-t, --timelimit` - Time filter (d=day, w=week, m=month, y=year)
- `-m, --max-results` - Maximum results (default: 10)
- `-b, --backend` - Search engine (duckduckgo, google, bing, brave, etc.)
- `-o, --output` - Save to file (JSON format)
- `-V, --verbose` - Show detailed errors

### 2. Image Search (⚠️ API deprecated)

```bash
# Will show error with alternatives
./sagent images "butterfly"
```

**Status:** DuckDuckGo Images API has been blocked. 
**Alternatives:**
- Use third-party APIs (SerpAPI, Google Custom Search)
- See `ENGINE_STATUS.md` for details

### 3. Video Search (⚠️ API deprecated)

```bash
# Will show error with alternatives
./sagent videos "tutorial"
```

**Status:** DuckDuckGo Videos API has been blocked.

### 4. News Search (⚠️ API deprecated)

```bash
# Will show error with alternatives
./sagent news "AI"
```

**Status:** DuckDuckGo News API has been blocked.

### 5. Book Search

```bash
# Search books
./sagent books "golang programming"
```

**Status:** Anna's Archive - needs HTML structure update.

## Engine Selection Guide

### For English Queries
```bash
# Best overall
./sagent text "query" -b google

# Privacy-focused
./sagent text "query" -b duckduckgo
./sagent text "query" -b brave

# Fast results
./sagent text "query" -b bing
```

### For Vietnamese Queries
```bash
# Best for Vietnamese
./sagent text "query" -r vn-vi -b duckduckgo
./sagent text "query" -r vn-vi -b google

# Example
./sagent text "học python cơ bản" -r vn-vi -b duckduckgo
```

### For Russian Queries
```bash
# Best for Russian
./sagent text "запрос" -r ru-ru -b yandex

# Example
./sagent text "программирование" -r ru-ru -b yandex
```

### For Encyclopedia Only
```bash
# Wikipedia only
./sagent text "artificial intelligence" -b wikipedia

# Grokipedia
./sagent text "quantum computing" -b grokipedia
```

## Advanced Usage

### Compare Multiple Engines

```bash
# Test all engines
for engine in duckduckgo google bing brave yahoo yandex wikipedia mojeek; do
  echo "=== $engine ==="
  ./sagent text "golang" -b $engine -m 2 2>&1 | head -10
  echo ""
done
```

### Batch Search

```bash
#!/bin/bash
# search.sh - Search with multiple engines

QUERY=$1
ENGINES="duckduckgo google bing brave"

for engine in $ENGINES; do
  echo "🔍 Searching with $engine..."
  ./sagent text "$QUERY" -b $engine -m 3 -o "results_${engine}.json"
done

echo "✅ Search complete!"
ls -lh results_*.json
```

### Extract + Search

```bash
# Search and extract top result
./sagent text "golang tutorial" -m 1 -o top.json
URL=$(jq -r '.[0].href' top.json)
./sagent extract "$URL" -f text_markdown
```

## Error Handling

### Engine Not Available
```bash
$ ./sagent text "query" -b unknown
Error: search failed: no engines found for category: text
```

### Engine Blocked (403)
```bash
$ ./sagent images "query"
Error: DuckDuckGo Images API has been deprecated (403 Forbidden).
Alternatives:
  1. Use Bing Images API (requires API key)
  2. Use Google Custom Search API (requires API key)
  3. Use SerpAPI or similar third-party services
```

### No Results
```bash
$ ./sagent text "xyzabc123"
Error: no results found
```

## Makefile Commands

```bash
# Test all engines
make test-text

# Test specific engine
./sagent text "query" -b duckduckgo -V

# Build
make build-cli

# Run example
make run
```

## Performance Comparison

| Engine | Avg Time | Result Quality | Privacy | Best For |
|--------|----------|----------------|---------|----------|
| DuckDuckGo | ~2s | Good | ✅✅✅ | Privacy |
| Google | ~1.5s | ✅✅✅ Best | ⚠️ | Overall |
| Bing | ~2s | Good | ⚠️ | English |
| Brave | ~2s | Good | ✅✅✅ | Privacy |
| Yandex | ~1.5s | Good | ⚠️ | Russian |
| Wikipedia | ~0.5s | ✅✅✅ | ✅✅✅ | Encyclopedia |
| Yahoo | ~1s | Fair | ⚠️ | Alternative |
| Mojeek | ~1.5s | Fair | ✅✅✅ | Independent |
| Grokipedia | ~2s | Fair | ⚠️ | Alternative |

## Tips & Tricks

### 1. Use region for better results
```bash
# Vietnamese
./sagent text "query" -r vn-vi

# Russian  
./sagent text "запрос" -r ru-ru

# UK English
./sagent text "query" -r uk-en
```

### 2. Filter by time
```bash
# Today only
./sagent text "AI news" -t d

# This week
./sagent text "golang" -t w

# This month
./sagent text "python" -t m

# This year
./sagent text "programming" -t y
```

### 3. Control SafeSearch
```bash
# Strict
./sagent text "query" -s on

# Moderate (default)
./sagent text "query" -s moderate

# Off
./sagent text "query" -s off
```

### 4. Save and process results
```bash
# Save to JSON
./sagent text "query" -m 20 -o results.json

# Extract titles
jq '.[].title' results.json

# Extract URLs
jq -r '.[].href' results.json

# Count results
jq 'length' results.json
```

## Supported Regions

| Code | Country/Language |
|------|------------------|
| us-en | United States (English) |
| uk-en | United Kingdom (English) |
| vn-vi | Vietnam (Vietnamese) |
| ru-ru | Russia (Russian) |
| de-de | Germany (German) |
| fr-fr | France (French) |
| es-es | Spain (Spanish) |
| ja-jp | Japan (Japanese) |
| zh-cn | China (Chinese) |

## Next Steps

- [ ] Add Images/Videos/News with third-party APIs
- [ ] Fix Books scraping
- [ ] Add more engines (Baidu, etc.)
- [ ] Implement result ranking algorithms

See `TEST_RESULTS.md` for detailed test reports.
