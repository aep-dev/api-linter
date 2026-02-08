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

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aep-dev/api-linter/internal"
	"github.com/aep-dev/api-linter/lint"
	lint_v2 "github.com/aep-dev/api-linter/lint/v2"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

type cli struct {
	ConfigPath                string
	FormatType                string
	OutputPath                string
	ExitStatusOnLintFailure   bool
	VersionFlag               bool
	ProtoImportPaths          []string
	ProtoFiles                []string
	ProtoDescPath             []string
	EnabledRules              []string
	DisabledRules             []string
	ListRulesFlag             bool
	DebugFlag                 bool
	IgnoreCommentDisablesFlag bool
}

// ExitForLintFailure indicates that a problem was found during linting.
//
//lint:ignore ST1012 modifying this variable name is a breaking change.
var ExitForLintFailure = errors.New("found problems during linting")

func newCli(args []string) *cli {
	// Define flag variables.
	var cfgFlag string
	var fmtFlag string
	var outFlag string
	var setExitStatusOnLintFailure bool
	var versionFlag bool
	var protoImportFlag []string
	var protoDescFlag []string
	var ruleEnableFlag []string
	var ruleDisableFlag []string
	var listRulesFlag bool
	var debugFlag bool
	var ignoreCommentDisablesFlag bool

	// Register flag variables.
	fs := pflag.NewFlagSet("api-linter", pflag.ExitOnError)
	fs.StringVar(&cfgFlag, "config", "", "The linter config file.")
	fs.StringVar(&fmtFlag, "output-format", "", "The format of the linting results.\nSupported formats include \"yaml\", \"json\",\"github\" and \"summary\" table.\nYAML is the default.")
	fs.StringVarP(&outFlag, "output-path", "o", "", "The output file path.\nIf not given, the linting results will be printed out to STDOUT.")
	fs.BoolVar(&setExitStatusOnLintFailure, "set-exit-status", false, "Return exit status 1 when lint errors are found.")
	fs.BoolVar(&versionFlag, "version", false, "Print version and exit.")
	fs.StringArrayVarP(&protoImportFlag, "proto-path", "I", nil, "The folder for searching proto imports.\nMay be specified multiple times; directories will be searched in order.\nThe current working directory is always used.")
	fs.StringArrayVar(&protoDescFlag, "descriptor-set-in", nil, "The file containing a FileDescriptorSet for searching proto imports.\nMay be specified multiple times.")
	fs.StringArrayVar(&ruleEnableFlag, "enable-rule", nil, "Enable a rule with the given name.\nMay be specified multiple times.")
	fs.StringArrayVar(&ruleDisableFlag, "disable-rule", nil, "Disable a rule with the given name.\nMay be specified multiple times.")
	fs.BoolVar(&listRulesFlag, "list-rules", false, "Print the rules and exit.  Honors the output-format flag.")
	fs.BoolVar(&debugFlag, "debug", false, "Run in debug mode. Panics will print stack.")
	fs.BoolVar(&ignoreCommentDisablesFlag, "ignore-comment-disables", false, "If set to true, disable comments will be ignored.\nThis is helpful when strict enforcement of AIPs are necessary and\nproto definitions should not be able to disable checks.")

	// Parse flags.
	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}

	return &cli{
		ConfigPath:                cfgFlag,
		FormatType:                fmtFlag,
		OutputPath:                outFlag,
		ExitStatusOnLintFailure:   setExitStatusOnLintFailure,
		ProtoImportPaths:          append(protoImportFlag, "."),
		ProtoDescPath:             protoDescFlag,
		EnabledRules:              ruleEnableFlag,
		DisabledRules:             ruleDisableFlag,
		ProtoFiles:                fs.Args(),
		VersionFlag:               versionFlag,
		ListRulesFlag:             listRulesFlag,
		DebugFlag:                 debugFlag,
		IgnoreCommentDisablesFlag: ignoreCommentDisablesFlag,
	}
}

func (c *cli) lint(rulesV1 lint.RuleRegistry, rulesV2 lint_v2.RuleRegistry, configs lint.Configs) error {
	// Print version and exit if asked.
	if c.VersionFlag {
		fmt.Printf("api-linter %s\n", internal.Version)
		return nil
	}

	if c.ListRulesFlag {
		return outputRules(c.FormatType)
	}

	// Pre-check if there are files to lint.
	if len(c.ProtoFiles) == 0 {
		return fmt.Errorf("no file to lint")
	}
	// Read linter config and append it to the default.
	if c.ConfigPath != "" {
		config, err := lint.ReadConfigsFromFile(c.ConfigPath)
		if err != nil {
			return err
		}
		configs = append(configs, config...)
	}
	// Add configs for the enabled rules.
	configs = append(configs, lint.Config{
		EnabledRules: c.EnabledRules,
	})
	// Add configs for the disabled rules.
	configs = append(configs, lint.Config{
		DisabledRules: c.DisabledRules,
	})

	// V1 Pipeline
	resultsV1, err := c.runV1(rulesV1, configs)
	if err != nil {
		return err
	}

	// V2 Pipeline
	var configsV2 lint_v2.Configs
	for _, cfg := range configs {
		configsV2 = append(configsV2, lint_v2.Config{
			IncludedPaths: cfg.IncludedPaths,
			ExcludedPaths: cfg.ExcludedPaths,
			EnabledRules:  cfg.EnabledRules,
			DisabledRules: cfg.DisabledRules,
		})
	}
	resultsV2, err := c.runV2(rulesV2, configsV2)
	if err != nil {
		return fmt.Errorf("V2 pipeline failed: %w", err)
	}

	// Combine results from both pipelines, preserving file order.
	results := combineResponses(resultsV1, resultsV2)

	// Determine the output for writing the results.
	// Stdout is the default output.
	w := os.Stdout
	if c.OutputPath != "" {
		var err error
		w, err = os.Create(c.OutputPath)
		if err != nil {
			return err
		}
		defer w.Close()
	}

	// Determine the format for printing the results.
	// YAML format is the default.
	marshal := getOutputFormatFunc(c.FormatType)

	// Print the results.
	b, err := marshal(results)
	if err != nil {
		return err
	}
	if _, err = w.Write(b); err != nil {
		return err
	}

	// Return error on lint failure which subsequently
	// exits with a non-zero status code
	if c.ExitStatusOnLintFailure && anyProblems(results) {
		return ExitForLintFailure
	}

	return nil
}

