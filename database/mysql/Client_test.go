package mysql_test

import (
	"database/sql"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/common-library/go/database/mysql"
	"github.com/google/uuid"
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

	if err := client.Initialize(`root:root@tcp(127.0.0.1)/`, 1); err != nil {
		return err
	}
	defer client.Finalize()

	if err := client.Execute(`CREATE DATABASE ` + database + `;`); err != nil {
		return err
	}

	if err := client.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1); err != nil {
		return err
	}

	if err := client.Execute(`CREATE TABLE ` + table + `(field int);`); err != nil {
		return err
	}

	if err := client.Execute(`INSERT INTO ` + table + `(field) VALUE(1);`); err != nil {
		return err
	}

	return nil
}

func deleteDatabase() error {
	client := mysql.Client{}

	if err := client.Initialize(`root:root@tcp(127.0.0.1)/`, 1); err != nil {
		return err
	}
	defer client.Finalize()

	if err := client.Execute(`DROP DATABASE ` + database + `;`); err != nil {
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

	if err := client.Initialize(`root:root@tcp(127.0.0.1)/`, 1); err != nil {
		t.Error(err)
	}
	defer client.Finalize()
}

func TestFinalize(t *testing.T) {
	client := mysql.Client{}

	client.Finalize()

	if err := client.Initialize(`root:root@tcp(127.0.0.1)/`, 1); err != nil {
		t.Error(err)
	}
	defer client.Finalize()
}

func TestQuery(t *testing.T) {
	if err := createTable(); err != nil {
		t.Error(err)
	}

	client := mysql.Client{}

	if _, err := client.Query(""); err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	if err := client.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1); err != nil {
		t.Error(err)
	}
	defer client.Finalize()

	check := func(query string, args ...any) {
		if rows, err := client.Query(query, args...); err != nil {
			t.Error(err)
		} else {
			defer rows.Close()

			for rows.Next() {
				field := 0
				if err := rows.Scan(&field); err != nil {
					t.Error(err)
				}

				if field != 1 {
					t.Errorf("invalid field : (%d)", field)
				}
			}
		}
	}

	check(`SELECT field FROM ` + table + `;`)
	check(`SELECT field FROM `+table+` WHERE field=?;`, 1)

	if err := deleteDatabase(); err != nil {
		t.Error(err)
	}
}

func TestQueryRow(t *testing.T) {
	if err := createTable(); err != nil {
		t.Error(err)
	}

	client := mysql.Client{}

	if err := client.QueryRow(``); err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	if err := client.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1); err != nil {
		t.Error(err)
	}
	defer client.Finalize()

	field := 0
	if err := client.QueryRow(`SELECT field FROM `+table+`;`, &field); err != nil {
		t.Error(err)
	} else if field != 1 {
		t.Errorf("invalid field : (%d)", field)
	}

	if err := deleteDatabase(); err != nil {
		t.Error(err)
	}
}

func TestExecute(t *testing.T) {
	if err := createTable(); err != nil {
		t.Error(err)
	}

	client := mysql.Client{}

	if err := client.Execute(""); err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	if err := client.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1); err != nil {
		t.Error(err)
	}
	defer client.Finalize()

	if err := client.Execute(``); err.Error() != `Error 1065 (42000): Query was empty` {
		t.Error(err)
	}

	if err := client.Execute(`UPDATE ` + table + ` SET field=2`); err != nil {
		t.Error(err)
	}

	field := 0
	if err := client.QueryRow(`SELECT field FROM `+table+`;`, &field); err != nil {
		t.Error(err)
	} else if field != 2 {
		t.Errorf("invalid field : (%d)", field)
	}

	if err := client.Execute(`INSERT INTO `+table+` VALUE(field=?);`, 1); err != nil {
		t.Error(err)
	}

	count := 0
	if err := client.QueryRow(`SELECT COUNT(*) FROM `+table+`;`, &count); err != nil {
		t.Error(err)
	} else if count != 2 {
		t.Errorf("invalid count : (%d)", count)
	}

	if err := deleteDatabase(); err != nil {
		t.Error(err)
	}
}

