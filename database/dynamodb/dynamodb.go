// Package dynamodb provides DynamoDB interface.
//
// used "github.com/aws/aws-sdk-go-v2".
package dynamodb

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// DynamoDB is object that provides DynamoDB interface.
type DynamoDB struct {
	ctx    context.Context
	client *dynamodb.Client
}

// CreateClient is create a client.
//
// See dynamodb_test.go for a detailed example.
//
// ex) err := dynamoDB.CreateClient(context.TODO(), ...)
func (this *DynamoDB) CreateClient(ctx context.Context, optionFunctions ...func(*config.LoadOptions) error) error {
	cfg, err := config.LoadDefaultConfig(ctx, optionFunctions...)
	if err != nil {
		return err
	}

	this.ctx = ctx
	this.client = dynamodb.NewFromConfig(cfg)

	return nil
}

// CreateTable is create the table.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := dynamoDB.CreateTable(&aws_dynamodb.CreateTableInput{...}, true, 10)
func (this *DynamoDB) CreateTable(request *dynamodb.CreateTableInput, wait bool, waitTimeout uint64, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error) {
	response, err := this.client.CreateTable(this.ctx, request, optionFunctions...)
	if err != nil {
		return nil, err
	}

	if wait {
		waiter := dynamodb.NewTableExistsWaiter(this.client)
		err := waiter.Wait(this.ctx, &dynamodb.DescribeTableInput{
			TableName: request.TableName}, time.Duration(waitTimeout)*time.Second)
		if err != nil {
			return nil, err
		}
	}

	return response, nil
}

// ListTables is get the list of the table.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := dynamoDB.ListTables(&aws_dynamodb.ListTablesInput{Limit: aws.Int32(10)})
func (this *DynamoDB) ListTables(request *dynamodb.ListTablesInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.ListTablesOutput, error) {
	return this.client.ListTables(this.ctx, request, optionFunctions...)
}

// DescribeTable is get the description of the table.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := dynamoDB.DescribeTable("table_name")
func (this *DynamoDB) DescribeTable(tableName string, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error) {
	return this.client.DescribeTable(this.ctx, &dynamodb.DescribeTableInput{TableName: aws.String(tableName)}, optionFunctions...)
}

// UpdateTable is update the table.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := dynamoDB.UpdateTable(&aws_dynamodb.UpdateTableInput{...})
func (this *DynamoDB) UpdateTable(request *dynamodb.UpdateTableInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.UpdateTableOutput, error) {
	return this.client.UpdateTable(this.ctx, request, optionFunctions...)
}

// DeleteTable is delete the table.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := dynamoDB.DeleteTable("table_name", true, 10)
func (this *DynamoDB) DeleteTable(tableName string, wait bool, waitTimeout uint64, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.DeleteTableOutput, error) {
	response, err := this.client.DeleteTable(this.ctx, &dynamodb.DeleteTableInput{TableName: aws.String(tableName)}, optionFunctions...)
	if err != nil {
		return nil, err
	}

	if wait {
		waiter := dynamodb.NewTableNotExistsWaiter(this.client)
		err := waiter.Wait(this.ctx, &dynamodb.DescribeTableInput{
			TableName: aws.String(tableName)}, time.Duration(waitTimeout)*time.Second)
		if err != nil {
			return nil, err
		}
	}

	return response, nil
}

// GetItem is get the item.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := dynamoDB.GetItem(&aws_dynamodb.GetItemInput{...})
func (this *DynamoDB) GetItem(request *dynamodb.GetItemInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	return this.client.GetItem(this.ctx, request, optionFunctions...)
}

// PutItem is put the item.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := dynamoDB.PutItem(&aws_dynamodb.PutItemInput{...})
func (this *DynamoDB) PutItem(request *dynamodb.PutItemInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return this.client.PutItem(this.ctx, request, optionFunctions...)
}

// UpdateItem is update the item.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := dynamoDB.UpdateItem(&aws_dynamodb.UpdateItemInput{...})
func (this *DynamoDB) UpdateItem(request *dynamodb.UpdateItemInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
	return this.client.UpdateItem(this.ctx, request, optionFunctions...)
}

// DeleteItem is delete the item.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := dynamoDB.DeleteItem(&aws_dynamodb.DeleteItemInput{...})
func (this *DynamoDB) DeleteItem(request *dynamodb.DeleteItemInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	return this.client.DeleteItem(this.ctx, request, optionFunctions...)
}

// Query is get the items based on primary key values.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := dynamoDB.Query(&aws_dynamodb.QueryInput{...})
func (this *DynamoDB) Query(request *dynamodb.QueryInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	return this.client.Query(this.ctx, request, optionFunctions...)
}

// Scan is get the every items in a table or a secondary index.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := dynamoDB.Scan(&aws_dynamodb.ScanInput{...})
func (this *DynamoDB) Scan(request *dynamodb.ScanInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	return this.client.Scan(this.ctx, request, optionFunctions...)
}

// DescribeTimeToLive is get the TTL information.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := dynamoDB.DescribeTimeToLive(TABLE_NAME)
func (this *DynamoDB) DescribeTimeToLive(tableName string, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.DescribeTimeToLiveOutput, error) {
	return this.client.DescribeTimeToLive(this.ctx, &dynamodb.DescribeTimeToLiveInput{TableName: aws.String(tableName)}, optionFunctions...)
}

// UpdateTimeToLive is update the TTL information.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := dynamoDB.UpdateTimeToLive(TABLE_NAME, TTL_NAME, true)
func (this *DynamoDB) UpdateTimeToLive(tableName, attributeName string, enabled bool, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.UpdateTimeToLiveOutput, error) {
	return this.client.UpdateTimeToLive(this.ctx, &dynamodb.UpdateTimeToLiveInput{TableName: aws.String(tableName), TimeToLiveSpecification: &types.TimeToLiveSpecification{AttributeName: aws.String(attributeName), Enabled: aws.Bool(enabled)}}, optionFunctions...)
}

// QueryPaginatorNextPage is fetch the next page using QueryPaginator.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := dynamoDB.QueryPaginatorNextPage(&aws_dynamodb.QueryInput{...})
func (this *DynamoDB) QueryPaginatorNextPage(request *dynamodb.QueryInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	paginator := dynamodb.NewQueryPaginator(this.client, request)

	return paginator.NextPage(this.ctx, optionFunctions...)
}

// ScanPaginatorNextPage is fetch the next page using ScanPaginator.
//
// See dynamodb_test.go for a detailed example.
//
// ex) response, err := dynamoDB.ScanPaginatorNextPage(&aws_dynamodb.ScanInput{...})
func (this *DynamoDB) ScanPaginatorNextPage(request *dynamodb.ScanInput, optionFunctions ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	paginator := dynamodb.NewScanPaginator(this.client, request)

	return paginator.NextPage(this.ctx, optionFunctions...)
}
