package dynamodb_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	aws_dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-library/go/aws/dynamodb"
)

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
	if pk, err := attributevalue.Marshal(this.PrimaryKey); err != nil {
		return nil, err
	} else if sk, err := attributevalue.Marshal(this.SortKey); err != nil {
		return nil, err
	} else {
		return map[string]types.AttributeValue{"primary-key": pk, "sort-key": sk}, nil
	}
}

func initialize(t *testing.T, client *dynamodb.Client, putItems bool) bool {
	t.Parallel()

	if len(os.Getenv("DYNAMODB_URL")) == 0 {
		return false
	}

	if err := client.CreateClient(
		context.TODO(), "dummy", "dummy", "dummy", "dummy",
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{URL: os.Getenv("DYNAMODB_URL")}, nil
			}))); err != nil {
		t.Fatal(err)
	}

	if response, err := client.CreateTable(
		&aws_dynamodb.CreateTableInput{
			TableName: aws.String(t.Name()),
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
	} else if *response.TableDescription.TableName != t.Name() {
		t.Fatal(*response.TableDescription.TableName, ",", t.Name())
	} else if response.TableDescription.TableStatus != types.TableStatusActive {
		t.Fatal(response.TableDescription.TableStatus)
	}

	{
		if putItems == false {
			return true
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
					TableName: aws.String(t.Name()), Item: item,
				}); err != nil {
				t.Fatal(err)
			}
		}
	}

	return true
}

func finalize(t *testing.T, client *dynamodb.Client) {
	if response, err := client.DeleteTable(t.Name(), true, 10); err != nil {
		t.Fatal(err)
	} else if *response.TableDescription.TableName != t.Name() {
		t.Fatal(*response.TableDescription.TableName, ",", t.Name())
	} else if response.TableDescription.TableStatus != types.TableStatusActive {
		t.Fatal(response.TableDescription.TableStatus)
	}
}

func TestCreateClient(t *testing.T) {
	client := dynamodb.Client{}

	if initialize(t, &client, false) == false {
		return
	}
	defer finalize(t, &client)
}

func TestCreateTable(t *testing.T) {
	client := dynamodb.Client{}

	if initialize(t, &client, false) == false {
		return
	}
	defer finalize(t, &client)
}

func TestListTables(t *testing.T) {
	client := dynamodb.Client{}

	if initialize(t, &client, false) == false {
		return
	}
	defer finalize(t, &client)

	response, err := client.ListTables(&aws_dynamodb.ListTablesInput{Limit: aws.Int32(10)})
	if err != nil {
		t.Fatal(err)
	}

	exist := false
	for _, name := range response.TableNames {
		if name == t.Name() {
			exist = true
		}
	}

	if exist == false {
		t.Fatal(response.TableNames, ",", t.Name())
	}
}

func TestDescribeTable(t *testing.T) {
	client := dynamodb.Client{}

	if initialize(t, &client, false) == false {
		return
	}
	defer finalize(t, &client)

	if response, err := client.DescribeTable(t.Name()); err != nil {
		t.Fatal(err)
	} else if *response.Table.TableName != t.Name() {
		t.Fatal(*response.Table.TableName, ",", t.Name())
	} else if response.Table.TableStatus != types.TableStatusActive {
		t.Fatal(response.Table.TableStatus)
	}
}

func TestUpdateTable(t *testing.T) {
	client := dynamodb.Client{}

	if initialize(t, &client, false) == false {
		return
	}
	defer finalize(t, &client)

	if response, err := client.DescribeTable(t.Name()); err != nil {
		t.Fatal(err)
	} else if len(response.Table.GlobalSecondaryIndexes) != 0 {
		t.Fatal(response.Table.GlobalSecondaryIndexes)
	}

	if response, err := client.UpdateTable(
		&aws_dynamodb.UpdateTableInput{
			TableName: aws.String(t.Name()),
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
		t.Fatal(*response.TableDescription.GlobalSecondaryIndexes[0].IndexName, ",", INDEX_NAME)
	}
}

func TestDeleteTable(t *testing.T) {
	client := dynamodb.Client{}

	if initialize(t, &client, false) == false {
		return
	}
	defer finalize(t, &client)
}

func TestGetItem(t *testing.T) {
	client := dynamodb.Client{}

	if initialize(t, &client, false) == false {
		return
	}
	defer finalize(t, &client)

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
				TableName: aws.String(t.Name()),
				Item:      item,
			}); err != nil {
			t.Fatal(err)
		}
	}

	{
		testItemForGet := TestItem{PrimaryKey: 1, SortKey: "a"}

		if key, err := testItemForGet.getKey(); err != nil {
			t.Fatal(err)
		} else if response, err := client.GetItem(&aws_dynamodb.GetItemInput{
			TableName: aws.String(t.Name()), Key: key}); err != nil {
			t.Fatal(err)
		} else if err := attributevalue.UnmarshalMap(response.Item, &testItemForGet); err != nil {
			t.Fatal(err)
		} else if testItemForGet.Field1 != testItemForPut.Field1 ||
			testItemForGet.Field2 != testItemForPut.Field2 ||
			testItemForGet.Field3 != testItemForPut.Field3 ||
			testItemForGet.Field4[0].SubField1 != testItemForPut.Field4[0].SubField1 ||
			testItemForGet.Field4[0].SubField2 != testItemForPut.Field4[0].SubField2 ||
			testItemForGet.Field4[1].SubField1 != testItemForPut.Field4[1].SubField1 ||
			testItemForGet.Field4[1].SubField2 != testItemForPut.Field4[1].SubField2 {
			t.Log(testItemForGet)
			t.Log(testItemForPut)
			t.Fatal("invalid")
		}
	}
}

