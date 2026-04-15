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
	"testing"
	"time"
)

func TestDecodeBool(t *testing.T) {
	value, err := DecodeBoolean(true)
	assertNilError(t, err)
	assertEqual(t, true, value)

	ptr, err := DecodeNullableBoolean(true)
	assertNilError(t, err)
	assertEqual(t, true, *ptr)

	ptr2, err := DecodeNullableBoolean(ptr)
	assertNilError(t, err)
	assertEqual(t, *ptr, *ptr2)

	_, err = DecodeBoolean(nil)
	assertError(t, err, "boolean value must not be null")

	t.Run("decode_string", func(t *testing.T) {
		for _, tc := range []struct {
			input    string
			expected bool
		}{
			{"true", true},
			{"false", false},
			{"1", true},
			{"0", false},
			{"TRUE", true},
			{"FALSE", false},
		} {
			v, err := DecodeBoolean(tc.input)
			assertNilError(t, err)
			assertEqual(t, tc.expected, v)
		}

		_, err := DecodeBoolean("failure")
		assertError(t, err, "malformed boolean")
	})

	t.Run("decode_string_pointer", func(t *testing.T) {
		s := "true"
		v, err := DecodeBoolean(&s)
		assertNilError(t, err)
		assertEqual(t, true, v)

		var nilStr *string
		_, err = DecodeBoolean(nilStr)
		assertError(t, err, "boolean value must not be null")

		bad := "bad"
		_, err = DecodeBoolean(&bad)
		assertError(t, err, "malformed boolean")
	})

	t.Run("decode_nullable_string", func(t *testing.T) {
		ptr, err := DecodeNullableBoolean("true")
		assertNilError(t, err)
		assertEqual(t, true, *ptr)

		ptr, err = DecodeNullableBoolean("false")
		assertNilError(t, err)
		assertEqual(t, false, *ptr)

		_, err = DecodeNullableBoolean("oops")
		assertError(t, err, "malformed boolean")
	})

	t.Run("decode_nullable_string_pointer", func(t *testing.T) {
		s := "1"
		ptr, err := DecodeNullableBoolean(&s)
		assertNilError(t, err)
		assertEqual(t, true, *ptr)

		var nilStr *string
		ptr, err = DecodeNullableBoolean(nilStr)
		assertNilError(t, err)
		assertEqual(t, (*bool)(nil), ptr)

		bad := "bad"
		_, err = DecodeNullableBoolean(&bad)
		assertError(t, err, "malformed boolean")
	})
}

func TestDecodeBooleanSlice(t *testing.T) {
	value, err := DecodeBooleanSlice([]bool{true, false})
	assertNilError(t, err)
	assertDeepEqual(t, []bool{true, false}, value)

	value, err = DecodeBooleanSlice([]any{false, true})
	assertNilError(t, err)
	assertDeepEqual(t, []bool{false, true}, value)

	value, err = DecodeBooleanSlice(&[]any{true})
	assertNilError(t, err)
	assertDeepEqual(t, []bool{true}, value)

	_, err = DecodeBooleanSlice(nil)
	assertError(t, err, "boolean slice must not be null")

	_, err = DecodeBooleanSlice("failure")
	assertError(t, err, "malformed boolean slice; got: string")

	_, err = DecodeBooleanSlice(time.Now())
	assertError(t, err, "malformed boolean slice; got: struct")

	_, err = DecodeBooleanSlice([]any{nil})
	assertError(t, err, "element 0: boolean value must not be null")

	_, err = DecodeBooleanSlice(&[]any{nil})
	assertError(t, err, "element 0: boolean value must not be null")
}

