package goutils

import (
	"encoding/json"
	"testing"

	"go.yaml.in/yaml/v4"
)

func TestNewRegexpMatcher(t *testing.T) {
	t.Run("text with space", func(t *testing.T) {
		rm := MustRegexpMatcher("hello world")

		if rm.text == nil {
			t.Error("expected text to be set for string with space")
		}

		if rm.regexp != nil {
			t.Error("expected regexp to be nil")
		}

		if *rm.text != "hello world" {
			t.Errorf("expected text 'hello world', got: %s", *rm.text)
		}
	})

	t.Run("text with only meta chars", func(t *testing.T) {
		// "hello" contains only meta characters, so it's treated as text
		rm, err := NewRegexpMatcher("hello")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if rm.text == nil {
			t.Error("expected text to be set")
		}

		if rm.regexp != nil {
			t.Error("expected regexp to be nil")
		}

		if *rm.text != "hello" {
			t.Errorf("expected text 'hello', got: %s", *rm.text)
		}
	})

	t.Run("regexp pattern with regex chars", func(t *testing.T) {
		rm, err := NewRegexpMatcher("^[a-z]+$")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if rm.text != nil {
			t.Error("expected text to be nil")
		}

		if rm.regexp == nil {
			t.Error("expected regexp to be set")
		}

		if rm.regexp.String() != "^[a-z]+$" {
			t.Errorf("expected regexp '^[a-z]+$', got: %s", rm.regexp.String())
		}
	})

	t.Run("invalid regexp", func(t *testing.T) {
		_, err := NewRegexpMatcher("[invalid")
		if err == nil {
			t.Error("expected error for invalid regexp")
		}
	})

	t.Run("empty string", func(t *testing.T) {
		rm, err := NewRegexpMatcher("")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Empty string has no special chars, treated as text
		if rm.text == nil {
			t.Error("expected text to be set for empty string")
		}

		if *rm.text != "" {
			t.Errorf("expected empty text, got: %s", *rm.text)
		}
	})

	t.Run("only meta characters - treated as text", func(t *testing.T) {
		// Only alphanumeric, underscore, hyphen -> treated as text
		rm, err := NewRegexpMatcher("hello_world-123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if rm.text == nil {
			t.Error("expected text to be set")
		}

		if rm.regexp != nil {
			t.Error("expected regexp to be nil")
		}

		if *rm.text != "hello_world-123" {
			t.Errorf("expected text 'hello_world-123', got: %s", *rm.text)
		}
	})

	t.Run("text with allowed special chars", func(t *testing.T) {
		// Characters in `~@#%&;  are treated as text
		rm, err := NewRegexpMatcher("hello@world")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if rm.text == nil {
			t.Error("expected text to be set")
		}

		if *rm.text != "hello@world" {
			t.Errorf("expected text 'hello@world', got: %s", *rm.text)
		}
	})
}

func TestRegexpMatcher_IsZero(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		var rm RegexpMatcher
		if !rm.IsZero() {
			t.Error("expected IsZero() to be true for zero value")
		}
	})

	t.Run("with text", func(t *testing.T) {
		rm, _ := NewRegexpMatcher("hello world")
		if rm.IsZero() {
			t.Error("expected IsZero() to be false with text")
		}
	})

	t.Run("with regexp", func(t *testing.T) {
		rm, _ := NewRegexpMatcher("^test$")
		if rm.IsZero() {
			t.Error("expected IsZero() to be false with regexp")
		}
	})
}

