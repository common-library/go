# JSON

Utilities for converting between JSON and Go structs.

## Overview

The json package provides convenient wrapper functions around Go's encoding/json package, simplifying common JSON operations like marshaling, unmarshaling, and file I/O with type-safe generic functions.

## Features

- **Struct to JSON** - Convert Go structs to JSON strings
- **Formatted Output** - Generate indented, human-readable JSON
- **Generic File Reading** - Type-safe JSON file deserialization
- **Generic String Parsing** - Type-safe JSON string deserialization
- **Simple API** - Minimal wrapper around standard library

## Installation

```bash
go get -u github.com/common-library/go/json
```

## Quick Start

### Marshal to JSON

```go
import "github.com/common-library/go/json"

type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}

user := User{
    Name:  "Alice",
    Email: "alice@example.com",
    Age:   30,
}

// Compact JSON
jsonStr, err := json.ToString(user)
// Output: {"name":"Alice","email":"alice@example.com","age":30}

// Formatted JSON
jsonStr, err := json.ToStringIndent(user, "", "  ")
// Output:
// {
//   "name": "Alice",
//   "email": "alice@example.com",
//   "age": 30
// }
```

### Unmarshal from JSON

```go
// From string
jsonData := `{"name":"Bob","email":"bob@example.com","age":25}`
user, err := json.ConvertFromString[User](jsonData)

// From file
user, err := json.ConvertFromFile[User]("user.json")
```

## API Reference

### ToString

```go
func ToString(input any) (string, error)
```

Converts a Go value to a compact JSON string.

**Parameters:**
- `input` - Any Go value to convert (typically a struct)

**Returns:**
- `string` - Compact JSON representation
- `error` - Error if marshaling fails

**Behavior:**
- Uses json.Marshal internally
- Produces compact output with no whitespace
- Exported struct fields only
- Respects json struct tags

### ToStringIndent

```go
func ToStringIndent(input any, prefix string, indent string) (string, error)
```

Converts a Go value to a formatted JSON string with custom indentation.

**Parameters:**
- `input` - Any Go value to convert
- `prefix` - String prepended to each line (typically "")
- `indent` - String used for each indentation level (e.g., "\t" or "  ")

**Returns:**
- `string` - Formatted JSON representation
- `error` - Error if marshaling fails

**Behavior:**
- Uses json.MarshalIndent internally
- Produces human-readable output
- Adds prefix to beginning of each line
- Adds indent for each nesting level

### ConvertFromFile

```go
func ConvertFromFile[T any](fileName string) (T, error)
```

Reads a JSON file and converts it to type T.

**Type Parameters:**
- `T` - Target type (struct, slice, map, or primitive)

**Parameters:**
- `fileName` - Path to JSON file

**Returns:**
- `T` - Unmarshaled data as type T
- `error` - Error if file reading or unmarshaling fails

**Behavior:**
- Reads entire file into memory
- Unmarshals JSON to specified type
- Type-safe with compile-time checking

### ConvertFromString

```go
func ConvertFromString[T any](data string) (T, error)
```

Converts a JSON string to type T.

**Type Parameters:**
- `T` - Target type (struct, slice, map, or primitive)

**Parameters:**
- `data` - JSON string to parse

**Returns:**
- `T` - Unmarshaled data as type T
- `error` - Error if unmarshaling fails

**Behavior:**
- Uses json.Unmarshal internally
- Type-safe with compile-time checking
- Handles nested structures

## Complete Examples

### Struct to JSON

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/common-library/go/json"
)

type Address struct {
    Street  string `json:"street"`
    City    string `json:"city"`
    Country string `json:"country"`
    ZipCode string `json:"zip_code"`
}

type User struct {
    ID      int     `json:"id"`
    Name    string  `json:"name"`
    Email   string  `json:"email"`
    Age     int     `json:"age"`
    Address Address `json:"address"`
}

