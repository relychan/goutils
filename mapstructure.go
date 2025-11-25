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
	case bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, *bool, *int, *int8, *int16, *int32, *int64, *uint, *uint8, *uint16, *uint32, *uint64, *float32, *float64, complex64, complex128, *complex64, *complex128:
		return nil, fmt.Errorf("%w, got: %v", ErrMalformedString, reflect.TypeOf(v))
	default:
		return DecodeNullableStringReflection(reflect.ValueOf(value))
	}
}

// DecodeNullableStringReflection a nullable string from reflection value.
func DecodeNullableStringReflection(value reflect.Value) (*string, error) {
	inferredValue, ok := UnwrapPointerFromReflectValue(value)
	if !ok {
		return nil, nil
	}

	result, err := DecodeStringReflection(inferredValue)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// DecodeNullableBooleanSlice decodes a nullable boolean slice from an unknown value.
func DecodeNullableBooleanSlice(value any) (*[]bool, error) {
	if value == nil {
		return nil, nil
	}

	reflectValue, ok := UnwrapPointerFromReflectValue(reflect.ValueOf(value))
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
		elem, err := DecodeNullableBooleanReflection(reflectValue.Index(i))
		if err != nil {
			return nil, fmt.Errorf("failed to decode boolean element at %d: %w", i, err)
		}

		if elem == nil {
			return nil, fmt.Errorf("element %d: %w", i, ErrBooleanNull)
		}

		results[i] = *elem
	}

	return &results, nil
}

// DecodeBooleanSlice decodes a boolean slice from an unknown value.
func DecodeBooleanSlice(value any) ([]bool, error) {
	results, err := DecodeNullableBooleanSlice(value)
	if err != nil {
		return nil, err
	}

	if results == nil {
		return nil, ErrBooleanSliceNull
	}

	return *results, nil
}

// DecodeStringReflection decodes a string from reflection value.
func DecodeStringReflection(value reflect.Value) (string, error) {
	switch value.Kind() {
	case reflect.String:
		return value.String(), nil
	case reflect.Interface:
		return fmt.Sprint(value.Interface()), nil
	default:
		return "", fmt.Errorf("%w, got: %v", ErrMalformedString, value)
	}
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
		return *v, nil
	case bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, *bool, *int, *int8, *int16, *int32, *int64, *uint, *uint8, *uint16, *uint32, *uint64, *float32, *float64:
		return "", fmt.Errorf("%w, got: %s", ErrMalformedString, reflect.TypeOf(v))
	case complex64, complex128, *complex64, *complex128, map[string]any, []any:
		return "", fmt.Errorf("%w, got: %s", ErrMalformedString, reflect.TypeOf(v))
	default:
		return DecodeStringReflection(reflect.ValueOf(value))
	}
}

// DecodeStringSlice decodes a string slice from an unknown value.
func DecodeStringSlice(value any) ([]string, error) {
	if value == nil {
		return nil, nil
	}

	switch vs := value.(type) {
	case []string:
		return vs, nil
	case []*string:
		results := make([]string, len(vs))

		for i, v := range vs {
			if v == nil {
				return nil, fmt.Errorf("failed to decode element at %d: %w", i, ErrStringNull)
			}

			results[i] = *v
		}

		return results, nil
	case []any:
		results := make([]string, len(vs))

		for i, v := range vs {
			s, err := DecodeString(v)
			if err != nil {
				return nil, fmt.Errorf("failed to decode string at %d: %w", i, err)
			}

			results[i] = s
		}

		return results, nil
	case bool, string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, *bool, *string, *int, *int8, *int16, *int32, *int64, *uint, *uint8, *uint16, *uint32, *uint64, *float32, *float64, complex64, complex128, *complex64, *complex128:
		return nil, fmt.Errorf("%w; got: %s", ErrMalformedStringSlice, reflect.TypeOf(vs))
	case []bool, []int, []int8, []int16, []int32, []int64, []uint, []uint8, []uint16, []uint32, []uint64, []float32, []float64, []complex64, []complex128, map[string]any:
		return nil, fmt.Errorf("%w; got: %s", ErrMalformedStringSlice, reflect.TypeOf(vs))
	default:
		return DecodeStringSliceRefection(reflect.ValueOf(value))
	}
}

