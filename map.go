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
	"cmp"
	"slices"
)

// GetSortedKeys extract and sort keys of a map.
func GetSortedKeys[K cmp.Ordered, V any](input map[K]V) []K {
	// To minimize allocations, we eschew iterators and pre-size the slice in
	// which we collect v's keys.
	keys := make([]K, len(input))

	var i int

	for k := range input {
		keys[i] = k
		i++
	}

	slices.Sort(keys)

	return keys
}

// ToAnyMap converts a typed map to a any map.
func ToAnyMap[K comparable, V any](input map[K]V) map[K]any {
	result := make(map[K]any, len(input))

	for key, value := range input {
		result[key] = value
	}

	return result
}
