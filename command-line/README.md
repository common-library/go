# Command-Line

Utilities for parsing and accessing command-line arguments and flags in Go.

## Overview

This package provides two complementary utilities for handling command-line input:

- **[arguments](arguments/)** - Simple access to positional command-line arguments
- **[flags](flags/)** - Type-safe parsing of named command-line flags

Both packages are built on Go's standard library and offer convenient wrappers with improved ergonomics.

## Packages

### Arguments Package

Provides simple access to command-line arguments via wrapper functions around `os.Args`.

**Features:**
- Individual argument access by index
- Retrieve all arguments as a slice
- Minimal overhead wrapper

**Documentation:** [arguments/README.md](arguments/)

**Quick Example:**
```go
import "github.com/common-library/go/command-line/arguments"

programName := arguments.Get(0)
firstArg := arguments.Get(1)
allArgs := arguments.GetAll()
```

### Flags Package

Provides type-safe command-line flag parsing with generic value retrieval.

**Features:**
- Type-safe flag parsing (8 supported types)
- Generic value retrieval with `Get[T]()`
- Declarative flag definition
- Built on Go's standard `flag` package

**Documentation:** [flags/README.md](flags/)

**Quick Example:**
```go
import "github.com/common-library/go/command-line/flags"

flags.Parse([]flags.FlagInfo{
    {FlagName: "port", Usage: "server port", DefaultValue: 8080},
    {FlagName: "host", Usage: "server host", DefaultValue: "localhost"},
})

port := flags.Get[int]("port")
host := flags.Get[string]("host")
```

## Installation

```bash
# Install both packages
go get -u github.com/common-library/go/command-line

# Or install individually
go get -u github.com/common-library/go/command-line/arguments
go get -u github.com/common-library/go/command-line/flags
```

## When to Use Which Package

### Use Arguments Package When:

✅ You have simple positional arguments  
✅ Argument order is meaningful and fixed  
✅ You need minimal overhead  
✅ You're building simple CLI tools  

**Example scenarios:**
```bash
# File operations
./copy source.txt destination.txt

# Simple commands
./program start
./program stop

# Basic utilities
./converter input.json output.xml
```

### Use Flags Package When:

✅ You need named parameters  
✅ You want default values  
✅ You need type safety (int, bool, duration, etc.)  
✅ You want automatic help generation  

**Example scenarios:**
```bash
# Server configuration
./server -port=8080 -host=localhost -debug=true

# Build tools
./build -output=dist -workers=4 -verbose

# Data processing
./processor -input=data.csv -format=json -timeout=30s
```

### Use Both Together:

```go
import (
    "github.com/common-library/go/command-line/arguments"
    "github.com/common-library/go/command-line/flags"
)

func main() {
    // Get command from arguments
    if len(arguments.GetAll()) < 2 {
        fmt.Println("Usage: program <command> [flags]")
        os.Exit(1)
    }
    command := arguments.Get(1)
    
    // Parse flags
    flags.Parse([]flags.FlagInfo{
        {FlagName: "verbose", Usage: "verbose output", DefaultValue: false},
        {FlagName: "config", Usage: "config file", DefaultValue: "config.json"},
    })
    
    verbose := flags.Get[bool]("verbose")
    config := flags.Get[string]("config")
    
    // Execute command with flags
    switch command {
    case "start":
        startServer(verbose, config)
    case "stop":
        stopServer(verbose)
    }
}
```

```bash
# Usage
./program start -verbose=true -config=prod.json
./program stop -verbose
```

## Complete Example

### Server Application

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    "github.com/common-library/go/command-line/arguments"
    "github.com/common-library/go/command-line/flags"
)

func main() {
    // Get program name from arguments
    programName := arguments.Get(0)
    fmt.Printf("Starting %s...\n", programName)
    
    // Parse configuration flags
    err := flags.Parse([]flags.FlagInfo{
        {FlagName: "port", Usage: "HTTP server port", DefaultValue: 8080},
        {FlagName: "host", Usage: "HTTP server host", DefaultValue: "0.0.0.0"},
        {FlagName: "timeout", Usage: "request timeout", DefaultValue: 30 * time.Second},
        {FlagName: "debug", Usage: "enable debug mode", DefaultValue: false},
    })
    if err != nil {
        log.Fatalf("Failed to parse flags: %v", err)
    }
    
    // Build configuration
    config := ServerConfig{
        Port:    flags.Get[int]("port"),
        Host:    flags.Get[string]("host"),
        Timeout: flags.Get[time.Duration]("timeout"),
        Debug:   flags.Get[bool]("debug"),
    }
    
    fmt.Printf("Server: %s:%d (timeout: %s, debug: %v)\n",
        config.Host, config.Port, config.Timeout, config.Debug)
    
    // Start server...
}

