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

func TestGetSortedKeys(t *testing.T) {
	t.Run("string keys", func(t *testing.T) {
		input := map[string]int{"banana": 2, "apple": 1, "cherry": 3}
		got := GetSortedKeys(input)
		want := []string{"apple", "banana", "cherry"}
		if len(got) != len(want) {
			t.Fatalf("got %v, want %v", got, want)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("got[%d] = %q, want %q", i, got[i], want[i])
			}
		}
	})

	t.Run("int keys", func(t *testing.T) {
		input := map[int]string{3: "c", 1: "a", 2: "b"}
		got := GetSortedKeys(input)
		want := []int{1, 2, 3}
		if len(got) != len(want) {
			t.Fatalf("got %v, want %v", got, want)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("got[%d] = %d, want %d", i, got[i], want[i])
			}
		}
	})

	t.Run("empty map", func(t *testing.T) {
		input := map[string]int{}
		got := GetSortedKeys(input)
		if len(got) != 0 {
			t.Errorf("expected empty slice, got %v", got)
		}
	})

	t.Run("single entry", func(t *testing.T) {
		input := map[string]bool{"only": true}
		got := GetSortedKeys(input)
		if len(got) != 1 || got[0] != "only" {
			t.Errorf("got %v, want [only]", got)
		}
	})

	t.Run("duplicate-value map with unique keys", func(t *testing.T) {
		input := map[string]int{"z": 1, "a": 1, "m": 1}
		got := GetSortedKeys(input)
		want := []string{"a", "m", "z"}
		if len(got) != len(want) {
			t.Fatalf("got %v, want %v", got, want)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("got[%d] = %q, want %q", i, got[i], want[i])
			}
		}
	})
}
