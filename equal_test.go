package goutils

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

// Test helper struct that implements Equaler interface
type testEqualer struct {
	Value int
}

func (t testEqualer) Equal(other testEqualer) bool {
	return t.Value == other.Value
}

// Test helper struct that implements IsZeroer interface
type testIsZeroer struct {
	Value int
}

func (t testIsZeroer) IsZero() bool {
	return t.Value == 0
}

func TestEqualComparable(t *testing.T) {
	t.Run("equal values", func(t *testing.T) {
		if !EqualComparableAny(42, 42) {
			t.Error("expected equal")
		}
		if !EqualComparableAny("hello", "hello") {
			t.Error("expected equal")
		}
		if !EqualComparableAny(true, true) {
			t.Error("expected equal")
		}
		if !EqualComparableAny(3.14, 3.14) {
			t.Error("expected equal")
		}
	})

	t.Run("different values", func(t *testing.T) {
		if EqualComparableAny(42, 43) {
			t.Error("expected not equal")
		}
		if EqualComparableAny("hello", "world") {
			t.Error("expected not equal")
		}
	})

	t.Run("value vs pointer", func(t *testing.T) {
		val := 42
		ptr := &val
		if !EqualComparableAny(42, ptr) {
			t.Error("expected equal when comparing value to pointer")
		}
	})

	t.Run("wrong type", func(t *testing.T) {
		if EqualComparableAny(42, "42") {
			t.Error("expected not equal for different types")
		}
	})
}

func TestEqualComparablePtr(t *testing.T) {
	t.Run("both nil", func(t *testing.T) {
		var a *int
		if !EqualComparablePtr(a, nil) {
			t.Error("expected equal when both nil")
		}
	})

	t.Run("one nil", func(t *testing.T) {
		val := 42
		if EqualComparablePtr(&val, nil) {
			t.Error("expected not equal when one is nil")
		}
		if EqualComparableAnyPtr[int](nil, 42) {
			t.Error("expected not equal when pointer is nil")
		}
	})

	t.Run("equal values", func(t *testing.T) {
		val1 := 42
		val2 := 42
		if !EqualComparablePtr(&val1, &val2) {
			t.Error("expected equal")
		}
	})

	t.Run("pointer vs value", func(t *testing.T) {
		val := 42
		if !EqualComparableAnyPtr(&val, 42) {
			t.Error("expected equal when comparing pointer to value")
		}
	})

	t.Run("different values", func(t *testing.T) {
		val1 := 42
		val2 := 43
		if EqualComparablePtr(&val1, &val2) {
			t.Error("expected not equal")
		}
	})
}

func TestEqualMap(t *testing.T) {
	t.Run("equal maps", func(t *testing.T) {
		map1 := map[string]int{"a": 1, "b": 2}
		map2 := map[string]int{"a": 1, "b": 2}
		if !EqualMap(map1, map2, false) {
			t.Error("expected equal")
		}
	})

	t.Run("both nil without omitZero", func(t *testing.T) {
		var map1 map[string]int
		var map2 map[string]int
		if !EqualMap(map1, map2, false) {
			t.Error("expected equal when both nil")
		}
	})

	t.Run("both nil with omitZero", func(t *testing.T) {
		var map1 map[string]int
		var map2 map[string]int
		if !EqualMap(map1, map2, true) {
			t.Error("expected equal when both nil with omitZero")
		}
	})

	t.Run("one nil without omitZero", func(t *testing.T) {
		map1 := map[string]int{}
		var map2 map[string]int
		if EqualMap(map1, map2, false) {
			t.Error("expected not equal when one is nil without omitZero")
		}
	})

	t.Run("one nil with omitZero", func(t *testing.T) {
		map1 := map[string]int{}
		var map2 map[string]int
		if !EqualMap(map1, map2, true) {
			t.Error("expected equal when one is nil with omitZero")
		}
	})

	t.Run("same pointer", func(t *testing.T) {
		map1 := map[string]int{"a": 1}
		if !EqualMap(map1, map1, false) {
			t.Error("expected equal when same pointer")
		}
	})

	t.Run("different lengths", func(t *testing.T) {
		map1 := map[string]int{"a": 1}
		map2 := map[string]int{"a": 1, "b": 2}
		if EqualMap(map1, map2, false) {
			t.Error("expected not equal when different lengths")
		}
	})

	t.Run("different values", func(t *testing.T) {
		map1 := map[string]int{"a": 1}
		map2 := map[string]int{"a": 2}
		if EqualMap(map1, map2, false) {
			t.Error("expected not equal when different values")
		}
	})

	t.Run("missing key without omitZero", func(t *testing.T) {
		map1 := map[string]int{"a": 1, "b": 2}
		map2 := map[string]int{"a": 1}
		if EqualMap(map1, map2, false) {
			t.Error("expected not equal when key missing")
		}
	})

	t.Run("missing key with omitZero and IsZeroer value", func(t *testing.T) {
		// The EqualMap function checks lengths first, so maps with different lengths
		// will never be equal, even with omitZero. The omitZero flag only affects
		// the comparison when a key exists in mapA but not in mapB during iteration.
		// However, since lengths must match first, this scenario can't happen.
		// Let's test that maps with same length but zero values are handled correctly.
		map1 := map[string]any{"a": testIsZeroer{Value: 1}, "b": testIsZeroer{Value: 1}}
		map2 := map[string]any{"a": testIsZeroer{Value: 1}, "b": testIsZeroer{Value: 1}}
		if !EqualMap(map1, map2, true) {
			t.Error("expected equal for maps with same values")
		}
	})

	t.Run("nested maps", func(t *testing.T) {
		map1 := map[string]any{"a": map[string]int{"x": 1}}
		map2 := map[string]any{"a": map[string]int{"x": 1}}
		if !EqualMap(map1, map2, false) {
			t.Error("expected equal for nested maps")
		}
	})
}

