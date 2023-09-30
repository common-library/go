package mysql_test

import (
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/heaven-chp/common-library-go/database/mysql"
)

var database string = strings.ReplaceAll(uuid.NewString(), "-", "")
var table string = strings.ReplaceAll(uuid.NewString(), "-", "")

func createTable() error {
	for i := 0; i < 10; i++ {
		if string(database[0]) == strconv.Itoa(i) {
			i = -1
			database = strings.ReplaceAll(uuid.NewString(), "-", "")
			continue
		}

		if string(table[0]) == strconv.Itoa(i) {
			i = -1
			table = strings.ReplaceAll(uuid.NewString(), "-", "")
			continue
		}
	}

	client := mysql.Client{}

	err := client.Initialize(`root:root@tcp(127.0.0.1)/`, 1)
	if err != nil {
		return err
	}
	defer client.Finalize()

	err = client.Execute(`CREATE DATABASE ` + database + `;`)
	if err != nil {
		return err
	}

	err = client.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1)
	if err != nil {
		return err
	}

	err = client.Execute(`CREATE TABLE ` + table + `(field int);`)
	if err != nil {
		return err
	}

	err = client.Execute(`INSERT INTO ` + table + `(field) VALUE(1);`)
	if err != nil {
		return err
	}

	return nil
}

func deleteDatabase() error {
	client := mysql.Client{}

	err := client.Initialize(`root:root@tcp(127.0.0.1)/`, 1)
	if err != nil {
		return err
	}
	defer client.Finalize()

	err = client.Execute(`DROP DATABASE ` + database + `;`)
	if err != nil {
		return err
	}

	return nil
}

func TestInitialize(t *testing.T) {
	client := mysql.Client{}

	err := client.Initialize(`root:root@tcp(127.0.0.1)`, 1)
	if err.Error() != `invalid DSN: missing the slash separating the database name` {
		t.Error(err)
	}

	err = client.Initialize(`root:root@tcp(127.0.0.1)/`, 1)
	if err != nil {
		t.Error(err)
	}
	defer client.Finalize()
}

func TestFinalize(t *testing.T) {
	client := mysql.Client{}

	client.Finalize()

	err := client.Initialize(`root:root@tcp(127.0.0.1)/`, 1)
	if err != nil {
		t.Error(err)
	}
	defer client.Finalize()
}

func TestQuery(t *testing.T) {
	err := createTable()
	if err != nil {
		t.Error(err)
	}

	client := mysql.Client{}

	_, err = client.Query(`SELECT field FROM ` + table + `;`)
	if err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	err = client.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1)
	if err != nil {
		t.Error(err)
	}
	defer client.Finalize()

	rows, err := client.Query(`SELECT field FROM ` + table + `;`)
	if err != nil {
		t.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		field := 0
		err := rows.Scan(&field)
		if err != nil {
			t.Error(err)
		}

		if field != 1 {
			t.Errorf("invalid field : (%d)", field)
		}
	}

	rows2, err := client.Query(`SELECT field FROM `+table+` WHERE field=?;`, 1)
	if err != nil {
		t.Error(err)
	}
	defer rows2.Close()
	for rows2.Next() {
		field := 0
		err := rows2.Scan(&field)
		if err != nil {
			t.Error(err)
		}

		if field != 1 {
			t.Errorf("invalid field : (%d)", field)
		}
	}

	err = deleteDatabase()
	if err != nil {
		t.Error(err)
	}
}

func TestQueryRow(t *testing.T) {
	err := createTable()
	if err != nil {
		t.Error(err)
	}

	client := mysql.Client{}

	err = client.QueryRow(``)
	if err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	err = client.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1)
	if err != nil {
		t.Error(err)
	}
	defer client.Finalize()

	field := 0
	err = client.QueryRow(`SELECT field FROM `+table+`;`, &field)
	if err != nil {
		t.Error(err)
	}
	if field != 1 {
		t.Errorf("invalid field : (%d)", field)
	}

	err = deleteDatabase()
	if err != nil {
		t.Error(err)
	}
}

func TestExecute(t *testing.T) {
	err := createTable()
	if err != nil {
		t.Error(err)
	}

	client := mysql.Client{}

	err = client.Execute(`UPDATE ` + table + ` SET field=2`)
	if err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	err = client.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1)
	if err != nil {
		t.Error(err)
	}
	defer client.Finalize()

	err = client.Execute(``)
	if err.Error() != `Error 1065 (42000): Query was empty` {
		t.Error(err)
	}

	err = client.Execute(`UPDATE ` + table + ` SET field=2`)
	if err != nil {
		t.Error(err)
	}

	field := 0
	err = client.QueryRow(`SELECT field FROM `+table+`;`, &field)
	if err != nil {
		t.Error(err)
	}
	if field != 2 {
		t.Errorf("invalid field : (%d)", field)
	}

	err = client.Execute(`INSERT INTO `+table+` VALUE(field=?);`, 1)
	if err != nil {
		t.Error(err)
	}

	count := 0
	err = client.QueryRow(`SELECT COUNT(*) FROM `+table+`;`, &count)
	if err != nil {
		t.Error(err)
	}
	if count != 2 {
		t.Errorf("invalid count : (%d)", count)
	}

	err = deleteDatabase()
	if err != nil {
		t.Error(err)
	}
}

