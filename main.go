package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/iancoleman/strcase"
	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	var flags flag.FlagSet
	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		var errors int
		for _, f := range gen.Files {
			for _, m := range f.Messages {
				for _, field := range m.Fields {
					name := string(field.Desc.Name())
					camel := strcase.ToLowerCamel(name)
					if name != camel {
						fmt.Fprintf(os.Stderr, "%s:Field name \"%s\" should be camelCase, such as \"%s\".\n", *f.Proto.Name, name, camel)
						errors++
					}
				}
			}
		}
		if errors > 0 {
			return fmt.Errorf("🐪: %d total", errors)
		}
		return nil
	})
}
