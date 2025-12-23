# GORM

The fantastic ORM library for Golang with full-featured associations, hooks, preloading, transactions, and auto-migration.

## Overview

GORM is a developer-friendly ORM that follows convention over configuration. It provides an ActiveRecord-style interface with automatic migrations, associations, hooks, and a rich plugin ecosystem. GORM is ideal for rapid development while maintaining flexibility for complex scenarios.

**Key Features:**
- Full-featured ORM with associations
- Auto migrations
- Hooks (Before/After Create/Save/Update/Delete/Find)
- Preloading with joins
- Transactions with nested support
- Context support
- SQL Builder for complex queries
- Logger with customizable output
- Support for MySQL, PostgreSQL, SQLite, SQL Server

## Installation

```bash
# Install GORM
go get -u gorm.io/gorm

# Install database drivers
go get -u gorm.io/driver/mysql
go get -u gorm.io/driver/postgres
go get -u gorm.io/driver/sqlite
go get -u gorm.io/driver/sqlserver
```

## Quick Start

### 1. Connect to Database

```go
import (
    "gorm.io/driver/mysql"
    "gorm.io/driver/postgres"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

// MySQL
dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

// PostgreSQL
dsn = "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable"
db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

// SQLite
db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
```

### 2. Define Models

```go
type User struct {
    ID        uint           `gorm:"primaryKey"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
    Name      string         `gorm:"size:100;not null"`
    Age       int            `gorm:"check:age > 0"`
    Email     string         `gorm:"uniqueIndex;not null"`
    Active    bool           `gorm:"default:true"`
}

// Or use gorm.Model (includes ID, CreatedAt, UpdatedAt, DeletedAt)
type Product struct {
    gorm.Model
    Code  string
    Price uint
}
```

### 3. Auto Migrate

```go
// Migrate the schema
db.AutoMigrate(&User{}, &Product{})
```

### 4. Basic CRUD

```go
// Create
user := User{Name: "Alice", Age: 30, Email: "alice@example.com"}
result := db.Create(&user)

// Read
var user User
db.First(&user, 1)                  // Find by primary key
db.First(&user, "name = ?", "Alice") // Find with condition

// Update
db.Model(&user).Update("Age", 31)
db.Model(&user).Updates(User{Age: 31, Active: true})

// Delete
db.Delete(&user, 1)
```

## Configuration

### Database Configuration

```go
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
    // Skip default transaction for write operations
    SkipDefaultTransaction: true,
    
    // Naming strategy
    NamingStrategy: schema.NamingStrategy{
        TablePrefix:   "t_",
        SingularTable: true,
    },
    
    // Logger
    Logger: logger.Default.LogMode(logger.Info),
    
    // Disable foreign key constraints
    DisableForeignKeyConstraintWhenMigrating: true,
    
    // Connection pool
    PrepareStmt: true,
})

// Get generic database object
sqlDB, err := db.DB()

// Connection pool settings
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

### Custom Logger

```go
newLogger := logger.New(
    log.New(os.Stdout, "\r\n", log.LstdFlags),
    logger.Config{
        SlowThreshold:             time.Second,
        LogLevel:                  logger.Info,
        IgnoreRecordNotFoundError: true,
        Colorful:                  true,
    },
)

db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
    Logger: newLogger,
})
```

## Model Definition

### Field Tags

```go
type User struct {
    // Primary key
    ID uint `gorm:"primaryKey"`
    
    // Column name
    Name string `gorm:"column:user_name"`
    
    // Size constraint
    Bio string `gorm:"size:500"`
    
    // Not null
    Email string `gorm:"not null"`
    
    // Unique
    Username string `gorm:"unique"`
    
    // Index
    Code string `gorm:"index"`
    
    // Composite index
    Field1 string `gorm:"index:idx_member"`
    Field2 string `gorm:"index:idx_member"`
    
    // Unique index
    Phone string `gorm:"uniqueIndex"`
    
    // Default value
    Active bool `gorm:"default:true"`
    
    // Check constraint
    Age int `gorm:"check:age > 0"`
    
    // Auto increment
    Counter int64 `gorm:"autoIncrement"`
    
    // Precision (for numeric types)
    Price float64 `gorm:"precision:2"`
    
    // Auto create/update time
    CreatedAt time.Time
    UpdatedAt time.Time
    
    // Soft delete
    DeletedAt gorm.DeletedAt `gorm:"index"`
    
    // Ignore field
    IgnoredField string `gorm:"-"`
    
    // Read/Write permission
    ReadOnlyField string `gorm:"<-:false"`  // Read-only
    CreateOnlyField string `gorm:"<-:create"` // Create-only
}
```

