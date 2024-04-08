// Package mongodb provides MongoDB client implementations.
package mongodb

import (
	"context"
	"errors"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Client is a struct that provides client related methods.
type Client struct {
	address string
	timeout time.Duration

	ctx           context.Context
	ctxCancelFunc context.CancelFunc

	client *mongo.Client
}

// connect is connect.
//
// ex) err := client.connect()
func (this *Client) connect() error {
	if this.client != nil && this.client.Ping(this.ctx, readpref.Primary()) == nil {
		return nil
	}

	this.disConnect()

	if client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://" + this.address)); err != nil {
		return err
	} else {
		this.client = client
	}

	this.ctx = context.TODO()
	if this.timeout > 0 {
		this.ctx, this.ctxCancelFunc = context.WithTimeout(context.Background(), this.timeout*time.Second)
	}

	if err := this.client.Connect(this.ctx); err != nil {
		return err
	}

	return this.client.Ping(this.ctx, readpref.Primary())
}

// disconnect is disconnect.
//
// ex) err := client.disconnect()
func (this *Client) disConnect() error {
	if this.client == nil {
		return nil
	}

	err := this.client.Disconnect(this.ctx)

	if this.ctxCancelFunc != nil {
		this.ctxCancelFunc()
	}

	this.client = nil

	return err
}

// Initialize is initialize.
//
// ex) err := client.Initialize("localhost:27017", 10)
func (this *Client) Initialize(address string, timeout time.Duration) error {
	this.address = address
	this.timeout = timeout

	return this.connect()
}

// Finalize is finalize.
//
// ex) err := client.Finalize()
func (this *Client) Finalize() error {
	return this.disConnect()
}

// "FindOne" is returns one result value corresponding to the filter argument as an "dataForm" argument type interface.
//
//	ex)
//
//	 result_interface, err := client.FindOne("test_database", "test_collection", bson.M{"value1": 1}, TestStruct{})
//
//	 result, ok := result_interface.(TestStruct)
func (this *Client) FindOne(databaseName string, collectionName string, filter any, dataForm any) (any, error) {
	if this.client == nil {
		return nil, errors.New("please call Initialize first")
	}

	if err := this.connect(); err != nil {
		return nil, err
	}

	collection := this.client.Database(databaseName).Collection(collectionName)

	document := reflect.New(reflect.TypeOf(dataForm))

	if err := collection.FindOne(this.ctx, filter).Decode(document.Interface()); err != nil {
		return nil, err
	}

	return document.Elem().Interface(), nil
}

// "Find" is returns results value corresponding to the filter argument as an "dataForm" argument array type interface.
//
// ex)
//
//	results_interface, err := client.Find("test_database", "test_collection", bson.M{}, TestStruct{})
//
//	results, ok := results_interface.([]TestStruct)
func (this *Client) Find(databaseName string, collectionName string, filter any, dataForm any) (any, error) {
	if this.client == nil {
		return nil, errors.New("please call Initialize first")
	}

	if err := this.connect(); err != nil {
		return nil, err
	}

	collection := this.client.Database(databaseName).Collection(collectionName)

	cursor, err := collection.Find(this.ctx, filter)
	if err != nil {
		return nil, err
	}

	dataType := reflect.TypeOf(dataForm)
	tempSlice := reflect.MakeSlice(reflect.SliceOf(dataType), 0, 1024)
	results := reflect.New(tempSlice.Type())
	results.Elem().Set(tempSlice)

	if err := cursor.All(this.ctx, results.Interface()); err != nil {
		return nil, err
	}

	if err := cursor.Close(this.ctx); err != nil {
		return nil, err
	}

	return results.Elem().Interface(), nil
}

// InsertOne is insert a one document.
//
// ex) err := client.InsertOne("test_database", "test_collection", TestStruct{})
func (this *Client) InsertOne(databaseName string, collectionName string, document any) error {
	if this.client == nil {
		return errors.New("please call Initialize first")
	}

	if err := this.connect(); err != nil {
		return err
	}

	collection := this.client.Database(databaseName).Collection(collectionName)
	if _, err := collection.InsertOne(this.ctx, document); err != nil {
		return err
	}

	return nil
}

// InsertMany is insert a array type documents.
//
// ex)
//
//	insertData := make([]any, 0)
//
//	insertData = append(insertData, TestStruct{Value1: 1, Value2: "abc"}, TestStruct{Value1: 2, Value2: "def"})
//
//	err := client.InsertMany("test_database", "test_collection", insertData)
func (this *Client) InsertMany(databaseName string, collectionName string, documents []any) error {
	if this.client == nil {
		return errors.New("please call Initialize first")
	}

	if err := this.connect(); err != nil {
		return err
	}

	collection := this.client.Database(databaseName).Collection(collectionName)
	if _, err := collection.InsertMany(this.ctx, documents); err != nil {
		return err
	}

	return nil
}

// UpdateOne is update the one value corresponding to the filter argument with the value of the "update" argument.
//
// ex) err := client.UpdateOne("test_database", "test_collection", bson.M{"value1": 1}, bson.D{{"$set", bson.D{{"value2", "update_value"}}}})
func (this *Client) UpdateOne(databaseName string, collectionName string, filter any, update any) error {
	if this.client == nil {
		return errors.New("please call Initialize first")
	}

	if err := this.connect(); err != nil {
		return err
	}

	collection := this.client.Database(databaseName).Collection(collectionName)
	if _, err := collection.UpdateOne(this.ctx, filter, update); err != nil {
		return err
	}

	return nil
}

// UpdateMany is update the value corresponding to the filter argument with the values of the "update" argument.
//
// ex) err := client.UpdateMany("test_database", "test_collection", bson.M{"value1": 1}, bson.D{{"$set", bson.D{{"value2", "update_value"}}}})
func (this *Client) UpdateMany(databaseName string, collectionName string, filter any, update any) error {
	if this.client == nil {
		return errors.New("please call Initialize first")
	}

	if err := this.connect(); err != nil {
		return err
	}

	collection := this.client.Database(databaseName).Collection(collectionName)
	if _, err := collection.UpdateMany(this.ctx, filter, update); err != nil {
		return err
	}

	return nil
}

// DeleteOne is delete one value corresponding to the filter argument.
//
// ex) err := client.DeleteOne("test_database", "test_collection", bson.M{"value1": 1})
func (this *Client) DeleteOne(databaseName string, collectionName string, filter any) error {
	if this.client == nil {
		return errors.New("please call Initialize first")
	}

	if err := this.connect(); err != nil {
		return err
	}

	collection := this.client.Database(databaseName).Collection(collectionName)
	if _, err := collection.DeleteOne(this.ctx, filter); err != nil {
		return err
	}

	return nil
}

// DeleteMany is delete the values corresponding to the filter argument.
//
// ex) err := client.DeleteMany("test_database", "test_collection", bson.M{})
func (this *Client) DeleteMany(databaseName string, collectionName string, filter any) error {
	if this.client == nil {
		return errors.New("please call Initialize first")
	}

	if err := this.connect(); err != nil {
		return err
	}

	collection := this.client.Database(databaseName).Collection(collectionName)
	if _, err := collection.DeleteMany(this.ctx, filter); err != nil {
		return err
	}

	return nil
}
