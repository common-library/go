package sql_test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"

	sql "github.com/common-library/go/database/sql"
	"github.com/common-library/go/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/localstack"
)

var (
	dynamoDBContainer *localstack.LocalStackContainer
	dynamoDBClient    *sql.Client
	dynamoDBTestDSN   string
	dynamoDBSetupOnce sync.Once
	dynamoDBSetupErr  error
)

func setupSharedDynamoDBContainer() (*sql.Client, error) {
	dynamoDBSetupOnce.Do(func() {
		os.Setenv("AWS_ACCESS_KEY_ID", "test")
		os.Setenv("AWS_SECRET_ACCESS_KEY", "test")
		os.Setenv("AWS_DEFAULT_REGION", "us-east-1")

		ctx := context.Background()

		container, err := localstack.Run(ctx, testutil.LocalstackImage)
		if err != nil {
			dynamoDBSetupErr = fmt.Errorf("failed to start LocalStack container: %w", err)
			return
		}

		dynamoDBContainer = container

		endpoint, err := container.PortEndpoint(ctx, "4566", "http")
		if err != nil {
			dynamoDBSetupErr = fmt.Errorf("failed to get DynamoDB endpoint: %w", err)
			return
		}

		dynamoDBTestDSN = fmt.Sprintf("region=us-east-1;endpoint=%s;access_key_id=test;secret_access_key=test;max_retries=1;timeout=10", endpoint)

		client := &sql.Client{}
		err = client.Open(sql.DriverAmazonDynamoDB, dynamoDBTestDSN, 1)
		if err != nil {
			dynamoDBSetupErr = fmt.Errorf("failed to open DynamoDB client: %w", err)
			return
		}

		dynamoDBClient = client
	})

	return dynamoDBClient, dynamoDBSetupErr
}

func getTestDynamoDBClient(t *testing.T) *sql.Client {
	client, err := setupSharedDynamoDBContainer()
	require.NoError(t, err)
	require.NotNil(t, client)

	return client
}

func getNewTestDynamoDBClient(t *testing.T) *sql.Client {
	_, err := setupSharedDynamoDBContainer()
	require.NoError(t, err)

	client := &sql.Client{}
	err = client.Open(sql.DriverAmazonDynamoDB, dynamoDBTestDSN, 1)
	require.NoError(t, err)

	t.Cleanup(func() {
		client.Close()
	})

	return client
}

func generateUniqueTableName(prefix string) string {
	return fmt.Sprintf("%s_%d_%d", prefix, time.Now().UnixNano(), rand.New(rand.NewSource(time.Now().UnixNano())).Intn(10000))
}

func cleanupTable(_ *testing.T, client *sql.Client, tableName string) {
	dropSQL := fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)
	_ = client.Execute(dropSQL)
}

func TestDynamoDBClient_Open(t *testing.T) {
	t.Parallel()

	client := &sql.Client{}

	_, err := setupSharedDynamoDBContainer()
	require.NoError(t, err)

	err = client.Open(sql.DriverAmazonDynamoDB, dynamoDBTestDSN, 1)
	assert.NoError(t, err)

	defer func() {
		err := client.Close()
		assert.NoError(t, err)
	}()

	assert.Equal(t, sql.DriverAmazonDynamoDB, client.GetDriver())
}

func TestDynamoDBClient_CreateTable(t *testing.T) {
	t.Parallel()

	client := getTestDynamoDBClient(t)
	tableName := generateUniqueTableName("test_users_create")

	defer cleanupTable(t, client, tableName)

	createTableSQL := fmt.Sprintf(`CREATE TABLE %s WITH PK=id:String`, tableName)

	err := client.Execute(createTableSQL)
	assert.NoError(t, err)
}

func TestDynamoDBClient_InsertAndQuery(t *testing.T) {
	t.Parallel()

	client := getNewTestDynamoDBClient(t)
	tableName := generateUniqueTableName("test_products_insert")

	defer cleanupTable(t, client, tableName)

	createTableSQL := fmt.Sprintf(`CREATE TABLE %s WITH PK=id:String`, tableName)
	err := client.Execute(createTableSQL)
	require.NoError(t, err)

	insertSQL := fmt.Sprintf(`INSERT INTO %s VALUE {'id': ?, 'name': ?, 'price': ?, 'category': ?}`, tableName)
	err = client.Execute(insertSQL, "prod1", "Laptop", 999.99, "Electronics")
	assert.NoError(t, err)

	err = client.Execute(insertSQL, "prod2", "Mouse", 29.99, "Electronics")
	assert.NoError(t, err)

	var id, name, category string
	var price float64
	querySQL := fmt.Sprintf(`SELECT id, name, price, category FROM %s WHERE id = 'prod1'`, tableName)
	err = client.QueryRow(querySQL, &id, &name, &price, &category)
	assert.NoError(t, err)
	assert.Equal(t, "prod1", id)
	assert.Equal(t, "Laptop", name)
	assert.Equal(t, 999.99, price)
	assert.Equal(t, "Electronics", category)

	rows, err := client.Query(fmt.Sprintf(`SELECT id, name, price FROM %s`, tableName))
	assert.NoError(t, err)
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		var rowID, rowName string
		var rowPrice float64
		err := rows.Scan(&rowID, &rowName, &rowPrice)
		assert.NoError(t, err)
		rowCount++
	}
	assert.Equal(t, 2, rowCount)
}

