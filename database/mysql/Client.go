// Package mysql provides MySQL client implementations.
package mysql

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// Client is a struct that provides client related methods.
type Client struct {
	tx     *sql.Tx
	txStmt *sql.Stmt

	stmt *sql.Stmt

	connection *sql.DB
}

// Initialize is initialize.
//
// ex)
//
//	err := client.Initialize(`id:password@tcp(address)/table`, 1)
//
//	defer client.Finalize()
func (this *Client) Initialize(dsn string, maxOpenConnection int) error {
	this.Finalize()

	if connection, err := sql.Open("mysql", dsn); err != nil {
		return err
	} else {
		this.connection = connection
	}

	this.connection.SetMaxOpenConns(maxOpenConnection)

	return this.connection.Ping()
}

// Finalize is finalize.
//
// ex)
//
//	err := client.Initialize(`id:password@tcp(address)/table`, 1)
//
//	defer client.Finalize()
func (this *Client) Finalize() error {
	if this.connection == nil {
		return nil
	}

	err := this.connection.Close()
	this.connection = nil

	return err
}

// Query is executes a query and returns the result rows.
//
// ex 1) rows, err := client.Query(`SELECT field ...;`)
//
// ex 2) rows, err := client.Query(`SELECT field ... WHERE field=? ...;`, "value")
//
// defer rows.Close()
//
//	for rows.Next() {
//	    field := 0
//	    err := rows.Scan(&field)
//	}
func (this *Client) Query(query string, args ...any) (*sql.Rows, error) {
	if this.connection == nil {
		return nil, errors.New(fmt.Sprintf("please call Initialize first"))
	}

	return this.connection.Query(query, args...)
}

// QueryRow is select row
//
// ex) err := client.QueryRow(`SELECT field ...;`, &field)
func (this *Client) QueryRow(query string, result ...any) error {
	if this.connection == nil {
		return errors.New(fmt.Sprintf("please call Initialize first"))
	}

	return this.connection.QueryRow(query).Scan(result...)
}

// Execute is executes a query.
//
// ex 1) err := client.Execute(`...`)
//
// ex 2) err := client.Execute(`... WHERE field=? ...;`, "value")
func (this *Client) Execute(query string, args ...any) error {
	if this.connection == nil {
		return errors.New(fmt.Sprintf("please call Initialize first"))
	}

	if result, err := this.connection.Exec(query, args...); err != nil {
		return err
	} else if _, err = result.RowsAffected(); err != nil {
		return err
	}

	return nil
}

// SetPrepare is set prepared statement.
//
// ex) err := client.SetPrepare(`SELECT field ... WHERE field=? ...;`)
func (this *Client) SetPrepare(query string) error {
	if this.connection == nil {
		return errors.New(fmt.Sprintf("please call Initialize first"))
	}

	if stmt, err := this.connection.Prepare(query); err != nil {
		return err
	} else {
		this.stmt = stmt
	}

	return nil
}

// QueryPrepare is query a prepared statement.
//
// ex)
//
//	err := client.SetPrepare(`SELECT field ... WHERE field=? ...;`)
//
//	rows, err := client.QueryPrepare("value")
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

// QueryRowPrepare is query a prepared statement about row.
//
// ex)
//
//	err := client.SetPrepare(`SELECT field ... WHERE field=? ...;`)
//
//	row, err := client.QueryRowPrepare("value")
//
//	field := 0
//
//	err := row.Scan(&field)
func (this *Client) QueryRowPrepare(args ...any) (*sql.Row, error) {
	if this.stmt == nil {
		return nil, errors.New(fmt.Sprintf("please call SetPrepare first"))
	}

	row := this.stmt.QueryRow(args...)

	return row, row.Err()
}

// ExecutePrepare is executes a prepared statement.
//
// ex)
//
//	err := client.SetPrepare(`INSERT INTO ` + table + ` VALUE(field=?);`)
//
//	err = client.ExecutePrepare(2)
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
		return errors.New(fmt.Sprintf("please call Initialize first"))
	}

	if tx, err := this.connection.Begin(); err != nil {
		return err
	} else {
		this.tx = tx
	}

	return nil
}

// EndTransaction is ends a transaction.
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

	if err != nil {
		return this.tx.Rollback()
	}

	return this.tx.Commit()
}

// QueryTransaction is executes a query and returns the result rows.
//
// ex 1) rows, err := client.QueryTransaction(`SELECT field ...;`)
//
// ex 2)
//
//	rows, err := client.QueryTransaction(`SELECT field ... WHERE field=? ...;`, "value")
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

// QueryRowTransaction is select row
//
// ex) err := client.QueryRowTransaction(`SELECT field ...;`, &field)
func (this *Client) QueryRowTransaction(query string, result ...any) error {
	if this.tx == nil {
		return errors.New(fmt.Sprintf("please call BeginTransaction first"))
	}

	return this.tx.QueryRow(query).Scan(result...)
}

// ExecuteTransaction is executes a query.
//
// ex 1) err := client.ExecuteTransaction(`...`)
//
// ex 2) err := client.ExecuteTransaction(`... WHERE field=? ...;`, "value")
func (this *Client) ExecuteTransaction(query string, args ...any) error {
	if this.tx == nil {
		return errors.New(fmt.Sprintf("please call BeginTransaction first"))
	}

	if result, err := this.tx.Exec(query, args...); err != nil {
		return err
	} else if _, err = result.RowsAffected(); err != nil {
		return err
	}

	return nil
}

// SetPrepareTransaction is set prepared statement.
//
// ex) err := client.SetPrepareTransaction(`SELECT field ... WHERE field=? ...;`)
func (this *Client) SetPrepareTransaction(query string) error {
	if this.tx == nil {
		return errors.New(fmt.Sprintf("please call BeginTransaction first"))
	}

	if txStmt, err := this.tx.Prepare(query); err != nil {
		return err
	} else {
		this.txStmt = txStmt
	}

	return nil
}

// QueryPrepareTransaction is query a prepared statement.
//
// ex)
//
//	err := client.SetPrepareTransaction(`SELECT field ... WHERE field=? ...;`)
//
//	rows, err := client.QueryPrepareTransaction("value")
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

// QueryPrepareTransaction is query a prepared statement about row.
//
// ex)
//
//	err := client.SetPrepareTransaction(`SELECT field ... WHERE field=? ...;`)
//
//	row, err := client.QueryRowPrepareTransaction("value")
//
//	field := 0
//
//	err := row.Scan(&field)
func (this *Client) QueryRowPrepareTransaction(args ...any) (*sql.Row, error) {
	if this.txStmt == nil {
		return nil, errors.New(fmt.Sprintf("please call SetPrepareTransaction first"))
	}

	row := this.txStmt.QueryRow(args...)

	return row, row.Err()
}

// ExecutePrepareTransaction is executes a prepared statement.
//
// ex)
//
//	err := client.SetPrepareTransaction(`INSERT INTO ` + table + ` VALUE(field=?);`)
//
//	err = client.ExecutePrepareTransaction(2)
func (this *Client) ExecutePrepareTransaction(args ...any) error {
	if this.txStmt == nil {
		return errors.New(fmt.Sprintf("please call SetPrepareTransaction first"))
	}

	_, err := this.txStmt.Exec(args...)
	return err
}
