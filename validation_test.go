package goutils

import (
	"context"
	"errors"
	"net"
	"net/url"
	"strings"
	"testing"
)

func TestValidateURL_AllowedSchemes(t *testing.T) {
	t.Run("scheme in allowed list passes", func(t *testing.T) {
		u := &url.URL{Scheme: "https", Host: "127.0.0.1"}
		err := ValidateURLWithOptions(context.Background(), u, &ValidateHTTPURLOptions{
			AllowedSchemes: []string{"http", "https"},
		})
		if err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	t.Run("scheme not in allowed list returns ErrInvalidURLScheme", func(t *testing.T) {
		u := &url.URL{Scheme: "ftp", Host: "127.0.0.1"}
		err := ValidateURLWithOptions(context.Background(), u, &ValidateHTTPURLOptions{
			AllowedSchemes: []string{"http", "https"},
		})
		if !errors.Is(err, ErrInvalidURLScheme) {
			t.Fatalf("expected ErrInvalidURLScheme, got: %v", err)
		}
	})

	t.Run("empty AllowedSchemes skips scheme check", func(t *testing.T) {
		u := &url.URL{Scheme: "ftp", Host: "127.0.0.1"}
		// No scheme restriction — only IP validation matters; 127.0.0.1 resolves so no DNS error
		err := ValidateURLWithOptions(context.Background(), u, &ValidateHTTPURLOptions{})
		// Should not get ErrInvalidURLScheme (may get ErrBlockedIP or nil depending on IP rules)
		if errors.Is(err, ErrInvalidURLScheme) {
			t.Fatalf("did not expect ErrInvalidURLScheme, got: %v", err)
		}
	})
}

func TestValidateURL_EmptyHost(t *testing.T) {
	u := &url.URL{Scheme: "https", Host: ""}
	err := ValidateURLWithOptions(context.Background(), u, &ValidateHTTPURLOptions{})
	if !strings.Contains(err.Error(), "invalid URI") {
		t.Fatalf("invalid URI, got: %v", err)
	}
}

func TestValidateURL_AllowedHosts(t *testing.T) {
	t.Run("host in allowed list passes", func(t *testing.T) {
		u := &url.URL{Scheme: "https", Host: "127.0.0.1"}
		err := ValidateURLWithOptions(context.Background(), u, &ValidateHTTPURLOptions{
			AllowedHosts: []string{"127.0.0.1"},
		})
		if err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	t.Run("host not in allowed list returns ErrInvalidURI", func(t *testing.T) {
		u := &url.URL{Scheme: "https", Host: "127.0.0.1"}
		err := ValidateURLWithOptions(context.Background(), u, &ValidateHTTPURLOptions{
			AllowedHosts: []string{"example.com"},
		})
		if !errors.Is(err, ErrInvalidURI) {
			t.Fatalf("expected ErrInvalidURI, got: %v", err)
		}
	})

	t.Run("host prefix match", func(t *testing.T) {
		u := &url.URL{Scheme: "https", Host: "api.example.com"}
		err := ValidateURLWithOptions(context.Background(), u, &ValidateHTTPURLOptions{
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
		err := ValidateURLWithOptions(context.Background(), u, &ValidateHTTPURLOptions{
			BlockedHosts: []string{"127.0.0.1"},
		})
		if !errors.Is(err, ErrInvalidURI) {
			t.Fatalf("expected ErrInvalidURI, got: %v", err)
		}
	})

	t.Run("non-blocked host proceeds to IP validation", func(t *testing.T) {
		u := &url.URL{Scheme: "https", Host: "127.0.0.1"}
		err := ValidateURLWithOptions(context.Background(), u, &ValidateHTTPURLOptions{
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
		err := ValidateURLWithOptions(context.Background(), u, &ValidateHTTPURLOptions{
			BlockedIPRanges: []string{"127.0.0.0/8"},
		})
		if !errors.Is(err, ErrBlockedIP) {
			t.Fatalf("expected ErrBlockedIP, got: %v", err)
		}
	})

	t.Run("IP not in blocked range and no allowed range", func(t *testing.T) {
		u := &url.URL{Scheme: "https", Host: "127.0.0.1"}
		err := ValidateURLWithOptions(context.Background(), u, &ValidateHTTPURLOptions{
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
		err := ValidateURLWithOptions(context.Background(), u, &ValidateHTTPURLOptions{
			AllowedIPRanges: []string{"127.0.0.0/8"},
			PublicIPOnly:    true,
		})
		if err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	t.Run("IP not in allowed range returns ErrBlockedIP", func(t *testing.T) {
		u := &url.URL{Scheme: "https", Host: "127.0.0.1"}
		err := ValidateURLWithOptions(context.Background(), u, &ValidateHTTPURLOptions{
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
		err := ValidateURLWithOptions(context.Background(), u, &ValidateHTTPURLOptions{
			BlockedIPRanges: []string{"not-a-cidr"},
		})
		if err == nil {
			t.Fatal("expected error for invalid CIDR, got nil")
		}
	})

	t.Run("invalid allowed IP range returns error", func(t *testing.T) {
		err := ValidateURLWithOptions(context.Background(), u, &ValidateHTTPURLOptions{
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

	t.Run("public IP with no allowlist is allowed", func(t *testing.T) {
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

	t.Run("allowlist takes precedence over blocklist", func(t *testing.T) {
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

	t.Run("custom LookupIP is used", func(t *testing.T) {
		called := false
		err := ValidateIPOrDomain(context.Background(), "example.com", ValidateIPOptions{
			LookupIP: func(_ context.Context, host string) ([]net.IP, error) {
				called = true
				return []net.IP{net.ParseIP("8.8.8.8")}, nil
			},
		})
		if !called {
			t.Fatal("expected custom LookupIP to be called")
		}
		if err != nil {
			t.Fatalf("expected nil, got: %v", err)
		}
	})
}

func TestValidateJSONPointer(t *testing.T) {
	valid := []string{
		"",
		"/",
		"/foo",
		"/foo/bar",
		"/foo/0",
		"/a~0b",
		"/a~1b",
		"/~01",
		"/ ",
	}
	for _, s := range valid {
		if err := ValidateJSONPointer(s); err != nil {
			t.Errorf("expected valid for %q, got: %v", s, err)
		}
	}

	invalid := []struct {
		input string
		msg   string
	}{
		{"foo", "not starting with /"},
		{"/~", "trailing ~ without 0 or 1"},
		{"/~2", "~ followed by 2"},
		{"/~x", "~ followed by x"},
	}
	for _, tc := range invalid {
		if err := ValidateJSONPointer(tc.input); err == nil {
			t.Errorf("expected invalid for %q (%s), got nil", tc.input, tc.msg)
		}
	}
}

func TestValidateUUID(t *testing.T) {
	valid := []string{
		"00000000-0000-0000-0000-000000000000",
		"550e8400-e29b-41d4-a716-446655440000",
		"6ba7b810-9dad-11d1-80b4-00c04fd430c8",
	}
	for _, s := range valid {
		if err := ValidateUUID(s); err != nil {
			t.Errorf("expected valid UUID for %q, got: %v", s, err)
		}
	}

	invalid := []string{
		"",
		"not-a-uuid",
		"550e8400-e29b-41d4-a716",
		"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
	}
	for _, s := range invalid {
		if err := ValidateUUID(s); err == nil {
			t.Errorf("expected invalid UUID for %q, got nil", s)
		}
	}
}

func TestValidateDurationRFC3339(t *testing.T) {
	valid := []string{
		"P1Y",
		"P1M",
		"P1W",
		"P1D",
		"PT1H",
		"PT1M",
		"PT1S",
		"P1Y2M3DT4H5M6S",
		"P1Y2M",
		"PT30S",
		"P0D",
		"P1DT12H",
	}
	for _, s := range valid {
		if err := ValidateDurationRFC3339(s); err != nil {
			t.Errorf("expected valid duration for %q, got: %v", s, err)
		}
	}

	invalid := []struct {
		input string
		desc  string
	}{
		{"", "empty string"},
		{"1Y", "missing P prefix"},
		{"P", "nothing after P"},
		{"PT", "T with no time elements"},
		{"P1YT", "T with no time elements after T"},
		{"PX", "invalid unit"},
		{"P1H", "H is time unit but before T"},
		{"P1DT1HTT1S", "more than one T"},
		{"P1MT1M2M", "out of order units"},
		{"PT1H2H", "repeated unit"},
	}
	for _, tc := range invalid {
		if err := ValidateDurationRFC3339(tc.input); err == nil {
			t.Errorf("expected invalid duration for %q (%s), got nil", tc.input, tc.desc)
		}
	}
}

func TestValidateIPV4(t *testing.T) {
	valid := []string{
		"0.0.0.0",
		"127.0.0.1",
		"192.168.1.1",
		"255.255.255.255",
		"10.0.0.1",
	}
	for _, s := range valid {
		if err := ValidateIPV4(s); err != nil {
			t.Errorf("expected valid IPv4 for %q, got: %v", s, err)
		}
	}

	invalid := []string{
		"",
		"abc",
		"::1",
		"2001:db8::1",
		"256.0.0.1",
		"1.2.3",
		"1.2.3.4.5",
	}
	for _, s := range invalid {
		if err := ValidateIPV4(s); err == nil {
			t.Errorf("expected invalid IPv4 for %q, got nil", s)
		}
	}
}

func TestValidateIPV6(t *testing.T) {
	valid := []string{
		"::1",
		"2001:db8::1",
		"fe80::1",
		"::ffff:192.0.2.1",
		"2001:0db8:0000:0000:0000:0000:0000:0001",
	}
	for _, s := range valid {
		if err := ValidateIPV6(s); err != nil {
			t.Errorf("expected valid IPv6 for %q, got: %v", s, err)
		}
	}

	invalid := []struct {
		input string
		desc  string
	}{
		{"", "empty"},
		{"127.0.0.1", "IPv4 address"},
		{"abc", "no colon"},
		{":::1", "malformed"},
		{"fe80::1%eth0", "zone ID not allowed"},
	}
	for _, tc := range invalid {
		if err := ValidateIPV6(tc.input); err == nil {
			t.Errorf("expected invalid IPv6 for %q (%s), got nil", tc.input, tc.desc)
		}
	}
}

func TestValidateHostname(t *testing.T) {
	valid := []string{
		"example.com",
		"foo.bar.baz",
		"localhost",
		"my-host",
		"a.b.c.d.e",
		"example.com.",
		"xn--nxasmq6b.com",
	}
	for _, s := range valid {
		if err := ValidateHostname(s); err != nil {
			t.Errorf("expected valid hostname for %q, got: %v", s, err)
		}
	}

	invalid := []struct {
		input string
		desc  string
	}{
		{"-example.com", "label starts with hyphen"},
		{"example-.com", "label ends with hyphen"},
		{"ex ample.com", "space in label"},
		{"ex@mple.com", "@ in label"},
		{"foo..bar", "empty label (consecutive dots)"},
		{"." + strings.Repeat("a", 64), "label too long (64 chars)"},
	}
	for _, tc := range invalid {
		if err := ValidateHostname(tc.input); err == nil {
			t.Errorf("expected invalid hostname for %q (%s), got nil", tc.input, tc.desc)
		}
	}
}

func TestValidateEmail(t *testing.T) {
	valid := []string{
		"user@example.com",
		"user.name+tag@example.co.uk",
		"user@[127.0.0.1]",
		"user@[IPv6:::1]",
		"\"user\"@example.com",
		"user123@sub.domain.org",
		"a@b.io",
		"!#$%&'*+-/=?^_`{|}~@example.com",
	}
	for _, s := range valid {
		if err := ValidateEmail(s); err != nil {
			t.Errorf("expected valid email for %q, got: %v", s, err)
		}
	}

	invalid := []struct {
		input string
		desc  string
	}{
		{"", "empty"},
		{"userexample.com", "missing @"},
		{".user@example.com", "local starts with dot"},
		{"user.@example.com", "local ends with dot"},
		{"user..name@example.com", "consecutive dots in local"},
		{"user @example.com", "space in local"},
		{"user@[999.0.0.1]", "invalid IPv4 in domain bracket"},
		{"user@[IPv6:invalid]", "invalid IPv6 in domain bracket"},
		{"user@-example.com", "domain label starts with hyphen"},
	}
	for _, tc := range invalid {
		if err := ValidateEmail(tc.input); err == nil {
			t.Errorf("expected invalid email for %q (%s), got nil", tc.input, tc.desc)
		}
	}
}

func TestValidateDate(t *testing.T) {
	valid := []string{
		"2024-01-01",
		"2000-12-31",
		"1999-06-15",
	}
	for _, s := range valid {
		if err := ValidateDate(s); err != nil {
			t.Errorf("expected valid date for %q, got: %v", s, err)
		}
	}

	invalid := []string{
		"",
		"2024-13-01",
		"2024-00-01",
		"2024-01-32",
		"20240101",
		"not-a-date",
		"2024/01/01",
	}
	for _, s := range invalid {
		if err := ValidateDate(s); err == nil {
			t.Errorf("expected invalid date for %q, got nil", s)
		}
	}
}

func TestValidateTime(t *testing.T) {
	valid := []string{
		"00:00:00Z",
		"23:59:59Z",
		"12:30:00+05:30",
		"08:00:00-07:00",
		"12:00:00.999Z",
	}
	for _, s := range valid {
		if err := ValidateTime(s); err != nil {
			t.Errorf("expected valid time for %q, got: %v", s, err)
		}
	}

	invalid := []string{
		"",
		"24:00:00Z",
		"12:60:00Z",
		"12:00:60Z",
		"not-a-time",
	}
	for _, s := range invalid {
		if err := ValidateTime(s); err == nil {
			t.Errorf("expected invalid time for %q, got nil", s)
		}
	}
}

func TestValidateDateTime(t *testing.T) {
	valid := []string{
		"2024-01-15T12:30:00Z",
		"2000-12-31T23:59:59+05:30",
		"1999-06-15T00:00:00-07:00",
	}
	for _, s := range valid {
		if err := ValidateDateTime(s); err != nil {
			t.Errorf("expected valid datetime for %q, got: %v", s, err)
		}
	}

	invalid := []string{
		"",
		"12:30:00Z",
		"2024-13-01T00:00:00Z",
		"not-a-datetime",
		"2024-01-15T25:00:00Z",
	}
	for _, s := range invalid {
		if err := ValidateDateTime(s); err == nil {
			t.Errorf("expected invalid datetime for %q, got nil", s)
		}
	}
}

func TestValidateURI(t *testing.T) {
	valid := []string{
		"https://example.com",
		"http://localhost:8080/path?q=1#frag",
		"ftp://files.example.com/pub",
	}
	for _, s := range valid {
		if err := ValidateURI(s); err != nil {
			t.Errorf("expected valid URI for %q, got: %v", s, err)
		}
	}

	invalid := []struct {
		input string
		desc  string
	}{
		{"/relative/path", "relative URI not allowed"},
		{"just-a-string", "no scheme"},
		{"://missing-scheme", "empty scheme"},
	}
	for _, tc := range invalid {
		if err := ValidateURI(tc.input); err == nil {
			t.Errorf("expected invalid URI for %q (%s), got nil", tc.input, tc.desc)
		}
	}
}

func BenchmarkValidateJSONPointer(b *testing.B) {
	for b.Loop() {
		ValidateJSONPointer("/a~0b/c~1d")
	}
}

func BenchmarkValidateUUID(b *testing.B) {
	for b.Loop() {
		ValidateUUID("550e8400-e29b-41d4-a716-446655440000")
	}
}

func BenchmarkValidateDurationRFC3339(b *testing.B) {
	for b.Loop() {
		ValidateDurationRFC3339("P1Y2M3DT4H5M6S")
	}
}

func BenchmarkValidateIPV4(b *testing.B) {
	for b.Loop() {
		ValidateIPV4("192.168.1.1")
	}
}

func BenchmarkValidateIPV6(b *testing.B) {
	b.Run("valid_full", func(b *testing.B) {
		for b.Loop() {
			ValidateIPV6("2001:0db8:0000:0000:0000:0000:0000:0001")
		}
	})
	b.Run("valid_compressed", func(b *testing.B) {
		for b.Loop() {
			ValidateIPV6("2001:db8::1")
		}
	})
}

func BenchmarkValidateHostname(b *testing.B) {
	b.Run("valid_simple", func(b *testing.B) {
		for b.Loop() {
			ValidateHostname("example.com")
		}
	})
	b.Run("valid_subdomain", func(b *testing.B) {
		for b.Loop() {
			ValidateHostname("api.sub.example.com")
		}
	})
}

func BenchmarkValidateEmail(b *testing.B) {
	b.Run("valid_simple", func(b *testing.B) {
		for b.Loop() {
			ValidateEmail("user@example.com")
		}
	})
	b.Run("valid_complex", func(b *testing.B) {
		for b.Loop() {
			ValidateEmail("user.name+tag@sub.example.co.uk")
		}
	})
	b.Run("valid_ip_domain", func(b *testing.B) {
		for b.Loop() {
			ValidateEmail("user@[127.0.0.1]")
		}
	})
}

func BenchmarkValidateDate(b *testing.B) {
	for b.Loop() {
		ValidateDate("2024-06-15")
	}
}

func BenchmarkValidateTime(b *testing.B) {
	b.Run("valid_utc", func(b *testing.B) {
		for b.Loop() {
			ValidateTime("12:30:00Z")
		}
	})
	b.Run("valid_offset", func(b *testing.B) {
		for b.Loop() {
			ValidateTime("08:00:00-07:00")
		}
	})
}

func BenchmarkValidateDateTime(b *testing.B) {
	b.Run("valid", func(b *testing.B) {
		for b.Loop() {
			ValidateDateTime("2024-01-15T12:30:00Z")
		}
	})
	b.Run("valid_offset", func(b *testing.B) {
		for b.Loop() {
			ValidateDateTime("2000-12-31T23:59:59+05:30")
		}
	})
}

func BenchmarkValidateURI(b *testing.B) {
	b.Run("valid_https", func(b *testing.B) {
		for b.Loop() {
			ValidateURI("https://example.com/path?q=1#frag")
		}
	})
}
