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
