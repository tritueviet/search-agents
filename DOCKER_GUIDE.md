# Search Agents Docker Setup

## Quick Start

### 1. Build and Run

```bash
# Build Docker image
docker-compose build

# Start API server
docker-compose up -d

# Check logs
docker-compose logs -f

# Stop server
docker-compose down
```

### 2. Test API

```bash
# Health check
curl http://localhost:8000/health

# Text search
curl "http://localhost:8000/search/text?q=golang&max_results=5" | jq .

# Images search
curl "http://localhost:8000/search/images?q=butterfly&max_results=5" | jq .

# News search
curl "http://localhost:8000/search/news?q=AI&max_results=5" | jq .
```

## Configuration

### Environment Variables

Create `.env` file:

```bash
# API Server
SAGENT_PORT=8000
SAGENT_TIMEOUT=10

# Proxy (optional)
SAGENT_PROXY=socks5h://127.0.0.1:9050
```

### docker-compose.yml

Default configuration:
```yaml
ports:
  - "8000:8000"  # Change to your preferred port

environment:
  - SAGENT_PORT=8000
  - SAGENT_TIMEOUT=10
```

## Commands

```bash
# Build
docker-compose build

# Start
docker-compose up -d

# Stop
docker-compose down

# View logs
docker-compose logs -f

# Restart
docker-compose restart

# Rebuild and restart
docker-compose up -d --build
```

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Health check |
| `/search/text` | GET, POST | Text search |
| `/search/images` | GET, POST | Images search |
| `/search/news` | GET, POST | News search |
| `/search/videos` | GET, POST | Videos search |
| `/search/books` | GET, POST | Books search |
| `/docs` | GET | API documentation |

## Examples

### Text Search

```bash
# Basic search
curl "http://localhost:8000/search/text?q=golang+tutorial" | jq .

# With options
curl "http://localhost:8000/search/text?q=python&region=us-en&max_results=10" | jq .

# Vietnamese
curl "http://localhost:8000/search/text?q=học+lập+trình&region=vn-vi" | jq .

# With extract
curl "http://localhost:8000/search/text?q=golang&extract=true&max_results=3" | jq .
```

### Images Search

```bash
curl "http://localhost:8000/search/images?q=butterfly&max_results=5" | jq .
```

### News Search

```bash
curl "http://localhost:8000/search/news?q=AI+technology&max_results=5" | jq .
```

## Troubleshooting

### Port Already in Use

```bash
# Change port in .env
echo "SAGENT_PORT=9000" > .env
docker-compose up -d
```

### Health Check Failing

```bash
# Check container status
docker-compose ps

# View logs
docker-compose logs api

# Restart
docker-compose restart
```

### Rebuild After Code Changes

```bash
docker-compose up -d --build
```

## Production Deployment

For production, consider:

1. Use reverse proxy (nginx/traefik)
2. Add SSL/TLS termination
3. Set up monitoring
4. Configure resource limits:

```yaml
deploy:
  resources:
    limits:
      cpus: '2'
      memory: 1G
    reservations:
      cpus: '0.5'
      memory: 256M
```
