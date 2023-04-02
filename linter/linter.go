package linter

import (
	"fmt"
	"io"
	"os"

	"github.com/iancoleman/strcase"
	"github.com/samber/lo"
	"google.golang.org/protobuf/compiler/protogen"
)

type FullFieldName string

type Config struct {
	FileName string
	Ignore   []FullFieldName
	Messages []*protogen.Message
	OutFile  io.Writer
}

type Violations struct {
	AllViolations []FullFieldName
	NotIgnored    []FullFieldName
}

// LintProtoFile takes a file name, proto file description, and a file.
// It checks the file for errors and writes them to the output file.
func LintProtoFile(config Config) Violations {
	violations := Violations{}
	for _, m := range config.Messages {
		for _, field := range m.Fields {
			name := string(field.Desc.Name())
			camel := strcase.ToLowerCamel(name)
			full := FullFieldName(field.Desc.FullName())
			if name == camel {
				continue
			}

			violations.AllViolations = append(violations.AllViolations, full)

			if !lo.Contains(config.Ignore, full) {
				fmt.Fprintf(os.Stderr, "%s:Field name \"%s\" should be camelCase, such as \"%s\".\n", config.FileName, name, camel)

				violations.NotIgnored = append(violations.NotIgnored, full)
			}
		}
	}

	return violations
}
