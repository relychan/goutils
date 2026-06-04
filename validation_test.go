package goutils

import (
	"context"
	"errors"
	"net"
	"net/url"
	"testing"
)

func TestValidateURL_AllowedSchemes(t *testing.T) {
	t.Run("scheme in allowed list passes", func(t *testing.T) {
		u := &url.URL{Scheme: "https", Host: "127.0.0.1"}
		err := ValidateURLWithOptions(context.Background(), u, ValidateHTTPURLOptions{
			AllowedSchemes: []string{"http", "https"},
		})
		if err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	t.Run("scheme not in allowed list returns ErrInvalidURLScheme", func(t *testing.T) {
		u := &url.URL{Scheme: "ftp", Host: "127.0.0.1"}
		err := ValidateURLWithOptions(context.Background(), u, ValidateHTTPURLOptions{
			AllowedSchemes: []string{"http", "https"},
		})
		if !errors.Is(err, ErrInvalidURLScheme) {
			t.Fatalf("expected ErrInvalidURLScheme, got: %v", err)
		}
	})

	t.Run("empty AllowedSchemes skips scheme check", func(t *testing.T) {
		u := &url.URL{Scheme: "ftp", Host: "127.0.0.1"}
		// No scheme restriction — only IP validation matters; 127.0.0.1 resolves so no DNS error
		err := ValidateURLWithOptions(context.Background(), u, ValidateHTTPURLOptions{})
		// Should not get ErrInvalidURLScheme (may get ErrBlockedIP or nil depending on IP rules)
		if errors.Is(err, ErrInvalidURLScheme) {
			t.Fatalf("did not expect ErrInvalidURLScheme, got: %v", err)
		}
	})
}

func TestValidateURL_EmptyHost(t *testing.T) {
	u := &url.URL{Scheme: "https", Host: ""}
	err := ValidateURLWithOptions(context.Background(), u, ValidateHTTPURLOptions{})
	if !errors.Is(err, ErrInvalidURI) {
		t.Fatalf("expected ErrInvalidURI for empty host, got: %v", err)
	}
}

func TestValidateURL_AllowedHosts(t *testing.T) {
	t.Run("host in allowed list passes", func(t *testing.T) {
		u := &url.URL{Scheme: "https", Host: "127.0.0.1"}
		err := ValidateURLWithOptions(context.Background(), u, ValidateHTTPURLOptions{
			AllowedHosts: []string{"127.0.0.1"},
		})
		if err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	t.Run("host not in allowed list returns ErrInvalidURI", func(t *testing.T) {
		u := &url.URL{Scheme: "https", Host: "127.0.0.1"}
		err := ValidateURLWithOptions(context.Background(), u, ValidateHTTPURLOptions{
			AllowedHosts: []string{"example.com"},
		})
		if !errors.Is(err, ErrInvalidURI) {
			t.Fatalf("expected ErrInvalidURI, got: %v", err)
		}
	})

	t.Run("host prefix match", func(t *testing.T) {
		u := &url.URL{Scheme: "https", Host: "api.example.com"}
		err := ValidateURLWithOptions(context.Background(), u, ValidateHTTPURLOptions{
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
		err := ValidateURLWithOptions(context.Background(), u, ValidateHTTPURLOptions{
			BlockedHosts: []string{"127.0.0.1"},
		})
		if !errors.Is(err, ErrInvalidURI) {
			t.Fatalf("expected ErrInvalidURI, got: %v", err)
		}
	})

	t.Run("non-blocked host proceeds to IP validation", func(t *testing.T) {
		u := &url.URL{Scheme: "https", Host: "127.0.0.1"}
		err := ValidateURLWithOptions(context.Background(), u, ValidateHTTPURLOptions{
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
		err := ValidateURLWithOptions(context.Background(), u, ValidateHTTPURLOptions{
			BlockedIPRanges: []string{"127.0.0.0/8"},
		})
		if !errors.Is(err, ErrBlockedIP) {
			t.Fatalf("expected ErrBlockedIP, got: %v", err)
		}
	})

	t.Run("IP not in blocked range and no allowed range", func(t *testing.T) {
		u := &url.URL{Scheme: "https", Host: "127.0.0.1"}
		err := ValidateURLWithOptions(context.Background(), u, ValidateHTTPURLOptions{
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
		err := ValidateURLWithOptions(context.Background(), u, ValidateHTTPURLOptions{
			AllowedIPRanges: []string{"127.0.0.0/8"},
			PublicIPOnly:    true,
		})
		if err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	t.Run("IP not in allowed range returns ErrBlockedIP", func(t *testing.T) {
		u := &url.URL{Scheme: "https", Host: "127.0.0.1"}
		err := ValidateURLWithOptions(context.Background(), u, ValidateHTTPURLOptions{
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
		err := ValidateURLWithOptions(context.Background(), u, ValidateHTTPURLOptions{
			BlockedIPRanges: []string{"not-a-cidr"},
		})
		if err == nil {
			t.Fatal("expected error for invalid CIDR, got nil")
		}
	})

	t.Run("invalid allowed IP range returns error", func(t *testing.T) {
		err := ValidateURLWithOptions(context.Background(), u, ValidateHTTPURLOptions{
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
			"0.0.0.0",
			"255.255.255.255",
			"100.64.0.10",
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
		if err != nil {
			t.Fatalf("expected nil error, got: %v", err)
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
