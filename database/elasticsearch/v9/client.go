// Package v9 provides Elasticsearch version 9 client implementation.
//
// This package wraps the official Elasticsearch 9.x Go client, providing a simplified
// interface for common operations while maintaining compatibility with the ClientInterface.
//
// Features:
//   - Full Elasticsearch 9.x API support
//   - Document operations (index, exists, delete, delete by query)
//   - Index management (create, delete, exists, force merge)
//   - Template management (put, delete, exists)
//   - Search operations with JSON responses
//   - Certificate fingerprint authentication
//   - Cloud ID support
//
// Example:
//
//	client := &v9.Client{}
//	err := client.Initialize([]string{"https://localhost:9200"}, 30*time.Second, "", "", "elastic", "password", "", nil)
//	err = client.Index("products", "1", `{"name":"Product 1","price":29.99}`)
package v9

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esapi"
	"github.com/thedevsaddam/gojsonq/v2"

	"github.com/common-library/go/database/elasticsearch/internal/eslock"
)

// Client is a struct that provides client related methods.
type Client struct {
	client *elasticsearch.Client
}

// Initialize initializes the Elasticsearch v9 client with connection configuration.
//
// Parameters:
//   - addresses: Elasticsearch node URLs (e.g., []string{"http://localhost:9200"})
//   - timeout: HTTP response timeout in seconds
//   - cloudID: Elastic Cloud deployment ID (optional, empty string if not used)
//   - apiKey: Base64-encoded API key for authentication (optional)
//   - username: Username for basic authentication (optional)
//   - password: Password for basic authentication (optional)
//   - certificateFingerprint: Server certificate fingerprint for HTTPS (optional)
//   - caCert: CA certificate bytes for HTTPS verification (optional)
//
// Returns error if client creation fails.
//
// The client is protected by a mutex during initialization to prevent data races
// in the underlying elastictransport library.
//
// Example:
//
//	err := client.Initialize(
//		[]string{"https://localhost:9200"},
//		30*time.Second,
//		"", "base64ApiKey", "", "", "", nil,
//	)
func (c *Client) Initialize(addresses []string, timeout time.Duration, cloudID, apiKey, username, password, certificateFingerprint string, caCert []byte) error {
	eslock.InitMu.Lock()
	defer eslock.InitMu.Unlock()

	config := elasticsearch.Config{
		CloudID:                cloudID,
		APIKey:                 apiKey,
		Username:               username,
		Password:               password,
		CertificateFingerprint: certificateFingerprint,

		Addresses:         addresses,
		EnableDebugLogger: true,
		Logger:            &elastictransport.ColorLogger{Output: os.Stdout},
		Transport: &http.Transport{
			ResponseHeaderTimeout: timeout * time.Second,
		},
	}
	if len(caCert) != 0 {
		config.CACert = caCert
	}

	if client, err := elasticsearch.NewClient(config); err != nil {
		return err
	} else {
		c.client = client
	}

	return nil
}

// Exists checks if a document exists in the specified index.
//
// Parameters:
//   - index: The name of the index
//   - documentID: The document identifier
//
// Returns true if the document exists, false if not found.
// Returns error if client is not initialized or request fails.
//
// The request uses refresh=true to ensure real-time visibility.
//
// Example:
//
//	exists, err := client.Exists("products", "product-123")
//	if exists {
//		fmt.Println("Document found")
//	}
func (c *Client) Exists(index, documentID string) (bool, error) {
	if c.client == nil {
		return false, errors.New("please call Initialize first")
	}

	refresh := true

	request := esapi.ExistsRequest{
		Index:      index,
		DocumentID: documentID,
		Refresh:    &refresh,
		Human:      true,
		ErrorTrace: true}

	if response, err := request.Do(context.Background(), c.client); err != nil {
		return false, err
	} else {
		defer response.Body.Close()

		if response.StatusCode == http.StatusNotFound {
			return false, nil
		} else if response.IsError() {
			return false, c.responseErrorToError(response.Status(), response.Body)
		}
	}

	return true, nil
}

