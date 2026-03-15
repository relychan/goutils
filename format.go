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
	"cmp"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"go.yaml.in/yaml/v4"
)

const nullStr = "null"

// ToString converts an arbitrary value to its string representation.
//
// It handles the following cases:
//   - For nil values, it returns the provided emptyValue.
//   - For primitive types (bool, string, integers, floats, complex), it uses the appropriate formatting.
//   - For time.Time, time.Duration, and their pointer types, it uses the standard time formatting.
//   - For types implementing fmt.Stringer, it uses their String() method.
//   - For pointers, it dereferences and formats the underlying value, or returns emptyValue if nil.
//   - For unsupported types, it attempts to marshal the value to JSON.
//
// If JSON marshaling fails, it returns an error.
//
// The emptyValue parameter specifies the string to return for nil values or nil pointers.
func ToString(value any, emptyValue string) string {
	var sb strings.Builder

	buildStringIndent(&sb, value, emptyValue, 0)

	return sb.String()
}

// PrintMap prints a map value to a human-readable format.
func PrintMap[K cmp.Ordered, V any](values map[K]V, indent int) string {
	if len(values) == 0 {
		return ""
	}

	var sb strings.Builder

	BuildMapToString(&sb, values, indent)

	return sb.String()
}

// BuildMapToString builds a map value to a human-readable format.
func BuildMapToString[K cmp.Ordered, V any](sb *strings.Builder, values map[K]V, indent int) {
	first := true
	prefix := strings.Repeat(" ", indent)

	for key, value := range values {
		strKey, ok := FormatScalar(key, "")
		if !ok || strKey == "" {
			continue
		}

		if !first {
			sb.WriteByte('\n')
			sb.WriteString(prefix)
		} else {
			first = false
		}

		sb.WriteString(strKey)
		sb.WriteString(": ")
		sb.WriteString(ToString(value, nullStr))
	}
}

// PrintSlice prints a slice value to a human-readable format.
func PrintSlice[V any](values []V, indent int) string {
	if len(values) == 0 {
		return ""
	}

	var sb strings.Builder

	BuildSliceToString(&sb, values, indent)

	return sb.String()
}

// BuildSliceToString builds a slice value to a human-readable format.
func BuildSliceToString[V any](sb *strings.Builder, values []V, indent int) {
	prefix := strings.Repeat(" ", indent)

	for i, value := range values {
		if i > 0 {
			sb.WriteByte('\n')
			sb.WriteString(prefix)
		}

		sb.WriteString("- ")
		sb.WriteString(ToString(value, nullStr))
	}
}

