// Copyright 2026 RelyChan Pte. Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

	if !EqualSliceSorted(sortedInts, unsortedInts) {
		t.Fatalf("expected equal")
	}

	if EqualSliceSorted([]int{0}, unsortedInts) {
		t.Fatalf("expected not equal")
	}

	if !EqualSliceSorted([]int{}, []int{}) {
		t.Fatalf("expected equal")
	}

	if EqualSliceSorted([]int{0, 1, 2}, []int{0, 2, 3}) {
		t.Fatalf("expected not equal")
	}
}

func TestMapSlice(t *testing.T) {
	inputs := []string{"Hello", "WORLD"}

	if !reflect.DeepEqual(Map(inputs, strings.ToLower), []string{"hello", "world"}) {
		t.Fatal("lowercase not equal")
	}
}

func TestToNumberSlice(t *testing.T) {
	intSlice := []int{1, 2, 3}
	floatSlice := ToNumberSlice[int, float64](intSlice)
	if !reflect.DeepEqual(floatSlice, []float64{1.0, 2.0, 3.0}) {
		t.Fatal("conversion failed")
	}
}

func TestPtrToNumberSlice(t *testing.T) {
	intSlice := []*int{new(1), new(2)}
	result, err := PtrToNumberSlice[int, float64](intSlice)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result, []float64{1.0, 2.0}) {
		t.Fatal("conversion failed")
	}

	// Test nil element
	nilSlice := []*int{nil}
	_, err = PtrToNumberSlice[int, float64](nilSlice)
	if err == nil {
		t.Fatal("expected error for nil element")
	}
}
