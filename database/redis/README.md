# Redis

Redis client with connection pooling and simplified operations.

## Overview

The redis package provides a convenient wrapper around the Redigo Redis client library, offering connection pooling, automatic health checks, and simplified method signatures for common Redis operations.

## Features

- **Connection Pooling** - Configurable pool size and idle timeout
- **Automatic Health Checks** - PING-based connection validation
- **String Operations** - GET, SET, SETEX, MGET, MSET
- **Key Management** - DEL, EXISTS, RENAME, TTL, RANDOMKEY
- **Database Selection** - SELECT command support
- **Database Operations** - FLUSHDB, FLUSHALL, DBSIZE
- **Server Info** - INFO command with category filtering

## Installation

```bash
go get -u github.com/common-library/go/database/redis
go get -u github.com/gomodule/redigo/redis
```

## Quick Start

```go
import (
    "time"
    "github.com/common-library/go/database/redis"
)

func main() {
    var client redis.Client
    
    // Initialize with connection pool
    err := client.Initialize(
        "localhost:6379",  // address
        "",                // password (empty for no auth)
        10,                // max idle connections
        60*time.Second,    // idle timeout
    )
    if err != nil {
        log.Fatal(err)
    }
    defer client.Finalize()
    
    // Set and get values
    err = client.Set("user:1:name", "Alice")
    name, err := client.Get("user:1:name")
    fmt.Println(name) // Alice
}
```

## Basic Operations

### String Operations

```go
// SET - Store a value
err := client.Set("key", "value")

// GET - Retrieve a value  
value, err := client.Get("key")
fmt.Println(value) // "value"

// SETEX - Set with expiration (in seconds)
err = client.Setex("session:123", 3600, "session_data")

// MGET - Get multiple values
values, err := client.MGet("key1", "key2", "key3")
for _, val := range values {
    fmt.Println(val)
}

// MSET - Set multiple values
err = client.MSet("key1", "value1", "key2", "value2", "key3", "value3")
```

### Key Management

```go
// DEL - Delete a key
err := client.Del("key")

// EXISTS - Check if key exists
exists, err := client.Exists("key")
if exists {
    fmt.Println("Key exists")
}

// EXISTS - Check multiple keys
exists, err = client.Exists("key1", "key2", "key3")

// RENAME - Rename a key
err = client.Rename("old_key", "new_key")

// RANDOMKEY - Get random key
key, err := client.RandomKey()

// TTL - Get time to live (in seconds)
ttl, err := client.Ttl("session:123")
if ttl == -2 {
    fmt.Println("Key does not exist")
} else if ttl == -1 {
    fmt.Println("Key has no expiration")
} else {
    fmt.Printf("Key expires in %d seconds\n", ttl)
}
```

### Database Operations

```go
// SELECT - Select database by index (0-15)
err := client.Select(1)

// DBSIZE - Get number of keys in current database
count, err := client.DBsize()
fmt.Printf("Database has %d keys\n", count)

// FLUSHDB - Delete all keys in current database
err = client.FlushDB()

// FLUSHALL - Delete all keys in all databases
err = client.FlushAll()
```

### Server Information

```go
// INFO - Get all server information
info, err := client.Info("ALL")
fmt.Println(info)

// INFO - Get specific category
serverInfo, err := client.Info("Server")
memoryInfo, err := client.Info("Memory")
statsInfo, err := client.Info("Stats")

// Available categories:
// - Server, Clients, Memory, Persistence
// - Stats, Replication, CPU, Cluster, Keyspace
```

### Connection Management

```go
// PING - Test connection
err := client.Ping()
if err != nil {
    fmt.Println("Redis connection failed")
}

// Finalize - Close all connections
defer client.Finalize()
```

## Complete Examples

### Session Storage

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    "github.com/common-library/go/database/redis"
)

