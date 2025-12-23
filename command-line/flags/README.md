# Flags

Type-safe command-line flag parsing with generic value retrieval for Go.

## Overview

The flags package provides a type-safe wrapper around Go's standard `flag` package, offering declarative flag definition and generic-based value retrieval. It supports 8 common data types and eliminates the need for type-specific getter functions.

For simple positional arguments, see the [arguments package](../arguments/).

## Features

- **Type-safe flag parsing** - Support for bool, int, string, duration, and more
- **Generic value retrieval** - Use `Get[T]()` instead of type-specific getters
- **Declarative definition** - Define all flags in a single `FlagInfo` slice
- **8 supported types** - bool, int, int64, uint, uint64, float64, string, time.Duration
- **Error handling** - Returns errors instead of panicking
- **Built on standard library** - Uses Go's proven `flag` package

## Installation

```bash
go get -u github.com/common-library/go/command-line/flags
```

## Quick Start

```go
import (
    "fmt"
    "log"
    "github.com/common-library/go/command-line/flags"
)

func main() {
    // Define flags
    err := flags.Parse([]flags.FlagInfo{
        {FlagName: "port", Usage: "server port", DefaultValue: 8080},
        {FlagName: "host", Usage: "server host", DefaultValue: "localhost"},
        {FlagName: "verbose", Usage: "enable verbose logging", DefaultValue: false},
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Retrieve values with type safety
    port := flags.Get[int]("port")
    host := flags.Get[string]("host")
    verbose := flags.Get[bool]("verbose")
    
    fmt.Printf("Server: %s:%d (verbose: %v)\n", host, port, verbose)
}
```

Run with:
```bash
./program -port=9000 -host=0.0.0.0 -verbose=true
```

## Supported Types

The flags package supports 8 common data types:

| Type | Example Default | Command-Line Example | Notes |
|------|----------------|---------------------|-------|
| `bool` | `true` | `-verbose=true` or `-verbose` | Boolean flag |
| `int` | `0` | `-port=8080` | Integer |
| `int64` | `int64(0)` | `-count=1000000` | 64-bit integer |
| `uint` | `uint(0)` | `-size=100` | Unsigned integer |
| `uint64` | `uint64(0)` | `-bytes=1024` | 64-bit unsigned |
| `float64` | `0.0` | `-rate=0.5` | Floating-point |
| `string` | `""` | `-name=server` | String |
| `time.Duration` | `0 * time.Second` | `-timeout=30s` | Duration (1s, 5m, 1h) |

## Usage Examples

### Basic Flag Parsing

```go
package main

import (
    "fmt"
    "log"
    "github.com/common-library/go/command-line/flags"
)

func main() {
    err := flags.Parse([]flags.FlagInfo{
        {FlagName: "name", Usage: "user name", DefaultValue: "guest"},
        {FlagName: "age", Usage: "user age", DefaultValue: 0},
        {FlagName: "admin", Usage: "admin privileges", DefaultValue: false},
    })
    if err != nil {
        log.Fatal(err)
    }
    
    name := flags.Get[string]("name")
    age := flags.Get[int]("age")
    admin := flags.Get[bool]("admin")
    
    fmt.Printf("User: %s, Age: %d, Admin: %v\n", name, age, admin)
}
```

```bash
./program -name=Alice -age=30 -admin=true
# Output: User: Alice, Age: 30, Admin: true
```

### Server Configuration

