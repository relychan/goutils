package goutils

import (
	"cmp"
	"fmt"
	"reflect"
	"slices"
)

// ToAnySlice converts the a typed slice to any slice.
func ToAnySlice[T any](inputs []T) []any {
	results := make([]any, len(inputs))

	for i, value := range inputs {
		results[i] = value
	}

	return results
}

// ToNumberSlice converts the element type of a number slice.
func ToNumberSlice[T1, T2 ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](
	inputs []T1,
) []T2 {
	results := make([]T2, len(inputs))

	for i, value := range inputs {
		results[i] = T2(value)
	}

	return results
}

// PtrToNumberSlice converts a slice of number pointers to a slice of values.
// Returns nil, nil if inputs is nil. Returns an error if any element in the slice is nil.
func PtrToNumberSlice[T1, T2 ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](
	inputs []*T1,
) ([]T2, error) {
	if inputs == nil {
		return nil, nil
	}

	results := make([]T2, len(inputs))

	for i, value := range inputs {
		if value == nil {
			return nil, fmt.Errorf("element %d: %w", i, ErrNumberNull)
		}

		results[i] = T2(*value)
	}

	return results, nil
}

// EqualSlice checks if both slices's elements are matched.
func EqualSlice[T any](sliceA, sliceB []T, omitZero bool) bool {
	if sliceA == nil || sliceB == nil {
		if omitZero {
			return len(sliceA) == len(sliceB)
		}

		return sliceA == nil && sliceB == nil
	}

	// the both maps have the same pointer, they should equal.
	if reflect.ValueOf(sliceA).UnsafePointer() == reflect.ValueOf(sliceB).UnsafePointer() {
		return true
	}

	if len(sliceA) != len(sliceB) {
		return false
	}

	for i, a := range sliceA {
		b := sliceB[i]
		if !DeepEqual(a, b, false) {
			return false
		}
	}

	return true
}

// EqualSlicePtr checks if both slices's pointer elements are matched.
func EqualSlicePtr[T Equaler[T]](sliceA, sliceB []*T) bool {
	if len(sliceA) != len(sliceB) {
		return false
	}

	if len(sliceA) == 0 {
		return true
	}

	for i, a := range sliceA {
		b := sliceB[i]

		if a == nil && b == nil {
			continue
		}

		if a == nil || b == nil {
			return false
		}

		if !(*a).Equal(*b) {
			return false
		}
	}

	return true
}

// EqualSliceSorted checks if both slices's elements are matched with sorted order.
func EqualSliceSorted[T cmp.Ordered](sliceA, sliceB []T) bool {
	if len(sliceA) != len(sliceB) {
		return false
	}

	if len(sliceA) == 0 {
		return true
	}

	slices.Sort(sliceA)
	slices.Sort(sliceB)

	for i, a := range sliceA {
		b := sliceB[i]
		if a != b {
			return false
		}
	}

	return true
}

// Map applies a function to each element of a slice and returns a new slice
// with the results.
// T is the type of the input slice elements.
// M is the type of the output slice elements.
func Map[T, M any](input []T, f func(T) M) []M {
	output := make([]M, len(input))

	for i, v := range input {
		output[i] = f(v)
	}

	return output
}
