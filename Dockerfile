# Build stage
FROM golang:1.23-alpine AS builder

# Install dependencies
RUN apk add --no-cache git make

WORKDIR /build

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build API server
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /build/sagent-api ./cmd/server

# Final stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache ca-certificates curl tzdata

# Create non-root user
RUN addgroup -g 1001 appgroup && \
    adduser -u 1001 -G appgroup -s /bin/sh -D appuser

# Copy binary
COPY --from=builder /build/sagent-api /usr/local/bin/sagent-api

# Set working directory
WORKDIR /home/appuser

# Copy startup script
COPY start_api.sh ./start_api.sh
RUN chmod +x ./start_api.sh && \
    chown -R appuser:appgroup /home/appuser

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8000

# Environment variables
ENV SAGENT_HOST=0.0.0.0
ENV SAGENT_PORT=8000
ENV SAGENT_TIMEOUT=10

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:8000/health || exit 1

# Run
CMD ["./start_api.sh"]
