# DynamoDB

AWS DynamoDB client utilities for Go.

Part of the [AWS package collection](../).

## Features

- Table management (create, list, describe, update, delete)
- Item operations (get, put, update, delete)
- Query and scan operations
- Pagination support
- TTL (Time To Live) management
- Table creation/deletion waiter support

## Installation

```bash
go get -u github.com/common-library/go/aws/dynamodb
```

## Usage

### Basic Setup

```go
import (
    "context"
    "github.com/common-library/go/aws/dynamodb"
    "github.com/aws/aws-sdk-go-v2/aws"
    aws_dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Initialize client
var client dynamodb.Client
err := client.CreateClient(
    context.TODO(),
    "us-east-1",
    "your-access-key",
    "your-secret-key",
    "",
)
```

### Table Operations

#### Create Table

```go
response, err := client.CreateTable(
    &aws_dynamodb.CreateTableInput{
        TableName: aws.String("users"),
        AttributeDefinitions: []types.AttributeDefinition{
            {
                AttributeName: aws.String("id"),
                AttributeType: types.ScalarAttributeTypeS,
            },
        },
        KeySchema: []types.KeySchemaElement{
            {
                AttributeName: aws.String("id"),
                KeyType:       types.KeyTypeHash,
            },
        },
        BillingMode: types.BillingModePayPerRequest,
    },
    true,  // wait for table to be active
    30,    // wait timeout in seconds
)
```

#### List Tables

```go
listOutput, err := client.ListTables(&aws_dynamodb.ListTablesInput{
    Limit: aws.Int32(10),
})
```

#### Describe Table

```go
describeOutput, err := client.DescribeTable("users")
fmt.Printf("Table status: %s\n", describeOutput.Table.TableStatus)
```

#### Update Table

```go
updateOutput, err := client.UpdateTable(&aws_dynamodb.UpdateTableInput{
    TableName: aws.String("users"),
    // ... update specifications
})
```

#### Delete Table

```go
deleteOutput, err := client.DeleteTable(
    "users",
    true,  // wait for deletion
    30,    // timeout seconds
)
```

### Item Operations

#### Put Item

```go
_, err = client.PutItem(&aws_dynamodb.PutItemInput{
    TableName: aws.String("users"),
    Item: map[string]types.AttributeValue{
        "id":    &types.AttributeValueMemberS{Value: "user-123"},
        "name":  &types.AttributeValueMemberS{Value: "John Doe"},
        "email": &types.AttributeValueMemberS{Value: "john@example.com"},
    },
})
```

#### Get Item

```go
output, err := client.GetItem(&aws_dynamodb.GetItemInput{
    TableName: aws.String("users"),
    Key: map[string]types.AttributeValue{
        "id": &types.AttributeValueMemberS{Value: "user-123"},
    },
})
```

#### Update Item

```go
updateOutput, err := client.UpdateItem(&aws_dynamodb.UpdateItemInput{
    TableName: aws.String("users"),
    Key: map[string]types.AttributeValue{
        "id": &types.AttributeValueMemberS{Value: "user-123"},
    },
    UpdateExpression: aws.String("SET #name = :name"),
    ExpressionAttributeNames: map[string]string{
        "#name": "name",
    },
    ExpressionAttributeValues: map[string]types.AttributeValue{
        ":name": &types.AttributeValueMemberS{Value: "Jane Doe"},
    },
})
```

#### Delete Item

```go
deleteOutput, err := client.DeleteItem(&aws_dynamodb.DeleteItemInput{
    TableName: aws.String("users"),
    Key: map[string]types.AttributeValue{
        "id": &types.AttributeValueMemberS{Value: "user-123"},
    },
})
```

### Query and Scan

#### Query Items

```go
queryOutput, err := client.Query(&aws_dynamodb.QueryInput{
    TableName:              aws.String("users"),
    KeyConditionExpression: aws.String("id = :id"),
    ExpressionAttributeValues: map[string]types.AttributeValue{
        ":id": &types.AttributeValueMemberS{Value: "user-123"},
    },
})
```

#### Scan Items

