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
	source_code_info:  {
		location:  {
			span:  0
			span:  0
			span:  8
			span:  1
		}
		location:  {
			path:  12
			span:  0
			span:  0
			span:  18
		}
		location:  {
			path:  2
			span:  2
			span:  0
			span:  14
		}
		location:  {
			path:  8
			span:  4
			span:  0
			span:  67
		}
		location:  {
			path:  8
			path:  11
			span:  4
			span:  0
			span:  67
		}
		location:  {
			path:  4
			path:  0
			span:  6
			span:  0
			span:  8
			span:  1
		}
		location:  {
			path:  4
			path:  0
			path:  1
			span:  6
			span:  8
			span:  11
		}
		location:  {
			path:  4
			path:  0
			path:  2
			path:  0
			span:  7
			span:  2
			span:  20
		}
		location:  {
			path:  4
			path:  0
			path:  2
			path:  0
			path:  5
			span:  7
			span:  2
			span:  7
		}
		location:  {
			path:  4
			path:  0
			path:  2
			path:  0
			path:  1
			span:  7
			span:  8
			span:  15
		}
		location:  {
			path:  4
			path:  0
			path:  2
			path:  0
			path:  3
			span:  7
			span:  18
			span:  19
		}
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
