package sqlx_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/common-library/go/testutil"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
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
var containers = []testcontainers.Container{}

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
	ctx := context.Background()

	setup := func() {
		mysqlContainer, err := mysql.Run(ctx, testutil.MySQLImage,
			mysql.WithDatabase("testdb"),
			mysql.WithUsername("testuser"),
			mysql.WithPassword("testpass"),
		)
		if err != nil {
			panic(err)
		}
		containers = append(containers, mysqlContainer)

		mysqlHost, err := mysqlContainer.Host(ctx)
		if err != nil {
			panic(err)
		}

		mysqlPort, err := mysqlContainer.MappedPort(ctx, "3306")
		if err != nil {
			panic(err)
		}

		mysqlDSN := fmt.Sprintf("testuser:testpass@tcp(%s:%s)/testdb?charset=utf8&parseTime=True&loc=Local", mysqlHost, mysqlPort.Port())

		var db *sqlx.DB
		maxRetries := 10
		for i := 0; i < maxRetries; i++ {
			db, err = sqlx.Connect("mysql", mysqlDSN)
			if err == nil {
				break
			}
			if i < maxRetries-1 {
				backoff := time.Duration(50<<uint(i)) * time.Millisecond
				if backoff > time.Second {
					backoff = time.Second
				}
				time.Sleep(backoff)
			}
		}
		if err != nil {
			panic(err)
		}
		dbs["mysql"] = db

		postgresContainer, err := postgres.Run(ctx, testutil.PostgresImage,
			postgres.WithDatabase("testdb"),
			postgres.WithUsername("testuser"),
			postgres.WithPassword("testpass"),
			postgres.BasicWaitStrategies(),
		)
		if err != nil {
			panic(err)
		}
		containers = append(containers, postgresContainer)

		postgresHost, err := postgresContainer.Host(ctx)
		if err != nil {
			panic(err)
		}

		postgresPort, err := postgresContainer.MappedPort(ctx, "5432")
		if err != nil {
			panic(err)
		}

		postgresDSN := fmt.Sprintf("host=%s user=testuser password=testpass dbname=testdb port=%s sslmode=disable TimeZone=Asia/Seoul", postgresHost, postgresPort.Port())

		for i := 0; i < maxRetries; i++ {
			db, err = sqlx.Connect("postgres", postgresDSN)
			if err == nil {
				break
			}
			if i < maxRetries-1 {
				backoff := time.Duration(50<<uint(i)) * time.Millisecond
				if backoff > time.Second {
					backoff = time.Second
				}
				time.Sleep(backoff)
			}
		}
		if err != nil {
			panic(err)
		}
		dbs["postgres"] = db

		for _, db := range dbs {
			db.MustExec(drop_schema)
			db.MustExec(create_schema)
		}
	}

	teardown := func() {
		for _, db := range dbs {
			db.MustExec(drop_schema)
		}

		for _, container := range containers {
			if err := container.Terminate(ctx); err != nil {
				panic(err)
			}
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
