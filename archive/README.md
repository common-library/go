# Archive

Utility package collection for file compression and decompression.

## Features

### Gzip
- Single file compression/decompression
- `.gz` format support
- Based on standard `compress/gzip` package

### Tar
- Multiple files/directories compression/decompression
- `.tar.gz` (gzip-compressed tar) format support
- Preserves directory structure
- Maintains file permissions and metadata

### Zip
- Multiple files/directories compression/decompression
- `.zip` format support
- Preserves directory structure
- Cross-platform compatibility

## Installation

```bash
go get -u github.com/common-library/go/archive/gzip
go get -u github.com/common-library/go/archive/tar
go get -u github.com/common-library/go/archive/zip
```

## Usage

### Gzip

```go
import "github.com/common-library/go/archive/gzip"

// Compress
err := gzip.Compress("output.gz", "input.txt")

// Decompress
err := gzip.Decompress("input.gz", "output.txt", "./output-directory")
```

**Key Functions:**
- `Compress(name, path)` - Compress a single file with gzip
- `Decompress(gzipName, fileName, outputPath)` - Decompress gzip file

### Tar

```go
import "github.com/common-library/go/archive/tar"

// Compress (multiple files/directories)
err := tar.Compress("archive.tar.gz", []string{"./dir1", "./file.txt"})

// Decompress
err := tar.Decompress("archive.tar.gz", "./output-directory")
```

**Key Functions:**
- `Compress(name, paths)` - Compress multiple files/directories to tar.gz
- `Decompress(name, outputPath)` - Decompress tar.gz file
- Recursive directory traversal
- Preserves file permissions

### Zip

```go
import "github.com/common-library/go/archive/zip"

// Compress (multiple files/directories)
err := zip.Compress("archive.zip", []string{"./dir1", "./file.txt"})

// Decompress
err := zip.Decompress("archive.zip", "./output-directory")
```

**Key Functions:**
- `Compress(name, paths)` - Compress multiple files/directories to zip
- `Decompress(name, outputPath)` - Decompress zip file
- Recursive directory traversal
- Preserves file mode

## Key Differences

| Feature | Gzip | Tar | Zip |
|---------|------|-----|-----|
| File Count | Single | Multiple | Multiple |
| Directory Support | ❌ | ✅ | ✅ |
| Compression Ratio | High | High (uses gzip) | Medium |
| Permission Preservation | ❌ | ✅ | ✅ |
| Cross-platform | ✅ | ✅ | ✅ |
| Use Cases | Log compression | Backup, deployment | Deployment, sharing |

## Implementation Details

### Common Features
- Automatic output directory creation
- Error handling
- Automatic resource cleanup (defer)
- Dependency: `github.com/common-library/go/file`

### Gzip Characteristics
- Uses `compress/gzip` standard library
- Single file stream processing
- Maximum compression ratio

### Tar Characteristics
- Combination of `archive/tar` + `compress/gzip`
- Preserves file metadata (FileInfoHeader)
- Distinguishes directory/file types (TypeReg, TypeDir)
- Recursive file list generation

### Zip Characteristics
- Uses `archive/zip` standard library
- Preserves file mode
- Excellent Windows compatibility

## Error Handling

All functions return `error`, which occurs in the following cases:
- Input file/directory not found
- Output path creation failure
- Insufficient read/write permissions
- Corrupted archive file
- Insufficient disk space

## Examples

### Log File Compression
```go
// Compress single log file
gzip.Compress("app.log.gz", "app.log")
```

### Project Backup
```go
// Backup entire project to tar.gz
tar.Compress("backup.tar.gz", []string{"./src", "./config", "./data"})
```

### Release Package Creation
```go
// Create zip file for deployment
zip.Compress("release.zip", []string{"./bin", "./config", "./README.md"})
```

### Decompression
```go
// Restore tar.gz backup
tar.Decompress("backup.tar.gz", "./restore")

// Extract zip package
zip.Decompress("release.zip", "./deploy")

// Restore gzip log
gzip.Decompress("app.log.gz", "app.log", "./logs")
```

## Best Practices

1. **Choose Appropriate Format**
   - Single file → Gzip
   - Backup/deployment (Unix) → Tar
   - Cross-platform sharing → Zip

2. **Error Handling**
   ```go
   if err := tar.Compress("backup.tar.gz", paths); err != nil {
       log.Fatalf("Compression failed: %v", err)
   }
   ```

3. **Path Validation**
   - Verify file/directory existence before compression
   - Check write permissions for output path

4. **Resource Management**
   - All functions automatically clean up resources with defer
   - No need for explicit Close() calls

## Dependencies

- `github.com/common-library/go/file` - File/directory utilities
- `compress/gzip` - Gzip compression
- `archive/tar` - Tar archive
- `archive/zip` - Zip archive
