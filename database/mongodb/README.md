# MongoDB

MongoDB client wrapper with simplified operations and automatic reconnection.

## Overview

The mongodb package provides a convenient wrapper around the official MongoDB Go driver, offering simplified method signatures for common operations and automatic connection management.

## Features

- **Automatic Reconnection** - Handles connection issues transparently
- **CRUD Operations** - FindOne, Find, InsertOne, InsertMany, UpdateOne, UpdateMany, DeleteOne, DeleteMany
- **Aggregation** - Pipeline-based aggregation support
- **Bulk Operations** - Efficient bulk write operations
- **Index Management** - Create and manage indexes
- **Type-Safe Results** - Generic-based type conversion
- **Context Timeout** - Configurable operation timeouts

## Installation

```bash
go get -u github.com/common-library/go/database/mongodb
go get -u go.mongodb.org/mongo-driver/mongo
```

## Quick Start

```go
import (
    "time"
    "github.com/common-library/go/database/mongodb"
    "go.mongodb.org/mongo-driver/bson"
)

type User struct {
    ID   int    `bson:"_id"`
    Name string `bson:"name"`
    Age  int    `bson:"age"`
}

func main() {
    var client mongodb.Client
    
    // Initialize with timeout
    err := client.Initialize(
        "localhost:27017",
        10*time.Second, // operation timeout
    )
    if err != nil {
        log.Fatal(err)
    }
    defer client.Finalize()
    
    // Insert document
    user := User{ID: 1, Name: "Alice", Age: 30}
    err = client.InsertOne("mydb", "users", user)
    
    // Find document
    result, err := client.FindOne(
        "mydb",
        "users",
        bson.M{"name": "Alice"},
        User{},
    )
    
    if foundUser, ok := result.(User); ok {
        fmt.Printf("Found: %+v\n", foundUser)
    }
}
```

## Basic Operations

### Insert Operations

```go
// InsertOne - Insert single document
user := User{ID: 1, Name: "Alice", Age: 30}
err := client.InsertOne("mydb", "users", user)

// InsertOne with bson.M
err = client.InsertOne("mydb", "users", bson.M{
    "_id": 2,
    "name": "Bob",
    "age": 25,
})

// InsertMany - Insert multiple documents
users := []any{
    User{ID: 3, Name: "Charlie", Age: 35},
    User{ID: 4, Name: "Diana", Age: 28},
    bson.M{"_id": 5, "name": "Eve", "age": 32},
}
err = client.InsertMany("mydb", "users", users)
```

### Find Operations

```go
// FindOne - Find single document
result, err := client.FindOne(
    "mydb",
    "users",
    bson.M{"name": "Alice"},
    User{}, // template for result type
)

if user, ok := result.(User); ok {
    fmt.Printf("User: %s, Age: %d\n", user.Name, user.Age)
}

// Find - Find multiple documents
results, err := client.Find(
    "mydb",
    "users",
    bson.M{"age": bson.M{"$gte": 25}}, // age >= 25
    User{},
)

if users, ok := results.([]User); ok {
    for _, user := range users {
        fmt.Printf("User: %s, Age: %d\n", user.Name, user.Age)
    }
}

// Find all documents
allResults, err := client.Find("mydb", "users", bson.M{}, User{})
```

### Update Operations

```go
// UpdateOne - Update single document
err := client.UpdateOne(
    "mydb",
    "users",
    bson.M{"name": "Alice"},
    bson.D{{"$set", bson.D{{"age", 31}}}},
)

// UpdateMany - Update multiple documents
err = client.UpdateMany(
    "mydb",
    "users",
    bson.M{"age": bson.M{"$lt": 30}}, // age < 30
    bson.D{{"$inc", bson.D{{"age", 1}}}}, // increment age by 1
)

// Update with multiple fields
err = client.UpdateOne(
    "mydb",
    "users",
    bson.M{"_id": 1},
    bson.D{
        {"$set", bson.D{
            {"name", "Alice Smith"},
            {"age", 32},
            {"email", "alice@example.com"},
        }},
    },
)
```

### Delete Operations