// Index indexes a document in the specified index.
//
// Parameters:
//   - index: The name of the index
//   - documentID: The document identifier (use empty string for auto-generated ID)
//   - body: JSON document body as a string
//
// Returns error if client is not initialized or indexing fails.
//
// The document is immediately available for search (refresh=true).
//
// Example:
//
//	err := client.Index("products", "1", `{
//		"name": "Laptop",
//		"price": 999.99
//	}`)
func (c *Client) Index(index, documentID, body string) error {
	if c.client == nil {
		return errors.New("please call Initialize first")
	}

	builder := strings.Builder{}
	if _, err := builder.WriteString(body); err != nil {
		return err
	}

	request := esapi.IndexRequest{
		Index:      index,
		DocumentID: documentID,
		Body:       strings.NewReader(builder.String()),
		Refresh:    "true",
		Human:      true,
		ErrorTrace: true}

	if response, err := request.Do(context.Background(), c.client); err != nil {
		return err
	} else {
		defer response.Body.Close()

		if response.IsError() {
			return c.responseErrorToError(response.Status(), response.Body)
		}
	}

	return nil
}

// Delete deletes a document from the specified index.
//
// Parameters:
//   - index: The name of the index
//   - documentID: The document identifier to delete
//
// Returns error if client is not initialized or deletion fails.
//
// The deletion is immediately visible (refresh=true).
//
// Example:
//
//	err := client.Delete("products", "product-123")
func (c *Client) Delete(index, documentID string) error {
	if c.client == nil {
		return errors.New("please call Initialize first")
	}

	request := esapi.DeleteRequest{
		Index:      index,
		DocumentID: documentID,
		Refresh:    "true",
		Human:      true,
		ErrorTrace: true}

	if response, err := request.Do(context.Background(), c.client); err != nil {
		return err
	} else {
		defer response.Body.Close()

		if response.IsError() {
			return c.responseErrorToError(response.Status(), response.Body)
		}
	}

	return nil
}

// DeleteByQuery deletes all documents matching a query.
//
// Parameters:
//   - indices: List of index names to search
//   - body: JSON query body as a string
//
// Returns error if client is not initialized or deletion fails.
//
// The deletions are immediately visible (refresh=true).
//
// Example:
//
//	err := client.DeleteByQuery(
//		[]string{"products"},
//		`{"query": {"range": {"price": {"lt": 10}}}}`
//	)
func (c *Client) DeleteByQuery(indices []string, body string) error {
	if c.client == nil {
		return errors.New("please call Initialize first")
	}

	builder := strings.Builder{}
	if _, err := builder.WriteString(body); err != nil {
		return err
	}

	refresh := true

	request := esapi.DeleteByQueryRequest{
		Index:      indices,
		Body:       strings.NewReader(builder.String()),
		Refresh:    &refresh,
		Human:      true,
		ErrorTrace: true}

	if response, err := request.Do(context.Background(), c.client); err != nil {
		return err
	} else {
		defer response.Body.Close()

		if response.IsError() {
			return c.responseErrorToError(response.Status(), response.Body)
		}
	}

	return nil
}

// IndicesExists checks if one or more indices exist.
//
// Parameters:
//   - indices: List of index names to check
//
// Returns true if all specified indices exist, false otherwise.
// Returns error if client is not initialized or request fails.
//
// Example:
//
//	exists, err := client.IndicesExists([]string{"products", "orders"})
//	if !exists {
//		// Create indices
//	}
func (c *Client) IndicesExists(indices []string) (bool, error) {
	if c.client == nil {
		return false, errors.New("please call Initialize first")
	}

	request := esapi.IndicesExistsRequest{
		Index:      indices,
		Human:      true,
		ErrorTrace: true}

	if response, err := request.Do(context.Background(), c.client); err != nil {
		return false, err
	} else {
		defer response.Body.Close()

		if response.StatusCode == http.StatusNotFound {
			return false, nil
		} else if response.IsError() {
			return false, c.responseErrorToError(response.Status(), response.Body)
		}
	}

	return true, nil
}

