package sql_test

import (
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/common-library/go/database/sql"
	"github.com/common-library/go/file"
	"github.com/google/uuid"
)

var SQLiteFileName = uuid.NewString() + ".db"
var Prepare = map[sql.Driver]string{
	sql.DriverAmazonDynamoDB: "?",
	sql.DriverMySQL:          "?",
	sql.DriverPostgres:       "$1",
	sql.DriverSQLite:         "?",
}

var Drivers = []sql.Driver{
	sql.DriverAmazonDynamoDB,
	sql.DriverMySQL,
	sql.DriverPostgres,
	sql.DriverSQLite,
}

func setUp() {
}

func tearDown() {
	file.RemoveAll(SQLiteFileName)
}

func TestMain(m *testing.M) {
	setUp()

	code := m.Run()

	tearDown()

	os.Exit(code)
}

func TestOpen(t *testing.T) {
	job := func(*testing.T, sql.Client, sql.Driver, string) {}

	test(t, job)
}

func TestClose(t *testing.T) {
	job := func(*testing.T, sql.Client, sql.Driver, string) {}

	test(t, job)
}

func TestQuery(t *testing.T) {
	if _, err := (&sql.Client{}).Query(""); err.Error() != `please call Open first` {
		t.Fatal(err)
	}

	check := func(client sql.Client, query string, args ...any) {
		if rows, err := client.Query(query, args...); err != nil {
			t.Fatal(err)
		} else {
			defer rows.Close()

			for rows.Next() {
				field := 0
				if err := rows.Scan(&field); err != nil {
					t.Fatal(err)
				}

				if field != 1 {
					t.Fatal("invalid -", field)
				}
			}
		}
	}

	job := func(t *testing.T, client sql.Client, driver sql.Driver, table string) {
		check(client, `SELECT field FROM `+table+`;`)
		check(client, `SELECT field FROM `+table+` WHERE field=`+Prepare[driver]+`;`, 1)
	}

	test(t, job)
}

func TestQueryRow(t *testing.T) {
	if err := (&sql.Client{}).QueryRow(``); err.Error() != `please call Open first` {
		t.Fatal(err)
	}

	job := func(t *testing.T, client sql.Client, driver sql.Driver, table string) {
		field := 0
		if err := client.QueryRow(`SELECT field FROM `+table+`;`, &field); err != nil {
			t.Fatal(err)
		} else if field != 1 {
			t.Fatal("invalid -", field)
		}
	}

	test(t, job)
}

func TestExecute(t *testing.T) {
	if err := (&sql.Client{}).Execute(""); err.Error() != `please call Open first` {
		t.Fatal(err)
	}

	job := func(t *testing.T, client sql.Client, driver sql.Driver, table string) {
		errorString := map[sql.Driver]string{
			sql.DriverAmazonDynamoDB: `invalid query:`,
			sql.DriverMySQL:          `Error 1065 (42000): Query was empty`,
			sql.DriverPostgres:       `no RowsAffected available after the empty statement`,
			sql.DriverSQLite:         ``,
		}
		if err := client.Execute(``); err != nil && err.Error() != errorString[driver] {
			t.Fatal(err)
		}

		if err := client.Execute(`UPDATE ` + table + ` SET field=2`); err != nil {
			t.Fatal(err)
		}

		field := 0
		if err := client.QueryRow(`SELECT field FROM `+table+`;`, &field); err != nil {
			t.Fatal(err)
		} else if field != 2 {
			t.Fatal("invalid -", field)
		}

		if err := client.Execute(`INSERT INTO `+table+`(field) VALUES(`+Prepare[driver]+`);`, 1); err != nil {
			t.Fatal(err)
		}

		count := 0
		if err := client.QueryRow(`SELECT COUNT(*) FROM `+table+`;`, &count); err != nil {
			t.Fatal(err)
		} else if count != 2 {
			t.Fatal("invalid -", count)
		}
	}

	test(t, job)
}

