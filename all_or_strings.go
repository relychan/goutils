package goutils

import (
	"encoding/json"
	"fmt"

	"go.yaml.in/yaml/v4"
)

// AllOrListString represents a list of strings or a wildcard.
type AllOrListString struct {
	// All represents the wildcard that match all items.
	all bool
	// List of string items. Ignored if all is true.
	list []string
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
