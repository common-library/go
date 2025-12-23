# File

Simplified utilities for file and directory operations in Go.

## Overview

The file package provides convenient wrapper functions around Go's standard `os` and `path/filepath` packages, offering a simplified API for common file system operations. It reduces boilerplate code and provides consistent error handling for reading, writing, listing, creating, and removing files and directories.

## Features

- **File I/O** - Read and write entire files with simple function calls
- **Directory Listing** - List files and directories recursively or non-recursively
- **Directory Creation** - Create single or nested directory structures
- **File/Directory Removal** - Delete files, empty directories, or entire directory trees
- **Permission Management** - Set file and directory permissions easily
- **Cross-Platform** - Works on Unix/Linux, macOS, and Windows

## Installation

```bash
go get -u github.com/common-library/go/file
```

## Quick Start

```go
import "github.com/common-library/go/file"

func main() {
    // Read file
    data, err := file.Read("config.txt")
    if err != nil {
        log.Fatal(err)
    }
    
    // Write file
    err = file.Write("output.txt", "Hello, World!", 0644)
    if err != nil {
        log.Fatal(err)
    }
    
    // List files
    files, err := file.List("./data", true)
    for _, f := range files {
        fmt.Println(f)
    }
    
    // Create directory
    err = file.CreateDirectoryAll("data/2024/logs", 0755)
    
    // Remove directory
    err = file.RemoveAll("temp")
}
```

## API Reference

### File Operations

#### Read

```go
func Read(fileName string) (string, error)
```

Reads and returns the entire contents of a file as a string.

**Parameters:**
- `fileName` - Path to the file (absolute or relative)

**Returns:**
- `string` - File contents
- `error` - Error if file doesn't exist or cannot be read

**Example:**

```go
// Read configuration file
config, err := file.Read("config.json")
if err != nil {
    log.Fatal(err)
}
fmt.Println(config)

// Read with error handling
data, err := file.Read("data.txt")
if err != nil {
    if os.IsNotExist(err) {
        log.Println("File does not exist")
    } else {
        log.Printf("Error reading file: %v", err)
    }
    return
}
```

#### Write

```go
func Write(fileName string, data string, fileMode os.FileMode) error
```

Writes data to a file, creating it if it doesn't exist or truncating it if it does.

**Parameters:**
- `fileName` - Path to the file
- `data` - String content to write
- `fileMode` - File permissions (e.g., 0644, 0600, os.ModePerm)

**Returns:**
- `error` - Error if file cannot be written

**Common File Modes:**
- `0644` - `rw-r--r--` (owner read/write, group/others read)
- `0600` - `rw-------` (owner read/write only)
- `0666` - `rw-rw-rw-` (all read/write)
- `os.ModePerm` - `0777` (all permissions)

**Example:**

```go
// Write with standard permissions
err := file.Write("output.txt", "Hello, World!", 0644)
if err != nil {
    log.Fatal(err)
}

// Write sensitive data with restricted permissions
secret := "API_KEY=12345"
err = file.Write(".env", secret, 0600)

// Overwrite existing file
err = file.Write("log.txt", "New log entry\n", 0644)

// Write JSON data
jsonData := `{"name": "Alice", "age": 30}`
err = file.Write("user.json", jsonData, 0644)
```

### Directory Operations

#### List

```go
func List(path string, recursive bool) ([]string, error)
```

Returns a list of files and directories in the specified path.

**Parameters:**
- `path` - Directory path to list
- `recursive` - If true, lists all subdirectories; if false, lists only immediate children

**Returns:**
- `[]string` - Slice of file and directory paths (directories end with separator)
- `error` - Error if path doesn't exist or cannot be read

**Example:**

