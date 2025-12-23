# Ent ORM

Graph-based entity framework for Go with type-safe schema definition and code generation.

## Overview

Ent is a modern ORM framework that treats your data as a graph. It provides compile-time type safety through code generation and enables powerful graph traversal queries. The framework uses a code-first approach where you define schemas in Go, and Ent generates all necessary code including clients, builders, and predicates.

**Key Features:**
- Schema as code with field builders
- Automatic code generation for type-safe operations
- Graph-based query traversal with edges
- Built-in schema migration support
- Privacy layer for access control
- Hook system for lifecycle events
- Support for MySQL, PostgreSQL, SQLite, Gremlin

## Architecture

```
ent/
├── schema/          # Schema definitions
│   └── user.go      # Entity schema
├── generate.go      # Code generation directive
├── client.go        # Generated client
├── user/            # Generated user files
│   ├── user.go      # User constants and predicates
│   └── where.go     # Query predicates
├── user_create.go   # Create builder
├── user_update.go   # Update builder
├── user_query.go    # Query builder
└── user_delete.go   # Delete builder
```

## Installation

```bash
# Install Ent CLI
go get entgo.io/ent/cmd/ent

# Install Ent library
go get entgo.io/ent

# Install database drivers
go get github.com/go-sql-driver/mysql
go get github.com/lib/pq
go get github.com/mattn/go-sqlite3
```

## Quick Start

### 1. Initialize Ent

```bash
# Create ent directory
go run entgo.io/ent/cmd/ent new User
```

This generates:
```
ent/
├── generate.go
└── schema/
    └── user.go
```

### 2. Define Schema

```go
// ent/schema/user.go
package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/edge"
)

type User struct {
    ent.Schema
}

// Fields defines the entity fields
func (User) Fields() []ent.Field {
    return []ent.Field{
        field.String("name").
            NotEmpty().
            MaxLen(100),
        field.Int("age").
            Positive(),
        field.String("email").
            Optional().
            Unique(),
        field.Time("created_at").
            Default(time.Now),
    }
}

// Edges defines relationships
func (User) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("posts", Post.Type),
        edge.From("groups", Group.Type).
            Ref("users"),
    }
}
```

### 3. Generate Code

```bash
# Generate entities
go generate ./ent

# Or with custom template
go generate -mod=mod ./ent
```

### 4. Create Client

```go
package main

import (
    "context"
    "log"

    "<your-project>/ent"
    
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    client, err := ent.Open("mysql", "user:pass@tcp(localhost:3306)/dbname?parseTime=True")
    if err != nil {
        log.Fatalf("failed opening connection: %v", err)
    }
    defer client.Close()
    
    // Run migration
    if err := client.Schema.Create(context.Background()); err != nil {
        log.Fatalf("failed creating schema: %v", err)
    }
}
```

## Core Operations

### Create Entity

```go
// Simple create
user, err := client.User.
    Create().
    SetName("Alice").
    SetAge(30).
    SetEmail("alice@example.com").
    Save(ctx)

// Bulk create
users, err := client.User.CreateBulk(
    client.User.Create().SetName("Bob").SetAge(25),
    client.User.Create().SetName("Charlie").SetAge(35),
).Save(ctx)

// Upsert (create or update)
err := client.User.
    Create().
    SetName("Alice").
    SetAge(30).
    OnConflict().
    UpdateNewValues().
    Exec(ctx)
```

### Query Entity

```go
// Get by ID
user, err := client.User.Get(ctx, id)

// Query with predicates
users, err := client.User.
    Query().
    Where(user.AgeGT(25)).
    Where(user.NameHasPrefix("A")).
    All(ctx)

// Query with pagination
users, err := client.User.
    Query().
    Limit(10).
    Offset(20).
    Order(ent.Asc(user.FieldName)).
    All(ctx)

// Count
count, err := client.User.
    Query().
    Where(user.AgeGT(25)).
    Count(ctx)

// Exists
exists, err := client.User.
    Query().
    Where(user.Email("alice@example.com")).
    Exist(ctx)

// First/Only
user, err := client.User.
    Query().
    Where(user.Name("Alice")).
    First(ctx)  // Returns first or error

user, err := client.User.
    Query().
    Where(user.Name("Alice")).
    Only(ctx)  // Returns error if != 1 result
```

### Update Entity

