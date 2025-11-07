package goutils

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"go.yaml.in/yaml/v4"
)

var (
	errFilePathRequired             = errors.New("file path is required")
	errUnsupportedFilePathExtension = errors.New("only {json,yaml,yml} extension is supported")
)

// ReadJSONOrYAMLFile reads and decodes the json or yaml file from the path.
func ReadJSONOrYAMLFile[T any](filePath string) (*T, error) {
	filePath = strings.TrimSpace(filePath)
	if filePath == "" {
		return nil, errFilePathRequired
	}

	filePath = filepath.Clean(filePath)

	var result T

	ext := filepath.Ext(filePath)

	switch ext {
	case ".json":
		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}

		defer CatchWarnError(file.Close)

		err = json.NewDecoder(file).Decode(&result)

		return &result, err
	case ".yaml", ".yml":
		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}

		defer CatchWarnError(file.Close)

		err = yaml.NewDecoder(file).Decode(&result)

		return &result, err
	default:
		return nil, errUnsupportedFilePathExtension
	}
}
