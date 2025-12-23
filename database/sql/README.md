# SQL

Unified SQL database client supporting multiple database drivers with consistent API.

## Overview

The sql package provides a single, unified interface for interacting with various SQL databases. It abstracts the differences between database drivers while maintaining compatibility with Go's standard `database/sql` package.

## Supported Databases

The package supports 7 database drivers out of the box:

| Database | Driver Constant | DSN Format |
|----------|----------------|------------|
| **MySQL** | `DriverMySQL` | `user:pass@tcp(host:port)/database` |
| **PostgreSQL** | `DriverPostgreSQL` | `postgres://user:pass@host:port/database` |
| **SQLite** | `DriverSQLite` | `file:path/to/database.db` |
| **ClickHouse** | `DriverClickHouse` | `clickhouse://host:port/database` |
| **DynamoDB** | `DriverAmazonDynamoDB` | `Region=region/TablePrefix=prefix` |
| **SQL Server** | `DriverMicrosoftSQLServer` | `sqlserver://user:pass@host:port?database=db` |
| **Oracle** | `DriverOracle` | `oracle://user:pass@host:port/service` |

## Installation

```bash
go get -u github.com/common-library/go/database/sql
```

## Quick Start

```go
import "github.com/common-library/go/database/sql"

func main() {
    var client sql.Client
    
    // Open connection
    err := client.Open(
        sql.DriverMySQL,
        "user:password@tcp(localhost:3306)/mydb",
        10, // max connections
    )
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()
    
    // Execute query
    err = client.Execute("INSERT INTO users (name, age) VALUES (?, ?)", "Alice", 30)
    if err != nil {
        log.Fatal(err)
    }
    
    // Query data
    rows, err := client.Query("SELECT id, name FROM users WHERE age > ?", 18)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    
    for rows.Next() {
        var id int
        var name string
        rows.Scan(&id, &name)
        fmt.Printf("ID: %d, Name: %s\n", id, name)
    }
}
```

## Features

### Connection Management

```go
// Open with connection pooling
err := client.Open(sql.DriverPostgreSQL, dsn, 20)

// Close connection
defer client.Close()

// Get current driver
driver := client.GetDriver()
```

### Basic Operations

#### Execute (INSERT, UPDATE, DELETE)

```go
// Insert
err := client.Execute(
    "INSERT INTO products (name, price) VALUES (?, ?)",
    "Laptop", 999.99,
)

// Update
err = client.Execute(
    "UPDATE products SET price = ? WHERE name = ?",
    899.99, "Laptop",
)

// Delete
err = client.Execute(
    "DELETE FROM products WHERE price < ?",
    100.0,
)
```

#### Query (SELECT multiple rows)

```go
rows, err := client.Query("SELECT id, name, price FROM products")
if err != nil {
    log.Fatal(err)
}
defer rows.Close()

for rows.Next() {
    var id int
    var name string
    var price float64
    
    err := rows.Scan(&id, &name, &price)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Product: %s - $%.2f\n", name, price)
}
```

#### QueryRow (SELECT single row)

```go
var name string
var age int

err := client.QueryRow(
    "SELECT name, age FROM users WHERE id = ?",
    &name, &age,
    1,
)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("User: %s, Age: %d\n", name, age)
```

### Prepared Statements

Prepared statements improve performance for repeated queries and provide SQL injection protection:

```go
// Create prepared statement
err := client.SetPrepare("INSERT INTO logs (level, message) VALUES (?, ?)")
if err != nil {
    log.Fatal(err)
}

// Execute multiple times with different parameters
client.ExecutePrepare("INFO", "Server started")
client.ExecutePrepare("DEBUG", "Connection established")
client.ExecutePrepare("WARN", "High memory usage")
```

#### Query with Prepared Statements