func TestSetPrepare(t *testing.T) {
	err := createTable()
	if err != nil {
		t.Error(err)
	}

	client := mysql.Client{}

	err = client.SetPrepare(`INSERT INTO ` + table + ` VALUE(field=?);`)
	if err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	err = client.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1)
	if err != nil {
		t.Error(err)
	}
	defer client.Finalize()

	err = client.SetPrepare(`INSERT INTO ` + table + ` VALUE(field=?);`)
	if err != nil {
		t.Error(err)
	}

	err = client.ExecutePrepare(2)
	if err != nil {
		t.Error(err)
	}

	count := 0
	err = client.QueryRow(`SELECT COUNT(*) FROM `+table+`;`, &count)
	if err != nil {
		t.Error(err)
	}
	if count != 2 {
		t.Errorf("invalid count : (%d)", count)
	}

	err = deleteDatabase()
	if err != nil {
		t.Error(err)
	}
}

func TestQueryPrepare(t *testing.T) {
	err := createTable()
	if err != nil {
		t.Error(err)
	}

	client := mysql.Client{}

	_, err = client.QueryPrepare(1)
	if err.Error() != `please call SetPrepare first` {
		t.Error(err)
	}

	err = client.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1)
	if err != nil {
		t.Error(err)
	}
	defer client.Finalize()

	err = client.SetPrepare(`SELECT field FROM ` + table + ` WHERE field=?;`)
	if err != nil {
		t.Error(err)
	}

	rows, err := client.QueryPrepare(1)
	if err != nil {
		t.Error(err)
	}
	defer rows.Close()

	for rows.Next() {
		field := 0
		err := rows.Scan(&field)
		if err != nil {
			t.Error(err)
		}

		if field != 1 {
			t.Errorf("invalid field : (%d)", field)
		}
	}

	err = deleteDatabase()
	if err != nil {
		t.Error(err)
	}
}

func TestQueryRowPrepare(t *testing.T) {
	err := createTable()
	if err != nil {
		t.Error(err)
	}

	client := mysql.Client{}

	_, err = client.QueryRowPrepare(1)
	if err.Error() != `please call SetPrepare first` {
		t.Error(err)
	}

	err = client.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1)
	if err != nil {
		t.Error(err)
	}
	defer client.Finalize()

	err = client.SetPrepare(`SELECT field FROM ` + table + ` WHERE field=?;`)
	if err != nil {
		t.Error(err)
	}

	row, err := client.QueryRowPrepare()
	if err.Error() != `sql: expected 1 arguments, got 0` {
		t.Error(err)
	}

	row, err = client.QueryRowPrepare(1)
	if err != nil {
		t.Error(err)
	}

	field := 0
	err = row.Scan(&field)
	if err != nil {
		t.Error(err)
	}
	if field != 1 {
		t.Errorf("invalid field : (%d)", field)
	}

	err = deleteDatabase()
	if err != nil {
		t.Error(err)
	}
}

func TestExecutePrepare(t *testing.T) {
	err := createTable()
	if err != nil {
		t.Error(err)
	}

	client := mysql.Client{}

	err = client.ExecutePrepare(2)
	if err.Error() != `please call SetPrepare first` {
		t.Error(err)
	}

	err = client.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1)
	if err != nil {
		t.Error(err)
	}
	defer client.Finalize()

	err = client.SetPrepare(`INSERT INTO ` + table + ` VALUE(field=?);`)
	if err != nil {
		t.Error(err)
	}

	err = client.ExecutePrepare(2)
	if err != nil {
		t.Error(err)
	}

	err = client.ExecutePrepare(3)
	if err != nil {
		t.Error(err)
	}

	count := 0
	err = client.QueryRow(`SELECT COUNT(*) FROM `+table+`;`, &count)
	if err != nil {
		t.Error(err)
	}
	if count != 3 {
		t.Errorf("invalid count : (%d)", count)
	}

	err = deleteDatabase()
	if err != nil {
		t.Error(err)
	}
}

