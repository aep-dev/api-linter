package main

import (
	"context"
	"errors"

	"buf.build/go/bufplugin/check"
	"buf.build/go/bufplugin/descriptor"
	"github.com/aep-dev/api-linter/lint"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type fileDescriptorsContextKey struct{}

func newRuleHandlerV1(protoRule lint.ProtoRule) check.RuleHandler {
	return check.RuleHandlerFunc(
		func(ctx context.Context, responseWriter check.ResponseWriter, request check.Request) error {
			fileDescriptors, _ := ctx.Value(fileDescriptorsContextKey{}).([]*desc.FileDescriptor)
			for _, fileDescriptor := range fileDescriptors {
				for _, problem := range protoRule.Lint(fileDescriptor) {
					if err := addProblemV1(responseWriter, problem); err != nil {
						return err
					}
				}
			}
			return nil
		},
	)
}

func addProblemV1(responseWriter check.ResponseWriter, problem lint.Problem) error {
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
	fileDescriptors, err := nonImportFileDescriptorsForFileDescriptors(request.FileDescriptors())
	if err != nil {
		return nil, nil, err
	}
	ctx = context.WithValue(ctx, fileDescriptorsContextKey{}, fileDescriptors)
	return ctx, request, nil
}

func nonImportFileDescriptorsForFileDescriptors(fileDescriptors []descriptor.FileDescriptor) ([]*desc.FileDescriptor, error) {
	if len(fileDescriptors) == 0 {
		return nil, nil
	}
	reflectFileDescriptors := make([]protoreflect.FileDescriptor, 0, len(fileDescriptors))
	for _, fileDescriptor := range fileDescriptors {
		if fileDescriptor.IsImport() {
			continue
		}
		reflectFileDescriptors = append(reflectFileDescriptors, fileDescriptor.ProtoreflectFileDescriptor())
	}
	return desc.WrapFiles(reflectFileDescriptors)
}
