package goutils

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"go.yaml.in/yaml/v4"
)

var (
	dateLength     = len("2006-01-02")
	dateTimeLength = len("2006-01-02T15:04:05")
)

var (
	errBadValue          = errors.New("bad value for field") // placeholder not passed to user
	errInvalidJSONString = errors.New("input is not a JSON string")
)

// Time wraps time.Time. It is used to parse and format custom date/time strings
// from YAML, JSON, and text formats.
type Time time.Time //nolint:recvcheck

// IsZero reports whether t represents the zero time instant,
// January 1, year 1, 00:00:00 UTC.
func (t Time) IsZero() bool {
	return time.Time(t).IsZero()
}

// After reports whether the time instant t is after u.
func (t Time) After(u Time) bool {
	return time.Time(t).After(time.Time(u))
}

// Before reports whether the time instant t is before u.
func (t Time) Before(u Time) bool {
	return time.Time(t).Before(time.Time(u))
}

// Compare compares the time instant t with u. If t is before u, it returns -1;
// if t is after u, it returns +1; if they're the same, it returns 0.
func (t Time) Compare(u Time) int {
	return time.Time(t).Compare(time.Time(u))
}

// Equal reports whether t and u represent the same time instant.
// Two times can be equal even if they are in different locations.
// For example, 6:00 +0200 and 4:00 UTC are Equal.
// See the documentation on the Time type for the pitfalls of using == with
// Time values; most code should use Equal instead.
func (t Time) Equal(u Time) bool {
	return time.Time(t).Equal(time.Time(u))
}

// Add returns the time t+d.
func (t Time) Add(d time.Duration) Time {
	r := time.Time(t).Add(d)

	return Time(r)
}

// AddDate returns the time corresponding to adding the given number of years, months, and days to t.
// For example, AddDate(-1, 2, 3) applied to January 1, 2011 returns March 4, 2010.
//
// Note that dates are fundamentally coupled to timezones, and calendrical periods like days don't have fixed durations.
// AddDate uses the Location of the Time value to determine these durations.
// That means that the same AddDate arguments can produce a different shift in absolute time depending on the base Time value and its Location.
// For example, AddDate(0, 0, 1) applied to 12:00 on March 27 always returns 12:00 on March 28. At some locations and in some years this is a 24 hour shift.
// In others it's a 23 hour shift due to daylight savings time transitions.
//
// AddDate normalizes its result in the same way that Date does, so, for example,
// adding one month to October 31 yields December 1, the normalized form for November 31.
func (t Time) AddDate(years int, months int, days int) Time {
	r := time.Time(t).AddDate(years, months, days)

	return Time(r)
}

// AppendBinary implements the encoding.BinaryAppender interface.
func (t Time) AppendBinary(b []byte) ([]byte, error) {
	return time.Time(t).AppendBinary(b)
}

// AppendFormat is like [Time.Format] but appends the textual representation to b and returns the extended buffer.
func (t Time) AppendFormat(b []byte, layout string) []byte {
	return time.Time(t).AppendFormat(b, layout)
}

// AppendText implements the encoding.TextAppender interface.
// The time is formatted in RFC 3339 format with sub-second precision.
// If the timestamp cannot be represented as valid RFC 3339 (e.g., the year is out of range), then an error is returned.
func (t Time) AppendText(b []byte) ([]byte, error) {
	return time.Time(t).AppendText(b)
}

// Clock returns the hour, minute, and second within the day specified by t.
func (t Time) Clock() (int, int, int) {
	return time.Time(t).Clock()
}

// Date returns the year, month, and day in which t occurs.
func (t Time) Date() (int, time.Month, int) {
	return time.Time(t).Date()
}

// Day returns the day of the month specified by t.
func (t Time) Day() int {
	return time.Time(t).Day()
}

