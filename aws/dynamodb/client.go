// Package dynamodb provides utilities for working with AWS DynamoDB.
//
// This package wraps the AWS SDK v2 DynamoDB client to provide simplified
// functions for table management, item operations, and query/scan functionality.
//
// Features:
//   - Table operations (create, list, describe, update, delete)
//   - Item operations (get, put, update, delete)
//   - Query and scan with pagination support
//   - TTL (Time To Live) management
//   - Waiter support for table creation/deletion
//
// Example usage:
//
//	var client dynamodb.Client
//	err := client.CreateClient(ctx, "us-east-1", "key", "secret", "")
//	response, err := client.PutItem(&dynamodb.PutItemInput{...})
package dynamodb

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Client is a struct that provides client related methods.
type Client struct {
	ctx    context.Context
	client *dynamodb.Client
}

// CreateClient creates a DynamoDB client with service-specific options.
//
// This is the recommended way to create a DynamoDB client with custom endpoint
// and other service-specific configurations.
//
// Parameters:
//   - ctx: context for the client lifecycle
//   - region: AWS region (e.g., "us-east-1")
//   - accessKey: AWS access key ID
//   - secretAccessKey: AWS secret access key
//   - sessionToken: optional session token (empty string if not using temporary credentials)
//   - dynamodbOptionsFuncs: optional service-specific configuration functions
//
// Example:
//
//	err := client.CreateClient(context.TODO(), "us-east-1", "access-key", "secret-key", "",
//	    func(o *dynamodb.Options) {
//	        o.BaseEndpoint = aws.String("http://localhost:8000")
//	    })
func (c *Client) CreateClient(ctx context.Context, region, accessKey, secretAccessKey, sessionToken string, dynamodbOptionsFuncs ...func(*dynamodb.Options)) error {
	c.ctx = ctx

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretAccessKey, sessionToken)))
	if err != nil {
		return err
	}

	c.client = dynamodb.NewFromConfig(cfg, dynamodbOptionsFuncs...)

	return nil
}

// CreateTable creates a new DynamoDB table.
//
// Parameters:
//   - request: CreateTableInput containing table definition
//   - wait: if true, waits for table to be active before returning
//   - waitTimeout: timeout duration in seconds for table creation wait
//   - optionFunctions: optional service-specific configuration functions
//
// Returns the CreateTableOutput and any error encountered.
// See client_test.go for detailed examples.
//
// Example:
//
//	response, err := client.CreateTable(&dynamodb.CreateTableInput{...}, true, 10)
func (c *Client) CreateTable(request *dynamodb.CreateTableInput, wait bool, waitTimeout time.Duration, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error) {
	response, err := c.client.CreateTable(c.ctx, request, optionFunctions...)
	if err != nil {
		return nil, err
	}

	if wait {
		waiter := dynamodb.NewTableExistsWaiter(c.client)
		if err := waiter.Wait(
			c.ctx,
			&dynamodb.DescribeTableInput{TableName: request.TableName},
			waitTimeout*time.Second); err != nil {
			return nil, err
		}
	}

	return response, nil
}

// ListTables retrieves the list of table names in the current region.
//
// Parameters:
//   - request: ListTablesInput with optional pagination parameters
//   - optionFunctions: optional service-specific configuration functions
//
// Returns the ListTablesOutput containing table names and any error encountered.
//
// Example:
//
//	response, err := client.ListTables(&dynamodb.ListTablesInput{Limit: aws.Int32(10)})
func (c *Client) ListTables(request *dynamodb.ListTablesInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.ListTablesOutput, error) {
	return c.client.ListTables(c.ctx, request, optionFunctions...)
}

// DescribeTable retrieves information about a table.
//
// Parameters:
//   - tableName: name of the table to describe
//   - optionFunctions: optional service-specific configuration functions
//
// Returns the DescribeTableOutput containing table metadata and any error encountered.
//
// Example:
//
//	response, err := client.DescribeTable("users")
func (c *Client) DescribeTable(tableName string, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error) {
	return c.client.DescribeTable(c.ctx, &dynamodb.DescribeTableInput{TableName: aws.String(tableName)}, optionFunctions...)
}