func TestSetPrepare(t *testing.T) {
	if err := createTable(); err != nil {
		t.Error(err)
	}

	client := mysql.Client{}

	if err := client.SetPrepare(""); err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	if err := client.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1); err != nil {
		t.Error(err)
	}
	defer client.Finalize()

	if err := client.SetPrepare(`INSERT INTO ` + table + ` VALUE(field=?);`); err != nil {
		t.Error(err)
	}

	if err := client.ExecutePrepare(2); err != nil {
		t.Error(err)
	}

	count := 0
	if err := client.QueryRow(`SELECT COUNT(*) FROM `+table+`;`, &count); err != nil {
		t.Error(err)
	} else if count != 2 {
		t.Errorf("invalid count : (%d)", count)
	}

	if err := deleteDatabase(); err != nil {
		t.Error(err)
	}
}

func TestQueryPrepare(t *testing.T) {
	if err := createTable(); err != nil {
		t.Error(err)
	}

	client := mysql.Client{}

	if _, err := client.QueryPrepare(1); err.Error() != `please call SetPrepare first` {
		t.Error(err)
	}

	if err := client.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1); err != nil {
		t.Error(err)
	}
	defer client.Finalize()

	if err := client.SetPrepare(`SELECT field FROM ` + table + ` WHERE field=?;`); err != nil {
		t.Error(err)
	}

	if rows, err := client.QueryPrepare(1); err != nil {
		t.Error(err)
	} else {
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
	}

	if err := deleteDatabase(); err != nil {
		t.Error(err)
	}
}

func TestQueryRowPrepare(t *testing.T) {
	if err := createTable(); err != nil {
		t.Error(err)
	}

	client := mysql.Client{}

	if _, err := client.QueryRowPrepare(1); err.Error() != `please call SetPrepare first` {
		t.Error(err)
	}

	if err := client.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1); err != nil {
		t.Error(err)
	}
	defer client.Finalize()

	if err := client.SetPrepare(`SELECT field FROM ` + table + ` WHERE field=?;`); err != nil {
		t.Error(err)
	}

	if _, err := client.QueryRowPrepare(); err.Error() != `sql: expected 1 arguments, got 0` {
		t.Error(err)
	}

	if row, err := client.QueryRowPrepare(1); err != nil {
		t.Error(err)
	} else {
		field := 0
		err = row.Scan(&field)
		if err != nil {
			t.Error(err)
		}
		if field != 1 {
			t.Errorf("invalid field : (%d)", field)
		}
	}

	if err := deleteDatabase(); err != nil {
		t.Error(err)
	}
}

func TestExecutePrepare(t *testing.T) {
	if err := createTable(); err != nil {
		t.Error(err)
	}

	client := mysql.Client{}

	if err := client.ExecutePrepare(2); err.Error() != `please call SetPrepare first` {
		t.Error(err)
	}

	if err := client.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1); err != nil {
		t.Error(err)
	}
	defer client.Finalize()

	if err := client.SetPrepare(`INSERT INTO ` + table + ` VALUE(field=?);`); err != nil {
		t.Error(err)
	}

	if err := client.ExecutePrepare(2); err != nil {
		t.Error(err)
	}

	if err := client.ExecutePrepare(3); err != nil {
		t.Error(err)
	}

	count := 0
	if err := client.QueryRow(`SELECT COUNT(*) FROM `+table+`;`, &count); err != nil {
		t.Error(err)
	} else if count != 3 {
		t.Errorf("invalid count : (%d)", count)
	}

	if err := deleteDatabase(); err != nil {
		t.Error(err)
	}
}

func TestBeginTransaction(t *testing.T) {
	if err := createTable(); err != nil {
		t.Error(err)
	}

	client := mysql.Client{}

	if err := client.BeginTransaction(); err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	if err := client.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1); err != nil {
		t.Error(err)
	}
	defer client.Finalize()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()

		err := new(error)
		defer func(errForTransaction *error) {
			if err := client.EndTransaction(*errForTransaction); err != nil {
				t.Error(err)
			}
		}(err)

		if *err = client.BeginTransaction(); *err != nil {
			t.Error(*err)
			return
		}

		if *err = client.ExecuteTransaction(`INSERT INTO `+table+` VALUE(field=?);`, 1); *err != nil {
			t.Error(*err)
			return
		}
	}()
	wg.Wait()

	count := 0
	if err := client.QueryRow(`SELECT COUNT(*) FROM `+table+`;`, &count); err != nil {
		t.Error(err)
	} else if count != 2 {
		t.Errorf("invalid count : (%d)", count)
	}

	if err := deleteDatabase(); err != nil {
		t.Error(err)
	}
}

