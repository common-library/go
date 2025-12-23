# SQLx

Extensions to Go's database/sql package for easier querying and scanning.

## Overview

SQLx is a library that provides a set of extensions on top of the excellent built-in database/sql package. It maintains the standard database/sql interfaces while adding convenient methods for common operations like scanning into structs, named queries, and easier transaction handling. SQLx is perfect when you want minimal abstraction over SQL while reducing boilerplate.

**Key Features:**
- Thin layer over database/sql
- Scan query results into structs
- Named parameter support
- In queries with slices
- Safer connection handling
- Minimal learning curve
- Zero dependencies
- Support for MySQL, PostgreSQL, SQLite, SQL Server, Oracle

## Installation

```bash
go get github.com/jmoiron/sqlx

# Database drivers (same as database/sql)
go get github.com/go-sql-driver/mysql
go get github.com/lib/pq
go get github.com/mattn/go-sqlite3
```

## Quick Start

### 1. Connect to Database

```go
import (
    "github.com/jmoiron/sqlx"
    _ "github.com/go-sql-driver/mysql"
)

// Connect and verify
db, err := sqlx.Connect("mysql", "user:pass@tcp(localhost:3306)/dbname?parseTime=true")
if err != nil {
    log.Fatal(err)
}
defer db.Close()

// Or open without ping
db, err := sqlx.Open("mysql", dsn)
defer db.Close()

// Verify connection
if err := db.Ping(); err != nil {
    log.Fatal(err)
}
```

### 2. Define Struct

```go
type User struct {
    ID        int64     `db:"id"`
    Name      string    `db:"name"`
    Email     string    `db:"email"`
    Age       int       `db:"age"`
    CreatedAt time.Time `db:"created_at"`
}
```

### 3. Query Examples

```go
// Get single row
var user User
err := db.Get(&user, "SELECT * FROM users WHERE id = ?", 1)

// Get multiple rows
var users []User
err := db.Select(&users, "SELECT * FROM users WHERE age > ?", 25)

// Execute query
result, err := db.Exec("INSERT INTO users (name, email, age) VALUES (?, ?, ?)", 
    "Alice", "alice@example.com", 30)

// Named query
user := User{Name: "Bob", Email: "bob@example.com", Age: 25}
result, err := db.NamedExec(`INSERT INTO users (name, email, age) VALUES (:name, :email, :age)`, user)
```

## Core Operations

### Query Single Row

```go
// Get scans a single row into a struct
var user User
err := db.Get(&user, "SELECT * FROM users WHERE id = ?", 1)
if err != nil {
    if err == sql.ErrNoRows {
        // Handle not found
    }
    return err
}

// GetContext with timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

err = db.GetContext(ctx, &user, "SELECT * FROM users WHERE email = ?", "alice@example.com")

// QueryRowx for manual scanning
row := db.QueryRowx("SELECT * FROM users WHERE id = ?", 1)
err = row.StructScan(&user)
```

### Query Multiple Rows

```go
// Select scans multiple rows into a slice
var users []User
err := db.Select(&users, "SELECT * FROM users WHERE age > ? ORDER BY created_at DESC", 25)

// Empty slice if no results
fmt.Println(len(users)) // 0 if no rows

// SelectContext with context
err = db.SelectContext(ctx, &users, "SELECT * FROM users WHERE active = ?", true)

// Queryx for manual scanning
rows, err := db.Queryx("SELECT * FROM users WHERE age > ?", 25)
defer rows.Close()

for rows.Next() {
    var user User
    if err := rows.StructScan(&user); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%+v\n", user)
}

// MapScan for map results
for rows.Next() {
    result := make(map[string]interface{})
    if err := rows.MapScan(result); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%v\n", result)
}
```

### Execute Queries

```go
// Exec for INSERT, UPDATE, DELETE
result, err := db.Exec("INSERT INTO users (name, email, age) VALUES (?, ?, ?)",
    "Alice", "alice@example.com", 30)

// Get last insert ID
id, err := result.LastInsertId()

// Get affected rows
affected, err := result.RowsAffected()

// ExecContext with context
result, err = db.ExecContext(ctx, "UPDATE users SET age = ? WHERE id = ?", 31, 1)

// MustExec panics on error (useful for setup)
db.MustExec("CREATE TABLE IF NOT EXISTS users (id INT PRIMARY KEY)")
```

