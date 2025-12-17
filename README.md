# Common Library for Go

[![CI](https://github.com/common-library/go/workflows/CI/badge.svg)](https://github.com/common-library/go/actions)
[![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/heaven-chp/c7e11ff6ca6c490bd028e4a6d9b79c92/raw/coverage.json)](https://github.com/common-library/go/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/common-library/go)](https://goreportcard.com/report/github.com/common-library/go)
[![Go Version](https://img.shields.io/github/go-mod/go-version/common-library/go?logo=go)](https://github.com/common-library/go)
[![Reference](https://pkg.go.dev/badge/github.com/common-library/go.svg)](https://pkg.go.dev/github.com/common-library/go)
[![License](https://img.shields.io/github/license/common-library/go)](https://github.com/common-library/go/blob/main/LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/common-library/go)](https://github.com/common-library/go/stargazers)

## Installation
```bash
go get -u github.com/common-library/go
```

<br/>

## Features
 - archive
   - gzip
   - tar
   - zip
 - aws
   - Amazon DynamoDB
   - Amazon S3
 - command line
   - arguments
   - flag
 - data structure
   - Deque
   - Queue
 - database
   - dbmate
   - Elasticsearch v7/v8
   - MongoDB
   - Prometheus
   - Redis
   - ORM
     - beego
     - ent
     - GORM
     - sqlc
     - sqlx
   - SQL
     - Amazon DynamoDB
     - ClickHouse
     - Microsoft SQL Server
     - MySQL
     - Oracle
     - Postgres
     - SQLite
 - event
   - cloudevents
 - file
 - grpc
 - http
 - json
 - kubernetes
   - resource
     - client
     - custom resource
     - custom resource definition
 - lock
 - log
 - long polling
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
 - Test
   - `go clean -testcache && go test -cover ./...`
 - Coverage
   - make coverage file
     - `go clean -testcache && go test -coverprofile=coverage.out -cover ./...`
   - convert coverage file to html file
     - `go tool cover -html=./coverage.out -o ./coverage.html`

<br/>

## How to add a ent schema
 - Assuming the schema name is `Xxx`
 - `go get entgo.io/ent/cmd/ent`
 - `go run entgo.io/ent/cmd/ent new --target ./database/orm/ent/schema Xxx`
 - Modify `./database/orm/ent/schema/xxx.go`
 - `go run entgo.io/ent/cmd/ent generate --feature sql/upsert ./database/orm/ent/schema`

<br/>

## How to use sqlc
 - Add or modify query file to `./database/orm/sqlc/queries`
 - Add or modify schema file to `./database/orm/sqlc/schema`
 - Modify `./database/orm/sqlc/sqlc.json`
 - `go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0`
 - `sqlc generate --file ./database/orm/sqlc/sqlc.json`

<br/>

## How to add grpc
 - create protobuf IDL(Interface Definition Language) file
   - see [grpc/sample/sample.proto](https://github.com/common-library/go/blob/main/grpc/sample/sample.proto)
 - convert IDL file to code
   - `go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.10`
   - `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1`
   - `wget https://github.com/protocolbuffers/protobuf/releases/download/v32.1/protoc-32.1-linux-x86_64.zip`
   - `unzip protoc-32.1-linux-x86_64.zip -d protoc/`
   - `protoc/bin/protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative grpc/sample/sample.proto`
  - implement functions defined in IDL file
    - implement to satisfy [implementServer interface](https://github.com/common-library/go/blob/main/grpc/server.go)
    - see [grpc/sample/Server.go](https://github.com/common-library/go/blob/main/grpc/sample/Server.go)
    - see [grpc/sample/Server_test.go](https://github.com/common-library/go/blob/main/grpc/sample/Server_test.go)
