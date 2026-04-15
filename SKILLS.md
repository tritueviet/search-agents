# Search Agents - Curl Usage Guide

Hướng dẫn sử dụng `curl` và CLI để tương tác với Search Agents.

## Quick Start

### 1. Start API Server

```bash
# Build và start server
go build -o sagent-api ./cmd/server
./sagent-api --host 0.0.0.0 --port 8000

# Hoặc dùng script
./start_api.sh
```

---

## Search Engines Overview

### Available Engines (12/15 Working)

| Category | Engine | Status | Command |
|----------|--------|--------|---------|
| **Text** | DuckDuckGo, Bing, Google, Brave, Yahoo, Yandex, Wikipedia, Grokipedia, Mojeek | ✅ 9/9 | `./sagent text "query" -b ENGINE` |
| **Images** | Bing Images | ✅ 1/1 | `./sagent images "query"` |
| **Videos** | DuckDuckGo Videos | ⚠️ VQD issue | `./sagent videos "query"` |
| **News** | Bing News, DuckDuckGo News, Yahoo News | ✅ 2/3 | `./sagent news "query"` |
| **Books** | Anna's Archive | ⚠️ HTML issue | `./sagent books "query"` |

---

## 1. Text Search (9 Engines - ALL WORKING ✅)

### Basic Search (GET)

```bash
# Search đơn giản
curl "http://localhost:8000/search/text?q=golang+tutorial"

# Lưu kết quả vào file
curl "http://localhost:8000/search/text?q=golang+tutorial" -o results.json

# Pretty print JSON
curl "http://localhost:8000/search/text?q=golang+tutorial" | jq .
```

### Search với Specific Engine

```bash
# DuckDuckGo (privacy-focused)
./sagent text "golang" -b duckduckgo

# Google (best overall)
./sagent text "golang" -b google

# Bing (good for English)
./sagent text "golang" -b bing

# Brave (privacy)
./sagent text "golang" -b brave

# Yandex (best for Russian)
./sagent text "программирование" -b yandex -r ru-ru

# Wikipedia (encyclopedia only)
./sagent text "artificial intelligence" -b wikipedia
```

### Vietnamese Queries

```bash
# Best engines for Vietnamese
./sagent text "học lập trình python" -r vn-vi -b duckduckgo
./sagent text "tổng bí thư việt nam" -r vn-vi -b google
./sagent text "tin tức mới nhất" -r vn-vi
```

### Advanced Options

```bash
# With filters
./sagent text "AI programming" \
  -r us-en \
  -s moderate \
  -t w \
  -m 10 \
  -b google

# Save results
./sagent text "golang" -m 5 -o results.json

# Verbose mode
./sagent text "query" -V
```

### Search với Extract Data

```bash
# API: Search và tự động extract nội dung từ URLs
curl "http://localhost:8000/search/text?q=golang+tutorial&extract=true" | jq .

# CLI: Search, extract và lưu vào file
./sagent text "python tutorial" -e -o results.json

# Extract với định dạng Markdown
./sagent text "AI programming" --extract --extract-format text_markdown
```

---

## 2. Images Search (1 Engine - WORKING ✅)

### Bing Images Search

```bash
# CLI: Basic images search
./sagent images "butterfly" -m 5

# With filters
./sagent images "mountain landscape" \
  -m 10 \
  -o images.json

# With timelimit (recent images)
./sagent images "technology" -t w -m 5
```

### API: Images Search

```bash
# Basic search
curl "http://localhost:8000/search/images?q=butterfly&max_results=5" | jq .

# Save to file
curl "http://localhost:8000/search/images?q=nature&max_results=10" -o images.json

# With time filter
curl "http://localhost:8000/search/images?q=AI&timelimit=w&max_results=5" | jq .
```

### Response Format

```json
{
  "query": "butterfly",
  "results": [
    {
      "title": "Butterfly - Macro Photography",
      "image": "https://example.com/butterfly.jpg",
      "thumbnail": "https://ts1.mm.bing.net/th?id=...",
      "url": "https://example.com/page",
      "width": "1920",
      "height": "1080",
      "source": "example.com"
    }
  ],
  "count": 5
}
```

### Extract Image URLs

```bash
# Get only image URLs
curl "http://localhost:8000/search/images?q=butterfly" | jq -r '.results[].image'

# Get thumbnails
curl "http://localhost:8000/search/images?q=butterfly" | jq -r '.results[].thumbnail'

# Download images
curl "http://localhost:8000/search/images?q=butterfly&max_results=5" | \
  jq -r '.results[].image' | xargs -I {} wget {}
```
  -o ai_plain.txt