func TestSetPrepare(t *testing.T) {
	if err := (&sql.Client{}).SetPrepare(""); err.Error() != `please call Open first` {
		t.Fatal(err)
	}

	job := func(t *testing.T, client sql.Client, driver sql.Driver, table string) {
		if err := client.SetPrepare(`INSERT INTO ` + table + `(field) VALUES(` + Prepare[driver] + `);`); err != nil {
			t.Fatal(err)
		}

		if err := client.ExecutePrepare(2); err != nil {
			t.Fatal(err)
		}

		count := 0
		if err := client.QueryRow(`SELECT COUNT(*) FROM `+table+`;`, &count); err != nil {
			t.Fatal(err)
		} else if count != 2 {
			t.Fatal("invalid -", count)
		}
	}

	test(t, job)
}

func TestQueryPrepare(t *testing.T) {
	if _, err := (&sql.Client{}).QueryPrepare(1); err.Error() != `please call SetPrepare first` {
		t.Fatal(err)
	}

	job := func(t *testing.T, client sql.Client, driver sql.Driver, table string) {
		if err := client.SetPrepare(`SELECT field FROM ` + table + ` WHERE field=` + Prepare[driver] + `;`); err != nil {
			t.Fatal(err)
		}

		if rows, err := client.QueryPrepare(1); err != nil {
			t.Fatal(err)
		} else {
			defer rows.Close()

			for rows.Next() {
				field := 0
				if err := rows.Scan(&field); err != nil {
					t.Fatal(err)
				} else if field != 1 {
					t.Fatal("invalid -", field)
				}
			}
		}
	}

	test(t, job)
}

func TestQueryRowPrepare(t *testing.T) {
	if _, err := (&sql.Client{}).QueryRowPrepare(1); err.Error() != `please call SetPrepare first` {
		t.Fatal(err)
	}

	job := func(t *testing.T, client sql.Client, driver sql.Driver, table string) {
		if err := client.SetPrepare(`SELECT field FROM ` + table + ` WHERE field=` + Prepare[driver] + `;`); err != nil {
			t.Fatal(err)
		}

		if _, err := client.QueryRowPrepare(); err.Error() != `sql: expected 1 arguments, got 0` {
			t.Fatal(err)
		}

		if row, err := client.QueryRowPrepare(1); err != nil {
			t.Fatal(err)
		} else {
			field := 0
			if err := row.Scan(&field); err != nil {
				t.Fatal(err)
			} else if field != 1 {
				t.Fatal("invalid -", field)
			}
		}
	}

	test(t, job)
}

func TestExecutePrepare(t *testing.T) {
	if err := (&sql.Client{}).ExecutePrepare(2); err.Error() != `please call SetPrepare first` {
		t.Fatal(err)
	}

	job := func(t *testing.T, client sql.Client, driver sql.Driver, table string) {
		if err := client.SetPrepare(`INSERT INTO ` + table + `(field) VALUES(` + Prepare[driver] + `);`); err != nil {
			t.Fatal(err)
		}

		if err := client.ExecutePrepare(2); err != nil {
			t.Fatal(err)
		}

		if err := client.ExecutePrepare(3); err != nil {
			t.Fatal(err)
		}

		count := 0
		if err := client.QueryRow(`SELECT COUNT(*) FROM `+table+`;`, &count); err != nil {
			t.Fatal(err)
		} else if count != 3 {
			t.Fatal("invalid -", count)
		}
	}

	test(t, job)
}

func TestBeginTransaction(t *testing.T) {
	if err := (&sql.Client{}).BeginTransaction(); err.Error() != `please call Open first` {
		t.Fatal(err)
	}

	job := func(t *testing.T, client sql.Client, driver sql.Driver, table string) {
		wg := new(sync.WaitGroup)
		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := client.BeginTransaction(); err != nil {
				t.Fatal(err)
			}

			execute := func() error {
				if err := client.ExecuteTransaction(`INSERT INTO `+table+`(field) VALUES(`+Prepare[driver]+`);`, 1); err != nil {
					return err
				} else {
					return nil
				}
			}

			errExecute := execute()
			if errExecute != nil {
				t.Error(errExecute)
			}

			if err := client.EndTransaction(errExecute); err != nil {
				t.Fatal(err)
			}
		}()
		wg.Wait()

		count := 0
		if err := client.QueryRow(`SELECT COUNT(*) FROM `+table+`;`, &count); err != nil {
			t.Fatal(err)
		} else if count != 2 {
			t.Fatal("invalid -", count)
		}
	}

	test(t, job)
}

