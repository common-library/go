# S3

AWS S3 (Simple Storage Service) client utilities for Go.

Part of the [AWS package collection](../).

## Features

- Bucket management (create, list, delete)
- Object operations (put, get, delete)
- Custom endpoint support (MinIO, LocalStack compatible)
- Path-style and virtual-hosted-style access

## Installation

```bash
go get -u github.com/common-library/go/aws/s3
```

## Usage

### Basic Setup

```go
import (
    "context"
    "github.com/common-library/go/aws/s3"
    "github.com/aws/aws-sdk-go-v2/aws"
)

// Initialize client
var client s3.Client
err := client.CreateClient(
    context.TODO(),
    "us-east-1",
    "your-access-key",
    "your-secret-key",
    "",
)
```

### Bucket Operations

#### Create Bucket

```go
// Create bucket in default region
_, err = client.CreateBucket("my-bucket", "")

// Create bucket in specific region
_, err = client.CreateBucket("my-bucket", "us-west-2")
```

#### List Buckets

```go
output, err := client.ListBuckets()
for _, bucket := range output.Buckets {
    fmt.Printf("Bucket: %s, Created: %s\n", 
        *bucket.Name, bucket.CreationDate)
}
```

#### Delete Bucket

```go
// Bucket must be empty before deletion
_, err = client.DeleteBucket("my-bucket")
```

### Object Operations

#### Put Object (Upload)

```go
// Upload text content
_, err := client.PutObject("my-bucket", "docs/file.txt", "Hello, World!")

// Upload JSON
jsonData := `{"name": "John", "age": 30}`
_, err = client.PutObject("my-bucket", "data/user.json", jsonData)
```

#### Get Object (Download)

```go
output, err := client.GetObject("my-bucket", "docs/file.txt")
if err != nil {
    log.Fatal(err)
}
defer output.Body.Close()

content, err := io.ReadAll(output.Body)
fmt.Printf("Content: %s\n", string(content))
```

#### Delete Object

```go
_, err = client.DeleteObject("my-bucket", "docs/file.txt")
```

## API Reference

### Client Management
- `CreateClient(ctx, region, accessKey, secretKey, sessionToken, options...)` - Initialize client

### Bucket Operations
- `CreateBucket(name, region)` - Create new bucket
- `ListBuckets()` - List all buckets
- `DeleteBucket(name)` - Delete empty bucket

### Object Operations
- `PutObject(bucketName, key, body)` - Upload object
- `GetObject(bucketName, key)` - Download object
- `DeleteObject(bucketName, key)` - Delete object

## Custom Endpoint Configuration

### MinIO (S3-compatible)

For local development with MinIO:

```go
var client s3.Client
err := client.CreateClient(
    context.TODO(),
    "us-east-1",
    "minioadmin",
    "minioadmin",
    "",
    func(o *aws_s3.Options) {
        o.BaseEndpoint = aws.String("http://localhost:9000")
        o.UsePathStyle = true  // Required for MinIO
    },
)
```

Run MinIO with Docker:

```bash
docker run -p 9000:9000 -p 9001:9001 \
  -e "MINIO_ROOT_USER=minioadmin" \
  -e "MINIO_ROOT_PASSWORD=minioadmin" \
  minio/minio server /data --console-address ":9001"
```

### LocalStack

For local AWS service emulation with LocalStack:

```go
var client s3.Client
err := client.CreateClient(
    context.TODO(),
    "us-east-1",
    "test",
    "test",
    "",
    func(o *aws_s3.Options) {
        o.BaseEndpoint = aws.String("http://localhost:4566")
        o.UsePathStyle = true
    },
)
```

Run LocalStack with Docker:

```bash
docker run -p 4566:4566 localstack/localstack
```

## Implementation Details

- Uses `github.com/aws/aws-sdk-go-v2/service/s3`
- AWS SDK v2 based implementation
- Context support for all operations
- Path-style and virtual-hosted-style URL support
- Compatible with S3-compatible services (MinIO, LocalStack)
- Simple string-based object upload (PutObject)
- Stream-based download (GetObject returns io.ReadCloser)

## Error Handling

All functions return errors that should be checked:

```go
output, err := client.GetObject("bucket", "key")
if err != nil {
    log.Fatalf("Failed to get object: %v", err)
}
defer output.Body.Close()  // Always close response body
```

Common error scenarios:
- Invalid credentials (authentication failure)
- Bucket does not exist
- Object not found
- Insufficient permissions
- Network connectivity issues
- Service throttling (rate limits)
- Bucket name conflicts (must be globally unique)

## Best Practices

1. **Always Close Response Body**
   ```go
   output, err := client.GetObject("bucket", "key")
   if err != nil {
       return err
   }
   defer output.Body.Close()  // Critical for resource cleanup
   ```

2. **Use Context with Timeout**
   ```go
   ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
   defer cancel()
   
   err := client.CreateClient(ctx, region, accessKey, secretKey, "")
   ```

