package goutils

import (
	"reflect"
	"time"

	"github.com/google/uuid"
)

// Equaler abstracts an interface to check the equality.
type Equaler[T any] interface {
	// Equal checks if the target value is equal.
	Equal(target T) bool
}

// DeepEqual checks if both values are equal recursively.
func DeepEqual[T any](x, y T, omitZero bool) bool { //nolint:cyclop,funlen,gocyclo,gocognit,maintidx
	switch vx := any(x).(type) {
	case bool, string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, complex64, complex128, time.Duration, time.Time, uuid.UUID:
		return EqualComparable(vx, y)
	case *bool:
		return EqualComparablePtr(vx, y)
	case *string:
		return EqualComparablePtr(vx, y)
	case *int:
		return EqualComparablePtr(vx, y)
	case *int8:
		return EqualComparablePtr(vx, y)
	case *int16:
		return EqualComparablePtr(vx, y)
	case *int32:
		return EqualComparablePtr(vx, y)
	case *int64:
		return EqualComparablePtr(vx, y)
	case *uint:
		return EqualComparablePtr(vx, y)
	case *uint8:
		return EqualComparablePtr(vx, y)
	case *uint16:
		return EqualComparablePtr(vx, y)
	case *uint32:
		return EqualComparablePtr(vx, y)
	case *uint64:
		return EqualComparablePtr(vx, y)
	case *float32:
		return EqualComparablePtr(vx, y)
	case *float64:
		return EqualComparablePtr(vx, y)
	case *complex64:
		return EqualComparablePtr(vx, y)
	case *complex128:
		return EqualComparablePtr(vx, y)
	case *time.Duration:
		return EqualComparablePtr(vx, y)
	case *time.Time:
		return EqualComparablePtr(vx, y)
	case *uuid.UUID:
		return EqualComparablePtr(vx, y)
	case Equaler[T]:
		vy := any(y)

		if vx == nil || vy == nil {
			return vx == vy
		}

		return vx.Equal(y)
	case map[string]any:
		mapY, ok := any(y).(map[string]any)
		if !ok {
			return false
		}

		return EqualMap(vx, mapY, omitZero)
	case map[string]string:
		mapY, ok := any(y).(map[string]string)
		if !ok {
			return false
		}

		return EqualMap(vx, mapY, omitZero)
	case map[bool]any:
		mapY, ok := any(y).(map[bool]any)
		if !ok {
			return false
		}

		return EqualMap(vx, mapY, omitZero)
	case map[int]any:
		mapY, ok := any(y).(map[int]any)
		if !ok {
			return false
		}

		return EqualMap(vx, mapY, omitZero)
	case map[int8]any:
		mapY, ok := any(y).(map[int8]any)
		if !ok {
			return false
		}

		return EqualMap(vx, mapY, omitZero)
	case map[int16]any:
		mapY, ok := any(y).(map[int16]any)
		if !ok {
			return false
		}

		return EqualMap(vx, mapY, omitZero)
	case map[int32]any:
		mapY, ok := any(y).(map[int32]any)
		if !ok {
			return false
		}

		return EqualMap(vx, mapY, omitZero)
	case map[int64]any:
		mapY, ok := any(y).(map[int64]any)
		if !ok {
			return false
		}

		return EqualMap(vx, mapY, omitZero)
	case map[uint]any:
		mapY, ok := any(y).(map[uint]any)
		if !ok {
			return false
		}

		return EqualMap(vx, mapY, omitZero)
	case map[uint8]any:
		mapY, ok := any(y).(map[uint8]any)
		if !ok {
			return false
		}

		return EqualMap(vx, mapY, omitZero)
	case map[uint16]any:
		mapY, ok := any(y).(map[uint16]any)
		if !ok {
			return false
		}

		return EqualMap(vx, mapY, omitZero)
	case map[uint32]any:
		mapY, ok := any(y).(map[uint32]any)
		if !ok {
			return false
		}

		return EqualMap(vx, mapY, omitZero)
	case map[uint64]any:
		mapY, ok := any(y).(map[uint64]any)
		if !ok {
			return false
		}

		return EqualMap(vx, mapY, omitZero)
	case map[float32]any:
		mapY, ok := any(y).(map[float32]any)
		if !ok {
			return false
		}

		return EqualMap(vx, mapY, omitZero)
	case map[float64]any:
		mapY, ok := any(y).(map[float64]any)
		if !ok {
			return false
		}

		return EqualMap(vx, mapY, omitZero)
	case map[complex64]any:
		mapY, ok := any(y).(map[complex64]any)
		if !ok {
			return false
		}

		return EqualMap(vx, mapY, omitZero)
	case map[complex128]any:
		mapY, ok := any(y).(map[complex128]any)
		if !ok {
			return false
		}

		return EqualMap(vx, mapY, omitZero)
	case map[any]any:
		mapY, ok := any(y).(map[any]any)
		if !ok {
			return false
		}

		return EqualMap(vx, mapY, omitZero)
	case []bool:
		return EqualComparableSlice(vx, y, omitZero)
	case []string:
		return EqualComparableSlice(vx, y, omitZero)
	case []int:
		return EqualComparableSlice(vx, y, omitZero)
	case []int8:
		return EqualComparableSlice(vx, y, omitZero)
	case []int16:
		return EqualComparableSlice(vx, y, omitZero)
	case []int32:
		return EqualComparableSlice(vx, y, omitZero)
	case []int64:
		return EqualComparableSlice(vx, y, omitZero)
	case []uint:
		return EqualComparableSlice(vx, y, omitZero)
	case []uint8:
		return EqualComparableSlice(vx, y, omitZero)
	case []uint16:
		return EqualComparableSlice(vx, y, omitZero)
	case []uint32:
		return EqualComparableSlice(vx, y, omitZero)
	case []uint64:
		return EqualComparableSlice(vx, y, omitZero)
	case []float32:
		return EqualComparableSlice(vx, y, omitZero)
	case []float64:
		return EqualComparableSlice(vx, y, omitZero)
	case []complex64:
		return EqualComparableSlice(vx, y, omitZero)
	case []complex128:
		return EqualComparableSlice(vx, y, omitZero)
	case []time.Time:
		return EqualComparableSlice(vx, y, omitZero)
	case []time.Duration:
		return EqualComparableSlice(vx, y, omitZero)
	case []uuid.UUID:
		return EqualComparableSlice(vx, y, omitZero)
	case []any:
		sliceY, ok := any(y).([]any)
		if !ok {
			return false
		}

		return EqualSlice(vx, sliceY, omitZero)
	default:
		return reflect.DeepEqual(x, y)
	}
}

