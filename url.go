// File: url.go
/*
Package rawurlparser provides URL parsing functionality that preserves exact URL paths.

Note on HTTP Requests:
When using parsed URLs with Go's http.Client, you'll need to use URL.Opaque to preserve
the exact path encoding. Example:

    parsedURL := rawurlparser.RawURLParse(rawURL)
    req := &http.Request{
        Method: "GET",
        URL: &url.URL{
            Scheme: parsedURL.Scheme,
            Host:   parsedURL.Host,
            Opaque: parsedURL.Path,  // Use Opaque to prevent path normalization
        },
    }
*/
package rawurlparser

import (
	"errors"
	"strings"
)

// URL represents a raw URL with no normalization or encoding.
// It preserves the exact format of the original URL string,
// including any percent-encoding or special characters.
type URL struct {
	Original string    // The original, unmodified URL string
	Scheme   string    // The URL scheme (e.g., "http", "https")
	Opaque   string    // For non-hierarchical URLs (e.g., mailto:user@example.com)
	User     *Userinfo // username and password information
	Host     string    // The host component
	Path     string    // The path component, exactly as provided
	Query    string    // The query string without the leading '?'
	Fragment string    // The fragment without the leading '#'
}

// Userinfo stores username and password info
type Userinfo struct {
	username    string
	password    string
	passwordSet bool
}

// String returns the original URL string
func (u *URL) String() string {
	return u.Original
}

// RawURLParseWithError is like RawURLParse but returns an error if URL is invalid
func RawURLParseWithError(rawURL string) (*URL, error) {
	if len(rawURL) == 0 {
		return nil, errors.New("empty URL")
	}

	result := &URL{
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
	} else {
		return nil, errors.New("missing scheme (e.g., http:// or https://)")
	}

	// Get userinfo if present
	if idx := strings.Index(remaining, "@"); idx != -1 {
		userinfo := remaining[:idx]
		remaining = remaining[idx+1:]

		// Split username and password
		if pwIdx := strings.Index(userinfo, ":"); pwIdx != -1 {
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
	return result, nil
}

// Keep the original function for backward compatibility
func RawURLParse(rawURL string) *URL {
	result, _ := RawURLParseWithError(rawURL)
	return result
}

// Helper methods
//
// FullString reconstructs the URL from its components
func (u *URL) FullString() string {
	var buf strings.Builder

	if u.Scheme != "" {
		buf.WriteString(u.Scheme)
		buf.WriteString("://")
	}

	if u.User != nil {
		buf.WriteString(u.User.username)
		if u.User.passwordSet {
			buf.WriteByte(':')
			buf.WriteString(u.User.password)
		}
		buf.WriteByte('@')
	}

	buf.WriteString(u.Host)
	buf.WriteString(u.Path)

	if u.Query != "" {
		buf.WriteByte('?')
		buf.WriteString(u.Query)
	}

	if u.Fragment != "" {
		buf.WriteByte('#')
		buf.WriteString(u.Fragment)
	}

	return buf.String()
}

func (u *URL) Port() string {
	if i := strings.LastIndex(u.Host, ":"); i != -1 {
		return u.Host[i+1:]
	}
	return ""
}

func (u *URL) Hostname() string {
	if i := strings.LastIndex(u.Host, ":"); i != -1 {
		return u.Host[:i]
	}
	return u.Host
}

func (u *URL) QueryValues() map[string][]string {
	values := make(map[string][]string)
	for _, pair := range strings.Split(u.Query, "&") {
		if pair == "" {
			continue
		}
		kv := strings.SplitN(pair, "=", 2)
		key := kv[0]
		value := ""
		if len(kv) == 2 {
			value = kv[1]
		}
		values[key] = append(values[key], value)
	}
	return values
}