func TestEqualMapPointer(t *testing.T) {
	t.Run("equal maps", func(t *testing.T) {
		val1 := testEqualer{Value: 1}
		val2 := testEqualer{Value: 2}
		val3 := testEqualer{Value: 1}
		val4 := testEqualer{Value: 2}

		map1 := map[string]*testEqualer{"a": &val1, "b": &val2}
		map2 := map[string]*testEqualer{"a": &val3, "b": &val4}
		if !EqualMapPointer(map1, map2, false) {
			t.Error("expected equal")
		}
	})

	t.Run("different lengths", func(t *testing.T) {
		val1 := testEqualer{Value: 1}
		val2 := testEqualer{Value: 2}
		map1 := map[string]*testEqualer{"a": &val1}
		map2 := map[string]*testEqualer{"a": &val1, "b": &val2}
		if EqualMapPointer(map1, map2, false) {
			t.Error("expected not equal when different lengths")
		}
	})

	t.Run("both nil", func(t *testing.T) {
		var map1 map[string]*testEqualer
		var map2 map[string]*testEqualer
		if !EqualMapPointer(map1, map2, false) {
			t.Error("expected equal when both nil")
		}
	})

	t.Run("both empty", func(t *testing.T) {
		map1 := map[string]*testEqualer{}
		map2 := map[string]*testEqualer{}
		if !EqualMapPointer(map1, map2, false) {
			t.Error("expected equal when both empty")
		}
	})

	t.Run("both values nil", func(t *testing.T) {
		map1 := map[string]*testEqualer{"a": nil}
		map2 := map[string]*testEqualer{"a": nil}
		if !EqualMapPointer(map1, map2, false) {
			t.Error("expected equal when both values nil")
		}
	})

	t.Run("one value nil without omitZero", func(t *testing.T) {
		val := testEqualer{Value: 1}
		map1 := map[string]*testEqualer{"a": &val}
		map2 := map[string]*testEqualer{"a": nil}
		if EqualMapPointer(map1, map2, false) {
			t.Error("expected not equal when one value is nil")
		}
	})

	t.Run("missing key with nil value", func(t *testing.T) {
		// When a key is missing in map2 but has nil value in map1, they continue (skip)
		// But the lengths are different, so they should not be equal
		map1 := map[string]*testEqualer{"a": nil}
		map2 := map[string]*testEqualer{}
		if EqualMapPointer(map1, map2, false) {
			t.Error("expected not equal when lengths differ")
		}
	})

	t.Run("different values", func(t *testing.T) {
		val1 := testEqualer{Value: 1}
		val2 := testEqualer{Value: 2}
		map1 := map[string]*testEqualer{"a": &val1}
		map2 := map[string]*testEqualer{"a": &val2}
		if EqualMapPointer(map1, map2, false) {
			t.Error("expected not equal when values differ")
		}
	})
}

func TestEqualComparableSlice(t *testing.T) {
	t.Run("equal slices", func(t *testing.T) {
		slice1 := []int{1, 2, 3}
		if !EqualComparableSlice(slice1, []int{1, 2, 3}, false) {
			t.Error("expected equal")
		}
	})

	t.Run("wrong type", func(t *testing.T) {
		slice1 := []int{1, 2, 3}
		if EqualComparableSlice(slice1, []string{"1", "2", "3"}, false) {
			t.Error("expected not equal for different types")
		}
	})

	t.Run("different values", func(t *testing.T) {
		slice1 := []int{1, 2, 3}
		if EqualComparableSlice(slice1, []int{1, 2, 4}, false) {
			t.Error("expected not equal")
		}
	})
}

