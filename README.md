# RawURLParse

A Go package that parses URLs in their raw form.

- Preserves all characters exactly as provided
- No normalization of paths
- No encoding/decoding
- No validation

Unlike the standard library's url.Parse which uses Opaque for non-hierarchical URLs,
this package preserves the exact path encoding for all URLs. When using with http.Client,
the raw path should be assigned to URL.Opaque to prevent normalization.

## Features

- Preserves all characters exactly as provided
- No normalization of paths
- No encoding/decoding
- Raw preservation of special characters
- Simple and fast parsing
- Helper methods for port, hostname, and query parsing
- Optional error handling with RawURLParseWithError



# Understanding URL/URI 

URIs, URLs, and URNs

- A Uniform Resource Identifier (URI) is a string of characters that uniquely identify a name or a resource on the internet. A URI identifies a resource by name, location, or both. URIs have two specializations known as Uniform Resource Locator (URL), and Uniform Resource Name (URN).

-  A Uniform Resource Locator (URL) is a type of URI that specifies not only a resource, but how to reach it on the internet—like http://, ftp://, or mailto://.

- A Uniform Resource Name (URN) is a type of URI that uses the specific naming scheme of urn:—like urn:isbn:0-486-27557-4 or urn:isbn:0-395-36341-1.

- So a URI or URN is like your name, and a URL is a specific subtype of URI that’s like your name combined with your address. 

All URLs are URIs, but not all URIs are URLs.

# URL Structures

![url-structure-and-scheme](./images/url-structure-and-scheme.jpg)

 - A URI is an identifier of a specific resource. Examples: Books, Documents

 - A URL is special type of identifier that also tells you how to access it. Examples: HTTP, FTP, MAILTO

- If the protocol (https, ftp, etc.) is either present or implied for a domain, you should call it a URL—even though it’s also a URI. 

![full-uri-breakdown](./images/full-uri-breakdown.jpg)



## Important Notice

### Go's http.Client
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

### Other Problematic Methods

I found that Go's `to.String()` also applies encodings, do not use it on URLs or when debugging URLs.

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

## Helper Methods

The URL struct provides several helper methods:

- `Port()` - Returns the port number from the host if present
- `Hostname()` - Returns the hostname without the port
- `QueryValues()` - Parses query string into a map[string][]string
- `FullString()` - Reconstructs the full URL from its components

## Error Handling

The package provides two parsing functions:

- `RawURLParse()` - Returns a URL struct or nil if error
- `RawURLParseWithError()` - Returns a URL struct or an error

Example with error handling:

```go
url, err := RawURLParseWithError("example.com/path")
if err != nil {
    log.Fatal("Invalid URL:", err) // Will error: missing scheme
}
```

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