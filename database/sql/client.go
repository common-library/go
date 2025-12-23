// Package sql provides a unified SQL database client supporting multiple database drivers.
//
// This package offers a consistent interface for interacting with various SQL databases
// including MySQL, PostgreSQL, SQLite, ClickHouse, DynamoDB, SQL Server, and Oracle.
//
// Features:
//   - Support for 7 database drivers (MySQL, PostgreSQL, SQLite, ClickHouse, DynamoDB, SQL Server, Oracle)
//   - Transaction management with Begin/End pattern
//   - Prepared statement support for both regular and transactional queries
//   - Connection pooling configuration
//   - Consistent API across all database types
//
// Example:
//
//	var client sql.Client
//	client.Open(sql.DriverMySQL, "user:pass@tcp(localhost)/db", 10)
//	defer client.Close()
//	client.Execute("INSERT INTO users (name) VALUES (?)", "Alice")
package sql

import (
	"database/sql"
	"errors"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	_ "github.com/btnguyen2k/godynamo"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/microsoft/go-mssqldb"
	_ "github.com/sijms/go-ora"
	_ "modernc.org/sqlite"
)

type Driver string

const (
	DriverAmazonDynamoDB     = Driver("godynamo")
	DriverClickHouse         = Driver("clickhouse")
	DriverMicrosoftSQLServer = Driver("sqlserver")
	DriverMySQL              = Driver("mysql")
	DriverOracle             = Driver("oracle")
	DriverPostgreSQL         = Driver("postgres")
	DriverSQLite             = Driver("sqlite")
)

// Client is a struct that provides client related methods.
type Client struct {
	driver Driver

	tx     *sql.Tx
	txStmt *sql.Stmt

	stmt *sql.Stmt

	connection *sql.DB
}

// Open establishes a connection to the specified database.
// This method initializes the database connection pool and validates the connection with a ping.
//
// Parameters:
//   - driver: Database driver (DriverMySQL, DriverPostgreSQL, DriverSQLite, etc.)
//   - dsn: Data Source Name - connection string specific to the driver
//   - maxOpenConnection: Maximum number of open connections in the pool
//
// Returns:
//   - error: Returns an error if the connection cannot be established or ping fails
//
// Example:
//
//	err := client.Open(sql.DriverMySQL, "user:pass@tcp(localhost:3306)/database", 10)
//	if err != nil {
//		log.Fatal(err)
//	}
func (c *Client) Open(driver Driver, dsn string, maxOpenConnection int) error {
	c.driver = driver

	if err := c.Close(); err != nil {
		return err
	} else if connection, err := sql.Open(string(driver), dsn); err != nil {
		return err
	} else {
		c.connection = connection
	}

	c.connection.SetMaxOpenConns(maxOpenConnection)

	return c.connection.Ping()
}

// Close closes the database connection and releases all resources.
// It is safe to call Close multiple times.
//
// Returns:
//   - error: Returns an error if closing the connection fails, nil otherwise
//
// Example:
//
//	defer client.Close()
func (c *Client) Close() error {
	if c.connection == nil {
		return nil
	}

	err := c.connection.Close()
	c.connection = nil

	return err
}

// Query executes a SQL query and returns the result rows.
// The caller is responsible for closing the returned rows.
//
// Parameters:
//   - query: SQL query string (use ? for parameter placeholders)
//   - args: Optional query parameters
//
// Returns:
//   - *sql.Rows: Result set that can be iterated
//   - error: Returns an error if the query fails
//
// Example:
//
//	rows, err := client.Query("SELECT id, name FROM users WHERE age > ?", 18)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer rows.Close()
//
//	for rows.Next() {
//		var id int
//		var name string
//		err := rows.Scan(&id, &name)
//		if err != nil {
//			log.Fatal(err)
//		}
//		fmt.Printf("ID: %d, Name: %s\n", id, name)
//	}
func (c *Client) Query(query string, args ...any) (*sql.Rows, error) {
	if c.connection == nil {
		return nil, errors.New("please call Open first")
	}

	return c.connection.Query(query, args...)
}

// QueryRow executes a query and scans the first row into the provided variables.
// This is a convenience method for queries that return a single row.
//
// Parameters:
//   - query: SQL query string
//   - result: Pointers to variables where column values will be stored
//
// Returns:
//   - error: Returns an error if the query fails or scanning fails
//
// Example:
//
//	var name string
//	var age int
//	err := client.QueryRow("SELECT name, age FROM users WHERE id = ?", &name, &age, 1)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Name: %s, Age: %d\n", name, age)
func (c *Client) QueryRow(query string, result ...any) error {
	if c.connection == nil {
		return errors.New("please call Open first")
	}

	return c.connection.QueryRow(query).Scan(result...)
}

