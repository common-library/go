package sql_test

import (
	"context"
	"os"
	"testing"
	"time"

	sqlclient "github.com/common-library/go/database/sql"
	"github.com/common-library/go/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type PostgreSQLTestSuite struct {
	suite.Suite
	container *postgres.PostgresContainer
	client    *sqlclient.Client
	dsn       string
}

func (suite *PostgreSQLTestSuite) SetupSuite() {
	ctx := context.Background()

	container, err := postgres.Run(ctx,
		testutil.PostgresImage,
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		postgres.BasicWaitStrategies(),
	)
	suite.Require().NoError(err)
	suite.container = container

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	suite.Require().NoError(err)
	suite.dsn = dsn
	suite.client = &sqlclient.Client{}

	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		testClient := &sqlclient.Client{}
		if err := testClient.Open(sqlclient.DriverPostgreSQL, suite.dsn, 1); err == nil {
			testClient.Close()
			break
		}
		if i == maxRetries-1 {
			suite.T().Fatalf("Failed to connect to PostgreSQL after %d retries", maxRetries)
		}
		backoff := time.Duration(50<<uint(i)) * time.Millisecond
		if backoff > time.Second {
			backoff = time.Second
		}
		time.Sleep(backoff)
	}

	suite.setupTestTable()
}

func (suite *PostgreSQLTestSuite) TearDownSuite() {
	if suite.client != nil {
		suite.client.Close()
	}

	if suite.container != nil {
		ctx := context.Background()
		suite.container.Terminate(ctx)
	}
}

func (suite *PostgreSQLTestSuite) SetupTest() {
	err := suite.client.Open(sqlclient.DriverPostgreSQL, suite.dsn, 5)
	suite.Require().NoError(err)
}

func (suite *PostgreSQLTestSuite) TearDownTest() {
	suite.cleanupTestData()
	suite.client.Close()
}

func (suite *PostgreSQLTestSuite) setupTestTable() {
	client := &sqlclient.Client{}
	err := client.Open(sqlclient.DriverPostgreSQL, suite.dsn, 5)
	suite.Require().NoError(err)
	defer client.Close()

	createTableQuery := `
		CREATE TABLE IF NOT EXISTS test_users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			age INTEGER,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`

	err = client.Execute(createTableQuery)
	suite.Require().NoError(err)
}

func (suite *PostgreSQLTestSuite) cleanupTestData() {
	if suite.client != nil {
		suite.client.Execute("TRUNCATE TABLE test_users RESTART IDENTITY")
	}
}

func (suite *PostgreSQLTestSuite) TestOpen() {
	assert.Equal(suite.T(), sqlclient.DriverPostgreSQL, suite.client.GetDriver())
}

func (suite *PostgreSQLTestSuite) TestExecute() {
	insertQuery := "INSERT INTO test_users (name, email, age) VALUES ($1, $2, $3)"
	err := suite.client.Execute(insertQuery, "John Doe", "john@example.com", 30)
	assert.NoError(suite.T(), err)

	var count int
	err = suite.client.QueryRow("SELECT COUNT(*) FROM test_users", &count)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, count)
}

func (suite *PostgreSQLTestSuite) TestQuery() {
	suite.client.Execute("INSERT INTO test_users (name, email, age) VALUES ($1, $2, $3)", "Alice", "alice@example.com", 25)
	suite.client.Execute("INSERT INTO test_users (name, email, age) VALUES ($1, $2, $3)", "Bob", "bob@example.com", 35)

	rows, err := suite.client.Query("SELECT name, email, age FROM test_users ORDER BY name")
	assert.NoError(suite.T(), err)
	defer rows.Close()

	var users []struct {
		Name  string
		Email string
		Age   int
	}

	for rows.Next() {
		var user struct {
			Name  string
			Email string
			Age   int
		}
		err := rows.Scan(&user.Name, &user.Email, &user.Age)
		assert.NoError(suite.T(), err)
		users = append(users, user)
	}

	assert.Len(suite.T(), users, 2)
	assert.Equal(suite.T(), "Alice", users[0].Name)
	assert.Equal(suite.T(), "Bob", users[1].Name)
}

func (suite *PostgreSQLTestSuite) TestTransaction() {
	err := suite.client.BeginTransaction()
	assert.NoError(suite.T(), err)

	err = suite.client.ExecuteTransaction("INSERT INTO test_users (name, email, age) VALUES ($1, $2, $3)", "Henry", "henry@example.com", 31)
	assert.NoError(suite.T(), err)

	err = suite.client.EndTransaction(nil)
	assert.NoError(suite.T(), err)

	var count int
	rows, err := suite.client.Query("SELECT COUNT(*) FROM test_users WHERE name = $1", "Henry")
	assert.NoError(suite.T(), err)
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&count)
		assert.NoError(suite.T(), err)
	}
	assert.Equal(suite.T(), 1, count)
}

func TestPostgreSQLSuite(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Integration tests are skipped")
	}

	t.Parallel()

	suite.Run(t, new(PostgreSQLTestSuite))
}
