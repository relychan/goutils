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
	"fmt"
	"reflect"
	"strconv"
)

var (
	trueValue  = reflect.ValueOf(true)
	falseValue = reflect.ValueOf(false)
)

// DecodeNullableString tries to convert an unknown value to a string pointer.
func DecodeNullableString(value any) (*string, error) {
	if value == nil {
		return nil, nil
	}

	switch v := value.(type) {
	case string:
		return &v, nil
	case *string:
		if v == nil {
			return nil, nil
		}

		return v, nil
	case bool,
		int,
		int8,
		int16,
		int32,
		int64,
		uint,
		uint8,
		uint16,
		uint32,
		uint64,
		float32,
		float64,
		*bool,
		*int,
		*int8,
		*int16,
		*int32,
		*int64,
		*uint,
		*uint8,
		*uint16,
		*uint32,
		*uint64,
		*float32,
		*float64,
		complex64,
		complex128,
		*complex64,
		*complex128:
		return nil, fmt.Errorf("%w, got: %v", ErrMalformedString, reflect.TypeOf(v))
	default:
		return DecodeNullableStringReflection(reflect.ValueOf(value))
	}
}

// DecodeNullableStringReflection decodes a nullable string from reflection value.
func DecodeNullableStringReflection(value reflect.Value) (*string, error) {
	inferredValue, ok := UnwrapPointerFromReflectValue(value)
	if !ok {
		return nil, nil
	}

	switch inferredValue.Kind() {
	case reflect.String:
		return new(inferredValue.String()), nil
	case reflect.Interface:
		str, ok := inferredValue.Interface().(string)
		if ok {
			return &str, nil
		}
	default:
	}

	return nil, fmt.Errorf("%w, got: %s", ErrMalformedString, value.Kind())
}

// DecodeStringReflection decodes a string from reflection value.
func DecodeStringReflection(value reflect.Value) (string, error) {
	result, err := DecodeNullableStringReflection(value)
	if err != nil {
		return "", err
	}

	if result == nil {
		return "", ErrStringNull
	}

	return *result, nil
}

// DecodeString tries to convert an unknown value to a string value.
func DecodeString(value any) (string, error) {
	if value == nil {
		return "", ErrStringNull
	}

	switch v := value.(type) {
	case string:
		return v, nil
	case *string:
		if v == nil {
			return "", ErrStringNull
		}

		return *v, nil
	case bool,
		int,
		int8,
		int16,
		int32,
		int64,
		uint,
		uint8,
		uint16,
		uint32,
		uint64,
		float32,
		float64,
		*bool,
		*int,
		*int8,
		*int16,
		*int32,
		*int64,
		*uint,
		*uint8,
		*uint16,
		*uint32,
		*uint64,
		*float32,
		*float64:
		return "", fmt.Errorf("%w, got: %s", ErrMalformedString, reflect.TypeOf(v))
	case complex64, complex128, *complex64, *complex128, map[string]any, []any:
		return "", fmt.Errorf("%w, got: %s", ErrMalformedString, reflect.TypeOf(v))
	default:
		return DecodeStringReflection(reflect.ValueOf(value))
	}
}

// DecodeStringSlice decodes a string slice from an unknown value.
func DecodeStringSlice(value any) ([]string, error) { //nolint:funlen
	if value == nil {
		return nil, nil
	}

	switch vs := value.(type) {
	case []string:
		return vs, nil
	case []*string:
		if vs == nil {
			return nil, nil
		}

		results := make([]string, len(vs))

		for i, v := range vs {
			if v == nil {
				return nil, fmt.Errorf("failed to decode element at %d: %w", i, ErrStringNull)
			}

			results[i] = *v
		}

		return results, nil
	case []any:
		if vs == nil {
			return nil, nil
		}

		results := make([]string, len(vs))

		for i, v := range vs {
			s, err := DecodeString(v)
			if err != nil {
				return nil, fmt.Errorf("failed to decode string at %d: %w", i, err)
			}

			results[i] = s
		}

		return results, nil
	case bool,
		string,
		int,
		int8,
		int16,
		int32,
		int64,
		uint,
		uint8,
		uint16,
		uint32,
		uint64,
		float32,
		float64,
		*bool,
		*string,
		*int,
		*int8,
		*int16,
		*int32,
		*int64,
		*uint,
		*uint8,
		*uint16,
		*uint32,
		*uint64,
		*float32,
		*float64,
		complex64,
		complex128,
		*complex64,
		*complex128:
		return nil, fmt.Errorf("%w; got: %s", ErrMalformedStringSlice, reflect.TypeOf(vs))
	case []bool,
		[]int,
		[]int8,
		[]int16,
		[]int32,
		[]int64,
		[]uint,
		[]uint8,
		[]uint16,
		[]uint32,
		[]uint64,
		[]float32,
		[]float64,
		[]complex64,
		[]complex128,
		map[string]any:
		return nil, fmt.Errorf("%w; got: %s", ErrMalformedStringSlice, reflect.TypeOf(vs))
	default:
		return DecodeStringSliceReflection(reflect.ValueOf(value))
	}
}