func TestDynamoDBClient_PreparedStatements(t *testing.T) {
	t.Parallel()

	client := getNewTestDynamoDBClient(t)
	tableName := generateUniqueTableName("test_orders_prepared")

	defer cleanupTable(t, client, tableName)

	createTableSQL := fmt.Sprintf(`CREATE TABLE %s WITH PK=order_id:String`, tableName)
	err := client.Execute(createTableSQL)
	require.NoError(t, err)

	insertSQL := fmt.Sprintf(`INSERT INTO %s VALUE {'order_id': ?, 'customer_id': ?, 'amount': ?, 'status': ?}`, tableName)
	err = client.SetPrepare(insertSQL)
	require.NoError(t, err)

	err = client.ExecutePrepare("order1", "cust1", 100.50, "pending")
	assert.NoError(t, err)

	selectSQL := fmt.Sprintf(`SELECT order_id, customer_id, amount, status FROM %s WHERE order_id = ?`, tableName)
	err = client.SetPrepare(selectSQL)
	require.NoError(t, err)

	row, err := client.QueryRowPrepare("order1")
	assert.NoError(t, err)

	var orderID, customerID, status string
	var amount float64
	err = row.Scan(&orderID, &customerID, &amount, &status)
	assert.NoError(t, err)
	assert.Equal(t, "order1", orderID)
	assert.Equal(t, "cust1", customerID)
	assert.Equal(t, 100.50, amount)
	assert.Equal(t, "pending", status)
}

func TestDynamoDBClient_Transactions(t *testing.T) {
	t.Parallel()

	client := getNewTestDynamoDBClient(t)
	tableName := generateUniqueTableName("test_accounts_tx")

	defer cleanupTable(t, client, tableName)

	createTableSQL := fmt.Sprintf(`CREATE TABLE %s WITH PK=account_id:String`, tableName)
	err := client.Execute(createTableSQL)
	require.NoError(t, err)

	err = client.BeginTransaction()
	assert.NoError(t, err)

	err = client.EndTransaction(nil)
	assert.NoError(t, err)

	insertSQL := fmt.Sprintf(`INSERT INTO %s VALUE {'account_id': ?, 'balance': ?, 'owner': ?}`, tableName)
	err = client.Execute(insertSQL, "acc1", 1000.0, "Alice")
	assert.NoError(t, err)

	var balance float64
	selectSQL := fmt.Sprintf(`SELECT balance FROM %s WHERE account_id = 'acc1'`, tableName)
	err = client.QueryRow(selectSQL, &balance)
	assert.NoError(t, err)
	assert.Equal(t, 1000.0, balance)
}

func TestDynamoDBClient_TransactionRollback(t *testing.T) {
	t.Parallel()

	client := getNewTestDynamoDBClient(t)
	tableName := generateUniqueTableName("test_inventory_rollback")

	defer cleanupTable(t, client, tableName)

	createTableSQL := fmt.Sprintf(`CREATE TABLE %s WITH PK=item_id:String`, tableName)
	err := client.Execute(createTableSQL)
	require.NoError(t, err)

	err = client.BeginTransaction()
	require.NoError(t, err)

	simulatedError := fmt.Errorf("simulated error")

	err = client.EndTransaction(simulatedError)
	assert.NoError(t, err)

	insertSQL := fmt.Sprintf(`INSERT INTO %s VALUE {'item_id': ?, 'quantity': ?, 'name': ?}`, tableName)
	err = client.Execute(insertSQL, "item1", 50, "Widget")
	assert.NoError(t, err)

	var quantity float64
	selectSQL := fmt.Sprintf(`SELECT quantity FROM %s WHERE item_id = 'item1'`, tableName)
	err = client.QueryRow(selectSQL, &quantity)
	assert.NoError(t, err)
	assert.Equal(t, 50.0, quantity)
}

