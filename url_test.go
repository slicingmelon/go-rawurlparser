// File: url_test.go
package rawurlparse

import (
	"testing"
)

func TestRawURLParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		scheme   string
		host     string
		path     string
		query    string
		fragment string
	}{
		{
			name:     "Basic URL",
			input:    "https://example.com/path",
			scheme:   "https",
			host:     "example.com",
			path:     "/path",
			query:    "",
			fragment: "",
		},
		{
			name:     "URL with everything",
			input:    "https://example.com/path1/path2%2f;/?q=1#f",
			scheme:   "https",
			host:     "example.com",
			path:     "/path1/path2%2f;/",
			query:    "q=1",
			fragment: "f",
		},
		{
			name:     "Path traversal raw",
			input:    "http://example.com/path1/..%2f/",
			scheme:   "http",
			host:     "example.com",
			path:     "/path1/..%2f/",
			query:    "",
			fragment: "",
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RawURLParse(tt.input)
			if got.Scheme != tt.scheme {
				t.Errorf("Scheme = %v, want %v", got.Scheme, tt.scheme)
			}
			if got.Host != tt.host {
				t.Errorf("Host = %v, want %v", got.Host, tt.host)
			}
			if got.Path != tt.path {
				t.Errorf("Path = %v, want %v", got.Path, tt.path)
			}
			if got.Query != tt.query {
				t.Errorf("Query = %v, want %v", got.Query, tt.query)
			}
			if got.Fragment != tt.fragment {
				t.Errorf("Fragment = %v, want %v", got.Fragment, tt.fragment)
			}
		})
	}
}
