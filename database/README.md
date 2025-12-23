# Database

Comprehensive database client libraries and tools for Go applications.

## Overview

This package provides a collection of database clients, ORM tools, and utilities for working with various databases and data stores. Each subpackage focuses on a specific database technology or use case.

## Packages

### Core Database Clients

#### [SQL](sql/)
Unified SQL database client supporting 7 database drivers with a consistent API.

**Supported Databases:**
- MySQL
- PostgreSQL
- SQLite
- ClickHouse
- Amazon DynamoDB
- Microsoft SQL Server
- Oracle

**Features:**
- Connection pooling
- Transaction management
- Prepared statements
- Consistent API across all databases

**Quick Example:**
```go
import "github.com/common-library/go/database/sql"

var client sql.Client
client.Open(sql.DriverMySQL, "user:pass@tcp(localhost)/db", 10)
defer client.Close()

client.Execute("INSERT INTO users (name) VALUES (?)", "Alice")
rows, _ := client.Query("SELECT id, name FROM users")
```

**[Full Documentation →](sql/)**

---

#### [Redis](redis/)
Redis client with connection pooling and simplified operations.

**Features:**
- Connection pooling
- String, Hash, List, Set, Sorted Set operations
- Key expiration (TTL)
- Database selection
- Batch operations (MGET, MSET)

**Quick Example:**
```go
import "github.com/common-library/go/database/redis"

var client redis.Client
client.Initialize("localhost:6379", "", 10, 60*time.Second)
defer client.Finalize()

client.Set("key", "value")
value, _ := client.Get("key")
```

**[Full Documentation →](redis/)**

---

#### [MongoDB](mongodb/)
MongoDB client wrapper with simplified operations.

**Features:**
- Automatic reconnection
- CRUD operations
- Aggregation support
- Index management
- Bulk operations

**Quick Example:**
```go
import "github.com/common-library/go/database/mongodb"

var client mongodb.Client
client.Initialize("localhost:27017", 10*time.Second)
defer client.Finalize()

client.InsertOne("mydb", "users", bson.M{"name": "Alice"})
result, _ := client.FindOne("mydb", "users", bson.M{"name": "Alice"}, User{})
```

**[Full Documentation →](mongodb/)**

---

#### [Elasticsearch](elasticsearch/)
Multi-version Elasticsearch client (v7, v8, v9) with unified interface.

**Features:**
- Support for Elasticsearch 7.x, 8.x, 9.x
- Document CRUD operations
- Index management
- Template management
- Search queries
- Cloud and on-premise support

**Quick Example:**
```go
import "github.com/common-library/go/database/elasticsearch/v8"

var client v8.Client
client.Initialize([]string{"http://localhost:9200"}, 10*time.Second, "", "", "user", "pass", "", nil)

client.Index("users", "1", `{"name":"Alice","age":30}`)
result, _ := client.Search("users", `{"query":{"match_all":{}}}`)
```

**[Full Documentation →](elasticsearch/)**

---

### ORM and Query Tools

#### [ORM](orm/)
Object-Relational Mapping tools and examples.

**Supported ORMs:**
- **Ent** - Entity framework with code generation
- **GORM** - Popular ORM with chainable API
- **SQLC** - Type-safe SQL from SQL queries
- **SQLx** - Extensions for database/sql
- **Beego ORM** - Beego framework ORM

**Quick Example (GORM):**
```go
import "gorm.io/gorm"

db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
db.AutoMigrate(&User{})
db.Create(&User{Name: "Alice"})
```

**[Full Documentation →](orm/)**

---

### Database Tools

#### [Dbmate](dbmate/)
Database migration management using dbmate.

**Features:**
- Migration versioning
- Support for MySQL, PostgreSQL, ClickHouse
- Up/down migrations
- Migration status tracking

**Quick Example:**
```bash
# Create migration
dbmate new create_users_table

# Apply migrations
dbmate up

# Rollback
dbmate down
```

**[Full Documentation →](dbmate/)**

---

#### [Prometheus](prometheus/)
Prometheus client and exporter utilities.

**Features:**
- Metrics client for querying Prometheus
- Custom exporter creation
- Gauge, Counter, Histogram, Summary support
- HTTP handler for /metrics endpoint

**Quick Example:**
```go
import "github.com/common-library/go/database/prometheus/client"

c, _ := client.NewClient("http://localhost:9090")
value, _, _ := c.Query("up", time.Now(), 10*time.Second)
```