```go
// DeleteOne - Delete single document
err := client.DeleteOne(
    "mydb",
    "users",
    bson.M{"name": "Alice"},
)

// DeleteMany - Delete multiple documents
err = client.DeleteMany(
    "mydb",
    "users",
    bson.M{"age": bson.M{"$lt": 25}}, // age < 25
)

// Delete all documents in collection
err = client.DeleteMany("mydb", "users", bson.M{})
```

## Advanced Operations

### Aggregation

```go
// Aggregation pipeline
pipeline := []bson.M{
    {"$match": bson.M{"age": bson.M{"$gte": 25}}},
    {"$group": bson.M{
        "_id": "$department",
        "avgAge": bson.M{"$avg": "$age"},
        "count": bson.M{"$sum": 1},
    }},
    {"$sort": bson.M{"avgAge": -1}},
}

type AggResult struct {
    Department string  `bson:"_id"`
    AvgAge     float64 `bson:"avgAge"`
    Count      int     `bson:"count"`
}

results, err := client.Aggregate(
    "mydb",
    "users",
    pipeline,
    AggResult{},
)

if aggResults, ok := results.([]AggResult); ok {
    for _, r := range aggResults {
        fmt.Printf("Dept: %s, Avg Age: %.1f, Count: %d\n",
            r.Department, r.AvgAge, r.Count)
    }
}
```

### Bulk Operations

```go
// Bulk write operations
operations := []mongo.WriteModel{
    mongo.NewInsertOneModel().SetDocument(bson.M{
        "_id": 10,
        "name": "User10",
    }),
    mongo.NewUpdateOneModel().
        SetFilter(bson.M{"_id": 1}).
        SetUpdate(bson.D{{"$set", bson.D{{"name", "Updated"}}}}),
    mongo.NewDeleteOneModel().
        SetFilter(bson.M{"_id": 2}),
}

err := client.BulkWrite("mydb", "users", operations)
```

### Index Management

```go
// Create single field index
err := client.CreateIndex(
    "mydb",
    "users",
    "email",
    true, // unique
)

// Create compound index
err = client.CreateIndexes(
    "mydb",
    "users",
    []string{"department", "age"},
    false, // not unique
)

// Create index with options
err = client.CreateIndex("mydb", "users", "username", true)
```

## Complete Examples

### User Management System

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    "github.com/common-library/go/database/mongodb"
    "go.mongodb.org/mongo-driver/bson"
)

type User struct {
    ID        int       `bson:"_id"`
    Username  string    `bson:"username"`
    Email     string    `bson:"email"`
    Age       int       `bson:"age"`
    CreatedAt time.Time `bson:"created_at"`
}

func main() {
    var client mongodb.Client
    
    err := client.Initialize("localhost:27017", 10*time.Second)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Finalize()
    
    // Create unique index on email
    client.CreateIndex("myapp", "users", "email", true)
    
    // Register new user
    newUser := User{
        ID:        1,
        Username:  "alice",
        Email:     "alice@example.com",
        Age:       30,
        CreatedAt: time.Now(),
    }
    
    err = client.InsertOne("myapp", "users", newUser)
    if err != nil {
        log.Fatal(err)
    }
    
    // Find user by email
    result, err := client.FindOne(
        "myapp",
        "users",
        bson.M{"email": "alice@example.com"},
        User{},
    )
    
    if user, ok := result.(User); ok {
        fmt.Printf("Found user: %s (%s)\n", user.Username, user.Email)
    }
    
    // Update user age
    client.UpdateOne(
        "myapp",
        "users",
        bson.M{"email": "alice@example.com"},
        bson.D{{"$set", bson.D{{"age", 31}}}},
    )
    
    // Find all users over 25
    results, err := client.Find(
        "myapp",
        "users",
        bson.M{"age": bson.M{"$gt": 25}},
        User{},
    )
    
    if users, ok := results.([]User); ok {
        fmt.Printf("Found %d users over 25\n", len(users))
        for _, u := range users {
            fmt.Printf("- %s: %d years old\n", u.Username, u.Age)
        }
    }
}
```

### Product Catalog

```go
type Product struct {
    ID          string   `bson:"_id"`
    Name        string   `bson:"name"`
    Category    string   `bson:"category"`
    Price       float64  `bson:"price"`
    Tags        []string `bson:"tags"`
    InStock     bool     `bson:"in_stock"`
}

