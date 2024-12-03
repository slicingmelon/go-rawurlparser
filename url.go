// File: url.go
/*
Package rawurlparser provides URL parsing functionality that preserves exact URL paths.

Will update wait!
*/
package rawurlparser

import (
	"errors"
	"strings"
)

// RawURL represents a raw URL with no normalization or encoding.
// It preserves the exact format of the original URL string
// including any percent-encoding, special characters, unicode chars, etc.
type RawURL struct {
	Original string    // The original, unmodified URL string
	Scheme   string    // The URL scheme (e.g., "http", "https")
	Opaque   string    // For non-hierarchical URLs (e.g., mailto:user@example.com)
	User     *Userinfo // username and password information
	Host     string    // The host component (hostname + port)
	Hostname string    // The hostname component
	Port     string    // The port component
	Path     string    // The path component, exactly as provided
	Query    string    // The query string without the leading '?'
	Fragment string    // The fragment without the leading '#'

	//other stuff
	RawRequestURI string // The exact URI as it appears in HTTP request line (will be needed to store fuzzed paths,payloads,etc)
}

// ParseOptions contains configuration options for URL parsing
type ParseOptions struct {
	FallbackScheme     string // Default scheme if none provided
	AllowMissingScheme bool   // If true, uses FallbackScheme when scheme is missing
}

// Userinfo stores username and password info
type Userinfo struct {
	username    string
	password    string
	passwordSet bool
}

// DefaultOptions returns the default parsing options
func DefaultOptions() *ParseOptions {
	return &ParseOptions{
		FallbackScheme:     "https",
		AllowMissingScheme: true,
	}
}

// String returns the Full URL string
func (u *RawURL) String() string {
	return u.GetRawFullURL()
}

// RawURLParseWithOptions parses URL with custom options
func RawURLParseWithOptions(rawURL string, opts *ParseOptions) (*RawURL, error) {
	if len(rawURL) == 0 {
		return nil, errors.New("empty URL")
	}

	result := &RawURL{
		Original: rawURL,
	}

	remaining := rawURL

	// Get scheme
	if idx := strings.Index(remaining, "://"); idx != -1 {
		result.Scheme = remaining[:idx]
		remaining = remaining[idx+3:]
	} else if idx := strings.Index(remaining, ":"); idx != -1 {
		// Handle opaque URLs (mailto:user@example.com)
		result.Scheme = remaining[:idx]
		result.Opaque = remaining[idx+1:]
		return result, nil
	} else if opts != nil && opts.AllowMissingScheme {
		// Apply fallback scheme
		result.Scheme = opts.FallbackScheme
	} else {
		return nil, errors.New("missing scheme (e.g., http:// or https://)")
	}

	// Get userinfo if present
	if idx := strings.IndexRune(remaining, '@'); idx != -1 {
		userinfo := remaining[:idx]
		remaining = remaining[idx+1:]

		// Split username and password
		if pwIdx := strings.IndexRune(userinfo, ':'); pwIdx != -1 {
			result.User = &Userinfo{
				username:    userinfo[:pwIdx],
				password:    userinfo[pwIdx+1:],
				passwordSet: true,
			}
		} else {
			result.User = &Userinfo{
				username: userinfo,
			}
		}
	}

	// Get fragment
	if idx := strings.IndexRune(remaining, '#'); idx != -1 {
		result.Fragment = remaining[idx+1:]
		remaining = remaining[:idx]
	}

	// Get query
	if idx := strings.IndexRune(remaining, '?'); idx != -1 {
		result.Query = remaining[idx+1:]
		remaining = remaining[:idx]
	}

	// Get path
	if idx := strings.IndexRune(remaining, '/'); idx != -1 {
		result.Path = remaining[idx:]
		remaining = remaining[:idx]
	}

	// What remains is the host
	result.Host = remaining

	// Parse hostname and port from host
	if idx := strings.LastIndex(result.Host, ":"); idx != -1 {
		result.Hostname = result.Host[:idx]
		result.Port = result.Host[idx+1:]
	} else {
		result.Hostname = result.Host
		result.Port = ""
	}

	return result, nil
}

// RawURLParse uses default options - core function
func RawURLParse(rawURL string) (*RawURL, error) {
	return RawURLParseWithOptions(rawURL, DefaultOptions())
}

// RawURLParseStrictScheme parses without fallback scheme
func RawURLParseStrictScheme(rawURL string) (*RawURL, error) {
	return RawURLParseWithOptions(rawURL, nil)
}
