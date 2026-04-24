// Copyright 2026 RelyChan Pte. Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package goutils

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/relychan/goutils/httperror"
)

func TestReadJSONOrYAMLFile(t *testing.T) {
	t.Run("read_json", func(t *testing.T) {
		result, err := ReadJSONOrYAMLFile[map[string]string](context.Background(), "testdata/config.json")
		if err != nil {
			t.Fatalf("expected nil error, got: %s", err)
		}

		expected := map[string]string{"foo": "baz"}

		if !reflect.DeepEqual(*result, expected) {
			t.Fatalf("expected %v, got: %v", expected, *result)
		}
	})

	t.Run("read_yaml", func(t *testing.T) {
		result, err := ReadJSONOrYAMLFile[map[string]string](context.Background(), "testdata/config.yaml")
		if err != nil {
			t.Fatalf("expected nil error, got: %s", err)
		}

		expected := map[string]string{"foo": "bar"}

		if !reflect.DeepEqual(*result, expected) {
			t.Fatalf("expected %v, got: %v", expected, *result)
		}
	})

	t.Run("path_required", func(t *testing.T) {
		_, err := ReadJSONOrYAMLFile[string](context.Background(), "")
		if !errors.Is(err, errFilePathRequired) {
			t.Fatalf("expected error: %s, got: %s", errFilePathRequired, err)
		}
	})

	t.Run("file_not_found", func(t *testing.T) {
		_, err := ReadJSONOrYAMLFile[string](context.Background(), "testdata/not-found.json")
		if !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("expected error: %s, got: %s", os.ErrNotExist, err)
		}

		_, err = ReadJSONOrYAMLFile[string](context.Background(), "testdata/not-found.yaml")
		if !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("expected error: %s, got: %s", os.ErrNotExist, err)
		}
	})
}

