package goutils

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"slices"
	"strings"
)

var (
	// RFC6598 Carrier-Grade NAT.
	cgNATSubnet = mustParseCIDR("100.64.0.0/10")
	httpSchemes = []string{"http", "https"}
)

// ParseRelativeOrHTTPURL validates and parses a path or HTTP URL.
func ParseRelativeOrHTTPURL(input string) (*url.URL, error) {
	parsedURL, err := ParsePathOrURL(input)
	if err != nil {
		return nil, err
	}

	// Returns if the parsedURL is a path.
	if parsedURL.Scheme == "" {
		return parsedURL, nil
	}

	return parsedURL, validateURLScheme(parsedURL, httpSchemes)
}

// ParsePathOrURL validates and parses a path or URL.
func ParsePathOrURL(input string) (*url.URL, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, ErrInvalidURI
	}

	schemeIndex := strings.Index(input, "://")
	if schemeIndex == 0 {
		return nil, ErrInvalidURLScheme
	}

	if schemeIndex > 0 {
		parsedURL, err := url.Parse(input)
		if err != nil {
			return nil, err
		}

		hostname := parsedURL.Hostname()
		if hostname == "" {
			return nil, ErrInvalidURI
		}

		return parsedURL, nil
	}

	if !filepath.IsAbs(input) && !filepath.IsLocal(input) {
		return nil, ErrInvalidURI
	}

	u, frag, _ := strings.Cut(input, "#")
	urlPath, query, _ := strings.Cut(u, "?")

	result := &url.URL{
		Path:       urlPath,
		RawQuery:   query,
		ForceQuery: query != "",
		Fragment:   frag,
	}

	return result, nil
}