func TestRegexpMatcher_Equal(t *testing.T) {
	t.Run("both zero", func(t *testing.T) {
		var rm1, rm2 RegexpMatcher
		if !rm1.Equal(rm2) {
			t.Error("expected zero values to be equal")
		}
	})

	t.Run("same text", func(t *testing.T) {
		rm1, _ := NewRegexpMatcher("hello world")
		rm2, _ := NewRegexpMatcher("hello world")
		if !rm1.Equal(*rm2) {
			t.Error("expected same text to be equal")
		}
	})

	t.Run("different text", func(t *testing.T) {
		rm1, _ := NewRegexpMatcher("hello world")
		rm2, _ := NewRegexpMatcher("goodbye world")
		if rm1.Equal(*rm2) {
			t.Error("expected different text to not be equal")
		}
	})

	t.Run("same regexp", func(t *testing.T) {
		rm1, _ := NewRegexpMatcher("^hello$")
		rm2, _ := NewRegexpMatcher("^hello$")
		if !rm1.Equal(*rm2) {
			t.Error("expected same regexp to be equal")
		}
	})

	t.Run("different regexp", func(t *testing.T) {
		rm1, _ := NewRegexpMatcher("^hello$")
		rm2, _ := NewRegexpMatcher("^world$")
		if rm1.Equal(*rm2) {
			t.Error("expected different regexp to not be equal")
		}
	})

	t.Run("text vs regexp", func(t *testing.T) {
		rm1, _ := NewRegexpMatcher("hello world")
		rm2, _ := NewRegexpMatcher("^hello$")
		if rm1.Equal(*rm2) {
			t.Error("expected text and regexp to not be equal")
		}
	})

	t.Run("zero vs non-zero", func(t *testing.T) {
		var rm1 RegexpMatcher
		rm2, _ := NewRegexpMatcher("hello")
		if rm1.Equal(*rm2) {
			t.Error("expected zero and non-zero to not be equal")
		}
	})
}

func TestRegexpMatcher_String(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		var rm RegexpMatcher
		if rm.String() != "" {
			t.Errorf("expected empty string, got: %s", rm.String())
		}
	})

	t.Run("with text", func(t *testing.T) {
		rm, _ := NewRegexpMatcher("hello world")
		if rm.String() != "hello world" {
			t.Errorf("expected 'hello world', got: %s", rm.String())
		}
	})

	t.Run("with regexp", func(t *testing.T) {
		rm, _ := NewRegexpMatcher("^[a-z]+$")
		if rm.String() != "^[a-z]+$" {
			t.Errorf("expected '^[a-z]+$', got: %s", rm.String())
		}
	})
}