func TestPutItem(t *testing.T) {
	client := dynamodb.Client{}

	if initialize(t, &client, false) == false {
		return
	}
	defer finalize(t, &client)

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
				TableName: aws.String(t.Name()),
				Item:      item}); err != nil {
			t.Fatal(err)
		}
	}

	{
		testItemForGet := TestItem{PrimaryKey: 1, SortKey: "a"}

		if key, err := testItemForGet.getKey(); err != nil {
			t.Fatal(err)
		} else if response, err := client.GetItem(
			&aws_dynamodb.GetItemInput{TableName: aws.String(t.Name()), Key: key}); err != nil {
			t.Fatal(err)
		} else if err := attributevalue.UnmarshalMap(response.Item, &testItemForGet); err != nil {
			t.Fatal(err)
		} else if testItemForGet.Field1 != testItemForPut.Field1 ||
			testItemForGet.Field2 != testItemForPut.Field2 ||
			testItemForGet.Field3 != testItemForPut.Field3 ||
			testItemForGet.Field4[0].SubField1 != testItemForPut.Field4[0].SubField1 ||
			testItemForGet.Field4[0].SubField2 != testItemForPut.Field4[0].SubField2 ||
			testItemForGet.Field4[1].SubField1 != testItemForPut.Field4[1].SubField1 ||
			testItemForGet.Field4[1].SubField2 != testItemForPut.Field4[1].SubField2 {
			t.Log(testItemForGet)
			t.Log(testItemForPut)
			t.Fatal("invalid")
		}
	}
}

func TestUpdateItem(t *testing.T) {
	const updateValue = 10

	client := dynamodb.Client{}

	if initialize(t, &client, false) == false {
		return
	}
	defer finalize(t, &client)

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
				TableName: aws.String(t.Name()),
				Item:      item}); err != nil {
			t.Fatal(err)
		}
	}

	{
		testItem := TestItem{PrimaryKey: 1, SortKey: "a"}
		update := expression.Set(expression.Name("field2"), expression.Value(aws.Int(updateValue)))

		if key, err := testItem.getKey(); err != nil {
			t.Fatal(err)
		} else if expr, err := expression.NewBuilder().WithUpdate(update).Build(); err != nil {
			t.Fatal(err)
		} else if _, err = client.UpdateItem(
			&aws_dynamodb.UpdateItemInput{
				TableName:                 aws.String(t.Name()),
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

		if key, err := testItem.getKey(); err != nil {
			t.Fatal(err)
		} else if response, err := client.GetItem(
			&aws_dynamodb.GetItemInput{
				TableName: aws.String(t.Name()),
				Key:       key}); err != nil {
			t.Fatal(err)
		} else if err := attributevalue.UnmarshalMap(response.Item, &testItem); err != nil {
			t.Fatal(err)
		} else if testItem.Field2 != updateValue {
			t.Fatal(testItem.Field2, ",", updateValue)
		}
	}
}

func TestDeleteItem(t *testing.T) {
	client := dynamodb.Client{}

	if initialize(t, &client, false) == false {
		return
	}
	defer finalize(t, &client)

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
				TableName: aws.String(t.Name()),
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
				TableName: aws.String(t.Name()),
				Key:       key}); err != nil {
			t.Fatal(err)
		}
	}

	{
		if response, err := client.DescribeTable(t.Name()); err != nil {
			t.Fatal(err)
		} else if *response.Table.ItemCount != 0 {
			t.Fatal(*response.Table.ItemCount)
		}
	}
}

func TestQuery(t *testing.T) {
	client := dynamodb.Client{}

	if initialize(t, &client, true) == false {
		return
	}
	defer finalize(t, &client)

	keyEx := expression.Key("primary-key").Equal(expression.Value(3)).And(expression.Key("sort-key").Equal(expression.Value("c-1")))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		t.Fatal(err)
	}

	response, err := client.Query(
		&aws_dynamodb.QueryInput{
			TableName:                 aws.String(t.Name()),
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
		t.Fatal(testItems)
	}
}

