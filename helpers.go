package rawurlparser

import (
	"strings"
	"unicode/utf8"
)

// Helper methods //

// FullString reconstructs the URL from its components
func (u *RawURL) FullString() string {
	var buf strings.Builder

	if u.Scheme != "" {
		buf.WriteString(u.Scheme)
		buf.WriteString("://")
	}

	if u.User != nil {
		buf.WriteString(u.User.username)
		if u.User.passwordSet {
			buf.WriteRune(':')
			buf.WriteString(u.User.password)
		}
		buf.WriteRune('@')
	}

	buf.WriteString(u.Host)
	buf.WriteString(u.Path)

	if u.Query != "" {
		buf.WriteRune('?')
		buf.WriteString(u.Query)
	}

	if u.Fragment != "" {
		buf.WriteRune('#')
		buf.WriteString(u.Fragment)
	}

	return buf.String()
}

// ...
// https://media.beehiiv.com/cdn-cgi/image/fit=scale-down,format=auto,onerror=redirect,quality=80/uploads/asset/file/fac6de4c-8a4f-4688-aa43-0fbcf784426e/url-structure-and-scheme-2022.png
func (u *RawURL) GetFullRawURI() string {
	//..
}

func (u *RawURL) GetRawRelativeURI() string {
	//..
}

func (u *RawURL) GetRawURIPath() string {
	//..
}

// lastIndexRune returns the last index} of a rune in a string
func lastIndexRune(s string, r rune) int {
	// Fast path for ASCII
	if r < utf8.RuneSelf {
		return strings.LastIndex(s, string(r))
	}

	// For non-ASCII runes, we need to scan backwards
	for i := len(s); i > 0; {
		r1, size := utf8.DecodeLastRuneInString(s[:i])
		if r1 == r {
			return i - size
		}
		i -= size
	}
	return -1
}

// Port returns the port of the URL
func (u *RawURL) Port() string {
	if i := strings.LastIndex(u.Host, ":"); i != -1 {
		return u.Host[i+1:]
	}
	return ""
}

// Hostname returns the hostname of the URL (without port)
func (u *RawURL) Hostname() string {
	if i := strings.LastIndex(u.Host, ":"); i != -1 {
		return u.Host[:i]
	}
	return u.Host
}

// QueryValues returns a map of query parameters
func (u *RawURL) QueryValues() map[string][]string {
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