// Add products
products := []any{
    Product{ID: "p1", Name: "Laptop", Category: "Electronics", Price: 999.99, InStock: true},
    Product{ID: "p2", Name: "Mouse", Category: "Electronics", Price: 29.99, InStock: true},
    Product{ID: "p3", Name: "Desk", Category: "Furniture", Price: 299.99, InStock: false},
}

client.InsertMany("store", "products", products)

// Find products by category
results, _ := client.Find(
    "store",
    "products",
    bson.M{"category": "Electronics"},
    Product{},
)

// Update price
client.UpdateOne(
    "store",
    "products",
    bson.M{"_id": "p1"},
    bson.D{{"$set", bson.D{{"price", 899.99}}}},
)

// Aggregate by category
pipeline := []bson.M{
    {"$group": bson.M{
        "_id": "$category",
        "count": bson.M{"$sum": 1},
        "avgPrice": bson.M{"$avg": "$price"},
    }},
}

aggResults, _ := client.Aggregate("store", "products", pipeline, bson.M{})
```

### Event Logging

```go
type Event struct {
    ID        string    `bson:"_id"`
    UserID    int       `bson:"user_id"`
    Action    string    `bson:"action"`
    Timestamp time.Time `bson:"timestamp"`
    Metadata  bson.M    `bson:"metadata"`
}

// Log events
events := []any{
    Event{
        ID: "evt1",
        UserID: 123,
        Action: "login",
        Timestamp: time.Now(),
        Metadata: bson.M{"ip": "192.168.1.1"},
    },
    Event{
        ID: "evt2",
        UserID: 123,
        Action: "view_page",
        Timestamp: time.Now(),
        Metadata: bson.M{"page": "/dashboard"},
    },
}

client.InsertMany("logs", "events", events)

// Find user's recent events
yesterday := time.Now().Add(-24 * time.Hour)
results, _ := client.Find(
    "logs",
    "events",
    bson.M{
        "user_id": 123,
        "timestamp": bson.M{"$gte": yesterday},
    },
    Event{},
)

// Aggregate event counts by action
pipeline := []bson.M{
    {"$match": bson.M{"user_id": 123}},
    {"$group": bson.M{
        "_id": "$action",
        "count": bson.M{"$sum": 1},
    }},
}

client.Aggregate("logs", "events", pipeline, bson.M{})
```

## API Reference

### Initialization

#### `Initialize(address string, timeout time.Duration) error`

Initialize MongoDB client with connection settings.

**Parameters:**
- `address` - MongoDB server address (host:port)
- `timeout` - Operation timeout duration

**Returns:** Error if connection fails

#### `Finalize() error`

Close connection and clean up resources.

### Insert Methods

#### `InsertOne(databaseName, collectionName string, document any) error`

Insert a single document into the collection.

#### `InsertMany(databaseName, collectionName string, documents []any) error`

Insert multiple documents into the collection.

### Find Methods

#### `FindOne(databaseName, collectionName string, filter any, dataForm any) (any, error)`

Find a single document matching the filter.

**Parameters:**
- `dataForm` - Template for result type (e.g., `User{}`)

**Returns:** Interface containing the found document (type assert to access)

#### `Find(databaseName, collectionName string, filter, dataForm any) (any, error)`

Find all documents matching the filter.

**Returns:** Interface containing slice of documents

### Update Methods

#### `UpdateOne(databaseName, collectionName string, filter, update any) error`

Update a single document matching the filter.

#### `UpdateMany(databaseName, collectionName string, filter, update any) error`

Update all documents matching the filter.

### Delete Methods

#### `DeleteOne(databaseName, collectionName string, filter any) error`

Delete a single document matching the filter.

#### `DeleteMany(databaseName, collectionName string, filter any) error`

Delete all documents matching the filter.

### Advanced Methods

#### `Aggregate(databaseName, collectionName string, pipeline any, dataForm any) (any, error)`

Execute aggregation pipeline.

#### `BulkWrite(databaseName, collectionName string, operations []mongo.WriteModel) error`

Execute multiple write operations in bulk.

#### `CreateIndex(databaseName, collectionName, field string, unique bool) error`

Create an index on a single field.

#### `CreateIndexes(databaseName, collectionName string, fields []string, unique bool) error`

Create a compound index on multiple fields.

## Query Operators

Common BSON query operators:

```go
// Comparison
bson.M{"age": bson.M{"$eq": 30}}   // Equal
bson.M{"age": bson.M{"$ne": 30}}   // Not equal
bson.M{"age": bson.M{"$gt": 30}}   // Greater than
bson.M{"age": bson.M{"$gte": 30}}  // Greater or equal
bson.M{"age": bson.M{"$lt": 30}}   // Less than
bson.M{"age": bson.M{"$lte": 30}}  // Less or equal
bson.M{"age": bson.M{"$in": []int{25, 30, 35}}} // In array

