package rawurlparser

import (
	"strings"
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
