# Beego ORM

Full-featured ORM integrated with the Beego web framework, supporting multiple databases and QueryBuilder.

## Overview

Beego ORM is part of the Beego web framework, providing a powerful ORM with support for multiple database backends, query builder, relationship management, and transaction handling. While designed for Beego applications, it can be used standalone. The ORM follows a pragmatic approach with good performance and extensive features.

**Key Features:**
- Multiple database support (MySQL, PostgreSQL, SQLite, Oracle)
- QueryBuilder for dynamic queries
- Model relationships and lazy loading
- Automatic table creation
- Transaction support
- Raw SQL execution
- Query caching
- Command-line tools for code generation

## Installation

```bash
# Install Beego v2
go get github.com/beego/beego/v2

# Database drivers
go get github.com/go-sql-driver/mysql
go get github.com/lib/pq
go get github.com/mattn/go-sqlite3
```

## Quick Start

### 1. Register Database

```go
package main

import (
    "github.com/beego/beego/v2/client/orm"
    _ "github.com/go-sql-driver/mysql"
)

func init() {
    // Register driver
    orm.RegisterDriver("mysql", orm.DRMySQL)
    
    // Register database
    orm.RegisterDataBase("default", "mysql", 
        "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local")
    
    // Register models
    orm.RegisterModel(new(User), new(Post))
}

func main() {
    // Run ORM
    o := orm.NewOrm()
    o.Using("default")
}
```

### 2. Define Models

```go
type User struct {
    Id        int       `orm:"auto"`
    Name      string    `orm:"size(100)"`
    Email     string    `orm:"unique"`
    Age       int       
    Created   time.Time `orm:"auto_now_add;type(datetime)"`
    Updated   time.Time `orm:"auto_now;type(datetime)"`
}

// TableName specifies custom table name
func (u *User) TableName() string {
    return "users"
}
```

### 3. CRUD Operations

```go
o := orm.NewOrm()

// Create
user := User{Name: "Alice", Email: "alice@example.com", Age: 30}
id, err := o.Insert(&user)

// Read
user = User{Id: 1}
err = o.Read(&user)

// Update
user.Age = 31
num, err := o.Update(&user)

// Delete
num, err = o.Delete(&user)
```

## Model Definition

### Field Tags

```go
type User struct {
    // Auto increment primary key
    Id int `orm:"auto"`
    
    // Column with specific size
    Name string `orm:"size(100)"`
    
    // Unique constraint
    Email string `orm:"unique"`
    
    // Not null
    Password string `orm:"null"`  // Allows null
    
    // Default value
    Status int `orm:"default(1)"`
    
    // Column name
    UserName string `orm:"column(user_name)"`
    
    // Type specification
    Created time.Time `orm:"type(datetime)"`
    Content string    `orm:"type(text)"`
    
    // Precision for decimals
    Price float64 `orm:"digits(12);decimals(2)"`
    
    // Auto now (updated on every save)
    Updated time.Time `orm:"auto_now"`
    
    // Auto now add (set on creation)
    Created time.Time `orm:"auto_now_add"`
    
    // Index
    Code string `orm:"index"`
    
    // Ignore field
    IgnoredField string `orm:"-"`
}
```

### Field Types

```go
type Example struct {
    Id int `orm:"auto"`
    
    // String types
    Name    string `orm:"size(100)"` // varchar(100)
    Content string `orm:"type(text)"` // text
    
    // Integer types
    Age     int
    Counter int64
    
    // Float types
    Price   float64 `orm:"digits(12);decimals(2)"`
    
    // Boolean
    Active bool
    
    // Date/Time
    Created   time.Time `orm:"type(datetime)"`
    Birthday  time.Time `orm:"type(date)"`
    
    // JSON (stored as text)
    Metadata string `orm:"type(text)"`
}
```

### Relationships