```go
package main

import (
    "fmt"
    "log"
    "time"
    "github.com/common-library/go/command-line/flags"
)

func main() {
    err := flags.Parse([]flags.FlagInfo{
        {
            FlagName:     "port",
            Usage:        "HTTP server port",
            DefaultValue: 8080,
        },
        {
            FlagName:     "host",
            Usage:        "HTTP server host",
            DefaultValue: "0.0.0.0",
        },
        {
            FlagName:     "read-timeout",
            Usage:        "HTTP read timeout",
            DefaultValue: 10 * time.Second,
        },
        {
            FlagName:     "write-timeout",
            Usage:        "HTTP write timeout",
            DefaultValue: 10 * time.Second,
        },
        {
            FlagName:     "max-connections",
            Usage:        "Maximum concurrent connections",
            DefaultValue: uint64(1000),
        },
        {
            FlagName:     "debug",
            Usage:        "Enable debug mode",
            DefaultValue: false,
        },
    })
    if err != nil {
        log.Fatalf("Failed to parse flags: %v", err)
    }
    
    config := ServerConfig{
        Port:           flags.Get[int]("port"),
        Host:           flags.Get[string]("host"),
        ReadTimeout:    flags.Get[time.Duration]("read-timeout"),
        WriteTimeout:   flags.Get[time.Duration]("write-timeout"),
        MaxConnections: flags.Get[uint64]("max-connections"),
        Debug:          flags.Get[bool]("debug"),
    }
    
    fmt.Printf("Server Configuration:\n")
    fmt.Printf("  Address: %s:%d\n", config.Host, config.Port)
    fmt.Printf("  Read Timeout: %s\n", config.ReadTimeout)
    fmt.Printf("  Write Timeout: %s\n", config.WriteTimeout)
    fmt.Printf("  Max Connections: %d\n", config.MaxConnections)
    fmt.Printf("  Debug: %v\n", config.Debug)
}

type ServerConfig struct {
    Port           int
    Host           string
    ReadTimeout    time.Duration
    WriteTimeout   time.Duration
    MaxConnections uint64
    Debug          bool
}
```

```bash
# Use defaults
./server

# Custom configuration
./server -port=9000 -host=localhost -debug=true

# With timeouts
./server -read-timeout=30s -write-timeout=30s -max-connections=5000

# Help
./server -h
```

### All Supported Types

```go
package main

import (
    "fmt"
    "log"
    "time"
    "github.com/common-library/go/command-line/flags"
)

func main() {
    err := flags.Parse([]flags.FlagInfo{
        // Boolean flags
        {FlagName: "verbose", Usage: "verbose output", DefaultValue: false},
        {FlagName: "dry-run", Usage: "perform dry run", DefaultValue: true},
        
        // Integer flags
        {FlagName: "count", Usage: "iteration count", DefaultValue: 100},
        {FlagName: "workers", Usage: "worker count", DefaultValue: int64(8)},
        
        // Unsigned integer flags
        {FlagName: "buffer-size", Usage: "buffer size in bytes", DefaultValue: uint(4096)},
        {FlagName: "max-size", Usage: "maximum file size", DefaultValue: uint64(1073741824)},
        
        // Floating-point flags
        {FlagName: "rate", Usage: "sampling rate", DefaultValue: 0.01},
        
        // String flags
        {FlagName: "output", Usage: "output file path", DefaultValue: "output.txt"},
        {FlagName: "format", Usage: "output format", DefaultValue: "json"},
        
        // Duration flags
        {FlagName: "interval", Usage: "polling interval", DefaultValue: 5 * time.Second},
        {FlagName: "deadline", Usage: "operation deadline", DefaultValue: 1 * time.Minute},
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Retrieve with type safety
    verbose := flags.Get[bool]("verbose")
    count := flags.Get[int]("count")
    rate := flags.Get[float64]("rate")
    output := flags.Get[string]("output")
    interval := flags.Get[time.Duration]("interval")
    
    fmt.Printf("Config: verbose=%v, count=%d, rate=%.2f, output=%s, interval=%s\n",
        verbose, count, rate, output, interval)
}
```

### Duration Flags

Duration flags support human-readable formats:

```go
err := flags.Parse([]flags.FlagInfo{
    {FlagName: "timeout", Usage: "request timeout", DefaultValue: 30 * time.Second},
    {FlagName: "interval", Usage: "polling interval", DefaultValue: 5 * time.Minute},
    {FlagName: "deadline", Usage: "operation deadline", DefaultValue: 1 * time.Hour},
})

timeout := flags.Get[time.Duration]("timeout")
```

