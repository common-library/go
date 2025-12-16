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
func (c *Client) connect() error {
	if c.client != nil && c.client.Ping(c.ctx, readpref.Primary()) == nil {
		return nil
	}

	c.disConnect()

	if client, err := mongo.Connect(c.ctx, options.Client().ApplyURI("mongodb://"+c.address)); err != nil {
		return err
	} else {
		c.client = client
	}

	c.ctx = context.TODO()
	if c.timeout > 0 {
		c.ctx, c.ctxCancelFunc = context.WithTimeout(context.Background(), c.timeout*time.Second)
	}

	return c.client.Ping(c.ctx, readpref.Primary())
}

// disconnect is disconnect.
//
// ex) err := client.disconnect()
func (c *Client) disConnect() error {
	if c.client == nil {
		return nil
	}

	err := c.client.Disconnect(c.ctx)

	if c.ctxCancelFunc != nil {
		c.ctxCancelFunc()
	}

	c.client = nil

	return err
}

// Initialize is initialize.
//
// ex) err := client.Initialize("localhost:27017", 10)
func (c *Client) Initialize(address string, timeout time.Duration) error {
	c.address = address
	c.timeout = timeout

	return c.connect()
}

// Finalize is finalize.
//
// ex) err := client.Finalize()
func (c *Client) Finalize() error {
	return c.disConnect()
}

// "FindOne" is returns one result value corresponding to the filter argument as an "dataForm" argument type interface.
//
//	ex)
//
//	 result_interface, err := client.FindOne("test_database", "test_collection", bson.M{"value1": 1}, TestStruct{})
//
//	 result, ok := result_interface.(TestStruct)
func (c *Client) FindOne(databaseName, collectionName string, filter any, dataForm any) (any, error) {
	if c.client == nil {
		return nil, errors.New("please call Initialize first")
	}

	if err := c.connect(); err != nil {
		return nil, err
	}

	collection := c.client.Database(databaseName).Collection(collectionName)

	document := reflect.New(reflect.TypeOf(dataForm))

	if err := collection.FindOne(c.ctx, filter).Decode(document.Interface()); err != nil {
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
func (c *Client) Find(databaseName, collectionName string, filter, dataForm any) (any, error) {
	if c.client == nil {
		return nil, errors.New("please call Initialize first")
	}

	if err := c.connect(); err != nil {
		return nil, err
	}

	collection := c.client.Database(databaseName).Collection(collectionName)

	cursor, err := collection.Find(c.ctx, filter)
	if err != nil {
		return nil, err
	}

	dataType := reflect.TypeOf(dataForm)
	tempSlice := reflect.MakeSlice(reflect.SliceOf(dataType), 0, 1024)
	results := reflect.New(tempSlice.Type())
	results.Elem().Set(tempSlice)

	if err := cursor.All(c.ctx, results.Interface()); err != nil {
		return nil, err
	}

	if err := cursor.Close(c.ctx); err != nil {
		return nil, err
	}

	return results.Elem().Interface(), nil
}

// InsertOne is insert a one document.
//
// ex) err := client.InsertOne("test_database", "test_collection", TestStruct{})
func (c *Client) InsertOne(databaseName, collectionName string, document any) error {
	if c.client == nil {
		return errors.New("please call Initialize first")
	}

	if err := c.connect(); err != nil {
		return err
	}

	collection := c.client.Database(databaseName).Collection(collectionName)
	if _, err := collection.InsertOne(c.ctx, document); err != nil {
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
func (c *Client) InsertMany(databaseName, collectionName string, documents []any) error {
	if c.client == nil {
		return errors.New("please call Initialize first")
	}

	if err := c.connect(); err != nil {
		return err
	}

	collection := c.client.Database(databaseName).Collection(collectionName)
	if _, err := collection.InsertMany(c.ctx, documents); err != nil {
		return err
	}

	return nil
}

// UpdateOne is update the one value corresponding to the filter argument with the value of the "update" argument.
//
// ex) err := client.UpdateOne("test_database", "test_collection", bson.M{"value1": 1}, bson.D{{"$set", bson.D{{"value2", "update_value"}}}})
func (c *Client) UpdateOne(databaseName, collectionName string, filter, update any) error {
	if c.client == nil {
		return errors.New("please call Initialize first")
	}

	if err := c.connect(); err != nil {
		return err
	}

	collection := c.client.Database(databaseName).Collection(collectionName)
	if _, err := collection.UpdateOne(c.ctx, filter, update); err != nil {
		return err
	}

	return nil
}

// UpdateMany is update the value corresponding to the filter argument with the values of the "update" argument.
//
// ex) err := client.UpdateMany("test_database", "test_collection", bson.M{"value1": 1}, bson.D{{"$set", bson.D{{"value2", "update_value"}}}})
func (c *Client) UpdateMany(databaseName, collectionName string, filter, update any) error {
	if c.client == nil {
		return errors.New("please call Initialize first")
	}

	if err := c.connect(); err != nil {
		return err
	}

	collection := c.client.Database(databaseName).Collection(collectionName)
	if _, err := collection.UpdateMany(c.ctx, filter, update); err != nil {
		return err
	}

	return nil
}

// DeleteOne is delete one value corresponding to the filter argument.
//
// ex) err := client.DeleteOne("test_database", "test_collection", bson.M{"value1": 1})
func (c *Client) DeleteOne(databaseName, collectionName string, filter any) error {
	if c.client == nil {
		return errors.New("please call Initialize first")
	}

	if err := c.connect(); err != nil {
		return err
	}

	collection := c.client.Database(databaseName).Collection(collectionName)
	if _, err := collection.DeleteOne(c.ctx, filter); err != nil {
		return err
	}

	return nil
}

// DeleteMany is delete the values corresponding to the filter argument.
//
// ex) err := client.DeleteMany("test_database", "test_collection", bson.M{})
func (c *Client) DeleteMany(databaseName, collectionName string, filter any) error {
	if c.client == nil {
		return errors.New("please call Initialize first")
	}

	if err := c.connect(); err != nil {
		return err
	}

	collection := c.client.Database(databaseName).Collection(collectionName)
	if _, err := collection.DeleteMany(c.ctx, filter); err != nil {
		return err
	}

	return nil
}