func TestDecodeNullableBooleanSlice(t *testing.T) {
	value, err := DecodeNullableBooleanSlice([]bool{true, false})
	assertNilError(t, err)
	assertDeepEqual(t, []bool{true, false}, *value)

	value, err = DecodeNullableBooleanSlice([]*bool{new(true), new(false)})
	assertNilError(t, err)
	assertDeepEqual(t, []bool{true, false}, *value)

	value, err = DecodeNullableBooleanSlice(&[]*bool{new(true), new(false)})
	assertNilError(t, err)
	assertDeepEqual(t, []bool{true, false}, *value)

	value, err = DecodeNullableBooleanSlice([]any{false, true})
	assertNilError(t, err)
	assertDeepEqual(t, []bool{false, true}, *value)

	value, err = DecodeNullableBooleanSlice(&[]any{true})
	assertNilError(t, err)
	assertDeepEqual(t, []bool{true}, *value)

	_, err = DecodeNullableBooleanSlice([]any{nil})
	assertError(t, err, "element 0: boolean value must not be null")

	value, err = DecodeNullableBooleanSlice(nil)
	assertNilError(t, err)
	assertEqual(t, nil, value)

	value, err = DecodeNullableBooleanSlice((*string)(nil))
	assertNilError(t, err)
	assertEqual(t, nil, value)

	_, err = DecodeNullableBooleanSlice([]any{"true"})
	assertError(
		t,
		err,
		"failed to decode boolean element at 0: malformed boolean; got: interface",
	)

	_, err = DecodeNullableBooleanSlice("failure")
	assertError(t, err, "malformed boolean slice; got: string")

	_, err = DecodeNullableBooleanSlice(time.Now())
	assertError(t, err, "malformed boolean slice; got: struct")
}

func TestDecodeString(t *testing.T) {
	value, err := DecodeString("success")
	assertNilError(t, err)
	assertEqual(t, "success", value)

	ptr, err := DecodeNullableString("pointer")
	assertNilError(t, err)
	assertEqual(t, "pointer", *ptr)

	ptr2, err := DecodeNullableString(ptr)
	assertNilError(t, err)
	assertEqual(t, *ptr, *ptr2)

	_, err = DecodeString(nil)
	assertEqual(t, err.Error(), ErrStringNull.Error())

	_, err = DecodeString(0)
	assertError(t, err, "malformed string, got: int")
}

func TestDecodeStringSlice(t *testing.T) {
	value, err := DecodeStringSlice([]string{"foo", "bar"})
	assertNilError(t, err)
	assertDeepEqual(t, []string{"foo", "bar"}, value)

	value, err = DecodeStringSlice([]*string{new("foo"), new("bar")})
	assertNilError(t, err)
	assertDeepEqual(t, []string{"foo", "bar"}, value)

	value, err = DecodeStringSlice(&[]*string{new("foo")})
	assertNilError(t, err)
	assertDeepEqual(t, []string{"foo"}, value)

	value, err = DecodeStringSlice([]any{"bar", "foo"})
	assertNilError(t, err)
	assertDeepEqual(t, []string{"bar", "foo"}, value)

	value, err = DecodeStringSlice(&[]any{"foo"})
	assertNilError(t, err)
	assertDeepEqual(t, []string{"foo"}, value)

	value, err = DecodeStringSlice(nil)
	assertNilError(t, err)
	assertDeepEqual(t, nil, value)

	_, err = DecodeStringSlice("failure")
	assertError(t, err, "malformed string slice; got: string")

	_, err = DecodeStringSlice(time.Now())
	assertError(t, err, "malformed string slice; got: struct")

	_, err = DecodeStringSlice([]any{nil})
	assertError(t, err, "failed to decode string at 0: string value must not be null")

	_, err = DecodeStringSlice(&[]any{nil})
	assertError(t, err, "failed to decode element at 0: string value must not be null")
}