func main() {
    user := User{
        ID:    1,
        Name:  "Alice Johnson",
        Email: "alice@example.com",
        Age:   30,
        Address: Address{
            Street:  "123 Main St",
            City:    "New York",
            Country: "USA",
            ZipCode: "10001",
        },
    }
    
    // Compact JSON
    compact, err := json.ToString(user)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Compact JSON:")
    fmt.Println(compact)
    
    // Formatted JSON
    formatted, err := json.ToStringIndent(user, "", "  ")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("\nFormatted JSON:")
    fmt.Println(formatted)
}
```

### JSON to Struct

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/common-library/go/json"
)

type Config struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    Database string `json:"database"`
    Username string `json:"username"`
    Password string `json:"password"`
}

func main() {
    // From string
    jsonData := `{
        "host": "localhost",
        "port": 5432,
        "database": "myapp",
        "username": "admin",
        "password": "secret"
    }`
    
    config, err := json.ConvertFromString[Config](jsonData)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Config: %+v\n", config)
}
```

### Read JSON File

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/common-library/go/json"
)

type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}

func main() {
    // users.json contains:
    // [
    //   {"name": "Alice", "email": "alice@example.com", "age": 30},
    //   {"name": "Bob", "email": "bob@example.com", "age": 25}
    // ]
    
    users, err := json.ConvertFromFile[[]User]("users.json")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Loaded %d users:\n", len(users))
    for _, user := range users {
        fmt.Printf("  %s (%s) - %d years old\n", user.Name, user.Email, user.Age)
    }
}
```

### Write JSON File

```go
package main

import (
    "log"
    
    "github.com/common-library/go/file"
    "github.com/common-library/go/json"
)

type Settings struct {
    Theme      string `json:"theme"`
    Language   string `json:"language"`
    AutoSave   bool   `json:"auto_save"`
    MaxBackups int    `json:"max_backups"`
}

func main() {
    settings := Settings{
        Theme:      "dark",
        Language:   "en",
        AutoSave:   true,
        MaxBackups: 5,
    }
    
    // Convert to formatted JSON
    jsonData, err := json.ToStringIndent(settings, "", "  ")
    if err != nil {
        log.Fatal(err)
    }
    
    // Write to file
    err = file.Write("settings.json", jsonData)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Println("Settings saved to settings.json")
}
```

### Working with Maps

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/common-library/go/json"
)

func main() {
    // Marshal map to JSON
    data := map[string]interface{}{
        "name": "Alice",
        "age":  30,
        "active": true,
        "tags": []string{"admin", "developer"},
    }
    
    jsonStr, err := json.ToStringIndent(data, "", "  ")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Map as JSON:")
    fmt.Println(jsonStr)
    
    // Unmarshal JSON to map
    jsonData := `{
        "server": "localhost",
        "port": 8080,
        "debug": true
    }`
    
    config, err := json.ConvertFromString[map[string]interface{}](jsonData)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("\nParsed config: %+v\n", config)
    fmt.Printf("Server: %s\n", config["server"])
    fmt.Printf("Port: %v\n", config["port"])
    fmt.Printf("Debug: %v\n", config["debug"])
}
```

### Slices and Arrays

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/common-library/go/json"
)

type Product struct {
    ID    int     `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}

func main() {
    // Marshal slice to JSON
    products := []Product{
        {ID: 1, Name: "Laptop", Price: 999.99},
        {ID: 2, Name: "Mouse", Price: 29.99},
        {ID: 3, Name: "Keyboard", Price: 79.99},
    }
    
    jsonStr, err := json.ToStringIndent(products, "", "  ")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Products as JSON:")
    fmt.Println(jsonStr)
    
    // Unmarshal JSON array
    jsonData := `[
        {"id": 10, "name": "Monitor", "price": 299.99},
        {"id": 11, "name": "Webcam", "price": 89.99}
    ]`
    
    newProducts, err := json.ConvertFromString[[]Product](jsonData)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("\nLoaded %d products:\n", len(newProducts))
    for _, p := range newProducts {
        fmt.Printf("  %s: $%.2f\n", p.Name, p.Price)
    }
}
```

### Nested Structures

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/common-library/go/json"
)

