package dynamodb_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	aws_dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-library/go/aws/dynamodb"
	"github.com/google/uuid"
)

var TABLE_NAME = strings.ReplaceAll(uuid.NewString(), "-", "")

const INDEX_NAME = "field1-index"
const TTL_NAME = "ttl-test"

type TestItem struct {
	PrimaryKey int    `dynamodbav:"primary-key"`
	SortKey    string `dynamodbav:"sort-key"`

	Field1 bool   `dynamodbav:"field1"`
	Field2 int    `dynamodbav:"field2"`
	Field3 string `dynamodbav:"field3"`

	Field4 []struct {
		SubField1 string `dynamodbav:"sub-field1"`
		SubField2 string `dynamodbav:"sub-field2"`
	} `dynamodbav:"field4"`

	TTLTest int64 `dynamodbav:"ttl-test,omitempty"`
}

func (this *TestItem) getKey() (map[string]types.AttributeValue, error) {
	pk, err := attributevalue.Marshal(this.PrimaryKey)
	if err != nil {
		return nil, err
	}

	sk, err := attributevalue.Marshal(this.SortKey)
	if err != nil {
		return nil, err
	}

	return map[string]types.AttributeValue{"primary-key": pk, "sort-key": sk}, nil
}

func initialize(client *dynamodb.Client, putItems bool, t *testing.T) {
	if err := client.CreateClient(
		context.TODO(), "dummy", "dummy", "dummy", "dummy",
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{URL: fmt.Sprintf("http://127.0.0.1:8000")}, nil
			}))); err != nil {
		t.Fatal(err)
	}

	if response, err := client.CreateTable(
		&aws_dynamodb.CreateTableInput{
			TableName: aws.String(TABLE_NAME),
			AttributeDefinitions: []types.AttributeDefinition{{
				AttributeName: aws.String("primary-key"),
				AttributeType: types.ScalarAttributeTypeN,
			}, {
				AttributeName: aws.String("sort-key"),
				AttributeType: types.ScalarAttributeTypeS,
			}},
			KeySchema: []types.KeySchemaElement{{
				AttributeName: aws.String("primary-key"),
				KeyType:       types.KeyTypeHash,
			}, {
				AttributeName: aws.String("sort-key"),
				KeyType:       types.KeyTypeRange,
			}},
			ProvisionedThroughput: &types.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(10),
				WriteCapacityUnits: aws.Int64(10),
			},
			//BillingMode: types.BillingModePayPerRequest,
		},
		true, 10); err != nil {
		t.Fatal(err)
	} else if *response.TableDescription.TableName != TABLE_NAME {
		t.Fatalf("invalid TableName - (%s)(%s)", *response.TableDescription.TableName, TABLE_NAME)
	} else if response.TableDescription.TableStatus != types.TableStatusActive {
		t.Fatalf("invalid TableStatus - (%s)", response.TableDescription.TableStatus)
	}

	{
		if putItems == false {
			return
		}

		testItems := []TestItem{
			{PrimaryKey: 1, SortKey: "a", Field1: true, Field2: 1, Field3: "value_for_1",
				Field4: []struct {
					SubField1 string `dynamodbav:"sub-field1"`
					SubField2 string `dynamodbav:"sub-field2"`
				}{
					{SubField1: "sub-value1-1_for_1", SubField2: "sub-value1-2_for_2"},
					{SubField1: "sub-value2-1_for_1", SubField2: "sub-value2-2_for_1"}},
				TTLTest: time.Now().Unix() + int64(10)},

			{PrimaryKey: 2, SortKey: "b", Field1: false, Field2: 2, Field3: "value_for_2",
				Field4: []struct {
					SubField1 string `dynamodbav:"sub-field1"`
					SubField2 string `dynamodbav:"sub-field2"`
				}{
					{SubField1: "sub-value1-1_for_2", SubField2: "sub-value1-2_for_2"},
					{SubField1: "sub-value2-1_for_2", SubField2: "sub-value2-2_for_2"}}},

			{PrimaryKey: 3, SortKey: "c-1", Field1: true, Field2: 31, Field3: "value_for_3",
				Field4: []struct {
					SubField1 string `dynamodbav:"sub-field1"`
					SubField2 string `dynamodbav:"sub-field2"`
				}{
					{SubField1: "sub-value1-1_for_3", SubField2: "sub-value1-2_for_3"},
					{SubField1: "sub-value2-1_for_3", SubField2: "sub-value2-2_for_3"}}},

			{PrimaryKey: 3, SortKey: "c-2", Field1: true, Field2: 32, Field3: "value_for_3",
				Field4: []struct {
					SubField1 string `dynamodbav:"sub-field1"`
					SubField2 string `dynamodbav:"sub-field2"`
				}{
					{SubField1: "sub-value1-1_for_3", SubField2: "sub-value1-2_for_3"},
					{SubField1: "sub-value2-1_for_3", SubField2: "sub-value2-2_for_3"}}},

			{PrimaryKey: 4, SortKey: "d", Field1: true, Field2: 4, Field3: "value_for_4",
				Field4: []struct {
					SubField1 string `dynamodbav:"sub-field1"`
					SubField2 string `dynamodbav:"sub-field2"`
				}{
					{SubField1: "sub-value1-1_for_4", SubField2: "sub-value1-2_for_4"},
					{SubField1: "sub-value2-1_for_4", SubField2: "sub-value2-2_for_4"}}},
		}

		for _, testItem := range testItems {
			if item, err := attributevalue.MarshalMap(testItem); err != nil {
				t.Fatal(err)
			} else if _, err = client.PutItem(
				&aws_dynamodb.PutItemInput{
					TableName: aws.String(TABLE_NAME), Item: item,
				}); err != nil {
				t.Fatal(err)
			}
		}
	}
}