// IndicesCreate creates an index with optional settings and mappings.
//
// Parameters:
//   - index: The name of the index to create
//   - body: JSON configuration with settings and mappings
//
// Returns error if client is not initialized, index already exists, or creation fails.
//
// Example:
//
//	err := client.IndicesCreate("products", `{
//		"settings": {"number_of_shards": 1},
//		"mappings": {
//			"properties": {
//				"name": {"type": "text"},
//				"price": {"type": "float"}
//			}
//		}
//	}`)
func (c *Client) IndicesCreate(index, body string) error {
	if c.client == nil {
		return errors.New("please call Initialize first")
	}

	builder := strings.Builder{}
	if _, err := builder.WriteString(body); err != nil {
		return err
	}

	request := esapi.IndicesCreateRequest{
		Index:      index,
		Body:       strings.NewReader(builder.String()),
		Human:      true,
		ErrorTrace: true}

	if response, err := request.Do(context.Background(), c.client); err != nil {
		return err
	} else {
		defer response.Body.Close()

		if response.IsError() {
			return c.responseErrorToError(response.Status(), response.Body)
		}
	}

	return nil
}

// IndicesDelete deletes one or more indices.
//
// Parameters:
//   - indices: List of index names to delete
//
// Returns error if client is not initialized or deletion fails.
//
// Warning: This permanently deletes all data in the specified indices.
//
// Example:
//
//	err := client.IndicesDelete([]string{"old-logs-2023", "old-logs-2022"})
func (c *Client) IndicesDelete(indices []string) error {
	if c.client == nil {
		return errors.New("please call Initialize first")
	}

	request := esapi.IndicesDeleteRequest{
		Index:      indices,
		Human:      true,
		ErrorTrace: true}

	if response, err := request.Do(context.Background(), c.client); err != nil {
		return err
	} else {
		defer response.Body.Close()

		if response.IsError() {
			return c.responseErrorToError(response.Status(), response.Body)
		}
	}

	return nil
}

// IndicesExistsTemplate checks if one or more index templates exist.
//
// Parameters:
//   - name: List of template names to check
//
// Returns true if all specified templates exist, false otherwise.
// Returns error if client is not initialized or request fails.
//
// Example:
//
//	exists, err := client.IndicesExistsTemplate([]string{"logs_template"})
//	if !exists {
//		// Create template
//	}
func (c *Client) IndicesExistsTemplate(name []string) (bool, error) {
	if c.client == nil {
		return false, errors.New("please call Initialize first")
	}

	request := esapi.IndicesExistsTemplateRequest{
		Name:       name,
		Human:      true,
		ErrorTrace: true}

	if response, err := request.Do(context.Background(), c.client); err != nil {
		return false, err
	} else {
		defer response.Body.Close()

		if response.StatusCode == http.StatusNotFound {
			return false, nil
		} else if response.IsError() {
			return false, c.responseErrorToError(response.Status(), response.Body)
		}
	}

	return true, nil
}

// IndicesPutTemplate creates or updates an index template.
//
// Parameters:
//   - name: The template name
//   - body: JSON template configuration with index patterns, settings, and mappings
//
// Returns error if client is not initialized or operation fails.
//
// Templates automatically apply settings and mappings to matching indices.
//
// Example:
//
//	err := client.IndicesPutTemplate("logs_template", `{
//		"index_patterns": ["logs-*"],
//		"settings": {"number_of_shards": 1},
//		"mappings": {
//			"properties": {
//				"timestamp": {"type": "date"}
//			}
//		}
//	}`)
func (c *Client) IndicesPutTemplate(name, body string) error {
	if c.client == nil {
		return errors.New("please call Initialize first")
	}

	builder := strings.Builder{}
	if _, err := builder.WriteString(body); err != nil {
		return err
	}

	request := esapi.IndicesPutTemplateRequest{
		Name:       name,
		Body:       strings.NewReader(builder.String()),
		Human:      true,
		ErrorTrace: true}

	if response, err := request.Do(context.Background(), c.client); err != nil {
		return err
	} else {
		defer response.Body.Close()

		if response.IsError() {
			return c.responseErrorToError(response.Status(), response.Body)
		}
	}

	return nil
}

