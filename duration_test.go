package goutils

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"go.yaml.in/yaml/v4"
)

func TestDuration(t *testing.T) {
	dur := Duration(time.Second * 10)
	assertEqual(t, "duration", dur.Type())
	assertNilError(t, dur.Set("10m"))
	assertEqual(t, Duration(10*time.Minute), dur)
}

func TestParseDuration(t *testing.T) {
	testCases := []struct {
		in             string
		out            time.Duration
		expectedString string
	}{
		{
			in:             "0",
			out:            0,
			expectedString: "0s",
		},
		{
			in:             "0w",
			out:            0,
			expectedString: "0s",
		},
		{
			in:             "0s",
			out:            0,
			expectedString: "",
		},
		{
			in:             "324ms",
			out:            324 * time.Millisecond,
			expectedString: "",
		},
		{
			in:             "3s",
			out:            3 * time.Second,
			expectedString: "",
		},
		{
			in:             "5m",
			out:            5 * time.Minute,
			expectedString: "",
		},
		{
			in:             "1h",
			out:            time.Hour,
			expectedString: "",
		},
		{
			in:             "4d",
			out:            4 * 24 * time.Hour,
			expectedString: "",
		},
		{
			in:             "4d1h",
			out:            4*24*time.Hour + time.Hour,
			expectedString: "",
		},
		{
			in:             "14d",
			out:            14 * 24 * time.Hour,
			expectedString: "2w",
		},
		{
			in:             "3w",
			out:            3 * 7 * 24 * time.Hour,
			expectedString: "",
		},
		{
			in:             "3w2d1h",
			out:            3*7*24*time.Hour + 2*24*time.Hour + time.Hour,
			expectedString: "23d1h",
		},
		{
			in:             "10y",
			out:            10 * 365 * 24 * time.Hour,
			expectedString: "",
		},
	}

	for _, c := range testCases {
		d, err := ParseDuration(c.in)

		if err != nil {
			t.Errorf("Unexpected error on input %q", c.in)
		}

		assertEqual(t, c.out, time.Duration(d))

		expectedString := c.expectedString
		if expectedString == "" {
			expectedString = c.in
		}

		assertEqual(t, expectedString, d.String())
	}
}

func TestDuration_UnmarshalTextAndYAML(t *testing.T) {
	cases := []struct {
		in             string
		out            time.Duration
		expectedString string
	}{
		{
			in:             "0",
			out:            0,
			expectedString: "0s",
		}, {
			in:             "0w",
			out:            0,
			expectedString: "0s",
		}, {
			in:  "0s",
			out: 0,
		}, {
			in:  "324ms",
			out: 324 * time.Millisecond,
		}, {
			in:  "3s",
			out: 3 * time.Second,
		}, {
			in:  "5m",
			out: 5 * time.Minute,
		}, {
			in:  "1h",
			out: time.Hour,
		}, {
			in:  "4d",
			out: 4 * 24 * time.Hour,
		}, {
			in:  "4d1h",
			out: 4*24*time.Hour + time.Hour,
		}, {
			in:             "14d",
			out:            14 * 24 * time.Hour,
			expectedString: "2w",
		}, {
			in:  "3w",
			out: 3 * 7 * 24 * time.Hour,
		}, {
			in:             "3w2d1h",
			out:            3*7*24*time.Hour + 2*24*time.Hour + time.Hour,
			expectedString: "23d1h",
		}, {
			in:  "10y",
			out: 10 * 365 * 24 * time.Hour,
		},
	}

	for _, c := range cases {
		var d Duration

		err := d.UnmarshalText([]byte(c.in))
		if err != nil {
			t.Errorf("Unexpected error on input %q", c.in)
		}

		assertEqual(t, c.out, time.Duration(d))

		expectedString := c.expectedString
		if c.expectedString == "" {
			expectedString = c.in
		}

		text, _ := d.MarshalText()
		assertEqual(t, expectedString, string(text))
	}

	for _, c := range cases {
		var d Duration

		err := yaml.Unmarshal([]byte(c.in), &d)
		if err != nil {
			t.Errorf("Unexpected error on input %q", c.in)
		}

		assertEqual(t, c.out, time.Duration(d))

		expectedString := c.expectedString
		if c.expectedString == "" {
			expectedString = c.in
		}

		text, err := yaml.Marshal(d)
		assertNilError(t, err)

		assertEqual(t, expectedString, strings.TrimSpace(string(text)))
	}
}

