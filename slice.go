package goutils

import (
	"cmp"
	"fmt"
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

// PtrToNumberSlice converts the pointer type of a number slice.
func PtrToNumberSlice[T1, T2 ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](
	inputs []*T1,
) ([]T2, error) {
	results := make([]T2, len(inputs))

	for i, value := range inputs {
		if value == nil {
			return nil, fmt.Errorf("element %d: %w", i, ErrNumberNull)
		}

		results[i] = T2(*value)
	}

	return results, nil
}

// SliceEqualSorted checks if both slices's elements are matched with sorted order.
func SliceEqualSorted[T cmp.Ordered](sliceA, sliceB []T) bool {
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
