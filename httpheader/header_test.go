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

package httpheader

import "testing"

func TestExtractBaseMediaType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "no params", input: "application/json", expected: "application/json"},
		{name: "with charset param", input: "application/json; charset=utf-8", expected: "application/json"},
		{name: "with multiple params", input: "text/html; charset=utf-8; boundary=something", expected: "text/html"},
		{name: "with spaces around semicolon", input: "application/xml ;charset=utf-8", expected: "application/xml"},
		{name: "empty string", input: "", expected: ""},
		{name: "only semicolon", input: ";charset=utf-8", expected: ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ExtractBaseMediaType(tc.input)
			if got != tc.expected {
				t.Errorf("ExtractBaseMediaType(%q) = %q, want %q", tc.input, got, tc.expected)
			}
		})
	}
}

func TestIsContentTypeJSON(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		expected    bool
	}{
		{name: "exact json", contentType: ContentTypeJSON, expected: true},
		{name: "graphql-response+json", contentType: ContentTypeGraphQLResponseJSON, expected: true},
		{name: "json with charset", contentType: "application/json; charset=utf-8", expected: true},
		{name: "uppercase JSON", contentType: "application/JSON", expected: true},
		{name: "custom +json vendor type", contentType: "application/vnd.api+json", expected: true},
		{name: "custom +json with params", contentType: "application/vnd.api+json; charset=utf-8", expected: true},
		{name: "uppercase +JSON suffix", contentType: "application/vnd.api+JSON", expected: true},
		{name: "xml is not json", contentType: ContentTypeXML, expected: false},
		{name: "plain text is not json", contentType: ContentTypeTextPlain, expected: false},
		{name: "empty string", contentType: "", expected: false},
		{name: "json substring in path", contentType: "text/html", expected: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := IsContentTypeJSON(tc.contentType)
			if got != tc.expected {
				t.Errorf("IsContentTypeJSON(%q) = %v, want %v", tc.contentType, got, tc.expected)
			}
		})
	}
}

func TestIsContentTypeXML(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		expected    bool
	}{
		{name: "application/xml", contentType: ContentTypeXML, expected: true},
		{name: "text/xml", contentType: ContentTypeTextXML, expected: true},
		{name: "xml with charset", contentType: "application/xml; charset=utf-8", expected: true},
		{name: "text/xml with charset", contentType: "text/xml; charset=utf-8", expected: true},
		{name: "uppercase XML", contentType: "application/XML", expected: true},
		{name: "custom +xml vendor type", contentType: "application/vnd.foo+xml", expected: true},
		{name: "custom +xml with params", contentType: "application/vnd.foo+xml; charset=utf-8", expected: true},
		{name: "uppercase +XML suffix", contentType: "application/vnd.foo+XML", expected: true},
		{name: "json is not xml", contentType: ContentTypeJSON, expected: false},
		{name: "plain text is not xml", contentType: ContentTypeTextPlain, expected: false},
		{name: "empty string", contentType: "", expected: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := IsContentTypeXML(tc.contentType)
			if got != tc.expected {
				t.Errorf("IsContentTypeXML(%q) = %v, want %v", tc.contentType, got, tc.expected)
			}
		})
	}
}

func TestIsContentTypeText(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		expected    bool
	}{
		{name: "text/plain", contentType: ContentTypeTextPlain, expected: true},
		{name: "text/html", contentType: ContentTypeTextHTML, expected: true},
		{name: "text/xml", contentType: ContentTypeTextXML, expected: true},
		{name: "uppercase TEXT/plain", contentType: "TEXT/plain", expected: true},
		{name: "Text/HTML mixed case", contentType: "Text/HTML", expected: true},
		{name: "application/json is not text", contentType: ContentTypeJSON, expected: false},
		{name: "application/xml is not text", contentType: ContentTypeXML, expected: false},
		{name: "empty string", contentType: "", expected: false},
		{name: "multipart/form-data is not text", contentType: "multipart/form-data", expected: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := IsContentTypeText(tc.contentType)
			if got != tc.expected {
				t.Errorf("IsContentTypeText(%q) = %v, want %v", tc.contentType, got, tc.expected)
			}
		})
	}
}

func TestIsContentTypeMultipartForm(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		expected    bool
	}{
		{name: "multipart/form-data", contentType: "multipart/form-data", expected: true},
		{name: "multipart/mixed", contentType: "multipart/mixed", expected: true},
		{name: "uppercase MULTIPART", contentType: "MULTIPART/form-data", expected: true},
		{name: "Mixed case", contentType: "Multipart/Form-Data", expected: true},
		{name: "application/json is not multipart", contentType: ContentTypeJSON, expected: false},
		{name: "text/plain is not multipart", contentType: ContentTypeTextPlain, expected: false},
		{name: "empty string", contentType: "", expected: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := IsContentTypeMultipartForm(tc.contentType)
			if got != tc.expected {
				t.Errorf("IsContentTypeMultipartForm(%q) = %v, want %v", tc.contentType, got, tc.expected)
			}
		})
	}
}