func TestDuration_UnmarshalJSON(t *testing.T) {
	cases := []struct {
		in  string
		out time.Duration

		expectedString string
	}{
		{
			in:             `"0"`,
			out:            0,
			expectedString: `"0s"`,
		},
		{
			in:             `"0w"`,
			out:            0,
			expectedString: `"0s"`,
		},
		{
			in:  `"0s"`,
			out: 0,
		},
		{
			in:  `"324ms"`,
			out: 324 * time.Millisecond,
		},
		{
			in:  `"3s"`,
			out: 3 * time.Second,
		},
		{
			in:  `"5m"`,
			out: 5 * time.Minute,
		},
		{
			in:  `"1h"`,
			out: time.Hour,
		},
		{
			in:  `"4d"`,
			out: 4 * 24 * time.Hour,
		},
		{
			in:  `"4d1h"`,
			out: 4*24*time.Hour + time.Hour,
		},
		{
			in:             `"14d"`,
			out:            14 * 24 * time.Hour,
			expectedString: `"2w"`,
		},
		{
			in:  `"3w"`,
			out: 3 * 7 * 24 * time.Hour,
		},
		{
			in:             `"3w2d1h"`,
			out:            3*7*24*time.Hour + 2*24*time.Hour + time.Hour,
			expectedString: `"23d1h"`,
		},
		{
			in:  `"10y"`,
			out: 10 * 365 * 24 * time.Hour,
		},
		{
			in:  `"289y"`,
			out: 289 * 365 * 24 * time.Hour,
		},
	}

	for _, c := range cases {
		var d Duration
		err := json.Unmarshal([]byte(c.in), &d)
		if err != nil {
			t.Errorf("Unexpected error on input %q", c.in)
		}
		if time.Duration(d) != c.out {
			t.Errorf("Expected %v but got %v", c.out, d)
		}

		expectedString := c.expectedString
		if c.expectedString == "" {
			expectedString = c.in
		}

		bytes, err := json.Marshal(d)
		if err != nil {
			t.Errorf("Unexpected error on marshal of %v: %s", d, err)
		}

		if string(bytes) != expectedString {
			t.Errorf("Expected duration string %q but got %q", c.in, d.String())
		}
	}
}

func TestParseBadDuration(t *testing.T) {
	cases := []string{
		"1",
		"1y1m1d",
		"1.5d",
		"d",
		"294y",
		"200y10400w",
		"107675d",
		"2584200h",
		"",
	}

	for _, c := range cases {
		_, err := ParseDuration(c)
		if err == nil {
			t.Errorf("Expected error on input %s", c)
		}
	}
}

func assertEqual[T comparable](t *testing.T, expected, value T) {
	t.Helper()

	if expected != value {
		t.Errorf("Expected: %v, got: %v", expected, value)
	}
}

func assertNilError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("Expected nil error, got: %v", err)
	}
}

var durationTruncateTests = []struct {
	d    Duration
	m    Duration
	want Duration
}{
	{0, Duration(time.Second), 0},
	{Duration(time.Minute), Duration(-7 * time.Second), Duration(time.Minute)},
}

func TestDurationTruncate(t *testing.T) {
	for _, tt := range durationTruncateTests {
		if got := tt.d.Truncate(tt.m); got != tt.want {
			t.Errorf("Duration(%s).Truncate(%s) = %s; want: %s", tt.d, tt.m, got, tt.want)
		}
	}
}

var durationRoundTests = []struct {
	d    Duration
	m    Duration
	want Duration
}{
	{0, Duration(time.Second), 0},
	{Duration(time.Minute), Duration(-11 * time.Second), Duration(time.Minute)},
	{8e18, 3e18, 9e18},
	{9e18, 5e18, 1<<63 - 1},
	{-8e18, 3e18, -9e18},
	{-9e18, 5e18, -1 << 63},
	{3<<61 - 1, 3 << 61, 3 << 61},
}

func TestDurationRound(t *testing.T) {
	for _, tt := range durationRoundTests {
		if got := tt.d.Round(tt.m); got != tt.want {
			t.Errorf("Duration(%s).Round(%s) = %s; want: %s", tt.d, tt.m, got, tt.want)
		}
	}
}

var durationAbsTests = []struct {
	d    Duration
	want Duration
}{
	{0, 0},
	{1, 1},
	{-1, 1},
	{minDuration, maxDuration},
	{minDuration + 1, maxDuration},
	{minDuration + 2, maxDuration - 1},
	{maxDuration, maxDuration},
	{maxDuration - 1, maxDuration - 1},
}

func TestDurationAbs(t *testing.T) {
	for _, tt := range durationAbsTests {
		if got := tt.d.Abs(); got != tt.want {
			t.Errorf("Duration(%s).Abs() = %s; want: %s", tt.d, got, tt.want)
		}
	}
}
