package goutils

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"

	"go.yaml.in/yaml/v4"
)

// ============================================================================
// AllOrListString Tests
// ============================================================================

func TestAllOrListString_NewAll(t *testing.T) {
	aos := NewAll()

	if !aos.IsAll() {
		t.Error("expected IsAll() to be true")
	}

	if aos.List() != nil {
		t.Error("expected List() to be nil")
	}

	if aos.IsZero() {
		t.Error("expected IsZero() to be false")
	}

	if !aos.Contains("bar") {
		t.Error("expected Contains() to be true")
	}

	if aos.String() != "*" {
		t.Errorf("expected String() to be '*', got: %s", aos.String())
	}
}

func TestAllOrListString_NewStringList(t *testing.T) {
	list := []string{"foo", "bar", "baz"}
	aos := NewStringList(list)

	if aos.IsAll() {
		t.Error("expected IsAll() to be false")
	}

	if aos.IsZero() {
		t.Error("expected IsZero() to be false")
	}

	if !aos.Contains("bar") {
		t.Error("expected Contains() to be true")
	}

	gotList := aos.List()
	if len(gotList) != len(list) {
		t.Errorf("expected list length %d, got: %d", len(list), len(gotList))
	}

	for i, v := range list {
		if gotList[i] != v {
			t.Errorf("expected list[%d] to be %s, got: %s", i, v, gotList[i])
		}
	}

	aos2 := aos.Map(func(s string, _ int) string {
		return s + "1"
	})

	if !slices.Equal(aos2.List(), []string{"foo1", "bar1", "baz1"}) {
		t.Error("expected equal")
	}

	aos3 := NewAll().Map(func(s string, _ int) string {
		return s + "1"
	})
	if !aos3.IsAll() {
		t.Error("expected IsAll() to be true")
	}
}

func TestAllOrListString_IsZero(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		var aos AllOrListString
		if !aos.IsZero() {
			t.Error("expected IsZero() to be true for zero value")
		}
	})

	t.Run("not zero with all", func(t *testing.T) {
		aos := NewAll()
		if aos.IsZero() {
			t.Error("expected IsZero() to be false when all is true")
		}
	})

	t.Run("not zero with list", func(t *testing.T) {
		aos := NewStringList([]string{"foo"})
		if aos.IsZero() {
			t.Error("expected IsZero() to be false when list is not empty")
		}
	})

	t.Run("zero with empty list", func(t *testing.T) {
		aos := NewStringList([]string{})
		if !aos.IsZero() {
			t.Error("expected IsZero() to be true when list is empty")
		}
	})

	t.Run("zero with nil list", func(t *testing.T) {
		aos := NewStringList(nil)
		if !aos.IsZero() {
			t.Error("expected IsZero() to be true when list is nil")
		}
	})
}

func TestAllOrListString_Equal(t *testing.T) {
	t.Run("both all", func(t *testing.T) {
		aos1 := NewAll()
		aos2 := NewAll()
		if !aos1.Equal(aos2) {
			t.Error("expected two 'all' instances to be equal")
		}
	})

	t.Run("both empty", func(t *testing.T) {
		aos1 := NewStringList([]string{})
		aos2 := NewStringList(nil)
		if !aos1.Equal(aos2) {
			t.Error("expected empty lists to be equal")
		}
	})

	t.Run("same list", func(t *testing.T) {
		aos1 := NewStringList([]string{"foo", "bar"})
		aos2 := NewStringList([]string{"foo", "bar"})
		if !aos1.Equal(aos2) {
			t.Error("expected same lists to be equal")
		}
	})

	t.Run("same list different order", func(t *testing.T) {
		aos1 := NewStringList([]string{"foo", "bar"})
		aos2 := NewStringList([]string{"bar", "foo"})
		if !aos1.Equal(aos2) {
			t.Error("expected lists with same elements in different order to be equal")
		}
	})

	t.Run("different lists", func(t *testing.T) {
		aos1 := NewStringList([]string{"foo", "bar"})
		aos2 := NewStringList([]string{"foo", "baz"})
		if aos1.Equal(aos2) {
			t.Error("expected different lists to not be equal")
		}
	})

	t.Run("all vs list", func(t *testing.T) {
		aos1 := NewAll()
		aos2 := NewStringList([]string{"foo"})
		if aos1.Equal(aos2) {
			t.Error("expected 'all' and list to not be equal")
		}
	})

	t.Run("different lengths", func(t *testing.T) {
		aos1 := NewStringList([]string{"foo"})
		aos2 := NewStringList([]string{"foo", "bar"})
		if aos1.Equal(aos2) {
			t.Error("expected lists of different lengths to not be equal")
		}
	})
}

