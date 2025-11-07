package goutils

import (
	"cmp"
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