func TestDecodeNumber(t *testing.T) {
	for _, expected := range []any{int(1), int8(2), int16(3), int32(4), int64(5), uint(6), uint8(7), uint16(8), uint32(9), uint64(10)} {
		t.Run(fmt.Sprintf("decode_%s", reflect.TypeOf(expected).String()), func(t *testing.T) {
			value, err := DecodeNumber[float64](expected)
			assertNilError(t, err)
			assertEqual(t, fmt.Sprint(expected), fmt.Sprintf("%.0f", value))

			ptr, err := DecodeNullableNumber[float64](&expected)
			assertNilError(t, err)
			assertEqual(t, fmt.Sprint(expected), fmt.Sprintf("%.0f", *ptr))

			ptr2, err := DecodeNullableNumber[float64](ptr)
			assertNilError(t, err)
			assertEqual(t, fmt.Sprintf("%.1f", *ptr), fmt.Sprintf("%.1f", *ptr2))
		})
	}

	for _, expected := range []any{float32(1.1), float64(2.2)} {
		t.Run(fmt.Sprintf("decode_%s", reflect.TypeOf(expected).String()), func(t *testing.T) {
			value, err := DecodeNumber[float64](expected)
			assertNilError(t, err)
			assertEqual(t, fmt.Sprintf("%.1f", expected), fmt.Sprintf("%.1f", value))

			ptr, err := DecodeNullableNumber[float64](&expected)
			assertNilError(t, err)
			assertEqual(t, fmt.Sprintf("%.1f", expected), fmt.Sprintf("%.1f", *ptr))

			ptr2, err := DecodeNullableNumber[float64](ptr)
			assertNilError(t, err)
			assertEqual(t, fmt.Sprintf("%.1f", *ptr), fmt.Sprintf("%.1f", *ptr2))
		})
	}

	t.Run("decode_string", func(t *testing.T) {
		expected := "0"
		value, err := DecodeNumberReflection[float64](reflect.ValueOf(expected))
		assertNilError(t, err)
		assertEqual(t, "0", fmt.Sprintf("%.0f", value))

		ptr, err := DecodeNullableNumberReflection[float64](reflect.ValueOf(expected))
		assertNilError(t, err)
		assertEqual(t, expected, fmt.Sprintf("%.0f", *ptr))

		ptr2, err := DecodeNullableNumberReflection[float64](reflect.ValueOf(ptr))
		assertNilError(t, err)
		assertEqual(t, fmt.Sprintf("%.1f", *ptr), fmt.Sprintf("%.1f", *ptr2))
	})

	t.Run("decode_pointers", func(t *testing.T) {
		vInt := int(1)
		ptr, err := DecodeNullableNumber[float64](&vInt)
		assertNilError(t, err)
		assertEqual(t, float64(vInt), *ptr)

		ptrInt := &vInt
		ptr, err = DecodeNullableNumber[float64](&ptrInt)
		assertNilError(t, err)
		assertEqual(t, float64(vInt), *ptr)

		vInt8 := int8(1)
		ptr, err = DecodeNullableNumber[float64](&vInt8)
		assertNilError(t, err)
		assertEqual(t, float64(vInt8), *ptr)

		vInt16 := int16(1)
		ptr, err = DecodeNullableNumber[float64](&vInt16)
		assertNilError(t, err)
		assertEqual(t, float64(vInt16), *ptr)

		vInt32 := int32(1)
		ptr, err = DecodeNullableNumber[float64](&vInt32)
		assertNilError(t, err)
		assertEqual(t, float64(vInt32), *ptr)

		vInt64 := int64(1)
		ptr, err = DecodeNullableNumber[float64](&vInt64)
		assertNilError(t, err)
		assertEqual(t, float64(vInt64), *ptr)

		vUint := uint(1)
		ptr, err = DecodeNullableNumber[float64](&vUint)
		assertNilError(t, err)
		assertEqual(t, float64(vUint), *ptr)

		ptrUint := &vUint
		ptr, err = DecodeNullableNumber[float64](&ptrUint)
		assertNilError(t, err)
		assertEqual(t, float64(vUint), *ptr)

		vUint8 := uint8(1)
		ptr, err = DecodeNullableNumber[float64](&vUint8)
		assertNilError(t, err)
		assertEqual(t, float64(vUint8), *ptr)

		vUint16 := uint16(1)
		ptr, err = DecodeNullableNumber[float64](&vUint16)
		assertNilError(t, err)
		assertEqual(t, float64(vUint16), *ptr)

		vUint32 := uint32(1)
		ptr, err = DecodeNullableNumber[float64](&vUint32)
		assertNilError(t, err)
		assertEqual(t, float64(vUint32), *ptr)

		vUint64 := uint64(1)
		ptr, err = DecodeNullableNumber[float64](&vUint64)
		assertNilError(t, err)
		assertEqual(t, float64(vUint64), *ptr)

		vFloat32 := float32(1)
		ptr, err = DecodeNullableNumber[float64](&vFloat32)
		assertNilError(t, err)
		assertEqual(t, float64(vFloat32), *ptr)

		ptrFloat32 := &vFloat32
		ptr, err = DecodeNullableNumber[float64](&ptrFloat32)
		assertNilError(t, err)
		assertEqual(t, float64(vFloat32), *ptr)

		vFloat64 := float64(2.2)
		ptr, err = DecodeNullableNumber[float64](&vFloat64)
		assertNilError(t, err)
		assertEqual(t, float64(vFloat64), *ptr)

		var vAny any = "test"
		_, err = DecodeNullableNumber[float64](&vAny)
		assertError(t, err, "failed to convert number, got: strconv.ParseFloat: parsing \"test\": invalid syntax")

		var vFn any = func() {}
		_, err = DecodeNullableNumber[float64](&vFn)
		assertError(t, err, "failed to convert number, got: strconv.ParseFloat")
	})

	t.Run("decode_nil", func(t *testing.T) {
		_, err := DecodeNumber[float64](nil)
		assertError(t, err, "number value must not be null")
	})

	t.Run("decode_invalid_type", func(t *testing.T) {
		_, err := DecodeNumber[float64](true)
		assertError(t, err, "malformed number")
	})

	t.Run("decode_string", func(t *testing.T) {
		value, err := DecodeNumber[int]("42")
		assertNilError(t, err)
		assertEqual(t, 42, value)

		value64, err := DecodeNumber[int64]("100")
		assertNilError(t, err)
		assertEqual(t, int64(100), value64)

		u64, err := DecodeNumber[uint64]("18446744073709551615")
		assertNilError(t, err)
		assertEqual(t, uint64(18446744073709551615), u64)

		fv, err := DecodeNumber[float64]("3.14")
		assertNilError(t, err)
		assertEqual(t, "3.14", fmt.Sprintf("%.2f", fv))

		_, err = DecodeNumber[int]("not-a-number")
		assertError(t, err, "malformed number")
	})

	t.Run("decode_string_pointer", func(t *testing.T) {
		s := "99"
		value, err := DecodeNumber[int](&s)
		assertNilError(t, err)
		assertEqual(t, 99, value)

		var nilStr *string
		_, err = DecodeNumber[int](nilStr)
		assertError(t, err, "number value must not be null")

		bad := "bad"
		_, err = DecodeNumber[int](&bad)
		assertError(t, err, "malformed number")
	})

	t.Run("decode_nullable_string", func(t *testing.T) {
		ptr, err := DecodeNullableNumber[int]("7")
		assertNilError(t, err)
		assertEqual(t, 7, *ptr)

		_, err = DecodeNullableNumber[int]("oops")
		assertError(t, err, "malformed number")
	})

	t.Run("decode_nullable_string_pointer", func(t *testing.T) {
		s := "55"
		ptr, err := DecodeNullableNumber[int](&s)
		assertNilError(t, err)
		assertEqual(t, 55, *ptr)

		var nilStr *string
		ptr, err = DecodeNullableNumber[int](nilStr)
		assertNilError(t, err)
		assertEqual(t, (*int)(nil), ptr)

		bad := "bad"
		_, err = DecodeNullableNumber[int](&bad)
		assertError(t, err, "malformed number")
	})
}

