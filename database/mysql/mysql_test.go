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

	mysql := mysql.MySQL{}

	err := mysql.Initialize(`root:root@tcp(127.0.0.1)/`, 1)
	if err != nil {
		return err
	}
	defer mysql.Finalize()

	err = mysql.Execute(`CREATE DATABASE ` + database + `;`)
	if err != nil {
		return err
	}

	err = mysql.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1)
	if err != nil {
		return err
	}

	err = mysql.Execute(`CREATE TABLE ` + table + `(field int);`)
	if err != nil {
		return err
	}

	err = mysql.Execute(`INSERT INTO ` + table + `(field) VALUE(1);`)
	if err != nil {
		return err
	}

	return nil
}

func deleteDatabase() error {
	mysql := mysql.MySQL{}

	err := mysql.Initialize(`root:root@tcp(127.0.0.1)/`, 1)
	if err != nil {
		return err
	}
	defer mysql.Finalize()

	err = mysql.Execute(`DROP DATABASE ` + database + `;`)
	if err != nil {
		return err
	}

	return nil
}

func TestInitialize(t *testing.T) {
	mysql := mysql.MySQL{}

	err := mysql.Initialize(`root:root@tcp(127.0.0.1)`, 1)
	if err.Error() != `invalid DSN: missing the slash separating the database name` {
		t.Error(err)
	}

	err = mysql.Initialize(`root:root@tcp(127.0.0.1)/`, 1)
	if err != nil {
		t.Error(err)
	}
	defer mysql.Finalize()
}

func TestFinalize(t *testing.T) {
	mysql := mysql.MySQL{}

	mysql.Finalize()

	err := mysql.Initialize(`root:root@tcp(127.0.0.1)/`, 1)
	if err != nil {
		t.Error(err)
	}
	defer mysql.Finalize()
}

