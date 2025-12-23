// Package mongodb provides a MongoDB client wrapper with simplified operations and automatic reconnection.
//
// This package offers a convenient wrapper around the official MongoDB Go driver,
// providing simplified method signatures for common operations and automatic connection management.
//
// Features:
//   - Automatic connection and reconnection handling
//   - CRUD operations with type-safe results
//   - Aggregation pipeline support
//   - Bulk write operations
//   - Index management
//   - Context-based timeout control
//
// Example:
//
//	var client mongodb.Client
//	client.Initialize("localhost:27017", 10*time.Second)
//	defer client.Finalize()
//	client.InsertOne("mydb", "users", bson.M{"name": "Alice", "age": 30})
//	result, _ := client.FindOne("mydb", "users", bson.M{"name": "Alice"}, User{})
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

// connect establishes a connection to the MongoDB server or reuses an existing connection.
//
// This method implements automatic reconnection by checking if an existing connection
// is alive via ping. If the connection is dead or doesn't exist, it creates a new one.
//
// Returns error if connection or ping fails.
//
// The method is called automatically by all database operations.
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

// disConnect closes the MongoDB connection and cancels the context.
//
// Returns error if disconnection fails.
//
// This method is called by Finalize and during reconnection attempts.
// It safely handles nil clients and cleans up context cancellation functions.
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

// Initialize initializes the MongoDB client with connection settings.
//
// Parameters:
//   - address: MongoDB server address in format "host:port" (e.g., "localhost:27017")
//   - timeout: Operation timeout duration for database operations
//
// Returns error if initial connection fails.
//
// The client automatically reconnects on subsequent operations if the connection is lost.
// Call Finalize() to properly close the connection when done.
//
// Example:
//
//	var client mongodb.Client
//	err := client.Initialize("localhost:27017", 10*time.Second)
//	defer client.Finalize()
func (c *Client) Initialize(address string, timeout time.Duration) error {
	c.address = address
	c.timeout = timeout

	return c.connect()
}

// Finalize closes the MongoDB connection and releases resources.
//
// Returns error if disconnection fails.
//
// This method should be called when the client is no longer needed,
// typically using defer after Initialize.
//
// Example:
//
//	var client mongodb.Client
//	client.Initialize("localhost:27017", 10*time.Second)
//	defer client.Finalize()
func (c *Client) Finalize() error {
	return c.disConnect()
}

// FindOne finds a single document matching the filter and returns it as the specified type.
//
// Parameters:
//   - databaseName: Name of the database
//   - collectionName: Name of the collection
//   - filter: BSON filter document (e.g., bson.M{"name": "Alice"})
//   - dataForm: Template for the result type (e.g., User{})
//
// Returns the found document as an interface that must be type-asserted to the dataForm type.
// Returns error if client is not initialized, connection fails, or no document is found.
//
// The method uses reflection to create a properly typed result from the dataForm template.
//
// Example:
//
//	result, err := client.FindOne("mydb", "users", bson.M{"name": "Alice"}, User{})
//	if err != nil {
//		return err
//	}
//	user, ok := result.(User)
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

// Find finds all documents matching the filter and returns them as a slice of the specified type.
//
// Parameters:
//   - databaseName: Name of the database
//   - collectionName: Name of the collection
//   - filter: BSON filter document (e.g., bson.M{"age": bson.M{"$gte": 25}})
//   - dataForm: Template for the result element type (e.g., User{})
//
// Returns a slice of documents as an interface that must be type-asserted to []dataForm type.
// Returns error if client is not initialized, connection fails, or query fails.
//
// Use an empty filter bson.M{} to retrieve all documents in the collection.
// The method uses reflection to create a properly typed slice from the dataForm template.
//
// Example:
//
//	results, err := client.Find("mydb", "users", bson.M{"age": bson.M{"$gte": 25}}, User{})
//	if err != nil {
//		return err
//	}
//	users, ok := results.([]User)
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

