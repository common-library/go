// Package mysql provides MySQL interface.
//
// used "github.com/go-sql-driver/mysql".
package mysql

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// MySQL is object that provides MySQL interface.
type MySQL struct {
	tx     *sql.Tx
	txStmt *sql.Stmt

	stmt *sql.Stmt

	connection *sql.DB
}

// Initialize is initialize.
//
// ex)
//
//	err := mysql.Initialize(`id:password@tcp(address)/table`, 1)
//
//	defer mysql.Finalize()
func (this *MySQL) Initialize(dsn string, maxOpenConnection int) error {
	this.Finalize()

	var err error
	this.connection, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	this.connection.SetMaxOpenConns(maxOpenConnection)

	return this.connection.Ping()
}

// Finalize is finalize.
//
// ex)
//
//	err := mysql.Initialize(`id:password@tcp(address)/table`, 1)
//
//	defer mysql.Finalize()
func (this *MySQL) Finalize() error {
	if this.connection != nil {
		err := this.connection.Close()
		if err != nil {
			return err
		}
		this.connection = nil
	}

	return nil
}

// Query is executes a query and returns the result rows.
//
// ex 1) rows, err := mysql.Query(`SELECT field ...;`)
//
// ex 2) rows, err := mysql.Query(`SELECT field ... WHERE field=? ...;`, "value")
//
// defer rows.Close()
//
//	for rows.Next() {
//	    field := 0
//	    err := rows.Scan(&field)
//	}
func (this *MySQL) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if this.connection == nil {
		return nil, errors.New(fmt.Sprintf("please call Initialize first"))
	}

	return this.connection.Query(query, args...)
}

// QueryRow is select row
//
// ex) err := mysql.QueryRow(`SELECT field ...;`, &field)
func (this *MySQL) QueryRow(query string, result ...interface{}) error {
	if this.connection == nil {
		return errors.New(fmt.Sprintf("please call Initialize first"))
	}

	return this.connection.QueryRow(query).Scan(result...)
}