func TestEndTransaction(t *testing.T) {
	if err := (&sql.Client{}).EndTransaction(nil); err.Error() != `please call BeginTransaction first` {
		t.Fatal(err)
	}

	job := func(t *testing.T, client sql.Client, driver sql.Driver, table string) {
		wg := new(sync.WaitGroup)
		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := client.BeginTransaction(); err != nil {
				t.Fatal(err)
			}

			execute := func() error {
				errorString := map[sql.Driver]string{
					sql.DriverMySQL:    `Error 1065 (42000): Query was empty`,
					sql.DriverPostgres: `no RowsAffected available after the empty statement`,
					sql.DriverSQLite:   ``,
				}

				if err := client.ExecuteTransaction(`INSERT INTO `+table+`(field) VALUES(`+Prepare[driver]+`);`, 1); err != nil {
					t.Error(err)
					return err
				} else if err := client.ExecuteTransaction(``); err != nil && err.Error() != errorString[driver] {
					t.Error(err)
					return err
				} else {
					return nil
				}
			}

			errExecute := execute()
			if errExecute != nil {
				t.Error(errExecute)
			}

			if err := client.EndTransaction(errExecute); err != nil {
				t.Fatal(err)
			}
		}()
		wg.Wait()

		count := 0
		if err := client.QueryRow(`SELECT COUNT(*) FROM `+table+`;`, &count); err != nil {
			t.Fatal(err)
		} else if count != 2 {
			t.Fatal("invalid -", count)
		}
	}

	test(t, job)
}

func TestQueryTransaction(t *testing.T) {
	if _, err := (&sql.Client{}).QueryTransaction(``); err.Error() != `please call BeginTransaction first` {
		t.Fatal(err)
	}

	job := func(t *testing.T, client sql.Client, driver sql.Driver, table string) {
		if err := client.BeginTransaction(); err != nil {
			t.Fatal(err)
		}

		query := func() error {
			if rows, err := client.QueryTransaction(`SELECT field FROM ` + table + `;`); err != nil {
				return err
			} else {
				defer rows.Close()

				for rows.Next() {
					field := 0
					if err := rows.Scan(&field); err != nil {
						return err
					} else if field != 1 {
						t.Fatal("invalid -", field)
					}
				}

				return nil
			}
		}

		errQuery := query()
		if errQuery != nil {
			t.Error(errQuery)
		}

		if err := client.EndTransaction(errQuery); err != nil {
			t.Fatal(err)
		}
	}

	test(t, job)
}

func TestQueryRowTransaction(t *testing.T) {
	if err := (&sql.Client{}).QueryRowTransaction(``); err.Error() != `please call BeginTransaction first` {
		t.Fatal(err)
	}

	job := func(t *testing.T, client sql.Client, driver sql.Driver, table string) {
		wg := new(sync.WaitGroup)
		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := client.BeginTransaction(); err != nil {
				t.Fatal(err)
			}

			query := func() error {
				field := 0
				if err := client.QueryRowTransaction(`SELECT field FROM `+table+`;`, &field); err != nil {
					return err
				} else if field != 1 {
					t.Fatal("invalid -", field)
				}

				return nil
			}

			errQuery := query()
			if errQuery != nil {
				t.Error(errQuery)
			}

			if err := client.EndTransaction(errQuery); err != nil {
				t.Fatal(err)
			}
		}()
		wg.Wait()
	}

	test(t, job)
}

func TestExecuteTransaction(t *testing.T) {
	if err := (&sql.Client{}).ExecuteTransaction(``); err.Error() != `please call BeginTransaction first` {
		t.Fatal(err)
	}

	job := func(t *testing.T, client sql.Client, driver sql.Driver, table string) {
		wg := new(sync.WaitGroup)
		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := client.BeginTransaction(); err != nil {
				t.Fatal(err)
			}

			execute := func() error {
				return client.ExecuteTransaction(`INSERT INTO `+table+`(field) VALUES(`+Prepare[driver]+`);`, 1)
			}

			errExecute := execute()
			if errExecute != nil {
				t.Error(errExecute)
			}

			if err := client.EndTransaction(errExecute); err != nil {
				t.Fatal(err)
			}
		}()
		wg.Wait()

		count := 0
		if err := client.QueryRow(`SELECT COUNT(*) FROM `+table+`;`, &count); err != nil {
			t.Fatal(err)
		} else if count != 2 {
			t.Fatal("invalid -", count)
		}
	}

	test(t, job)
}

