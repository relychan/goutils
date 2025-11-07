package goutils

import (
	"errors"
	"os"
	"reflect"
	"testing"
)

func TestReadJSONOrYAMLFile(t *testing.T) {
	t.Run("read_json", func(t *testing.T) {
		result, err := ReadJSONOrYAMLFile[map[string]string]("testdata/config.json")
		if err != nil {
			t.Fatalf("expected nil error, got: %s", err)
		}

		expected := map[string]string{"foo": "baz"}

		if !reflect.DeepEqual(*result, expected) {
			t.Fatalf("expected %v, got: %v", expected, *result)
		}
	})

	t.Run("read_yaml", func(t *testing.T) {
		result, err := ReadJSONOrYAMLFile[map[string]string]("testdata/config.yaml")
		if err != nil {
			t.Fatalf("expected nil error, got: %s", err)
		}

		expected := map[string]string{"foo": "bar"}

		if !reflect.DeepEqual(*result, expected) {
			t.Fatalf("expected %v, got: %v", expected, *result)
		}
	})

	t.Run("path_required", func(t *testing.T) {
		_, err := ReadJSONOrYAMLFile[string]("")
		if !errors.Is(err, errFilePathRequired) {
			t.Fatalf("expected error: %s, got: %s", errFilePathRequired, err)
		}
	})

	t.Run("unsupported_extension", func(t *testing.T) {
		_, err := ReadJSONOrYAMLFile[string]("testdata/not-found.txt")
		if !errors.Is(err, errUnsupportedFilePathExtension) {
			t.Fatalf("expected error: %s, got: %s", errUnsupportedFilePathExtension, err)
		}
	})

	t.Run("file_not_found", func(t *testing.T) {
		_, err := ReadJSONOrYAMLFile[string]("testdata/not-found.json")
		if !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("expected error: %s, got: %s", os.ErrNotExist, err)
		}

		_, err = ReadJSONOrYAMLFile[string]("testdata/not-found.yaml")
		if !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("expected error: %s, got: %s", os.ErrNotExist, err)
		}
	})
}
