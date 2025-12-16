// Package sql provides a database client implementation that uses SQL.
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

// Open opens the database.
//
// ex) err := client.Open(sql.DriverMySQL, `id:password@tcp(address)/table`, 1)
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

// Close closes the database.
//
// ex) err := client.Close()
func (c *Client) Close() error {
	if c.connection == nil {
		return nil
	}

	err := c.connection.Close()
	c.connection = nil

	return err
}

// Query executes a query and returns the result rows.
//
// ex)
//
//	rows, err := client.Query(`SELECT field ...;`)
//
//	defer rows.Close()
//
//	for rows.Next() {
//	    field := 0
//	    err := rows.Scan(&field)
//	}
func (c *Client) Query(query string, args ...any) (*sql.Rows, error) {
	if c.connection == nil {
		return nil, errors.New("please call Open first")
	}

	return c.connection.Query(query, args...)
}

// QueryRow executes a query and copies the values of matched rows.
//
// ex) err := client.QueryRow(`SELECT field FROM ...;`, &field)
func (c *Client) QueryRow(query string, result ...any) error {
	if c.connection == nil {
		return errors.New("please call Open first")
	}

	return c.connection.QueryRow(query).Scan(result...)
}

// Execute executes the query.
//
// ex 1) err := client.Execute(`INSERT INTO ... VALUES(value);`)
// ex 2) err := client.Execute(`INSERT INTO ... VALUES(?);`, value)
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

// SetPrepare is set prepared statement.
//
// ex)
//
//	err := client.SetPrepare(`SELECT field ... WHERE field=?;`)
//	err := client.ExecutePrepare(value)
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

// QueryPrepare query a prepared statement.
//
// ex)
//
//	err := client.SetPrepare(`SELECT field ... WHERE field=?;`)
//
//	rows, err := client.QueryPrepare(value)
//
//	defer rows.Close()
//
//	for rows.Next() {
//	    field := 0
//	    err := rows.Scan(&field)
//	}
func (c *Client) QueryPrepare(args ...any) (*sql.Rows, error) {
	if c.stmt == nil {
		return nil, errors.New("please call SetPrepare first")
	}

	return c.stmt.Query(args...)
}

// QueryRowPrepare query a prepared statement about row.
//
// ex)
//
//	err := client.SetPrepare(`SELECT field ... WHERE field=?;`)
//
//	row, err := client.QueryRowPrepare(value)
//
//	field := 0
//	err := row.Scan(&field)
func (c *Client) QueryRowPrepare(args ...any) (*sql.Row, error) {
	if c.stmt == nil {
		return nil, errors.New("please call SetPrepare first")
	}

	row := c.stmt.QueryRow(args...)

	return row, row.Err()
}

// ExecutePrepare  executes a prepared statement.
//
// ex)
//
//	err := client.SetPrepare(`INSERT INTO ` + table + `(field) VALUE(?);`)
//
//	err = client.ExecutePrepare(value)
func (c *Client) ExecutePrepare(args ...any) error {
	if c.stmt == nil {
		return errors.New("please call SetPrepare first")
	}

	_, err := c.stmt.Exec(args...)
	return err
}

// BeginTransaction begins a transaction.
//
// finally, you must call EndTransaction.
//
// ex)
//
//	err := client.BeginTransaction()
//
//	err = client.ExecuteTransaction(`...`)
//
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

// EndTransaction ends a transaction.
//
// if the argument is nil, commit is performed; otherwise, rollback is performed.
//
// ex)
//
//	err := client.BeginTransaction()
//
//	err = client.ExecuteTransaction(`...`)
//
//	err = client.EndTransaction(err)
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

// QueryTransaction executes a query and returns the result rows.
//
// ex 1) rows, err := client.QueryTransaction(`SELECT field ... WHERE field=value;`)
// ex 2) rows, err := client.QueryTransaction(`SELECT field ... WHERE field=?;`, "value")
//
//	defer rows.Close()
//
//	for rows.Next() {
//	    field := 0
//	    err := rows.Scan(&field)
//	}
func (c *Client) QueryTransaction(query string, args ...any) (*sql.Rows, error) {
	if c.tx == nil {
		return nil, errors.New("please call BeginTransaction first")
	}

	return c.tx.Query(query, args...)
}

// QueryRowTransaction executes a query and copies the values of matched rows.
//
// ex) err := client.QueryRowTransaction(`SELECT field ...;`, &field)
func (c *Client) QueryRowTransaction(query string, result ...any) error {
	if c.tx == nil {
		return errors.New("please call BeginTransaction first")
	}

	return c.tx.QueryRow(query).Scan(result...)
}

// ExecuteTransaction executes a query.
//
// ex 1) err := client.ExecuteTransaction(`...`)
//
// ex 2) err := client.ExecuteTransaction(`... WHERE field=?;`, value)
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

// SetPrepareTransaction is set prepared statement.
//
// ex) err := client.SetPrepareTransaction(`SELECT field ... WHERE field=?;`)
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

// QueryPrepareTransaction is query a prepared statement.
//
// ex)
//
//	err := client.SetPrepareTransaction(`SELECT field ... WHERE field=?;`)
//
//	rows, err := client.QueryPrepareTransaction(value)
//
//	defer rows.Close()
//
//	for rows.Next() {
//	    field := 0
//	    err := rows.Scan(&field)
//	}
func (c *Client) QueryPrepareTransaction(args ...any) (*sql.Rows, error) {
	if c.txStmt == nil {
		return nil, errors.New("please call SetPrepareTransaction first")
	}

	return c.txStmt.Query(args...)
}

// QueryPrepareTransaction query a prepared statement about row.
//
// ex)
//
//	err := client.SetPrepareTransaction(`SELECT field ... WHERE field=?;`)
//
//	row, err := client.QueryRowPrepareTransaction(value)
//
//	field := 0
//	err := row.Scan(&field)
func (c *Client) QueryRowPrepareTransaction(args ...any) (*sql.Row, error) {
	if c.txStmt == nil {
		return nil, errors.New("please call SetPrepareTransaction first")
	}

	row := c.txStmt.QueryRow(args...)

	return row, row.Err()
}

// ExecutePrepareTransaction executes a prepared statement.
//
// ex)
//
//	err := client.SetPrepareTransaction(`INSERT INTO ` + table + ` VALUE(field=?);`)
//
//	err = client.ExecutePrepareTransaction(value)
func (c *Client) ExecutePrepareTransaction(args ...any) error {
	if c.txStmt == nil {
		return errors.New("please call SetPrepareTransaction first")
	}

	_, err := c.txStmt.Exec(args...)
	return err
}

func (c *Client) GetDriver() Driver {
	return c.driver
}