### Table Configuration

```go
// TableName overrides default table name
func (User) TableName() string {
    return "users"
}

// BeforeCreate hook
func (u *User) BeforeCreate(tx *gorm.DB) error {
    u.ID = uuid.New().ID()
    return nil
}
```

## CRUD Operations

### Create

```go
// Create single record
user := User{Name: "Alice", Age: 30}
result := db.Create(&user)
// user.ID is populated with generated ID

// Check errors
if result.Error != nil {
    // Handle error
}
fmt.Println(result.RowsAffected) // Rows inserted

// Create multiple records
users := []User{
    {Name: "Alice", Age: 30},
    {Name: "Bob", Age: 25},
}
db.Create(&users)

// Create with selected fields
db.Select("Name", "Age").Create(&user)

// Create with omitted fields
db.Omit("Age").Create(&user)

// Batch insert with size
db.CreateInBatches(users, 100)

// Upsert
db.Clauses(clause.OnConflict{
    Columns:   []clause.Column{{Name: "email"}},
    DoUpdates: clause.AssignmentColumns([]string{"name", "age"}),
}).Create(&user)
```

### Query

```go
// Get first record
var user User
db.First(&user)                    // ORDER BY id LIMIT 1
db.First(&user, 10)                // Find by primary key
db.First(&user, "name = ?", "Alice") // With condition

// Get last record
db.Last(&user)

// Get all records
var users []User
db.Find(&users)

// Get record with conditions
db.Where("name = ?", "Alice").First(&user)
db.Where("name <> ?", "Alice").Find(&users)
db.Where("age IN ?", []int{20, 30, 40}).Find(&users)
db.Where("age BETWEEN ? AND ?", 20, 30).Find(&users)
db.Where("name LIKE ?", "%alice%").Find(&users)

// Struct conditions
db.Where(&User{Name: "Alice", Age: 30}).First(&user)

// Map conditions
db.Where(map[string]interface{}{"name": "Alice", "age": 30}).Find(&users)

// Not conditions
db.Not("name = ?", "Alice").Find(&users)

// Or conditions
db.Where("name = ?", "Alice").Or("age = ?", 30).Find(&users)

// Select specific fields
db.Select("name", "age").Find(&users)

// Order
db.Order("age desc, name").Find(&users)

// Limit & Offset
db.Limit(10).Offset(5).Find(&users)

// Group & Having
db.Model(&User{}).Select("name, sum(age) as total").Group("name").Having("total > ?", 100).Find(&result)

// Distinct
db.Distinct("name").Find(&users)

// Joins
db.Joins("Company").Find(&users)
db.Joins("LEFT JOIN companies ON companies.id = users.company_id").Find(&users)

// Scan into custom struct
type Result struct {
    Name string
    Age  int
}
var result Result
db.Model(&User{}).Select("name, age").Where("age > ?", 25).Scan(&result)

// Count
var count int64
db.Model(&User{}).Where("age > ?", 25).Count(&count)

// Pluck (single column)
var names []string
db.Model(&User{}).Pluck("name", &names)
```

### Update

```go
// Update single column
db.Model(&user).Update("Age", 31)

// Update multiple columns (struct)
db.Model(&user).Updates(User{Age: 31, Name: "Alice Updated"})

// Update multiple columns (map)
db.Model(&user).Updates(map[string]interface{}{"Age": 31, "Active": false})

// Update selected fields
db.Model(&user).Select("Age").Updates(User{Age: 31, Name: "ignored"})

// Update omitted fields
db.Model(&user).Omit("Name").Updates(User{Age: 31, Name: "ignored"})

// Update with conditions
db.Model(&User{}).Where("age > ?", 25).Update("Active", false)

// Update with expression
db.Model(&user).Update("Age", gorm.Expr("age + ?", 1))

// Batch updates
db.Model(&User{}).Where("age > ?", 25).Updates(map[string]interface{}{"Active": false})

// Update without hooks
db.Model(&user).UpdateColumn("Age", 31)
db.Model(&user).UpdateColumns(User{Age: 31, Active: false})
```

