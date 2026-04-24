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
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/relychan/goutils/httperror"
)

func TestHTTPError(t *testing.T) {
	testCases := []struct {
		Error *HTTPErrorWithExtensions
	}{
		{
			Error: NewHTTPErrorWithExtensions(*httperror.NewAlreadyExistsError(), nil),
		},
		{
			Error: NewHTTPErrorWithExtensions(*httperror.NewForbiddenError(httperror.ValidationError{
				Detail: "forbidden",
			}), nil),
		},
		{
			Error: NewHTTPErrorWithExtensions(*httperror.NewBadRequestError(), nil),
		},
		{
			Error: NewHTTPErrorWithExtensions(*httperror.NewBusinessRuleViolationError(), nil),
		},
		{
			Error: NewHTTPErrorWithExtensions(*httperror.NewForbiddenError(), nil),
		},
		{
			Error: NewHTTPErrorWithExtensions(*httperror.NewInvalidBodyPropertyFormatError(), nil),
		},
		{
			Error: NewHTTPErrorWithExtensions(*httperror.NewInvalidBodyPropertyValueError(), nil),
		},
		{
			Error: NewHTTPErrorWithExtensions(*httperror.NewInvalidRequestHeaderFormatError(), nil),
		},
		{
			Error: NewHTTPErrorWithExtensions(*httperror.NewInvalidRequestParameterFormatError(), nil),
		},
		{
			Error: NewHTTPErrorWithExtensions(*httperror.NewInvalidRequestParameterValueError(), nil),
		},
		{
			Error: NewHTTPErrorWithExtensions(*httperror.NewLicenseCancelledError(), nil),
		},
		{
			Error: NewHTTPErrorWithExtensions(*httperror.NewLicenseExpiredError(), nil),
		},
		{
			Error: NewHTTPErrorWithExtensions(*httperror.NewMissingBodyPropertyError(), nil),
		},
		{
			Error: NewHTTPErrorWithExtensions(*httperror.NewServiceUnavailableError(), nil),
		},
		{
			Error: NewHTTPErrorWithExtensions(*httperror.NewNotFoundError(), nil),
		},
		{
			Error: NewHTTPErrorWithExtensions(*httperror.NewUnauthorizedError(), nil),
		},
		{
			Error: NewHTTPErrorWithExtensions(*httperror.NewServerError(), nil),
		},
		{
			Error: NewHTTPErrorWithExtensions(*httperror.NewMissingRequestHeaderError(), nil),
		},
		{
			Error: NewHTTPErrorWithExtensions(*httperror.NewMissingRequestParameterError(), nil),
		},
		{
			Error: NewHTTPErrorWithExtensions(*httperror.NewValidationError(), map[string]any{
				"foo": "bar",
			}),
		},
		{
			Error: NewHTTPErrorWithExtensions(*httperror.NewHTTPError(http.StatusBadGateway, "bad gateway"), map[string]any{
				"message": "hello world",
			}),
		},
	}

	for _, tc := range testCases {
		rawBytes, err := json.Marshal(tc.Error)
		if err != nil {
			t.Fatal("expected nil error ,got: " + err.Error())
		}

		var result HTTPErrorWithExtensions

		err = json.Unmarshal(rawBytes, &result)
		if err != nil {
			t.Fatal("expected nil error ,got: " + err.Error())
		}

		if tc.Error.Code != result.Code {
			t.Errorf("expected Code=%v, got=%v", tc.Error.Code, result.Code)
		}

		if tc.Error.Detail != result.Detail {
			t.Errorf("expected Detail=%v, got=%v", tc.Error.Detail, result.Detail)
		}

		if tc.Error.Instance != result.Instance {
			t.Errorf("expected Instance=%v, got=%v", tc.Error.Instance, result.Instance)
		}

		if tc.Error.Status != result.Status {
			t.Errorf("expected Status=%v, got=%v", tc.Error.Status, result.Status)
		}

		if tc.Error.Title != result.Title {
			t.Errorf("expected Title=%v, got=%v", tc.Error.Title, result.Title)
		}

		if tc.Error.Type != result.Type {
			t.Errorf("expected Type=%v, got=%v", tc.Error.Type, result.Type)
		}

		if !reflect.DeepEqual(tc.Error.Errors, result.Errors) {
			t.Errorf("expected Errors=%v, got=%v", tc.Error.Errors, result.Errors)
		}

		for key, value := range tc.Error.Extensions {
			if !reflect.DeepEqual(value, result.Extensions[key]) {
				t.Errorf("expected Extensions[%s]=%v, got=%v", key, value, result.Extensions[key])
			}
		}
	}
}
