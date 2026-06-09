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
	"net/netip"
	"net/url"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/relychan/goutils/httperror"
)

var allDurationUnits = []string{"YMWD", "HMS"}

const (
	// ErrCodeInvalidJSONPointer represents an error code for invalid JSON pointer.
	ErrCodeInvalidJSONPointer = "invalid_json_pointer"
	// ErrCodeInvalidUUID represents an error code for invalid UUID string.
	ErrCodeInvalidUUID = "invalid_uuid"
	// ErrCodeInvalidDuration represents an error code for invalid duration.
	ErrCodeInvalidDuration = "invalid_duration"
	// ErrCodeInvalidDate represents an error code for invalid date.
	ErrCodeInvalidDate = "invalid_date"
	// ErrCodeInvalidTime represents an error code for invalid time.
	ErrCodeInvalidTime = "invalid_time"
	// ErrCodeInvalidDateTime represents an error code for invalid date time.
	ErrCodeInvalidDateTime = "invalid_date_time"
	// ErrCodeInvalidIPv4 represents an error code for invalid IPv4 string.
	ErrCodeInvalidIPv4 = "invalid_ipv4"
	// ErrCodeInvalidIPv6 represents an error code for invalid IPv6 string.
	ErrCodeInvalidIPv6 = "invalid_ipv6"
	// ErrCodeInvalidHostname represents an error code for invalid hostname.
	ErrCodeInvalidHostname = "invalid_hostname"
	// ErrCodeInvalidEmail represents an error code for invalid email.
	ErrCodeInvalidEmail = "invalid_email"
	// ErrCodeInvalidURI represents an error code for invalid URI.
	ErrCodeInvalidURI = "invalid_uri"
)

