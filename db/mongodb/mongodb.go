// Package mongodb provides mongodb interface.
// used "go.mongodb.org/mongo-driver/mongo".
package mongodb

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"reflect"
	"time"
)

// Mongodb is object that provides mongodb interface.
type Mongodb struct {
	address string
	timeout int

	ctx           context.Context
	ctxCancelFunc context.CancelFunc

	client *mongo.Client
}

// Initialize is initialize.
//  ex) err := mongodb.Initialize("localhost:27017", 10)
func (mongodb *Mongodb) Initialize(address string, timeout int) error {
	mongodb.address = address
	mongodb.timeout = timeout

	mongodb.ctx, mongodb.ctxCancelFunc = context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)

	var err error
	mongodb.client, err = mongo.Connect(mongodb.ctx, options.Client().ApplyURI("mongodb://"+address))
	if err != nil {
		return err
	}

	err = mongodb.client.Ping(mongodb.ctx, readpref.Primary())
	if err != nil {
		return err
	}

	return nil
}

// Finalize is finalize.
//  ex) err := mongodb.Finalize()
func (mongodb *Mongodb) Finalize() error {
	if mongodb.client != nil {
		err := mongodb.client.Disconnect(mongodb.ctx)
		mongodb.client = nil
		if err != nil {
			return err
		}
	}

	if mongodb.ctxCancelFunc != nil {
		mongodb.ctxCancelFunc()
	}

	return nil
}

// "FindOne" is returns one result value corresponding to the filter argument as an "dataForm" argument type interface.
//  ex)
//     result_interface, err := mongodb.FindOne("test_database", "test_collection", bson.M{"value1": 1}, TestStruct{})
//     result, ok := result_interface.(TestStruct)
func (mongodb *Mongodb) FindOne(databaseName string, collectionName string, filter interface{}, dataForm interface{}) (interface{}, error) {
	if mongodb.client == nil {
		return nil, errors.New("please call Initialize first")
	}

	collection := mongodb.client.Database(databaseName).Collection(collectionName)

	document := reflect.New(reflect.TypeOf(dataForm))
	err := collection.FindOne(mongodb.ctx, filter).Decode(document.Interface())
	if err != nil {
		return nil, err
	}

	return document.Elem().Interface(), nil
}

// "Find" is returns results value corresponding to the filter argument as an "dataForm" argument array type interface.
//  ex)
//      results_interface, err := mongodb.Find("test_database", "test_collection", bson.M{}, TestStruct{})
//      results, ok := results_interface.([]TestStruct)
func (mongodb *Mongodb) Find(databaseName string, collectionName string, filter interface{}, dataForm interface{}) (interface{}, error) {
	if mongodb.client == nil {
		return nil, errors.New("please call Initialize first")
	}

	collection := mongodb.client.Database(databaseName).Collection(collectionName)

	cursor, err := collection.Find(mongodb.ctx, filter)
	if err != nil {
		return nil, err
	}

	dataType := reflect.TypeOf(dataForm)
	results := reflect.MakeSlice(reflect.SliceOf(dataType), 0, 1024)
	document := reflect.New(dataType)

	for cursor.Next(mongodb.ctx) {
		err := cursor.Decode(document.Interface())
		if err != nil {
			return nil, err
		}

		results = reflect.Append(results, document.Elem())
	}

	return results.Interface(), nil
}

// InsertOne is insert a one document.
//  ex) err := mongodb.InsertOne("test_database", "test_collection", TestStruct{})
func (mongodb *Mongodb) InsertOne(databaseName string, collectionName string, document interface{}) error {
	if mongodb.client == nil {
		return errors.New("please call Initialize first")
	}

	collection := mongodb.client.Database(databaseName).Collection(collectionName)

	_, err := collection.InsertOne(mongodb.ctx, document)
	return err
}

// InsertMany is insert a array type documents.
//  ex)
//      insertData := make([]interface{}, 0)
//      insertData = append(insertData, TestStruct{Value1: 1, Value2: "abc"}, TestStruct{Value1: 2, Value2: "def"})
//      err := mongodb.InsertMany("test_database", "test_collection", insertData)
func (mongodb *Mongodb) InsertMany(databaseName string, collectionName string, documents []interface{}) error {
	if mongodb.client == nil {
		return errors.New("please call Initialize first")
	}

	collection := mongodb.client.Database(databaseName).Collection(collectionName)

	_, err := collection.InsertMany(mongodb.ctx, documents)
	return err
}

// UpdateOne is update the one value corresponding to the filter argument with the value of the "update" argument.
//  ex) err := mongodb.UpdateOne("test_database", "test_collection", bson.M{"value1": 1}, bson.D{{"$set", bson.D{{"value2", "update_value"}}}})
func (mongodb *Mongodb) UpdateOne(databaseName string, collectionName string, filter interface{}, update interface{}) error {
	if mongodb.client == nil {
		return errors.New("please call Initialize first")
	}

	collection := mongodb.client.Database(databaseName).Collection(collectionName)

	_, err := collection.UpdateOne(mongodb.ctx, filter, update)
	return err
}

// UpdateMany is update the value corresponding to the filter argument with the values of the "update" argument.
//  ex) err := mongodb.UpdateMany("test_database", "test_collection", bson.M{"value1": 1}, bson.D{{"$set", bson.D{{"value2", "update_value"}}}})
func (mongodb *Mongodb) UpdateMany(databaseName string, collectionName string, filter interface{}, update interface{}) error {
	if mongodb.client == nil {
		return errors.New("please call Initialize first")
	}

	collection := mongodb.client.Database(databaseName).Collection(collectionName)

	_, err := collection.UpdateMany(mongodb.ctx, filter, update)
	return err
}

// DeleteOne is delete one value corresponding to the filter argument.
//  ex) err := mongodb.DeleteOne("test_database", "test_collection", bson.M{"value1": 1})
func (mongodb *Mongodb) DeleteOne(databaseName string, collectionName string, filter interface{}) error {
	if mongodb.client == nil {
		return errors.New("please call Initialize first")
	}

	collection := mongodb.client.Database(databaseName).Collection(collectionName)

	_, err := collection.DeleteOne(mongodb.ctx, filter)
	return err
}

// DeleteMany is delete the values corresponding to the filter argument.
//  ex) err := mongodb.DeleteMany("test_database", "test_collection", bson.M{})
func (mongodb *Mongodb) DeleteMany(databaseName string, collectionName string, filter interface{}) error {
	if mongodb.client == nil {
		return errors.New("please call Initialize first")
	}

	collection := mongodb.client.Database(databaseName).Collection(collectionName)

	_, err := collection.DeleteMany(mongodb.ctx, filter)
	return err
}
