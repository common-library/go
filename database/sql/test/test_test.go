package test

import (
	"errors"
	"strings"
	"testing"

	"github.com/common-library/go/database/sql"
	"github.com/common-library/go/file"
)

var databases = map[sql.Driver]database{
	sql.DriverAmazonDynamoDB: &dynamodb{},
	sql.DriverMySQL:          &mysql{},
	sql.DriverPostgreSQL:     &postgresql{},
	sql.DriverSQLite:         &sqlite{},
}

type database interface {
	createDatabase() bool

	getClient(t *testing.T, databaseName string) (sql.Client, bool)
	getDriver() sql.Driver
	getPrepare(index uint64) string
	getCreateTableQuery(tableName string) []string
}

func TestOpen(t *testing.T) {
	for _, db := range databases {
		client, ok := db.getClient(t, "")
		if ok == false {
			continue
		}

		if err := client.Close(); err != nil {
			t.Fatal(err)
		}
	}
}

func TestClose(t *testing.T) {
	for _, db := range databases {
		client, ok := db.getClient(t, "")
		if ok == false {
			continue
		}

		if err := client.Close(); err != nil {
			t.Fatal(err)
		}
	}
}

func TestQuery(t *testing.T) {
	if _, err := (&sql.Client{}).Query(""); err.Error() != `please call Open first` {
		t.Fatal(err)
	}

	job := func(db database, client sql.Client, tableName string) {
		query := `SELECT field02 FROM ` + tableName + ` WHERE field01=1`
		if rows, err := client.Query(query); err != nil {
			t.Fatal(err)
		} else {
			defer rows.Close()

			for rows.Next() {
				field02 := ""
				if err := rows.Scan(&field02); err != nil {
					t.Fatal(err)
				}

				if field02 != "abc" {
					t.Fatal(field02)
				}
			}
		}
	}

	test(t, job)
}

func TestQueryRow(t *testing.T) {
	if err := (&sql.Client{}).QueryRow(``); err.Error() != `please call Open first` {
		t.Fatal(err)
	}

	job := func(db database, client sql.Client, tableName string) {
		field01 := 0
		field02 := ""

		if err := client.QueryRow(`SELECT field01, field02 FROM `+tableName, &field01, &field02); err != nil {
			t.Fatal(err)
		} else if field01 != 1 || field02 != "abc" {
			t.Log(field01)
			t.Log(field02)
			t.Fatal("invalid")
		}
	}

	test(t, job)
}

func TestExecute(t *testing.T) {
	if err := (&sql.Client{}).Execute(""); err.Error() != `please call Open first` {
		t.Fatal(err)
	}

	job := func(db database, client sql.Client, tableName string) {
		if err := client.Execute(`UPDATE ` + tableName + ` SET field02='123' WHERE field01=1`); err != nil {
			t.Fatal(err)
		}

		field02 := ""
		if err := client.QueryRow(`SELECT field02 FROM `+tableName, &field02); err != nil {
			t.Fatal(err)
		} else if field02 != "123" {
			t.Fatal(field02)
		}
	}

	test(t, job)
}

func TestSetPrepare(t *testing.T) {
	if err := (&sql.Client{}).SetPrepare(""); err.Error() != `please call Open first` {
		t.Fatal(err)
	}

	job := func(db database, client sql.Client, tableName string) {
		if len(db.getPrepare(1)) == 0 {
			return
		}

		if err := client.SetPrepare(`SELECT field02 FROM ` + tableName + ` WHERE field01=` + db.getPrepare(1)); err != nil {
			t.Log(`SELECT field02 FROM ` + tableName + ` WHERE field01=` + db.getPrepare(1))
			t.Fatal(err)
		} else if row, err := client.QueryRowPrepare(1); err != nil {
			t.Fatal(err)
		} else {
			field02 := ""
			if err := row.Scan(&field02); err != nil {
				t.Fatal(err)
			} else if field02 != "abc" {
				t.Fatal(field02)
			}
		}
	}

	test(t, job)
}

