# Arguments

Simple utilities for accessing command-line arguments in Go.

## Overview

The arguments package provides a straightforward wrapper around Go's `os.Args`, offering convenient methods to access command-line arguments by index or retrieve them all at once.

This is a lightweight utility package - if you need named flags with type safety, see the [flags package](../flags/).

## Features

- Individual argument access by index
- Retrieve all arguments as a slice
- Simple wrapper around `os.Args`
- No external dependencies (uses only standard library)

## Installation

```bash
go get -u github.com/common-library/go/command-line/arguments
```

## Quick Start

```go
import "github.com/common-library/go/command-line/arguments"

func main() {
    // Get program name
    programName := arguments.Get(0)
    
    // Get first argument
    firstArg := arguments.Get(1)
    
    // Get all arguments
    allArgs := arguments.GetAll()
    
    fmt.Printf("Program: %s\n", programName)
    fmt.Printf("First argument: %s\n", firstArg)
    fmt.Printf("Total arguments: %d\n", len(allArgs))
}
```

## Usage Examples

### Basic Argument Access

```go
package main

import (
    "fmt"
    "github.com/common-library/go/command-line/arguments"
)

func main() {
    // Command: ./program arg1 arg2 arg3
    
    programName := arguments.Get(0)  // "./program"
    firstArg := arguments.Get(1)     // "arg1"
    secondArg := arguments.Get(2)    // "arg2"
    thirdArg := arguments.Get(3)     // "arg3"
    
    fmt.Printf("Running: %s\n", programName)
    fmt.Printf("Arguments: %s, %s, %s\n", firstArg, secondArg, thirdArg)
}
```

### Get All Arguments

```go
package main

import (
    "fmt"
    "github.com/common-library/go/command-line/arguments"
)

func main() {
    // Command: ./program file1.txt file2.txt file3.txt
    
    args := arguments.GetAll()
    
    // Skip program name (index 0)
    files := args[1:]
    
    fmt.Printf("Processing %d files:\n", len(files))
    for i, file := range files {
        fmt.Printf("%d: %s\n", i+1, file)
    }
}
```

### Validate Argument Count

```go
package main

import (
    "fmt"
    "os"
    "github.com/common-library/go/command-line/arguments"
)

func main() {
    args := arguments.GetAll()
    
    if len(args) < 3 {
        fmt.Printf("Usage: %s <source> <destination>\n", args[0])
        os.Exit(1)
    }
    
    source := arguments.Get(1)
    destination := arguments.Get(2)
    
    fmt.Printf("Copying %s to %s\n", source, destination)
}
```

### Simple Command Dispatcher

```go
package main

import (
    "fmt"
    "os"
    "github.com/common-library/go/command-line/arguments"
)

func main() {
    args := arguments.GetAll()
    
    if len(args) < 2 {
        fmt.Println("Usage: program <command> [args...]")
        fmt.Println("Commands: start, stop, restart, status")
        os.Exit(1)
    }
    
    command := arguments.Get(1)
    
    switch command {
    case "start":
        fmt.Println("Starting service...")
    case "stop":
        fmt.Println("Stopping service...")
    case "restart":
        fmt.Println("Restarting service...")
    case "status":
        fmt.Println("Checking status...")
    default:
        fmt.Printf("Unknown command: %s\n", command)
        os.Exit(1)
    }
}
```

### File Batch Processing

```go
package main

import (
    "fmt"
    "path/filepath"
    "github.com/common-library/go/command-line/arguments"
)

func main() {
    // Command: ./processor *.txt
    
    args := arguments.GetAll()
    
    if len(args) < 2 {
        fmt.Println("Usage: processor <files...>")
        return
    }
    
    // Get all file arguments (skip program name)
    files := args[1:]
    
    for _, file := range files {
        ext := filepath.Ext(file)
        fmt.Printf("Processing %s (type: %s)\n", file, ext)
        // Process file...
    }
    
    fmt.Printf("Processed %d files\n", len(files))
}
```

