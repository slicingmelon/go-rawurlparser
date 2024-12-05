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
var (
	ErrEmptyURL   = errors.New("empty URL")
	ErrInvalidURL = errors.New("invalid URL format")
)

// RawURL represents a raw URL with no normalization or encoding
type RawURL struct {
	Original      string    // The original, unmodified URL string
	Scheme        string    // The URL scheme (e.g., "http", "https")
	Opaque        string    // For non-hierarchical URLs (e.g., mailto:user@example.com)
	User          *Userinfo // username and password information
	Host          string    // The host component (hostname + port)
	Path          string    // The path component, exactly as provided
	Query         string    // The query string without the leading '?'
	Fragment      string    // The fragment without the leading '#'
	RawRequestURI string    // Everything after host: /path?query#fragment
}

// Userinfo stores username and password info
type Userinfo struct {
	username    string
	password    string
	passwordSet bool
}

// ParseOptions contains configuration options for URL parsing
type ParseOptions struct {
	FallbackScheme     string // Default scheme if none provided
	AllowMissingScheme bool   // If true, uses FallbackScheme when scheme is missing
}

// DefaultOptions returns the default parsing options
func DefaultOptions() *ParseOptions {
	return &ParseOptions{
		FallbackScheme:     "https",
		AllowMissingScheme: true,
	}
}

// RawURLParseWithOptions parses URL with custom options
func RawURLParseWithOptions(rawURL string, opts *ParseOptions) (*RawURL, error) {
	if len(rawURL) == 0 {
		return nil, ErrEmptyURL
	}

	result := &RawURL{
		Original: rawURL,
	}

	// Handle scheme
	schemeEnd := strings.Index(rawURL, "://")
	remaining := rawURL

	if schemeEnd != -1 {
		result.Scheme = rawURL[:schemeEnd]
		remaining = rawURL[schemeEnd+3:]
	} else {
		// Check for scheme without //
		if colonIndex := strings.Index(rawURL, ":"); colonIndex != -1 {
			beforeColon := rawURL[:colonIndex]
			if !strings.Contains(beforeColon, "/") && !strings.Contains(beforeColon, ".") {
				result.Scheme = beforeColon
				result.Opaque = rawURL[colonIndex+1:]
				return result, nil
			}
		}

		// Apply fallback scheme if configured
		if opts != nil && opts.AllowMissingScheme {
			result.Scheme = opts.FallbackScheme
		}
	}

	// Split authority (host + optional userinfo) from path
	pathStart := strings.Index(remaining, "/")
	authority := remaining
	if pathStart != -1 {
		authority = remaining[:pathStart]
		remaining = remaining[pathStart:]
	} else {
		remaining = "/"
	}

	// Parse authority (user:pass@host:port)
	if atIndex := strings.Index(authority, "@"); atIndex != -1 {
		userinfo := authority[:atIndex]
		authority = authority[atIndex+1:]

		result.User = &Userinfo{}
		if colonIndex := strings.Index(userinfo, ":"); colonIndex != -1 {
			result.User.username = userinfo[:colonIndex]
			result.User.password = userinfo[colonIndex+1:]
			result.User.passwordSet = true
		} else {
			result.User.username = userinfo
		}
	}

	// Handle IPv6 addresses
	if strings.HasPrefix(authority, "[") {
		closeBracket := strings.LastIndex(authority, "]")
		if closeBracket == -1 {
			return nil, ErrInvalidURL
		}

		// Get the IPv6 address part
		result.Host = authority[:closeBracket+1]

		// Check for port after the IPv6 address
		if len(authority) > closeBracket+1 {
			if authority[closeBracket+1] == ':' {
				result.Host = authority // Include the full authority with port
			}
		}
	} else {
		// Handle IPv4 and regular hostnames
		result.Host = authority
	}

	// Parse path, query, and fragment
	if len(remaining) > 0 {
		// Extract fragment
		if hashIndex := strings.Index(remaining, "#"); hashIndex != -1 {
			result.Fragment = remaining[hashIndex+1:]
			remaining = remaining[:hashIndex]
		}

		// Extract query
		if queryIndex := strings.Index(remaining, "?"); queryIndex != -1 {
			result.Query = remaining[queryIndex+1:]
			remaining = remaining[:queryIndex]
		}

		// What's left is the path
		result.Path = remaining
	}

	// Build RawRequestURI
	result.RawRequestURI = result.Path
	if result.Query != "" {
		result.RawRequestURI += "?" + result.Query
	}
	if result.Fragment != "" {
		result.RawRequestURI += "#" + result.Fragment
	}

	return result, nil
}

// RawURLParse parses URL with default options
func RawURLParse(rawURL string) (*RawURL, error) {
	return RawURLParseWithOptions(rawURL, DefaultOptions())
}

// RawURLParseStrict parses URL without fallback scheme
func RawURLParseStrict(rawURL string) (*RawURL, error) {
	return RawURLParseWithOptions(rawURL, nil)
}

// Hostname() returns u.Host, stripping any valid port number if present.
// If the result is enclosed in square brackets, as literal IPv6 addresses are,
// the square brackets are removed from the result.
func (u *RawURL) Hostname() string {
	host, _ := splitHostPort(u.Host)
	return host
}

// Port() returns the port part of u.Host, without the leading colon.
// If u.Host doesn't contain a valid numeric port, Port returns an empty string.
func (u *RawURL) Port() string {
	_, port := splitHostPort(u.Host)
	return port
}

// splitHostPort() separates host and port. If the port is not valid, it returns
// the entire input as host, and it doesn't check the validity of the host.
// Unlike net.SplitHostPort, but per RFC 3986, it requires ports to be numeric.
func splitHostPort(hostPort string) (host, port string) {
	host = hostPort

	colon := strings.LastIndexByte(host, ':')
	if colon != -1 && validOptionalPort(host[colon:]) {
		host, port = host[:colon], host[colon+1:]
	}

	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		host = host[1 : len(host)-1]
	}

	return
}

// validOptionalPort reports whether port is either an empty string
// or matches /^:\d*$/
func validOptionalPort(port string) bool {
	if port == "" {
		return true
	}
	if port[0] != ':' {
		return false
	}
	for _, b := range port[1:] {
		if b < '0' || b > '9' {
			return false
		}
	}
	return true
}
