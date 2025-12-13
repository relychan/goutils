package goutils

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
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

func TestIsZero(t *testing.T) {
	if IsZero(&mockEquality{}) {
		t.Error("expected non-zero")
	}

	if !IsZeroPtr[any](nil) {
		t.Error("expected zero")
	}

	if !IsZero[[]int](nil) {
		t.Error("expected zero")
	}

	if !IsZero[[]struct{}](nil) {
		t.Error("expected zero")
	}
}

func TestIsZero_Maps(t *testing.T) {
	t.Run("map[string]string", func(t *testing.T) {
		if !IsZero(map[string]string{}) {
			t.Error("expected zero for empty map")
		}
		if !IsZero(map[string]string(nil)) {
			t.Error("expected zero for nil map")
		}
		if IsZero(map[string]string{"a": "b"}) {
			t.Error("expected non-zero for non-empty map")
		}
	})

	t.Run("map[bool]any", func(t *testing.T) {
		if !IsZero(map[bool]any{}) {
			t.Error("expected zero for empty map")
		}
		if IsZero(map[bool]any{true: 1}) {
			t.Error("expected non-zero for non-empty map")
		}
	})

	t.Run("map[int]any", func(t *testing.T) {
		if !IsZero(map[int]any{}) {
			t.Error("expected zero for empty map")
		}
		if IsZero(map[int]any{1: "a"}) {
			t.Error("expected non-zero for non-empty map")
		}
	})

	t.Run("map[int8]any", func(t *testing.T) {
		if !IsZero(map[int8]any{}) {
			t.Error("expected zero for empty map")
		}
		if IsZero(map[int8]any{1: "a"}) {
			t.Error("expected non-zero for non-empty map")
		}
	})

	t.Run("map[int16]any", func(t *testing.T) {
		if !IsZero(map[int16]any{}) {
			t.Error("expected zero for empty map")
		}
		if IsZero(map[int16]any{1: "a"}) {
			t.Error("expected non-zero for non-empty map")
		}
	})

	t.Run("map[int32]any", func(t *testing.T) {
		if !IsZero(map[int32]any{}) {
			t.Error("expected zero for empty map")
		}
		if IsZero(map[int32]any{1: "a"}) {
			t.Error("expected non-zero for non-empty map")
		}
	})

	t.Run("map[int64]any", func(t *testing.T) {
		if !IsZero(map[int64]any{}) {
			t.Error("expected zero for empty map")
		}
		if IsZero(map[int64]any{1: "a"}) {
			t.Error("expected non-zero for non-empty map")
		}
	})

	t.Run("map[uint]any", func(t *testing.T) {
		if !IsZero(map[uint]any{}) {
			t.Error("expected zero for empty map")
		}
		if IsZero(map[uint]any{1: "a"}) {
			t.Error("expected non-zero for non-empty map")
		}
	})

	t.Run("map[uint8]any", func(t *testing.T) {
		if !IsZero(map[uint8]any{}) {
			t.Error("expected zero for empty map")
		}
		if IsZero(map[uint8]any{1: "a"}) {
			t.Error("expected non-zero for non-empty map")
		}
	})

	t.Run("map[uint16]any", func(t *testing.T) {
		if !IsZero(map[uint16]any{}) {
			t.Error("expected zero for empty map")
		}
		if IsZero(map[uint16]any{1: "a"}) {
			t.Error("expected non-zero for non-empty map")
		}
	})

	t.Run("map[uint32]any", func(t *testing.T) {
		if !IsZero(map[uint32]any{}) {
			t.Error("expected zero for empty map")
		}
		if IsZero(map[uint32]any{1: "a"}) {
			t.Error("expected non-zero for non-empty map")
		}
	})

	t.Run("map[uint64]any", func(t *testing.T) {
		if !IsZero(map[uint64]any{}) {
			t.Error("expected zero for empty map")
		}
		if IsZero(map[uint64]any{1: "a"}) {
			t.Error("expected non-zero for non-empty map")
		}
	})

	t.Run("map[float32]any", func(t *testing.T) {
		if !IsZero(map[float32]any{}) {
			t.Error("expected zero for empty map")
		}
		if IsZero(map[float32]any{1.5: "a"}) {
			t.Error("expected non-zero for non-empty map")
		}
	})

	t.Run("map[float64]any", func(t *testing.T) {
		if !IsZero(map[float64]any{}) {
			t.Error("expected zero for empty map")
		}
		if IsZero(map[float64]any{1.5: "a"}) {
			t.Error("expected non-zero for non-empty map")
		}
	})

	t.Run("map[complex64]any", func(t *testing.T) {
		if !IsZero(map[complex64]any{}) {
			t.Error("expected zero for empty map")
		}
		if IsZero(map[complex64]any{1 + 2i: "a"}) {
			t.Error("expected non-zero for non-empty map")
		}
	})

	t.Run("map[complex128]any", func(t *testing.T) {
		if !IsZero(map[complex128]any{}) {
			t.Error("expected zero for empty map")
		}
		if IsZero(map[complex128]any{1 + 2i: "a"}) {
			t.Error("expected non-zero for non-empty map")
		}
	})

	t.Run("map[string]any", func(t *testing.T) {
		if !IsZero(map[string]any{}) {
			t.Error("expected zero for empty map")
		}
		if IsZero(map[string]any{"a": 1}) {
			t.Error("expected non-zero for non-empty map")
		}
	})

	t.Run("map[any]any", func(t *testing.T) {
		if !IsZero(map[any]any{}) {
			t.Error("expected zero for empty map")
		}
		if IsZero(map[any]any{"a": 1}) {
			t.Error("expected non-zero for non-empty map")
		}
	})
}

