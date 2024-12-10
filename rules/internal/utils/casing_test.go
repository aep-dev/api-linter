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

import "testing"

func TestToUpperCamelCase(t *testing.T) {
	for _, test := range []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "OneWord",
			input: "foo",
			want:  "foo",
		},
		{
			name:  "OneWordNoop",
			input: "Foo",
			want:  "foo",
		},
		{
			name:  "TwoWords",
			input: "bookShelf",
			want:  "book-shelf",
		},
		{
			name:  "WithDash",
			input: "book-shelf",
			want:  "book-shelf",
		},
		{
			name:  "WithNumbers",
			input: "universe42love",
			want:  "universe42love",
		},
		{
			name:  "WithUnderscore",
			input: "Book_shelf",
			want:  "book-shelf",
		},
		{
			name:  "WithUnderscore",
			input: "Book_shelf",
			want:  "book-shelf",
		},
		{
			name:  "WithSpaces",
			input: "Book shelf",
			want:  "book-shelf",
		},
		{
			name:  "WithPeriods",
			input: "book.shelf",
			want:  "book-shelf",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got := ToKebabCase(test.input)
			if got != test.want {
				t.Errorf("ToKebabCase(%q) = %q, got %q", test.input, test.want, got)
			}
		})
	}
}
