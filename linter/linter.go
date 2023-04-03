package linter

import (
	"fmt"
	"io"

	"github.com/iancoleman/strcase"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/descriptorpb"
)

type FullFieldName string

type Config struct {
	Proto   *descriptorpb.FileDescriptorProto
	Ignore  []FullFieldName
	OutFile io.Writer
}

type Violations struct {
	AllViolations []FullFieldName
	NotIgnored    []FullFieldName
}

// LintProtoFile checks a proto file description for camel case field name violations.
func LintProtoFile(config Config) Violations {
	violations := Violations{}
	for _, message := range config.Proto.GetMessageType() {
		for _, field := range message.GetField() {
			name := string(field.GetName())
			camel := strcase.ToLowerCamel(name)
			full := FullFieldName(fmt.Sprintf("%s.%s.%s", config.Proto.GetPackage(), message.GetName(), name))
			if name == camel {
				continue
			}

			violations.AllViolations = append(violations.AllViolations, full)

			if !lo.Contains(config.Ignore, full) {
				fmt.Fprintf(config.OutFile, "%s:Field name \"%s\" should be camelCase, such as \"%s\".\n", *config.Proto.Name, name, camel)

				violations.NotIgnored = append(violations.NotIgnored, full)
			}
		}
	}

	return violations
}
