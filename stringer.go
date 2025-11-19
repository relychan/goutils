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

		digit := int(c) - '0'
		if x > (maxValue-digit)/10 {
			return minValue, false
		}
		x = x*10 + digit
	}

	if x < minValue || maxValue < x {
		return minValue, false
	}

	return x, true
}
