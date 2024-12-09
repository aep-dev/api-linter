// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

// ToKebobCase returns the kebob-case of a word (book-edition).
func ToKebobCase(s string) string {
	asLower := make([]rune, 0, len(s))
	for i, r := range s {
		if isUpper(r) {
			r = r | ' ' // make lowercase

			// Only insert hypen after first word.
			if i != 0 {
				asLower = append(asLower, '-')
			}
			asLower = append(asLower, r)
		} else if r == '-' || r == '_' || r == ' ' || r == '.' {
			asLower = append(asLower, '-')
		} else {
			asLower = append(asLower, r)
		}
	}
	return string(asLower)
}

func isUpper(r rune) bool {
	return ('A' <= r && r <= 'Z')
}
