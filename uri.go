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
	"strconv"
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

		return nil, &httperror.ValidationError{
			Code:   ErrCodeInvalidURIScheme,
			Detail: "Invalid HTTP scheme. Expected http(s), got " + strconv.Quote(parsedURL.Scheme),
		}
	}

	return parsedURL, nil
}

// ParsePathOrURL validates and parses a path or URL.
func ParsePathOrURL(input string) (*url.URL, error) {
	schemeIndex := strings.IndexRune(input, ':')
	if schemeIndex == 0 {
		// The authority could be missing because of missing slashes.
		return nil, &httperror.ValidationError{
			Code:   ErrCodeInvalidURIScheme,
			Detail: "Invalid URL. Scheme is empty",
		}
	}

	if schemeIndex > 0 {
		return ParseURL(input)
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return new(url.URL), nil
	}

	if StringContainsCTLByte(input) {
		return nil, &httperror.ValidationError{
			Code:   ErrCodeInvalidPath,
			Detail: "Path contains invalid characters",
		}
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
func ParseHTTPURL(s string) (*url.URL, error) {
	input := strings.TrimSpace(s)
	if input == "" {
		return nil, &httperror.ValidationError{
			Code:   ErrCodeInvalidURL,
			Detail: "Invalid HTTP URL. The input string is empty",
		}
	}

	if !hasHTTPSchemePrefix(input) {
		return nil, &httperror.ValidationError{
			Code:   ErrCodeInvalidURIScheme,
			Detail: "Invalid HTTP URL scheme",
		}
	}

	parsedURL, err := parseNormalizedURL(input)
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
	if err != nil && err.Code == ErrCodeInvalidURI {
		err.Code = ErrCodeInvalidURL
	}

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
// Make sure that the options field is not null.
// Otherwise, the validation is skipped.
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
		return &httperror.ValidationError{
			Code:   ErrCodeInvalidHostname,
			Detail: "Invalid URL. Hostname is empty",
		}
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

	uri, _, err := parseURIAndHostname(input)

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
			return &httperror.ValidationError{
				Code:   ErrCodeInvalidHostname,
				Detail: "Hostname is blocked",
			}
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

	return &httperror.ValidationError{
		Code:   ErrCodeInvalidHostname,
		Detail: "Hostname is not allowed",
	}
}

func validateURLScheme(uri *url.URL, allowedSchemes []string) error {
	if len(allowedSchemes) > 0 && !slices.ContainsFunc(allowedSchemes, func(item string) bool {
		return strings.EqualFold(item, uri.Scheme)
	}) {
		return &httperror.ValidationError{
			Code: ErrCodeInvalidURIScheme,
			Detail: "Invalid URI scheme. Accept one of [" + strings.Join(allowedSchemes, ", ") +
				"], got " + strconv.Quote(uri.Scheme),
		}
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
			Code:   ErrCodeInvalidHostname,
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

func parseAndValidateURL(s string) (*url.URL, *httperror.ValidationError) {
	input := strings.TrimSpace(s)
	if input == "" {
		return nil, &httperror.ValidationError{
			Code:   ErrCodeInvalidURL,
			Detail: "Invalid URL. The input string is empty",
		}
	}

	schemeIndex := strings.Index(input, "://")
	if schemeIndex <= 0 || len(input)-schemeIndex <= 1 {
		return nil, &httperror.ValidationError{
			Code:   ErrCodeInvalidURL,
			Detail: "Invalid URL syntax",
		}
	}

	return parseNormalizedURL(input)
}

func parseNormalizedURL(input string) (*url.URL, *httperror.ValidationError) {
	parsedURI, hostname, err := parseURIAndHostname(input)
	if err != nil {
		return nil, err
	}

	if hostname == "" {
		return nil, &httperror.ValidationError{
			Code:   ErrCodeInvalidHostname,
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

func hasHTTPSchemePrefix(input string) bool {
	if len(input) < 7 || !strings.EqualFold(input[:4], "http") {
		return false
	}

	switch input[4] {
	case 's', 'S':
		return len(input) >= 8 && input[5:8] == "://"
	case ':':
		return input[5:7] == "//"
	default:
		return false
	}
}

func isHTTPScheme(scheme string) bool {
	switch len(scheme) {
	case 4:
		return strings.EqualFold(scheme, "http")
	case 5:
		return strings.EqualFold(scheme, "https")
	default:
		return false
	}
}
