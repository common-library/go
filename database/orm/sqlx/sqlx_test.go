package sqlx_test

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var create_schema = `
CREATE TABLE IF NOT EXISTS table01_for_sqlx (
	field01 VARCHAR(512) NOT NULL,
	field02 INTEGER NOT NULL DEFAULT 0,
	field03 BOOLEAN NOT NULL DEFAULT false,
	PRIMARY KEY (field01)
);
`

var drop_schema = `
DROP TABLE IF EXISTS table01_for_sqlx;
`

var dbs = map[string]*sqlx.DB{}

const (
	MySQL      string = "mysql"
	PostgreSQL string = "postgres"
)

type Table01 struct {
	Field01 string
	Field02 int
	Field03 bool
}

func getDbs(t *testing.T) map[string]*sqlx.DB {
	t.Parallel()

	return dbs
}

func TestMain(m *testing.M) {
	setup := func() {
		databaseInfo := map[string]string{}
		if len(os.Getenv("MYSQL_DSN")) != 0 {
			databaseInfo[MySQL] = os.Getenv("MYSQL_DSN") + "mysql"
		}
		if len(os.Getenv("POSTGRESQL_DSN")) != 0 {
			databaseInfo[PostgreSQL] = os.Getenv("POSTGRESQL_DSN") + " dbname=postgres"
		}

		for driverName, dataSource := range databaseInfo {
			if db, err := sqlx.Connect(driverName, dataSource); err != nil {
				panic(err)
			} else {
				dbs[driverName] = db

				db.MustExec(drop_schema)
				db.MustExec(create_schema)
			}
		}
	}

	teardown := func() {
		for _, db := range dbs {
			db.MustExec(drop_schema)
		}
	}

	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestMustExec(t *testing.T) {
	for driverName, db := range getDbs(t) {
		query := "INSERT INTO table01_for_sqlx(field01, field02, field03) VALUES"

		switch driverName {
		case MySQL:
			query += "(?, ?, ?)"
		case PostgreSQL:
			query += "($1, $2, $3)"
		}
		db.MustExec(query, t.Name(), 1, true)

		table01 := Table01{}
		if err := db.Get(&table01, "SELECT * FROM table01_for_sqlx WHERE field01='"+t.Name()+"'"); err != nil {
			t.Fatal(err)
		} else if table01.Field01 != t.Name() || table01.Field02 != 1 || table01.Field03 != true {
			t.Fatal(table01)
		}

		table01s := []Table01{}
		if err := db.Select(&table01s, "SELECT * FROM table01_for_sqlx WHERE field01='"+t.Name()+"'"); err != nil {
			t.Fatal(err)
		} else if len(table01s) != 1 {
			t.Fatal(len(table01s))
		}
	}
}

func TestGet(t *testing.T) {
	TestMustExec(t)
}

func TestSelect(t *testing.T) {
	TestMustExec(t)
}

func TestNamedExec_row(t *testing.T) {
	table01 := Table01{}

	for _, db := range getDbs(t) {
		const query = "INSERT INTO table01_for_sqlx(field01, field02, field03) VALUES (:field01, :field02, :field03)"

		table01 = Table01{Field01: t.Name(), Field02: 1, Field03: true}
		if _, err := db.NamedExec(query, &table01); err != nil {
			t.Fatal(err)
		}
		table01 = Table01{}
		if err := db.Get(&table01, "SELECT * FROM table01_for_sqlx WHERE field01='"+t.Name()+"'"); err != nil {
			t.Fatal(err)
		} else if table01.Field01 != t.Name() || table01.Field02 != 1 || table01.Field03 != true {
			t.Fatal(table01)
		}

		db.MustExec("DELETE FROM table01_for_sqlx WHERE field01='" + t.Name() + "'")

		if _, err := db.NamedExec(query, map[string]any{"field01": t.Name(), "field02": 1, "field03": true}); err != nil {
			t.Fatal(err)
		}
		table01 = Table01{}
		if err := db.Get(&table01, "SELECT * FROM table01_for_sqlx WHERE field01='"+t.Name()+"'"); err != nil {
			t.Fatal(err)
		} else if table01.Field01 != t.Name() || table01.Field02 != 1 || table01.Field03 != true {
			t.Fatal(table01)
		}
	}
}

func TestNamedExec_rows(t *testing.T) {
	table01s := []Table01{
		{Field01: t.Name() + "-1", Field02: 1, Field03: true},
		{Field01: t.Name() + "-2", Field02: 1, Field03: true},
		{Field01: t.Name() + "-3", Field02: 1, Field03: true},
	}

	table01sMap := []map[string]any{
		{"field01": t.Name() + "-1", "field02": 1, "field03": true},
		{"field01": t.Name() + "-2", "field02": 1, "field03": true},
		{"field01": t.Name() + "-3", "field02": 1, "field03": true},
	}

	for _, db := range getDbs(t) {
		const query = "INSERT INTO table01_for_sqlx(field01, field02, field03) VALUES (:field01, :field02, :field03)"

		if _, err := db.NamedExec(query, table01s); err != nil {
			t.Fatal(err)
		}

		for _, table01 := range table01s {
			db.MustExec("DELETE FROM table01_for_sqlx WHERE field01='" + table01.Field01 + "'")
		}

		if _, err := db.NamedExec(query, table01sMap); err != nil {
			t.Fatal(err)
		}
	}
}

func TestQueryx(t *testing.T) {
	table01 := Table01{}

	for _, db := range getDbs(t) {
		const query = "INSERT INTO table01_for_sqlx(field01, field02, field03) VALUES (:field01, :field02, :field03)"

		table01 = Table01{Field01: t.Name(), Field02: 1, Field03: true}
		if _, err := db.NamedExec(query, &table01); err != nil {
			t.Fatal(err)
		}

		table01 = Table01{}
		if rows, err := db.Queryx("SELECT * FROM table01_for_sqlx WHERE field01='" + t.Name() + "'"); err != nil {
			t.Fatal(err)
		} else {
			for rows.Next() {
				if err := rows.StructScan(&table01); err != nil {
					t.Fatal(err)
				} else if table01.Field01 != t.Name() || table01.Field02 != 1 || table01.Field03 != true {
					t.Fatal(table01)
				}
			}
		}
	}
}

func TestNamedQuery(t *testing.T) {
	table01 := Table01{}

	for _, db := range getDbs(t) {
		const query = "INSERT INTO table01_for_sqlx(field01, field02, field03) VALUES (:field01, :field02, :field03)"

		table01 = Table01{Field01: t.Name(), Field02: 1, Field03: true}
		if _, err := db.NamedExec(query, &table01); err != nil {
			t.Fatal(err)
		}

		table01 = Table01{Field01: t.Name()}
		if rows, err := db.NamedQuery("SELECT * FROM table01_for_sqlx WHERE field01=:field01", table01); err != nil {
			t.Fatal(err)
		} else {
			for rows.Next() {
				if err := rows.StructScan(&table01); err != nil {
					t.Fatal(err)
				} else if table01.Field01 != t.Name() || table01.Field02 != 1 || table01.Field03 != true {
					t.Fatal(table01)
				}
			}
		}

		if rows, err := db.NamedQuery("SELECT * FROM table01_for_sqlx WHERE field01=:f1", map[string]interface{}{"f1": t.Name()}); err != nil {
			t.Fatal(err)
		} else {
			for rows.Next() {
				if err := rows.StructScan(&table01); err != nil {
					t.Fatal(err)
				} else if table01.Field01 != t.Name() || table01.Field02 != 1 || table01.Field03 != true {
					t.Fatal(table01)
				}
			}
		}
	}
}

func TestMustBegin(t *testing.T) {
	for _, db := range getDbs(t) {
		tx := db.MustBegin()

		const query = "INSERT INTO table01_for_sqlx(field01, field02, field03) VALUES (:field01, :field02, :field03)"

		table01 := Table01{Field01: t.Name(), Field02: 1, Field03: true}
		if _, err := tx.NamedExec(query, &table01); err != nil {
			t.Fatal(err)
		}

		table01 = Table01{}
		if err := db.Get(&table01, "SELECT * FROM table01_for_sqlx WHERE field01='"+t.Name()+"'"); err != sql.ErrNoRows {
			t.Fatal(err)
		}

		tx.Commit()

		table01 = Table01{}
		if err := db.Get(&table01, "SELECT * FROM table01_for_sqlx WHERE field01='"+t.Name()+"'"); err != nil {
			t.Fatal(err)
		} else if table01.Field01 != t.Name() || table01.Field02 != 1 || table01.Field03 != true {
			t.Fatal(table01)
		}

	}
}
