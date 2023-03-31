# protoc-gen-camel

A plug-in for Google's [Protocol Buffers](https://github.com/google/protobuf)
compiler to check camel case field name violations.

## Installation
```sh
go install github.com/cobbinma/protoc-gen-camel@latest
```

## Usage
```sh
❯ protoc --camel_out=. $(find . -name '*.proto')
example/camel.proto:Field name "one_two" should be camelCase, such as "oneTwo".
--camel_out: 🐪: 1 total
```

### Configuration

```yaml
ignore:
    - camel.Foo.one_two
```

generate a configuration file with all violations ignored
```sh
protoc --camel_out=generate=true:. $(find . -name '*.proto')
```

use the configuration file
```sh
protoc --camel_out=config=camel.yml:. $(find . -name '*.proto')
```

### Buf

```yml
version: v1
plugins:
  - plugin: camel
    out: .
    strategy: all
    opt:
      - config=camel.yml
```
