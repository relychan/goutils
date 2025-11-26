package goutils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// IsMetaCharacter checks if the character is a word character.
// A word character is a character a-z, A-Z, 0-9, including _ (underscore) and - (hyphen).
func IsMetaCharacter[C byte | rune](c C) bool {
	return c == '-' || c == '_' ||
		IsDigit(c) ||
		IsLowerAlphabet(c) ||
		IsUpperAlphabet(c)
}

// IsDigit checks if the character is a digit.
func IsDigit[C byte | rune](c C) bool {
	return c >= '0' && c <= '9'
}

// IsLowerAlphabet checks if the character is a lowercase alphabet.
func IsLowerAlphabet[C byte | rune](c C) bool {
	return c >= 'a' && c <= 'z'
}

// IsUpperAlphabet checks if the character is an uppercase alphabet.
func IsUpperAlphabet[C byte | rune](c C) bool {
	return c >= 'A' && c <= 'Z'
}

// ParseIntInRange parses s as an integer and
// verifies that it is within some range.
// If it is invalid or out-of-range,
// it sets ok to false and returns the min value.
func ParseIntInRange[B []byte | string](s B, minValue int, maxValue int) (int, bool) {
	var x int

	for _, c := range []byte(s) {
		if !IsDigit(c) {
			return minValue, false
		}

		x = x*10 + int(c) - '0'
	}

	if x < minValue || maxValue < x {
		return minValue, false
	}

	return x, true
}

// ToString converts an arbitrary value to its string representation.
//
// It handles the following cases:
//   - For nil values, it returns the provided emptyValue.
//   - For primitive types (bool, string, integers, floats, complex), it uses the appropriate formatting.
//   - For time.Time and *time.Time, it uses the standard time formatting.
//   - For types implementing fmt.Stringer, it uses their String() method.
//   - For pointers, it dereferences and formats the underlying value, or returns emptyValue if nil.
//   - For unsupported types, it attempts to marshal the value to JSON.
//
// If JSON marshaling fails, it returns an error.
//
// The emptyValue parameter specifies the string to return for nil values or nil pointers.
func ToString(value any, emptyValue string) (string, error) {
	return ToStringWithCustomTypeFormatter(value, emptyValue, func(value any) (string, error) {
		jsonValue, err := json.Marshal(value)
		if err != nil {
			return "", err
		}

		return string(jsonValue), nil
	})
}

// ToDebugString returns a string representation of an arbitrary value for debugging or logging purposes.
// It wraps ToString, and if ToString returns an error, it falls back to fmt.Sprint to ensure a string is always returned.
// This function never returns an error; it always returns a string.
// The emptyValue parameter specifies the string to use if the input value is nil or otherwise empty.
func ToDebugString(value any, emptyValue string) string {
	result, err := ToString(value, emptyValue)
	if err != nil {
		return fmt.Sprint(value)
	}

	return result
}

