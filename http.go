package goutils

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"slices"
	"strings"
)

// ParseRelativeOrHTTPURL validates and parses relative or HTTP URL.
func ParseRelativeOrHTTPURL(input string) (*url.URL, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, ErrInvalidURI
	}

	schemeIndex := strings.Index(input, "://")
	if schemeIndex == 0 {
		return nil, ErrInvalidURLScheme
	}

	if schemeIndex > 0 {
		scheme := input[:schemeIndex]
		if scheme != "http" && scheme != "https" {
			return nil, ErrInvalidURLScheme
		}

		return url.Parse(input)
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
func ParseHTTPURL(input string, options ValidateHTTPURLOptions) (*url.URL, error) {
	if len(options.AllowedSchemes) == 0 {
		options.AllowedSchemes = []string{"http", "https"}
	} else {
		// remove schemes that are not http(s)
		schemes := []string{}

		for _, scheme := range options.AllowedSchemes {
			if scheme == "http" || scheme == "https" {
				schemes = append(schemes, scheme)
			}
		}

		options.AllowedSchemes = schemes
	}

	return ValidateURLString(input, options)
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
	ResolveDNS      bool
	PublicIPOnly    bool
	AllowedSchemes  []string
	AllowedIPRanges []string
	BlockedIPRanges []string
	AllowedHosts    []string
	BlockedHosts    []string
}

// ValidateURLString parses and validates URL from a string. Returns the parsed URL and an error.
func ValidateURLString(urlStr string, options ValidateHTTPURLOptions) (*url.URL, error) {
	urlStr = strings.TrimSpace(urlStr)
	if urlStr == "" {
		return nil, ErrInvalidURI
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	return parsedURL, ValidateURL(parsedURL, options)
}

// ValidateURL parses and validates URL.
func ValidateURL(uri *url.URL, options ValidateHTTPURLOptions) error {
	if len(options.AllowedSchemes) > 0 && !slices.Contains(options.AllowedSchemes, uri.Scheme) {
		return fmt.Errorf(
			"%w. Accept one of %v, got: %s",
			ErrInvalidURLScheme,
			options.AllowedSchemes,
			uri.Scheme,
		)
	}

	// Extract hostname without port
	hostname := uri.Hostname()
	if hostname == "" {
		return ErrInvalidURI
	}

	err := validateHost(uri.Host, hostname, &options)
	if err != nil {
		return err
	}

	for _, expr := range options.BlockedHosts {
		re, err := NewRegexpMatcher(expr)
		if err != nil {
			return fmt.Errorf("failed to parse allowed host rule: %w", err)
		}

		if re.MatchString(hostname) || re.MatchString(uri.Host) {
			return fmt.Errorf("%w: host is blocked", ErrInvalidURI)
		}
	}

	if !options.ResolveDNS && !options.PublicIPOnly &&
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

	return ValidateIP(hostname, options.PublicIPOnly, allowedIPRanges, blockedIPRanges)
}

// ValidateIP checks if the IP is valid.
func ValidateIP(
	hostname string,
	publicIPOnly bool,
	allowedIPRanges []*net.IPNet,
	blockedIPRanges []*net.IPNet,
) error {
	// Resolve IP addresses
	ips, err := net.LookupIP(hostname) //nolint:noctx
	if err != nil {
		// Block on DNS resolution failure
		return err
	}

	// Check each IP against blocked ranges
	for _, ip := range ips {
		if publicIPOnly && (ip.IsPrivate() || ip.IsLoopback() || ip.IsUnspecified() ||
			ip.IsLinkLocalUnicast() ||
			ip.IsLinkLocalMulticast()) {
			return fmt.Errorf("%w: allow public IP only", ErrBlockedIP)
		}

		for _, subnet := range blockedIPRanges {
			if subnet.Contains(ip) {
				return ErrBlockedIP
			}
		}
	}

	if len(allowedIPRanges) == 0 {
		return nil
	}

	// Check each IP against allowed ranges
	for _, ip := range ips {
		for _, subnet := range allowedIPRanges {
			if subnet.Contains(ip) {
				return nil
			}
		}
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