```bash
# Seconds
./program -timeout=45s

# Minutes
./program -timeout=5m

# Hours
./program -timeout=2h

# Combined
./program -timeout=1h30m45s

# Milliseconds, microseconds, nanoseconds
./program -timeout=500ms
./program -timeout=100us
./program -timeout=1000ns
```

### Error Handling

```go
package main

import (
    "fmt"
    "log"
    "github.com/common-library/go/command-line/flags"
)

func main() {
    // Try to use unsupported type
    err := flags.Parse([]flags.FlagInfo{
        {FlagName: "config", Usage: "config map", DefaultValue: map[string]string{}},
    })
    
    if err != nil {
        // Error: "this data type is not supported. - (map[string]string)"
        log.Fatalf("Flag parsing error: %v", err)
    }
}
```

### Type Safety with Generics

```go
// Correct: Type matches default value
{FlagName: "port", DefaultValue: 8080}
port := flags.Get[int]("port")  // ✅ OK

// Incorrect: Type mismatch
{FlagName: "port", DefaultValue: 8080}
port := flags.Get[int64]("port")  // ❌ Runtime panic!

// Explicit type to avoid ambiguity
{FlagName: "port", DefaultValue: int(8080)}     // Explicit int
{FlagName: "count", DefaultValue: int64(1000)}  // Explicit int64
```

## API Reference

### `FlagInfo` Struct

Defines a single command-line flag.

```go
type FlagInfo struct {
    FlagName     string  // Name of the flag (used in -flagName=value)
    Usage        string  // Help text displayed with -h
    DefaultValue any     // Default value (type determines flag type)
}
```

**Fields:**
- `FlagName` - The name used on command line (e.g., "port" for `-port=8080`)
- `Usage` - Help text shown when user runs `program -h`
- `DefaultValue` - Default value and type indicator (must be one of 8 supported types)

**Example:**
```go
flagInfo := flags.FlagInfo{
    FlagName:     "timeout",
    Usage:        "request timeout duration (e.g., 30s, 5m)",
    DefaultValue: 30 * time.Second,
}
```

### `Parse(flagInfos []FlagInfo) error`

Parses command-line flags based on the provided configuration.

**Parameters:**
- `flagInfos` - Slice of `FlagInfo` structs defining each flag

**Returns:**
- `error` - Returns an error if an unsupported data type is encountered, nil otherwise

**Supported Types:**
- `bool`, `int`, `int64`, `uint`, `uint64`, `float64`, `string`, `time.Duration`

**Example:**
```go
err := flags.Parse([]flags.FlagInfo{
    {FlagName: "port", Usage: "server port", DefaultValue: 8080},
    {FlagName: "host", Usage: "server host", DefaultValue: "localhost"},
})
if err != nil {
    log.Fatal(err)
}
```

### `Get[T any](flagName string) T`

Retrieves a parsed flag value with type safety using Go generics.

**Type Parameters:**
- `T` - The expected type of the flag value (must match the type used in `Parse`)

**Parameters:**
- `flagName` - The name of the flag to retrieve

**Returns:**
- The flag value cast to type `T`

**Panics:**
- If `flagName` doesn't exist
- If type parameter `T` doesn't match the actual flag type

**Example:**
```go
port := flags.Get[int]("port")
host := flags.Get[string]("host")
verbose := flags.Get[bool]("verbose")
timeout := flags.Get[time.Duration]("timeout")
```

## Advanced Usage

### Dynamic Flag Configuration

```go
package main

import (
    "fmt"
    "log"
    "os"
    "github.com/common-library/go/command-line/flags"
)

func main() {
    // Base flags
    flagInfos := []flags.FlagInfo{
        {FlagName: "config", Usage: "config file path", DefaultValue: "config.json"},
    }
    
    // Add debug flags in development
    if os.Getenv("ENV") == "development" {
        flagInfos = append(flagInfos, []flags.FlagInfo{
            {FlagName: "debug", Usage: "enable debug mode", DefaultValue: true},
            {FlagName: "profile", Usage: "enable profiling", DefaultValue: false},
        }...)
    }
    
    err := flags.Parse(flagInfos)
    if err != nil {
        log.Fatal(err)
    }
    
    config := flags.Get[string]("config")
    fmt.Printf("Using config: %s\n", config)
}
```

