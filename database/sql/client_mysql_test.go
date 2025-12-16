package sql_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	sqlclient "github.com/common-library/go/database/sql"
	"github.com/common-library/go/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/wait"
)

type MySQLTestSuite struct {
	suite.Suite
	container *mysql.MySQLContainer
	client    *sqlclient.Client
	dsn       string
}

func (suite *MySQLTestSuite) SetupSuite() {
	ctx := context.Background()

	mysqlContainer, err := mysql.Run(ctx,
		testutil.MySQLImage,
		mysql.WithDatabase("testdb"),
		mysql.WithUsername("testuser"),
		mysql.WithPassword("testpass"),
		testcontainers.WithEnv(map[string]string{
			"MYSQL_ROOT_PASSWORD":      "rootpass",
			"MYSQL_INITDB_SKIP_TZINFO": "1",
		}),
		testcontainers.WithWaitStrategy(
			wait.ForLog("ready for connections").
				WithOccurrence(2).
				WithStartupTimeout(90*time.Second).
				WithPollInterval(500*time.Millisecond),
		),
	)
	require.NoError(suite.T(), err)

	suite.container = mysqlContainer

	host, err := mysqlContainer.Host(ctx)
	require.NoError(suite.T(), err)

	port, err := mysqlContainer.MappedPort(ctx, "3306")
	require.NoError(suite.T(), err)

	suite.dsn = fmt.Sprintf("testuser:testpass@tcp(%s:%s)/testdb?parseTime=true&timeout=30s&readTimeout=30s&writeTimeout=30s", host, port.Port())

	suite.client = &sqlclient.Client{}

	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		testClient := &sqlclient.Client{}
		if err := testClient.Open(sqlclient.DriverMySQL, suite.dsn, 1); err == nil {
			testClient.Close()
			break
		}
		if i == maxRetries-1 {
			suite.T().Fatalf("Failed to connect to MySQL after %d retries", maxRetries)
		}
		backoff := time.Duration(50<<uint(i)) * time.Millisecond
		if backoff > time.Second {
			backoff = time.Second
		}
		time.Sleep(backoff)
	}
}

func (suite *MySQLTestSuite) TearDownSuite() {
	if suite.client != nil {
		suite.client.Close()
	}
	if suite.container != nil {
		suite.container.Terminate(context.Background())
	}
}