type Company struct {
    Name      string       `json:"name"`
    Founded   int          `json:"founded"`
    Employees []Employee   `json:"employees"`
    Addresses []Address    `json:"addresses"`
}

type Employee struct {
    Name     string `json:"name"`
    Position string `json:"position"`
    Salary   int    `json:"salary"`
}

type Address struct {
    Type    string `json:"type"` // e.g., "headquarters", "branch"
    Street  string `json:"street"`
    City    string `json:"city"`
    Country string `json:"country"`
}

func main() {
    company := Company{
        Name:    "Tech Corp",
        Founded: 2010,
        Employees: []Employee{
            {Name: "Alice", Position: "CEO", Salary: 200000},
            {Name: "Bob", Position: "CTO", Salary: 180000},
        },
        Addresses: []Address{
            {Type: "headquarters", Street: "100 Tech Blvd", City: "San Francisco", Country: "USA"},
            {Type: "branch", Street: "200 Innovation Way", City: "Austin", Country: "USA"},
        },
    }
    
    // Convert to JSON
    jsonStr, err := json.ToStringIndent(company, "", "  ")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Company as JSON:")
    fmt.Println(jsonStr)
    
    // Parse back
    parsed, err := json.ConvertFromString[Company](jsonStr)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("\nParsed: %s with %d employees\n", parsed.Name, len(parsed.Employees))
}
```

### API Response Handling

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/common-library/go/json"
)

type APIResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}

type UserData struct {
    ID       int      `json:"id"`
    Username string   `json:"username"`
    Roles    []string `json:"roles"`
}

func main() {
    // Parse API response
    responseJSON := `{
        "success": true,
        "message": "User retrieved successfully",
        "data": {
            "id": 123,
            "username": "alice",
            "roles": ["admin", "user"]
        }
    }`
    
    response, err := json.ConvertFromString[APIResponse](responseJSON)
    if err != nil {
        log.Fatal(err)
    }
    
    if response.Success {
        fmt.Printf("Message: %s\n", response.Message)
        
        // Convert data field to specific type
        dataJSON, _ := json.ToString(response.Data)
        userData, err := json.ConvertFromString[UserData](dataJSON)
        if err != nil {
            log.Fatal(err)
        }
        
        fmt.Printf("User: %s (ID: %d)\n", userData.Username, userData.ID)
        fmt.Printf("Roles: %v\n", userData.Roles)
    }
}
```

## Best Practices

### 1. Use Struct Tags

```go
// Good: Use struct tags to control JSON field names
type User struct {
    ID        int    `json:"id"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Email     string `json:"email"`
}

// Avoid: Relying on default field names
type User struct {
    ID        int    // Becomes "ID" in JSON
    FirstName string // Becomes "FirstName" in JSON
}
```

### 2. Handle Errors Properly

```go
// Good: Always check errors
jsonStr, err := json.ToString(data)
if err != nil {
    log.Printf("Failed to marshal data: %v", err)
    return err
}

// Avoid: Ignoring errors
jsonStr, _ := json.ToString(data)
```

### 3. Use Type Parameters Correctly

```go
// Good: Explicit type parameter
users, err := json.ConvertFromFile[[]User]("users.json")

config, err := json.ConvertFromString[Config](jsonData)

// Avoid: Wrong type parameter
// users, err := json.ConvertFromFile[User]("users.json")
// Will fail if file contains array
```

### 4. Validate JSON Data

```go
// Good: Validate after unmarshaling
user, err := json.ConvertFromString[User](jsonData)
if err != nil {
    return fmt.Errorf("invalid JSON: %w", err)
}

if user.Email == "" {
    return fmt.Errorf("email is required")
}

if user.Age < 0 {
    return fmt.Errorf("invalid age: %d", user.Age)
}
```

### 5. Use Appropriate Indentation

```go
// Good: Use 2 or 4 spaces for readability
jsonStr, _ := json.ToStringIndent(data, "", "  ")

