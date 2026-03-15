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
	"math"
	"strings"
	"testing"
	"time"
)

type customStringer struct {
	val string
}

func (c customStringer) String() string {
	return c.val
}

func TestToString_Primitives(t *testing.T) {
	now := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	duration := 2 * time.Hour

	tests := []struct {
		name    string
		value   any
		want    string
		wantErr bool
	}{
		{"nil", nil, NullStr, false},
		{"bool true", true, "true", false},
		{"bool false", false, "false", false},
		{"string", "abc", "abc", false},
		{"int", int(-42), "-42", false},
		{"int8", int8(-8), "-8", false},
		{"int16", int16(-16), "-16", false},
		{"int32", int32(-32), "-32", false},
		{"int64", int64(-64), "-64", false},
		{"uint", uint(42), "42", false},
		{"uint8", uint8(8), "8", false},
		{"uint16", uint16(16), "16", false},
		{"uint32", uint32(32), "32", false},
		{"uint64", uint64(64), "64", false},
		{"float32", float32(3.14), "3.14", false},
		{"float64", float64(-2.718), "-2.718", false},
		{"complex64", complex64(1 + 2i), "(1+2i)", false},
		{"complex128", complex128(-3 + 4i), "(-3+4i)", false},
		{"time.Time", now, now.Format(time.RFC3339), false},
		{"time.Duration", duration, duration.String(), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToString(tt.value)
			if got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToString_Pointers(t *testing.T) {
	s := "hello"
	i := 123
	f := 1.23
	c := complex(1, 2)
	now := time.Now()
	d := time.Minute

	tests := []struct {
		name  string
		value any
		want  string
	}{
		{"*string", &s, "hello"},
		{"*string nil", (*string)(nil), NullStr},
		{"*int", &i, "123"},
		{"*int nil", (*int)(nil), NullStr},
		{"*float64", &f, "1.23"},
		{"*float64 nil", (*float64)(nil), NullStr},
		{"*complex128", &c, "(1+2i)"},
		{"*complex128 nil", (*complex128)(nil), NullStr},
		{"*time.Time", &now, now.Format(time.RFC3339)},
		{"*time.Time nil", (*time.Time)(nil), NullStr},
		{"*time.Duration", &d, d.String()},
		{"*time.Duration nil", (*time.Duration)(nil), NullStr},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToString(tt.value)
			if got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToString_StringerAndFallback(t *testing.T) {
	cs := customStringer{"custom"}
	var nilStringer fmt.Stringer
	type testStruct struct {
		A int    `yaml:"a"`
		B string `yaml:"b,omitempty"`
	}
	ts := testStruct{A: 1, B: "b"}

	tests := []struct {
		name  string
		value any
		want  string
	}{
		{"fmt.Stringer", cs, "custom"},
		{"fmt.Stringer nil", nilStringer, NullStr},
		{"customStringer nil", any((*customStringer)(nil)), NullStr},
		{"struct fallback", ts, "\na: 1\nb: b"},
		{"slice fallback", []int{1, 2, 3}, "- 1\n- 2\n- 3"},
		{"map fallback", map[string]int{"a": 1}, "\na: 1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToString(tt.value)
			if got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToString_SpecialCases(t *testing.T) {
	// Test NaN, +Inf, -Inf for float64
	tests := []struct {
		name  string
		value any
		want  string
	}{
		{"float64 NaN", math.NaN(), "NaN"},
		{"float64 +Inf", math.Inf(1), "+Inf"},
		{"float64 -Inf", math.Inf(-1), "-Inf"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToString(tt.value)
			if got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToString_NilPointers(t *testing.T) {
	tests := []struct {
		name  string
		value any
	}{
		{"*bool nil", (*bool)(nil)},
		{"*int8 nil", (*int8)(nil)},
		{"*int16 nil", (*int16)(nil)},
		{"*int32 nil", (*int32)(nil)},
		{"*int64 nil", (*int64)(nil)},
		{"*uint nil", (*uint)(nil)},
		{"*uint8 nil", (*uint8)(nil)},
		{"*uint16 nil", (*uint16)(nil)},
		{"*uint32 nil", (*uint32)(nil)},
		{"*uint64 nil", (*uint64)(nil)},
		{"*float32 nil", (*float32)(nil)},
		{"*complex64 nil", (*complex64)(nil)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToString(tt.value)
			if got != NullStr {
				t.Errorf("ToString() = %v, want empty", got)
			}
		})
	}
}

func TestToString_TypedSlices(t *testing.T) {
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name  string
		value any
		want  string
	}{
		{"[]bool", []bool{true, false}, "- true\n- false"},
		{"[]string", []string{"a", "b"}, "- a\n- b"},
		{"[]int", []int{1, 2}, "- 1\n- 2"},
		{"[]int8", []int8{1, 2}, "- 1\n- 2"},
		{"[]int16", []int16{1, 2}, "- 1\n- 2"},
		{"[]int32", []int32{1, 2}, "- 1\n- 2"},
		{"[]int64", []int64{1, 2}, "- 1\n- 2"},
		{"[]uint", []uint{1, 2}, "- 1\n- 2"},
		{"[]uint8", []uint8{1, 2}, "- 1\n- 2"},
		{"[]uint16", []uint16{1, 2}, "- 1\n- 2"},
		{"[]uint32", []uint32{1, 2}, "- 1\n- 2"},
		{"[]uint64", []uint64{1, 2}, "- 1\n- 2"},
		{"[]float32", []float32{1.5, 2.5}, "- 1.5\n- 2.5"},
		{"[]float64", []float64{1.5, 2.5}, "- 1.5\n- 2.5"},
		{"[]complex64", []complex64{1 + 2i}, "- (1+2i)"},
		{"[]complex128", []complex128{3 + 4i}, "- (3+4i)"},
		{"[]time.Time", []time.Time{now}, "- " + now.Format(time.RFC3339)},
		{"[]fmt.Stringer", []fmt.Stringer{customStringer{"x"}}, "- x"},
		{"[]any", []any{1, "two"}, "- 1\n- two"},
		{"empty []int", []int{}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToString(tt.value)
			if got != tt.want {
				t.Errorf("ToString() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestToString_Maps(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  string
	}{
		{"map[string]string", map[string]string{"k": "v"}, "k: v"},
		{"map[string]any", map[string]any{"k": 1}, "k: 1"},
		{"empty map[string]string", map[string]string{}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToString(tt.value)
			if got != tt.want {
				t.Errorf("ToString() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestToString_StructTags(t *testing.T) {
	type withYAML struct {
		Name string `yaml:"name"`
		Age  int    `yaml:"age"`
	}
	type withJSON struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	type withUnexported struct {
		Name    string `yaml:"name"`
		private string //nolint:unused
	}

	t.Run("yaml tags", func(t *testing.T) {
		got := ToString(withYAML{Name: "alice", Age: 30})
		if !strings.Contains(got, "name: alice") || !strings.Contains(got, "age: 30") {
			t.Errorf("unexpected output: %q", got)
		}
	})

	t.Run("json tags fallback", func(t *testing.T) {
		got := ToString(withJSON{Name: "bob", Age: 25})
		if !strings.Contains(got, "name: bob") || !strings.Contains(got, "age: 25") {
			t.Errorf("unexpected output: %q", got)
		}
	})

	t.Run("unexported fields omitted", func(t *testing.T) {
		got := ToString(withUnexported{Name: "carol", private: "secret"})
		if strings.Contains(got, "secret") || strings.Contains(got, "private") {
			t.Errorf("unexported field leaked into output: %q", got)
		}
		if !strings.Contains(got, "name: carol") {
			t.Errorf("unexpected output: %q", got)
		}
	})
}

func TestToString_NestedStruct(t *testing.T) {
	type Inner struct {
		Val int `yaml:"val"`
	}
	type Outer struct {
		Name  string `yaml:"name"`
		Inner Inner  `yaml:"inner"`
	}

	got := ToString(Outer{Name: "outer", Inner: Inner{Val: 42}})
	if !strings.Contains(got, "name: outer") {
		t.Errorf("missing name field in output: %q", got)
	}
	if !strings.Contains(got, "val: 42") {
		t.Errorf("missing nested val field in output: %q", got)
	}
}

func TestStringContainsCTLByte(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"empty string", "", false},
		{"normal ascii", "hello world", false},
		{"with newline", "hello\nworld", true},
		{"with tab", "hello\tworld", true},
		{"with carriage return", "hello\rworld", true},
		{"with null byte", "hello\x00world", true},
		{"with DEL (0x7f)", "hello\x7fworld", true},
		{"with space (boundary)", "hello world", false},
		{"with tilde (printable)", "hello~world", false},
		{"control at start", "\x01hello", true},
		{"control at end", "hello\x1f", true},
		{"all printable", "!\"#$%&'()*+,-./0-9:;<=>?@A-Z[\\]^_`a-z{|}~", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StringContainsCTLByte(tt.input)
			if got != tt.want {
				t.Errorf("StringContainsCTLByte(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
