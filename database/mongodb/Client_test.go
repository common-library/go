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

const databaseName string = "testDatabase"
const collectionName string = "testCollection"

func getClient(t *testing.T) (mongodb.Client, bool) {
	client := mongodb.Client{}

	if len(os.Getenv("MONGODB_ADDRESS")) == 0 {
		return client, false
	}

	if err := client.Initialize(os.Getenv("MONGODB_ADDRESS"), 10); err != nil {
		t.Fatal(err)
	}

	return client, true
}

func TestInitialize(t *testing.T) {
	client, ok := getClient(t)
	if ok == false {
		return
	}

	if err := client.Finalize(); err != nil {
		t.Fatal(err)
	}
}

func TestFinalize(t *testing.T) {
	client, ok := getClient(t)
	if ok == false {
		return
	}

	if err := client.Finalize(); err != nil {
		t.Fatal(err)
	}
}

func TestFindOne(t *testing.T) {
	filter := bson.M{"value1": 1}

	client := mongodb.Client{}
	_, err := client.FindOne(databaseName, collectionName, filter, TestStruct{})
	if err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	client, ok := getClient(t)
	if ok == false {
		return
	}

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

	if err := client.Finalize(); err != nil {
		t.Fatal(err)
	}
}

func TestFind(t *testing.T) {
	filter := bson.M{}

	client := mongodb.Client{}
	_, err := client.Find(databaseName, collectionName, filter, TestStruct{})
	if err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	client, ok := getClient(t)
	if ok == false {
		return
	}

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

	if err := client.Finalize(); err != nil {
		t.Fatal(err)
	}
}

func TestInsertOne(t *testing.T) {
	client := mongodb.Client{}
	err := client.InsertOne(databaseName, collectionName, TestStruct{})
	if err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	client, ok := getClient(t)
	if ok == false {
		return
	}

	insertData := TestStruct{Value1: 1, Value2: "abc"}
	if err := client.InsertOne(databaseName, collectionName, insertData); err != nil {
		t.Fatal(err)
	}

	if err := client.DeleteMany(databaseName, collectionName, bson.M{}); err != nil {
		t.Fatal(err)
	}

	if err := client.Finalize(); err != nil {
		t.Fatal(err)
	}
}

func TestInsertMany(t *testing.T) {
	insertData := make([]any, 0)
	insertData = append(insertData, TestStruct{Value1: 1, Value2: "abc"}, TestStruct{Value1: 2, Value2: "def"})

	client := mongodb.Client{}
	err := client.InsertMany(databaseName, collectionName, insertData)
	if err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	client, ok := getClient(t)
	if ok == false {
		return
	}

	if err := client.InsertMany(databaseName, collectionName, insertData); err != nil {
		t.Fatal(err)
	}

	if err := client.DeleteMany(databaseName, collectionName, bson.M{}); err != nil {
		t.Fatal(err)
	}

	if err := client.Finalize(); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateOne(t *testing.T) {
	filter := bson.M{"value1": 1}
	update := bson.D{{"$set", bson.D{{"value2", "update_value"}}}}

	client := mongodb.Client{}
	err := client.UpdateOne(databaseName, collectionName, filter, update)
	if err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	client, ok := getClient(t)
	if ok == false {
		return
	}

	if err := client.UpdateOne(databaseName, collectionName, filter, update); err != nil {
		t.Fatal(err)
	}

	_, err = client.FindOne(databaseName, collectionName, filter, TestStruct{})
	if err.Error() != "mongo: no documents in result" {
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

	if err := client.Finalize(); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateMany(t *testing.T) {
	filter := bson.M{"value1": 1}
	update := bson.D{{"$set", bson.D{{"value2", "update_value"}}}}

	client := mongodb.Client{}
	err := client.UpdateMany(databaseName, collectionName, filter, update)
	if err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	client, ok := getClient(t)
	if ok == false {
		return
	}

	if err := client.UpdateMany(databaseName, collectionName, filter, update); err != nil {
		t.Fatal(err)
	}

	_, err = client.FindOne(databaseName, collectionName, filter, TestStruct{})
	if err.Error() != "mongo: no documents in result" {
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

	if err := client.Finalize(); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteOne(t *testing.T) {
	filter := bson.M{}

	client := mongodb.Client{}
	err := client.DeleteOne(databaseName, collectionName, filter)
	if err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	client, ok := getClient(t)
	if ok == false {
		return
	}

	if err := client.DeleteOne(databaseName, collectionName, filter); err != nil {
		t.Fatal(err)
	}

	_, err = client.FindOne(databaseName, collectionName, filter, TestStruct{})
	if err.Error() != "mongo: no documents in result" {
		t.Fatal(err)
	}

	if err := client.Finalize(); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteMany(t *testing.T) {
	filter := bson.M{}

	client := mongodb.Client{}
	err := client.DeleteMany(databaseName, collectionName, filter)
	if err.Error() != "please call Initialize first" {
		t.Fatal(err)
	}

	client, ok := getClient(t)
	if ok == false {
		return
	}

	if err = client.DeleteMany(databaseName, collectionName, filter); err != nil {
		t.Fatal(err)
	}

	_, err = client.FindOne(databaseName, collectionName, filter, TestStruct{})
	if err.Error() != "mongo: no documents in result" {
		t.Fatal(err)
	}

	if err := client.Finalize(); err != nil {
		t.Fatal(err)
	}
}