## Named Queries

### Named Parameters

```go
// Struct-based named query
user := User{
    Name:  "Alice",
    Email: "alice@example.com",
    Age:   30,
}

result, err := db.NamedExec(`
    INSERT INTO users (name, email, age)
    VALUES (:name, :email, :age)
`, user)

// Map-based named query
params := map[string]interface{}{
    "name":  "Bob",
    "email": "bob@example.com",
    "age":   25,
}

result, err = db.NamedExec(`
    INSERT INTO users (name, email, age)
    VALUES (:name, :email, :age)
`, params)

// Named query for SELECT
var users []User
err = db.NamedSelect(&users, `
    SELECT * FROM users
    WHERE age > :min_age AND age < :max_age
`, map[string]interface{}{
    "min_age": 20,
    "max_age": 40,
})
```

### Prepared Named Statements

```go
// Prepare named statement
stmt, err := db.PrepareNamed(`
    INSERT INTO users (name, email, age)
    VALUES (:name, :email, :age)
`)
defer stmt.Close()

// Execute multiple times
users := []User{
    {Name: "Alice", Email: "alice@example.com", Age: 30},
    {Name: "Bob", Email: "bob@example.com", Age: 25},
}

for _, user := range users {
    result, err := stmt.Exec(user)
    if err != nil {
        log.Fatal(err)
    }
}

// Query with prepared statement
stmt, err = db.PrepareNamed("SELECT * FROM users WHERE age > :age")
defer stmt.Close()

var users []User
err = stmt.Select(&users, map[string]interface{}{"age": 25})
```

## IN Queries

### Slice Expansion

```go
// IN query with slice
ids := []int{1, 2, 3, 4, 5}

query, args, err := sqlx.In("SELECT * FROM users WHERE id IN (?)", ids)
if err != nil {
    log.Fatal(err)
}

// Rebind for your database (?, $1, etc.)
query = db.Rebind(query)

var users []User
err = db.Select(&users, query, args...)

// Named IN query
arg := map[string]interface{}{
    "ids": []int{1, 2, 3, 4, 5},
}

query, args, err = sqlx.Named("SELECT * FROM users WHERE id IN (:ids)", arg)
query, args, err = sqlx.In(query, args...)
query = db.Rebind(query)

err = db.Select(&users, query, args...)
```

### Multiple IN Clauses

```go
// Multiple IN conditions
ids := []int{1, 2, 3}
names := []string{"Alice", "Bob", "Charlie"}

query, args, err := sqlx.In(`
    SELECT * FROM users
    WHERE id IN (?) AND name IN (?)
`, ids, names)

query = db.Rebind(query)
err = db.Select(&users, query, args...)
```

## Transactions

### Basic Transaction

```go
// Begin transaction
tx, err := db.Beginx()
if err != nil {
    return err
}
defer tx.Rollback() // Rollback if not committed

// Execute in transaction
_, err = tx.Exec("INSERT INTO users (name, email) VALUES (?, ?)", "Alice", "alice@example.com")
if err != nil {
    return err
}

var user User
err = tx.Get(&user, "SELECT * FROM users WHERE email = ?", "alice@example.com")
if err != nil {
    return err
}

// Commit transaction
if err = tx.Commit(); err != nil {
    return err
}
```

### Transaction with Context

```go
ctx := context.Background()

tx, err := db.BeginTxx(ctx, &sql.TxOptions{
    Isolation: sql.LevelSerializable,
})
if err != nil {
    return err
}
defer tx.Rollback()

// Operations
_, err = tx.ExecContext(ctx, "UPDATE accounts SET balance = balance - ? WHERE id = ?", 100, 1)
if err != nil {
    return err
}

_, err = tx.ExecContext(ctx, "UPDATE accounts SET balance = balance + ? WHERE id = ?", 100, 2)
if err != nil {
    return err
}

return tx.Commit()
```

