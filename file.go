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

// ReadJSONOrYAMLFile reads and decodes the json or yaml file from the path.
func ReadJSONOrYAMLFile[T any](filePath string) (*T, error) {
	var result T

	filePath = strings.TrimSpace(filePath)
	if filePath == "" {
		return nil, errFilePathRequired
	}

	ext := filepath.Ext(filePath)

	switch ext {
	case ".json":
		file, err := FileReaderFromPath(filePath)
		if err != nil {
			return nil, err
		}

		defer CatchWarnErrorFunc(file.Close)

		err = json.NewDecoder(file).Decode(&result)

		return &result, err
	case ".yaml", ".yml":
		file, err := FileReaderFromPath(filePath)
		if err != nil {
			return nil, err
		}

		defer CatchWarnErrorFunc(file.Close)

		err = yaml.NewDecoder(file).Decode(&result)

		return &result, err
	default:
		return nil, errUnsupportedFilePathExtension
	}
}

// FileReaderFromPath reads content from either file path or URL.
// Returns a ReadCloser instance.
func FileReaderFromPath(filePath string) (io.ReadCloser, error) {
	filePath = strings.TrimSpace(filePath)
	if filePath == "" {
		return nil, errFilePathRequired
	}

	fileURL, err := url.Parse(filePath)
	if err == nil && slices.Contains([]string{"http", "https"}, strings.ToLower(fileURL.Scheme)) {
		req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, filePath, nil)
		if err != nil {
			return nil, err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			respError := NewRFC9457ErrorFromResponse(resp)
			respError.Title = "Read File Failure"

			return nil, respError
		}

		if resp.Body == nil {
			return nil, errFileNoContent
		}

		return resp.Body, nil
	}

	filePath = filepath.Clean(filePath)

	return os.Open(filePath)
}