```go
// Update one
err := client.User.
    UpdateOneID(id).
    SetAge(31).
    AddAge(1).  // Increment
    Exec(ctx)

// Update many
affected, err := client.User.
    Update().
    Where(user.AgeGT(30)).
    SetEmail("old@example.com").
    Save(ctx)

// Update with mutation
err := client.User.
    UpdateOne(user).
    SetAge(user.Age + 1).
    Exec(ctx)
```

### Delete Entity

```go
// Delete one
err := client.User.
    DeleteOneID(id).
    Exec(ctx)

// Delete many
affected, err := client.User.
    Delete().
    Where(user.AgeGT(100)).
    Exec(ctx)
```

## Graph Traversal

### Define Relationships

```go
// Post schema
type Post struct {
    ent.Schema
}

func (Post) Fields() []ent.Field {
    return []ent.Field{
        field.String("title"),
        field.Text("content"),
    }
}

func (Post) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("author", User.Type).
            Ref("posts").
            Unique().
            Required(),
        edge.To("comments", Comment.Type),
    }
}
```

### Query with Edges

```go
// Query with edge loading
users, err := client.User.
    Query().
    Where(user.AgeGT(25)).
    WithPosts().  // Eager load posts
    All(ctx)

// Access loaded edges
for _, u := range users {
    posts := u.Edges.Posts
    fmt.Printf("%s has %d posts\n", u.Name, len(posts))
}

// Query through edges
posts, err := client.Post.
    Query().
    Where(post.HasAuthorWith(user.Name("Alice"))).
    All(ctx)

// Nested edge loading
users, err := client.User.
    Query().
    WithPosts(func(q *ent.PostQuery) {
        q.WithComments()  // Load posts with comments
    }).
    All(ctx)

// Query edges directly
posts, err := user.
    QueryPosts().
    Where(post.TitleContains("Go")).
    All(ctx)
```

### Complex Graph Queries

```go
// Multi-level traversal
users, err := client.User.
    Query().
    Where(user.HasPostsWith(
        post.HasCommentsWith(
            comment.ContentContains("awesome"),
        ),
    )).
    All(ctx)

// Query aggregation
var v []struct {
    Name  string `json:"name"`
    Count int    `json:"count"`
}

err := client.User.
    Query().
    GroupBy(user.FieldName).
    Aggregate(ent.Count()).
    Scan(ctx, &v)
```

## Schema Definition

### Field Types

```go
func (User) Fields() []ent.Field {
    return []ent.Field{
        // String types
        field.String("name").
            NotEmpty().
            MaxLen(100).
            MinLen(2),
            
        // Numeric types
        field.Int("age").
            Positive().
            Min(0).
            Max(150),
        field.Float("balance").
            Default(0.0),
            
        // Boolean
        field.Bool("active").
            Default(true),
            
        // Time
        field.Time("created_at").
            Default(time.Now).
            Immutable(),
        field.Time("updated_at").
            Default(time.Now).
            UpdateDefault(time.Now),
            
        // Enum
        field.Enum("status").
            Values("active", "inactive", "suspended").
            Default("active"),
            
        // JSON
        field.JSON("metadata", map[string]interface{}{}).
            Optional(),
            
        // UUID
        field.UUID("id", uuid.UUID{}).
            Default(uuid.New),
            
        // Binary
        field.Bytes("avatar").
            Optional(),
    }
}
```

### Field Validators

```go
func (User) Fields() []ent.Field {
    return []ent.Field{
        field.String("email").
            Validate(func(s string) error {
                if !strings.Contains(s, "@") {
                    return errors.New("invalid email")
                }
                return nil
            }),
            
        field.Int("age").
            Validate(func(i int) error {
                if i < 0 || i > 150 {
                    return errors.New("invalid age")
                }
                return nil
            }),
    }
}
```

### Indexes

```go
func (User) Indexes() []ent.Index {
    return []ent.Index{
        // Single field index
        index.Fields("email").
            Unique(),
            
        // Composite index
        index.Fields("first_name", "last_name"),
        
        // Index with edges
        index.Fields("created_at").
            Edges("posts"),
    }
}
```

## Migration

### Automatic Migration

```go
// Create all tables
err := client.Schema.Create(ctx)

// Create with options
err := client.Schema.Create(
    ctx,
    migrate.WithDropIndex(true),
    migrate.WithDropColumn(true),
)

// Dry run
err := client.Schema.WriteTo(ctx, os.Stdout)
```

