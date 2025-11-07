package goutils

import (
	"encoding/json"
	"fmt"

	"go.yaml.in/yaml/v4"
)

// Slug represents a url-encoded string that allows alphabet, digits, hyphens and underscores only.
type Slug string

// NewSlug creates a new slug instance from string.
func NewSlug(value string) Slug {
	return Slug(value)
}

// Validate checks if the slug is valid.
func (s Slug) Validate() error {
	for _, c := range s {
		if !IsMetaCharacter(c) {
			return fmt.Errorf("%w, character '%s' is not allowed", ErrInvalidSlug, string(c))
		}
	}

	return nil
}

// String implements the fmt.Stringer interface.
func (s Slug) String() string {
	return string(s)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *Slug) UnmarshalJSON(bytes []byte) error {
	var rawValue string

	err := json.Unmarshal(bytes, &rawValue)
	if err != nil {
		return err
	}

	slug := Slug(rawValue)

	err = slug.Validate()
	if err != nil {
		return err
	}

	*s = slug

	return nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface to decode value.
func (s *Slug) UnmarshalYAML(value *yaml.Node) error {
	var rawValue string

	err := value.Decode(&rawValue)
	if err != nil {
		return fmt.Errorf("invalid slug, %w", err)
	}

	slug := Slug(value.Value)

	err = slug.Validate()
	if err != nil {
		return err
	}

	*s = slug

	return nil
}