// FormatScalar converts an arbitrary value to its string representation.
//
// It handles the following cases:
//   - For nil values, it returns the provided emptyValue.
//   - For primitive types (bool, string, integers, floats, complex), it uses the appropriate formatting.
//   - For time.Time and *time.Time, it uses the standard time formatting.
//   - For types implementing fmt.Stringer, it uses their String() method.
//   - For pointers, it dereferences and formats the underlying value, or returns emptyValue if nil.
//   - For unsupported types, return an empty string and a false value.
//
// The emptyValue parameter specifies the string to return for nil values or nil pointers.
func FormatScalar( //nolint:cyclop,gocognit,gocyclo,funlen,maintidx
	value any,
	emptyValue string,
) (string, bool) {
	if value == nil {
		return emptyValue, true
	}

	switch typedValue := value.(type) {
	case bool:
		return strconv.FormatBool(typedValue), true
	case string:
		return typedValue, true
	case int:
		return strconv.FormatInt(int64(typedValue), 10), true
	case int8:
		return strconv.FormatInt(int64(typedValue), 10), true
	case int16:
		return strconv.FormatInt(int64(typedValue), 10), true
	case int32:
		return strconv.FormatInt(int64(typedValue), 10), true
	case int64:
		return strconv.FormatInt(typedValue, 10), true
	case uint:
		return strconv.FormatUint(uint64(typedValue), 10), true
	case uint8:
		return strconv.FormatUint(uint64(typedValue), 10), true
	case uint16:
		return strconv.FormatUint(uint64(typedValue), 10), true
	case uint32:
		return strconv.FormatUint(uint64(typedValue), 10), true
	case uint64:
		return strconv.FormatUint(typedValue, 10), true
	case float32:
		return strconv.FormatFloat(float64(typedValue), 'f', -1, 32), true
	case float64:
		return strconv.FormatFloat(typedValue, 'f', -1, 64), true
	case complex64:
		return strconv.FormatComplex(complex128(typedValue), 'f', -1, 64), true
	case complex128:
		return strconv.FormatComplex(typedValue, 'f', -1, 128), true
	case time.Time:
		return typedValue.Format(time.RFC3339), true
	case time.Duration:
		return typedValue.String(), true
	case *bool:
		if typedValue == nil {
			return emptyValue, true
		}

		return strconv.FormatBool(*typedValue), true
	case *string:
		if typedValue == nil {
			return emptyValue, true
		}

		return *typedValue, true
	case *int:
		if typedValue == nil {
			return emptyValue, true
		}

		return strconv.FormatInt(int64(*typedValue), 10), true
	case *int8:
		if typedValue == nil {
			return emptyValue, true
		}

		return strconv.FormatInt(int64(*typedValue), 10), true
	case *int16:
		if typedValue == nil {
			return emptyValue, true
		}

		return strconv.FormatInt(int64(*typedValue), 10), true
	case *int32:
		if typedValue == nil {
			return emptyValue, true
		}

		return strconv.FormatInt(int64(*typedValue), 10), true
	case *int64:
		if typedValue == nil {
			return emptyValue, true
		}

		return strconv.FormatInt(*typedValue, 10), true
	case *uint:
		if typedValue == nil {
			return emptyValue, true
		}

		return strconv.FormatUint(uint64(*typedValue), 10), true
	case *uint8:
		if typedValue == nil {
			return emptyValue, true
		}

		return strconv.FormatUint(uint64(*typedValue), 10), true
	case *uint16:
		if typedValue == nil {
			return emptyValue, true
		}

		return strconv.FormatUint(uint64(*typedValue), 10), true
	case *uint32:
		if typedValue == nil {
			return emptyValue, true
		}

		return strconv.FormatUint(uint64(*typedValue), 10), true
	case *uint64:
		if typedValue == nil {
			return emptyValue, true
		}

		return strconv.FormatUint(*typedValue, 10), true
	case *float32:
		if typedValue == nil {
			return emptyValue, true
		}

		return strconv.FormatFloat(float64(*typedValue), 'f', -1, 32), true
	case *float64:
		if typedValue == nil {
			return emptyValue, true
		}

		return strconv.FormatFloat(*typedValue, 'f', -1, 64), true
	case *complex64:
		if typedValue == nil {
			return emptyValue, true
		}

		return strconv.FormatComplex(complex128(*typedValue), 'f', -1, 64), true
	case *complex128:
		if typedValue == nil {
			return emptyValue, true
		}

		return strconv.FormatComplex(*typedValue, 'f', -1, 128), true
	case *time.Time:
		if typedValue == nil {
			return emptyValue, true
		}

		return typedValue.Format(time.RFC3339), true
	case *time.Duration:
		if typedValue == nil {
			return emptyValue, true
		}

		return typedValue.String(), true
	case fmt.Stringer:
		if IsNil(typedValue) {
			return emptyValue, true
		}

		return typedValue.String(), true
	default:
		return "", false
	}
}

// StringContainsCTLByte reports whether s contains any ASCII control character.
func StringContainsCTLByte(s string) bool {
	for i := range len(s) {
		b := s[i]
		if b < ' ' || b == 0x7f {
			return true
		}
	}

	return false
}

