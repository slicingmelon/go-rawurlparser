package rawurlparser

import (
	"bytes"
	"encoding/hex"
	"strconv"
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

// GetRawScheme reconstructs the scheme from its components
func GetRawSche(u *RawURL) string {
	// xx
	//return buf.string()
}

// GetRawUserInfo reconstructs the userinfo from its components
func GetRawUserInfo(u *RawURL) string {
	var buf strings.Builder
	// xx
	//return buf.string()
}

// GetRawAuthority reconstructs the authority from its components
func GetRawAuthority(u *RawURL) string {
	return u.Host
}

// GetRawHostname reconstructs the hostname of the URL (without port)
func GetRawHostname(u *RawURL) string {
	if i := strings.LastIndex(u.Host, ":"); i != -1 {
		return u.Host[:i]
	}
	return u.Host
}

// GetRawPort reconstructs the port of the URL
func GetRawPort(u *RawURL) string {
	if i := strings.LastIndex(u.Host, ":"); i != -1 {
		return u.Host[i+1:]
	}
	return ""
}

// GetRawPath reconstructs the path from its components
func GetRawPath(u *RawURL) string {
	// xx
}

// GetRawPathUnsafe reconstructs the path from its components
// Similar to GetRawPath but will omit first / in path
func GetRawPathUnsafe(u *RawURL) string {
	// xx
}

// GetRawQuery reconstructs the query from its components
func GetRawQuery(u *RawURL) string {
	// xx
}

// GetRawFragment reconstructs the fragment from its components
func GetRawFragment(u *RawURL) string {
	// xx
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

	             userinfo      host       port		path		       query		     fragment
	            |------ ||--------------||---||--------------||-------------------------||---|
		https://john.doe@www.example.com:123/forum/questions/?tag=networking&order=newest#top
		|----|  |---------------------------|
		scheme         authority
*/
func (u *RawURL) GetFullRawURL() string {
	//..
}

// ...
// scheme://host/path
func (u *RawURL) GetFullRawURI() string {
	//..
	return u.GetFullRawURL()
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
