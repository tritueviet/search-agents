# Search Agents - Curl Quick Reference

Các curl commands thông dụng nhất để sử dụng Search Agents API.

## Setup

```bash
# 1. Start API server
./sagent-api --port 8000

# 2. Kiểm tra server chạy
curl http://localhost:8000/health
```

## Text Search

```bash
# Cơ bản
curl "http://localhost:8000/search/text?q=QUERY"

# Lưu vào file
curl "http://localhost:8000/search/text?q=QUERY" -o results.json

# Giới hạn kết quả
curl "http://localhost:8000/search/text?q=QUERY&max_results=5"

# Tìm trong tuần này
curl "http://localhost:8000/search/text?q=QUERY&timelimit=w"

# Theo vùng
curl "http://localhost:8000/search/text?q=QUERY&region=uk-en"

# NEW: Search + Extract tự động
curl "http://localhost:8000/search/text?q=QUERY&extract=true" -o extracted.json

# NEW: Search + Extract với Markdown
curl "http://localhost:8000/search/text?q=QUERY&extract=true&extract_format=text_markdown"

# NEW: CLI extract
./sagent text "QUERY" --extract
./sagent text "QUERY" -e -m 5 --extract-format text_plain

# POST với JSON
curl -X POST http://localhost:8000/search/text \
  -H "Content-Type: application/json" \
  -d '{"q":"QUERY","max_results":10}'
```

## News/Videos/Images Search

```bash
# Tin tức
curl "http://localhost:8000/search/news?q=QUERY" -o news.json

# Video  
curl "http://localhost:8000/search/videos?q=QUERY" -o videos.json

# Images
curl "http://localhost:8000/search/images?q=QUERY" -o images.json
```

## Content Extraction

```bash
# Trích xuất nội dung trang web
curl "http://localhost:8000/extract?url=https://example.com"

# Chỉ lấy text
curl "http://localhost:8000/extract?url=https://example.com&format=text_plain"

# Lấy Markdown
curl "http://localhost:8000/extract?url=https://example.com&format=text_markdown" -o content.md
```

## Useful Commands

```bash
# Chỉ lấy titles
curl "http://localhost:8000/search/text?q=QUERY" | jq '.results[].title'

# Chỉ lấy URLs
curl "http://localhost:8000/search/text?q=QUERY" | jq -r '.results[].href'

# Đếm kết quả
curl "http://localhost:8000/search/text?q=QUERY" | jq '.count'

# Pretty print
curl "http://localhost:8000/search/text?q=QUERY" | jq .

# Search nhiều queries (song song)
curl "http://localhost:8000/search/text?q=golang" -o golang.json &
curl "http://localhost:8000/search/text?q=rust" -o rust.json &
curl "http://localhost:8000/search/text?q=python" -o python.json &
wait
```

## Script Mẫu

```bash
#!/bin/bash
# search.sh - Search và lưu kết quả

QUERY=$1
FILE="${QUERY// /_}.json"

curl "http://localhost:8000/search/text?q=$(echo $QUERY | sed 's/ /+/g')&max_results=10" | jq . > $FILE
echo "Saved to $FILE"
echo "Found $(jq '.count' $FILE) results"
```

```bash
#!/bin/bash
# extract.sh - Trích xuất nội dung trang web

URL=$1
curl "http://localhost:8000/extract?url=$URL" | jq -r '.content' > extracted.md
echo "Extracted to extracted.md"
```

## Response Format

### Normal Search
```json
{
  "query": "golang",
  "results": [
    {
      "title": "...",
      "href": "...",
      "body": "..."
    }
  ],
  "count": 10
}
```

### With Extract
```json
{
  "query": "golang",
  "extract": true,
  "extract_format": "text_markdown",
  "results": [
    {
      "title": "...",
      "href": "...",
      "body": "...",
      "index": 1,
      "extracted_content": "# Full markdown content here...",
      "extract_format": "text_markdown"
    }
  ],
  "count": 10
}
```

## Errors

```json
{
  "error": "query parameter 'q' is required"
}
```

---

📖 **Xem đầy đủ:** [SKILLS.md](./SKILLS.md)
