package goutils

import (
	"encoding/json"
	"regexp"
	"strings"

	"go.yaml.in/yaml/v4"
)

// RegexpMatcher wraps the regexp.Regexp with a string for simple matching.
type RegexpMatcher struct { //nolint:recvcheck
	text   *string
	regexp *regexp.Regexp
}

// NewRegexpMatcher creates a [RegexpMatcher] from a raw string.
func NewRegexpMatcher(input string) (*RegexpMatcher, error) {
	result := &RegexpMatcher{}

	return result, result.UnmarshalText([]byte(input))
}

// IsZero returns true if the current instance is in its zero state.
func (j RegexpMatcher) IsZero() bool {
	return j.text == nil && j.regexp == nil
}

// Equal checks if the target value is equal.
func (j RegexpMatcher) Equal(target RegexpMatcher) bool {
	return EqualComparablePtr(j.text, target.text) &&
		(j.regexp == target.regexp || (j.regexp != nil && target.regexp != nil &&
			j.regexp.String() == target.regexp.String()))
}

// String implements the fmt.Stringer interface.
func (j RegexpMatcher) String() string {
	if j.regexp != nil {
		return j.regexp.String()
	}

	if j.text != nil {
		return *j.text
	}

	return ""
}

// UnmarshalText implements [encoding.TextUnmarshaler] by calling
// [Compile] on the encoded value.
func (j *RegexpMatcher) UnmarshalText(bs []byte) error {
	text := string(bs)

	if !Every(bs, IsMetaCharacter) {
		j.text = &text
		j.regexp = nil
	} else {
		re, err := regexp.Compile(text)
		if err != nil {
			return err
		}

		j.regexp = re
		j.text = nil
	}

	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *RegexpMatcher) UnmarshalJSON(data []byte) error {
	var raw string

	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	return j.UnmarshalText([]byte(raw))
}

// MarshalJSON implements json.Marshaler.
func (j RegexpMatcher) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.String())
}

// UnmarshalYAML implements custom deserialization for the yaml.Unmarshaler interface.
func (j *RegexpMatcher) UnmarshalYAML(value *yaml.Node) error {
	var raw string

	err := value.Decode(&raw)
	if err != nil {
		return err
	}

	return j.UnmarshalText([]byte(raw))
}

// MarshalYAML implements the custom behavior for the yaml.Marshaler interface.
func (j RegexpMatcher) MarshalYAML() (any, error) {
	return j.String(), nil
}

// Match reports whether the byte slice b contains any match of the regular expression re.
func (j *RegexpMatcher) Match(b []byte) bool {
	if j.text != nil {
		return strings.Contains(string(b), *j.text)
	}

	if j.regexp != nil {
		return j.regexp.Match(b)
	}

	return false
}

// MatchString reports whether the string s contains any match of the regular expression re.
func (j *RegexpMatcher) MatchString(s string) bool {
	if j.text != nil {
		return strings.Contains(s, *j.text)
	}

	if j.regexp != nil {
		return j.regexp.MatchString(s)
	}

	return false
}