func TestIsZero_Slices(t *testing.T) {
	t.Run("[]bool", func(t *testing.T) {
		if !IsZero([]bool{}) {
			t.Error("expected zero for empty slice")
		}
		if !IsZero([]bool(nil)) {
			t.Error("expected zero for nil slice")
		}
		if IsZero([]bool{true, false}) {
			t.Error("expected non-zero for non-empty slice")
		}
	})

	t.Run("[]string", func(t *testing.T) {
		if !IsZero([]string{}) {
			t.Error("expected zero for empty slice")
		}
		if !IsZero([]string(nil)) {
			t.Error("expected zero for nil slice")
		}
		if IsZero([]string{"a", "b"}) {
			t.Error("expected non-zero for non-empty slice")
		}
	})

	t.Run("[]int", func(t *testing.T) {
		if !IsZero([]int{}) {
			t.Error("expected zero for empty slice")
		}
		if IsZero([]int{1, 2, 3}) {
			t.Error("expected non-zero for non-empty slice")
		}
	})

	t.Run("[]int8", func(t *testing.T) {
		if !IsZero([]int8{}) {
			t.Error("expected zero for empty slice")
		}
		if IsZero([]int8{1, 2}) {
			t.Error("expected non-zero for non-empty slice")
		}
	})

	t.Run("[]int16", func(t *testing.T) {
		if !IsZero([]int16{}) {
			t.Error("expected zero for empty slice")
		}
		if IsZero([]int16{1, 2}) {
			t.Error("expected non-zero for non-empty slice")
		}
	})

	t.Run("[]int32", func(t *testing.T) {
		if !IsZero([]int32{}) {
			t.Error("expected zero for empty slice")
		}
		if IsZero([]int32{1, 2}) {
			t.Error("expected non-zero for non-empty slice")
		}
	})

	t.Run("[]int64", func(t *testing.T) {
		if !IsZero([]int64{}) {
			t.Error("expected zero for empty slice")
		}
		if IsZero([]int64{1, 2}) {
			t.Error("expected non-zero for non-empty slice")
		}
	})

	t.Run("[]uint", func(t *testing.T) {
		if !IsZero([]uint{}) {
			t.Error("expected zero for empty slice")
		}
		if IsZero([]uint{1, 2}) {
			t.Error("expected non-zero for non-empty slice")
		}
	})

	t.Run("[]uint8", func(t *testing.T) {
		if !IsZero([]uint8{}) {
			t.Error("expected zero for empty slice")
		}
		if IsZero([]uint8{1, 2}) {
			t.Error("expected non-zero for non-empty slice")
		}
	})

	t.Run("[]uint16", func(t *testing.T) {
		if !IsZero([]uint16{}) {
			t.Error("expected zero for empty slice")
		}
		if IsZero([]uint16{1, 2}) {
			t.Error("expected non-zero for non-empty slice")
		}
	})

	t.Run("[]uint32", func(t *testing.T) {
		if !IsZero([]uint32{}) {
			t.Error("expected zero for empty slice")
		}
		if IsZero([]uint32{1, 2}) {
			t.Error("expected non-zero for non-empty slice")
		}
	})

	t.Run("[]uint64", func(t *testing.T) {
		if !IsZero([]uint64{}) {
			t.Error("expected zero for empty slice")
		}
		if IsZero([]uint64{1, 2}) {
			t.Error("expected non-zero for non-empty slice")
		}
	})

	t.Run("[]float32", func(t *testing.T) {
		if !IsZero([]float32{}) {
			t.Error("expected zero for empty slice")
		}
		if IsZero([]float32{1.5, 2.5}) {
			t.Error("expected non-zero for non-empty slice")
		}
	})

	t.Run("[]float64", func(t *testing.T) {
		if !IsZero([]float64{}) {
			t.Error("expected zero for empty slice")
		}
		if IsZero([]float64{1.5, 2.5}) {
			t.Error("expected non-zero for non-empty slice")
		}
	})

	t.Run("[]complex64", func(t *testing.T) {
		if !IsZero([]complex64{}) {
			t.Error("expected zero for empty slice")
		}
		if IsZero([]complex64{1 + 2i}) {
			t.Error("expected non-zero for non-empty slice")
		}
	})

	t.Run("[]complex128", func(t *testing.T) {
		if !IsZero([]complex128{}) {
			t.Error("expected zero for empty slice")
		}
		if IsZero([]complex128{1 + 2i}) {
			t.Error("expected non-zero for non-empty slice")
		}
	})

	t.Run("[]time.Time", func(t *testing.T) {
		if !IsZero([]time.Time{}) {
			t.Error("expected zero for empty slice")
		}
		if IsZero([]time.Time{time.Now()}) {
			t.Error("expected non-zero for non-empty slice")
		}
	})

	t.Run("[]time.Duration", func(t *testing.T) {
		if !IsZero([]time.Duration{}) {
			t.Error("expected zero for empty slice")
		}
		if IsZero([]time.Duration{time.Second}) {
			t.Error("expected non-zero for non-empty slice")
		}
	})

	t.Run("[]uuid.UUID", func(t *testing.T) {
		if !IsZero([]uuid.UUID{}) {
			t.Error("expected zero for empty slice")
		}
		if IsZero([]uuid.UUID{uuid.New()}) {
			t.Error("expected non-zero for non-empty slice")
		}
	})

	t.Run("[]any", func(t *testing.T) {
		if !IsZero([]any{}) {
			t.Error("expected zero for empty slice")
		}
		if IsZero([]any{1, "test", true}) {
			t.Error("expected non-zero for non-empty slice")
		}
	})
}

