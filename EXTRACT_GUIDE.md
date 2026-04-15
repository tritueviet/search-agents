# Search Agents - Extract Feature Guide

## Tổng quan

Tính năng **Extract Data** cho phép tự động trích xuất nội dung từ các URL trong kết quả search, giúp bạn:
- ✅ Đọc nội dung đầy đủ của từng kết quả
- ✅ Chuyển đổi HTML sang Markdown/Plain text
- ✅ Thu thập dữ liệu tự động cho research/aggregation
- ✅ Tiết kiệm thời gian không phải mở từng URL

## Usage

### API (Curl)

```bash
# Search + Extract tự động
curl "http://localhost:8000/search/text?q=QUERY&extract=true"

# Với định dạng Markdown
curl "http://localhost:8000/search/text?q=QUERY&extract=true&extract_format=text_markdown"

# Với Plain Text
curl "http://localhost:8000/search/text?q=QUERY&extract=true&extract_format=text_plain"

# Giới hạn + Extract
curl "http://localhost:8000/search/text?q=QUERY&extract=true&max_results=5"

# Lưu kết quả
curl "http://localhost:8000/search/text?q=QUERY&extract=true" -o results.json
```

### CLI

```bash
# Search + Extract
./sagent text "QUERY" --extract

# Short form
./sagent text "QUERY" -e

# Với định dạng
./sagent text "QUERY" --extract --extract-format text_markdown

# Giới hạn + Extract
./sagent text "QUERY" -e -m 5

# Lưu file
./sagent text "QUERY" -e -o results.json
```

## Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `extract` | bool | false | Bật/tắt extract từ result URLs |
| `extract_format` | string | text_markdown | Định dạng extract: `text_markdown`, `text_plain`, `text` |

## Response Format

### Không có Extract
```json
{
  "query": "golang",
  "results": [{"title": "...", "href": "...", "body": "..."}],
  "count": 10
}
```

### Có Extract
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
      "extracted_content": "# Full content here...",
      "extract_format": "text_markdown"
    }
  ],
  "count": 10
}
```

## Use Cases

### 1. Research tự động
```bash
# Search + đọc nội dung research papers
curl "http://localhost:8000/search/text?q=AI+research+paper&extract=true&max_results=10" \
  | jq -r '.results[].extracted_content' > research.md
```

### 2. Content Aggregation
```bash
# Thu thập tin tức từ nhiều nguồn
curl "http://localhost:8000/search/text?q=AI+news&extract=true&timelimit=w" \
  | jq '.results[] | {title, content: .extracted_content}' > news.json
```

### 3. Documentation Mining
```bash
# Lấy documentation từ sites
./sagent text "rust documentation site:docs.rs" -e -m 5 \
  --extract-format text_markdown > rust_docs.md
```

### 4. Blog Extraction
```bash
# Extract blog posts
curl "http://localhost:8000/search/text?q=go+blog+golang.org&extract=true" \
  | jq -r '.results[].extracted_content' > blog_posts.md
```

### 5. Tutorial Collection
```bash
# Sưu tập tutorials
./sagent text "python tutorial beginners" -e -m 10 \
  --extract-format text_plain > tutorials.txt
```

## Tips

1. **Timeout cao hơn khi extract:**
   ```bash
   # Tăng timeout cho extract
   ./sagent text "QUERY" -e -T 30
   ```

2. **Giới hạn số lượng để tránh overload:**
   ```bash
   curl "http://localhost:8000/search/text?q=QUERY&extract=true&max_results=3"
   ```

3. **Chỉ lấy content với jq:**
   ```bash
   curl "..." | jq -r '.results[].extracted_content'
   ```

4. **Extract song song (CLI):**
   ```bash
   ./sagent text "query1" -e -o q1.json &
   ./sagent text "query2" -e -o q2.json &
   wait
   ```

## Errors

| Error | Cause | Solution |
|-------|-------|----------|
| Failed to extract | URL không accessible | Kiểm tra URL, giảm max_results |
| Timeout | Extract lâu quá | Tăng timeout `-T 30` |
| Empty content | Page không có nội dung | Thử URL khác |

## Examples Chi Tiết

Xem thêm:
- 📖 [SKILLS.md](./SKILLS.md) - 50+ curl examples
- ⚡ [CURL_EXAMPLES.md](./CURL_EXAMPLES.md) - Quick cheatsheet