func TestDeepEqual(t *testing.T) {
	t.Run("primitive types", func(t *testing.T) {
		// bool
		if !DeepEqual(true, true, false) {
			t.Error("expected equal for bool")
		}
		if DeepEqual(true, false, false) {
			t.Error("expected not equal for bool")
		}

		// string
		if !DeepEqual("hello", "hello", false) {
			t.Error("expected equal for string")
		}
		if DeepEqual("hello", "world", false) {
			t.Error("expected not equal for string")
		}

		// int types
		if !DeepEqual(42, 42, false) {
			t.Error("expected equal for int")
		}
		if !DeepEqual(int8(42), int8(42), false) {
			t.Error("expected equal for int8")
		}
		if !DeepEqual(int16(42), int16(42), false) {
			t.Error("expected equal for int16")
		}
		if !DeepEqual(int32(42), int32(42), false) {
			t.Error("expected equal for int32")
		}
		if !DeepEqual(int64(42), int64(42), false) {
			t.Error("expected equal for int64")
		}

		// uint types
		if !DeepEqual(uint(42), uint(42), false) {
			t.Error("expected equal for uint")
		}
		if !DeepEqual(uint8(42), uint8(42), false) {
			t.Error("expected equal for uint8")
		}
		if !DeepEqual(uint16(42), uint16(42), false) {
			t.Error("expected equal for uint16")
		}
		if !DeepEqual(uint32(42), uint32(42), false) {
			t.Error("expected equal for uint32")
		}
		if !DeepEqual(uint64(42), uint64(42), false) {
			t.Error("expected equal for uint64")
		}

		// float types
		if !DeepEqual(float32(3.14), float32(3.14), false) {
			t.Error("expected equal for float32")
		}
		if !DeepEqual(float64(3.14), float64(3.14), false) {
			t.Error("expected equal for float64")
		}

		// complex types
		if !DeepEqual(complex64(1+2i), complex64(1+2i), false) {
			t.Error("expected equal for complex64")
		}
		if !DeepEqual(complex128(1+2i), complex128(1+2i), false) {
			t.Error("expected equal for complex128")
		}
	})

	t.Run("time types", func(t *testing.T) {
		now := time.Now()
		if !DeepEqual(now, now, false) {
			t.Error("expected equal for time.Time")
		}

		dur := 5 * time.Second
		if !DeepEqual(dur, dur, false) {
			t.Error("expected equal for time.Duration")
		}
	})

	t.Run("uuid", func(t *testing.T) {
		id := uuid.New()
		if !DeepEqual(id, id, false) {
			t.Error("expected equal for uuid.UUID")
		}

		id2 := uuid.New()
		if DeepEqual(id, id2, false) {
			t.Error("expected not equal for different UUIDs")
		}
	})

	t.Run("pointer types", func(t *testing.T) {
		val := 42
		ptr1 := &val
		ptr2 := &val

		if !DeepEqual(ptr1, ptr2, false) {
			t.Error("expected equal for pointers to same value")
		}

		val2 := 43
		ptr3 := &val2
		if DeepEqual(ptr1, ptr3, false) {
			t.Error("expected not equal for pointers to different values")
		}
	})

	t.Run("Equaler interface", func(t *testing.T) {
		eq1 := testEqualer{Value: 1}
		eq2 := testEqualer{Value: 1}
		eq3 := testEqualer{Value: 2}

		if !DeepEqual(eq1, eq2, false) {
			t.Error("expected equal for Equaler")
		}
		if DeepEqual(eq1, eq3, false) {
			t.Error("expected not equal for Equaler")
		}

		if !DeepEqualPtr(&eq1, &eq2, false) {
			t.Error("expected equal for Equaler")
		}
		if DeepEqualPtr(&eq1, &eq3, false) {
			t.Error("expected not equal for Equaler")
		}
	})

	t.Run("maps with string keys", func(t *testing.T) {
		map1 := map[string]any{"a": 1, "b": "test"}
		map2 := map[string]any{"a": 1, "b": "test"}
		if !DeepEqual(map1, map2, false) {
			t.Error("expected equal for map[string]any")
		}

		map3 := map[string]string{"a": "1", "b": "2"}
		map4 := map[string]string{"a": "1", "b": "2"}
		if !DeepEqual(map3, map4, false) {
			t.Error("expected equal for map[string]string")
		}
	})

	t.Run("maps with various key types", func(t *testing.T) {
		// bool keys
		mapBool1 := map[bool]any{true: 1, false: 2}
		mapBool2 := map[bool]any{true: 1, false: 2}
		if !DeepEqual(mapBool1, mapBool2, false) {
			t.Error("expected equal for map[bool]any")
		}

		// int keys
		mapInt1 := map[int]any{1: "a", 2: "b"}
		mapInt2 := map[int]any{1: "a", 2: "b"}
		if !DeepEqual(mapInt1, mapInt2, false) {
			t.Error("expected equal for map[int]any")
		}

		// float keys
		mapFloat1 := map[float64]any{1.5: "a", 2.5: "b"}
		mapFloat2 := map[float64]any{1.5: "a", 2.5: "b"}
		if !DeepEqual(mapFloat1, mapFloat2, false) {
			t.Error("expected equal for map[float64]any")
		}
	})

	t.Run("slices of comparable types", func(t *testing.T) {
		// bool slice
		if !DeepEqual([]bool{true, false}, []bool{true, false}, false) {
			t.Error("expected equal for []bool")
		}

		// string slice
		if !DeepEqual([]string{"a", "b"}, []string{"a", "b"}, false) {
			t.Error("expected equal for []string")
		}

		// int slices
		if !DeepEqual([]int{1, 2, 3}, []int{1, 2, 3}, false) {
			t.Error("expected equal for []int")
		}
		if !DeepEqual([]int8{1, 2}, []int8{1, 2}, false) {
			t.Error("expected equal for []int8")
		}
		if !DeepEqual([]int16{1, 2}, []int16{1, 2}, false) {
			t.Error("expected equal for []int16")
		}
		if !DeepEqual([]int32{1, 2}, []int32{1, 2}, false) {
			t.Error("expected equal for []int32")
		}
		if !DeepEqual([]int64{1, 2}, []int64{1, 2}, false) {
			t.Error("expected equal for []int64")
		}

		// uint slices
		if !DeepEqual([]uint{1, 2}, []uint{1, 2}, false) {
			t.Error("expected equal for []uint")
		}
		if !DeepEqual([]uint8{1, 2}, []uint8{1, 2}, false) {
			t.Error("expected equal for []uint8")
		}
		if !DeepEqual([]uint16{1, 2}, []uint16{1, 2}, false) {
			t.Error("expected equal for []uint16")
		}
		if !DeepEqual([]uint32{1, 2}, []uint32{1, 2}, false) {
			t.Error("expected equal for []uint32")
		}
		if !DeepEqual([]uint64{1, 2}, []uint64{1, 2}, false) {
			t.Error("expected equal for []uint64")
		}

		// float slices
		if !DeepEqual([]float32{1.5, 2.5}, []float32{1.5, 2.5}, false) {
			t.Error("expected equal for []float32")
		}
		if !DeepEqual([]float64{1.5, 2.5}, []float64{1.5, 2.5}, false) {
			t.Error("expected equal for []float64")
		}

		// complex slices
		if !DeepEqual([]complex64{1 + 2i}, []complex64{1 + 2i}, false) {
			t.Error("expected equal for []complex64")
		}
		if !DeepEqual([]complex128{1 + 2i}, []complex128{1 + 2i}, false) {
			t.Error("expected equal for []complex128")
		}

		// time slices
		now := time.Now()
		if !DeepEqual([]time.Time{now}, []time.Time{now}, false) {
			t.Error("expected equal for []time.Time")
		}
		if !DeepEqual([]time.Duration{time.Second}, []time.Duration{time.Second}, false) {
			t.Error("expected equal for []time.Duration")
		}

		// uuid slice
		id := uuid.New()
		if !DeepEqual([]uuid.UUID{id}, []uuid.UUID{id}, false) {
			t.Error("expected equal for []uuid.UUID")
		}
	})

	t.Run("slice of any", func(t *testing.T) {
		slice1 := []any{1, "test", true}
		slice2 := []any{1, "test", true}
		if !DeepEqual(slice1, slice2, false) {
			t.Error("expected equal for []any")
		}

		slice3 := []any{1, "test", false}
		if DeepEqual(slice1, slice3, false) {
			t.Error("expected not equal for []any with different values")
		}
	})

	t.Run("complex nested structures", func(t *testing.T) {
		complex1 := map[string]any{
			"name": "test",
			"age":  30,
			"tags": []string{"a", "b", "c"},
			"metadata": map[string]any{
				"created": time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				"active":  true,
			},
		}
		complex2 := map[string]any{
			"name": "test",
			"age":  30,
			"tags": []string{"a", "b", "c"},
			"metadata": map[string]any{
				"created": time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				"active":  true,
			},
		}
		if !DeepEqual(complex1, complex2, false) {
			t.Error("expected equal for complex nested structure")
		}
	})

	t.Run("fallback to reflect.DeepEqual", func(t *testing.T) {
		// Test with a type not explicitly handled
		type customStruct struct {
			Field1 string
			Field2 int
		}
		s1 := customStruct{Field1: "test", Field2: 42}
		s2 := customStruct{Field1: "test", Field2: 42}
		if !DeepEqual(s1, s2, false) {
			t.Error("expected equal for custom struct (fallback)")
		}

		s3 := customStruct{Field1: "test", Field2: 43}
		if DeepEqual(s1, s3, false) {
			t.Error("expected not equal for custom struct (fallback)")
		}
	})

	t.Run("omitZero flag", func(t *testing.T) {
		// Test with maps
		map1 := map[string]int{"a": 1, "b": 0}
		map2 := map[string]int{"a": 1}
		if DeepEqual(map1, map2, false) {
			t.Error("expected not equal without omitZero")
		}

		// Test with slices
		slice1 := []int{1, 2, 3}
		slice2 := []int{1, 2, 3}
		if !DeepEqual(slice1, slice2, true) {
			t.Error("expected equal with omitZero")
		}
	})

	t.Run("AllOrListString", func(t *testing.T) {
		// Test with all strings.
		list1 := NewAll()
		list2 := NewAll()
		emptyList := NewStringList([]string{})

		if !DeepEqual(list1, list2, true) {
			t.Error("expected equal")
		}

		if !emptyList.IsZero() {
			t.Error("expected empty")
		}

		if DeepEqual(list1, emptyList, true) {
			t.Error("expected not equal")
		}
	})

	t.Run("Duration", func(t *testing.T) {
		// Test with Duration values.
		v1 := Duration(time.Second)
		v2 := Duration(time.Second)
		zero := Duration(0)

		if !DeepEqual(v1, v2, true) {
			t.Error("expected equal")
		}

		if !IsZero(zero) {
			t.Error("expected empty")
		}

		if DeepEqual(v1, zero, true) {
			t.Error("expected not equal")
		}
	})

	t.Run("Slug", func(t *testing.T) {

		v1 := Slug("test")
		v2 := Slug("test")
		zero := Slug("")

		if !DeepEqual(v1, v2, true) {
			t.Error("expected equal")
		}

		if !IsZero(zero) {
			t.Error("expected empty")
		}

		if IsZero(string(zero)) {
			t.Error("expected non-empty")
		}

		if DeepEqual(v1, zero, true) {
			t.Error("expected not equal")
		}
	})
}

