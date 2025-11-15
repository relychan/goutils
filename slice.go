package goutils

import (
	"cmp"
	"slices"
	"strings"
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

// ToLowerStrings transform string elements in the slice to lowercase.
func ToLowerStrings(inputs []string) []string {
	output := make([]string, len(inputs))

	for i, v := range inputs {
		output[i] = strings.ToLower(v)
	}

	return output
}

// ToUpperStrings transform string elements in the slice to UPPERCASE.
func ToUpperStrings(inputs []string) []string {
	output := make([]string, len(inputs))

	for i, v := range inputs {
		output[i] = strings.ToUpper(v)
	}

	return output
}
