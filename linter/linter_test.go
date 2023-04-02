package linter_test

import (
	"io"
	"testing"

	"github.com/cobbinma/protoc-gen-camel/linter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestLintProtoFile(t *testing.T) {
	descriptor := &descriptorpb.FileDescriptorProto{}
	require.NoError(t, prototext.Unmarshal([]byte(`name:  "example/camel.proto"
	package:  "camel"
	message_type:  {
		name:  "Foo"
		field:  {
			name:  "one_two"
			number:  1
			label:  LABEL_OPTIONAL
			type:  TYPE_INT64
			json_name:  "oneTwo"
		}
		field:  {
			name:  "twoThree"
			number:  2
			label:  LABEL_OPTIONAL
			type:  TYPE_INT64
			json_name:  "twoThree"
		}
	}
	options:  {
		go_package:  "github.com/cobbinma/protoc-gen-camel/example"
	}
	syntax:  "proto3"`), descriptor))

	t.Run("violation without ignore", func(t *testing.T) {
		assert.Equal(t, linter.Violations{
			AllViolations: []linter.FullFieldName{
				"camel.Foo.one_two",
			},
			NotIgnored: []linter.FullFieldName{
				"camel.Foo.one_two",
			},
		}, linter.LintProtoFile(linter.Config{
			Proto:   descriptor,
			OutFile: io.Discard,
			Ignore:  []linter.FullFieldName{},
		}))
	})

	t.Run("ignored violation", func(t *testing.T) {
		assert.Equal(t, linter.Violations{
			AllViolations: []linter.FullFieldName{
				"camel.Foo.one_two",
			},
			NotIgnored: nil,
		}, linter.LintProtoFile(linter.Config{
			Proto:   descriptor,
			OutFile: io.Discard,
			Ignore: []linter.FullFieldName{
				"camel.Foo.one_two",
			},
		}))
	})
}