func finalize(client *dynamodb.Client, t *testing.T) {
	if response, err := client.DeleteTable(TABLE_NAME, true, 10); err != nil {
		t.Fatal(err)
	} else if *response.TableDescription.TableName != TABLE_NAME {
		t.Fatalf("invalid TableName - (%s)(%s)", *response.TableDescription.TableName, TABLE_NAME)
	} else if response.TableDescription.TableStatus != types.TableStatusActive {
		t.Fatalf("invalid TableStatus - (%s)", response.TableDescription.TableStatus)
	}
}

func TestCreateClient(t *testing.T) {
	client := dynamodb.Client{}

	initialize(&client, false, t)
	defer finalize(&client, t)
}

func TestCreateTable(t *testing.T) {
	client := dynamodb.Client{}

	initialize(&client, false, t)
	defer finalize(&client, t)
}

func TestListTables(t *testing.T) {
	client := dynamodb.Client{}

	initialize(&client, false, t)
	defer finalize(&client, t)

	response, err := client.ListTables(&aws_dynamodb.ListTablesInput{Limit: aws.Int32(10)})
	if err != nil {
		t.Fatal(err)
	}

	exist := false
	for _, name := range response.TableNames {
		if name == TABLE_NAME {
			exist = true
		}
	}

	if exist == false {
		t.Fatalf("invalid ListTables - (%#v)", response.TableNames)
	}
}

func TestDescribeTable(t *testing.T) {
	client := dynamodb.Client{}

	initialize(&client, false, t)
	defer finalize(&client, t)

	if response, err := client.DescribeTable(TABLE_NAME); err != nil {
		t.Fatal(err)
	} else if *response.Table.TableName != TABLE_NAME {
		t.Fatalf("invalid TableName - (%s)(%s)", *response.Table.TableName, TABLE_NAME)
	} else if response.Table.TableStatus != types.TableStatusActive {
		t.Fatalf("invalid TableStatus - (%s)", response.Table.TableStatus)
	}
}