// DecodeStringSliceRefection decodes a string slice from a reflection value.
func DecodeStringSliceRefection(reflectValue reflect.Value) ([]string, error) {
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

// DecodeNumber tries to convert an unknown value to a typed number.
func DecodeNumber[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64]( //nolint:cyclop,funlen
	value any,
) (T, error) {
	if value == nil {
		return 0, ErrNumberNull
	}

	switch v := value.(type) {
	case int:
		return T(v), nil
	case int8:
		return (T(v)), nil
	case int16:
		return (T(v)), nil
	case int32:
		return (T(v)), nil
	case int64:
		return (T(v)), nil
	case uint:
		return (T(v)), nil
	case uint8:
		return (T(v)), nil
	case uint16:
		return (T(v)), nil
	case uint32:
		return (T(v)), nil
	case uint64:
		return (T(v)), nil
	case float32:
		return (T(v)), nil
	case float64:
		return (T(v)), nil
	case *int:
		return (T(*v)), nil
	case *int8:
		return (T(*v)), nil
	case *int16:
		return (T(*v)), nil
	case *int32:
		return (T(*v)), nil
	case *int64:
		return (T(*v)), nil
	case *uint:
		return (T(*v)), nil
	case *uint8:
		return (T(*v)), nil
	case *uint16:
		return (T(*v)), nil
	case *uint32:
		return (T(*v)), nil
	case *uint64:
		return (T(*v)), nil
	case *float32:
		return (T(*v)), nil
	case *float64:
		return (T(*v)), nil
	case bool, string, complex64, complex128, *bool, *string, *complex64, *complex128, map[string]any, []any, []float64:
		return 0, fmt.Errorf("%w; got: %s", ErrMalformedNumber, reflect.TypeOf(value))
	default:
		return DecodeNumberReflection[T](reflect.ValueOf(value))
	}
}

// DecodeNullableNumber tries to convert an unknown value to a typed number.
func DecodeNullableNumber[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64]( //nolint:cyclop,funlen
	value any,
) (*T, error) {
	if value == nil {
		return nil, nil
	}

	switch v := value.(type) {
	case int:
		return ToPtr(T(v)), nil
	case int8:
		return ToPtr(T(v)), nil
	case int16:
		return ToPtr(T(v)), nil
	case int32:
		return ToPtr(T(v)), nil
	case int64:
		return ToPtr(T(v)), nil
	case uint:
		return ToPtr(T(v)), nil
	case uint8:
		return ToPtr(T(v)), nil
	case uint16:
		return ToPtr(T(v)), nil
	case uint32:
		return ToPtr(T(v)), nil
	case uint64:
		return ToPtr(T(v)), nil
	case float32:
		return ToPtr(T(v)), nil
	case float64:
		return ToPtr(T(v)), nil
	case *int:
		return ToPtr(T(*v)), nil
	case *int8:
		return ToPtr(T(*v)), nil
	case *int16:
		return ToPtr(T(*v)), nil
	case *int32:
		return ToPtr(T(*v)), nil
	case *int64:
		return ToPtr(T(*v)), nil
	case *uint:
		return ToPtr(T(*v)), nil
	case *uint8:
		return ToPtr(T(*v)), nil
	case *uint16:
		return ToPtr(T(*v)), nil
	case *uint32:
		return ToPtr(T(*v)), nil
	case *uint64:
		return ToPtr(T(*v)), nil
	case *float32:
		return ToPtr(T(*v)), nil
	case *float64:
		return ToPtr(T(*v)), nil
	case bool, string, complex64, complex128, *bool, *string, *complex64, *complex128, map[string]any, []any, []float64:
		return nil, fmt.Errorf("%w; got: %s", ErrMalformedNumber, reflect.TypeOf(value))
	default:
		return DecodeNullableNumberReflection[T](reflect.ValueOf(value))
	}
}

// DecodeNullableNumberReflection decodes the nullable floating-point value using reflection.
func DecodeNullableNumberReflection[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](
	value reflect.Value,
) (*T, error) {
	inferredValue, ok := UnwrapPointerFromReflectValue(value)
	if !ok {
		return nil, nil
	}

	result, err := DecodeNumberReflection[T](inferredValue)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// DecodeNumberReflection decodes the number value using reflection.
func DecodeNumberReflection[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](
	value reflect.Value,
) (T, error) {
	kind := value.Kind()

	var result T

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		result = T(value.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		result = T(value.Uint())
	case reflect.Float32, reflect.Float64:
		result = T(value.Float())
	case reflect.String:
		v := value.String()

		newVal, parseErr := strconv.ParseFloat(v, 64)
		if parseErr != nil {
			return T(0), fmt.Errorf("failed to convert number: %w", parseErr)
		}

		result = T(newVal)
	case reflect.Interface:
		v := fmt.Sprint(value.Interface())

		newVal, parseErr := strconv.ParseFloat(v, 64)
		if parseErr != nil {
			return T(0), fmt.Errorf("failed to convert number, got: %w", parseErr)
		}

		result = T(newVal)
	default:
		return T(0), fmt.Errorf("%w, got: %s <%s>", ErrMalformedNumber, value.Type(), kind)
	}

	return result, nil
}

// DecodeNumberSlice decodes a number slice from an unknown value.
func DecodeNumberSlice[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64]( //nolint:cyclop,funlen,gocyclo
	value any,
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
			n, err := DecodeNumber[T](v)
			if err != nil {
				return nil, fmt.Errorf("failed to decode number at %d: %w", i, err)
			}

			results[i] = n
		}

		return results, nil
	case bool, string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, *bool, *string, *int, *int8, *int16, *int32, *int64, *uint, *uint8, *uint16, *uint32, *uint64, *float32, *float64, complex64, complex128, *complex64, *complex128:
		return nil, fmt.Errorf("%w; got: %s", ErrMalformedNumberSlice, reflect.TypeOf(vs))
	case []bool, []string, []complex64, []complex128, map[string]any:
		return nil, fmt.Errorf("%w; got: %s", ErrMalformedNumberSlice, reflect.TypeOf(vs))
	default:
		return DecodeNumberSliceRefection[T](reflect.ValueOf(value))
	}
}

