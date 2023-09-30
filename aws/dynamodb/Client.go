// Package dynamodb provides DynamoDB interface.
//
// used "github.com/aws/aws-sdk-go-v2".
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

// Client is object that provides DynamoDB interface.
type Client struct {
	ctx    context.Context
	client *dynamodb.Client
}

// CreateClient is create a client.
//
// See dynamodb_test.go for a detailed example.
//
// ex) err := client.CreateClient(context.TODO(), region, accessKey, secretAccessKey, sessionToken)
func (this *Client) CreateClient(ctx context.Context, region, accessKey, secretAccessKey, sessionToken string, loadOptionFunctions ...func(*config.LoadOptions) error) error {
	loadOptionFunctions = append(loadOptionFunctions, config.WithRegion(region))
	loadOptionFunctions = append(loadOptionFunctions, config.WithCredentialsProvider(
		credentials.NewStaticCredentialsProvider(accessKey, secretAccessKey, sessionToken)))

	if cfg, err := config.LoadDefaultConfig(ctx, loadOptionFunctions...); err != nil {
		return err
	} else {
		this.ctx = ctx
		this.client = dynamodb.NewFromConfig(cfg)
	}

	return nil
}

// CreateTable is create the table.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := client.CreateTable(&aws_dynamodb.CreateTableInput{...}, true, 10)
func (this *Client) CreateTable(request *dynamodb.CreateTableInput, wait bool, waitTimeout uint64, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error) {
	response, err := this.client.CreateTable(this.ctx, request, optionFunctions...)
	if err != nil {
		return nil, err
	}

	if wait {
		waiter := dynamodb.NewTableExistsWaiter(this.client)
		if err := waiter.Wait(
			this.ctx,
			&dynamodb.DescribeTableInput{TableName: request.TableName},
			time.Duration(waitTimeout)*time.Second); err != nil {
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
func (this *Client) ListTables(request *dynamodb.ListTablesInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.ListTablesOutput, error) {
	return this.client.ListTables(this.ctx, request, optionFunctions...)
}

// DescribeTable is get the description of the table.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := client.DescribeTable("table_name")
func (this *Client) DescribeTable(tableName string, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error) {
	return this.client.DescribeTable(this.ctx, &dynamodb.DescribeTableInput{TableName: aws.String(tableName)}, optionFunctions...)
}

// UpdateTable is update the table.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := client.UpdateTable(&aws_dynamodb.UpdateTableInput{...})
func (this *Client) UpdateTable(request *dynamodb.UpdateTableInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.UpdateTableOutput, error) {
	return this.client.UpdateTable(this.ctx, request, optionFunctions...)
}

// DeleteTable is delete the table.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := client.DeleteTable("table_name", true, 10)
func (this *Client) DeleteTable(tableName string, wait bool, waitTimeout uint64, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.DeleteTableOutput, error) {
	response, err := this.client.DeleteTable(this.ctx, &dynamodb.DeleteTableInput{TableName: aws.String(tableName)}, optionFunctions...)
	if err != nil {
		return nil, err
	}

	if wait {
		waiter := dynamodb.NewTableNotExistsWaiter(this.client)
		if err := waiter.Wait(
			this.ctx,
			&dynamodb.DescribeTableInput{TableName: aws.String(tableName)},
			time.Duration(waitTimeout)*time.Second); err != nil {
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
func (this *Client) GetItem(request *dynamodb.GetItemInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	return this.client.GetItem(this.ctx, request, optionFunctions...)
}

// PutItem is put the item.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := client.PutItem(&aws_dynamodb.PutItemInput{...})
func (this *Client) PutItem(request *dynamodb.PutItemInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return this.client.PutItem(this.ctx, request, optionFunctions...)
}

// UpdateItem is update the item.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := client.UpdateItem(&aws_dynamodb.UpdateItemInput{...})
func (this *Client) UpdateItem(request *dynamodb.UpdateItemInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
	return this.client.UpdateItem(this.ctx, request, optionFunctions...)
}

// DeleteItem is delete the item.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := client.DeleteItem(&aws_dynamodb.DeleteItemInput{...})
func (this *Client) DeleteItem(request *dynamodb.DeleteItemInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	return this.client.DeleteItem(this.ctx, request, optionFunctions...)
}

// Query is get the items based on primary key values.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := client.Query(&aws_dynamodb.QueryInput{...})
func (this *Client) Query(request *dynamodb.QueryInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	return this.client.Query(this.ctx, request, optionFunctions...)
}

// Scan is get the every items in a table or a secondary index.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := client.Scan(&aws_dynamodb.ScanInput{...})
func (this *Client) Scan(request *dynamodb.ScanInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	return this.client.Scan(this.ctx, request, optionFunctions...)
}

// DescribeTimeToLive is get the TTL information.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := client.DescribeTimeToLive(TABLE_NAME)
func (this *Client) DescribeTimeToLive(tableName string, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.DescribeTimeToLiveOutput, error) {
	return this.client.DescribeTimeToLive(
		this.ctx,
		&dynamodb.DescribeTimeToLiveInput{TableName: aws.String(tableName)},
		optionFunctions...)
}

// UpdateTimeToLive is update the TTL information.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := client.UpdateTimeToLive(TABLE_NAME, TTL_NAME, true)
func (this *Client) UpdateTimeToLive(tableName, attributeName string, enabled bool, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.UpdateTimeToLiveOutput, error) {
	return this.client.UpdateTimeToLive(
		this.ctx,
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
func (this *Client) QueryPaginatorNextPage(request *dynamodb.QueryInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	paginator := dynamodb.NewQueryPaginator(this.client, request)

	return paginator.NextPage(this.ctx, optionFunctions...)
}

// ScanPaginatorNextPage is fetch the next page using ScanPaginator.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := client.ScanPaginatorNextPage(&aws_dynamodb.ScanInput{...})
func (this *Client) ScanPaginatorNextPage(request *dynamodb.ScanInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	paginator := dynamodb.NewScanPaginator(this.client, request)

	return paginator.NextPage(this.ctx, optionFunctions...)
}