### Delete

```go
// Delete record
db.Delete(&user, 1)

// Delete with conditions
db.Where("age > ?", 100).Delete(&User{})

// Soft delete (if model has DeletedAt field)
db.Delete(&user) // Sets DeletedAt to current time

// Find soft deleted records
db.Unscoped().Where("age > ?", 25).Find(&users)

// Permanently delete
db.Unscoped().Delete(&user)

// Delete in batch
db.Where("age > ?", 100).Delete(&User{})
```

## Associations

### Belongs To

```go
type User struct {
    gorm.Model
    Name      string
    CompanyID uint
    Company   Company `gorm:"foreignKey:CompanyID"`
}

type Company struct {
    gorm.Model
    Name string
}

// Query with association
var user User
db.Preload("Company").First(&user)

// Create with association
db.Create(&User{
    Name: "Alice",
    Company: Company{Name: "Tech Corp"},
})
```

### Has One

```go
type User struct {
    gorm.Model
    Name    string
    Profile Profile
}

type Profile struct {
    gorm.Model
    UserID uint
    Bio    string
}

// Preload
var user User
db.Preload("Profile").First(&user)
```

### Has Many

```go
type User struct {
    gorm.Model
    Name  string
    Posts []Post
}

type Post struct {
    gorm.Model
    UserID uint
    Title  string
}

// Preload
var user User
db.Preload("Posts").First(&user)

// Preload with conditions
db.Preload("Posts", "published = ?", true).First(&user)

// Nested preload
db.Preload("Posts.Comments").First(&user)
```

### Many To Many

```go
type User struct {
    gorm.Model
    Name      string
    Languages []Language `gorm:"many2many:user_languages;"`
}

type Language struct {
    gorm.Model
    Name string
}

// Create with association
user := User{
    Name: "Alice",
    Languages: []Language{
        {Name: "Go"},
        {Name: "Python"},
    },
}
db.Create(&user)

// Query with association
var user User
db.Preload("Languages").First(&user)

// Append association
var language Language
db.First(&language, "name = ?", "JavaScript")
db.Model(&user).Association("Languages").Append(&language)

// Replace association
db.Model(&user).Association("Languages").Replace(languages)

// Delete association
db.Model(&user).Association("Languages").Delete(&language)

// Clear association
db.Model(&user).Association("Languages").Clear()

// Count association
count := db.Model(&user).Association("Languages").Count()
```

### Association Mode

```go
// Find associations
var languages []Language
db.Model(&user).Association("Languages").Find(&languages)

// Append associations
db.Model(&user).Association("Languages").Append(&language1, &language2)

// Replace associations
db.Model(&user).Association("Languages").Replace(&language1, &language2)

// Delete associations
db.Model(&user).Association("Languages").Delete(&language1, &language2)

// Clear associations
db.Model(&user).Association("Languages").Clear()
```

## Hooks

### Available Hooks

```go
type User struct {
    gorm.Model
    Name string
}

// Create hooks
func (u *User) BeforeCreate(tx *gorm.DB) error {
    // Modify u before insert
    return nil
}

func (u *User) AfterCreate(tx *gorm.DB) error {
    // After insert
    return nil
}

// Update hooks
func (u *User) BeforeUpdate(tx *gorm.DB) error {
    if u.Name == "" {
        return errors.New("name cannot be empty")
    }
    return nil
}

func (u *User) AfterUpdate(tx *gorm.DB) error {
    return nil
}

// Save hooks (create or update)
func (u *User) BeforeSave(tx *gorm.DB) error {
    return nil
}

func (u *User) AfterSave(tx *gorm.DB) error {
    return nil
}

// Delete hooks
func (u *User) BeforeDelete(tx *gorm.DB) error {
    return nil
}

func (u *User) AfterDelete(tx *gorm.DB) error {
    return nil
}

// Query hooks
func (u *User) AfterFind(tx *gorm.DB) error {
    // After query
    return nil
}
```

