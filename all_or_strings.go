package goutils

import (
	"encoding/json"
	"fmt"

	"go.yaml.in/yaml/v4"
)

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
		return json.Marshal("*")
	}

	return json.Marshal(j.list)
}

// UnmarshalYAML implements custom deserialization for the yaml.Unmarshaler interface.
// If the YAML value is the string "*", it is treated as a wildcard and sets 'all' to true and 'list' to nil.
// Otherwise, it expects a list of strings and sets 'list' accordingly, with 'all' set to false.
func (j *AllOrListString) UnmarshalYAML(value *yaml.Node) error {
	if value.Value == "*" {
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
		return "*", nil
	}

	return j.list, nil
}

// String implements the custom behavior for the fmt.Stringer interface.
func (j AllOrListString) String() string {
	if j.all {
		return "*"
	}

	return fmt.Sprintf("%v", j.list)
}