```go
// List immediate children only
files, err := file.List("./data", false)
if err != nil {
    log.Fatal(err)
}

for _, f := range files {
    if strings.HasSuffix(f, string(filepath.Separator)) {
        fmt.Printf("Directory: %s\n", f)
    } else {
        fmt.Printf("File: %s\n", f)
    }
}

// List all files recursively
allFiles, err := file.List("./project", true)
if err != nil {
    log.Fatal(err)
}

for _, f := range allFiles {
    fmt.Println(f)
}

// Count files by extension
txtCount := 0
for _, f := range allFiles {
    if strings.HasSuffix(f, ".txt") {
        txtCount++
    }
}
fmt.Printf("Found %d .txt files\n", txtCount)
```

#### CreateDirectory

```go
func CreateDirectory(name string, fileMode os.FileMode) error
```

Creates a single directory.

**Parameters:**
- `name` - Directory path to create
- `fileMode` - Directory permissions (e.g., 0755, 0700)

**Returns:**
- `error` - Error if directory exists or parent doesn't exist

**Common Directory Modes:**
- `0755` - `rwxr-xr-x` (owner full, group/others read/execute)
- `0700` - `rwx------` (owner full access only)
- `0775` - `rwxrwxr-x` (owner/group full, others read/execute)
- `os.ModePerm` - `0777` (all permissions)

**Example:**

```go
// Create directory with standard permissions
err := file.CreateDirectory("uploads", 0755)
if err != nil {
    log.Fatal(err)
}

// Create private directory
err = file.CreateDirectory("secrets", 0700)

// Note: Parent must exist
// This will fail if "data" doesn't exist:
err = file.CreateDirectory("data/2024", 0755) // Error!

// Use CreateDirectoryAll for nested directories
err = file.CreateDirectoryAll("data/2024", 0755) // Success
```

#### CreateDirectoryAll

```go
func CreateDirectoryAll(path string, fileMode os.FileMode) error
```

Creates a directory along with all necessary parent directories.

**Parameters:**
- `path` - Directory path to create (can include multiple levels)
- `fileMode` - Directory permissions for all created directories

**Returns:**
- `error` - Error if creation fails (nil if already exists)

**Example:**

```go
// Create nested directory structure
err := file.CreateDirectoryAll("logs/2024/December/app", 0755)
if err != nil {
    log.Fatal(err)
}

// Create project structure
err = file.CreateDirectoryAll("project/src/components", 0755)
err = file.CreateDirectoryAll("project/tests/unit", 0755)
err = file.CreateDirectoryAll("project/build/output", 0755)

// Safe to call multiple times
err = file.CreateDirectoryAll("data", 0755) // Creates if not exists
err = file.CreateDirectoryAll("data", 0755) // No error if already exists
```

#### Remove

```go
func Remove(name string) error
```

Deletes a single file or empty directory.

**Parameters:**
- `name` - Path to the file or empty directory

**Returns:**
- `error` - Error if doesn't exist or directory not empty

**Example:**

```go
// Remove a file
err := file.Remove("temp.txt")
if err != nil {
    if os.IsNotExist(err) {
        log.Println("File already deleted")
    } else {
        log.Fatal(err)
    }
}

// Remove empty directory
err = file.Remove("empty_dir")
if err != nil {
    log.Println("Directory not empty or doesn't exist")
}

// This will fail if directory has contents
err = file.Remove("data") // Error if "data" has files!

// Use RemoveAll for non-empty directories
err = file.RemoveAll("data") // Removes directory and contents
```

#### RemoveAll

```go
func RemoveAll(path string) error
```

Recursively removes a path and all its contents.

**Parameters:**
- `path` - Path to remove

