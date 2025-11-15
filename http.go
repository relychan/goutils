package goutils

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// ParseRelativeOrHTTPURL validates and parses relative or HTTP URL.
func ParseRelativeOrHTTPURL(input string) (*url.URL, error) {
	if strings.HasPrefix(input, "/") || !strings.HasPrefix(input, "http") {
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

	return url.Parse(input)
}

// ParseHTTPURL parses and validate the input string to have http(s) scheme.
func ParseHTTPURL(input string) (*url.URL, error) {
	u, err := url.Parse(input)
	if err != nil {
		return nil, err
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("%w, got: %s", ErrInvalidHTTPScheme, u.Scheme)
	}

	return u, nil
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