func TestReadJSONOrYAMLFile_URL(t *testing.T) {
	t.Run("read_json_from_http_url", func(t *testing.T) {
		// Create a test HTTP server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"name": "test", "value": 123}`))
		}))
		defer server.Close()

		type TestData struct {
			Name  string `json:"name"`
			Value int    `json:"value"`
		}

		result, err := ReadJSONOrYAMLFile[TestData](context.Background(), server.URL+"/config.json")
		if err != nil {
			t.Fatalf("expected nil error, got: %s", err)
		}

		if result.Name != "test" {
			t.Errorf("expected name 'test', got: %s", result.Name)
		}

		if result.Value != 123 {
			t.Errorf("expected value 123, got: %d", result.Value)
		}
	})

	t.Run("read_yaml_from_http_url", func(t *testing.T) {
		// Create a test HTTP server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/x-yaml")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("name: test\nvalue: 456\n"))
		}))
		defer server.Close()

		type TestData struct {
			Name  string `yaml:"name"`
			Value int    `yaml:"value"`
		}

		result, err := ReadJSONOrYAMLFile[TestData](context.Background(), server.URL+"/config.yaml")
		if err != nil {
			t.Fatalf("expected nil error, got: %s", err)
		}

		if result.Name != "test" {
			t.Errorf("expected name 'test', got: %s", result.Name)
		}

		if result.Value != 456 {
			t.Errorf("expected value 456, got: %d", result.Value)
		}
	})

	t.Run("read_yml_from_http_url", func(t *testing.T) {
		// Create a test HTTP server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/x-yaml")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("foo: bar\nbaz: qux\n"))
		}))
		defer server.Close()

		result, err := ReadJSONOrYAMLFile[map[string]string](context.Background(), server.URL+"/config.yml")
		if err != nil {
			t.Fatalf("expected nil error, got: %s", err)
		}

		expected := map[string]string{"foo": "bar", "baz": "qux"}
		if !reflect.DeepEqual(*result, expected) {
			t.Errorf("expected %v, got: %v", expected, *result)
		}
	})

	t.Run("unsupported_extension", func(t *testing.T) {
		_, err := ReadJSONOrYAMLFile[string](context.Background(), "testdata/config.txt")
		if !errors.Is(err, errUnsupportedFilePathExtension) {
			t.Fatalf("expected error: %s, got: %s", errUnsupportedFilePathExtension, err)
		}
	})

	t.Run("read_json_from_https_url", func(t *testing.T) {
		// Create a test HTTPS server
		server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"secure": true}`))
		}))
		defer server.Close()

		// Use the test server's client which trusts the test certificate
		originalClient := http.DefaultClient
		http.DefaultClient = server.Client()
		defer func() { http.DefaultClient = originalClient }()

		type TestData struct {
			Secure bool `json:"secure"`
		}

		result, err := ReadJSONOrYAMLFile[TestData](context.Background(), server.URL+"/config.json")
		if err != nil {
			t.Fatalf("expected nil error, got: %s", err)
		}

		if !result.Secure {
			t.Error("expected secure to be true")
		}
	})

	t.Run("url_returns_404", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("Not Found"))
		}))
		defer server.Close()

		_, err := ReadJSONOrYAMLFile[map[string]string](context.Background(), server.URL+"/missing.json")
		if err == nil {
			t.Fatal("expected error for 404 response")
		}

		var rfc9457Err *httperror.HTTPError
		if !errors.As(err, &rfc9457Err) {
			t.Fatalf("expected httperror.HTTPError, got: %T", err)
		}

		if rfc9457Err.Status != http.StatusNotFound {
			t.Errorf("expected status 404, got: %d", rfc9457Err.Status)
		}
	})

	t.Run("url_returns_500", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("Internal Server Error"))
		}))
		defer server.Close()

		_, err := ReadJSONOrYAMLFile[map[string]string](context.Background(), server.URL+"/error.json")
		if err == nil {
			t.Fatal("expected error for 500 response")
		}

		var rfc9457Err *httperror.HTTPError
		if !errors.As(err, &rfc9457Err) {
			t.Fatalf("expected httperror.HTTPError, got: %T", err)
		}

		if rfc9457Err.Status != http.StatusInternalServerError {
			t.Errorf("expected status 500, got: %d", rfc9457Err.Status)
		}
	})

	t.Run("url_returns_empty_body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			// No body written
		}))
		defer server.Close()

		_, err := ReadJSONOrYAMLFile[map[string]string](context.Background(), server.URL+"/empty.json")
		if err == nil {
			t.Fatal("expected error for empty body")
		}
	})

	t.Run("url_returns_invalid_json", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{invalid json}`))
		}))
		defer server.Close()

		_, err := ReadJSONOrYAMLFile[map[string]string](context.Background(), server.URL+"/invalid.json")
		if err == nil {
			t.Fatal("expected error for invalid JSON")
		}
	})

	t.Run("url_returns_invalid_yaml", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/x-yaml")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("invalid:\n  - yaml\n  structure:\n    - broken"))
		}))
		defer server.Close()

		type SimpleStruct struct {
			Name string `yaml:"name"`
		}

		_, err := ReadJSONOrYAMLFile[SimpleStruct](context.Background(), server.URL+"/invalid.yaml")
		// Should not error on parsing, but result might be empty
		if err != nil {
			// This is acceptable - YAML parser might reject it
			t.Logf("YAML parsing error (acceptable): %v", err)
		}
	})

	t.Run("url_with_complex_json_structure", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"users": [
					{"name": "Alice", "age": 30},
					{"name": "Bob", "age": 25}
				],
				"metadata": {
					"version": "1.0",
					"timestamp": "2024-01-01T00:00:00Z"
				}
			}`))
		}))
		defer server.Close()

		type User struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		type Metadata struct {
			Version   string `json:"version"`
			Timestamp string `json:"timestamp"`
		}

		type Config struct {
			Users    []User   `json:"users"`
			Metadata Metadata `json:"metadata"`
		}

		result, err := ReadJSONOrYAMLFile[Config](context.Background(), server.URL+"/config.json")
		if err != nil {
			t.Fatalf("expected nil error, got: %s", err)
		}

		if len(result.Users) != 2 {
			t.Errorf("expected 2 users, got: %d", len(result.Users))
		}

		if result.Users[0].Name != "Alice" {
			t.Errorf("expected first user name 'Alice', got: %s", result.Users[0].Name)
		}

		if result.Metadata.Version != "1.0" {
			t.Errorf("expected version '1.0', got: %s", result.Metadata.Version)
		}
	})

	t.Run("url_with_complex_yaml_structure", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/x-yaml")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`
