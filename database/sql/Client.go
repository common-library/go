// Package sql provides a database client implementation that uses SQL.
package sql

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/btnguyen2k/godynamo"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/microsoft/go-mssqldb"
	_ "github.com/sijms/go-ora"
)

type Driver string

const (
	DriverAmazonDynamoDB     = Driver("godynamo")
	DriverMicrosoftSQLServer = Driver("sqlserver")
	DriverMySQL              = Driver("mysql")
	DriverOracle             = Driver("oracle")
	DriverPostgreSQL         = Driver("postgres")
	DriverSQLite             = Driver("sqlite3")
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
func (this *Client) Open(driver Driver, dsn string, maxOpenConnection int) error {
	this.driver = driver

	if err := this.Close(); err != nil {
		return err
	} else if connection, err := sql.Open(string(driver), dsn); err != nil {
		return err
	} else {
		this.connection = connection
	}

	this.connection.SetMaxOpenConns(maxOpenConnection)

	return this.connection.Ping()
}

// Close closes the database.
//
// ex) err := client.Close()
func (this *Client) Close() error {
	if this.connection == nil {
		return nil
	}

	err := this.connection.Close()
	this.connection = nil

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
func (this *Client) Query(query string, args ...any) (*sql.Rows, error) {
	if this.connection == nil {
		return nil, errors.New(fmt.Sprintf("please call Open first"))
	}

	return this.connection.Query(query, args...)
}

// QueryRow executes a query and copies the values of matched rows.
//
// ex) err := client.QueryRow(`SELECT field FROM ...;`, &field)
func (this *Client) QueryRow(query string, result ...any) error {
	if this.connection == nil {
		return errors.New(fmt.Sprintf("please call Open first"))
	}

	return this.connection.QueryRow(query).Scan(result...)
}

// Execute executes the query.
//
// ex 1) err := client.Execute(`INSERT INTO ... VALUES(value);`)
// ex 2) err := client.Execute(`INSERT INTO ... VALUES(?);`, value)
func (this *Client) Execute(query string, args ...any) error {
	if this.connection == nil {
		return errors.New(fmt.Sprintf("please call Open first"))
	}

	if result, err := this.connection.Exec(query, args...); err != nil {
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
func (this *Client) SetPrepare(query string) error {
	if this.connection == nil {
		return errors.New(fmt.Sprintf("please call Open first"))
	}

	if stmt, err := this.connection.Prepare(query); err != nil {
		return err
	} else {
		this.stmt = stmt
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
func (this *Client) QueryPrepare(args ...any) (*sql.Rows, error) {
	if this.stmt == nil {
		return nil, errors.New(fmt.Sprintf("please call SetPrepare first"))
	}

	return this.stmt.Query(args...)
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
func (this *Client) QueryRowPrepare(args ...any) (*sql.Row, error) {
	if this.stmt == nil {
		return nil, errors.New(fmt.Sprintf("please call SetPrepare first"))
	}

	row := this.stmt.QueryRow(args...)

	return row, row.Err()
}

// ExecutePrepare  executes a prepared statement.
//
// ex)
//
//	err := client.SetPrepare(`INSERT INTO ` + table + `(field) VALUE(?);`)
//
//	err = client.ExecutePrepare(value)
func (this *Client) ExecutePrepare(args ...any) error {
	if this.stmt == nil {
		return errors.New(fmt.Sprintf("please call SetPrepare first"))
	}

	_, err := this.stmt.Exec(args...)
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
func (this *Client) BeginTransaction() error {
	if this.connection == nil {
		return errors.New(fmt.Sprintf("please call Open first"))
	}

	if tx, err := this.connection.Begin(); err != nil {
		return err
	} else {
		this.tx = tx
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
func (this *Client) EndTransaction(err error) error {
	if this.tx == nil {
		return errors.New(fmt.Sprintf("please call BeginTransaction first"))
	}

	if err == nil {
		return this.tx.Commit()
	} else {
		return this.tx.Rollback()
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
func (this *Client) QueryTransaction(query string, args ...any) (*sql.Rows, error) {
	if this.tx == nil {
		return nil, errors.New(fmt.Sprintf("please call BeginTransaction first"))
	}

	return this.tx.Query(query, args...)
}

// QueryRowTransaction executes a query and copies the values of matched rows.
//
// ex) err := client.QueryRowTransaction(`SELECT field ...;`, &field)
func (this *Client) QueryRowTransaction(query string, result ...any) error {
	if this.tx == nil {
		return errors.New(fmt.Sprintf("please call BeginTransaction first"))
	}

	return this.tx.QueryRow(query).Scan(result...)
}

// ExecuteTransaction executes a query.
//
// ex 1) err := client.ExecuteTransaction(`...`)
//
// ex 2) err := client.ExecuteTransaction(`... WHERE field=?;`, value)
func (this *Client) ExecuteTransaction(query string, args ...any) error {
	if this.tx == nil {
		return errors.New(fmt.Sprintf("please call BeginTransaction first"))
	}

	if result, err := this.tx.Exec(query, args...); err != nil {
		return err
	} else {
		_, err = result.RowsAffected()
		return err
	}
}

// SetPrepareTransaction is set prepared statement.
//
// ex) err := client.SetPrepareTransaction(`SELECT field ... WHERE field=?;`)
func (this *Client) SetPrepareTransaction(query string) error {
	if this.tx == nil {
		return errors.New(fmt.Sprintf("please call BeginTransaction first"))
	}

	if txStmt, err := this.tx.Prepare(query); err != nil {
		return err
	} else {
		this.txStmt = txStmt
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
func (this *Client) QueryPrepareTransaction(args ...any) (*sql.Rows, error) {
	if this.txStmt == nil {
		return nil, errors.New(fmt.Sprintf("please call SetPrepareTransaction first"))
	}

	return this.txStmt.Query(args...)
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
func (this *Client) QueryRowPrepareTransaction(args ...any) (*sql.Row, error) {
	if this.txStmt == nil {
		return nil, errors.New(fmt.Sprintf("please call SetPrepareTransaction first"))
	}

	row := this.txStmt.QueryRow(args...)

	return row, row.Err()
}

// ExecutePrepareTransaction executes a prepared statement.
//
// ex)
//
//	err := client.SetPrepareTransaction(`INSERT INTO ` + table + ` VALUE(field=?);`)
//
//	err = client.ExecutePrepareTransaction(value)
func (this *Client) ExecutePrepareTransaction(args ...any) error {
	if this.txStmt == nil {
		return errors.New(fmt.Sprintf("please call SetPrepareTransaction first"))
	}

	_, err := this.txStmt.Exec(args...)
	return err
}

func (this *Client) GetDriver() Driver {
	return this.driver
}
