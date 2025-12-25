package goutils

import (
	"encoding/json"
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

	gotList := aos.List()
	if len(gotList) != len(list) {
		t.Errorf("expected list length %d, got: %d", len(list), len(gotList))
	}

	for i, v := range list {
		if gotList[i] != v {
			t.Errorf("expected list[%d] to be %s, got: %s", i, v, gotList[i])
		}
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