```go
// Prepare query
err := client.SetPrepare("SELECT id, name FROM users WHERE age > ?")

// Execute with parameters
rows, err := client.QueryPrepare(21)
defer rows.Close()

for rows.Next() {
    var id int
    var name string
    rows.Scan(&id, &name)
    fmt.Printf("Adult: %s (ID: %d)\n", name, id)
}
```

#### Query Single Row with Prepared Statement

```go
err := client.SetPrepare("SELECT balance FROM accounts WHERE id = ?")

row, err := client.QueryRowPrepare(123)
if err != nil {
    log.Fatal(err)
}

var balance float64
row.Scan(&balance)
```

### Transactions

Transactions ensure atomicity - all operations succeed or all fail:

```go
// Begin transaction
err := client.BeginTransaction()
if err != nil {
    log.Fatal(err)
}

// Execute operations
err = client.ExecuteTransaction(
    "UPDATE accounts SET balance = balance - ? WHERE id = ?",
    100.0, 1,
)
if err != nil {
    client.EndTransaction(err) // Rollback
    return err
}

err = client.ExecuteTransaction(
    "UPDATE accounts SET balance = balance + ? WHERE id = ?",
    100.0, 2,
)

// Commit if no errors, rollback if error
err = client.EndTransaction(err)
```

#### Transaction with Queries

```go
err := client.BeginTransaction()

// Query within transaction
rows, err := client.QueryTransaction(
    "SELECT id, balance FROM accounts WHERE user_id = ? FOR UPDATE",
    123,
)
if err != nil {
    client.EndTransaction(err)
    return err
}

// Process rows...
for rows.Next() {
    // ...
}
rows.Close()

// Update within transaction
err = client.ExecuteTransaction(
    "UPDATE accounts SET last_accessed = NOW() WHERE user_id = ?",
    123,
)

err = client.EndTransaction(err)
```

#### Transaction with Prepared Statements

```go
err := client.BeginTransaction()

// Prepare statement within transaction
err = client.SetPrepareTransaction(
    "INSERT INTO audit_log (action, user_id, timestamp) VALUES (?, ?, ?)",
)

// Execute multiple times
err = client.ExecutePrepareTransaction("LOGIN", 1, time.Now())
err = client.ExecutePrepareTransaction("VIEW_PROFILE", 1, time.Now())
err = client.ExecutePrepareTransaction("UPDATE_SETTINGS", 1, time.Now())

// Commit all or rollback all
err = client.EndTransaction(err)
```

## Database-Specific Examples

### MySQL

```go
var client sql.Client

err := client.Open(
    sql.DriverMySQL,
    "root:password@tcp(localhost:3306)/myapp?parseTime=true&charset=utf8mb4",
    50,
)
defer client.Close()

// MySQL-specific: REPLACE INTO
err = client.Execute(
    "REPLACE INTO cache (key_name, value, expires) VALUES (?, ?, ?)",
    "session:123", "data", time.Now().Add(1*time.Hour),
)
```

### PostgreSQL

```go
var client sql.Client

err := client.Open(
    sql.DriverPostgreSQL,
    "postgres://user:pass@localhost:5432/mydb?sslmode=disable",
    25,
)
defer client.Close()

// PostgreSQL-specific: RETURNING clause
rows, err := client.Query(
    "INSERT INTO users (name) VALUES ($1) RETURNING id",
    "Alice",
)
defer rows.Close()

var newID int
if rows.Next() {
    rows.Scan(&newID)
    fmt.Printf("New user ID: %d\n", newID)
}
```

### SQLite

```go
var client sql.Client

err := client.Open(
    sql.DriverSQLite,
    "file:./myapp.db?cache=shared&mode=rwc",
    1, // SQLite: single connection recommended
)
defer client.Close()

// SQLite-specific: Enable foreign keys
err = client.Execute("PRAGMA foreign_keys = ON")
```

### ClickHouse