// ParseHTTPURL parses and validate the input string to have http(s) scheme.
func ParseHTTPURL(input string) (*url.URL, error) {
	urlStr := strings.TrimSpace(input)
	if urlStr == "" {
		return nil, ErrInvalidURI
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	hostname := parsedURL.Hostname()
	if hostname == "" {
		return nil, ErrInvalidURI
	}

	return parsedURL, validateURLScheme(parsedURL, httpSchemes)
}

// ExtractHeaders converts the http.Header to string map with lowercase header names.
func ExtractHeaders(headers http.Header) map[string]string {
	result := make(map[string]string)

	for key, header := range headers {
		if len(header) == 0 {
			continue
		}

		result[strings.ToLower(key)] = header[0]
	}

	return result
}

// ValidateHTTPURLOptions represent URL validation options.
type ValidateHTTPURLOptions struct {
	AllowedSchemes  []string
	AllowedHosts    []string
	BlockedHosts    []string
	PublicIPOnly    bool
	AllowedIPRanges []string
	BlockedIPRanges []string
	// Custom lookup IP function.
	LookupIP func(ctx context.Context, host string) ([]net.IP, error)
}

// ValidateURLString parses and validates URL from a string. Returns the parsed URL and an error.
func ValidateURLString(
	ctx context.Context,
	urlStr string,
	options ValidateHTTPURLOptions,
) (*url.URL, error) {
	urlStr = strings.TrimSpace(urlStr)
	if urlStr == "" {
		return nil, ErrInvalidURI
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	return parsedURL, ValidateURL(ctx, parsedURL, options)
}

// ValidateURL parses and validates URL.
func ValidateURL(ctx context.Context, uri *url.URL, options ValidateHTTPURLOptions) error {
	err := validateURLScheme(uri, options.AllowedSchemes)
	if err != nil {
		return err
	}

	// Extract hostname without port
	hostname := uri.Hostname()
	if hostname == "" {
		return ErrInvalidURI
	}

	err = validateHost(uri.Host, hostname, &options)
	if err != nil {
		return err
	}

	if !options.PublicIPOnly &&
		len(options.AllowedIPRanges) == 0 && len(options.BlockedIPRanges) == 0 {
		return nil
	}

	allowedIPRanges, err := parseIPRanges(options.AllowedIPRanges)
	if err != nil {
		return err
	}

	blockedIPRanges, err := parseIPRanges(options.BlockedIPRanges)
	if err != nil {
		return err
	}

	return ValidateIPOrDomain(ctx, hostname, ValidateIPOptions{
		PublicIPOnly:    options.PublicIPOnly,
		AllowedIPRanges: allowedIPRanges,
		BlockedIPRanges: blockedIPRanges,
		LookupIP:        options.LookupIP,
	})
}

// ValidateIPOptions represent URL validation options.
type ValidateIPOptions struct {
	// Block all private IPs.
	PublicIPOnly bool
	// IP ranges to allow.
	AllowedIPRanges []*net.IPNet
	// IP ranges to block.
	BlockedIPRanges []*net.IPNet
	// Custom lookup IP function.
	LookupIP func(ctx context.Context, host string) ([]net.IP, error)
}

// ValidateIPOrDomain checks if the IP string or IP of domain is valid for SSRF protection.
// If the input string is a domain, lookup the IP from it before validation.
func ValidateIPOrDomain(
	ctx context.Context,
	domainOrIP string,
	options ValidateIPOptions,
) error {
	// Resolve IP addresses
	var ips []net.IP

	var err error

	if options.LookupIP != nil {
		ips, err = options.LookupIP(ctx, domainOrIP)
	} else {
		ips, err = net.DefaultResolver.LookupIP(ctx, "ip", domainOrIP)
	}

	if err != nil {
		// Block on DNS resolution failure
		return err
	}

	// Check each IP against blocked ranges
	for _, ip := range ips {
		err := ValidateIP(ip, options)
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateIP checks if the IP is valid for SSRF protection.
// Note: the allowed ranges option is the highest priority to bypass other rules.
func ValidateIP(ip net.IP, options ValidateIPOptions) error {
	for _, subnet := range options.AllowedIPRanges {
		if subnet.Contains(ip) {
			return nil
		}
	}

	if options.PublicIPOnly && (ip.IsPrivate() ||
		!ip.IsGlobalUnicast() ||
		ip.IsLinkLocalMulticast() ||
		cgNATSubnet.Contains(ip)) {
		return ErrBlockedIP
	}

	for _, subnet := range options.BlockedIPRanges {
		if subnet.Contains(ip) {
			return ErrBlockedIP
		}
	}

	// The IP is valid if allowed IP ranges are empty.
	if len(options.AllowedIPRanges) == 0 {
		return nil
	}

	return ErrBlockedIP
}

// ParseSubnet parses the subnet from a raw string.
func ParseSubnet(value string) (*net.IPNet, error) {
	if value == "" {
		return nil, ErrInvalidSubnet
	}

	if !strings.Contains(value, "/") {
		ip := net.ParseIP(value)
		if ip == nil {
			return nil, ErrInvalidSubnet
		}

		if ip.To4() != nil {
			value += "/32"
		} else {
			value += "/128"
		}
	}

	_, subnet, err := net.ParseCIDR(value)
	if err != nil {
		return nil, err
	}

	return subnet, err
}

func parseIPRanges(ipRanges []string) ([]*net.IPNet, error) {
	results := make([]*net.IPNet, len(ipRanges))

	for i, rawIPRange := range ipRanges {
		ip, err := ParseSubnet(rawIPRange)
		if err != nil {
			return nil, fmt.Errorf("failed to parse IP range %s: %w", rawIPRange, err)
		}

		results[i] = ip
	}

	return results, nil
}

func validateHost(host, hostname string, options *ValidateHTTPURLOptions) error {
	for _, expr := range options.BlockedHosts {
		re, err := NewRegexpMatcher(expr)
		if err != nil {
			return fmt.Errorf("failed to parse allowed host rule: %w", err)
		}

		if re.MatchString(hostname) || re.MatchString(host) {
			return fmt.Errorf("%w: host is blocked", ErrInvalidURI)
		}
	}

	if len(options.AllowedHosts) == 0 {
		return nil
	}

	for _, expr := range options.AllowedHosts {
		re, err := NewRegexpMatcher(expr)
		if err != nil {
			return fmt.Errorf("failed to parse allowed host rule: %w", err)
		}

		if re.MatchString(hostname) || re.MatchString(host) {
			return nil
		}
	}

	return fmt.Errorf("%w: host is not allowed", ErrInvalidURI)
}

func validateURLScheme(uri *url.URL, allowedSchemes []string) error {
	if len(allowedSchemes) > 0 && !slices.Contains(allowedSchemes, uri.Scheme) {
		return fmt.Errorf(
			"%w. Accept one of %v, got: %s",
			ErrInvalidURLScheme,
			allowedSchemes,
			uri.Scheme,
		)
	}

	return nil
}

func mustParseCIDR(cidr string) *net.IPNet {
	_, network, err := net.ParseCIDR(cidr)
	if err != nil {
		panic(fmt.Sprintf("invalid CIDR %q: %v", cidr, err))
	}

	return network
}