# Extract và chỉ lấy content
curl "http://localhost:8000/search/text?q=golang&extract=true&max_results=1" \
  | jq -r '.results[0].extracted_content'
```

### CLI Extract Data (NEW!)

```bash
# Search và extract tự động
./sagent text "golang tutorial" --extract

# Search, extract và lưu vào file
./sagent text "python tutorial" -e -o results.json

# Extract với định dạng cụ thể
./sagent text "AI programming" --extract --extract-format text_markdown

# Search, extract và giới hạn kết quả
./sagent text "Rust programming" -e -m 3 --extract-format text_plain
```

### Advanced Search Options

```bash
# Với region (us-en, uk-en, ru-ru, etc.)
curl "http://localhost:8000/search/text?q=AI+programming&region=us-en" \
  -o search_us_en.json

# Với safesearch (on, moderate, off)
curl "http://localhost:8000/search/text?q=tutorial&safesearch=off" \
  -o search_nsfw.json

# Giới hạn thời gian (d=day, w=week, m=month, y=year)
curl "http://localhost:8000/search/text?q=golang&timelimit=w" \
  -o search_this_week.json

# Số lượng kết quả tối đa
curl "http://localhost:8000/search/text?q=python&max_results=20" \
  -o search_20_results.json

# Page number (pagination)
curl "http://localhost:8000/search/text?q=rust&page=2" \
  -o search_page2.json
```

### POST Request

```bash
# POST với JSON body
curl -X POST "http://localhost:8000/search/text" \
  -H "Content-Type: application/json" \
  -d '{
    "q": "machine learning",
    "region": "us-en",
    "safesearch": "moderate",
    "timelimit": "m",
    "max_results": 15,
    "page": 1
  }' | jq . > ml_results.json
```

### Search với Filetype

```bash
# Tìm PDF files
curl "http://localhost:8000/search/text?q=golang+filetype:pdf" \
  -o golang_pdfs.json

# Tìm documentation
curl "http://localhost:8000/search/text?q=python+documentation+site:docs.python.org" \
  -o python_docs.json
```

---

## Image Search

```bash
# Search images
curl "http://localhost:8000/search/images?q=butterfly" \
  -o butterfly_images.json

# Với filters
curl "http://localhost:8000/search/images?q=nature&size=Large&color=Monochrome" \
  -o nature_bw.json
```

---

## News Search

```bash
# Latest news
curl "http://localhost:8000/search/news?q=artificial+intelligence" \
  -o ai_news.json

# News with timelimit
curl "http://localhost:8000/search/news?q=technology&timelimit=d" \
  -o today_tech_news.json

# News by region
curl "http://localhost:8000/search/news?q=politics&region=uk-en" \
  -o uk_politics.json
```

---

## Video Search (⚠️ API Blocked - Use Workaround)

### ⚠️ Problem

DuckDuckGo Videos API (`/v.js`) đã bị **chặn vĩnh viễn** và trả về kết quả rỗng.

```bash
# ❌ Cách này KHÔNG hoạt động
curl "http://localhost:8000/search/videos?q=python+tutorial"
Error: search failed: no results found
```

### ✅ Workaround: Dùng Text Search để tìm Video

Thay vì dùng video API, bạn có thể dùng text search với các operators đặc biệt.

#### Method 1: Site-Specific Search (YouTube, Vimeo, etc.)

```bash
# CLI: Tìm video trên YouTube
./sagent text "python tutorial site:youtube.com" -m 10
./sagent text "golang crash course site:youtube.com" -m 5

# CLI: Tìm trên Vimeo
./sagent text "documentary site:vimeo.com" -m 5

# CLI: Tìm trên Dailymotion
./sagent text "tutorial site:dailymotion.com" -m 5
```

#### Method 2: API Endpoints

```bash
# API: Search YouTube videos
curl "http://localhost:8000/search/text?q=python+tutorial+site:youtube.com&max_results=10" | jq .

# Extract chỉ YouTube URLs
curl "http://localhost:8000/search/text?q=golang+tutorial+site:youtube.com" | \
  jq -r '.results[].href | select(contains("youtube"))'
```

#### Method 3: Advanced Filters

```bash
# Video mới (trong tuần này)
./sagent text "AI tutorial site:youtube.com" -t w -m 10

# Video dài (>20 phút)
./sagent text "full course programming site:youtube.com" -m 5

# Video tiếng Việt
./sagent text "học lập trình python site:youtube.com" -r vn-vi -m 10