func TestAllOrListString_JSON(t *testing.T) {
	t.Run("unmarshal all", func(t *testing.T) {
		var aos AllOrListString
		err := json.Unmarshal([]byte(`"*"`), &aos)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !aos.IsAll() {
			t.Error("expected IsAll() to be true")
		}

		if aos.List() != nil {
			t.Error("expected List() to be nil")
		}
	})

	t.Run("marshal all", func(t *testing.T) {
		aos := NewAll()
		data, err := json.Marshal(aos)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if string(data) != `"*"` {
			t.Errorf("expected '*', got: %s", string(data))
		}
	})

	t.Run("unmarshal list", func(t *testing.T) {
		var aos AllOrListString
		err := json.Unmarshal([]byte(`["foo","bar"]`), &aos)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if aos.IsAll() {
			t.Error("expected IsAll() to be false")
		}

		list := aos.List()
		if len(list) != 2 {
			t.Errorf("expected list length 2, got: %d", len(list))
		}

		if list[0] != "foo" || list[1] != "bar" {
			t.Errorf("unexpected list values: %v", list)
		}
	})

	t.Run("marshal list", func(t *testing.T) {
		aos := NewStringList([]string{"foo", "bar"})
		data, err := json.Marshal(aos)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if string(data) != `["foo","bar"]` {
			t.Errorf("expected [\"foo\",\"bar\"], got: %s", string(data))
		}
	})

	t.Run("unmarshal empty list", func(t *testing.T) {
		var aos AllOrListString
		err := json.Unmarshal([]byte(`[]`), &aos)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if aos.IsAll() {
			t.Error("expected IsAll() to be false")
		}

		if !aos.IsZero() {
			t.Error("expected IsZero() to be true")
		}
	})

	t.Run("unmarshal invalid", func(t *testing.T) {
		var aos AllOrListString
		err := json.Unmarshal([]byte(`{}`), &aos)
		if err == nil {
			t.Error("expected error for invalid JSON")
		}
	})

	t.Run("round trip all", func(t *testing.T) {
		original := NewAll()
		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var decoded AllOrListString
		err = json.Unmarshal(data, &decoded)
		if err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}

		if !original.Equal(decoded) {
			t.Error("round trip failed for 'all'")
		}
	})

	t.Run("round trip list", func(t *testing.T) {
		original := NewStringList([]string{"foo", "bar", "baz"})
		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var decoded AllOrListString
		err = json.Unmarshal(data, &decoded)
		if err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}

		if !original.Equal(decoded) {
			t.Error("round trip failed for list")
		}
	})
}

func TestAllOrListString_YAML(t *testing.T) {
	t.Run("unmarshal all", func(t *testing.T) {
		var aos AllOrListString
		err := yaml.Unmarshal([]byte(`"*"`), &aos)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !aos.IsAll() {
			t.Error("expected IsAll() to be true")
		}

		if aos.List() != nil {
			t.Error("expected List() to be nil")
		}
	})

	t.Run("marshal all", func(t *testing.T) {
		aos := NewAll()
		data, err := yaml.Marshal(aos)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if strings.TrimSpace(string(data)) != `'*'` {
			t.Errorf("expected '*', got: %s", string(data))
		}
	})

	t.Run("unmarshal list", func(t *testing.T) {
		var aos AllOrListString
		err := yaml.Unmarshal([]byte("- foo\n- bar"), &aos)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if aos.IsAll() {
			t.Error("expected IsAll() to be false")
		}

		list := aos.List()
		if len(list) != 2 {
			t.Errorf("expected list length 2, got: %d", len(list))
		}

		if list[0] != "foo" || list[1] != "bar" {
			t.Errorf("unexpected list values: %v", list)
		}
	})

	t.Run("marshal list", func(t *testing.T) {
		aos := NewStringList([]string{"foo", "bar"})
		data, err := yaml.Marshal(aos)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := "- foo\n- bar"
		if strings.TrimSpace(string(data)) != expected {
			t.Errorf("expected %s, got: %s", expected, string(data))
		}
	})

	t.Run("round trip all", func(t *testing.T) {
		original := NewAll()
		data, err := yaml.Marshal(original)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var decoded AllOrListString
		err = yaml.Unmarshal(data, &decoded)
		if err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}

		if !original.Equal(decoded) {
			t.Error("round trip failed for 'all'")
		}
	})

	t.Run("round trip list", func(t *testing.T) {
		original := NewStringList([]string{"foo", "bar", "baz"})
		data, err := yaml.Marshal(original)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var decoded AllOrListString
		err = yaml.Unmarshal(data, &decoded)
		if err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}

		if !original.Equal(decoded) {
			t.Error("round trip failed for list")
		}
	})
}

