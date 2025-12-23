# Common Library for Go

[![CI](https://github.com/common-library/go/workflows/CI/badge.svg)](https://github.com/common-library/go/actions)
[![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/heaven-chp/c7e11ff6ca6c490bd028e4a6d9b79c92/raw/coverage.json)](https://github.com/common-library/go/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/common-library/go)](https://goreportcard.com/report/github.com/common-library/go)
[![Go Version](https://img.shields.io/github/go-mod/go-version/common-library/go?logo=go)](https://github.com/common-library/go)
[![Reference](https://pkg.go.dev/badge/github.com/common-library/go.svg)](https://pkg.go.dev/github.com/common-library/go)
[![License](https://img.shields.io/github/license/common-library/go)](https://github.com/common-library/go/blob/main/LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/common-library/go)](https://github.com/common-library/go/stargazers)

A comprehensive, production-ready Go library providing utilities and wrappers for common development tasks. This library simplifies complex operations while maintaining type safety, performance, and idiomatic Go patterns.

## ✨ Highlights

- 🚀 **Production-Ready** - Battle-tested in production environments
- 📦 **Comprehensive** - 20+ packages covering common development needs
- 🔒 **Type-Safe** - Leveraging Go generics for compile-time safety
- 📚 **Well-Documented** - Extensive documentation and examples for every package
- ⚡ **High Performance** - Optimized implementations with minimal overhead
- 🧪 **Thoroughly Tested** - Comprehensive test coverage with integration tests

## 📦 Installation

```bash
go get -u github.com/common-library/go
```

## 🚀 Quick Start

```go
package main

import (
    "fmt"
    "github.com/common-library/go/http"
    "github.com/common-library/go/json"
)

func main() {
    // HTTP request made easy
    response, _ := http.Get("https://api.example.com/data", nil, "", "", nil)
    
    // Type-safe JSON conversion
    data, _ := json.ConvertFromString[map[string]interface{}](response)
    
    fmt.Println(data)
}
```



## 📚 Features

### 📁 Archive & Compression
- **[gzip](archive/gzip)** - Gzip compression and decompression
- **[tar](archive/tar)** - TAR archive creation and extraction
- **[zip](archive/zip)** - ZIP archive operations

### ☁️ AWS Services
- **[DynamoDB](aws/dynamodb)** - AWS DynamoDB client with simplified operations
- **[S3](aws/s3)** - Amazon S3 object storage operations

### 💻 Command Line
- **[arguments](command-line/arguments)** - Command-line argument parsing
- **[flags](command-line/flags)** - Type-safe flag parsing with generics

### 🗂️ Data Structures
- **[collection](collection)** - Generic data structures (Deque, Queue)

### 🗄️ Databases

#### NoSQL & Document Stores
- **[Elasticsearch v7/v8](database/elasticsearch)** - Full-text search and analytics
- **[MongoDB](database/mongodb)** - Document database operations
- **[Redis](database/redis)** - In-memory data structure store

#### SQL Databases
- **[sql](database/sql)** - Unified SQL client supporting:
  - Amazon DynamoDB
  - ClickHouse
  - Microsoft SQL Server
  - MySQL
  - Oracle Database
  - PostgreSQL
  - SQLite

#### ORM & Query Builders
- **[beego ORM](database/orm/beego)** - Beego ORM wrapper
- **[ent](database/orm/ent)** - Entity framework for Go
- **[GORM](database/orm/gorm)** - Feature-rich ORM library
- **[sqlc](database/orm/sqlc)** - Compile-time SQL query generator
- **[sqlx](database/orm/sqlx)** - Extensions to database/sql

#### Database Tools
- **[dbmate](database/dbmate)** - Database migration tool
- **[Prometheus client](database/prometheus)** - Prometheus metrics querying

### 📡 Events & Messaging
- **[CloudEvents](event/cloudevents)** - CloudEvents client and server

### 📄 File Operations
- **[file](file)** - File and directory utilities

### 🌐 Network & Communication
- **[gRPC](grpc)** - Simplified gRPC client and server ([Documentation](grpc/README.md))
- **[HTTP](http)** - HTTP client with retry logic and utilities
- **[Socket](socket)** - TCP/UDP socket server and client
- **[Long Polling](long-polling)** - HTTP long polling implementation

### 📦 Data Formats
- **[JSON](json)** - Type-safe JSON marshaling with generics

### ☸️ Kubernetes
- **[resource/client](kubernetes/resource/client)** - Kubernetes resource management
- **[Custom Resources](kubernetes/resource/custom-resource)** - CRD operations
- **[CRD](kubernetes/resource/custom-resource-definition)** - Custom Resource Definitions

### 🔒 Concurrency & Synchronization
- **[lock](lock)** - Mutex utilities with key-based locking

### 📝 Logging
- **[klog](log/klog)** - Kubernetes-style structured logging
- **[slog](log/slog)** - Structured logging with context

### 🔐 Security
- **[crypto/dsa](security/crypto/dsa)** - DSA signatures (deprecated, legacy support)
- **[crypto/ecdsa](security/crypto/ecdsa)** - ECDSA elliptic curve signatures
- **[crypto/ed25519](security/crypto/ed25519)** - Ed25519 signatures (recommended)
- **[crypto/rsa](security/crypto/rsa)** - RSA encryption and signatures

### 💾 Storage
- **[MinIO](storage/minio)** - S3-compatible object storage client

### 🛠️ Utilities
- **[utility](utility)** - Runtime introspection, type info, and CIDR utilities

### 🧪 Testing
- **[testutil](testutil)** - Testing utilities and container image constants



## 🧪 Testing & Development

### Running Tests

Run all tests with coverage:
```bash
go clean -testcache && go test -cover ./...
```

### Coverage Reports

Generate coverage profile:
```bash
go clean -testcache && go test -coverprofile=coverage.out -cover ./...
```

Convert to HTML report:
```bash
go tool cover -html=./coverage.out -o ./coverage.html
```

View in browser:
```bash
open coverage.html  # macOS
xdg-open coverage.html  # Linux
```

## 🔧 Development Guides

### Working with Ent ORM

Add a new entity schema (e.g., `User`):

```bash
# Install ent CLI
go get entgo.io/ent/cmd/ent

# Create new schema
go run entgo.io/ent/cmd/ent new --target ./database/orm/ent/schema User

# Edit schema file
vim ./database/orm/ent/schema/user.go

# Generate code with upsert support
go run entgo.io/ent/cmd/ent generate --feature sql/upsert ./database/orm/ent/schema
```

### Working with sqlc

Configure and generate type-safe SQL code:

```bash
# Install sqlc
go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0

# Add queries to ./database/orm/sqlc/queries
# Add schemas to ./database/orm/sqlc/schema
# Configure ./database/orm/sqlc/sqlc.json

# Generate Go code
sqlc generate --file ./database/orm/sqlc/sqlc.json
```

### Working with gRPC

Create and implement gRPC services:

```bash
# Install protoc compiler and plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.10
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1

# Download protoc (example for Linux)
wget https://github.com/protocolbuffers/protobuf/releases/download/v32.1/protoc-32.1-linux-x86_64.zip
unzip protoc-32.1-linux-x86_64.zip -d protoc/

# Create your .proto file (see grpc/sample/sample.proto for example)
# Generate Go code
protoc/bin/protoc \
  --go_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_out=. \
  --go-grpc_opt=paths=source_relative \
  grpc/sample/sample.proto
```

**Implementation Steps:**
1. Create `.proto` file with service definition ([example](grpc/sample/sample.proto))
2. Generate Go code using protoc
3. Implement the service interface ([example](grpc/sample/Server.go))
4. Register service with gRPC server (implement `RegisterServer` method)
5. Write tests ([example](grpc/sample/Server_test.go))

## 📖 Documentation

Each package includes comprehensive documentation:
- Package-level godoc comments
- Function/method documentation with examples
- README.md files for complex packages
- Integration test examples

Browse the [pkg.go.dev documentation](https://pkg.go.dev/github.com/common-library/go) or explore individual package directories.

## 🤝 Contributing

Contributions are welcome! Please ensure:
- All tests pass: `go test ./...`
- Code is formatted: `go fmt ./...`
- Linting passes: `golangci-lint run`
- Documentation is updated
- Examples are provided for new features

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## ⭐ Support

If you find this library helpful, please consider giving it a star on [GitHub](https://github.com/common-library/go)!

## 🔗 Related Projects

- [Go Standard Library](https://pkg.go.dev/std)
- [Awesome Go](https://github.com/avelino/awesome-go)

## 📞 Contact & Support

- 🐛 **Issues**: [GitHub Issues](https://github.com/common-library/go/issues)
- 📖 **Documentation**: [pkg.go.dev](https://pkg.go.dev/github.com/common-library/go)
- ⭐ **Star us**: [GitHub Repository](https://github.com/common-library/go)

