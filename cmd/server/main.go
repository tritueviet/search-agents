// Package main is the API server entry point.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tritueviet/search-agents/api"
	"github.com/tritueviet/search-agents/internal/extractor"
	"github.com/tritueviet/search-agents/internal/httpclient"
	"github.com/tritueviet/search-agents/pkg/searchagents"
)

func main() {
	host := flag.String("host", "0.0.0.0", "Host to bind")
	port := flag.Int("port", 8000, "Port to bind")
	proxy := flag.String("proxy", "", "Proxy URL")
	timeout := flag.Int("timeout", 5, "Timeout in seconds")
	flag.Parse()

	// Create SearchAgents client
	client, err := searchagents.New(searchagents.Options{
		Proxy:   *proxy,
		Timeout: *timeout,
		Verify:  true,
	})
	if err != nil {
		log.Fatalf("Failed to create SearchAgents client: %v", err)
	}

	// Create HTTP client for extractor
	httpClient, err := httpclient.NewClient(httpclient.Options{
		Proxy:   *proxy,
		Timeout: time.Duration(*timeout) * time.Second,
		Verify:  true,
	})
	if err != nil {
		log.Fatalf("Failed to create HTTP client: %v", err)
	}

	// Create extractor
	ext := extractor.New(httpClient)

	// Create API server
	server := api.NewServer(client, ext)

	addr := fmt.Sprintf("%s:%d", *host, *port)
	fmt.Printf("Starting Search Agents API on http://%s\n", addr)
	fmt.Printf("Health check: http://%s/health\n", addr)
	fmt.Printf("API docs: http://%s/docs\n", addr)

	// Handle graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		fmt.Println("\nShutting down server...")
		os.Exit(0)
	}()

	// Start server
	if err := server.Run(addr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
