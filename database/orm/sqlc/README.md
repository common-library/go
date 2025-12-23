# SQLC

Generate type-safe Go code from SQL queries with compile-time verification and zero runtime overhead.

## Overview

SQLC is a SQL-first code generation tool that reads your SQL queries and schema, then generates fully type-safe Go code. It enables you to write actual SQL while benefiting from Go's type system, catching errors at compile time rather than runtime. SQLC produces zero-dependency code with no reflection and minimal runtime overhead.

**Key Features:**
- Write actual SQL, get type-safe Go
- Compile-time query verification
- Zero runtime overhead
- No reflection, no ORM magic
- Support for MySQL, PostgreSQL, SQLite
- Database-specific features and optimizations
- Named and positional parameters
- Nullable types handled correctly
- Transactions and prepared statements

## Installation

```bash
# Install sqlc CLI
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Verify installation
sqlc version

# Or use Docker
docker run --rm -v $(pwd):/src -w /src sqlcdev/sqlc generate
```

## Quick Start

### 1. Initialize Configuration

Create `sqlc.yaml` in your project root:

```yaml
version: "2"
sql:
  - engine: "mysql"
    queries: "queries/"
    schema: "schema/"
    gen:
      go:
        package: "db"
        out: "db"
        emit_json_tags: true
        emit_prepared_queries: true
        emit_interface: true
        emit_exact_table_names: false
```

### 2. Define Schema

Create `schema/schema.sql`:

```sql
-- MySQL
CREATE TABLE users (
  id            BIGINT PRIMARY KEY AUTO_INCREMENT,
  name          VARCHAR(100) NOT NULL,
  email         VARCHAR(255) NOT NULL UNIQUE,
  age           INT NOT NULL CHECK (age > 0),
  created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE posts (
  id            BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id       BIGINT NOT NULL,
  title         VARCHAR(255) NOT NULL,
  content       TEXT,
  published     BOOLEAN NOT NULL DEFAULT FALSE,
  created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  INDEX idx_user_id (user_id),
  INDEX idx_published (published)
);
```

### 3. Write Queries

Create `queries/users.sql`:

```sql
-- name: GetUser :one
SELECT * FROM users
WHERE id = ?;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC
LIMIT ?;

-- name: CreateUser :execresult
INSERT INTO users (name, email, age)
VALUES (?, ?, ?);

-- name: UpdateUser :exec
UPDATE users
SET name = ?, email = ?, age = ?
WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ?;

-- name: ListUsersByAge :many
SELECT * FROM users
WHERE age >= ?
ORDER BY age ASC;
```

### 4. Generate Code

```bash
sqlc generate
```

This generates:
```
db/
├── db.go          # DBTX interface
├── models.go      # Struct definitions
├── querier.go     # Query interface
└── users.sql.go   # Query implementations
```

### 5. Use Generated Code

```go
package main

import (
    "context"
    "database/sql"
    "log"
    
    "your-project/db"
    
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    ctx := context.Background()
    
    // Open database
    conn, err := sql.Open("mysql", "user:pass@tcp(localhost:3306)/dbname?parseTime=true")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()
    
    // Create queries
    queries := db.New(conn)
    
    // Create user
    result, err := queries.CreateUser(ctx, db.CreateUserParams{
        Name:  "Alice",
        Email: "alice@example.com",
        Age:   30,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    id, _ := result.LastInsertId()
    
    // Get user
    user, err := queries.GetUser(ctx, id)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("User: %+v", user)
}
```

## Configuration

### MySQL Configuration

```yaml
version: "2"
sql:
  - engine: "mysql"
    queries: "queries/"
    schema: "schema/"
    gen:
      go:
        package: "db"
        out: "db"
        # Generate JSON tags for structs
        emit_json_tags: true
        # Generate prepared queries methods
        emit_prepared_queries: true
        # Generate interface for all queries
        emit_interface: true
        # Use exact table names for struct types
        emit_exact_table_names: false
        # Generate methods on structs
        emit_methods_with_db_argument: false
        # Generate empty slices instead of nil
        emit_empty_slices: true
        # Generate pointers for nullable columns
        emit_pointers_for_null_types: false
        # Use specific package for SQL types
        sql_package: "database/sql"
        # Rename imports
        rename:
          id: "user_id"
        # Override types
        overrides:
          - db_type: "varchar"
            go_type: "string"
          - db_type: "text"
            go_type: "string"
            nullable: true
```

