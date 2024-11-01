// File: url.go
/*
Package rawurlparser provides URL parsing functionality that preserves exact URL paths.

Note on HTTP Requests:
When using parsed URLs with Go's http.Client, you'll need to use URL.Opaque to preserve
the exact path encoding. Example:

    parsedURL := rawurlparser.Parse(rawURL)
    req := &http.Request{
        Method: "GET",
        URL: &url.URL{
            Scheme: parsedURL.Scheme,
            Host:   parsedURL.Host,
            Opaque: parsedURL.Path,  // Use Opaque to prevent path normalization
        },
    }
*/
package rawurlparse

import "strings"

// URL represents a raw URL with no normalization or encoding
type URL struct {
	Original string
	Scheme   string
	Host     string
	Path     string
	Query    string
	Fragment string
}

// String returns the original URL string
func (u *URL) String() string {
	return u.Original
}

// RawURLParse parses a URL string with no normalization or encoding
func RawURLParse(rawURL string) *URL {
	result := &URL{
		Original: rawURL,
	}

	remaining := rawURL

	// Get scheme
	if idx := strings.Index(remaining, "://"); idx != -1 {
		result.Scheme = remaining[:idx]
		remaining = remaining[idx+3:]
	}

	// Get fragment
	if idx := strings.Index(remaining, "#"); idx != -1 {
		result.Fragment = remaining[idx+1:]
		remaining = remaining[:idx]
	}

	// Get query
	if idx := strings.Index(remaining, "?"); idx != -1 {
		result.Query = remaining[idx+1:]
		remaining = remaining[:idx]
	}

	// Get path
	if idx := strings.Index(remaining, "/"); idx != -1 {
		result.Path = remaining[idx:]
		remaining = remaining[:idx]
	}

	// What remains is the host
	result.Host = remaining

	return result
}