// DecodeStringSliceReflection decodes a string slice from a reflection value.
func DecodeStringSliceReflection(reflectValue reflect.Value) ([]string, error) {
	reflectValue, ok := UnwrapPointerFromReflectValue(reflectValue)
	if !ok {
		return nil, nil
	}

	valueKind := reflectValue.Kind()
	if valueKind != reflect.Slice {
		return nil, fmt.Errorf("%w; got: %s", ErrMalformedStringSlice, valueKind)
	}

	valueLen := reflectValue.Len()
	results := make([]string, valueLen)

	for i := range valueLen {
		elem, err := DecodeNullableStringReflection(reflectValue.Index(i))
		if err != nil {
			return nil, fmt.Errorf("failed to decode string element at %d: %w", i, err)
		}

		if elem == nil {
			return nil, fmt.Errorf("failed to decode element at %d: %w", i, ErrStringNull)
		}

		results[i] = *elem
	}

	return results, nil
}

// DecodeNullableBooleanSlice decodes a nullable boolean slice from an unknown value.
func DecodeNullableBooleanSlice(value any) (*[]bool, error) {
	return decodeNullableBooleanSlice(value, false)
}

// AsNullableBooleanSlice tries to cast a nullable boolean slice from an unknown value.
func AsNullableBooleanSlice(value any) (*[]bool, error) {
	return decodeNullableBooleanSlice(value, true)
}

// DecodeBooleanSlice decodes a boolean slice from an unknown value.
func DecodeBooleanSlice(value any) ([]bool, error) {
	return decodeBooleanSlice(value, false)
}

// AsBooleanSlice tries to cast a boolean slice from an unknown value.
func AsBooleanSlice(value any) ([]bool, error) {
	return decodeBooleanSlice(value, true)
}

// DecodeNumber tries to convert an unknown value to a typed number.
func DecodeNumber[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](
	value any,
) (T, error) {
	return decodeNumber[T](value, false, false)
}

// AsNumber tries to cast an unknown value to a typed number.
func AsNumber[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](
	value any,
) (T, error) {
	return decodeNumber[T](value, true, false)
}

// DecodeNullableNumber tries to convert an unknown value to a typed number.
func DecodeNullableNumber[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](
	value any,
) (*T, error) {
	return decodeNullableNumber[T](value, false, false)
}

// AsNullableNumber tries to cast an unknown value to a typed number.
func AsNullableNumber[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](
	value any,
) (*T, error) {
	return decodeNullableNumber[T](value, true, false)
}

// DecodeNullableNumberReflection decodes a nullable numeric value (int, uint, or float) using reflection.
func DecodeNullableNumberReflection[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](
	value reflect.Value,
) (*T, error) {
	return decodeNullableNumberReflection[T](value, false, false)
}

// AsNullableNumberReflection tries to cast a nullable numeric value (int, uint, or float) using reflection.
func AsNullableNumberReflection[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](
	value reflect.Value,
) (*T, error) {
	return decodeNullableNumberReflection[T](value, true, false)
}

// DecodeNumberReflection decodes the number value using reflection.
func DecodeNumberReflection[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](
	value reflect.Value,
) (T, error) {
	return decodeNumberReflection[T](value, false, false)
}

// AsNumberReflection tries to cast the number value using reflection.
func AsNumberReflection[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](
	value reflect.Value,
) (T, error) {
	return decodeNumberReflection[T](value, true, false)
}

