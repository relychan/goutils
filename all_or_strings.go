package goutils

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"go.yaml.in/yaml/v4"
)

const wildcardSymbol = "*"

// AllOrListString is a type that represents either a wildcard ("*") meaning "all items",
// or a specific list of strings. This is useful for configuration fields where you want
// to allow users to specify either all possible values (using "*") or a subset of values.
type AllOrListString struct {
	// All represents the wildcard that matches all items.
	all bool
	// List of string items. Ignored if all is true.
	list []string
}

// NewAll creates an [AllOrListString] that accepts all values.
func NewAll() AllOrListString {
	return AllOrListString{
		all: true,
	}
}

// NewStringList creates an [AllOrListString] with a list of static strings.
func NewStringList(list []string) AllOrListString {
	return AllOrListString{
		list: list,
	}
}

// IsZero returns true if the current instance is in its zero state (neither all nor list is set).
func (j AllOrListString) IsZero() bool {
	return !j.all && len(j.list) == 0
}

// Equal checks if the target value is equal.
func (j AllOrListString) Equal(target AllOrListString) bool {
	return j.all == target.all && EqualSliceSorted(j.list, target.list)
}

// Contains reports whether the input value is accepted. It returns true if the
// wildcard "all" is set or if the input value is present in the list.
func (j AllOrListString) Contains(input string) bool {
	if j.all {
		return true
	}

	return slices.Contains(j.list, input)
}

// IsAll returns true if the value represents the wildcard ("all").
func (j AllOrListString) IsAll() bool {
	return j.all
}

// List returns the list of strings, or nil if IsAll() is true.
func (j AllOrListString) List() []string {
	return j.list
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *AllOrListString) UnmarshalJSON(data []byte) error {
	if string(data) == `"*"` {
		j.all = true
		j.list = nil

		return nil
	}

	err := json.Unmarshal(data, &j.list)
	if err != nil {
		return err
	}

	j.all = false

	return nil
}

// MarshalJSON implements json.Marshaler.
func (j AllOrListString) MarshalJSON() ([]byte, error) {
	if j.all {
		return json.Marshal(wildcardSymbol)
	}

	return json.Marshal(j.list)
}

// UnmarshalYAML implements custom deserialization for the yaml.Unmarshaler interface.
// If the YAML value is the string "*", it is treated as a wildcard and sets 'all' to true and 'list' to nil.
// Otherwise, it expects a list of strings and sets 'list' accordingly, with 'all' set to false.
func (j *AllOrListString) UnmarshalYAML(value *yaml.Node) error {
	if value.Value == wildcardSymbol {
		j.all = true
		j.list = nil

		return nil
	}

	err := value.Decode(&j.list)
	if err != nil {
		return err
	}

	j.all = false

	return nil
}

// MarshalYAML implements the custom behavior for the yaml.Marshaler interface.
// If the wildcard state is set (all == true), it serializes as the string "*".
// Otherwise, it serializes the list as a YAML sequence.
func (j AllOrListString) MarshalYAML() (any, error) {
	if j.all {
		return wildcardSymbol, nil
	}

	return j.list, nil
}

// String implements the custom behavior for the fmt.Stringer interface.
func (j AllOrListString) String() string {
	if j.all {
		return wildcardSymbol
	}

	return fmt.Sprintf("%v", j.list)
}

// Wildcard represents a string with a wildcard pattern.
type Wildcard struct {
	prefix string
	suffix string
}

// NewWildcard creates a new Wildcard instance from string.
func NewWildcard(input string) (Wildcard, bool) {
	if input == "" {
		return Wildcard{}, false
	}

	before, after, found := strings.Cut(input, wildcardSymbol)
	result := Wildcard{
		prefix: before,
		suffix: after,
	}

	if found {
		result.suffix = strings.TrimLeft(after, wildcardSymbol)
	}

	return result, found
}

// IsZero returns true if the current instance is in its zero state (both prefix and suffix are empty).
func (w Wildcard) IsZero() bool {
	return w.prefix == "" && w.suffix == ""
}

// Equal checks if the target value is equal.
func (w Wildcard) Equal(target Wildcard) bool {
	return w.prefix == target.prefix &&
		w.suffix == target.suffix
}

// Match checks if the input string matches the wildcard pattern.
func (w Wildcard) Match(s string) bool {
	return len(s) >= (len(w.prefix)+len(w.suffix)) &&
		strings.HasPrefix(s, w.prefix) && strings.HasSuffix(s, w.suffix)
}

// String implements the fmt.Stringer interface.
func (w Wildcard) String() string {
	return w.prefix + wildcardSymbol + w.suffix
}

// AllOrListWildcardString is a type that represents either a wildcard ("*") meaning "all items",
// or a specific list of strings. Accept a star (*) in any string for sub-matching patterns.
type AllOrListWildcardString struct { //nolint:recvcheck
	AllOrListString

	// List of wildcard items. Ignored if all is true.
	wildcards []Wildcard
}

