package goutils

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestUUIDv7(t *testing.T) {
	value := NewUUIDv7()
	ptr := ToPtr(value)

	if value.String() != ptr.String() {
		t.Error("expected equal")
		t.FailNow()
	}
}

func TestIsNil(t *testing.T) {
	var jsonValue any

	err := json.Unmarshal([]byte("null"), &jsonValue)
	if err != nil {
		t.Error(err)
	}

	if !IsNil(jsonValue) {
		t.Errorf("expected nil, got: %v", jsonValue)
	}

	if IsNil(any((*int)(ToPtr(1)))) {
		t.Error("expected not nil, got: nil")
	}
}

type mockEquality struct{}

func (m mockEquality) Equal(target mockEquality) bool {
	return m == target
}

func TestEqualPtr(t *testing.T) {
	if !EqualPtr(&mockEquality{}, &mockEquality{}) {
		t.Error("expected equal")
	}

	if EqualPtr(&mockEquality{}, nil) {
		t.Error("expected not equal")
	}

	if !EqualPtr[mockEquality](nil, nil) {
		t.Error("expected equal")
	}
}

type equalStruct struct {
	Value int
}

var complexObject = map[string]any{
	"foo":     "bar",
	"boolean": false,
	"number":  10,
	"ints":    []int{1, 2, 3, 4},
	"object": map[int]any{
		1: "test",
		2: false,
		3: []int{1, 2, 3, 4},
		4: map[string]any{
			"foo": "bar",
		},
	},
	"slice": []any{
		map[string]any{
			"foo":     "bar",
			"boolean": false,
			"number":  10,
			"object": map[int]any{
				1: "test",
				2: false,
				3: []int{1, 2, 3, 4},
				4: map[string]any{
					"foo": "bar",
				},
			},
		},
	},
	"equaler": equalStruct{
		Value: 1,
	},
}

var complexObject2 = map[string]any{
	"foo":     "bar",
	"boolean": false,
	"number":  10,
	"ints":    []int{1, 2, 3, 4},
	"object": map[int]any{
		1: "test",
		2: false,
		3: []int{1, 2, 3, 4},
		4: map[string]any{
			"foo": "bar",
		},
	},
	"slice": []any{
		map[string]any{
			"foo":     "bar",
			"boolean": false,
			"number":  10,
			"object": map[int]any{
				1: "test",
				2: false,
				3: []int{1, 2, 3, 4},
				4: map[string]any{
					"foo": "bar",
				},
			},
		},
	},
	"equaler": equalStruct{
		Value: 1,
	},
}

// BenchmarkDeepEqual-11    	  320522	      3720 ns/op	    5752 B/op	      64 allocs/op
func BenchmarkDeepEqual(b *testing.B) {
	for b.Loop() {
		reflect.DeepEqual(complexObject, complexObject2)
	}
}

// BenchmarkDeepEqual2-11    	  758126	      1556 ns/op	     416 B/op	      30 allocs/op
func BenchmarkDeepEqual2(b *testing.B) {
	for b.Loop() {
		DeepEqual(complexObject, complexObject2, false)
	}
}
