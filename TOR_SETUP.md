# Search Agents - Setup Tor Proxy

Để sử dụng các tính năng Images, Videos, News, Books search mà không bị chặn (403), bạn cần cài Tor proxy.

## Quick Setup

### 1. Install Tor

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install tor
sudo systemctl start tor
sudo systemctl enable tor
```

**Arch Linux:**
```bash
sudo pacman -S tor
sudo systemctl start tor
sudo systemctl enable tor
```

**macOS (Homebrew):**
```bash
brew install tor
brew services start tor
```

**Windows:**
Download từ https://www.torproject.org/download/

### 2. Verify Tor is Running

```bash
# Check if Tor is listening on port 9150
ss -tlnp | grep 9150
# or
netstat -tlnp | grep 9150
```

Expected output:
```
LISTEN  0  4096  127.0.0.1:9150  0.0.0.0:*
```

### 3. Test with Search Agents

```bash
# Images search (will use Tor automatically)
./sagent images "butterfly" -m 3

# Videos search
./sagent videos "programming tutorial" -m 3

# News search
./sagent news "artificial intelligence" -m 5

# Books search
./sagent books "golang programming" -m 3
```

## How It Works

Search Agents sẽ **tự động sử dụng Tor proxy** cho các engines:
- ✅ `duckduckgo_images` - Always uses Tor
- ✅ `duckduckgo_videos` - Always uses Tor  
- ✅ `duckduckgo_news` - Always uses Tor
- ✅ `annasarchive` - Always uses Tor

Nếu Tor không khả dụng, sẽ fallback về normal client (có thể bị 403).

## Custom Tor Proxy

Nếu bạn dùng Tor proxy khác port:

```bash
# Set environment variable
export TOR_PROXY_URL="socks5h://127.0.0.1:9050"

# Or use CLI flag
./sagent images "query" -P socks5h://127.0.0.1:9050
```

## Troubleshooting

### Tor not starting

```bash
# Check Tor status
sudo systemctl status tor

# View logs
sudo journalctl -u tor -f
```

### Port 9150 not listening

```bash
# Tor default port is 9050, not 9150
# Check which port Tor uses
grep -i "SocksPort" /etc/tor/torrc

# Update TOR_PROXY_URL accordingly
export TOR_PROXY_URL="socks5h://127.0.0.1:9050"
```

### Still getting 403

```bash
# Restart Tor
sudo systemctl restart tor

# Wait a few seconds for circuit establishment
sleep 5

# Try again
./sagent images "query" -V
```

## Without Tor (Limited Functionality)

Nếu không có Tor, bạn vẫn có thể dùng:

```bash
# ✅ Text search - Works without Tor
./sagent text "golang tutorial"

# ⚠️ Other categories - May return 403
./sagent images "query"  # Will fail without Tor
```

## Docker Setup (Alternative)

```bash
# Run with Tor proxy
docker run -d --name tor dperson/torproxy

# Use with search-agents
./sagent images "query" -P socks5h://172.17.0.2:9050
```
