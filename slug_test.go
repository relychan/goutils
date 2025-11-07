package goutils

import (
	"encoding/json"
	"strings"
	"testing"

	"go.yaml.in/yaml/v4"
)

func TestSlug(t *testing.T) {
	t.Run("unmarshal_json", func(t *testing.T) {
		var slug Slug

		err := json.Unmarshal([]byte(`"hello-world"`), &slug)
		if err != nil {
			t.Fatalf("expected nil error, got: %s", err)
		}

		if slug.String() != "hello-world" {
			t.Fatalf("slug not equal, expected hello-world, got: %s", slug)
		}

		err = json.Unmarshal([]byte(`"hello world"`), &slug)
		if !strings.Contains(err.Error(), "invalid slug") {
			t.Fatalf("expected invalid slug error, got: %s", err)
		}

		err = json.Unmarshal([]byte(`1`), &slug)
		if !strings.Contains(err.Error(), "cannot unmarshal number into Go value of type string") {
			t.Fatalf("expected unmarshal error, got: %s", err)
		}
	})

	t.Run("unmarshal_yaml", func(t *testing.T) {
		var slug Slug

		err := yaml.Unmarshal([]byte(`"hello-world"`), &slug)
		if err != nil {
			t.Fatalf("expected nil error, got: %s", err)
		}

		if slug.String() != "hello-world" {
			t.Fatalf("slug not equal, expected hello-world, got: %s", slug)
		}

		err = yaml.Unmarshal([]byte(`"hello world"`), &slug)
		if !strings.Contains(err.Error(), "invalid slug") {
			t.Fatalf("expected invalid slug error, got: %s", err)
		}

		err = yaml.Unmarshal([]byte(`{}`), &slug)
		if !strings.Contains(err.Error(), "cannot unmarshal !!map into string") {
			t.Fatalf("expected unmarshal error, got: %s", err)
		}
	})
}
