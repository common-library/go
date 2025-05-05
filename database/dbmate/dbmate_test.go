package main

import (
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/clickhouse"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/mysql"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
)

func getDb(t *testing.T, rawURL string) *dbmate.DB {
	t.Parallel()

	u, err := url.Parse(rawURL)
	if err != nil {
		t.Fatal(err)
	}

	return dbmate.New(u)
}

func test(t *testing.T, url, migrationsDir string) {
	if len(url) == 0 {
		return
	}

	url = strings.Replace(url, "${database}", "dbmate_test", 1)
	db := getDb(t, url)
	db.AutoDumpSchema = false
	db.MigrationsDir = []string{migrationsDir}

	query := func(query string) error {
		if driver, err := db.Driver(); err != nil {
			return err
		} else if sqlDB, err := driver.Open(); err != nil {
			return err
		} else if rows, err := sqlDB.Query(query); err != nil {
			return err
		} else if err := rows.Close(); err != nil {
			return err
		} else if err := sqlDB.Close(); err != nil {
			return err
		} else {
			return nil
		}
	}

	validation := func() error {
		if err := query("SELECT field03 FROM test_01;"); err != nil {
			return err
		} else if err := db.Rollback(); err != nil {
			return err
		} else if err := query("SELECT field02 FROM test_01;"); err != nil {
			return err
		} else {
			return nil
		}
	}

	defer func() {
		if err := db.Drop(); err != nil {
			t.Fatal(err)
		}
	}()

	if err := db.Create(); err != nil {
		t.Fatal(err)
	} else if err := db.Migrate(); err != nil {
		t.Fatal(err)
	} else if err := validation(); err != nil {
		t.Fatal(err)
	} else if err := db.Drop(); err != nil {
		t.Fatal(err)
	}

	if err := db.CreateAndMigrate(); err != nil {
		t.Fatal(err)
	} else if err := validation(); err != nil {
		t.Fatal(err)
	}
}

func TestClickHouse(t *testing.T) {
	test(t, os.Getenv("CLICKHOUSE_URL"), "./clickhouse/migrations")
}

func TestMySQL(t *testing.T) {
	test(t, os.Getenv("MYSQL_URL"), "./mysql/migrations")
}

func TestPostgresql(t *testing.T) {
	test(t, os.Getenv("POSTGRESQL_URL"), "./postgresql/migrations")
}