// DecodeNumberSlice decodes a number slice from an unknown value.
func DecodeNumberSlice[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](
	value any,
) ([]T, error) {
	return decodeNumberSlice[T](value, false)
}

// AsNumberSlice tries to cast a number slice from an unknown value.
func AsNumberSlice[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](
	value any,
) ([]T, error) {
	return decodeNumberSlice[T](value, true)
}

// DecodeNumberSliceReflection decodes a number slice from a reflection value.
func DecodeNumberSliceReflection[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](
	reflectValue reflect.Value,
) ([]T, error) {
	return decodeNumberSliceReflection[T](reflectValue, false)
}

// AsNumberSliceReflection tries to cast a number slice from a reflection value.
func AsNumberSliceReflection[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](
	reflectValue reflect.Value,
) ([]T, error) {
	return decodeNumberSliceReflection[T](reflectValue, true)
}

// DecodeNullableBoolean tries to convert an unknown value to a bool pointer.
func DecodeNullableBoolean(value any) (*bool, error) {
	return decodeNullableBoolean(value, false)
}

// AsNullableBoolean tries to cast an unknown value to a bool pointer.
func AsNullableBoolean(value any) (*bool, error) {
	return decodeNullableBoolean(value, true)
}

// DecodeBoolean tries to convert an unknown value to a bool value.
func DecodeBoolean(value any) (bool, error) {
	return decodeBoolean(value, false)
}

// AsBoolean tries to cast an unknown value to a bool value.
func AsBoolean(value any) (bool, error) {
	return decodeBoolean(value, true)
}

// DecodeNullableBooleanReflection decodes a nullable boolean value from reflection.
func DecodeNullableBooleanReflection(value reflect.Value) (*bool, error) {
	return decodeNullableBooleanReflection(value, false)
}

// AsNullableBooleanReflection tries to cast a nullable boolean value from reflection.
func AsNullableBooleanReflection(value reflect.Value) (*bool, error) {
	return decodeNullableBooleanReflection(value, true)
}

// DecodeBooleanReflection decodes a boolean value from reflection.
func DecodeBooleanReflection(value reflect.Value) (bool, error) {
	return decodeBooleanReflection(value, false)
}

// AsBooleanReflection tries to cast a boolean value from reflection.
func AsBooleanReflection(value reflect.Value) (bool, error) {
	return decodeBooleanReflection(value, true)
}

func decodeBooleanSlice(value any, strict bool) ([]bool, error) {
	results, err := decodeNullableBooleanSlice(value, strict)
	if err != nil {
		return nil, err
	}

	if results == nil {
		return nil, ErrBooleanSliceNull
	}

	return *results, nil
}

func decodeNullableBooleanSliceReflection(value reflect.Value, strict bool) (*[]bool, error) {
	reflectValue, ok := UnwrapPointerFromReflectValue(value)
	if !ok {
		return nil, nil
	}

	valueKind := reflectValue.Kind()
	if valueKind != reflect.Slice {
		return nil, fmt.Errorf("%w; got: %s", ErrMalformedBooleanSlice, valueKind)
	}

	valueLen := reflectValue.Len()
	results := make([]bool, valueLen)

	for i := range valueLen {
		elem, err := decodeNullableBooleanReflection(reflectValue.Index(i), strict)
		if err != nil {
			return nil, fmt.Errorf("failed to decode boolean element at %d: %w", i, err)
		}

		if elem == nil {
			return nil, fmt.Errorf("failed to decode boolean element at %d: %w", i, ErrBooleanNull)
		}

		results[i] = *elem
	}

	return &results, nil
}