```go
// One to One
type User struct {
    Id      int
    Name    string
    Profile *Profile `orm:"rel(one)"`
}

type Profile struct {
    Id     int
    User   *User  `orm:"reverse(one)"`
    Bio    string
}

// One to Many
type User struct {
    Id    int
    Name  string
    Posts []*Post `orm:"reverse(many)"`
}

type Post struct {
    Id     int
    User   *User  `orm:"rel(fk)"`
    Title  string
}

// Many to Many
type User struct {
    Id     int
    Name   string
    Groups []*Group `orm:"rel(m2m)"`
}

type Group struct {
    Id    int
    Name  string
    Users []*User `orm:"reverse(many)"`
}
```

## CRUD Operations

### Create

```go
o := orm.NewOrm()

// Insert single record
user := User{Name: "Alice", Email: "alice@example.com"}
id, err := o.Insert(&user)

// Insert multiple records
users := []User{
    {Name: "Bob", Email: "bob@example.com"},
    {Name: "Charlie", Email: "charlie@example.com"},
}
successNum, err := o.InsertMulti(100, users) // Batch size 100

// Insert or update
user := User{Id: 1, Name: "Alice Updated"}
id, err = o.InsertOrUpdate(&user, "name") // Update if Id exists
```

### Read

```go
o := orm.NewOrm()

// Read by primary key
user := User{Id: 1}
err := o.Read(&user)

// Read with specific fields
err = o.Read(&user, "Name", "Email")

// Read by other field
user = User{Email: "alice@example.com"}
err = o.Read(&user, "Email")

// Read or create
user = User{Email: "new@example.com"}
created, id, err := o.ReadOrCreate(&user, "Email")
if created {
    fmt.Println("Created new user")
}
```

### Update

```go
o := orm.NewOrm()

// Update all fields
user := User{Id: 1}
o.Read(&user)
user.Age = 31
num, err := o.Update(&user)

// Update specific fields
num, err = o.Update(&user, "Age", "Name")

// Raw update
num, err = o.QueryTable("user").Filter("id", 1).Update(orm.Params{
    "age":  31,
    "name": "Alice Updated",
})
```

### Delete

```go
o := orm.NewOrm()

// Delete single record
user := User{Id: 1}
num, err := o.Delete(&user)

// Delete with filter
num, err = o.QueryTable("user").Filter("age__lt", 18).Delete()
```

## QuerySeter

### Basic Queries

```go
o := orm.NewOrm()
qs := o.QueryTable("user")

// Filter
qs = qs.Filter("name", "Alice")
qs = qs.Filter("age__gt", 25)     // age > 25
qs = qs.Filter("age__gte", 25)    // age >= 25
qs = qs.Filter("age__lt", 40)     // age < 40
qs = qs.Filter("age__lte", 40)    // age <= 40
qs = qs.Filter("name__contains", "ali")      // LIKE %ali%
qs = qs.Filter("name__startswith", "Ali")    // LIKE Ali%
qs = qs.Filter("name__endswith", "ice")      // LIKE %ice
qs = qs.Filter("name__in", "Alice", "Bob")   // IN
qs = qs.Filter("age__between", 20, 30)       // BETWEEN

// Exclude
qs = qs.Exclude("age", 30)

// Or condition
cond := orm.NewCondition()
cond1 := cond.And("age__gt", 25)
cond2 := cond.Or("name", "Alice")
qs = qs.SetCond(cond1.Or(cond2))

// Limit and Offset
qs = qs.Limit(10).Offset(20)

// Order by
qs = qs.OrderBy("age")         // ASC
qs = qs.OrderBy("-created")    // DESC

// Distinct
qs = qs.Distinct()

// Execute
var users []*User
num, err := qs.All(&users)

// Count
count, err := qs.Count()

// Exists
exist := qs.Exist()

// One
var user User
err = qs.One(&user)
```

### Advanced Queries

```go
o := orm.NewOrm()

// Select specific fields
var maps []orm.Params
num, err := o.QueryTable("user").
    Filter("age__gt", 25).
    Values(&maps, "name", "email")

// Values list (single column)
var names []string
num, err = o.QueryTable("user").
    ValuesFlat(&names, "name")

// Values with aggregate
var results []orm.Params
num, err = o.QueryTable("user").
    GroupBy("age").
    Values(&results, "age", "COUNT(id)")

// Related object loading
qs.RelatedSel()  // Load all foreign keys
qs.RelatedSel("profile")  // Load specific relation

// Prepare for update/delete
num, err = o.QueryTable("user").
    Filter("age__lt", 18).
    Update(orm.Params{"active": false})
```

