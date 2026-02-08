package main

import (
	"context"
	"errors"

	"buf.build/go/bufplugin/check"
	lint_v2 "github.com/aep-dev/api-linter/lint/v2"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func newRuleHandlerV2(protoRule lint_v2.ProtoRule) check.RuleHandler {
	return check.RuleHandlerFunc(
		func(ctx context.Context, responseWriter check.ResponseWriter, request check.Request) error {
			for _, fileDescriptor := range request.FileDescriptors() {
				if fileDescriptor.IsImport() {
					continue
				}
				// bufplugin provides protoreflect descriptors natively!
				for _, problem := range protoRule.Lint(fileDescriptor.ProtoreflectFileDescriptor()) {
					if err := addProblemV2(responseWriter, problem); err != nil {
						return err
					}
				}
			}
			return nil
		},
	)
}

func addProblemV2(responseWriter check.ResponseWriter, problem lint_v2.Problem) error {
	addAnnotationOptions := []check.AddAnnotationOption{
		check.WithMessage(problem.Message),
	}
	descriptor := problem.Descriptor
	if descriptor == nil {
		return errors.New("got nil problem.Descriptor")
	}
	fileDescriptor := descriptor.ParentFile()
	if fileDescriptor == nil {
		responseWriter.AddAnnotation(addAnnotationOptions...)
		return nil
	}

	if location := problem.Location; location != nil {
		addAnnotationOptions = append(
			addAnnotationOptions,
			check.WithFileNameAndSourcePath(
				fileDescriptor.Path(),
				protoreflect.SourcePath(location.GetPath()),
			),
		)
	} else {
		// Use V2 SourceLocations
		loc := fileDescriptor.SourceLocations().ByDescriptor(descriptor)
		addAnnotationOptions = append(
			addAnnotationOptions,
			check.WithFileNameAndSourcePath(
				fileDescriptor.Path(),
				protoreflect.SourcePath(loc.Path),
			),
		)
	}
	responseWriter.AddAnnotation(addAnnotationOptions...)
	return nil
}