### Named Queries in Transactions

```go
tx, err := db.Beginx()
defer tx.Rollback()

user := User{Name: "Alice", Email: "alice@example.com", Age: 30}

result, err := tx.NamedExec(`
    INSERT INTO users (name, email, age)
    VALUES (:name, :email, :age)
`, user)
if err != nil {
    return err
}

id, _ := result.LastInsertId()

var createdUser User
err = tx.Get(&createdUser, "SELECT * FROM users WHERE id = ?", id)
if err != nil {
    return err
}

return tx.Commit()
```

## Prepared Statements

### Prepared Queries

```go
// Prepare statement
stmt, err := db.Preparex("SELECT * FROM users WHERE age > ?")
if err != nil {
    return err
}
defer stmt.Close()

// Use multiple times
var users []User
err = stmt.Select(&users, 25)

var moreUsers []User
err = stmt.Select(&moreUsers, 30)

// Get single row
var user User
err = stmt.Get(&user, 25)
```

### Prepared Statements in Transactions

```go
tx, err := db.Beginx()
defer tx.Rollback()

stmt, err := tx.Preparex("INSERT INTO users (name, email) VALUES (?, ?)")
defer stmt.Close()

for _, user := range users {
    _, err = stmt.Exec(user.Name, user.Email)
    if err != nil {
        return err
    }
}

return tx.Commit()
```

## Advanced Features

### Unsafe Mode

```go
// Skip field validation (matches columns by position)
db = db.Unsafe()

// Useful for SELECT * queries
var users []User
err := db.Select(&users, "SELECT * FROM users")

// Or per-query
err = db.Unsafe().Select(&users, "SELECT id, name, email FROM users")
```

### Custom Mapping

```go
// Custom mapper function
db.MapperFunc(strings.ToUpper)

// Now matches uppercase struct tags
type User struct {
    ID   int    `db:"ID"`
    Name string `db:"NAME"`
}
```

### StructScan with Embedded Structs

```go
type Base struct {
    ID        int64     `db:"id"`
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}

type User struct {
    Base
    Name  string `db:"name"`
    Email string `db:"email"`
}

// Scans into embedded struct fields
var user User
err := db.Get(&user, "SELECT * FROM users WHERE id = ?", 1)
fmt.Println(user.ID, user.Name) // Embedded field accessible
```

### LoadFile for Batch Operations

```go
// Execute SQL file
schema, err := ioutil.ReadFile("schema.sql")
if err != nil {
    log.Fatal(err)
}

db.MustExec(string(schema))
```

## Connection Pool Configuration

```go
db, err := sqlx.Connect("mysql", dsn)

// Set max open connections
db.SetMaxOpenConns(25)

// Set max idle connections
db.SetMaxIdleConns(5)

// Set connection max lifetime
db.SetConnMaxLifetime(5 * time.Minute)

// Set connection max idle time
db.SetConnMaxIdleTime(5 * time.Minute)

// Get underlying *sql.DB if needed
sqlDB := db.DB
```

## Error Handling

```go
// Check for specific errors
var user User
err := db.Get(&user, "SELECT * FROM users WHERE id = ?", 999)

if err != nil {
    switch err {
    case sql.ErrNoRows:
        // Handle not found
        return nil, ErrUserNotFound
    case sql.ErrConnDone:
        // Connection closed
        return nil, ErrDatabaseClosed
    default:
        // Other errors
        return nil, err
    }
}

// Check for constraint violations (driver-specific)
if err := db.Exec("INSERT INTO users (email) VALUES (?)", "duplicate@example.com"); err != nil {
    if strings.Contains(err.Error(), "Duplicate entry") {
        return ErrDuplicateEmail
    }
    return err
}
```

## Testing

