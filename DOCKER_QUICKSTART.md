# Search Agents - Docker Quick Start

## 1. Build & Run

```bash
# Build image
docker-compose build

# Start server
docker-compose up -d

# Check status
docker-compose ps
```

## 2. Test API

```bash
# Health check
curl http://localhost:8000/health

# Text search
curl "http://localhost:8000/search/text?q=golang&max_results=3" | jq .

# Images
curl "http://localhost:8000/search/images?q=butterfly&max_results=3" | jq .

# News
curl "http://localhost:8000/search/news?q=AI&max_results=3" | jq .
```

## 3. Common Commands

```bash
# View logs
docker-compose logs -f

# Stop
docker-compose down

# Rebuild after code changes
docker-compose up -d --build

# Restart
docker-compose restart
```

## Configuration

Edit `.env` file:
```bash
SAGENT_PORT=8000
SAGENT_TIMEOUT=10
SAGENT_PROXY=  # Optional: socks5h://127.0.0.1:9050
```

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Health check |
| `/search/text` | GET/POST | Text search (9 engines) |
| `/search/images` | GET/POST | Images (Bing) |
| `/search/news` | GET/POST | News (Bing, DDG, Yahoo) |
| `/search/videos` | GET/POST | Videos (needs fix) |
| `/search/books` | GET/POST | Books (needs fix) |
| `/docs` | GET | Documentation |

## Examples

### Text Search
```bash
curl "http://localhost:8000/search/text?q=python+tutorial&region=us-en&max_results=5"
```

### Images Search
```bash
curl "http://localhost:8000/search/images?q=nature&max_results=5"
```

### Vietnamese Search
```bash
curl "http://localhost:8000/search/text?q=học+lập+trình&region=vn-vi&max_results=5"
```

### Extract Content
```bash
curl "http://localhost:8000/extract?url=https://golang.org&format=text_markdown"
```

## Troubleshooting

### Port in use
```bash
# Change port in .env
echo "SAGENT_PORT=9000" > .env
docker-compose up -d
```

### Rebuild after changes
```bash
docker-compose up -d --build
```

### View logs
```bash
docker-compose logs -f api
```

See [DOCKER_GUIDE.md](./DOCKER_GUIDE.md) for full documentation.
