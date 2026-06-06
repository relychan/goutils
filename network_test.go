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
	"errors"
	"net"
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
		err   error
	}{
		{
			name:  "invalid scheme",
			input: "ftp://example.com",
			err:   ErrInvalidURLScheme,
		},
		{
			name:  "scheme prefix only",
			input: "://example.com",
			err:   ErrInvalidURLScheme,
		},
		{
			name:  "postgresql scheme",
			input: "postgresql://localhost/db",
			err:   ErrInvalidURLScheme,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParsePathOrHTTPURL(tc.input)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}

			if !errors.Is(err, tc.err) {
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
			Error: "invalid url scheme. Accept one of [http https], got: postgresql",
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

func TestParseHTTPURL_AllowedSchemes(t *testing.T) {
	t.Run("filters non-http schemes from AllowedSchemes", func(t *testing.T) {
		// ftp and ws are not http/https, they should be removed; only https passes
		_, err := ParseAndValidateURL(context.Background(), "https://localhost/path", ValidateHTTPURLOptions{
			AllowedSchemes: []string{"ftp", "https", "ws"},
		})
		if err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	t.Run("http blocked when only https allowed", func(t *testing.T) {
		_, err := ParseAndValidateURL(context.Background(), "http://localhost/path", ValidateHTTPURLOptions{
			AllowedSchemes: []string{"ftp", "https", "ws"},
		})
		if err == nil || !errors.Is(err, ErrInvalidURLScheme) {
			t.Fatalf("expected ErrInvalidURLScheme, got: %v", err)
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
		_, err := ParseAndValidateURL(context.Background(), "", ValidateHTTPURLOptions{})
		if !errors.Is(err, ErrInvalidURI) {
			t.Fatalf("expected ErrInvalidURI, got: %v", err)
		}
	})

	t.Run("whitespace-only returns ErrInvalidURI", func(t *testing.T) {
		_, err := ParseAndValidateURL(context.Background(), "   ", ValidateHTTPURLOptions{})
		if !errors.Is(err, ErrInvalidURI) {
			t.Fatalf("expected ErrInvalidURI, got: %v", err)
		}
	})

	t.Run("valid http URL", func(t *testing.T) {
		u, err := ParseAndValidateURL(context.Background(), "http://127.0.0.1/path", ValidateHTTPURLOptions{})
		if err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}

		if u.Host != "127.0.0.1" {
			t.Fatalf("unexpected host: %s", u.Host)
		}
	})
}

func TestParseSubnet(t *testing.T) {
	t.Run("valid CIDR", func(t *testing.T) {
		subnet, err := ParseSubnet("192.168.1.0/24")
		if err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}

		if subnet == nil {
			t.Fatal("expected non-nil subnet")
		}
	})

	t.Run("IPv4 address without prefix gets /32", func(t *testing.T) {
		subnet, err := ParseSubnet("10.0.0.1")
		if err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}

		ones, bits := subnet.Mask.Size()
		if ones != 32 || bits != 32 {
			t.Fatalf("expected /32, got /%d/%d", ones, bits)
		}
	})

	t.Run("IPv6 address without prefix gets /128", func(t *testing.T) {
		subnet, err := ParseSubnet("::1")
		if err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}

		ones, bits := subnet.Mask.Size()
		if ones != 128 || bits != 128 {
			t.Fatalf("expected /128, got /%d/%d", ones, bits)
		}
	})

	t.Run("empty string returns ErrInvalidSubnet", func(t *testing.T) {
		_, err := ParseSubnet("")
		if !errors.Is(err, ErrInvalidSubnet) {
			t.Fatalf("expected ErrInvalidSubnet, got: %v", err)
		}
	})

	t.Run("invalid IP string returns ErrInvalidSubnet", func(t *testing.T) {
		_, err := ParseSubnet("not-an-ip")
		if !errors.Is(err, ErrInvalidSubnet) {
			t.Fatalf("expected ErrInvalidSubnet, got: %v", err)
		}
	})

	t.Run("invalid CIDR notation", func(t *testing.T) {
		_, err := ParseSubnet("999.999.999.999/24")
		if err == nil {
			t.Fatal("expected error for invalid CIDR, got nil")
		}
	})

	t.Run("valid IPv6 CIDR", func(t *testing.T) {
		subnet, err := ParseSubnet("2001:db8::/32")
		if err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}

		if subnet == nil {
			t.Fatal("expected non-nil subnet")
		}
	})
}

// parseNetCIDR is a helper that wraps net.ParseCIDR for test use.
func parseNetCIDR(s string) (net.IP, *net.IPNet, error) {
	return net.ParseCIDR(s)
}