### Modify Values in Hooks

```go
func (u *User) BeforeCreate(tx *gorm.DB) error {
    u.CreatedAt = time.Now()
    u.ID = uuid.New().ID()
    
    // Update using tx
    tx.Statement.SetColumn("Name", strings.ToLower(u.Name))
    
    return nil
}
```

## Transactions

### Manual Transaction

```go
// Begin transaction
tx := db.Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

// Operations
if err := tx.Create(&user).Error; err != nil {
    tx.Rollback()
    return err
}

if err := tx.Create(&post).Error; err != nil {
    tx.Rollback()
    return err
}

// Commit
tx.Commit()
```

### Transaction Function

```go
err := db.Transaction(func(tx *gorm.DB) error {
    // Operations within transaction
    if err := tx.Create(&user).Error; err != nil {
        return err
    }
    
    if err := tx.Create(&post).Error; err != nil {
        return err
    }
    
    // Return nil commits, return error rolls back
    return nil
})
```

### Nested Transactions

```go
db.Transaction(func(tx *gorm.DB) error {
    tx.Create(&user1)
    
    tx.Transaction(func(tx2 *gorm.DB) error {
        tx2.Create(&user2)
        return errors.New("rollback user2")
    })
    
    tx.Transaction(func(tx2 *gorm.DB) error {
        tx2.Create(&user3)
        return nil
    })
    
    return nil
})
```

### SavePoint

```go
tx := db.Begin()

tx.Create(&user1)

tx.SavePoint("sp1")
tx.Create(&user2)
tx.RollbackTo("sp1") // Rollback user2

tx.Commit() // Commits user1
```

## Advanced Features

### Scopes

```go
// Define reusable query logic
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        offset := (page - 1) * pageSize
        return db.Offset(offset).Limit(pageSize)
    }
}

func Active() func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("active = ?", true)
    }
}

// Use scopes
db.Scopes(Paginate(1, 10), Active()).Find(&users)
```

### Context Support

```go
// With context
db.WithContext(ctx).First(&user)

// Timeout context
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

db.WithContext(ctx).Find(&users)
```

### Session

```go
// Create new session
sess := db.Session(&gorm.Session{
    PrepareStmt: true,
    Context:     ctx,
})

sess.First(&user)
```

### Dry Run

```go
// Get SQL without executing
stmt := db.Session(&gorm.Session{DryRun: true}).First(&user).Statement
fmt.Println(stmt.SQL.String())
fmt.Println(stmt.Vars)
```

### Raw SQL

```go
// Query
var users []User
db.Raw("SELECT * FROM users WHERE age > ?", 25).Scan(&users)

// Exec
db.Exec("UPDATE users SET age = ? WHERE name = ?", 31, "Alice")

// Named arguments
db.Raw("SELECT * FROM users WHERE name = @name AND age = @age",
    sql.Named("name", "Alice"),
    sql.Named("age", 30),
).Scan(&users)
```

### SQL Builder

```go
// Where
db.Where("name = ?", "Alice").Where("age > ?", 25).Find(&users)

// Or
db.Where("name = ?", "Alice").Or("age = ?", 30).Find(&users)

// Not
db.Not("name = ?", "Alice").Find(&users)

// Select
db.Select("name, age").Find(&users)

// Order
db.Order("age desc").Find(&users)

// Limit
db.Limit(10).Find(&users)

// Offset
db.Offset(5).Find(&users)

// Group
db.Group("name").Having("count(*) > ?", 1).Find(&users)
```

## Performance Optimization

### Indexes

```go
type User struct {
    Name  string `gorm:"index"`
    Email string `gorm:"uniqueIndex"`
    
    // Composite index
    Field1 string `gorm:"index:idx_member"`
    Field2 string `gorm:"index:idx_member"`
    
    // Index with options
    Code string `gorm:"index:,unique,sort:desc,priority:2"`
}

// Create index in migration
db.Migrator().CreateIndex(&User{}, "Name")
db.Migrator().CreateIndex(&User{}, "idx_member")
```

### Prepared Statements

```go
// Enable globally
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
    PrepareStmt: true,
})

// Use per query
stmt := db.Session(&gorm.Session{PrepareStmt: true})
```

### Select Specific Fields

