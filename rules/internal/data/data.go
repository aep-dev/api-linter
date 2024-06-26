// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package data contains constants used in multiple AEP rules.
package data

import "bitbucket.org/creachadair/stringset"

// Conjunctions is a set of conjunctions.
var Conjunctions = stringset.New("and", "or")

// Prepositions is a set of prepositions.
var Prepositions = stringset.New(
	"after", "at", "before", "between", "but", "by", "except",
	"for", "from", "in", "including", "into", "of", "over", "since", "to",
	"toward", "under", "upon", "with", "within", "without",
)

// ----------------------------------------------------------------------------
// IMPORTANT: Make sure you update docs/_includes/prepositions.md if you
// update the set of prepositions.
// ----------------------------------------------------------------------------
