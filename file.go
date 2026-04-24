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
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/relychan/goutils/httperror"
	"go.yaml.in/yaml/v4"
)

var (
	errFilePathRequired             = errors.New("file path is required")
	errDisallowedFilePath           = errors.New("file path is not allowed to read")
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
		// Ensure the reader read all data.
		_, _ = io.Copy(io.Discard, file) //nolint:errcheck

		return &result, err
	case ".yaml", ".yml":
		var result T

		loader, err := yaml.NewLoader(file)
		if err != nil {
			return nil, err
		}

		err = loader.Load(&result)
		if err != nil {
			return nil, err
		}

		// Ensure the reader read all data.
		_, _ = io.Copy(io.Discard, file) //nolint:errcheck

		return &result, nil
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
		if opt == nil {
			continue
		}

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

	if slices.ContainsFunc(httpSchemes, func(scheme string) bool {
		return strings.EqualFold(fileURL.Scheme, scheme)
	}) {
		return fileReaderFromURL(ctx, fileURL, filePath, defaultOptions)
	}

	filePath = filepath.Clean(filePath)

	err = validateFilePath(filePath, filepath.Match, defaultOptions)
	if err != nil {
		return nil, "", err
	}

	ext := filepath.Ext(filePath)
	reader, err := os.Open(filePath)

	return reader, strings.ToLower(ext), err
}

func fileReaderFromURL(
	ctx context.Context,
	fileURL *url.URL,
	filePath string,
	options *downloadFileOptions,
) (io.ReadCloser, string, error) {
	if len(options.AllowedHosts) > 0 || len(options.BlockedHosts) > 0 {
		err := validateHost(fileURL.Host, fileURL.Hostname(), &ValidateHTTPURLOptions{
			AllowedHosts: options.AllowedHosts,
			BlockedHosts: options.BlockedHosts,
		})
		if err != nil {
			return nil, "", err
		}
	}

	if fileURL.Path != "" && fileURL.Path[0] != '/' {
		fileURL.Path = "/" + fileURL.Path
	}

	err := validateFilePath(fileURL.Path, path.Match, options)
	if err != nil {
		return nil, "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, filePath, nil)
	if err != nil {
		return nil, "", err
	}

	resp, err := options.HTTPClient.Do(req)
	if err != nil {
		return nil, "", err
	}

	if resp.StatusCode != http.StatusOK {
		respError := httperror.NewHTTPErrorFromResponse(resp)
		respError.Title = "Read File Failure"

		return nil, "", respError
	}

	if resp.Body == nil || resp.Body == http.NoBody {
		return nil, "", errFileNoContent
	}

	ext := strings.ToLower(filepath.Ext(filePath))

	return resp.Body, ext, nil
}

func validateFilePath(
	filePath string,
	matchFunc func(pattern string, name string) (matched bool, err error),
	options *downloadFileOptions,
) error {
	isMatched := len(options.IncludePaths) == 0

	for _, includedPath := range options.IncludePaths {
		matched, err := matchFunc(includedPath, filePath)
		if err != nil {
			return err
		}

		if matched {
			isMatched = true

			break
		}
	}

	if !isMatched {
		return errDisallowedFilePath
	}

	for _, excludedPath := range options.ExcludePaths {
		matched, err := filepath.Match(excludedPath, filePath)
		if err != nil {
			return err
		}

		if matched {
			return errDisallowedFilePath
		}
	}

	return nil
}

type downloadFileOptions struct {
	HTTPClient   Doer
	IncludePaths []string
	ExcludePaths []string
	AllowedHosts []string
	BlockedHosts []string
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

// DownloadFileIncludingPaths creates an option to set a list of paths to be included.
func DownloadFileIncludingPaths(paths []string) DownloadFileOption {
	return func(opts *downloadFileOptions) {
		opts.IncludePaths = paths
	}
}

// DownloadFileExcludingPaths creates an option to set a list of paths to be excluded.
func DownloadFileExcludingPaths(paths []string) DownloadFileOption {
	return func(opts *downloadFileOptions) {
		opts.ExcludePaths = paths
	}
}

// DownloadFileWithAllowedHosts creates an option to set a list of allowed hosts for URL.
func DownloadFileWithAllowedHosts(hosts []string) DownloadFileOption {
	return func(opts *downloadFileOptions) {
		opts.AllowedHosts = hosts
	}
}

// DownloadFileWithBlockedHosts creates an option to set a list of blocked hosts for URL.
func DownloadFileWithBlockedHosts(hosts []string) DownloadFileOption {
	return func(opts *downloadFileOptions) {
		opts.BlockedHosts = hosts
	}
}