// UpdateTable modifies the provisioned throughput settings, global secondary indexes,
// or DynamoDB Streams settings for a table.
//
// Parameters:
//   - request: UpdateTableInput containing the modifications
//   - optionFunctions: optional service-specific configuration functions
//
// Returns the UpdateTableOutput and any error encountered.
//
// Example:
//
//	response, err := client.UpdateTable(&dynamodb.UpdateTableInput{...})
func (c *Client) UpdateTable(request *dynamodb.UpdateTableInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.UpdateTableOutput, error) {
	return c.client.UpdateTable(c.ctx, request, optionFunctions...)
}

// DeleteTable deletes a table and all of its items.
//
// Parameters:
//   - tableName: name of the table to delete
//   - wait: if true, waits for table deletion to complete before returning
//   - waitTimeout: timeout duration in seconds for deletion wait
//   - optionFunctions: optional service-specific configuration functions
//
// Returns the DeleteTableOutput and any error encountered.
//
// Example:
//
//	response, err := client.DeleteTable("users", true, 10)
func (c *Client) DeleteTable(tableName string, wait bool, waitTimeout time.Duration, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.DeleteTableOutput, error) {
	response, err := c.client.DeleteTable(c.ctx, &dynamodb.DeleteTableInput{TableName: aws.String(tableName)}, optionFunctions...)
	if err != nil {
		return nil, err
	}

	if wait {
		waiter := dynamodb.NewTableNotExistsWaiter(c.client)
		if err := waiter.Wait(
			c.ctx,
			&dynamodb.DescribeTableInput{TableName: aws.String(tableName)},
			waitTimeout*time.Second); err != nil {
			return nil, err
		}
	}

	return response, nil
}

// GetItem retrieves a single item from a table.
//
// Parameters:
//   - request: GetItemInput containing the table name and primary key
//   - optionFunctions: optional service-specific configuration functions
//
// Returns the GetItemOutput containing the item and any error encountered.
//
// Example:
//
//	response, err := client.GetItem(&dynamodb.GetItemInput{...})
func (c *Client) GetItem(request *dynamodb.GetItemInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	return c.client.GetItem(c.ctx, request, optionFunctions...)
}

// PutItem creates a new item or replaces an existing item.
//
// Parameters:
//   - request: PutItemInput containing the table name and item attributes
//   - optionFunctions: optional service-specific configuration functions
//
// Returns the PutItemOutput and any error encountered.
//
// Example:
//
//	response, err := client.PutItem(&dynamodb.PutItemInput{...})
func (c *Client) PutItem(request *dynamodb.PutItemInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return c.client.PutItem(c.ctx, request, optionFunctions...)
}

// UpdateItem modifies the attributes of an existing item.
//
// Parameters:
//   - request: UpdateItemInput containing the table name, key, and update expression
//   - optionFunctions: optional service-specific configuration functions
//
// Returns the UpdateItemOutput and any error encountered.
//
// Example:
//
//	response, err := client.UpdateItem(&dynamodb.UpdateItemInput{...})
func (c *Client) UpdateItem(request *dynamodb.UpdateItemInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
	return c.client.UpdateItem(c.ctx, request, optionFunctions...)
}

// DeleteItem deletes a single item from a table.
//
// Parameters:
//   - request: DeleteItemInput containing the table name and primary key
//   - optionFunctions: optional service-specific configuration functions
//
// Returns the DeleteItemOutput and any error encountered.
//
// Example:
//
//	response, err := client.DeleteItem(&dynamodb.DeleteItemInput{...})
func (c *Client) DeleteItem(request *dynamodb.DeleteItemInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	return c.client.DeleteItem(c.ctx, request, optionFunctions...)
}