## QueryBuilder

### SQL Builder

```go
o := orm.NewOrm()
qb, _ := orm.NewQueryBuilder("mysql")

// Build query
qb.Select("id", "name", "age").
    From("users").
    Where("age > ?").
    And("active = ?").
    OrderBy("created_at").
    Desc().
    Limit(10).
    Offset(20)

// Get SQL
sql := qb.String()

// Execute
var users []User
o.Raw(sql, 25, true).QueryRows(&users)
```

### Complex Queries

```go
// Join query
qb.Select("u.id", "u.name", "p.title").
    From("users u").
    LeftJoin("posts p").On("u.id = p.user_id").
    Where("u.age > ?").
    GroupBy("u.id").
    Having("COUNT(p.id) > ?")

sql := qb.String()
var results []orm.Params
o.Raw(sql, 25, 5).Values(&results)

// Subquery
subQb, _ := orm.NewQueryBuilder("mysql")
subQb.Select("user_id").
    From("posts").
    Where("published = ?")

qb.Select("*").
    From("users").
    Where("id IN (" + subQb.String() + ")")

sql = qb.String()
```

## Raw SQL

### Raw Queries

```go
o := orm.NewOrm()

// Query rows
var users []User
num, err := o.Raw("SELECT * FROM users WHERE age > ?", 25).QueryRows(&users)

// Query row
var user User
err = o.Raw("SELECT * FROM users WHERE id = ?", 1).QueryRow(&user)

// Values (map slice)
var maps []orm.Params
num, err = o.Raw("SELECT name, age FROM users").Values(&maps)

// Values list (single column)
var names []string
num, err = o.Raw("SELECT name FROM users").ValuesFlat(&names)

// Exec
res, err := o.Raw("UPDATE users SET age = ? WHERE id = ?", 31, 1).Exec()
num := res.RowsAffected()

// Prepare
stmt, err := o.Raw("INSERT INTO users (name, email) VALUES (?, ?)").Prepare()
defer stmt.Close()

for _, user := range users {
    stmt.Exec(user.Name, user.Email)
}
```

### Named Parameters

```go
// Not directly supported, use placeholders
params := map[string]interface{}{
    "name": "Alice",
    "age":  30,
}

// Build query manually
o.Raw("INSERT INTO users (name, age) VALUES (?, ?)", 
    params["name"], params["age"]).Exec()
```

## Transactions

### Basic Transaction

```go
o := orm.NewOrm()

// Begin transaction
err := o.Begin()

// Operations
user := User{Name: "Alice"}
_, err = o.Insert(&user)
if err != nil {
    o.Rollback()
    return err
}

post := Post{UserId: user.Id, Title: "Hello"}
_, err = o.Insert(&post)
if err != nil {
    o.Rollback()
    return err
}

// Commit
o.Commit()
```

### Transaction with Closure

```go
o := orm.NewOrm()

err := o.DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
    // All operations use txOrm
    user := User{Name: "Alice"}
    _, err := txOrm.Insert(&user)
    if err != nil {
        return err // Auto rollback
    }
    
    post := Post{UserId: user.Id, Title: "Hello"}
    _, err = txOrm.Insert(&post)
    return err // Auto commit if nil
})
```

### Transaction Options

```go
err := o.DoTxWithCtx(ctx, &sql.TxOptions{
    Isolation: sql.LevelSerializable,
    ReadOnly:  false,
}, func(ctx context.Context, txOrm orm.TxOrmer) error {
    // Operations
    return nil
})
```

## Relationships

### Load Related

```go
o := orm.NewOrm()

// Load foreign key
user := User{Id: 1}
o.Read(&user)
o.LoadRelated(&user, "Posts")

// Load reverse relation
post := Post{Id: 1}
o.Read(&post)
o.LoadRelated(&post, "User")

// Load many to many
user = User{Id: 1}
o.Read(&user)
o.LoadRelated(&user, "Groups")

// With conditions
o.LoadRelated(&user, "Posts", 
    orm.WithOrder("-created"),
    orm.WithLimit(10))
```

