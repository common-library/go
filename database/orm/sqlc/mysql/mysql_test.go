package mysql_test

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/common-library/go/database/orm/sqlc/mysql/pkg"
	_ "github.com/go-sql-driver/mysql"
)

func getQueries(t *testing.T) (*pkg.Queries, error) {
	if t != nil {
		t.Parallel()
	}

	dsn := os.Getenv("MYSQL_DSN") + "mysql"
	if connection, err := sql.Open("mysql", dsn); err != nil {
		return nil, err
	} else {
		return pkg.New(connection), nil
	}
}

func TestMain(m *testing.M) {
	if len(os.Getenv("MYSQL_DSN")) == 0 {
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
	queries, err := getQueries(t)
	if err != nil {
		t.Fatal(err)
	}

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
	queries, err := getQueries(t)
	if err != nil {
		t.Fatal(err)
	}

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
	queries, err := getQueries(t)
	if err != nil {
		t.Fatal(err)
	}

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