```go
// Bad - loads all fields
db.Find(&users)

// Good - only loads needed fields
db.Select("id", "name").Find(&users)
```

### Batch Loading

```go
// Bad - N+1 queries
var users []User
db.Find(&users)
for _, user := range users {
    var posts []Post
    db.Where("user_id = ?", user.ID).Find(&posts)
}

// Good - 2 queries with preload
db.Preload("Posts").Find(&users)
```

### FindInBatches

```go
result := db.Where("age > ?", 25).FindInBatches(&users, 100, func(tx *gorm.DB, batch int) error {
    for _, user := range users {
        // Process user
    }
    return nil
})
```

## Testing

```go
package main_test

import (
    "testing"
    
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        t.Fatalf("failed to connect database: %v", err)
    }
    
    db.AutoMigrate(&User{})
    return db
}

func TestCreateUser(t *testing.T) {
    db := setupTestDB(t)
    
    user := User{Name: "Alice", Age: 30}
    result := db.Create(&user)
    
    if result.Error != nil {
        t.Errorf("failed to create user: %v", result.Error)
    }
    
    if user.ID == 0 {
        t.Error("expected user ID to be set")
    }
}

func TestQueryUser(t *testing.T) {
    db := setupTestDB(t)
    
    // Create test data
    db.Create(&User{Name: "Alice", Age: 30})
    
    // Query
    var user User
    db.Where("name = ?", "Alice").First(&user)
    
    if user.Name != "Alice" {
        t.Errorf("expected name Alice, got %s", user.Name)
    }
}
```

## Best Practices

### 1. Always Check Errors

```go
// Bad
db.Create(&user)

// Good
if err := db.Create(&user).Error; err != nil {
    return err
}
```

### 2. Use Transactions for Related Operations

```go
db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(&user).Error; err != nil {
        return err
    }
    if err := tx.Create(&profile).Error; err != nil {
        return err
    }
    return nil
})
```

### 3. Preload Associations

```go
// Avoid N+1
db.Preload("Posts").Preload("Profile").Find(&users)
```

### 4. Use Scopes for Reusable Logic

```go
func ActiveUsers(db *gorm.DB) *gorm.DB {
    return db.Where("active = ?", true)
}

db.Scopes(ActiveUsers).Find(&users)
```

### 5. Configure Connection Pool

```go
sqlDB, _ := db.DB()
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

## Common Patterns

### Soft Delete

```go
type User struct {
    gorm.Model // Includes DeletedAt
}

// Soft delete
db.Delete(&user) // Sets DeletedAt

// Query excludes soft deleted
db.Find(&users)

// Include soft deleted
db.Unscoped().Find(&users)

// Permanently delete
db.Unscoped().Delete(&user)
```

### Audit Fields

```go
type User struct {
    ID        uint
    Name      string
    CreatedBy uint
    UpdatedBy uint
    CreatedAt time.Time
    UpdatedAt time.Time
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
    u.CreatedBy = getCurrentUserID()
    return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
    tx.Statement.SetColumn("UpdatedBy", getCurrentUserID())
    return nil
}
```

### Polymorphic Associations

```go
type Image struct {
    gorm.Model
    URL          string
    ImageableID   uint
    ImageableType string
}

type Product struct {
    gorm.Model
    Name   string
    Images []Image `gorm:"polymorphic:Imageable;"`
}

type User struct {
    gorm.Model
    Name   string
    Avatar []Image `gorm:"polymorphic:Imageable;"`
}
```

## Troubleshooting

### RecordNotFound Error

```go
// Check error type
err := db.First(&user).Error
if errors.Is(err, gorm.ErrRecordNotFound) {
    // Handle not found
}
```

### Slow Queries

```bash
# Enable query logging
db.Debug().Find(&users)

# Check for missing indexes
# Check for N+1 queries
```

### Migration Issues

```go
// Check if table exists
db.Migrator().HasTable(&User{})

// Drop and recreate (development only!)
db.Migrator().DropTable(&User{})
db.AutoMigrate(&User{})
```

## References

- [Official Documentation](https://gorm.io/docs/)
- [GitHub](https://github.com/go-gorm/gorm)
- [Community](https://github.com/go-gorm/gorm/discussions)