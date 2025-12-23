# MinIO Storage Client

S3-compatible object storage client for MinIO and AWS S3.

## Overview

The minio package provides a high-level interface for MinIO and S3-compatible object storage systems. It wraps the MinIO Go SDK with simplified methods for common operations including bucket management, object upload/download, and metadata queries.

## Features

- **Bucket Management** - Create, list, check existence, remove buckets
- **Object Upload** - Upload from files, readers, or streams
- **Object Download** - Download to files or readers
- **Object Operations** - Copy, delete, stat objects
- **Bulk Operations** - Remove multiple objects efficiently
- **Metadata Access** - Get object size, modified time, content type
- **S3 Compatibility** - Works with MinIO, AWS S3, and S3-compatible services

## Installation

```bash
go get -u github.com/common-library/go/storage/minio
```

## Quick Start

```go
import "github.com/common-library/go/storage/minio"

client := &minio.Client{}

// Connect to MinIO server
err := client.CreateClient("localhost:9000", "minioadmin", "minioadmin", false)

// Create bucket
err = client.MakeBucket("mybucket", "us-east-1", false)

// Upload file
err = client.FPutObject("mybucket", "document.pdf", "/path/to/file.pdf", "application/pdf")

// Download file
err = client.FGetObject("mybucket", "document.pdf", "/tmp/download.pdf")

// List objects
objects, err := client.ListObjects("mybucket", "", true)
for _, obj := range objects {
    fmt.Printf("%s - %d bytes\n", obj.Key, obj.Size)
}
```

## API Reference

### Client Type

```go
type Client struct {
    client *minio.Client
}
```

MinIO client wrapper with simplified operations.

### Client Methods

#### CreateClient

```go
func (c *Client) CreateClient(endpoint, accessKeyID, secretAccessKey string, secure bool) error
```

Initializes the MinIO client connection.

**Parameters:**
- `endpoint` - Server address (e.g., "localhost:9000", "s3.amazonaws.com")
- `accessKeyID` - Access key for authentication
- `secretAccessKey` - Secret key for authentication
- `secure` - Use HTTPS if true, HTTP if false

#### MakeBucket

```go
func (c *Client) MakeBucket(bucketName, region string, objectLocking bool) error
```

Creates a new bucket.

#### ListBuckets

```go
func (c *Client) ListBuckets() ([]minio.BucketInfo, error)
```

Returns all buckets owned by the user.

#### BucketExists

```go
func (c *Client) BucketExists(bucketName string) (bool, error)
```

Checks if a bucket exists.

#### RemoveBucket

```go
func (c *Client) RemoveBucket(bucketName string) error
```

Deletes a bucket (must be empty).

#### ListObjects

```go
func (c *Client) ListObjects(bucketName, prefix string, recursive bool) ([]minio.ObjectInfo, error)
```

Lists objects in a bucket with optional prefix filtering.

#### GetObject

```go
func (c *Client) GetObject(bucketName, objectName string) (*minio.Object, error)
```

Retrieves an object as a reader.

#### PutObject

```go
func (c *Client) PutObject(bucketName, objectName, contentType string, reader io.Reader, objectSize int64) error
```

Uploads an object from a reader.

#### CopyObject

```go
func (c *Client) CopyObject(sourceBucketName, sourceObjectName, destinationBucketName, destinationObjectName string) error
```

Copies an object within or between buckets.

#### StatObject

```go
func (c *Client) StatObject(bucketName, objectName string) (minio.ObjectInfo, error)
```

Retrieves object metadata.

#### RemoveObject

```go
func (c *Client) RemoveObject(bucketName, objectName string, forceDelete bool, governanceBypass bool, versionID string) error
```

Deletes a single object.

#### RemoveObjects

```go
func (c *Client) RemoveObjects(bucketName string, objectInfos []minio.ObjectInfo, governanceBypass bool) []minio.RemoveObjectError
```

Deletes multiple objects in a single operation.

#### FPutObject

```go
func (c *Client) FPutObject(bucketName, objectName, filePath, contentType string) error
```

Uploads a file to a bucket.

#### FGetObject

```go
func (c *Client) FGetObject(bucketName, objectName, filePath string) error
```

Downloads an object to a file.

## Complete Examples