type ServerConfig struct {
    Port    int
    Host    string
    Timeout time.Duration
    Debug   bool
}
```

**Usage:**
```bash
# Use defaults
./server

# Custom configuration
./server -port=9000 -host=localhost -debug=true

# With help
./server -h
```

## Quick Reference

### Arguments Package

| Function | Description | Example |
|----------|-------------|---------|
| `Get(index int)` | Get argument by index | `arg := arguments.Get(1)` |
| `GetAll()` | Get all arguments | `args := arguments.GetAll()` |

[Full documentation →](arguments/)

### Flags Package

| Function | Description | Example |
|----------|-------------|---------|
| `Parse([]FlagInfo)` | Parse flags | `flags.Parse(flagInfos)` |
| `Get[T](name)` | Get flag value | `port := flags.Get[int]("port")` |

**Supported Types:** `bool`, `int`, `int64`, `uint`, `uint64`, `float64`, `string`, `time.Duration`

[Full documentation →](flags/)

## Comparison with Standard Library

| Feature | arguments | flags | `os.Args` | `flag` |
|---------|-----------|-------|-----------|--------|
| Positional arguments | ✅ | ❌ | ✅ | ❌ |
| Named flags | ❌ | ✅ | ❌ | ✅ |
| Type safety | ❌ | ✅ | ❌ | ⚠️ |
| Generic retrieval | N/A | ✅ | N/A | ❌ |
| Declarative API | N/A | ✅ | N/A | ❌ |
| Default values | N/A | ✅ | N/A | ✅ |

## Best Practices

### 1. Validate Input

Always validate argument count and flag values:

```go
// Arguments
args := arguments.GetAll()
if len(args) < 2 {
    fmt.Println("Usage: program <command>")
    os.Exit(1)
}

// Flags
port := flags.Get[int]("port")
if port < 1024 || port > 65535 {
    log.Fatalf("Invalid port: %d", port)
}
```

### 2. Provide Clear Help

```go
// Show usage for arguments
if len(arguments.GetAll()) < 2 {
    fmt.Println("Usage: program <source> <destination>")
    fmt.Println("  source      - Input file path")
    fmt.Println("  destination - Output file path")
    os.Exit(1)
}

// Flags automatically support -h
flags.Parse([]flags.FlagInfo{
    {FlagName: "port", Usage: "server port (1024-65535)", DefaultValue: 8080},
})
```

### 3. Use Explicit Types

For flags, always use explicit type conversions:

```go
// Good
{FlagName: "port", DefaultValue: int(8080)}

// Avoid (ambiguous)
{FlagName: "port", DefaultValue: 8080}
```

### 4. Combine Wisely

Use arguments for commands, flags for options:

```bash
./tool <command> -flag1=value1 -flag2=value2
```

```go
command := arguments.Get(1)  // start, stop, restart
flags.Parse(configFlags)     // -verbose, -config, etc.
```

## Common Patterns

### Subcommand Pattern

```go
func main() {
    args := arguments.GetAll()
    if len(args) < 2 {
        printUsage()
        os.Exit(1)
    }
    
    command := arguments.Get(1)
    
    switch command {
    case "serve":
        serveCommand()
    case "build":
        buildCommand()
    case "test":
        testCommand()
    default:
        fmt.Printf("Unknown command: %s\n", command)
        os.Exit(1)
    }
}

func serveCommand() {
    flags.Parse([]flags.FlagInfo{
        {FlagName: "port", Usage: "server port", DefaultValue: 8080},
    })
    // Serve...
}
```

### Configuration Override Pattern

```go
func main() {
    // Load defaults from config file
    config := loadConfig("config.json")
    
    // Override with flags
    flags.Parse([]flags.FlagInfo{
        {FlagName: "port", DefaultValue: config.Port},
        {FlagName: "host", DefaultValue: config.Host},
    })
    
    config.Port = flags.Get[int]("port")
    config.Host = flags.Get[string]("host")
    
    // Use merged config
}
```

## Dependencies

- **arguments**: `os` (Go standard library)
- **flags**: `flag`, `fmt`, `time` (Go standard library), `github.com/common-library/go/utility`

## Further Reading

- [arguments package documentation](arguments/)
- [flags package documentation](flags/)
- [Go flag package](https://pkg.go.dev/flag)
- [os.Args documentation](https://pkg.go.dev/os#pkg-variables)
- [Command-line applications in Go](https://gobyexample.com/command-line-arguments)
