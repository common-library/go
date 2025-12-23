# Utility

General-purpose utility functions for Go applications.

## Overview

The utility package provides helper functions for common tasks including runtime introspection, type information extraction, and network utilities. These utilities simplify everyday programming tasks that would otherwise require boilerplate code.

## Features

- **Caller Information** - Retrieve file, line, function, and goroutine details
- **Type Introspection** - Get type names from any value
- **CIDR Utilities** - Check IP containment and enumerate IP ranges
- **Runtime Reflection** - Access stack and goroutine information

## Installation

```bash
go get -u github.com/common-library/go/utility
```

## Quick Start

```go
import "github.com/common-library/go/utility"

// Get caller information
callerInfo, _ := utility.GetCallerInfo(1)
fmt.Printf("Called from %s:%d\n", callerInfo.FileName, callerInfo.Line)

// Get type name
typeName := utility.GetTypeName(42)
fmt.Println(typeName) // Output: "int"

// Check CIDR containment
contains, _ := utility.WhetherCidrContainsIp("192.168.1.0/24", "192.168.1.100")
fmt.Println(contains) // Output: true
```

## API Reference

### CallerInfo Type

```go
type CallerInfo struct {
    PackageName  string
    FileName     string
    FunctionName string
    Line         int
    GoroutineID  int
}
```

Contains details about a function call location.

### GetCallerInfo

```go
func GetCallerInfo(numberOfStackFramesToAscend int) (CallerInfo, error)
```

Retrieves caller information from the stack.