### Connecting to Different Services

```go
package main

import (
    "log"
    "github.com/common-library/go/storage/minio"
)

func main() {
    client := &minio.Client{}
    
    // Local MinIO
    err := client.CreateClient("localhost:9000", "minioadmin", "minioadmin", false)
    
    // AWS S3
    // err := client.CreateClient("s3.amazonaws.com", "ACCESS_KEY", "SECRET_KEY", true)
    
    // DigitalOcean Spaces
    // err := client.CreateClient("nyc3.digitaloceanspaces.com", "KEY", "SECRET", true)
    
    if err != nil {
        log.Fatal(err)
    }
}
```

### Bucket Operations

```go
package main

import (
    "fmt"
    "log"
    "github.com/common-library/go/storage/minio"
)

func main() {
    client := &minio.Client{}
    client.CreateClient("localhost:9000", "minioadmin", "minioadmin", false)
    
    // Create bucket
    err := client.MakeBucket("photos", "us-east-1", false)
    if err != nil {
        log.Printf("Make bucket error: %v", err)
    }
    
    // Check if bucket exists
    exists, err := client.BucketExists("photos")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Bucket exists: %v\n", exists)
    
    // List all buckets
    buckets, err := client.ListBuckets()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Buckets:")
    for _, bucket := range buckets {
        fmt.Printf("  - %s (created: %v)\n", bucket.Name, bucket.CreationDate)
    }
}
```

### File Upload and Download

```go
package main

import (
    "log"
    "github.com/common-library/go/storage/minio"
)

func main() {
    client := &minio.Client{}
    client.CreateClient("localhost:9000", "minioadmin", "minioadmin", false)
    
    // Upload file
    err := client.FPutObject(
        "photos",
        "vacation/beach.jpg",
        "/home/user/pictures/beach.jpg",
        "image/jpeg",
    )
    if err != nil {
        log.Fatalf("Upload failed: %v", err)
    }
    log.Println("Upload successful")
    
    // Download file
    err = client.FGetObject(
        "photos",
        "vacation/beach.jpg",
        "/tmp/beach.jpg",
    )
    if err != nil {
        log.Fatalf("Download failed: %v", err)
    }
    log.Println("Download successful")
}
```

### Stream Upload/Download

```go
package main

import (
    "bytes"
    "io"
    "log"
    "github.com/common-library/go/storage/minio"
)

func main() {
    client := &minio.Client{}
    client.CreateClient("localhost:9000", "minioadmin", "minioadmin", false)
    
    // Upload from buffer
    data := []byte("Hello, MinIO!")
    reader := bytes.NewReader(data)
    
    err := client.PutObject(
        "mybucket",
        "hello.txt",
        "text/plain",
        reader,
        int64(len(data)),
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Download to buffer
    object, err := client.GetObject("mybucket", "hello.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer object.Close()
    
    content, err := io.ReadAll(object)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Content: %s", content)
}
```

### Listing and Filtering Objects

```go
package main

import (
    "fmt"
    "log"
    "github.com/common-library/go/storage/minio"
)

func main() {
    client := &minio.Client{}
    client.CreateClient("localhost:9000", "minioadmin", "minioadmin", false)
    
    // List all objects recursively
    objects, err := client.ListObjects("photos", "", true)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("All objects:")
    for _, obj := range objects {
        fmt.Printf("  %s - %d bytes (modified: %v)\n",
            obj.Key, obj.Size, obj.LastModified)
    }
    
    // List objects with prefix
    vacation, err := client.ListObjects("photos", "vacation/", true)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("\nVacation photos: %d\n", len(vacation))
}
```

### Object Metadata

```go
package main

import (
    "fmt"
    "log"
    "github.com/common-library/go/storage/minio"
)

func main() {
    client := &minio.Client{}
    client.CreateClient("localhost:9000", "minioadmin", "minioadmin", false)
    
    // Get object metadata
    info, err := client.StatObject("photos", "vacation/beach.jpg")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Object: %s\n", info.Key)
    fmt.Printf("Size: %d bytes\n", info.Size)
    fmt.Printf("Content-Type: %s\n", info.ContentType)
    fmt.Printf("Last Modified: %v\n", info.LastModified)
    fmt.Printf("ETag: %s\n", info.ETag)
}
```

