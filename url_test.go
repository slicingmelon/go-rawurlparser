package rawurlparser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"
)

const (
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	colorReset = "\033[0m"
)

// readTestURLs reads test URLs from the specified file
func readTestURLs(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error opening test URLs file: %w", err)
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			urls = append(urls, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading test URLs file: %w", err)
	}

	return urls, nil
}

// readPayloads reads payloads from the specified file
func readPayloads(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error opening payloads file: %w", err)
	}
	defer file.Close()

	var payloads []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			payloads = append(payloads, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading payloads file: %w", err)
	}

	return payloads, nil
}

func TestRawURLParse(t *testing.T) {
	urls, err := readTestURLs("data/test_urls.txt")
	if err != nil {
		t.Fatalf("Failed to read test URLs: %v", err)
	}

	for _, url := range urls {
		t.Run(url, func(t *testing.T) {
			parsedURL, err := RawURLParse(url)
			if err != nil {
				t.Errorf("Error parsing URL: %s", err)
				return
			}

			// Compare with original instead of reconstructing
			if parsedURL.Original != url {
				fmt.Printf("%sURL mismatch:\nOriginal: %s\nStored:   %s%s\n",
					colorRed, url, parsedURL.Original, colorReset)
			}

			fmt.Printf("\nTesting URL: %s\n", url)
			fmt.Printf("Parsed Scheme: %s\n", parsedURL.Scheme)
			fmt.Printf("Parsed Host: %s\n", parsedURL.Host)
			fmt.Printf("Parsed Path: %s\n", parsedURL.Path)
			if parsedURL.User != nil {
				fmt.Printf("Parsed Username: %s\n", parsedURL.User.username)
				if parsedURL.User.passwordSet {
					fmt.Printf("Parsed Password: %s\n", parsedURL.User.password)
				}
			}
			if parsedURL.Opaque != "" {
				fmt.Printf("Parsed Opaque: %s\n", parsedURL.Opaque)
			}
			fmt.Printf("Parsed Query: %s\n", parsedURL.Query)
			fmt.Printf("Parsed Fragment: %s\n", parsedURL.Fragment)
			fmt.Println("---")
		})
	}
}

func TestURLHelperMethods(t *testing.T) {
	url := "https://www.example.com:8080/path?key=value&foo=bar"
	parsedURL, err := RawURLParse(url)
	if err != nil {
		t.Fatalf("Failed to parse URL: %v", err)
	}

	port := parsedURL.GetPort()
	hostname := parsedURL.GetHostname()
	params := parsedURL.GetQueryValues()

	fmt.Printf("Parsed Port: %s\n", port)
	fmt.Printf("Parsed Hostname: %s\n", hostname)
	fmt.Printf("Parsed Params: %v\n", params)
}

func TestMidPathPayloads(t *testing.T) {
	baseURL := "https://test-go-bypass-403-new.com"
	payloads, err := readPayloads("data/mid_paths_payloads.txt")
	if err != nil {
		t.Fatalf("Failed to read payloads: %v", err)
	}

	fmt.Printf("\nTesting Mid-Path Payloads with base URL: %s\n", baseURL)
	fmt.Println("----------------------------------------")

	for _, payload := range payloads {
		// Ensure payload starts with / if it doesn't already
		if !strings.HasPrefix(payload, "/") {
			payload = "/" + payload
		}

		// Create test URL by combining base URL and payload
		testURL := baseURL + payload

		parsedURL, err := RawURLParse(testURL)
		if err != nil {
			t.Errorf("Error parsing URL with payload %q: %s", payload, err)
			continue
		}

		// Verify the payload is in the path, not the host
		if strings.Contains(parsedURL.Host, payload) {
			t.Errorf("%sFAILED - Payload incorrectly parsed as part of host:\nPayload: %s\nHost: %s%s\n",
				colorRed, payload, parsedURL.Host, colorReset)
			continue
		}

		// Verify the path contains the payload
		if !strings.Contains(parsedURL.Path, strings.TrimPrefix(payload, "/")) {
			t.Errorf("%sFAILED - Payload not found in path:\nPayload: %s\nPath: %s%s\n",
				colorRed, payload, parsedURL.Path, colorReset)
			continue
		}

		parsedUrlForComparison := fmt.Sprintf("%s://%s%s", parsedURL.Scheme, parsedURL.Host, parsedURL.Path)
		if parsedURL.Query != "" {
			parsedUrlForComparison += "?" + parsedURL.Query
		}
		if parsedURL.Fragment != "" {
			parsedUrlForComparison += "#" + parsedURL.Fragment
		}

		fmt.Printf("%sPASSED - Payload correctly parsed in path:\nPayload: %s\ntestURL: %s\nParsed Path: %s%s\n\n",
			colorGreen, payload, testURL, parsedURL.Path, colorReset)
	}
}

func TestIPAddressURLs(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		wantHost string
		wantPath string
		wantPort string
	}{
		{
			name:     "IPv4 with scheme and path",
			input:    "http://192.168.1.1/test",
			wantHost: "192.168.1.1",
			wantPath: "/test",
			wantPort: "",
		},
		{
			name:     "IPv4 with port and path",
			input:    "192.168.10.10:443/x/y",
			wantHost: "192.168.10.10:443",
			wantPath: "/x/y",
			wantPort: "443",
		},
		{
			name:     "IPv4 with scheme, port and path",
			input:    "https://10.0.0.1:8080/admin",
			wantHost: "10.0.0.1:8080",
			wantPath: "/admin",
			wantPort: "8080",
		},
		{
			name:     "IPv6 with scheme",
			input:    "http://[2001:db8::1]/test",
			wantHost: "[2001:db8::1]",
			wantPath: "/test",
			wantPort: "",
		},
		{
			name:     "IPv6 with port",
			input:    "[2001:db8::1]:8443/secure",
			wantHost: "[2001:db8::1]:8443",
			wantPath: "/secure",
			wantPort: "8443",
		},
		{
			name:     "IPv4 localhost with port",
			input:    "127.0.0.1:3000/api/v1",
			wantHost: "127.0.0.1:3000",
			wantPath: "/api/v1",
			wantPort: "3000",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parsedURL, err := RawURLParse(tc.input)
			if err != nil {
				t.Errorf("RawURLParse(%q) returned error: %v", tc.input, err)
				return
			}

			if parsedURL.Host != tc.wantHost {
				t.Errorf("Host = %q, want %q", parsedURL.Host, tc.wantHost)
			}

			if parsedURL.Path != tc.wantPath {
				t.Errorf("Path = %q, want %q", parsedURL.Path, tc.wantPath)
			}

			gotPort := parsedURL.GetPort()
			if gotPort != tc.wantPort {
				t.Errorf("Port = %q, want %q", gotPort, tc.wantPort)
			}

			// Print the results for visual inspection
			fmt.Printf("\nTesting: %s\n", tc.input)
			fmt.Printf("----------------------------------------\n")
			fmt.Printf("Host: %s\n", parsedURL.Host)
			fmt.Printf("Path: %s\n", parsedURL.Path)
			fmt.Printf("Port: %s\n", gotPort)
			fmt.Printf("Full URL: %s\n", parsedURL.GetRawFullURL())
			fmt.Printf("----------------------------------------\n")
		})
	}
}