// Execute executes a SQL statement (INSERT, UPDATE, DELETE, etc.) that doesn't return rows.
// This method is used for data modification operations.
//
// Parameters:
//   - query: SQL statement string (use ? for parameter placeholders)
//   - args: Optional statement parameters
//
// Returns:
//   - error: Returns an error if the statement execution fails
//
// Example:
//
//	// Insert
//	err := client.Execute("INSERT INTO users (name, age) VALUES (?, ?)", "Alice", 30)
//
//	// Update
//	err = client.Execute("UPDATE users SET age = ? WHERE name = ?", 31, "Alice")
//
//	// Delete
//	err = client.Execute("DELETE FROM users WHERE age < ?", 18)
func (c *Client) Execute(query string, args ...any) error {
	if c.connection == nil {
		return errors.New("please call Open first")
	}

	if result, err := c.connection.Exec(query, args...); err != nil {
		return err
	} else {
		_, err = result.RowsAffected()
		return err
	}
}

// SetPrepare creates a prepared statement for later execution.
// Prepared statements improve performance when executing the same query multiple times
// and provide protection against SQL injection.
//
// Parameters:
//   - query: SQL query string with ? placeholders for parameters
//
// Returns:
//   - error: Returns an error if statement preparation fails
//
// Example:
//
//	err := client.SetPrepare("SELECT id, name FROM users WHERE age > ?")
//	if err != nil {
//		log.Fatal(err)
//	}
func (c *Client) SetPrepare(query string) error {
	if c.connection == nil {
		return errors.New("please call Open first")
	}

	if stmt, err := c.connection.Prepare(query); err != nil {
		return err
	} else {
		c.stmt = stmt
		return nil
	}
}

// QueryPrepare executes a prepared statement and returns the result rows.
// Must be called after SetPrepare.
//
// Parameters:
//   - args: Parameters to bind to the prepared statement
//
// Returns:
//   - *sql.Rows: Result set that can be iterated
//   - error: Returns an error if the query fails
//
// Example:
//
//	err := client.SetPrepare("SELECT id, name FROM users WHERE age > ?")
//	rows, err := client.QueryPrepare(18)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer rows.Close()
//
//	for rows.Next() {
//		var id int
//		var name string
//		rows.Scan(&id, &name)
//		fmt.Printf("ID: %d, Name: %s\n", id, name)
//	}
func (c *Client) QueryPrepare(args ...any) (*sql.Rows, error) {
	if c.stmt == nil {
		return nil, errors.New("please call SetPrepare first")
	}

	return c.stmt.Query(args...)
}

// QueryRowPrepare executes a prepared statement and returns a single row.
// Must be called after SetPrepare.
//
// Parameters:
//   - args: Parameters to bind to the prepared statement
//
// Returns:
//   - *sql.Row: Single row result
//   - error: Returns an error if the query fails
//
// Example:
//
//	err := client.SetPrepare("SELECT name, age FROM users WHERE id = ?")
//	row, err := client.QueryRowPrepare(1)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	var name string
//	var age int
//	err = row.Scan(&name, &age)
func (c *Client) QueryRowPrepare(args ...any) (*sql.Row, error) {
	if c.stmt == nil {
		return nil, errors.New("please call SetPrepare first")
	}

	row := c.stmt.QueryRow(args...)

	return row, row.Err()
}

// ExecutePrepare executes a prepared statement that doesn't return rows.
// Must be called after SetPrepare. Use for INSERT, UPDATE, DELETE operations.
//
// Parameters:
//   - args: Parameters to bind to the prepared statement
//
// Returns:
//   - error: Returns an error if the execution fails
//
// Example:
//
//	err := client.SetPrepare("INSERT INTO users (name, age) VALUES (?, ?)")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	err = client.ExecutePrepare("Alice", 30)
//	err = client.ExecutePrepare("Bob", 25)
//	err = client.ExecutePrepare("Charlie", 35)
func (c *Client) ExecutePrepare(args ...any) error {
	if c.stmt == nil {
		return errors.New("please call SetPrepare first")
	}

	_, err := c.stmt.Exec(args...)
	return err
}