```go
var client sql.Client

err := client.Open(
    sql.DriverClickHouse,
    "clickhouse://localhost:9000/analytics",
    10,
)
defer client.Close()

// ClickHouse-specific: Batch insert
err = client.BeginTransaction()
err = client.SetPrepareTransaction(
    "INSERT INTO events (timestamp, user_id, event_type) VALUES (?, ?, ?)",
)

for _, event := range events {
    client.ExecutePrepareTransaction(
        event.Timestamp,
        event.UserID,
        event.Type,
    )
}

client.EndTransaction(err)
```

## API Reference

### Connection Methods

#### `Open(driver Driver, dsn string, maxOpenConnection int) error`

Establishes a connection to the database with connection pooling.

**Parameters:**
- `driver` - Database driver constant (e.g., DriverMySQL)
- `dsn` - Data Source Name (connection string)
- `maxOpenConnection` - Maximum number of open connections

**Returns:** Error if connection fails

#### `Close() error`

Closes the database connection. Safe to call multiple times.

**Returns:** Error if closing fails

#### `GetDriver() Driver`

Returns the current database driver being used.

**Returns:** Driver constant

### Query Methods

#### `Query(query string, args ...any) (*sql.Rows, error)`

Executes a query that returns multiple rows.

**Parameters:**
- `query` - SQL query string (use ? for placeholders)
- `args` - Query parameters

**Returns:** Rows and error

#### `QueryRow(query string, result ...any) error`

Executes a query and scans the first row into provided variables.

**Parameters:**
- `query` - SQL query string
- `result` - Pointers to destination variables

**Returns:** Error if query or scan fails

#### `Execute(query string, args ...any) error`

Executes a statement that doesn't return rows (INSERT, UPDATE, DELETE).

**Parameters:**
- `query` - SQL statement string
- `args` - Statement parameters

**Returns:** Error if execution fails

### Prepared Statement Methods

#### `SetPrepare(query string) error`

Creates a prepared statement for repeated execution.

#### `QueryPrepare(args ...any) (*sql.Rows, error)`

Executes prepared statement and returns rows.

#### `QueryRowPrepare(args ...any) (*sql.Row, error)`

Executes prepared statement and returns single row.

#### `ExecutePrepare(args ...any) error`

Executes prepared statement (INSERT/UPDATE/DELETE).

### Transaction Methods

#### `BeginTransaction() error`

Starts a new transaction.

#### `EndTransaction(err error) error`

Commits transaction if err is nil, otherwise rolls back.

#### `QueryTransaction(query string, args ...any) (*sql.Rows, error)`

Executes query within current transaction.

#### `QueryRowTransaction(query string, result ...any) error`

Executes query and scans first row within transaction.

#### `ExecuteTransaction(query string, args ...any) error`

Executes statement within current transaction.

### Transaction Prepared Statement Methods

#### `SetPrepareTransaction(query string) error`

Creates prepared statement within transaction.

#### `QueryPrepareTransaction(args ...any) (*sql.Rows, error)`

Executes transaction prepared statement and returns rows.

#### `QueryRowPrepareTransaction(args ...any) (*sql.Row, error)`

Executes transaction prepared statement and returns single row.

#### `ExecutePrepareTransaction(args ...any) error`

Executes transaction prepared statement.

## Best Practices

### 1. Always Close Resources

```go
// Close client
defer client.Close()

// Close rows
rows, err := client.Query("SELECT * FROM users")
defer rows.Close()
```

### 2. Use Prepared Statements for Repeated Queries

```go
// Good: Reuse prepared statement
client.SetPrepare("INSERT INTO logs (message) VALUES (?)")
for _, msg := range messages {
    client.ExecutePrepare(msg)
}

// Avoid: Preparing each time
for _, msg := range messages {
    client.Execute("INSERT INTO logs (message) VALUES (?)", msg)
}
```

### 3. Use Transactions for Related Operations

