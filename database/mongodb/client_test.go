package mongodb_test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/common-library/go/database/mongodb"
	"github.com/common-library/go/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson"
)

type TestStruct struct {
	ID    string `bson:"_id,omitempty"`
	Name  string `bson:"name"`
	Age   int    `bson:"age"`
	Email string `bson:"email"`
}

type MongoDBTestSuite struct {
	suite.Suite
	client   *mongodb.Client
	address  string
	dbName   string
	collName string
}

var (
	mongoContainer testcontainers.Container
	mongoAddress   string
	containerOnce  sync.Once
	cleanupOnce    sync.Once
)

func setupMongoDBContainer() error {
	var err error
	containerOnce.Do(func() {
		ctx := context.Background()

		req := testcontainers.ContainerRequest{
			Image:        testutil.MongoImage,
			ExposedPorts: []string{"27017/tcp"},
			Env: map[string]string{
				"MONGO_INITDB_ROOT_USERNAME": "testuser",
				"MONGO_INITDB_ROOT_PASSWORD": "testpass",
				"MONGO_INITDB_DATABASE":      "testdb",
			},
			WaitingFor: wait.ForAll(
				wait.ForListeningPort("27017/tcp"),
				wait.ForLog("Waiting for connections"),
			).WithDeadline(60 * time.Second),
		}

		mongoContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})

		if err == nil {
			host, err2 := mongoContainer.Host(ctx)
			if err2 != nil {
				err = err2
				return
			}

			port, err2 := mongoContainer.MappedPort(ctx, "27017")
			if err2 != nil {
				err = err2
				return
			}

			mongoAddress = fmt.Sprintf("testuser:testpass@%s:%s", host, port.Port())

			maxRetries := 5
			for i := 0; i < maxRetries; i++ {
				testClient := &mongodb.Client{}
				if initErr := testClient.Initialize(mongoAddress, 30); initErr == nil {
					testClient.Finalize()
					break
				}
				if i == maxRetries-1 {
					err = fmt.Errorf("failed to connect to MongoDB after %d retries", maxRetries)
					return
				}
				time.Sleep(1 * time.Second)
			}
		}
	})
	return err
}

func teardownMongoDBContainer() {
	cleanupOnce.Do(func() {
		if mongoContainer != nil {
			_ = mongoContainer.Terminate(context.Background())
		}
	})
}

func TestMain(m *testing.M) {
	if err := setupMongoDBContainer(); err != nil {
		fmt.Printf("Failed to setup MongoDB container: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()

	teardownMongoDBContainer()

	os.Exit(code)
}

func (suite *MongoDBTestSuite) SetupSuite() {
	suite.dbName = "testdb"
	suite.collName = "test_collection"
	suite.address = mongoAddress
	suite.client = &mongodb.Client{}
}

func (suite *MongoDBTestSuite) TearDownSuite() {
	if suite.client != nil {
		suite.client.Finalize()
	}
}

func (suite *MongoDBTestSuite) SetupTest() {
	err := suite.client.Initialize(suite.address, 30)
	assert.NoError(suite.T(), err)

	suite.cleanupTestData()

	suite.insertTestData()
}

func (suite *MongoDBTestSuite) TearDownTest() {
	suite.cleanupTestData()
	suite.client.Finalize()
}

func (suite *MongoDBTestSuite) cleanupTestData() {
	if suite.client != nil {
		suite.client.DeleteMany(suite.dbName, suite.collName, bson.M{})
	}
}

func (suite *MongoDBTestSuite) insertTestData() {
	testData := []any{
		TestStruct{Name: "John Doe", Age: 30, Email: "john@example.com"},
		TestStruct{Name: "Jane Smith", Age: 25, Email: "jane@example.com"},
	}

	err := suite.client.InsertMany(suite.dbName, suite.collName, testData)
	assert.NoError(suite.T(), err)
}

func (suite *MongoDBTestSuite) TestInitializeAndFinalize() {
	client := &mongodb.Client{}

	err := client.Initialize(suite.address, 30)
	assert.NoError(suite.T(), err)

	err = client.Finalize()
	assert.NoError(suite.T(), err)
}

func (suite *MongoDBTestSuite) TestInsertOne() {
	newDoc := TestStruct{
		Name:  "Alice Brown",
		Age:   28,
		Email: "alice@example.com",
	}

	err := suite.client.InsertOne(suite.dbName, suite.collName, newDoc)
	assert.NoError(suite.T(), err)

	result, err := suite.client.FindOne(suite.dbName, suite.collName, bson.M{"name": "Alice Brown"}, TestStruct{})
	assert.NoError(suite.T(), err)

	foundDoc, ok := result.(TestStruct)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), "Alice Brown", foundDoc.Name)
	assert.Equal(suite.T(), 28, foundDoc.Age)
	assert.Equal(suite.T(), "alice@example.com", foundDoc.Email)
}