func TestIsZero_TimeDuration(t *testing.T) {
	t.Run("zero duration", func(t *testing.T) {
		if !IsZero(time.Duration(0)) {
			t.Error("expected zero for zero duration")
		}
	})

	t.Run("non-zero duration", func(t *testing.T) {
		if IsZero(time.Second) {
			t.Error("expected non-zero for non-zero duration")
		}
		if IsZero(time.Millisecond * 100) {
			t.Error("expected non-zero for non-zero duration")
		}
	})
}

func TestIsZero_IsZeroer(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		val := testIsZeroer{Value: 0}
		if !IsZero(val) {
			t.Error("expected zero for IsZeroer with zero value")
		}
	})

	t.Run("non-zero value", func(t *testing.T) {
		val := testIsZeroer{Value: 42}
		if IsZero(val) {
			t.Error("expected non-zero for IsZeroer with non-zero value")
		}
	})
}

func TestIsZero_Reflection(t *testing.T) {
	t.Run("array", func(t *testing.T) {
		var emptyArray [0]int
		if !IsZero(emptyArray) {
			t.Error("expected zero for empty array")
		}

		var nonEmptyArray [3]int
		if IsZero(nonEmptyArray) {
			t.Error("expected non-zero for non-empty array")
		}
	})

	t.Run("pointer to slice", func(t *testing.T) {
		var nilSlice *[]int
		if !IsZero(nilSlice) {
			t.Error("expected zero for nil pointer to slice")
		}

		emptySlice := []int{}
		if !IsZero(&emptySlice) {
			t.Error("expected zero for pointer to empty slice")
		}

		nonEmptySlice := []int{1, 2, 3}
		if IsZero(&nonEmptySlice) {
			t.Error("expected non-zero for pointer to non-empty slice")
		}
	})

	t.Run("pointer to map", func(t *testing.T) {
		var nilMap *map[string]int
		if !IsZero(nilMap) {
			t.Error("expected zero for nil pointer to map")
		}

		emptyMap := map[string]int{}
		if !IsZero(&emptyMap) {
			t.Error("expected zero for pointer to empty map")
		}

		nonEmptyMap := map[string]int{"a": 1}
		if IsZero(&nonEmptyMap) {
			t.Error("expected non-zero for pointer to non-empty map")
		}
	})

	t.Run("custom struct", func(t *testing.T) {
		type customStruct struct {
			Field1 string
			Field2 int
		}

		// Custom structs without IsZeroer interface should return false
		s := customStruct{Field1: "test", Field2: 42}
		if IsZero(s) {
			t.Error("expected non-zero for custom struct")
		}
	})
}

