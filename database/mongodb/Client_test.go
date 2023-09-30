package mongodb_test

import (
	"strings"
	"testing"

	"github.com/heaven-chp/common-library-go/database/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

type TestStruct struct {
	Value1 int
	Value2 string
}

const address string = "localhost:27017"
const timeout uint64 = 3
const database_name string = "testDatabase"
const collection_name string = "testCollection"

func TestInitialize(t *testing.T) {
	client := mongodb.Client{}

	err := client.Initialize("invalid_address", timeout)
	if strings.HasPrefix(err.Error(), "server selection error:") == false {
		t.Error(err)
	}

	err = client.Initialize(address, 0)
	if err != nil {
		t.Error(err)
	}

	err = client.Initialize(address, timeout)
	if err != nil {
		t.Error(err)
	}

	err = client.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestFinalize(t *testing.T) {
	client := mongodb.Client{}

	err := client.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestFindOne(t *testing.T) {
	client := mongodb.Client{}

	filter := bson.M{"value1": 1}

	_, err := client.FindOne(database_name, collection_name, filter, TestStruct{})
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = client.Initialize(address, timeout)
	if err != nil {
		t.Error(err)
	}

	insertData := TestStruct{Value1: 1, Value2: "abc"}
	err = client.InsertOne(database_name, collection_name, insertData)
	if err != nil {
		t.Error(err)
	}

	result_interface, err := client.FindOne(database_name, collection_name, filter, TestStruct{})
	if err != nil {
		t.Error(err)
	}

	result, ok := result_interface.(TestStruct)
	if ok == false {
		t.Error("Type Assertions error")
	}

	if result.Value1 != 1 || result.Value2 != "abc" {
		t.Errorf("invalid data - result : (%#v)", result)
	}

	err = client.DeleteOne(database_name, collection_name, filter)
	if err != nil {
		t.Error(err)
	}

	err = client.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestFind(t *testing.T) {
	client := mongodb.Client{}

	filter := bson.M{}

	_, err := client.Find(database_name, collection_name, filter, TestStruct{})
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = client.Initialize(address, timeout)
	if err != nil {
		t.Error(err)
	}

	insertData := make([]interface{}, 0)
	insertData = append(insertData, TestStruct{Value1: 1, Value2: "abc"}, TestStruct{Value1: 2, Value2: "def"})

	err = client.InsertMany(database_name, collection_name, insertData)
	if err != nil {
		t.Error(err)
	}

	results_interface, err := client.Find(database_name, collection_name, filter, TestStruct{})
	if err != nil {
		t.Error(err)
	}

	results, ok := results_interface.([]TestStruct)
	if ok == false {
		t.Error("Type Assertions error")
	}

	if results != nil &&
		(results[0].Value1 != 1 ||
			results[0].Value2 != "abc" ||
			results[1].Value1 != 2 ||
			results[1].Value2 != "def") {
		t.Errorf("invalid data - results : (%#v)", results)
	}

	err = client.DeleteMany(database_name, collection_name, filter)
	if err != nil {
		t.Error(err)
	}

	err = client.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestInsertOne(t *testing.T) {
	client := mongodb.Client{}

	err := client.InsertOne(database_name, collection_name, TestStruct{})
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = client.Initialize(address, timeout)
	if err != nil {
		t.Error(err)
	}

	insertData := TestStruct{Value1: 1, Value2: "abc"}
	err = client.InsertOne(database_name, collection_name, insertData)
	if err != nil {
		t.Error(err)
	}

	err = client.DeleteMany(database_name, collection_name, bson.M{})
	if err != nil {
		t.Error(err)
	}

	err = client.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestInsertMany(t *testing.T) {
	client := mongodb.Client{}

	insertData := make([]interface{}, 0)
	insertData = append(insertData, TestStruct{Value1: 1, Value2: "abc"}, TestStruct{Value1: 2, Value2: "def"})

	err := client.InsertMany(database_name, collection_name, insertData)
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = client.Initialize(address, timeout)
	if err != nil {
		t.Error(err)
	}

	err = client.InsertMany(database_name, collection_name, insertData)
	if err != nil {
		t.Error(err)
	}

	err = client.DeleteMany(database_name, collection_name, bson.M{})
	if err != nil {
		t.Error(err)
	}

	err = client.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateOne(t *testing.T) {
	client := mongodb.Client{}

	filter := bson.M{"value1": 1}
	update := bson.D{{"$set", bson.D{{"value2", "update_value"}}}}

	err := client.UpdateOne(database_name, collection_name, filter, update)
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = client.Initialize(address, timeout)
	if err != nil {
		t.Error(err)
	}

	err = client.UpdateOne(database_name, collection_name, filter, update)
	if err != nil {
		t.Error(err)
	}

	_, err = client.FindOne(database_name, collection_name, filter, TestStruct{})
	if err.Error() != "mongo: no documents in result" {
		t.Error(err)
	}

	insertData := make([]interface{}, 0)
	insertData = append(insertData, TestStruct{Value1: 1, Value2: "abc"}, TestStruct{Value1: 1, Value2: "abc"})

	err = client.InsertMany(database_name, collection_name, insertData)
	if err != nil {
		t.Error(err)
	}

	err = client.UpdateOne(database_name, collection_name, filter, update)
	if err != nil {
		t.Error(err)
	}

	results_interface, err := client.Find(database_name, collection_name, filter, TestStruct{})
	if err != nil {
		t.Error(err)
	}

	results, ok := results_interface.([]TestStruct)
	if ok == false {
		t.Error("Type Assertions error")
	}

	if results != nil &&
		(results[0].Value1 != 1 ||
			results[0].Value2 != "update_value" ||
			results[1].Value1 != 1 ||
			results[1].Value2 != "abc") {
		t.Errorf("invalid data - results : (%#v)", results)
	}

	err = client.DeleteMany(database_name, collection_name, bson.M{})
	if err != nil {
		t.Error(err)
	}

	err = client.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateMany(t *testing.T) {
	client := mongodb.Client{}

	filter := bson.M{"value1": 1}
	update := bson.D{{"$set", bson.D{{"value2", "update_value"}}}}

	err := client.UpdateMany(database_name, collection_name, filter, update)
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = client.Initialize(address, timeout)
	if err != nil {
		t.Error(err)
	}

	err = client.UpdateMany(database_name, collection_name, filter, update)
	if err != nil {
		t.Error(err)
	}

	_, err = client.FindOne(database_name, collection_name, filter, TestStruct{})
	if err.Error() != "mongo: no documents in result" {
		t.Error(err)
	}

	insertData := make([]interface{}, 0)
	insertData = append(insertData, TestStruct{Value1: 1, Value2: "abc"}, TestStruct{Value1: 1, Value2: "abc"})

	err = client.InsertMany(database_name, collection_name, insertData)
	if err != nil {
		t.Error(err)
	}

	err = client.UpdateMany(database_name, collection_name, filter, update)
	if err != nil {
		t.Error(err)
	}

	results_interface, err := client.Find(database_name, collection_name, filter, TestStruct{})
	if err != nil {
		t.Error(err)
	}

	results, ok := results_interface.([]TestStruct)
	if ok == false {
		t.Error("Type Assertions error")
	}

	if results != nil &&
		(results[0].Value1 != 1 ||
			results[0].Value2 != "update_value" ||
			results[1].Value1 != 1 ||
			results[1].Value2 != "update_value") {
		t.Errorf("invalid data - results : (%#v)", results)
	}

	err = client.DeleteMany(database_name, collection_name, bson.M{})
	if err != nil {
		t.Error(err)
	}

	err = client.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteOne(t *testing.T) {
	client := mongodb.Client{}

	filter := bson.M{}

	err := client.DeleteOne(database_name, collection_name, filter)
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = client.Initialize(address, timeout)
	if err != nil {
		t.Error(err)
	}

	err = client.DeleteOne(database_name, collection_name, filter)
	if err != nil {
		t.Error(err)
	}

	_, err = client.FindOne(database_name, collection_name, filter, TestStruct{})
	if err.Error() != "mongo: no documents in result" {
		t.Error(err)
	}

	err = client.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteMany(t *testing.T) {
	client := mongodb.Client{}

	filter := bson.M{}

	err := client.DeleteMany(database_name, collection_name, filter)
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = client.Initialize(address, timeout)
	if err != nil {
		t.Error(err)
	}

	err = client.DeleteMany(database_name, collection_name, filter)
	if err != nil {
		t.Error(err)
	}

	_, err = client.FindOne(database_name, collection_name, filter, TestStruct{})
	if err.Error() != "mongo: no documents in result" {
		t.Error(err)
	}

	err = client.Finalize()
	if err != nil {
		t.Error(err)
	}
}
