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

	port := GetRawPort(parsedURL)
	hostname := GetRawHostname(parsedURL)
	params := parsedURL.GetRawQueryValues()

	fmt.Printf("Parsed Port: %s\n", port)
	fmt.Printf("Parsed Hostname: %s\n", hostname)
	fmt.Printf("Parsed Params: %v\n", params)
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

		parsedURL, err := RawURLParse(testURL)
		if err != nil {
			t.Errorf("Error parsing URL with payload %q: %s", payload, err)
			continue
		}

		// Compare with original
		if parsedURL.Original != testURL {
			fmt.Printf("%sFAILED - Payload: %s\nOriginal: %s\nStored: %s\nPath: %s\nRawPathUnsafe: %s%s\n\n",
				colorRed, payload, testURL, parsedURL.Original, parsedURL.Path, GetRawPathUnsafe(parsedURL), colorReset)
		} else {
			fmt.Printf("%sPASSED - Payload: %s\nOriginal: %s\nStored: %s\nPath: %s\nRawPathUnsafe: %s%s\n\n",
				colorGreen, payload, testURL, parsedURL.Original, parsedURL.Path, GetRawPathUnsafe(parsedURL), colorReset)
		}
	}
}
