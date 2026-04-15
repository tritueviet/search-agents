// Package models contains shared data structures for DDGS.
package models

// TextResult represents a text search result.
type TextResult struct {
	Title string `json:"title"`
	Href  string `json:"href"`
	Body  string `json:"body"`
}

// ImagesResult represents an image search result.
type ImagesResult struct {
	Title     string `json:"title"`
	Image     string `json:"image"`
	Thumbnail string `json:"thumbnail"`
	URL       string `json:"url"`
	Height    string `json:"height"`
	Width     string `json:"width"`
	Source    string `json:"source"`
}

// NewsResult represents a news search result.
type NewsResult struct {
	Date   string `json:"date"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	URL    string `json:"url"`
	Image  string `json:"image"`
	Source string `json:"source"`
}

// VideosResult represents a video search result.
type VideosResult struct {
	Title       string            `json:"title"`
	Content     string            `json:"content"`
	Description string            `json:"description"`
	Duration    string            `json:"duration"`
	EmbedHTML   string            `json:"embed_html"`
	EmbedURL    string            `json:"embed_url"`
	ImageToken  string            `json:"image_token"`
	Images      map[string]string `json:"images"`
	Provider    string            `json:"provider"`
	Published   string            `json:"published"`
	Publisher   string            `json:"publisher"`
	Statistics  map[string]string `json:"statistics"`
	Uploader    string            `json:"uploader"`
}

// BooksResult represents a book search result.
type BooksResult struct {
	Title     string `json:"title"`
	Author    string `json:"author"`
	Publisher string `json:"publisher"`
	Info      string `json:"info"`
	URL       string `json:"url"`
	Thumbnail string `json:"thumbnail"`
}
