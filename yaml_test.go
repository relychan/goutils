package goutils

import (
	"errors"
	"testing"

	"go.yaml.in/yaml/v4"
)

func buildYAMLMappingNode(pairs ...string) *yaml.Node {
	node := &yaml.Node{Kind: yaml.MappingNode, Tag: YAMLMapTag}
	for i := 0; i < len(pairs)-1; i += 2 {
		key := &yaml.Node{Kind: yaml.ScalarNode, Tag: YAMLStrTag, Value: pairs[i]}
		value := &yaml.Node{Kind: yaml.ScalarNode, Tag: YAMLStrTag, Value: pairs[i+1]}
		node.Content = append(node.Content, key, value)
	}
	return node
}

var rawYamlStr = `foo: bar
baz: qux`

type yamlStringMap struct {
	Foo string `yaml:"foo"`
	Baz string `yaml:"baz"`
}

func (ysm *yamlStringMap) UnmarshalYAML(value *yaml.Node) error {
	foo, err := GetStringValueFromYAMLMap(value, "foo")
	if err != nil {
		return err
	}

	if foo != nil {
		ysm.Foo = *foo
	}

	baz, err := GetStringValueFromYAMLMap(value, "baz")
	if err != nil {
		return err
	}

	if foo != nil {
		ysm.Baz = *baz
	}

	return nil
}

func TestGetStringValueFromYAMLMap(t *testing.T) {
	t.Run("returns string value for matching key", func(t *testing.T) {
		node := buildYAMLMappingNode("foo", "bar", "baz", "qux")
		result, err := GetStringValueFromYAMLMap(node, "foo")
		assertNilError(t, err)
		if result == nil {
			t.Fatal("expected non-nil result")
		}
		assertEqual(t, "bar", *result)
	})

	t.Run("returns second key value", func(t *testing.T) {
		node := buildYAMLMappingNode("foo", "bar", "baz", "qux")
		result, err := GetStringValueFromYAMLMap(node, "baz")
		assertNilError(t, err)
		if result == nil {
			t.Fatal("expected non-nil result")
		}
		assertEqual(t, "qux", *result)
	})

	t.Run("returns nil for missing key", func(t *testing.T) {
		node := buildYAMLMappingNode("foo", "bar")
		result, err := GetStringValueFromYAMLMap(node, "missing")
		assertNilError(t, err)
		if result != nil {
			t.Fatalf("expected nil result, got %q", *result)
		}
	})

	t.Run("returns nil for empty mapping node", func(t *testing.T) {
		node := &yaml.Node{Kind: yaml.MappingNode, Tag: YAMLMapTag}
		result, err := GetStringValueFromYAMLMap(node, "key")
		assertNilError(t, err)
		if result != nil {
			t.Fatalf("expected nil result, got %q", *result)
		}
	})

	t.Run("returns nil when value tag is null", func(t *testing.T) {
		node := &yaml.Node{Kind: yaml.MappingNode, Tag: YAMLMapTag}
		node.Content = []*yaml.Node{
			{Kind: yaml.ScalarNode, Tag: YAMLStrTag, Value: "key"},
			{Kind: yaml.ScalarNode, Tag: YAMLNullTag, Value: "null"},
		}
		result, err := GetStringValueFromYAMLMap(node, "key")
		assertNilError(t, err)
		if result != nil {
			t.Fatalf("expected nil result, got %q", *result)
		}
	})

	t.Run("returns nil when node is nil", func(t *testing.T) {
		result, err := GetStringValueFromYAMLMap(nil, "key")
		assertNilError(t, err)
		if result != nil {
			t.Fatalf("expected nil result, got %q", *result)
		}
	})

	t.Run("returns error when node is not a mapping node", func(t *testing.T) {
		node := &yaml.Node{Kind: yaml.ScalarNode, Tag: YAMLStrTag, Value: "scalar"}
		_, err := GetStringValueFromYAMLMap(node, "key")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, ErrInvalidYAMLSyntax) {
			t.Fatalf("expected ErrInvalidYAMLSyntax, got %v", err)
		}
	})

	t.Run("returns error when value is not a string tag", func(t *testing.T) {
		node := &yaml.Node{Kind: yaml.MappingNode, Tag: YAMLMapTag}
		node.Content = []*yaml.Node{
			{Kind: yaml.ScalarNode, Tag: YAMLStrTag, Value: "count"},
			{Kind: yaml.ScalarNode, Tag: YAMLIntTag, Value: "42"},
		}
		_, err := GetStringValueFromYAMLMap(node, "count")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, ErrInvalidYAMLSyntax) {
			t.Fatalf("expected ErrInvalidYAMLSyntax, got %v", err)
		}
	})

	t.Run("returns error when value is a bool tag", func(t *testing.T) {
		node := &yaml.Node{Kind: yaml.MappingNode, Tag: YAMLMapTag}
		node.Content = []*yaml.Node{
			{Kind: yaml.ScalarNode, Tag: YAMLStrTag, Value: "flag"},
			{Kind: yaml.ScalarNode, Tag: YAMLBoolTag, Value: "true"},
		}
		_, err := GetStringValueFromYAMLMap(node, "flag")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, ErrInvalidYAMLSyntax) {
			t.Fatalf("expected ErrInvalidYAMLSyntax, got %v", err)
		}
	})

	t.Run("returns string value for matching key", func(t *testing.T) {
		var result yamlStringMap
		err := yaml.Load([]byte(rawYamlStr), &result)
		assertNilError(t, err)
		assertEqual(t, "bar", result.Foo)
		assertEqual(t, "qux", result.Baz)
	})
}