// IndicesDeleteTemplate deletes an index template.
//
// Parameters:
//   - name: The template name to delete
//
// Returns error if client is not initialized or deletion fails.
//
// Note: Deleting a template does not affect existing indices created from it.
//
// Example:
//
//	err := client.IndicesDeleteTemplate("old_logs_template")
func (c *Client) IndicesDeleteTemplate(name string) error {
	if c.client == nil {
		return errors.New("please call Initialize first")
	}

	request := esapi.IndicesDeleteTemplateRequest{
		Name:       name,
		Human:      true,
		ErrorTrace: true}

	if response, err := request.Do(context.Background(), c.client); err != nil {
		return err
	} else {
		defer response.Body.Close()

		if response.IsError() {
			return c.responseErrorToError(response.Status(), response.Body)
		}
	}

	return nil
}

// IndicesForcemerge performs a force merge operation on indices.
//
// Parameters:
//   - indices: List of index names to force merge
//
// Returns error if client is not initialized or operation fails.
//
// Force merge optimizes index storage by merging segments and removing deleted documents.
// This operation is I/O intensive and should be used carefully on production systems.
//
// Example:
//
//	// Optimize old read-only indices
//	err := client.IndicesForcemerge([]string{"logs-2023-12"})
func (c *Client) IndicesForcemerge(indices []string) error {
	if c.client == nil {
		return errors.New("please call Initialize first")
	}

	onlyExpungeDeletes := true

	request := esapi.IndicesForcemergeRequest{
		Index:              indices,
		OnlyExpungeDeletes: &onlyExpungeDeletes,
		Human:              true,
		ErrorTrace:         true}

	if response, err := request.Do(context.Background(), c.client); err != nil {
		return err
	} else {
		defer response.Body.Close()

		if response.IsError() {
			return c.responseErrorToError(response.Status(), response.Body)
		}
	}

	return nil
}

// Search executes a search query and returns the raw JSON response.
//
// Parameters:
//   - index: The index name to search
//   - body: JSON search query body
//
// Returns the search results as a JSON string.
// Returns error if client is not initialized or search fails.
//
// The response includes hits, aggregations, and metadata.
// Parse the JSON string to extract specific results.
//
// Example:
//
//	result, err := client.Search("products", `{
//		"query": {
//			"match": {"category": "electronics"}
//		},
//		"size": 10
//	}`)
func (c *Client) Search(index, body string) (string, error) {
	if c.client == nil {
		return "", errors.New("please call Initialize first")
	}

	builder := strings.Builder{}
	if _, err := builder.WriteString(body); err != nil {
		return "", err
	}

	if response, err := c.client.Search(
		c.client.Search.WithContext(context.Background()),
		c.client.Search.WithIndex(index),
		c.client.Search.WithBody(strings.NewReader(builder.String())),
		c.client.Search.WithTrackTotalHits(true),
		c.client.Search.WithPretty()); err != nil {
		return "", err
	} else {
		defer response.Body.Close()

		if response.IsError() {
			return "", c.responseErrorToError(response.Status(), response.Body)
		}

		result, err := io.ReadAll(response.Body)
		return string(result), err
	}
}

func (c *Client) responseErrorToError(status string, reader io.Reader) error {
	buffer := new(bytes.Buffer)
	buffer.ReadFrom(reader)

	if len(buffer.String()) == 0 {
		return fmt.Errorf("response error - status : (%s)", status)
	}

	return fmt.Errorf("response error - status : (%s), type : (%s), reason : (%s)",
		status,
		gojsonq.New().FromString(buffer.String()).Find("error.type").(string),
		gojsonq.New().FromString(buffer.String()).Find("error.reason").(string))
}
