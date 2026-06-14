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
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/relychan/goutils/httperror"
)

// ParseFilePathOrHTTPURL validates and parses a file path or HTTP URL.
func ParseFilePathOrHTTPURL(input string) (*url.URL, error) {
	parsedURL, err := ParseFilePathOrURL(input)
	if err != nil {
		return nil, err
	}

	// Returns if the parsedURL is a path.
	if parsedURL == nil {
		return nil, nil
	}

	if !isHTTPScheme(parsedURL.Scheme) {
		return nil, &httperror.ValidationError{
			Code:   ErrCodeInvalidURIScheme,
			Detail: "Invalid HTTP scheme. Expected http(s), got " + strconv.Quote(parsedURL.Scheme),
		}
	}

	return parsedURL, nil
}

// ParseFilePathOrURL validates and parses a file path or URL.
// If the input string is a URL, returns the parsed URL.
// Otherwise returns null.
func ParseFilePathOrURL(input string) (*url.URL, error) {
	input = strings.TrimSpace(input)
	if input == "" || input == "." ||
		(filepath.Separator == '/' && input == string(filepath.Separator)) {
		return nil, nil
	}

	urlSepIndex := strings.IndexAny(input, ":?#")
	if urlSepIndex != -1 {
		if urlSepIndex == 1 && input[urlSepIndex] == ':' &&
			IsAlphabet(input[0]) &&
			(len(input) == 2 || input[2] == '\\') {
			// It is likely a Windows style's path
			return nil, validateFilePath(input)
		}

		return ParseAbsoluteURL(input)
	}

	return nil, validateFilePath(input)
}

// SplitPathQueryFragment splits path, query and fragment from string.
func SplitPathQueryFragment(input string) (string, string, string) {
	if input == "" {
		return "", "", ""
	}

	u, fragment, _ := strings.Cut(input, "#")
	uriPath, query, _ := strings.Cut(u, "?")

	return uriPath, query, fragment
}

// ParseAbsoluteURI parses and validate the input string to be a valid absolute URI.
//
//	scheme:[//authority]/path[?query][#fragment]
func ParseAbsoluteURI(input string) (*url.URL, error) {
	parsedURI, _, err := parseURIAndHostname(input)
	if err != nil {
		return nil, err
	}

	return parsedURI, nil
}