func TestAllOrListString_String(t *testing.T) {
	t.Run("all", func(t *testing.T) {
		aos := NewAll()
		if aos.String() != "*" {
			t.Errorf("expected '*', got: %s", aos.String())
		}
	})

	t.Run("empty list", func(t *testing.T) {
		aos := NewStringList([]string{})
		if aos.String() != "[]" {
			t.Errorf("expected '[]', got: %s", aos.String())
		}
	})

	t.Run("single item", func(t *testing.T) {
		aos := NewStringList([]string{"foo"})
		if aos.String() != "[foo]" {
			t.Errorf("expected '[foo]', got: %s", aos.String())
		}
	})

	t.Run("multiple items", func(t *testing.T) {
		aos := NewStringList([]string{"foo", "bar"})
		if aos.String() != "[foo bar]" {
			t.Errorf("expected '[foo bar]', got: %s", aos.String())
		}
	})
}

// ============================================================================
// Wildcard Tests
// ============================================================================

func TestWildcard_NewWildcard(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		w, ok := NewWildcard("")
		if ok {
			t.Error("expected ok to be false for empty string")
		}
		if !w.IsZero() {
			t.Error("expected wildcard to be zero")
		}
	})

	t.Run("no wildcard", func(t *testing.T) {
		w, ok := NewWildcard("test")
		if ok {
			t.Error("expected ok to be false for string without wildcard")
		}
		if w.prefix != "test" {
			t.Errorf("expected prefix 'test', got: %s", w.prefix)
		}
	})

	t.Run("prefix wildcard", func(t *testing.T) {
		w, ok := NewWildcard("test*")
		if !ok {
			t.Error("expected ok to be true")
		}
		if w.prefix != "test" {
			t.Errorf("expected prefix 'test', got: %s", w.prefix)
		}
		if w.suffix != "" {
			t.Errorf("expected empty suffix, got: %s", w.suffix)
		}
	})

	t.Run("suffix wildcard", func(t *testing.T) {
		w, ok := NewWildcard("*test")
		if !ok {
			t.Error("expected ok to be true")
		}
		if w.prefix != "" {
			t.Errorf("expected empty prefix, got: %s", w.prefix)
		}
		if w.suffix != "test" {
			t.Errorf("expected suffix 'test', got: %s", w.suffix)
		}
	})

	t.Run("middle wildcard", func(t *testing.T) {
		w, ok := NewWildcard("pre*suf")
		if !ok {
			t.Error("expected ok to be true")
		}
		if w.prefix != "pre" {
			t.Errorf("expected prefix 'pre', got: %s", w.prefix)
		}
		if w.suffix != "suf" {
			t.Errorf("expected suffix 'suf', got: %s", w.suffix)
		}
	})

	t.Run("only wildcard", func(t *testing.T) {
		w, ok := NewWildcard("*")
		if !ok {
			t.Error("expected ok to be true")
		}
		if w.prefix != "" {
			t.Errorf("expected empty prefix, got: %s", w.prefix)
		}
		if w.suffix != "" {
			t.Errorf("expected empty suffix, got: %s", w.suffix)
		}
	})

	t.Run("multiple wildcards", func(t *testing.T) {
		w, ok := NewWildcard("pre**suf")
		if !ok {
			t.Error("expected ok to be true")
		}
		if w.prefix != "pre" {
			t.Errorf("expected prefix 'pre', got: %s", w.prefix)
		}
		// Multiple wildcards should be treated as one, with suffix trimming leading *
		if w.suffix != "suf" {
			t.Errorf("expected suffix 'suf', got: %s", w.suffix)
		}
	})
}

func TestWildcard_IsZero(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		var w Wildcard
		if !w.IsZero() {
			t.Error("expected IsZero() to be true for zero value")
		}
	})

	t.Run("with prefix", func(t *testing.T) {
		w, _ := NewWildcard("test*")
		if w.IsZero() {
			t.Error("expected IsZero() to be false with prefix")
		}
	})

	t.Run("with suffix", func(t *testing.T) {
		w, _ := NewWildcard("*test")
		if w.IsZero() {
			t.Error("expected IsZero() to be false with suffix")
		}
	})

	t.Run("with both", func(t *testing.T) {
		w, _ := NewWildcard("pre*suf")
		if w.IsZero() {
			t.Error("expected IsZero() to be false with both prefix and suffix")
		}
	})
}

func TestWildcard_Equal(t *testing.T) {
	t.Run("both zero", func(t *testing.T) {
		var w1, w2 Wildcard
		if !w1.Equal(w2) {
			t.Error("expected zero wildcards to be equal")
		}
	})

	t.Run("same prefix wildcard", func(t *testing.T) {
		w1, _ := NewWildcard("test*")
		w2, _ := NewWildcard("test*")
		if !w1.Equal(w2) {
			t.Error("expected same prefix wildcards to be equal")
		}
	})

	t.Run("same suffix wildcard", func(t *testing.T) {
		w1, _ := NewWildcard("*test")
		w2, _ := NewWildcard("*test")
		if !w1.Equal(w2) {
			t.Error("expected same suffix wildcards to be equal")
		}
	})

	t.Run("same middle wildcard", func(t *testing.T) {
		w1, _ := NewWildcard("pre*suf")
		w2, _ := NewWildcard("pre*suf")
		if !w1.Equal(w2) {
			t.Error("expected same middle wildcards to be equal")
		}
	})

	t.Run("different prefix", func(t *testing.T) {
		w1, _ := NewWildcard("test1*")
		w2, _ := NewWildcard("test2*")
		if w1.Equal(w2) {
			t.Error("expected different prefix wildcards to not be equal")
		}
	})

	t.Run("different suffix", func(t *testing.T) {
		w1, _ := NewWildcard("*test1")
		w2, _ := NewWildcard("*test2")
		if w1.Equal(w2) {
			t.Error("expected different suffix wildcards to not be equal")
		}
	})
}

