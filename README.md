# RawURLParse

A Go package that parses URLs in their raw form, with:
- No normalization
- No encoding/decoding
- No validation
- Preserves exact characters as provided

Unlike the standard library's url.Parse which uses Opaque for non-hierarchical URLs,
this package preserves the exact path encoding for all URLs. When using with http.Client,
the raw path should be assigned to URL.Opaque to prevent normalization.

# Important Notice

## Go's http.Client
When using parsed URLs with Go's `http.Client`, you'll need to use URL.Opaque to preserve
the exact path encoding, otherwise http.Client will perform encodings, normalization, etc. 

Example Code to preserve and send raw URLs

```go
parsedURL := rawurlparser.RawURLParse(rawURL)
req := &http.Request{
    Method: "GET",
    URL: &url.URL{
        Scheme: parsedURL.Scheme,
        Host:   parsedURL.Host,
        Opaque: parsedURL.Path,  // Use Opaque to prevent path normalization
    },
}
```

In other words, this achieves the same thing as sending a request with curl using `--path-as-is`

## Other Problematic Methods

I found that Go's `to.String()` also applies encodings, do not use on URLs or when debugging URLs.

Use something like this:
```go
log.Printf("Debug - Request Components - Scheme: %s, Host: %s, Path: %s, RawPath: %s, Opaque: %s",
    req.URL.Scheme,
    req.URL.Host,
    req.URL.Path,
    req.URL.RawPath,
    req.URL.Opaque,
)
```

## Installation

```bash
go get github.com/slicingmelon/go-rawurlparser
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/slicingmelon/go-rawurlparser"
)

func main() {
    url := "https://example.com/path1/..%2f/test?q=1#fragment"
    parsed := rawurlparser.RawURLParse(url)
    
    fmt.Printf("Scheme:   %q\n", parsed.Scheme)
    fmt.Printf("Host:     %q\n", parsed.Host)
    fmt.Printf("Path:     %q\n", parsed.Path)
    fmt.Printf("Query:    %q\n", parsed.Query)
    fmt.Printf("Fragment: %q\n", parsed.Fragment)
}
```

## Features

- Preserves all characters exactly as provided
- No normalization of paths
- No encoding/decoding
- Raw preservation of special characters
- Simple and fast parsing

## Tests

```bash
go run .\examples\main.go

Testing URL: https://example.com/x/。。;//
----------------------------------------
url.Parse:
Full URL: https://example.com/x/%E3%80%82%E3%80%82;//
Path: /x/。。;//
RawPath: /x/。。;//

rawurlparser.RawURLParse:
Full URL: https://example.com/x/。。;//
Path: /x/。。;//

## Closer URLs Comparison ##
Standard UrlParser: https://example.com/x/%E3%80%82%E3%80%82;//
RawUrlPaser:      https://example.com/x/。。;//

========================================

Testing URL: https://example.com/x/。。;//
----------------------------------------
url.Parse:
Full URL: https://example.com/x/%E3%80%82%E3%80%82;//
Path: /x/。。;//
RawPath: /x/。。;//

rawurlparser.RawURLParse:
Full URL: https://example.com/x/。。;//
Path: /x/。。;//

## Closer URLs Comparison ##
Standard UrlParser: https://example.com/x/%E3%80%82%E3%80%82;//
RawUrlPaser:      https://example.com/x/。。;//

========================================

Testing URL: https://example.com\..\.\
----------------------------------------
url.Parse error: parse "https://example.com\\..\\.\\": invalid character "\\" in host name

rawurlparser.RawURLParse:
Full URL: https://example.com\..\.\
Path:
========================================

Testing URL: https://example.com#
----------------------------------------
url.Parse:
Full URL: https://example.com
Path:
RawPath:

rawurlparser.RawURLParse:
Full URL: https://example.com#
Path:

## Closer URLs Comparison ##
Standard UrlParser: https://example.com
RawUrlPaser:      https://example.com#
```

## Author

Petru Surugiu<br>
https://twitter.com/pedro_infosec

## License

MIT License