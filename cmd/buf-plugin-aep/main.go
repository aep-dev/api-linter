package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"buf.build/go/bufplugin/check"
	"github.com/aep-dev/api-linter/lint"
	"github.com/aep-dev/api-linter/rules"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	aepCategoryID     = "AEP"
	aepCoreCategoryID = "AEP_CORE"
)

type fileDescriptorsContextKey struct{}

func main() {
	spec, err := newSpec()
	if err != nil {
		log.Fatalln(err)
	}
	// AEP rules cannot be run in parallel as there is thread-unsafe code in
	// this repository that causes concurrent read and write access to a map.
	check.Main(spec, check.MainWithParallelism(1))
}

func newSpec() (*check.Spec, error) {
	ruleRegistry := lint.NewRuleRegistry()
	if err := rules.Add(ruleRegistry); err != nil {
		return nil, err
	}
	ruleSpecs := make([]*check.RuleSpec, 0, len(ruleRegistry))
	for _, protoRule := range ruleRegistry {
		ruleSpec, err := newRuleSpec(protoRule)
		if err != nil {
			return nil, err
		}
		ruleSpecs = append(ruleSpecs, ruleSpec)
	}
	return &check.Spec{
		Rules: ruleSpecs,
		Categories: []*check.CategorySpec{
			{
				ID:      aepCategoryID,
				Purpose: "Checks all API Enhancement proposals as specified at https://aep.dev.",
			},
			{
				ID:      aepCoreCategoryID,
				Purpose: "Checks all core API Improvement proposals as specified at https://aep.dev.",
			},
		},
		Before: before,
	}, nil
}

func newRuleSpec(protoRule lint.ProtoRule) (*check.RuleSpec, error) {
	ruleName := protoRule.GetName()
	if !ruleName.IsValid() {
		return nil, fmt.Errorf("lint.RuleName is invalid: %q", ruleName)
	}

	split := strings.Split(string(ruleName), "::")
	if len(split) != 3 {
		return nil, fmt.Errorf("unknown lint.RuleName format, expected three parts split by '::' : %q", ruleName)
	}
	categoryIDs := []string{aepCategoryID}
	switch extraCategoryID := split[0]; extraCategoryID {
	case "core":
		categoryIDs = append(categoryIDs, aepCoreCategoryID)
	default:
		return nil, fmt.Errorf("unknown lint.RuleName format: unknown category %q : %q", extraCategoryID, ruleName)
	}

	// The allowed characters for RuleName are a-z, 0-9, -.
	// The separator :: is also allowed.
	// We do a translation of these into valid check.Rule IDs.
	ruleID := "AEP_" + strings.Join(split[1:3], "_")
	ruleID = strings.ReplaceAll(ruleID, "-", "_")
	ruleID = strings.ToUpper(ruleID)

	return &check.RuleSpec{
		ID:          ruleID,
		CategoryIDs: categoryIDs,
		Default:     true,
		Purpose:     fmt.Sprintf("Checks AEP rule %s.", ruleName),
		Type:        check.RuleTypeLint,
		Handler:     newRuleHandler(protoRule),
	}, nil
}

func newRuleHandler(protoRule lint.ProtoRule) check.RuleHandler {
	return check.RuleHandlerFunc(
		func(ctx context.Context, responseWriter check.ResponseWriter, request check.Request) error {
			fileDescriptors, _ := ctx.Value(fileDescriptorsContextKey{}).([]*desc.FileDescriptor)
			for _, fileDescriptor := range fileDescriptors {
				for _, problem := range protoRule.Lint(fileDescriptor) {
					if err := addProblem(responseWriter, problem); err != nil {
						return err
					}
				}
			}
			return nil
		},
	)
}

func addProblem(responseWriter check.ResponseWriter, problem lint.Problem) error {
	addAnnotationOptions := []check.AddAnnotationOption{
		check.WithMessage(problem.Message),
	}
	descriptor := problem.Descriptor
	if descriptor == nil {
		// This should never happen.
		return errors.New("got nil problem.Descriptor")
	}
	fileDescriptor := descriptor.GetFile()
	if fileDescriptor == nil {
		// If we do not have a FileDescriptor, we cannot report a location.
		responseWriter.AddAnnotation(addAnnotationOptions...)
		return nil
	}
	// If a location is available from the problem, we use that directly.
	if location := problem.Location; location != nil {
		addAnnotationOptions = append(
			addAnnotationOptions,
			check.WithFileNameAndSourcePath(
				fileDescriptor.GetName(),
				protoreflect.SourcePath(location.GetPath()),
			),
		)
	} else {
		// Otherwise we check the source info for the descriptor from the problem.
		if location := descriptor.GetSourceInfo(); location != nil {
			addAnnotationOptions = append(
				addAnnotationOptions,
				check.WithFileNameAndSourcePath(
					fileDescriptor.GetName(),
					protoreflect.SourcePath(location.GetPath()),
				),
			)
		}
	}
	responseWriter.AddAnnotation(addAnnotationOptions...)
	return nil
}

func before(ctx context.Context, request check.Request) (context.Context, check.Request, error) {
	fileDescriptors, err := nonImportFileDescriptorsForFiles(request.Files())
	if err != nil {
		return nil, nil, err
	}
	ctx = context.WithValue(ctx, fileDescriptorsContextKey{}, fileDescriptors)
	return ctx, request, nil
}

func nonImportFileDescriptorsForFiles(files []check.File) ([]*desc.FileDescriptor, error) {
	if len(files) == 0 {
		return nil, nil
	}
	reflectFileDescriptors := make([]protoreflect.FileDescriptor, 0, len(files))
	for _, file := range files {
		if file.IsImport() {
			continue
		}
		reflectFileDescriptors = append(reflectFileDescriptors, file.FileDescriptor())
	}
	return desc.WrapFiles(reflectFileDescriptors)
}
