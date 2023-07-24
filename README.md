# common-library-go

## Installation
```bash
go get -u github.com/heaven-chp/common-library-go
```

<br/>

## Features
 - command-line-argument interface
 - db interface
   - dynamodb
   - elasticsearch
   - mongodb
   - mysql
   - redis
 - file interface
 - grpc interface
 - json interface
 - log interface
 - socket interface
 - utility interface

<br/>

## Test and Coverage
 - prepare
   - dynamodb
     - `docker run --name dynamodb -d -p 8000:8000 -e "-jar DynamoDBLocal.jar -sharedDb -inMemory" amazon/dynamodb-local:2.0.0`
   - Elasticsearch v7
     - `docker run --name elasticsearch-v7 -d -p 19200:9200 -p 19300:9300 -e discovery.type=single-node -e ES_JAVA_OPTS="-Xms500m -Xmx500m" elasticsearch:7.17.10`
   - Elasticsearch v8
     - `docker network create elastic`
     - `docker run --name elasticsearch-v8 --net elastic -d -p 29200:9200 -p 29300:9300 -e discovery.type=single-node -e ES_JAVA_OPTS="-Xms500m -Xmx500m" -e xpack.security.enabled=false elasticsearch:8.8.1`
   - MongoDB
     - `docker run --name mongodb -d -p 27017:27017 mongo:6.0.8`
   - MySQL
     - `docker run --name mysql -d -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root mysql:8.0.33`
   - Redis
     - `docker run --name redis -d -p 6379:6379 redis:7.0.12`
 - Test
   - `go clean -testcache && go test -cover ./...`
 - Coverage
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