### Query Related

```go
// Query through relation
o := orm.NewOrm()

user := User{Id: 1}
qs := o.QueryM2M(&user, "Groups")

// Add relation
group := Group{Id: 1}
num, err := qs.Add(&group)

// Remove relation
num, err = qs.Remove(&group)

// Clear all relations
num, err = qs.Clear()

// Count relations
count, err := qs.Count()
```

## Schema Migration

### Auto Create Tables

```go
// Create tables
err := orm.RunSyncdb("default", false, true)
// Parameters: database alias, force, verbose

// Drop and recreate (development only)
err = orm.RunSyncdb("default", true, true)
```

### Schema Commands

```go
// Generate create table SQL
sqls, err := orm.GetDB("default").DumpTables()
for _, sql := range sqls {
    fmt.Println(sql)
}
```

## Advanced Features

### Model Validation

```go
type User struct {
    Id    int
    Email string
}

func (u *User) TableName() string {
    return "users"
}

func (u *User) Valid(v *validation.Validation) {
    v.Required(u.Email, "email")
    v.Email(u.Email, "email")
}

// Use with validation
o := orm.NewOrm()
user := User{Email: "invalid"}

valid := validation.Validation{}
b, err := valid.Valid(&user)
if err != nil {
    return err
}
if !b {
    for _, err := range valid.Errors {
        log.Println(err.Key, err.Message)
    }
}
```

### Custom QuerySeter

```go
type UserQuery struct {
    qs orm.QuerySeter
}

func NewUserQuery(o orm.Ormer) *UserQuery {
    return &UserQuery{
        qs: o.QueryTable("user"),
    }
}

func (uq *UserQuery) Active() *UserQuery {
    uq.qs = uq.qs.Filter("active", true)
    return uq
}

func (uq *UserQuery) AgeGreaterThan(age int) *UserQuery {
    uq.qs = uq.qs.Filter("age__gt", age)
    return uq
}

func (uq *UserQuery) All(users *[]*User) error {
    _, err := uq.qs.All(users)
    return err
}

// Usage
o := orm.NewOrm()
var users []*User
err := NewUserQuery(o).Active().AgeGreaterThan(25).All(&users)
```

## Configuration

### Database Configuration

```go
// Max idle connections
orm.SetMaxIdleConns("default", 30)

// Max open connections
orm.SetMaxOpenConns("default", 100)

// Connection max lifetime
orm.SetConnMaxLifetime("default", time.Hour)

// Get DB object
db, err := orm.GetDB("default")
```

### Debug Mode

```go
// Enable debug (prints SQL)
orm.Debug = true

// Log queries
orm.DebugLog = orm.NewLog(os.Stdout)
```

## Testing

```go
package main_test

import (
    "testing"
    
    "github.com/beego/beego/v2/client/orm"
    _ "github.com/mattn/go-sqlite3"
)

func init() {
    orm.RegisterDriver("sqlite3", orm.DRSqlite)
    orm.RegisterDataBase("default", "sqlite3", ":memory:")
    orm.RegisterModel(new(User))
}

func TestCRUD(t *testing.T) {
    // Create tables
    orm.RunSyncdb("default", false, false)
    
    o := orm.NewOrm()
    
    // Create
    user := User{Name: "Alice", Email: "alice@example.com"}
    id, err := o.Insert(&user)
    if err != nil {
        t.Fatal(err)
    }
    
    // Read
    user = User{Id: int(id)}
    err = o.Read(&user)
    if err != nil {
        t.Fatal(err)
    }
    
    if user.Name != "Alice" {
        t.Errorf("expected Alice, got %s", user.Name)
    }
    
    // Update
    user.Age = 31
    _, err = o.Update(&user)
    if err != nil {
        t.Fatal(err)
    }
    
    // Delete
    _, err = o.Delete(&user)
    if err != nil {
        t.Fatal(err)
    }
}

func TestTransaction(t *testing.T) {
    orm.RunSyncdb("default", false, false)
    
    o := orm.NewOrm()
    
    err := o.DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
        user := User{Name: "Bob", Email: "bob@example.com"}
        _, err := txOrm.Insert(&user)
        return err
    })
    
    if err != nil {
        t.Fatal(err)
    }
}
```

