package dynamodb_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/localstack"

	dynamodb_client "github.com/common-library/go/aws/dynamodb"
	"github.com/common-library/go/testutil"
)

const (
	TEST_TABLE_NAME = "test-table"
	TTL_ATTRIBUTE   = "ttl"
)

type TestUser struct {
	ID       string `dynamodbav:"id"`
	Name     string `dynamodbav:"name"`
	Email    string `dynamodbav:"email"`
	Age      int    `dynamodbav:"age"`
	TTL      int64  `dynamodbav:"ttl,omitempty"`
	CreateAt string `dynamodbav:"create_at"`
}

var (
	sharedContainer *localstack.LocalStackContainer
	sharedClient    *dynamodb_client.Client
	setupOnce       sync.Once
	setupErr        error
)

func setupSharedLocalStack() (*dynamodb_client.Client, error) {
	setupOnce.Do(func() {
		ctx := context.Background()

		container, err := localstack.Run(ctx, testutil.LocalstackImage)
		if err != nil {
			setupErr = fmt.Errorf("failed to start LocalStack container: %v", err)
			return
		}
		sharedContainer = container

		mappedPort, err := container.MappedPort(ctx, "4566")
		if err != nil {
			setupErr = fmt.Errorf("failed to get LocalStack port: %v", err)
			return
		}

		host, err := container.Host(ctx)
		if err != nil {
			setupErr = fmt.Errorf("failed to get LocalStack host: %v", err)
			return
		}

		endpoint := fmt.Sprintf("%s:%s", host, mappedPort.Port())

		client := &dynamodb_client.Client{}
		err = client.CreateClient(ctx, "us-east-1", "test", "test", "",
			func(o *dynamodb.Options) {
				o.BaseEndpoint = aws.String(fmt.Sprintf("http://%s", endpoint))
			})
		if err != nil {
			setupErr = fmt.Errorf("failed to create DynamoDB client: %v", err)
			return
		}

		sharedClient = client
	})

	return sharedClient, setupErr
}

func TestMain(m *testing.M) {
	_, err := setupSharedLocalStack()
	if err != nil {
		panic(fmt.Sprintf("Failed to setup shared LocalStack: %v", err))
	}

	code := m.Run()

	if sharedContainer != nil {
		sharedContainer.Terminate(context.Background())
	}

	if code != 0 {
		panic(fmt.Sprintf("Tests failed with code: %d", code))
	}
}

func createUniqueTable(t *testing.T, client *dynamodb_client.Client, baseTableName string) string {
	tableName := fmt.Sprintf("%s_%s_%d", baseTableName, t.Name(), time.Now().UnixNano())

	_, err := client.CreateTable(&dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       types.KeyTypeHash,
			},
		},
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	}, true, 5)
	require.NoError(t, err, "Failed to create test table")

	t.Cleanup(func() {
		client.DeleteTable(tableName, false, 0)
	})

	return tableName
}

func TestClient_CreateTable(t *testing.T) {
	t.Parallel()

	client, err := setupSharedLocalStack()
	require.NoError(t, err, "Failed to setup shared LocalStack")

	tableName := fmt.Sprintf("create_test_%d", time.Now().UnixNano())

	response, err := client.CreateTable(&dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       types.KeyTypeHash,
			},
		},
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	}, true, 5)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, tableName, *response.TableDescription.TableName)

	t.Cleanup(func() {
		client.DeleteTable(tableName, false, 0)
	})
}

func TestClient_ListTables(t *testing.T) {
	t.Parallel()

	client, err := setupSharedLocalStack()
	require.NoError(t, err, "Failed to setup shared LocalStack")

	tableName := createUniqueTable(t, client, "list_test")

	response, err := client.ListTables(&dynamodb.ListTablesInput{})

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Contains(t, response.TableNames, tableName)
}

func TestClient_DescribeTable(t *testing.T) {
	t.Parallel()

	client, err := setupSharedLocalStack()
	require.NoError(t, err, "Failed to setup shared LocalStack")

	tableName := createUniqueTable(t, client, "describe_test")

	response, err := client.DescribeTable(tableName)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, tableName, *response.Table.TableName)
	assert.Equal(t, types.TableStatusActive, response.Table.TableStatus)
}

func TestClient_DeleteTable(t *testing.T) {
	t.Parallel()

	client, err := setupSharedLocalStack()
	require.NoError(t, err, "Failed to setup shared LocalStack")

	tableName := fmt.Sprintf("delete_test_%d", time.Now().UnixNano())

	_, err = client.CreateTable(&dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       types.KeyTypeHash,
			},
		},
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	}, true, 5)
	require.NoError(t, err)

	response, err := client.DeleteTable(tableName, true, 5)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, tableName, *response.TableDescription.TableName)

	_, err = client.DescribeTable(tableName)
	assert.Error(t, err)
}

