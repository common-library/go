// Package testutil provides testing utilities and container image constants.
//
// This package centralizes Docker container image versions used across all
// database and service integration tests. Maintaining images here ensures
// consistency and simplifies version updates.
//
// # Features
//
//   - Database container images (ClickHouse, MySQL, Postgres, Mongo, Redis)
//   - Elasticsearch images (v7, v8, v9)
//   - AWS service images (LocalStack)
//   - Monitoring images (Prometheus)
//   - Centralized version management
//
// # Basic Example
//
//	import "github.com/common-library/go/testutil"
//
//	container := testcontainers.ContainerRequest{
//	    Image: testutil.PostgresImage,
//	    // ...
//	}
package testutil

// Container images used across all database tests.
// Centralizing image versions here makes it easy to update and maintain consistency.
const (
	// Database images
	ClickHouseImage = "clickhouse/clickhouse-server:25.12.1-alpine"
	MySQLImage      = "mysql:9.5.0"
	PostgresImage   = "postgres:18.1-alpine"
	MongoImage      = "mongo:8.2.3"
	RedisImage      = "redis:8.4.0-alpine3.22"

	// Elasticsearch images (versioned)
	ElasticsearchV7Image = "elasticsearch:7.17.28"
	ElasticsearchV8Image = "elasticsearch:8.19.9"
	ElasticsearchV9Image = "elasticsearch:9.2.3"

	// AWS service images
	LocalstackImage = "localstack/localstack:4.12.0"

	// Monitoring images
	PrometheusImage = "prom/prometheus:v3.8.1"
)