func TestDecodeNumberSlice(t *testing.T) {
	testNumbers := []any{int(1), int8(2), int16(3), int32(4), int64(5), uint(6), uint8(7), uint16(8), uint32(9), uint64(10)}
	result, err := DecodeNumberSlice[int](testNumbers)
	assertNilError(t, err)
	assertDeepEqual(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, result)

	t.Run("string_slice", func(t *testing.T) {
		strNums := []string{"1", "2", "3"}
		res, err := DecodeNumberSlice[int](strNums)
		assertNilError(t, err)
		assertDeepEqual(t, []int{1, 2, 3}, res)

		floats, err := DecodeNumberSlice[float64]([]string{"1.1", "2.2"})
		assertNilError(t, err)
		assertEqual(t, "1.10", fmt.Sprintf("%.2f", floats[0]))
		assertEqual(t, "2.20", fmt.Sprintf("%.2f", floats[1]))

		_, err = DecodeNumberSlice[int]([]string{"a", "b"})
		assertError(t, err, "malformed number")
	})
}

func TestAsBoolean(t *testing.T) {
	// Direct bool values work in both strict and non-strict
	value, err := AsBoolean(true)
	assertNilError(t, err)
	assertEqual(t, true, value)

	value, err = AsBoolean(false)
	assertNilError(t, err)
	assertEqual(t, false, value)

	boolPtr := true
	value, err = AsBoolean(&boolPtr)
	assertNilError(t, err)
	assertEqual(t, true, value)

	// nil returns error
	_, err = AsBoolean(nil)
	assertError(t, err, "boolean value must not be null")

	// nil *bool returns error
	var nilBool *bool
	_, err = AsBoolean(nilBool)
	assertError(t, err, "boolean value must not be null")

	// strings are rejected in strict mode
	_, err = AsBoolean("true")
	assertError(t, err, "malformed boolean")

	_, err = AsBoolean("false")
	assertError(t, err, "malformed boolean")

	s := "true"
	_, err = AsBoolean(&s)
	assertError(t, err, "malformed boolean")

	// invalid type
	_, err = AsBoolean(42)
	assertError(t, err, "malformed boolean")
}

