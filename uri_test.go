// Copyright 2026 RelyChan Pte. Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package goutils

import (
	"context"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestParseRelativeOrHttpURL(t *testing.T) {
	testCases := []struct {
		URL string
	}{
		{
			URL: "",
		},
		{
			URL: "/",
		},
		{
			URL: "/healthz",
		},
		{
			URL: "../healthz",
		},
		{
			URL: "https://localhost:8080/hello?foo=bar#about",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.URL, func(t *testing.T) {
			result, err := ParsePathOrHTTPURL(tc.URL)
			if err != nil {
				t.Fatalf("expected nil error, got: %s", err)
			}

			if result.String() != tc.URL {
				t.Fatalf("expected equal, got: %s", result)
			}
		})
	}
}

func TestParseRelativeOrHttpURL_Errors(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		err   string
	}{
		{
			name:  "invalid scheme",
			input: "ftp://example.com",
			err:   `Invalid HTTP scheme. Expected http(s), got "ftp"`,
		},
		{
			name:  "scheme prefix only",
			input: "://example.com",
			err:   `Invalid URL. Scheme is empty`,
		},
		{
			name:  "postgresql scheme",
			input: "postgresql://localhost/db",
			err:   `Invalid HTTP scheme. Expected http(s), got "postgresql"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParsePathOrHTTPURL(tc.input)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}

			if !strings.Contains(err.Error(), tc.err) {
				t.Fatalf("expected error %v, got: %v", tc.err, err)
			}
		})
	}
}

func TestParseRelativeOrHTTPURL_RelativePaths(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		wantPath string
	}{
		{
			name:     "absolute path with query and fragment",
			input:    "/api/v1/resource?key=value#section",
			wantPath: "/api/v1/resource",
		},
		{
			name:     "http URL",
			input:    "http://example.com/path",
			wantPath: "/path",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ParsePathOrHTTPURL(tc.input)
			if err != nil {
				t.Fatalf("expected nil error, got: %s", err)
			}

			if result.Path != tc.wantPath {
				t.Fatalf("expected path %q, got %q", tc.wantPath, result.Path)
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
			Error: "Invalid HTTP URL scheme",
		},
		{
			URL:   "!@#$$%",
			Error: "Invalid HTTP URL scheme",
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

func TestParseHTTPURL_AllowedSchemes(t *testing.T) {
	t.Run("filters non-http schemes from AllowedSchemes", func(t *testing.T) {
		// ftp and ws are not http/https, they should be removed; only https passes
		_, err := ParseAndValidateURLWithOptions(context.Background(), "https://localhost/path", &ValidateHTTPURLOptions{
			AllowedSchemes: []string{"ftp", "https", "ws"},
		})
		if err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	t.Run("http blocked when only https allowed", func(t *testing.T) {
		_, err := ParseAndValidateURLWithOptions(context.Background(), "http://localhost/path", &ValidateHTTPURLOptions{
			AllowedSchemes: []string{"ftp", "https", "ws"},
		})
		if err == nil || !strings.Contains(err.Error(), `Invalid URI scheme. Accept one of [ftp, https, ws], got "http"`) {
			t.Fatalf("expected invalid URI scheme validation error, got: %v", err)
		}
	})
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

func TestExtractHeaders_EdgeCases(t *testing.T) {
	t.Run("empty headers", func(t *testing.T) {
		result := ExtractHeaders(http.Header{})
		if len(result) != 0 {
			t.Fatalf("expected empty map, got: %v", result)
		}
	})

	t.Run("header with empty value slice is skipped", func(t *testing.T) {
		result := ExtractHeaders(http.Header{
			"X-Empty": []string{},
			"X-Set":   []string{"value"},
		})
		if _, ok := result["x-empty"]; ok {
			t.Fatal("expected x-empty to be absent")
		}

		if result["x-set"] != "value" {
			t.Fatalf("expected x-set=value, got: %v", result["x-set"])
		}
	})

	t.Run("multiple values: last is kept", func(t *testing.T) {
		result := ExtractHeaders(http.Header{
			"Accept": []string{"text/html", "application/json"},
		})
		if result["accept"] != "application/json" {
			t.Fatalf("expected last value application/json, got: %s", result["accept"])
		}
	})
}

func TestParseAndValidateURL(t *testing.T) {
	t.Run("empty string returns ErrInvalidURI", func(t *testing.T) {
		_, err := ParseAndValidateURLWithOptions(context.Background(), "", &ValidateHTTPURLOptions{})
		if !strings.Contains(err.Error(), "Invalid URL. The input string is empty") {
			t.Fatalf("expected ErrInvalidURI, got: %v", err)
		}
	})

	t.Run("whitespace-only returns ErrInvalidURI", func(t *testing.T) {
		_, err := ParseAndValidateURLWithOptions(context.Background(), "   ", &ValidateHTTPURLOptions{})
		if !strings.Contains(err.Error(), "Invalid URL. The input string is empty") {
			t.Fatalf("expected ErrInvalidURI, got: %v", err)
		}
	})

	t.Run("valid http URL", func(t *testing.T) {
		u, err := ParseAndValidateURLWithOptions(context.Background(), "http://127.0.0.1/path", &ValidateHTTPURLOptions{})
		if err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}

		if u.Host != "127.0.0.1" {
			t.Fatalf("unexpected host: %s", u.Host)
		}
	})
}

func TestParseURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  string
		wantHost string
		wantPath string
	}{
		{
			name:     "valid http URL",
			input:    "http://example.com/path",
			wantHost: "example.com",
			wantPath: "/path",
		},
		{
			name:     "valid https URL with port",
			input:    "https://example.com:8443/api?q=1#frag",
			wantHost: "example.com:8443",
			wantPath: "/api",
		},
		{
			name:     "valid postgresql URL",
			input:    "postgresql://db.example.com/mydb",
			wantHost: "db.example.com",
			wantPath: "/mydb",
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: "Invalid URL. The input string is empty",
		},
		{
			name:    "whitespace only",
			input:   "   ",
			wantErr: "Invalid URL. The input string is empty",
		},
		{
			name:    "no scheme",
			input:   "example.com/path",
			wantErr: "Invalid URL syntax",
		},
		{
			name:    "colon without slashes",
			input:   "example:path",
			wantErr: "Invalid URL syntax",
		},
		{
			name:    "scheme with no host",
			input:   "http:///path",
			wantErr: "Invalid URL. Hostname is empty",
		},
		{
			name:     "IPv4 host",
			input:    "http://192.168.1.1/resource",
			wantHost: "192.168.1.1",
		},
		{
			name:     "IPv6 host",
			input:    "http://[::1]/path",
			wantHost: "[::1]",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, err := ParseURL(tc.input)
			if tc.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tc.wantErr)
				}
				if !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("expected error %q, got: %v", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.wantHost != "" && u.Host != tc.wantHost {
				t.Fatalf("expected host %q, got %q", tc.wantHost, u.Host)
			}
			if tc.wantPath != "" && u.Path != tc.wantPath {
				t.Fatalf("expected path %q, got %q", tc.wantPath, u.Path)
			}
		})
	}
}

func TestParseURI(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  string
		wantHost string
		wantPath string
	}{
		{
			name:     "standard http URL",
			input:    "http://example.com/path",
			wantHost: "example.com",
			wantPath: "/path",
		},
		{
			name:  "URI without authority (opaque URI)",
			input: "urn:isbn:0451450523",
		},
		{
			name:  "mailto URI",
			input: "mailto:user@example.com",
		},
		{
			// ParseURI delegates to url.Parse which accepts scheme-less strings;
			// strict scheme enforcement only applies via ValidateURI / parseAndValidateURI.
			name:     "scheme-less string is accepted by ParseURI",
			input:    "example.com/path",
			wantPath: "example.com/path",
		},
		{
			name:     "IPv6 host",
			input:    "http://[::1]:8080/path",
			wantHost: "[::1]:8080",
		},
		{
			// url.Parse reports "invalid port" rather than an IPv6-specific message
			// when an unbracketed IPv6 address is used as the host.
			name:    "invalid IPv6 without brackets",
			input:   "http://::1/path",
			wantErr: "invalid port",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, err := ParseURI(tc.input)
			if tc.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tc.wantErr)
				}
				if !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("expected error %q, got: %v", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.wantHost != "" && u.Host != tc.wantHost {
				t.Fatalf("expected host %q, got %q", tc.wantHost, u.Host)
			}
			if tc.wantPath != "" && u.Path != tc.wantPath {
				t.Fatalf("expected path %q, got %q", tc.wantPath, u.Path)
			}
		})
	}
}

func TestValidateURI_ErrorCode(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "valid http URL", input: "http://example.com/path", wantErr: false},
		{name: "valid urn", input: "urn:isbn:0451450523", wantErr: false},
		{name: "empty string", input: "", wantErr: true},
		{name: "missing scheme", input: "example.com/path", wantErr: true},
		{name: "only colon prefix", input: ":bad", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateURI(tc.input)
			if tc.wantErr && err == nil {
				t.Fatal("expected validation error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("expected nil error, got: %v", err)
			}
			if err != nil && err.Code != ErrCodeInvalidURI {
				t.Fatalf("expected code %q, got %q", ErrCodeInvalidURI, err.Code)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		wantMsg string
	}{
		{name: "valid http URL", input: "http://example.com", wantErr: false},
		{name: "valid https with port", input: "https://api.example.com:443/v1", wantErr: false},
		{name: "empty string", input: "", wantErr: true, wantMsg: "empty"},
		{name: "no scheme", input: "example.com", wantErr: true, wantMsg: "Invalid URL syntax"},
		{name: "scheme missing slashes", input: "http:example.com", wantErr: true, wantMsg: "Invalid URL syntax"},
		{name: "no host after scheme", input: "http:///path", wantErr: true, wantMsg: "Hostname is empty"},
		{name: "no double slashes", input: "http://example.com//path", wantErr: true, wantMsg: "Invalid double slashes in the URL path syntax"},
		{name: "wildcard", input: "http://example.com/*", wantErr: true, wantMsg: "Invalid URL path syntax"},
		{name: "traversal", input: "http://example.com/foo/../bar", wantErr: true, wantMsg: "Invalid URL path syntax"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateURL(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected validation error, got nil")
				}
				if tc.wantMsg != "" && !strings.Contains(err.Error(), tc.wantMsg) {
					t.Fatalf("expected error containing %q, got: %v", tc.wantMsg, err)
				}
			} else if err != nil {
				t.Fatalf("expected nil error, got: %v", err)
			}
		})
	}
}

func TestParsePathOrURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  string
		wantPath string
		wantHost string
	}{
		{
			name:     "empty string returns empty URL",
			input:    "",
			wantPath: "",
			wantHost: "",
		},
		{
			name:     "absolute path",
			input:    "/api/v1",
			wantPath: "/api/v1",
		},
		{
			name:     "path with query and fragment",
			input:    "/search?q=test#results",
			wantPath: "/search",
		},
		{
			name:     "relative path",
			input:    "../up/one",
			wantPath: "../up/one",
		},
		{
			name:     "valid http URL",
			input:    "http://example.com/path",
			wantHost: "example.com",
			wantPath: "/path",
		},
		{
			name:    "colon at position 0",
			input:   "://host",
			wantErr: "Invalid URL. Scheme is empty",
		},
		{
			name:    "colon without double slash",
			input:   "http:example.com",
			wantErr: "Invalid URL syntax",
		},
		{
			name:    "single slash after colon",
			input:   "http:/example.com",
			wantErr: "Invalid URL syntax",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, err := ParsePathOrURL(tc.input)
			if tc.wantErr != "" {
				if !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("expected error %v, got: %v", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.wantHost != "" && u.Host != tc.wantHost {
				t.Fatalf("expected host %q, got %q", tc.wantHost, u.Host)
			}
			if tc.wantPath != "" && u.Path != tc.wantPath {
				t.Fatalf("expected path %q, got %q", tc.wantPath, u.Path)
			}
		})
	}
}

func TestParsePathOrURL_CTLBytes(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"null byte", "/path\x00end"},
		{"control byte 0x01", "/path\x01end"},
		{"DEL byte 0x7f", "/path\x7fend"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParsePathOrURL(tc.input)
			if !strings.Contains(err.Error(), "Path contains invalid characters") {
				t.Fatalf("expected ErrInvalidURI for CTL byte input, got: %v", err)
			}
		})
	}
}

func TestParsePathOrURL_QueryFragment(t *testing.T) {
	u, err := ParsePathOrURL("/path?key=val&other=2#section")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.Path != "/path" {
		t.Fatalf("expected path /path, got %q", u.Path)
	}
	if u.RawQuery != "key=val&other=2" {
		t.Fatalf("expected query key=val&other=2, got %q", u.RawQuery)
	}
	if u.Fragment != "section" {
		t.Fatalf("expected fragment section, got %q", u.Fragment)
	}
}

func TestParseURL_IPv6(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
	}{
		{
			name:  "valid IPv6 loopback",
			input: "http://[::1]/path",
		},
		{
			name:  "valid IPv6 full address",
			input: "http://[2001:db8::1]/resource",
		},
		{
			// url.Parse reports "invalid port" for unbracketed IPv6; parseURIAndHostname
			// checks for brackets and delegates the IPv6 validation message accordingly.
			name:    "IPv6 without brackets",
			input:   "http://::1/path",
			wantErr: "invalid port",
		},
		{
			name:    "invalid IPv6 in brackets",
			input:   "http://[not:valid:ipv6:addr:foo:bar:baz:qux:extra]/path",
			wantErr: "invalid host",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseURL(tc.input)
			if tc.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tc.wantErr)
				}
				if !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("expected error %q, got: %v", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestParseAndValidateURLWithOptions_InvalidURL(t *testing.T) {
	ctx := context.Background()

	_, err := ParseAndValidateURLWithOptions(ctx, "not-a-url", &ValidateHTTPURLOptions{})
	if err == nil {
		t.Fatal("expected error for invalid URL, got nil")
	}
	if !strings.Contains(err.Error(), "Invalid URL syntax") {
		t.Fatalf("expected 'Invalid URL syntax', got: %v", err)
	}
}

func mustParseURL(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		panic(err)
	}

	return u
}

func TestAppendURL(t *testing.T) {
	tests := []struct {
		name    string
		base    string
		uriPath string
		wantURL string
		wantErr string
	}{

		// No-op cases
		{
			name:    "empty path returns base unchanged",
			base:    "https://example.com/api",
			uriPath: "",
			wantURL: "https://example.com/api",
		},
		{
			name:    "slash-only path returns base unchanged",
			base:    "https://example.com/api",
			uriPath: "/",
			wantURL: "https://example.com/api",
		},

		// Path-only cases
		{
			name:    "absolute path appended to empty base path",
			base:    "https://example.com",
			uriPath: "/users",
			wantURL: "https://example.com/users",
		},
		{
			name:    "absolute path appended to root base path",
			base:    "https://example.com/",
			uriPath: "/users",
			wantURL: "https://example.com/users",
		},
		{
			name:    "absolute path appended to existing base path",
			base:    "https://example.com/api",
			uriPath: "/v1/users",
			wantURL: "https://example.com/api/v1/users",
		},
		{
			name:    "relative path appended with separator",
			base:    "https://example.com/api",
			uriPath: "v1/users",
			wantURL: "https://example.com/api/v1/users",
		},
		{
			name:    "relative path appended to empty base path",
			base:    "https://example.com",
			uriPath: "users",
			wantURL: "https://example.com/users",
		},

		// Query-only cases
		{
			name:    "query appended to base with no query",
			base:    "https://example.com/api",
			uriPath: "?limit=10",
			wantURL: "https://example.com/api?limit=10",
		},
		{
			name:    "query merged with existing base query",
			base:    "https://example.com/api?page=1",
			uriPath: "?limit=10",
			wantURL: "https://example.com/api?page=1&limit=10",
		},

		// Fragment cases
		{
			name:    "fragment set from uri path",
			base:    "https://example.com/api",
			uriPath: "#section",
			wantURL: "https://example.com/api#section",
		},

		// Combined cases
		{
			name:    "path and query combined",
			base:    "https://example.com/api",
			uriPath: "/v1/users?limit=10",
			wantURL: "https://example.com/api/v1/users?limit=10",
		},
		{
			name:    "path, query, and fragment combined",
			base:    "https://example.com/api",
			uriPath: "/v1/users?limit=10#section",
			wantURL: "https://example.com/api/v1/users?limit=10#section",
		},
		{
			name:    "path and fragment without query",
			base:    "https://example.com/api",
			uriPath: "/v1/users#section",
			wantURL: "https://example.com/api/v1/users#section",
		},
		{
			name:    "query merged and fragment set",
			base:    "https://example.com/api?page=1",
			uriPath: "?limit=10#section",
			wantURL: "https://example.com/api?page=1&limit=10#section",
		},

		// Base URL with trailing slash
		{
			name:    "absolute path appended to base with trailing slash",
			base:    "https://example.com/api/",
			uriPath: "/v2/items",
			wantURL: "https://example.com/api/v2/items",
		},

		// Base URL with existing query
		{
			name:    "path appended and base query preserved",
			base:    "https://example.com/api?version=2",
			uriPath: "/users",
			wantURL: "https://example.com/api/users?version=2",
		},
		{
			name:    "invalid path segment returns error",
			base:    "https://example.com/api",
			uriPath: "/..",
			wantErr: "Invalid URL path syntax",
		},
		{
			name:    "invalid double slashes in appended path returns error",
			base:    "https://example.com/api",
			uriPath: "/v1//users",
			wantErr: "Invalid double slashes",
		},
		{
			name:    "invalid CTL byte in query returns error",
			base:    "https://example.com/api",
			uriPath: "?a=\x00",
			wantErr: "Invalid URL query syntax",
		},
		{
			name:    "invalid CTL byte in fragment returns error",
			base:    "https://example.com/api",
			uriPath: "#\x7f",
			wantErr: "Invalid URL fragment syntax",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			base := mustParseURL(tc.base)
			err := AppendURL(base, tc.uriPath)

			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("expected error containing %q, got: %v", tc.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected nil error, got: %v", err)
			}
			if got := base.String(); got != tc.wantURL {
				t.Fatalf("unexpected URL: got %q, want %q", got, tc.wantURL)
			}
		})
	}
}
