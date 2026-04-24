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
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/relychan/goutils/httperror"
)

// Doer abstracts an interface for sending HTTP requests.
type Doer interface {
	// Do sends an HTTP request and returns an HTTP response, following policy
	// (such as redirects, cookies, auth) as configured on the client.
	Do(req *http.Request) (*http.Response, error)
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

// CloseResponse gracefully closes the HTTP response and tries to drain the body if it exists.
// It makes a best effort to reuse the HTTP connection.
func CloseResponse(resp *http.Response) {
	if resp == nil || resp.Body == nil || resp.Body == http.NoBody {
		return
	}

	contentLength := resp.ContentLength
	if contentLength == -1 {
		rawContentLength := resp.Header["Content-Length"]
		if len(rawContentLength) > 0 {
			intContentLength, err := strconv.ParseInt(rawContentLength[0], 10, 64)
			if err == nil {
				contentLength = intContentLength
			}
		}
	}

	if contentLength == -1 || contentLength <= httperror.MaxPostCloseReadBytes {
		_, _ = io.CopyN(io.Discard, resp.Body, httperror.MaxPostCloseReadBytes+1) //nolint:errcheck
	}

	CatchWarnErrorFunc(resp.Body.Close)
}