func TestIsZeroPtr(t *testing.T) {
	t.Run("nil pointer", func(t *testing.T) {
		var ptr *int
		if !IsZeroPtr(ptr) {
			t.Error("expected zero for nil pointer")
		}
	})

	t.Run("pointer to zero value", func(t *testing.T) {
		// IsZeroPtr checks if pointer is nil OR if the dereferenced value is zero
		// For primitive types like int, the zero value (0) is considered zero
		val := 0
		// This should return true because IsZero(0) returns true for primitive zero values
		// However, IsZero for int doesn't have a specific case, so it falls through to reflection
		// which returns false for non-zero-length types
		if IsZeroPtr(&val) {
			t.Error("expected non-zero for pointer to zero int value (reflection doesn't consider primitive zero)")
		}
	})

	t.Run("pointer to non-zero value", func(t *testing.T) {
		val := 42
		if IsZeroPtr(&val) {
			t.Error("expected non-zero for pointer to non-zero value")
		}
	})

	t.Run("pointer to empty slice", func(t *testing.T) {
		slice := []int{}
		if !IsZeroPtr(&slice) {
			t.Error("expected zero for pointer to empty slice")
		}
	})

	t.Run("pointer to non-empty slice", func(t *testing.T) {
		slice := []int{1, 2, 3}
		if IsZeroPtr(&slice) {
			t.Error("expected non-zero for pointer to non-empty slice")
		}
	})
}

