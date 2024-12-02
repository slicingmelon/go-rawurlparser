package rawurlparser

import (
	"bytes"
	"encoding/hex"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Helper methods //

// URLComponent represents different parts of a URL that can be updated
type URLComponent int

const (
	Scheme URLComponent = iota
	Username
	Password
	Host
	Port
	Path
	Query
	Fragment
)

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

// GetRawScheme reconstructs the scheme from its components
func GetRawScheme(u *RawURL) string {
	if u.Scheme == "" {
		return ""
	}
	var buf strings.Builder
	buf.WriteString(u.Scheme)
	buf.WriteString("://")
	return buf.String()
}

// GetRawUserInfo reconstructs the userinfo from its components
func GetRawUserInfo(u *RawURL) string {
	if u.User == nil {
		return ""
	}
	var buf strings.Builder
	buf.WriteString(u.User.username)
	if u.User.passwordSet {
		buf.WriteRune(':')
		buf.WriteString(u.User.password)
	}
	buf.WriteRune('@')
	return buf.String()
}

// GetRawAuthority reconstructs the authority from its components
func GetRawAuthority(u *RawURL) string {
	var buf strings.Builder
	if u.User != nil {
		buf.WriteString(GetRawUserInfo(u))
	}
	buf.WriteString(u.Host)
	return buf.String()
}

// GetRawHost reconstructs the host of the URL (with port)
func GetRawHost(u *RawURL) string {
	var buf strings.Builder
	buf.WriteString(u.Host)
	return buf.String()
}

// GetRawHostname reconstructs the hostname of the URL (without port)
func GetRawHostname(u *RawURL) string {
	var buf strings.Builder
	if i := strings.LastIndex(u.Host, ":"); i != -1 {
		buf.WriteString(u.Host[:i])
	} else {
		buf.WriteString(u.Host)
	}
	return buf.String()
}

// GetRawPort reconstructs the port of the URL
func GetRawPort(u *RawURL) string {
	var buf strings.Builder
	if i := strings.LastIndex(u.Host, ":"); i != -1 {
		buf.WriteString(u.Host[i+1:])
		return buf.String()
	}
	return ""
}

// GetRawPath reconstructs the path from its components
func GetRawPath(u *RawURL) string {
	var buf strings.Builder
	if u.Path == "" {
		buf.WriteString("/")
		return buf.String()
	}

	// Check first byte directly for '/'
	if len(u.Path) > 0 && u.Path[0] != '/' {
		buf.WriteString("/")
	}
	buf.WriteString(u.Path)
	return buf.String()
}

// GetRawPathUnsafe reconstructs the path from its components
// Similar to GetRawPath but will omit first / in path
// Might be needed when fuzzing full paths
func GetRawPathUnsafe(u *RawURL) string {
	if u.Path == "" {
		return ""
	}

	var buf strings.Builder
	// Skip first char if it's a '/'
	if len(u.Path) > 0 {
		if u.Path[0] == '/' {
			buf.WriteString(u.Path[1:])
		} else {
			buf.WriteString(u.Path)
		}
	}
	return buf.String()
}

// GetRawQuery reconstructs the query from its components
func GetRawQuery(u *RawURL) string {
	if u.Query == "" {
		return ""
	}
	var buf strings.Builder
	buf.WriteRune('?')
	buf.WriteString(u.Query)
	return buf.String()
}

// GetRawFragment reconstructs the fragment from its components
func GetRawFragment(u *RawURL) string {
	if u.Fragment == "" {
		return ""
	}
	var buf strings.Builder
	buf.WriteRune('#')
	buf.WriteString(u.Fragment)
	return buf.String()
}

// QueryValues returns a map of query parameters
func (u *RawURL) GetRawQueryValues() map[string][]string {
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

/*
GetFullRawURL reconstructs the full URL from its components

--->  scheme://host/path?query#fragment

	             userinfo      host      port    path		       query		            fragment
	            |------| |-------------| |--||---------------| |-------------------------| |-----------|
		https://john.doe@www.example.com:8092/forum/questions/?tag=networking&order=newest#fragmentation
		|----|  |---------------------------|
		scheme         authority
*/
func (u *RawURL) GetFullRawURL() string {
	var buf strings.Builder

	// Scheme
	if u.Scheme != "" {
		buf.WriteString(u.Scheme)
		buf.WriteString("://")
	}

	// Authority (userinfo + host)
	buf.WriteString(GetRawAuthority(u))

	// Path
	buf.WriteString(GetRawPath(u))

	// Query
	if u.Query != "" {
		buf.WriteRune('?')
		buf.WriteString(u.Query)
	}

	// Fragment
	if u.Fragment != "" {
		buf.WriteRune('#')
		buf.WriteString(u.Fragment)
	}

	return buf.String()
}

// GetFullRawURI returns scheme://host/path
func (u *RawURL) GetFullRawURI() string {
	var buf strings.Builder

	// Scheme
	if u.Scheme != "" {
		buf.WriteString(u.Scheme)
		buf.WriteString("://")
	}

	// Authority
	buf.WriteString(GetRawAuthority(u))

	// Path
	buf.WriteString(GetRawPath(u))

	return buf.String()
}

// GetRawRelativeURI returns /path?query#fragment
func (u *RawURL) GetRawRelativeURI() string {
	var buf strings.Builder

	// Path
	buf.WriteString(GetRawPath(u))

	// Query
	if u.Query != "" {
		buf.WriteRune('?')
		buf.WriteString(u.Query)
	}

	// Fragment
	if u.Fragment != "" {
		buf.WriteRune('#')
		buf.WriteString(u.Fragment)
	}

	return buf.String()
}

// GetRawURIPath returns just the path component
func (u *RawURL) GetRawURIPath() string {
	return GetRawPath(u)
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

// UpdateRawURL updates a specific component of the URL with a new value
func (u *RawURL) UpdateRawURL(component URLComponent, newValue string) {
	switch component {
	case Scheme:
		u.Scheme = newValue
	case Username:
		if u.User == nil {
			u.User = &Userinfo{}
		}
		u.User.username = newValue
	case Password:
		if u.User == nil {
			u.User = &Userinfo{}
		}
		u.User.password = newValue
		u.User.passwordSet = true
	case Host:
		// Update host without affecting port
		if port := GetRawPort(u); port != "" {
			u.Host = newValue + ":" + port
		} else {
			u.Host = newValue
		}
	case Port:
		// Update port without affecting host
		hostname := GetRawHostname(u)
		if newValue != "" {
			u.Host = hostname + ":" + newValue
		} else {
			u.Host = hostname
		}
	case Path:
		u.Path = newValue
	case Query:
		u.Query = newValue
	case Fragment:
		u.Fragment = newValue
	}
}

// GetAsciiHex returns hex value of ascii char
func GetAsciiHex(r rune) string {
	val := strconv.FormatInt(int64(r), 16)
	if len(val) == 1 {
		// append 0 formatInt skips it by default
		val = "0" + val
	}
	return strings.ToUpper(val)
}

// GetUTF8Hex returns hex value of utf-8 non-ascii char
func GetUTF8Hex(r rune) string {
	// Percent Encoding is only done in hexadecimal values and in ASCII Range only
	// other UTF-8 chars (chinese etc) can be used by utf-8 encoding and byte conversion
	// let golang do utf-8 encoding of rune
	var buff bytes.Buffer
	utfchar := string(r)
	hexencstr := hex.EncodeToString([]byte(utfchar))
	for k, v := range hexencstr {
		if k != 0 && k%2 == 0 {
			buff.WriteRune('%')
		}
		buff.WriteRune(v)
	}
	return buff.String()
}

// GetRuneMap returns a map of runes
func GetRuneMap(runes []rune) map[rune]struct{} {
	x := map[rune]struct{}{}
	for _, v := range runes {
		x[v] = struct{}{}
	}
	return x
}
