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

func TestHasStringPrefixFold(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		prefix   string
		expected bool
	}{
		{name: "exact match", input: "Hello", prefix: "Hello", expected: true},
		{name: "case insensitive match", input: "Hello", prefix: "hello", expected: true},
		{name: "uppercase prefix", input: "hello world", prefix: "HELLO", expected: true},
		{name: "mixed case input and prefix", input: "TEXT/html", prefix: "text/", expected: true},
		{name: "no match", input: "hello", prefix: "world", expected: false},
		{name: "prefix longer than input", input: "hi", prefix: "hello", expected: false},
		{name: "empty prefix", input: "hello", prefix: "", expected: true},
		{name: "empty input and prefix", input: "", prefix: "", expected: true},
		{name: "empty input non-empty prefix", input: "", prefix: "x", expected: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := HasStringPrefixFold(tc.input, tc.prefix)
			if got != tc.expected {
				t.Errorf("HasStringPrefixFold(%q, %q) = %v, want %v", tc.input, tc.prefix, got, tc.expected)
			}
		})
	}
}

func TestHasStringSuffixFold(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		suffix   string
		expected bool
	}{
		{name: "exact match", input: "Hello", suffix: "Hello", expected: true},
		{name: "case insensitive match", input: "Hello", suffix: "hello", expected: true},
		{name: "uppercase suffix", input: "say hello", suffix: "HELLO", expected: true},
		{name: "mixed case", input: "application/graphql+JSON", suffix: "+json", expected: true},
		{name: "no match", input: "hello", suffix: "world", expected: false},
		{name: "suffix longer than input", input: "hi", suffix: "hello", expected: false},
		{name: "empty suffix", input: "hello", suffix: "", expected: true},
		{name: "empty input and suffix", input: "", suffix: "", expected: true},
		{name: "empty input non-empty suffix", input: "", suffix: "x", expected: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := HasStringSuffixFold(tc.input, tc.suffix)
			if got != tc.expected {
				t.Errorf("HasStringSuffixFold(%q, %q) = %v, want %v", tc.input, tc.suffix, got, tc.expected)
			}
		})
	}
}

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
