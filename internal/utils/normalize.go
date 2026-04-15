// Package utils provides utility functions for DDGS.
package utils

import (
	"html"
	"net/url"
	"regexp"
	"strings"
	"unicode"
)

var stripTagsRegex = regexp.MustCompile("<.*?>")

// NormalizeURL unquotes a URL and replaces spaces with '+'.
func NormalizeURL(rawURL string) string {
	if rawURL == "" {
		return ""
	}
	decoded, err := url.QueryUnescape(rawURL)
	if err != nil {
		decoded = rawURL
	}
	return strings.ReplaceAll(decoded, " ", "+")
}

// NormalizeText normalizes text by stripping HTML tags, unescaping HTML entities,
// normalizing Unicode, removing control characters, and collapsing whitespace.
func NormalizeText(raw string) string {
	if raw == "" {
		return ""
	}

	// Strip HTML tags
	text := stripTagsRegex.ReplaceAllString(raw, "")

	// Unescape HTML entities
	text = html.UnescapeString(text)

	// Remove control characters and collapse whitespace
	var builder strings.Builder
	for _, ch := range text {
		if unicode.IsControl(ch) {
			continue
		}
		if unicode.IsSpace(ch) {
			if builder.Len() > 0 && builder.String()[builder.Len()-1] != ' ' {
				builder.WriteRune(' ')
			}
			continue
		}
		builder.WriteRune(ch)
	}

	return strings.TrimSpace(builder.String())
}
