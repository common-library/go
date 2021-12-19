# common-library-go

## Installation
```bash
go get -u github.com/heaven-chp/common-library-go
```

## Features
 - db interface
   - elasticsearch
   - mongodb
   - redis
 - file interface
 - grpc interface
 - json interface
 - log interface
 - socket interface

## Test
```bash
go test -cover github.com/heaven-chp/common-library-go/...
```

## Coverage
```bash
go test -cover -coverprofile=coverage.out github.com/heaven-chp/common-library-go/...
go tool cover -html=./coverage.out -o ./coverage.html
```

## How to add grpc
 - create protobuf IDL(Interface Definition Language) file
   - see grpc/sample/sample.proto
 - convert IDL file to code
   - wget https://github.com/protocolbuffers/protobuf/releases/download/v3.13.0/protoc-3.13.0-linux-x86_64.zip
   - unzip protoc-3.12.3-linux-x86_64.zip -d bin/protoc
   - go get -u google.golang.org/grpc
   - go get -u github.com/golang/protobuf/protoc-gen-go
   - ./bin/protoc/bin/protoc src/github.com/heaven-chp/common-library-go/grpc/sample_server/sample.proto --go_out=plugins=grpc:src/github.com/heaven-chp/common-library-go/grpc/sample_server/ --plugin=bin/protoc-gen-go
  - implement functions defined in IDL file
    - implement to satisfy serverDetail interface
    - see grpc/sample/sample.go