```go
scanOutput, err := client.Scan(&aws_dynamodb.ScanInput{
    TableName: aws.String("users"),
    FilterExpression: aws.String("age > :minAge"),
    ExpressionAttributeValues: map[string]types.AttributeValue{
        ":minAge": &types.AttributeValueMemberN{Value: "18"},
    },
})
```

#### Query with Pagination

```go
// First page
output, err := client.QueryPaginatorNextPage(&aws_dynamodb.QueryInput{
    TableName:              aws.String("users"),
    KeyConditionExpression: aws.String("status = :status"),
    ExpressionAttributeValues: map[string]types.AttributeValue{
        ":status": &types.AttributeValueMemberS{Value: "active"},
    },
    Limit: aws.Int32(100),
})

// Process results
for _, item := range output.Items {
    // Process each item
}

// Check if more pages exist
if output.LastEvaluatedKey != nil {
    // More pages available
}
```

#### Scan with Pagination

```go
for hasMorePages := true; hasMorePages; {
    output, err := client.ScanPaginatorNextPage(&aws_dynamodb.ScanInput{
        TableName: aws.String("users"),
        Limit: aws.Int32(100),
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Process output
    for _, item := range output.Items {
        // Process each item
    }
    
    hasMorePages = output.LastEvaluatedKey != nil
}
```

### TTL Management

#### Enable TTL

```go
_, err := client.UpdateTimeToLive("users", "expirationTime", true)
```

#### Check TTL Status

```go
ttlOutput, err := client.DescribeTimeToLive("users")
fmt.Printf("TTL Status: %s\n", ttlOutput.TimeToLiveDescription.TimeToLiveStatus)
```

#### Disable TTL

```go
_, err = client.UpdateTimeToLive("users", "expirationTime", false)
```

## API Reference

### Client Management
- `CreateClient(ctx, region, accessKey, secretKey, sessionToken, options...)` - Initialize client

### Table Operations
- `CreateTable(request, wait, waitTimeout, options...)` - Create table with optional wait
- `ListTables(request, options...)` - List all tables
- `DescribeTable(tableName, options...)` - Get table information
- `UpdateTable(request, options...)` - Modify table settings
- `DeleteTable(tableName, wait, waitTimeout, options...)` - Delete table with optional wait

### Item Operations
- `GetItem(request, options...)` - Retrieve single item
- `PutItem(request, options...)` - Create or replace item
- `UpdateItem(request, options...)` - Modify item attributes
- `DeleteItem(request, options...)` - Delete item

### Query and Scan
- `Query(request, options...)` - Query items by primary key
- `Scan(request, options...)` - Scan all items
- `QueryPaginatorNextPage(request, options...)` - Paginated query
- `ScanPaginatorNextPage(request, options...)` - Paginated scan

### TTL Management
- `DescribeTimeToLive(tableName, options...)` - Get TTL settings
- `UpdateTimeToLive(tableName, attributeName, enabled, options...)` - Enable/disable TTL

## Custom Endpoint Configuration

### DynamoDB Local

For local development and testing, use DynamoDB Local:

```go
var client dynamodb.Client
err := client.CreateClient(
    context.TODO(),
    "us-east-1",
    "dummy",
    "dummy",
    "",
    func(o *aws_dynamodb.Options) {
        o.BaseEndpoint = aws.String("http://localhost:8000")
    },
)
```

Run DynamoDB Local with Docker:

```bash
docker run -p 8000:8000 amazon/dynamodb-local
```

See [DynamoDB Local documentation](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/DynamoDBLocal.html) for more details.

## Implementation Details

- Uses `github.com/aws/aws-sdk-go-v2/service/dynamodb`
- AWS SDK v2 based implementation
- Context support for all operations
- Waiter support for table operations
- Pagination support for large result sets
- TTL management capabilities
- 1MB limit per query/scan (use paginators for larger datasets)

## Error Handling

All functions return errors that should be checked:

```go
response, err := client.GetItem(&aws_dynamodb.GetItemInput{...})
if err != nil {
    var nfe *types.ResourceNotFoundException
    if errors.As(err, &nfe) {
        // Handle resource not found
        log.Printf("Table or item not found")
    } else {
        log.Fatalf("Failed to get item: %v", err)
    }
}
```

Common error scenarios:
- Invalid credentials (authentication failure)
- Table does not exist
- Item not found
- Insufficient permissions
- Network connectivity issues
- Service throttling (rate limits)
- Validation errors (invalid attribute types, missing required fields)

## Best Practices

1. **Use Context with Timeout**
   ```go
   ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
   defer cancel()
   
   err := client.CreateClient(ctx, region, accessKey, secretKey, "")
   ```

2. **Pagination for Large Datasets**
   ```go
   // DynamoDB limits query/scan results to 1MB
   for hasMorePages := true; hasMorePages; {
       output, err := client.QueryPaginatorNextPage(request)
       // Process output
       hasMorePages = output.LastEvaluatedKey != nil
   }
   ```

3. **Conditional Writes**
   ```go
   _, err := client.PutItem(&aws_dynamodb.PutItemInput{
       TableName: aws.String("users"),
       Item: item,
       ConditionExpression: aws.String("attribute_not_exists(id)"),
   })
   ```

4. **Batch Operations**
   ```go
   // Use BatchWriteItem for multiple puts/deletes
   // Use BatchGetItem for multiple gets
   // Both operations support up to 25 items
   ```

5. **Error Handling with Retries**
   ```go
   if err != nil {
       var throttleErr *types.ProvisionedThroughputExceededException
       if errors.As(err, &throttleErr) {
           // Implement exponential backoff retry
       }
   }
   ```

6. **Credential Security**
   - Never hardcode credentials in source code
   - Use environment variables or AWS IAM roles
   - Rotate credentials regularly
   - Use session tokens for temporary access

## Examples

### Complete Table Lifecycle

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/common-library/go/aws/dynamodb"
    "github.com/aws/aws-sdk-go-v2/aws"
    aws_dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func main() {
    var client dynamodb.Client
    
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
    
    // Create table
    _, err = client.CreateTable(
        &aws_dynamodb.CreateTableInput{
            TableName: aws.String("users"),
            AttributeDefinitions: []types.AttributeDefinition{
                {
                    AttributeName: aws.String("id"),
                    AttributeType: types.ScalarAttributeTypeS,
                },
            },
            KeySchema: []types.KeySchemaElement{
                {
                    AttributeName: aws.String("id"),
                    KeyType:       types.KeyTypeHash,
                },
            },
            BillingMode: types.BillingModePayPerRequest,
        },
        true,  // wait for active
        30,    // timeout seconds
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Put item
    _, err = client.PutItem(&aws_dynamodb.PutItemInput{
        TableName: aws.String("users"),
        Item: map[string]types.AttributeValue{
            "id":    &types.AttributeValueMemberS{Value: "user-123"},
            "name":  &types.AttributeValueMemberS{Value: "John Doe"},
            "email": &types.AttributeValueMemberS{Value: "john@example.com"},
        },
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Get item
    output, err := client.GetItem(&aws_dynamodb.GetItemInput{
        TableName: aws.String("users"),
        Key: map[string]types.AttributeValue{
            "id": &types.AttributeValueMemberS{Value: "user-123"},
        },
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Retrieved item: %+v\n", output.Item)
    
    // Delete table
    _, err = client.DeleteTable("users", true, 30)
    if err != nil {
        log.Fatal(err)
    }
}
```

## Dependencies

- `github.com/aws/aws-sdk-go-v2/aws` - Core AWS SDK
- `github.com/aws/aws-sdk-go-v2/config` - Configuration loading
- `github.com/aws/aws-sdk-go-v2/credentials` - Credentials management
- `github.com/aws/aws-sdk-go-v2/service/dynamodb` - DynamoDB service

## Further Reading

- [AWS SDK for Go v2 Documentation](https://aws.github.io/aws-sdk-go-v2/)
- [DynamoDB Developer Guide](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/)
- [DynamoDB Best Practices](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/best-practices.html)
- [AWS package overview](../)