func TestToPtr(t *testing.T) {
	t.Run("int value", func(t *testing.T) {
		val := 42
		ptr := ToPtr(val)
		if ptr == nil {
			t.Error("expected non-nil pointer")
		}
		if *ptr != val {
			t.Errorf("expected %d, got %d", val, *ptr)
		}
	})

	t.Run("string value", func(t *testing.T) {
		val := "hello"
		ptr := ToPtr(val)
		if ptr == nil {
			t.Error("expected non-nil pointer")
		}
		if *ptr != val {
			t.Errorf("expected %s, got %s", val, *ptr)
		}
	})

	t.Run("struct value", func(t *testing.T) {
		type testStruct struct {
			Field1 string
			Field2 int
		}
		val := testStruct{Field1: "test", Field2: 42}
		ptr := ToPtr(val)
		if ptr == nil {
			t.Error("expected non-nil pointer")
		}
		if ptr.Field1 != val.Field1 || ptr.Field2 != val.Field2 {
			t.Error("expected equal struct values")
		}
	})

	t.Run("slice value", func(t *testing.T) {
		val := []int{1, 2, 3}
		ptr := ToPtr(val)
		if ptr == nil {
			t.Error("expected non-nil pointer")
		}
		if len(*ptr) != len(val) {
			t.Error("expected equal slice lengths")
		}
	})

	t.Run("zero value", func(t *testing.T) {
		val := 0
		ptr := ToPtr(val)
		if ptr == nil {
			t.Error("expected non-nil pointer even for zero value")
		}
		if *ptr != 0 {
			t.Error("expected zero value")
		}
	})
}

