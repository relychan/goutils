package goutils

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestRFC9457Error(t *testing.T) {
	testCases := []struct {
		Error RFC9457ErrorWithExtensions
	}{
		{
			Error: NewRFC9457ErrorWithExtensions(NewAlreadyExistsError(), nil),
		},
		{
			Error: NewRFC9457ErrorWithExtensions(NewForbiddenError(ErrorDetail{
				Detail: "forbidden",
			}), nil),
		},
		{
			Error: NewRFC9457ErrorWithExtensions(NewBadRequestError(), nil),
		},
		{
			Error: NewRFC9457ErrorWithExtensions(NewBusinessRuleViolationError(), nil),
		},
		{
			Error: NewRFC9457ErrorWithExtensions(NewForbiddenError(), nil),
		},
		{
			Error: NewRFC9457ErrorWithExtensions(NewInvalidBodyPropertyFormatError(), nil),
		},
		{
			Error: NewRFC9457ErrorWithExtensions(NewInvalidBodyPropertyValueError(), nil),
		},
		{
			Error: NewRFC9457ErrorWithExtensions(NewInvalidRequestHeaderFormatError(), nil),
		},
		{
			Error: NewRFC9457ErrorWithExtensions(NewInvalidRequestParameterFormatError(), nil),
		},
		{
			Error: NewRFC9457ErrorWithExtensions(NewInvalidRequestParameterValueError(), nil),
		},
		{
			Error: NewRFC9457ErrorWithExtensions(NewLicenseCancelledError(), nil),
		},
		{
			Error: NewRFC9457ErrorWithExtensions(NewLicenseExpiredError(), nil),
		},
		{
			Error: NewRFC9457ErrorWithExtensions(NewMissingBodyPropertyError(), nil),
		},
		{
			Error: NewRFC9457ErrorWithExtensions(NewServiceUnavailableError(), nil),
		},
		{
			Error: NewRFC9457ErrorWithExtensions(NewNotFoundError(), nil),
		},
		{
			Error: NewRFC9457ErrorWithExtensions(NewUnauthorizedError(), nil),
		},
		{
			Error: NewRFC9457ErrorWithExtensions(NewServerError(), nil),
		},
		{
			Error: NewRFC9457ErrorWithExtensions(NewMissingRequestHeaderError(), nil),
		},
		{
			Error: NewRFC9457ErrorWithExtensions(NewMissingRequestParameterError(), nil),
		},
		{
			Error: NewRFC9457ErrorWithExtensions(NewValidationError(), map[string]any{
				"foo": "bar",
			}),
		},
		{
			Error: NewRFC9457ErrorWithExtensions(NewRFC9457Error(http.StatusBadGateway, "bad gateway"), map[string]any{
				"message": "hello world",
			}),
		},
	}

	for _, tc := range testCases {
		rawBytes, err := json.Marshal(tc.Error)
		assertNilError(t, err)

		var result RFC9457ErrorWithExtensions

		err = json.Unmarshal(rawBytes, &result)
		assertNilError(t, err)

		assertEqual(t, tc.Error.Code, result.Code)
		assertEqual(t, tc.Error.Detail, result.Detail)
		assertEqual(t, tc.Error.Instance, result.Instance)
		assertEqual(t, tc.Error.Status, result.Status)
		assertEqual(t, tc.Error.Title, result.Title)
		assertEqual(t, tc.Error.Type, result.Type)
		assertDeepEqual(t, tc.Error.Errors, result.Errors)

		for key, value := range tc.Error.Extensions {
			assertDeepEqual(t, value, result.Extensions[key])
		}
	}
}
