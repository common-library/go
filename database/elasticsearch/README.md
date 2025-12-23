# Elasticsearch

Multi-version Elasticsearch client with unified interface supporting v7, v8, and v9.

## Overview

The elasticsearch package provides a consistent interface for working with Elasticsearch across multiple major versions (7.x, 8.x, and 9.x). This allows applications to switch between Elasticsearch versions without changing their code, or maintain compatibility with different deployment environments.

## Features

- **Multi-Version Support** - Compatible with Elasticsearch 7, 8, and 9
- **Unified Interface** - Single API across all versions
- **Document Operations** - Index, exists, delete, delete by query
- **Index Management** - Create, delete, exists, force merge
- **Template Management** - Put, delete, exists templates
- **Search Operations** - Full-text search with JSON DSL
- **Authentication** - Username/password, API key, certificate fingerprint
- **Cloud Support** - Elastic Cloud ID configuration
- **Thread-Safe** - Protected initialization to prevent data races

## Installation

```bash
# Install the interface and version you need
go get -u github.com/common-library/go/database/elasticsearch

# Elasticsearch v7
go get -u github.com/elastic/go-elasticsearch/v7

# Elasticsearch v8
go get -u github.com/elastic/go-elasticsearch/v8

# Elasticsearch v9
go get -u github.com/elastic/go-elasticsearch/v9
```

## Quick Start

### Elasticsearch v7

```go
import (
    "time"
    v7 "github.com/common-library/go/database/elasticsearch/v7"
)

func main() {
    client := &v7.Client{}
    
    // Initialize with basic auth
    err := client.Initialize(
        []string{"http://localhost:9200"},
        30*time.Second,
        "",          // cloudID
        "",          // apiKey
        "elastic",   // username
        "password",  // password
        "",          // certificateFingerprint
        nil,         // caCert
    )
    
    // Index a document
    err = client.Index("products", "1", `{
        "name": "Laptop",
        "price": 999.99,
        "category": "electronics"
    }`)
    
    // Check if document exists
    exists, err := client.Exists("products", "1")
    
    // Search
    result, err := client.Search("products", `{
        "query": {
            "match": {
                "category": "electronics"
            }
        }
    }`)
}
```

### Elasticsearch v8

```go
import (
    "time"
    v8 "github.com/common-library/go/database/elasticsearch/v8"
)

func main() {
    client := &v8.Client{}
    
    // Initialize with HTTPS and certificate fingerprint
    err := client.Initialize(
        []string{"https://localhost:9200"},
        30*time.Second,
        "",                                      // cloudID
        "",                                      // apiKey
        "elastic",                               // username
        "password",                              // password
        "a1b2c3...",                            // certificateFingerprint
        nil,                                     // caCert
    )
    
    // Use same API as v7
    client.Index("products", "1", `{"name":"Product"}`)
}
```

### Elasticsearch v9

```go
import (
    "time"
    v9 "github.com/common-library/go/database/elasticsearch/v9"
)

func main() {
    client := &v9.Client{}
    
    // Initialize with API key
    err := client.Initialize(
        []string{"https://localhost:9200"},
        30*time.Second,
        "",                    // cloudID
        "base64EncodedApiKey", // apiKey
        "",                    // username
        "",                    // password
        "",                    // certificateFingerprint
        nil,                   // caCert
    )
    
    // Use same API as v7 and v8
    client.Index("products", "1", `{"name":"Product"}`)
}
```

## Using the Interface

```go
import (
    "github.com/common-library/go/database/elasticsearch"
    v8 "github.com/common-library/go/database/elasticsearch/v8"
)

func processData(client elasticsearch.ClientInterface) error {
    // Works with any version
    return client.Index("data", "1", `{"value":100}`)
}

func main() {
    var client elasticsearch.ClientInterface
    client = &v8.Client{}
    client.Initialize(addresses, timeout, "", "", "", "", "", nil)
    
    processData(client)
}
```

## Document Operations

### Index Document

```go
// Index with explicit ID
err := client.Index("products", "product-123", `{
    "name": "Wireless Mouse",
    "price": 29.99,
    "stock": 150
}`)

// Index with generated ID (use empty string)
err = client.Index("logs", "", `{
    "timestamp": "2024-01-17T10:00:00Z",
    "level": "INFO",
    "message": "Application started"
}`)
```

