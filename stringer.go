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

// IsUpperAlphabet checks if the character is a uppercase alphabet.
func IsUpperAlphabet[C byte | rune](c C) bool {
	return c >= 'A' && c <= 'Z'
}