func TestEndTransaction(t *testing.T) {
	if err := createTable(); err != nil {
		t.Error(err)
	}

	client := mysql.Client{}

	if err := client.EndTransaction(nil); err.Error() != `please call BeginTransaction first` {
		t.Error(err)
	}

	if err := client.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1); err != nil {
		t.Error(err)
	}
	defer client.Finalize()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()

		err := new(error)
		defer func(errForTransaction *error) {
			if err := client.EndTransaction(*errForTransaction); err != nil {
				t.Error(err)
			}
		}(err)

		if *err = client.BeginTransaction(); *err != nil {
			t.Error(*err)
			return
		}

		if *err = client.ExecuteTransaction(`INSERT INTO `+table+` VALUE(field=?);`, 1); *err != nil {
			t.Error(*err)
			return
		}

		if *err = client.ExecuteTransaction(``); (*err).Error() != `Error 1065 (42000): Query was empty` {
			t.Error(*err)
		}
	}()
	wg.Wait()

	count := 0
	if err := client.QueryRow(`SELECT COUNT(*) FROM `+table+`;`, &count); err != nil {
		t.Error(err)
	} else if count != 1 {
		t.Errorf("invalid count : (%d)", count)
	}

	if err := deleteDatabase(); err != nil {
		t.Error(err)
	}
}

func TestQueryTransaction(t *testing.T) {
	if err := createTable(); err != nil {
		t.Error(err)
	}

	client := mysql.Client{}

	if _, err := client.QueryTransaction(``); err.Error() != `please call BeginTransaction first` {
		t.Error(err)
	}

	if err := client.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1); err != nil {
		t.Error(err)
	}
	defer client.Finalize()

	func() {
		err := new(error)
		defer func(errForTransaction *error) {
			if err := client.EndTransaction(*errForTransaction); err != nil {
				t.Error(err)
			}
		}(err)

		if *err = client.BeginTransaction(); *err != nil {
			t.Error(*err)
			return
		}

		rows := new(sql.Rows)
		if rows, *err = client.QueryTransaction(`SELECT field FROM ` + table + `;`); *err != nil {
			t.Error(*err)
		} else {
			defer rows.Close()

			for rows.Next() {
				field := 0
				if *err = rows.Scan(&field); *err != nil {
					t.Error(*err)
					break
				} else if field != 1 {
					t.Errorf("invalid field : (%d)", field)
				}
			}
		}
	}()

	if err := deleteDatabase(); err != nil {
		t.Error(err)
	}
}

func TestQueryRowTransaction(t *testing.T) {
	if err := createTable(); err != nil {
		t.Error(err)
	}

	client := mysql.Client{}

	if err := client.QueryRowTransaction(``); err.Error() != `please call BeginTransaction first` {
		t.Error(err)
	}

	if err := client.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1); err != nil {
		t.Error(err)
	}
	defer client.Finalize()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()

		err := new(error)
		defer func(errForTransaction *error) {
			if err := client.EndTransaction(*errForTransaction); err != nil {
				t.Error(err)
			}
		}(err)

		if *err = client.BeginTransaction(); *err != nil {
			t.Error(*err)
			return
		}

		field := 0
		if *err = client.QueryRowTransaction(`SELECT field FROM `+table+`;`, &field); *err != nil {
			t.Error(*err)
			return
		} else if field != 1 {
			t.Errorf("invalid field : (%d)", field)
		}
	}()
	wg.Wait()

	if err := deleteDatabase(); err != nil {
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

		err := new(error)
		defer func(errForTransaction *error) {
			if err := client.EndTransaction(*errForTransaction); err != nil {
				t.Error(err)
			}
		}(err)

		if *err = client.BeginTransaction(); *err != nil {
			t.Error(*err)
		}

		if *err = client.ExecuteTransaction(`INSERT INTO `+table+` VALUE(field=?);`, 1); *err != nil {
			t.Error(*err)
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

	if err := client.SetPrepareTransaction(``); err.Error() != `please call BeginTransaction first` {
		t.Error(err)
	}

	if err := client.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1); err != nil {
		t.Error(err)
	}
	defer client.Finalize()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()

		err := new(error)
		defer func(errForTransaction *error) {
			if err := client.EndTransaction(*errForTransaction); err != nil {
				t.Error(err)
			}
		}(err)

		if *err = client.BeginTransaction(); *err != nil {
			t.Error(*err)
		}

		if *err = client.SetPrepareTransaction(`INSERT INTO ` + table + ` VALUE(field=?);`); *err != nil {
			t.Error(*err)
		}

		if *err = client.ExecutePrepareTransaction(2); *err != nil {
			t.Error(err)
		}
	}()
	wg.Wait()

	count := 0
	if err := client.QueryRow(`SELECT COUNT(*) FROM `+table+`;`, &count); err != nil {
		t.Error(err)
	} else if count != 2 {
		t.Errorf("invalid count : (%d)", count)
	}

	if err := deleteDatabase(); err != nil {
		t.Error(err)
	}
}