### Copying Objects

```go
package main

import (
    "log"
    "github.com/common-library/go/storage/minio"
)

func main() {
    client := &minio.Client{}
    client.CreateClient("localhost:9000", "minioadmin", "minioadmin", false)
    
    // Copy within same bucket
    err := client.CopyObject(
        "photos", "original.jpg",
        "photos", "backup/original.jpg",
    )
    if err != nil {
        log.Fatal(err)
    }
    log.Println("Copied within bucket")
    
    // Copy to different bucket
    err = client.CopyObject(
        "photos", "important.jpg",
        "backup-bucket", "important.jpg",
    )
    if err != nil {
        log.Fatal(err)
    }
    log.Println("Copied to backup bucket")
}
```

### Deleting Objects

```go
package main

import (
    "log"
    "github.com/common-library/go/storage/minio"
)

func main() {
    client := &minio.Client{}
    client.CreateClient("localhost:9000", "minioadmin", "minioadmin", false)
    
    // Delete single object
    err := client.RemoveObject("photos", "old-photo.jpg", false, false, "")
    if err != nil {
        log.Fatal(err)
    }
    log.Println("Object deleted")
    
    // Bulk delete
    objects, _ := client.ListObjects("temp-bucket", "", true)
    errors := client.RemoveObjects("temp-bucket", objects, false)
    
    if len(errors) > 0 {
        log.Printf("Some deletions failed:")
        for _, e := range errors {
            log.Printf("  - %s: %v", e.ObjectName, e.Err)
        }
    } else {
        log.Printf("Deleted %d objects", len(objects))
    }
}
```

### Photo Gallery Manager

```go
package main

import (
    "fmt"
    "log"
    "path/filepath"
    "strings"
    "github.com/common-library/go/storage/minio"
)

type PhotoGallery struct {
    client *minio.Client
    bucket string
}

func NewPhotoGallery(endpoint, accessKey, secretKey, bucket string) (*PhotoGallery, error) {
    client := &minio.Client{}
    err := client.CreateClient(endpoint, accessKey, secretKey, false)
    if err != nil {
        return nil, err
    }
    
    // Create bucket if not exists
    exists, _ := client.BucketExists(bucket)
    if !exists {
        err = client.MakeBucket(bucket, "", false)
        if err != nil {
            return nil, err
        }
    }
    
    return &PhotoGallery{client: client, bucket: bucket}, nil
}

func (g *PhotoGallery) UploadPhoto(localPath, album string) error {
    filename := filepath.Base(localPath)
    objectName := fmt.Sprintf("%s/%s", album, filename)
    
    contentType := "image/jpeg"
    if strings.HasSuffix(localPath, ".png") {
        contentType = "image/png"
    }
    
    return g.client.FPutObject(g.bucket, objectName, localPath, contentType)
}

func (g *PhotoGallery) ListAlbums() ([]string, error) {
    objects, err := g.client.ListObjects(g.bucket, "", false)
    if err != nil {
        return nil, err
    }
    
    albums := make(map[string]bool)
    for _, obj := range objects {
        parts := strings.Split(obj.Key, "/")
        if len(parts) > 1 {
            albums[parts[0]] = true
        }
    }
    
    result := make([]string, 0, len(albums))
    for album := range albums {
        result = append(result, album)
    }
    return result, nil
}

func (g *PhotoGallery) GetPhotosInAlbum(album string) ([]string, error) {
    objects, err := g.client.ListObjects(g.bucket, album+"/", false)
    if err != nil {
        return nil, err
    }
    
    photos := make([]string, 0, len(objects))
    for _, obj := range objects {
        photos = append(photos, obj.Key)
    }
    return photos, nil
}

func main() {
    gallery, err := NewPhotoGallery("localhost:9000", "minioadmin", "minioadmin", "photos")
    if err != nil {
        log.Fatal(err)
    }
    
    // Upload photos
    gallery.UploadPhoto("/path/to/photo1.jpg", "vacation")
    gallery.UploadPhoto("/path/to/photo2.jpg", "vacation")
    
    // List albums
    albums, _ := gallery.ListAlbums()
    fmt.Printf("Albums: %v\n", albums)
    
    // List photos in album
    photos, _ := gallery.GetPhotosInAlbum("vacation")
    fmt.Printf("Vacation photos: %v\n", photos)
}
```