### PostgreSQL Configuration

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "queries/"
    schema: "schema/"
    gen:
      go:
        package: "pgdb"
        out: "pgdb"
        emit_json_tags: true
        emit_prepared_queries: true
        emit_interface: true
```

### Multiple Databases

```yaml
version: "2"
sql:
  - engine: "mysql"
    queries: "mysql/queries/"
    schema: "mysql/schema/"
    gen:
      go:
        package: "mysqldb"
        out: "mysqldb"
        
  - engine: "postgresql"
    queries: "postgres/queries/"
    schema: "postgres/schema/"
    gen:
      go:
        package: "pgdb"
        out: "pgdb"
```

## Query Annotations

### Return Types

```sql
-- :one - Returns single row
-- name: GetUser :one
SELECT * FROM users WHERE id = ?;

-- :many - Returns multiple rows
-- name: ListUsers :many
SELECT * FROM users;

-- :exec - Executes query, returns error
-- name: UpdateUser :exec
UPDATE users SET name = ? WHERE id = ?;

-- :execresult - Returns sql.Result
-- name: CreateUser :execresult
INSERT INTO users (name, email) VALUES (?, ?);

-- :execrows - Returns number of affected rows
-- name: DeleteOldUsers :execrows
DELETE FROM users WHERE created_at < ?;

-- :copyfrom - Batch insert (PostgreSQL only)
-- name: InsertUsers :copyfrom
INSERT INTO users (name, email, age) VALUES ($1, $2, $3);
```

### Parameters

```sql
-- Positional parameters (MySQL: ?, PostgreSQL: $1, $2)

-- MySQL
-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ?;

-- PostgreSQL
-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- Named parameters (for code clarity)
-- name: CreateUser :execresult
INSERT INTO users (
  name,
  email,
  age
) VALUES (
  ?,  -- name
  ?,  -- email
  ?   -- age
);
```

### Nullable Columns

```sql
-- Optional fields become *Type in Go
CREATE TABLE users (
  id      BIGINT PRIMARY KEY,
  name    VARCHAR(100) NOT NULL,
  bio     TEXT,              -- Nullable, becomes *string
  age     INT                -- Nullable, becomes *int32
);

-- name: GetUser :one
SELECT id, name, bio, age FROM users WHERE id = ?;

-- Generated struct:
// type User struct {
//     ID   int64
//     Name string
//     Bio  *string  // Nullable
//     Age  *int32   // Nullable
// }
```

## Generated Code

### Models

```go
// db/models.go
package db

import "time"

type User struct {
    ID        int64     `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Age       int32     `json:"age"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type Post struct {
    ID        int64     `json:"id"`
    UserID    int64     `json:"user_id"`
    Title     string    `json:"title"`
    Content   *string   `json:"content"`  // Nullable
    Published bool      `json:"published"`
    CreatedAt time.Time `json:"created_at"`
}
```

### Querier Interface

```go
// db/querier.go
package db

import "context"

type Querier interface {
    GetUser(ctx context.Context, id int64) (User, error)
    ListUsers(ctx context.Context, limit int32) ([]User, error)
    CreateUser(ctx context.Context, arg CreateUserParams) (sql.Result, error)
    UpdateUser(ctx context.Context, arg UpdateUserParams) error
    DeleteUser(ctx context.Context, id int64) error
}

var _ Querier = (*Queries)(nil)
```

### Query Implementation

```go
// db/users.sql.go
package db

import (
    "context"
    "database/sql"
)

const getUser = `-- name: GetUser :one
SELECT id, name, email, age, created_at, updated_at FROM users
WHERE id = ?
`

func (q *Queries) GetUser(ctx context.Context, id int64) (User, error) {
    row := q.db.QueryRowContext(ctx, getUser, id)
    var i User
    err := row.Scan(
        &i.ID,
        &i.Name,
        &i.Email,
        &i.Age,
        &i.CreatedAt,
        &i.UpdatedAt,
    )
    return i, err
}

type CreateUserParams struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int32  `json:"age"`
}

