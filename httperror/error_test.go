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

package httperror

import (
	"testing"
)

func TestValidationError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  ValidationError
		want string
	}{
		{
			name: "detail only",
			err:  ValidationError{Detail: "field is required"},
			want: "detail: field is required",
		},
		{
			name: "detail and pointer",
			err:  ValidationError{Detail: "must be positive", Pointer: "/count"},
			want: "detail: must be positive\npointer: /count",
		},
		{
			name: "detail and location",
			err:  ValidationError{Detail: "missing value", Location: "query"},
			want: "detail: missing value\nlocation: query",
		},
		{
			name: "detail and code",
			err:  ValidationError{Detail: "invalid format", Code: "ERR001"},
			want: "detail: invalid format\ncode: ERR001",
		},
		{
			name: "detail and hint",
			err:  ValidationError{Detail: "too short", Hint: "minimum length is 8"},
			want: "detail: too short\nhint: minimum length is 8",
		},
		{
			name: "all fields",
			err: ValidationError{
				Detail:   "invalid value",
				Pointer:  "/name",
				Location: "body",
				Code:     "ERR042",
				Hint:     "use alphanumeric characters",
			},
			want: "detail: invalid value\ncode: ERR042\npointer: /name\nlocation: body\nhint: use alphanumeric characters",
		},
		{
			name: "empty struct",
			err:  ValidationError{},
			want: "",
		},
		{
			name: "pointer and location without detail",
			err:  ValidationError{Pointer: "/id", Location: "path"},
			want: "pointer: /id\nlocation: path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestHTTPError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  HTTPError
		want string
	}{
		{
			name: "title only",
			err:  HTTPError{Title: "Not Found"},
			want: "title: Not Found",
		},
		{
			name: "status only",
			err:  HTTPError{Status: 404},
			want: "status: 404",
		},
		{
			name: "title and status",
			err:  HTTPError{Title: "Not Found", Status: 404},
			want: "title: Not Found\nstatus: 404",
		},
		{
			name: "title, status, and detail",
			err:  HTTPError{Title: "Bad Request", Status: 400, Detail: "The request is invalid."},
			want: "title: Bad Request\nstatus: 400\ndetail: The request is invalid.",
		},
		{
			name: "title, type, status, and detail",
			err: HTTPError{
				Title:  "Not Found",
				Type:   "https://example.com/not-found",
				Status: 404,
				Detail: "Resource not found.",
			},
			want: "title: Not Found\ntype: https://example.com/not-found\nstatus: 404\ndetail: Resource not found.",
		},
		{
			name: "with instance",
			err: HTTPError{
				Title:    "Not Found",
				Status:   404,
				Detail:   "Resource not found.",
				Instance: "/resources/123",
			},
			want: "title: Not Found\nstatus: 404\ninstance: /resources/123\ndetail: Resource not found.",
		},
		{
			name: "with single validation error",
			err: HTTPError{
				Title:  "Bad Request",
				Status: 400,
				Detail: "The request is invalid.",
				Errors: []ValidationError{
					{Detail: "field is required", Pointer: "/name"},
				},
			},
			want: "title: Bad Request\nstatus: 400\ndetail: The request is invalid.\nerrors:\n  - detail: field is required\n    pointer: /name",
		},
		{
			name: "with multiple validation errors",
			err: HTTPError{
				Title:  "Bad Request",
				Status: 400,
				Detail: "The request is invalid.",
				Errors: []ValidationError{
					{Detail: "field is required", Pointer: "/name"},
					{Detail: "must be positive", Pointer: "/count", Code: "ERR001"},
				},
			},
			want: "title: Bad Request\nstatus: 400\ndetail: The request is invalid.\nerrors:\n  - detail: field is required\n    pointer: /name\n  - detail: must be positive\n    code: ERR001\n    pointer: /count",
		},
		{
			name: "empty struct",
			err:  HTTPError{},
			want: "",
		},
		{
			name: "validation error with all subfields",
			err: HTTPError{
				Title:  "Unprocessable Entity",
				Status: 422,
				Detail: "Validation failed.",
				Errors: []ValidationError{
					{
						Detail:   "invalid email",
						Pointer:  "/email",
						Location: "body",
						Code:     "ERR010",
						Hint:     "use a valid email address",
					},
				},
			},
			want: "title: Unprocessable Entity\nstatus: 422\ndetail: Validation failed.\nerrors:\n  - detail: invalid email\n    code: ERR010\n    pointer: /email\n    location: body\n    hint: use a valid email address",
		},
		{
			name: "zero status is omitted",
			err:  HTTPError{Title: "Custom Error", Detail: "something went wrong"},
			want: "title: Custom Error\ndetail: something went wrong",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.want {
				t.Errorf("Error() =\n%q\nwant:\n%q", got, tt.want)
			}
		})
	}
}
