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
			parsed, err := RawURLParse(url)
			if err != nil {
				t.Errorf("Error parsing URL: %s", err)
				return
			}

			// Compare with original instead of reconstructing
			if parsed.Original != url {
				fmt.Printf("%sURL mismatch:\nOriginal: %s\nStored:   %s%s\n",
					colorRed, url, parsed.Original, colorReset)
			}

			fmt.Printf("\nTesting URL: %s\n", url)
			fmt.Printf("Scheme: %s\n", parsed.Scheme)
			fmt.Printf("Host: %s\n", parsed.Host)
			fmt.Printf("Path: %s\n", parsed.Path)
			if parsed.User != nil {
				fmt.Printf("Username: %s\n", parsed.User.username)
				if parsed.User.passwordSet {
					fmt.Printf("Password: %s\n", parsed.User.password)
				}
			}
			if parsed.Opaque != "" {
				fmt.Printf("Opaque: %s\n", parsed.Opaque)
			}
			fmt.Printf("Query: %s\n", parsed.Query)
			fmt.Printf("Fragment: %s\n", parsed.Fragment)
			fmt.Println("---")
		})
	}
}

func TestURLHelperMethods(t *testing.T) {
	url := "https://www.example.com:8080/path?key=value&foo=bar"
	parsed, err := RawURLParse(url)
	if err != nil {
		t.Fatalf("Failed to parse URL: %v", err)
	}

	port := GetRawPort(parsed)
	hostname := GetRawHostname(parsed)
	params := parsed.GetRawQueryValues()

	fmt.Printf("Port: %s\n", port)
	fmt.Printf("Hostname: %s\n", hostname)
	fmt.Printf("Params: %v\n", params)
}

func TestMidPathPayloads(t *testing.T) {
	baseURL := "https://test-go-bypass-403.new.com"
	payloads, err := readPayloads("data/mid_paths_payloads.txt")
	if err != nil {
		t.Fatalf("Failed to read payloads: %v", err)
	}

	fmt.Printf("\nTesting Mid-Path Payloads with base URL: %s\n", baseURL)
	fmt.Println("----------------------------------------")

	for _, payload := range payloads {
		// Create test URL by combining base URL and payload
		testURL := baseURL + payload

		parsed, err := RawURLParse(testURL)
		if err != nil {
			t.Errorf("Error parsing URL with payload %q: %s", payload, err)
			continue
		}

		// Compare with original
		if parsed.Original != testURL {
			fmt.Printf("%sFAILED - Payload: %s\nOriginal: %s\nStored: %s\nPath: %s\nRawPathUnsafe: %s%s\n\n",
				colorRed, payload, testURL, parsed.Original, parsed.Path, GetRawPathUnsafe(parsed), colorReset)
		} else {
			fmt.Printf("%sPASSED - Payload: %s\nOriginal: %s\nStored: %s\nPath: %s\nRawPathUnsafe: %s%s\n\n",
				colorGreen, payload, testURL, parsed.Original, parsed.Path, GetRawPathUnsafe(parsed), colorReset)
		}
	}
}