```go
package main_test

import (
    "testing"
    
    "github.com/jmoiron/sqlx"
    _ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sqlx.DB {
    db, err := sqlx.Connect("sqlite3", ":memory:")
    if err != nil {
        t.Fatal(err)
    }
    
    schema := `
    CREATE TABLE users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        email TEXT UNIQUE NOT NULL,
        age INTEGER NOT NULL
    )`
    
    db.MustExec(schema)
    return db
}

func TestCreateUser(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()
    
    result, err := db.Exec("INSERT INTO users (name, email, age) VALUES (?, ?, ?)",
        "Alice", "alice@example.com", 30)
    if err != nil {
        t.Fatal(err)
    }
    
    id, _ := result.LastInsertId()
    if id == 0 {
        t.Error("expected ID > 0")
    }
}

func TestGetUser(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()
    
    // Insert test data
    db.MustExec("INSERT INTO users (name, email, age) VALUES (?, ?, ?)",
        "Alice", "alice@example.com", 30)
    
    // Query
    var user User
    err := db.Get(&user, "SELECT * FROM users WHERE email = ?", "alice@example.com")
    if err != nil {
        t.Fatal(err)
    }
    
    if user.Name != "Alice" {
        t.Errorf("expected name Alice, got %s", user.Name)
    }
}

func TestTransaction(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()
    
    tx, err := db.Beginx()
    if err != nil {
        t.Fatal(err)
    }
    defer tx.Rollback()
    
    _, err = tx.Exec("INSERT INTO users (name, email, age) VALUES (?, ?, ?)",
        "Bob", "bob@example.com", 25)
    if err != nil {
        t.Fatal(err)
    }
    
    if err := tx.Commit(); err != nil {
        t.Fatal(err)
    }
    
    var count int
    db.Get(&count, "SELECT COUNT(*) FROM users")
    if count != 1 {
        t.Errorf("expected 1 user, got %d", count)
    }
}
```

## Best Practices

### 1. Always Use Context

```go
// Good: with context
ctx := context.Background()
err := db.GetContext(ctx, &user, query, args...)

// Bad: no context
err := db.Get(&user, query, args...)
```

### 2. Close Resources

```go
// Close DB
defer db.Close()

// Close rows
rows, err := db.Queryx(query)
defer rows.Close()

// Close statements
stmt, err := db.Preparex(query)
defer stmt.Close()
```

### 3. Handle sql.ErrNoRows

```go
var user User
err := db.Get(&user, "SELECT * FROM users WHERE id = ?", id)
if err != nil {
    if err == sql.ErrNoRows {
        return nil, ErrNotFound
    }
    return nil, err
}
```

### 4. Use Transactions for Multiple Operations

```go
func transferMoney(db *sqlx.DB, from, to int, amount int) error {
    tx, err := db.Beginx()
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    _, err = tx.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?", amount, from)
    if err != nil {
        return err
    }
    
    _, err = tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, to)
    if err != nil {
        return err
    }
    
    return tx.Commit()
}
```

### 5. Prepare Statements for Repeated Queries

```go
// Prepare once, use many times
stmt, err := db.Preparex("SELECT * FROM users WHERE age > ?")
defer stmt.Close()

for _, age := range ages {
    var users []User
    stmt.Select(&users, age)
    // Process users
}
```

## Common Patterns

### Repository Pattern

```go
type UserRepository struct {
    db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*User, error) {
    var user User
    err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = ?", id)
    if err == sql.ErrNoRows {
        return nil, ErrNotFound
    }
    return &user, err
}

func (r *UserRepository) Create(ctx context.Context, user *User) error {
    result, err := r.db.NamedExecContext(ctx, `
        INSERT INTO users (name, email, age)
        VALUES (:name, :email, :age)
    `, user)
    if err != nil {
        return err
    }
    
    id, _ := result.LastInsertId()
    user.ID = id
    return nil
}

func (r *UserRepository) Update(ctx context.Context, user *User) error {
    _, err := r.db.NamedExecContext(ctx, `
        UPDATE users
        SET name = :name, email = :email, age = :age
        WHERE id = :id
    `, user)
    return err
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
    _, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = ?", id)
    return err
}

func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]User, error) {
    var users []User
    err := r.db.SelectContext(ctx, &users, `
        SELECT * FROM users
        ORDER BY created_at DESC
        LIMIT ? OFFSET ?
    `, limit, offset)
    return users, err
}
```

### Pagination Helper