func buildStringIndentRefection( //nolint:cyclop,funlen
	sb *strings.Builder,
	value reflect.Value,
	indent int,
) bool {
	reflectValue, notNull := UnwrapPointerFromReflectValue(value)
	if !notNull {
		sb.WriteString(nullStr)

		return true
	}

	switch reflectValue.Kind() {
	case reflect.Bool:
		sb.WriteString(strconv.FormatBool(reflectValue.Bool()))
	case reflect.String:
		sb.WriteString(reflectValue.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		sb.WriteString(strconv.FormatInt(reflectValue.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		sb.WriteString(strconv.FormatUint(reflectValue.Uint(), 10))
	case reflect.Float32:
		sb.WriteString(strconv.FormatFloat(reflectValue.Float(), 'f', -1, 32))
	case reflect.Float64:
		sb.WriteString(strconv.FormatFloat(reflectValue.Float(), 'f', -1, 64))
	case reflect.Complex64:
		sb.WriteString(strconv.FormatComplex(reflectValue.Complex(), 'f', -1, 64))
	case reflect.Complex128:
		sb.WriteString(strconv.FormatComplex(reflectValue.Complex(), 'f', -1, 128))
	case reflect.Slice, reflect.Array:
		valueLength := reflectValue.Len()
		prefix := strings.Repeat(" ", indent)

		for i := range valueLength {
			elem := reflectValue.Index(i)

			sb.WriteByte('\n')
			sb.WriteString(prefix)
			sb.WriteString("- ")
			buildStringIndentRefection(sb, elem, indent+2)
		}
	case reflect.Map:
		prefix := strings.Repeat(" ", indent)
		keys := reflectValue.MapKeys()

		for _, key := range keys {
			sb.WriteByte('\n')
			sb.WriteString(prefix)
			buildStringIndentRefection(sb, key, 0)
			sb.WriteString(": ")

			mapValue := reflectValue.MapIndex(key)
			buildStringIndentRefection(sb, mapValue, indent+2)
		}
	case reflect.Struct:
		prefix := strings.Repeat(" ", indent)

		for field, itemValue := range reflectValue.Fields() {
			if !field.IsExported() {
				continue
			}

			tag, ok := field.Tag.Lookup("yaml")
			if !ok || tag == "" {
				tag, _ = field.Tag.Lookup("json")
			}

			key := field.Name

			if tag != "" {
				parts := strings.Split(tag, ",")
				if len(parts) > 0 {
					key = parts[0]
				}
			}

			sb.WriteByte('\n')
			sb.WriteString(prefix)
			sb.WriteString(key)
			sb.WriteString(": ")

			buildStringIndentRefection(sb, itemValue, indent+2)
		}
	case reflect.Func, reflect.Chan, reflect.Invalid:
		return false
	default:
		rawBytes, err := yaml.Dump(value.Interface())
		if err != nil {
			return false
		}

		sb.Write(rawBytes)
	}

	return true
}

func buildStringIndent( //nolint:cyclop,funlen
	sb *strings.Builder,
	value any,
	emptyValue string,
	indent int,
) bool {
	result, ok := FormatScalar(value, emptyValue)
	if ok {
		sb.WriteString(result)

		return true
	}

	switch typedValue := value.(type) {
	case []any:
		sb.WriteString(PrintSlice(typedValue, indent))

		return true
	case []bool:
		sb.WriteString(PrintSlice(typedValue, indent))

		return true
	case []string:
		sb.WriteString(PrintSlice(typedValue, indent))

		return true
	case []int:
		sb.WriteString(PrintSlice(typedValue, indent))

		return true
	case []int8:
		sb.WriteString(PrintSlice(typedValue, indent))

		return true
	case []int16:
		sb.WriteString(PrintSlice(typedValue, indent))

		return true
	case []int32:
		sb.WriteString(PrintSlice(typedValue, indent))

		return true
	case []int64:
		sb.WriteString(PrintSlice(typedValue, indent))

		return true
	case []uint:
		sb.WriteString(PrintSlice(typedValue, indent))

		return true
	case []uint8:
		sb.WriteString(PrintSlice(typedValue, indent))

		return true
	case []uint16:
		sb.WriteString(PrintSlice(typedValue, indent))

		return true
	case []uint32:
		sb.WriteString(PrintSlice(typedValue, indent))

		return true
	case []uint64:
		sb.WriteString(PrintSlice(typedValue, indent))

		return true
	case []float32:
		sb.WriteString(PrintSlice(typedValue, indent))

		return true
	case []float64:
		sb.WriteString(PrintSlice(typedValue, indent))

		return true
	case []complex64:
		sb.WriteString(PrintSlice(typedValue, indent))

		return true
	case []complex128:
		sb.WriteString(PrintSlice(typedValue, indent))

		return true
	case []time.Time:
		sb.WriteString(PrintSlice(typedValue, indent))

		return true
	case []fmt.Stringer:
		sb.WriteString(PrintSlice(typedValue, indent))

		return true
	case map[string]string:
		sb.WriteString(PrintMap(typedValue, indent))

		return true
	case map[string]any:
		sb.WriteString(PrintMap(typedValue, indent))

		return true
	default:
		return buildStringIndentRefection(sb, reflect.ValueOf(value), indent)
	}
}