func TestSetPrepareTransaction(t *testing.T) {
	if err := (&sql.Client{}).SetPrepareTransaction(``); err.Error() != `please call BeginTransaction first` {
		t.Fatal(err)
	}

	job := func(t *testing.T, client sql.Client, driver sql.Driver, table string) {
		wg := new(sync.WaitGroup)
		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := client.BeginTransaction(); err != nil {
				t.Fatal(err)
			}

			execute := func() error {
				if err := client.SetPrepareTransaction(`INSERT INTO ` + table + `(field) VALUES(` + Prepare[driver] + `);`); err != nil {
					return err
				} else if err := client.ExecutePrepareTransaction(2); err != nil {
					return err
				} else {
					return nil
				}
			}

			errExecute := execute()
			if errExecute != nil {
				t.Error(errExecute)
			}

			if err := client.EndTransaction(errExecute); err != nil {
				t.Fatal(err)
			}
		}()
		wg.Wait()

		count := 0
		if err := client.QueryRow(`SELECT COUNT(*) FROM `+table+`;`, &count); err != nil {
			t.Fatal(err)
		} else if count != 2 {
			t.Fatal("invalid -", count)
		}
	}

	test(t, job)
}

func TestQueryPrepareTransaction(t *testing.T) {
	if _, err := (&sql.Client{}).QueryPrepareTransaction(``); err.Error() != `please call SetPrepareTransaction first` {
		t.Fatal(err)
	}

	job := func(t *testing.T, client sql.Client, driver sql.Driver, table string) {
		if err := client.BeginTransaction(); err != nil {
			t.Fatal(err)
		}

		query := func() error {
			if err := client.SetPrepareTransaction(`SELECT field FROM ` + table + ` WHERE field=` + Prepare[driver] + `;`); err != nil {
				return err
			}

			if rows, err := client.QueryPrepareTransaction(1); err != nil {
				return err
			} else {
				defer rows.Close()

				for rows.Next() {
					field := 0
					if err := rows.Scan(&field); err != nil {
						return err
					} else if field != 1 {
						t.Fatal("invalid -", field)
					}
				}
			}

			return nil
		}

		errQuery := query()
		if errQuery != nil {
			t.Error(errQuery)
		}

		if err := client.EndTransaction(errQuery); err != nil {
			t.Fatal(err)
		}
	}

	test(t, job)
}

func TestQueryRowPrepareTransaction(t *testing.T) {
	if _, err := (&sql.Client{}).QueryRowPrepareTransaction(); err.Error() != `please call SetPrepareTransaction first` {
		t.Fatal(err)
	}

	job := func(t *testing.T, client sql.Client, driver sql.Driver, table string) {
		if err := client.BeginTransaction(); err != nil {
			t.Fatal(err)
		}

		query := func() error {
			if err := client.SetPrepareTransaction(`SELECT field FROM ` + table + ` WHERE field=` + Prepare[driver] + `;`); err != nil {
				return err
			}

			if _, err := client.QueryRowPrepareTransaction(); err.Error() != `sql: expected 1 arguments, got 0` {
				return err
			} else if row, err := client.QueryRowPrepareTransaction(1); err != nil {
				return err
			} else {
				field := 0
				if err := row.Scan(&field); err != nil {
					return err
				} else if field != 1 {
					t.Fatal("invalid -", field)
				}

				return nil
			}
		}

		errQuery := query()
		if errQuery != nil {
			t.Error(errQuery)
		}

		if err := client.EndTransaction(errQuery); err != nil {
			t.Fatal(err)
		}
	}

	test(t, job)
}