### Versioned Migration

```go
// Generate migration files
err := client.Schema.Create(
    ctx,
    schema.WithDir("migrations"),
    schema.WithMigrationMode(schema.ModeReplay),
)

// Apply migrations
err := client.Schema.Create(
    ctx,
    schema.WithDir("migrations"),
    schema.WithMigrationMode(schema.ModeReplay),
    schema.WithDropColumn(true),
)
```

### Custom Migration

```go
// Add custom SQL
func (User) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{
            Table: "users",
            Checks: map[string]string{
                "age_check": "age >= 0 AND age <= 150",
            },
        },
    }
}
```

## Transactions

### Basic Transaction

```go
tx, err := client.Tx(ctx)
if err != nil {
    return err
}

user, err := tx.User.
    Create().
    SetName("Alice").
    Save(ctx)
if err != nil {
    return rollback(tx, err)
}

post, err := tx.Post.
    Create().
    SetTitle("Hello").
    SetAuthor(user).
    Save(ctx)
if err != nil {
    return rollback(tx, err)
}

return tx.Commit()

// Helper
func rollback(tx *ent.Tx, err error) error {
    if rerr := tx.Rollback(); rerr != nil {
        err = fmt.Errorf("%w: %v", err, rerr)
    }
    return err
}
```

### Transaction Function

```go
func WithTx(ctx context.Context, client *ent.Client, fn func(tx *ent.Tx) error) error {
    tx, err := client.Tx(ctx)
    if err != nil {
        return err
    }
    defer func() {
        if v := recover(); v != nil {
            tx.Rollback()
            panic(v)
        }
    }()
    if err := fn(tx); err != nil {
        if rerr := tx.Rollback(); rerr != nil {
            err = fmt.Errorf("%w: rolling back transaction: %v", err, rerr)
        }
        return err
    }
    if err := tx.Commit(); err != nil {
        return fmt.Errorf("committing transaction: %w", err)
    }
    return nil
}

// Usage
err := WithTx(ctx, client, func(tx *ent.Tx) error {
    _, err := tx.User.Create().SetName("Alice").Save(ctx)
    return err
})
```

## Hooks

### Mutation Hooks

```go
// Schema hook
func (User) Hooks() []ent.Hook {
    return []ent.Hook{
        // On create
        hook.On(
            func(next ent.Mutator) ent.Mutator {
                return hook.UserFunc(func(ctx context.Context, m *ent.UserMutation) (ent.Value, error) {
                    name, _ := m.Name()
                    m.SetName(strings.ToLower(name))
                    return next.Mutate(ctx, m)
                })
            },
            ent.OpCreate,
        ),
        
        // On update
        hook.On(
            func(next ent.Mutator) ent.Mutator {
                return hook.UserFunc(func(ctx context.Context, m *ent.UserMutation) (ent.Value, error) {
                    m.SetUpdatedAt(time.Now())
                    return next.Mutate(ctx, m)
                })
            },
            ent.OpUpdate|ent.OpUpdateOne,
        ),
    }
}
```

### Client Hooks

```go
// Global hooks
client.Use(func(next ent.Mutator) ent.Mutator {
    return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
        start := time.Now()
        defer func() {
            log.Printf("Op=%s Type=%s Time=%s", m.Op(), m.Type(), time.Since(start))
        }()
        return next.Mutate(ctx, m)
    })
})
```

## Privacy Layer

### Define Privacy Rules

```go
func (User) Policy() ent.Policy {
    return privacy.Policy{
        Mutation: privacy.MutationPolicy{
            privacy.OnMutationOperation(
                privacy.DenyIfNoViewer(),
                ent.OpCreate|ent.OpUpdate|ent.OpUpdateOne,
            ),
            privacy.AlwaysAllowRule(),
        },
        Query: privacy.QueryPolicy{
            privacy.AlwaysAllowRule(),
        },
    }
}
```

### Custom Privacy Rules

