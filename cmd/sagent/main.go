// Package main is the CLI entry point for Search Agents (sagent).
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/tritueviet/search-agents/internal/engine"
	"github.com/tritueviet/search-agents/internal/extractor"
	"github.com/tritueviet/search-agents/internal/httpclient"
	"github.com/tritueviet/search-agents/mcp"
	"github.com/tritueviet/search-agents/pkg/searchagents"
)

var (
	version   = "0.1.0"
	proxy     string
	timeout   int
	verifySSL bool
	verbose   bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "sagent",
		Short: "sagent - Search Agents CLI",
		Long:  "A metasearch library that aggregates results from diverse web search services.",
	}

	// Version command
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}

	// Text search command
	textCmd := &cobra.Command{
		Use:   "text [query]",
		Short: "Perform text search",
		Args:  cobra.MinimumNArgs(1),
		RunE:  runTextSearch,
	}

	addSearchFlags(textCmd)

	// Extract command
	extractCmd := &cobra.Command{
		Use:   "extract [url]",
		Short: "Extract content from URL",
		Args:  cobra.ExactArgs(1),
		RunE:  runExtract,
	}

	extractCmd.Flags().StringP("format", "f", "text_markdown", "Output format")
	extractCmd.Flags().StringP("output", "o", "", "Output file")

	// MCP command
	mcpCmd := &cobra.Command{
		Use:   "mcp",
		Short: "Start MCP server (stdio transport)",
		RunE:  runMCP,
	}

	// Category commands
	imagesCmd := &cobra.Command{
		Use:   "images [query]",
		Short: "Perform image search",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSearch(cmd, args, "images")
		},
	}
	addSearchFlags(imagesCmd)

	videosCmd := &cobra.Command{
		Use:   "videos [query]",
		Short: "Perform video search",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSearch(cmd, args, "videos")
		},
	}
	addSearchFlags(videosCmd)

	newsCmd := &cobra.Command{
		Use:   "news [query]",
		Short: "Perform news search",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSearch(cmd, args, "news")
		},
	}
	addSearchFlags(newsCmd)

	booksCmd := &cobra.Command{
		Use:   "books [query]",
		Short: "Perform book search",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSearch(cmd, args, "books")
		},
	}
	addSearchFlags(booksCmd)

	// Global flags
	rootCmd.PersistentFlags().StringVarP(&proxy, "proxy", "P", "", "Proxy URL")
	rootCmd.PersistentFlags().IntVarP(&timeout, "timeout", "T", 10, "Timeout in seconds")
	rootCmd.PersistentFlags().BoolVarP(&verifySSL, "verify", "v", true, "Verify SSL")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "V", false, "Verbose output (show errors)")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(textCmd)
	rootCmd.AddCommand(extractCmd)
	rootCmd.AddCommand(mcpCmd)
	rootCmd.AddCommand(imagesCmd)
	rootCmd.AddCommand(videosCmd)
	rootCmd.AddCommand(newsCmd)
	rootCmd.AddCommand(booksCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func addSearchFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("region", "r", "us-en", "Region (e.g., us-en, ru-ru)")
	cmd.Flags().StringP("safesearch", "s", "moderate", "SafeSearch (on, moderate, off)")
	cmd.Flags().StringP("timelimit", "t", "", "Time limit (d, w, m, y)")
	cmd.Flags().IntP("max-results", "m", 10, "Maximum number of results")
	cmd.Flags().IntP("page", "p", 1, "Page number")
	cmd.Flags().StringP("backend", "b", "auto", "Backend(s) to use")
	cmd.Flags().StringP("output", "o", "", "Output file (JSON or CSV)")
	cmd.Flags().BoolP("extract", "e", false, "Extract content from result URLs")
	cmd.Flags().String("extract-format", "text_markdown", "Extract format")
}

func runTextSearch(cmd *cobra.Command, args []string) error {
	return runSearch(cmd, args, "text")
}

func runSearch(cmd *cobra.Command, args []string, category string) error {
	query := strings.Join(args, " ")

	region, _ := cmd.Flags().GetString("region")
	safesearch, _ := cmd.Flags().GetString("safesearch")
	timelimit, _ := cmd.Flags().GetString("timelimit")
	maxResults, _ := cmd.Flags().GetInt("max-results")
	page, _ := cmd.Flags().GetInt("page")
	output, _ := cmd.Flags().GetString("output")

	client, err := searchagents.New(searchagents.Options{
		Proxy:   proxy,
		Timeout: timeout,
		Verify:  verifySSL,
	})
	if err != nil {
		return fmt.Errorf("failed to create SearchAgents: %w", err)
	}

	opts := engine.DefaultSearchOptions()
	opts.Region = region
	opts.SafeSearch = safesearch
	opts.TimeLimit = timelimit
	opts.Page = page
	opts.Extra["max_results"] = fmt.Sprintf("%d", maxResults)

	ctx := cmd.Context()
	var results []map[string]string

	switch category {
	case "text":
		results, err = client.Text(ctx, query, opts)
	case "images":
		results, err = client.Images(ctx, query, opts)
	case "videos":
		results, err = client.Videos(ctx, query, opts)
	case "news":
		results, err = client.News(ctx, query, opts)
	case "books":
		results, err = client.Books(ctx, query, opts)
	default:
		return fmt.Errorf("unknown category: %s", category)
	}

	if err != nil {
		if verbose {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			fmt.Fprintf(os.Stderr, "Query: %s, Category: %s, Region: %s\n", query, category, region)
		}
		return fmt.Errorf("search failed: %v", err)
	}

	if output != "" {
		return saveResults(output, results)
	}

	return printResults(results)
}

func runExtract(cmd *cobra.Command, args []string) error {
	url := args[0]
	format, _ := cmd.Flags().GetString("format")
	output, _ := cmd.Flags().GetString("output")

	httpOpts := httpclient.Options{
		Proxy:   proxy,
		Timeout: time.Duration(timeout) * time.Second,
		Verify:  verifySSL,
	}
	client, err := httpclient.NewClient(httpOpts)
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}

	ext := extractor.New(client)
	ctx := cmd.Context()
	result, err := ext.Extract(ctx, url, format)
	if err != nil {
		return fmt.Errorf("extract failed: %w", err)
	}

	if output != "" {
		return saveResults(output, []map[string]string{
			{"url": result["url"].(string), "content": fmt.Sprintf("%v", result["content"])},
		})
	}

	fmt.Printf("URL: %s\n\n", result["url"])
	fmt.Printf("%v\n", result["content"])
	return nil
}

func runMCP(cmd *cobra.Command, args []string) error {
	client, err := searchagents.New(searchagents.Options{
		Proxy:   proxy,
		Timeout: timeout,
		Verify:  verifySSL,
	})
	if err != nil {
		return fmt.Errorf("failed to create SearchAgents: %w", err)
	}

	httpOpts := httpclient.Options{
		Proxy:   proxy,
		Timeout: time.Duration(timeout) * time.Second,
		Verify:  verifySSL,
	}
	httpClient, err := httpclient.NewClient(httpOpts)
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}

	ext := extractor.New(httpClient)
	server := mcp.NewServer(client, ext)
	return server.RunStdioAsync()
}

func printResults(results []map[string]string) error {
	for i, result := range results {
		fmt.Printf("%d.\t%s\n", i+1, strings.Repeat("=", 78))
		for key, value := range result {
			if value != "" {
				fmt.Printf("%-15s%s\n", key+":", value)
			}
		}
		fmt.Println()
	}
	return nil
}

func saveResults(filename string, results []map[string]string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(results)
}
