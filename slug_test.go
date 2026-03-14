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
	"strings"
	"testing"

	"go.yaml.in/yaml/v4"
)

func TestSlug(t *testing.T) {
	t.Run("unmarshal_json", func(t *testing.T) {
		var slug Slug

		err := json.Unmarshal([]byte(`"hello-world"`), &slug)
		if err != nil {
			t.Fatalf("expected nil error, got: %s", err)
		}

		if slug.String() != "hello-world" {
			t.Fatalf("slug not equal, expected hello-world, got: %s", slug)
		}

		err = json.Unmarshal([]byte(`"hello world"`), &slug)
		if !strings.Contains(err.Error(), "invalid slug") {
			t.Fatalf("expected invalid slug error, got: %s", err)
		}

		err = json.Unmarshal([]byte(`1`), &slug)
		if !strings.Contains(err.Error(), "cannot unmarshal number into Go value of type string") {
			t.Fatalf("expected unmarshal error, got: %s", err)
		}
	})

	t.Run("unmarshal_yaml", func(t *testing.T) {
		var slug Slug

		err := yaml.Unmarshal([]byte(`"hello-world"`), &slug)
		if err != nil {
			t.Fatalf("expected nil error, got: %s", err)
		}

		if slug.String() != "hello-world" {
			t.Fatalf("slug not equal, expected hello-world, got: %s", slug)
		}

		err = yaml.Unmarshal([]byte(`"hello world"`), &slug)
		if !strings.Contains(err.Error(), "invalid slug") {
			t.Fatalf("expected invalid slug error, got: %s", err)
		}

		err = yaml.Load([]byte(`{}`), &slug)
		if !strings.Contains(err.Error(), "cannot construct !!map into goutils.Slug") {
			t.Fatalf("expected unmarshal error, got: %s", err)
		}
	})
}
