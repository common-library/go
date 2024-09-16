package mongodb_test

import (
	"os"
	"testing"

	"github.com/common-library/go/database/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

type TestStruct struct {
	Value1 int
	Value2 string
}

func getClient(t *testing.T) (*mongodb.Client, bool) {
	t.Parallel()

	client := &mongodb.Client{}
	address := os.Getenv("MONGODB_ADDRESS")
	if len(address) == 0 {
		return nil, true
	} else if err := client.Initialize(address, 10); err != nil {
		t.Fatal(err)
	}

	return client, false
}

func finalize(t *testing.T, client *mongodb.Client) {
	if err := client.Finalize(); err != nil {
		t.Fatal(err)
	}
}

func getDatabaseName(t *testing.T) string {
	return t.Name()
}

func getCollectionName(t *testing.T) string {
	return t.Name()
}

func TestInitialize(t *testing.T) {
	if client, stop := getClient(t); stop {
		return
	} else if err := client.Finalize(); err != nil {
		t.Fatal(err)
	}
}

func TestFinalize(t *testing.T) {
	TestInitialize(t)
}

func TestFindOne(t *testing.T) {
	databaseName := getDatabaseName(t)
	collectionName := getCollectionName(t)

	filter := bson.M{"value1": 1}

	if _, err := (&mongodb.Client{}).FindOne(databaseName, collectionName, filter, TestStruct{}); err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	client, stop := getClient(t)
	if stop {
		return
	}
	defer finalize(t, client)

	insertData := TestStruct{Value1: 1, Value2: "abc"}
	if err := client.InsertOne(databaseName, collectionName, insertData); err != nil {
		t.Fatal(err)
	}

	if result_interface, err := client.FindOne(databaseName, collectionName, filter, TestStruct{}); err != nil {
		t.Fatal(err)
	} else if result, ok := result_interface.(TestStruct); ok == false {
		t.Fatal("Type Assertions error")
	} else if result.Value1 != 1 || result.Value2 != "abc" {
		t.Fatal(result)
	}

	if err := client.DeleteOne(databaseName, collectionName, filter); err != nil {
		t.Fatal(err)
	}
}

func TestFind(t *testing.T) {
	databaseName := getDatabaseName(t)
	collectionName := getCollectionName(t)

	filter := bson.M{}

	if _, err := (&mongodb.Client{}).Find(databaseName, collectionName, filter, TestStruct{}); err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	client, stop := getClient(t)
	if stop {
		return
	}
	defer finalize(t, client)

	insertData := make([]any, 0)
	insertData = append(insertData, TestStruct{Value1: 1, Value2: "abc"}, TestStruct{Value1: 2, Value2: "def"})
	if err := client.InsertMany(databaseName, collectionName, insertData); err != nil {
		t.Fatal(err)
	}

	if results_interface, err := client.Find(databaseName, collectionName, filter, TestStruct{}); err != nil {
		t.Fatal(err)
	} else if results, ok := results_interface.([]TestStruct); ok == false {
		t.Fatal("Type Assertions error")
	} else if results != nil &&
		(results[0].Value1 != 1 ||
			results[0].Value2 != "abc" ||
			results[1].Value1 != 2 ||
			results[1].Value2 != "def") {
		t.Fatal(results)
	}

	if err := client.DeleteMany(databaseName, collectionName, filter); err != nil {
		t.Fatal(err)
	}
}

func TestInsertOne(t *testing.T) {
	databaseName := getDatabaseName(t)
	collectionName := getCollectionName(t)

	if err := (&mongodb.Client{}).InsertOne(databaseName, collectionName, TestStruct{}); err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	client, stop := getClient(t)
	if stop {
		return
	}
	defer finalize(t, client)

	insertData := TestStruct{Value1: 1, Value2: "abc"}
	if err := client.InsertOne(databaseName, collectionName, insertData); err != nil {
		t.Fatal(err)
	}

	if err := client.DeleteMany(databaseName, collectionName, bson.M{}); err != nil {
		t.Fatal(err)
	}
}

func TestInsertMany(t *testing.T) {
	databaseName := getDatabaseName(t)
	collectionName := getCollectionName(t)

	insertData := make([]any, 0)
	insertData = append(insertData, TestStruct{Value1: 1, Value2: "abc"}, TestStruct{Value1: 2, Value2: "def"})

	if err := (&mongodb.Client{}).InsertMany(databaseName, collectionName, insertData); err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	client, stop := getClient(t)
	if stop {
		return
	}
	defer finalize(t, client)

	if err := client.InsertMany(databaseName, collectionName, insertData); err != nil {
		t.Fatal(err)
	}

	if err := client.DeleteMany(databaseName, collectionName, bson.M{}); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateOne(t *testing.T) {
	databaseName := getDatabaseName(t)
	collectionName := getCollectionName(t)

	filter := bson.M{"value1": 1}
	update := bson.D{{"$set", bson.D{{"value2", "update_value"}}}}

	if err := (&mongodb.Client{}).UpdateOne(databaseName, collectionName, filter, update); err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	client, stop := getClient(t)
	if stop {
		return
	}
	defer finalize(t, client)

	if err := client.UpdateOne(databaseName, collectionName, filter, update); err != nil {
		t.Fatal(err)
	}

	if _, err := client.FindOne(databaseName, collectionName, filter, TestStruct{}); err.Error() != "mongo: no documents in result" {
		t.Fatal(err)
	}

	insertData := make([]any, 0)
	insertData = append(insertData, TestStruct{Value1: 1, Value2: "abc"}, TestStruct{Value1: 1, Value2: "abc"})
	if err := client.InsertMany(databaseName, collectionName, insertData); err != nil {
		t.Fatal(err)
	}

	if err := client.UpdateOne(databaseName, collectionName, filter, update); err != nil {
		t.Fatal(err)
	}

	if results_interface, err := client.Find(databaseName, collectionName, filter, TestStruct{}); err != nil {
		t.Fatal(err)
	} else if results, ok := results_interface.([]TestStruct); ok == false {
		t.Fatal("Type Assertions error")
	} else if results != nil &&
		(results[0].Value1 != 1 ||
			results[0].Value2 != "update_value" ||
			results[1].Value1 != 1 ||
			results[1].Value2 != "abc") {
		t.Fatal(results)
	}

	if err := client.DeleteMany(databaseName, collectionName, bson.M{}); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateMany(t *testing.T) {
	databaseName := getDatabaseName(t)
	collectionName := getCollectionName(t)

	filter := bson.M{"value1": 1}
	update := bson.D{{"$set", bson.D{{"value2", "update_value"}}}}

	if err := (&mongodb.Client{}).UpdateMany(databaseName, collectionName, filter, update); err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	client, stop := getClient(t)
	if stop {
		return
	}
	defer finalize(t, client)

	if err := client.UpdateMany(databaseName, collectionName, filter, update); err != nil {
		t.Fatal(err)
	}

	if _, err := client.FindOne(databaseName, collectionName, filter, TestStruct{}); err.Error() != "mongo: no documents in result" {
		t.Fatal(err)
	}

	insertData := make([]any, 0)
	insertData = append(insertData, TestStruct{Value1: 1, Value2: "abc"}, TestStruct{Value1: 1, Value2: "abc"})
	if err := client.InsertMany(databaseName, collectionName, insertData); err != nil {
		t.Fatal(err)
	}

	if err := client.UpdateMany(databaseName, collectionName, filter, update); err != nil {
		t.Fatal(err)
	}

	if results_interface, err := client.Find(databaseName, collectionName, filter, TestStruct{}); err != nil {
		t.Fatal(err)
	} else if results, ok := results_interface.([]TestStruct); ok == false {
		t.Fatal("Type Assertions error")
	} else if results != nil &&
		(results[0].Value1 != 1 ||
			results[0].Value2 != "update_value" ||
			results[1].Value1 != 1 ||
			results[1].Value2 != "update_value") {
		t.Fatal(results)
	}

	if err := client.DeleteMany(databaseName, collectionName, bson.M{}); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteOne(t *testing.T) {
	databaseName := getDatabaseName(t)
	collectionName := getCollectionName(t)

	filter := bson.M{}

	if err := (&mongodb.Client{}).DeleteOne(databaseName, collectionName, filter); err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	client, stop := getClient(t)
	if stop {
		return
	}
	defer finalize(t, client)

	if err := client.DeleteOne(databaseName, collectionName, filter); err != nil {
		t.Fatal(err)
	}

	if _, err := client.FindOne(databaseName, collectionName, filter, TestStruct{}); err.Error() != "mongo: no documents in result" {
		t.Fatal(err)
	}
}

func TestDeleteMany(t *testing.T) {
	databaseName := getDatabaseName(t)
	collectionName := getCollectionName(t)

	filter := bson.M{}

	if err := (&mongodb.Client{}).DeleteMany(databaseName, collectionName, filter); err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	client, stop := getClient(t)
	if stop {
		return
	}
	defer finalize(t, client)

	if err := client.DeleteMany(databaseName, collectionName, filter); err != nil {
		t.Fatal(err)
	}

	if _, err := client.FindOne(databaseName, collectionName, filter, TestStruct{}); err.Error() != "mongo: no documents in result" {
		t.Fatal(err)
	}
}
