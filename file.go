package goutils

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
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
func ReadJSONOrYAMLFile[T any](
	ctx context.Context,
	filePath string,
	options ...DownloadFileOption,
) (*T, error) {
	file, ext, err := FileReaderFromPath(ctx, filePath, options...)
	if err != nil {
		return nil, err
	}

	defer CatchWarnErrorFunc(file.Close)

	switch ext {
	case ".json":
		var result T

		err = json.NewDecoder(file).Decode(&result)

		return &result, err
	case ".yaml", ".yml":
		var result T

		loader, err := yaml.NewLoader(file)
		if err != nil {
			return nil, err
		}

		return &result, loader.Load(&result)
	default:
		return nil, errUnsupportedFilePathExtension
	}
}

// ReadMultiFromJSONOrYAMLFile reads and decodes multiple JSON or YAML documents from the given source,
// which may be a local file path or an HTTP/HTTPS URL.
func ReadMultiFromJSONOrYAMLFile[T any](
	ctx context.Context,
	filePath string,
	options ...DownloadFileOption,
) ([]T, error) {
	file, ext, err := FileReaderFromPath(ctx, filePath, options...)
	if err != nil {
		return nil, err
	}

	defer CatchWarnErrorFunc(file.Close)

	switch ext {
	case ".json":
		return LoadMultiJSONDocumentStream[T](file)
	case ".yaml", ".yml":
		return LoadMultiYAMLDocumentStream[T](file)
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
func FileReaderFromPath(
	ctx context.Context,
	filePath string,
	options ...DownloadFileOption,
) (io.ReadCloser, string, error) {
	defaultOptions := &downloadFileOptions{
		HTTPClient: http.DefaultClient,
	}

	for _, opt := range options {
		opt(defaultOptions)
	}

	filePath = strings.TrimSpace(filePath)
	if filePath == "" {
		return nil, "", errFilePathRequired
	}

	fileURL, err := ParsePathOrURL(filePath)
	if err != nil {
		return nil, "", err
	}

	if slices.Contains(httpSchemes, strings.ToLower(fileURL.Scheme)) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, filePath, nil)
		if err != nil {
			return nil, "", err
		}

		resp, err := defaultOptions.HTTPClient.Do(req)
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

	filePath = filepath.Clean(fileURL.Path)
	ext := filepath.Ext(filePath)
	reader, err := os.Open(filePath) //nolint:gosec // user knows if the path is safe to read

	return reader, ext, err
}

type downloadFileOptions struct {
	HTTPClient   Doer
	IncludePaths []string
	ExcludePaths []string
}

// DownloadFileOption abstracts a function to configure options for loading files.
type DownloadFileOption func(opts *downloadFileOptions)

// DownloadFileWithHTTPClient creates an option to set a custom HTTP client to load file.
func DownloadFileWithHTTPClient(client Doer) DownloadFileOption {
	return func(opts *downloadFileOptions) {
		if client == nil {
			opts.HTTPClient = http.DefaultClient
		} else {
			opts.HTTPClient = client
		}
	}
}