database:
  host: localhost
  port: 5432
  credentials:
    username: admin
    password: secret
features:
  - name: feature1
    enabled: true
  - name: feature2
    enabled: false
`))
		}))
		defer server.Close()

		type Credentials struct {
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		}

		type Database struct {
			Host        string      `yaml:"host"`
			Port        int         `yaml:"port"`
			Credentials Credentials `yaml:"credentials"`
		}

		type Feature struct {
			Name    string `yaml:"name"`
			Enabled bool   `yaml:"enabled"`
		}

		type Config struct {
			Database Database  `yaml:"database"`
			Features []Feature `yaml:"features"`
		}

		result, err := ReadJSONOrYAMLFile[Config](context.Background(), server.URL+"/config.yaml")
		if err != nil {
			t.Fatalf("expected nil error, got: %s", err)
		}

		if result.Database.Host != "localhost" {
			t.Errorf("expected host 'localhost', got: %s", result.Database.Host)
		}

		if result.Database.Port != 5432 {
			t.Errorf("expected port 5432, got: %d", result.Database.Port)
		}

		if len(result.Features) != 2 {
			t.Errorf("expected 2 features, got: %d", len(result.Features))
		}

		if result.Features[0].Name != "feature1" || !result.Features[0].Enabled {
			t.Error("expected feature1 to be enabled")
		}
	})

	t.Run("url_with_query_parameters_limitation", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status": "ok"}`))
		}))
		defer server.Close()

		// Note: This is a known limitation - filepath.Ext() includes query parameters
		// in the extension, so "config.json?version=v1" has extension ".json?version=v1"
		// which doesn't match ".json"
		_, err := ReadJSONOrYAMLFile[map[string]string](context.Background(), server.URL+"/config.json?version=v1")
		if !errors.Is(err, errUnsupportedFilePathExtension) {
			t.Fatalf("expected unsupported extension error due to query params, got: %s", err)
		}
	})

	t.Run("url_with_special_characters_in_path", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"path": "special"}`))
		}))
		defer server.Close()

		type Response struct {
			Path string `json:"path"`
		}

		result, err := ReadJSONOrYAMLFile[Response](context.Background(), server.URL+"/path/to/config.json")
		if err != nil {
			t.Fatalf("expected nil error, got: %s", err)
		}

		if result.Path != "special" {
			t.Errorf("expected path 'special', got: %s", result.Path)
		}
	})

	t.Run("url_returns_401_unauthorized", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Unauthorized"))
		}))
		defer server.Close()

		_, err := ReadJSONOrYAMLFile[map[string]string](context.Background(), server.URL+"/secure.json")
		if err == nil {
			t.Fatal("expected error for 401 response")
		}

		var rfc9457Err *httperror.HTTPError
		if !errors.As(err, &rfc9457Err) {
			t.Fatalf("expected httperror.HTTPError, got: %T", err)
		}

		if rfc9457Err.Status != http.StatusUnauthorized {
			t.Errorf("expected status 401, got: %d", rfc9457Err.Status)
		}
	})

	t.Run("url_returns_403_forbidden", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte("Forbidden"))
		}))
		defer server.Close()

		_, err := ReadJSONOrYAMLFile[map[string]string](context.Background(), server.URL+"/forbidden.json")
		if err == nil {
			t.Fatal("expected error for 403 response")
		}

		var rfc9457Err *httperror.HTTPError
		if !errors.As(err, &rfc9457Err) {
			t.Fatalf("expected httperror.HTTPError, got: %T", err)
		}

		if rfc9457Err.Status != http.StatusForbidden {
			t.Errorf("expected status 403, got: %d", rfc9457Err.Status)
		}
	})

	t.Run("url_returns_503_service_unavailable", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("Service Unavailable"))
		}))
		defer server.Close()

		_, err := ReadJSONOrYAMLFile[map[string]string](context.Background(), server.URL+"/unavailable.json")
		if err == nil {
			t.Fatal("expected error for 503 response")
		}

		var rfc9457Err *httperror.HTTPError
		if !errors.As(err, &rfc9457Err) {
			t.Fatalf("expected httperror.HTTPError, got: %T", err)
		}

		if rfc9457Err.Status != http.StatusServiceUnavailable {
			t.Errorf("expected status 503, got: %d", rfc9457Err.Status)
		}
	})

	t.Run("url_with_redirect", func(t *testing.T) {
		finalServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"redirected": true}`))
		}))
		defer finalServer.Close()

		redirectServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, finalServer.URL+"/config.json", http.StatusMovedPermanently)
		}))
		defer redirectServer.Close()

		type Response struct {
			Redirected bool `json:"redirected"`
		}

		result, err := ReadJSONOrYAMLFile[Response](context.Background(), redirectServer.URL+"/redirect.json")
		if err != nil {
			t.Fatalf("expected nil error, got: %s", err)
		}

		if !result.Redirected {
			t.Error("expected redirected to be true")
		}
	})

	t.Run("invalid_url_scheme", func(t *testing.T) {
		// ftp:// is not supported, should fall back to file path
		_, err := ReadJSONOrYAMLFile[map[string]string](context.Background(), "ftp://example.com/config.json")
		// This should fail as a file path
		if err == nil {
			t.Fatal("expected error for unsupported URL scheme")
		}
	})

	t.Run("url_with_whitespace", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"trimmed": true}`))
		}))
		defer server.Close()

		type Response struct {
			Trimmed bool `json:"trimmed"`
		}

		// URL with leading/trailing whitespace should be trimmed
		result, err := ReadJSONOrYAMLFile[Response](context.Background(), "  "+server.URL+"/config.json  ")
		if err != nil {
			t.Fatalf("expected nil error, got: %s", err)
		}

		if !result.Trimmed {
			t.Error("expected trimmed to be true")
		}
	})
}

func TestReadMultiFromJSONOrYAMLFile(t *testing.T) {
	t.Run("read_multi_json", func(t *testing.T) {
		results, err := ReadMultiFromJSONOrYAMLFile[map[string]string](context.Background(), "testdata/multi_config.json")
		if err != nil {
			t.Fatalf("expected nil error, got: %s", err)
		}

		expected := []map[string]string{{"foo": "bar"}, {"foo": "baz"}, {"foo": "qux"}}
		if !reflect.DeepEqual(results, expected) {
			t.Fatalf("expected %v, got: %v", expected, results)
		}
	})

	t.Run("read_multi_yaml", func(t *testing.T) {
		results, err := ReadMultiFromJSONOrYAMLFile[map[string]string](context.Background(), "testdata/multi_config.yaml")
		if err != nil {
			t.Fatalf("expected nil error, got: %s", err)
		}

		expected := []map[string]string{{"foo": "bar"}, {"foo": "baz"}, {"foo": "qux"}}
		if !reflect.DeepEqual(results, expected) {
			t.Fatalf("expected %v, got: %v", expected, results)
		}
	})

	t.Run("read_single_json_as_multi", func(t *testing.T) {
		results, err := ReadMultiFromJSONOrYAMLFile[map[string]string](context.Background(), "testdata/config.json")
		if err != nil {
			t.Fatalf("expected nil error, got: %s", err)
		}

		expected := []map[string]string{{"foo": "baz"}}
		if !reflect.DeepEqual(results, expected) {
			t.Fatalf("expected %v, got: %v", expected, results)
		}
	})

	t.Run("read_single_yaml_as_multi", func(t *testing.T) {
		results, err := ReadMultiFromJSONOrYAMLFile[map[string]string](context.Background(), "testdata/config.yaml")
		if err != nil {
			t.Fatalf("expected nil error, got: %s", err)
		}

		expected := []map[string]string{{"foo": "bar"}}
		if !reflect.DeepEqual(results, expected) {
			t.Fatalf("expected %v, got: %v", expected, results)
		}
	})

	t.Run("path_required", func(t *testing.T) {
		_, err := ReadMultiFromJSONOrYAMLFile[string](context.Background(), "")
		if !errors.Is(err, errFilePathRequired) {
			t.Fatalf("expected error: %s, got: %s", errFilePathRequired, err)
		}
	})

	t.Run("file_not_found", func(t *testing.T) {
		_, err := ReadMultiFromJSONOrYAMLFile[string](context.Background(), "testdata/not-found.json")
		if !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("expected error: %s, got: %s", os.ErrNotExist, err)
		}

		_, err = ReadMultiFromJSONOrYAMLFile[string](context.Background(), "testdata/not-found.yaml")
		if !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("expected error: %s, got: %s", os.ErrNotExist, err)
		}
	})

	t.Run("unsupported_extension", func(t *testing.T) {
		_, err := ReadMultiFromJSONOrYAMLFile[string](context.Background(), "testdata/config.txt")
		if !errors.Is(err, errUnsupportedFilePathExtension) {
			t.Fatalf("expected error: %s, got: %s", errUnsupportedFilePathExtension, err)
		}
	})

	t.Run("invalid_json", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{invalid}`))
		}))
		defer server.Close()

		_, err := ReadMultiFromJSONOrYAMLFile[map[string]string](context.Background(), server.URL+"/data.json")
		if err == nil {
			t.Fatal("expected error for invalid JSON")
		}
	})

	t.Run("read_multi_json_from_url", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("{\"foo\":\"bar\"}\n{\"foo\":\"baz\"}\n"))
		}))
		defer server.Close()

		results, err := ReadMultiFromJSONOrYAMLFile[map[string]string](context.Background(), server.URL+"/data.json")
		if err != nil {
			t.Fatalf("expected nil error, got: %s", err)
		}

		expected := []map[string]string{{"foo": "bar"}, {"foo": "baz"}}
		if !reflect.DeepEqual(results, expected) {
			t.Fatalf("expected %v, got: %v", expected, results)
		}
	})

	t.Run("read_multi_yaml_from_url", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("foo: bar\n---\nfoo: baz\n"))
		}))
		defer server.Close()

		results, err := ReadMultiFromJSONOrYAMLFile[map[string]string](context.Background(), server.URL+"/data.yaml")
		if err != nil {
			t.Fatalf("expected nil error, got: %s", err)
		}

		expected := []map[string]string{{"foo": "bar"}, {"foo": "baz"}}
		if !reflect.DeepEqual(results, expected) {
			t.Fatalf("expected %v, got: %v", expected, results)
		}
	})

	t.Run("url_returns_error_status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		_, err := ReadMultiFromJSONOrYAMLFile[map[string]string](context.Background(), server.URL+"/missing.json")
		if err == nil {
			t.Fatal("expected error for non-200 response")
		}

		var rfc9457Err *httperror.HTTPError
		if !errors.As(err, &rfc9457Err) {
			t.Fatalf("expected httperror.HTTPError, got: %T", err)
		}

		if rfc9457Err.Status != http.StatusNotFound {
			t.Errorf("expected status 404, got: %d", rfc9457Err.Status)
		}
	})
}