**[Full Documentation →](prometheus/)**

---

## Choosing the Right Package

### By Use Case

| Use Case | Recommended Package | Why |
|----------|-------------------|-----|
| SQL database access | [sql](sql/) | Unified API for 7 databases |
| Object-relational mapping | [orm](orm/) | Type-safe models, code generation |
| Caching | [redis](redis/) | Fast in-memory operations |
| Document storage | [mongodb](mongodb/) | Flexible schema, JSON documents |
| Full-text search | [elasticsearch](elasticsearch/) | Advanced search capabilities |
| Database migrations | [dbmate](dbmate/) | Version control for database schema |
| Metrics storage | [prometheus](prometheus/) | Time-series data |

### By Database Type

| Database | Package | Notes |
|----------|---------|-------|
| MySQL | [sql](sql/) or [orm](orm/) | Use sql for raw queries, orm for models |
| PostgreSQL | [sql](sql/) or [orm](orm/) | Use sql for raw queries, orm for models |
| SQLite | [sql](sql/) or [orm](orm/) | Great for embedded/testing |
| MongoDB | [mongodb](mongodb/) | NoSQL document database |
| Redis | [redis](redis/) | Key-value store, caching |
| Elasticsearch | [elasticsearch](elasticsearch/) | Search engine |
| ClickHouse | [sql](sql/) | Analytics, time-series |
| DynamoDB | [sql](sql/) | AWS serverless NoSQL |

### By Performance Needs

**High Throughput:**
- Redis (in-memory, microsecond latency)
- ClickHouse (columnar, analytics)

**ACID Transactions:**
- PostgreSQL (strong consistency)
- MySQL (InnoDB engine)

**Horizontal Scaling:**
- DynamoDB (serverless, automatic scaling)
- MongoDB (sharding support)
- Elasticsearch (distributed by design)

**Embedded/Serverless:**
- SQLite (no server needed)

## Quick Start Guide

### 1. Install Required Drivers

```bash
# SQL databases
go get github.com/go-sql-driver/mysql
go get github.com/lib/pq
go get modernc.org/sqlite

# Redis
go get github.com/gomodule/redigo/redis

# MongoDB
go get go.mongodb.org/mongo-driver/mongo

# Elasticsearch
go get github.com/elastic/go-elasticsearch/v8
```

### 2. Import the Package

```go
import "github.com/common-library/go/database/sql"
// or
import "github.com/common-library/go/database/redis"
// or
import "github.com/common-library/go/database/mongodb"
```

### 3. Initialize Client

```go
// SQL
var sqlClient sql.Client
sqlClient.Open(sql.DriverPostgreSQL, dsn, 25)
defer sqlClient.Close()

// Redis
var redisClient redis.Client
redisClient.Initialize("localhost:6379", "", 10, 60*time.Second)
defer redisClient.Finalize()

// MongoDB
var mongoClient mongodb.Client
mongoClient.Initialize("localhost:27017", 10*time.Second)
defer mongoClient.Finalize()
```

### 4. Perform Operations

```go
// SQL: Query
rows, err := sqlClient.Query("SELECT * FROM users WHERE age > ?", 18)

// Redis: Set/Get
redisClient.Set("user:123", "Alice")
value, err := redisClient.Get("user:123")

// MongoDB: Find
result, err := mongoClient.FindOne("mydb", "users", bson.M{"name": "Alice"}, User{})
```

## Common Patterns

### Connection Pooling

All database clients support connection pooling for optimal performance:

```go
// SQL: Set max open connections
sqlClient.Open(driver, dsn, 100) // 100 max connections

// Redis: Configure pool
redisClient.Initialize(address, password, 50, 60*time.Second) // 50 idle connections

// MongoDB: Timeout for operations
mongoClient.Initialize(address, 30*time.Second) // 30s timeout
```

### Error Handling

```go
if err := client.Execute("INSERT INTO users (name) VALUES (?)", name); err != nil {
    // Handle database error
    log.Printf("Database error: %v", err)
    return err
}
```

### Resource Cleanup

Always close database connections:

```go
defer sqlClient.Close()
defer redisClient.Finalize()
defer mongoClient.Finalize()
```

### Transactions (SQL Only)

