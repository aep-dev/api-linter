// Copyright 2022 Google LLC
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

package aep0127

import (
	"testing"

	"github.com/aep-dev/api-linter/rules/internal/testutils"
)

func TestHttpTemplatePattern_PatternMatching(t *testing.T) {
	tests := []struct {
		name            string
		HTTPAnnotation  string
		ResourcePattern string
		problems        testutils.Problems
	}{
		// HTTP variable uses literals
		{"LiteralMatchesSameLiteralInPattern", "/v1/{path=shelves}", "shelves", nil},
		{"LiteralDoesNotMatchDifferentLiteral", "/v1/{path=shelves}", "books", testutils.Problems{{Message: "does not match"}}},
		{"SuffixAfterHttpVariableIgnoredForMatch", "/v1/{path=shelves}/books", "shelves", nil},

		// HTTP variable uses single wildcard
		{"SingleWildcardMatchesAnyLiteralSegment", "/v1/{path=*}", "shelves", nil},
		{"SingleWildcardMatchesAnyVariableSegment", "/v1/{path=*}", "{shelf}", nil},
		{"SingleWildcardDoesNotMatchMultipleUrlSegments", "/v1/{path=*}", "shelves/{shelf}", testutils.Problems{{Message: "does not match"}}},
		{"LiteralAndWildcardMatch", "/v1/{path=shelves/*}", "shelves/{shelf}", nil},
		// This case is only theoretical, as "{shelf}" represents a resource ID,
		// rather than a resource path. It should not be observed in practice.
		{"ImplicitWildcardMatches", "/v1/{path}/books", "{shelf}", nil},
		{"MulitpleWildcardsMatches", "/v1/{path=shelves/*/books/*}", "shelves/{shelf}/books/{book}", nil},

		// HTTP variable uses double wildcard
		{"DoubleWildcardMatchesZeroSegments", "/v1/{path=**}", "", nil},
		{"DoubleWildcardMatchesOneLiteralSegment", "/v1/{path=**}", "shelves", nil},
		{"DoubleWildcardMatchesMultipleLiteralSegments", "/v1/{path=**}", "my/shelves", nil},
		{"DoubleWildcardMatchesOneVariableSegment", "/v1/{path=**}", "{shelf}", nil},
		{"DoubleWildcardMatchesMultipleVariableSegments", "/v1/{path=**}", "{shelf}/{book}", nil},
		{"DoubleWildcardMatchesMixedLiteralVariableSegments", "/v1/{path=**}", "shelves/{shelf}", nil},
		{"DoubleWildcardPrecededByLiteralMatches", "/v1/{path=shelves/**}", "shelves/{shelf}", nil},
		{"DoubleWildcardFollowedByLiteralMatches", "/v1/{path=**/shelves/*}", "my/shelves/{shelf}", nil},
		{"DoubleWildcardMissingLiteralDoesNotMatch", "/v1/{path=shelves/**}", "{shelf}", testutils.Problems{{Message: "does not match"}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/annotations.proto";
				import "google/api/resource.proto";
				service Library {
					rpc GetBook(GetBookRequest) returns (Book) {
						option (google.api.http) = {
							get: "{{.HTTPAnnotation}}"
						};
					}
				}
				message GetBookRequest {
					// Format: shelves/{shelf}/books/{book}
					string path = 1 [(google.api.resource_reference).type = "library.googleapis.com/Book"];
				}
				message Book {
					option (google.api.resource) = {
						type: "library.googleapis.com/Book"
						pattern: "{{.ResourcePattern}}"
					};
					string path = 1;
				}
			`, test)
			m := f.GetServices()[0].GetMethods()[0]
			if diff := test.problems.SetDescriptor(m).Diff(httpTemplatePattern.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestHttpTemplatePattern_MultiplePatterns(t *testing.T) {
	tests := []struct {
		name             string
		HTTPAnnotation   string
		ResourcePattern1 string
		ResourcePattern2 string
		problems         testutils.Problems
	}{
		{"MatchesIfFirstPatternMatches", "/v1/{path=shelves}", "shelves", "books", nil},
		{"MatchesIfSecondPatternMatches", "/v1/{path=shelves}", "books", "shelves", nil},
		{"FailsIfNeitherPatternMatches", "/v1/{path=shelves}", "books", "bins", testutils.Problems{{Message: "does not match"}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := testutils.ParseProto3Tmpl(t, `
				import "google/api/annotations.proto";
				import "google/api/resource.proto";
				service Library {
					rpc GetBook(GetBookRequest) returns (Book) {
						option (google.api.http) = {
							get: "{{.HTTPAnnotation}}"
						};
					}
				}
				message GetBookRequest {
					string path = 1 [(google.api.resource_reference).type = "library.googleapis.com/Book"];
				}
				message Book {
					option (google.api.resource) = {
						type: "library.googleapis.com/Book"
						pattern: "{{.ResourcePattern1}}"
						pattern: "{{.ResourcePattern2}}"
					};
					string path = 1;
				}
			`, test)
			m := f.GetServices()[0].GetMethods()[0]
			if diff := test.problems.SetDescriptor(m).Diff(httpTemplatePattern.Lint(f)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestHttpTemplatePattern_SkipCheckIfNoHTTPRules(t *testing.T) {
	f := testutils.ParseProto3String(t, `
			import "google/api/annotations.proto";
			import "google/api/resource.proto";
			service Library {
				rpc GetBook(GetBookRequest) returns (Book) {}
			}
			message GetBookRequest {
				string path = 1 [(google.api.resource_reference).type = "library.googleapis.com/Book"];
			}
			message Book {
				option (google.api.resource) = {
					type: "library.googleapis.com/Book"
				};
				string path = 1;
			}
		`)
	if problems := httpTemplatePattern.Lint(f); len(problems) > 0 {
		t.Errorf("%v", problems)
	}
}

func TestHttpTemplatePattern_SkipCheckIfHTTPRuleHasNoVariables(t *testing.T) {
	f := testutils.ParseProto3String(t, `
			import "google/api/annotations.proto";
			import "google/api/resource.proto";
			service Library {
				rpc GetBook(GetBookRequest) returns (Book) {
					option (google.api.http) = {
						get: "/v1/book"
					};
				}
			}
			message GetBookRequest {
				string path = 1 [(google.api.resource_reference).type = "library.googleapis.com/Book"];
			}
			message Book {
				option (google.api.resource) = {
					type: "library.googleapis.com/Book"
				};
				string path = 1;
			}
		`)
	if problems := httpTemplatePattern.Lint(f); len(problems) > 0 {
		t.Errorf("%v", problems)
	}
}

func TestHttpTemplatePattern_SkipCheckIfFieldPathMissingResourceAnnotation(t *testing.T) {
	f := testutils.ParseProto3String(t, `
			import "google/api/annotations.proto";
			import "google/api/resource.proto";
			service Library {
				rpc GetBook(GetBookRequest) returns (Book) {
					option (google.api.http) = {
						get: "/v1/{path=shelves}"
					};
				}
			}
			message GetBookRequest {
				string path = 1;
			}
			message Book {
				option (google.api.resource) = {
					type: "library.googleapis.com/Book"
				};
				string path = 1;
			}
		`)
	if problems := httpTemplatePattern.Lint(f); len(problems) > 0 {
		t.Errorf("%v", problems)
	}
}
