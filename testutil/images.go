// Package testutil provides utilities for database testing, including container image definitions.
package testutil

// Container images used across all database tests.
// Centralizing image versions here makes it easy to update and maintain consistency.
const (
	// Database images
	ClickHouseImage = "clickhouse/clickhouse-server:25.10.2-alpine"
	MySQLImage      = "mysql:9.5.0"
	PostgresImage   = "postgres:18.1-alpine"
	MongoImage      = "mongo:8.0.16"
	RedisImage      = "redis:8.4.0-alpine3.22"

	// Elasticsearch images (versioned)
	ElasticsearchV7Image = "elasticsearch:7.17.28"
	ElasticsearchV8Image = "elasticsearch:8.19.7"
	ElasticsearchV9Image = "elasticsearch:9.2.1"

	// AWS service images
	LocalstackImage = "localstack/localstack:4.10.0"

	// Monitoring images
	PrometheusImage = "prom/prometheus:v3.7.3"
)