```go
type Page struct {
    Items      interface{}
    Page       int
    PageSize   int
    TotalItems int64
    TotalPages int
}

func Paginate(db *sqlx.DB, query string, dest interface{}, page, pageSize int, args ...interface{}) (*Page, error) {
    // Count total
    countQuery := "SELECT COUNT(*) FROM (" + query + ") AS count_table"
    var total int64
    err := db.Get(&total, countQuery, args...)
    if err != nil {
        return nil, err
    }
    
    // Get page
    offset := (page - 1) * pageSize
    paginatedQuery := query + " LIMIT ? OFFSET ?"
    args = append(args, pageSize, offset)
    
    err = db.Select(dest, paginatedQuery, args...)
    if err != nil {
        return nil, err
    }
    
    totalPages := int(total) / pageSize
    if int(total)%pageSize != 0 {
        totalPages++
    }
    
    return &Page{
        Items:      dest,
        Page:       page,
        PageSize:   pageSize,
        TotalItems: total,
        TotalPages: totalPages,
    }, nil
}
```

### Query Builder

```go
type QueryBuilder struct {
    query  strings.Builder
    args   []interface{}
    argIdx int
}

func NewQueryBuilder() *QueryBuilder {
    return &QueryBuilder{}
}

func (qb *QueryBuilder) Select(columns string) *QueryBuilder {
    qb.query.WriteString("SELECT " + columns)
    return qb
}

func (qb *QueryBuilder) From(table string) *QueryBuilder {
    qb.query.WriteString(" FROM " + table)
    return qb
}

func (qb *QueryBuilder) Where(condition string, args ...interface{}) *QueryBuilder {
    if qb.argIdx == 0 {
        qb.query.WriteString(" WHERE ")
    } else {
        qb.query.WriteString(" AND ")
    }
    qb.query.WriteString(condition)
    qb.args = append(qb.args, args...)
    qb.argIdx += len(args)
    return qb
}

func (qb *QueryBuilder) Build() (string, []interface{}) {
    return qb.query.String(), qb.args
}

// Usage
qb := NewQueryBuilder().
    Select("*").
    From("users").
    Where("age > ?", 25).
    Where("active = ?", true)

query, args := qb.Build()
var users []User
db.Select(&users, query, args...)
```

## Performance Tips

1. **Use connection pooling:** Configure `SetMaxOpenConns` and `SetMaxIdleConns`
2. **Prepare repeated queries:** Use `Preparex` for queries executed multiple times
3. **Batch operations:** Use transactions for bulk inserts/updates
4. **Select specific columns:** Avoid `SELECT *` when possible
5. **Use indexes:** Create indexes on frequently queried columns
6. **Profile queries:** Use database query profiling tools

## Troubleshooting

### Field not found

```go
// Ensure db tags match column names
type User struct {
    ID   int    `db:"id"`     // Column name in database
    Name string `db:"name"`
}

// Or use Unsafe() to skip validation
db.Unsafe().Get(&user, query)
```

### Connection closed

```go
// Check connection before use
if err := db.Ping(); err != nil {
    // Reconnect
}

// Or use connection pool properly
db.SetConnMaxLifetime(5 * time.Minute)
```

### Transaction already committed

```go
// Don't commit twice
tx, _ := db.Beginx()
defer tx.Rollback() // Safe even after Commit

if err := doWork(tx); err != nil {
    return err // Rollback via defer
}

return tx.Commit() // Only commits once
```

## Migration from database/sql

SQLx is a drop-in replacement:

```go
// database/sql
import "database/sql"
db, err := sql.Open("mysql", dsn)

rows, err := db.Query("SELECT id, name FROM users")
defer rows.Close()
for rows.Next() {
    var id int
    var name string
    rows.Scan(&id, &name)
}

// SQLx - same interface + extensions
import "github.com/jmoiron/sqlx"
db, err := sqlx.Open("mysql", dsn)

var users []User
err = db.Select(&users, "SELECT id, name FROM users")
// No manual scanning!
```

## References

- [Documentation](http://jmoiron.github.io/sqlx/)
- [GitHub](https://github.com/jmoiron/sqlx)
- [Illustrated Guide](http://jmoiron.github.io/sqlx/)