**Parameters:**
- `numberOfStackFramesToAscend` - Stack frames to skip (0=self, 1=caller, 2=caller's caller)

**Returns:**
- `CallerInfo` - Caller details
- `error` - Error if stack retrieval fails

### GetTypeName

```go
func GetTypeName(value any) string
```

Returns the type name of any value.

**Parameters:**
- `value` - Any value to inspect

**Returns:**
- `string` - Type name (e.g., "int", "*User", "[]string")

### WhetherCidrContainsIp

```go
func WhetherCidrContainsIp(cidr, ip string) (bool, error)
```

Checks if an IP is within a CIDR range.

**Parameters:**
- `cidr` - CIDR notation (e.g., "192.168.1.0/24")
- `ip` - IP address to check

**Returns:**
- `bool` - true if IP is in range
- `error` - Error if CIDR parsing fails

### GetAllIpsOfCidr

```go
func GetAllIpsOfCidr(cidr string) ([]string, error)
```

Returns all usable IPs in a CIDR range (excludes network and broadcast).

**Parameters:**
- `cidr` - CIDR notation

**Returns:**
- `[]string` - List of IP addresses
- `error` - Error if CIDR parsing fails

## Complete Examples

### Caller Information for Logging

```go
package main

import (
    "fmt"
    "log"
    "github.com/common-library/go/utility"
)

func logWithCaller(message string) {
    info, err := utility.GetCallerInfo(1)
    if err != nil {
        log.Println(message)
        return
    }
    
    log.Printf("[%s:%d %s] %s",
        info.FileName,
        info.Line,
        info.FunctionName,
        message,
    )
}

func processData() {
    logWithCaller("Starting data processing")
    // ... process data ...
    logWithCaller("Data processing complete")
}

func main() {
    processData()
    // Output: [main.go:21 main.processData] Starting data processing
}
```

### Stack Trace Debugging

```go
package main

import (
    "fmt"
    "github.com/common-library/go/utility"
)

func printStack(depth int) {
    fmt.Println("Stack trace:")
    for i := 0; i < depth; i++ {
        info, err := utility.GetCallerInfo(i)
        if err != nil {
            break
        }
        
        fmt.Printf("  #%d %s:%d %s (goroutine %d)\n",
            i,
            info.FileName,
            info.Line,
            info.FunctionName,
            info.GoroutineID,
        )
    }
}

func level3() {
    printStack(5)
}

func level2() {
    level3()
}

func level1() {
    level2()
}

func main() {
    level1()
}
```

### Type-Based Processing

```go
package main

import (
    "fmt"
    "github.com/common-library/go/utility"
)

func processValue(value any) {
    typeName := utility.GetTypeName(value)
    
    switch typeName {
    case "int":
        fmt.Printf("Processing integer: %d\n", value)
    case "string":
        fmt.Printf("Processing string: %s\n", value)
    case "[]int":
        fmt.Printf("Processing int slice: %v\n", value)
    default:
        fmt.Printf("Processing %s: %v\n", typeName, value)
    }
}

func main() {
    processValue(42)
    processValue("hello")
    processValue([]int{1, 2, 3})
    processValue(map[string]int{"a": 1})
}
```

### CIDR IP Validation

```go
package main

import (
    "fmt"
    "log"
    "github.com/common-library/go/utility"
)

func validateIPInNetwork(cidr, ip string) {
    contains, err := utility.WhetherCidrContainsIp(cidr, ip)
    if err != nil {
        log.Fatalf("Invalid CIDR: %v", err)
    }
    
    if contains {
        fmt.Printf("✓ %s is in %s\n", ip, cidr)
    } else {
        fmt.Printf("✗ %s is NOT in %s\n", ip, cidr)
    }
}

func main() {
    cidr := "192.168.1.0/24"
    
    validateIPInNetwork(cidr, "192.168.1.100")  // ✓
    validateIPInNetwork(cidr, "192.168.1.1")    // ✓
    validateIPInNetwork(cidr, "192.168.2.100")  // ✗
    validateIPInNetwork(cidr, "10.0.0.1")       // ✗
}
```

### IP Range Enumeration

```go
package main

import (
    "fmt"
    "log"
    "github.com/common-library/go/utility"
)

func main() {
    // Small subnet
    cidr := "192.168.1.0/29"
    
    ips, err := utility.GetAllIpsOfCidr(cidr)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("CIDR: %s\n", cidr)
    fmt.Printf("Usable IPs: %d\n", len(ips))
    fmt.Println("IP addresses:")
    for _, ip := range ips {
        fmt.Printf("  - %s\n", ip)
    }
    
    // Output:
    // CIDR: 192.168.1.0/29
    // Usable IPs: 6
    // IP addresses:
    //   - 192.168.1.1
    //   - 192.168.1.2
    //   - 192.168.1.3
    //   - 192.168.1.4
    //   - 192.168.1.5
    //   - 192.168.1.6
}
```

### IP Allocation Manager

```go
package main

import (
    "fmt"
    "log"
    "github.com/common-library/go/utility"
)

type IPPool struct {
    cidr      string
    allocated map[string]bool
    available []string
}

func NewIPPool(cidr string) (*IPPool, error) {
    ips, err := utility.GetAllIpsOfCidr(cidr)
    if err != nil {
        return nil, err
    }
    
    return &IPPool{
        cidr:      cidr,
        allocated: make(map[string]bool),
        available: ips,
    }, nil
}

func (p *IPPool) Allocate() (string, error) {
    if len(p.available) == 0 {
        return "", fmt.Errorf("no available IPs")
    }
    
    ip := p.available[0]
    p.available = p.available[1:]
    p.allocated[ip] = true
    
    return ip, nil
}

func (p *IPPool) Release(ip string) {
    if p.allocated[ip] {
        delete(p.allocated, ip)
        p.available = append(p.available, ip)
    }
}

func (p *IPPool) Stats() {
    fmt.Printf("CIDR: %s\n", p.cidr)
    fmt.Printf("Allocated: %d\n", len(p.allocated))
    fmt.Printf("Available: %d\n", len(p.available))
}

func main() {
    pool, err := NewIPPool("10.0.0.0/28")
    if err != nil {
        log.Fatal(err)
    }
    
    // Allocate IPs
    ip1, _ := pool.Allocate()
    ip2, _ := pool.Allocate()
    ip3, _ := pool.Allocate()
    
    fmt.Printf("Allocated: %s, %s, %s\n", ip1, ip2, ip3)
    pool.Stats()
    
    // Release one
    pool.Release(ip2)
    pool.Stats()
}
```

### Generic Type Checker

```go
package main

import (
    "fmt"
    "github.com/common-library/go/utility"
)

func isPointerType(value any) bool {
    typeName := utility.GetTypeName(value)
    return len(typeName) > 0 && typeName[0] == '*'
}

func isSliceType(value any) bool {
    typeName := utility.GetTypeName(value)
    return len(typeName) > 1 && typeName[0:2] == "[]"
}

func isMapType(value any) bool {
    typeName := utility.GetTypeName(value)
    return len(typeName) > 2 && typeName[0:3] == "map"
}

func main() {
    var x int = 42
    var ptr *int = &x
    var slice []int = []int{1, 2, 3}
    var m map[string]int = map[string]int{"a": 1}
    
    fmt.Printf("x is pointer: %v\n", isPointerType(x))     // false
    fmt.Printf("ptr is pointer: %v\n", isPointerType(ptr)) // true
    fmt.Printf("slice is slice: %v\n", isSliceType(slice)) // true
    fmt.Printf("m is map: %v\n", isMapType(m))             // true
}
```

### Goroutine Tracking

```go
package main

import (
    "fmt"
    "sync"
    "time"
    "github.com/common-library/go/utility"
)

func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done()
    
    info, _ := utility.GetCallerInfo(0)
    fmt.Printf("Worker %d running in goroutine %d\n", id, info.GoroutineID)
    
    time.Sleep(100 * time.Millisecond)
}

func main() {
    var wg sync.WaitGroup
    
    mainInfo, _ := utility.GetCallerInfo(0)
    fmt.Printf("Main goroutine ID: %d\n", mainInfo.GoroutineID)
    
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go worker(i, &wg)
    }
    
    wg.Wait()
}
```

## Best Practices

### 1. Handle Errors from GetCallerInfo

```go
// Good: Check error
info, err := utility.GetCallerInfo(1)
if err != nil {
    log.Printf("Failed to get caller info: %v", err)
    return
}

// Avoid: Ignore error
info, _ := utility.GetCallerInfo(1)
```

### 2. Use Appropriate Stack Depth

```go
// Good: Skip correct number of frames
func logHelper(msg string) {
    info, _ := utility.GetCallerInfo(1) // Skip logHelper itself
    log.Printf("[%s:%d] %s", info.FileName, info.Line, msg)
}

// Avoid: Wrong depth
func logHelper(msg string) {
    info, _ := utility.GetCallerInfo(0) // Points to logHelper, not caller
}
```

### 3. Be Cautious with Large CIDR Ranges

```go
// Good: Check size before enumerating
cidr := "10.0.0.0/16" // 65,534 IPs
// Consider: Do you really need all IPs in memory?

// Risky: Very large range
cidr := "10.0.0.0/8" // 16,777,214 IPs - may exhaust memory
ips, _ := utility.GetAllIpsOfCidr(cidr) // Dangerous!
```

### 4. Cache Type Names for Performance

```go
// Good: Cache type names if checking repeatedly
typeCache := make(map[any]string)

func getCachedTypeName(value any) string {
    if name, ok := typeCache[value]; ok {
        return name
    }
    name := utility.GetTypeName(value)
    typeCache[value] = name
    return name
}
```

## Performance Tips

1. **GetCallerInfo** - Has runtime overhead, use sparingly in hot paths
2. **GetTypeName** - Uses reflection, cache results when possible
3. **GetAllIpsOfCidr** - Memory intensive for large ranges, consider pagination
4. **CIDR Checks** - WhetherCidrContainsIp is efficient, prefer over full enumeration

## Testing

```go
func TestGetCallerInfo(t *testing.T) {
    info, err := utility.GetCallerInfo(0)
    if err != nil {
        t.Fatalf("GetCallerInfo failed: %v", err)
    }
    
    if info.FunctionName == "" {
        t.Error("Function name should not be empty")
    }
    
    if info.Line <= 0 {
        t.Error("Line number should be positive")
    }
}

func TestGetTypeName(t *testing.T) {
    tests := []struct {
        value    any
        expected string
    }{
        {42, "int"},
        {"hello", "string"},
        {[]int{}, "[]int"},
        {map[string]int{}, "map[string]int"},
    }
    
    for _, tt := range tests {
        result := utility.GetTypeName(tt.value)
        if result != tt.expected {
            t.Errorf("Expected %s, got %s", tt.expected, result)
        }
    }
}

func TestCIDR(t *testing.T) {
    contains, err := utility.WhetherCidrContainsIp("192.168.1.0/24", "192.168.1.100")
    if err != nil {
        t.Fatal(err)
    }
    
    if !contains {
        t.Error("IP should be in range")
    }
}
```

## Dependencies

- `runtime` - Go standard library
- `reflect` - Go standard library
- `net` - Go standard library

## Further Reading

- [Go runtime package](https://pkg.go.dev/runtime)
- [Go reflect package](https://pkg.go.dev/reflect)
- [CIDR Notation](https://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing)
