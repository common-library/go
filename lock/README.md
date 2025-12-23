# Lock

Mutex implementations for thread-safe synchronization in Go.

## Overview

The lock package provides enhanced mutex functionality including basic mutex operations and key-based mutex management. It offers convenient wrappers around Go's sync.Mutex with additional features like TryLock and per-key locking mechanisms.

## Features

- **Basic Mutex** - Wrapper around sync.Mutex with TryLock support
- **Key-Based Mutexes** - Manage multiple independent locks by key
- **Non-Blocking Locks** - TryLock for conditional locking
- **Thread-Safe** - Built on sync.Map for concurrent access
- **On-Demand Creation** - Mutexes created automatically when needed
- **Memory Management** - UnlockAndDelete for cleanup

## Installation

```bash
go get -u github.com/common-library/go/lock
```

## Quick Start

### Basic Mutex

```go
import "github.com/common-library/go/lock"

var mu lock.Mutex

func criticalSection() {
    mu.Lock()
    defer mu.Unlock()
    
    // Only one goroutine executes this at a time
    sharedResource++
}
```

### Key-Based Mutex

```go
import "github.com/common-library/go/lock"

var mutexes lock.MutexByKey

func updateUser(userID string) {
    mutexes.Lock(userID)
    defer mutexes.Unlock(userID)
    
    // Different users can be updated concurrently
    // Same user operations are serialized
    database.Update(userID, newData)
}
```

## API Reference

### Mutex Type

```go
type Mutex struct {
    // Unexported fields
}
```

Wrapper around sync.Mutex providing lock operations.

### Mutex.Lock

```go
func (m *Mutex) Lock()
```

Acquires the mutex, blocking until available.

**Behavior:**
- Blocks until lock is acquired
- Must be paired with Unlock
- Calling on already-locked mutex causes deadlock

### Mutex.TryLock

```go
func (m *Mutex) TryLock() bool
```

Attempts to acquire the mutex without blocking.

**Returns:**
- `bool` - true if lock acquired, false otherwise

**Behavior:**
- Returns immediately (non-blocking)
- Returns false if mutex is locked
- Must call Unlock if returns true

### Mutex.Unlock

```go
func (m *Mutex) Unlock()
```

Releases the mutex.

**Behavior:**
- Releases the lock
- Panics if called on unlocked mutex
- Should be called by same goroutine that locked

### MutexByKey Type

```go
type MutexByKey struct {
    // Unexported fields
}
```

Manages multiple mutexes indexed by keys.

### MutexByKey.Lock

```go
func (mbk *MutexByKey) Lock(key any)
```

Acquires the mutex for the given key.

**Parameters:**
- `key` - Any comparable value to identify the mutex

**Behavior:**
- Creates mutex on-demand if doesn't exist
- Blocks until lock acquired
- Different keys can be locked concurrently

### MutexByKey.TryLock

```go
func (mbk *MutexByKey) TryLock(key any) bool
```

Attempts to acquire the mutex for the key without blocking.

**Parameters:**
- `key` - Any comparable value to identify the mutex

**Returns:**
- `bool` - true if lock acquired, false otherwise

### MutexByKey.Unlock

```go
func (mbk *MutexByKey) Unlock(key any)
```

Releases the mutex for the given key.

**Parameters:**
- `key` - The key identifying the mutex

**Behavior:**
- Unlocks the mutex
- Mutex remains in map after unlock
- Use UnlockAndDelete to remove mutex

### MutexByKey.UnlockAndDelete

```go
func (mbk *MutexByKey) UnlockAndDelete(key any)
```

Releases and removes the mutex for the key.

**Parameters:**
- `key` - The key identifying the mutex

**Behavior:**
- Unlocks the mutex
- Removes mutex from internal map
- Frees memory

### MutexByKey.Delete

```go
func (mbk *MutexByKey) Delete(key any)
```

Removes the mutex without unlocking.

**Parameters:**
- `key` - The key identifying the mutex

**Warning:** Only use when certain mutex is not locked.