func TestUpdateTable(t *testing.T) {
	client := dynamodb.Client{}

	initialize(&client, false, t)
	defer finalize(&client, t)

	if response, err := client.DescribeTable(TABLE_NAME); err != nil {
		t.Fatal(err)
	} else if len(response.Table.GlobalSecondaryIndexes) != 0 {
		for _, index := range response.Table.GlobalSecondaryIndexes {
			t.Log(*index.IndexName)
		}

		t.Fatalf("invalid indexes size - (%d)", len(response.Table.GlobalSecondaryIndexes))
	}

	if response, err := client.UpdateTable(
		&aws_dynamodb.UpdateTableInput{
			TableName: aws.String(TABLE_NAME),
			AttributeDefinitions: []types.AttributeDefinition{{
				AttributeName: aws.String("field1"),
				AttributeType: types.ScalarAttributeTypeN,
			}},
			GlobalSecondaryIndexUpdates: []types.GlobalSecondaryIndexUpdate{
				{
					Create: &types.CreateGlobalSecondaryIndexAction{
						IndexName: aws.String(INDEX_NAME),
						KeySchema: []types.KeySchemaElement{
							{
								AttributeName: aws.String("field1"),
								KeyType:       types.KeyTypeHash,
							},
						},
						Projection: &types.Projection{ProjectionType: types.ProjectionTypeAll},
					},
				},
			},
			BillingMode: types.BillingModePayPerRequest,
			/*
			   ProvisionedThroughput: &types.ProvisionedThroughput{
			       ReadCapacityUnits:  aws.Int64(10),
			       WriteCapacityUnits: aws.Int64(10),
			   },
			*/
		}); err != nil {
		t.Fatal(err)
	} else if *response.TableDescription.GlobalSecondaryIndexes[0].IndexName != INDEX_NAME {
		t.Fatalf("invalid IndexName - (%s)(%s)", *response.TableDescription.GlobalSecondaryIndexes[0].IndexName, INDEX_NAME)
	}
}

func TestDeleteTable(t *testing.T) {
	client := dynamodb.Client{}

	initialize(&client, false, t)
	defer finalize(&client, t)
}

func TestGetItem(t *testing.T) {
	client := dynamodb.Client{}

	initialize(&client, false, t)
	defer finalize(&client, t)

	testItemForPut := TestItem{
		PrimaryKey: 1, SortKey: "a", Field1: false, Field2: 1,
		Field3: "value_for_1", Field4: []struct {
			SubField1 string `dynamodbav:"sub-field1"`
			SubField2 string `dynamodbav:"sub-field2"`
		}{
			{SubField1: "sub-value1-1_for_1", SubField2: "sub-value1-2_for_2"},
			{SubField1: "sub-value2-1_for_1", SubField2: "sub-value2-2_for_1"},
		}}

	{
		if item, err := attributevalue.MarshalMap(testItemForPut); err != nil {
			t.Fatal(err)
		} else if _, err = client.PutItem(
			&aws_dynamodb.PutItemInput{
				TableName: aws.String(TABLE_NAME),
				Item:      item,
			}); err != nil {
			t.Fatal(err)
		}
	}

	{
		testItemForGet := TestItem{PrimaryKey: 1, SortKey: "a"}

		key, err := testItemForGet.getKey()
		if err != nil {
			t.Fatal(err)
		}

		response, err := client.GetItem(&aws_dynamodb.GetItemInput{
			TableName: aws.String(TABLE_NAME), Key: key})
		if err != nil {
			t.Fatal(err)
		}

		if err := attributevalue.UnmarshalMap(response.Item, &testItemForGet); err != nil {
			t.Fatal(err)
		} else if testItemForGet.Field1 != testItemForPut.Field1 ||
			testItemForGet.Field2 != testItemForPut.Field2 ||
			testItemForGet.Field3 != testItemForPut.Field3 ||
			testItemForGet.Field4[0].SubField1 != testItemForPut.Field4[0].SubField1 ||
			testItemForGet.Field4[0].SubField2 != testItemForPut.Field4[0].SubField2 ||
			testItemForGet.Field4[1].SubField1 != testItemForPut.Field4[1].SubField1 ||
			testItemForGet.Field4[1].SubField2 != testItemForPut.Field4[1].SubField2 {
			t.Logf("testItemForGet : (%#v)", testItemForGet)
			t.Logf("testItemForPut : (%#v)", testItemForPut)
			t.Fatal("invalid GetItem")
		}
	}
}

