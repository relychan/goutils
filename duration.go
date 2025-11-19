package goutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"go.yaml.in/yaml/v4"
)

var (
	// ErrDurationOutOfRange occurs when the duration int64 value is out of range.
	ErrDurationOutOfRange = errors.New("duration out of range")
	// ErrDurationEmpty occurs when the duration string is empty.
	ErrDurationEmpty = errors.New("empty duration string")
	// ErrInvalidDurationString occurs when the duration string is invalid.
	ErrInvalidDurationString = errors.New("not a valid duration string")
	// ErrUnknownDurationUnit occurs when the duration unit is unknown.
	ErrUnknownDurationUnit = errors.New("unknown duration unit")
	// ErrInvalidDateTimeString occurs when the date time string is invalid.
	ErrInvalidDateTimeString = errors.New("not a valid date time string")
)

// Duration wraps time.Duration. It is used to parse and format custom duration strings
// from YAML, JSON, and text formats.
type Duration time.Duration //nolint:recvcheck

// Set implements pflag/flag.Value.
func (d *Duration) Set(s string) error {
	var err error

	*d, err = ParseDuration(s)

	return err
}

// Type implements pflag.Value.
func (*Duration) Type() string {
	return "duration"
}

// Units are required to go in order from biggest to smallest.
// This guards against confusion from "1m1d" being 1 minute + 1 day, not 1 month + 1 day.
var unitMap = map[string]struct {
	pos  int
	mult uint64
}{
	"ms": {7, uint64(time.Millisecond)},
	"s":  {6, uint64(time.Second)},
	"m":  {5, uint64(time.Minute)},
	"h":  {4, uint64(time.Hour)},
	"d":  {3, uint64(24 * time.Hour)},
	"w":  {2, uint64(7 * 24 * time.Hour)},
	"y":  {1, uint64(365 * 24 * time.Hour)},
}

// ParseDuration parses a string into a time.Duration, assuming that a year
// always has 365d, a week always has 7d, and a day always has 24h.
// Negative durations are not supported.
func ParseDuration(s string) (Duration, error) {
	switch s {
	case "0":
		// Allow 0 without a unit.
		return 0, nil
	case "":
		return 0, ErrDurationEmpty
	}

	orig := s
	lastUnitPos := 0

	var dur uint64

	for s != "" {
		if !IsDigit(s[0]) {
			return 0, fmt.Errorf("%w: %q", ErrInvalidDurationString, orig)
		}

		// Consume [0-9]*
		i := 0

		for ; i < len(s) && IsDigit(s[i]); i++ { //nolint:revive
		}

		v, err := strconv.ParseUint(s[:i], 10, 0)
		if err != nil {
			return 0, fmt.Errorf("%w: %q", ErrInvalidDurationString, orig)
		}

		s = s[i:]

		// Consume unit.
		for i = 0; i < len(s) && !IsDigit(s[i]); i++ { //nolint:revive
		}

		if i == 0 {
			return 0, fmt.Errorf("%w: %q", ErrInvalidDurationString, orig)
		}

		u := s[:i]
		s = s[i:]

		unit, ok := unitMap[u]
		if !ok {
			return 0, fmt.Errorf("%w %q in %q", ErrUnknownDurationUnit, u, orig)
		}

		if unit.pos <= lastUnitPos { // Units must go in order from biggest to smallest.
			return 0, fmt.Errorf("%w: %q", ErrInvalidDurationString, orig)
		}

		lastUnitPos = unit.pos
		// Check if the provided duration overflows time.Duration (> ~ 290years).
		if v > ((1 << 63) / unit.mult) {
			return 0, ErrDurationOutOfRange
		}

		dur += v * unit.mult
		if dur > 1<<63-1 {
			return 0, ErrDurationOutOfRange
		}
	}

	return Duration(dur), nil //nolint:gosec
}

// Abs returns the absolute value of the current duration.
func (d Duration) Abs() Duration {
	r := time.Duration(d).Abs()

	return Duration(r)
}

// Hours returns the duration as a floating point number of hours.
func (d Duration) Hours() float64 {
	return time.Duration(d).Hours()
}

// Microseconds returns the duration as an integer microsecond count.
func (d Duration) Microseconds() int64 {
	return time.Duration(d).Microseconds()
}

// Milliseconds returns the duration as an integer millisecond count.
func (d Duration) Milliseconds() int64 {
	return time.Duration(d).Milliseconds()
}

// Minutes returns the duration as a floating point number of minutes.
func (d Duration) Minutes() float64 {
	return time.Duration(d).Minutes()
}

// Nanoseconds returns the duration as an integer nanosecond count.
func (d Duration) Nanoseconds() int64 {
	return time.Duration(d).Nanoseconds()
}

// Round returns the result of rounding d to the nearest multiple of m.
// The rounding behavior for halfway values is to round away from zero.
// If the result exceeds the maximum (or minimum) value that can be stored in a [Duration],
// Round returns the maximum (or minimum) duration.
// If m <= 0, Round returns d unchanged.
func (d Duration) Round(m Duration) Duration {
	r := time.Duration(d).Round(time.Duration(m))

	return Duration(r)
}

// Seconds returns the duration as a floating point number of seconds.
func (d Duration) Seconds() float64 {
	return time.Duration(d).Seconds()
}

// Truncate returns the result of rounding d toward zero to a multiple of m.
// If m <= 0, Truncate returns d unchanged.
func (d Duration) Truncate(m Duration) Duration {
	r := time.Duration(d).Truncate(time.Duration(m))

	return Duration(r)
}

// String implements the fmt.Stringer interface.
func (d Duration) String() string {
	var (
		ms   = int64(time.Duration(d) / time.Millisecond)
		r    = ""
		sign = ""
	)

	if ms == 0 {
		return "0s"
	}

	if ms < 0 {
		sign, ms = "-", -ms
	}

	f := func(unit string, mult int64, exact bool) {
		if exact && ms%mult != 0 {
			return
		}

		if v := ms / mult; v > 0 {
			r += fmt.Sprintf("%d%s", v, unit)
			ms -= v * mult
		}
	}

	// Only format years and weeks if the remainder is zero, as it is often
	// easier to read 90d than 12w6d.
	f("y", 1000*60*60*24*365, true)
	f("w", 1000*60*60*24*7, true)

	f("d", 1000*60*60*24, false)
	f("h", 1000*60*60, false)
	f("m", 1000*60, false)
	f("s", 1000, false)
	f("ms", 1, false)

	return sign + r
}

// MarshalJSON implements the json.Marshaler interface.
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *Duration) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	// Properly unescape a JSON string.
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return fmt.Errorf("Duration.UnmarshalJSON: %w", errInvalidJSONString)
	}

	data = data[len(`"`) : len(data)-len(`"`)]

	dur, err := ParseDuration(string(data))
	if err != nil {
		return err
	}

	*d = dur

	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
func (d Duration) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (d *Duration) UnmarshalText(text []byte) error {
	var err error

	*d, err = ParseDuration(string(text))

	return err
}

// MarshalYAML implements the yaml.Marshaler interface.
func (d Duration) MarshalYAML() (any, error) {
	return d.String(), nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (d *Duration) UnmarshalYAML(value *yaml.Node) error {
	var s string

	err := value.Decode(&s)
	if err != nil {
		return err
	}

	dur, err := ParseDuration(s)
	if err != nil {
		return err
	}

	*d = dur

	return nil
}