const createUser = `-- name: CreateUser :execresult
INSERT INTO users (name, email, age)
VALUES (?, ?, ?)
`

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (sql.Result, error) {
    return q.db.ExecContext(ctx, createUser, arg.Name, arg.Email, arg.Age)
}
```

## Common Patterns

### Transactions

```sql
-- Queries work the same in transactions

-- name: TransferMoney :exec
UPDATE accounts SET balance = balance - ? WHERE id = ?;

-- name: ReceiveMoney :exec
UPDATE accounts SET balance = balance + ? WHERE id = ?;
```

```go
func transfer(ctx context.Context, db *sql.DB, from, to int64, amount int) error {
    tx, err := db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    qtx := queries.WithTx(tx)
    
    if err := qtx.TransferMoney(ctx, TransferMoneyParams{
        Amount: amount,
        ID:     from,
    }); err != nil {
        return err
    }
    
    if err := qtx.ReceiveMoney(ctx, ReceiveMoneyParams{
        Amount: amount,
        ID:     to,
    }); err != nil {
        return err
    }
    
    return tx.Commit()
}
```

### Prepared Statements

```go
// If emit_prepared_queries is true
queries := db.New(conn)
defer queries.Close() // Close prepared statements

// Use prepared queries
user, err := queries.GetUser(ctx, 1)
```

### Pagination

```sql
-- name: ListUsersPaginated :many
SELECT * FROM users
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;
```

```go
func getPage(ctx context.Context, q *db.Queries, page, pageSize int32) ([]db.User, int64, error) {
    offset := (page - 1) * pageSize
    
    users, err := q.ListUsersPaginated(ctx, db.ListUsersPaginatedParams{
        Limit:  pageSize,
        Offset: offset,
    })
    if err != nil {
        return nil, 0, err
    }
    
    total, err := q.CountUsers(ctx)
    if err != nil {
        return nil, 0, err
    }
    
    return users, total, nil
}
```

### Complex Queries

```sql
-- name: SearchUsers :many
SELECT 
  u.*,
  COUNT(p.id) as post_count
FROM users u
LEFT JOIN posts p ON p.user_id = u.id
WHERE 
  u.name LIKE ?
  AND u.age >= ?
GROUP BY u.id
HAVING COUNT(p.id) > ?
ORDER BY post_count DESC
LIMIT ?;
```

```go
type SearchUsersRow struct {
    User      db.User
    PostCount int64
}

// Custom result scanning
func searchUsers(ctx context.Context, q *db.Queries, name string, minAge, minPosts, limit int32) ([]SearchUsersRow, error) {
    // Generated query handles complex SQL
    rows, err := q.SearchUsers(ctx, db.SearchUsersParams{
        Name:     "%" + name + "%",
        Age:      minAge,
        MinPosts: minPosts,
        Limit:    limit,
    })
    return rows, err
}
```

### Batch Operations

```sql
-- name: BatchInsertUsers :exec
INSERT INTO users (name, email, age) VALUES (?, ?, ?);
```

```go
// PostgreSQL CopyFrom (high performance)
// name: BulkInsertUsers :copyfrom
INSERT INTO users (name, email, age) VALUES ($1, $2, $3);

// Usage (PostgreSQL)
users := []db.BulkInsertUsersParams{
    {Name: "Alice", Email: "alice@example.com", Age: 30},
    {Name: "Bob", Email: "bob@example.com", Age: 25},
}
count, err := queries.BulkInsertUsers(ctx, users)
```

### Joins

```sql
-- name: GetUserWithPosts :many
SELECT 
  u.id as user_id,
  u.name as user_name,
  u.email as user_email,
  p.id as post_id,
  p.title as post_title,
  p.content as post_content
FROM users u
LEFT JOIN posts p ON p.user_id = u.id
WHERE u.id = ?;
```

```go
type UserWithPost struct {
    UserID      int64
    UserName    string
    UserEmail   string
    PostID      *int64   // Nullable from LEFT JOIN
    PostTitle   *string
    PostContent *string
}

