package goutils

import (
	"fmt"

	"go.yaml.in/yaml/v4"
)

const (
	// YAMLNullTag represents a constant for a YAML null tag.
	YAMLNullTag = "!!null"
	// YAMLBoolTag represents a constant for a YAML boolean tag.
	YAMLBoolTag = "!!bool"
	// YAMLStrTag represents a constant for a YAML string tag.
	YAMLStrTag = "!!str"
	// YAMLIntTag represents a constant for a YAML integer tag.
	YAMLIntTag = "!!int"
	// YAMLFloatTag represents a constant for a YAML float tag.
	YAMLFloatTag = "!!float"
	// YAMLTimestampTag represents a constant for a YAML timestamp tag.
	YAMLTimestampTag = "!!timestamp"
	// YAMLSeqTag represents a constant for a YAML sequence tag.
	YAMLSeqTag = "!!seq"
	// YAMLMapTag represents a constant for a YAML map tag.
	YAMLMapTag = "!!map"
	// YAMLBinaryTag represents a constant for a YAML binary tag.
	YAMLBinaryTag = "!!binary"
	// YAMLMergeTag represents a constant for a YAML merge tag.
	YAMLMergeTag = "!!merge"
)

// GetStringValueFromYAMLMapNode gets the string value from a YAML map node.
func GetStringValueFromYAMLMapNode(node *yaml.Node, key string) (*string, error) {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("%w. Expected an object, got %s", ErrInvalidYAMLSyntax, node.Tag)
	}

	i := 0
	contentLength := len(node.Content)

	for ; i < contentLength; i++ {
		if i == contentLength-1 {
			return nil, nil
		}

		keyNode := node.Content[i]
		if keyNode.Kind != yaml.ScalarNode && keyNode.Tag != "!!str" {
			return nil, fmt.Errorf(
				"%w. Expected a key string, got %s",
				ErrInvalidYAMLSyntax,
				keyNode.Tag,
			)
		}

		if keyNode.Value != key {
			continue
		}

		i++

		valueNode := node.Content[i]

		switch valueNode.Tag {
		case YAMLStrTag:
			return &valueNode.Value, nil
		case YAMLNullTag:
			return nil, nil
		default:
			return nil, fmt.Errorf(
				"%w. Expected ref is a string, got %s",
				ErrInvalidYAMLSyntax,
				valueNode.Tag,
			)
		}
	}

	return nil, nil
}