```go
// Transfer money between accounts - atomic operation
err := client.BeginTransaction()
err = client.ExecuteTransaction("UPDATE accounts SET balance = balance - ? WHERE id = ?", 100, 1)
err = client.ExecuteTransaction("UPDATE accounts SET balance = balance + ? WHERE id = ?", 100, 2)
err = client.EndTransaction(err) // Commits both or rolls back both
```

### 4. Handle Errors Properly

```go
err := client.BeginTransaction()
if err != nil {
    return fmt.Errorf("failed to begin transaction: %w", err)
}

err = client.ExecuteTransaction("INSERT INTO users (name) VALUES (?)", name)
if err != nil {
    client.EndTransaction(err) // Rollback
    return fmt.Errorf("failed to insert user: %w", err)
}

err = client.EndTransaction(nil) // Commit
if err != nil {
    return fmt.Errorf("failed to commit transaction: %w", err)
}
```

### 5. Use Parameter Placeholders to Prevent SQL Injection

```go
// Good: Parameterized query
client.Query("SELECT * FROM users WHERE name = ?", userName)

// Bad: String concatenation (SQL injection risk!)
// client.Query("SELECT * FROM users WHERE name = '" + userName + "'")
```

### 6. Configure Connection Pool Appropriately

```go
// Web application: Higher concurrency
client.Open(sql.DriverMySQL, dsn, 100)

// Background worker: Lower concurrency
client.Open(sql.DriverMySQL, dsn, 10)

// SQLite: Single connection
client.Open(sql.DriverSQLite, dsn, 1)
```

## Error Handling

Common errors and how to handle them:

```go
// Connection not opened
err := client.Query("SELECT * FROM users")
// Error: "please call Open first"

// Prepared statement not set
err := client.ExecutePrepare("value")
// Error: "please call SetPrepare first"

// Transaction not started
err := client.ExecuteTransaction("INSERT ...")
// Error: "please call BeginTransaction first"
```

## Performance Tips

1. **Use Prepared Statements** - 10-50% faster for repeated queries
2. **Batch Operations in Transactions** - Reduces network round-trips
3. **Set Appropriate Pool Size** - Balance between resource usage and concurrency
4. **Close Rows Promptly** - Prevents connection pool exhaustion
5. **Use QueryRow for Single Row** - More efficient than Query + Next

## Migration from Standard library

```go
// Standard library
db, err := sql.Open("mysql", dsn)
stmt, err := db.Prepare("INSERT INTO users (name) VALUES (?)")
stmt.Exec("Alice")

// This package
var client sql.Client
client.Open(sql.DriverMySQL, dsn, 10)
client.SetPrepare("INSERT INTO users (name) VALUES (?)")
client.ExecutePrepare("Alice")
```

## Limitations

1. **No ORM Features** - Raw SQL only, no object mapping
2. **Limited to Supported Drivers** - Cannot add custom drivers without code changes
3. **Single Active Transaction** - One transaction per client instance
4. **No Query Builder** - Manual SQL string construction required

## Dependencies

- `database/sql` - Go standard library
- `github.com/ClickHouse/clickhouse-go/v2` - ClickHouse driver
- `github.com/btnguyen2k/godynamo` - DynamoDB driver
- `github.com/go-sql-driver/mysql` - MySQL driver
- `github.com/lib/pq` - PostgreSQL driver
- `github.com/microsoft/go-mssqldb` - SQL Server driver
- `github.com/sijms/go-ora` - Oracle driver
- `modernc.org/sqlite` - SQLite driver

## Related Packages

- [../orm/](../orm/) - ORM tools (GORM, Ent, SQLC, SQLx)
- [../dbmate/](../dbmate/) - Database migration management

## Further Reading

- [Go database/sql tutorial](https://go.dev/doc/database/overview)
- [Prepared statements explained](https://go.dev/doc/database/prepared-statements)
- [Managing connections](https://go.dev/doc/database/manage-connections)
- [SQL transactions in Go](https://go.dev/doc/database/execute-transactions)
