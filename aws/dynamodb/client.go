// Package dynamodb provides DynamoDB client implementations.
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

// CreateClient creates S3 client with service-specific options.
// This is the recommended way to create S3 client with custom endpoint.
//
// Example:
//
//	err := client.CreateClient(context.TODO(), "us-east-1", "access-key", "secret-key", "",
//	    func(o *dynamodb.Options) {
//	        o.BaseEndpoint = aws.String("http://localhost:9090")
//	        o.UsePathStyle = true
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

// CreateTable is create the table.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := client.CreateTable(&aws_dynamodb.CreateTableInput{...}, true, 10)
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

// ListTables is get the list of the table.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := client.ListTables(&aws_dynamodb.ListTablesInput{Limit: aws.Int32(10)})
func (c *Client) ListTables(request *dynamodb.ListTablesInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.ListTablesOutput, error) {
	return c.client.ListTables(c.ctx, request, optionFunctions...)
}

// DescribeTable is get the description of the table.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := client.DescribeTable("table_name")
func (c *Client) DescribeTable(tableName string, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error) {
	return c.client.DescribeTable(c.ctx, &dynamodb.DescribeTableInput{TableName: aws.String(tableName)}, optionFunctions...)
}

// UpdateTable is update the table.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := client.UpdateTable(&aws_dynamodb.UpdateTableInput{...})
func (c *Client) UpdateTable(request *dynamodb.UpdateTableInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.UpdateTableOutput, error) {
	return c.client.UpdateTable(c.ctx, request, optionFunctions...)
}

// DeleteTable is delete the table.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := client.DeleteTable("table_name", true, 10)
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

// GetItem is get the item.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := client.GetItem(&aws_dynamodb.GetItemInput{...})
func (c *Client) GetItem(request *dynamodb.GetItemInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	return c.client.GetItem(c.ctx, request, optionFunctions...)
}

// PutItem is put the item.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := client.PutItem(&aws_dynamodb.PutItemInput{...})
func (c *Client) PutItem(request *dynamodb.PutItemInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return c.client.PutItem(c.ctx, request, optionFunctions...)
}

// UpdateItem is update the item.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := client.UpdateItem(&aws_dynamodb.UpdateItemInput{...})
func (c *Client) UpdateItem(request *dynamodb.UpdateItemInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
	return c.client.UpdateItem(c.ctx, request, optionFunctions...)
}

// DeleteItem is delete the item.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := client.DeleteItem(&aws_dynamodb.DeleteItemInput{...})
func (c *Client) DeleteItem(request *dynamodb.DeleteItemInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	return c.client.DeleteItem(c.ctx, request, optionFunctions...)
}

// Query is get the items based on primary key values.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := client.Query(&aws_dynamodb.QueryInput{...})
func (c *Client) Query(request *dynamodb.QueryInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	return c.client.Query(c.ctx, request, optionFunctions...)
}

// Scan is get the every items in a table or a secondary index.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := client.Scan(&aws_dynamodb.ScanInput{...})
func (c *Client) Scan(request *dynamodb.ScanInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	return c.client.Scan(c.ctx, request, optionFunctions...)
}

// DescribeTimeToLive is get the TTL information.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := client.DescribeTimeToLive(TABLE_NAME)
func (c *Client) DescribeTimeToLive(tableName string, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.DescribeTimeToLiveOutput, error) {
	return c.client.DescribeTimeToLive(
		c.ctx,
		&dynamodb.DescribeTimeToLiveInput{TableName: aws.String(tableName)},
		optionFunctions...)
}

// UpdateTimeToLive is update the TTL information.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := client.UpdateTimeToLive(TABLE_NAME, TTL_NAME, true)
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

// QueryPaginatorNextPage is fetch the next page using QueryPaginator.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := client.QueryPaginatorNextPage(&aws_dynamodb.QueryInput{...})
func (c *Client) QueryPaginatorNextPage(request *dynamodb.QueryInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	paginator := dynamodb.NewQueryPaginator(c.client, request)

	return paginator.NextPage(c.ctx, optionFunctions...)
}

// ScanPaginatorNextPage is fetch the next page using ScanPaginator.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := client.ScanPaginatorNextPage(&aws_dynamodb.ScanInput{...})
func (c *Client) ScanPaginatorNextPage(request *dynamodb.ScanInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	paginator := dynamodb.NewScanPaginator(c.client, request)

	return paginator.NextPage(c.ctx, optionFunctions...)
}
