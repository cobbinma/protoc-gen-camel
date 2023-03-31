package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/iancoleman/strcase"
	"github.com/samber/lo"
	"google.golang.org/protobuf/compiler/protogen"
	"gopkg.in/yaml.v3"
)

const DEFAULT_CONFIG_FILE = "camel.yml"

type Config struct {
	Ignore []string
}

func main() {
	var flags flag.FlagSet
	generate := flags.Bool("generate", false, "generate a configuration file")
	configuration := flags.String("config", DEFAULT_CONFIG_FILE, "path for an optional configuration file")
	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		generated := &Config{Ignore: []string{}}
		config, err := func(conf *string) (Config, error) {
			data := Config{}
			path := DEFAULT_CONFIG_FILE
			if configuration != nil && *configuration != "" {
				path = *configuration
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
		}(configuration)
		if err != nil {
			return err
		}
		var names []string
		for _, f := range gen.Files {
			for _, m := range f.Messages {
				for _, field := range m.Fields {
					name := string(field.Desc.Name())
					camel := strcase.ToLowerCamel(name)
					full := string(field.Desc.FullName())
					if name == camel {
						continue
					}

					generated.Ignore = append(generated.Ignore, full)

					if name != camel && !lo.Contains(config.Ignore, full) {
						fmt.Fprintf(os.Stderr, "%s:Field name \"%s\" should be camelCase, such as \"%s\".\n", *f.Proto.Name, name, camel)
						names = append(names, full)
					}
				}
			}
		}

		if generate != nil && *generate {
			data, err := yaml.Marshal(generated)
			if err != nil {
				return err
			}

			if err := os.WriteFile(DEFAULT_CONFIG_FILE, data, 0o666); err != nil {
				return err
			}
		}

		if len(names) > 0 {
			return fmt.Errorf("🐪: %d total", len(names))
		}

		return nil
	})
}
