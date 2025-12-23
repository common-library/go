package dbmate_test

import (
	"context"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/clickhouse"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/mysql"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	"github.com/common-library/go/testutil"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/wait"
)

func getDb(t *testing.T, rawURL string) *dbmate.DB {
	u, err := url.Parse(rawURL)
	if err != nil {
		t.Fatal(err)
	}

	return dbmate.New(u)
}

func test(t *testing.T, dbURL, migrationsDir string) {
	db := getDb(t, dbURL)
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
	t.Parallel()

	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        testutil.ClickHouseImage,
		ExposedPorts: []string{"8123/tcp", "9000/tcp"},
		Env: map[string]string{
			"CLICKHOUSE_DB":                        "default",
			"CLICKHOUSE_USER":                      "default",
			"CLICKHOUSE_DEFAULT_ACCESS_MANAGEMENT": "1",
		},
		WaitingFor: wait.ForHTTP("/ping").WithPort("8123").WithStartupTimeout(60 * time.Second),
	}

	clickhouseContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer clickhouseContainer.Terminate(ctx)

	host, err := clickhouseContainer.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}

	port, err := clickhouseContainer.MappedPort(ctx, "9000")
	if err != nil {
		t.Fatal(err)
	}

	dbURL := fmt.Sprintf("clickhouse://%s:%s/testdb", host, port.Port())

	test(t, dbURL, "./clickhouse/migrations")
}

func TestMySQL(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mysqlContainer, err := mysql.Run(ctx,
		testutil.MySQLImage,
		mysql.WithDatabase("tmp_db"),
		mysql.WithUsername("root"),
		mysql.WithPassword("password"),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer mysqlContainer.Terminate(ctx)

	host, err := mysqlContainer.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}

	port, err := mysqlContainer.MappedPort(ctx, "3306")
	if err != nil {
		t.Fatal(err)
	}

	dbURL := fmt.Sprintf("mysql://root:password@%s:%s/testdb?parseTime=true", host, port.Port())

	test(t, dbURL, "./mysql/migrations")
}

func TestPostgresql(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        testutil.PostgresImage,
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "password",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(60 * time.Second),
	}

	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer postgresContainer.Terminate(ctx)

	host, err := postgresContainer.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}

	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatal(err)
	}

	dbURL := fmt.Sprintf("postgres://postgres:password@%s:%s/testdb?sslmode=disable", host, port.Port())

	test(t, dbURL, "./postgresql/migrations")
}