func TestBeginTransaction(t *testing.T) {
	err := createTable()
	if err != nil {
		t.Error(err)
	}

	client := mysql.Client{}

	err = client.BeginTransaction()
	if err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	err = client.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1)
	if err != nil {
		t.Error(err)
	}
	defer client.Finalize()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()

		err = client.BeginTransaction()
		if err != nil {
			t.Error(err)
		}

		err = client.ExecuteTransaction(`INSERT INTO `+table+` VALUE(field=?);`, 1)
		if err != nil {
			t.Error(err)
		}

		err = client.EndTransaction(err)
		if err != nil {
			t.Error(err)
		}
	}()
	wg.Wait()

	count := 0
	err = client.QueryRow(`SELECT COUNT(*) FROM `+table+`;`, &count)
	if err != nil {
		t.Error(err)
	}
	if count != 2 {
		t.Errorf("invalid count : (%d)", count)
	}

	err = deleteDatabase()
	if err != nil {
		t.Error(err)
	}
}

func TestEndTransaction(t *testing.T) {
	err := createTable()
	if err != nil {
		t.Error(err)
	}

	client := mysql.Client{}

	err = client.EndTransaction(nil)
	if err.Error() != `please call BeginTransaction first` {
		t.Error(err)
	}

	err = client.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1)
	if err != nil {
		t.Error(err)
	}
	defer client.Finalize()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()

		err = client.BeginTransaction()
		if err != nil {
			t.Error(err)
		}

		err = client.ExecuteTransaction(`INSERT INTO `+table+` VALUE(field=?);`, 1)
		if err != nil {
			t.Error(err)
		}

		err = client.ExecuteTransaction(``)
		if err.Error() != `Error 1065 (42000): Query was empty` {
			t.Error(err)
		}

		err = client.EndTransaction(err)
		if err != nil {
			t.Error(err)
		}
	}()
	wg.Wait()

	count := 0
	err = client.QueryRow(`SELECT COUNT(*) FROM `+table+`;`, &count)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Errorf("invalid count : (%d)", count)
	}

	err = deleteDatabase()
	if err != nil {
		t.Error(err)
	}
}

func TestQueryTransaction(t *testing.T) {
	err := createTable()
	if err != nil {
		t.Error(err)
	}

	client := mysql.Client{}

	_, err = client.QueryTransaction(``)
	if err.Error() != `please call BeginTransaction first` {
		t.Error(err)
	}

	err = client.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1)
	if err != nil {
		t.Error(err)
	}
	defer client.Finalize()

	err = client.BeginTransaction()
	if err != nil {
		t.Error(err)
	}

	rows, err := client.QueryTransaction(`SELECT field FROM ` + table + `;`)
	if err != nil {
		t.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		field := 0
		err := rows.Scan(&field)
		if err != nil {
			t.Error(err)
		}

		if field != 1 {
			t.Errorf("invalid field : (%d)", field)
		}
	}

	err = client.EndTransaction(err)
	if err != nil {
		t.Error(err)
	}

	err = deleteDatabase()
	if err != nil {
		t.Error(err)
	}
}

func TestQueryRowTransaction(t *testing.T) {
	err := createTable()
	if err != nil {
		t.Error(err)
	}

	client := mysql.Client{}

	err = client.QueryRowTransaction(``)
	if err.Error() != `please call BeginTransaction first` {
		t.Error(err)
	}

	err = client.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1)
	if err != nil {
		t.Error(err)
	}
	defer client.Finalize()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()

		err = client.BeginTransaction()
		if err != nil {
			t.Error(err)
		}

		field := 0
		err = client.QueryRowTransaction(`SELECT field FROM `+table+`;`, &field)
		if err != nil {
			t.Error(err)
		}
		if field != 1 {
			t.Errorf("invalid field : (%d)", field)
		}

		err = client.EndTransaction(err)
		if err != nil {
			t.Error(err)
		}
	}()
	wg.Wait()

	err = deleteDatabase()
	if err != nil {
		t.Error(err)
	}
}

func TestExecuteTransaction(t *testing.T) {
	err := createTable()
	if err != nil {
		t.Error(err)
	}

	client := mysql.Client{}

	err = client.ExecuteTransaction(``)
	if err.Error() != `please call BeginTransaction first` {
		t.Error(err)
	}

	err = client.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1)
	if err != nil {
		t.Error(err)
	}
	defer client.Finalize()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()

		err = client.BeginTransaction()
		if err != nil {
			t.Error(err)
		}

		err = client.ExecuteTransaction(`INSERT INTO `+table+` VALUE(field=?);`, 1)
		if err != nil {
			t.Error(err)
		}

		err = client.EndTransaction(err)
		if err != nil {
			t.Error(err)
		}
	}()
	wg.Wait()

	count := 0
	err = client.QueryRow(`SELECT COUNT(*) FROM `+table+`;`, &count)
	if err != nil {
		t.Error(err)
	}
	if count != 2 {
		t.Errorf("invalid count : (%d)", count)
	}

	err = deleteDatabase()
	if err != nil {
		t.Error(err)
	}
}