// InsertOne inserts a single document into the collection.
//
// Parameters:
//   - databaseName: Name of the database
//   - collectionName: Name of the collection
//   - document: Document to insert (struct, bson.M, or bson.D)
//
// Returns error if client is not initialized, connection fails, or insertion fails.
//
// The document can be a struct with bson tags or a bson.M/bson.D map.
//
// Example:
//
//	err := client.InsertOne("mydb", "users", User{ID: 1, Name: "Alice", Age: 30})
//
//	// Or with bson.M
//	err = client.InsertOne("mydb", "users", bson.M{"_id": 1, "name": "Alice", "age": 30})
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

// InsertMany inserts multiple documents into the collection in a single operation.
//
// Parameters:
//   - databaseName: Name of the database
//   - collectionName: Name of the collection
//   - documents: Slice of documents to insert ([]any containing structs, bson.M, or bson.D)
//
// Returns error if client is not initialized, connection fails, or insertion fails.
//
// All documents are inserted in a single batch operation for efficiency.
//
// Example:
//
//	docs := []any{
//		User{ID: 1, Name: "Alice", Age: 30},
//		User{ID: 2, Name: "Bob", Age: 25},
//		bson.M{"_id": 3, "name": "Charlie", "age": 35},
//	}
//	err := client.InsertMany("mydb", "users", docs)
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

// UpdateOne updates a single document matching the filter.
//
// Parameters:
//   - databaseName: Name of the database
//   - collectionName: Name of the collection
//   - filter: BSON filter to match the document (e.g., bson.M{"_id": 1})
//   - update: Update operations using MongoDB update operators (e.g., bson.D{{"$set", bson.D{{"age", 31}}}})
//
// Returns error if client is not initialized, connection fails, or update fails.
//
// Only the first matching document is updated.
// Use update operators like $set, $inc, $push, etc.
//
// Example:
//
//	err := client.UpdateOne(
//		"mydb", "users",
//		bson.M{"name": "Alice"},
//		bson.D{{"$set", bson.D{{"age", 31}}}},
//	)
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

// UpdateMany updates all documents matching the filter.
//
// Parameters:
//   - databaseName: Name of the database
//   - collectionName: Name of the collection
//   - filter: BSON filter to match documents (e.g., bson.M{"age": bson.M{"$lt": 30}})
//   - update: Update operations using MongoDB update operators (e.g., bson.D{{"$inc", bson.D{{"age", 1}}}})
//
// Returns error if client is not initialized, connection fails, or update fails.
//
// All matching documents are updated in a single operation.
// Use update operators like $set, $inc, $push, etc.
//
// Example:
//
//	err := client.UpdateMany(
//		"mydb", "users",
//		bson.M{"age": bson.M{"$lt": 30}},
//		bson.D{{"$inc", bson.D{{"age", 1}}}},
//	)
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

// DeleteOne deletes a single document matching the filter.
//
// Parameters:
//   - databaseName: Name of the database
//   - collectionName: Name of the collection
//   - filter: BSON filter to match the document (e.g., bson.M{"_id": 1})
//
// Returns error if client is not initialized, connection fails, or deletion fails.
//
// Only the first matching document is deleted.
//
// Example:
//
//	err := client.DeleteOne("mydb", "users", bson.M{"name": "Alice"})
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

// DeleteMany deletes all documents matching the filter.
//
// Parameters:
//   - databaseName: Name of the database
//   - collectionName: Name of the collection
//   - filter: BSON filter to match documents (e.g., bson.M{"age": bson.M{"$lt": 25}})
//
// Returns error if client is not initialized, connection fails, or deletion fails.
//
// All matching documents are deleted in a single operation.
// Use an empty filter bson.M{} to delete all documents in the collection.
//
// Example:
//
//	err := client.DeleteMany("mydb", "users", bson.M{"age": bson.M{"$lt": 25}})
//
//	// Delete all documents
//	err = client.DeleteMany("mydb", "users", bson.M{})
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