# Video có phụ đề
./sagent text "tutorial subtitles site:youtube.com" -m 5
```

#### Method 4: Save & Extract Links

```bash
# Search và lưu kết quả
./sagent text "golang tutorial site:youtube.com" -m 10 -o videos.json

# Extract chỉ YouTube URLs
jq -r '.[] | select(.href | contains("youtube")) | .href' videos.json

# Mở video đầu tiên trong trình duyệt
jq -r '.[0].href' videos.json | xargs open      # macOS
jq -r '.[0].href' videos.json | xargs xdg-open  # Linux

# Download danh sách URLs
jq -r '.[].href' videos.json | grep youtube > youtube_links.txt
```

### Example: Tìm Video Tutorials

```bash
# Python tutorials trên YouTube
./sagent text "python full course site:youtube.com" -m 5

# Output mẫu:
# 1. Python Tutorial for Beginners - freeCodeCamp
#    https://youtube.com/watch?v=rfscVS0vtbw
# 2. Learn Python in 1 Hour - Programming with Mosh
#    https://youtube.com/watch?v=kqtD5dpn9C8
# 3. ...

# Lưu và extract links
./sagent text "javascript tutorial site:youtube.com" -m 10 -o js_videos.json
jq -r '.[].href' js_videos.json | grep youtube
```

### Alternative Video APIs

Nếu cần video search API chuyên dụng:

| Service | Cost | Free Tier | Notes |
|---------|------|-----------|-------|
| **YouTube Data API** | Paid | 10,000 units/day | Tốt nhất cho YouTube |
| **SerpAPI Video** | $50/mo | 100 searches/mo | Đa nguồn |
| **Vimeo API** | Free | Limited | Chỉ Vimeo |
| **Dailymotion API** | Free | Limited | Chỉ Dailymotion |

---

## Books Search

```bash
# Search books
curl "http://localhost:8000/search/books?q=golang+programming" \
  -o golang_books.json

# Book by author
curl "http://localhost:8000/search/books?q=sea+wolf+jack+london" \
  -o jack_london.json
```

---

## Content Extraction

### Extract URL Content

```bash
# Extract as Markdown (default)
curl "http://localhost:8000/extract?url=https://example.com" \
  -o example_markdown.json

# Extract as Plain Text
curl "http://localhost:8000/extract?url=https://example.com&format=text_plain" \
  -o example_plain.txt

# Extract as Raw HTML
curl "http://localhost:8000/extract?url=https://example.com&format=text" \
  -o example_html.html

# Extract from blog
curl "https://localhost:8000/extract?url=https://blog.golang.org/defer&format=text_markdown" \
  -o blog_post.md
```

### POST Extract

```bash
curl -X POST "http://localhost:8000/extract" \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://golang.org/doc/tutorial/getting-started",
    "format": "text_markdown"
  }' | jq . > tutorial.md
```

---

## Batch Operations

### Multiple Searches (Parallel)

```bash
# Chạy nhiều searches song song
curl "http://localhost:8000/search/text?q=golang" -o golang.json &
curl "http://localhost:8000/search/text?q=rust" -o rust.json &
curl "http://localhost:8000/search/text?q=python" -o python.json &
wait

echo "All searches completed!"
```

### Search từ File

```bash
# Đọc queries từ file và search
while IFS= read -r query; do
  safe_query=$(echo "$query" | sed 's/ /+/g')
  curl "http://localhost:8000/search/text?q=$safe_query" \
    -o "result_$(echo $query | tr ' ' '_').json"
done < queries.txt
```

### Extract Multiple URLs

```bash
# Extract danh sách URLs
urls=(
  "https://golang.org"
  "https://rust-lang.org"
  "https://python.org"
)

for url in "${urls[@]}"; do
  encoded=$(echo "$url" | sed 's|/|\\/|g' | sed 's|:|%3A|g')
  curl "http://localhost:8000/extract?url=$url" \
    -o "extract_${url##*/}.json"
  sleep 1
done
```

---

## Health Check & Documentation

```bash
# Health check
curl "http://localhost:8000/health" | jq .

# API documentation
curl "http://localhost:8000/docs" | jq .
```

---

## Advanced Usage

### Với Proxy

```bash
# Search qua proxy
curl -x socks5h://127.0.0.1:9150 \
  "http://localhost:8000/search/text?q=privacy+tools" \
  -o tor_search.json
```

### Custom Headers

```bash
# Với authentication header
curl "http://localhost:8000/search/text?q=secret" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -o authenticated_search.json
```

### Timeout & Retry

```bash
# Set timeout 30s
curl --max-time 30 \
  "http://localhost:8000/search/text?q=slow+query" \
  -o result.json