func decodeNumber[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64]( //nolint:cyclop,funlen,gocyclo,gocognit
	value any,
	strict bool,
	isRecursive bool,
) (T, error) {
	if value == nil {
		return 0, ErrNumberNull
	}

	switch v := value.(type) {
	case int:
		return T(v), nil
	case int8:
		return T(v), nil
	case int16:
		return T(v), nil
	case int32:
		return T(v), nil
	case int64:
		return T(v), nil
	case uint:
		return T(v), nil
	case uint8:
		return T(v), nil
	case uint16:
		return T(v), nil
	case uint32:
		return T(v), nil
	case uint64:
		return T(v), nil
	case float32:
		return T(v), nil
	case float64:
		return T(v), nil
	case *int:
		if v == nil {
			return 0, ErrNumberNull
		}

		return T(*v), nil
	case *int8:
		if v == nil {
			return 0, ErrNumberNull
		}

		return T(*v), nil
	case *int16:
		if v == nil {
			return 0, ErrNumberNull
		}

		return T(*v), nil
	case *int32:
		if v == nil {
			return 0, ErrNumberNull
		}

		return T(*v), nil
	case *int64:
		if v == nil {
			return 0, ErrNumberNull
		}

		return T(*v), nil
	case *uint:
		if v == nil {
			return 0, ErrNumberNull
		}

		return T(*v), nil
	case *uint8:
		if v == nil {
			return 0, ErrNumberNull
		}

		return T(*v), nil
	case *uint16:
		if v == nil {
			return 0, ErrNumberNull
		}

		return T(*v), nil
	case *uint32:
		if v == nil {
			return 0, ErrNumberNull
		}

		return T(*v), nil
	case *uint64:
		if v == nil {
			return 0, ErrNumberNull
		}

		return T(*v), nil
	case *float32:
		if v == nil {
			return 0, ErrNumberNull
		}

		return T(*v), nil
	case *float64:
		if v == nil {
			return 0, ErrNumberNull
		}

		return T(*v), nil
	case string:
		if strict {
			return 0, fmt.Errorf("%w; got: %s", ErrMalformedNumber, reflect.TypeOf(value))
		}

		return parseNumber[T](v)
	case *string:
		if strict {
			return 0, fmt.Errorf("%w; got: %s", ErrMalformedNumber, reflect.TypeOf(value))
		}

		if v == nil {
			return 0, ErrNumberNull
		}

		return parseNumber[T](*v)
	case bool,
		complex64,
		complex128,
		*bool,
		*complex64,
		*complex128,
		map[string]any,
		[]any,
		[]float64:
		return 0, fmt.Errorf("%w; got: %s", ErrMalformedNumber, reflect.TypeOf(value))
	default:
		return decodeNumberReflection[T](reflect.ValueOf(value), strict, isRecursive)
	}
}

func decodeNullableNumber[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64]( //nolint:cyclop,funlen,gocyclo,gocognit
	value any,
	strict bool,
	isRecursive bool,
) (*T, error) {
	if value == nil {
		return nil, nil
	}

	switch v := value.(type) {
	case int:
		return new(T(v)), nil
	case int8:
		return new(T(v)), nil
	case int16:
		return new(T(v)), nil
	case int32:
		return new(T(v)), nil
	case int64:
		return new(T(v)), nil
	case uint:
		return new(T(v)), nil
	case uint8:
		return new(T(v)), nil
	case uint16:
		return new(T(v)), nil
	case uint32:
		return new(T(v)), nil
	case uint64:
		return new(T(v)), nil
	case float32:
		return new(T(v)), nil
	case float64:
		return new(T(v)), nil
	case *int:
		if v == nil {
			return nil, nil
		}

		return new(T(*v)), nil
	case *int8:
		if v == nil {
			return nil, nil
		}

		return new(T(*v)), nil
	case *int16:
		if v == nil {
			return nil, nil
		}

		return new(T(*v)), nil
	case *int32:
		if v == nil {
			return nil, nil
		}

		return new(T(*v)), nil
	case *int64:
		if v == nil {
			return nil, nil
		}

		return new(T(*v)), nil
	case *uint:
		if v == nil {
			return nil, nil
		}

		return new(T(*v)), nil
	case *uint8:
		if v == nil {
			return nil, nil
		}

		return new(T(*v)), nil
	case *uint16:
		if v == nil {
			return nil, nil
		}

		return new(T(*v)), nil
	case *uint32:
		if v == nil {
			return nil, nil
		}

		return new(T(*v)), nil
	case *uint64:
		if v == nil {
			return nil, nil
		}

		return new(T(*v)), nil
	case *float32:
		if v == nil {
			return nil, nil
		}

		return new(T(*v)), nil
	case *float64:
		if v == nil {
			return nil, nil
		}

		return new(T(*v)), nil
	case string:
		if strict {
			return nil, fmt.Errorf("%w; got: %s", ErrMalformedNumber, reflect.TypeOf(value))
		}

		result, err := parseNumber[T](v)
		if err != nil {
			return nil, err
		}

		return &result, nil
	case *string:
		if strict {
			return nil, fmt.Errorf("%w; got: %s", ErrMalformedNumber, reflect.TypeOf(value))
		}

		if v == nil {
			return nil, nil
		}

		result, err := parseNumber[T](*v)
		if err != nil {
			return nil, err
		}

		return &result, nil
	case bool,
		complex64,
		complex128,
		*bool,
		*complex64,
		*complex128,
		map[string]any,
		[]any,
		[]float64:
		return nil, fmt.Errorf("%w; got: %s", ErrMalformedNumber, reflect.TypeOf(value))
	default:
		return decodeNullableNumberReflection[T](reflect.ValueOf(value), strict, isRecursive)
	}
}

