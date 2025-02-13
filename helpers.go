package rawurlparser

import (
	"bytes"
	"encoding/hex"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Helper methods //

// URLBuilder represents a mutable URL structure for manipulation
// WIP
// type RawURLBuilder struct {
// 	*RawURL           // Embed the original RawURL
// 	workingURI string // Working copy of RequestURI
// }

// // NewURLBuilder creates a new builder from RawURL
// func NewRawURLBuilder(u *RawURL) *RawURLBuilder {
// 	return &RawURLBuilder{
// 		RawURL:     u,
// 		workingURI: u.RawRequestURI,
// 	}
// }

// GetScheme reconstructs the scheme from its components and returns a string representation
func GetScheme(u *RawURL) string {
	if u.Scheme == "" {
		return ""
	}
	var buf strings.Builder
	buf.WriteString(u.Scheme)
	buf.WriteString("://")
	return buf.String()
}

// GetUserInfo reconstructs the userinfo from its components
func GetUserInfo(u *RawURL) string {
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

// GetAuthority reconstructs the authority from its components and returns a string representation
func GetAuthority(u *RawURL) string {
	var buf strings.Builder

	// Add userinfo if present
	if u.User != nil {
		buf.WriteString(u.User.username)
		if u.User.passwordSet {
			buf.WriteByte(':')
			buf.WriteString(u.User.password)
		}
		buf.WriteByte('@')
	}

	// Add host
	buf.WriteString(u.Host)

	return buf.String()
}

// GetHost returns the string representation of the Host
func GetHost(u *RawURL) string {
	var buf strings.Builder
	buf.WriteString(u.Host)
	return buf.String()
}

// QueryValues returns a map of query parameters
func (u *RawURL) GetQueryValues() map[string][]string {
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

// GetRawRequestURI returns the exact URI as it would appear in an HTTP request line
// It could be "//a/b/c../;/?x=test", "\\a\\bb\\..\\test//..//..//users.json", "@collaboratorhost", etc
func (u *RawURL) GetRawRequestURI() string {
	if u.RawRequestURI != "" {
		return u.RawRequestURI
	}

	var buf strings.Builder

	// If no custom RawRequestURI is set, construct from Path
	if u.Path != "" {
		buf.WriteString(u.Path)
	}

	// Add query if present
	if u.Query != "" {
		buf.WriteRune('?')
		buf.WriteString(u.Query)
	}

	// Add fragment if present
	if u.Fragment != "" {
		buf.WriteRune('#')
		buf.WriteString(u.Fragment)
	}

	return buf.String()
}

// SplitHostPort() separates host and port. If the port is not valid, it returns
// the entire input as host, and it doesn't check the validity of the host.
// Unlike net.SplitHostPort, but per RFC 3986, it requires ports to be numeric.
// splitHostPort separates host and port while handling IPv6 addresses
func SplitHostPort(hostPort string) (host, port string) {
	// Handle IPv6 addresses
	if strings.HasPrefix(hostPort, "[") {
		closeBracket := strings.Index(hostPort, "]")
		if closeBracket != -1 {
			host = hostPort[1:closeBracket]
			if len(hostPort) > closeBracket+1 && hostPort[closeBracket+1] == ':' {
				port = hostPort[closeBracket+2:]
			}
			return
		}
	}

	// Handle regular host:port
	colon := strings.LastIndex(hostPort, ":")
	if colon != -1 && validOptionalPort(hostPort[colon:]) {
		host = hostPort[:colon]
		port = hostPort[colon+1:]
	} else {
		host = hostPort
	}
	return
}

// validOptionalPort reports whether port is either an empty string
// or matches /^:\d+$/
func validOptionalPort(port string) bool {
	if port == "" {
		return true
	}
	if port[0] != ':' {
		return false
	}
	for _, c := range port[1:] {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
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

// GetRuneMap returns a map of runes
func GetRuneMap(runes []rune) map[rune]struct{} {
	x := map[rune]struct{}{}
	for _, v := range runes {
		x[v] = struct{}{}
	}
	return x
}
