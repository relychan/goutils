package goutils

import (
	"fmt"
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