func decodeNumberSlice[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64]( //nolint:cyclop,funlen,gocyclo
	value any,
	strict bool,
) ([]T, error) {
	if value == nil {
		return nil, nil
	}

	switch vs := value.(type) {
	case []int:
		return ToNumberSlice[int, T](vs), nil
	case []int8:
		return ToNumberSlice[int8, T](vs), nil
	case []int16:
		return ToNumberSlice[int16, T](vs), nil
	case []int32:
		return ToNumberSlice[int32, T](vs), nil
	case []int64:
		return ToNumberSlice[int64, T](vs), nil
	case []uint:
		return ToNumberSlice[uint, T](vs), nil
	case []uint8:
		return ToNumberSlice[uint8, T](vs), nil
	case []uint16:
		return ToNumberSlice[uint16, T](vs), nil
	case []uint32:
		return ToNumberSlice[uint32, T](vs), nil
	case []uint64:
		return ToNumberSlice[uint64, T](vs), nil
	case []float32:
		return ToNumberSlice[float32, T](vs), nil
	case []float64:
		return ToNumberSlice[float64, T](vs), nil
	case []*int:
		return PtrToNumberSlice[int, T](vs)
	case []*int8:
		return PtrToNumberSlice[int8, T](vs)
	case []*int16:
		return PtrToNumberSlice[int16, T](vs)
	case []*int32:
		return PtrToNumberSlice[int32, T](vs)
	case []*int64:
		return PtrToNumberSlice[int64, T](vs)
	case []*uint:
		return PtrToNumberSlice[uint, T](vs)
	case []*uint8:
		return PtrToNumberSlice[uint8, T](vs)
	case []*uint16:
		return PtrToNumberSlice[uint16, T](vs)
	case []*uint32:
		return PtrToNumberSlice[uint32, T](vs)
	case []*uint64:
		return PtrToNumberSlice[uint64, T](vs)
	case []*float32:
		return PtrToNumberSlice[float32, T](vs)
	case []*float64:
		return PtrToNumberSlice[float64, T](vs)
	case []any:
		results := make([]T, len(vs))

		for i, v := range vs {
			n, err := decodeNumber[T](v, strict, false)
			if err != nil {
				return nil, fmt.Errorf("failed to decode number at %d: %w", i, err)
			}

			results[i] = n
		}

		return results, nil
	case []string:
		if strict {
			return nil, fmt.Errorf(
				"%w; got: %s",
				ErrMalformedNumberSlice,
				reflect.TypeOf(value),
			)
		}

		results := make([]T, len(vs))

		for i, v := range vs {
			num, err := parseNumber[T](v)
			if err != nil {
				return nil, fmt.Errorf("failed to decode number at %d: %w", i, err)
			}

			results[i] = num
		}

		return results, nil
	case bool,
		string,
		int,
		int8,
		int16,
		int32,
		int64,
		uint,
		uint8,
		uint16,
		uint32,
		uint64,
		float32,
		float64,
		*bool,
		*string,
		*int,
		*int8,
		*int16,
		*int32,
		*int64,
		*uint,
		*uint8,
		*uint16,
		*uint32,
		*uint64,
		*float32,
		*float64,
		complex64,
		complex128,
		*complex64,
		*complex128:
		return nil, fmt.Errorf("%w; got: %s", ErrMalformedNumberSlice, reflect.TypeOf(vs))
	case []bool, []complex64, []complex128, map[string]any:
		return nil, fmt.Errorf("%w; got: %s", ErrMalformedNumberSlice, reflect.TypeOf(vs))
	default:
		return decodeNumberSliceReflection[T](reflect.ValueOf(value), strict)
	}
}

