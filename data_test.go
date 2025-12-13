package goutils

import (
	"encoding/json"
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