```go
// Viewer context
func AllowIfAdmin() privacy.MutationRule {
    return privacy.MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
        viewer := viewer.FromContext(ctx)
        if viewer.Admin {
            return privacy.Allow
        }
        return privacy.Deny
    })
}

// Field-level privacy
func DenyIfNotOwner() privacy.QueryRule {
    return privacy.QueryRuleFunc(func(ctx context.Context, q ent.Query) error {
        viewer := viewer.FromContext(ctx)
        if uq, ok := q.(*ent.UserQuery); ok {
            uq.Where(user.ID(viewer.UserID))
        }
        return privacy.Skip
    })
}
```

## Testing

```go
package ent_test

import (
    "context"
    "testing"
    
    "<your-project>/ent"
    "<your-project>/ent/enttest"
    
    _ "github.com/mattn/go-sqlite3"
)

func TestUser(t *testing.T) {
    // Open in-memory SQLite
    client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
    defer client.Close()
    
    ctx := context.Background()
    
    // Create user
    user := client.User.
        Create().
        SetName("Alice").
        SetAge(30).
        SaveX(ctx)  // X suffix panics on error
    
    // Query
    users := client.User.Query().AllX(ctx)
    if len(users) != 1 {
        t.Errorf("expected 1 user, got %d", len(users))
    }
}
```

## Best Practices

### 1. Use Code Generation

Always regenerate after schema changes:
```bash
go generate ./ent
```

### 2. Handle Errors Properly

```go
// Bad
user := client.User.GetX(ctx, id)  // Panics on error

// Good
user, err := client.User.Get(ctx, id)
if err != nil {
    if ent.IsNotFound(err) {
        // Handle not found
    }
    return err
}
```

### 3. Use Transactions for Related Operations

```go
err := WithTx(ctx, client, func(tx *ent.Tx) error {
    user, err := tx.User.Create().SetName("Alice").Save(ctx)
    if err != nil {
        return err
    }
    _, err = tx.Post.Create().SetAuthor(user).Save(ctx)
    return err
})
```

### 4. Eager Load Related Entities

```go
// Bad - N+1 queries
users, _ := client.User.Query().All(ctx)
for _, u := range users {
    posts, _ := u.QueryPosts().All(ctx)  // N queries
}

// Good - 2 queries
users, _ := client.User.Query().WithPosts().All(ctx)
for _, u := range users {
    posts := u.Edges.Posts  // Already loaded
}
```

### 5. Use Predicates for Complex Queries

```go
predicates := []predicate.User{
    user.AgeGT(25),
}

if nameFilter != "" {
    predicates = append(predicates, user.NameContains(nameFilter))
}

users, err := client.User.Query().Where(predicates...).All(ctx)
```

## Performance Tips

1. **Use Batch Operations:** `CreateBulk` for multiple inserts
2. **Select Specific Fields:** `Select(user.FieldName, user.FieldAge)`
3. **Limit Query Depth:** Avoid deep nested edge loading
4. **Use Indexes:** Define indexes on frequently queried fields
5. **Connection Pooling:** Configure `SetMaxOpenConns` and `SetMaxIdleConns`

## Common Patterns

### Soft Delete

```go
func (User) Mixin() []ent.Mixin {
    return []ent.Mixin{
        mixin.Time{},
        SoftDeleteMixin{},
    }
}

// Query non-deleted
users, _ := client.User.Query().Where(user.DeletedAtIsNil()).All(ctx)
```

### Audit Log

```go
func AuditHook() ent.Hook {
    return func(next ent.Mutator) ent.Mutator {
        return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
            // Log mutation
            log.Printf("Op=%s Type=%s", m.Op(), m.Type())
            return next.Mutate(ctx, m)
        })
    }
}
```

## Troubleshooting

### Code Generation Fails

```bash
# Clean and regenerate
rm -rf ent/
go run entgo.io/ent/cmd/ent new User
go generate ./ent
```

### Migration Issues

```bash
# Drop and recreate (development only)
client.Schema.Create(ctx, migrate.WithDropColumn(true), migrate.WithDropIndex(true))
```

### Performance Issues

- Enable query logging: `client.Debug()`
- Check for N+1 queries
- Add indexes on WHERE clause fields
- Use `Select()` to limit returned fields

## References

- [Official Documentation](https://entgo.io/docs/getting-started)
- [Schema Guide](https://entgo.io/docs/schema-def)
- [Migration](https://entgo.io/docs/migrate)
- [Hooks](https://entgo.io/docs/hooks)
- [Privacy](https://entgo.io/docs/privacy)
- [GitHub](https://github.com/ent/ent)