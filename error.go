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
	"encoding/json"
	"errors"
	"log/slog"
	"maps"

	"github.com/relychan/goutils/httperror"
)

var (
	// ErrInvalidURI represents an invalid uri error.
	ErrInvalidURI = errors.New("invalid URI")
	// ErrBlockedIP occurs when the IP is blocked.
	ErrBlockedIP = errors.New("ip is blocked")
	// ErrInvalidURLScheme represents an invalid URL scheme error.
	ErrInvalidURLScheme = errors.New("invalid url scheme")
	// ErrInvalidSubnet occurs when the subnet string is invalid.
	ErrInvalidSubnet = errors.New("invalid IP or subnet")
	// ErrInvalidSlug represents an invalid slug error.
	ErrInvalidSlug = errors.New("invalid slug")
	// ErrMalformedJSON occurs when the JSON syntax or value is malformed.
	ErrMalformedJSON = errors.New("malformed JSON")
	// ErrStringNull occurs when the string value is nil.
	ErrStringNull = errors.New("string value must not be null")
	// ErrMalformedString occurs when the string value is malformed.
	ErrMalformedString = errors.New("malformed string")
	// ErrMalformedStringSlice occurs when the string slice is malformed.
	ErrMalformedStringSlice = errors.New("malformed string slice")
	// ErrStringSliceNull occurs when the string slice is nil.
	ErrStringSliceNull = errors.New("string slice must not be null")
	// ErrNumberNull occurs when the number value is nil.
	ErrNumberNull = errors.New("number value must not be null")
	// ErrMalformedNumber occurs when the number value is malformed.
	ErrMalformedNumber = errors.New("malformed number")
	// ErrMalformedNumberSlice occurs when the number slice is malformed.
	ErrMalformedNumberSlice = errors.New("malformed number slice")
	// ErrBooleanNull occurs when the boolean value is nil.
	ErrBooleanNull = errors.New("boolean value must not be null")
	// ErrMalformedBoolean occurs when the boolean value is malformed.
	ErrMalformedBoolean = errors.New("malformed boolean")
	// ErrMalformedBooleanSlice occurs when the boolean slice is malformed.
	ErrMalformedBooleanSlice = errors.New("malformed boolean slice")
	// ErrBooleanSliceNull occurs when the boolean slice is nil.
	ErrBooleanSliceNull = errors.New("boolean slice must not be null")
	// ErrMalformedYAML occurs when the YAML syntax or structure is malformed.
	ErrMalformedYAML = errors.New("malformed YAML")
)

// CatchWarnErrorFunc catches the closer function and prints error with the WARN level.
func CatchWarnErrorFunc(fn func() error) {
	err := fn()
	if err != nil {
		slog.Warn(err.Error())
	}
}

// CatchWarnContextErrorFunc catches the closer function with context and prints error with the WARN level.
func CatchWarnContextErrorFunc(fn func(ctx context.Context) error) {
	err := fn(context.TODO())
	if err != nil {
		slog.Warn(err.Error())
	}
}

// HTTPErrorWithExtensions is the data structure of an HTTP error with extensions.
// that follows the [RFC 9457] specification.
// The schema is inspired by [Swagger API] specification.
//
// [RFC 9457]: https://www.rfc-editor.org/rfc/rfc9457.html
// [Swagger API]: https://swagger.io/blog/problem-details-rfc9457-api-error-handling/
type HTTPErrorWithExtensions struct { //nolint:errname
	httperror.HTTPError

	// Additional members that are specific to that problem type.
	Extensions map[string]any
}

// NewHTTPErrorWithExtensions creates a new RFC 9457 error with extensions.
func NewHTTPErrorWithExtensions(
	err httperror.HTTPError,
	extensions map[string]any,
) *HTTPErrorWithExtensions {
	return &HTTPErrorWithExtensions{
		HTTPError:  err,
		Extensions: extensions,
	}
}

// Error implements the error interface for HTTPErrorWithExtensions.
func (e HTTPErrorWithExtensions) Error() string {
	sb := httperror.NewHTTPErrorStringBuilder(e.HTTPError)

	if len(e.Extensions) > 0 {
		if sb.Len() > 0 {
			sb.WriteByte('\n')
		}

		buildMapToString(sb, e.Extensions, 0)
	}

	return sb.String()
}

// MarshalJSON implements the json.Marshaler interface.
func (e HTTPErrorWithExtensions) MarshalJSON() ([]byte, error) {
	if len(e.Extensions) == 0 {
		return json.Marshal(e.HTTPError)
	}

	result := maps.Clone(e.Extensions)

	if e.Type != "" {
		result["type"] = e.Type
	}

	if e.Status > 0 {
		result["status"] = e.Status
	}

	if e.Title != "" {
		result["title"] = e.Title
	}

	if e.Detail != "" {
		result["detail"] = e.Detail
	}

	if e.Instance != "" {
		result["instance"] = e.Instance
	}

	if e.Code != "" {
		result["code"] = e.Code
	}

	if e.Errors != nil {
		result["errors"] = e.Errors
	}

	return json.Marshal(result)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (e *HTTPErrorWithExtensions) UnmarshalJSON(
	data []byte,
) error {
	var baseError httperror.HTTPError

	err := json.Unmarshal(data, &baseError)
	if err != nil {
		return err
	}

	e.HTTPError = baseError

	extensions := map[string]any{}

	err = json.Unmarshal(data, &extensions)
	if err != nil {
		return err
	}

	for _, key := range []string{"type", "status", "title", "detail", "instance", "code", "errors"} {
		delete(extensions, key)
	}

	if e.Detail == "" {
		message, ok := extensions["message"]
		if ok && message != nil {
			msg, ok := message.(string)
			if ok {
				e.Detail = msg
			}

			delete(extensions, "message")
		}
	}

	e.Extensions = extensions

	return nil
}
