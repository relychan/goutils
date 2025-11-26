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
	assertError(t, err, "the boolean value must not be null")

	_, err = DecodeBoolean("failure")
	assertError(t, err, "malformed boolean; got: string")
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
	assertError(t, err, "element 0: the boolean value must not be null")

	_, err = DecodeBooleanSlice(&[]any{nil})
	assertError(t, err, "element 0: the boolean value must not be null")
}

func TestDecodeNullableBooleanSlice(t *testing.T) {
	value, err := DecodeNullableBooleanSlice([]bool{true, false})
	assertNilError(t, err)
	assertDeepEqual(t, []bool{true, false}, *value)

	value, err = DecodeNullableBooleanSlice([]*bool{ToPtr(true), ToPtr(false)})
	assertNilError(t, err)
	assertDeepEqual(t, []bool{true, false}, *value)

	value, err = DecodeNullableBooleanSlice(&[]*bool{ToPtr(true), ToPtr(false)})
	assertNilError(t, err)
	assertDeepEqual(t, []bool{true, false}, *value)

	value, err = DecodeNullableBooleanSlice([]any{false, true})
	assertNilError(t, err)
	assertDeepEqual(t, []bool{false, true}, *value)

	value, err = DecodeNullableBooleanSlice(&[]any{true})
	assertNilError(t, err)
	assertDeepEqual(t, []bool{true}, *value)

	_, err = DecodeNullableBooleanSlice([]any{nil})
	assertError(t, err, "element 0: the boolean value must not be null")

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

	value, err = DecodeStringSlice([]*string{ToPtr("foo"), ToPtr("bar")})
	assertNilError(t, err)
	assertDeepEqual(t, []string{"foo", "bar"}, value)

	value, err = DecodeStringSlice(&[]*string{ToPtr("foo")})
	assertNilError(t, err)
	assertDeepEqual(t, []string{"foo"}, value)

	value, err = DecodeStringSlice([]any{"bar", "foo"})
	assertNilError(t, err)
	assertDeepEqual(t, []string{"bar", "foo"}, value)

	value, err = DecodeStringSlice(&[]any{"foo"})
	assertNilError(t, err)
	assertDeepEqual(t, []string{"foo"}, value)

	value, err = DecodeStringSlice(nil)
	assertDeepEqual(t, nil, value)

	_, err = DecodeStringSlice("failure")
	assertError(t, err, "malformed string slice; got: string")

	_, err = DecodeStringSlice(time.Now())
	assertError(t, err, "malformed string slice; got: struct")

	_, err = DecodeStringSlice([]any{nil})
	assertError(t, err, "failed to decode string at 0: the string value must not be null")

	_, err = DecodeStringSlice(&[]any{nil})
	assertError(t, err, "failed to decode element at 0: the string value must not be null")
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
		assertError(t, err, "the number value must not be null")
	})

	t.Run("decode_invalid_type", func(t *testing.T) {
		_, err := DecodeNumber[float64]("failure")
		assertError(t, err, "malformed number; got: string")
	})
}

func TestDecodeNumberSlice(t *testing.T) {
	testNumbers := []any{int(1), int8(2), int16(3), int32(4), int64(5), uint(6), uint8(7), uint16(8), uint32(9), uint64(10)}
	result, err := DecodeNumberSlice[int](testNumbers)
	assertNilError(t, err)
	assertDeepEqual(t, []int{(1), (2), (3), (4), (5), (6), (7), (8), (9), (10)}, result)
}