func TestExecutePrepareTransaction(t *testing.T) {
	if err := (&sql.Client{}).ExecutePrepareTransaction(); err.Error() != `please call SetPrepareTransaction first` {
		t.Fatal(err)
	}

	job := func(t *testing.T, client sql.Client, driver sql.Driver, table string) {
		wg := new(sync.WaitGroup)
		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := client.BeginTransaction(); err != nil {
				t.Fatal(err)
			}

			execute := func() error {
				if err := client.SetPrepareTransaction(`INSERT INTO ` + table + `(field) VALUES(` + Prepare[driver] + `);`); err != nil {
					return err
				} else if err := client.ExecutePrepareTransaction(2); err != nil {
					return err
				} else {
					return nil
				}
			}

			errExecute := execute()
			if errExecute != nil {
				t.Error(errExecute)
			}

			if err := client.EndTransaction(errExecute); err != nil {
				t.Fatal(err)
			}
		}()
		wg.Wait()

		count := 0
		if err := client.QueryRow(`SELECT COUNT(*) FROM `+table+`;`, &count); err != nil {
			t.Fatal(err)
		} else if count != 2 {
			t.Fatal("invalid -", count)
		}
	}

	test(t, job)
}

func getOpenInfo(driver sql.Driver, database string) (sql.Driver, string, int) {
	switch driver {
	case sql.DriverAmazonDynamoDB:
		return driver, "Region=dummy;AkId=dummy;SecretKey=dummy;Endpoint=http://127.0.0.1:8000", 10
	case sql.DriverMySQL:
		return driver, "root:root@tcp(127.0.0.1)/" + database, 10
	case sql.DriverPostgres:
		return driver, "host=localhost port=5432 user=postgres password=postgres sslmode=disable dbname=" + database, 10
	case sql.DriverSQLite:
		return driver, SQLiteFileName, 10
	default:
		return driver, "", -1
	}
}

func test(t *testing.T, job func(*testing.T, sql.Client, sql.Driver, string)) {
	for _, driver := range Drivers {
		testDetail(t, driver, job)
	}
}

func testDetail(t *testing.T, driver sql.Driver, job func(*testing.T, sql.Client, sql.Driver, string)) {
	database := strings.ToLower(t.Name())
	table := strings.ToLower(t.Name())

	client := sql.Client{}

	if err := client.Open(getOpenInfo(driver, "")); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	defer func() {
		switch driver {
		case sql.DriverAmazonDynamoDB:
			fallthrough
		case sql.DriverSQLite:
			return
		}

		if err := client.Open(getOpenInfo(driver, "")); err != nil {
			t.Fatal(err)
		} else if err := client.Execute(`DROP DATABASE IF EXISTS ` + database + `;`); err != nil {
			t.Fatal(err)
		}
	}()
	func() {
		switch driver {
		case sql.DriverAmazonDynamoDB:
			fallthrough
		case sql.DriverSQLite:
			return
		}

		if err := client.Execute(`CREATE DATABASE ` + database + `;`); err != nil {
			t.Fatal(err)
		}

		if err := client.Open(getOpenInfo(driver, database)); err != nil {
			t.Fatal(err)
		}
	}()

	defer func() {
		if err := client.Execute(`DROP TABLE IF EXISTS ` + table); err != nil {
			t.Fatal(err)
		}
	}()
	createTableQuery := `CREATE TABLE ` + table
	switch driver {
	case sql.DriverAmazonDynamoDB:
		createTableQuery += ` WITH PK=field:number WITH rcu=3 WITH wcu=5`
	default:
		createTableQuery += `(field int);`
	}
	if err := client.Execute(createTableQuery); err != nil {
		t.Fatal(err)
	}

	insertQuery := `INSERT INTO ` + table
	switch driver {
	case sql.DriverAmazonDynamoDB:
		insertQuery += ` VALUE {'field': 1};`
	default:
		insertQuery += `(field) VALUES(1);`
	}
	if err := client.Execute(insertQuery); err != nil {
		t.Log(insertQuery)
		t.Fatal(err)
	}

	switch driver {
	case sql.DriverAmazonDynamoDB:
	default:
		job(t, client, driver, table)
	}
}