## API Reference

### `Get(index int) string`

Returns the command-line argument at the specified index.

**Parameters:**
- `index` - Position of the argument (0 is the program name)

**Returns:**
- The argument string at the specified index

**Panics:**
- If index is out of bounds

**Example:**
```go
// Command: ./program arg1 arg2
programName := arguments.Get(0)  // "./program"
firstArg := arguments.Get(1)     // "arg1"
secondArg := arguments.Get(2)    // "arg2"
```

### `GetAll() []string`

Returns all command-line arguments including the program name.

**Returns:**
- A slice containing all command-line arguments

**Example:**
```go
// Command: ./program arg1 arg2 --flag=value
args := arguments.GetAll()
// args = ["./program", "arg1", "arg2", "--flag=value"]

fmt.Printf("Total arguments: %d\n", len(args))
```

## Best Practices

### 1. Always Validate Argument Count

```go
args := arguments.GetAll()
if len(args) < 2 {
    fmt.Println("Usage: program <argument>")
    os.Exit(1)
}
```

### 2. Handle Index Out of Bounds

```go
// Good: Check length first
args := arguments.GetAll()
if len(args) > 1 {
    firstArg := arguments.Get(1)
    // Use firstArg...
}

// Avoid: Direct access without checking
// firstArg := arguments.Get(1)  // May panic if no arguments
```

### 3. Skip Program Name When Iterating

```go
args := arguments.GetAll()
userArgs := args[1:]  // Skip program name at index 0

for _, arg := range userArgs {
    // Process user arguments only
}
```

### 4. Provide Usage Information

```go
if len(arguments.GetAll()) < 2 {
    fmt.Printf("Usage: %s <input> <output>\n", arguments.Get(0))
    fmt.Println("  input  - Source file path")
    fmt.Println("  output - Destination file path")
    os.Exit(1)
}
```

### 5. Consider Using Flags for Complex Arguments

If you need named parameters, default values, or type safety, use the [flags package](../flags/) instead:

```go
// Simple positional arguments - Use arguments package
// ./program input.txt output.txt

// Named flags with defaults - Use flags package
// ./program -input=file.txt -output=result.txt -verbose=true
```

## When to Use This Package

**Use the arguments package when:**
- ✅ You have simple positional arguments
- ✅ Argument order is meaningful
- ✅ You need minimal overhead
- ✅ You're wrapping os.Args for convenience

**Use the flags package when:**
- ❌ You need named parameters
- ❌ You want default values
- ❌ You need type safety (int, bool, duration, etc.)
- ❌ You want automatic help generation

## Comparison with Standard Library

| Feature | This Package | Standard `os.Args` |
|---------|--------------|-------------------|
| Get by index | `arguments.Get(i)` | `os.Args[i]` |
| Get all | `arguments.GetAll()` | `os.Args` |
| API clarity | ✅ Explicit function names | ⚠️ Direct array access |
| Error handling | ❌ Panics on out of bounds | ❌ Panics on out of bounds |

This package is essentially a thin wrapper that may improve code readability in some contexts.

## Limitations

1. **No Type Conversion**: All arguments are strings
2. **No Validation**: No built-in validation for argument values
3. **No Default Values**: No support for optional arguments with defaults
4. **Panic on Invalid Index**: No graceful error handling for out-of-bounds access
5. **No Named Parameters**: Only supports positional arguments

For type-safe, named parameters with validation, see the [flags package](../flags/).

## Dependencies

- `os` - Go standard library

## Related Packages

- [flags](../flags/) - Type-safe command-line flag parsing
- [os package](https://pkg.go.dev/os) - Go standard library

## Further Reading

- [os.Args documentation](https://pkg.go.dev/os#pkg-variables)
- [Command-line arguments in Go](https://gobyexample.com/command-line-arguments)
- [Flag package for named parameters](https://pkg.go.dev/flag)