// Process joined results
func getUserPosts(ctx context.Context, q *db.Queries, userID int64) (db.User, []db.Post, error) {
    rows, err := q.GetUserWithPosts(ctx, userID)
    if err != nil {
        return db.User{}, nil, err
    }
    
    if len(rows) == 0 {
        return db.User{}, nil, sql.ErrNoRows
    }
    
    user := db.User{
        ID:    rows[0].UserID,
        Name:  rows[0].UserName,
        Email: rows[0].UserEmail,
    }
    
    posts := make([]db.Post, 0)
    for _, row := range rows {
        if row.PostID != nil {
            posts = append(posts, db.Post{
                ID:      *row.PostID,
                Title:   *row.PostTitle,
                Content: row.PostContent,
            })
        }
    }
    
    return user, posts, nil
}
```

## Database-Specific Features

### MySQL

```sql
-- Auto increment
CREATE TABLE users (
  id BIGINT PRIMARY KEY AUTO_INCREMENT
);

-- name: CreateUser :execresult
INSERT INTO users (name) VALUES (?);

-- Get last insert ID
result, _ := q.CreateUser(ctx, "Alice")
id, _ := result.LastInsertId()

-- ON DUPLICATE KEY UPDATE
-- name: UpsertUser :exec
INSERT INTO users (email, name, age)
VALUES (?, ?, ?)
ON DUPLICATE KEY UPDATE
  name = VALUES(name),
  age = VALUES(age);

-- JSON operations
-- name: GetUserSettings :one
SELECT settings->>'$.theme' as theme FROM users WHERE id = ?;
```

### PostgreSQL

```sql
-- RETURNING clause
-- name: CreateUser :one
INSERT INTO users (name, email, age)
VALUES ($1, $2, $3)
RETURNING *;

-- Array types
CREATE TABLE users (
  id      SERIAL PRIMARY KEY,
  tags    TEXT[]
);

-- name: GetUserTags :one
SELECT tags FROM users WHERE id = $1;

-- JSONB operations
-- name: GetUserByMetadata :many
SELECT * FROM users WHERE metadata @> $1::jsonb;

-- Window functions
-- name: RankUsers :many
SELECT 
  *,
  ROW_NUMBER() OVER (ORDER BY score DESC) as rank
FROM users;

-- CTE (Common Table Expressions)
-- name: GetActiveUsers :many
WITH active_users AS (
  SELECT * FROM users WHERE last_login > NOW() - INTERVAL '30 days'
)
SELECT * FROM active_users ORDER BY last_login DESC;
```

## Type Overrides

### Custom Types

```yaml
# sqlc.yaml
version: "2"
sql:
  - engine: "mysql"
    queries: "queries/"
    schema: "schema/"
    gen:
      go:
        package: "db"
        out: "db"
        overrides:
          # Use custom type for UUID columns
          - column: "users.id"
            go_type: "github.com/google/uuid.UUID"
          
          # Use time.Duration for interval columns
          - db_type: "bigint"
            go_type: "time.Duration"
            
          # Use custom enum type
          - column: "users.role"
            go_type:
              import: "your-project/types"
              type: "UserRole"
```

### Nullable Types

```yaml
overrides:
  # Use sql.NullString instead of *string
  - db_type: "text"
    nullable: true
    go_type: "database/sql.NullString"
    
  # Use custom nullable type
  - db_type: "varchar"
    nullable: true
    go_type:
      import: "github.com/guregu/null"
      type: "String"
```

## Testing

```go
package db_test

import (
    "context"
    "database/sql"
    "testing"
    
    "your-project/db"
    
    _ "github.com/go-sql-driver/mysql"
)

func setupTestDB(t *testing.T) *sql.DB {
    conn, err := sql.Open("mysql", "test:test@tcp(localhost:3306)/testdb?parseTime=true")
    if err != nil {
        t.Fatal(err)
    }
    
    // Run migrations
    schema, _ := os.ReadFile("schema/schema.sql")
    _, err = conn.Exec(string(schema))
    if err != nil {
        t.Fatal(err)
    }
    
    return conn
}

