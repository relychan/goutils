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
	"fmt"
	"net"
	"net/url"
	"strings"
)

var (
	// RFC6598 Carrier-Grade NAT.
	cgNATSubnet = mustParseCIDR("100.64.0.0/10")
	httpSchemes = []string{"http", "https"}
)

// ParsePathOrHTTPURL validates and parses a path or HTTP URL.
func ParsePathOrHTTPURL(input string) (*url.URL, error) {
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
		return new(url.URL), nil
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

	if StringContainsCTLByte(input) {
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

// ParseURL parses and validate the input string to be a valid URI.
func ParseURL(input string) (*url.URL, error) {
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

	if strings.Contains(hostname, ":") {
		if !strings.Contains(parsedURL.Host, "[") || !strings.Contains(parsedURL.Host, "]") {
			return nil, ErrInvalidURI
		}

		err := ValidateIPV6(hostname)
		if err != nil {
			return nil, err
		}
	}

	return parsedURL, nil
}

// ParseHTTPURL parses and validate the input string to be a valid URI and have http(s) scheme.
func ParseHTTPURL(input string) (*url.URL, error) {
	parsedURL, err := ParseURL(input)
	if err != nil {
		return nil, err
	}

	err = validateURLScheme(parsedURL, httpSchemes)
	if err != nil {
		return nil, err
	}

	return parsedURL, nil
}

// ParseAndValidateURL parses and validates URL from a string. Returns the parsed URL and an error.
func ParseAndValidateURL(
	ctx context.Context,
	urlStr string,
	options ValidateHTTPURLOptions,
) (*url.URL, error) {
	parsedURL, err := ParseURL(urlStr)
	if err != nil {
		return nil, err
	}

	err = ValidateURLWithOptions(ctx, parsedURL, options)
	if err != nil {
		return nil, err
	}

	return parsedURL, nil
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
			return nil, fmt.Errorf("failed to parse IP range %q: %w", rawIPRange, err)
		}

		results[i] = ip
	}

	return results, nil
}

func mustParseCIDR(cidr string) *net.IPNet {
	_, network, err := net.ParseCIDR(cidr)
	if err != nil {
		panic(fmt.Sprintf("invalid CIDR %q: %v", cidr, err))
	}

	return network
}
