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
	"slices"
	"strings"

	"github.com/relychan/goutils/httperror"
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

	err = validateURLScheme(parsedURL, httpSchemes)
	if err != nil {
		return nil, err
	}

	return parsedURL, nil
}

// ParsePathOrURL validates and parses a path or URL.
func ParsePathOrURL(input string) (*url.URL, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return new(url.URL), nil
	}

	schemeIndex := strings.IndexRune(input, ':')
	if schemeIndex > 0 {
		// If ':' appears after a path/query/fragment delimiter, it's part of the path, not a scheme.
		if delim := strings.IndexAny(input, "/?#"); delim != -1 && delim < schemeIndex {
			schemeIndex = -1
		}
	}

	if schemeIndex > 0 {
		// Treat Windows drive paths (e.g. C:\\path or C:/path) as paths, not URL schemes.
		if len(input) >= 3 && schemeIndex == 1 && (input[2] == '\\' || input[2] == '/') &&
			((input[0] >= 'A' && input[0] <= 'Z') || (input[0] >= 'a' && input[0] <= 'z')) {
			schemeIndex = -1
		}
	}

	if schemeIndex > 0 {
		if !strings.HasPrefix(input[schemeIndex+1:], "//") {
			// The authority could be missing because of missing slashes.
			return nil, ErrInvalidURLScheme
		}

		return ParseURL(input)
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

// ParseURI parses and validate the input string to be a valid URI.
//
//	scheme:[//authority]/path[?query][#fragment]
func ParseURI(input string) (*url.URL, error) {
	parsedURI, _, err := parseURIAndHostname(input)
	if err != nil {
		return nil, err
	}

	return parsedURI, nil
}

// ParseURL parses and validate the input string to be a valid URL.
// Unlike URI, the URL requires an explicit authority.
//
//	scheme://authority[/path][?query][#fragment]
func ParseURL(input string) (*url.URL, error) {
	result, err := parseAndValidateURL(input)
	if err != nil {
		return nil, err
	}

	return result, nil
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

// ParseAndValidateURLWithOptions parses and validates URL from a string with options.
// Returns the parsed URL and an error.
func ParseAndValidateURLWithOptions(
	ctx context.Context,
	urlStr string,
	options *ValidateHTTPURLOptions,
) (*url.URL, error) {
	parsedURL, err := ParseURL(urlStr)
	if err != nil {
		return nil, err
	}

	validatedErr := ValidateURLWithOptions(ctx, parsedURL, options)
	if validatedErr != nil {
		return nil, validatedErr
	}

	return parsedURL, nil
}

// ValidateURI checks if the input string is a valid URI.
func ValidateURI(s string) *httperror.ValidationError {
	_, err := parseAndValidateURI(s)

	return err
}

// ValidateURL checks if the input string is a valid URL.
func ValidateURL(s string) *httperror.ValidationError {
	_, err := parseAndValidateURL(s)

	return err
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

// ValidateURLWithOptions parses and validates URL with options.
func ValidateURLWithOptions(
	ctx context.Context,
	uri *url.URL,
	options *ValidateHTTPURLOptions,
) error {
	if options == nil {
		return nil
	}

	err := validateURLScheme(uri, options.AllowedSchemes)
	if err != nil {
		return err
	}

	// Extract hostname without port
	hostname := uri.Hostname()
	if hostname == "" {
		return ErrInvalidURI
	}

	err = validateHostWithOptions(uri.Host, hostname, options)
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

func parseAndValidateURI(s string) (*url.URL, *httperror.ValidationError) {
	input := strings.TrimSpace(s)

	schemeIndex := strings.IndexByte(input, ':')
	if schemeIndex < 1 || len(input)-schemeIndex <= 1 {
		return nil, &httperror.ValidationError{
			Code:   ErrCodeInvalidURI,
			Detail: "Invalid URI syntax",
		}
	}

	uri, _, err := parseURIAndHostname(s)

	return uri, err
}

func validateHostWithOptions(host, hostname string, options *ValidateHTTPURLOptions) error {
	validatedError := ValidateHostname(hostname)
	if validatedError != nil {
		return validatedError
	}

	for _, expr := range options.BlockedHosts {
		re, err := NewRegexpMatcher(expr)
		if err != nil {
			return fmt.Errorf("failed to parse blocked host rule: %w", err)
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
	if len(allowedSchemes) > 0 && !slices.ContainsFunc(allowedSchemes, func(item string) bool {
		return strings.EqualFold(item, uri.Scheme)
	}) {
		return fmt.Errorf(
			"%w. Accept one of %v, got: %q",
			ErrInvalidURLScheme,
			allowedSchemes,
			uri.Scheme,
		)
	}

	return nil
}

func parseURIAndHostname(input string) (*url.URL, string, *httperror.ValidationError) {
	parsedURI, err := url.Parse(input)
	if err != nil {
		return nil, "", &httperror.ValidationError{
			Code:   ErrCodeInvalidURI,
			Detail: err.Error(),
		}
	}

	if parsedURI.Host == "" {
		return parsedURI, "", nil
	}

	hostname := parsedURI.Hostname()
	if hostname == "" {
		return nil, "", &httperror.ValidationError{
			Code:   ErrCodeInvalidURI,
			Detail: "Invalid URI. Hostname is empty",
		}
	}

	if strings.Contains(hostname, ":") {
		if !strings.Contains(parsedURI.Host, "[") || !strings.Contains(parsedURI.Host, "]") {
			return nil, "", &httperror.ValidationError{
				Code:   ErrCodeInvalidURI,
				Detail: "Invalid URI. IPv6 in hostname is invalid",
			}
		}

		err := ValidateIPV6(hostname)
		if err != nil {
			err.Code = ErrCodeInvalidURI
			err.Detail = "Invalid URI. " + err.Detail

			return nil, "", err
		}

		return parsedURI, hostname, nil
	}

	validatedError := ValidateHostname(hostname)
	if validatedError != nil {
		return nil, "", validatedError
	}

	return parsedURI, hostname, nil
}

func parseAndValidateURL(input string) (*url.URL, *httperror.ValidationError) { //nolint:funlen
	uriStr := strings.TrimSpace(input)
	if uriStr == "" {
		return nil, &httperror.ValidationError{
			Code:   ErrCodeInvalidURI,
			Detail: "Invalid URL. The input string is empty",
		}
	}

	schemeIndex := strings.Index(input, "://")
	if schemeIndex <= 0 || len(uriStr)-schemeIndex <= 1 {
		return nil, &httperror.ValidationError{
			Code:   ErrCodeInvalidURI,
			Detail: "Invalid URL syntax",
		}
	}

	parsedURI, hostname, err := parseURIAndHostname(input)
	if err != nil {
		return nil, err
	}

	if hostname == "" {
		return nil, &httperror.ValidationError{
			Code:   ErrCodeInvalidURI,
			Detail: "Invalid URL. Hostname is empty",
		}
	}

	if parsedURI.Path == "" || parsedURI.Path == "/" {
		return parsedURI, nil
	}

	// validate invalid path patterns
	uriPath := parsedURI.Path
	if uriPath[0] == '/' {
		uriPath = uriPath[1:]
	}

	for uriPath != "" {
		slashIndex := strings.IndexByte(uriPath, '/')
		if slashIndex == 0 {
			return nil, &httperror.ValidationError{
				Code:   ErrCodeInvalidURI,
				Detail: "Invalid double slashes in the URL path syntax",
			}
		}

		part := uriPath

		if slashIndex != -1 {
			part = uriPath[:slashIndex]
			uriPath = uriPath[slashIndex+1:]
		} else {
			uriPath = ""
		}

		if part == "*" || StringAllRune(part, '.') {
			return nil, &httperror.ValidationError{
				Code:   ErrCodeInvalidURI,
				Detail: "Invalid URL path syntax",
			}
		}
	}

	if parsedURI.Path[0] != '/' {
		parsedURI.Path = "/" + parsedURI.Path
	}

	return parsedURI, nil
}