func anyProblems(results []combinedResponse) bool {
	for _, r := range results {
		if r.hasProblems() {
			return true
		}
	}
	return false
}

// combinedResponse holds lint results from both V1 and V2 pipelines for a
// single file. It preserves the concrete problem types so that custom
// MarshalJSON/MarshalYAML methods on each Problem type are invoked correctly.
type combinedResponse struct {
	FilePath   string
	ProblemsV1 []lint.Problem
	ProblemsV2 []lint_v2.Problem
}

// MarshalJSON produces output compatible with the original lint.Response format.
func (r combinedResponse) MarshalJSON() ([]byte, error) {
	problems := make([]interface{}, 0, len(r.ProblemsV1)+len(r.ProblemsV2))
	for _, p := range r.ProblemsV1 {
		problems = append(problems, p)
	}
	for _, p := range r.ProblemsV2 {
		problems = append(problems, p)
	}
	return json.Marshal(struct {
		FilePath string        `json:"file_path"`
		Problems []interface{} `json:"problems"`
	}{r.FilePath, problems})
}

// MarshalYAML produces output compatible with the original lint.Response format.
func (r combinedResponse) MarshalYAML() (interface{}, error) {
	problems := make([]interface{}, 0, len(r.ProblemsV1)+len(r.ProblemsV2))
	for _, p := range r.ProblemsV1 {
		problems = append(problems, p)
	}
	for _, p := range r.ProblemsV2 {
		problems = append(problems, p)
	}
	return struct {
		FilePath string        `yaml:"file_path"`
		Problems []interface{} `yaml:"problems"`
	}{r.FilePath, problems}, nil
}

func (r combinedResponse) hasProblems() bool {
	return len(r.ProblemsV1) > 0 || len(r.ProblemsV2) > 0
}

// problemInfo provides a unified view of a problem for format functions
// (github, summary) that need to access problem fields directly.
type problemInfo struct {
	RuleID  string
	Message string
	Span    []int32 // From Location.Span; nil if no location.
	RuleURI string
}

// allProblems returns a unified list of problem info from both V1 and V2 problems.
func (r combinedResponse) allProblems() []problemInfo {
	result := make([]problemInfo, 0, len(r.ProblemsV1)+len(r.ProblemsV2))
	for _, p := range r.ProblemsV1 {
		var span []int32
		if p.Location != nil {
			span = p.Location.Span
		}
		result = append(result, problemInfo{
			RuleID:  string(p.RuleID),
			Message: p.Message,
			Span:    span,
			RuleURI: p.GetRuleURI(),
		})
	}
	for _, p := range r.ProblemsV2 {
		var span []int32
		if p.Location != nil {
			span = p.Location.Span
		}
		result = append(result, problemInfo{
			RuleID:  string(p.RuleID),
			Message: p.Message,
			Span:    span,
			RuleURI: p.GetRuleURI(),
		})
	}
	return result
}

// combineResponses merges V1 and V2 responses by file path, preserving
// the input file ordering (V1 order first, then any V2-only files).
func combineResponses(v1 []lint.Response, v2 []lint_v2.Response) []combinedResponse {
	order := make([]string, 0)
	byFile := make(map[string]*combinedResponse)

	for _, r := range v1 {
		if _, ok := byFile[r.FilePath]; !ok {
			order = append(order, r.FilePath)
			byFile[r.FilePath] = &combinedResponse{FilePath: r.FilePath}
		}
		byFile[r.FilePath].ProblemsV1 = append(byFile[r.FilePath].ProblemsV1, r.Problems...)
	}

	for _, r := range v2 {
		if _, ok := byFile[r.FilePath]; !ok {
			order = append(order, r.FilePath)
			byFile[r.FilePath] = &combinedResponse{FilePath: r.FilePath}
		}
		byFile[r.FilePath].ProblemsV2 = append(byFile[r.FilePath].ProblemsV2, r.Problems...)
	}

	result := make([]combinedResponse, 0, len(order))
	for _, fp := range order {
		result = append(result, *byFile[fp])
	}
	return result
}

var outputFormatFuncs = map[string]formatFunc{
	"yaml": yaml.Marshal,
	"yml":  yaml.Marshal,
	"json": json.Marshal,
	"github": func(i interface{}) ([]byte, error) {
		switch v := i.(type) {
		case []combinedResponse:
			return formatGitHubActionOutput(v), nil
		default:
			return json.Marshal(v)
		}
	},
	"summary": func(i interface{}) ([]byte, error) {
		switch v := i.(type) {
		case []combinedResponse:
			return printSummaryTable(v)
		case listedRules:
			return v.printSummaryTable()
		default:
			return json.Marshal(v)
		}
	},
}

type formatFunc func(interface{}) ([]byte, error)

func getOutputFormatFunc(formatType string) formatFunc {
	if f, found := outputFormatFuncs[strings.ToLower(formatType)]; found {
		return f
	}
	return yaml.Marshal
}