// ValidateJSONPointer validates the JSON pointer string according to the [RFC 6901] specification.
//
// [RFC 6901]: https://www.rfc-editor.org/rfc/rfc6901#section-3
func ValidateJSONPointer(s string) *httperror.ValidationError {
	if s == "" {
		return nil
	}

	if s[0] != '/' {
		return &httperror.ValidationError{
			Code:   ErrCodeInvalidJSONPointer,
			Detail: "Invalid JSON pointer; the value is not starting with /",
		}
	}

	for tok := range strings.SplitSeq(s[1:], "/") {
		err := validateJSONPointerToken(tok)
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateUUID validates if the string is a valid UUID.
func ValidateUUID(s string) *httperror.ValidationError {
	err := uuid.Validate(s)
	if err != nil {
		return &httperror.ValidationError{
			Code:   ErrCodeInvalidUUID,
			Detail: err.Error(),
		}
	}

	return nil
}

// ValidateDurationRFC3339 validates an RFC 3339 duration in format:
//
//	P[Y]Y[M]M[W]W[D]DT[H]H[M]M[S]S
//
// with:
//
//	P: Duration designator.
//	T: Time designator (separates the date parts from the time parts).
//
// see https://datatracker.ietf.org/doc/html/rfc3339#appendix-A
func ValidateDurationRFC3339( //nolint:gocognit,funlen
	value string,
) *httperror.ValidationError {
	// must start with 'P'
	if value == "" || value[0] != 'P' {
		return &httperror.ValidationError{
			Code:   ErrCodeInvalidDuration,
			Detail: "Invalid duration; the string must start with P",
		}
	}

	if len(value) == 1 {
		return &httperror.ValidationError{
			Code:   ErrCodeInvalidDuration,
			Detail: "Invalid duration; nothing after P",
		}
	}

	for i, s := range strings.Split(value[1:], "T") {
		if i != 0 && s == "" {
			return &httperror.ValidationError{
				Code:   ErrCodeInvalidDuration,
				Detail: "Invalid duration; no time elements",
			}
		}

		if i >= len(allDurationUnits) {
			return &httperror.ValidationError{
				Code:   ErrCodeInvalidDuration,
				Detail: "Invalid duration; more than one T",
			}
		}

		units := allDurationUnits[i]

		for s != "" {
			digitCount := 0

			for _, ch := range s {
				if !IsDigit(ch) {
					break
				}

				digitCount++
			}

			if digitCount == 0 {
				return &httperror.ValidationError{
					Code:   ErrCodeInvalidDuration,
					Detail: "Invalid duration; missing number",
				}
			}

			s = s[digitCount:]
			if s == "" {
				return &httperror.ValidationError{
					Code:   ErrCodeInvalidDuration,
					Detail: "Invalid duration; missing unit",
				}
			}

			unit := s[0]

			j := strings.IndexByte(units, unit)
			if j == -1 {
				if strings.IndexByte(allDurationUnits[i], unit) != -1 {
					return &httperror.ValidationError{
						Code:   ErrCodeInvalidDuration,
						Detail: fmt.Sprintf("unit %q out of order", unit),
					}
				}

				return &httperror.ValidationError{
					Code:   ErrCodeInvalidDuration,
					Detail: "invalid unit " + string(unit),
				}
			}

			units = units[j+1:]
			s = s[1:]
		}
	}

	return nil
}

// ValidateIPV4 validates if the string satisfies the IPv4 format.
func ValidateIPV4(s string) *httperror.ValidationError {
	ip := net.ParseIP(s)
	if ip == nil || ip.To4() == nil || !strings.ContainsRune(s, '.') {
		return &httperror.ValidationError{
			Code:   ErrCodeInvalidIPv4,
			Detail: "Invalid IPv4",
		}
	}

	return nil
}

// ValidateIPV6 validates if the string satisfies the IPv6 format.
func ValidateIPV6(s string) *httperror.ValidationError {
	if !strings.ContainsRune(s, ':') {
		return &httperror.ValidationError{
			Code:   ErrCodeInvalidIPv6,
			Detail: "Invalid IPv6; missing colon",
		}
	}

	addr, err := netip.ParseAddr(s)
	if err != nil {
		return &httperror.ValidationError{
			Code:   ErrCodeInvalidIPv6,
			Detail: "Invalid IPv6; " + err.Error(),
		}
	}

	if addr.Zone() != "" {
		return &httperror.ValidationError{
			Code:   ErrCodeInvalidIPv6,
			Detail: "Invalid IPv6; zone id is not a part of ipv6 address",
		}
	}

	return nil
}

// ValidateHostname validates if the string is a [valid hostname].
//
// [valid hostname]: https://en.wikipedia.org/wiki/Hostname#Restrictions_on_valid_host_names
func ValidateHostname(s string) *httperror.ValidationError {
	// entire hostname (including the delimiting dots but not a trailing dot) has a maximum of 253 ASCII characters
	s = strings.TrimSuffix(s, ".")
	if len(s) > 253 {
		return &httperror.ValidationError{
			Code:   ErrCodeInvalidHostname,
			Detail: "Hostname must not be more than 253 characters long",
		}
	}

	// Hostnames are composed of series of labels concatenated with dots, as are all domain names
	for label := range strings.SplitSeq(s, ".") {
		// Each label must be from 1 to 63 characters long
		if len(label) < 1 || len(label) > 63 {
			return &httperror.ValidationError{
				Code:   ErrCodeInvalidHostname,
				Detail: "Invalid hostname; label must be 1 to 63 characters long",
			}
		}

		// labels must not start or end with a hyphen
		if strings.HasPrefix(label, "-") {
			return &httperror.ValidationError{
				Code:   ErrCodeInvalidHostname,
				Detail: "Invalid hostname; label must not start with hyphen",
			}
		}

		if strings.HasSuffix(label, "-") {
			return &httperror.ValidationError{
				Code:   ErrCodeInvalidHostname,
				Detail: "Invalid hostname; label must not end with hyphen",
			}
		}

		// labels may contain only the ASCII letters 'a' through 'z' (in a case-insensitive manner),
		// the digits '0' through '9', and the hyphen ('-')
		for _, c := range label {
			if c != '-' && !IsDigit(c) && !IsLowerAlphabet(c) && !IsUpperAlphabet(c) {
				return &httperror.ValidationError{
					Code:   ErrCodeInvalidHostname,
					Detail: "Invalid hostname; invalid character " + string(c),
				}
			}
		}
	}

	return nil
}

// ValidateEmail validates if the string is a valid [email address].
//
// [email address]: https://en.wikipedia.org/wiki/Email_address
func ValidateEmail(s string) *httperror.ValidationError { //nolint:cyclop,funlen
	// entire email address to be no more than 254 characters long
	if len(s) > 254 {
		return &httperror.ValidationError{
			Code:   ErrCodeInvalidEmail,
			Detail: "Email must not be more than 254 characters long",
		}
	}

	// email address is generally recognized as having two parts joined with an at-sign
	at := strings.LastIndexByte(s, '@')
	if at == -1 {
		return &httperror.ValidationError{
			Code:   ErrCodeInvalidEmail,
			Detail: "Invalid email; @ character must exist",
		}
	}

	local, domain := s[:at], s[at+1:]

	// local part may be up to 64 characters long
	if len(local) > 64 {
		return &httperror.ValidationError{
			Code:   ErrCodeInvalidEmail,
			Detail: "Invalid email; local part must not be more than 64 characters long",
		}
	}

	if len(local) > 1 && local[0] == '"' && local[len(local)-1] == '"' { //nolint:nestif
		// quoted
		local := local[1 : len(local)-1]
		if strings.IndexByte(local, '\\') != -1 || strings.IndexByte(local, '"') != -1 {
			return &httperror.ValidationError{
				Code:   ErrCodeInvalidEmail,
				Detail: "Invalid email; backslash and quote are not allowed within quoted local part",
			}
		}
	} else {
		// unquoted
		if strings.HasPrefix(local, ".") {
			return &httperror.ValidationError{
				Code:   ErrCodeInvalidEmail,
				Detail: "Email must not start with dot",
			}
		}

		if strings.HasSuffix(local, ".") {
			return &httperror.ValidationError{
				Code:   ErrCodeInvalidEmail,
				Detail: "Email must not end with dot",
			}
		}

		// consecutive dots not allowed
		if strings.Contains(local, "..") {
			return &httperror.ValidationError{
				Code:   ErrCodeInvalidEmail,
				Detail: "Email must not contain consecutive dots",
			}
		}

		// check allowed chars
		for _, c := range local {
			if !IsDigit(c) && !IsLowerAlphabet(c) && !IsUpperAlphabet(c) &&
				!strings.ContainsRune(".!#$%&'*+-/=?^_`{|}~", c) {
				return &httperror.ValidationError{
					Code:   ErrCodeInvalidEmail,
					Detail: "Invalid email; invalid character " + string(c),
				}
			}
		}
	}

	// domain if enclosed in brackets, must match an IP address
	if strings.HasPrefix(domain, "[") && strings.HasSuffix(domain, "]") {
		domain = domain[1 : len(domain)-1]

		rem, ok := strings.CutPrefix(domain, "IPv6:")
		if ok {
			err := ValidateIPV6(rem)
			if err != nil {
				err.Code = ErrCodeInvalidEmail
				err.Detail = "Invalid email address: " + err.Detail

				return err
			}

			return nil
		}

		err := ValidateIPV4(domain)
		if err != nil {
			err.Code = ErrCodeInvalidEmail
			err.Detail = "Invalid email address: " + err.Detail

			return err
		}

		return nil
	}

	// domain must match the requirements for a hostname
	err := ValidateHostname(domain)
	if err != nil {
	err.Code = ErrCodeInvalidEmail
	err.Detail = "Invalid email address: " + err.Detail

		return err
	}

	return nil
}

// ValidateDate validates if the string is a valid date.
func ValidateDate(s string) *httperror.ValidationError {
	_, _, _, err := parseDateString(s) //nolint:dogsled
	if err != nil {
		return &httperror.ValidationError{
			Code:   ErrCodeInvalidDate,
			Detail: err.Error(),
		}
	}

	return nil
}

// ValidateTime validates time with [RFC 3339] specification.
//
// [RFC 3339]: https://datatracker.ietf.org/doc/html/rfc3339#section-5.6
func ValidateTime(str string) *httperror.ValidationError {
	_, _, _, _, _, err := parseTimeString(str) //nolint:dogsled
	if err != nil {
		return &httperror.ValidationError{
			Code:   ErrCodeInvalidTime,
			Detail: err.Error(),
		}
	}

	return nil
}

// ValidateDateTime validates date time with [RFC 3339] specification.
//
// [RFC 3339]: https://datatracker.ietf.org/doc/html/rfc3339#section-5.6
func ValidateDateTime(input string) *httperror.ValidationError {
	if len(input) < dateTimeLength {
		return &httperror.ValidationError{
			Code:   ErrCodeInvalidDateTime,
			Detail: ErrInvalidDateTimeString.Error(),
		}
	}

	_, _, _, err := parseDateString(input[:dateLength]) //nolint:dogsled
	if err != nil {
		return &httperror.ValidationError{
			Code:   ErrCodeInvalidDateTime,
			Detail: err.Error(),
		}
	}

	if input[10] != 't' && input[10] != 'T' && input[10] != ' ' {
		return &httperror.ValidationError{
			Code:   ErrCodeInvalidDateTime,
			Detail: "Invalid date time format. Missing T between date and time",
		}
	}

	_, _, _, _, _, err = parseTimeString(input[11:]) //nolint:dogsled
	if err != nil {
		return &httperror.ValidationError{
			Code:   ErrCodeInvalidDateTime,
			Detail: err.Error(),
		}
	}

	return nil
}

// ValidateURI checks if the input string is a valid URI.
func ValidateURI(s string) *httperror.ValidationError {
	u, err := ParseURL(s)
	if err != nil {
		return &httperror.ValidationError{
			Code:   ErrCodeInvalidURI,
			Detail: err.Error(),
		}
	}

	if !u.IsAbs() {
		return &httperror.ValidationError{
			Code:   ErrCodeInvalidURI,
			Detail: "Relative URL is not allowed",
		}
	}

	return nil
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
	options ValidateHTTPURLOptions,
) error {
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

func validateHost(host, hostname string, options *ValidateHTTPURLOptions) error {
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

func validateJSONPointerToken(tok string) *httperror.ValidationError {
	escape := false

	for _, ch := range tok {
		if escape {
			escape = false

			if ch != '0' && ch != '1' {
				return &httperror.ValidationError{
					Code:   ErrCodeInvalidJSONPointer,
					Detail: "Invalid JSON pointer; ~ must be followed by 0 or 1",
				}
			}

			continue
		}

		if ch == '~' {
			escape = true

			continue
		}

		switch {
		case ch >= '\x00' && ch <= '\x2E':
		case ch >= '\x30' && ch <= '\x7D':
		case ch >= '\x7F' && ch <= '\U0010FFFF':
		default:
			return &httperror.ValidationError{
				Detail: "Invalid JSON pointer; invalid character " + string(ch),
			}
		}
	}

	if escape {
		return &httperror.ValidationError{
			Detail: "Invalid JSON pointer; ~ must be followed by 0 or 1",
		}
	}

	return nil
}