func (suite *MongoDBTestSuite) TestInsertMany() {
	suite.cleanupTestData()

	newDocs := []any{
		TestStruct{Name: "Charlie Wilson", Age: 32, Email: "charlie@example.com"},
		TestStruct{Name: "Diana Prince", Age: 29, Email: "diana@example.com"},
	}

	err := suite.client.InsertMany(suite.dbName, suite.collName, newDocs)
	assert.NoError(suite.T(), err)

	results, err := suite.client.Find(suite.dbName, suite.collName, bson.M{}, TestStruct{})
	assert.NoError(suite.T(), err)

	foundDocs, ok := results.([]TestStruct)
	assert.True(suite.T(), ok)
	assert.Len(suite.T(), foundDocs, 2)

	names := make([]string, len(foundDocs))
	for i, doc := range foundDocs {
		names[i] = doc.Name
	}
	assert.Contains(suite.T(), names, "Charlie Wilson")
	assert.Contains(suite.T(), names, "Diana Prince")
}

func (suite *MongoDBTestSuite) TestFindOne() {
	result, err := suite.client.FindOne(suite.dbName, suite.collName, bson.M{"name": "John Doe"}, TestStruct{})
	assert.NoError(suite.T(), err)

	foundDoc, ok := result.(TestStruct)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), "John Doe", foundDoc.Name)
	assert.Equal(suite.T(), 30, foundDoc.Age)
	assert.Equal(suite.T(), "john@example.com", foundDoc.Email)
}

func (suite *MongoDBTestSuite) TestFind() {
	results, err := suite.client.Find(suite.dbName, suite.collName, bson.M{}, TestStruct{})
	assert.NoError(suite.T(), err)

	foundDocs, ok := results.([]TestStruct)
	assert.True(suite.T(), ok)
	assert.Len(suite.T(), foundDocs, 2)

	results, err = suite.client.Find(suite.dbName, suite.collName, bson.M{"age": bson.M{"$gte": 30}}, TestStruct{})
	assert.NoError(suite.T(), err)

	foundDocs, ok = results.([]TestStruct)
	assert.True(suite.T(), ok)
	assert.Len(suite.T(), foundDocs, 1)

	for _, doc := range foundDocs {
		assert.GreaterOrEqual(suite.T(), doc.Age, 30)
	}
}

func (suite *MongoDBTestSuite) TestUpdateOne() {
	filter := bson.M{"name": "John Doe"}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "age", Value: 35}}}}

	err := suite.client.UpdateOne(suite.dbName, suite.collName, filter, update)
	assert.NoError(suite.T(), err)

	result, err := suite.client.FindOne(suite.dbName, suite.collName, filter, TestStruct{})
	assert.NoError(suite.T(), err)

	updatedDoc, ok := result.(TestStruct)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), "John Doe", updatedDoc.Name)
	assert.Equal(suite.T(), 35, updatedDoc.Age)
	assert.Equal(suite.T(), "john@example.com", updatedDoc.Email)
}

func (suite *MongoDBTestSuite) TestUpdateMany() {
	filter := bson.M{"age": bson.M{"$gte": 30}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "email", Value: "updated@newdomain.com"}}}}

	err := suite.client.UpdateMany(suite.dbName, suite.collName, filter, update)
	assert.NoError(suite.T(), err)

	results, err := suite.client.Find(suite.dbName, suite.collName, filter, TestStruct{})
	assert.NoError(suite.T(), err)

	updatedDocs, ok := results.([]TestStruct)
	assert.True(suite.T(), ok)
	assert.Len(suite.T(), updatedDocs, 1)

	for _, doc := range updatedDocs {
		assert.Equal(suite.T(), "updated@newdomain.com", doc.Email)
		assert.GreaterOrEqual(suite.T(), doc.Age, 30)
	}
}

func (suite *MongoDBTestSuite) TestDeleteOne() {
	filter := bson.M{"name": "Jane Smith"}

	err := suite.client.DeleteOne(suite.dbName, suite.collName, filter)
	assert.NoError(suite.T(), err)

	_, err = suite.client.FindOne(suite.dbName, suite.collName, filter, TestStruct{})
	assert.Error(suite.T(), err)

	results, err := suite.client.Find(suite.dbName, suite.collName, bson.M{}, TestStruct{})
	assert.NoError(suite.T(), err)

	remainingDocs, ok := results.([]TestStruct)
	assert.True(suite.T(), ok)
	assert.Len(suite.T(), remainingDocs, 1)
}