// Format returns a textual representation of the time value formatted according to the layout defined by the argument.
// See the documentation for the constant called [Layout] to see how to represent the layout format.
// The executable example for [Time.Format] demonstrates the working of the layout string in detail and is a good reference.
func (t Time) Format(layout string) string {
	return time.Time(t).Format(layout)
}

// GoString implements fmt.GoStringer and formats t to be printed in Go source code.
func (t Time) GoString() string {
	return time.Time(t).GoString()
}

// GobDecode implements the gob.GobDecoder interface.
func (t *Time) GobDecode(data []byte) error {
	var raw time.Time

	err := raw.GobDecode(data)
	if err != nil {
		return err
	}

	*t = Time(raw)

	return nil
}

// GobEncode implements the gob.GobEncoder interface.
func (t Time) GobEncode() ([]byte, error) {
	return time.Time(t).GobEncode()
}

// Hour returns the hour within the day specified by t, in the range [0, 23].
func (t Time) Hour() int {
	return time.Time(t).Hour()
}

// ISOWeek returns the ISO 8601 year and week number in which t occurs.
// Week ranges from 1 to 53. Jan 01 to Jan 03 of year n might belong to week 52 or 53 of year n-1,
// and Dec 29 to Dec 31 might belong to week 1 of year n+1.
func (t Time) ISOWeek() (int, int) {
	return time.Time(t).ISOWeek()
}

// In returns a copy of t representing the same time instant,
// but with the copy's location information set to loc for display purposes.
func (t Time) In(loc *time.Location) Time {
	return Time(time.Time(t).In(loc))
}

// IsDST reports whether the time in the configured location is in Daylight Savings Time.
func (t Time) IsDST() bool {
	return time.Time(t).IsDST()
}

// Local returns t with the location set to local time.
func (t Time) Local() Time {
	return Time(time.Time(t).Local()) //nolint:gosmopolitan
}