func TestAsNullableBoolean(t *testing.T) {
	// Direct bool values work
	ptr, err := AsNullableBoolean(true)
	assertNilError(t, err)
	assertEqual(t, true, *ptr)

	boolPtr := false
	ptr, err = AsNullableBoolean(&boolPtr)
	assertNilError(t, err)
	assertEqual(t, false, *ptr)

	// nil returns nil pointer, no error
	ptr, err = AsNullableBoolean(nil)
	assertNilError(t, err)
	assertEqual(t, (*bool)(nil), ptr)

	var nilBool *bool
	ptr, err = AsNullableBoolean(nilBool)
	assertNilError(t, err)
	assertEqual(t, (*bool)(nil), ptr)

	// strings are rejected
	_, err = AsNullableBoolean("true")
	assertError(t, err, "malformed boolean")

	s := "true"
	_, err = AsNullableBoolean(&s)
	assertError(t, err, "malformed boolean")
}

func TestAsBooleanReflection(t *testing.T) {
	// Bool kind works
	value, err := AsBooleanReflection(reflect.ValueOf(true))
	assertNilError(t, err)
	assertEqual(t, true, value)

	// String kind is rejected in strict mode
	_, err = AsBooleanReflection(reflect.ValueOf("true"))
	assertError(t, err, "malformed boolean")

	// Invalid kind rejected
	_, err = AsBooleanReflection(reflect.ValueOf(42))
	assertError(t, err, "malformed boolean")
}

func TestAsNullableBooleanReflection(t *testing.T) {
	value, err := AsNullableBooleanReflection(reflect.ValueOf(true))
	assertNilError(t, err)
	assertEqual(t, true, *value)

	_, err = AsNullableBooleanReflection(reflect.ValueOf("true"))
	assertError(t, err, "malformed boolean")
}

func TestAsBooleanSlice(t *testing.T) {
	// Direct bool slices work
	value, err := AsBooleanSlice([]bool{true, false})
	assertNilError(t, err)
	assertDeepEqual(t, []bool{true, false}, value)

	// []any with actual bools works
	value, err = AsBooleanSlice([]any{true, false})
	assertNilError(t, err)
	assertDeepEqual(t, []bool{true, false}, value)

	// nil returns error (slice must not be null)
	_, err = AsBooleanSlice(nil)
	assertError(t, err, "boolean slice must not be null")

	// non-slice type rejected
	_, err = AsBooleanSlice("true")
	assertError(t, err, "malformed boolean slice; got: string")

	// []any with strings rejected in strict mode
	_, err = AsBooleanSlice([]any{"true"})
	assertError(t, err, "malformed boolean")

	// nil element in slice rejected
	_, err = AsBooleanSlice([]any{nil})
	assertError(t, err, "element 0: boolean value must not be null")
}

func TestAsNullableBooleanSlice(t *testing.T) {
	value, err := AsNullableBooleanSlice([]bool{true, false})
	assertNilError(t, err)
	assertDeepEqual(t, []bool{true, false}, *value)

	// nil returns nil pointer, no error
	value, err = AsNullableBooleanSlice(nil)
	assertNilError(t, err)
	assertEqual(t, (*[]bool)(nil), value)

	// []any with actual bools works
	value, err = AsNullableBooleanSlice([]any{true, false})
	assertNilError(t, err)
	assertDeepEqual(t, []bool{true, false}, *value)

	// []any with strings rejected
	_, err = AsNullableBooleanSlice([]any{"true"})
	assertError(t, err, "malformed boolean")

	// non-slice rejected
	_, err = AsNullableBooleanSlice("true")
	assertError(t, err, "malformed boolean slice; got: string")
}