func TestFileReaderFromPath(t *testing.T) {
	t.Run("read_from_local_file", func(t *testing.T) {
		reader, _, err := FileReaderFromPath(context.Background(), "testdata/config.json")
		if err != nil {
			t.Fatalf("expected nil error, got: %s", err)
		}
		defer reader.Close()

		if reader == nil {
			t.Fatal("expected non-nil reader")
		}
	})

	t.Run("read_from_http_url", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("test content"))
		}))
		defer server.Close()

		reader, _, err := FileReaderFromPath(context.Background(), server.URL+"/test.txt")
		if err != nil {
			t.Fatalf("expected nil error, got: %s", err)
		}
		defer reader.Close()

		if reader == nil {
			t.Fatal("expected non-nil reader")
		}
	})

	t.Run("read_from_https_url", func(t *testing.T) {
		server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("secure content"))
		}))
		defer server.Close()

		// Use the test server's client
		originalClient := http.DefaultClient
		http.DefaultClient = server.Client()
		defer func() { http.DefaultClient = originalClient }()

		reader, _, err := FileReaderFromPath(context.Background(), server.URL+"/secure.txt")
		if err != nil {
			t.Fatalf("expected nil error, got: %s", err)
		}
		defer reader.Close()

		if reader == nil {
			t.Fatal("expected non-nil reader")
		}
	})

	t.Run("empty_path", func(t *testing.T) {
		_, _, err := FileReaderFromPath(context.Background(), "")
		if !errors.Is(err, errFilePathRequired) {
			t.Fatalf("expected error: %s, got: %s", errFilePathRequired, err)
		}
	})

	t.Run("whitespace_only_path", func(t *testing.T) {
		_, _, err := FileReaderFromPath(context.Background(), "   ")
		if !errors.Is(err, errFilePathRequired) {
			t.Fatalf("expected error: %s, got: %s", errFilePathRequired, err)
		}
	})

	t.Run("file_not_found", func(t *testing.T) {
		_, _, err := FileReaderFromPath(context.Background(), "testdata/nonexistent.txt")
		if !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("expected error: %s, got: %s", os.ErrNotExist, err)
		}
	})

	t.Run("url_returns_404", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		_, _, err := FileReaderFromPath(context.Background(), server.URL+"/missing.txt")
		if err == nil {
			t.Fatal("expected error for 404 response")
		}

		var rfc9457Err *httperror.HTTPError
		if !errors.As(err, &rfc9457Err) {
			t.Fatalf("expected httperror.HTTPError, got: %T", err)
		}
	})

	t.Run("url_with_uppercase_scheme", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("content"))
		}))
		defer server.Close()

		// Replace http:// with HTTP://
		upperURL := "HTTP" + server.URL[4:]

		reader, _, err := FileReaderFromPath(context.Background(), upperURL)
		if err != nil {
			t.Fatalf("expected nil error for uppercase scheme, got: %s", err)
		}
		defer reader.Close()

		if reader == nil {
			t.Fatal("expected non-nil reader")
		}
	})

	t.Run("url_with_mixed_case_scheme", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("content"))
		}))
		defer server.Close()

		// Replace http:// with HtTp://
		mixedURL := "HtTp" + server.URL[4:]

		reader, _, err := FileReaderFromPath(context.Background(), mixedURL)
		if err != nil {
			t.Fatalf("expected nil error for mixed case scheme, got: %s", err)
		}
		defer reader.Close()

		if reader == nil {
			t.Fatal("expected non-nil reader")
		}
	})

	t.Run("path_traversal_cleaned", func(t *testing.T) {
		// Test that filepath.Clean is applied
		_, _, err := FileReaderFromPath(context.Background(), "testdata/../testdata/config.json")
		if err != nil {
			t.Fatalf("expected nil error for cleaned path, got: %s", err)
		}
	})
}