## Complete Examples

### Protecting Shared Resources

```go
package main

import (
    "fmt"
    "sync"
    
    "github.com/common-library/go/lock"
)

var (
    counter int
    mu      lock.Mutex
)

func increment() {
    mu.Lock()
    defer mu.Unlock()
    
    counter++
}

func main() {
    var wg sync.WaitGroup
    
    // Start 1000 goroutines
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            increment()
        }()
    }
    
    wg.Wait()
    fmt.Printf("Final counter: %d\n", counter) // Output: 1000
}
```

### Non-Blocking Lock with TryLock

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/common-library/go/lock"
)

var mu lock.Mutex

func processWithFallback() {
    if mu.TryLock() {
        defer mu.Unlock()
        
        fmt.Println("Processing primary task...")
        time.Sleep(1 * time.Second)
    } else {
        fmt.Println("Primary task busy, doing alternative task...")
        doAlternativeTask()
    }
}

func doAlternativeTask() {
    fmt.Println("Alternative task completed")
}

func main() {
    // First call acquires lock
    go processWithFallback()
    
    // Second call immediately falls back
    time.Sleep(10 * time.Millisecond)
    processWithFallback()
    
    time.Sleep(2 * time.Second)
}
```

### Retry Pattern with TryLock

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    "github.com/common-library/go/lock"
)

var mu lock.Mutex

func processWithRetry(maxRetries int) error {
    for i := 0; i < maxRetries; i++ {
        if mu.TryLock() {
            defer mu.Unlock()
            
            fmt.Println("Processing...")
            return nil
        }
        
        fmt.Printf("Retry %d/%d...\n", i+1, maxRetries)
        time.Sleep(100 * time.Millisecond)
    }
    
    return fmt.Errorf("failed to acquire lock after %d retries", maxRetries)
}

func main() {
    if err := processWithRetry(3); err != nil {
        log.Fatal(err)
    }
}
```

### Per-User Locking

```go
package main

import (
    "fmt"
    "sync"
    
    "github.com/common-library/go/lock"
)

var mutexes lock.MutexByKey

type User struct {
    ID      string
    Balance int
}

var users = map[string]*User{
    "user1": {ID: "user1", Balance: 100},
    "user2": {ID: "user2", Balance: 200},
}

func updateBalance(userID string, amount int) {
    mutexes.Lock(userID)
    defer mutexes.Unlock(userID)
    
    user := users[userID]
    user.Balance += amount
    
    fmt.Printf("User %s balance: %d\n", userID, user.Balance)
}

func main() {
    var wg sync.WaitGroup
    
    // Update different users concurrently (no blocking)
    wg.Add(2)
    go func() {
        defer wg.Done()
        updateBalance("user1", 50)
    }()
    go func() {
        defer wg.Done()
        updateBalance("user2", 75)
    }()
    
    // Update same user (serialized)
    wg.Add(2)
    go func() {
        defer wg.Done()
        updateBalance("user1", 10)
    }()
    go func() {
        defer wg.Done()
        updateBalance("user1", 20)
    }()
    
    wg.Wait()
}
```

### Resource Pool Management

```go
package main

import (
    "fmt"
    "sync"
    
    "github.com/common-library/go/lock"
)

var mutexes lock.MutexByKey

type ResourcePool struct {
    resources map[string]*Resource
}

type Resource struct {
    ID   string
    Data string
}

func (p *ResourcePool) UpdateResource(id string, data string) {
    mutexes.Lock(id)
    defer mutexes.Unlock(id)
    
    if resource, exists := p.resources[id]; exists {
        resource.Data = data
        fmt.Printf("Updated resource %s\n", id)
    }
}

func (p *ResourcePool) TryUpdateResource(id string, data string) bool {
    if mutexes.TryLock(id) {
        defer mutexes.Unlock(id)
        
        if resource, exists := p.resources[id]; exists {
            resource.Data = data
            fmt.Printf("Updated resource %s\n", id)
            return true
        }
    }
    
    fmt.Printf("Resource %s is busy\n", id)
    return false
}

func main() {
    pool := &ResourcePool{
        resources: map[string]*Resource{
            "res1": {ID: "res1", Data: "initial"},
            "res2": {ID: "res2", Data: "initial"},
        },
    }
    
    var wg sync.WaitGroup
    
    wg.Add(3)
    go func() {
        defer wg.Done()
        pool.UpdateResource("res1", "data1")
    }()
    go func() {
        defer wg.Done()
        pool.UpdateResource("res2", "data2")
    }()
    go func() {
        defer wg.Done()
        pool.TryUpdateResource("res1", "data3")
    }()
    
    wg.Wait()
}
```