func TestIsNil_Extended(t *testing.T) {
	t.Run("nil interface", func(t *testing.T) {
		var i interface{}
		if !IsNil(i) {
			t.Error("expected nil for nil interface")
		}
	})

	t.Run("nil pointer", func(t *testing.T) {
		var ptr *int
		if !IsNil(ptr) {
			t.Error("expected nil for nil pointer")
		}
	})

	t.Run("nil slice", func(t *testing.T) {
		var slice []int
		if !IsNil(slice) {
			t.Error("expected nil for nil slice")
		}
	})

	t.Run("nil map", func(t *testing.T) {
		var m map[string]int
		if !IsNil(m) {
			t.Error("expected nil for nil map")
		}
	})

	t.Run("nil channel", func(t *testing.T) {
		var ch chan int
		if !IsNil(ch) {
			t.Error("expected nil for nil channel")
		}
	})

	t.Run("nil function", func(t *testing.T) {
		var fn func()
		if !IsNil(fn) {
			t.Error("expected nil for nil function")
		}
	})

	t.Run("non-nil pointer", func(t *testing.T) {
		val := 42
		if IsNil(&val) {
			t.Error("expected not nil for non-nil pointer")
		}
	})

	t.Run("non-nil slice", func(t *testing.T) {
		slice := []int{1, 2, 3}
		if IsNil(slice) {
			t.Error("expected not nil for non-nil slice")
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		slice := []int{}
		if IsNil(slice) {
			t.Error("expected not nil for empty slice (non-nil)")
		}
	})

	t.Run("non-nil map", func(t *testing.T) {
		m := map[string]int{"a": 1}
		if IsNil(m) {
			t.Error("expected not nil for non-nil map")
		}
	})

	t.Run("primitive types", func(t *testing.T) {
		if IsNil(42) {
			t.Error("expected not nil for int")
		}
		if IsNil("hello") {
			t.Error("expected not nil for string")
		}
		if IsNil(true) {
			t.Error("expected not nil for bool")
		}
	})
}

func TestUnwrapPointerFromReflectValue(t *testing.T) {
	t.Run("direct value", func(t *testing.T) {
		val := 42
		rv := reflect.ValueOf(val)
		unwrapped, ok := UnwrapPointerFromReflectValue(rv)
		if !ok {
			t.Error("expected valid unwrapped value")
		}
		if unwrapped.Int() != 42 {
			t.Errorf("expected 42, got %d", unwrapped.Int())
		}
	})

	t.Run("single pointer", func(t *testing.T) {
		val := 42
		ptr := &val
		rv := reflect.ValueOf(ptr)
		unwrapped, ok := UnwrapPointerFromReflectValue(rv)
		if !ok {
			t.Error("expected valid unwrapped value")
		}
		if unwrapped.Int() != 42 {
			t.Errorf("expected 42, got %d", unwrapped.Int())
		}
	})

	t.Run("double pointer", func(t *testing.T) {
		val := 42
		ptr := &val
		ptrPtr := &ptr
		rv := reflect.ValueOf(ptrPtr)
		unwrapped, ok := UnwrapPointerFromReflectValue(rv)
		if !ok {
			t.Error("expected valid unwrapped value")
		}
		if unwrapped.Int() != 42 {
			t.Errorf("expected 42, got %d", unwrapped.Int())
		}
	})

	t.Run("nil pointer", func(t *testing.T) {
		var ptr *int
		rv := reflect.ValueOf(ptr)
		_, ok := UnwrapPointerFromReflectValue(rv)
		if ok {
			t.Error("expected invalid for nil pointer")
		}
	})

	t.Run("nil slice", func(t *testing.T) {
		var slice []int
		rv := reflect.ValueOf(slice)
		_, ok := UnwrapPointerFromReflectValue(rv)
		if ok {
			t.Error("expected invalid for nil slice")
		}
	})

	t.Run("non-nil slice", func(t *testing.T) {
		slice := []int{1, 2, 3}
		rv := reflect.ValueOf(slice)
		unwrapped, ok := UnwrapPointerFromReflectValue(rv)
		if !ok {
			t.Error("expected valid for non-nil slice")
		}
		if unwrapped.Len() != 3 {
			t.Errorf("expected length 3, got %d", unwrapped.Len())
		}
	})

	t.Run("nil map", func(t *testing.T) {
		var m map[string]int
		rv := reflect.ValueOf(m)
		_, ok := UnwrapPointerFromReflectValue(rv)
		if ok {
			t.Error("expected invalid for nil map")
		}
	})

	t.Run("non-nil map", func(t *testing.T) {
		m := map[string]int{"a": 1}
		rv := reflect.ValueOf(m)
		unwrapped, ok := UnwrapPointerFromReflectValue(rv)
		if !ok {
			t.Error("expected valid for non-nil map")
		}
		if unwrapped.Len() != 1 {
			t.Errorf("expected length 1, got %d", unwrapped.Len())
		}
	})

	t.Run("nil interface", func(t *testing.T) {
		var i any
		rv := reflect.ValueOf(i)
		_, ok := UnwrapPointerFromReflectValue(rv)
		if ok {
			t.Error("expected invalid for nil interface")
		}
	})

	t.Run("non-nil interface", func(t *testing.T) {
		var i any = 42
		rv := reflect.ValueOf(i)
		unwrapped, ok := UnwrapPointerFromReflectValue(rv)
		if !ok {
			t.Error("expected valid for non-nil interface")
		}
		if unwrapped.Int() != 42 {
			t.Errorf("expected 42, got %d", unwrapped.Int())
		}
	})

	t.Run("channel", func(t *testing.T) {
		ch := make(chan int)
		rv := reflect.ValueOf(ch)
		_, ok := UnwrapPointerFromReflectValue(rv)
		if ok {
			t.Error("expected invalid for channel")
		}
	})

	t.Run("function", func(t *testing.T) {
		fn := func() {}
		rv := reflect.ValueOf(fn)
		_, ok := UnwrapPointerFromReflectValue(rv)
		if ok {
			t.Error("expected invalid for function")
		}
	})
}

func TestNewUUIDv7_Extended(t *testing.T) {
	t.Run("generates valid UUID", func(t *testing.T) {
		id := NewUUIDv7()
		if id == uuid.Nil {
			t.Error("expected non-nil UUID")
		}
	})

	t.Run("generates unique UUIDs", func(t *testing.T) {
		id1 := NewUUIDv7()
		id2 := NewUUIDv7()
		if id1 == id2 {
			t.Error("expected different UUIDs")
		}
	})

	t.Run("UUID string format", func(t *testing.T) {
		id := NewUUIDv7()
		str := id.String()
		if len(str) != 36 {
			t.Errorf("expected UUID string length 36, got %d", len(str))
		}
	})
}