// Execute is executes a query.
//
// ex 1) err := mysql.Execute(`...`)
//
// ex 2) err := mysql.Execute(`... WHERE field=? ...;`, "value")
func (this *MySQL) Execute(query string, args ...interface{}) error {
	if this.connection == nil {
		return errors.New(fmt.Sprintf("please call Initialize first"))
	}

	result, err := this.connection.Exec(query, args...)
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

// SetPrepare is set prepared statement.
//
// ex) err := mysql.SetPrepare(`SELECT field ... WHERE field=? ...;`)
func (this *MySQL) SetPrepare(query string) error {
	if this.connection == nil {
		return errors.New(fmt.Sprintf("please call Initialize first"))
	}

	var err error
	this.stmt, err = this.connection.Prepare(query)
	return err
}

// QueryPrepare is query a prepared statement.
//
// ex)
//
//	err := mysql.SetPrepare(`SELECT field ... WHERE field=? ...;`)
//
//	rows, err := mysql.QueryPrepare("value")
//
//	defer rows.Close()
//
//	for rows.Next() {
//	    field := 0
//	    err := rows.Scan(&field)
//	}
func (this *MySQL) QueryPrepare(args ...interface{}) (*sql.Rows, error) {
	if this.stmt == nil {
		return nil, errors.New(fmt.Sprintf("please call SetPrepare first"))
	}

	return this.stmt.Query(args...)
}

// QueryRowPrepare is query a prepared statement about row.
//
// ex)
//
//	err := mysql.SetPrepare(`SELECT field ... WHERE field=? ...;`)
//
//	row, err := mysql.QueryRowPrepare("value")
//
//	field := 0
//
//	err := row.Scan(&field)
func (this *MySQL) QueryRowPrepare(args ...interface{}) (*sql.Row, error) {
	if this.stmt == nil {
		return nil, errors.New(fmt.Sprintf("please call SetPrepare first"))
	}

	row := this.stmt.QueryRow(args...)

	if row.Err() != nil {
		return nil, row.Err()
	}

	return row, nil
}

// ExecutePrepare is executes a prepared statement.
//
// ex)
//
//	err := mysql.SetPrepare(`INSERT INTO ` + table + ` VALUE(field=?);`)
//
//	err = mysql.ExecutePrepare(2)
func (this *MySQL) ExecutePrepare(args ...interface{}) error {
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
//	err := mysql.BeginTransaction()
//
//	err = mysql.ExecuteTransaction(`...`)
//
//	err = mysql.EndTransaction(err)
func (this *MySQL) BeginTransaction() error {
	if this.connection == nil {
		return errors.New(fmt.Sprintf("please call Initialize first"))
	}

	var err error
	this.tx, err = this.connection.Begin()

	return err
}

// EndTransaction is ends a transaction.
//
// if the argument is nil, commit is performed; otherwise, rollback is performed.
//
// ex)
//
//	err := mysql.BeginTransaction()
//
//	err = mysql.ExecuteTransaction(`...`)
//
//	err = mysql.EndTransaction(err)
func (this *MySQL) EndTransaction(err error) error {
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
// ex 1) rows, err := mysql.QueryTransaction(`SELECT field ...;`)
//
// ex 2)
//
//	rows, err := mysql.QueryTransaction(`SELECT field ... WHERE field=? ...;`, "value")
//
//	defer rows.Close()
//
//	for rows.Next() {
//	    field := 0
//	    err := rows.Scan(&field)
//	}
func (this *MySQL) QueryTransaction(query string, args ...interface{}) (*sql.Rows, error) {
	if this.tx == nil {
		return nil, errors.New(fmt.Sprintf("please call BeginTransaction first"))
	}

	return this.tx.Query(query, args...)
}

// QueryRowTransaction is select row
//
// ex) err := mysql.QueryRowTransaction(`SELECT field ...;`, &field)
func (this *MySQL) QueryRowTransaction(query string, result ...interface{}) error {
	if this.tx == nil {
		return errors.New(fmt.Sprintf("please call BeginTransaction first"))
	}

	return this.tx.QueryRow(query).Scan(result...)
}

// ExecuteTransaction is executes a query.
//
// ex 1) err := mysql.ExecuteTransaction(`...`)
//
// ex 2) err := mysql.ExecuteTransaction(`... WHERE field=? ...;`, "value")
func (this *MySQL) ExecuteTransaction(query string, args ...interface{}) error {
	if this.tx == nil {
		return errors.New(fmt.Sprintf("please call BeginTransaction first"))
	}

	result, err := this.tx.Exec(query, args...)
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

// SetPrepareTransaction is set prepared statement.
//
// ex) err := mysql.SetPrepareTransaction(`SELECT field ... WHERE field=? ...;`)
func (this *MySQL) SetPrepareTransaction(query string) error {
	if this.tx == nil {
		return errors.New(fmt.Sprintf("please call BeginTransaction first"))
	}

	var err error
	this.txStmt, err = this.tx.Prepare(query)
	return err
}

// QueryPrepareTransaction is query a prepared statement.
//
// ex)
//
//	err := mysql.SetPrepareTransaction(`SELECT field ... WHERE field=? ...;`)
//
//	rows, err := mysql.QueryPrepareTransaction("value")
//
//	defer rows.Close()
//
//	for rows.Next() {
//	    field := 0
//	    err := rows.Scan(&field)
//	}
func (this *MySQL) QueryPrepareTransaction(args ...interface{}) (*sql.Rows, error) {
	if this.txStmt == nil {
		return nil, errors.New(fmt.Sprintf("please call SetPrepareTransaction first"))
	}

	return this.txStmt.Query(args...)
}

// QueryPrepareTransaction is query a prepared statement about row.
//
// ex)
//
//	err := mysql.SetPrepareTransaction(`SELECT field ... WHERE field=? ...;`)
//
//	row, err := mysql.QueryRowPrepareTransaction("value")
//
//	field := 0
//
//	err := row.Scan(&field)
func (this *MySQL) QueryRowPrepareTransaction(args ...interface{}) (*sql.Row, error) {
	if this.txStmt == nil {
		return nil, errors.New(fmt.Sprintf("please call SetPrepareTransaction first"))
	}

	row := this.txStmt.QueryRow(args...)

	if row.Err() != nil {
		return nil, row.Err()
	}

	return row, nil
}

// ExecutePrepareTransaction is executes a prepared statement.
//
// ex)
//
//	err := mysql.SetPrepareTransaction(`INSERT INTO ` + table + ` VALUE(field=?);`)
//
//	err = mysql.ExecutePrepareTransaction(2)
func (this *MySQL) ExecutePrepareTransaction(args ...interface{}) error {
	if this.txStmt == nil {
		return errors.New(fmt.Sprintf("please call SetPrepareTransaction first"))
	}

	_, err := this.txStmt.Exec(args...)
	return err
}