### Temporary Locks with Cleanup

```go
package main

import (
    "fmt"
    "sync"
    
    "github.com/common-library/go/lock"
)

var mutexes lock.MutexByKey

func processRequest(requestID string) {
    // Lock for this specific request
    mutexes.Lock(requestID)
    defer mutexes.UnlockAndDelete(requestID) // Clean up when done
    
    fmt.Printf("Processing request %s\n", requestID)
    // Process...
}

func main() {
    var wg sync.WaitGroup
    
    // Each request gets its own lock, cleaned up after completion
    for i := 0; i < 10; i++ {
        wg.Add(1)
        requestID := fmt.Sprintf("req-%d", i)
        
        go func(id string) {
            defer wg.Done()
            processRequest(id)
        }(requestID)
    }
    
    wg.Wait()
}
```

### Cache Access Control

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/common-library/go/lock"
)

var mutexes lock.MutexByKey

type Cache struct {
    data map[string]string
}

func NewCache() *Cache {
    return &Cache{
        data: make(map[string]string),
    }
}

func (c *Cache) Get(key string) (string, bool) {
    mutexes.Lock(key)
    defer mutexes.Unlock(key)
    
    value, exists := c.data[key]
    return value, exists
}

func (c *Cache) Set(key string, value string) {
    mutexes.Lock(key)
    defer mutexes.Unlock(key)
    
    c.data[key] = value
}

func (c *Cache) Delete(key string) {
    mutexes.Lock(key)
    defer mutexes.UnlockAndDelete(key)
    
    delete(c.data, key)
}

func main() {
    cache := NewCache()
    
    // Concurrent access to different keys
    go cache.Set("user:1", "Alice")
    go cache.Set("user:2", "Bob")
    
    time.Sleep(100 * time.Millisecond)
    
    if value, exists := cache.Get("user:1"); exists {
        fmt.Printf("Found: %s\n", value)
    }
    
    cache.Delete("user:1")
}
```

### Timeout Pattern

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/common-library/go/lock"
)

var mu lock.Mutex

func processWithTimeout(timeout time.Duration) error {
    deadline := time.After(timeout)
    ticker := time.NewTicker(10 * time.Millisecond)
    defer ticker.Stop()
    
    for {
        select {
        case <-deadline:
            return fmt.Errorf("timeout acquiring lock")
            
        case <-ticker.C:
            if mu.TryLock() {
                defer mu.Unlock()
                
                fmt.Println("Lock acquired, processing...")
                time.Sleep(100 * time.Millisecond)
                return nil
            }
        }
    }
}

func main() {
    // Simulate long-running lock holder
    go func() {
        mu.Lock()
        time.Sleep(500 * time.Millisecond)
        mu.Unlock()
    }()
    
    time.Sleep(50 * time.Millisecond)
    
    if err := processWithTimeout(1 * time.Second); err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Success!")
    }
}
```

## Best Practices

### 1. Always Use defer with Unlock

```go
// Good: Ensures unlock even if panic occurs
mu.Lock()
defer mu.Unlock()
// Critical section

// Avoid: Manual unlock can be missed
mu.Lock()
// Critical section
mu.Unlock() // Might not be called if panic occurs
```

### 2. Keep Critical Sections Small

