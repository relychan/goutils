package goutils

import (
	"reflect"
	"time"

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
func IsZero[T any](value T) bool { //nolint:cyclop,funlen,gocyclo
	if any(value) == nil {
		return true
	}

	switch v := any(value).(type) {
	case map[string]string:
		return len(v) == 0
	case map[bool]any:
		return len(v) == 0
	case map[int]any:
		return len(v) == 0
	case map[int8]any:
		return len(v) == 0
	case map[int16]any:
		return len(v) == 0
	case map[int32]any:
		return len(v) == 0
	case map[int64]any:
		return len(v) == 0
	case map[uint]any:
		return len(v) == 0
	case map[uint8]any:
		return len(v) == 0
	case map[uint16]any:
		return len(v) == 0
	case map[uint32]any:
		return len(v) == 0
	case map[uint64]any:
		return len(v) == 0
	case map[float32]any:
		return len(v) == 0
	case map[float64]any:
		return len(v) == 0
	case map[complex64]any:
		return len(v) == 0
	case map[complex128]any:
		return len(v) == 0
	case map[string]any:
		return len(v) == 0
	case map[any]any:
		return len(v) == 0
	case []bool:
		return len(v) == 0
	case []string:
		return len(v) == 0
	case []int:
		return len(v) == 0
	case []int8:
		return len(v) == 0
	case []int16:
		return len(v) == 0
	case []int32:
		return len(v) == 0
	case []int64:
		return len(v) == 0
	case []uint:
		return len(v) == 0
	case []uint8:
		return len(v) == 0
	case []uint16:
		return len(v) == 0
	case []uint32:
		return len(v) == 0
	case []uint64:
		return len(v) == 0
	case []float32:
		return len(v) == 0
	case []float64:
		return len(v) == 0
	case []complex64:
		return len(v) == 0
	case []complex128:
		return len(v) == 0
	case []time.Time:
		return len(v) == 0
	case []time.Duration:
		return len(v) == 0
	case []uuid.UUID:
		return len(v) == 0
	case []any:
		return len(v) == 0
	case time.Duration:
		return v == 0
	case IsZeroer:
		return v.IsZero()
	default:
		reflectValue := reflect.ValueOf(value)

		return isZeroReflection(reflectValue)
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

func isZeroReflection(reflectValue reflect.Value) bool {
	switch reflectValue.Kind() {
	case reflect.Array:
		return reflectValue.Len() == 0
	case reflect.Slice, reflect.Map:
		return reflectValue.IsNil() || reflectValue.Len() == 0
	case reflect.Pointer:
		if reflectValue.IsNil() {
			return true
		}

		elem := reflectValue.Elem()

		return isZeroReflection(elem)
	default:
		return false
	}
}
