# Common Library for Go

## Installation
```bash
go get -u github.com/common-library/go
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
 - command line
   - arguments
   - flag
 - data structure
   - Deque
   - Queue
 - database
   - Elasticsearch v7/v8
   - MongoDB
   - Prometheus
   - Redis
   - ORM
     - beego
     - ent
     - GORM
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
 - prepare
   - Amazon
     - DynamoDB
       - `docker run --name dynamodb --detach --publish 8000:8000 --env "-jar DynamoDBLocal.jar -sharedDb -inMemory" amazon/dynamodb-local:2.4.0`
       - `export DYNAMODB_URL=http://127.0.0.1:8000`
     - S3
       - `docker run --name s3mock --detach --publish 9090:9090 --publish 9191:9191 adobe/s3mock:3.7.3`
       - `export S3_URL=http://127.0.0.1:9090`
   - Elasticsearch
     - v7
       - `docker run --name elasticsearch-v7 --detach --publish 19200:9200 --publish 19300:9300 --env discovery.type=single-node --env ES_JAVA_OPTS="-Xms500m -Xmx500m" elasticsearch:7.17.21`
       - `export ELASTICSEARCH_ADDRESS_V7=http://:19200`
     - v8
       - `docker network create elastic`
       - `docker run --name elasticsearch-v8 --net elastic --detach --publish 29200:9200 --publish 29300:9300 --env discovery.type=single-node --env ES_JAVA_OPTS="-Xms500m -Xmx500m" --env xpack.security.enabled=false elasticsearch:8.13.0`
       - `export ELASTICSEARCH_ADDRESS_V8=http://:29200`
   - GEMINI
     - [Set up API key](https://ai.google.dev/gemini-api/docs/get-started/tutorial?lang=go&hl=ko#set-up-api-key)
     - `export GEMINI_API_KEY=${API_KEY}`
   - MongoDB
     - `docker run --name mongodb --detach --publish 27017:27017 mongo:7.0.9`
     - `export MONGODB_ADDRESS=:27017`
   - Prometheus
     - `wget https://github.com/prometheus/prometheus/releases/download/v2.52.0/prometheus-2.52.0.linux-amd64.tar.gz`
     - `tar xvzf prometheus-2.52.0.linux-amd64.tar.gz`
     - `cd prometheus-2.52.0.linux-amd64`
     - `yq -i '.scrape_configs[0].static_configs[0].targets = [":9095"]' prometheus.yml`
     - `./prometheus --config.file=prometheus.yml --web.listen-address=:9095 &`
     - `export PROMETHEUS_ADDRESS=:9095`
   - Redis
     - `docker run --name redis --detach --publish 6379:6379 redis:7.2.4`
     - `export REDIS_ADDRESS=:6379`
   - SQL
     - Amazon DynamoDB
       - `docker run --name dynamodb --detach --publish 8000:8000 --env "-jar DynamoDBLocal.jar -sharedDb -inMemory" amazon/dynamodb-local:2.4.0`
       - `export DYNAMODB_URL=http://127.0.0.1:8000`
     - ClickHouse
       - `docker run --name clickhouse --detach --publish 19000:9000 --ulimit nofile=262144:262144 clickhouse/clickhouse-server:24.4.1`
       - `export CLICKHOUSE_DSN='clickhouse://default:@127.0.0.1:19000'`
     - MySQL
       - `docker run --name mysql --detach --publish 3306:3306 --env MYSQL_ROOT_PASSWORD=root mysql:8.4.0`
       - `export MYSQL_DSN='root:root@tcp(127.0.0.1)/'`
     - PostgreSQL
       - `docker run --name postgres --detach --publish 5432:5432 --env POSTGRES_PASSWORD=postgres postgres:16.2-alpine`
       - `export POSTGRESQL_DSN='host=localhost port=5432 user=postgres password=postgres sslmode=disable'`
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

## How to add grpc
 - create protobuf IDL(Interface Definition Language) file
   - see [grpc/sample/sample.proto](https://github.com/common-library/go/blob/main/grpc/sample/sample.proto)
 - convert IDL file to code
   - `go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0`
   - `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0`
   - `wget https://github.com/protocolbuffers/protobuf/releases/download/v3.20.3/protoc-3.20.3-linux-x86_64.zip`
   - `unzip protoc-3.20.3-linux-x86_64.zip -d protoc/`
   - `protoc/bin/protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative grpc/sample/sample.proto`
  - implement functions defined in IDL file
    - implement to satisfy [implementServer interface](https://github.com/common-library/go/blob/main/grpc/server.go)
    - see [grpc/sample/Server.go](https://github.com/common-library/go/blob/main/grpc/sample/Server.go)
    - see [grpc/sample/Server_test.go](https://github.com/common-library/go/blob/main/grpc/sample/Server_test.go)
