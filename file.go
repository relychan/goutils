package goutils

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"go.yaml.in/yaml/v4"
)

var (
	errFilePathRequired             = errors.New("file path is required")
	errFileNoContent                = errors.New("file has no content")
	errUnsupportedFilePathExtension = errors.New("only {json,yaml,yml} extension is supported")
)

// ReadJSONOrYAMLFile reads and decodes a JSON or YAML document from the given source,
// which may be a local file path or an HTTP/HTTPS URL.
func ReadJSONOrYAMLFile[T any](filePath string) (*T, error) {
	var result T

	file, ext, err := FileReaderFromPath(filePath)
	if err != nil {
		return nil, err
	}

	defer CatchWarnErrorFunc(file.Close)

	switch ext {
	case ".json":
		err = json.NewDecoder(file).Decode(&result)

		return &result, err
	case ".yaml", ".yml":
		err = yaml.NewDecoder(file).Decode(&result)

		return &result, err
	default:
		return nil, errUnsupportedFilePathExtension
	}
}

// FileReaderFromPath reads content from either a local filesystem path or an HTTP/HTTPS URL.
//
// Supported URL schemes are "http" and "https". If the provided path parses as a URL
// with one of these schemes, an HTTP GET request is issued using http.DefaultClient
// without an explicit timeout configured. Callers that require timeouts or custom
// HTTP behavior should arrange this outside of this helper.
//
// For other schemes, or if the input does not represent an http/https URL, the value
// is treated as a filesystem path, cleaned with filepath.Clean, and opened via os.Open.
//
// The caller is responsible for closing the returned io.ReadCloser when finished with it.
func FileReaderFromPath(filePath string) (io.ReadCloser, string, error) {
	filePath = strings.TrimSpace(filePath)
	if filePath == "" {
		return nil, "", errFilePathRequired
	}

	fileURL, err := url.Parse(filePath)
	if err == nil && slices.Contains([]string{"http", "https"}, strings.ToLower(fileURL.Scheme)) {
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, filePath, nil)
		if err != nil {
			return nil, "", err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, "", err
		}

		if resp.StatusCode != http.StatusOK {
			respError := NewRFC9457ErrorFromResponse(resp)
			respError.Title = "Read File Failure"

			return nil, "", respError
		}

		if resp.Body == nil {
			return nil, "", errFileNoContent
		}

		ext := filepath.Ext(filePath)

		return resp.Body, ext, nil
	}

	filePath = filepath.Clean(filePath)
	ext := filepath.Ext(filePath)
	reader, err := os.Open(filePath)

	return reader, ext, err
}