func TestPutItem(t *testing.T) {
	client := dynamodb.Client{}

	initialize(&client, false, t)
	defer finalize(&client, t)

	testItemForPut := TestItem{
		PrimaryKey: 1, SortKey: "a", Field1: false, Field2: 1, Field3: "value_for_1",
		Field4: []struct {
			SubField1 string `dynamodbav:"sub-field1"`
			SubField2 string `dynamodbav:"sub-field2"`
		}{
			{SubField1: "sub-value1-1_for_1", SubField2: "sub-value1-2_for_2"},
			{SubField1: "sub-value2-1_for_1", SubField2: "sub-value2-2_for_1"}}}

	{
		if item, err := attributevalue.MarshalMap(testItemForPut); err != nil {
			t.Fatal(err)
		} else if _, err = client.PutItem(
			&aws_dynamodb.PutItemInput{
				TableName: aws.String(TABLE_NAME),
				Item:      item}); err != nil {
			t.Fatal(err)
		}
	}

	{
		testItemForGet := TestItem{PrimaryKey: 1, SortKey: "a"}

		key, err := testItemForGet.getKey()
		if err != nil {
			t.Fatal(err)
		}

		response, err := client.GetItem(
			&aws_dynamodb.GetItemInput{TableName: aws.String(TABLE_NAME), Key: key})
		if err != nil {
			t.Fatal(err)
		}

		if err := attributevalue.UnmarshalMap(response.Item, &testItemForGet); err != nil {
			t.Fatal(err)
		} else if testItemForGet.Field1 != testItemForPut.Field1 ||
			testItemForGet.Field2 != testItemForPut.Field2 ||
			testItemForGet.Field3 != testItemForPut.Field3 ||
			testItemForGet.Field4[0].SubField1 != testItemForPut.Field4[0].SubField1 ||
			testItemForGet.Field4[0].SubField2 != testItemForPut.Field4[0].SubField2 ||
			testItemForGet.Field4[1].SubField1 != testItemForPut.Field4[1].SubField1 ||
			testItemForGet.Field4[1].SubField2 != testItemForPut.Field4[1].SubField2 {
			t.Logf("testItemForGet : (%#v)", testItemForGet)
			t.Logf("testItemForPut : (%#v)", testItemForPut)
			t.Fatal("invalid GetItem")
		}
	}
}

func TestUpdateItem(t *testing.T) {
	const updateValue = 10

	client := dynamodb.Client{}

	initialize(&client, false, t)
	defer finalize(&client, t)

	{
		testItem := TestItem{
			PrimaryKey: 1, SortKey: "a", Field1: false, Field2: 1, Field3: "value_for_1",
			Field4: []struct {
				SubField1 string `dynamodbav:"sub-field1"`
				SubField2 string `dynamodbav:"sub-field2"`
			}{
				{SubField1: "sub-value1-1_for_1", SubField2: "sub-value1-2_for_2"},
				{SubField1: "sub-value2-1_for_1", SubField2: "sub-value2-2_for_1"}}}

		if item, err := attributevalue.MarshalMap(testItem); err != nil {
			t.Fatal(err)
		} else if _, err := client.PutItem(
			&aws_dynamodb.PutItemInput{
				TableName: aws.String(TABLE_NAME),
				Item:      item}); err != nil {
			t.Fatal(err)
		}
	}

	{
		testItem := TestItem{PrimaryKey: 1, SortKey: "a"}

		key, err := testItem.getKey()
		if err != nil {
			t.Fatal(err)
		}

		update := expression.Set(expression.Name("field2"), expression.Value(aws.Int(updateValue)))
		expr, err := expression.NewBuilder().WithUpdate(update).Build()
		if err != nil {
			t.Fatal(err)
		}

		if _, err = client.UpdateItem(
			&aws_dynamodb.UpdateItemInput{
				TableName:                 aws.String(TABLE_NAME),
				Key:                       key,
				ExpressionAttributeNames:  expr.Names(),
				ExpressionAttributeValues: expr.Values(),
				UpdateExpression:          expr.Update(),
				ReturnValues:              types.ReturnValueUpdatedNew,
			}); err != nil {
			t.Fatal(err)
		}
	}

	{
		testItem := TestItem{PrimaryKey: 1, SortKey: "a"}

		key, err := testItem.getKey()
		if err != nil {
			t.Fatal(err)
		}

		response, err := client.GetItem(
			&aws_dynamodb.GetItemInput{
				TableName: aws.String(TABLE_NAME),
				Key:       key})
		if err != nil {
			t.Fatal(err)
		}

		if err := attributevalue.UnmarshalMap(response.Item, &testItem); err != nil {
			t.Fatal(err)
		} else if testItem.Field2 != updateValue {
			t.Fatalf("invalid field2 - (%d)(%d)", testItem.Field2, updateValue)
		}
	}
}

