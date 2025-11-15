package goutils

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestToAnySlice(t *testing.T) {
	sliceInt := []int{0, 1, 2, 3, 4}
	anySlice := ToAnySlice(sliceInt)

	if fmt.Sprint(sliceInt) != fmt.Sprint(anySlice) {
		t.Fatalf("expected equal")
	}
}

func TestSliceEqualSorted(t *testing.T) {
	sortedInts := []int{0, 1, 2, 3, 4}
	unsortedInts := []int{1, 3, 0, 2, 4}

	if !SliceEqualSorted(sortedInts, unsortedInts) {
		t.Fatalf("expected equal")
	}

	if SliceEqualSorted([]int{0}, unsortedInts) {
		t.Fatalf("expected not equal")
	}

	if !SliceEqualSorted([]int{}, []int{}) {
		t.Fatalf("expected equal")
	}

	if SliceEqualSorted([]int{0, 1, 2}, []int{0, 2, 3}) {
		t.Fatalf("expected not equal")
	}
}

func TestMapSlice(t *testing.T) {
	inputs := []string{"Hello", "WORLD"}

	if !reflect.DeepEqual(Map(inputs, strings.ToLower), ToLowerStrings(inputs)) {
		t.Fatal("lowercase not equal")
	}

	if !reflect.DeepEqual(Map(inputs, strings.ToUpper), ToUpperStrings(inputs)) {
		t.Fatal("uppercase not equal")
	}
}
