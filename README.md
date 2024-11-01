# RawURLParse

A Go package that parses URLs in their raw form, with:
- No normalization
- No encoding/decoding
- No validation
- Preserves exact characters as provided

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
    parsed := rawurlparse.RawURLParse(url)
    
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

## License

MIT License --