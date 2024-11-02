package rawurlparser

import (
	"fmt"
	"testing"
)

func TestRawURLParse(t *testing.T) {
	tests := []string{
		"https://example.com/path?query=123#fragment",
		"http://user:pass@example.com:8080/path",
		"mailto:user@example.com",
		"https://api.example.com/path with spaces/file.txt",
	}

	for _, url := range tests {
		parsed := RawURLParse(url)
		fmt.Printf("\nTesting URL: %s\n", url)
		fmt.Printf("Scheme: %s\n", parsed.Scheme)
		fmt.Printf("Host: %s\n", parsed.Host)
		fmt.Printf("Path: %s\n", parsed.Path)
		if parsed.User != nil {
			fmt.Printf("Username: %s\n", parsed.User.username)
		}
		if parsed.Opaque != "" {
			fmt.Printf("Opaque: %s\n", parsed.Opaque)
		}
		fmt.Printf("Query: %s\n", parsed.Query)
		fmt.Printf("Fragment: %s\n", parsed.Fragment)
		fmt.Println("---")
	}
}

func TestURLHelperMethods(t *testing.T) {
	url2 := RawURLParse("https://www.example.com:8080/path?key=value&foo=bar")
	port := url2.Port()
	hostname := url2.Hostname()
	params := url2.QueryValues()

	fmt.Printf("Port: %s\n", port)
	fmt.Printf("Hostname: %s\n", hostname)
	fmt.Printf("Params: %v\n", params)
}