```go
err := sqlClient.BeginTransaction()
err = sqlClient.ExecuteTransaction("UPDATE accounts SET balance = balance - 100 WHERE id = 1")
err = sqlClient.ExecuteTransaction("UPDATE accounts SET balance = balance + 100 WHERE id = 2")
err = sqlClient.EndTransaction(err) // Commit or rollback
```

## Testing Support

Most packages include test utilities:

```go
// Elasticsearch: Test container
import "github.com/common-library/go/database/elasticsearch/testutil"

container := testutil.Container{}
container.Run("8.0.0")
defer container.Terminate()
```

## Best Practices

### 1. Use Connection Pooling

```go
// Good: Configure pool size
client.Open(driver, dsn, 25)

// Avoid: Creating new connection per request
// for each request {
//     client.Open(driver, dsn, 1)
//     client.Close()
// }
```

### 2. Close Resources

```go
// Good: Defer close
defer client.Close()
defer rows.Close()

// Avoid: Forgetting to close
// rows, _ := client.Query("SELECT ...")
// // rows never closed, connection leak!
```

### 3. Handle Errors

```go
// Good: Check all errors
if err := client.Execute("INSERT ..."); err != nil {
    return fmt.Errorf("insert failed: %w", err)
}

// Avoid: Ignoring errors
// client.Execute("INSERT ...") // Error silently ignored
```

### 4. Use Appropriate Database

```go
// Caching: Use Redis
redisClient.Setex("session:123", 3600, sessionData)

// Transactions: Use SQL
sqlClient.BeginTransaction()
// ...
sqlClient.EndTransaction(err)

// Full-text search: Use Elasticsearch
esClient.Search("products", searchQuery)
```

### 5. Secure Credentials

```go
// Good: Use environment variables
dsn := os.Getenv("DATABASE_URL")

// Avoid: Hardcoded credentials
// dsn := "user:password@tcp(localhost)/db"
```

## Migration Guide

### From database/sql

```go
// Before
db, _ := sql.Open("mysql", dsn)
stmt, _ := db.Prepare("SELECT * FROM users WHERE id = ?")
rows, _ := stmt.Query(123)

// After
var client sql.Client
client.Open(sql.DriverMySQL, dsn, 10)
client.SetPrepare("SELECT * FROM users WHERE id = ?")
rows, _ := client.QueryPrepare(123)
```

### From go-redis/redis

```go
// Before
rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
rdb.Set(ctx, "key", "value", 0)

// After
var client redis.Client
client.Initialize("localhost:6379", "", 10, 60*time.Second)
client.Set("key", "value")
```

## Performance Considerations

| Database | Best For | Throughput | Latency |
|----------|----------|------------|---------|
| Redis | Cache, sessions | 100K+ ops/sec | < 1ms |
| PostgreSQL | OLTP, complex queries | 10K+ TPS | 1-10ms |
| ClickHouse | Analytics, time-series | 1M+ rows/sec | 10-100ms |
| MongoDB | Document storage | 10K+ ops/sec | 1-10ms |
| Elasticsearch | Search, logs | 10K+ docs/sec | 10-100ms |
| SQLite | Embedded, testing | 100K+ ops/sec | < 1ms |

## Dependencies

Each package has its own dependencies. See individual package documentation for details.

**Common dependencies:**
- Go 1.18+ (for generics support in some packages)
- Database drivers (package-specific)

## Troubleshooting

### Connection Errors

```go
// Check connection
if err := client.Open(driver, dsn, 10); err != nil {
    log.Fatal("Connection failed:", err)
}
// Error: "dial tcp: connection refused" → Database not running
// Error: "access denied" → Wrong credentials
```

### Pool Exhaustion

```go
// Symptom: "connection pool exhausted"
// Solution: Increase pool size or close connections
client.Open(driver, dsn, 100) // Increase from 10 to 100
defer rows.Close() // Always close rows
```

### Transaction Deadlocks

```go
// Use shorter transactions
err := client.BeginTransaction()
// Do minimal work here
err = client.EndTransaction(err)
```

## Further Reading

- [Go database/sql tutorial](https://go.dev/doc/database/overview)
- [Redis documentation](https://redis.io/documentation)
- [MongoDB Go driver](https://www.mongodb.com/docs/drivers/go/current/)
- [Elasticsearch Go client](https://www.elastic.co/guide/en/elasticsearch/client/go-api/current/index.html)
- [GORM documentation](https://gorm.io/docs/)
- [Ent documentation](https://entgo.io/docs/getting-started)