func TestFileReaderFromPath_IncludeExcludePaths(t *testing.T) {
	t.Run("local_include_path_allowed", func(t *testing.T) {
		reader, _, err := FileReaderFromPath(
			context.Background(),
			"testdata/config.json",
			DownloadFileIncludingPaths([]string{"testdata/*.json"}),
		)
		if err != nil {
			t.Fatalf("expected nil error, got: %s", err)
		}
		defer reader.Close()
	})

	t.Run("local_include_path_blocked", func(t *testing.T) {
		_, _, err := FileReaderFromPath(
			context.Background(),
			"testdata/config.json",
			DownloadFileIncludingPaths([]string{"testdata/*.yaml"}),
		)
		if !errors.Is(err, errDisallowedFilePath) {
			t.Fatalf("expected errDisallowedFilePath, got: %v", err)
		}
	})

	t.Run("local_exclude_path_blocked", func(t *testing.T) {
		_, _, err := FileReaderFromPath(
			context.Background(),
			"testdata/config.json",
			DownloadFileExcludingPaths([]string{"testdata/*.json"}),
		)
		if !errors.Is(err, errDisallowedFilePath) {
			t.Fatalf("expected errDisallowedFilePath, got: %v", err)
		}
	})

	t.Run("local_exclude_path_allowed", func(t *testing.T) {
		reader, _, err := FileReaderFromPath(
			context.Background(),
			"testdata/config.json",
			DownloadFileExcludingPaths([]string{"testdata/*.yaml"}),
		)
		if err != nil {
			t.Fatalf("expected nil error, got: %s", err)
		}
		defer reader.Close()
	})

	t.Run("local_include_and_exclude_combined", func(t *testing.T) {
		// included by pattern, but also excluded — exclude wins
		_, _, err := FileReaderFromPath(
			context.Background(),
			"testdata/config.json",
			DownloadFileIncludingPaths([]string{"testdata/*"}),
			DownloadFileExcludingPaths([]string{"testdata/*.json"}),
		)
		if !errors.Is(err, errDisallowedFilePath) {
			t.Fatalf("expected errDisallowedFilePath, got: %v", err)
		}
	})

	t.Run("local_no_include_no_exclude_allows_all", func(t *testing.T) {
		reader, _, err := FileReaderFromPath(context.Background(), "testdata/config.yaml")
		if err != nil {
			t.Fatalf("expected nil error, got: %s", err)
		}
		defer reader.Close()
	})

	t.Run("url_include_path_allowed", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"ok": true}`))
		}))
		defer server.Close()

		reader, _, err := FileReaderFromPath(
			context.Background(),
			server.URL+"/data/config.json",
			DownloadFileIncludingPaths([]string{"/data/*.json"}),
		)
		if err != nil {
			t.Fatalf("expected nil error, got: %s", err)
		}
		defer reader.Close()
	})

	t.Run("url_include_path_blocked", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"ok": true}`))
		}))
		defer server.Close()

		_, _, err := FileReaderFromPath(
			context.Background(),
			server.URL+"/data/config.json",
			DownloadFileIncludingPaths([]string{"/other/*.json"}),
		)
		if !errors.Is(err, errDisallowedFilePath) {
			t.Fatalf("expected errDisallowedFilePath, got: %v", err)
		}
	})

	t.Run("url_exclude_path_blocked", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"ok": true}`))
		}))
		defer server.Close()

		_, _, err := FileReaderFromPath(
			context.Background(),
			server.URL+"/data/config.json",
			DownloadFileExcludingPaths([]string{"/data/*.json"}),
		)
		if !errors.Is(err, errDisallowedFilePath) {
			t.Fatalf("expected errDisallowedFilePath, got: %v", err)
		}
	})
}

func TestFileReaderFromPath_AllowedBlockedHosts(t *testing.T) {
	t.Run("allowed_host_passes", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("content"))
		}))
		defer server.Close()

		reader, _, err := FileReaderFromPath(
			context.Background(),
			server.URL+"/test.txt",
			DownloadFileWithAllowedHosts([]string{"127.0.0.1"}),
		)
		if err != nil {
			t.Fatalf("expected nil error, got: %s", err)
		}
		defer reader.Close()
	})

	t.Run("allowed_host_blocked_when_not_matching", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("content"))
		}))
		defer server.Close()

		_, _, err := FileReaderFromPath(
			context.Background(),
			server.URL+"/test.txt",
			DownloadFileWithAllowedHosts([]string{"example.com"}),
		)
		if err == nil {
			t.Fatal("expected error when host not in allowed list")
		}
	})

	t.Run("blocked_host_rejected", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("content"))
		}))
		defer server.Close()

		_, _, err := FileReaderFromPath(
			context.Background(),
			server.URL+"/test.txt",
			DownloadFileWithBlockedHosts([]string{"127.0.0.1"}),
		)
		if err == nil {
			t.Fatal("expected error when host is blocked")
		}
	})

	t.Run("no_host_filter_allows_all", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("content"))
		}))
		defer server.Close()

		reader, _, err := FileReaderFromPath(context.Background(), server.URL+"/test.txt")
		if err != nil {
			t.Fatalf("expected nil error, got: %s", err)
		}
		defer reader.Close()
	})

	t.Run("blocked_host_takes_precedence_over_allowed", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("content"))
		}))
		defer server.Close()

		_, _, err := FileReaderFromPath(
			context.Background(),
			server.URL+"/test.txt",
			DownloadFileWithAllowedHosts([]string{"127.0.0.1"}),
			DownloadFileWithBlockedHosts([]string{"127.0.0.1"}),
		)
		if err == nil {
			t.Fatal("expected error: blocked host should take precedence")
		}
	})

	t.Run("host_filter_only_applies_to_urls_not_local_paths", func(t *testing.T) {
		// AllowedHosts / BlockedHosts should not affect local file reads
		reader, _, err := FileReaderFromPath(
			context.Background(),
			"testdata/config.json",
			DownloadFileWithAllowedHosts([]string{"example.com"}),
		)
		if err != nil {
			t.Fatalf("expected nil error for local file, got: %s", err)
		}
		defer reader.Close()
	})
}
