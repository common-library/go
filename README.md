# common-library-go

## Installation
```bash
go get -u github.com/heaven-chp/common-library-go
```

<br/>

## Features
 - ai
   - gemini
 - archive
   - gzip
   - tar
   - zip
 - aws
   - Amazon DynamoDB
   - Amazon S3
 - command-line-argument
 - database
   - Elasticsearch v7/v8
   - MongoDB
   - MySQL
   - Redis
 - file
 - grpc
 - http
 - json
 - kubernetes
   - resource
     - client
     - custom-resource
     - custom-resource-definition
 - log
 - long-polling
 - security
   - crypto
     - dsa
     - ecdsa
     - ed25519
     - rsa
 - socket
 - storage
   - MinIO
 - utility

<br/>

## Test and Coverage
 - prepare
   - Amazon DynamoDB
     - `docker run --name dynamodb -d -p 8000:8000 -e "-jar DynamoDBLocal.jar -sharedDb -inMemory" amazon/dynamodb-local:2.2.1`
   - Amazon S3
     - `docker run --name s3mock -d -p 9090:9090 -p 9191:9191 adobe/s3mock:3.4.0`
   - Elasticsearch v7
     - `docker run --name elasticsearch-v7 -d -p 19200:9200 -p 19300:9300 -e discovery.type=single-node -e ES_JAVA_OPTS="-Xms500m -Xmx500m" elasticsearch:7.17.18`
   - Elasticsearch v8
     - `docker network create elastic`
     - `docker run --name elasticsearch-v8 --net elastic -d -p 29200:9200 -p 29300:9300 -e discovery.type=single-node -e ES_JAVA_OPTS="-Xms500m -Xmx500m" -e xpack.security.enabled=false elasticsearch:8.12.1`
   - MongoDB
     - `docker run --name mongodb -d -p 27017:27017 mongo:7.0.5`
   - MySQL
     - `docker run --name mysql -d -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root mysql:8.3.0`
   - Redis
     - `docker run --name redis -d -p 6379:6379 redis:7.2.4`
 - Test
   - `go clean -testcache && go test -cover ./...`
 - Coverage
   - make coverage file
     - `go clean -testcache && go test -coverprofile=coverage.out -cover ./...`
   - convert coverage file to html file
     - `go tool cover -html=./coverage.out -o ./coverage.html`

<br/>

## How to add grpc
 - create protobuf IDL(Interface Definition Language) file
   - see [grpc/sample/sample.proto](https://github.com/heaven-chp/common-library-go/blob/main/grpc/sample/sample.proto)
 - convert IDL file to code
   - `go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0`
   - `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0`
   - `wget https://github.com/protocolbuffers/protobuf/releases/download/v3.20.3/protoc-3.20.3-linux-x86_64.zip`
   - `unzip protoc-3.20.3-linux-x86_64.zip -d protoc/`
   - `protoc/bin/protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative grpc/sample/sample.proto`
  - implement functions defined in IDL file
    - implement to satisfy [implementServer interface](https://github.com/heaven-chp/common-library-go/blob/main/grpc/server.go)
    - see [grpc/sample/Server.go](https://github.com/heaven-chp/common-library-go/blob/main/grpc/sample/Server.go)
    - see [grpc/sample/Server_test.go](https://github.com/heaven-chp/common-library-go/blob/main/grpc/sample/Server_test.go)