// ToStringWithCustomTypeFormatter converts an arbitrary value to its string representation.
//
// It handles the following cases:
//   - For nil values, it returns the provided emptyValue.
//   - For primitive types (bool, string, integers, floats, complex), it uses the appropriate formatting.
//   - For time.Time and *time.Time, it uses the standard time formatting.
//   - For types implementing fmt.Stringer, it uses their String() method.
//   - For pointers, it dereferences and formats the underlying value, or returns emptyValue if nil.
//   - For unsupported types, it invokes the customTypeFormatter function to format the value.
//
// If JSON marshaling fails, it returns an error.
//
// The emptyValue parameter specifies the string to return for nil values or nil pointers.
func ToStringWithCustomTypeFormatter( //nolint:cyclop,gocognit,gocyclo,funlen,maintidx
	value any,
	emptyValue string,
	customTypeFormatter func(any) (string, error),
) (string, error) {
	if value == nil {
		return emptyValue, nil
	}

	switch typedValue := value.(type) {
	case bool:
		return strconv.FormatBool(typedValue), nil
	case string:
		return typedValue, nil
	case int:
		return strconv.FormatInt(int64(typedValue), 10), nil
	case int8:
		return strconv.FormatInt(int64(typedValue), 10), nil
	case int16:
		return strconv.FormatInt(int64(typedValue), 10), nil
	case int32:
		return strconv.FormatInt(int64(typedValue), 10), nil
	case int64:
		return strconv.FormatInt(typedValue, 10), nil
	case uint:
		return strconv.FormatUint(uint64(typedValue), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(typedValue), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(typedValue), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(typedValue), 10), nil
	case uint64:
		return strconv.FormatUint(typedValue, 10), nil
	case float32:
		return strconv.FormatFloat(float64(typedValue), 'f', -1, 32), nil
	case float64:
		return strconv.FormatFloat(float64(typedValue), 'f', -1, 64), nil
	case complex64:
		return strconv.FormatComplex(complex128(typedValue), 'f', -1, 64), nil
	case complex128:
		return strconv.FormatComplex(complex128(typedValue), 'f', -1, 128), nil
	case time.Time:
		return typedValue.Format(time.RFC3339), nil
	case time.Duration:
		return typedValue.String(), nil
	case *bool:
		if typedValue == nil {
			return emptyValue, nil
		}

		return strconv.FormatBool(*typedValue), nil
	case *string:
		if typedValue == nil {
			return emptyValue, nil
		}

		return *typedValue, nil
	case *int:
		if typedValue == nil {
			return emptyValue, nil
		}

		return strconv.FormatInt(int64(*typedValue), 10), nil
	case *int8:
		if typedValue == nil {
			return emptyValue, nil
		}

		return strconv.FormatInt(int64(*typedValue), 10), nil
	case *int16:
		if typedValue == nil {
			return emptyValue, nil
		}

		return strconv.FormatInt(int64(*typedValue), 10), nil
	case *int32:
		if typedValue == nil {
			return emptyValue, nil
		}

		return strconv.FormatInt(int64(*typedValue), 10), nil
	case *int64:
		if typedValue == nil {
			return emptyValue, nil
		}

		return strconv.FormatInt(*typedValue, 10), nil
	case *uint:
		if typedValue == nil {
			return emptyValue, nil
		}

		return strconv.FormatUint(uint64(*typedValue), 10), nil
	case *uint8:
		if typedValue == nil {
			return emptyValue, nil
		}

		return strconv.FormatUint(uint64(*typedValue), 10), nil
	case *uint16:
		if typedValue == nil {
			return emptyValue, nil
		}

		return strconv.FormatUint(uint64(*typedValue), 10), nil
	case *uint32:
		if typedValue == nil {
			return emptyValue, nil
		}

		return strconv.FormatUint(uint64(*typedValue), 10), nil
	case *uint64:
		if typedValue == nil {
			return emptyValue, nil
		}

		return strconv.FormatUint(*typedValue, 10), nil
	case *float32:
		if typedValue == nil {
			return emptyValue, nil
		}

		return strconv.FormatFloat(float64(*typedValue), 'f', -1, 32), nil
	case *float64:
		if typedValue == nil {
			return emptyValue, nil
		}

		return strconv.FormatFloat(float64(*typedValue), 'f', -1, 64), nil
	case *complex64:
		if typedValue == nil {
			return emptyValue, nil
		}

		return strconv.FormatComplex(complex128(*typedValue), 'f', -1, 64), nil
	case *complex128:
		if typedValue == nil {
			return emptyValue, nil
		}

		return strconv.FormatComplex(complex128(*typedValue), 'f', -1, 128), nil
	case *time.Time:
		if typedValue == nil {
			return emptyValue, nil
		}

		return typedValue.Format(time.RFC3339), nil
	case *time.Duration:
		if typedValue == nil {
			return emptyValue, nil
		}

		return typedValue.String(), nil
	case fmt.Stringer:
		if IsNil(typedValue) {
			return emptyValue, nil
		}

		return typedValue.String(), nil
	default:
		return customTypeFormatter(value)
	}
}