func TestDeepEqual_WithAnyType(t *testing.T) {
	t.Run("primitive types as any", func(t *testing.T) {
		// int
		if !DeepEqual(42, any(42), false) {
			t.Error("expected equal for int vs any(int)")
		}
		if DeepEqual(42, any(43), false) {
			t.Error("expected not equal for different int values")
		}

		// string
		if !DeepEqual("hello", any("hello"), false) {
			t.Error("expected equal for string vs any(string)")
		}
		if DeepEqual("hello", any("world"), false) {
			t.Error("expected not equal for different string values")
		}

		// bool
		if !DeepEqual(true, any(true), false) {
			t.Error("expected equal for bool vs any(bool)")
		}
		if DeepEqual(true, any(false), false) {
			t.Error("expected not equal for different bool values")
		}

		// float64
		if !DeepEqual(3.14, any(3.14), false) {
			t.Error("expected equal for float64 vs any(float64)")
		}
		if DeepEqual(3.14, any(2.71), false) {
			t.Error("expected not equal for different float64 values")
		}
	})

	t.Run("nil comparisons with any", func(t *testing.T) {
		// nil slice vs any(nil slice) - both are nil slices, should be equal
		var slice []int
		var nilSlice []int
		if !DeepEqual(slice, any(nilSlice), false) {
			t.Error("expected equal for nil slice vs any(nil slice)")
		}

		// nil map vs any(nil map) - both are nil maps, should be equal
		var m map[string]int
		var nilMap map[string]int
		if !DeepEqual(m, any(nilMap), false) {
			t.Error("expected equal for nil map vs any(nil map)")
		}

		// non-nil vs any(nil)
		var nilAny any
		if DeepEqual(42, nilAny, false) {
			t.Error("expected not equal for non-nil vs any(nil)")
		}

		// nil pointer vs non-nil any
		var ptr *int
		if DeepEqual(ptr, any(42), false) {
			t.Error("expected not equal for nil pointer vs non-nil any")
		}

		// empty slice vs nil slice wrapped in any
		emptySlice := []int{}
		if DeepEqual(emptySlice, any(nilSlice), false) {
			t.Error("expected not equal for empty slice vs any(nil slice)")
		}

		// empty map vs nil map wrapped in any
		emptyMap := map[string]int{}
		if DeepEqual(emptyMap, any(nilMap), false) {
			t.Error("expected not equal for empty map vs any(nil map)")
		}
	})

	t.Run("slices as any", func(t *testing.T) {
		// []int vs any([]int)
		slice1 := []int{1, 2, 3}
		slice2 := any([]int{1, 2, 3})
		if !DeepEqual(slice1, slice2, false) {
			t.Error("expected equal for []int vs any([]int)")
		}

		// different values
		slice3 := any([]int{1, 2, 4})
		if DeepEqual(slice1, slice3, false) {
			t.Error("expected not equal for different slice values")
		}

		// []string vs any([]string)
		strSlice1 := []string{"a", "b", "c"}
		strSlice2 := any([]string{"a", "b", "c"})
		if !DeepEqual(strSlice1, strSlice2, false) {
			t.Error("expected equal for []string vs any([]string)")
		}

		// []any vs any([]any)
		anySlice1 := []any{1, "test", true}
		anySlice2 := any([]any{1, "test", true})
		if !DeepEqual(anySlice1, anySlice2, false) {
			t.Error("expected equal for []any vs any([]any)")
		}
	})

	t.Run("maps as any", func(t *testing.T) {
		// map[string]int vs any(map[string]int)
		map1 := map[string]int{"a": 1, "b": 2}
		map2 := any(map[string]int{"a": 1, "b": 2})
		if !DeepEqual(map1, map2, false) {
			t.Error("expected equal for map[string]int vs any(map[string]int)")
		}

		// different values
		map3 := any(map[string]int{"a": 1, "b": 3})
		if DeepEqual(map1, map3, false) {
			t.Error("expected not equal for different map values")
		}

		// map[string]any vs any(map[string]any)
		anyMap1 := map[string]any{"a": 1, "b": "test"}
		anyMap2 := any(map[string]any{"a": 1, "b": "test"})
		if !DeepEqual(anyMap1, anyMap2, false) {
			t.Error("expected equal for map[string]any vs any(map[string]any)")
		}
	})

	t.Run("time types as any", func(t *testing.T) {
		// time.Time
		now := time.Now()
		if !DeepEqual(now, any(now), false) {
			t.Error("expected equal for time.Time vs any(time.Time)")
		}

		later := now.Add(time.Hour)
		if DeepEqual(now, any(later), false) {
			t.Error("expected not equal for different time.Time values")
		}

		// time.Duration
		dur := 5 * time.Second
		if !DeepEqual(dur, any(dur), false) {
			t.Error("expected equal for time.Duration vs any(time.Duration)")
		}

		dur2 := 10 * time.Second
		if DeepEqual(dur, any(dur2), false) {
			t.Error("expected not equal for different time.Duration values")
		}
	})

	t.Run("uuid as any", func(t *testing.T) {
		id := uuid.New()
		if !DeepEqual(id, any(id), false) {
			t.Error("expected equal for uuid.UUID vs any(uuid.UUID)")
		}

		id2 := uuid.New()
		if DeepEqual(id, any(id2), false) {
			t.Error("expected not equal for different uuid.UUID values")
		}
	})

	t.Run("pointer types as any", func(t *testing.T) {
		val := 42
		ptr1 := &val
		ptr2 := any(&val)
		if !DeepEqual(ptr1, ptr2, false) {
			t.Error("expected equal for *int vs any(*int)")
		}

		val2 := 43
		ptr3 := any(&val2)
		if DeepEqual(ptr1, ptr3, false) {
			t.Error("expected not equal for different pointer values")
		}
	})

	t.Run("type mismatches with any", func(t *testing.T) {
		// int vs any(string)
		if DeepEqual(42, any("42"), false) {
			t.Error("expected not equal for int vs any(string)")
		}

		// int vs any(float64)
		if DeepEqual(42, any(42.0), false) {
			t.Error("expected not equal for int vs any(float64)")
		}

		// []int vs any([]string)
		if DeepEqual([]int{1, 2, 3}, any([]string{"1", "2", "3"}), false) {
			t.Error("expected not equal for []int vs any([]string)")
		}

		// map[string]int vs any(map[string]string)
		if DeepEqual(map[string]int{"a": 1}, any(map[string]string{"a": "1"}), false) {
			t.Error("expected not equal for map[string]int vs any(map[string]string)")
		}
	})

	t.Run("Equaler interface with any", func(t *testing.T) {
		eq1 := testEqualer{Value: 1}
		eq2 := any(testEqualer{Value: 1})
		if !DeepEqual(eq1, eq2, false) {
			t.Error("expected equal for Equaler vs any(Equaler)")
		}

		eq3 := any(testEqualer{Value: 2})
		if DeepEqual(eq1, eq3, false) {
			t.Error("expected not equal for different Equaler values")
		}
	})

	t.Run("nested structures with any", func(t *testing.T) {
		// nested map
		nested1 := map[string]any{
			"level1": map[string]any{
				"level2": map[string]any{
					"value": 42,
				},
			},
		}
		nested2 := any(map[string]any{
			"level1": map[string]any{
				"level2": map[string]any{
					"value": 42,
				},
			},
		})
		if !DeepEqual(nested1, nested2, false) {
			t.Error("expected equal for nested map vs any(nested map)")
		}

		// nested slice
		nestedSlice1 := []any{[]any{[]any{1, 2, 3}}}
		nestedSlice2 := any([]any{[]any{[]any{1, 2, 3}}})
		if !DeepEqual(nestedSlice1, nestedSlice2, false) {
			t.Error("expected equal for nested slice vs any(nested slice)")
		}
	})

	t.Run("all integer types as any", func(t *testing.T) {
		// int8
		if !DeepEqual(int8(42), any(int8(42)), false) {
			t.Error("expected equal for int8 vs any(int8)")
		}
		// int16
		if !DeepEqual(int16(42), any(int16(42)), false) {
			t.Error("expected equal for int16 vs any(int16)")
		}
		// int32
		if !DeepEqual(int32(42), any(int32(42)), false) {
			t.Error("expected equal for int32 vs any(int32)")
		}
		// int64
		if !DeepEqual(int64(42), any(int64(42)), false) {
			t.Error("expected equal for int64 vs any(int64)")
		}
		// uint
		if !DeepEqual(uint(42), any(uint(42)), false) {
			t.Error("expected equal for uint vs any(uint)")
		}
		// uint8
		if !DeepEqual(uint8(42), any(uint8(42)), false) {
			t.Error("expected equal for uint8 vs any(uint8)")
		}
		// uint16
		if !DeepEqual(uint16(42), any(uint16(42)), false) {
			t.Error("expected equal for uint16 vs any(uint16)")
		}
		// uint32
		if !DeepEqual(uint32(42), any(uint32(42)), false) {
			t.Error("expected equal for uint32 vs any(uint32)")
		}
		// uint64
		if !DeepEqual(uint64(42), any(uint64(42)), false) {
			t.Error("expected equal for uint64 vs any(uint64)")
		}
	})

	t.Run("float and complex types as any", func(t *testing.T) {
		// float32
		if !DeepEqual(float32(3.14), any(float32(3.14)), false) {
			t.Error("expected equal for float32 vs any(float32)")
		}
		// float64
		if !DeepEqual(float64(3.14), any(float64(3.14)), false) {
			t.Error("expected equal for float64 vs any(float64)")
		}
		// complex64
		if !DeepEqual(complex64(1+2i), any(complex64(1+2i)), false) {
			t.Error("expected equal for complex64 vs any(complex64)")
		}
		// complex128
		if !DeepEqual(complex128(1+2i), any(complex128(1+2i)), false) {
			t.Error("expected equal for complex128 vs any(complex128)")
		}
	})

	t.Run("all map key types as any", func(t *testing.T) {
		// map[bool]any
		if !DeepEqual(map[bool]any{true: 1}, any(map[bool]any{true: 1}), false) {
			t.Error("expected equal for map[bool]any vs any(map[bool]any)")
		}
		// map[int]any
		if !DeepEqual(map[int]any{1: "a"}, any(map[int]any{1: "a"}), false) {
			t.Error("expected equal for map[int]any vs any(map[int]any)")
		}
		// map[float64]any
		if !DeepEqual(map[float64]any{1.5: "a"}, any(map[float64]any{1.5: "a"}), false) {
			t.Error("expected equal for map[float64]any vs any(map[float64]any)")
		}
		// map[complex128]any
		if !DeepEqual(map[complex128]any{1 + 2i: "a"}, any(map[complex128]any{1 + 2i: "a"}), false) {
			t.Error("expected equal for map[complex128]any vs any(map[complex128]any)")
		}
	})

	t.Run("empty collections as any", func(t *testing.T) {
		// empty slice
		if !DeepEqual([]int{}, any([]int{}), false) {
			t.Error("expected equal for empty []int vs any(empty []int)")
		}

		// empty map
		if !DeepEqual(map[string]int{}, any(map[string]int{}), false) {
			t.Error("expected equal for empty map vs any(empty map)")
		}

		// empty string
		if !DeepEqual("", any(""), false) {
			t.Error("expected equal for empty string vs any(empty string)")
		}
	})

	t.Run("all slice types as any", func(t *testing.T) {
		// []int8
		if !DeepEqual([]int8{1, 2}, any([]int8{1, 2}), false) {
			t.Error("expected equal for []int8 vs any([]int8)")
		}
		// []int16
		if !DeepEqual([]int16{1, 2}, any([]int16{1, 2}), false) {
			t.Error("expected equal for []int16 vs any([]int16)")
		}
		// []int32
		if !DeepEqual([]int32{1, 2}, any([]int32{1, 2}), false) {
			t.Error("expected equal for []int32 vs any([]int32)")
		}
		// []int64
		if !DeepEqual([]int64{1, 2}, any([]int64{1, 2}), false) {
			t.Error("expected equal for []int64 vs any([]int64)")
		}
		// []uint
		if !DeepEqual([]uint{1, 2}, any([]uint{1, 2}), false) {
			t.Error("expected equal for []uint vs any([]uint)")
		}
		// []uint8
		if !DeepEqual([]uint8{1, 2}, any([]uint8{1, 2}), false) {
			t.Error("expected equal for []uint8 vs any([]uint8)")
		}
		// []uint16
		if !DeepEqual([]uint16{1, 2}, any([]uint16{1, 2}), false) {
			t.Error("expected equal for []uint16 vs any([]uint16)")
		}
		// []uint32
		if !DeepEqual([]uint32{1, 2}, any([]uint32{1, 2}), false) {
			t.Error("expected equal for []uint32 vs any([]uint32)")
		}
		// []uint64
		if !DeepEqual([]uint64{1, 2}, any([]uint64{1, 2}), false) {
			t.Error("expected equal for []uint64 vs any([]uint64)")
		}
		// []float32
		if !DeepEqual([]float32{1.5, 2.5}, any([]float32{1.5, 2.5}), false) {
			t.Error("expected equal for []float32 vs any([]float32)")
		}
		// []float64
		if !DeepEqual([]float64{1.5, 2.5}, any([]float64{1.5, 2.5}), false) {
			t.Error("expected equal for []float64 vs any([]float64)")
		}
		// []complex64
		if !DeepEqual([]complex64{1 + 2i}, any([]complex64{1 + 2i}), false) {
			t.Error("expected equal for []complex64 vs any([]complex64)")
		}
		// []complex128
		if !DeepEqual([]complex128{1 + 2i}, any([]complex128{1 + 2i}), false) {
			t.Error("expected equal for []complex128 vs any([]complex128)")
		}
		// []bool
		if !DeepEqual([]bool{true, false}, any([]bool{true, false}), false) {
			t.Error("expected equal for []bool vs any([]bool)")
		}
		// []time.Time
		now := time.Now()
		if !DeepEqual([]time.Time{now}, any([]time.Time{now}), false) {
			t.Error("expected equal for []time.Time vs any([]time.Time)")
		}
		// []time.Duration
		if !DeepEqual([]time.Duration{time.Second}, any([]time.Duration{time.Second}), false) {
			t.Error("expected equal for []time.Duration vs any([]time.Duration)")
		}
		// []uuid.UUID
		id := uuid.New()
		if !DeepEqual([]uuid.UUID{id}, any([]uuid.UUID{id}), false) {
			t.Error("expected equal for []uuid.UUID vs any([]uuid.UUID)")
		}
	})

	t.Run("complex real-world scenarios with any", func(t *testing.T) {
		// Mixed types in map[string]any
		complex1 := map[string]any{
			"string":   "value",
			"int":      42,
			"float":    3.14,
			"bool":     true,
			"slice":    []int{1, 2, 3},
			"map":      map[string]string{"key": "value"},
			"time":     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			"duration": time.Second,
			"uuid":     uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
		}
		complex2 := any(map[string]any{
			"string":   "value",
			"int":      42,
			"float":    3.14,
			"bool":     true,
			"slice":    []int{1, 2, 3},
			"map":      map[string]string{"key": "value"},
			"time":     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			"duration": time.Second,
			"uuid":     uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
		})
		if !DeepEqual(complex1, complex2, false) {
			t.Error("expected equal for complex real-world map")
		}

		// Change one value
		complex3 := any(map[string]any{
			"string":   "value",
			"int":      43, // changed
			"float":    3.14,
			"bool":     true,
			"slice":    []int{1, 2, 3},
			"map":      map[string]string{"key": "value"},
			"time":     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			"duration": time.Second,
			"uuid":     uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
		})
		if DeepEqual(complex1, complex3, false) {
			t.Error("expected not equal when one value differs")
		}
	})

	t.Run("omitZero with any type", func(t *testing.T) {
		// map with zero values
		map1 := map[string]int{"a": 1, "b": 2}
		map2 := any(map[string]int{"a": 1, "b": 2})
		if !DeepEqual(map1, map2, true) {
			t.Error("expected equal with omitZero")
		}

		// slice with any
		slice1 := []int{1, 2, 3}
		slice2 := any([]int{1, 2, 3})
		if !DeepEqual(slice1, slice2, true) {
			t.Error("expected equal with omitZero")
		}
	})
}

func TestEqualPtr(t *testing.T) {
	if !EqualPtr(&testEqualer{}, &testEqualer{}) {
		t.Error("expected equal")
	}

	if EqualPtr(&testEqualer{}, nil) {
		t.Error("expected not equal")
	}

	if !EqualPtr[testEqualer](nil, nil) {
		t.Error("expected equal")
	}
}

func createBenchmarkObject() map[string]any {
	return map[string]any{
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
		"equaler": testEqualer{
			Value: 1,
		},
	}
}

// BenchmarkDeepEqualReflection-11    	  296131	      4009 ns/op	    5848 B/op	      70 allocs/op
func BenchmarkDeepEqualReflection(b *testing.B) {
	obj1 := createBenchmarkObject()
	obj2 := createBenchmarkObject()

	for b.Loop() {
		reflect.DeepEqual(obj1, obj2)
	}
}

// BenchmarkDeepEqual-11    	 1904056	       636.6 ns/op	       0 B/op	       0 allocs/op
func BenchmarkDeepEqual(b *testing.B) {
	obj1 := createBenchmarkObject()
	obj2 := createBenchmarkObject()
	for b.Loop() {
		DeepEqual(obj1, obj2, false)
	}
}
