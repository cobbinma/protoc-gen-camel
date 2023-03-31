# protoc-gen-camel

A plug-in for Google's [Protocol Buffers](https://github.com/google/protobuf)
compiler to check camel case field name violations.

## Installation
```sh
go install github.com/cobbinma/protoc-gen-camel
```

## Usage
```sh
protoc --camel_out=. example/*.proto
```

``sh
example/camel.proto:Field name "one_two" should be camelCase, such as "oneTwo".
--camel_out: 🐪: 1 total
``