func (suite *MySQLTestSuite) SetupTest() {
	if suite.client != nil {
		suite.client.Close()
	}

	err := suite.client.Open(sqlclient.DriverMySQL, suite.dsn, 5)
	require.NoError(suite.T(), err)

	suite.client.Execute("DROP TABLE IF EXISTS test_users")

	err = suite.client.Execute(`
		CREATE TABLE test_users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	require.NoError(suite.T(), err)

	err = suite.client.Execute(`
		INSERT INTO test_users (name, email) VALUES
		('John Doe', 'john@example.com'),
		('Jane Smith', 'jane@example.com')
	`)
	require.NoError(suite.T(), err)
}

func (suite *MySQLTestSuite) TearDownTest() {
	if suite.client != nil {
		suite.client.Execute("DROP TABLE IF EXISTS test_users")
	}
}

func (suite *MySQLTestSuite) TestOpen() {
	client := &sqlclient.Client{}
	err := client.Open(sqlclient.DriverMySQL, suite.dsn, 5)
	assert.NoError(suite.T(), err)
	defer client.Close()

	assert.Equal(suite.T(), sqlclient.DriverMySQL, client.GetDriver())
}

func (suite *MySQLTestSuite) TestQuery() {
	rows, err := suite.client.Query("SELECT id, name, email FROM test_users ORDER BY id")
	require.NoError(suite.T(), err)
	defer rows.Close()

	var users []struct {
		ID    int
		Name  string
		Email string
	}

	for rows.Next() {
		var user struct {
			ID    int
			Name  string
			Email string
		}
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		require.NoError(suite.T(), err)
		users = append(users, user)
	}

	assert.Len(suite.T(), users, 2)
	assert.Equal(suite.T(), "John Doe", users[0].Name)
	assert.Equal(suite.T(), "jane@example.com", users[1].Email)
}

func (suite *MySQLTestSuite) TestQueryRow() {
	var name string
	var email string
	err := suite.client.QueryRow("SELECT name, email FROM test_users WHERE id = 1", &name, &email)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "John Doe", name)
	assert.Equal(suite.T(), "john@example.com", email)
}

func (suite *MySQLTestSuite) TestExecute() {
	err := suite.client.Execute("INSERT INTO test_users (name, email) VALUES (?, ?)", "Test User", "test@example.com")
	assert.NoError(suite.T(), err)

	var count int
	err = suite.client.QueryRow("SELECT COUNT(*) FROM test_users", &count)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 3, count)
}

func (suite *MySQLTestSuite) TestPreparedStatements() {
	err := suite.client.SetPrepare("SELECT name, email FROM test_users WHERE id = ?")
	assert.NoError(suite.T(), err)

	rows, err := suite.client.QueryPrepare(1)
	require.NoError(suite.T(), err)
	defer rows.Close()

	require.True(suite.T(), rows.Next())
	var name, email string
	err = rows.Scan(&name, &email)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "John Doe", name)
	rows.Close()

	err = suite.client.SetPrepare("SELECT name, email FROM test_users WHERE id = ?")
	assert.NoError(suite.T(), err)

	row, err := suite.client.QueryRowPrepare(2)
	require.NoError(suite.T(), err)

	var name2, email2 string
	err = row.Scan(&name2, &email2)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Jane Smith", name2)

	err = suite.client.SetPrepare("INSERT INTO test_users (name, email) VALUES (?, ?)")
	assert.NoError(suite.T(), err)

	err = suite.client.ExecutePrepare("Prepared User", "prepared@example.com")
	assert.NoError(suite.T(), err)
}

func (suite *MySQLTestSuite) TestTransaction() {
	err := suite.client.BeginTransaction()
	assert.NoError(suite.T(), err)

	err = suite.client.ExecuteTransaction("INSERT INTO test_users (name, email) VALUES (?, ?)", "Transaction User", "transaction@example.com")
	assert.NoError(suite.T(), err)

	var count int
	err = suite.client.QueryRowTransaction("SELECT COUNT(*) FROM test_users", &count)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 3, count)

	err = suite.client.EndTransaction(nil)
	assert.NoError(suite.T(), err)

	err = suite.client.QueryRow("SELECT COUNT(*) FROM test_users", &count)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 3, count)
}

func (suite *MySQLTestSuite) TestTransactionRollback() {
	var initialCount int
	err := suite.client.QueryRow("SELECT COUNT(*) FROM test_users", &initialCount)
	assert.NoError(suite.T(), err)

	err = suite.client.BeginTransaction()
	assert.NoError(suite.T(), err)

	err = suite.client.ExecuteTransaction("INSERT INTO test_users (name, email) VALUES (?, ?)", "Rollback User", "rollback@example.com")
	assert.NoError(suite.T(), err)

	err = suite.client.EndTransaction(fmt.Errorf("intentional rollback"))
	assert.NoError(suite.T(), err)

	var finalCount int
	err = suite.client.QueryRow("SELECT COUNT(*) FROM test_users", &finalCount)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), initialCount, finalCount)
}

func (suite *MySQLTestSuite) TestTransactionPreparedStatements() {
	err := suite.client.BeginTransaction()
	assert.NoError(suite.T(), err)

	err = suite.client.SetPrepareTransaction("INSERT INTO test_users (name, email) VALUES (?, ?)")
	assert.NoError(suite.T(), err)

	err = suite.client.ExecutePrepareTransaction("TX Prepared User", "txprepared@example.com")
	assert.NoError(suite.T(), err)

	err = suite.client.SetPrepareTransaction("SELECT name FROM test_users WHERE email = ?")
	assert.NoError(suite.T(), err)

	rows, err := suite.client.QueryPrepareTransaction("txprepared@example.com")
	require.NoError(suite.T(), err)
	defer rows.Close()

	require.True(suite.T(), rows.Next())
	var name string
	err = rows.Scan(&name)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "TX Prepared User", name)
	rows.Close()

	err = suite.client.SetPrepareTransaction("SELECT name FROM test_users WHERE email = ?")
	assert.NoError(suite.T(), err)

	row, err := suite.client.QueryRowPrepareTransaction("txprepared@example.com")
	require.NoError(suite.T(), err)

	var name2 string
	err = row.Scan(&name2)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "TX Prepared User", name2)

	err = suite.client.EndTransaction(nil)
	assert.NoError(suite.T(), err)
}

func (suite *MySQLTestSuite) TestErrorCases() {
	client := &sqlclient.Client{}

	_, err := client.Query("SELECT 1")
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "please call Open first")

	err = client.QueryRow("SELECT 1", new(int))
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "please call Open first")

	err = client.Execute("SELECT 1")
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "please call Open first")

	err = client.Open(sqlclient.DriverMySQL, "invalid_dsn", 1)
	assert.Error(suite.T(), err)
}

func (suite *MySQLTestSuite) TestClose() {
	client := &sqlclient.Client{}

	err := client.Close()
	assert.NoError(suite.T(), err)

	err = client.Open(sqlclient.DriverMySQL, suite.dsn, 1)
	require.NoError(suite.T(), err)

	err = client.Close()
	assert.NoError(suite.T(), err)

	_, err = client.Query("SELECT 1")
	assert.Error(suite.T(), err)
}

func TestMySQLClientSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(MySQLTestSuite))
}
