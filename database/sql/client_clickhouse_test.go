package sql_test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/common-library/go/database/sql"
	"github.com/common-library/go/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/clickhouse"
)

var (
	clickhouseContainer *clickhouse.ClickHouseContainer
	clickhouseDSN       string
	containerOnce       sync.Once
	cleanupOnce         sync.Once
)

func setupClickHouseContainer() error {
	var err error
	containerOnce.Do(func() {
		ctx := context.Background()

		clickhouseContainer, err = clickhouse.Run(ctx,
			testutil.ClickHouseImage,
			clickhouse.WithUsername("testuser"),
			clickhouse.WithPassword("testpass"),
			clickhouse.WithDatabase("testdb"),
		)
		if err != nil {
			return
		}

		clickhouseDSN, err = clickhouseContainer.ConnectionString(ctx)
		if err != nil {
			return
		}
	})
	return err
}

func teardownClickHouseContainer() {
	cleanupOnce.Do(func() {
		if clickhouseContainer != nil {
			_ = clickhouseContainer.Terminate(context.Background())
		}
	})
}

func TestMain(m *testing.M) {
	if err := setupClickHouseContainer(); err != nil {
		fmt.Printf("Failed to setup ClickHouse container: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()

	teardownClickHouseContainer()

	os.Exit(code)
}

func getTestClient(t *testing.T) *sql.Client {
	client := &sql.Client{}
	err := client.Open(sql.DriverClickHouse, clickhouseDSN, 10)
	require.NoError(t, err)
	return client
}

func createTestTable(t *testing.T, client *sql.Client, tableName string) {
	dropQuery := fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)
	err := client.Execute(dropQuery)
	require.NoError(t, err)

	createQuery := fmt.Sprintf(`
		CREATE TABLE %s (
			id UInt32,
			name String,
			created_at DateTime
		) ENGINE = Memory
	`, tableName)
	err = client.Execute(createQuery)
	require.NoError(t, err)
}

func TestClickHouseClient_OpenAndClose(t *testing.T) {
	t.Parallel()

	client := &sql.Client{}

	err := client.Open(sql.DriverClickHouse, clickhouseDSN, 10)
	assert.NoError(t, err)

	assert.Equal(t, sql.DriverClickHouse, client.GetDriver())

	err = client.Close()
	assert.NoError(t, err)
}

func TestClickHouseClient_Execute(t *testing.T) {
	t.Parallel()

	client := getTestClient(t)
	defer client.Close()

	tableName := "test_execute"
	createTestTable(t, client, tableName)

	insertQuery := fmt.Sprintf("INSERT INTO %s (id, name, created_at) VALUES (?, ?, ?)", tableName)
	err := client.Execute(insertQuery, 1, "test_user", time.Now())
	assert.NoError(t, err)

	for i := 2; i <= 3; i++ {
		err = client.Execute(insertQuery, i, fmt.Sprintf("user_%d", i), time.Now())
		assert.NoError(t, err)
	}
}

func TestClickHouseClient_Query(t *testing.T) {
	t.Parallel()

	client := getTestClient(t)
	defer client.Close()

	tableName := "test_query"
	createTestTable(t, client, tableName)

	insertQuery := fmt.Sprintf("INSERT INTO %s (id, name, created_at) VALUES (?, ?, ?)", tableName)
	for i := 1; i <= 2; i++ {
		err := client.Execute(insertQuery, i, fmt.Sprintf("user_%d", i), time.Now())
		require.NoError(t, err)
	}

	selectQuery := fmt.Sprintf("SELECT id, name FROM %s ORDER BY id", tableName)
	rows, err := client.Query(selectQuery)
	require.NoError(t, err)
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id uint32
		var name string
		err := rows.Scan(&id, &name)
		require.NoError(t, err)
		count++

		assert.Equal(t, uint32(count), id)
		assert.Equal(t, fmt.Sprintf("user_%d", count), name)
	}
	assert.Equal(t, 2, count)
}

func TestClickHouseClient_QueryRow(t *testing.T) {
	t.Parallel()

	client := getTestClient(t)
	defer client.Close()

	tableName := "test_query_row"
	createTestTable(t, client, tableName)

	insertQuery := fmt.Sprintf("INSERT INTO %s (id, name, created_at) VALUES (?, ?, ?)", tableName)
	err := client.Execute(insertQuery, 1, "single_user", time.Now())
	require.NoError(t, err)

	selectQuery := fmt.Sprintf("SELECT id, name FROM %s WHERE id = 1", tableName)
	var id uint32
	var name string
	err = client.QueryRow(selectQuery, &id, &name)
	require.NoError(t, err)

	assert.Equal(t, uint32(1), id)
	assert.Equal(t, "single_user", name)
}

func TestClickHouseClient_PreparedStatements(t *testing.T) {
	t.Parallel()

	t.Skip("ClickHouse prepared statements have different behavior compared to standard SQL databases")

}

