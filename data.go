package goutils

import (
	"reflect"

	"github.com/google/uuid"
)

// IsZeroer abstracts an interface to check if the instance is zero.
type IsZeroer interface {
	IsZero() bool
}

// IsZeroPtr checks if the pointer is zero.
func IsZeroPtr[T any](ptr *T) bool {
	if ptr == nil {
		return true
	}

	return IsZero(*ptr)
}

// IsZero checks if the value is zero.
func IsZero[T any](value T) bool {
	switch v := any(value).(type) {
	case IsZeroer:
		return v.IsZero()
	case []any:
		return len(v) == 0
	case map[string]any:
		return len(v) == 0
	case map[any]any:
		return len(v) == 0
	default:
		return false
	}
}

// EqualPtr checks if the value of both pointers are equal.
func EqualPtr[T Equaler[T]](a, b *T) bool {
	if a == nil && b == nil {
		return true
	}

	return a != nil && b != nil && (*a).Equal(*b)
}

// NewUUIDv7 creates a random UUID version 7. Fallback to v4 if there is error.
func NewUUIDv7() uuid.UUID {
	value, err := uuid.NewV7()
	if err == nil {
		return value
	}

	return uuid.New()
}

// ToPtr converts the a typed value to its pointer.
func ToPtr[T any](value T) *T {
	return &value
}

// IsNil safely checks if a value is nil.
func IsNil(value any) bool {
	if value == nil {
		return true
	}

	v := reflect.ValueOf(value)
	// Check for all kinds that can be nil
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
}

// UnwrapPointerFromReflectValue recursively unwraps pointers from a reflect value.
// Returns the unwrapped value and true if the value is valid and not nil,
// or false if the value is nil or invalid.
func UnwrapPointerFromReflectValue(reflectValue reflect.Value) (reflect.Value, bool) {
	switch reflectValue.Kind() {
	case reflect.Chan, reflect.Func, reflect.Invalid:
		return reflectValue, false
	case reflect.Pointer:
		if reflectValue.IsNil() {
			return reflectValue, false
		}

		return UnwrapPointerFromReflectValue(reflectValue.Elem())
	case reflect.Slice, reflect.Interface, reflect.Map:
		return reflectValue, !reflectValue.IsNil()
	default:
		return reflectValue, true
	}
}
