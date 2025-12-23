# ORM

Object-Relational Mapping (ORM) tools and query builders for Go.

## Overview

The orm package provides working examples and comprehensive documentation for five popular Go database abstraction libraries. Each tool represents a different philosophy for database interaction, from full-featured ORMs to lightweight query builders.

## Supported Tools

| Tool | Type | Documentation | Best For |
|------|------|---------------|----------|
| **[Ent](ent/)** | ORM | [📖 Guide](ent/README.md) | Complex relationships, type safety |
| **[GORM](gorm/)** | ORM | [📖 Guide](gorm/README.md) | Rapid development, familiar API |
| **[SQLC](sqlc/)** | Code Generator | [📖 Guide](sqlc/README.md) | Type-safe SQL, compile-time safety |
| **[SQLx](sqlx/)** | Query Builder | [📖 Guide](sqlx/README.md) | Minimal abstraction, SQL control |
| **[Beego](beego/)** | ORM | [📖 Guide](beego/README.md) | Beego framework integration |

## Quick Comparison

### Philosophy

| Tool | Approach | Schema Definition | Query Style |
|------|----------|-------------------|-------------|
| **Ent** | Code-first | Go structs with field builders | Graph-based, type-safe |
| **GORM** | Convention over config | Go structs with tags | ActiveRecord-style |
| **SQLC** | SQL-first | SQL schema files | Raw SQL with codegen |
| **SQLx** | Database/sql wrapper | Manual SQL | Raw SQL with scanning |
| **Beego** | Framework-integrated | Go structs with tags | QueryBuilder + ORM |

### Feature Matrix

| Feature | Ent | GORM | SQLC | SQLx | Beego |
|---------|-----|------|------|------|-------|
| Type Safety | ✅✅✅ | ✅✅ | ✅✅✅ | ✅ | ✅✅ |
| Learning Curve | Steep | Easy | Medium | Easy | Medium |
| Code Generation | ✅ | ❌ | ✅ | ❌ | ❌ |
| Auto Migration | ✅ | ✅ | ❌ | ❌ | ✅ |
| Relationships | ✅✅✅ | ✅✅✅ | Manual | Manual | ✅✅ |
| Performance | ✅✅ | ✅ | ✅✅✅ | ✅✅✅ | ✅✅ |
| SQL Control | ✅ | ✅ | ✅✅✅ | ✅✅✅ | ✅✅ |
| Community Size | Medium | Large | Medium | Large | Small |

## Decision Guide

### Choose Ent When

- **Complex relationships:** Your domain has many entity relationships
- **Type safety critical:** Compile-time guarantees are essential
- **Graph queries:** You need to traverse relationships efficiently
- **Schema evolution:** Automatic migration generation is valuable
- **Greenfield projects:** Starting fresh with code-first approach

[📖 Read the Ent Guide](ent/README.md)

### Choose GORM When

- **Rapid development:** Need to build quickly with minimal boilerplate
- **Familiar API:** Coming from Rails/ActiveRecord background
- **Rich features:** Need hooks, associations, auto-migration out of the box
- **Large community:** Want extensive documentation and plugins
- **Learning curve:** Need something easy to pick up

[📖 Read the GORM Guide](gorm/README.md)

### Choose SQLC When

- **SQL expertise:** You know SQL and want to write it directly
- **Performance critical:** Need zero runtime overhead
- **Type safety:** Want compile-time query verification
- **Database-specific:** Need to use database-specific features
- **Simple queries:** Mostly straightforward CRUD operations

[📖 Read the SQLC Guide](sqlc/README.md)

### Choose SQLx When

- **Minimal abstraction:** Want thin layer over database/sql
- **Full control:** Need complete control over SQL
- **Small projects:** Don't need heavy ORM machinery
- **Learning SQL:** Want to improve SQL skills
- **Existing code:** Migrating from database/sql

[📖 Read the SQLx Guide](sqlx/README.md)

### Choose Beego ORM When

- **Beego framework:** Already using Beego web framework
- **Framework integration:** Want tight web framework coupling
- **QueryBuilder needed:** Mix of ORM and query builder
- **Multi-database:** Need to support multiple databases easily
- **China-focused:** Development team is China-based

[📖 Read the Beego Guide](beego/README.md)

## Getting Started

See individual guides for detailed installation, configuration, and usage:

- **[Ent Guide](ent/README.md)** - Schema definition, code generation, graph queries
- **[GORM Guide](gorm/README.md)** - Models, associations, migrations, hooks
- **[SQLC Guide](sqlc/README.md)** - SQL queries, code generation, configuration
- **[SQLx Guide](sqlx/README.md)** - Struct scanning, named queries, transactions
- **[Beego Guide](beego/README.md)** - ORM features, QueryBuilder, framework integration

## Performance Considerations

| Tool | Query Performance | Memory Usage | Startup Time |
|------|------------------|--------------|--------------|
| Ent | Good | Medium | Medium (codegen) |
| GORM | Fair | High | Fast |
| SQLC | Excellent | Low | Medium (codegen) |
| SQLx | Excellent | Low | Fast |
| Beego | Good | Medium | Fast |

**Recommendations:**
- **High performance:** SQLC or SQLx for minimal overhead
- **Balanced:** Ent for type safety with good performance
- **Rapid development:** GORM for quick iteration
- **Control:** SQLx for manual optimization

## Common Patterns

### Connection Management

```go
// Ent
client, err := ent.Open("mysql", dsn)
defer client.Close()

// GORM
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
sqlDB, _ := db.DB()
defer sqlDB.Close()

// SQLC/SQLx
db, err := sqlx.Connect("mysql", dsn)
defer db.Close()

// Beego
orm.RegisterDataBase("default", "mysql", dsn)
```

### Transaction Handling

```go
// Ent
tx, err := client.Tx(ctx)
if err := txFunc(ctx, tx); err != nil {
    tx.Rollback()
} else {
    tx.Commit()
}

// GORM
db.Transaction(func(tx *gorm.DB) error {
    // operations
    return nil
})

// SQLx
tx, _ := db.Beginx()
defer tx.Rollback()
// operations
tx.Commit()
```

## Migration Between Tools

### From database/sql to SQLx
Easiest migration - SQLx extends database/sql interface.

### From GORM to Ent
Requires schema redesign but improves type safety.

### From Raw SQL to SQLC
Keep your SQL, gain type safety through codegen.

### From Any Tool to SQLx
Always possible - SQLx works with any database/sql driver.

## Next Steps

Choose a tool from the decision guide above and read its detailed documentation. Each guide includes installation, complete API reference, best practices, testing examples, and troubleshooting.

## Resources

- **Ent:** [Documentation](https://entgo.io) | [GitHub](https://github.com/ent/ent)
- **GORM:** [Documentation](https://gorm.io) | [GitHub](https://github.com/go-gorm/gorm)
- **SQLC:** [Documentation](https://docs.sqlc.dev) | [GitHub](https://github.com/sqlc-dev/sqlc)
- **SQLx:** [Documentation](http://jmoiron.github.io/sqlx/) | [GitHub](https://github.com/jmoiron/sqlx)
- **Beego:** [Documentation](https://beego.wiki) | [GitHub](https://github.com/beego/beego)