// Query retrieves items based on primary key values.
//
// You must provide the partition key value. You can optionally provide a sort key value
// and use a comparison operator to refine the search results.
//
// Parameters:
//   - request: QueryInput containing the table name and key conditions
//   - optionFunctions: optional service-specific configuration functions
//
// Returns the QueryOutput containing matching items and any error encountered.
//
// Example:
//
//	response, err := client.Query(&dynamodb.QueryInput{...})
func (c *Client) Query(request *dynamodb.QueryInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	return c.client.Query(c.ctx, request, optionFunctions...)
}

// Scan retrieves all items in a table or a secondary index.
//
// A Scan operation reads every item in a table or index. You can optionally
// apply a filter expression to return only matching items.
//
// Parameters:
//   - request: ScanInput containing the table name and optional filter expressions
//   - optionFunctions: optional service-specific configuration functions
//
// Returns the ScanOutput containing all items and any error encountered.
//
// Example:
//
//	response, err := client.Scan(&dynamodb.ScanInput{...})
func (c *Client) Scan(request *dynamodb.ScanInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	return c.client.Scan(c.ctx, request, optionFunctions...)
}

// DescribeTimeToLive retrieves the Time To Live (TTL) settings for a table.
//
// Parameters:
//   - tableName: name of the table
//   - optionFunctions: optional service-specific configuration functions
//
// Returns the DescribeTimeToLiveOutput containing TTL status and any error encountered.
//
// Example:
//
//	response, err := client.DescribeTimeToLive("users")
func (c *Client) DescribeTimeToLive(tableName string, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.DescribeTimeToLiveOutput, error) {
	return c.client.DescribeTimeToLive(
		c.ctx,
		&dynamodb.DescribeTimeToLiveInput{TableName: aws.String(tableName)},
		optionFunctions...)
}

// UpdateTimeToLive enables or disables Time To Live (TTL) for a table.
//
// TTL allows you to define a per-item timestamp to determine when an item is no longer needed.
//
// Parameters:
//   - tableName: name of the table
//   - attributeName: name of the TTL attribute
//   - enabled: true to enable TTL, false to disable
//   - optionFunctions: optional service-specific configuration functions
//
// Returns the UpdateTimeToLiveOutput and any error encountered.
//
// Example:
//
//	response, err := client.UpdateTimeToLive("users", "expirationTime", true)
func (c *Client) UpdateTimeToLive(tableName, attributeName string, enabled bool, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.UpdateTimeToLiveOutput, error) {
	return c.client.UpdateTimeToLive(
		c.ctx,
		&dynamodb.UpdateTimeToLiveInput{
			TableName: aws.String(tableName),
			TimeToLiveSpecification: &types.TimeToLiveSpecification{
				AttributeName: aws.String(attributeName),
				Enabled:       aws.Bool(enabled)}},
		optionFunctions...)
}

// QueryPaginatorNextPage fetches the next page of results using QueryPaginator.
//
// This is useful for processing large result sets that exceed the 1MB limit per query.
//
// Parameters:
//   - request: QueryInput containing the query parameters
//   - optionFunctions: optional service-specific configuration functions
//
// Returns the QueryOutput for the next page and any error encountered.
//
// Example:
//
//	response, err := client.QueryPaginatorNextPage(&dynamodb.QueryInput{...})
func (c *Client) QueryPaginatorNextPage(request *dynamodb.QueryInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	paginator := dynamodb.NewQueryPaginator(c.client, request)

	return paginator.NextPage(c.ctx, optionFunctions...)
}

// ScanPaginatorNextPage fetches the next page of results using ScanPaginator.
//
// This is useful for processing large result sets that exceed the 1MB limit per scan.
//
// Parameters:
//   - request: ScanInput containing the scan parameters
//   - optionFunctions: optional service-specific configuration functions
//
// Returns the ScanOutput for the next page and any error encountered.
//
// Example:
//
//	response, err := client.ScanPaginatorNextPage(&dynamodb.ScanInput{...})
func (c *Client) ScanPaginatorNextPage(request *dynamodb.ScanInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	paginator := dynamodb.NewScanPaginator(c.client, request)

	return paginator.NextPage(c.ctx, optionFunctions...)
}
