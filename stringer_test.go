package goutils

import (
	"errors"
	"fmt"
	"math"
	"testing"
	"time"
)

type customStringer struct {
	val string
}

func (c customStringer) String() string {
	return c.val
}

type customError struct{}

func (c customError) Error() string {
	return "custom error"
}

func TestToString_Primitives(t *testing.T) {
	now := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	duration := 2 * time.Hour

	tests := []struct {
		name       string
		value      any
		emptyValue string
		want       string
		wantErr    bool
	}{
		{"nil", nil, "empty", "empty", false},
		{"bool true", true, "empty", "true", false},
		{"bool false", false, "empty", "false", false},
		{"string", "abc", "empty", "abc", false},
		{"int", int(-42), "empty", "-42", false},
		{"int8", int8(-8), "empty", "-8", false},
		{"int16", int16(-16), "empty", "-16", false},
		{"int32", int32(-32), "empty", "-32", false},
		{"int64", int64(-64), "empty", "-64", false},
		{"uint", uint(42), "empty", "42", false},
		{"uint8", uint8(8), "empty", "8", false},
		{"uint16", uint16(16), "empty", "16", false},
		{"uint32", uint32(32), "empty", "32", false},
		{"uint64", uint64(64), "empty", "64", false},
		{"float32", float32(3.14), "empty", "3.14", false},
		{"float64", float64(-2.718), "empty", "-2.718", false},
		{"complex64", complex64(1 + 2i), "empty", "(1+2i)", false},
		{"complex128", complex128(-3 + 4i), "empty", "(-3+4i)", false},
		{"time.Time", now, "empty", now.Format(time.RFC3339), false},
		{"time.Duration", duration, "empty", duration.String(), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToString(tt.value, tt.emptyValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
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
		name       string
		value      any
		emptyValue string
		want       string
	}{
		{"*string", &s, "empty", "hello"},
		{"*string nil", (*string)(nil), "empty", "empty"},
		{"*int", &i, "empty", "123"},
		{"*int nil", (*int)(nil), "empty", "empty"},
		{"*float64", &f, "empty", "1.23"},
		{"*float64 nil", (*float64)(nil), "empty", "empty"},
		{"*complex128", &c, "empty", "(1+2i)"},
		{"*complex128 nil", (*complex128)(nil), "empty", "empty"},
		{"*time.Time", &now, "empty", now.Format(time.RFC3339)},
		{"*time.Time nil", (*time.Time)(nil), "empty", "empty"},
		{"*time.Duration", &d, "empty", d.String()},
		{"*time.Duration nil", (*time.Duration)(nil), "empty", "empty"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToString(tt.value, tt.emptyValue)
			if err != nil {
				t.Errorf("ToString() error = %v", err)
				return
			}
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
		A int
		B string
	}
	ts := testStruct{A: 1, B: "b"}

	tests := []struct {
		name       string
		value      any
		emptyValue string
		want       string
	}{
		{"fmt.Stringer", cs, "empty", "custom"},
		{"fmt.Stringer nil", nilStringer, "empty", "empty"},
		{"customStringer nil", any((*customStringer)(nil)), "empty", "empty"},
		{"struct fallback", ts, "empty", `{"A":1,"B":"b"}`},
		{"slice fallback", []int{1, 2, 3}, "empty", `[1,2,3]`},
		{"map fallback", map[string]int{"a": 1}, "empty", `{"a":1}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToString(tt.value, tt.emptyValue)
			if err != nil {
				t.Errorf("ToString() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToString_JSONMarshalError(t *testing.T) {
	ch := make(chan int)
	_, err := ToString(ch, "empty")
	if err == nil {
		t.Error("ToString() expected error for unmarshalable type, got nil")
	}
}

func TestToDebugString(t *testing.T) {

	ch := make(chan int)
	tests := []struct {
		name       string
		value      any
		emptyValue string
		want       string
	}{
		{"nil", nil, "empty", "empty"},
		{"string", "abc", "empty", "abc"},
		{"struct fallback", struct{ X int }{X: 1}, "empty", `{"X":1}`},
		{"unmarshalable fallback", ch, "empty", fmt.Sprint(ch)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToDebugString(tt.value, tt.emptyValue)
			if got != tt.want {
				t.Errorf("ToDebugString() = %v, want %v", got, tt.want)
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
			got, err := ToString(tt.value, "")
			if err != nil {
				t.Errorf("ToString() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToString_ErrorType(t *testing.T) {
	errVal := errors.New("err")
	got, err := ToString(errVal, "empty")
	if err != nil {
		t.Errorf("ToString() error = %v", err)
	}
	if got != "{}" {
		t.Errorf("ToString() = %v, want {}", got)
	}
}