func TestQuery(t *testing.T) {
	err := createTable()
	if err != nil {
		t.Error(err)
	}

	mysql := mysql.MySQL{}

	_, err = mysql.Query(`SELECT field FROM ` + table + `;`)
	if err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	err = mysql.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1)
	if err != nil {
		t.Error(err)
	}
	defer mysql.Finalize()

	rows, err := mysql.Query(`SELECT field FROM ` + table + `;`)
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

	rows2, err := mysql.Query(`SELECT field FROM `+table+` WHERE field=?;`, 1)
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

	mysql := mysql.MySQL{}

	err = mysql.QueryRow(``)
	if err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	err = mysql.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1)
	if err != nil {
		t.Error(err)
	}
	defer mysql.Finalize()

	field := 0
	err = mysql.QueryRow(`SELECT field FROM `+table+`;`, &field)
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

	mysql := mysql.MySQL{}

	err = mysql.Execute(`UPDATE ` + table + ` SET field=2`)
	if err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	err = mysql.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1)
	if err != nil {
		t.Error(err)
	}
	defer mysql.Finalize()

	err = mysql.Execute(``)
	if err.Error() != `Error 1065 (42000): Query was empty` {
		t.Error(err)
	}

	err = mysql.Execute(`UPDATE ` + table + ` SET field=2`)
	if err != nil {
		t.Error(err)
	}

	field := 0
	err = mysql.QueryRow(`SELECT field FROM `+table+`;`, &field)
	if err != nil {
		t.Error(err)
	}
	if field != 2 {
		t.Errorf("invalid field : (%d)", field)
	}

	err = mysql.Execute(`INSERT INTO `+table+` VALUE(field=?);`, 1)
	if err != nil {
		t.Error(err)
	}

	count := 0
	err = mysql.QueryRow(`SELECT COUNT(*) FROM `+table+`;`, &count)
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

	mysql := mysql.MySQL{}

	err = mysql.SetPrepare(`INSERT INTO ` + table + ` VALUE(field=?);`)
	if err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	err = mysql.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1)
	if err != nil {
		t.Error(err)
	}
	defer mysql.Finalize()

	err = mysql.SetPrepare(`INSERT INTO ` + table + ` VALUE(field=?);`)
	if err != nil {
		t.Error(err)
	}

	err = mysql.ExecutePrepare(2)
	if err != nil {
		t.Error(err)
	}

	count := 0
	err = mysql.QueryRow(`SELECT COUNT(*) FROM `+table+`;`, &count)
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

	mysql := mysql.MySQL{}

	_, err = mysql.QueryPrepare(1)
	if err.Error() != `please call SetPrepare first` {
		t.Error(err)
	}

	err = mysql.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1)
	if err != nil {
		t.Error(err)
	}
	defer mysql.Finalize()

	err = mysql.SetPrepare(`SELECT field FROM ` + table + ` WHERE field=?;`)
	if err != nil {
		t.Error(err)
	}

	rows, err := mysql.QueryPrepare(1)
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

	mysql := mysql.MySQL{}

	_, err = mysql.QueryRowPrepare(1)
	if err.Error() != `please call SetPrepare first` {
		t.Error(err)
	}

	err = mysql.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1)
	if err != nil {
		t.Error(err)
	}
	defer mysql.Finalize()

	err = mysql.SetPrepare(`SELECT field FROM ` + table + ` WHERE field=?;`)
	if err != nil {
		t.Error(err)
	}

	row, err := mysql.QueryRowPrepare()
	if err.Error() != `sql: expected 1 arguments, got 0` {
		t.Error(err)
	}

	row, err = mysql.QueryRowPrepare(1)
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

	mysql := mysql.MySQL{}

	err = mysql.ExecutePrepare(2)
	if err.Error() != `please call SetPrepare first` {
		t.Error(err)
	}

	err = mysql.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1)
	if err != nil {
		t.Error(err)
	}
	defer mysql.Finalize()

	err = mysql.SetPrepare(`INSERT INTO ` + table + ` VALUE(field=?);`)
	if err != nil {
		t.Error(err)
	}

	err = mysql.ExecutePrepare(2)
	if err != nil {
		t.Error(err)
	}

	err = mysql.ExecutePrepare(3)
	if err != nil {
		t.Error(err)
	}

	count := 0
	err = mysql.QueryRow(`SELECT COUNT(*) FROM `+table+`;`, &count)
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

	mysql := mysql.MySQL{}

	err = mysql.BeginTransaction()
	if err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	err = mysql.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1)
	if err != nil {
		t.Error(err)
	}
	defer mysql.Finalize()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()

		err = mysql.BeginTransaction()
		if err != nil {
			t.Error(err)
		}

		err = mysql.ExecuteTransaction(`INSERT INTO `+table+` VALUE(field=?);`, 1)
		if err != nil {
			t.Error(err)
		}

		err = mysql.EndTransaction(err)
		if err != nil {
			t.Error(err)
		}
	}()
	wg.Wait()

	count := 0
	err = mysql.QueryRow(`SELECT COUNT(*) FROM `+table+`;`, &count)
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

	mysql := mysql.MySQL{}

	err = mysql.EndTransaction(nil)
	if err.Error() != `please call BeginTransaction first` {
		t.Error(err)
	}

	err = mysql.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1)
	if err != nil {
		t.Error(err)
	}
	defer mysql.Finalize()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()

		err = mysql.BeginTransaction()
		if err != nil {
			t.Error(err)
		}

		err = mysql.ExecuteTransaction(`INSERT INTO `+table+` VALUE(field=?);`, 1)
		if err != nil {
			t.Error(err)
		}

		err = mysql.ExecuteTransaction(``)
		if err.Error() != `Error 1065 (42000): Query was empty` {
			t.Error(err)
		}

		err = mysql.EndTransaction(err)
		if err != nil {
			t.Error(err)
		}
	}()
	wg.Wait()

	count := 0
	err = mysql.QueryRow(`SELECT COUNT(*) FROM `+table+`;`, &count)
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

	mysql := mysql.MySQL{}

	_, err = mysql.QueryTransaction(``)
	if err.Error() != `please call BeginTransaction first` {
		t.Error(err)
	}

	err = mysql.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1)
	if err != nil {
		t.Error(err)
	}
	defer mysql.Finalize()

	err = mysql.BeginTransaction()
	if err != nil {
		t.Error(err)
	}

	rows, err := mysql.QueryTransaction(`SELECT field FROM ` + table + `;`)
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

	err = mysql.EndTransaction(err)
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

	mysql := mysql.MySQL{}

	err = mysql.QueryRowTransaction(``)
	if err.Error() != `please call BeginTransaction first` {
		t.Error(err)
	}

	err = mysql.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1)
	if err != nil {
		t.Error(err)
	}
	defer mysql.Finalize()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()

		err = mysql.BeginTransaction()
		if err != nil {
			t.Error(err)
		}

		field := 0
		err = mysql.QueryRowTransaction(`SELECT field FROM `+table+`;`, &field)
		if err != nil {
			t.Error(err)
		}
		if field != 1 {
			t.Errorf("invalid field : (%d)", field)
		}

		err = mysql.EndTransaction(err)
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

	mysql := mysql.MySQL{}

	err = mysql.ExecuteTransaction(``)
	if err.Error() != `please call BeginTransaction first` {
		t.Error(err)
	}

	err = mysql.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1)
	if err != nil {
		t.Error(err)
	}
	defer mysql.Finalize()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()

		err = mysql.BeginTransaction()
		if err != nil {
			t.Error(err)
		}

		err = mysql.ExecuteTransaction(`INSERT INTO `+table+` VALUE(field=?);`, 1)
		if err != nil {
			t.Error(err)
		}

		err = mysql.EndTransaction(err)
		if err != nil {
			t.Error(err)
		}
	}()
	wg.Wait()

	count := 0
	err = mysql.QueryRow(`SELECT COUNT(*) FROM `+table+`;`, &count)
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

	mysql := mysql.MySQL{}

	err = mysql.SetPrepareTransaction(``)
	if err.Error() != `please call BeginTransaction first` {
		t.Error(err)
	}

	err = mysql.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1)
	if err != nil {
		t.Error(err)
	}
	defer mysql.Finalize()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()

		err = mysql.BeginTransaction()
		if err != nil {
			t.Error(err)
		}

		err = mysql.SetPrepareTransaction(`INSERT INTO ` + table + ` VALUE(field=?);`)
		if err != nil {
			t.Error(err)
		}

		err = mysql.ExecutePrepareTransaction(2)
		if err != nil {
			t.Error(err)
		}

		err = mysql.EndTransaction(err)
		if err != nil {
			t.Error(err)
		}
	}()
	wg.Wait()

	count := 0
	err = mysql.QueryRow(`SELECT COUNT(*) FROM `+table+`;`, &count)
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

	mysql := mysql.MySQL{}

	_, err = mysql.QueryPrepareTransaction(``)
	if err.Error() != `please call SetPrepareTransaction first` {
		t.Error(err)
	}

	err = mysql.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1)
	if err != nil {
		t.Error(err)
	}
	defer mysql.Finalize()

	err = mysql.BeginTransaction()
	if err != nil {
		t.Error(err)
	}

	err = mysql.SetPrepareTransaction(`SELECT field FROM ` + table + ` WHERE field=?;`)
	if err != nil {
		t.Error(err)
	}

	rows, err := mysql.QueryPrepareTransaction(1)
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

	err = mysql.EndTransaction(err)
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

	mysql := mysql.MySQL{}

	_, err = mysql.QueryRowPrepareTransaction()
	if err.Error() != `please call SetPrepareTransaction first` {
		t.Error(err)
	}

	err = mysql.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1)
	if err != nil {
		t.Error(err)
	}
	defer mysql.Finalize()

	err = mysql.BeginTransaction()
	if err != nil {
		t.Error(err)
	}

	err = mysql.SetPrepareTransaction(`SELECT field FROM ` + table + ` WHERE field=?;`)
	if err != nil {
		t.Error(err)
	}

	row, err := mysql.QueryRowPrepareTransaction()
	if err.Error() != `sql: expected 1 arguments, got 0` {
		t.Error(err)
	}

	row, err = mysql.QueryRowPrepareTransaction(1)
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

	err = mysql.EndTransaction(err)
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

	mysql := mysql.MySQL{}

	err = mysql.ExecutePrepareTransaction()
	if err.Error() != `please call SetPrepareTransaction first` {
		t.Error(err)
	}

	err = mysql.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1)
	if err != nil {
		t.Error(err)
	}
	defer mysql.Finalize()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()

		err = mysql.BeginTransaction()
		if err != nil {
			t.Error(err)
		}

		err = mysql.SetPrepareTransaction(`INSERT INTO ` + table + ` VALUE(field=?);`)
		if err != nil {
			t.Error(err)
		}

		err = mysql.ExecutePrepareTransaction(2)
		if err != nil {
			t.Error(err)
		}

		err = mysql.EndTransaction(err)
		if err != nil {
			t.Error(err)
		}
	}()
	wg.Wait()

	count := 0
	err = mysql.QueryRow(`SELECT COUNT(*) FROM `+table+`;`, &count)
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
