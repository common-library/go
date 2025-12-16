package mysql_test

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"os"
	"testing"
	"time"

	"github.com/common-library/go/database/orm/sqlc/mysql/pkg"
	"github.com/common-library/go/testutil"
	_ "github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	container testcontainers.Container
	dsn       string
)

func getQueries(t *testing.T) (*pkg.Queries, func(), error) {
	if t != nil {
		t.Parallel()
	}

	connection, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		connection.Close()
	}

	return pkg.New(connection), cleanup, nil
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	setup := func() {
		var err error
		container, err = mysql.Run(ctx, testutil.MySQLImage,
			mysql.WithDatabase("testdb"),
			mysql.WithUsername("testuser"),
			mysql.WithPassword("testpass"),
			testcontainers.WithWaitStrategy(
				wait.ForLog("ready for connections").
					WithOccurrence(2).
					WithStartupTimeout(2*time.Minute)),
		)
		if err != nil {
			panic(err)
		}

		host, err := container.Host(ctx)
		if err != nil {
			panic(err)
		}

		port, err := container.MappedPort(ctx, "3306")
		if err != nil {
			panic(err)
		}

		dsn = fmt.Sprintf("testuser:testpass@tcp(%s:%s)/testdb?charset=utf8&parseTime=True&loc=Local", host, port.Port())

		maxRetries := 20
		for i := 0; i < maxRetries; i++ {
			var queries *pkg.Queries
			var cleanup func()
			queries, cleanup, err = getQueries(nil)
			if err == nil {
				err = queries.CreateTable01(ctx)
				if err == nil {
					cleanup()
					break
				}
				cleanup()
			} else if cleanup != nil {
				cleanup()
			}

			if i < maxRetries-1 {
				baseMs := 100 << uint(i)
				backoffMs := math.Min(float64(baseMs), 2000)
				backoff := time.Duration(backoffMs) * time.Millisecond
				time.Sleep(backoff)
			}
		}
		if err != nil {
			panic(fmt.Sprintf("Failed to setup MySQL after %d retries: %v", maxRetries, err))
		}
	}

	teardown := func() {

		queries, cleanup, err := getQueries(nil)
		if err == nil && queries != nil {
			defer cleanup()
			_ = queries.DropTable01(ctx)
		} else if cleanup != nil {
			cleanup()
		}

		if container != nil {
			if err := container.Terminate(ctx); err != nil {

				fmt.Fprintf(os.Stderr, "Failed to terminate container: %v\n", err)
			}
		}
	}

	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestGetTable01(t *testing.T) {
	queries, cleanup, err := getQueries(t)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	ctx := context.Background()

	insertTable01Params := pkg.InsertTable01Params{
		Field01: t.Name(),
		Field02: 1,
	}
	if result, err := queries.InsertTable01(ctx, insertTable01Params); err != nil {
		t.Fatal(err)
	} else if rowsAffected, err := result.RowsAffected(); err != nil {
		t.Fatal(err)
	} else if rowsAffected != 1 {
		t.Fatal(rowsAffected)
	}

	if table01, err := queries.GetTable01(ctx, t.Name()); err != nil {
		t.Fatal(err)
	} else if table01.Field01 != t.Name() || table01.Field02 != 1 {
		t.Fatal(table01)
	}
}

func TestListTable01(t *testing.T) {
	queries, cleanup, err := getQueries(t)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	ctx := context.Background()

	insertTable01Params := pkg.InsertTable01Params{
		Field01: t.Name(),
		Field02: 1,
	}
	if _, err := queries.InsertTable01(ctx, insertTable01Params); err != nil {
		t.Fatal(err)
	}

	if table01s, err := queries.ListTable01(ctx); err != nil {
		t.Fatal(err)
	} else {
		for _, table01 := range table01s {
			t.Log(table01)
		}
	}
}

func TestInsertTable01(t *testing.T) {
	TestGetTable01(t)
}

func TestUpdateTable01(t *testing.T) {
	queries, cleanup, err := getQueries(t)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	ctx := context.Background()

	insertTable01Params := pkg.InsertTable01Params{
		Field01: t.Name(),
		Field02: 1,
	}
	if _, err := queries.InsertTable01(ctx, insertTable01Params); err != nil {
		t.Fatal(err)
	}

	updateTable01Params := pkg.UpdateTable01Params{
		Field01: t.Name(),
		Field02: 2,
	}
	if err := queries.UpdateTable01(ctx, updateTable01Params); err != nil {
		t.Fatal(err)
	}

	if table01, err := queries.GetTable01(ctx, t.Name()); err != nil {
		t.Fatal(err)
	} else if table01.Field01 != t.Name() || table01.Field02 != 2 {
		t.Fatal(table01)
	}
}

func TestDeleteTable01(t *testing.T) {
	queries, cleanup, err := getQueries(t)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	ctx := context.Background()

	insertTable01Params := pkg.InsertTable01Params{
		Field01: t.Name(),
		Field02: 1,
	}
	if _, err := queries.InsertTable01(ctx, insertTable01Params); err != nil {
		t.Fatal(err)
	}

	if err := queries.DeleteTable01(ctx, t.Name()); err != nil {
		t.Fatal(err)
	}

	if _, err := queries.GetTable01(ctx, t.Name()); err != sql.ErrNoRows {
		t.Fatal(err)
	}
}