## Best Practices

### 1. Register Models at Init

```go
func init() {
    orm.RegisterModel(new(User), new(Post), new(Comment))
}
```

### 2. Use Transactions for Related Operations

```go
err := o.DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
    _, err := txOrm.Insert(&user)
    if err != nil {
        return err
    }
    _, err = txOrm.Insert(&profile)
    return err
})
```

### 3. Load Related Efficiently

```go
// Bad - N+1 queries
users := []*User{}
o.QueryTable("user").All(&users)
for _, user := range users {
    o.LoadRelated(user, "Posts")
}

// Better - use RelatedSel
qs := o.QueryTable("user").RelatedSel("Posts")
qs.All(&users)
```

### 4. Use QueryBuilder for Complex Queries

```go
qb, _ := orm.NewQueryBuilder("mysql")
qb.Select("u.*, COUNT(p.id) as post_count").
    From("users u").
    LeftJoin("posts p").On("u.id = p.user_id").
    GroupBy("u.id")
    
var results []orm.Params
o.Raw(qb.String()).Values(&results)
```

### 5. Validate Input

```go
func (u *User) Valid(v *validation.Validation) {
    v.Required(u.Name, "name")
    v.MaxSize(u.Name, 100, "name")
    v.Email(u.Email, "email")
    v.Range(u.Age, 0, 150, "age")
}
```

## Common Patterns

### Repository Pattern

```go
type UserRepository struct {
    o orm.Ormer
}

func NewUserRepository() *UserRepository {
    return &UserRepository{o: orm.NewOrm()}
}

func (r *UserRepository) Create(user *User) (int64, error) {
    return r.o.Insert(user)
}

func (r *UserRepository) GetByID(id int) (*User, error) {
    user := &User{Id: id}
    err := r.o.Read(user)
    return user, err
}

func (r *UserRepository) Update(user *User) error {
    _, err := r.o.Update(user)
    return err
}

func (r *UserRepository) Delete(id int) error {
    user := &User{Id: id}
    _, err := r.o.Delete(user)
    return err
}

func (r *UserRepository) List(page, pageSize int) ([]*User, error) {
    var users []*User
    _, err := r.o.QueryTable("user").
        Limit(pageSize).
        Offset((page - 1) * pageSize).
        All(&users)
    return users, err
}
```

### Soft Delete

```go
type SoftDeleteModel struct {
    DeletedAt *time.Time `orm:"null"`
}

type User struct {
    Id int `orm:"auto"`
    Name string
    SoftDeleteModel
}

func (r *UserRepository) SoftDelete(id int) error {
    now := time.Now()
    _, err := r.o.QueryTable("user").
        Filter("id", id).
        Update(orm.Params{"deleted_at": now})
    return err
}

func (r *UserRepository) ListActive() ([]*User, error) {
    var users []*User
    _, err := r.o.QueryTable("user").
        Filter("deleted_at__isnull", true).
        All(&users)
    return users, err
}
```

## Troubleshooting

### Table not found

```go
// Ensure model is registered
orm.RegisterModel(new(User))

// Run sync
orm.RunSyncdb("default", false, true)
```

### Relationship loading fails

```go
// Check reverse relation is defined
type User struct {
    Posts []*Post `orm:"reverse(many)"`
}

type Post struct {
    User *User `orm:"rel(fk)"`
}
```

### Connection issues

```go
// Test connection
db, err := orm.GetDB("default")
if err != nil {
    log.Fatal(err)
}

err = db.Ping()
if err != nil {
    log.Fatal(err)
}
```

## References

- [Beego Documentation](https://beego.wiki)
- [ORM Documentation](https://beego.wiki/docs/mvc/model/overview/)
- [GitHub](https://github.com/beego/beego)
- [Chinese Documentation](https://beego.wiki/docs/intro/)