// BeginTransaction starts a new database transaction.
// All subsequent operations using *Transaction methods will be part of this transaction
// until EndTransaction is called. Transactions ensure atomicity - all operations succeed or all fail.
//
// Returns:
//   - error: Returns an error if beginning the transaction fails
//
// Example:
//
//	err := client.BeginTransaction()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	err = client.ExecuteTransaction("INSERT INTO users (name) VALUES (?)", "Alice")
//	err = client.ExecuteTransaction("INSERT INTO logs (message) VALUES (?)", "User created")
//
//	// Commit if no errors, rollback otherwise
//	err = client.EndTransaction(err)
func (c *Client) BeginTransaction() error {
	if c.connection == nil {
		return errors.New("please call Open first")
	}

	if tx, err := c.connection.Begin(); err != nil {
		return err
	} else {
		c.tx = tx
		return nil
	}
}

// EndTransaction commits or rolls back a transaction based on the error argument.
// If err is nil, the transaction is committed. If err is not nil, the transaction is rolled back.
// This method must be called after BeginTransaction to complete the transaction.
//
// Parameters:
//   - err: If nil, commits the transaction; otherwise rolls back
//
// Returns:
//   - error: Returns an error if commit or rollback fails
//
// Example:
//
//	err := client.BeginTransaction()
//	err = client.ExecuteTransaction("INSERT INTO users (name) VALUES (?)", "Alice")
//	if err != nil {
//		client.EndTransaction(err) // Rollback
//		return err
//	}
//	err = client.ExecuteTransaction("UPDATE counters SET value = value + 1")
//	err = client.EndTransaction(err) // Commit if no error, rollback if error
func (c *Client) EndTransaction(err error) error {
	if c.tx == nil {
		return errors.New("please call BeginTransaction first")
	}

	if err == nil {
		return c.tx.Commit()
	} else {
		return c.tx.Rollback()
	}
}

// QueryTransaction executes a query within the current transaction and returns result rows.
// Must be called after BeginTransaction and before EndTransaction.
//
// Parameters:
//   - query: SQL query string (use ? for parameter placeholders)
//   - args: Optional query parameters
//
// Returns:
//   - *sql.Rows: Result set that can be iterated
//   - error: Returns an error if the query fails
//
// Example:
//
//	err := client.BeginTransaction()
//	rows, err := client.QueryTransaction("SELECT id, balance FROM accounts WHERE user_id = ?", 123)
//	if err != nil {
//		client.EndTransaction(err)
//		return err
//	}
//	defer rows.Close()
//
//	for rows.Next() {
//		var id int
//		var balance float64
//		err := rows.Scan(&id, &balance)
//		// Process row...
//	}
//	err = client.EndTransaction(err)
func (c *Client) QueryTransaction(query string, args ...any) (*sql.Rows, error) {
	if c.tx == nil {
		return nil, errors.New("please call BeginTransaction first")
	}

	return c.tx.Query(query, args...)
}

// QueryRowTransaction executes a query within the current transaction and scans the first row.
// Must be called after BeginTransaction and before EndTransaction.
//
// Parameters:
//   - query: SQL query string
//   - result: Pointers to variables where column values will be stored
//
// Returns:
//   - error: Returns an error if the query fails or scanning fails
//
// Example:
//
//	err := client.BeginTransaction()
//	var balance float64
//	err = client.QueryRowTransaction("SELECT balance FROM accounts WHERE id = ?", &balance, 1)
//	if err != nil {
//		client.EndTransaction(err)
//		return err
//	}
//	err = client.EndTransaction(err)
func (c *Client) QueryRowTransaction(query string, result ...any) error {
	if c.tx == nil {
		return errors.New("please call BeginTransaction first")
	}

	return c.tx.QueryRow(query).Scan(result...)
}

// ExecuteTransaction executes a SQL statement within the current transaction.
// Must be called after BeginTransaction and before EndTransaction.
// Use for INSERT, UPDATE, DELETE operations that should be part of a transaction.
//
// Parameters:
//   - query: SQL statement string (use ? for parameter placeholders)
//   - args: Optional statement parameters
//
// Returns:
//   - error: Returns an error if the execution fails
//
// Example:
//
//	err := client.BeginTransaction()
//	// Debit from one account
//	err = client.ExecuteTransaction("UPDATE accounts SET balance = balance - ? WHERE id = ?", 100.0, 1)
//	if err != nil {
//		client.EndTransaction(err)
//		return err
//	}
//	// Credit to another account
//	err = client.ExecuteTransaction("UPDATE accounts SET balance = balance + ? WHERE id = ?", 100.0, 2)
//	err = client.EndTransaction(err) // Commit both or rollback both
func (c *Client) ExecuteTransaction(query string, args ...any) error {
	if c.tx == nil {
		return errors.New("please call BeginTransaction first")
	}

	if result, err := c.tx.Exec(query, args...); err != nil {
		return err
	} else {
		_, err = result.RowsAffected()
		return err
	}
}