func main() {
    var client redis.Client
    
    err := client.Initialize("localhost:6379", "", 20, 120*time.Second)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Finalize()
    
    // Store session with 1 hour expiration
    sessionID := "session:user123"
    sessionData := `{"user_id": 123, "name": "Alice", "role": "admin"}`
    
    err = client.Setex(sessionID, 3600, sessionData)
    if err != nil {
        log.Fatal(err)
    }
    
    // Check if session exists
    exists, err := client.Exists(sessionID)
    if exists {
        // Retrieve session
        data, err := client.Get(sessionID)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("Session data: %s\n", data)
        
        // Check remaining time
        ttl, _ := client.Ttl(sessionID)
        fmt.Printf("Session expires in %d seconds\n", ttl)
    }
    
    // Delete session on logout
    err = client.Del(sessionID)
}
```

### Caching

```go
func getCachedUser(client *redis.Client, userID int) (User, error) {
    cacheKey := fmt.Sprintf("user:%d", userID)
    
    // Try to get from cache
    cached, err := client.Get(cacheKey)
    if err == nil && cached != "" {
        var user User
        json.Unmarshal([]byte(cached), &user)
        return user, nil
    }
    
    // Cache miss - fetch from database
    user, err := fetchUserFromDB(userID)
    if err != nil {
        return User{}, err
    }
    
    // Store in cache with 5 minute expiration
    userData, _ := json.Marshal(user)
    client.Setex(cacheKey, 300, string(userData))
    
    return user, nil
}
```

### Rate Limiting

```go
func checkRateLimit(client *redis.Client, userID int, limit int) (bool, error) {
    key := fmt.Sprintf("rate_limit:user:%d", userID)
    
    // Check current count
    countStr, err := client.Get(key)
    if err != nil {
        // First request - initialize counter
        client.Setex(key, 60, "1") // 1 minute window
        return true, nil
    }
    
    count, _ := strconv.Atoi(countStr)
    if count >= limit {
        return false, nil // Rate limit exceeded
    }
    
    // Increment counter
    newCount := count + 1
    ttl, _ := client.Ttl(key)
    client.Setex(key, ttl, strconv.Itoa(newCount))
    
    return true, nil
}
```

### Batch Operations

```go
// Store multiple user preferences
err := client.MSet(
    "pref:user1:theme", "dark",
    "pref:user1:lang", "en",
    "pref:user1:timezone", "UTC",
)

// Retrieve multiple preferences
prefs, err := client.MGet(
    "pref:user1:theme",
    "pref:user1:lang",
    "pref:user1:timezone",
)

fmt.Printf("Theme: %s, Language: %s, Timezone: %s\n", prefs[0], prefs[1], prefs[2])
```

### Database Maintenance

```go
// Switch to analytics database
client.Select(2)

// Get statistics
count, _ := client.DBsize()
fmt.Printf("Analytics database has %d keys\n", count)

// Clear old analytics data
client.FlushDB()

// Switch back to main database
client.Select(0)
```

## API Reference

### Initialization

#### `Initialize(address, password string, maxConnection int, timeout time.Duration) error`

Initialize the Redis client with connection pool settings.

**Parameters:**
- `address` - Redis server address (host:port)
- `password` - Authentication password (empty string for no auth)
- `maxConnection` - Maximum number of idle connections in pool
- `timeout` - Idle connection timeout duration

**Returns:** Error if initialization fails

### Connection Methods

#### `Ping() error`
Test connection to Redis server.

#### `Finalize() error`
Close all connections and clean up resources.

#### `Select(index int) error`
Select database by index (0-15 for default Redis configuration).

### String Operations

#### `Get(key any) (string, error)`
Retrieve value for the given key.

#### `Set(key, value any) error`
Set key to hold the string value.

#### `Setex(key any, second int, value any) error`
Set key with expiration time in seconds.

#### `MGet(keys ...any) ([]string, error)`
Get values of multiple keys.

#### `MSet(args ...any) error`
Set multiple keys to multiple values. Arguments: key1, value1, key2, value2, ...

### Key Management

#### `Del(key any) error`
Delete the specified key.

#### `Exists(keys ...any) (bool, error)`
Check if one or more keys exist.

#### `Rename(currentKey, newKey any) error`
Rename a key.

#### `RandomKey() (string, error)`
Return a random key from the current database.

#### `Ttl(key any) (int, error)`
Get time to live for a key in seconds.

**Return values:**
- `-2` if key does not exist
- `-1` if key has no expiration
- Positive integer for remaining seconds

### Database Operations

#### `DBsize() (int, error)`
Return the number of keys in the current database.

#### `FlushDB() error`
Delete all keys in the current database.

#### `FlushAll() error`
Delete all keys in all databases.

### Server Information

#### `Info(info string) (string, error)`
Get information and statistics about the server.

**Parameters:**
- `info` - Category name (case-insensitive) or "ALL" for everything

**Categories:**
- `Server` - General server information
- `Clients` - Client connections section
- `Memory` - Memory consumption
- `Persistence` - RDB and AOF information
- `Stats` - General statistics
- `Replication` - Master/replica replication info
- `CPU` - CPU consumption statistics
- `Cluster` - Cluster section
- `Keyspace` - Database related statistics

## Best Practices

### 1. Configure Pool Size Appropriately

```go
// Web application with high concurrency
client.Initialize(addr, pass, 100, 300*time.Second)