// Logical
bson.M{"$and": []bson.M{
    {"age": bson.M{"$gte": 25}},
    {"age": bson.M{"$lte": 35}},
}}

bson.M{"$or": []bson.M{
    {"status": "active"},
    {"status": "pending"},
}}

// Field existence
bson.M{"email": bson.M{"$exists": true}}

// Array
bson.M{"tags": bson.M{"$all": []string{"electronics", "sale"}}}
bson.M{"tags": "electronics"} // Contains element
```

## Best Practices

### 1. Use Struct Tags

```go
type User struct {
    ID    int    `bson:"_id"`
    Name  string `bson:"name"`
    Email string `bson:"email,omitempty"`
}
```

### 2. Handle Type Assertions

```go
result, err := client.FindOne("db", "coll", filter, User{})
if err != nil {
    return err
}

user, ok := result.(User)
if !ok {
    return errors.New("type assertion failed")
}
```

### 3. Use Indexes for Performance

```go
// Create indexes on frequently queried fields
client.CreateIndex("mydb", "users", "email", true)
client.CreateIndexes("mydb", "orders", []string{"user_id", "created_at"}, false)
```

### 4. Set Appropriate Timeouts

```go
// Short timeout for simple queries
client.Initialize(addr, 5*time.Second)

// Longer timeout for complex aggregations
client.Initialize(addr, 30*time.Second)
```

### 5. Use Bulk Operations for Efficiency

```go
// Good: Single bulk operation
client.BulkWrite("db", "coll", operations)

// Avoid: Multiple individual operations
for _, op := range operations {
    client.InsertOne("db", "coll", op)
}
```

## Error Handling

```go
// Not initialized
err := client.FindOne("db", "coll", filter, User{})
// Error: "please call Initialize first"

// Document not found
result, err := client.FindOne("db", "coll", bson.M{"_id": 999}, User{})
// err: mongo: no documents in result

// Duplicate key (unique index violation)
err = client.InsertOne("db", "users", User{ID: 1, Email: "alice@example.com"})
err = client.InsertOne("db", "users", User{ID: 2, Email: "alice@example.com"})
// Error: duplicate key error
```

## Limitations

1. **No Transaction Support** - Cannot execute multi-document ACID transactions
2. **Automatic Reconnection Only** - No manual connection control
3. **Type Assertion Required** - Results returned as `any`, must type assert
4. **Limited Index Options** - Only basic unique/non-unique index creation
5. **No Change Streams** - Cannot watch for real-time changes

## Dependencies

- `go.mongodb.org/mongo-driver/mongo` - Official MongoDB Go driver
- `go.mongodb.org/mongo-driver/mongo/options` - MongoDB options
- `go.mongodb.org/mongo-driver/bson` - BSON encoding/decoding

## Further Reading

- [MongoDB Go Driver Documentation](https://www.mongodb.com/docs/drivers/go/current/)
- [BSON Package](https://pkg.go.dev/go.mongodb.org/mongo-driver/bson)
- [MongoDB Query Operators](https://www.mongodb.com/docs/manual/reference/operator/query/)
- [MongoDB Aggregation](https://www.mongodb.com/docs/manual/aggregation/)
