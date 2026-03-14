// Copyright 2026 RelyChan Pte. Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

// UnmarshalText implements [encoding.TextUnmarshaler] by parsing the
// provided text into a Slug and validating that it contains only allowed characters.
func (s *Slug) UnmarshalText(bs []byte) error {
	slug := Slug(string(bs))

	err := slug.Validate()
	if err != nil {
		return err
	}

	*s = slug

	return nil
}
