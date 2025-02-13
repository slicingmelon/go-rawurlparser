// File: url.go
/*
Package rawurlparser provides URL parsing functionality that preserves exact URL paths.

Will update wait!
*/
package rawurlparser

import (
	"errors"
	"fmt"
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
	Hostname      string    // Just the hostname/domain (without port)
	Port          string    // Just the port (if specified)
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

	// Split host into hostname and port
	if result.Host != "" {
		if strings.HasPrefix(result.Host, "[") {
			// Handle IPv6 addresses
			closeBracket := strings.LastIndex(result.Host, "]")
			if closeBracket != -1 {
				result.Hostname = result.Host[:closeBracket+1] // Preserve brackets
				if len(result.Host) > closeBracket+1 && result.Host[closeBracket+1] == ':' {
					result.Port = result.Host[closeBracket+2:]
				}
			} else {
				result.Hostname = result.Host // Malformed IPv6, keep as-is
			}
		} else {
			// Handle IPv4 and regular hostnames
			host, port := SplitHostPort(result.Host)
			result.Hostname = host
			result.Port = port
		}
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
	} else {
		// Ensure Path is always set to "/" when empty
		result.Path = "/"
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

// The most basic function to quickly get the base URL using fmt.Sprintf
func (u *RawURL) BaseURL() string {
	return fmt.Sprintf("%s://%s", u.Scheme, u.Host)
}

// GetHostname returns the hostname without port.
// For IPv6 addresses, the square brackets are preserved.
func (u *RawURL) GetHostname() string {
	host := u.Host

	// Handle IPv6 addresses
	if strings.HasPrefix(host, "[") {
		if closeBracket := strings.LastIndex(host, "]"); closeBracket != -1 {
			// Return the IPv6 address with brackets
			if len(host) > closeBracket+1 && host[closeBracket+1] == ':' {
				return host[:closeBracket+1]
			}
			return host
		}
		return host // Malformed IPv6, return as-is
	}

	// Handle IPv4 and regular hostnames
	if i := strings.LastIndex(host, ":"); i != -1 {
		return host[:i]
	}
	return host
}

// GetPort returns the port part of the host.
// Returns empty string if no port is present.
func (u *RawURL) GetPort() string {
	host := u.Host

	// Handle IPv6 addresses
	if strings.HasPrefix(host, "[") {
		if closeBracket := strings.LastIndex(host, "]"); closeBracket != -1 {
			if len(host) > closeBracket+1 && host[closeBracket+1] == ':' {
				return host[closeBracket+2:] // Return everything after ]:
			}
			return ""
		}
		return ""
	}

	// Handle IPv4 and regular hostnames
	if i := strings.LastIndex(host, ":"); i != -1 {
		port := host[i+1:]
		// Validate port is numeric
		for _, b := range port {
			if b < '0' || b > '9' {
				return ""
			}
		}
		return port
	}
	return ""
}

/*
String() reconstructs the full URL from its components and returns a string representation

--->  scheme://host/path?query#fragment

	             userinfo      host      port    path		       query		            fragment
	            |------| |-------------| |--||---------------| |-------------------------| |-----------|
		https://john.doe@www.example.com:8092/forum/questions/?tag=networking&order=newest#fragmentation
		|----|  |---------------------------|
		scheme         authority
*/
func (u *RawURL) String() string {
	var buf strings.Builder

	// Scheme
	if u.Scheme != "" {
		buf.WriteString(u.Scheme)
		buf.WriteString("://")
	}

	// Authority (userinfo + host)
	buf.WriteString(GetAuthority(u))

	// Path
	buf.WriteString(u.Path)

	// Query
	if u.Query != "" {
		buf.WriteByte('?') // Use WriteByte for single-byte characters
		buf.WriteString(u.Query)
	}

	// Fragment
	if u.Fragment != "" {
		buf.WriteByte('#') // Use WriteByte for single-byte characters
		buf.WriteString(u.Fragment)
	}

	return buf.String()
}