func TestQueryPrepareTransaction(t *testing.T) {
	if err := createTable(); err != nil {
		t.Error(err)
	}

	client := mysql.Client{}

	if _, err := client.QueryPrepareTransaction(``); err.Error() != `please call SetPrepareTransaction first` {
		t.Error(err)
	}

	if err := client.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1); err != nil {
		t.Error(err)
	}
	defer client.Finalize()

	func() {
		err := new(error)
		defer func(errForTransaction *error) {
			if err := client.EndTransaction(*errForTransaction); err != nil {
				t.Error(err)
			}
		}(err)

		if *err = client.BeginTransaction(); *err != nil {
			t.Error(*err)
			return
		}

		if *err = client.SetPrepareTransaction(`SELECT field FROM ` + table + ` WHERE field=?;`); *err != nil {
			t.Error(*err)
			return
		}

		rows := new(sql.Rows)
		if rows, *err = client.QueryPrepareTransaction(1); *err != nil {
			t.Error(*err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			field := 0
			if *err = rows.Scan(&field); *err != nil {
				t.Error(*err)
				return
			} else if field != 1 {
				t.Errorf("invalid field : (%d)", field)
			}
		}
	}()

	if err := deleteDatabase(); err != nil {
		t.Error(err)
	}
}

func TestQueryRowPrepareTransaction(t *testing.T) {
	if err := createTable(); err != nil {
		t.Error(err)
	}

	client := mysql.Client{}

	if _, err := client.QueryRowPrepareTransaction(); err.Error() != `please call SetPrepareTransaction first` {
		t.Error(err)
	}

	if err := client.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1); err != nil {
		t.Error(err)
	}
	defer client.Finalize()

	func() {
		err := new(error)
		defer func(errForTransaction *error) {
			if err := client.EndTransaction(*errForTransaction); err != nil {
				t.Error(err)
			}
		}(err)

		if *err = client.BeginTransaction(); *err != nil {
			t.Error(*err)
			return
		}

		if *err = client.SetPrepareTransaction(`SELECT field FROM ` + table + ` WHERE field=?;`); *err != nil {
			t.Error(*err)
			return
		}

		row := new(sql.Row)
		if row, *err = client.QueryRowPrepareTransaction(); (*err).Error() != `sql: expected 1 arguments, got 0` {
			t.Error(*err)
			return
		}

		if row, *err = client.QueryRowPrepareTransaction(1); *err != nil {
			t.Error(err)
			return
		}

		field := 0
		if *err = row.Scan(&field); *err != nil {
			t.Error(*err)
		} else if field != 1 {
			t.Errorf("invalid field : (%d)", field)
		}
	}()

	if err := deleteDatabase(); err != nil {
		t.Error(err)
	}
}

func TestExecutePrepareTransaction(t *testing.T) {
	if err := createTable(); err != nil {
		t.Error(err)
	}

	client := mysql.Client{}

	if err := client.ExecutePrepareTransaction(); err.Error() != `please call SetPrepareTransaction first` {
		t.Error(err)
	}

	if err := client.Initialize(`root:root@tcp(127.0.0.1)/`+database, 1); err != nil {
		t.Error(err)
	}
	defer client.Finalize()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()

		err := new(error)
		defer func(errForTransaction *error) {
			if err := client.EndTransaction(*errForTransaction); err != nil {
				t.Error(err)
			}
		}(err)

		if *err = client.BeginTransaction(); *err != nil {
			t.Error(*err)
			return
		}

		if *err = client.SetPrepareTransaction(`INSERT INTO ` + table + ` VALUE(field=?);`); *err != nil {
			t.Error(*err)
			return
		}

		if *err = client.ExecutePrepareTransaction(2); *err != nil {
			t.Error(*err)
			return
		}
	}()
	wg.Wait()

	count := 0
	if err := client.QueryRow(`SELECT COUNT(*) FROM `+table+`;`, &count); err != nil {
		t.Error(err)
	} else if count != 2 {
		t.Errorf("invalid count : (%d)", count)
	}

	if err := deleteDatabase(); err != nil {
		t.Error(err)
	}
}
