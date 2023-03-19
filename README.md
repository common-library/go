# common-library-go

## Installation
```bash
go get -u github.com/heaven-chp/common-library-go
```

<br/>

## Features
 - db interface
   - elasticsearch
   - mongodb
   - mysql
   - redis
 - file interface
 - grpc interface
 - json interface
 - log interface
 - socket interface

<br/>

## Test
 - prepare
   - Elasticsearch
     - `docker run --name elasticsearch -d -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" docker.elastic.co/elasticsearch/elasticsearch:7.17.9`
   - MongoDB
     - `docker run --name mongodb -d -p 27017:27017 mongo:6.0.5`
   - Redis
     - `docker run --name redis -d -p 6379:6379 redis:7`
   - MySQL
     - `docker run --name mysql -e MYSQL_ROOT_PASSWORD=root -d -p 3306:3306 mysql:8.0.32`
 - run
   - `go clean -testcache && go test -cover ./...`

<br/>

## Coverage
 - make coverage file
   - `go clean -testcache && go test -cover -coverprofile=coverage.out ./...`
 - convert coverage file to html file
   - `go tool cover -html=./coverage.out -o ./coverage.html`

<br/>

## How to add grpc
 - create protobuf IDL(Interface Definition Language) file
   - see [grpc/sample/sample.proto](https://github.com/heaven-chp/common-library-go/blob/main/grpc/sample/sample.proto)
 - convert IDL file to code
   - `wget https://github.com/protocolbuffers/protobuf/releases/download/v3.13.0/protoc-3.13.0-linux-x86_64.zip`
   - `unzip protoc-3.12.3-linux-x86_64.zip -d bin/protoc`
   - `go get -u google.golang.org/grpc`
   - `go get -u github.com/golang/protobuf/protoc-gen-go`
   - `./bin/protoc/bin/protoc src/github.com/heaven-chp/common-library-go/grpc/sample_server/sample.proto --go_out=plugins=grpc:src/github.com/heaven-chp/common-library-go/grpc/sample_server/ --plugin=bin/protoc-gen-go`
  - implement functions defined in IDL file
    - implement to satisfy serverDetail interface
    - see [grpc/sample/sample.go](https://github.com/heaven-chp/common-library-go/blob/main/grpc/sample/sample.go)
