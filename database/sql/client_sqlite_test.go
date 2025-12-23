package sql_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	sqlclient "github.com/common-library/go/database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SQLiteTestSuite struct {
	suite.Suite
	client *sqlclient.Client
	dsn    string
	dbPath string
}

func (suite *SQLiteTestSuite) SetupSuite() {
	tempDir := os.TempDir()
	suite.dbPath = filepath.Join(tempDir, "test_sqlite.db")
	suite.dsn = suite.dbPath

	suite.client = &sqlclient.Client{}
}

func (suite *SQLiteTestSuite) TearDownSuite() {
	if suite.client != nil {
		suite.client.Close()
	}
	if suite.dbPath != "" {
		os.Remove(suite.dbPath)
	}
}

func (suite *SQLiteTestSuite) SetupTest() {
	tempDir := os.TempDir()
	suite.dbPath = filepath.Join(tempDir, fmt.Sprintf("test_sqlite_%d.db", time.Now().UnixNano()))
	suite.dsn = suite.dbPath

	if suite.client != nil {
		suite.client.Close()
	}

	err := suite.client.Open(sqlclient.DriverSQLite, suite.dsn, 1)
	require.NoError(suite.T(), err)

	err = suite.client.Execute(`
		CREATE TABLE test_users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
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

func (suite *SQLiteTestSuite) TearDownTest() {
	if suite.client != nil {
		suite.client.Execute("DELETE FROM test_users")
	}
	if suite.dbPath != "" {
		os.Remove(suite.dbPath)
	}
}

func (suite *SQLiteTestSuite) TestOpen() {
	tempPath := filepath.Join(os.TempDir(), "test_open_sqlite.db")
	defer os.Remove(tempPath)

	client := &sqlclient.Client{}
	err := client.Open(sqlclient.DriverSQLite, tempPath, 1)
	assert.NoError(suite.T(), err)
	defer client.Close()

	assert.Equal(suite.T(), sqlclient.DriverSQLite, client.GetDriver())
}

func (suite *SQLiteTestSuite) TestQuery() {
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

func (suite *SQLiteTestSuite) TestQueryRow() {
	var name string
	var email string
	err := suite.client.QueryRow("SELECT name, email FROM test_users WHERE id = 1", &name, &email)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "John Doe", name)
	assert.Equal(suite.T(), "john@example.com", email)
}

func (suite *SQLiteTestSuite) TestExecute() {
	err := suite.client.Execute("INSERT INTO test_users (name, email) VALUES (?, ?)", "Test User", "test@example.com")
	assert.NoError(suite.T(), err)

	var count int
	err = suite.client.QueryRow("SELECT COUNT(*) FROM test_users", &count)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 3, count)
}

func (suite *SQLiteTestSuite) TestPreparedStatements() {
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

func (suite *SQLiteTestSuite) TestTransaction() {
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

func (suite *SQLiteTestSuite) TestTransactionRollback() {
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

func (suite *SQLiteTestSuite) TestTransactionPreparedStatements() {
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

func (suite *SQLiteTestSuite) TestErrorCases() {
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

	err = client.Open(sqlclient.DriverSQLite, "/nonexistent/path/database.db", 1)
	assert.Error(suite.T(), err)
}

func (suite *SQLiteTestSuite) TestClose() {
	tempPath := filepath.Join(os.TempDir(), "test_close_sqlite.db")
	defer os.Remove(tempPath)

	client := &sqlclient.Client{}

	err := client.Close()
	assert.NoError(suite.T(), err)

	err = client.Open(sqlclient.DriverSQLite, tempPath, 1)
	require.NoError(suite.T(), err)

	err = client.Close()
	assert.NoError(suite.T(), err)

	_, err = client.Query("SELECT 1")
	assert.Error(suite.T(), err)
}

func (suite *SQLiteTestSuite) TestInMemoryDatabase() {
	client := &sqlclient.Client{}
	defer client.Close()

	err := client.Open(sqlclient.DriverSQLite, ":memory:", 1)
	assert.NoError(suite.T(), err)

	err = client.Execute(`
		CREATE TABLE memory_test (
			id INTEGER PRIMARY KEY,
			value TEXT
		)
	`)
	assert.NoError(suite.T(), err)

	err = client.Execute("INSERT INTO memory_test (value) VALUES (?)", "test value")
	assert.NoError(suite.T(), err)

	var value string
	err = client.QueryRow("SELECT value FROM memory_test WHERE id = 1", &value)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "test value", value)
}

func (suite *SQLiteTestSuite) TestConcurrentAccess() {
	tempPath := filepath.Join(os.TempDir(), "test_concurrent_sqlite.db")
	defer os.Remove(tempPath)

	client1 := &sqlclient.Client{}
	defer client1.Close()

	err := client1.Open(sqlclient.DriverSQLite, tempPath, 1)
	assert.NoError(suite.T(), err)

	err = client1.Execute(`
		CREATE TABLE concurrent_test (
			id INTEGER PRIMARY KEY,
			value TEXT
		)
	`)
	assert.NoError(suite.T(), err)

	err = client1.Execute("INSERT INTO concurrent_test (value) VALUES (?)", "concurrent value")
	assert.NoError(suite.T(), err)

	client1.Close()

	client2 := &sqlclient.Client{}
	defer client2.Close()

	err = client2.Open(sqlclient.DriverSQLite, tempPath, 1)
	assert.NoError(suite.T(), err)

	var value string
	err = client2.QueryRow("SELECT value FROM concurrent_test WHERE id = 1", &value)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "concurrent value", value)
}

func (suite *SQLiteTestSuite) TestSQLiteSpecificFeatures() {
	var userVersion int
	err := suite.client.QueryRow("PRAGMA user_version", &userVersion)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 0, userVersion)

	err = suite.client.Execute("PRAGMA user_version = 1")
	assert.NoError(suite.T(), err)

	err = suite.client.QueryRow("PRAGMA user_version", &userVersion)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, userVersion)

	err = suite.client.Execute(`
		CREATE TABLE json_test (
			id INTEGER PRIMARY KEY,
			data TEXT
		)
	`)
	assert.NoError(suite.T(), err)

	jsonData := `{"name": "test", "value": 123}`
	err = suite.client.Execute("INSERT INTO json_test (data) VALUES (?)", jsonData)
	assert.NoError(suite.T(), err)

	var retrievedData string
	err = suite.client.QueryRow("SELECT data FROM json_test WHERE id = 1", &retrievedData)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), jsonData, retrievedData)
}

func TestSQLiteClientSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(SQLiteTestSuite))
}