// DecodeNumberSliceRefection decodes a number slice from a reflection value.
func DecodeNumberSliceRefection[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](
	reflectValue reflect.Value,
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
		elem, err := DecodeNullableNumberReflection[T](reflectValue.Index(i))
		if err != nil {
			return nil, fmt.Errorf("failed to decode number element at %d: %w", i, err)
		}

		if elem == nil {
			return nil, fmt.Errorf("failed to number element at %d: %w", i, ErrNumberNull)
		}

		results[i] = *elem
	}

	return results, nil
}

// DecodeNullableBoolean tries to convert an unknown value to a bool pointer.
func DecodeNullableBoolean(value any) (*bool, error) {
	if value == nil {
		return nil, nil
	}

	switch v := value.(type) {
	case bool:
		return &v, nil
	case *bool:
		return v, nil
	default:
		return DecodeNullableBooleanReflection(reflect.ValueOf(value))
	}
}

// DecodeBoolean tries to convert an unknown value to a bool value.
func DecodeBoolean(value any) (bool, error) {
	if value == nil {
		return false, ErrBooleanNull
	}

	switch v := value.(type) {
	case bool:
		return v, nil
	case *bool:
		return *v, nil
	default:
		return DecodeBooleanReflection(reflect.ValueOf(value))
	}
}

// DecodeNullableBooleanReflection decodes a nullable boolean value from reflection.
func DecodeNullableBooleanReflection(value reflect.Value) (*bool, error) {
	inferredValue, ok := UnwrapPointerFromReflectValue(value)
	if !ok {
		return nil, nil
	}

	result, err := DecodeBooleanReflection(inferredValue)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// DecodeBooleanReflection decodes a boolean value from reflection.
func DecodeBooleanReflection(value reflect.Value) (bool, error) {
	kind := value.Kind()

	switch kind {
	case reflect.Bool:
		result := value.Bool()

		return result, nil
	case reflect.Interface:
		if value.Equal(trueValue) {
			return true, nil
		}

		if value.Equal(falseValue) {
			return false, nil
		}
	default:
	}

	return false, fmt.Errorf("%w; got: %v", ErrMalformedBoolean, kind)
}

// GetAny get an unknown value from object by key.
func GetAny(object map[string]any, key string) (any, bool) {
	if object == nil {
		return nil, false
	}

	value, ok := object[key]

	return value, ok
}