func TestQueryPrepare(t *testing.T) {
	if _, err := (&sql.Client{}).QueryPrepare(1); err.Error() != `please call SetPrepare first` {
		t.Fatal(err)
	}

	job := func(db database, client sql.Client, tableName string) {
		if len(db.getPrepare(1)) == 0 {
			return
		}

		if err := client.SetPrepare(`SELECT field02 FROM ` + tableName + ` WHERE field01=` + db.getPrepare(1)); err != nil {
			t.Fatal(err)
		} else if rows, err := client.QueryPrepare(1); err != nil {
			t.Fatal(err)
		} else {
			defer rows.Close()

			for rows.Next() {
				field02 := ""
				if err := rows.Scan(&field02); err != nil {
					t.Fatal(err)
				} else if field02 != "abc" {
					t.Fatal(field02)
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

	job := func(db database, client sql.Client, tableName string) {
		if len(db.getPrepare(1)) == 0 {
			return
		}

		if err := client.SetPrepare(`SELECT field02 FROM ` + tableName + ` WHERE field01=` + db.getPrepare(1)); err != nil {
			t.Fatal(err)
		} else if row, err := client.QueryRowPrepare(1); err != nil {
			t.Fatal(err)
		} else {
			field02 := ""
			if err := row.Scan(&field02); err != nil {
				t.Fatal(err)
			} else if field02 != "abc" {
				t.Fatal(field02)
			}
		}
	}

	test(t, job)
}

func TestExecutePrepare(t *testing.T) {
	if err := (&sql.Client{}).ExecutePrepare(2); err.Error() != `please call SetPrepare first` {
		t.Fatal(err)
	}

	job := func(db database, client sql.Client, tableName string) {
		if len(db.getPrepare(1)) == 0 {
			return
		}

		if err := client.SetPrepare(`UPDATE ` + tableName + ` SET field02=` + db.getPrepare(1) + ` WHERE field01=` + db.getPrepare(2)); err != nil {
			t.Fatal(err)
		} else if err := client.ExecutePrepare("123", 1); err != nil {
			t.Fatal(err)
		}

		field02 := ""
		if err := client.QueryRow(`SELECT field02 FROM `+tableName, &field02); err != nil {
			t.Fatal(err)
		} else if field02 != "123" {
			t.Fatal(field02)
		}
	}

	test(t, job)
}

func TestBeginTransaction(t *testing.T) {
	if err := (&sql.Client{}).BeginTransaction(); err.Error() != `please call Open first` {
		t.Fatal(err)
	}

	job := func(db database, client sql.Client, tableName string) {
		switch client.GetDriver() {
		case sql.DriverAmazonDynamoDB:
			return
		}

		check := func(compare string) error {
			field02 := ""
			if err := client.QueryRow(`SELECT field02 FROM `+tableName, &field02); err != nil {
				return err
			} else if field02 != compare {
				t.Log(field02)
				t.Log(compare)
				return errors.New("invalid")
			} else {
				return nil
			}
		}

		if err := client.BeginTransaction(); err != nil {
			t.Fatal(err)
		}

		errExecute := client.ExecuteTransaction(`UPDATE ` + tableName + ` SET field02='123' WHERE field01=1`)
		if errExecute != nil {
			t.Fatal(errExecute)
		}

		if err := check("abc"); err != nil {
			t.Fatal(err)
		}

		if err := client.EndTransaction(errExecute); err != nil {
			t.Fatal(err)
		}

		if err := check("123"); err != nil {
			t.Fatal(err)
		}
	}

	test(t, job)
}

func TestEndTransaction(t *testing.T) {
	if err := (&sql.Client{}).EndTransaction(nil); err.Error() != `please call BeginTransaction first` {
		t.Fatal(err)
	}

	TestBeginTransaction(t)
}

func TestQueryTransaction(t *testing.T) {
	if _, err := (&sql.Client{}).QueryTransaction(``); err.Error() != `please call BeginTransaction first` {
		t.Fatal(err)
	}

	job := func(db database, client sql.Client, tableName string) {
		switch client.GetDriver() {
		case sql.DriverAmazonDynamoDB:
			return
		}

		if err := client.BeginTransaction(); err != nil {
			t.Fatal(err)
		}

		errQuery := func() error {
			if rows, err := client.QueryTransaction(`SELECT field02 FROM ` + tableName); err != nil {
				t.Fatal(err)
			} else {
				defer rows.Close()

				for rows.Next() {
					field02 := ""
					if err := rows.Scan(&field02); err != nil {
						t.Fatal(err)
					} else if field02 != "abc" {
						t.Fatal(field02)
					}
				}
			}

			return nil
		}()

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

	job := func(db database, client sql.Client, tableName string) {
		switch client.GetDriver() {
		case sql.DriverAmazonDynamoDB:
			return
		}

		if err := client.BeginTransaction(); err != nil {
			t.Fatal(err)
		}

		errQuery := func() error {
			field02 := ""
			if err := client.QueryRowTransaction(`SELECT field02 FROM `+tableName+` WHERE field01=1`, &field02); err != nil {
				t.Fatal(err)
			} else if field02 != "abc" {
				t.Fatal(field02)
			}

			return nil
		}()

		if err := client.EndTransaction(errQuery); err != nil {
			t.Fatal(err)
		}
	}

	test(t, job)
}

func TestExecuteTransaction(t *testing.T) {
	if err := (&sql.Client{}).ExecuteTransaction(``); err.Error() != `please call BeginTransaction first` {
		t.Fatal(err)
	}

	if err := (&sql.Client{}).BeginTransaction(); err.Error() != `please call Open first` {
		t.Fatal(err)
	}

	job := func(db database, client sql.Client, tableName string) {
		switch client.GetDriver() {
		case sql.DriverAmazonDynamoDB:
			return
		}

		check := func(compare string) error {
			field02 := ""
			if err := client.QueryRow(`SELECT field02 FROM `+tableName, &field02); err != nil {
				return err
			} else if field02 != compare {
				t.Log(field02)
				t.Log(compare)
				return errors.New("invalid")
			} else {
				return nil
			}
		}

		if err := client.BeginTransaction(); err != nil {
			t.Fatal(err)
		}

		errExecute := client.ExecuteTransaction(`UPDATE ` + tableName + ` SET field02='123' WHERE field01=1`)
		if errExecute != nil {
			t.Fatal(errExecute)
		}

		if err := check("abc"); err != nil {
			t.Fatal(err)
		}

		if err := client.EndTransaction(errExecute); err != nil {
			t.Fatal(err)
		}

		if err := check("123"); err != nil {
			t.Fatal(err)
		}
	}

	test(t, job)
}

func TestSetPrepareTransaction(t *testing.T) {
	if err := (&sql.Client{}).SetPrepareTransaction(``); err.Error() != `please call BeginTransaction first` {
		t.Fatal(err)
	}

	job := func(db database, client sql.Client, tableName string) {
		switch client.GetDriver() {
		case sql.DriverAmazonDynamoDB:
			return
		}

		if err := client.BeginTransaction(); err != nil {
			t.Fatal(err)
		}

		errExecute := func() error {
			if err := client.SetPrepareTransaction(`SELECT field02 FROM ` + tableName + ` WHERE field01=` + db.getPrepare(1)); err != nil {
				t.Fatal(err)
			} else if row, err := client.QueryRowPrepareTransaction(1); err != nil {
				t.Fatal(err)
			} else {
				field02 := ""
				if err := row.Scan(&field02); err != nil {
					t.Fatal(err)
				} else if field02 != "abc" {
					t.Fatal(field02)
				}
			}

			return nil
		}()

		if err := client.EndTransaction(errExecute); err != nil {
			t.Fatal(err)
		}
	}

	test(t, job)
}

func TestQueryPrepareTransaction(t *testing.T) {
	if _, err := (&sql.Client{}).QueryPrepareTransaction(``); err.Error() != `please call SetPrepareTransaction first` {
		t.Fatal(err)
	}

	job := func(db database, client sql.Client, tableName string) {
		switch client.GetDriver() {
		case sql.DriverAmazonDynamoDB:
			return
		}

		if err := client.BeginTransaction(); err != nil {
			t.Fatal(err)
		}

		errQuery := func() error {
			if err := client.SetPrepareTransaction(`SELECT field02 FROM ` + tableName + ` WHERE field01=` + db.getPrepare(1)); err != nil {
				t.Fatal(err)
			} else if rows, err := client.QueryPrepareTransaction(1); err != nil {
				t.Fatal(err)
			} else {
				defer rows.Close()

				for rows.Next() {
					field02 := ""
					if err := rows.Scan(&field02); err != nil {
						t.Fatal(err)
					} else if field02 != "abc" {
						t.Fatal(field02)
					}
				}
			}

			return nil
		}()

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

	job := func(db database, client sql.Client, tableName string) {
		switch client.GetDriver() {
		case sql.DriverAmazonDynamoDB:
			return
		}

		if err := client.BeginTransaction(); err != nil {
			t.Fatal(err)
		}

		errQuery := func() error {
			if err := client.SetPrepareTransaction(`SELECT field02 FROM ` + tableName + ` WHERE field01=` + db.getPrepare(1)); err != nil {
				t.Fatal(err)
			} else if row, err := client.QueryRowPrepareTransaction(1); err != nil {
				t.Fatal(err)
			} else {
				field02 := ""
				if err := row.Scan(&field02); err != nil {
					t.Fatal(err)
				} else if field02 != "abc" {
					t.Fatal(field02)
				}
			}

			return nil
		}()

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

	job := func(db database, client sql.Client, tableName string) {
		switch client.GetDriver() {
		case sql.DriverAmazonDynamoDB:
			return
		}

		check := func(compare string) error {
			field02 := ""
			if err := client.QueryRow(`SELECT field02 FROM `+tableName, &field02); err != nil {
				return err
			} else if field02 != compare {
				t.Log(field02)
				t.Log(compare)
				return errors.New("invalid")
			} else {
				return nil
			}
		}

		if err := client.BeginTransaction(); err != nil {
			t.Fatal(err)
		}

		errExecute := func() error {
			if err := client.SetPrepareTransaction(`UPDATE ` + tableName + ` SET field02=` + db.getPrepare(1) + ` WHERE field01=` + db.getPrepare(2)); err != nil {
				t.Fatal(err)
			} else if err := client.ExecutePrepareTransaction("123", 1); err != nil {
				t.Fatal(err)
			}

			return nil
		}()

		if err := check("abc"); err != nil {
			t.Fatal(err)
		}

		if err := client.EndTransaction(errExecute); err != nil {
			t.Fatal(err)
		}

		if err := check("123"); err != nil {
			t.Fatal(err)
		}
	}

	test(t, job)
}

func TestGetDriver(t *testing.T) {
	for driver, db := range databases {
		client, ok := db.getClient(t, "")
		if ok == false {
			continue
		}

		if driver != client.GetDriver() {
			t.Log(driver)
			t.Log(client.GetDriver())
			t.Fatal("invalid")
		}

		if err := client.Close(); err != nil {
			t.Fatal(err)
		}
	}
}

func test(t *testing.T, job func(db database, client sql.Client, tableName string)) {
	finalJob := func(db database) {
		databaseName := strings.ToLower(t.Name())
		tableName := strings.ToLower(t.Name())

		client, ok := db.getClient(t, "")
		if ok == false {
			return
		}
		if db.createDatabase() {
			if err := client.Execute(`DROP DATABASE IF EXISTS ` + databaseName); err != nil {
				t.Fatal(err)
			}

			if err := client.Execute(`CREATE DATABASE ` + databaseName); err != nil {
				t.Fatal(err)
			}
		}
		if err := client.Close(); err != nil {
			t.Fatal(err)
		}

		client, _ = db.getClient(t, databaseName)
		defer func() {
			if err := client.Close(); err != nil {
				t.Fatal(err)
			}

			client, _ := db.getClient(t, "")
			defer func() {
				if err := client.Close(); err != nil {
					t.Fatal(err)
				}
			}()

			if db.createDatabase() == false {
				return
			} else if err := client.Execute(`DROP DATABASE IF EXISTS ` + databaseName); err != nil {
				t.Fatal(err)
			}
		}()

		defer func() {
			if err := client.Execute(`DROP TABLE IF EXISTS ` + tableName); err != nil {
				t.Fatal(err)
			}

			switch client.GetDriver() {
			case sql.DriverSQLite:
				file.RemoveAll(tableName)
			}
		}()

		for _, query := range db.getCreateTableQuery(tableName) {
			if err := client.Execute(query); err != nil {
				t.Fatal(err)
			}
		}

		insertQuery := `INSERT INTO ` + tableName
		switch db.getDriver() {
		case sql.DriverAmazonDynamoDB:
			insertQuery += ` VALUE {'field01' : 1, 'field02' : 'abc'};`
		default:
			insertQuery += `(field01, field02) VALUES(1, 'abc');`
		}

		if err := client.Execute(insertQuery); err != nil {
			t.Fatal(err)
		}

		job(db, client, tableName)
	}

	for _, db := range databases {
		finalJob(db)
	}
}
