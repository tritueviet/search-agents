#!/bin/bash
# Start Search Agents API Server

set -e

HOST=${SAGENT_HOST:-0.0.0.0}
PORT=${SAGENT_PORT:-8000}
PROXY=${SAGENT_PROXY:-}
TIMEOUT=${SAGENT_TIMEOUT:-5}

echo "==================================="
echo "  Search Agents API Server"
echo "==================================="
echo "Host: $HOST"
echo "Port: $PORT"
[ -n "$PROXY" ] && echo "Proxy: $PROXY"
echo "Timeout: ${TIMEOUT}s"
echo "==================================="
echo ""

exec ./sagent-api --host "$HOST" --port "$PORT" --proxy "$PROXY" --timeout "$TIMEOUT"
