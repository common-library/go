package postgresql_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/common-library/go/database/orm/sqlc/postgresql/pkg"
	"github.com/common-library/go/testutil"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var container testcontainers.Container
var dsn string

func getQueries(t *testing.T) (*pkg.Queries, error) {
	if t != nil {
		t.Parallel()
	}

	if connection, err := sql.Open("postgres", dsn); err != nil {
		return nil, err
	} else {
		return pkg.New(connection), nil
	}
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	setup := func() {
		var err error
		container, err = postgres.Run(ctx, testutil.PostgresImage,
			postgres.WithDatabase("testdb"),
			postgres.WithUsername("testuser"),
			postgres.WithPassword("testpass"),
			postgres.BasicWaitStrategies(),
		)
		if err != nil {
			panic(err)
		}

		host, err := container.Host(ctx)
		if err != nil {
			panic(err)
		}

		port, err := container.MappedPort(ctx, "5432")
		if err != nil {
			panic(err)
		}

		dsn = fmt.Sprintf("host=%s user=testuser password=testpass dbname=testdb port=%s sslmode=disable TimeZone=Asia/Seoul", host, port.Port())

		maxRetries := 10
		for i := 0; i < maxRetries; i++ {
			var queries *pkg.Queries
			queries, err = getQueries(nil)
			if err == nil {
				err = queries.CreateTable01(ctx)
				if err == nil {
					break
				}
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
	}

	teardown := func() {
		if queries, err := getQueries(nil); err != nil {
			panic(err)
		} else if err := queries.DropTable01(ctx); err != nil {
			panic(err)
		}

		if err := container.Terminate(ctx); err != nil {
			panic(err)
		}
	}

	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestGetTable01(t *testing.T) {
	queries, err := getQueries(t)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	insertTable01Params := pkg.InsertTable01Params{
		Field01: t.Name(),
		Field02: 1,
	}
	if table01, err := queries.InsertTable01(ctx, insertTable01Params); err != nil {
		t.Fatal(err)
	} else if table01.Field01 != t.Name() || table01.Field02 != 1 {
		t.Fatal(table01)
	}

	if table01, err := queries.GetTable01(ctx, t.Name()); err != nil {
		t.Fatal(err)
	} else if table01.Field01 != t.Name() || table01.Field02 != 1 {
		t.Fatal(table01)
	}
}

func TestListTable01(t *testing.T) {
	queries, err := getQueries(t)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	insertTable01Params := pkg.InsertTable01Params{
		Field01: t.Name(),
		Field02: 1,
	}
	if table01, err := queries.InsertTable01(ctx, insertTable01Params); err != nil {
		t.Fatal(err)
	} else if table01.Field01 != t.Name() || table01.Field02 != 1 {
		t.Fatal(table01)
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
	queries, err := getQueries(t)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	insertTable01Params := pkg.InsertTable01Params{
		Field01: t.Name(),
		Field02: 1,
	}
	if table01, err := queries.InsertTable01(ctx, insertTable01Params); err != nil {
		t.Fatal(err)
	} else if table01.Field01 != t.Name() || table01.Field02 != 1 {
		t.Fatal(table01)
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
	queries, err := getQueries(t)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	insertTable01Params := pkg.InsertTable01Params{
		Field01: t.Name(),
		Field02: 1,
	}
	if table01, err := queries.InsertTable01(ctx, insertTable01Params); err != nil {
		t.Fatal(err)
	} else if table01.Field01 != t.Name() || table01.Field02 != 1 {
		t.Fatal(table01)
	}

	if err := queries.DeleteTable01(ctx, t.Name()); err != nil {
		t.Fatal(err)
	}

	if _, err := queries.GetTable01(ctx, t.Name()); err != sql.ErrNoRows {
		t.Fatal(err)
	}
}
