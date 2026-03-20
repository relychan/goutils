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
	"testing"
)

func TestQuoteBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []byte
	}{
		{
			name:     "simple string",
			input:    "hello",
			expected: []byte(`"hello"`),
		},
		{
			name:     "empty string",
			input:    "",
			expected: []byte(`""`),
		},
		{
			name:     "string with spaces",
			input:    "hello world",
			expected: []byte(`"hello world"`),
		},
		{
			name:     "string with special characters",
			input:    "foo\nbar",
			expected: []byte("\"foo\nbar\""),
		},
		{
			name:     "string already containing quotes",
			input:    `say "hi"`,
			expected: []byte(`"say "hi""`),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name+"/string", func(t *testing.T) {
			got := QuoteBytes(tc.input)
			if string(got) != string(tc.expected) {
				t.Fatalf("QuoteBytes(%q) = %q, want %q", tc.input, got, tc.expected)
			}
		})

		t.Run(tc.name+"/bytes", func(t *testing.T) {
			got := QuoteBytes([]byte(tc.input))
			if string(got) != string(tc.expected) {
				t.Fatalf("QuoteBytes(%q) = %q, want %q", tc.input, got, tc.expected)
			}
		})
	}
}
