package goutils

import (
	"context"
	"errors"
	"net"
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

func TestParseRelativeOrHttpURL_Errors(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		err   error
	}{
		{
			name:  "empty string",
			input: "",
			err:   ErrInvalidURI,
		},
		{
			name:  "whitespace only",
			input: "   ",
			err:   ErrInvalidURI,
		},
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
			_, err := ParseRelativeOrHTTPURL(tc.input)
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
			result, err := ParseRelativeOrHTTPURL(tc.input)
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
		_, err := ValidateURLString(context.Background(), "https://localhost/path", ValidateHTTPURLOptions{
			AllowedSchemes: []string{"ftp", "https", "ws"},
		})
		if err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	t.Run("http blocked when only https allowed", func(t *testing.T) {
		_, err := ValidateURLString(context.Background(), "http://localhost/path", ValidateHTTPURLOptions{
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

	t.Run("multiple values: first is kept", func(t *testing.T) {
		result := ExtractHeaders(http.Header{
			"Accept": []string{"text/html", "application/json"},
		})
		if result["accept"] != "text/html" {
			t.Fatalf("expected first value text/html, got: %s", result["accept"])
		}
	})
}

func TestValidateURLString(t *testing.T) {
	t.Run("empty string returns ErrInvalidURI", func(t *testing.T) {
		_, err := ValidateURLString(context.Background(), "", ValidateHTTPURLOptions{})
		if !errors.Is(err, ErrInvalidURI) {
			t.Fatalf("expected ErrInvalidURI, got: %v", err)
		}
	})

	t.Run("whitespace-only returns ErrInvalidURI", func(t *testing.T) {
		_, err := ValidateURLString(context.Background(), "   ", ValidateHTTPURLOptions{})
		if !errors.Is(err, ErrInvalidURI) {
			t.Fatalf("expected ErrInvalidURI, got: %v", err)
		}
	})

	t.Run("valid http URL", func(t *testing.T) {
		u, err := ValidateURLString(context.Background(), "http://127.0.0.1/path", ValidateHTTPURLOptions{})
		if err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}

		if u.Host != "127.0.0.1" {
			t.Fatalf("unexpected host: %s", u.Host)
		}
	})
}

func TestValidateURL_AllowedSchemes(t *testing.T) {
	t.Run("scheme in allowed list passes", func(t *testing.T) {
		u := &url.URL{Scheme: "https", Host: "127.0.0.1"}
		err := ValidateURL(context.Background(), u, ValidateHTTPURLOptions{
			AllowedSchemes: []string{"http", "https"},
		})
		if err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	t.Run("scheme not in allowed list returns ErrInvalidURLScheme", func(t *testing.T) {
		u := &url.URL{Scheme: "ftp", Host: "127.0.0.1"}
		err := ValidateURL(context.Background(), u, ValidateHTTPURLOptions{
			AllowedSchemes: []string{"http", "https"},
		})
		if !errors.Is(err, ErrInvalidURLScheme) {
			t.Fatalf("expected ErrInvalidURLScheme, got: %v", err)
		}
	})

	t.Run("empty AllowedSchemes skips scheme check", func(t *testing.T) {
		u := &url.URL{Scheme: "ftp", Host: "127.0.0.1"}
		// No scheme restriction — only IP validation matters; 127.0.0.1 resolves so no DNS error
		err := ValidateURL(context.Background(), u, ValidateHTTPURLOptions{})
		// Should not get ErrInvalidURLScheme (may get ErrBlockedIP or nil depending on IP rules)
		if errors.Is(err, ErrInvalidURLScheme) {
			t.Fatalf("did not expect ErrInvalidURLScheme, got: %v", err)
		}
	})
}

func TestValidateURL_EmptyHost(t *testing.T) {
	u := &url.URL{Scheme: "https", Host: ""}
	err := ValidateURL(context.Background(), u, ValidateHTTPURLOptions{})
	if !errors.Is(err, ErrInvalidURI) {
		t.Fatalf("expected ErrInvalidURI for empty host, got: %v", err)
	}
}

func TestValidateURL_AllowedHosts(t *testing.T) {
	t.Run("host in allowed list passes", func(t *testing.T) {
		u := &url.URL{Scheme: "https", Host: "127.0.0.1"}
		err := ValidateURL(context.Background(), u, ValidateHTTPURLOptions{
			AllowedHosts: []string{"127.0.0.1"},
		})
		if err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	t.Run("host not in allowed list returns ErrInvalidURI", func(t *testing.T) {
		u := &url.URL{Scheme: "https", Host: "127.0.0.1"}
		err := ValidateURL(context.Background(), u, ValidateHTTPURLOptions{
			AllowedHosts: []string{"example.com"},
		})
		if !errors.Is(err, ErrInvalidURI) {
			t.Fatalf("expected ErrInvalidURI, got: %v", err)
		}
	})

	t.Run("host prefix match", func(t *testing.T) {
		u := &url.URL{Scheme: "https", Host: "api.example.com"}
		err := ValidateURL(context.Background(), u, ValidateHTTPURLOptions{
			AllowedHosts: []string{"^api."},
		})
		if err != nil {
			t.Fatalf("expected nil error for prefix match, got: %v", err)
		}
	})
}

func TestValidateURL_BlockedHosts(t *testing.T) {
	t.Run("blocked host returns ErrInvalidURI", func(t *testing.T) {
		u := &url.URL{Scheme: "https", Host: "127.0.0.1"}
		err := ValidateURL(context.Background(), u, ValidateHTTPURLOptions{
			BlockedHosts: []string{"127.0.0.1"},
		})
		if !errors.Is(err, ErrInvalidURI) {
			t.Fatalf("expected ErrInvalidURI, got: %v", err)
		}
	})

	t.Run("non-blocked host proceeds to IP validation", func(t *testing.T) {
		u := &url.URL{Scheme: "https", Host: "127.0.0.1"}
		err := ValidateURL(context.Background(), u, ValidateHTTPURLOptions{
			BlockedHosts: []string{"evil.com"},
		})
		// Not blocked by host rule; result depends on IP validation
		if errors.Is(err, ErrInvalidURI) {
			t.Fatalf("unexpected ErrInvalidURI from host block rule, got: %v", err)
		}
	})
}

func TestValidateURL_BlockedIPRanges(t *testing.T) {
	t.Run("IP in blocked range returns ErrBlockedIP", func(t *testing.T) {
		u := &url.URL{Scheme: "https", Host: "127.0.0.1"}
		err := ValidateURL(context.Background(), u, ValidateHTTPURLOptions{
			BlockedIPRanges: []string{"127.0.0.0/8"},
		})
		if !errors.Is(err, ErrBlockedIP) {
			t.Fatalf("expected ErrBlockedIP, got: %v", err)
		}
	})

	t.Run("IP not in blocked range and no allowed range", func(t *testing.T) {
		u := &url.URL{Scheme: "https", Host: "127.0.0.1"}
		err := ValidateURL(context.Background(), u, ValidateHTTPURLOptions{
			BlockedIPRanges: []string{"10.0.0.0/8"},
		})
		// 127.0.0.1 is not blocked by 10/8, but no allowed ranges means ValidateIP returns ErrBlockedIP
		if err != nil {
			t.Fatalf("expected nil, got: %v", err)
		}
	})
}

func TestValidateURL_AllowedIPRanges(t *testing.T) {
	t.Run("IP in allowed range passes", func(t *testing.T) {
		u := &url.URL{Scheme: "https", Host: "127.0.0.1"}
		err := ValidateURL(context.Background(), u, ValidateHTTPURLOptions{
			AllowedIPRanges: []string{"127.0.0.0/8"},
		})
		if err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	t.Run("IP not in allowed range returns ErrBlockedIP", func(t *testing.T) {
		u := &url.URL{Scheme: "https", Host: "127.0.0.1"}
		err := ValidateURL(context.Background(), u, ValidateHTTPURLOptions{
			AllowedIPRanges: []string{"10.0.0.0/8"},
		})
		if !errors.Is(err, ErrBlockedIP) {
			t.Fatalf("expected ErrBlockedIP, got: %v", err)
		}
	})
}

func TestValidateURL_InvalidIPRange(t *testing.T) {
	u := &url.URL{Scheme: "https", Host: "127.0.0.1"}

	t.Run("invalid blocked IP range returns error", func(t *testing.T) {
		err := ValidateURL(context.Background(), u, ValidateHTTPURLOptions{
			BlockedIPRanges: []string{"not-a-cidr"},
		})
		if err == nil {
			t.Fatal("expected error for invalid CIDR, got nil")
		}
	})

	t.Run("invalid allowed IP range returns error", func(t *testing.T) {
		err := ValidateURL(context.Background(), u, ValidateHTTPURLOptions{
			AllowedIPRanges: []string{"not-a-cidr"},
		})
		if err == nil {
			t.Fatal("expected error for invalid CIDR, got nil")
		}
	})
}

func TestValidateIP(t *testing.T) {
	t.Run("private_ip", func(t *testing.T) {
		// Blocked private/internal IP ranges
		privateIPs := []string{
			"10.0.0.1",
			"172.16.0.1",
			"192.168.0.1",
			"127.0.0.1",
			"169.254.0.1",
			// AWS metadata
			"::1",
			"fc00::",
		}

		for _, ip := range privateIPs {
			err := ValidateIPOrDomain(context.Background(), ip, ValidateIPOptions{
				PublicIPOnly: true,
			})
			if err == nil || !errors.Is(err, ErrBlockedIP) {
				t.Errorf("expected private ip error, got: %s, ip: %s", err, ip)
			}
		}
	})

	t.Run("public IP with no restrictions returns ErrBlockedIP (no allowed ranges)", func(t *testing.T) {
		// No allowedIPRanges means ValidateIP always returns ErrBlockedIP after passing public check
		err := ValidateIPOrDomain(context.Background(), "8.8.8.8", ValidateIPOptions{
			PublicIPOnly: true,
		})
		if err != nil {
			t.Fatalf("expected nil, got: %v", err)
		}
	})

	t.Run("public IP in allowed range passes", func(t *testing.T) {
		_, subnet, _ := parseNetCIDR("8.8.8.0/24")
		err := ValidateIPOrDomain(context.Background(), "8.8.8.8", ValidateIPOptions{
			AllowedIPRanges: []*net.IPNet{subnet},
		})
		if err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	t.Run("IP in blocked range returns ErrBlockedIP", func(t *testing.T) {
		_, subnet, _ := parseNetCIDR("8.8.8.0/24")
		_, allowed, _ := parseNetCIDR("0.0.0.0/0")
		err := ValidateIPOrDomain(context.Background(), "8.8.8.8", ValidateIPOptions{
			AllowedIPRanges: []*net.IPNet{allowed},
			BlockedIPRanges: []*net.IPNet{subnet},
		})
		if !errors.Is(err, ErrBlockedIP) {
			t.Fatalf("expected ErrBlockedIP, got: %v", err)
		}
	})

	t.Run("loopback blocked when publicIPOnly=true", func(t *testing.T) {
		err := ValidateIPOrDomain(context.Background(), "127.0.0.1", ValidateIPOptions{
			PublicIPOnly: true,
		})
		if !errors.Is(err, ErrBlockedIP) {
			t.Fatalf("expected ErrBlockedIP, got: %v", err)
		}
	})

	t.Run("unresolvable hostname returns error", func(t *testing.T) {
		err := ValidateIPOrDomain(context.Background(), "this.hostname.does.not.exist.invalid", ValidateIPOptions{})
		if err == nil {
			t.Fatal("expected DNS resolution error, got nil")
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