# Retry 3 times
for i in {1..3}; do
  curl "http://localhost:8000/search/text?q=test" -o result.json && break
  echo "Retry $i..."
  sleep 2
done
```

### Streaming Large Results

```bash
# Stream và save
curl -N "http://localhost:8000/search/text?q=big+data&max_results=100" \
  | tee big_data_results.json
```

---

## Response Format

### Text Search Response (Without Extract)

```json
{
  "query": "golang tutorial",
  "results": [
    {
      "title": "Getting Started with Go",
      "href": "https://golang.org/doc/tutorial/getting-started",
      "body": "This tutorial is a brief introduction to writing Go programs..."
    }
  ],
  "count": 1
}
```

### Text Search Response (With Extract)

```json
{
  "query": "golang tutorial",
  "extract": true,
  "extract_format": "text_markdown",
  "results": [
    {
      "title": "Getting Started with Go",
      "href": "https://golang.org/doc/tutorial/getting-started",
      "body": "This tutorial is a brief introduction...",
      "index": 1,
      "extracted_content": "# Getting Started with Go\n\nThis tutorial is a brief introduction to writing Go programs...",
      "extract_format": "text_markdown"
    }
  ],
  "count": 1
}
```

---

## Error Handling

```bash
# Missing query parameter
curl "http://localhost:8000/search/text"
# Response: {"error": "query parameter 'q' is required"}

# Invalid URL
curl "http://localhost:8000/extract?url=invalid-url"
# Response: {"error": "failed to fetch URL: ..."}
```

---

## Scripts Examples

### search.sh - Reusable Search Script

```bash
#!/bin/bash
# Usage: ./search.sh "query" [max_results] [output_file]

QUERY=${1:?Usage: ./search.sh "query" [max_results] [output_file]}
MAX_RESULTS=${2:-10}
OUTPUT=${3:-search_results.json}

curl "http://localhost:8000/search/text?q=$(echo $QUERY | sed 's/ /+/g')&max_results=$MAX_RESULTS" \
  | jq . > "$OUTPUT"

echo "Results saved to $OUTPUT"
cat "$OUTPUT" | jq '.count'
```

### extract.sh - URL Extraction Script

```bash
#!/bin/bash
# Usage: ./extract.sh <url> [format]

URL=${1:?Usage: ./extract.sh <url> [format]}
FORMAT=${2:-text_markdown}
OUTPUT="extract_$(basename $URL).json"

curl "http://localhost:8000/extract?url=$URL&format=$FORMAT" \
  | jq -r '.content' > "$OUTPUT"

echo "Content extracted to $OUTPUT"
```

---

## Tips & Tricks

### Extract Data Use Cases

```bash
# 1. Research: Search + đọc nội dung tự động
curl "http://localhost:8000/search/text?q=golang+best+practices&extract=true&max_results=5" \
  | jq -r '.results[].extracted_content' > research.md

# 2. Content aggregation: Thu thập nội dung từ nhiều nguồn
curl "http://localhost:8000/search/text?q=AI+news+2024&extract=true&timelimit=w" \
  | jq '.results[] | {title, content: .extracted_content}' > news_aggregated.json

# 3. Documentation mining
curl "http://localhost:8000/search/text?q=rust+documentation+site:docs.rs&extract=true" \
  | jq -r '.results[0].extracted_content' > rust_docs.md

# 4. Blog post extraction
curl "http://localhost:8000/search/text?q=go+blog+golang.org&extract=true" \
  | jq -r '.results[].extracted_content' > go_blog_posts.md

# 5. Tutorial extraction với CLI
./sagent text "python tutorial" -e -m 3 --extract-format text_markdown | \
  grep -A 1000 "Result #1" > tutorial.md
```

1. **Sử dụng `jq` để format và filter JSON:**
   ```bash
   curl "http://localhost:8000/search/text?q=golang" | jq '.results[].title'
   ```

2. **Chỉ lấy URLs:**
   ```bash
   curl "http://localhost:8000/search/text?q=golang" | jq -r '.results[].href'
   ```

3. **Đếm số kết quả:**
   ```bash
   curl "http://localhost:8000/search/text?q=golang" | jq '.count'
   ```

4. **Search và mở trình duyệt:**
   ```bash
   curl "http://localhost:8000/search/text?q=golang&max_results=1" | jq -r '.results[0].href' | xargs open
   ```

5. **Download tất cả URLs:**
   ```bash
   curl "http://localhost:8000/search/text?q=golang" | jq -r '.results[].href' | xargs -I {} wget {}
   ```