func TestScan(t *testing.T) {
	client := dynamodb.Client{}

	if initialize(t, &client, true) == false {
		return
	}
	defer finalize(t, &client)

	filtEx := expression.Name("primary-key").Between(expression.Value(2), expression.Value(3))
	projEx := expression.NamesList(expression.Name("primary-key"), expression.Name("sort-key"), expression.Name("field2"))
	expr, err := expression.NewBuilder().WithFilter(filtEx).WithProjection(projEx).Build()

	response, err := client.Scan(
		&aws_dynamodb.ScanInput{
			TableName:                 aws.String(t.Name()),
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
		t.Fatal(testItems)
	}
}

func TestDescribeTimeToLive(t *testing.T) {
	client := dynamodb.Client{}

	if initialize(t, &client, true) == false {
		return
	}
	defer finalize(t, &client)

	if response, err := client.DescribeTimeToLive(t.Name()); err != nil {
		t.Fatal(err)
	} else if response.TimeToLiveDescription.AttributeName != nil {
		t.Fatal(response.TimeToLiveDescription.AttributeName)
	}

	if response, err := client.UpdateTimeToLive(t.Name(), TTL_NAME, true); err != nil {
		t.Fatal(err)
	} else if *response.TimeToLiveSpecification.AttributeName != TTL_NAME {
		t.Fatal(*response.TimeToLiveSpecification.AttributeName)
	} else if *response.TimeToLiveSpecification.Enabled != true {
		t.Fatal(*response.TimeToLiveSpecification.Enabled)
	}

	if response, err := client.DescribeTimeToLive(t.Name()); err != nil {
		t.Fatal(err)
	} else if *response.TimeToLiveDescription.AttributeName != TTL_NAME {
		t.Fatal(*response.TimeToLiveDescription.AttributeName)
	} else if response.TimeToLiveDescription.TimeToLiveStatus != types.TimeToLiveStatusEnabled {
		t.Fatal(response.TimeToLiveDescription.TimeToLiveStatus)
	}
}

func TestUpdateTimeToLive(t *testing.T) {
	client := dynamodb.Client{}

	if initialize(t, &client, true) == false {
		return
	}
	defer finalize(t, &client)

	if response, err := client.DescribeTimeToLive(t.Name()); err != nil {
		t.Fatal(err)
	} else if response.TimeToLiveDescription.AttributeName != nil {
		t.Fatal(response.TimeToLiveDescription.AttributeName)
	}

	if response, err := client.UpdateTimeToLive(t.Name(), TTL_NAME, true); err != nil {
		t.Fatal(err)
	} else if *response.TimeToLiveSpecification.AttributeName != TTL_NAME {
		t.Fatal(*response.TimeToLiveSpecification.AttributeName)
	} else if *response.TimeToLiveSpecification.Enabled != true {
		t.Fatal(*response.TimeToLiveSpecification.Enabled)
	}

	if response, err := client.DescribeTimeToLive(t.Name()); err != nil {
		t.Fatal(err)
	} else if *response.TimeToLiveDescription.AttributeName != TTL_NAME {
		t.Fatal(*response.TimeToLiveDescription.AttributeName)
	} else if response.TimeToLiveDescription.TimeToLiveStatus != types.TimeToLiveStatusEnabled {
		t.Fatal(response.TimeToLiveDescription.TimeToLiveStatus)
	}
}

func TestQueryPaginatorNextPage(t *testing.T) {
	client := dynamodb.Client{}

	if initialize(t, &client, true) == false {
		return
	}
	defer finalize(t, &client)

	keyEx := expression.Key("primary-key").Equal(expression.Value(3))
	if expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build(); err != nil {
		t.Fatal(err)
	} else if response, err := client.QueryPaginatorNextPage(
		&aws_dynamodb.QueryInput{
			TableName:                 aws.String(t.Name()),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			KeyConditionExpression:    expr.KeyCondition(),
			Limit:                     aws.Int32(1),
			ExclusiveStartKey:         nil,
		}); err != nil {
		t.Fatal(err)
	} else if response.Count != 1 || response.ScannedCount != 1 ||
		len(response.Items) != 1 || response.LastEvaluatedKey == nil {
		t.Fatal(response)
	}
}

func TestScanPaginatorNextPage(t *testing.T) {
	client := dynamodb.Client{}

	if initialize(t, &client, true) == false {
		return
	}
	defer finalize(t, &client)

	if response, err := client.ScanPaginatorNextPage(
		&aws_dynamodb.ScanInput{
			TableName:         aws.String(t.Name()),
			Limit:             aws.Int32(2),
			ExclusiveStartKey: nil,
		}); err != nil {
		t.Fatal(err)
	} else if response.Count != 2 || response.ScannedCount != 2 ||
		len(response.Items) != 2 || response.LastEvaluatedKey == nil {
		t.Fatal(response)
	}
}