func TestDeleteItem(t *testing.T) {
	client := dynamodb.Client{}

	initialize(&client, false, t)
	defer finalize(&client, t)

	{
		testItem := TestItem{
			PrimaryKey: 1, SortKey: "a", Field1: false, Field2: 1, Field3: "value_for_1",
			Field4: []struct {
				SubField1 string `dynamodbav:"sub-field1"`
				SubField2 string `dynamodbav:"sub-field2"`
			}{
				{SubField1: "sub-value1-1_for_1", SubField2: "sub-value1-2_for_2"},
				{SubField1: "sub-value2-1_for_1", SubField2: "sub-value2-2_for_1"}}}

		if item, err := attributevalue.MarshalMap(testItem); err != nil {
			t.Fatal(err)
		} else if _, err := client.PutItem(
			&aws_dynamodb.PutItemInput{
				TableName: aws.String(TABLE_NAME),
				Item:      item}); err != nil {
			t.Fatal(err)
		}
	}

	{
		testItem := TestItem{PrimaryKey: 1, SortKey: "a"}
		if key, err := testItem.getKey(); err != nil {
			t.Fatal(err)
		} else if _, err := client.DeleteItem(
			&aws_dynamodb.DeleteItemInput{
				TableName: aws.String(TABLE_NAME),
				Key:       key}); err != nil {
			t.Fatal(err)
		}
	}

	{
		if response, err := client.DescribeTable(TABLE_NAME); err != nil {
			t.Fatal(err)
		} else if *response.Table.ItemCount != 0 {
			t.Fatalf("invalid ItemCount - (%d)", *response.Table.ItemCount)
		}
	}
}

func TestQuery(t *testing.T) {
	client := dynamodb.Client{}

	initialize(&client, true, t)
	defer finalize(&client, t)

	keyEx := expression.Key("primary-key").Equal(expression.Value(3)).And(expression.Key("sort-key").Equal(expression.Value("c-1")))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		t.Fatal(err)
	}

	response, err := client.Query(
		&aws_dynamodb.QueryInput{
			TableName:                 aws.String(TABLE_NAME),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			KeyConditionExpression:    expr.KeyCondition()})
	if err != nil {
		t.Fatal(err)
	}

	testItems := []TestItem{}
	if err := attributevalue.UnmarshalListOfMaps(response.Items, &testItems); err != nil {
		t.Fatal(err)
	} else if len(testItems) != 1 ||
		testItems[0].Field1 != true ||
		testItems[0].Field2 != 31 ||
		testItems[0].Field3 != "value_for_3" ||
		testItems[0].Field4[0].SubField1 != "sub-value1-1_for_3" ||
		testItems[0].Field4[0].SubField2 != "sub-value1-2_for_3" ||
		testItems[0].Field4[1].SubField1 != "sub-value2-1_for_3" ||
		testItems[0].Field4[1].SubField2 != "sub-value2-2_for_3" {
		t.Fatalf("invalid Query - (%#v)", testItems)
	}
}

func TestScan(t *testing.T) {
	client := dynamodb.Client{}

	initialize(&client, true, t)
	defer finalize(&client, t)

	filtEx := expression.Name("primary-key").Between(expression.Value(2), expression.Value(3))
	projEx := expression.NamesList(expression.Name("primary-key"), expression.Name("sort-key"), expression.Name("field2"))
	expr, err := expression.NewBuilder().WithFilter(filtEx).WithProjection(projEx).Build()

	response, err := client.Scan(
		&aws_dynamodb.ScanInput{
			TableName:                 aws.String(TABLE_NAME),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			FilterExpression:          expr.Filter(),
			ProjectionExpression:      expr.Projection()})
	if err != nil {
		t.Fatal(err)
	}

	testItems := []TestItem{}
	if err := attributevalue.UnmarshalListOfMaps(response.Items, &testItems); err != nil {
		t.Fatal(err)
	} else if len(testItems) != 3 ||
		testItems[0].Field2 != 2 ||
		testItems[1].Field2 != 31 ||
		testItems[2].Field2 != 32 {
		t.Fatalf("invalid Scan - (%#v)", testItems)
	}
}