// ParseAbsoluteURL parses and validate the input string to be a valid absolute URL.
// Unlike URI, the URL requires an explicit authority.
//
//	scheme://authority[/path][?query][#fragment]
func ParseAbsoluteURL(input string) (*url.URL, error) {
	result, err := parseAndValidateURL(input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ParseURL parses and validate the input string to be a valid URL.
// If the URL is relative, it must start with a slash.
// Use [ParseAbsoluteURL] if you expect strict absolute URLs.
func ParseURL(input string) (*url.URL, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return new(url.URL), nil
	}

	if isRelativeURL(input) {
		return ParseAbsoluteURL(input)
	}

	result := &url.URL{}

	err := AppendURL(result, input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ParseHTTPURL parses and validate the input string to be a valid HTTP URL.
// If the URL is relative, it must start with a slash.
// Use [ParseAbsoluteHTTPURL] if you expect strict absolute URLs.
func ParseHTTPURL(input string) (*url.URL, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return new(url.URL), nil
	}

	if isRelativeURL(input) {
		return ParseAbsoluteHTTPURL(input)
	}

	result := &url.URL{}

	err := AppendURL(result, input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ParseAbsoluteHTTPURL parses and validate the input string to be a valid absolute URL and have http(s) scheme.
func ParseAbsoluteHTTPURL(s string) (*url.URL, error) {
	input := strings.TrimSpace(s)
	if input == "" {
		return nil, &httperror.ValidationError{
			Code:   ErrCodeInvalidURI,
			Detail: "Invalid HTTP URL. The input string is empty",
		}
	}

	if !IsURLSchemePrefixHTTP(input) {
		return nil, &httperror.ValidationError{
			Code:   ErrCodeInvalidURI,
			Detail: "Invalid HTTP URL scheme",
		}
	}

	parsedURL, err := parseNormalizedURL(input)
	if err != nil {
		return nil, err
	}

	return parsedURL, nil
}

// ParseRelativeURI parses and validate the input string to be a valid relative URI.
func ParseRelativeURI(input string) (*url.URL, error) {
	input = strings.TrimSpace(input)
	if input == "" || input == "/" {
		return &url.URL{
			Path: "/",
		}, nil
	}

	if !isRelativeURL(input) {
		return nil, &httperror.ValidationError{
			Code:   ErrCodeInvalidURI,
			Detail: "Invalid relative URL syntax. The input string must start with a slash or question mark.",
		}
	}

	result := &url.URL{}

	err := AppendURL(result, input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ParseAndValidateURLWithOptions parses and validates URL from a string with options.
// Returns the parsed URL and an error.
func ParseAndValidateURLWithOptions(
	ctx context.Context,
	urlStr string,
	options *ValidateHTTPURLOptions,
) (*url.URL, error) {
	parsedURL, err := ParseAbsoluteURL(urlStr)
	if err != nil {
		return nil, err
	}

	validatedErr := ValidateURLWithOptions(ctx, parsedURL, options)
	if validatedErr != nil {
		return nil, validatedErr
	}

	return parsedURL, nil
}

// ValidateAbsoluteURI checks if the input string is a valid absolute URI.
func ValidateAbsoluteURI(s string) *httperror.ValidationError {
	_, err := parseAndValidateURI(s)

	return err
}

// ValidateAbsoluteURL checks if the input string is a valid absolute URL.
func ValidateAbsoluteURL(s string) *httperror.ValidationError {
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

// IsURLSchemePrefixHTTP reports whether input begins with "http://" or "https://" (case-insensitive for the scheme).
func IsURLSchemePrefixHTTP(input string) bool {
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

// AppendURL appends uriPath's path, query, and fragment components to uri in-place.
// uriPath may contain a path with optional ?query and #fragment.
func AppendURL(uri *url.URL, uriPath string) error { //nolint:cyclop
	uriPath = strings.TrimSpace(uriPath)
	if uriPath == "" || uriPath == "/" {
		return nil
	}

	path, query, fragment := SplitPathQueryFragment(uriPath)

	if fragment != "" {
		if StringContainsCTLByte(fragment) {
			return &httperror.ValidationError{
				Code:   ErrCodeInvalidURI,
				Detail: "Invalid URL fragment syntax",
			}
		}

		uri.Fragment = fragment
	}

	if query != "" {
		if StringContainsCTLByte(query) {
			return &httperror.ValidationError{
				Code:   ErrCodeInvalidURI,
				Detail: "Invalid URL query syntax",
			}
		}

		switch {
		case uri.RawQuery == "":
			uri.RawQuery = query
		case strings.HasSuffix(uri.RawQuery, "&") || strings.HasPrefix(query, "&"):
			uri.RawQuery += query
		default:
			uri.RawQuery += "&" + query
		}
	}

	if path != "" && path != "/" {
		err := ValidateURLPath(path)
		if err != nil {
			return err
		}

		uri.Path = strings.TrimRight(uri.Path, "/")

		switch {
		case uri.Path == "" || uri.Path == "/":
			uri.Path = path
		case path[0] == '/':
			uri.Path += path
		default:
			uri.Path += "/" + path
		}

		// Keep url.URL.Path normalized for absolute URLs.
		if uri.Host != "" && uri.Path != "" && uri.Path[0] != '/' {
			uri.Path = "/" + uri.Path
		}
	}

	return nil
}

// ValidateURLPath checks if the URL path is valid.
func ValidateURLPath(input string) *httperror.ValidationError {
	if input == "" || input == "/" {
		return nil
	}

	if strings.Contains(input, "://") {
		return &httperror.ValidationError{
			Code:   ErrCodeInvalidPath,
			Detail: "URL path must not be an absolute URL",
		}
	}

	// validate invalid path patterns
	if input[0] == '/' {
		input = input[1:]
	}

	for input != "" {
		slashIndex := strings.IndexByte(input, '/')
		if slashIndex == 0 {
			return &httperror.ValidationError{
				Code:   ErrCodeInvalidPath,
				Detail: "Invalid double slashes in the URL path syntax",
			}
		}

		part := input

		if slashIndex != -1 {
			part = input[:slashIndex]
			input = input[slashIndex+1:]
		} else {
			input = ""
		}

		if part == "*" || StringAllRune(part, '.') {
			return &httperror.ValidationError{
				Code:   ErrCodeInvalidPath,
				Detail: "Wildcard and traversal paths are not allowed in URL path",
			}
		}

		if StringContainsCTLByte(part) {
			return &httperror.ValidationError{
				Code:   ErrCodeInvalidPath,
				Detail: "URL path contains invalid characters",
			}
		}
	}

	return nil
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

	if parsedURI.Scheme == "" {
		return nil, "", &httperror.ValidationError{
			Code:   ErrCodeInvalidURI,
			Detail: "URI Scheme is empty",
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

	err = ValidateURLPath(parsedURI.Path)
	if err != nil {
		return nil, err
	}

	if parsedURI.Path != "" && parsedURI.Path[0] != '/' {
		parsedURI.Path = "/" + parsedURI.Path
	}

	return parsedURI, nil
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

func isRelativeURL(input string) bool {
	return input[0] != '/' && input[0] != '?' && input[0] != '#'
}