```go
// Good: Minimal lock time
mu.Lock()
data := sharedResource
mu.Unlock()

// Process data without holding lock
result := processData(data)

// Avoid: Holding lock during expensive operations
mu.Lock()
data := sharedResource
result := processData(data) // Expensive operation
mu.Unlock()
```

### 3. Use TryLock for Optional Processing

```go
// Good: Fall back if lock unavailable
if mu.TryLock() {
    defer mu.Unlock()
    doExpensiveWork()
} else {
    doQuickAlternative()
}

// Avoid: Blocking on optional work
mu.Lock() // Might block unnecessarily
defer mu.Unlock()
doExpensiveWork()
```

### 4. Use UnlockAndDelete for Temporary Locks

```go
// Good: Clean up one-time locks
mutexes.Lock(requestID)
defer mutexes.UnlockAndDelete(requestID)

// Avoid: Accumulating unused mutexes
mutexes.Lock(requestID)
defer mutexes.Unlock(requestID) // Mutex stays in memory
```

### 5. Use Consistent Lock Ordering

```go
// Good: Always lock in same order
mutexes.Lock("resource-A")
defer mutexes.Unlock("resource-A")
mutexes.Lock("resource-B")
defer mutexes.Unlock("resource-B")

// Avoid: Inconsistent ordering (can cause deadlock)
// Thread 1: Lock A then B
// Thread 2: Lock B then A
```

## Common Use Cases

### 1. Protecting Counters

```go
var (
    counter int
    mu      lock.Mutex
)

func increment() {
    mu.Lock()
    defer mu.Unlock()
    counter++
}
```

### 2. Cache Management

```go
var (
    cache   map[string]string
    mutexes lock.MutexByKey
)

func getCached(key string) string {
    mutexes.Lock(key)
    defer mutexes.Unlock(key)
    return cache[key]
}
```

### 3. Resource Pooling

```go
var mutexes lock.MutexByKey

func borrowResource(id string) {
    mutexes.Lock(id)
    // Use resource
}

func returnResource(id string) {
    mutexes.Unlock(id)
}
```

### 4. Rate Limiting

```go
var mu lock.Mutex

func rateLimitedAction() bool {
    if mu.TryLock() {
        defer mu.Unlock()
        doAction()
        return true
    }
    return false // Too busy
}
```

## Performance Tips

1. **Minimize Lock Duration** - Hold locks for shortest time possible
2. **Use TryLock** - Avoid blocking when work is optional
3. **Clean Up** - Use UnlockAndDelete for temporary locks
4. **Fine-Grained Locking** - Use MutexByKey instead of global mutex
5. **Lock Ordering** - Prevent deadlocks with consistent ordering

## Testing

### Testing Concurrent Access

```go
func TestConcurrentIncrement(t *testing.T) {
    var (
        counter int
        mu      lock.Mutex
        wg      sync.WaitGroup
    )
    
    iterations := 1000
    goroutines := 10
    
    for i := 0; i < goroutines; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for j := 0; j < iterations; j++ {
                mu.Lock()
                counter++
                mu.Unlock()
            }
        }()
    }
    
    wg.Wait()
    
    expected := goroutines * iterations
    if counter != expected {
        t.Errorf("Expected %d, got %d", expected, counter)
    }
}
```

### Testing TryLock

```go
func TestTryLock(t *testing.T) {
    var mu lock.Mutex
    
    // First lock succeeds
    if !mu.TryLock() {
        t.Fatal("Expected first TryLock to succeed")
    }
    
    // Second lock fails (already locked)
    if mu.TryLock() {
        t.Fatal("Expected second TryLock to fail")
    }
    
    mu.Unlock()
    
    // Third lock succeeds (after unlock)
    if !mu.TryLock() {
        t.Fatal("Expected third TryLock to succeed")
    }
    mu.Unlock()
}
```

## Dependencies

- `sync` - Go standard library

## Further Reading

- [Go sync package](https://pkg.go.dev/sync)
- [Mutexes in Go](https://go.dev/tour/concurrency/9)
- [Go Memory Model](https://go.dev/ref/mem)