func TestClient_PutAndGetItem(t *testing.T) {
	t.Parallel()

	client, err := setupSharedLocalStack()
	require.NoError(t, err, "Failed to setup shared LocalStack")

	tableName := createUniqueTable(t, client, "put_get_test")

	testUser := TestUser{
		ID:       "user-1",
		Name:     "John Doe",
		Email:    "john@example.com",
		Age:      30,
		CreateAt: time.Now().Format(time.RFC3339),
	}

	item, err := attributevalue.MarshalMap(testUser)
	require.NoError(t, err)

	putResponse, err := client.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	assert.NoError(t, err)
	assert.NotNil(t, putResponse)

	getResponse, err := client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: "user-1"},
		},
	})

	assert.NoError(t, err)
	assert.NotNil(t, getResponse)
	assert.NotEmpty(t, getResponse.Item)

	var retrievedUser TestUser
	err = attributevalue.UnmarshalMap(getResponse.Item, &retrievedUser)
	assert.NoError(t, err)
	assert.Equal(t, testUser.ID, retrievedUser.ID)
	assert.Equal(t, testUser.Name, retrievedUser.Name)
	assert.Equal(t, testUser.Email, retrievedUser.Email)
	assert.Equal(t, testUser.Age, retrievedUser.Age)
}

func TestClient_UpdateItem(t *testing.T) {
	t.Parallel()

	client, err := setupSharedLocalStack()
	require.NoError(t, err, "Failed to setup shared LocalStack")

	tableName := createUniqueTable(t, client, "update_test")

	testUser := TestUser{
		ID:       "user-1",
		Name:     "John Doe",
		Email:    "john@example.com",
		Age:      30,
		CreateAt: time.Now().Format(time.RFC3339),
	}

	item, err := attributevalue.MarshalMap(testUser)
	require.NoError(t, err)

	_, err = client.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	require.NoError(t, err)

	response, err := client.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: "user-1"},
		},
		UpdateExpression: aws.String("SET #name = :name, age = :age"),
		ExpressionAttributeNames: map[string]string{
			"#name": "name",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":name": &types.AttributeValueMemberS{Value: "Jane Doe"},
			":age":  &types.AttributeValueMemberN{Value: "31"},
		},
		ReturnValues: types.ReturnValueAllNew,
	})

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Attributes)

	var updatedUser TestUser
	err = attributevalue.UnmarshalMap(response.Attributes, &updatedUser)
	assert.NoError(t, err)
	assert.Equal(t, "Jane Doe", updatedUser.Name)
	assert.Equal(t, 31, updatedUser.Age)
}

func TestClient_DeleteItem(t *testing.T) {
	t.Parallel()

	client, err := setupSharedLocalStack()
	require.NoError(t, err, "Failed to setup shared LocalStack")

	tableName := createUniqueTable(t, client, "delete_item_test")

	testUser := TestUser{
		ID:       "user-1",
		Name:     "John Doe",
		Email:    "john@example.com",
		Age:      30,
		CreateAt: time.Now().Format(time.RFC3339),
	}

	item, err := attributevalue.MarshalMap(testUser)
	require.NoError(t, err)

	_, err = client.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	require.NoError(t, err)

	response, err := client.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: "user-1"},
		},
		ReturnValues: types.ReturnValueAllOld,
	})

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Attributes)

	getResponse, err := client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: "user-1"},
		},
	})

	assert.NoError(t, err)
	assert.Empty(t, getResponse.Item)
}

func TestClient_Scan(t *testing.T) {
	t.Parallel()

	client, err := setupSharedLocalStack()
	require.NoError(t, err, "Failed to setup shared LocalStack")

	tableName := createUniqueTable(t, client, "scan_test")

	users := []TestUser{
		{ID: "user-1", Name: "John Doe", Email: "john@example.com", Age: 30, CreateAt: time.Now().Format(time.RFC3339)},
		{ID: "user-2", Name: "Jane Smith", Email: "jane@example.com", Age: 25, CreateAt: time.Now().Format(time.RFC3339)},
		{ID: "user-3", Name: "Bob Johnson", Email: "bob@example.com", Age: 35, CreateAt: time.Now().Format(time.RFC3339)},
	}

	for _, user := range users {
		item, err := attributevalue.MarshalMap(user)
		require.NoError(t, err)

		_, err = client.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item:      item,
		})
		require.NoError(t, err)
	}

	response, err := client.Scan(&dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Len(t, response.Items, 3)
	assert.Equal(t, int32(3), response.Count)
}

