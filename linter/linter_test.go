package linter_test

import (
	"io"
	"os"
	"testing"

	"github.com/cobbinma/protoc-gen-camel/linter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestLintProtoFile(t *testing.T) {
	data, err := os.ReadFile("test/test.textproto")
	require.NoError(t, err)

	descriptor := &descriptorpb.FileDescriptorProto{}
	require.NoError(t, prototext.Unmarshal(data, descriptor))

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