func TestCreateAndGetUser(t *testing.T) {
    conn := setupTestDB(t)
    defer conn.Close()
    
    ctx := context.Background()
    queries := db.New(conn)
    
    // Create user
    result, err := queries.CreateUser(ctx, db.CreateUserParams{
        Name:  "Alice",
        Email: "alice@example.com",
        Age:   30,
    })
    if err != nil {
        t.Fatal(err)
    }
    
    id, _ := result.LastInsertId()
    
    // Get user
    user, err := queries.GetUser(ctx, id)
    if err != nil {
        t.Fatal(err)
    }
    
    if user.Name != "Alice" {
        t.Errorf("expected name Alice, got %s", user.Name)
    }
}

func TestTransaction(t *testing.T) {
    conn := setupTestDB(t)
    defer conn.Close()
    
    ctx := context.Background()
    
    tx, _ := conn.BeginTx(ctx, nil)
    defer tx.Rollback()
    
    qtx := db.New(tx)
    
    // Operations in transaction
    _, err := qtx.CreateUser(ctx, db.CreateUserParams{
        Name:  "Bob",
        Email: "bob@example.com",
        Age:   25,
    })
    if err != nil {
        t.Fatal(err)
    }
    
    // Commit
    if err := tx.Commit(); err != nil {
        t.Fatal(err)
    }
}
```

## Best Practices

### 1. Write Good SQL

```sql
-- Good: Specific columns
-- name: GetUserInfo :one
SELECT id, name, email FROM users WHERE id = ?;

-- Bad: SELECT *
-- name: GetUser :one
SELECT * FROM users WHERE id = ?;
```

### 2. Use Indexes

```sql
CREATE INDEX idx_user_email ON users(email);
CREATE INDEX idx_post_user_published ON posts(user_id, published);
```

### 3. Handle NULL Properly

```sql
-- Explicit NULL handling
-- name: GetUserBio :one
SELECT COALESCE(bio, '') as bio FROM users WHERE id = ?;
```

### 4. Validate in Application

```go
func createUser(ctx context.Context, q *db.Queries, name, email string, age int32) error {
    // Validate before query
    if age < 0 || age > 150 {
        return errors.New("invalid age")
    }
    
    _, err := q.CreateUser(ctx, db.CreateUserParams{
        Name:  name,
        Email: email,
        Age:   age,
    })
    return err
}
```

### 5. Use Transactions Wisely

```go
// Don't overuse transactions
func updateUser(ctx context.Context, db *sql.DB, id int64, name string) error {
    // Single operation doesn't need transaction
    return queries.UpdateUser(ctx, UpdateUserParams{ID: id, Name: name})
}

// Use for related operations
func transferMoney(ctx context.Context, db *sql.DB, from, to int64, amount int) error {
    tx, _ := db.BeginTx(ctx, nil)
    defer tx.Rollback()
    
    qtx := queries.WithTx(tx)
    qtx.Debit(ctx, DebitParams{ID: from, Amount: amount})
    qtx.Credit(ctx, CreditParams{ID: to, Amount: amount})
    
    return tx.Commit()
}
```

## Troubleshooting

### sqlc generate fails

```bash
# Verify SQL syntax
mysql < schema/schema.sql

# Check query syntax
sqlc verify

# Enable debug
sqlc generate --experimental
```

### Type mismatch errors

```yaml
# Add type override
overrides:
  - column: "users.status"
    go_type: "string"
```

### NULL handling issues

```sql
-- Use COALESCE for non-null results
SELECT COALESCE(bio, '') as bio FROM users;

-- Or handle in Go
if user.Bio != nil {
    fmt.Println(*user.Bio)
}
```

## Migration from Other Tools

### From GORM

1. Extract SQL from GORM queries
2. Define schema in SQL files
3. Write queries in SQL
4. Generate code with sqlc

### From Raw database/sql

1. Collect existing queries
2. Add sqlc annotations
3. Generate type-safe code
4. Replace manual scanning

## References

- [Official Documentation](https://docs.sqlc.dev)
- [GitHub](https://github.com/sqlc-dev/sqlc)
- [Playground](https://play.sqlc.dev)
- [Examples](https://github.com/sqlc-dev/sqlc/tree/main/examples)