func TestWildcard_Match(t *testing.T) {
	t.Run("prefix wildcard matches", func(t *testing.T) {
		w, _ := NewWildcard("test*")
		tests := []string{"test", "testing", "test123", "testABC"}
		for _, s := range tests {
			if !w.Match(s) {
				t.Errorf("expected '%s' to match 'test*'", s)
			}
		}
	})

	t.Run("prefix wildcard no match", func(t *testing.T) {
		w, _ := NewWildcard("test*")
		tests := []string{"tes", "atest", "TEST", ""}
		for _, s := range tests {
			if w.Match(s) {
				t.Errorf("expected '%s' to not match 'test*'", s)
			}
		}
	})

	t.Run("suffix wildcard matches", func(t *testing.T) {
		w, _ := NewWildcard("*test")
		tests := []string{"test", "mytest", "123test", "ABCtest"}
		for _, s := range tests {
			if !w.Match(s) {
				t.Errorf("expected '%s' to match '*test'", s)
			}
		}
	})

	t.Run("suffix wildcard no match", func(t *testing.T) {
		w, _ := NewWildcard("*test")
		tests := []string{"est", "testa", "TEST", ""}
		for _, s := range tests {
			if w.Match(s) {
				t.Errorf("expected '%s' to not match '*test'", s)
			}
		}
	})

	t.Run("middle wildcard matches", func(t *testing.T) {
		w, _ := NewWildcard("pre*suf")
		tests := []string{"presuf", "pre123suf", "preANYTHINGsuf"}
		for _, s := range tests {
			if !w.Match(s) {
				t.Errorf("expected '%s' to match 'pre*suf'", s)
			}
		}
	})

	t.Run("middle wildcard no match", func(t *testing.T) {
		w, _ := NewWildcard("pre*suf")
		tests := []string{"pre", "suf", "presuff", "ppresuf", ""}
		for _, s := range tests {
			if w.Match(s) {
				t.Errorf("expected '%s' to not match 'pre*suf'", s)
			}
		}
	})

	t.Run("wildcard only matches all", func(t *testing.T) {
		w, _ := NewWildcard("*")
		tests := []string{"", "a", "test", "anything goes"}
		for _, s := range tests {
			if !w.Match(s) {
				t.Errorf("expected '%s' to match '*'", s)
			}
		}
	})

	t.Run("minimum length check", func(t *testing.T) {
		w, _ := NewWildcard("abc*xyz")
		// Minimum length is len("abc") + len("xyz") = 6
		if w.Match("abxyz") { // length 5
			t.Error("expected string shorter than minimum to not match")
		}
		if !w.Match("abcxyz") { // length 6
			t.Error("expected string at minimum length to match")
		}
		if !w.Match("abc123xyz") { // length 9
			t.Error("expected string longer than minimum to match")
		}
	})

	t.Run("empty prefix and suffix", func(t *testing.T) {
		w, _ := NewWildcard("*")
		if !w.Match("anything") {
			t.Error("expected wildcard with empty prefix and suffix to match anything")
		}
	})
}

func TestWildcard_String(t *testing.T) {
	t.Run("prefix wildcard", func(t *testing.T) {
		w, _ := NewWildcard("test*")
		if w.String() != "test*" {
			t.Errorf("expected 'test*', got: %s", w.String())
		}
	})

	t.Run("suffix wildcard", func(t *testing.T) {
		w, _ := NewWildcard("*test")
		if w.String() != "*test" {
			t.Errorf("expected '*test', got: %s", w.String())
		}
	})

	t.Run("middle wildcard", func(t *testing.T) {
		w, _ := NewWildcard("pre*suf")
		if w.String() != "pre*suf" {
			t.Errorf("expected 'pre*suf', got: %s", w.String())
		}
	})

	t.Run("only wildcard", func(t *testing.T) {
		w, _ := NewWildcard("*")
		if w.String() != "*" {
			t.Errorf("expected '*', got: %s", w.String())
		}
	})

	t.Run("zero value", func(t *testing.T) {
		var w Wildcard
		if w.String() != "*" {
			t.Errorf("expected '*', got: %s", w.String())
		}
	})
}

