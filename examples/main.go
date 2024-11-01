// File: examples/main.go
package main

import (
	"fmt"
	"net/url"

	"github.com/slicingmelon/go-rawurlparser"
)

const (
	redColor   = "\033[31m"
	resetColor = "\033[0m"
)

func main() {
	// Test URLs -- use string literals
	urls := []string{
		`https://example.com/path1/..%2f/path2`,
		`https://example.com/path1;/%2e%2e/path2`,
		`https://example.com/path1/path2;/。。/path3`,
		`https://example.com/path1%2f..%2f/path2#fragment?query=1`,
		`https://example.com/x;/%2e%2e`,
		`https://example.com/x;/%2e%2e/`,
		`https://example.com/x/;/..;/`,
		`https://example.com/x/;/../`,
		`https://example.com/x/..;/`,
		`https://example.com/x/..;/;/`,
		`https://example.com/x/..;//`,
		`https://example.com/x/../`,
		`https://example.com/x/../;/`,
		`https://example.com/x/..//`,
		`https://example.com/x/。。;//`,
		`https://example.com/x//..;/`,
		`https://example.com/x//../`,
		`https://example.com/x;/%2e%2e`,
		`https://example.com/x;/%2e%2e/`,
		`https://example.com/x/;/..;/`,
		`https://example.com/x/;/../`,
		`https://example.com/x/..;/`,
		`https://example.com/x/..;/;/`,
		`https://example.com/x/..;//`,
		`https://example.com/x/../`,
		`https://example.com/x/../;/`,
		`https://example.com/x/..//`,
		`https://example.com/x/。。;//`,
		`https://example.com/x//..;/`,
		`https://example.com/x//../`,
		`https://example.com\..\.\`,
		`https://example.com&`,
		`https://example.com#`,
		`https://example.com#?`,
	}

	for _, testURL := range urls {
		fmt.Printf("\nTesting URL: %s\n", testURL)
		fmt.Printf("----------------------------------------\n")

		// Standard url.Parse
		stdURL, err := url.Parse(testURL)
		stdFullURL := ""
		if err != nil {
			fmt.Printf("%surl.Parse error: %v%s\n", redColor, err, resetColor)
		} else {
			stdFullURL = stdURL.String()
			fmt.Printf("url.Parse:\n")
			fmt.Printf("Full URL: %s\n", stdFullURL)
			fmt.Printf("Path: %s\n", stdURL.Path)
			fmt.Printf("RawPath: %s\n", stdURL.RawPath)
		}

		// rawurlparser
		rawURL := rawurlparser.RawURLParse(testURL)
		fmt.Printf("\nrawurlparser.RawURLParse:\n")
		fmt.Printf("Full URL: %s\n", rawURL)
		fmt.Printf("Path: %s\n", rawURL.Path)

		// Compare full URLs if standard parsing succeeded
		if err == nil && stdFullURL != rawURL.String() {
			fmt.Printf("\n%s## Closer URLs Comparison ##%s\n", redColor, resetColor)
			fmt.Printf("%sStandard UrlParser: %s%s\n", redColor, stdFullURL, resetColor)
			fmt.Printf("%sRawUrlParser:      %s%s\n", redColor, rawURL, resetColor)
		}

		fmt.Printf("========================================\n")
	}
}