func TestSetPrepareTransaction(t *testing.T) {
	err := createTable()
	if err != nil {
		t.Error(err)
	}

	client := mysql.Client{}

	err = client.SetPrepareTransaction(``)
	if err.Error() != `please call BeginTransaction first` {
		t.Error(err)
	}

	err = client.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1)
	if err != nil {
		t.Error(err)
	}
	defer client.Finalize()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()

		err = client.BeginTransaction()
		if err != nil {
			t.Error(err)
		}

		err = client.SetPrepareTransaction(`INSERT INTO ` + table + ` VALUE(field=?);`)
		if err != nil {
			t.Error(err)
		}

		err = client.ExecutePrepareTransaction(2)
		if err != nil {
			t.Error(err)
		}

		err = client.EndTransaction(err)
		if err != nil {
			t.Error(err)
		}
	}()
	wg.Wait()

	count := 0
	err = client.QueryRow(`SELECT COUNT(*) FROM `+table+`;`, &count)
	if err != nil {
		t.Error(err)
	}
	if count != 2 {
		t.Errorf("invalid count : (%d)", count)
	}

	err = deleteDatabase()
	if err != nil {
		t.Error(err)
	}
}

func TestQueryPrepareTransaction(t *testing.T) {
	err := createTable()
	if err != nil {
		t.Error(err)
	}

	client := mysql.Client{}

	_, err = client.QueryPrepareTransaction(``)
	if err.Error() != `please call SetPrepareTransaction first` {
		t.Error(err)
	}

	err = client.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1)
	if err != nil {
		t.Error(err)
	}
	defer client.Finalize()

	err = client.BeginTransaction()
	if err != nil {
		t.Error(err)
	}

	err = client.SetPrepareTransaction(`SELECT field FROM ` + table + ` WHERE field=?;`)
	if err != nil {
		t.Error(err)
	}

	rows, err := client.QueryPrepareTransaction(1)
	if err != nil {
		t.Error(err)
	}
	defer rows.Close()

	for rows.Next() {
		field := 0
		err := rows.Scan(&field)
		if err != nil {
			t.Error(err)
		}

		if field != 1 {
			t.Errorf("invalid field : (%d)", field)
		}
	}

	err = client.EndTransaction(err)
	if err != nil {
		t.Error(err)
	}

	err = deleteDatabase()
	if err != nil {
		t.Error(err)
	}
}

func TestQueryRowPrepareTransaction(t *testing.T) {
	err := createTable()
	if err != nil {
		t.Error(err)
	}

	client := mysql.Client{}

	_, err = client.QueryRowPrepareTransaction()
	if err.Error() != `please call SetPrepareTransaction first` {
		t.Error(err)
	}

	err = client.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1)
	if err != nil {
		t.Error(err)
	}
	defer client.Finalize()

	err = client.BeginTransaction()
	if err != nil {
		t.Error(err)
	}

	err = client.SetPrepareTransaction(`SELECT field FROM ` + table + ` WHERE field=?;`)
	if err != nil {
		t.Error(err)
	}

	row, err := client.QueryRowPrepareTransaction()
	if err.Error() != `sql: expected 1 arguments, got 0` {
		t.Error(err)
	}

	row, err = client.QueryRowPrepareTransaction(1)
	if err != nil {
		t.Error(err)
	}

	field := 0
	err = row.Scan(&field)
	if err != nil {
		t.Error(err)
	}
	if field != 1 {
		t.Errorf("invalid field : (%d)", field)
	}

	err = client.EndTransaction(err)
	if err != nil {
		t.Error(err)
	}

	err = deleteDatabase()
	if err != nil {
		t.Error(err)
	}
}

func TestExecutePrepareTransaction(t *testing.T) {
	err := createTable()
	if err != nil {
		t.Error(err)
	}

	client := mysql.Client{}

	err = client.ExecutePrepareTransaction()
	if err.Error() != `please call SetPrepareTransaction first` {
		t.Error(err)
	}

	err = client.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1)
	if err != nil {
		t.Error(err)
	}
	defer client.Finalize()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()

		err = client.BeginTransaction()
		if err != nil {
			t.Error(err)
		}

		err = client.SetPrepareTransaction(`INSERT INTO ` + table + ` VALUE(field=?);`)
		if err != nil {
			t.Error(err)
		}

		err = client.ExecutePrepareTransaction(2)
		if err != nil {
			t.Error(err)
		}

		err = client.EndTransaction(err)
		if err != nil {
			t.Error(err)
		}
	}()
	wg.Wait()

	count := 0
	err = client.QueryRow(`SELECT COUNT(*) FROM `+table+`;`, &count)
	if err != nil {
		t.Error(err)
	}
	if count != 2 {
		t.Errorf("invalid count : (%d)", count)
	}

	err = deleteDatabase()
	if err != nil {
		t.Error(err)
	}
}