// ============================================================================
// AllOrListWildcardString Tests
// ============================================================================

func TestAllOrListWildcardString_ParseStrings(t *testing.T) {
	t.Run("constructors", func(t *testing.T) {
		aos := NewAllWildcard()
		if !aos.IsAll() {
			t.Error("expected IsAll() to be true")
		}

		aos = NewAllOrListWildcardStringFromStrings([]string{"foo", "bar", "baz"})
		if len(aos.List()) != 3 {
			t.Errorf("expected 3 static strings, got: %d", len(aos.List()))
		}

		if len(aos.Wildcards()) != 0 {
			t.Errorf("expected 0 wildcards, got: %d", len(aos.Wildcards()))
		}

		if !aos.Contains("bar") {
			t.Error("expected Contains() to be true")
		}
	})

	t.Run("only static strings", func(t *testing.T) {
		var aos AllOrListWildcardString
		aos.parseStrings([]string{"foo", "bar", "baz"})

		if aos.IsAll() {
			t.Error("expected IsAll() to be false")
		}

		if len(aos.List()) != 3 {
			t.Errorf("expected 3 static strings, got: %d", len(aos.List()))
		}

		if len(aos.Wildcards()) != 0 {
			t.Errorf("expected 0 wildcards, got: %d", len(aos.Wildcards()))
		}

		if !aos.Contains("bar") {
			t.Error("expected Contains() to be true")
		}
	})

	t.Run("only wildcards", func(t *testing.T) {
		var aos AllOrListWildcardString
		aos.parseStrings([]string{"test*", "*value", "pre*suf"})

		if aos.IsAll() {
			t.Error("expected IsAll() to be false")
		}

		if len(aos.List()) != 0 {
			t.Errorf("expected 0 static strings, got: %d", len(aos.List()))
		}

		if len(aos.Wildcards()) != 3 {
			t.Errorf("expected 3 wildcards, got: %d", len(aos.Wildcards()))
		}
	})

	t.Run("mixed static and wildcards", func(t *testing.T) {
		var aos AllOrListWildcardString
		aos.parseStrings([]string{"static", "test*", "another", "*value"})

		if aos.IsAll() {
			t.Error("expected IsAll() to be false")
		}

		if len(aos.List()) != 2 {
			t.Errorf("expected 2 static strings, got: %d", len(aos.List()))
		}

		if len(aos.Wildcards()) != 2 {
			t.Errorf("expected 2 wildcards, got: %d", len(aos.Wildcards()))
		}

		if !aos.Contains("testvalue") {
			t.Error("expected Contains() to be true")
		}

		if aos.Contains("not contains") {
			t.Error("expected Contains() to be false")
		}
	})

	t.Run("star becomes all", func(t *testing.T) {
		var aos AllOrListWildcardString
		aos.parseStrings([]string{"foo", "*", "bar"})

		if !aos.IsAll() {
			t.Error("expected IsAll() to be true when '*' is present")
		}

		if aos.List() != nil {
			t.Error("expected List() to be nil")
		}

		if aos.Wildcards() != nil {
			t.Error("expected Wildcards() to be nil")
		}
	})

	t.Run("empty string", func(t *testing.T) {
		var aos AllOrListWildcardString
		aos.parseStrings([]string{""})

		if aos.IsAll() {
			t.Error("expected IsAll() to be false")
		}

		list := aos.List()
		if len(list) != 1 || list[0] != "" {
			t.Errorf("expected empty string in list, got: %v", list)
		}
	})

	t.Run("empty list", func(t *testing.T) {
		var aos AllOrListWildcardString
		aos.parseStrings([]string{})

		if aos.IsAll() {
			t.Error("expected IsAll() to be false")
		}

		if !aos.IsZero() {
			t.Error("expected IsZero() to be true")
		}
	})

	t.Run("nil list", func(t *testing.T) {
		var aos AllOrListWildcardString
		aos.parseStrings(nil)

		if aos.IsAll() {
			t.Error("expected IsAll() to be false")
		}

		if !aos.IsZero() {
			t.Error("expected IsZero() to be true")
		}
	})
}

func TestAllOrListWildcardString_IsZero(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		var aos AllOrListWildcardString
		if !aos.IsZero() {
			t.Error("expected IsZero() to be true for zero value")
		}
	})

	t.Run("not zero with all", func(t *testing.T) {
		var aos AllOrListWildcardString
		aos.all = true
		if aos.IsZero() {
			t.Error("expected IsZero() to be false when all is true")
		}
	})

	t.Run("not zero with static list", func(t *testing.T) {
		var aos AllOrListWildcardString
		aos.parseStrings([]string{"foo"})
		if aos.IsZero() {
			t.Error("expected IsZero() to be false when list is not empty")
		}
	})

	t.Run("not zero with wildcards", func(t *testing.T) {
		var aos AllOrListWildcardString
		aos.parseStrings([]string{"test*"})
		if aos.IsZero() {
			t.Error("expected IsZero() to be false when wildcards exist")
		}
	})
}

