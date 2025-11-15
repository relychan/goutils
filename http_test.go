package goutils

import (
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestParseRelativeOrHttpURL(t *testing.T) {
	testCases := []struct {
		URL string
	}{
		{
			URL: "/healthz",
		},
		{
			URL: "https://localhost:8080/hello?foo=bar#about",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.URL, func(t *testing.T) {
			result, err := ParseRelativeOrHTTPURL(tc.URL)
			if err != nil {
				t.Fatalf("expected nil error, got: %s", err)
			}

			if result.String() != tc.URL {
				t.Fatalf("expected equal, got: %s", result)
			}
		})
	}
}

func TestParseHttpURL(t *testing.T) {
	testCases := []struct {
		URL   string
		Error string
	}{
		{
			URL: "http://127.0.0.1/healthz",
		},
		{
			URL: "https://localhost:8080/hello?foo=bar#about",
		},
		{
			URL:   "postgresql://localhost:8080/hello?foo=bar#about",
			Error: "invalid http(s) scheme, got: postgresql",
		},
		{
			URL:   "!@#$$%",
			Error: "invalid URL escape",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.URL, func(t *testing.T) {
			result, err := ParseHTTPURL(tc.URL)
			if tc.Error == "" {
				if err != nil {
					t.Fatalf("expected nil error, got: %s", err)
				}

				if result.String() != tc.URL {
					t.Fatalf("expected equal, got: %s", result)
				}
			} else if err == nil || !strings.Contains(err.Error(), tc.Error) {
				t.Fatalf("expected error contains: %s, got: %s", tc.Error, err)
			}

		})
	}
}

func TestExtractHeaders(t *testing.T) {
	testCases := []struct {
		Input    http.Header
		Expected map[string]string
	}{
		{
			Input: http.Header{
				"Content-Type": []string{"application/json"},
				"FOO":          []string{"BAR"},
			},
			Expected: map[string]string{
				"content-type": "application/json",
				"foo":          "BAR",
			},
		},
	}

	for _, tc := range testCases {
		result := ExtractHeaders(tc.Input)

		if !reflect.DeepEqual(tc.Expected, result) {
			t.Fatalf("not equal, expected: %v, got: %v", tc.Expected, result)
		}
	}
}
