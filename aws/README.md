# AWS

AWS services integration for Go.

## Packages

This package collection provides simplified clients for AWS services:

- **[DynamoDB](dynamodb/)** - NoSQL database operations
- **[S3](s3/)** - Object storage operations

Each service has its own dedicated documentation with detailed usage examples and best practices.

## Installation

```bash
# Install DynamoDB client
go get -u github.com/common-library/go/aws/dynamodb

# Install S3 client
go get -u github.com/common-library/go/aws/s3
```

## Quick Start

### DynamoDB

```go
import "github.com/common-library/go/aws/dynamodb"

var client dynamodb.Client
err := client.CreateClient(ctx, "us-east-1", "key", "secret", "")

// Table and item operations
_, err = client.CreateTable(&dynamodb.CreateTableInput{...}, true, 30)
_, err = client.PutItem(&dynamodb.PutItemInput{...})
output, err := client.Query(&dynamodb.QueryInput{...})
```

See [DynamoDB documentation](dynamodb/) for complete examples.

### S3

```go
import "github.com/common-library/go/aws/s3"

var client s3.Client
err := client.CreateClient(ctx, "us-east-1", "key", "secret", "")

// Bucket and object operations
_, err = client.CreateBucket("my-bucket", "us-west-2")
_, err = client.PutObject("my-bucket", "key", "data")
output, err := client.GetObject("my-bucket", "key")
```

See [S3 documentation](s3/) for complete examples.

## Service Comparison

| Feature | [DynamoDB](dynamodb/) | [S3](s3/) |
|---------|----------|-----|
| Purpose | NoSQL Database | Object Storage |
| Data Model | Tables with items | Buckets with objects |
| Query Support | ✅ Primary key, indexes | ❌ Key-based only |
| Use Cases | Structured data, real-time apps | Files, backups, media |
| Local Testing | DynamoDB Local | MinIO, LocalStack |

## Common Features

All AWS service clients share:

- **AWS SDK v2 based** - Modern AWS SDK implementation
- **Context support** - Full context.Context integration
- **Custom endpoints** - Support for local testing (DynamoDB Local, MinIO, LocalStack)
- **Static credentials** - Simple credential configuration
- **Service-specific options** - Flexible configuration via functional options

## Error Handling

All functions return standard Go errors:

```go
output, err := client.GetItem(&dynamodb.GetItemInput{...})
if err != nil {
    var nfe *types.ResourceNotFoundException
    if errors.As(err, &nfe) {
        // Handle resource not found
    }
    return err
}
```

Common error scenarios:
- Invalid credentials
- Resource not found
- Insufficient permissions
- Network connectivity issues
- Service throttling

## Local Testing

### DynamoDB Local

```bash
docker run -p 8000:8000 amazon/dynamodb-local
```

```go
client.CreateClient(ctx, "us-east-1", "dummy", "dummy", "",
    func(o *dynamodb.Options) {
        o.BaseEndpoint = aws.String("http://localhost:8000")
    },
)
```

### MinIO (S3-compatible)

```bash
docker run -p 9000:9000 minio/minio server /data
```

```go
client.CreateClient(ctx, "us-east-1", "minioadmin", "minioadmin", "",
    func(o *s3.Options) {
        o.BaseEndpoint = aws.String("http://localhost:9000")
        o.UsePathStyle = true
    },
)
```

### LocalStack

```bash
docker run -p 4566:4566 localstack/localstack
```

Supports both DynamoDB and S3 with endpoint `http://localhost:4566`.

## Best Practices

1. **Use context with timeout**
   ```go
   ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
   defer cancel()
   ```

2. **Close resources** (S3)
   ```go
   output, err := client.GetObject("bucket", "key")
   defer output.Body.Close()
   ```

3. **Handle errors appropriately**
   ```go
   if err != nil {
       // Type assertion for specific error handling
       var throttleErr *types.ProvisionedThroughputExceededException
       if errors.As(err, &throttleErr) {
           // Implement retry logic
       }
   }
   ```

4. **Never hardcode credentials**
   - Use environment variables
   - Use AWS IAM roles
   - Rotate credentials regularly

## Dependencies

- `github.com/aws/aws-sdk-go-v2/aws` - Core AWS SDK
- `github.com/aws/aws-sdk-go-v2/config` - Configuration loading
- `github.com/aws/aws-sdk-go-v2/credentials` - Credentials management
- `github.com/aws/aws-sdk-go-v2/service/dynamodb` - DynamoDB service
- `github.com/aws/aws-sdk-go-v2/service/s3` - S3 service

## Further Reading

- [DynamoDB Package Documentation](dynamodb/)
- [S3 Package Documentation](s3/)
- [AWS SDK for Go v2](https://aws.github.io/aws-sdk-go-v2/)
- [AWS Documentation](https://docs.aws.amazon.com/)
