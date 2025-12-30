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
		if rm.MatchString("anything") {
			t.Error("expected empty pattern to not match any string")
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

func TestMustRegexpMatcher(t *testing.T) {
	t.Run("valid pattern", func(t *testing.T) {
		rm := MustRegexpMatcher("hello world")
		if rm == nil {
			t.Error("expected non-nil result")
		}
		if rm.text == nil || *rm.text != "hello world" {
			t.Error("expected text to be 'hello world'")
		}
	})

	t.Run("valid regexp", func(t *testing.T) {
		rm := MustRegexpMatcher("^[a-z]+$")
		if rm == nil {
			t.Error("expected non-nil result")
		}
		if rm.regexp == nil {
			t.Error("expected regexp to be set")
		}
	})

	t.Run("panic on invalid regexp", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for invalid regexp")
			}
		}()
		MustRegexpMatcher("[invalid")
	})
}

func TestRegexpMatcher_PrefixOp(t *testing.T) {
	t.Run("prefix match with caret", func(t *testing.T) {
		rm := MustRegexpMatcher("^hello")
		tests := []struct {
			input    string
			expected bool
		}{
			{"hello", true},
			{"hello world", true},
			{"helloworld", true},
			{"say hello", false},
			{"HELLO", false},
			{"", false},
		}

		for _, tt := range tests {
			result := rm.MatchString(tt.input)
			if result != tt.expected {
				t.Errorf("MatchString(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		}
	})

	t.Run("prefix with meta characters only", func(t *testing.T) {
		rm := MustRegexpMatcher("^test123")
		if !rm.MatchString("test123") {
			t.Error("expected to match 'test123'")
		}
		if !rm.MatchString("test123abc") {
			t.Error("expected to match 'test123abc'")
		}
		if rm.MatchString("abc test123") {
			t.Error("expected not to match 'abc test123'")
		}
	})

	t.Run("prefix with special chars becomes regexp", func(t *testing.T) {
		rm := MustRegexpMatcher("^hello.world")
		// "." is regex syntax, so this becomes a regexp
		if rm.regexp == nil {
			t.Error("expected regexp to be set")
		}
	})
}

func TestRegexpMatcher_SuffixOp(t *testing.T) {
	t.Run("suffix match with dollar", func(t *testing.T) {
		rm := MustRegexpMatcher("world$")
		tests := []struct {
			input    string
			expected bool
		}{
			{"world", true},
			{"hello world", true},
			{"helloworld", true},
			{"world hello", false},
			{"WORLD", false},
			{"", false},
		}

		for _, tt := range tests {
			result := rm.MatchString(tt.input)
			if result != tt.expected {
				t.Errorf("MatchString(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		}
	})

	t.Run("suffix with meta characters only", func(t *testing.T) {
		rm := MustRegexpMatcher("test123$")
		if !rm.MatchString("test123") {
			t.Error("expected to match 'test123'")
		}
		if !rm.MatchString("abc test123") {
			t.Error("expected to match 'abc test123'")
		}
		if rm.MatchString("test123 abc") {
			t.Error("expected not to match 'test123 abc'")
		}
	})

	t.Run("suffix with special chars becomes regexp", func(t *testing.T) {
		rm := MustRegexpMatcher("hello.world$")
		// "." is regex syntax, so this becomes a regexp
		if rm.regexp == nil {
			t.Error("expected regexp to be set")
		}
	})
}

func TestRegexpMatcher_EqualOp(t *testing.T) {
	t.Run("exact match with anchors", func(t *testing.T) {
		rm := MustRegexpMatcher("^hello$")
		tests := []struct {
			input    string
			expected bool
		}{
			{"hello", true},
			{"hello world", false},
			{"say hello", false},
			{"helloworld", false},
			{"HELLO", false},
			{"", false},
		}

		for _, tt := range tests {
			result := rm.MatchString(tt.input)
			if result != tt.expected {
				t.Errorf("MatchString(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		}
	})

	t.Run("exact match with spaces", func(t *testing.T) {
		rm := MustRegexpMatcher("^hello world$")
		if !rm.MatchString("hello world") {
			t.Error("expected to match 'hello world'")
		}
		if rm.MatchString("hello world!") {
			t.Error("expected not to match 'hello world!'")
		}
		if rm.MatchString("say hello world") {
			t.Error("expected not to match 'say hello world'")
		}
	})

	t.Run("exact match with meta characters", func(t *testing.T) {
		rm := MustRegexpMatcher("^test_value-123$")
		if !rm.MatchString("test_value-123") {
			t.Error("expected to match 'test_value-123'")
		}
		if rm.MatchString("test_value-123 ") {
			t.Error("expected not to match with trailing space")
		}
		if rm.MatchString(" test_value-123") {
			t.Error("expected not to match with leading space")
		}
	})

	t.Run("exact match with regex chars becomes regexp", func(t *testing.T) {
		rm := MustRegexpMatcher("^hello.world$")
		// "." is regex syntax, so this becomes a regexp
		if rm.regexp == nil {
			t.Error("expected regexp to be set")
		}
	})
}

func TestRegexpMatcher_ContainOp(t *testing.T) {
	t.Run("contain match default", func(t *testing.T) {
		rm := MustRegexpMatcher("hello")
		tests := []struct {
			input    string
			expected bool
		}{
			{"hello", true},
			{"hello world", true},
			{"say hello there", true},
			{"helloworld", true},
			{"HELLO", false},
			{"hell", false},
			{"", false},
		}

		for _, tt := range tests {
			result := rm.MatchString(tt.input)
			if result != tt.expected {
				t.Errorf("MatchString(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		}
	})

	t.Run("contain with spaces", func(t *testing.T) {
		rm := MustRegexpMatcher("hello world")
		if !rm.MatchString("hello world") {
			t.Error("expected to match 'hello world'")
		}
		if !rm.MatchString("say hello world there") {
			t.Error("expected to match 'say hello world there'")
		}
		if rm.MatchString("hello  world") {
			t.Error("expected not to match with double space")
		}
	})

	t.Run("contain with special allowed chars", func(t *testing.T) {
		rm := MustRegexpMatcher("user@example")
		if !rm.MatchString("user@example") {
			t.Error("expected to match 'user@example'")
		}
		if !rm.MatchString("contact: user@example.com") {
			t.Error("expected to match as substring")
		}
	})
}

func TestRegexpMatcher_SingleCharacter(t *testing.T) {
	t.Run("single meta character", func(t *testing.T) {
		rm := MustRegexpMatcher("a")
		if !rm.MatchString("a") {
			t.Error("expected to match 'a'")
		}
		if !rm.MatchString("abc") {
			t.Error("expected to match 'abc'")
		}
		if rm.MatchString("bc") {
			t.Error("expected not to match 'bc'")
		}
	})

	t.Run("single regex char", func(t *testing.T) {
		rm := MustRegexpMatcher(".")
		// "." is regex syntax, becomes regexp
		if rm.regexp == nil {
			t.Error("expected regexp to be set")
		}
		if !rm.MatchString("a") {
			t.Error("expected to match any character")
		}
		if !rm.MatchString("x") {
			t.Error("expected to match any character")
		}
	})

	t.Run("single special allowed char", func(t *testing.T) {
		rm := MustRegexpMatcher("@")
		if !rm.MatchString("@") {
			t.Error("expected to match '@'")
		}
		if !rm.MatchString("user@example") {
			t.Error("expected to match as substring")
		}
	})
}

func TestRegexpMatcher_MatchOperations(t *testing.T) {
	t.Run("Match with prefix op", func(t *testing.T) {
		rm := MustRegexpMatcher("^hello")
		tests := []struct {
			input    []byte
			expected bool
		}{
			{[]byte("hello"), true},
			{[]byte("hello world"), true},
			{[]byte("say hello"), false},
		}

		for _, tt := range tests {
			result := rm.Match(tt.input)
			if result != tt.expected {
				t.Errorf("Match(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		}
	})

	t.Run("Match with suffix op", func(t *testing.T) {
		rm := MustRegexpMatcher("world$")
		tests := []struct {
			input    []byte
			expected bool
		}{
			{[]byte("world"), true},
			{[]byte("hello world"), true},
			{[]byte("world hello"), false},
		}

		for _, tt := range tests {
			result := rm.Match(tt.input)
			if result != tt.expected {
				t.Errorf("Match(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		}
	})

	t.Run("Match with equal op", func(t *testing.T) {
		rm := MustRegexpMatcher("^hello$")
		tests := []struct {
			input    []byte
			expected bool
		}{
			{[]byte("hello"), true},
			{[]byte("hello world"), false},
			{[]byte("say hello"), false},
		}

		for _, tt := range tests {
			result := rm.Match(tt.input)
			if result != tt.expected {
				t.Errorf("Match(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		}
	})

	t.Run("Match with contain op", func(t *testing.T) {
		rm := MustRegexpMatcher("hello")
		tests := []struct {
			input    []byte
			expected bool
		}{
			{[]byte("hello"), true},
			{[]byte("hello world"), true},
			{[]byte("say hello"), true},
			{[]byte("hell"), false},
		}

		for _, tt := range tests {
			result := rm.Match(tt.input)
			if result != tt.expected {
				t.Errorf("Match(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		}
	})
}

func TestRegexpMatcher_ComplexPatterns(t *testing.T) {
	t.Run("prefix with allowed special chars", func(t *testing.T) {
		rm := MustRegexpMatcher("^user@")
		if !rm.MatchString("user@example") {
			t.Error("expected to match 'user@example'")
		}
		if !rm.MatchString("user@") {
			t.Error("expected to match 'user@'")
		}
		if rm.MatchString("admin@example") {
			t.Error("expected not to match 'admin@example'")
		}
	})

	t.Run("suffix with allowed special chars", func(t *testing.T) {
		rm := MustRegexpMatcher("@example$")
		if !rm.MatchString("user@example") {
			t.Error("expected to match 'user@example'")
		}
		if !rm.MatchString("@example") {
			t.Error("expected to match '@example'")
		}
		if rm.MatchString("@example.com") {
			t.Error("expected not to match '@example.com'")
		}
	})

	t.Run("exact match with allowed special chars", func(t *testing.T) {
		rm := MustRegexpMatcher("^user@example$")
		if !rm.MatchString("user@example") {
			t.Error("expected to match 'user@example'")
		}
		if rm.MatchString("user@example.com") {
			t.Error("expected not to match 'user@example.com'")
		}
		if rm.MatchString("admin@example") {
			t.Error("expected not to match 'admin@example'")
		}
	})

	t.Run("pattern with multiple special chars", func(t *testing.T) {
		rm := MustRegexpMatcher("test@#%")
		if !rm.MatchString("test@#%") {
			t.Error("expected to match 'test@#%'")
		}
		if !rm.MatchString("prefix test@#% suffix") {
			t.Error("expected to match as substring")
		}
	})

	t.Run("empty pattern with anchors", func(t *testing.T) {
		rm := MustRegexpMatcher("^$")
		if !rm.MatchString("") {
			t.Error("expected to match empty string")
		}
		if rm.MatchString("a") {
			t.Error("expected not to match non-empty string")
		}
	})

	t.Run("only caret", func(t *testing.T) {
		rm := MustRegexpMatcher("^")
		// "^" alone is regex syntax
		if rm.regexp == nil {
			t.Error("expected regexp to be set")
		}
	})

	t.Run("only dollar", func(t *testing.T) {
		rm := MustRegexpMatcher("$")
		// "$" alone is regex syntax
		if rm.regexp == nil {
			t.Error("expected regexp to be set")
		}
	})
}

func TestRegexpMatcher_StringRepresentation(t *testing.T) {
	t.Run("prefix op string", func(t *testing.T) {
		rm := MustRegexpMatcher("^hello")
		// String() reconstructs the original pattern with anchors
		if rm.String() != "^hello" {
			t.Errorf("expected '^hello', got: %s", rm.String())
		}
	})

	t.Run("suffix op string", func(t *testing.T) {
		rm := MustRegexpMatcher("world$")
		// String() reconstructs the original pattern with anchors
		if rm.String() != "world$" {
			t.Errorf("expected 'world$', got: %s", rm.String())
		}
	})

	t.Run("equal op string", func(t *testing.T) {
		rm := MustRegexpMatcher("^hello$")
		// String() reconstructs the original pattern with anchors
		if rm.String() != "^hello$" {
			t.Errorf("expected '^hello$', got: %s", rm.String())
		}
	})

	t.Run("contain op string", func(t *testing.T) {
		rm := MustRegexpMatcher("hello")
		if rm.String() != "hello" {
			t.Errorf("expected 'hello', got: %s", rm.String())
		}
	})

	t.Run("regexp op string", func(t *testing.T) {
		rm := MustRegexpMatcher("^[a-z]+$")
		if rm.String() != "^[a-z]+$" {
			t.Errorf("expected '^[a-z]+$', got: %s", rm.String())
		}
	})
}

func TestRegexpMatcher_JSONRoundTrip(t *testing.T) {
	t.Run("round trip prefix", func(t *testing.T) {
		original := MustRegexpMatcher("^hello")
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
			t.Error("round trip failed for prefix")
		}
	})

	t.Run("round trip suffix", func(t *testing.T) {
		original := MustRegexpMatcher("world$")
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
			t.Error("round trip failed for suffix")
		}
	})

	t.Run("round trip equal", func(t *testing.T) {
		original := MustRegexpMatcher("^hello$")
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
			t.Error("round trip failed for equal")
		}
	})
}

func TestRegexpMatcher_YAMLRoundTrip(t *testing.T) {
	t.Run("round trip prefix", func(t *testing.T) {
		original := MustRegexpMatcher("^hello")
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
			t.Error("round trip failed for prefix")
		}
	})

	t.Run("round trip suffix", func(t *testing.T) {
		original := MustRegexpMatcher("world$")
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
			t.Error("round trip failed for suffix")
		}
	})

	t.Run("round trip equal", func(t *testing.T) {
		original := MustRegexpMatcher("^hello$")
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
			t.Error("round trip failed for equal")
		}
	})
}

func TestRegexpMatcher_BoundaryConditions(t *testing.T) {
	t.Run("very long text pattern", func(t *testing.T) {
		longText := "this_is_a_very_long_text_pattern_with_only_meta_characters_0123456789"
		rm := MustRegexpMatcher(longText)
		if !rm.MatchString(longText) {
			t.Error("expected to match long text")
		}
	})

	t.Run("very long regexp pattern", func(t *testing.T) {
		pattern := "^test.*with.*many.*wildcards.*and.*patterns$"
		rm := MustRegexpMatcher(pattern)
		if !rm.MatchString("test with many wildcards and patterns") {
			t.Error("expected to match long regexp")
		}
	})

	t.Run("unicode in text", func(t *testing.T) {
		rm := MustRegexpMatcher("hello世界")
		if !rm.MatchString("hello世界") {
			t.Error("expected to match unicode text")
		}
		if !rm.MatchString("say hello世界 there") {
			t.Error("expected to match as substring")
		}
	})

	t.Run("unicode in prefix", func(t *testing.T) {
		rm := MustRegexpMatcher("^hello世界")
		if !rm.MatchString("hello世界") {
			t.Error("expected to match unicode prefix")
		}
		if !rm.MatchString("hello世界 test") {
			t.Error("expected to match with suffix")
		}
		if rm.MatchString("test hello世界") {
			t.Error("expected not to match with prefix")
		}
	})

	t.Run("unicode in suffix", func(t *testing.T) {
		rm := MustRegexpMatcher("世界$")
		if !rm.MatchString("世界") {
			t.Error("expected to match unicode suffix")
		}
		if !rm.MatchString("hello 世界") {
			t.Error("expected to match with prefix")
		}
		if rm.MatchString("世界 test") {
			t.Error("expected not to match with suffix")
		}
	})

	t.Run("pattern with newline character", func(t *testing.T) {
		// Newline is not a meta character, so it triggers regexp
		rm := MustRegexpMatcher("hello\nworld")
		if rm.regexp == nil {
			t.Error("expected regexp to be set for pattern with newline")
		}
	})

	t.Run("pattern with tab character", func(t *testing.T) {
		// Tab is not a meta character, so it triggers regexp
		rm := MustRegexpMatcher("hello\tworld")
		if rm.regexp == nil {
			t.Error("expected regexp to be set for pattern with tab")
		}
	})
}

func TestRegexpMatcher_EqualityChecks(t *testing.T) {
	t.Run("equal prefix patterns", func(t *testing.T) {
		rm1 := MustRegexpMatcher("^hello")
		rm2 := MustRegexpMatcher("^hello")
		if !rm1.Equal(*rm2) {
			t.Error("expected equal prefix patterns to be equal")
		}
	})

	t.Run("different prefix patterns", func(t *testing.T) {
		rm1 := MustRegexpMatcher("^hello")
		rm2 := MustRegexpMatcher("^world")
		if rm1.Equal(*rm2) {
			t.Error("expected different prefix patterns to not be equal")
		}
	})

	t.Run("equal suffix patterns", func(t *testing.T) {
		rm1 := MustRegexpMatcher("hello$")
		rm2 := MustRegexpMatcher("hello$")
		if !rm1.Equal(*rm2) {
			t.Error("expected equal suffix patterns to be equal")
		}
	})

	t.Run("equal exact patterns", func(t *testing.T) {
		rm1 := MustRegexpMatcher("^hello$")
		rm2 := MustRegexpMatcher("^hello$")
		if !rm1.Equal(*rm2) {
			t.Error("expected equal exact patterns to be equal")
		}
	})

	t.Run("prefix vs suffix - different ops", func(t *testing.T) {
		rm1 := MustRegexpMatcher("^hello")
		rm2 := MustRegexpMatcher("hello$")
		// Both have text="hello" but different ops, Equal checks op field
		if rm1.Equal(*rm2) {
			t.Error("expected prefix and suffix to not be equal (different ops)")
		}
	})

	t.Run("prefix vs contain - different ops", func(t *testing.T) {
		rm1 := MustRegexpMatcher("^hello")
		rm2 := MustRegexpMatcher("hello")
		// Both have text="hello" but different ops, Equal checks op field
		if rm1.Equal(*rm2) {
			t.Error("expected prefix and contain to not be equal (different ops)")
		}
	})

	t.Run("exact vs contain - different ops", func(t *testing.T) {
		rm1 := MustRegexpMatcher("^hello$")
		rm2 := MustRegexpMatcher("hello")
		// Both have text="hello" but different ops, Equal checks op field
		if rm1.Equal(*rm2) {
			t.Error("expected exact and contain to not be equal (different ops)")
		}
	})

	t.Run("different text values", func(t *testing.T) {
		rm1 := MustRegexpMatcher("^hello")
		rm2 := MustRegexpMatcher("^world")
		if rm1.Equal(*rm2) {
			t.Error("expected different text values to not be equal")
		}
	})
}

func TestRegexpMatcher_SpecialCases(t *testing.T) {
	t.Run("caret in middle of text", func(t *testing.T) {
		rm := MustRegexpMatcher("hello^world")
		// "^" in middle is regex syntax
		if rm.regexp == nil {
			t.Error("expected regexp to be set")
		}
	})

	t.Run("dollar in middle of text", func(t *testing.T) {
		rm := MustRegexpMatcher("hello$world")
		// "$" in middle is regex syntax
		if rm.regexp == nil {
			t.Error("expected regexp to be set")
		}
	})

	t.Run("multiple carets", func(t *testing.T) {
		rm := MustRegexpMatcher("^^hello")
		// Multiple "^" is regex syntax
		if rm.regexp == nil {
			t.Error("expected regexp to be set")
		}
	})

	t.Run("multiple dollars", func(t *testing.T) {
		rm := MustRegexpMatcher("hello$$")
		// Multiple "$" is regex syntax
		if rm.regexp == nil {
			t.Error("expected regexp to be set")
		}
	})

	t.Run("escaped characters in pattern", func(t *testing.T) {
		rm := MustRegexpMatcher(`hello\.world`)
		// Backslash is regex syntax
		if rm.regexp == nil {
			t.Error("expected regexp to be set")
		}
	})

	t.Run("parentheses in pattern", func(t *testing.T) {
		rm := MustRegexpMatcher("hello(world)")
		// Parentheses are regex syntax
		if rm.regexp == nil {
			t.Error("expected regexp to be set")
		}
	})

	t.Run("pipe in pattern", func(t *testing.T) {
		rm := MustRegexpMatcher("hello|world")
		// Pipe is regex syntax
		if rm.regexp == nil {
			t.Error("expected regexp to be set")
		}
	})

	t.Run("plus in pattern", func(t *testing.T) {
		rm := MustRegexpMatcher("hello+")
		// Plus is regex syntax
		if rm.regexp == nil {
			t.Error("expected regexp to be set")
		}
	})

	t.Run("asterisk in pattern", func(t *testing.T) {
		rm := MustRegexpMatcher("hello*")
		// Asterisk is regex syntax
		if rm.regexp == nil {
			t.Error("expected regexp to be set")
		}
	})

	t.Run("question mark in pattern", func(t *testing.T) {
		rm := MustRegexpMatcher("hello?")
		// Question mark is regex syntax
		if rm.regexp == nil {
			t.Error("expected regexp to be set")
		}
	})

	t.Run("curly braces in pattern", func(t *testing.T) {
		rm := MustRegexpMatcher("hello{2}")
		// Curly braces are regex syntax
		if rm.regexp == nil {
			t.Error("expected regexp to be set")
		}
	})
}