func TestAllOrListWildcardString_Equal(t *testing.T) {
	t.Run("both all", func(t *testing.T) {
		var aos1, aos2 AllOrListWildcardString
		aos1.all = true
		aos2.all = true
		if !aos1.Equal(aos2) {
			t.Error("expected two 'all' instances to be equal")
		}
	})

	t.Run("both empty", func(t *testing.T) {
		var aos1, aos2 AllOrListWildcardString
		if !aos1.Equal(aos2) {
			t.Error("expected empty instances to be equal")
		}
	})

	t.Run("same static lists", func(t *testing.T) {
		var aos1, aos2 AllOrListWildcardString
		aos1.parseStrings([]string{"foo", "bar"})
		aos2.parseStrings([]string{"foo", "bar"})
		if !aos1.Equal(aos2) {
			t.Error("expected same static lists to be equal")
		}
	})

	t.Run("same wildcards", func(t *testing.T) {
		var aos1, aos2 AllOrListWildcardString
		aos1.parseStrings([]string{"test*", "*value"})
		aos2.parseStrings([]string{"test*", "*value"})
		if !aos1.Equal(aos2) {
			t.Error("expected same wildcards to be equal")
		}
	})

	t.Run("same mixed", func(t *testing.T) {
		var aos1, aos2 AllOrListWildcardString
		aos1.parseStrings([]string{"static", "test*", "another"})
		aos2.parseStrings([]string{"static", "test*", "another"})
		if !aos1.Equal(aos2) {
			t.Error("expected same mixed lists to be equal")
		}
	})

	t.Run("different static lists", func(t *testing.T) {
		var aos1, aos2 AllOrListWildcardString
		aos1.parseStrings([]string{"foo"})
		aos2.parseStrings([]string{"bar"})
		if aos1.Equal(aos2) {
			t.Error("expected different static lists to not be equal")
		}
	})

	t.Run("different wildcards", func(t *testing.T) {
		var aos1, aos2 AllOrListWildcardString
		aos1.parseStrings([]string{"test*"})
		aos2.parseStrings([]string{"*test"})
		if aos1.Equal(aos2) {
			t.Error("expected different wildcards to not be equal")
		}
	})

	t.Run("all vs list", func(t *testing.T) {
		var aos1, aos2 AllOrListWildcardString
		aos1.all = true
		aos2.parseStrings([]string{"foo"})
		if aos1.Equal(aos2) {
			t.Error("expected 'all' and list to not be equal")
		}
	})
}

func TestAllOrListWildcardString_JSON(t *testing.T) {
	t.Run("unmarshal all", func(t *testing.T) {
		var aos AllOrListWildcardString
		err := json.Unmarshal([]byte(`"*"`), &aos)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !aos.IsAll() {
			t.Error("expected IsAll() to be true")
		}
	})

	t.Run("marshal all", func(t *testing.T) {
		var aos AllOrListWildcardString
		aos.all = true
		data, err := json.Marshal(aos)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if string(data) != `"*"` {
			t.Errorf("expected '*', got: %s", string(data))
		}
	})

	t.Run("unmarshal static list", func(t *testing.T) {
		var aos AllOrListWildcardString
		err := json.Unmarshal([]byte(`["foo","bar"]`), &aos)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if aos.IsAll() {
			t.Error("expected IsAll() to be false")
		}

		list := aos.List()
		if len(list) != 2 {
			t.Errorf("expected 2 items, got: %d", len(list))
		}
	})

	t.Run("marshal static list", func(t *testing.T) {
		var aos AllOrListWildcardString
		aos.parseStrings([]string{"foo", "bar"})
		data, err := json.Marshal(aos)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if string(data) != `["foo","bar"]` {
			t.Errorf("expected [\"foo\",\"bar\"], got: %s", string(data))
		}
	})

	t.Run("unmarshal wildcards", func(t *testing.T) {
		var aos AllOrListWildcardString
		err := json.Unmarshal([]byte(`["test*","*value"]`), &aos)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if aos.IsAll() {
			t.Error("expected IsAll() to be false")
		}

		wildcards := aos.Wildcards()
		if len(wildcards) != 2 {
			t.Errorf("expected 2 wildcards, got: %d", len(wildcards))
		}
	})

	t.Run("marshal wildcards", func(t *testing.T) {
		var aos AllOrListWildcardString
		aos.parseStrings([]string{"test*", "*value"})
		data, err := json.Marshal(aos)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if string(data) != `["test*","*value"]` {
			t.Errorf("expected [\"test*\",\"*value\"], got: %s", string(data))
		}
	})

	t.Run("unmarshal mixed", func(t *testing.T) {
		var aos AllOrListWildcardString
		err := json.Unmarshal([]byte(`["static","test*","another"]`), &aos)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		list := aos.List()
		if len(list) != 2 {
			t.Errorf("expected 2 static items, got: %d", len(list))
		}

		wildcards := aos.Wildcards()
		if len(wildcards) != 1 {
			t.Errorf("expected 1 wildcard, got: %d", len(wildcards))
		}
	})

	t.Run("marshal mixed", func(t *testing.T) {
		var aos AllOrListWildcardString
		aos.parseStrings([]string{"static", "test*", "another"})
		data, err := json.Marshal(aos)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Order: static items first, then wildcards
		if string(data) != `["static","another","test*"]` {
			t.Errorf("expected [\"static\",\"another\",\"test*\"], got: %s", string(data))
		}
	})

	t.Run("round trip", func(t *testing.T) {
		var original AllOrListWildcardString
		original.parseStrings([]string{"foo", "test*", "bar", "*value"})

		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var decoded AllOrListWildcardString
		err = json.Unmarshal(data, &decoded)
		if err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}

		if !original.Equal(decoded) {
			t.Error("round trip failed")
		}
	})
}