### Check Document Exists

```go
exists, err := client.Exists("products", "product-123")
if err != nil {
    log.Fatal(err)
}

if exists {
    fmt.Println("Document found")
} else {
    fmt.Println("Document not found")
}
```

### Delete Document

```go
// Delete by ID
err := client.Delete("products", "product-123")

// Delete by query
err = client.DeleteByQuery(
    []string{"products"},
    `{
        "query": {
            "range": {
                "price": {"lt": 10}
            }
        }
    }`,
)
```

## Index Management

### Create Index

```go
// Create with mappings
err := client.IndicesCreate("products", `{
    "settings": {
        "number_of_shards": 3,
        "number_of_replicas": 1
    },
    "mappings": {
        "properties": {
            "name": {"type": "text"},
            "price": {"type": "float"},
            "category": {"type": "keyword"},
            "created_at": {"type": "date"}
        }
    }
}`)
```

### Check Index Exists

```go
exists, err := client.IndicesExists([]string{"products"})
```

### Delete Index

```go
// Delete single index
err := client.IndicesDelete([]string{"products"})

// Delete multiple indices
err = client.IndicesDelete([]string{"logs-2024-01", "logs-2024-02"})
```

### Force Merge

```go
// Force merge to optimize index
err := client.IndicesForcemerge([]string{"products"})
```

## Template Management

### Put Template

```go
err := client.IndicesPutTemplate("logs_template", `{
    "index_patterns": ["logs-*"],
    "settings": {
        "number_of_shards": 1
    },
    "mappings": {
        "properties": {
            "timestamp": {"type": "date"},
            "level": {"type": "keyword"},
            "message": {"type": "text"}
        }
    }
}`)
```

### Check Template Exists

```go
exists, err := client.IndicesExistsTemplate([]string{"logs_template"})
```

### Delete Template

```go
err := client.IndicesDeleteTemplate("logs_template")
```

## Search Operations

### Basic Search

```go
result, err := client.Search("products", `{
    "query": {
        "match_all": {}
    }
}`)

fmt.Println(result) // JSON response string
```

### Match Query

```go
result, err := client.Search("products", `{
    "query": {
        "match": {
            "name": "laptop"
        }
    }
}`)
```

### Range Query

```go
result, err := client.Search("products", `{
    "query": {
        "range": {
            "price": {
                "gte": 100,
                "lte": 500
            }
        }
    },
    "sort": [
        {"price": "asc"}
    ]
}`)
```

### Aggregations

```go
result, err := client.Search("products", `{
    "size": 0,
    "aggs": {
        "categories": {
            "terms": {
                "field": "category"
            }
        },
        "avg_price": {
            "avg": {
                "field": "price"
            }
        }
    }
}`)
```

## Complete Examples

### Product Catalog

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "time"
    
    v8 "github.com/common-library/go/database/elasticsearch/v8"
)

type Product struct {
    Name     string  `json:"name"`
    Price    float64 `json:"price"`
    Category string  `json:"category"`
    Stock    int     `json:"stock"`
}

