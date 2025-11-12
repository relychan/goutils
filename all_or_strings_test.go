package goutils

import (
	"encoding/json"
	"strings"
	"testing"

	"go.yaml.in/yaml/v4"
)

func TestAllOrListString(t *testing.T) {
	t.Run("json", func(t *testing.T) {
		var aos AllOrListString

		err := json.Unmarshal([]byte(`"*"`), &aos)
		if err != nil {
			t.Fatal(err)
		}

		if aos.String() != "*" {
			t.Fatalf("expected *, got: %s", aos)
		}

		rawAOS, err := json.Marshal(aos)
		if err != nil {
			t.Fatal(err)
		}

		if string(rawAOS) != `"*"` {
			t.Fatalf("expected *, got: %s", aos)
		}

		listBytes := []byte(`["foo","bar"]`)

		err = json.Unmarshal(listBytes, &aos)
		if err != nil {
			t.Fatal(err)
		}

		expected := "[foo bar]"

		if aos.String() != expected {
			t.Fatalf("expected %s, got: %s", expected, aos)
		}

		rawAOS, err = json.Marshal(aos)
		if err != nil {
			t.Fatal(err)
		}

		if string(rawAOS) != string(listBytes) {
			t.Fatalf("expected %s, got: %s", string(listBytes), string(rawAOS))
		}
	})

	t.Run("yaml", func(t *testing.T) {
		var aos AllOrListString

		err := yaml.Unmarshal([]byte(`"*"`), &aos)
		if err != nil {
			t.Fatal(err)
		}

		if aos.String() != "*" {
			t.Fatalf("expected *, got: %s", aos)
		}

		rawAOS, err := yaml.Marshal(aos)
		if err != nil {
			t.Fatal(err)
		}

		if strings.TrimSpace(string(rawAOS)) != `'*'` {
			t.Fatalf("expected *, got: %s", string(rawAOS))
		}

		listBytes := []byte("- foo\n- bar")

		err = yaml.Unmarshal(listBytes, &aos)
		if err != nil {
			t.Fatal(err)
		}

		expected := "[foo bar]"

		if aos.String() != expected {
			t.Fatalf("expected %s, got: %s", expected, aos)
		}

		rawAOS, err = yaml.Marshal(aos)
		if err != nil {
			t.Fatal(err)
		}

		if strings.TrimSpace(string(rawAOS)) != string(listBytes) {
			t.Fatalf("expected %s, got: %s", string(listBytes), string(rawAOS))
		}
	})
}
