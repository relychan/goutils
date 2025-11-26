package goutils

import (
	"reflect"

	"github.com/google/uuid"
)

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

// IsNil a safe function to check null value.
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

// UnwrapPointerFromReflectValue unwraps pointers from the reflect value.
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