func TestAllOrListWildcardString_YAML(t *testing.T) {
	t.Run("unmarshal all", func(t *testing.T) {
		var aos AllOrListWildcardString
		err := yaml.Unmarshal([]byte(`"*"`), &aos)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !aos.IsAll() {
			t.Error("expected IsAll() to be true")
		}
	})

	t.Run("marshal all", func(t *testing.T) {
		var aos AllOrListWildcardString
		aos.all = true
		data, err := yaml.Marshal(aos)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if strings.TrimSpace(string(data)) != `'*'` {
			t.Errorf("expected '*', got: %s", string(data))
		}
	})

	t.Run("unmarshal static list", func(t *testing.T) {
		var aos AllOrListWildcardString
		err := yaml.Unmarshal([]byte("- foo\n- bar"), &aos)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if aos.IsAll() {
			t.Error("expected IsAll() to be false")
		}

		list := aos.List()
		if len(list) != 2 {
			t.Errorf("expected 2 items, got: %d", len(list))
		}
	})

	t.Run("marshal static list", func(t *testing.T) {
		var aos AllOrListWildcardString
		aos.parseStrings([]string{"foo", "bar"})
		data, err := yaml.Marshal(aos)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := "- foo\n- bar"
		if strings.TrimSpace(string(data)) != expected {
			t.Errorf("expected %s, got: %s", expected, string(data))
		}
	})

	t.Run("unmarshal wildcards", func(t *testing.T) {
		var aos AllOrListWildcardString
		err := yaml.Unmarshal([]byte("- test*\n- '*value'"), &aos)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		wildcards := aos.Wildcards()
		if len(wildcards) != 2 {
			t.Errorf("expected 2 wildcards, got: %d", len(wildcards))
		}
	})

	t.Run("marshal wildcards", func(t *testing.T) {
		var aos AllOrListWildcardString
		aos.parseStrings([]string{"test*", "*value"})
		data, err := yaml.Marshal(aos)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should contain both wildcards
		dataStr := string(data)
		if !strings.Contains(dataStr, "test*") || !strings.Contains(dataStr, "*value") {
			t.Errorf("expected wildcards in output, got: %s", dataStr)
		}
	})

	t.Run("unmarshal star becomes all", func(t *testing.T) {
		var aos AllOrListWildcardString
		err := yaml.Unmarshal([]byte("- foo\n- '*'\n- bar"), &aos)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !aos.IsAll() {
			t.Error("expected IsAll() to be true when '*' is in list")
		}
	})

	t.Run("round trip", func(t *testing.T) {
		var original AllOrListWildcardString
		original.parseStrings([]string{"foo", "test*", "bar"})

		data, err := yaml.Marshal(original)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}

		var decoded AllOrListWildcardString
		err = yaml.Unmarshal(data, &decoded)
		if err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}

		if !original.Equal(decoded) {
			t.Error("round trip failed")
		}
	})
}