func decodeNumberSliceReflection[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](
	reflectValue reflect.Value,
	strict bool,
) ([]T, error) {
	reflectValue, ok := UnwrapPointerFromReflectValue(reflectValue)
	if !ok {
		return nil, nil
	}

	valueKind := reflectValue.Kind()
	if valueKind != reflect.Slice {
		return nil, fmt.Errorf("%w; got: %s", ErrMalformedNumberSlice, valueKind)
	}

	valueLen := reflectValue.Len()
	results := make([]T, valueLen)

	for i := range valueLen {
		elem, err := decodeNullableNumberReflection[T](reflectValue.Index(i), strict, false)
		if err != nil {
			return nil, fmt.Errorf("failed to decode number element at %d: %w", i, err)
		}

		if elem == nil {
			return nil, fmt.Errorf("failed to decode number element at %d: %w", i, ErrNumberNull)
		}

		results[i] = *elem
	}

	return results, nil
}

func decodeNullableNumberReflection[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](
	value reflect.Value,
	strict bool,
	isRecursive bool,
) (*T, error) {
	inferredValue, ok := UnwrapPointerFromReflectValue(value)
	if !ok {
		return nil, nil
	}

	kind := inferredValue.Kind()

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return new(T(inferredValue.Int())), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return new(T(inferredValue.Uint())), nil
	case reflect.Float32, reflect.Float64:
		return new(T(inferredValue.Float())), nil
	case reflect.String:
		if strict {
			return nil, fmt.Errorf("%w, got: %s <%s>", ErrMalformedNumber, value.Type(), kind)
		}

		result, err := parseNumber[T](inferredValue.String())
		if err != nil {
			return nil, err
		}

		return &result, nil
	case reflect.Interface:
		// guard against infinite loop.
		if isRecursive {
			return nil, fmt.Errorf("%w, got: %s <%s>", ErrMalformedNumber, value.Type(), kind)
		}

		result, err := decodeNullableNumber[T](inferredValue.Interface(), strict, true)
		if err != nil {
			return nil, err
		}

		return result, nil
	default:
		return nil, fmt.Errorf("%w, got: %s <%s>", ErrMalformedNumber, value.Type(), kind)
	}
}

func decodeNumberReflection[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](
	value reflect.Value,
	strict bool,
	isRecursive bool,
) (T, error) {
	result, err := decodeNullableNumberReflection[T](value, strict, isRecursive)
	if err != nil {
		return 0, err
	}

	if result == nil {
		return 0, ErrNumberNull
	}

	return *result, nil
}

func decodeNullableBoolean(value any, strict bool) (*bool, error) {
	if value == nil {
		return nil, nil
	}

	switch v := value.(type) {
	case bool:
		return &v, nil
	case *bool:
		return v, nil
	case string:
		if strict {
			return nil, ErrMalformedBoolean
		}

		result, err := parseBool(v)
		if err != nil {
			return nil, err
		}

		return &result, nil
	case *string:
		if strict {
			return nil, ErrMalformedBoolean
		}

		if v == nil {
			return nil, nil
		}

		result, err := parseBool(*v)
		if err != nil {
			return nil, err
		}

		return &result, nil
	default:
		return decodeNullableBooleanReflection(reflect.ValueOf(value), strict)
	}
}

func decodeNullableBooleanReflection(reflectValue reflect.Value, strict bool) (*bool, error) {
	value, ok := UnwrapPointerFromReflectValue(reflectValue)
	if !ok {
		return nil, nil
	}

	kind := value.Kind()

	switch kind {
	case reflect.Bool:
		return new(value.Bool()), nil
	case reflect.String:
		if !strict {
			result, err := parseBool(value.String())
			if err != nil {
				return nil, err
			}

			return &result, nil
		}
	case reflect.Interface:
		if value.Equal(trueValue) {
			return new(true), nil
		}

		if value.Equal(falseValue) {
			return new(false), nil
		}
	default:
	}

	return nil, fmt.Errorf("%w; got: %v", ErrMalformedBoolean, kind)
}

