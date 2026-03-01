package goutils

import (
	"encoding/json"
	"fmt"
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

// UnmarshalText implements [encoding.TextUnmarshaler] by compiling the
// encoded value with [regexp.Compile] when treating it as a regular expression.
func (s *Slug) UnmarshalText(bs []byte) error {
	slug := Slug(string(bs))

	err := slug.Validate()
	if err != nil {
		return err
	}

	*s = slug

	return nil
}
