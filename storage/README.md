# Storage

Object storage and file management utilities.

## Overview

The storage package provides interfaces for object storage systems, currently supporting MinIO for S3-compatible object storage operations.

## Subpackages

### minio

MinIO client for S3-compatible object storage.

[📖 Documentation](minio/README.md)

**Features:**
- Bucket management (create, list, delete)
- Object operations (upload, download, copy, delete)
- File-based operations (FPutObject, FGetObject)
- Bulk object removal
- Object metadata and statistics

**Quick Example:**
```go
import "github.com/common-library/go/storage/minio"

// Create client
client := &minio.Client{}
err := client.CreateClient("localhost:9000", "accessKey", "secretKey", false)

// Create bucket
err = client.MakeBucket("mybucket", "us-east-1", false)

// Upload file
err = client.FPutObject("mybucket", "document.pdf", "/path/to/file.pdf", "application/pdf")

// Download file
err = client.FGetObject("mybucket", "document.pdf", "/path/to/download.pdf")

// List objects
objects, err := client.ListObjects("mybucket", "documents/", true)
```

## Use Cases

- **File Storage** - Store and retrieve user uploads
- **Backup Systems** - Backup application data to object storage
- **Media Management** - Store images, videos, and documents
- **Data Lake** - Build data lakes for analytics
- **Static Content** - Serve static website assets

## Installation

```bash
go get -u github.com/common-library/go/storage/minio
```

## Common Operations

### Bucket Management

```go
client := &minio.Client{}
client.CreateClient("s3.amazonaws.com", "ACCESS_KEY", "SECRET_KEY", true)

// Create bucket
err := client.MakeBucket("photos", "us-east-1", false)

// Check if bucket exists
exists, err := client.BucketExists("photos")

// List all buckets
buckets, err := client.ListBuckets()

// Remove bucket
err = client.RemoveBucket("photos")
```

### Object Upload

```go
// Upload from file
err := client.FPutObject("photos", "vacation.jpg", "/home/user/vacation.jpg", "image/jpeg")

// Upload from reader
file, _ := os.Open("/path/to/file.pdf")
defer file.Close()

stat, _ := file.Stat()
err = client.PutObject("documents", "file.pdf", "application/pdf", file, stat.Size())
```

### Object Download

```go
// Download to file
err := client.FGetObject("photos", "vacation.jpg", "/tmp/vacation.jpg")

// Download to reader
object, err := client.GetObject("documents", "file.pdf")
defer object.Close()

data, err := io.ReadAll(object)
```

### Object Management

```go
// List objects in bucket
objects, err := client.ListObjects("photos", "2024/", true)
for _, obj := range objects {
    fmt.Printf("%s - %d bytes\n", obj.Key, obj.Size)
}

// Get object metadata
info, err := client.StatObject("photos", "vacation.jpg")
fmt.Printf("Size: %d, Modified: %v\n", info.Size, info.LastModified)

// Copy object
err = client.CopyObject("photos", "original.jpg", "backup", "original.jpg")

// Delete object
err = client.RemoveObject("photos", "old-photo.jpg", false, false, "")

// Bulk delete
objectsToDelete := []minio.ObjectInfo{...}
errors := client.RemoveObjects("photos", objectsToDelete, false)
```

## Best Practices

1. **Connection Pooling** - Reuse client instances across requests
2. **Error Handling** - Check all error returns from operations
3. **Bucket Naming** - Follow S3 bucket naming conventions (lowercase, no special chars)
4. **Object Keys** - Use meaningful, hierarchical object keys (e.g., "2024/01/photo.jpg")
5. **Cleanup** - Always close readers from GetObject
6. **Versioning** - Enable bucket versioning for critical data
7. **Lifecycle Policies** - Configure automatic cleanup of old objects

## Performance Tips

- **Multipart Upload** - Use for large files (>5MB)
- **Concurrent Operations** - Upload/download multiple objects in parallel
- **Content Type** - Set accurate content types for proper handling
- **Compression** - Compress data before upload when appropriate
- **Regional Buckets** - Create buckets close to your application

## S3 Compatibility

The MinIO client is compatible with:
- MinIO Server
- Amazon S3
- Google Cloud Storage (S3 API)
- DigitalOcean Spaces
- Wasabi
- Backblaze B2 (S3 API)
- Any S3-compatible storage

## Dependencies

- `github.com/minio/minio-go/v7` - MinIO Go SDK

## Further Reading

- [MinIO Package Documentation](minio/README.md)
- [MinIO Server Documentation](https://min.io/docs/minio/linux/index.html)
- [S3 API Reference](https://docs.aws.amazon.com/AmazonS3/latest/API/Welcome.html)