func decodeBooleanReflection(value reflect.Value, strict bool) (bool, error) {
	result, err := decodeNullableBooleanReflection(value, strict)
	if err != nil {
		return false, err
	}

	if result == nil {
		return false, ErrBooleanNull
	}

	return *result, nil
}

func decodeNullableBooleanSlice(value any, strict bool) (*[]bool, error) { //nolint:cyclop,funlen
	if value == nil {
		return nil, nil
	}

	switch vs := value.(type) {
	case []bool:
		return &vs, nil
	case *[]bool:
		return vs, nil
	case []*bool:
		results := make([]bool, len(vs))

		for i, v := range vs {
			if v == nil {
				return nil, fmt.Errorf(
					"failed to decode boolean element at %d: %w",
					i,
					ErrBooleanNull,
				)
			}

			results[i] = *v
		}

		return &results, nil
	case []any:
		results := make([]bool, len(vs))

		for i, v := range vs {
			n, err := decodeBoolean(v, strict)
			if err != nil {
				return nil, fmt.Errorf("failed to decode boolean element at %d: %w", i, err)
			}

			results[i] = n
		}

		return &results, nil
	case []string:
		if strict {
			return nil, fmt.Errorf(
				"%w; got: %s",
				ErrMalformedBooleanSlice,
				reflect.TypeOf(value),
			)
		}

		results := make([]bool, len(vs))

		for i, v := range vs {
			num, err := parseBool(v)
			if err != nil {
				return nil, fmt.Errorf("failed to decode boolean element at %d: %w", i, err)
			}

			results[i] = num
		}

		return &results, nil
	case bool,
		string,
		int,
		int8,
		int16,
		int32,
		int64,
		uint,
		uint8,
		uint16,
		uint32,
		uint64,
		float32,
		float64,
		*bool,
		*string,
		*int,
		*int8,
		*int16,
		*int32,
		*int64,
		*uint,
		*uint8,
		*uint16,
		*uint32,
		*uint64,
		*float32,
		*float64,
		complex64,
		complex128,
		*complex64,
		*complex128:
		return nil, fmt.Errorf("%w; got: %s", ErrMalformedBooleanSlice, reflect.TypeOf(vs))
	case []int,
		[]int8,
		[]int16,
		[]int32,
		[]int64,
		[]uint,
		[]uint8,
		[]uint16,
		[]uint32,
		[]uint64,
		[]float32,
		[]float64,
		[]complex64,
		[]complex128,
		map[string]any:
		return nil, fmt.Errorf("%w; got: %s", ErrMalformedBooleanSlice, reflect.TypeOf(vs))
	default:
		return decodeNullableBooleanSliceReflection(reflect.ValueOf(value), strict)
	}
}

func parseNumber[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](
	rawValue string,
) (T, error) {
	var empty T

	switch any(empty).(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32:
		result, err := strconv.ParseInt(rawValue, 10, 64)
		if err != nil {
			return 0, ErrMalformedNumber
		}

		return T(result), nil
	case uint64:
		result, err := strconv.ParseUint(rawValue, 10, 64)
		if err != nil {
			return 0, ErrMalformedNumber
		}

		return T(result), nil
	default:
		fResult, err := strconv.ParseFloat(rawValue, 64)
		if err != nil {
			return 0, ErrMalformedNumber
		}

		return T(fResult), nil
	}
}

func decodeBoolean(value any, strict bool) (bool, error) {
	if value == nil {
		return false, ErrBooleanNull
	}

	switch v := value.(type) {
	case bool:
		return v, nil
	case *bool:
		if v == nil {
			return false, ErrBooleanNull
		}

		return *v, nil
	case string:
		if strict {
			return false, ErrMalformedBoolean
		}

		return parseBool(v)
	case *string:
		if strict {
			return false, ErrMalformedBoolean
		}

		if v == nil {
			return false, ErrBooleanNull
		}

		return parseBool(*v)
	default:
		return decodeBooleanReflection(reflect.ValueOf(value), strict)
	}
}

func parseBool(value string) (bool, error) {
	result, err := strconv.ParseBool(value)
	if err != nil {
		return false, ErrMalformedBoolean
	}

	return result, nil
}