func TestClickHouseClient_Transaction(t *testing.T) {
	t.Parallel()

	t.Skip("ClickHouse has limited transaction support with Memory engine")

	client := getTestClient(t)
	defer client.Close()

	tableName := "test_transaction"
	createTestTable(t, client, tableName)

	t.Run("Successful Transaction", func(t *testing.T) {
		t.Parallel()

		err := client.BeginTransaction()
		require.NoError(t, err)

		insertQuery := fmt.Sprintf("INSERT INTO %s (id, name, created_at) VALUES (?, ?, ?)", tableName)
		err = client.ExecuteTransaction(insertQuery, 1, "tx_user", time.Now())
		assert.NoError(t, err)

		err = client.EndTransaction(nil)
		assert.NoError(t, err)

		selectQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE id = 1", tableName)
		var count uint64
		err = client.QueryRow(selectQuery, &count)
		require.NoError(t, err)
		assert.Equal(t, uint64(1), count)
	})
}

func TestClickHouseClient_TransactionQuery(t *testing.T) {
	t.Parallel()

	t.Skip("ClickHouse has limited transaction support with Memory engine")

	client := getTestClient(t)
	defer client.Close()

	tableName := "test_transaction_query"
	createTestTable(t, client, tableName)

	err := client.BeginTransaction()
	require.NoError(t, err)

	insertQuery := fmt.Sprintf("INSERT INTO %s (id, name, created_at) VALUES (?, ?, ?)", tableName)
	err = client.ExecuteTransaction(insertQuery, 1, "tx_query_user", time.Now())
	require.NoError(t, err)

	selectQuery := fmt.Sprintf("SELECT id, name FROM %s WHERE id = 1", tableName)
	rows, err := client.QueryTransaction(selectQuery)
	require.NoError(t, err)
	defer rows.Close()

	assert.True(t, rows.Next())
	var id uint32
	var name string
	err = rows.Scan(&id, &name)
	require.NoError(t, err)

	assert.Equal(t, uint32(1), id)
	assert.Equal(t, "tx_query_user", name)

	err = client.QueryRowTransaction(selectQuery, &id, &name)
	require.NoError(t, err)

	assert.Equal(t, uint32(1), id)
	assert.Equal(t, "tx_query_user", name)

	err = client.EndTransaction(nil)
	assert.NoError(t, err)
}

func TestClickHouseClient_TransactionPreparedStatements(t *testing.T) {
	t.Parallel()

	t.Skip("ClickHouse has limited transaction support with Memory engine")

	client := getTestClient(t)
	defer client.Close()

	tableName := "test_transaction_prepared"
	createTestTable(t, client, tableName)

	err := client.BeginTransaction()
	require.NoError(t, err)

	insertQuery := fmt.Sprintf("INSERT INTO %s (id, name, created_at) VALUES (?, ?, ?)", tableName)
	err = client.SetPrepareTransaction(insertQuery)
	require.NoError(t, err)

	err = client.ExecutePrepareTransaction(1, "tx_prepared_user", time.Now())
	assert.NoError(t, err)

	err = client.EndTransaction(nil)
	assert.NoError(t, err)
}

