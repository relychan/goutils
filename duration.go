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
	ErrUnknownDurationUnit = errors.New("unknown unit in duration")
)

// Duration wraps time.Duration. It is used to parse the custom duration format
// from YAML.
// This type should not propagate beyond the scope of input/output processing.
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
			return 0, fmt.Errorf("%q is %w %q", u, ErrUnknownDurationUnit, orig)
		}

		if unit.pos <= lastUnitPos { // Units must go in order from biggest to smallest.
			return 0, fmt.Errorf("%w: %q", ErrInvalidDurationString, orig)
		}

		lastUnitPos = unit.pos
		// Check if the provided duration overflows time.Duration (> ~ 290years).
		if v > 1<<63/unit.mult {
			return 0, ErrDurationOutOfRange
		}

		dur += v * unit.mult
		if dur > 1<<63-1 {
			return 0, ErrDurationOutOfRange
		}
	}

	return Duration(dur), nil //nolint:gosec
}

var _ fmt.Stringer = (*Duration)(nil)

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
func (d *Duration) UnmarshalJSON(bytes []byte) error {
	var s string

	err := json.Unmarshal(bytes, &s)
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

// MarshalText implements the encoding.TextMarshaler interface.
func (d *Duration) MarshalText() ([]byte, error) {
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