// Location returns the time zone information associated with t.
func (t Time) Location() *time.Location {
	return time.Time(t).Location()
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (t Time) MarshalBinary() ([]byte, error) {
	return time.Time(t).MarshalBinary()
}

// MarshalJSON implements the encoding/json.Marshaler interface.
// The time is a quoted string in the RFC 3339 format with sub-second precision.
// If the timestamp cannot be represented as valid RFC 3339 (e.g., the year is out of range), then an error is reported.
func (t Time) MarshalJSON() ([]byte, error) {
	return time.Time(t).MarshalJSON()
}

// MarshalText implements the encoding.TextMarshaler interface.
// The output matches that of calling the [Time.AppendText] method.
func (t Time) MarshalText() ([]byte, error) {
	return time.Time(t).MarshalText()
}

// Minute returns the minute offset within the hour specified by t, in the range [0, 59].
func (t Time) Minute() int {
	return time.Time(t).Minute()
}

// Month returns the month of the year specified by t.
func (t Time) Month() time.Month {
	return time.Time(t).Month()
}

// Nanosecond returns the nanosecond offset within the second specified by t, in the range [0, 999999999].
func (t Time) Nanosecond() int {
	return time.Time(t).Nanosecond()
}

// Round returns the result of rounding t to the nearest multiple of d (since the zero time).
// The rounding behavior for halfway values is to round up.
// If d <= 0, Round returns t stripped of any monotonic clock reading but otherwise unchanged.
// Round operates on the time as an absolute duration since the zero time; it does not operate on the presentation form of the time.
// Thus, Round(Hour) may return a time with a non-zero minute, depending on the time's Location.
func (t Time) Round(d time.Duration) Time {
	return Time(time.Time(t).Round(d))
}

// Second returns the second offset within the minute specified by t, in the range [0, 59].
func (t Time) Second() int {
	return time.Time(t).Second()
}

// String returns the time formatted using the format string
//
// "2006-01-02 15:04:05.999999999 -0700 MST"
// If the time has a monotonic clock reading, the returned string includes a final field "m=±<value>", where value is the monotonic clock reading formatted as a decimal number of seconds.
//
// The returned string is meant for debugging; for a stable serialized representation, use t.MarshalText, t.MarshalBinary, or t.Format with an explicit format string.
func (t Time) String() string {
	return time.Time(t).String()
}

// Sub returns the duration t-u.
// If the result exceeds the maximum (or minimum) value that can be stored in a [time.Duration],
// the maximum (or minimum) duration will be returned.
// To compute t-d for a duration d, use t.Add(-d).
func (t Time) Sub(u Time) time.Duration {
	return time.Time(t).Sub(time.Time(u))
}

// Truncate returns the result of rounding t down to a multiple of d (since the zero time).
// If d <= 0, Truncate returns t stripped of any monotonic clock reading but otherwise unchanged.
// Truncate operates on the time as an absolute duration since the zero time; it does not operate on the presentation form of the time.
// Thus, Truncate(Hour) may return a time with a non-zero minute, depending on the time's Location.
func (t Time) Truncate(d time.Duration) Time {
	return Time(time.Time(t).Truncate(d))
}

// UTC returns t with the location set to UTC.
func (t Time) UTC() Time {
	return Time(time.Time(t).UTC())
}

// Unix returns t as a Unix time, the number of seconds elapsed since January 1, 1970 UTC.
// The result does not depend on the location associated with t.
// Unix-like operating systems often record time as a 32-bit count of seconds,
// but since the method here returns a 64-bit value it is valid for billions of years into the past or future.
func (t Time) Unix() int64 {
	return time.Time(t).Unix()
}

// UnixMicro returns t as a Unix time, the number of microseconds elapsed since January 1, 1970 UTC.
// The result is undefined if the Unix time in microseconds cannot be represented by an int64 (a date before year -290307 or after year 294246).
// The result does not depend on the location associated with t.
func (t Time) UnixMicro() int64 {
	return time.Time(t).UnixMicro()
}

// UnixMilli returns t as a Unix time, the number of milliseconds elapsed since January 1, 1970 UTC.
// The result is undefined if the Unix time in milliseconds cannot be represented by an int64 (a date more than 292 million years before or after 1970).
// The result does not depend on the location associated with t.
func (t Time) UnixMilli() int64 {
	return time.Time(t).UnixMilli()
}

// UnixNano returns t as a Unix time, the number of nanoseconds elapsed since January 1, 1970 UTC.
// The result is undefined if the Unix time in nanoseconds cannot be represented by an int64 (a date before the year 1678 or after 2262).
// Note that this means the result of calling UnixNano on the zero Time is undefined.
// The result does not depend on the location associated with t.
func (t Time) UnixNano() int64 {
	return time.Time(t).UnixNano()
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (t *Time) UnmarshalBinary(data []byte) error {
	var raw time.Time

	err := raw.UnmarshalBinary(data)
	if err != nil {
		return err
	}

	*t = Time(raw)

	return nil
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
// The time must be a quoted string in the RFC 3339 or ISO 8601 format.
func (t *Time) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	// Properly unescape a JSON string.
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return fmt.Errorf("Time.UnmarshalJSON: %w", errInvalidJSONString)
	}

	data = data[len(`"`) : len(data)-len(`"`)]

	return t.UnmarshalText(data)
}

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
// The time must be in the RFC 3339 or ISO 8601 format.
func (t *Time) UnmarshalText(data []byte) error {
	var err error

	*t, err = ParseDateTime(data)

	return err
}

// Weekday returns the day of the week specified by t.
func (t Time) Weekday() time.Weekday {
	return time.Time(t).Weekday()
}

// Year returns the year in which t occurs.
func (t Time) Year() int {
	return time.Time(t).Year()
}

// YearDay returns the day of the year specified by t,
// in the range [1,365] for non-leap years, and [1,366] in leap years.
func (t Time) YearDay() int {
	return time.Time(t).YearDay()
}

// Zone computes the time zone in effect at time t,
// returning the abbreviated name of the zone (such as "CET") and its offset in seconds east of UTC.
func (t Time) Zone() (string, int) {
	return time.Time(t).Zone()
}

// ZoneBounds returns the bounds of the time zone in effect at time t.
// The zone begins at start and the next zone begins at end.
// If the zone begins at the beginning of time, start will be returned as a zero Time.
// If the zone goes on forever, end will be returned as a zero Time.
// The Location of the returned times will be the same as t.
func (t Time) ZoneBounds() (Time, Time) {
	start, end := time.Time(t).ZoneBounds()

	return Time(start), Time(end)
}

// MarshalYAML implements the yaml.Marshaler interface.
func (t Time) MarshalYAML() (any, error) {
	return time.Time(t), nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (t *Time) UnmarshalYAML(value *yaml.Node) error {
	var s string

	err := value.Decode(&s)
	if err != nil {
		return err
	}

	*t, err = ParseDateTime(s)

	return err
}

// Now returns the current local time.
func Now() Time {
	return Time(time.Now())
}

// Date returns the Time corresponding to
// yyyy-mm-dd hh:mm:ss + nsec nanoseconds
// in the appropriate zone for that time in the given location.
// The month, day, hour, min, sec, and nsec values may be outside their usual ranges and will be normalized during the conversion. For example, October 32 converts to November 1.
// A daylight savings time transition skips or repeats times. For example, in the United States, March 13, 2011 2:15am never occurred, while November 6, 2011 1:15am occurred twice.
// In such cases, the choice of time zone, and therefore the time, is not well-defined. Date returns a time that is correct in one of the two zones involved in the transition, but it does not guarantee which.
// Date panics if loc is nil.
func Date(
	year int,
	month time.Month,
	day int,
	hour int,
	minutes int,
	sec int,
	nsec int,
	loc *time.Location,
) Time {
	return Time(time.Date(year, month, day, hour, minutes, sec, nsec, loc))
}

// Unix returns the local Time corresponding to the given Unix time, sec seconds and nsec nanoseconds since January 1, 1970 UTC.
// It is valid to pass nsec outside the range [0, 999999999]. Not all sec values have a corresponding time value.
// One such value is 1<<63-1 (the largest int64 value).
func Unix(sec int64, nsec int64) Time {
	return Time(time.Unix(sec, nsec))
}

// ParseDateTime parses date time in RFC 3339 or ISO 8601 format formats.
func ParseDateTime[B []byte | string](input B) (Time, error) {
	result, err := ParseDateTimeNative(input)

	return Time(result), err
}

// ParseDateTimeNative parses date time in RFC 3339 or ISO 8601 format formats.
func ParseDateTimeNative[B []byte | string](input B) (time.Time, error) {
	result, ok := parseDateTimeString(input)
	if ok {
		return result, nil
	}

	return time.Time{}, fmt.Errorf("%w: %s", ErrInvalidDateTimeString, input)
}

func parseDateTimeString[B []byte | string](s B) (time.Time, bool) { //nolint:cyclop,funlen
	// Parse the date and time.
	if (len(s) < dateLength) || s[4] != '-' || s[7] != '-' {
		return time.Time{}, false
	}

	year, ok := ParseIntInRange(s[0:4], 0, 9999) // e.g., 2006
	if !ok {
		return time.Time{}, false
	}

	month, ok := ParseIntInRange(s[5:7], 1, 12) // e.g., 01
	if !ok {
		return time.Time{}, false
	}

	day, ok := ParseIntInRange(s[8:10], 1, daysIn(time.Month(month), year)) // e.g., 02
	if !ok {
		return time.Time{}, false
	}

	if len(s) == dateLength {
		t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

		return t, true
	}

	if (len(s) < dateTimeLength) || (s[10] != 'T' && s[10] != ' ') || s[13] != ':' ||
		s[16] != ':' {
		return time.Time{}, false
	}

	hour, ok := ParseIntInRange(s[11:13], 0, 23) // e.g., 15
	if !ok {
		return time.Time{}, false
	}

	minute, ok := ParseIntInRange(s[14:16], 0, 59) // e.g., 04
	if !ok {
		return time.Time{}, false
	}

	sec, ok := ParseIntInRange(s[17:19], 0, 59) // e.g., 05
	if !ok {
		return time.Time{}, false
	}

	s = s[19:]

	// Parse the fractional second.
	var nsec int

	if len(s) >= 2 && s[0] == '.' && IsDigit(s[1]) {
		n := 2

		for ; n < len(s) && IsDigit(s[n]); n++ { //nolint:revive
		}

		nsec, _, _ = parseNanoseconds(s, n) //nolint:errcheck

		s = s[n:]
	}

	// Parse the time zone.
	t := time.Date(year, time.Month(month), day, hour, minute, sec, nsec, time.UTC)

	if len(s) == 0 || (len(s) == 1 && s[0] == 'Z') {
		return t, true
	}

	hr, mm, ok := parseTimeZoneOffset(s)
	if !ok {
		return time.Time{}, false
	}

	zoneOffset := (hr*60 + mm) * 60
	if s[0] == '-' {
		zoneOffset *= -1
	}

	t = t.Add(time.Duration(-zoneOffset) * time.Second)

	// Use local zone with the given offset if possible.
	tz := t.Local() //nolint:gosmopolitan

	_, offset := tz.Zone()
	if offset == zoneOffset {
		return tz, true
	}

	t = t.In(time.FixedZone("", zoneOffset))

	return t, true
}

func parseTimeZoneOffset[B []byte | string](s B) (int, int, bool) {
	var rawHour, rawMinute B

	if s[0] != 'Z' && s[0] != '+' && s[0] != '-' {
		return 0, 0, false
	}

	switch {
	case len(s) == 5:
		rawHour = s[1:3]
		rawMinute = s[3:5]
	case len(s) == 6 && s[3] == ':':
		rawHour = s[1:3]
		rawMinute = s[4:6]
	default:
		return 0, 0, false
	}

	hr, ok := ParseIntInRange(rawHour, 0, 23) // e.g., 07
	if !ok {
		return 0, 0, false
	}

	mm, ok := ParseIntInRange(rawMinute, 0, 59) // e.g., 00
	if !ok {
		return 0, 0, false
	}

	return hr, mm, true
}

func isLeap(year int) bool {
	// year%4 == 0 && (year%100 != 0 || year%400 == 0)
	// Bottom 2 bits must be clear.
	// For multiples of 25, bottom 4 bits must be clear.
	// Thanks to Cassio Neri for this trick.
	mask := 0xf
	if year%25 != 0 {
		mask = 3
	}

	return year&mask == 0
}

func daysIn(m time.Month, year int) int {
	if m == time.February {
		if isLeap(year) {
			return 29
		}

		return 28
	}
	// With the special case of February eliminated, the pattern is
	//	31 30 31 30 31 30 31 31 30 31 30 31
	// Adding m&1 produces the basic alternation;
	// adding (m>>3)&1 inverts the alternation starting in August.
	return 30 + int((m+m>>3)&1)
}

func commaOrPeriod(b byte) bool {
	return b == '.' || b == ','
}

func parseNanoseconds[bytes []byte | string](
	value bytes,
	nbytes int,
) (int, string, error) {
	if !commaOrPeriod(value[0]) {
		return 0, "", errBadValue
	}

	if nbytes > 10 {
		value = value[:10]
		nbytes = 10
	}

	ns, err := strconv.Atoi(string(value[1:nbytes]))
	if err != nil {
		return ns, "", err
	}

	if ns < 0 {
		return ns, "fractional second", nil
	}

	// We need nanoseconds, which means scaling by the number
	// of missing digits in the format, maximum length 10.
	scaleDigits := 10 - nbytes

	for range scaleDigits {
		ns *= 10
	}

	return ns, "", err
}