func TestClient_Query(t *testing.T) {
	t.Parallel()

	client, err := setupSharedLocalStack()
	require.NoError(t, err, "Failed to setup shared LocalStack")

	tableName := fmt.Sprintf("query_test_%d", time.Now().UnixNano())
	_, err = client.CreateTable(&dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("pk"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("sk"),
				KeyType:       types.KeyTypeRange,
			},
		},
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("pk"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("sk"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	}, true, 5)
	require.NoError(t, err)

	t.Cleanup(func() {
		client.DeleteTable(tableName, false, 0)
	})

	items := []map[string]types.AttributeValue{
		{
			"pk":   &types.AttributeValueMemberS{Value: "USER#123"},
			"sk":   &types.AttributeValueMemberS{Value: "PROFILE"},
			"name": &types.AttributeValueMemberS{Value: "John Doe"},
		},
		{
			"pk":     &types.AttributeValueMemberS{Value: "USER#123"},
			"sk":     &types.AttributeValueMemberS{Value: "ORDER#001"},
			"amount": &types.AttributeValueMemberN{Value: "100"},
		},
		{
			"pk":     &types.AttributeValueMemberS{Value: "USER#123"},
			"sk":     &types.AttributeValueMemberS{Value: "ORDER#002"},
			"amount": &types.AttributeValueMemberN{Value: "200"},
		},
	}

	for _, item := range items {
		_, err = client.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item:      item,
		})
		require.NoError(t, err)
	}

	response, err := client.Query(&dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("pk = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "USER#123"},
		},
	})

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Len(t, response.Items, 3)
	assert.Equal(t, int32(3), response.Count)
}

func TestClient_TimeToLive(t *testing.T) {
	t.Parallel()

	client, err := setupSharedLocalStack()
	require.NoError(t, err, "Failed to setup shared LocalStack")

	tableName := fmt.Sprintf("ttl_test_%d", time.Now().UnixNano())
	_, err = client.CreateTable(&dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       types.KeyTypeHash,
			},
		},
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	}, true, 5)
	require.NoError(t, err)

	t.Cleanup(func() {
		client.DeleteTable(tableName, false, 0)
	})

	updateResponse, err := client.UpdateTimeToLive(tableName, TTL_ATTRIBUTE, true)
	assert.NoError(t, err)
	assert.NotNil(t, updateResponse)

	describeResponse, err := client.DescribeTimeToLive(tableName)
	assert.NoError(t, err)
	assert.NotNil(t, describeResponse)
	assert.Equal(t, types.TimeToLiveStatusEnabled, describeResponse.TimeToLiveDescription.TimeToLiveStatus)
	assert.Equal(t, TTL_ATTRIBUTE, *describeResponse.TimeToLiveDescription.AttributeName)
}

func TestClient_Pagination(t *testing.T) {
	t.Parallel()

	client, err := setupSharedLocalStack()
	require.NoError(t, err, "Failed to setup shared LocalStack")

	tableName := createUniqueTable(t, client, "pagination_test")

	for i := 0; i < 10; i++ {
		testUser := TestUser{
			ID:       fmt.Sprintf("user-%d", i),
			Name:     fmt.Sprintf("User %d", i),
			Email:    fmt.Sprintf("user%d@example.com", i),
			Age:      20 + i,
			CreateAt: time.Now().Format(time.RFC3339),
		}

		item, err := attributevalue.MarshalMap(testUser)
		require.NoError(t, err)

		_, err = client.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item:      item,
		})
		require.NoError(t, err)
	}

	response, err := client.ScanPaginatorNextPage(&dynamodb.ScanInput{
		TableName: aws.String(tableName),
		Limit:     aws.Int32(5),
	})

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Len(t, response.Items, 5)
	assert.NotNil(t, response.LastEvaluatedKey)
}

func TestClient_CRUDOperations(t *testing.T) {
	t.Parallel()

	client, err := setupSharedLocalStack()
	require.NoError(t, err, "Failed to setup shared LocalStack")

	tableName := createUniqueTable(t, client, "crud_test")

	testUser := TestUser{
		ID:       "integration-test-user",
		Name:     "Integration Test User",
		Email:    "integration@example.com",
		Age:      28,
		CreateAt: time.Now().Format(time.RFC3339),
	}

	item, err := attributevalue.MarshalMap(testUser)
	require.NoError(t, err)

	_, err = client.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	assert.NoError(t, err)

	getResponse, err := client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: testUser.ID},
		},
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, getResponse.Item)

	_, err = client.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: testUser.ID},
		},
		UpdateExpression: aws.String("SET age = :age"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":age": &types.AttributeValueMemberN{Value: "29"},
		},
	})
	assert.NoError(t, err)

	updatedResponse, err := client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: testUser.ID},
		},
	})
	assert.NoError(t, err)

	var updatedUser TestUser
	err = attributevalue.UnmarshalMap(updatedResponse.Item, &updatedUser)
	assert.NoError(t, err)
	assert.Equal(t, 29, updatedUser.Age)

	_, err = client.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: testUser.ID},
		},
	})
	assert.NoError(t, err)

	deletedResponse, err := client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: testUser.ID},
		},
	})
	assert.NoError(t, err)
	assert.Empty(t, deletedResponse.Item)
}
