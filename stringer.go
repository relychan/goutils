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

// IsMetaCharacter checks if the character is a word character.
// A word character is a character a-z, A-Z, 0-9, including _ (underscore) and - (hyphen).
func IsMetaCharacter[C byte | rune](c C) bool {
	return c == '-' || c == '_' ||
		IsDigit(c) ||
		IsLowerAlphabet(c) ||
		IsUpperAlphabet(c)
}

// IsDigit checks if the character is a digit.
func IsDigit[C byte | rune](c C) bool {
	return c >= '0' && c <= '9'
}

// IsLowerAlphabet checks if the character is a lowercase alphabet.
func IsLowerAlphabet[C byte | rune](c C) bool {
	return c >= 'a' && c <= 'z'
}

// IsUpperAlphabet checks if the character is an uppercase alphabet.
func IsUpperAlphabet[C byte | rune](c C) bool {
	return c >= 'A' && c <= 'Z'
}

// ParseIntInRange parses s as an integer and
// verifies that it is within some range.
// If it is invalid or out-of-range,
// it sets ok to false and returns the min value.
func ParseIntInRange[B []byte | string](s B, minValue int, maxValue int) (int, bool) {
	var x int

	for _, c := range []byte(s) {
		if !IsDigit(c) {
			return minValue, false
		}

		x = x*10 + int(c) - '0'
	}

	if x < minValue || maxValue < x {
		return minValue, false
	}

	return x, true
}