func TestAsNumber(t *testing.T) {
	// Direct numeric types work
	value, err := AsNumber[int](42)
	assertNilError(t, err)
	assertEqual(t, 42, value)

	value64, err := AsNumber[int64](int64(100))
	assertNilError(t, err)
	assertEqual(t, int64(100), value64)

	fv, err := AsNumber[float64](3.14)
	assertNilError(t, err)
	assertEqual(t, "3.14", fmt.Sprintf("%.2f", fv))

	// nil rejected
	_, err = AsNumber[int](nil)
	assertError(t, err, "number value must not be null")

	// strings are rejected in strict mode
	_, err = AsNumber[int]("42")
	assertError(t, err, "malformed number")

	s := "42"
	_, err = AsNumber[int](&s)
	assertError(t, err, "malformed number")

	// invalid type rejected
	_, err = AsNumber[int](true)
	assertError(t, err, "malformed number")
}

func TestAsNullableNumber(t *testing.T) {
	// Direct numeric types work
	ptr, err := AsNullableNumber[int](42)
	assertNilError(t, err)
	assertEqual(t, 42, *ptr)

	// nil returns nil, no error
	ptr, err = AsNullableNumber[int](nil)
	assertNilError(t, err)
	assertEqual(t, (*int)(nil), ptr)

	// strings rejected
	_, err = AsNullableNumber[int]("42")
	assertError(t, err, "malformed number")

	s := "42"
	_, err = AsNullableNumber[int](&s)
	assertError(t, err, "malformed number")
}

func TestAsNumberReflection(t *testing.T) {
	// Numeric kinds work
	value, err := AsNumberReflection[int](reflect.ValueOf(42))
	assertNilError(t, err)
	assertEqual(t, 42, value)

	value64, err := AsNumberReflection[float64](reflect.ValueOf(3.14))
	assertNilError(t, err)
	assertEqual(t, "3.14", fmt.Sprintf("%.2f", value64))

	// String kind rejected in strict mode
	_, err = AsNumberReflection[int](reflect.ValueOf("42"))
	assertError(t, err, "malformed number")

	// Invalid kind rejected
	_, err = AsNumberReflection[int](reflect.ValueOf(true))
	assertError(t, err, "malformed number")
}

func TestAsNullableNumberReflection(t *testing.T) {
	ptr, err := AsNullableNumberReflection[int](reflect.ValueOf(42))
	assertNilError(t, err)
	assertEqual(t, 42, *ptr)

	// String kind rejected
	_, err = AsNullableNumberReflection[int](reflect.ValueOf("42"))
	assertError(t, err, "malformed number")
}

func TestAsNumberSlice(t *testing.T) {
	// Direct int slices work
	result, err := AsNumberSlice[int]([]int{1, 2, 3})
	assertNilError(t, err)
	assertDeepEqual(t, []int{1, 2, 3}, result)

	// []any with numeric values works
	result, err = AsNumberSlice[int]([]any{int(1), int64(2), float64(3)})
	assertNilError(t, err)
	assertDeepEqual(t, []int{1, 2, 3}, result)

	// nil returns nil, no error
	result, err = AsNumberSlice[int](nil)
	assertNilError(t, err)
	assertDeepEqual(t, ([]int)(nil), result)

	// []string rejected in strict mode
	_, err = AsNumberSlice[int]([]string{"1", "2", "3"})
	assertError(t, err, "malformed number slice")

	// invalid type rejected
	_, err = AsNumberSlice[int](true)
	assertError(t, err, "malformed number slice")
}

func TestAsNumberSliceReflection(t *testing.T) {
	result, err := AsNumberSliceReflection[int](reflect.ValueOf([]int{1, 2, 3}))
	assertNilError(t, err)
	assertDeepEqual(t, []int{1, 2, 3}, result)

	// []string rejected
	_, err = AsNumberSliceReflection[int](reflect.ValueOf([]string{"1", "2"}))
	assertError(t, err, "malformed number")
}