func (suite *MongoDBTestSuite) TestDeleteMany() {
	filter := bson.M{"age": bson.M{"$gte": 30}}

	err := suite.client.DeleteMany(suite.dbName, suite.collName, filter)
	assert.NoError(suite.T(), err)

	results, err := suite.client.Find(suite.dbName, suite.collName, bson.M{}, TestStruct{})
	assert.NoError(suite.T(), err)

	remainingDocs, ok := results.([]TestStruct)
	assert.True(suite.T(), ok)
	assert.Len(suite.T(), remainingDocs, 1)

	assert.Equal(suite.T(), "Jane Smith", remainingDocs[0].Name)
	assert.Equal(suite.T(), 25, remainingDocs[0].Age)
}

func (suite *MongoDBTestSuite) TestErrorHandling() {
	client := &mongodb.Client{}

	_, err := client.FindOne(suite.dbName, suite.collName, bson.M{}, TestStruct{})
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "please call Initialize first")

	err = client.InsertOne(suite.dbName, suite.collName, TestStruct{})
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "please call Initialize first")

	err = client.UpdateOne(suite.dbName, suite.collName, bson.M{}, bson.M{})
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "please call Initialize first")

	err = client.DeleteOne(suite.dbName, suite.collName, bson.M{})
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "please call Initialize first")

	err = client.Initialize("invalid:address", 1)
	assert.Error(suite.T(), err)
}

func (suite *MongoDBTestSuite) TestComplexQueries() {
	complexData := []any{
		TestStruct{Name: "David Lee", Age: 40, Email: "david@tech.com"},
		TestStruct{Name: "Sarah Kim", Age: 22, Email: "sarah@tech.com"},
		TestStruct{Name: "Mike Brown", Age: 28, Email: "mike@business.com"},
		TestStruct{Name: "Lisa White", Age: 33, Email: "lisa@business.com"},
	}

	err := suite.client.InsertMany(suite.dbName, suite.collName, complexData)
	assert.NoError(suite.T(), err)

	filter := bson.M{
		"age": bson.M{
			"$gte": 25,
			"$lte": 35,
		},
		"email": bson.M{
			"$regex": "business\\.com$",
		},
	}

	results, err := suite.client.Find(suite.dbName, suite.collName, filter, TestStruct{})
	assert.NoError(suite.T(), err)

	foundDocs, ok := results.([]TestStruct)
	assert.True(suite.T(), ok)
	assert.Len(suite.T(), foundDocs, 2)

	filter2 := bson.M{
		"$or": []bson.M{
			{"age": bson.M{"$lte": 30}},
			{"email": bson.M{"$regex": "business\\.com$"}},
		},
	}

	results2, err := suite.client.Find(suite.dbName, suite.collName, filter2, TestStruct{})
	assert.NoError(suite.T(), err)

	foundDocs2, ok := results2.([]TestStruct)
	assert.True(suite.T(), ok)
	assert.Greater(suite.T(), len(foundDocs2), 0)
}

func (suite *MongoDBTestSuite) TestConcurrentOperations() {
	const numGoroutines = 10
	const operationsPerGoroutine = 5

	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(routineID int) {
			defer wg.Done()

			client := &mongodb.Client{}
			err := client.Initialize(suite.address, 30)
			if err != nil {
				suite.T().Errorf("Failed to initialize client in goroutine %d: %v", routineID, err)
				return
			}
			defer client.Finalize()

			for j := 0; j < operationsPerGoroutine; j++ {
				doc := TestStruct{
					Name:  fmt.Sprintf("Concurrent User %d-%d", routineID, j),
					Age:   20 + (routineID*operationsPerGoroutine + j),
					Email: fmt.Sprintf("user%d_%d@concurrent.com", routineID, j),
				}

				err := client.InsertOne(suite.dbName, "concurrent_test", doc)
				if err != nil {
					suite.T().Errorf("Failed to insert in goroutine %d: %v", routineID, err)
				}
			}
		}(i)
	}

	wg.Wait()

	results, err := suite.client.Find(suite.dbName, "concurrent_test", bson.M{}, TestStruct{})
	assert.NoError(suite.T(), err)

	foundDocs, ok := results.([]TestStruct)
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), numGoroutines*operationsPerGoroutine, len(foundDocs))

	err = suite.client.DeleteMany(suite.dbName, "concurrent_test", bson.M{})
	assert.NoError(suite.T(), err)
}

func TestMongoDBSuite(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Integration tests are skipped")
	}

	t.Parallel()

	suite.Run(t, new(MongoDBTestSuite))
}