func TestDynamoDBClient_PreparedTransactions(t *testing.T) {
	t.Parallel()

	client := getNewTestDynamoDBClient(t)
	tableName := generateUniqueTableName("test_logs_prepared_tx")

	defer cleanupTable(t, client, tableName)

	createTableSQL := fmt.Sprintf(`CREATE TABLE %s WITH PK=log_id:String`, tableName)
	err := client.Execute(createTableSQL)
	require.NoError(t, err)

	err = client.BeginTransaction()
	require.NoError(t, err)

	err = client.EndTransaction(nil)
	assert.NoError(t, err)

	insertSQL := fmt.Sprintf(`INSERT INTO %s VALUE {'log_id': ?, 'message': ?, 'timestamp': ?, 'level': ?}`, tableName)
	err = client.SetPrepare(insertSQL)
	require.NoError(t, err)

	err = client.ExecutePrepare("log1", "System started", 1640995200, "INFO")
	assert.NoError(t, err)

	err = client.ExecutePrepare("log2", "User logged in", 1640995260, "INFO")
	assert.NoError(t, err)

	rows, err := client.Query(fmt.Sprintf(`SELECT log_id FROM %s`, tableName))
	assert.NoError(t, err)
	defer rows.Close()

	logCount := 0
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		assert.NoError(t, err)
		logCount++
	}
	assert.Equal(t, 2, logCount)

	assert.True(t, logCount > 0, "준비된 문장으로 데이터가 삽입되어야 합니다")
}

func TestDynamoDBClient_ErrorHandling(t *testing.T) {
	t.Parallel()

	client := &sql.Client{}

	_, err := client.Query("SELECT * FROM test_table")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "please call Open first")

	err = client.QueryRow("SELECT * FROM test_table")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "please call Open first")

	err = client.Execute("INSERT INTO test_table VALUES (1)")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "please call Open first")

	_, err = client.QueryPrepare("value")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "please call SetPrepare first")

	_, err = client.QueryRowPrepare("value")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "please call SetPrepare first")

	err = client.ExecutePrepare("value")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "please call SetPrepare first")

	_, err = client.QueryTransaction("SELECT * FROM test_table")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "please call BeginTransaction first")

	err = client.QueryRowTransaction("SELECT * FROM test_table")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "please call BeginTransaction first")

	err = client.ExecuteTransaction("INSERT INTO test_table VALUES (1)")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "please call BeginTransaction first")

	_, err = client.QueryPrepareTransaction("value")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "please call SetPrepareTransaction first")

	_, err = client.QueryRowPrepareTransaction("value")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "please call SetPrepareTransaction first")

	err = client.ExecutePrepareTransaction("value")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "please call SetPrepareTransaction first")
}

func TestDynamoDBClient_MultipleClients(t *testing.T) {

	_, err := setupSharedDynamoDBContainer()
	require.NoError(t, err)

	client1 := &sql.Client{}
	client2 := &sql.Client{}

	err = client1.Open(sql.DriverAmazonDynamoDB, dynamoDBTestDSN, 1)
	require.NoError(t, err)
	defer client1.Close()

	err = client2.Open(sql.DriverAmazonDynamoDB, dynamoDBTestDSN, 1)
	require.NoError(t, err)
	defer client2.Close()

	tableName := generateUniqueTableName("test_multi_clients")

	createTableSQL := fmt.Sprintf(`CREATE TABLE %s WITH PK=id:String`, tableName)
	err = client1.Execute(createTableSQL)
	require.NoError(t, err)

	defer cleanupTable(t, client1, tableName)

	insertSQL := fmt.Sprintf(`INSERT INTO %s VALUE {'id': ?, 'data': ?}`, tableName)
	err = client1.Execute(insertSQL, "key1", "data1")
	assert.NoError(t, err)

	var id, data string
	selectSQL := fmt.Sprintf(`SELECT id, data FROM %s WHERE id = 'key1'`, tableName)
	err = client2.QueryRow(selectSQL, &id, &data)
	assert.NoError(t, err)
	assert.Equal(t, "key1", id)
	assert.Equal(t, "data1", data)

	err = client2.Execute(insertSQL, "key2", "data2")
	assert.NoError(t, err)

	selectSQL = fmt.Sprintf(`SELECT id, data FROM %s WHERE id = 'key2'`, tableName)
	err = client1.QueryRow(selectSQL, &id, &data)
	assert.NoError(t, err)
	assert.Equal(t, "key2", id)
	assert.Equal(t, "data2", data)

}
