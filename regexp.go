package goutils

import (
	"encoding/json"
	"regexp"
	"strings"

	"go.yaml.in/yaml/v4"
)

type matchOp int8

const (
	equalOp matchOp = iota
	prefixOp
	suffixOp
	containOp
	regexpOp
)

// RegexpMatcher wraps the regexp.Regexp with a string for simple matching.
type RegexpMatcher struct { //nolint:recvcheck
	text   *string
	regexp *regexp.Regexp
	op     matchOp
}

// NewRegexpMatcher creates a [RegexpMatcher] from a raw string.
func NewRegexpMatcher(input string) (*RegexpMatcher, error) {
	result := &RegexpMatcher{}

	return result, result.UnmarshalText([]byte(input))
}

// MustRegexpMatcher creates a [RegexpMatcher] from a raw string. Panic if error.
func MustRegexpMatcher(input string) *RegexpMatcher {
	result := &RegexpMatcher{}

	err := result.UnmarshalText([]byte(input))
	if err != nil {
		panic(err)
	}

	return result
}

// IsZero returns true if the current instance is in its zero state.
func (j RegexpMatcher) IsZero() bool {
	return j.text == nil && j.regexp == nil
}

// Equal checks if the target value is equal.
func (j RegexpMatcher) Equal(target RegexpMatcher) bool {
	return EqualComparablePtr(j.text, target.text) &&
		j.op == target.op &&
		(j.regexp == target.regexp || (j.regexp != nil && target.regexp != nil &&
			j.regexp.String() == target.regexp.String()))
}

// String implements the fmt.Stringer interface.
func (j RegexpMatcher) String() string {
	if j.regexp != nil {
		return j.regexp.String()
	}

	if j.text == nil {
		return ""
	}

	switch j.op {
	case equalOp:
		return "^" + *j.text + "$"
	case prefixOp:
		return "^" + *j.text
	case suffixOp:
		return *j.text + "$"
	default:
		return *j.text
	}
}

// UnmarshalText implements [encoding.TextUnmarshaler] by calling
// [Compile] on the encoded value.
func (j *RegexpMatcher) UnmarshalText(bs []byte) error {
	text := string(bs)
	op := containOp

	switch {
	case len(bs) == 0:
	case len(bs) == 1:
		if isRegexSyntaxChar(rune(bs[0])) {
			op = regexpOp
		}
	case bs[0] == '^' && bs[len(bs)-1] == '$':
		subtext := text[1 : len(bs)-1]
		if isRegexSyntax(subtext) {
			op = regexpOp
		} else {
			op = equalOp
			text = subtext
		}
	case bs[0] == '^':
		subtext := text[1:]
		if isRegexSyntax(subtext) {
			op = regexpOp
		} else {
			op = prefixOp
			text = subtext
		}
	case bs[len(bs)-1] == '$':
		subtext := text[:len(text)-1]
		if isRegexSyntax(subtext) {
			op = regexpOp
		} else {
			op = suffixOp
			text = subtext
		}
	default:
		if isRegexSyntax(text) {
			op = regexpOp
		}
	}

	if op == regexpOp {
		re, err := regexp.Compile(text)
		if err != nil {
			return err
		}

		j.text = nil
		j.regexp = re
		j.op = op

		return nil
	}

	j.op = op
	j.text = &text
	j.regexp = nil

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
	switch j.op {
	case equalOp:
		return j.text != nil && *j.text == string(b)
	case prefixOp:
		return j.text != nil && strings.HasPrefix(string(b), *j.text)
	case suffixOp:
		return j.text != nil && strings.HasSuffix(string(b), *j.text)
	case containOp:
		return j.text != nil && strings.Contains(string(b), *j.text)
	default:
		return j.regexp != nil && j.regexp.Match(b)
	}
}

// MatchString reports whether the string s contains any match of the regular expression re.
func (j *RegexpMatcher) MatchString(s string) bool {
	switch j.op {
	case equalOp:
		return j.text != nil && *j.text == s
	case prefixOp:
		return j.text != nil && strings.HasPrefix(s, *j.text)
	case suffixOp:
		return j.text != nil && strings.HasSuffix(s, *j.text)
	case containOp:
		return j.text != nil && strings.Contains(s, *j.text)
	default:
		return j.regexp != nil && j.regexp.MatchString(s)
	}
}

func isRegexSyntaxChar(r rune) bool {
	return !IsMetaCharacter(r) && !strings.ContainsRune("`~@#%&; ", r)
}

func isRegexSyntax(s string) bool {
	return strings.ContainsFunc(s, isRegexSyntaxChar)
}
