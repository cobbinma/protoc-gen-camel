package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cobbinma/protoc-gen-camel/linter"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
	"gopkg.in/yaml.v3"
)

const DEFAULT_CONFIG_FILE = "camel.yml"

type Config struct {
	Ignore []linter.FullFieldName
}

func main() {
	var flags flag.FlagSet
	generate := flags.Bool("generate", false, "generate a configuration file")
	path := flags.String("config", DEFAULT_CONFIG_FILE, "path for an optional configuration file")
	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

		generated := &Config{Ignore: []linter.FullFieldName{}}

		config, err := ReadConfig(path)
		if err != nil {
			return err
		}

		var violations []linter.FullFieldName
		for _, f := range gen.Files {

			v := linter.LintProtoFile(linter.Config{
				Proto:   f.Proto,
				Ignore:  config.Ignore,
				OutFile: os.Stderr,
			})

			generated.Ignore = append(generated.Ignore, v.AllViolations...)
			violations = append(violations, v.NotIgnored...)
		}

		if generate != nil && *generate {
			if err := WriteConfig(generated); err != nil {
				return err
			}
		}

		if len(violations) > 0 {
			return fmt.Errorf("🐪: %d total", len(violations))
		}

		return nil
	})
}

func ReadConfig(custom *string) (*Config, error) {
	data := &Config{}
	path := DEFAULT_CONFIG_FILE
	if custom != nil && *custom != "" {
		path = *custom
	}

	f, err := os.ReadFile(path)
	if err != nil {
		return data, nil
	}

	if f == nil {
		return data, nil
	}

	if err := yaml.Unmarshal(f, &data); err != nil {
		return data, fmt.Errorf("unable to unmarshal yaml : %w", err)
	}

	return data, nil
}

func WriteConfig(config *Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(DEFAULT_CONFIG_FILE, data, 0o666)
}