// Good: Use tabs for compact formatting
jsonStr, _ := json.ToStringIndent(data, "", "\t")

// Avoid: Inconsistent indentation
jsonStr, _ := json.ToStringIndent(data, "", "   ") // 3 spaces
```

## Common Use Cases

### Configuration Files

```go
type Config struct {
    Database DatabaseConfig `json:"database"`
    Server   ServerConfig   `json:"server"`
    Logging  LogConfig      `json:"logging"`
}

// Load config
config, err := json.ConvertFromFile[Config]("config.json")

// Save config
jsonData, _ := json.ToStringIndent(config, "", "  ")
file.Write("config.json", jsonData)
```

### API Communication

```go
// Request payload
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

request := CreateUserRequest{Name: "Alice", Email: "alice@example.com"}
jsonBody, _ := json.ToString(request)

// Send HTTP request with jsonBody
// ...

// Parse response
type CreateUserResponse struct {
    ID      int    `json:"id"`
    Message string `json:"message"`
}

response, err := json.ConvertFromString[CreateUserResponse](responseBody)
```

### Data Export

```go
// Export data to JSON file
type Report struct {
    Title     string    `json:"title"`
    Generated time.Time `json:"generated"`
    Data      []Record  `json:"data"`
}

report := Report{
    Title:     "Monthly Report",
    Generated: time.Now(),
    Data:      records,
}

jsonData, _ := json.ToStringIndent(report, "", "  ")
file.Write("report.json", jsonData)
```

## Error Handling

### Marshaling Errors

```go
jsonStr, err := json.ToString(data)
if err != nil {
    // Possible errors:
    // - Unsupported type (e.g., channel, function)
    // - Circular reference
    // - Invalid UTF-8
    log.Printf("Marshal error: %v", err)
}
```

### Unmarshaling Errors

```go
user, err := json.ConvertFromString[User](jsonData)
if err != nil {
    // Possible errors:
    // - Invalid JSON syntax
    // - Type mismatch (e.g., string where number expected)
    // - Out of range values
    log.Printf("Unmarshal error: %v", err)
}
```

### File Errors

```go
config, err := json.ConvertFromFile[Config]("config.json")
if err != nil {
    // Possible errors:
    // - File not found
    // - Permission denied
    // - Invalid JSON in file
    log.Printf("File error: %v", err)
}
```

## Performance Tips

1. **Reuse Buffers** - For frequent marshaling, consider using encoding/json directly with buffer reuse
2. **Streaming** - For large files, use json.Decoder/Encoder for streaming instead of loading entire file
3. **Avoid String Conversion** - When working with bytes, use encoding/json directly to avoid string conversions
4. **Struct Tags** - Use omitempty to reduce output size: `json:"field,omitempty"`

## Testing

### Unit Test Example

```go
func TestUserSerialization(t *testing.T) {
    user := User{
        ID:    1,
        Name:  "Alice",
        Email: "alice@example.com",
    }
    
    // Marshal
    jsonStr, err := json.ToString(user)
    if err != nil {
        t.Fatalf("Marshal failed: %v", err)
    }
    
    // Unmarshal
    parsed, err := json.ConvertFromString[User](jsonStr)
    if err != nil {
        t.Fatalf("Unmarshal failed: %v", err)
    }
    
    // Verify
    if parsed.ID != user.ID {
        t.Errorf("Expected ID %d, got %d", user.ID, parsed.ID)
    }
    
    if parsed.Name != user.Name {
        t.Errorf("Expected name %s, got %s", user.Name, parsed.Name)
    }
}
```

## Dependencies

- `encoding/json` - Go standard library
- `github.com/common-library/go/file` - File operations

## Further Reading

- [Go JSON encoding package](https://pkg.go.dev/encoding/json)
- [JSON and Go](https://go.dev/blog/json)
- [JSON Specification](https://www.json.org/)
