# Test Utilities

Container image constants for integration testing.

## Overview

The testutil package provides centralized Docker container image definitions used across all integration tests. This ensures consistency and simplifies version management.

## Features

- **Database Images** - ClickHouse, MySQL, Postgres, MongoDB, Redis
- **Elasticsearch Images** - Versions 7, 8, and 9
- **AWS Services** - LocalStack for AWS service emulation
- **Monitoring** - Prometheus for metrics testing
- **Centralized Management** - Single source of truth for image versions

## Installation

```bash
go get -u github.com/common-library/go/testutil
```

## Quick Start

```go
import "github.com/common-library/go/testutil"

container := testcontainers.ContainerRequest{
    Image: testutil.PostgresImage,
    ExposedPorts: []string{"5432/tcp"},
    // ...
}
```

## Available Images

### Database Images

- `ClickHouseImage` - ClickHouse OLAP database
- `MySQLImage` - MySQL relational database
- `PostgresImage` - PostgreSQL database
- `MongoImage` - MongoDB document database
- `RedisImage` - Redis key-value store

### Elasticsearch Images

- `ElasticsearchV7Image` - Elasticsearch 7.x
- `ElasticsearchV8Image` - Elasticsearch 8.x
- `ElasticsearchV9Image` - Elasticsearch 9.x

### Service Images

- `LocalstackImage` - LocalStack for AWS services
- `PrometheusImage` - Prometheus monitoring

## Example

```go
package mytest

import (
    "testing"
    "github.com/testcontainers/testcontainers-go"
    "github.com/common-library/go/testutil"
)

func TestWithPostgres(t *testing.T) {
    req := testcontainers.ContainerRequest{
        Image: testutil.PostgresImage,
        ExposedPorts: []string{"5432/tcp"},
        Env: map[string]string{
            "POSTGRES_PASSWORD": "password",
        },
    }
    
    container, _ := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started: true,
    })
    defer container.Terminate(ctx)
    
    // Run tests...
}
```

## Dependencies

None - constants only.

## Further Reading

- [Testcontainers Go](https://golang.testcontainers.org/)