3. **Handle Large Files**
   ```go
   // For large files, consider using multipart upload
   // Or stream directly without loading into memory
   output, err := client.GetObject("bucket", "large-file.zip")
   if err != nil {
       return err
   }
   defer output.Body.Close()
   
   // Stream to file
   file, err := os.Create("local-file.zip")
   if err != nil {
       return err
   }
   defer file.Close()
   
   _, err = io.Copy(file, output.Body)
   ```

4. **Bucket Naming**
   ```go
   // Bucket names must be globally unique
   // Use DNS-compliant names
   // 3-63 characters, lowercase, numbers, hyphens
   bucketName := "my-company-app-data-2024"
   ```

5. **Error Handling**
   ```go
   if err != nil {
       var noBucket *types.NoSuchBucket
       if errors.As(err, &noBucket) {
           // Handle bucket not found
           log.Printf("Bucket does not exist")
       }
       return err
   }
   ```

6. **Credential Security**
   - Never hardcode credentials in source code
   - Use environment variables or AWS IAM roles
   - Rotate credentials regularly
   - Use session tokens for temporary access

## Examples

### Complete Workflow

```go
package main

import (
    "context"
    "fmt"
    "io"
    "log"
    
    "github.com/common-library/go/aws/s3"
)

func main() {
    var client s3.Client
    
    // Initialize client
    err := client.CreateClient(
        context.TODO(),
        "us-east-1",
        "your-access-key",
        "your-secret-key",
        "",
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Create bucket
    bucketName := "my-test-bucket-12345"
    _, err = client.CreateBucket(bucketName, "us-west-2")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Created bucket: %s\n", bucketName)
    
    // Upload object
    objectKey := "docs/readme.txt"
    content := "This is a test file"
    _, err = client.PutObject(bucketName, objectKey, content)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Uploaded object: %s\n", objectKey)
    
    // Download object
    output, err := client.GetObject(bucketName, objectKey)
    if err != nil {
        log.Fatal(err)
    }
    defer output.Body.Close()
    
    data, err := io.ReadAll(output.Body)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Downloaded content: %s\n", string(data))
    
    // Delete object
    _, err = client.DeleteObject(bucketName, objectKey)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Deleted object: %s\n", objectKey)
    
    // Delete bucket
    _, err = client.DeleteBucket(bucketName)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Deleted bucket: %s\n", bucketName)
}
```

### File Upload and Download

```go
// Upload file from disk
func uploadFile(client *s3.Client, bucketName, objectKey, filePath string) error {
    data, err := os.ReadFile(filePath)
    if err != nil {
        return err
    }
    
    _, err = client.PutObject(bucketName, objectKey, string(data))
    return err
}

// Download file to disk
func downloadFile(client *s3.Client, bucketName, objectKey, filePath string) error {
    output, err := client.GetObject(bucketName, objectKey)
    if err != nil {
        return err
    }
    defer output.Body.Close()
    
    file, err := os.Create(filePath)
    if err != nil {
        return err
    }
    defer file.Close()
    
    _, err = io.Copy(file, output.Body)
    return err
}
```

### List and Process Objects

```go
// Note: Current client doesn't have ListObjects
// This is an example of what you might implement
// using the AWS SDK directly with the client's credentials

import (
    aws_s3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

func listObjects(bucketName string) error {
    // You would need to access the underlying AWS S3 client
    // or add a ListObjects method to the Client struct
    
    // Example with direct SDK usage:
    // output, err := awsClient.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
    //     Bucket: aws.String(bucketName),
    // })
    
    return nil
}
```

## Path-Style vs Virtual-Hosted-Style

### Virtual-Hosted-Style (Default)
```
https://bucket-name.s3.region.amazonaws.com/object-key
```

### Path-Style (Required for MinIO)
```
https://s3.region.amazonaws.com/bucket-name/object-key
```

Enable path-style in client creation:
```go
func(o *aws_s3.Options) {
    o.UsePathStyle = true
}
```

## Dependencies

- `github.com/aws/aws-sdk-go-v2/aws` - Core AWS SDK
- `github.com/aws/aws-sdk-go-v2/config` - Configuration loading
- `github.com/aws/aws-sdk-go-v2/credentials` - Credentials management
- `github.com/aws/aws-sdk-go-v2/service/s3` - S3 service

## Further Reading

- [AWS SDK for Go v2 Documentation](https://aws.github.io/aws-sdk-go-v2/)
- [Amazon S3 Developer Guide](https://docs.aws.amazon.com/s3/)
- [S3 Best Practices](https://docs.aws.amazon.com/AmazonS3/latest/userguide/best-practices.html)
- [MinIO Documentation](https://min.io/docs/minio/linux/index.html)
- [LocalStack Documentation](https://docs.localstack.cloud/)
- [AWS package overview](../)