func TestClickHouseClient_ClickHouseSpecific(t *testing.T) {
	t.Parallel()

	client := getTestClient(t)
	defer client.Close()

	t.Run("Basic Operations", func(t *testing.T) {
		t.Parallel()

		client := getTestClient(t)
		defer client.Close()
		rows, err := client.Query("SELECT 1 as test_value")
		require.NoError(t, err)
		defer rows.Close()

		assert.True(t, rows.Next())
		var value int
		err = rows.Scan(&value)
		require.NoError(t, err)
		assert.Equal(t, 1, value)
	})

	t.Run("Database Info", func(t *testing.T) {
		t.Parallel()

		client := getTestClient(t)
		defer client.Close()
		var version string
		err := client.QueryRow("SELECT version()", &version)
		require.NoError(t, err)
		assert.NotEmpty(t, version)
		t.Logf("ClickHouse version: %s", version)
	})

	t.Run("Table Operations", func(t *testing.T) {
		t.Parallel()

		client := getTestClient(t)
		defer client.Close()
		tableName := "test_clickhouse_specific"

		createQuery := fmt.Sprintf(`
			CREATE TABLE %s (
				id UInt32,
				timestamp DateTime,
				value Float64,
				status String
			) ENGINE = Memory
		`, tableName)
		err := client.Execute(createQuery)
		require.NoError(t, err)

		insertQuery := fmt.Sprintf("INSERT INTO %s (id, timestamp, value, status) VALUES (?, ?, ?, ?)", tableName)
		for i := 1; i <= 100; i++ {
			err = client.Execute(insertQuery, i, time.Now(), float64(i)*1.5, fmt.Sprintf("status_%d", i%5))
			require.NoError(t, err)
		}

		var count uint64
		var avgValue float64
		err = client.QueryRow(fmt.Sprintf("SELECT COUNT(*), AVG(value) FROM %s", tableName), &count, &avgValue)
		require.NoError(t, err)

		assert.Equal(t, uint64(100), count)
		assert.Greater(t, avgValue, float64(0))

		rows, err := client.Query(fmt.Sprintf("SELECT status, COUNT(*) as cnt FROM %s GROUP BY status ORDER BY status", tableName))
		require.NoError(t, err)
		defer rows.Close()

		statusCount := 0
		for rows.Next() {
			var status string
			var cnt uint64
			err = rows.Scan(&status, &cnt)
			require.NoError(t, err)
			assert.Equal(t, uint64(20), cnt)
			statusCount++
		}
		assert.Equal(t, 5, statusCount)

		err = client.Execute(fmt.Sprintf("DROP TABLE %s", tableName))
		require.NoError(t, err)
	})

	t.Run("Data Types", func(t *testing.T) {
		t.Parallel()

		client := getTestClient(t)
		defer client.Close()
		tableName := "test_data_types"

		createQuery := fmt.Sprintf(`
			CREATE TABLE %s (
				int_val UInt32,
				string_val String,
				float_val Float64,
				date_val Date,
				datetime_val DateTime,
				bool_val UInt8
			) ENGINE = Memory
		`, tableName)
		err := client.Execute(createQuery)
		require.NoError(t, err)

		insertQuery := fmt.Sprintf("INSERT INTO %s VALUES (?, ?, ?, ?, ?, ?)", tableName)
		testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		testDateTime := time.Now()

		err = client.Execute(insertQuery, 42, "test_string", 3.14159, testDate, testDateTime, 1)
		require.NoError(t, err)

		selectQuery := fmt.Sprintf("SELECT int_val, string_val, float_val, bool_val FROM %s", tableName)
		var intVal uint32
		var stringVal string
		var floatVal float64
		var boolVal uint8

		err = client.QueryRow(selectQuery, &intVal, &stringVal, &floatVal, &boolVal)
		require.NoError(t, err)

		assert.Equal(t, uint32(42), intVal)
		assert.Equal(t, "test_string", stringVal)
		assert.InDelta(t, 3.14159, floatVal, 0.00001)
		assert.Equal(t, uint8(1), boolVal)

		err = client.Execute(fmt.Sprintf("DROP TABLE %s", tableName))
		require.NoError(t, err)
	})

	t.Run("Concurrent Operations", func(t *testing.T) {
		t.Parallel()

		client := getTestClient(t)
		defer client.Close()
		tableName := "test_concurrent"
		createTestTable(t, client, tableName)

		const numGoroutines = 10
		const insertsPerGoroutine = 10

		var wg sync.WaitGroup
		insertQuery := fmt.Sprintf("INSERT INTO %s (id, name, created_at) VALUES (?, ?, ?)", tableName)

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(startID int) {
				defer wg.Done()

				goroutineClient := getTestClient(t)
				defer goroutineClient.Close()

				for j := 0; j < insertsPerGoroutine; j++ {
					id := startID*insertsPerGoroutine + j + 1
					err := goroutineClient.Execute(insertQuery, id, fmt.Sprintf("user_%d", id), time.Now())
					if err != nil {
						t.Errorf("Failed to insert data: %v", err)
					}
				}
			}(i)
		}

		wg.Wait()

		var count uint64
		err := client.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName), &count)
		require.NoError(t, err)
		assert.Equal(t, uint64(numGoroutines*insertsPerGoroutine), count)
	})
}

func TestClickHouseClient_ErrorHandling(t *testing.T) {
	t.Parallel()

	client := &sql.Client{}

	t.Run("Operations without connection", func(t *testing.T) {
		t.Parallel()

		client := sql.Client{}

		_, err := client.Query("SELECT 1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "please call Open first")

		err = client.QueryRow("SELECT 1", new(int))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "please call Open first")

		err = client.Execute("INSERT INTO test VALUES(1)")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "please call Open first")

		err = client.SetPrepare("SELECT 1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "please call Open first")

		err = client.BeginTransaction()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "please call Open first")
	})

	t.Run("Invalid DSN", func(t *testing.T) {
		t.Parallel()

		err := client.Open(sql.DriverClickHouse, "invalid://dsn", 1)
		assert.Error(t, err)
	})

	t.Run("Prepared statement operations without prepare", func(t *testing.T) {
		t.Parallel()

		client := getTestClient(t)
		defer client.Close()

		_, err := client.QueryPrepare()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "please call SetPrepare first")

		_, err = client.QueryRowPrepare()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "please call SetPrepare first")

		err = client.ExecutePrepare()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "please call SetPrepare first")
	})

	t.Run("Transaction operations without begin", func(t *testing.T) {
		t.Parallel()

		client := getTestClient(t)
		defer client.Close()

		_, err := client.QueryTransaction("SELECT 1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "please call BeginTransaction first")

		err = client.QueryRowTransaction("SELECT 1", new(int))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "please call BeginTransaction first")

		err = client.ExecuteTransaction("INSERT INTO test VALUES(1)")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "please call BeginTransaction first")

		err = client.SetPrepareTransaction("SELECT 1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "please call BeginTransaction first")

		err = client.EndTransaction(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "please call BeginTransaction first")
	})
}