// SetPrepareTransaction creates a prepared statement within the current transaction.
// Must be called after BeginTransaction and before EndTransaction.
//
// Parameters:
//   - query: SQL query string with ? placeholders for parameters
//
// Returns:
//   - error: Returns an error if statement preparation fails
//
// Example:
//
//	err := client.BeginTransaction()
//	err = client.SetPrepareTransaction("INSERT INTO logs (message, timestamp) VALUES (?, ?)")
//	if err != nil {
//		client.EndTransaction(err)
//		return err
//	}
func (c *Client) SetPrepareTransaction(query string) error {
	if c.tx == nil {
		return errors.New("please call BeginTransaction first")
	}

	if txStmt, err := c.tx.Prepare(query); err != nil {
		return err
	} else {
		c.txStmt = txStmt
		return nil
	}
}

// QueryPrepareTransaction executes a prepared statement within a transaction and returns result rows.
// Must be called after SetPrepareTransaction.
//
// Parameters:
//   - args: Parameters to bind to the prepared statement
//
// Returns:
//   - *sql.Rows: Result set that can be iterated
//   - error: Returns an error if the query fails
//
// Example:
//
//	err := client.BeginTransaction()
//	err = client.SetPrepareTransaction("SELECT id, name FROM users WHERE age > ?")
//	rows, err := client.QueryPrepareTransaction(18)
//	if err != nil {
//		client.EndTransaction(err)
//		return err
//	}
//	defer rows.Close()
//
//	for rows.Next() {
//		var id int
//		var name string
//		rows.Scan(&id, &name)
//	}
//	err = client.EndTransaction(err)
func (c *Client) QueryPrepareTransaction(args ...any) (*sql.Rows, error) {
	if c.txStmt == nil {
		return nil, errors.New("please call SetPrepareTransaction first")
	}

	return c.txStmt.Query(args...)
}

// QueryRowPrepareTransaction executes a prepared statement within a transaction and returns a single row.
// Must be called after SetPrepareTransaction.
//
// Parameters:
//   - args: Parameters to bind to the prepared statement
//
// Returns:
//   - *sql.Row: Single row result
//   - error: Returns an error if the query fails
//
// Example:
//
//	err := client.BeginTransaction()
//	err = client.SetPrepareTransaction("SELECT balance FROM accounts WHERE id = ?")
//	row, err := client.QueryRowPrepareTransaction(1)
//	if err != nil {
//		client.EndTransaction(err)
//		return err
//	}
//
//	var balance float64
//	err = row.Scan(&balance)
//	err = client.EndTransaction(err)
func (c *Client) QueryRowPrepareTransaction(args ...any) (*sql.Row, error) {
	if c.txStmt == nil {
		return nil, errors.New("please call SetPrepareTransaction first")
	}

	row := c.txStmt.QueryRow(args...)

	return row, row.Err()
}

// ExecutePrepareTransaction executes a prepared statement within a transaction.
// Must be called after SetPrepareTransaction. Use for batch INSERT, UPDATE, DELETE operations.
//
// Parameters:
//   - args: Parameters to bind to the prepared statement
//
// Returns:
//   - error: Returns an error if the execution fails
//
// Example:
//
//	err := client.BeginTransaction()
//	err = client.SetPrepareTransaction("INSERT INTO logs (level, message) VALUES (?, ?)")
//	if err != nil {
//		client.EndTransaction(err)
//		return err
//	}
//
//	err = client.ExecutePrepareTransaction("INFO", "Server started")
//	err = client.ExecutePrepareTransaction("DEBUG", "Connection established")
//	err = client.ExecutePrepareTransaction("INFO", "Request processed")
//	err = client.EndTransaction(err)
func (c *Client) ExecutePrepareTransaction(args ...any) error {
	if c.txStmt == nil {
		return errors.New("please call SetPrepareTransaction first")
	}

	_, err := c.txStmt.Exec(args...)
	return err
}

// GetDriver returns the database driver currently in use.
// This can be used to implement driver-specific behavior.
//
// Returns:
//   - Driver: The current database driver
//
// Example:
//
//	driver := client.GetDriver()
//	if driver == sql.DriverMySQL {
//		// MySQL-specific logic
//	} else if driver == sql.DriverPostgreSQL {
//		// PostgreSQL-specific logic
//	}
func (c *Client) GetDriver() Driver {
	return c.driver
}
