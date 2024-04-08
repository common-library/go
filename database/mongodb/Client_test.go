package mongodb_test

import (
	"strings"
	"testing"

	"github.com/common-library/go/database/mongodb"
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

	if err := client.Initialize(address, 0); err != nil {
		t.Error(err)
	}

	if err := client.Initialize(address, timeout); err != nil {
		t.Error(err)
	}

	if err := client.Finalize(); err != nil {
		t.Error(err)
	}
}

func TestFinalize(t *testing.T) {
	client := mongodb.Client{}

	if err := client.Finalize(); err != nil {
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

	if err := client.Initialize(address, timeout); err != nil {
		t.Error(err)
	}

	insertData := TestStruct{Value1: 1, Value2: "abc"}
	if err := client.InsertOne(database_name, collection_name, insertData); err != nil {
		t.Error(err)
	}

	if result_interface, err := client.FindOne(database_name, collection_name, filter, TestStruct{}); err != nil {
		t.Error(err)
	} else if result, ok := result_interface.(TestStruct); ok == false {
		t.Error("Type Assertions error")
	} else if result.Value1 != 1 || result.Value2 != "abc" {
		t.Errorf("invalid data - result : (%#v)", result)
	}

	if err := client.DeleteOne(database_name, collection_name, filter); err != nil {
		t.Error(err)
	}

	if err := client.Finalize(); err != nil {
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

	if err = client.Initialize(address, timeout); err != nil {
		t.Error(err)
	}

	insertData := make([]any, 0)
	insertData = append(insertData, TestStruct{Value1: 1, Value2: "abc"}, TestStruct{Value1: 2, Value2: "def"})
	if err := client.InsertMany(database_name, collection_name, insertData); err != nil {
		t.Error(err)
	}

	if results_interface, err := client.Find(database_name, collection_name, filter, TestStruct{}); err != nil {
		t.Error(err)
	} else if results, ok := results_interface.([]TestStruct); ok == false {
		t.Error("Type Assertions error")
	} else if results != nil &&
		(results[0].Value1 != 1 ||
			results[0].Value2 != "abc" ||
			results[1].Value1 != 2 ||
			results[1].Value2 != "def") {
		t.Errorf("invalid data - results : (%#v)", results)
	}

	if err := client.DeleteMany(database_name, collection_name, filter); err != nil {
		t.Error(err)
	}

	if err := client.Finalize(); err != nil {
		t.Error(err)
	}
}

func TestInsertOne(t *testing.T) {
	client := mongodb.Client{}

	err := client.InsertOne(database_name, collection_name, TestStruct{})
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	if err := client.Initialize(address, timeout); err != nil {
		t.Error(err)
	}

	insertData := TestStruct{Value1: 1, Value2: "abc"}
	if err := client.InsertOne(database_name, collection_name, insertData); err != nil {
		t.Error(err)
	}

	if err := client.DeleteMany(database_name, collection_name, bson.M{}); err != nil {
		t.Error(err)
	}

	if err := client.Finalize(); err != nil {
		t.Error(err)
	}
}

func TestInsertMany(t *testing.T) {
	client := mongodb.Client{}

	insertData := make([]any, 0)
	insertData = append(insertData, TestStruct{Value1: 1, Value2: "abc"}, TestStruct{Value1: 2, Value2: "def"})

	err := client.InsertMany(database_name, collection_name, insertData)
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	if err := client.Initialize(address, timeout); err != nil {
		t.Error(err)
	}

	if err := client.InsertMany(database_name, collection_name, insertData); err != nil {
		t.Error(err)
	}

	if err := client.DeleteMany(database_name, collection_name, bson.M{}); err != nil {
		t.Error(err)
	}

	if err := client.Finalize(); err != nil {
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

	if err := client.Initialize(address, timeout); err != nil {
		t.Error(err)
	}

	if err := client.UpdateOne(database_name, collection_name, filter, update); err != nil {
		t.Error(err)
	}

	_, err = client.FindOne(database_name, collection_name, filter, TestStruct{})
	if err.Error() != "mongo: no documents in result" {
		t.Error(err)
	}

	insertData := make([]any, 0)
	insertData = append(insertData, TestStruct{Value1: 1, Value2: "abc"}, TestStruct{Value1: 1, Value2: "abc"})
	if err := client.InsertMany(database_name, collection_name, insertData); err != nil {
		t.Error(err)
	}

	if err := client.UpdateOne(database_name, collection_name, filter, update); err != nil {
		t.Error(err)
	}

	if results_interface, err := client.Find(database_name, collection_name, filter, TestStruct{}); err != nil {
		t.Error(err)
	} else if results, ok := results_interface.([]TestStruct); ok == false {
		t.Error("Type Assertions error")
	} else if results != nil &&
		(results[0].Value1 != 1 ||
			results[0].Value2 != "update_value" ||
			results[1].Value1 != 1 ||
			results[1].Value2 != "abc") {
		t.Errorf("invalid data - results : (%#v)", results)
	}

	if err := client.DeleteMany(database_name, collection_name, bson.M{}); err != nil {
		t.Error(err)
	}

	if err := client.Finalize(); err != nil {
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

	if err := client.Initialize(address, timeout); err != nil {
		t.Error(err)
	}

	if err := client.UpdateMany(database_name, collection_name, filter, update); err != nil {
		t.Error(err)
	}

	_, err = client.FindOne(database_name, collection_name, filter, TestStruct{})
	if err.Error() != "mongo: no documents in result" {
		t.Error(err)
	}

	insertData := make([]any, 0)
	insertData = append(insertData, TestStruct{Value1: 1, Value2: "abc"}, TestStruct{Value1: 1, Value2: "abc"})
	if err := client.InsertMany(database_name, collection_name, insertData); err != nil {
		t.Error(err)
	}

	if err := client.UpdateMany(database_name, collection_name, filter, update); err != nil {
		t.Error(err)
	}

	if results_interface, err := client.Find(database_name, collection_name, filter, TestStruct{}); err != nil {
		t.Error(err)
	} else if results, ok := results_interface.([]TestStruct); ok == false {
		t.Error("Type Assertions error")
	} else if results != nil &&
		(results[0].Value1 != 1 ||
			results[0].Value2 != "update_value" ||
			results[1].Value1 != 1 ||
			results[1].Value2 != "update_value") {
		t.Errorf("invalid data - results : (%#v)", results)
	}

	if err := client.DeleteMany(database_name, collection_name, bson.M{}); err != nil {
		t.Error(err)
	}

	if err := client.Finalize(); err != nil {
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

	if err := client.Initialize(address, timeout); err != nil {
		t.Error(err)
	}

	if err := client.DeleteOne(database_name, collection_name, filter); err != nil {
		t.Error(err)
	}

	_, err = client.FindOne(database_name, collection_name, filter, TestStruct{})
	if err.Error() != "mongo: no documents in result" {
		t.Error(err)
	}

	if err := client.Finalize(); err != nil {
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

	if err = client.Initialize(address, timeout); err != nil {
		t.Error(err)
	}

	if err = client.DeleteMany(database_name, collection_name, filter); err != nil {
		t.Error(err)
	}

	_, err = client.FindOne(database_name, collection_name, filter, TestStruct{})
	if err.Error() != "mongo: no documents in result" {
		t.Error(err)
	}

	if err := client.Finalize(); err != nil {
		t.Error(err)
	}
}
