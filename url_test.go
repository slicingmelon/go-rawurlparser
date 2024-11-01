// File: url_test.go
package rawurlparse

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"testing"
)

func TestRawURLParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		scheme   string
		host     string
		path     string
		query    string
		fragment string
	}{
		{
			name:     "Basic URL",
			input:    "https://example.com/path",
			scheme:   "https",
			host:     "example.com",
			path:     "/path",
			query:    "",
			fragment: "",
		},
		{
			name:     "URL with everything",
			input:    "https://example.com/path1/path2%2f;/?q=1#f",
			scheme:   "https",
			host:     "example.com",
			path:     "/path1/path2%2f;/",
			query:    "q=1",
			fragment: "f",
		},
		{
			name:     "Path traversal raw",
			input:    "http://example.com/path1/..%2f/",
			scheme:   "http",
			host:     "example.com",
			path:     "/path1/..%2f/",
			query:    "",
			fragment: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RawURLParse(tt.input)
			if got.Scheme != tt.scheme {
				t.Errorf("Scheme = %v, want %v", got.Scheme, tt.scheme)
			}
			if got.Host != tt.host {
				t.Errorf("Host = %v, want %v", got.Host, tt.host)
			}
			if got.Path != tt.path {
				t.Errorf("Path = %v, want %v", got.Path, tt.path)
			}
			if got.Query != tt.query {
				t.Errorf("Query = %v, want %v", got.Query, tt.query)
			}
			if got.Fragment != tt.fragment {
				t.Errorf("Fragment = %v, want %v", got.Fragment, tt.fragment)
			}
		})
	}
}

func ExampleRawURLParse_httpRequest() {
	// Example of using RawURLParse with http.Request while preserving raw paths
	rawURLs := []string{
		"https://example.com/path1/..%2f/path2",
		"https://example.com/path1;/%2e%2e/path2",
		"https://example.com/path1/path2;/。。/path3",
		"https://example.com/path1%2f..%2f/path2#fragment?query=1",
	}

	for _, rawURL := range rawURLs {
		fmt.Printf("\nTesting URL: %s\n", rawURL)
		fmt.Printf("----------------------------------------\n")

		// Parse with rawurlparser
		// Don't use URL.String() for debugging, as it will apply encodings and mess up your url
		parsedURL := RawURLParse(rawURL)
		log.Printf("RawURLParse result - Scheme: %s, Host: %s, Path: %s, Query: %s, Fragment: %s",
			parsedURL.Scheme,
			parsedURL.Host,
			parsedURL.Path,
			parsedURL.Query,
			parsedURL.Fragment,
		)

		// Create http.Request to demonstrate usage
		req := &http.Request{
			Method: "GET",
			URL: &url.URL{
				Scheme:   parsedURL.Scheme,
				Host:     parsedURL.Host,
				Path:     parsedURL.Path,     // Standard path
				RawPath:  parsedURL.Path,     // Raw path version
				Opaque:   parsedURL.Path,     // Opaque to prevent normalization
				RawQuery: parsedURL.Query,    // Raw query string
				Fragment: parsedURL.Fragment, // Fragment identifier
			},
		}

		// Log the request components for comparison
		log.Printf("http.Request URL components:")
		log.Printf("- Scheme:   %s", req.URL.Scheme)
		log.Printf("- Host:     %s", req.URL.Host)
		log.Printf("- Path:     %s", req.URL.Path)
		log.Printf("- RawPath:  %s", req.URL.RawPath)
		log.Printf("- Opaque:   %s", req.URL.Opaque)
		log.Printf("- RawQuery: %s", req.URL.RawQuery)
		log.Printf("- Fragment: %s", req.URL.Fragment)

		// Compare with standard url.Parse
		stdURL, err := url.Parse(rawURL)
		if err != nil {
			log.Printf("Standard url.Parse error: %v", err)
		} else {
			log.Printf("\nComparison with standard url.Parse:")
			log.Printf("- Path:     %s", stdURL.Path)     // Will be normalized
			log.Printf("- RawPath:  %s", stdURL.RawPath)  // May preserve some raw encoding
			log.Printf("- Opaque:   %s", stdURL.Opaque)   // Opaque part
			log.Printf("- String(): %s", stdURL.String()) // Full URL string
		}
	}
}

// Example of comparing different URL parsing approaches
func ExampleRawURLParse_comparison() {
	rawURL := "https://example.com/path1/..%2f/path2;/%2e%2e/test"

	// Our raw parser
	raw := RawURLParse(rawURL)
	fmt.Printf("RawURLParse:\n")
	fmt.Printf("Path: %s\n", raw.Path)

	// Standard parser
	std, _ := url.Parse(rawURL)
	fmt.Printf("\nstandard url.Parse:\n")
	fmt.Printf("Path: %s\n", std.Path)
	fmt.Printf("RawPath: %s\n", std.RawPath)

	// http.Request usage (raw)
	req := &http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme: raw.Scheme,
			Host:   raw.Host,
			Path:   raw.Path,
			Opaque: raw.Path, // Use Opaque to prevent normalization
		},
	}

	fmt.Printf("\nhttp.Request URL:\n")
	fmt.Printf("Path: %s\n", req.URL.Path)
	fmt.Printf("Opaque: %s\n", req.URL.Opaque)
}