// NewAllWildcard creates an AllOrListWildcardString that represents all items ("*").
func NewAllWildcard() AllOrListWildcardString {
	return AllOrListWildcardString{
		AllOrListString: NewAll(),
	}
}

// NewAllOrListWildcardStringFromStrings constructs an AllOrListWildcardString from
// the provided list of strings, parsing any wildcard patterns using the same logic
// as JSON/YAML unmarshaling.
func NewAllOrListWildcardStringFromStrings(values []string) AllOrListWildcardString {
	var res AllOrListWildcardString

	res.parseStrings(values)

	return res
}

// IsZero returns true if the current instance is in its zero state (neither all nor list is set).
func (j AllOrListWildcardString) IsZero() bool {
	return j.AllOrListString.IsZero() &&
		len(j.wildcards) == 0
}

// Equal checks if the target value is equal.
func (j AllOrListWildcardString) Equal(target AllOrListWildcardString) bool {
	return j.AllOrListString.Equal(target.AllOrListString) &&
		EqualSlice(j.wildcards, target.wildcards, false)
}

// Wildcards returns the list of wildcards, or nil if IsAll() is true.
func (j AllOrListWildcardString) Wildcards() []Wildcard {
	return j.wildcards
}

// Contains reports whether the input value is contained in this set.
// It returns true if the embedded AllOrListString is in the "all" state,
// or if the input value is present in its static list, or if the input
// matches any of the configured wildcard patterns.
func (j AllOrListWildcardString) Contains(input string) bool {
	if j.AllOrListString.Contains(input) {
		return true
	}

	for _, w := range j.wildcards {
		if w.Match(input) {
			return true
		}
	}

	return false
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *AllOrListWildcardString) UnmarshalJSON(data []byte) error {
	if string(data) == `"*"` {
		j.all = true
		j.list = nil
		j.wildcards = nil

		return nil
	}

	var list []string

	err := json.Unmarshal(data, &list)
	if err != nil {
		return err
	}

	j.all = false
	j.list = nil
	j.wildcards = nil

	j.parseStrings(list)

	return nil
}

// MarshalJSON implements json.Marshaler.
func (j AllOrListWildcardString) MarshalJSON() ([]byte, error) {
	if j.all {
		return json.Marshal(wildcardSymbol)
	}

	return json.Marshal(j.toStrings())
}

// UnmarshalYAML implements custom deserialization for the yaml.Unmarshaler interface.
// If the YAML value is the string "*", it is treated as a wildcard and sets 'all' to true and 'list' to nil.
// Otherwise, it expects a list of strings and sets 'list' accordingly, with 'all' set to false.
func (j *AllOrListWildcardString) UnmarshalYAML(value *yaml.Node) error {
	if value.Value == wildcardSymbol {
		j.all = true
		j.list = nil
		j.wildcards = nil

		return nil
	}

	var list []string

	err := value.Decode(&list)
	if err != nil {
		return err
	}

	j.all = false
	j.list = nil
	j.wildcards = nil

	j.parseStrings(list)

	return nil
}

// MarshalYAML implements the custom behavior for the yaml.Marshaler interface.
// If the wildcard state is set (all == true), it serializes as the string "*".
// Otherwise, it serializes the list as a YAML sequence.
func (j AllOrListWildcardString) MarshalYAML() (any, error) {
	if j.all {
		return wildcardSymbol, nil
	}

	return j.toStrings(), nil
}

// String implements the custom behavior for the fmt.Stringer interface.
func (j AllOrListWildcardString) String() string {
	if j.all {
		return wildcardSymbol
	}

	var sb strings.Builder
	sb.WriteByte('[')

	for i, str := range j.list {
		sb.WriteString(str)

		if i < len(j.list)-1 {
			sb.WriteString(", ")
		}
	}

	for _, w := range j.wildcards {
		sb.WriteString(", ")
		sb.WriteString(w.String())
	}

	sb.WriteByte(']')

	return sb.String()
}

func (j AllOrListWildcardString) toStrings() []string {
	result := make([]string, 0, len(j.list)+len(j.wildcards))
	result = append(result, j.list...)

	for _, w := range j.wildcards {
		result = append(result, w.String())
	}

	return result
}

func (j *AllOrListWildcardString) parseStrings(inputs []string) {
	for _, input := range inputs {
		j.parseString(input)

		if j.all {
			return
		}
	}
}

func (j *AllOrListWildcardString) parseString(input string) {
	if input == "" {
		j.list = append(j.list, "")

		return
	}

	if input == wildcardSymbol {
		// If "*" is present in the list, turn the whole list into a match all
		j.all = true
		j.list = nil
		j.wildcards = nil

		return
	}

	w, ok := NewWildcard(input)
	if ok {
		j.wildcards = append(j.wildcards, w)

		return
	}

	j.list = append(j.list, input)
}
