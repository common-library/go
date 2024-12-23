package postgresql_test

import (
	"context"
	"database/sql"
	"os"
	"strings"
	"testing"

	"github.com/common-library/go/database/orm/sqlc/postgresql/pkg"
	_ "github.com/lib/pq"
)

func getQueries(t *testing.T) (*pkg.Queries, error) {
	if t != nil {
		t.Parallel()
	}

	dsn := strings.Replace(os.Getenv("POSTGRESQL_DSN"), "${database}", "postgres", 1)
	if connection, err := sql.Open("postgres", dsn); err != nil {
		return nil, err
	} else {
		return pkg.New(connection), nil
	}
}

func TestMain(m *testing.M) {
	if len(os.Getenv("POSTGRESQL_DSN")) == 0 {
		return
	}

	ctx := context.Background()

	setup := func() {
		if queries, err := getQueries(nil); err != nil {
			panic(err)
		} else if err := queries.CreateTable01(ctx); err != nil {
			panic(err)
		}
	}

	teardown := func() {
		if queries, err := getQueries(nil); err != nil {
			panic(err)
		} else if err := queries.DropTable01(ctx); err != nil {
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