func TestAllOrListWildcardString_String(t *testing.T) {
	t.Run("all", func(t *testing.T) {
		var aos AllOrListWildcardString
		aos.all = true
		if aos.String() != "*" {
			t.Errorf("expected '*', got: %s", aos.String())
		}
	})

	t.Run("empty", func(t *testing.T) {
		var aos AllOrListWildcardString
		if aos.String() != "[]" {
			t.Errorf("expected '[]', got: %s", aos.String())
		}
	})

	t.Run("only static", func(t *testing.T) {
		var aos AllOrListWildcardString
		aos.parseStrings([]string{"foo", "bar"})
		expected := "[foo, bar]"
		if aos.String() != expected {
			t.Errorf("expected '%s', got: %s", expected, aos.String())
		}
	})

	t.Run("only wildcards", func(t *testing.T) {
		var aos AllOrListWildcardString
		aos.parseStrings([]string{"test*", "*value"})
		expected := "[, test*, *value]"
		if aos.String() != expected {
			t.Errorf("expected '%s', got: %s", expected, aos.String())
		}
	})

	t.Run("mixed", func(t *testing.T) {
		var aos AllOrListWildcardString
		aos.parseStrings([]string{"static", "test*", "another"})
		// Static items first, then wildcards
		expected := "[static, another, test*]"
		if aos.String() != expected {
			t.Errorf("expected '%s', got: %s", expected, aos.String())
		}
	})

	t.Run("single static", func(t *testing.T) {
		var aos AllOrListWildcardString
		aos.parseStrings([]string{"foo"})
		expected := "[foo]"
		if aos.String() != expected {
			t.Errorf("expected '%s', got: %s", expected, aos.String())
		}
	})

	t.Run("single wildcard", func(t *testing.T) {
		var aos AllOrListWildcardString
		aos.parseStrings([]string{"test*"})
		expected := "[, test*]"
		if aos.String() != expected {
			t.Errorf("expected '%s', got: %s", expected, aos.String())
		}
	})
}

func TestAllOrListWildcardString_ToStrings(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var aos AllOrListWildcardString
		result := aos.toStrings()
		if len(result) != 0 {
			t.Errorf("expected empty slice, got: %v", result)
		}
	})

	t.Run("only static", func(t *testing.T) {
		var aos AllOrListWildcardString
		aos.parseStrings([]string{"foo", "bar"})
		result := aos.toStrings()
		if len(result) != 2 {
			t.Errorf("expected 2 items, got: %d", len(result))
		}
		if result[0] != "foo" || result[1] != "bar" {
			t.Errorf("unexpected values: %v", result)
		}
	})

	t.Run("only wildcards", func(t *testing.T) {
		var aos AllOrListWildcardString
		aos.parseStrings([]string{"test*", "*value"})
		result := aos.toStrings()
		if len(result) != 2 {
			t.Errorf("expected 2 items, got: %d", len(result))
		}
		if result[0] != "test*" || result[1] != "*value" {
			t.Errorf("unexpected values: %v", result)
		}
	})

	t.Run("mixed", func(t *testing.T) {
		var aos AllOrListWildcardString
		aos.parseStrings([]string{"static", "test*", "another"})
		result := aos.toStrings()
		if len(result) != 3 {
			t.Errorf("expected 3 items, got: %d", len(result))
		}
		// Static items come first, then wildcards
		if result[0] != "static" || result[1] != "another" || result[2] != "test*" {
			t.Errorf("unexpected values: %v", result)
		}
	})
}

func TestAllOrListWildcardString_EdgeCases(t *testing.T) {
	t.Run("duplicate static strings", func(t *testing.T) {
		var aos AllOrListWildcardString
		aos.parseStrings([]string{"foo", "foo", "bar"})
		list := aos.List()
		// Duplicates are preserved
		if len(list) != 3 {
			t.Errorf("expected 3 items (with duplicate), got: %d", len(list))
		}
	})

	t.Run("duplicate wildcards", func(t *testing.T) {
		var aos AllOrListWildcardString
		aos.parseStrings([]string{"test*", "test*"})
		wildcards := aos.Wildcards()
		// Duplicates are preserved
		if len(wildcards) != 2 {
			t.Errorf("expected 2 wildcards (with duplicate), got: %d", len(wildcards))
		}
	})

	t.Run("empty string in list", func(t *testing.T) {
		var aos AllOrListWildcardString
		aos.parseStrings([]string{"", "foo"})
		list := aos.List()
		if len(list) != 2 {
			t.Errorf("expected 2 items, got: %d", len(list))
		}
		if list[0] != "" {
			t.Error("expected first item to be empty string")
		}
	})

	t.Run("multiple stars in pattern", func(t *testing.T) {
		var aos AllOrListWildcardString
		aos.parseStrings([]string{"pre**suf"})
		wildcards := aos.Wildcards()
		if len(wildcards) != 1 {
			t.Errorf("expected 1 wildcard, got: %d", len(wildcards))
		}
		// Multiple stars should be handled by NewWildcard
		if wildcards[0].prefix != "pre" || wildcards[0].suffix != "suf" {
			t.Errorf("unexpected wildcard: %v", wildcards[0])
		}
	})

	t.Run("star at beginning and end", func(t *testing.T) {
		var aos AllOrListWildcardString
		aos.parseStrings([]string{"*test*"})
		wildcards := aos.Wildcards()
		if len(wildcards) != 1 {
			t.Errorf("expected 1 wildcard, got: %d", len(wildcards))
		}
		if wildcards[0].prefix != "" || wildcards[0].suffix != "test*" {
			t.Errorf("unexpected wildcard: prefix='%s', suffix='%s'", wildcards[0].prefix, wildcards[0].suffix)
		}
	})
}
