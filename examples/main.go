// File: examples/main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/slicingmelon/go-rawurlparser"
)

func main() {
	// Example URLs to test
	rawURLs := []string{
		"https://example.com/path1/..%2f/path2",
		"https://example.com/path1;/%2e%2e/path2",
		"https://example.com/path1/path2;/。。/path3",
		"https://example.com/path1%2f..%2f/path2#fragment?query=1",
	}

	for _, rawURL := range rawURLs {
		fmt.Printf("\nTesting URL: %s\n", rawURL)
		fmt.Printf("----------------------------------------\n")

		// Parse with our raw parser
		parsedURL := rawurlparser.RawURLParse(rawURL)
		log.Printf("RawURLParse result - Scheme: %s, Host: %s, Path: %s, Query: %s, Fragment: %s",
			parsedURL.Scheme,
			parsedURL.Host,
			parsedURL.Path,
			parsedURL.Query,
			parsedURL.Fragment,
		)

		// Compare with standard parser
		stdURL, _ := url.Parse(rawURL)
		log.Printf("\nStandard url.Parse result:")
		log.Printf("Path: %s", stdURL.Path)
		log.Printf("RawPath: %s", stdURL.RawPath)

		// Show in http.Request
		req := &http.Request{
			Method: "GET",
			URL: &url.URL{
				Scheme:  parsedURL.Scheme,
				Host:    parsedURL.Host,
				Path:    parsedURL.Path,
				RawPath: parsedURL.Path,
				Opaque:  parsedURL.Path,
			},
		}

		log.Printf("\nIn http.Request:")
		log.Printf("URL: %s", req.URL)
		log.Printf("Raw Path: %s", req.URL.RawPath)
		log.Printf("Opaque: %s", req.URL.Opaque)
	}
}