func main() {
    client := &v8.Client{}
    
    err := client.Initialize(
        []string{"http://localhost:9200"},
        30*time.Second,
        "", "", "elastic", "password", "", nil,
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Create index
    indexMapping := `{
        "settings": {
            "number_of_shards": 1
        },
        "mappings": {
            "properties": {
                "name": {"type": "text"},
                "price": {"type": "float"},
                "category": {"type": "keyword"},
                "stock": {"type": "integer"}
            }
        }
    }`
    
    exists, _ := client.IndicesExists([]string{"products"})
    if !exists {
        client.IndicesCreate("products", indexMapping)
    }
    
    // Index products
    products := []Product{
        {Name: "Laptop", Price: 999.99, Category: "electronics", Stock: 50},
        {Name: "Mouse", Price: 29.99, Category: "electronics", Stock: 200},
        {Name: "Desk", Price: 299.99, Category: "furniture", Stock: 30},
    }
    
    for i, p := range products {
        data, _ := json.Marshal(p)
        client.Index("products", fmt.Sprintf("%d", i+1), string(data))
    }
    
    // Search electronics
    result, _ := client.Search("products", `{
        "query": {
            "term": {
                "category": "electronics"
            }
        }
    }`)
    
    fmt.Println(result)
}
```

### Log Management

```go
package main

import (
    "fmt"
    "time"
    
    v7 "github.com/common-library/go/database/elasticsearch/v7"
)

func main() {
    client := &v7.Client{}
    client.Initialize(
        []string{"http://localhost:9200"},
        30*time.Second,
        "", "", "", "", "", nil,
    )
    
    // Create template for daily log indices
    template := `{
        "index_patterns": ["logs-*"],
        "settings": {
            "number_of_shards": 2,
            "number_of_replicas": 0
        },
        "mappings": {
            "properties": {
                "timestamp": {"type": "date"},
                "level": {"type": "keyword"},
                "message": {"type": "text"},
                "service": {"type": "keyword"}
            }
        }
    }`
    
    client.IndicesPutTemplate("logs_template", template)
    
    // Index log entries
    today := time.Now().Format("2006-01-02")
    indexName := fmt.Sprintf("logs-%s", today)
    
    client.Index(indexName, "", `{
        "timestamp": "2024-01-17T10:00:00Z",
        "level": "INFO",
        "message": "Application started",
        "service": "api"
    }`)
    
    client.Index(indexName, "", `{
        "timestamp": "2024-01-17T10:05:00Z",
        "level": "ERROR",
        "message": "Connection timeout",
        "service": "api"
    }`)
    
    // Search error logs
    result, _ := client.Search(indexName, `{
        "query": {
            "term": {
                "level": "ERROR"
            }
        },
        "sort": [
            {"timestamp": "desc"}
        ]
    }`)
    
    fmt.Println(result)
    
    // Delete old indices
    oldDate := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
    oldIndex := fmt.Sprintf("logs-%s", oldDate)
    
    exists, _ := client.IndicesExists([]string{oldIndex})
    if exists {
        client.IndicesDelete([]string{oldIndex})
    }
}
```

### E-commerce Search

```go
package main

import (
    "time"
    
    v9 "github.com/common-library/go/database/elasticsearch/v9"
)

func main() {
    client := &v9.Client{}
    client.Initialize(
        []string{"http://localhost:9200"},
        30*time.Second,
        "", "", "", "", "", nil,
    )
    
    // Multi-field search
    searchQuery := `{
        "query": {
            "bool": {
                "must": [
                    {
                        "multi_match": {
                            "query": "laptop",
                            "fields": ["name^2", "description"]
                        }
                    }
                ],
                "filter": [
                    {
                        "range": {
                            "price": {
                                "gte": 500,
                                "lte": 2000
                            }
                        }
                    },
                    {
                        "term": {
                            "in_stock": true
                        }
                    }
                ]
            }
        },
        "sort": [
            {"_score": "desc"},
            {"price": "asc"}
        ],
        "from": 0,
        "size": 20
    }`
    
    result, _ := client.Search("products", searchQuery)
    
    // Aggregations for faceted search
    facetQuery := `{
        "size": 0,
        "aggs": {
            "price_ranges": {
                "range": {
                    "field": "price",
                    "ranges": [
                        {"to": 500},
                        {"from": 500, "to": 1000},
                        {"from": 1000, "to": 2000},
                        {"from": 2000}
                    ]
                }
            },
            "brands": {
                "terms": {
                    "field": "brand",
                    "size": 10
                }
            }
        }
    }`
    
    facets, _ := client.Search("products", facetQuery)
    _ = facets
}
```

## Authentication Methods

### Username/Password

```go
client.Initialize(
    []string{"http://localhost:9200"},
    timeout,
    "",        // cloudID
    "",        // apiKey
    "elastic", // username
    "password", // password
    "",        // certificateFingerprint
    nil,       // caCert
)
```

### API Key

```go
client.Initialize(
    []string{"http://localhost:9200"},
    timeout,
    "",                      // cloudID
    "base64EncodedApiKey",   // apiKey
    "",                      // username
    "",                      // password
    "",                      // certificateFingerprint
    nil,                     // caCert
)
```

### Certificate Fingerprint (v8+)

```go
client.Initialize(
    []string{"https://localhost:9200"},
    timeout,
    "",                 // cloudID
    "",                 // apiKey
    "elastic",          // username
    "password",         // password
    "a1:b2:c3:...",    // certificateFingerprint
    nil,                // caCert
)
```

### CA Certificate

```go
caCert, _ := os.ReadFile("/path/to/ca.crt")

client.Initialize(
    []string{"https://localhost:9200"},
    timeout,
    "",        // cloudID
    "",        // apiKey
    "elastic", // username
    "password", // password
    "",        // certificateFingerprint
    caCert,    // caCert
)
```

### Elastic Cloud

```go
client.Initialize(
    nil,                    // addresses (not needed with cloudID)
    timeout,
    "my-cloud-id:...",     // cloudID
    "",                    // apiKey
    "elastic",             // username
    "password",            // password
    "",                    // certificateFingerprint
    nil,                   // caCert
)
```

## API Reference

### ClientInterface

All version-specific clients implement this interface:

#### `Initialize(addresses []string, timeout time.Duration, cloudID, apiKey, username, password, certificateFingerprint string, caCert []byte) error`

Initialize Elasticsearch client with connection settings.

**Parameters:**
- `addresses` - Elasticsearch node URLs (e.g., []string{"http://localhost:9200"})
- `timeout` - HTTP response timeout duration
- `cloudID` - Elastic Cloud deployment ID (optional)
- `apiKey` - Base64-encoded API key (optional)
- `username` - Basic auth username (optional)
- `password` - Basic auth password (optional)
- `certificateFingerprint` - Server certificate fingerprint for HTTPS (optional, v8+)
- `caCert` - CA certificate bytes for HTTPS (optional)

#### `Exists(index, documentID string) (bool, error)`

Check if a document exists in an index.

#### `Index(index, documentID, body string) error`

Index a document. Use empty documentID for auto-generated ID.

#### `Delete(index, documentID string) error`

Delete a document by ID.

#### `DeleteByQuery(indices []string, body string) error`

Delete documents matching a query.

#### `IndicesExists(indices []string) (bool, error)`

Check if indices exist.

#### `IndicesCreate(index, body string) error`

Create an index with settings and mappings.

#### `IndicesDelete(indices []string) error`

Delete one or more indices.

#### `IndicesExistsTemplate(name []string) (bool, error)`

Check if index templates exist.

#### `IndicesPutTemplate(name, body string) error`

Create or update an index template.

#### `IndicesDeleteTemplate(name string) error`

Delete an index template.

#### `IndicesForcemerge(indices []string) error`

Force merge indices to optimize storage.

#### `Search(index, body string) (string, error)`

Execute a search query. Returns JSON response as string.

## Version Differences

### Elasticsearch 7 vs 8 vs 9

**All Versions Support:**
- Document CRUD operations
- Index management
- Template management
- Search queries
- Aggregations

**Version-Specific Features:**

| Feature | v7 | v8 | v9 |
|---------|----|----|-----|
| Certificate Fingerprint | ❌ | ✅ | ✅ |
| Improved Security | Basic | Enhanced | Enhanced |
| Performance | Good | Better | Best |
| Type Mappings | `_doc` | Removed | Removed |

**Migration Notes:**
- v7 → v8: Update certificate handling, remove `_doc` type
- v8 → v9: Minimal changes, mostly performance improvements
- All versions share the same `ClientInterface`

## Best Practices

### 1. Use Appropriate Timeouts

```go
// Short timeout for simple operations
client.Initialize(addresses, 10*time.Second, ...)

// Longer timeout for complex searches
client.Initialize(addresses, 60*time.Second, ...)
```

### 2. Handle Errors Properly

```go
exists, err := client.Exists("index", "id")
if err != nil {
    // Check error type
    if strings.Contains(err.Error(), "status : (404)") {
        // Index doesn't exist
    } else {
        // Connection or other error
        log.Fatal(err)
    }
}
```

### 3. Use Templates for Time-Series Data

```go
// Define template once
client.IndicesPutTemplate("logs_template", templateBody)

// Indices auto-apply settings
client.Index("logs-2024-01-17", "", logEntry)
client.Index("logs-2024-01-18", "", logEntry)
```

### 4. Batch Operations

```go
// Use bulk API for multiple documents
for _, doc := range documents {
    client.Index("index", doc.ID, doc.JSON())
}

// Better: Use Elasticsearch Bulk API directly
```

### 5. Optimize Index Settings

```go
indexSettings := `{
    "settings": {
        "number_of_shards": 1,      // Reduce for small indices
        "refresh_interval": "30s",   // Increase for bulk indexing
        "number_of_replicas": 0      // Disable for development
    }
}`
```

### 6. Clean Up Old Indices

```go
// Delete indices older than 30 days
cutoff := time.Now().AddDate(0, 0, -30)
for d := cutoff; d.Before(time.Now()); d = d.AddDate(0, 0, 1) {
    indexName := fmt.Sprintf("logs-%s", d.Format("2006-01-02"))
    client.IndicesDelete([]string{indexName})
}
```

## Testing

### Using Testutil

```go
import (
    "testing"
    "github.com/common-library/go/database/elasticsearch/testutil"
)

func TestElasticsearch(t *testing.T) {
    // Automatically creates testcontainer
    client := testutil.GetTestClient(t, "v8", []string{"http://localhost:9200"})
    
    err := client.Index("test", "1", `{"value":100}`)
    if err != nil {
        t.Fatal(err)
    }
    
    exists, err := client.Exists("test", "1")
    if !exists {
        t.Error("Document should exist")
    }
}
```

## Troubleshooting

### Connection Refused

```go
// Check if Elasticsearch is running
// Verify addresses are correct
client.Initialize(
    []string{"http://localhost:9200"}, // Not https:// for local dev
    30*time.Second,
    "", "", "", "", "", nil,
)
```

### Certificate Verification Failed

```go
// Use certificate fingerprint (v8+)
client.Initialize(
    []string{"https://localhost:9200"},
    timeout,
    "", "", "elastic", "password",
    "fingerprint-from-elasticsearch-setup", // Get from ES startup
    nil,
)
```

### Index Already Exists

```go
exists, _ := client.IndicesExists([]string{"myindex"})
if !exists {
    client.IndicesCreate("myindex", mappings)
} else {
    // Delete and recreate, or skip
    client.IndicesDelete([]string{"myindex"})
    client.IndicesCreate("myindex", mappings)
}
```

### Document Not Found

```go
exists, err := client.Exists("index", "id")
if err != nil {
    log.Fatal(err)
}

if !exists {
    // Document doesn't exist, index it
    client.Index("index", "id", documentJSON)
}
```

## Limitations

1. **No Bulk API** - Must use individual Index calls (use official client for bulk)
2. **No Scroll API** - Cannot paginate through large result sets
3. **String-Based Queries** - No query builder, must construct JSON manually
4. **No Response Parsing** - Search returns raw JSON string
5. **Limited Template Support** - Only basic template operations
6. **No Async Operations** - All operations are synchronous

## Dependencies

**Common:**
- `github.com/thedevsaddam/gojsonq/v2` - JSON query for error parsing

**Version-Specific:**
- `github.com/elastic/go-elasticsearch/v7` - ES 7.x client
- `github.com/elastic/go-elasticsearch/v8` - ES 8.x client
- `github.com/elastic/go-elasticsearch/v9` - ES 9.x client
- `github.com/elastic/elastic-transport-go/v8` - Transport layer (v8, v9)

**Testing:**
- `github.com/testcontainers/testcontainers-go` - Test containers
- `github.com/stretchr/testify` - Test assertions

## Further Reading

- [Elasticsearch Official Documentation](https://www.elastic.co/guide/en/elasticsearch/reference/current/index.html)
- [Go Elasticsearch Client v7](https://github.com/elastic/go-elasticsearch/tree/7.x)
- [Go Elasticsearch Client v8](https://github.com/elastic/go-elasticsearch/tree/8.x)
- [Go Elasticsearch Client v9](https://github.com/elastic/go-elasticsearch)
- [Query DSL](https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl.html)
- [Mapping](https://www.elastic.co/guide/en/elasticsearch/reference/current/mapping.html)