### Validation After Parsing

```go
package main

import (
    "fmt"
    "log"
    "github.com/common-library/go/command-line/flags"
)

func main() {
    err := flags.Parse([]flags.FlagInfo{
        {FlagName: "port", Usage: "server port", DefaultValue: 8080},
        {FlagName: "workers", Usage: "worker count", DefaultValue: 4},
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Validate port range
    port := flags.Get[int]("port")
    if port < 1024 || port > 65535 {
        log.Fatalf("Port must be between 1024 and 65535, got %d", port)
    }
    
    // Validate worker count
    workers := flags.Get[int]("workers")
    if workers < 1 || workers > 100 {
        log.Fatalf("Workers must be between 1 and 100, got %d", workers)
    }
    
    fmt.Printf("Valid configuration: port=%d, workers=%d\n", port, workers)
}
```

### Building Configuration Structs

```go
package main

import (
    "fmt"
    "log"
    "time"
    "github.com/common-library/go/command-line/flags"
)

type AppConfig struct {
    Server   ServerConfig
    Database DatabaseConfig
    Logging  LoggingConfig
}

type ServerConfig struct {
    Port    int
    Host    string
    Timeout time.Duration
}

type DatabaseConfig struct {
    URL         string
    MaxConns    uint
    IdleTimeout time.Duration
}

type LoggingConfig struct {
    Level   string
    Verbose bool
}

func main() {
    err := flags.Parse([]flags.FlagInfo{
        // Server flags
        {FlagName: "server-port", Usage: "server port", DefaultValue: 8080},
        {FlagName: "server-host", Usage: "server host", DefaultValue: "0.0.0.0"},
        {FlagName: "server-timeout", Usage: "server timeout", DefaultValue: 30 * time.Second},
        
        // Database flags
        {FlagName: "db-url", Usage: "database URL", DefaultValue: "postgres://localhost/db"},
        {FlagName: "db-max-conns", Usage: "max database connections", DefaultValue: uint(10)},
        {FlagName: "db-idle-timeout", Usage: "connection idle timeout", DefaultValue: 5 * time.Minute},
        
        // Logging flags
        {FlagName: "log-level", Usage: "logging level", DefaultValue: "info"},
        {FlagName: "log-verbose", Usage: "verbose logging", DefaultValue: false},
    })
    if err != nil {
        log.Fatal(err)
    }
    
    config := AppConfig{
        Server: ServerConfig{
            Port:    flags.Get[int]("server-port"),
            Host:    flags.Get[string]("server-host"),
            Timeout: flags.Get[time.Duration]("server-timeout"),
        },
        Database: DatabaseConfig{
            URL:         flags.Get[string]("db-url"),
            MaxConns:    flags.Get[uint]("db-max-conns"),
            IdleTimeout: flags.Get[time.Duration]("db-idle-timeout"),
        },
        Logging: LoggingConfig{
            Level:   flags.Get[string]("log-level"),
            Verbose: flags.Get[bool]("log-verbose"),
        },
    }
    
    fmt.Printf("Configuration loaded:\n")
    fmt.Printf("  Server: %s:%d (timeout: %s)\n", 
        config.Server.Host, config.Server.Port, config.Server.Timeout)
    fmt.Printf("  Database: %s (max conns: %d)\n", 
        config.Database.URL, config.Database.MaxConns)
    fmt.Printf("  Logging: level=%s, verbose=%v\n", 
        config.Logging.Level, config.Logging.Verbose)
}
```

## Best Practices

### 1. Use Explicit Types

Always use explicit type conversions for numeric literals to avoid ambiguity:

```go
// Good: Explicit types
{FlagName: "port", DefaultValue: int(8080)}
{FlagName: "count", DefaultValue: int64(1000)}
{FlagName: "size", DefaultValue: uint64(1024)}

// Avoid: Untyped constants (may be int or int64)
{FlagName: "port", DefaultValue: 8080}
```

### 2. Provide Clear Usage Messages

Include units, formats, and examples in usage strings:

```go
{
    FlagName:     "timeout",
    Usage:        "request timeout duration (e.g., 30s, 5m, 1h)",
    DefaultValue: 30 * time.Second,
}

{
    FlagName:     "rate",
    Usage:        "sampling rate between 0.0 and 1.0",
    DefaultValue: 0.01,
}
```

### 3. Match Types in Get

The type parameter in `Get[T]()` must match the type in `DefaultValue`:

```go
// Defined as int
{FlagName: "port", DefaultValue: int(8080)}

// Retrieve as int (correct)
port := flags.Get[int]("port")  // ✅

// Retrieve as int64 (wrong - runtime panic!)
port := flags.Get[int64]("port")  // ❌
```

### 4. Always Check Parse Errors

```go
err := flags.Parse(flagInfos)
if err != nil {
    log.Fatalf("Failed to parse flags: %v", err)
}
```

### 5. Validate Values After Parsing

The package doesn't validate ranges or formats - add your own validation:

```go
port := flags.Get[int]("port")
if port < 1024 || port > 65535 {
    log.Fatalf("Invalid port: %d", port)
}
```

### 6. Group Related Flags

Use consistent naming prefixes for related flags:

```go
// Server flags
{FlagName: "server-port", ...}
{FlagName: "server-host", ...}
{FlagName: "server-timeout", ...}

// Database flags
{FlagName: "db-url", ...}
{FlagName: "db-max-conns", ...}
```

## Comparison with Standard Library

| Feature | This Package | Standard `flag` |
|---------|--------------|-----------------|
| Type-safe retrieval | ✅ Generic `Get[T]()` | ❌ Type-specific getters (`flag.Int()`, `flag.String()`, etc.) |
| Declarative definition | ✅ Single `FlagInfo` slice | ❌ Imperative calls for each flag |
| Error handling | ✅ Returns error | ❌ Panics on some errors |
| All types in one call | ✅ Single `Parse()` call | ❌ Multiple variable declarations |
| Custom types | ❌ Only 8 predefined types | ✅ Supports `flag.Value` interface |

**Example comparison:**

```go
// This package (declarative)
flags.Parse([]flags.FlagInfo{
    {FlagName: "port", DefaultValue: 8080},
    {FlagName: "host", DefaultValue: "localhost"},
})
port := flags.Get[int]("port")
host := flags.Get[string]("host")

// Standard library (imperative)
var port = flag.Int("port", 8080, "server port")
var host = flag.String("host", "localhost", "server host")
flag.Parse()
// Use *port and *host (pointers)
```

## Limitations

1. **No Custom Types** - Only supports 8 predefined types (no `flag.Value` interface)
2. **No Validation** - No built-in value validation (range checks, format validation)
3. **No Subcommands** - No native support for subcommand patterns (like `git commit`)
4. **No Short Flags** - No automatic short flag aliases (e.g., `-p` for `--port`)
5. **Global State** - Uses package-level storage for parsed flags
6. **Runtime Type Panic** - Type mismatch in `Get[T]()` causes panic, not compile error
7. **No Flag Aliases** - Cannot define multiple names for the same flag

## Dependencies

- `flag` - Go standard library
- `fmt` - Go standard library
- `time` - Go standard library
- `github.com/common-library/go/utility` - Utility functions (`GetTypeName`)

## Related Packages

- [arguments](../arguments/) - Simple positional argument access
- [flag](https://pkg.go.dev/flag) - Go standard library flag package
- Popular alternatives:
  - [spf13/pflag](https://github.com/spf13/pflag) - POSIX/GNU-style flags
  - [spf13/cobra](https://github.com/spf13/cobra) - CLI framework with subcommands
  - [urfave/cli](https://github.com/urfave/cli) - Full CLI application framework

## Further Reading

- [Go flag package documentation](https://pkg.go.dev/flag)
- [Command-line flags in Go](https://gobyexample.com/command-line-flags)
- [Go generics tutorial](https://go.dev/doc/tutorial/generics)
- [time.Duration format](https://pkg.go.dev/time#ParseDuration)