func TestDescribeTimeToLive(t *testing.T) {
	client := dynamodb.Client{}

	initialize(&client, true, t)
	defer finalize(&client, t)

	if response, err := client.DescribeTimeToLive(TABLE_NAME); err != nil {
		t.Fatal(err)
	} else if response.TimeToLiveDescription.AttributeName != nil {
		t.Fatal("invalid DescribeTimeToLive")
	}

	if response, err := client.UpdateTimeToLive(TABLE_NAME, TTL_NAME, true); err != nil {
		t.Fatal(err)
	} else if *response.TimeToLiveSpecification.AttributeName != TTL_NAME {
		t.Fatalf("invalid AttributeName - (%s)", *response.TimeToLiveSpecification.AttributeName)
	} else if *response.TimeToLiveSpecification.Enabled != true {
		t.Fatalf("invalid Enabled - (%t)", *response.TimeToLiveSpecification.Enabled)
	}

	if response, err := client.DescribeTimeToLive(TABLE_NAME); err != nil {
		t.Fatal(err)
	} else if *response.TimeToLiveDescription.AttributeName != TTL_NAME {
		t.Fatalf("invalid AttributeName - (%s)", *response.TimeToLiveDescription.AttributeName)
	} else if response.TimeToLiveDescription.TimeToLiveStatus != types.TimeToLiveStatusEnabled {
		t.Fatalf("invalid TimeToLiveStatus - (%s)", response.TimeToLiveDescription.TimeToLiveStatus)
	}
}

func TestUpdateTimeToLive(t *testing.T) {
	client := dynamodb.Client{}

	initialize(&client, true, t)
	defer finalize(&client, t)

	if response, err := client.DescribeTimeToLive(TABLE_NAME); err != nil {
		t.Fatal(err)
	} else if response.TimeToLiveDescription.AttributeName != nil {
		t.Fatal("invalid DescribeTimeToLive")
	}

	if response, err := client.UpdateTimeToLive(TABLE_NAME, TTL_NAME, true); err != nil {
		t.Fatal(err)
	} else if *response.TimeToLiveSpecification.AttributeName != TTL_NAME {
		t.Fatalf("invalid AttributeName - (%s)", *response.TimeToLiveSpecification.AttributeName)
	} else if *response.TimeToLiveSpecification.Enabled != true {
		t.Fatalf("invalid Enabled - (%t)", *response.TimeToLiveSpecification.Enabled)
	}

	if response, err := client.DescribeTimeToLive(TABLE_NAME); err != nil {
		t.Fatal(err)
	} else if *response.TimeToLiveDescription.AttributeName != TTL_NAME {
		t.Fatalf("invalid AttributeName - (%s)", *response.TimeToLiveDescription.AttributeName)
	} else if response.TimeToLiveDescription.TimeToLiveStatus != types.TimeToLiveStatusEnabled {
		t.Fatalf("invalid TimeToLiveStatus - (%s)", response.TimeToLiveDescription.TimeToLiveStatus)
	}
}

func TestQueryPaginatorNextPage(t *testing.T) {
	client := dynamodb.Client{}

	initialize(&client, true, t)
	defer finalize(&client, t)

	keyEx := expression.Key("primary-key").Equal(expression.Value(3))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		t.Fatal(err)
	}

	if response, err := client.QueryPaginatorNextPage(
		&aws_dynamodb.QueryInput{
			TableName:                 aws.String(TABLE_NAME),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			KeyConditionExpression:    expr.KeyCondition(),
			Limit:                     aws.Int32(1),
			ExclusiveStartKey:         nil,
		}); err != nil {
		t.Fatal(err)
	} else if response.Count != 1 || response.ScannedCount != 1 ||
		len(response.Items) != 1 || response.LastEvaluatedKey == nil {
		t.Fatalf("invalid response - (%#v)", response)
	}
}

func TestScanPaginatorNextPage(t *testing.T) {
	client := dynamodb.Client{}

	initialize(&client, true, t)
	defer finalize(&client, t)

	if response, err := client.ScanPaginatorNextPage(
		&aws_dynamodb.ScanInput{
			TableName:         aws.String(TABLE_NAME),
			Limit:             aws.Int32(2),
			ExclusiveStartKey: nil,
		}); err != nil {
		t.Fatal(err)
	} else if response.Count != 2 || response.ScannedCount != 2 ||
		len(response.Items) != 2 || response.LastEvaluatedKey == nil {
		t.Fatalf("invalid response - (%#v)", response)
	}
}