// EqualComparable checks if the y is comparable and equal x.
func EqualComparable[T comparable](x T, y any) bool {
	vy, ok := y.(T)
	if ok {
		return x == vy
	}

	py, ok := y.(*T)
	if ok {
		return py != nil && *py == x
	}

	return false
}

// EqualComparablePtr checks if the y is comparable and equal the pointer x.
func EqualComparablePtr[T comparable](x *T, y any) bool {
	if x == nil && y == nil {
		return true
	}

	if x == nil {
		return false
	}

	vy, ok := y.(T)
	if ok {
		return *x == vy
	}

	py, ok := y.(*T)
	if ok {
		return py != nil && *py == *x
	}

	return false
}

// EqualMap checks if both maps' elements are matched.
func EqualMap[K comparable, V any](mapA, mapB map[K]V, omitZero bool) bool {
	if mapA == nil || mapB == nil {
		if omitZero {
			return len(mapA) == len(mapB)
		}

		return mapA == nil && mapB == nil
	}

	// the both maps have the same pointer, they should equal.
	if reflect.ValueOf(mapA).UnsafePointer() == reflect.ValueOf(mapB).UnsafePointer() {
		return true
	}

	if len(mapA) != len(mapB) {
		return false
	}

	for key, valueA := range mapA {
		valueB, ok := mapB[key]
		if !ok {
			if omitZero {
				izA, ok := any(valueA).(IsZeroer)
				if ok && izA.IsZero() {
					continue
				}
			}

			return false
		}

		if !DeepEqual(valueA, valueB, omitZero) {
			return false
		}
	}

	return true
}

// EqualMapPointer checks if both maps' pointer elements are matched.
func EqualMapPointer[K comparable, V Equaler[V]]( //nolint:cyclop
	mapA, mapB map[K]*V,
	omitZero bool,
) bool {
	if len(mapA) != len(mapB) {
		return false
	}

	if (mapA == nil && mapB == nil) || len(mapA) == 0 {
		return true
	}

	for key, valueA := range mapA {
		valueB, ok := mapB[key]
		if !ok {
			if valueA == nil {
				continue
			}

			if omitZero && IsZero(*valueA) {
				continue
			}

			return false
		}

		if valueA == nil && valueB == nil {
			continue
		}

		if valueA == nil || valueB == nil {
			if omitZero && IsZeroPtr(valueA) && IsZeroPtr(valueB) {
				continue
			}

			return false
		}

		if !(*valueA).Equal(*valueB) {
			return false
		}
	}

	return true
}

// EqualComparableSlice checks if the y is comparable and equal x.
func EqualComparableSlice[T comparable](x []T, y any, omitZero bool) bool {
	vy, ok := y.([]T)
	if !ok {
		return false
	}

	return EqualSlice(x, vy, omitZero)
}