func TestRegexpMatcher_UnmarshalText(t *testing.T) {
	t.Run("text with special chars", func(t *testing.T) {
		var rm RegexpMatcher
		err := rm.UnmarshalText([]byte("hello world"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if rm.text == nil || *rm.text != "hello world" {
			t.Error("expected text to be 'hello world'")
		}

		if rm.regexp != nil {
			t.Error("expected regexp to be nil")
		}
	})

	t.Run("only meta chars - treated as text", func(t *testing.T) {
		var rm RegexpMatcher
		err := rm.UnmarshalText([]byte("hello"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if rm.text == nil || *rm.text != "hello" {
			t.Error("expected text to be 'hello'")
		}

		if rm.regexp != nil {
			t.Error("expected regexp to be nil")
		}
	})

	t.Run("regexp pattern with regex chars", func(t *testing.T) {
		var rm RegexpMatcher
		err := rm.UnmarshalText([]byte("^[a-z]+$"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if rm.text != nil {
			t.Error("expected text to be nil")
		}

		if rm.regexp == nil {
			t.Error("expected regexp to be set")
		}
	})

	t.Run("invalid regexp", func(t *testing.T) {
		var rm RegexpMatcher
		err := rm.UnmarshalText([]byte("[invalid"))
		if err == nil {
			t.Error("expected error for invalid regexp")
		}
	})

	t.Run("empty bytes - treated as text", func(t *testing.T) {
		var rm RegexpMatcher
		err := rm.UnmarshalText([]byte(""))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if rm.text == nil || *rm.text != "" {
			t.Error("expected empty text")
		}
	})

	t.Run("text with spaces", func(t *testing.T) {
		var rm RegexpMatcher
		err := rm.UnmarshalText([]byte("hello world"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if rm.text == nil || *rm.text != "hello world" {
			t.Error("expected text to be 'hello world'")
		}
	})

	t.Run("text with allowed special chars", func(t *testing.T) {
		var rm RegexpMatcher
		err := rm.UnmarshalText([]byte("test@example"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if rm.text == nil || *rm.text != "test@example" {
			t.Error("expected text to be 'test@example'")
		}
	})

	t.Run("regexp with dot", func(t *testing.T) {
		var rm RegexpMatcher
		err := rm.UnmarshalText([]byte("test.example"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// "." is not a meta character and not in allowed set, so it's compiled as regexp
		if rm.regexp == nil {
			t.Error("expected regexp to be set")
		}

		if rm.text != nil {
			t.Error("expected text to be nil")
		}
	})
}

func TestRegexpMatcher_JSON(t *testing.T) {
	t.Run("unmarshal text", func(t *testing.T) {
		var rm RegexpMatcher
		err := json.Unmarshal([]byte(`"hello world"`), &rm)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if rm.text == nil || *rm.text != "hello world" {
			t.Error("expected text to be 'hello world'")
		}
	})

	t.Run("marshal text", func(t *testing.T) {
		rm, _ := NewRegexpMatcher("hello world")
		data, err := json.Marshal(rm)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if string(data) != `"hello world"` {
			t.Errorf("expected '\"hello world\"', got: %s", string(data))
		}
	})

	t.Run("unmarshal regexp", func(t *testing.T) {
		var rm RegexpMatcher
		err := json.Unmarshal([]byte(`"^[a-z]+$"`), &rm)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if rm.regexp == nil {
			t.Error("expected regexp to be set")
		}
	})

	t.Run("marshal regexp", func(t *testing.T) {
		rm, _ := NewRegexpMatcher("^[a-z]+$")
		data, err := json.Marshal(rm)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if string(data) != `"^[a-z]+$"` {
			t.Errorf("expected '\"^[a-z]+$\"', got: %s", string(data))
		}
	})

	t.Run("unmarshal invalid JSON", func(t *testing.T) {
		var rm RegexpMatcher
		err := json.Unmarshal([]byte(`{invalid}`), &rm)
		if err == nil {
			t.Error("expected error for invalid JSON")
		}
	})

	t.Run("round trip text", func(t *testing.T) {
		original, _ := NewRegexpMatcher("hello world!")
		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var decoded RegexpMatcher
		err = json.Unmarshal(data, &decoded)
		if err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}

		if !original.Equal(decoded) {
			t.Error("round trip failed for text")
		}
	})

	t.Run("round trip regexp", func(t *testing.T) {
		original, _ := NewRegexpMatcher("^test.*end$")
		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var decoded RegexpMatcher
		err = json.Unmarshal(data, &decoded)
		if err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}

		if !original.Equal(decoded) {
			t.Error("round trip failed for regexp")
		}
	})
}

func TestRegexpMatcher_YAML(t *testing.T) {
	t.Run("unmarshal text with spaces", func(t *testing.T) {
		var rm RegexpMatcher
		err := yaml.Unmarshal([]byte(`"hello world"`), &rm)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if rm.text == nil || *rm.text != "hello world" {
			t.Error("expected text to be 'hello world'")
		}
	})

	t.Run("marshal text", func(t *testing.T) {
		rm, _ := NewRegexpMatcher("hello world")
		data, err := yaml.Marshal(rm)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := "hello world\n"
		if string(data) != expected {
			t.Errorf("expected %q, got: %q", expected, string(data))
		}
	})

	t.Run("unmarshal regexp", func(t *testing.T) {
		var rm RegexpMatcher
		err := yaml.Unmarshal([]byte(`"^[a-z]+$"`), &rm)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if rm.regexp == nil {
			t.Error("expected regexp to be set")
		}
	})

	t.Run("marshal regexp", func(t *testing.T) {
		rm, _ := NewRegexpMatcher("^[a-z]+$")
		data, err := yaml.Marshal(rm)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := "^[a-z]+$\n"
		if string(data) != expected {
			t.Errorf("expected %q, got: %q", expected, string(data))
		}
	})

	t.Run("round trip text", func(t *testing.T) {
		original, _ := NewRegexpMatcher("hello world!")
		data, err := yaml.Marshal(original)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var decoded RegexpMatcher
		err = yaml.Unmarshal(data, &decoded)
		if err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}

		if !original.Equal(decoded) {
			t.Error("round trip failed for text")
		}
	})

	t.Run("round trip regexp", func(t *testing.T) {
		original, _ := NewRegexpMatcher("^test.*end$")
		data, err := yaml.Marshal(original)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var decoded RegexpMatcher
		err = yaml.Unmarshal(data, &decoded)
		if err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}

		if !original.Equal(decoded) {
			t.Error("round trip failed for regexp")
		}
	})
}

func TestRegexpMatcher_Match(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		var rm RegexpMatcher
		if rm.Match([]byte("anything")) {
			t.Error("expected zero value to not match")
		}
	})

	t.Run("text match", func(t *testing.T) {
		rm, _ := NewRegexpMatcher("hello world")
		tests := []struct {
			input    []byte
			expected bool
		}{
			{[]byte("hello world"), true},
			{[]byte("say hello world there"), true},
			{[]byte("hello world!"), true},
			{[]byte("HELLO WORLD"), false},
			{[]byte("hello"), false},
			{[]byte(""), false},
		}

		for _, tt := range tests {
			result := rm.Match(tt.input)
			if result != tt.expected {
				t.Errorf("Match(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		}
	})

	t.Run("regexp match - simple pattern", func(t *testing.T) {
		rm, _ := NewRegexpMatcher("^[a-z]+$")
		tests := []struct {
			input    []byte
			expected bool
		}{
			{[]byte("hello"), true},
			{[]byte("world"), true},
			{[]byte("abc"), true},
			{[]byte("Hello"), false},
			{[]byte("hello123"), false},
			{[]byte("hello world"), false},
			{[]byte(""), false},
		}

		for _, tt := range tests {
			result := rm.Match(tt.input)
			if result != tt.expected {
				t.Errorf("Match(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		}
	})

	t.Run("regexp with special patterns", func(t *testing.T) {
		rm, _ := NewRegexpMatcher("test.*end")
		tests := []struct {
			input    []byte
			expected bool
		}{
			{[]byte("testend"), true},
			{[]byte("test123end"), true},
			{[]byte("test anything here end"), true},
			{[]byte("prefix testend"), true},
			{[]byte("testend suffix"), true},
			{[]byte("test"), false},
			{[]byte("end"), false},
		}

		for _, tt := range tests {
			result := rm.Match(tt.input)
			if result != tt.expected {
				t.Errorf("Match(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		}
	})
}

func TestRegexpMatcher_MatchString(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		var rm RegexpMatcher
		if rm.MatchString("anything") {
			t.Error("expected zero value to not match")
		}
	})

	t.Run("text match", func(t *testing.T) {
		rm, _ := NewRegexpMatcher("hello world")
		tests := []struct {
			input    string
			expected bool
		}{
			{"hello world", true},
			{"say hello world there", true},
			{"hello world!", true},
			{"HELLO WORLD", false},
			{"hello", false},
			{"", false},
		}

		for _, tt := range tests {
			result := rm.MatchString(tt.input)
			if result != tt.expected {
				t.Errorf("MatchString(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		}
	})

	t.Run("regexp match - simple pattern", func(t *testing.T) {
		rm, _ := NewRegexpMatcher("^[a-z]+$")
		tests := []struct {
			input    string
			expected bool
		}{
			{"hello", true},
			{"world", true},
			{"abc", true},
			{"Hello", false},
			{"hello123", false},
			{"hello world", false},
			{"", false},
		}

		for _, tt := range tests {
			result := rm.MatchString(tt.input)
			if result != tt.expected {
				t.Errorf("MatchString(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		}
	})

	t.Run("regexp with digits pattern", func(t *testing.T) {
		rm, _ := NewRegexpMatcher(`^\d{3}-\d{4}$`)
		tests := []struct {
			input    string
			expected bool
		}{
			{"123-4567", true},
			{"000-0000", true},
			{"999-9999", true},
			{"12-4567", false},
			{"123-456", false},
			{"abc-defg", false},
			{"", false},
		}

		for _, tt := range tests {
			result := rm.MatchString(tt.input)
			if result != tt.expected {
				t.Errorf("MatchString(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		}
	})

	t.Run("text with underscores and hyphens", func(t *testing.T) {
		rm, _ := NewRegexpMatcher("test_value-123")
		tests := []struct {
			input    string
			expected bool
		}{
			{"test_value-123", true},
			{"prefix test_value-123", true},
			{"test_value-123 suffix", true},
			{"test_value-124", false},
			{"test-value-123", false},
			{"", false},
		}

		for _, tt := range tests {
			result := rm.MatchString(tt.input)
			if result != tt.expected {
				t.Errorf("MatchString(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		}
	})
}

func TestRegexpMatcher_EdgeCases(t *testing.T) {
	t.Run("empty pattern text", func(t *testing.T) {
		rm, _ := NewRegexpMatcher("")
		// Empty string is treated as text
		if !rm.MatchString("") {
			t.Error("expected empty pattern to match empty string")
		}
		if !rm.MatchString("anything") {
			t.Error("expected empty pattern to match any string (contains empty)")
		}
	})

	t.Run("regexp with brackets", func(t *testing.T) {
		// "[" and "]" are not meta characters and not in allowed set, so it's compiled as regexp
		rm, _ := NewRegexpMatcher(`hello [world]`)
		tests := []struct {
			input    string
			expected bool
		}{
			{"hello w", true},        // [world] matches any single char from "world"
			{"hello o", true},        // [world] matches "o"
			{"hello x", false},       // [world] doesn't match "x"
			{"hello [world]", false}, // literal brackets don't match
			{"", false},
		}

		for _, tt := range tests {
			result := rm.MatchString(tt.input)
			if result != tt.expected {
				t.Errorf("MatchString(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		}
	})

	t.Run("unicode text", func(t *testing.T) {
		rm, _ := NewRegexpMatcher("hello 世界")
		if !rm.MatchString("hello 世界") {
			t.Error("expected to match unicode text")
		}
		if !rm.MatchString("say hello 世界 there") {
			t.Error("expected to match unicode text in substring")
		}
	})

	t.Run("case sensitivity", func(t *testing.T) {
		rm, _ := NewRegexpMatcher("Hello World")
		if !rm.MatchString("Hello World") {
			t.Error("expected exact match")
		}
		if rm.MatchString("hello world") {
			t.Error("expected case-sensitive match to fail")
		}
	})

	t.Run("regexp with anchors", func(t *testing.T) {
		rm, _ := NewRegexpMatcher("^start.*end$")
		tests := []struct {
			input    string
			expected bool
		}{
			{"startend", true},
			{"start middle end", true},
			{"start", false},
			{"end", false},
			{"prefix start end", false},
			{"start end suffix", false},
		}

		for _, tt := range tests {
			result := rm.MatchString(tt.input)
			if result != tt.expected {
				t.Errorf("MatchString(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		}
	})

	t.Run("text with backtick", func(t *testing.T) {
		rm, _ := NewRegexpMatcher("hello`world")
		if rm.text == nil || *rm.text != "hello`world" {
			t.Error("expected text with backtick to be stored as text")
		}
		if !rm.MatchString("hello`world") {
			t.Error("expected to match text with backtick")
		}
	})

	t.Run("text with tilde", func(t *testing.T) {
		rm, _ := NewRegexpMatcher("hello~world")
		if rm.text == nil || *rm.text != "hello~world" {
			t.Error("expected text with tilde to be stored as text")
		}
	})

	t.Run("text with at sign", func(t *testing.T) {
		rm, _ := NewRegexpMatcher("user@example")
		if rm.text == nil || *rm.text != "user@example" {
			t.Error("expected text with @ to be stored as text")
		}
	})

	t.Run("regexp with dot", func(t *testing.T) {
		// "." is not a meta character and not in allowed set, so it's compiled as regexp
		rm, _ := NewRegexpMatcher("user.example")
		if rm.regexp == nil {
			t.Error("expected regexp to be set for text with dot")
		}
		if rm.text != nil {
			t.Error("expected text to be nil")
		}
	})

	t.Run("text with hash", func(t *testing.T) {
		rm, _ := NewRegexpMatcher("hello#world")
		if rm.text == nil || *rm.text != "hello#world" {
			t.Error("expected text with hash to be stored as text")
		}
	})

	t.Run("text with percent", func(t *testing.T) {
		rm, _ := NewRegexpMatcher("50%")
		if rm.text == nil || *rm.text != "50%" {
			t.Error("expected text with percent to be stored as text")
		}
	})

	t.Run("text with ampersand", func(t *testing.T) {
		rm, _ := NewRegexpMatcher("hello&world")
		if rm.text == nil || *rm.text != "hello&world" {
			t.Error("expected text with ampersand to be stored as text")
		}
	})

	t.Run("text with semicolon", func(t *testing.T) {
		rm, _ := NewRegexpMatcher("hello;world")
		if rm.text == nil || *rm.text != "hello;world" {
			t.Error("expected text with semicolon to be stored as text")
		}
	})
}
