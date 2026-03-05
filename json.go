package goutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// LoadMultiJSONDocumentStream loads multi-document JSON from a reader stream.
func LoadMultiJSONDocumentStream[T any](reader io.Reader) ([]T, error) {
	var results []T

	decoder := json.NewDecoder(reader)

	for {
		var doc T

		err := decoder.Decode(&doc)
		if err == io.EOF || errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return nil, fmt.Errorf(
				"failed to decode multi-documents from JSON at %d: %w",
				len(results),
				err,
			)
		}

		results = append(results, doc)
	}

	return results, nil
}