**Returns:**
- `error` - Error if removal fails (nil if doesn't exist)

**Example:**

```go
// Remove directory tree
err := file.RemoveAll("build/output")
if err != nil {
    log.Fatal(err)
}

// Clean up project artifacts
file.RemoveAll("node_modules")
file.RemoveAll("dist")
file.RemoveAll("coverage")

// Safe to call on non-existent paths
err = file.RemoveAll("nonexistent") // Returns nil, no error

// Remove everything in temp directory
tempDir := "temp"
file.RemoveAll(tempDir)
file.CreateDirectory(tempDir, 0755)
```

## Complete Examples

### Configuration File Manager

```go
package main

import (
    "encoding/json"
    "log"
    
    "github.com/common-library/go/file"
)

type Config struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    Database string `json:"database"`
}

func SaveConfig(cfg Config) error {
    data, err := json.MarshalIndent(cfg, "", "  ")
    if err != nil {
        return err
    }
    
    return file.Write("config.json", string(data), 0644)
}

func LoadConfig() (Config, error) {
    var cfg Config
    
    data, err := file.Read("config.json")
    if err != nil {
        return cfg, err
    }
    
    err = json.Unmarshal([]byte(data), &cfg)
    return cfg, err
}

func main() {
    // Save configuration
    cfg := Config{
        Host:     "localhost",
        Port:     5432,
        Database: "myapp",
    }
    
    if err := SaveConfig(cfg); err != nil {
        log.Fatal(err)
    }
    
    // Load configuration
    loadedCfg, err := LoadConfig()
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Config: %+v\n", loadedCfg)
}
```

### Log File Manager

```go
package main

import (
    "fmt"
    "log"
    "path/filepath"
    "strings"
    "time"
    
    "github.com/common-library/go/file"
)

type Logger struct {
    logDir string
}

func NewLogger(dir string) (*Logger, error) {
    if err := file.CreateDirectoryAll(dir, 0755); err != nil {
        return nil, err
    }
    
    return &Logger{logDir: dir}, nil
}

func (l *Logger) Log(message string) error {
    today := time.Now().Format("2006-01-02")
    logFile := filepath.Join(l.logDir, fmt.Sprintf("%s.log", today))
    
    // Read existing logs
    existing, _ := file.Read(logFile)
    
    // Append new log
    timestamp := time.Now().Format("15:04:05")
    newLog := fmt.Sprintf("[%s] %s\n", timestamp, message)
    
    return file.Write(logFile, existing+newLog, 0644)
}

func (l *Logger) GetLogs(date string) (string, error) {
    logFile := filepath.Join(l.logDir, fmt.Sprintf("%s.log", date))
    return file.Read(logFile)
}

func (l *Logger) CleanOldLogs(daysToKeep int) error {
    files, err := file.List(l.logDir, false)
    if err != nil {
        return err
    }
    
    cutoff := time.Now().AddDate(0, 0, -daysToKeep)
    
    for _, f := range files {
        if strings.HasSuffix(f, ".log") {
            // Parse date from filename
            basename := filepath.Base(f)
            dateStr := strings.TrimSuffix(basename, ".log")
            
            fileDate, err := time.Parse("2006-01-02", dateStr)
            if err != nil {
                continue
            }
            
            if fileDate.Before(cutoff) {
                file.Remove(f)
            }
        }
    }
    
    return nil
}

func main() {
    logger, err := NewLogger("logs")
    if err != nil {
        log.Fatal(err)
    }
    
    // Write logs
    logger.Log("Application started")
    logger.Log("Processing request")
    logger.Log("Request completed")
    
    // Read today's logs
    today := time.Now().Format("2006-01-02")
    logs, _ := logger.GetLogs(today)
    fmt.Println(logs)
    
    // Clean logs older than 7 days
    logger.CleanOldLogs(7)
}
```

### Backup Utility

```go
package main

import (
    "fmt"
    "log"
    "path/filepath"
    "time"
    
    "github.com/common-library/go/file"
)

func BackupDirectory(sourceDir, backupDir string) error {
    timestamp := time.Now().Format("20060102-150405")
    backupPath := filepath.Join(backupDir, timestamp)
    
    // Create backup directory
    if err := file.CreateDirectoryAll(backupPath, 0755); err != nil {
        return err
    }
    
    // Get all files recursively
    files, err := file.List(sourceDir, true)
    if err != nil {
        return err
    }
    
    // Copy each file
    for _, srcFile := range files {
        // Skip directories
        if filepath.Ext(srcFile) == "" {
            continue
        }
        
        // Read source file
        data, err := file.Read(srcFile)
        if err != nil {
            log.Printf("Failed to read %s: %v", srcFile, err)
            continue
        }
        
        // Create destination path
        relPath, _ := filepath.Rel(sourceDir, srcFile)
        dstFile := filepath.Join(backupPath, relPath)
        
        // Create parent directory
        dstDir := filepath.Dir(dstFile)
        file.CreateDirectoryAll(dstDir, 0755)
        
        // Write to backup
        if err := file.Write(dstFile, data, 0644); err != nil {
            log.Printf("Failed to write %s: %v", dstFile, err)
        }
    }
    
    return nil
}

func main() {
    err := BackupDirectory("./data", "./backups")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Backup completed successfully")
}
```

### File Tree Walker

```go
package main

import (
    "fmt"
    "path/filepath"
    "strings"
    
    "github.com/common-library/go/file"
)

type FileStats struct {
    TotalFiles      int
    TotalDirectories int
    FilesByExtension map[string]int
    TotalSize        int64
}

func AnalyzeDirectory(path string) (FileStats, error) {
    stats := FileStats{
        FilesByExtension: make(map[string]int),
    }
    
    files, err := file.List(path, true)
    if err != nil {
        return stats, err
    }
    
    for _, f := range files {
        if strings.HasSuffix(f, string(filepath.Separator)) {
            stats.TotalDirectories++
        } else {
            stats.TotalFiles++
            
            ext := filepath.Ext(f)
            if ext == "" {
                ext = "(no extension)"
            }
            stats.FilesByExtension[ext]++
        }
    }
    
    return stats, nil
}

func main() {
    stats, err := AnalyzeDirectory("./project")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    
    fmt.Printf("Total Files: %d\n", stats.TotalFiles)
    fmt.Printf("Total Directories: %d\n", stats.TotalDirectories)
    fmt.Println("\nFiles by extension:")
    
    for ext, count := range stats.FilesByExtension {
        fmt.Printf("  %s: %d\n", ext, count)
    }
}
```

## Best Practices

### 1. Check File Existence

```go
// Good: Check before reading
data, err := file.Read("config.txt")
if err != nil {
    if os.IsNotExist(err) {
        // Create default config
        file.Write("config.txt", defaultConfig, 0644)
    } else {
        log.Fatal(err)
    }
}

// Or use os.Stat
if _, err := os.Stat("config.txt"); os.IsNotExist(err) {
    // File doesn't exist
}
```

### 2. Use Appropriate Permissions

```go
// Good: Restrictive permissions for sensitive data
file.Write(".env", secrets, 0600) // Owner only

// Good: Standard permissions for regular files
file.Write("data.json", jsonData, 0644) // Owner RW, others R

// Good: Standard permissions for directories
file.CreateDirectory("uploads", 0755) // Owner RWX, others RX

// Avoid: Overly permissive
// file.Write("secret.txt", password, 0777) // Too open!
```

### 3. Clean Up Resources

```go
// Good: Remove temporary files
tempFile := "temp_" + uuid.New().String() + ".txt"
file.Write(tempFile, data, 0644)
defer file.Remove(tempFile) // Ensure cleanup

// Good: Remove temporary directories
tempDir := "temp_processing"
file.CreateDirectory(tempDir, 0755)
defer file.RemoveAll(tempDir)
```

### 4. Handle Errors Properly

```go
// Good: Specific error handling
if err := file.Write("log.txt", data, 0644); err != nil {
    if os.IsPermission(err) {
        log.Println("Permission denied")
    } else if os.IsExist(err) {
        log.Println("File already exists")
    } else {
        log.Printf("Unknown error: %v", err)
    }
}

// Good: Wrap errors with context
if err := file.CreateDirectoryAll(path, 0755); err != nil {
    return fmt.Errorf("failed to create directory %s: %w", path, err)
}
```

### 5. Use Relative Paths Carefully

```go
// Good: Use absolute paths when needed
absPath, _ := filepath.Abs("data/config.json")
data, err := file.Read(absPath)

// Good: Be aware of working directory
// Current directory may change
cwd, _ := os.Getwd()
configPath := filepath.Join(cwd, "config", "app.json")
```

### 6. Validate Input Paths

```go
// Good: Sanitize user input
func SafeReadFile(userPath string) (string, error) {
    // Prevent path traversal
    cleanPath := filepath.Clean(userPath)
    if strings.Contains(cleanPath, "..") {
        return "", errors.New("invalid path")
    }
    
    return file.Read(cleanPath)
}
```

## Error Handling

### Common Errors

```go
// File not found
data, err := file.Read("missing.txt")
if os.IsNotExist(err) {
    log.Println("File does not exist")
}

// Permission denied
err = file.Write("/root/file.txt", "data", 0644)
if os.IsPermission(err) {
    log.Println("Permission denied")
}

// Directory not empty
err = file.Remove("data")
// Error: "directory not empty"

// Parent directory doesn't exist
err = file.CreateDirectory("a/b/c", 0755)
// Error: "no such file or directory"
```

### Error Recovery

```go
// Retry logic
func WriteWithRetry(filename, data string, maxRetries int) error {
    for i := 0; i < maxRetries; i++ {
        err := file.Write(filename, data, 0644)
        if err == nil {
            return nil
        }
        
        if os.IsPermission(err) {
            return err // Don't retry permission errors
        }
        
        time.Sleep(time.Second * time.Duration(i+1))
    }
    
    return fmt.Errorf("failed after %d retries", maxRetries)
}
```

## Performance Considerations

### Reading Large Files

```go
// Not ideal: Reads entire file into memory
data, err := file.Read("large.log") // May cause OOM

// Better: Use streaming for large files
f, _ := os.Open("large.log")
defer f.Close()
scanner := bufio.NewScanner(f)
for scanner.Scan() {
    line := scanner.Text()
    // Process line by line
}
```

### Batch Operations

```go
// Not ideal: Multiple separate operations
for i := 0; i < 1000; i++ {
    file.Write(fmt.Sprintf("file%d.txt", i), data, 0644)
}

// Better: Batch create directory structure first
file.CreateDirectoryAll("output", 0755)
for i := 0; i < 1000; i++ {
    file.Write(fmt.Sprintf("output/file%d.txt", i), data, 0644)
}
```

## Platform Differences

### Path Separators

```go
// Cross-platform: Use filepath.Separator
path := "data" + string(filepath.Separator) + "file.txt"

// Or use filepath.Join
path = filepath.Join("data", "file.txt")
```

### File Permissions (Windows vs Unix)

```go
// Unix/Linux: Full permission support
file.Write("file.txt", data, 0644) // rw-r--r--

// Windows: Limited permission support
// Only 0200 (read-only) and 0600 (read-write) are meaningful
file.Write("file.txt", data, 0644) // Works but limited effect
```

## Testing

### Unit Testing Example

```go
func TestConfigManager(t *testing.T) {
    // Setup
    testDir := "test_" + uuid.New().String()
    file.CreateDirectory(testDir, 0755)
    defer file.RemoveAll(testDir)
    
    configFile := filepath.Join(testDir, "config.json")
    
    // Test write
    testData := `{"key": "value"}`
    err := file.Write(configFile, testData, 0644)
    if err != nil {
        t.Fatalf("Write failed: %v", err)
    }
    
    // Test read
    data, err := file.Read(configFile)
    if err != nil {
        t.Fatalf("Read failed: %v", err)
    }
    
    if data != testData {
        t.Errorf("Expected %s, got %s", testData, data)
    }
}
```

## Dependencies

- `os` - Go standard library
- `path/filepath` - Go standard library

## Further Reading

- [os package documentation](https://pkg.go.dev/os)
- [filepath package documentation](https://pkg.go.dev/path/filepath)
- [File I/O in Go](https://gobyexample.com/reading-files)
- [Working with Files](https://go.dev/blog/gos-declaration-syntax)