// Background worker
client.Initialize(addr, pass, 10, 60*time.Second)
```

### 2. Always Close Connections

```go
defer client.Finalize()
```

### 3. Use SETEX for Temporary Data

```go
// Good: Automatic expiration
client.Setex("temp:token", 3600, token)

// Avoid: Manual cleanup needed
client.Set("temp:token", token)
// Remember to delete later!
```

### 4. Batch Operations When Possible

```go
// Good: Single round trip
values, err := client.MGet("key1", "key2", "key3")

// Avoid: Multiple round trips
val1, _ := client.Get("key1")
val2, _ := client.Get("key2")
val3, _ := client.Get("key3")
```

### 5. Handle TTL Return Values

```go
ttl, err := client.Ttl("key")
if err != nil {
    return err
}

switch ttl {
case -2:
    fmt.Println("Key does not exist")
case -1:
    fmt.Println("Key has no expiration")
default:
    fmt.Printf("Expires in %d seconds\n", ttl)
}
```

### 6. Use Key Naming Conventions

```go
// Good: Hierarchical naming
"user:123:profile"
"session:abc123:data"
"cache:product:456"

// Avoid: Flat naming
"user123profile"
"sessionabc123"
```

## Error Handling

Common errors and solutions:

```go
// Not initialized
err := client.Get("key")
// Error: "please call Initialize first"
// Solution: Call Initialize before any operation

// Connection failed
err := client.Initialize("localhost:6379", "", 10, 60*time.Second)
// Error: dial tcp connection refused
// Solution: Ensure Redis server is running

// Authentication failed
err := client.Initialize("localhost:6379", "wrong_password", 10, 60*time.Second)
// Error: NOAUTH Authentication required / ERR invalid password
// Solution: Use correct password or empty string for no auth
```

## Performance Tips

1. **Use Connection Pooling** - Reuse connections instead of creating new ones
2. **Batch Operations** - Use MGET/MSET to reduce network round trips
3. **Set Appropriate TTLs** - Prevent memory bloat from stale keys
4. **Monitor Memory** - Use INFO Memory to track usage
5. **Use SELECT Sparingly** - Minimize database switching overhead

## Limitations

1. **Limited Data Structures** - Only basic string operations (no Lists, Sets, Hashes, etc.)
2. **No Pipeline Support** - Cannot batch commands into pipeline
3. **No Pub/Sub** - No support for publish/subscribe messaging
4. **No Lua Scripts** - Cannot execute Lua scripts
5. **No Transactions** - No MULTI/EXEC transaction support

For advanced features, consider using the Redigo library directly.

## Dependencies

- `github.com/gomodule/redigo/redis` - Redis client library

## Related Packages

- [Redigo](https://github.com/gomodule/redigo) - Full-featured Redis client
- [go-redis](https://github.com/redis/go-redis) - Alternative Redis client

## Further Reading

- [Redis Commands](https://redis.io/commands)
- [Redis Data Types](https://redis.io/topics/data-types)
- [Redis Persistence](https://redis.io/topics/persistence)
- [Redigo Documentation](https://pkg.go.dev/github.com/gomodule/redigo/redis)
