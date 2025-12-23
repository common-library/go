# Dbmate

Database migration management using dbmate for ClickHouse, MySQL, and PostgreSQL.

## Overview

The dbmate package provides integration examples and test utilities for managing database schema migrations using [dbmate](https://github.com/amacneil/dbmate), a lightweight database migration tool. This package includes migration examples for three popular databases and demonstrates best practices for automated migration testing.

## Features

- **Multi-Database Support** - ClickHouse, MySQL, PostgreSQL migration examples
- **Version Control** - Track schema changes with timestamped migration files
- **Up/Down Migrations** - Support for both applying and rolling back migrations
- **Automated Testing** - Integration tests using testcontainers
- **Simple Workflow** - Create, migrate, rollback, and drop databases
- **Schema Validation** - Verify migration results programmatically

## Installation

```bash
go get -u github.com/amacneil/dbmate/v2
go get -u github.com/amacneil/dbmate/v2/pkg/dbmate
go get -u github.com/amacneil/dbmate/v2/pkg/driver/clickhouse
go get -u github.com/amacneil/dbmate/v2/pkg/driver/mysql
go get -u github.com/amacneil/dbmate/v2/pkg/driver/postgres
```

## Migration File Structure

Migrations are organized by database type:

```
dbmate/
├── clickhouse/
│   └── migrations/
│       ├── 20240929074201_create_table_test_01.sql
│       └── 20240929074254_alter_table_test_01.sql
├── mysql/
│   └── migrations/
│       ├── 20240929074201_create_table_test_01.sql
│       └── 20240929074254_alter_table_test_01.sql
└── postgresql/
    └── migrations/
        ├── 20240929074201_create_table_test_01.sql
        └── 20240929074254_alter_table_test_01.sql
```

### Migration File Format

Each migration file contains both `up` and `down` SQL:

```sql
-- migrate:up
CREATE TABLE IF NOT EXISTS users(
    id INTEGER PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE
);

-- migrate:down
DROP TABLE IF EXISTS users;
```

## Quick Start

### 1. Initialize Dbmate

```go
import (
    "net/url"
    "github.com/amacneil/dbmate/v2/pkg/dbmate"
    _ "github.com/amacneil/dbmate/v2/pkg/driver/mysql"
)

// Parse database URL
u, _ := url.Parse("mysql://user:pass@localhost:3306/mydb")

// Create dbmate instance
db := dbmate.New(u)
db.AutoDumpSchema = false
db.MigrationsDir = []string{"./migrations"}
```

### 2. Create Database and Run Migrations

```go
// Create database
if err := db.Create(); err != nil {
    log.Fatal(err)
}

// Run all pending migrations
if err := db.Migrate(); err != nil {
    log.Fatal(err)
}
```

### 3. Rollback Migration

```go
// Rollback last migration
if err := db.Rollback(); err != nil {
    log.Fatal(err)
}
```

### 4. Drop Database

```go
// Drop database
if err := db.Drop(); err != nil {
    log.Fatal(err)
}
```

## Database-Specific Examples

### ClickHouse

**Connection URL:**
```
clickhouse://host:port/database
clickhouse://localhost:9000/mydb
```

**Migration Example:**
```sql
-- migrate:up
CREATE TABLE IF NOT EXISTS events(
    event_id String,
    user_id UInt32,
    event_time DateTime
)
ENGINE = MergeTree
PRIMARY KEY (event_id)
ORDER BY (event_id, event_time);

-- migrate:down
DROP TABLE IF EXISTS events;
```

**Usage:**
```go
import _ "github.com/amacneil/dbmate/v2/pkg/driver/clickhouse"

u, _ := url.Parse("clickhouse://localhost:9000/analytics")
db := dbmate.New(u)
db.MigrationsDir = []string{"./clickhouse/migrations"}
db.CreateAndMigrate()
```

### MySQL

**Connection URL:**
```
mysql://user:password@host:port/database?parseTime=true
mysql://root:secret@localhost:3306/myapp?parseTime=true
```

**Migration Example:**
```sql
-- migrate:up
CREATE TABLE IF NOT EXISTS users(
    id INTEGER AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_username ON users(username);

-- migrate:down
DROP TABLE IF EXISTS users;
```

**Usage:**
```go
import _ "github.com/amacneil/dbmate/v2/pkg/driver/mysql"

u, _ := url.Parse("mysql://root:password@localhost:3306/myapp?parseTime=true")
db := dbmate.New(u)
db.MigrationsDir = []string{"./mysql/migrations"}
db.CreateAndMigrate()
```

### PostgreSQL

**Connection URL:**
```
postgres://user:password@host:port/database?sslmode=disable
postgres://postgres:secret@localhost:5432/myapp?sslmode=disable
```

**Migration Example:**
```sql
-- migrate:up
CREATE TABLE IF NOT EXISTS orders(
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    total_price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_user_orders ON orders(user_id);

-- migrate:down
DROP TABLE IF EXISTS orders;
```

**Usage:**
```go
import _ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"

u, _ := url.Parse("postgres://postgres:password@localhost:5432/myapp?sslmode=disable")
db := dbmate.New(u)
db.MigrationsDir = []string{"./postgresql/migrations"}
db.CreateAndMigrate()
```

## Complete Examples

### E-commerce Schema Migration

**migrations/20240101000001_create_products.sql:**
```sql
-- migrate:up
CREATE TABLE IF NOT EXISTS products(
    id INTEGER AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    stock INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- migrate:down
DROP TABLE IF EXISTS products;
```

**migrations/20240101000002_create_orders.sql:**
```sql
-- migrate:up
CREATE TABLE IF NOT EXISTS orders(
    id INTEGER AUTO_INCREMENT PRIMARY KEY,
    product_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    total DECIMAL(10,2) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id)
);

-- migrate:down
DROP TABLE IF EXISTS orders;
```

**Migration Script:**
```go
package main

import (
    "log"
    "net/url"
    
    "github.com/amacneil/dbmate/v2/pkg/dbmate"
    _ "github.com/amacneil/dbmate/v2/pkg/driver/mysql"
)

func main() {
    // Database connection
    u, err := url.Parse("mysql://root:password@localhost:3306/ecommerce?parseTime=true")
    if err != nil {
        log.Fatal(err)
    }
    
    // Initialize dbmate
    db := dbmate.New(u)
    db.AutoDumpSchema = false
    db.MigrationsDir = []string{"./migrations"}
    
    // Create database and run migrations
    if err := db.CreateAndMigrate(); err != nil {
        log.Fatal(err)
    }
    
    log.Println("Database migrated successfully")
}
```

### Multi-Tenant Schema

**migrations/20240101000001_create_tenants.sql:**
```sql
-- migrate:up
CREATE TABLE IF NOT EXISTS tenants(
    id INTEGER AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    subdomain VARCHAR(50) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- migrate:down
DROP TABLE IF EXISTS tenants;
```

**migrations/20240101000002_create_tenant_users.sql:**
```sql
-- migrate:up
CREATE TABLE IF NOT EXISTS tenant_users(
    id INTEGER AUTO_INCREMENT PRIMARY KEY,
    tenant_id INTEGER NOT NULL,
    email VARCHAR(100) NOT NULL,
    role VARCHAR(20) DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    UNIQUE KEY unique_tenant_email (tenant_id, email)
);

CREATE INDEX idx_tenant_users ON tenant_users(tenant_id);

-- migrate:down
DROP TABLE IF EXISTS tenant_users;
```

### Analytics Schema (ClickHouse)

**clickhouse/migrations/20240101000001_create_events.sql:**
```sql
-- migrate:up
CREATE TABLE IF NOT EXISTS page_views(
    timestamp DateTime,
    user_id UInt32,
    page_url String,
    referrer String,
    country String,
    device String
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(timestamp)
ORDER BY (timestamp, user_id);

-- migrate:down
DROP TABLE IF EXISTS page_views;
```

**clickhouse/migrations/20240101000002_create_aggregation.sql:**
```sql
-- migrate:up
CREATE MATERIALIZED VIEW IF NOT EXISTS daily_stats
ENGINE = SummingMergeTree
PARTITION BY toYYYYMM(date)
ORDER BY (date, country)
AS SELECT
    toDate(timestamp) as date,
    country,
    count() as views,
    uniq(user_id) as unique_users
FROM page_views
GROUP BY date, country;

-- migrate:down
DROP VIEW IF EXISTS daily_stats;
```

## Testing Migrations

### Integration Test Example

```go
package main

import (
    "context"
    "fmt"
    "net/url"
    "testing"
    
    "github.com/amacneil/dbmate/v2/pkg/dbmate"
    _ "github.com/amacneil/dbmate/v2/pkg/driver/mysql"
    "github.com/testcontainers/testcontainers-go/modules/mysql"
)

func TestMigration(t *testing.T) {
    ctx := context.Background()
    
    // Start MySQL container
    mysqlContainer, err := mysql.Run(ctx,
        "mysql:8.0",
        mysql.WithDatabase("testdb"),
        mysql.WithUsername("root"),
        mysql.WithPassword("password"),
    )
    if err != nil {
        t.Fatal(err)
    }
    defer mysqlContainer.Terminate(ctx)
    
    // Get connection details
    host, _ := mysqlContainer.Host(ctx)
    port, _ := mysqlContainer.MappedPort(ctx, "3306")
    
    // Create dbmate instance
    dbURL := fmt.Sprintf("mysql://root:password@%s:%s/testdb?parseTime=true",
        host, port.Port())
    u, _ := url.Parse(dbURL)
    db := dbmate.New(u)
    db.MigrationsDir = []string{"./migrations"}
    
    // Test migration
    if err := db.CreateAndMigrate(); err != nil {
        t.Fatal(err)
    }
    
    // Verify schema
    driver, _ := db.Driver()
    sqlDB, _ := driver.Open()
    defer sqlDB.Close()
    
    rows, err := sqlDB.Query("SHOW TABLES")
    if err != nil {
        t.Fatal(err)
    }
    defer rows.Close()
    
    // Test rollback
    if err := db.Rollback(); err != nil {
        t.Fatal(err)
    }
}
```

## API Reference

### Core Methods

#### `New(databaseURL *url.URL) *DB`

Create new dbmate instance with database URL.

#### `Create() error`

Create the database if it doesn't exist.

#### `Drop() error`

Drop the database.

#### `Migrate() error`

Run all pending migrations.

#### `Rollback() error`

Rollback the most recent migration.

#### `CreateAndMigrate() error`

Create database and run all migrations (convenience method).

### Configuration

```go
type DB struct {
    DatabaseURL      *url.URL
    MigrationsDir    []string  // Migration directories
    AutoDumpSchema   bool      // Auto-dump schema after migration
    SchemaFile       string    // Schema dump file path
    WaitBefore       bool      // Wait for user input before migration
    WaitAfter        bool      // Wait for user input after migration
}
```

### Driver Methods

```go
// Get driver instance
driver, err := db.Driver()

// Open database connection
sqlDB, err := driver.Open()

// Execute query
rows, err := sqlDB.Query("SELECT * FROM users")
```

## Migration Workflow

### Development Workflow

1. **Create Migration:**
   ```bash
   # Create new migration file
   dbmate new create_users_table
   # Creates: migrations/20240117120000_create_users_table.sql
   ```

2. **Edit Migration:**
   ```sql
   -- migrate:up
   CREATE TABLE users(id INTEGER PRIMARY KEY);
   
   -- migrate:down
   DROP TABLE users;
   ```

3. **Apply Migration:**
   ```bash
   dbmate up
   ```

4. **Rollback if Needed:**
   ```bash
   dbmate rollback
   ```

### CI/CD Integration

```yaml
# .github/workflows/migrate.yml
name: Database Migration

on:
  push:
    branches: [main]

jobs:
  migrate:
    runs-on: ubuntu-latest
    
    services:
      mysql:
        image: mysql:8.0
        env:
          MYSQL_ROOT_PASSWORD: password
          MYSQL_DATABASE: myapp
        ports:
          - 3306:3306
    
    steps:
      - uses: actions/checkout@v2
      
      - name: Install dbmate
        run: |
          sudo curl -fsSL -o /usr/local/bin/dbmate \
            https://github.com/amacneil/dbmate/releases/latest/download/dbmate-linux-amd64
          sudo chmod +x /usr/local/bin/dbmate
      
      - name: Run migrations
        run: dbmate up
        env:
          DATABASE_URL: mysql://root:password@localhost:3306/myapp
```

## Best Practices

### 1. Use Timestamped Migration Names

```
20240117120000_create_users_table.sql
20240117120100_add_email_to_users.sql
20240117120200_create_orders_table.sql
```

### 2. Always Provide Down Migrations

```sql
-- migrate:up
CREATE TABLE products(id INTEGER PRIMARY KEY);

-- migrate:down
DROP TABLE products;  -- Always include rollback
```

### 3. Make Migrations Idempotent

```sql
-- migrate:up
CREATE TABLE IF NOT EXISTS users(...);  -- Use IF NOT EXISTS

-- migrate:down
DROP TABLE IF EXISTS users;  -- Use IF EXISTS
```

### 4. Test Migrations Before Production

```go
// Test both up and down migrations
db.Migrate()    // Apply
db.Rollback()   // Rollback
db.Migrate()    // Re-apply
```

### 5. Keep Migrations Small and Focused

```
❌ Bad: 20240117_big_schema_changes.sql (multiple tables)
✅ Good: 20240117120000_create_users.sql
✅ Good: 20240117120100_create_orders.sql
```

### 6. Version Control Migrations

```bash
git add migrations/
git commit -m "Add user authentication migrations"
```

## Common Migration Patterns

### Add Column

```sql
-- migrate:up
ALTER TABLE users ADD COLUMN phone VARCHAR(20);

-- migrate:down
ALTER TABLE users DROP COLUMN phone;
```

### Create Index

```sql
-- migrate:up
CREATE INDEX idx_user_email ON users(email);

-- migrate:down
DROP INDEX idx_user_email ON users;
```

### Rename Table

```sql
-- migrate:up
RENAME TABLE old_users TO users;

-- migrate:down
RENAME TABLE users TO old_users;
```

### Add Foreign Key

```sql
-- migrate:up
ALTER TABLE orders 
ADD CONSTRAINT fk_user 
FOREIGN KEY (user_id) REFERENCES users(id);

-- migrate:down
ALTER TABLE orders DROP FOREIGN KEY fk_user;
```

## Troubleshooting

### Migration Failed

```go
// Check migration status
driver, _ := db.Driver()
sqlDB, _ := driver.Open()
rows, _ := sqlDB.Query("SELECT * FROM schema_migrations")
```

### Manual Rollback

```sql
-- Delete migration record
DELETE FROM schema_migrations 
WHERE version = '20240117120000';

-- Manually drop changes
DROP TABLE IF EXISTS problem_table;
```

### Reset Database

```go
// Drop and recreate
db.Drop()
db.CreateAndMigrate()
```

## Limitations

1. **No Programmatic Migrations** - Only SQL files supported
2. **No Migration Dependencies** - Migrations run in timestamp order only
3. **No Dry Run** - Cannot preview migrations without applying
4. **Limited Validation** - Syntax errors only caught during execution
5. **No Automatic Backups** - Manual backup recommended before migration

## Dependencies

- `github.com/amacneil/dbmate/v2` - Core migration tool
- `github.com/amacneil/dbmate/v2/pkg/driver/clickhouse` - ClickHouse driver
- `github.com/amacneil/dbmate/v2/pkg/driver/mysql` - MySQL driver
- `github.com/amacneil/dbmate/v2/pkg/driver/postgres` - PostgreSQL driver

## Further Reading

- [Dbmate Documentation](https://github.com/amacneil/dbmate)
- [Migration Best Practices](https://github.com/amacneil/dbmate#best-practices)
- [Database Migration Patterns](https://www.martinfowler.com/articles/evodb.html)
- [Testcontainers Go](https://golang.testcontainers.org/)