## Best Practices

### 1. Initialize Client Once

```go
// Good: Reuse client
var globalClient *minio.Client

func init() {
    globalClient = &minio.Client{}
    globalClient.CreateClient("localhost:9000", "key", "secret", false)
}

// Avoid: Creating new client for each request
func uploadFile() {
    client := &minio.Client{}
    client.CreateClient(...) // Wasteful
}
```

### 2. Handle Errors Properly

```go
// Good: Check all errors
err := client.FPutObject("bucket", "key", "file.txt", "text/plain")
if err != nil {
    return fmt.Errorf("upload failed: %w", err)
}

// Avoid: Ignore errors
client.FPutObject("bucket", "key", "file.txt", "text/plain")
```

### 3. Close Object Readers

```go
// Good: Always close
object, err := client.GetObject("bucket", "key")
if err != nil {
    return err
}
defer object.Close()

// Avoid: Forget to close
object, _ := client.GetObject("bucket", "key")
// Resource leak!
```

### 4. Use Meaningful Object Keys

```go
// Good: Hierarchical naming
"users/1234/profile.jpg"
"documents/2024/01/report.pdf"
"logs/2024-01-15/app.log"

// Avoid: Flat naming
"file1.jpg"
"doc.pdf"
```

### 5. Set Correct Content Types

```go
// Good: Accurate content types
client.FPutObject("bucket", "image.jpg", path, "image/jpeg")
client.FPutObject("bucket", "doc.pdf", path, "application/pdf")

// Avoid: Generic or wrong types
client.FPutObject("bucket", "image.jpg", path, "application/octet-stream")
```

## Performance Tips

1. **Concurrent Uploads** - Upload multiple files in parallel using goroutines
2. **Streaming** - Use PutObject with readers for large files to avoid memory issues
3. **Bulk Operations** - Use RemoveObjects instead of multiple RemoveObject calls
4. **Connection Pooling** - Reuse the same client instance across requests
5. **Regional Buckets** - Create buckets in regions close to your application

## Error Handling

```go
// Check for specific errors
exists, err := client.BucketExists("mybucket")
if err != nil {
    log.Printf("Error checking bucket: %v", err)
}

// Bulk operation errors
errors := client.RemoveObjects("bucket", objects, false)
for _, e := range errors {
    log.Printf("Failed to delete %s: %v", e.ObjectName, e.Err)
}
```

## S3 Compatibility

This client works with:
- **MinIO** - Self-hosted object storage
- **AWS S3** - Amazon Simple Storage Service
- **Google Cloud Storage** - S3-compatible API
- **DigitalOcean Spaces** - S3-compatible object storage
- **Wasabi** - Cloud object storage
- **Backblaze B2** - S3-compatible API
- **Any S3-compatible service**

## Testing

```go
func TestMinIO(t *testing.T) {
    client := &minio.Client{}
    err := client.CreateClient("localhost:9000", "minioadmin", "minioadmin", false)
    if err != nil {
        t.Fatal(err)
    }
    
    // Create test bucket
    testBucket := "test-bucket"
    err = client.MakeBucket(testBucket, "", false)
    if err != nil {
        t.Fatal(err)
    }
    defer client.RemoveBucket(testBucket)
    
    // Test upload
    err = client.FPutObject(testBucket, "test.txt", "testfile.txt", "text/plain")
    if err != nil {
        t.Errorf("Upload failed: %v", err)
    }
    
    // Test download
    err = client.FGetObject(testBucket, "test.txt", "/tmp/test.txt")
    if err != nil {
        t.Errorf("Download failed: %v", err)
    }
}
```

## Dependencies

- `github.com/minio/minio-go/v7` - MinIO Go SDK
- `context` - Go standard library

## Further Reading

- [MinIO Documentation](https://min.io/docs/minio/linux/index.html)
- [MinIO Go Client API Reference](https://min.io/docs/minio/linux/developers/go/API.html)
- [AWS S3 API Documentation](https://docs.aws.amazon.com/s3/)
- [S3 Best Practices](https://docs.aws.amazon.com/AmazonS3/latest/userguide/optimizing-performance.html)
