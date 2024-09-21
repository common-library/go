package ent_test

import (
	"context"
	"database/sql"
	"os"
	"strconv"
	"sync"
	"testing"

	ent_sql "entgo.io/ent/dialect/sql"
	"github.com/common-library/go/database/orm/ent"
	"github.com/common-library/go/database/orm/ent/table01forent"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var clients = []*ent.Client{}

var databaseInfo = map[string]string{}

func prepareData(t *testing.T, ctx context.Context) {
	for _, client := range clients {
		if table01, err := client.Table01ForEnt.Create().SetField01(t.Name()).SetField02(1).SetField03(true).Save(ctx); err != nil {
			t.Fatal(err)
		} else if table01.Field01 != t.Name() || table01.Field02 != 1 || table01.Field03 != true {
			t.Fatal(table01)
		}

		table01ForEntCreates := []*ent.Table01ForEntCreate{
			client.Table01ForEnt.Create().SetField01(t.Name() + "-1").SetField02(1).SetField05(t.Name()),
			client.Table01ForEnt.Create().SetField01(t.Name() + "-2").SetField02(2).SetField05(t.Name()),
			client.Table01ForEnt.Create().SetField01(t.Name() + "-3").SetField02(3).SetField05(t.Name()),
		}
		if err := client.Table01ForEnt.CreateBulk(table01ForEntCreates...).
			OnConflictColumns(table01forent.FieldField01).
			UpdateNewValues().
			Exec(ctx); err != nil {
			t.Fatal(err)
		}
	}

	t.Parallel()
}

func dropTables() error {
	tables := []string{
		"table01for_ents",

		"repository_for_ent_user_for_ents",
		"issue_for_ents",
		"user_for_ents",
		"repository_for_ents",
	}

	for _, table := range tables {
		for driverName, dataSourceName := range databaseInfo {
			if db, err := sql.Open(driverName, dataSourceName); err != nil {
				return err
			} else if _, err := db.Exec("DROP TABLE IF EXISTS " + table); err != nil {
				return err
			} else if err := db.Close(); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestMain(m *testing.M) {
	setup := func() {
		ctx := context.Background()

		if len(os.Getenv("MYSQL_DSN")) != 0 {
			databaseInfo["mysql"] = os.Getenv("MYSQL_DSN") + "mysql?parseTime=true"
		}
		if len(os.Getenv("POSTGRESQL_DSN")) != 0 {
			databaseInfo["postgres"] = os.Getenv("POSTGRESQL_DSN") + " dbname=postgres"
		}

		for driverName, dataSourceName := range databaseInfo {
			if client, err := ent.Open(driverName, dataSourceName); err != nil {
				panic(err)
			} else if err := client.Schema.Create(ctx); err != nil {
				panic(err)
			} else if _, err := client.Table01ForEnt.Delete().Exec(ctx); err != nil {
				panic(err)
			} else {
				clients = append(clients, client)
			}
		}
	}

	teardown := func() {
		if err := dropTables(); err != nil {
			panic(err)
		}

		for _, client := range clients {
			if err := client.Close(); err != nil {
				panic(err)
			}
		}
	}

	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestCreate(t *testing.T) {
	ctx := context.Background()

	prepareData(t, ctx)

	wg := new(sync.WaitGroup)
	defer wg.Wait()
	for _, client := range clients {
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func(client *ent.Client, i int) {
				defer wg.Done()

				field01 := t.Name() + strconv.Itoa(i)

				if table01, err := client.Table01ForEnt.Create().SetField01(field01).SetField02(1).Save(ctx); err != nil {
					t.Fatal(err)
				} else if table01.Field01 != field01 || table01.Field02 != 1 {
					t.Fatal(table01)
				}
			}(client, i)
		}
	}
}

func TestCreateBulk(t *testing.T) {
	ctx := context.Background()

	prepareData(t, ctx)

	for _, client := range clients {
		table01ForEntCreates := []*ent.Table01ForEntCreate{
			client.Table01ForEnt.Create().SetField01(t.Name() + "-5").SetField02(1),
			client.Table01ForEnt.Create().SetField01(t.Name() + "-6").SetField02(1),
			client.Table01ForEnt.Create().SetField01(t.Name() + "-7").SetField02(1),
		}

		if results, err := client.Table01ForEnt.CreateBulk(table01ForEntCreates...).Save(ctx); err != nil {
			t.Fatal(err)
		} else if len(results) != 3 {
			t.Fatal(results)
		}
	}
}

func TestMapCreateBulk(t *testing.T) {
	ctx := context.Background()

	prepareData(t, ctx)

	field01s := []string{t.Name() + "-5", t.Name() + "-6", t.Name() + "-7"}

	for _, client := range clients {
		if results, err := client.Table01ForEnt.MapCreateBulk(field01s, func(table01ForEntCreate *ent.Table01ForEntCreate, index int) {
			table01ForEntCreate.SetField01(field01s[index]).SetField02(1)
		}).Save(ctx); err != nil {
			t.Fatal(err)
		} else if len(results) != 3 {
			t.Fatal(results)
		}
	}
}

func TestUpdateOne(t *testing.T) {
	ctx := context.Background()

	prepareData(t, ctx)

	for _, client := range clients {
		if table01, err := client.Table01ForEnt.Query().Where(table01forent.Field01(t.Name())).Only(ctx); err != nil {
			t.Fatal(err)
		} else if table01, err := client.Table01ForEnt.UpdateOne(table01).SetField02(2).Save(ctx); err != nil {
			t.Fatal(err)
		} else if table01.Field02 != 2 {
			t.Fatal(table01)
		}

		if table01, err := client.Table01ForEnt.Query().Where(table01forent.Field01(t.Name())).Only(ctx); err != nil {
			t.Fatal(err)
		} else if table01.Field01 != t.Name() || table01.Field02 != 2 || table01.Field03 != true {
			t.Fatal(table01)
		}
	}
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()

	prepareData(t, ctx)

	for _, client := range clients {
		if affected, err := client.Table01ForEnt.Update().SetField02(2).Where(table01forent.Field01(t.Name())).Save(ctx); err != nil {
			t.Fatal(err)
		} else if affected != 1 {
			t.Fatal(affected)
		}

		if table01, err := client.Table01ForEnt.Query().Where(table01forent.Field01(t.Name())).Only(ctx); err != nil {
			t.Fatal(err)
		} else if table01.Field01 != t.Name() || table01.Field02 != 2 || table01.Field03 != true {
			t.Fatal(table01)
		}
	}
}

func TestUpsertOne(t *testing.T) {
	ctx := context.Background()

	prepareData(t, ctx)

	wg := new(sync.WaitGroup)
	defer wg.Wait()
	for _, client := range clients {
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func(client *ent.Client, i int) {
				defer wg.Done()

				field01 := t.Name() + strconv.Itoa(i)

				for j := 0; j < 5; j++ {
					if err := client.Table01ForEnt.Create().
						SetField01(field01).
						SetField02(j).
						OnConflictColumns(table01forent.FieldField01).
						UpdateNewValues().
						Exec(ctx); err != nil {
						t.Fatal(err)
					}

					if table01, err := client.Table01ForEnt.Query().Where(table01forent.Field01(field01)).Only(ctx); err != nil {
						t.Fatal(err)
					} else if table01.Field01 != field01 || table01.Field02 != j {
						t.Fatal(table01)
					}
				}
			}(client, i)
		}
	}
}

func TestUpsertBulk(t *testing.T) {
	ctx := context.Background()

	prepareData(t, ctx)

	for _, client := range clients {
		for i := 0; i < 10; i++ {
			field05 := t.Name() + strconv.Itoa(i)

			table01ForEntCreates := []*ent.Table01ForEntCreate{
				client.Table01ForEnt.Create().SetField01(t.Name() + "-5").SetField02(i).SetField05(field05),
				client.Table01ForEnt.Create().SetField01(t.Name() + "-6").SetField02(i).SetField05(field05),
				client.Table01ForEnt.Create().SetField01(t.Name() + "-7").SetField02(i).SetField05(field05),
			}
			if err := client.Table01ForEnt.CreateBulk(table01ForEntCreates...).
				OnConflictColumns(table01forent.FieldField01).
				UpdateNewValues().
				Exec(ctx); err != nil {
				t.Fatal(err)
			}

			if table01s, err := client.Table01ForEnt.Query().Where(table01forent.Field05(field05)).All(ctx); err != nil {
				t.Fatal(err)
			} else if len(table01s) != 3 {
				t.Fatal(table01s)
			} else {
				for _, table01 := range table01s {
					if table01.Field02 != i {
						t.Fatal(table01)
					}
				}
			}
		}
	}
}

func TestQuery(t *testing.T) {
	ctx := context.Background()

	prepareData(t, ctx)

	for _, client := range clients {
		if table01, err := client.Table01ForEnt.Query().Where(table01forent.Field01(t.Name())).Only(ctx); err != nil {
			t.Fatal(err)
		} else if table01.Field01 != t.Name() || table01.Field02 != 1 || table01.Field03 != true {
			t.Fatal(table01)
		}

		if results, err := client.Table01ForEnt.Query().Where(table01forent.Field05(t.Name())).All(ctx); err != nil {
			t.Fatal(err)
		} else if len(results) != 3 {
			t.Fatal(len(results))
		}
	}
}

func TestSelect(t *testing.T) {
	ctx := context.Background()

	prepareData(t, ctx)

	for _, client := range clients {
		if field05s, err := client.Table01ForEnt.Query().
			Where(table01forent.Field05(t.Name())).
			Select(table01forent.FieldField05).
			Strings(ctx); err != nil {
			t.Fatal(err)
		} else if len(field05s) != 3 {
			t.Fatal(field05s)
		}

		if field05s, err := client.Table01ForEnt.Query().
			Where(table01forent.Field05(t.Name())).
			Unique(true).
			Select(table01forent.FieldField05).
			Strings(ctx); err != nil {
			t.Fatal(err)
		} else if len(field05s) != 1 {
			t.Fatal(field05s)
		}

		if count, err := client.Table01ForEnt.Query().
			Where(table01forent.Field05(t.Name())).
			Unique(true).
			Select(table01forent.FieldField05).
			Count(ctx); err != nil {
			t.Fatal(err)
		} else if count != 1 {
			t.Fatal(count)
		}

		result := []struct {
			Field01 string `json:"field01"`
			Field02 string `json:"field02"`
			Field05 string `json:"field05"`
		}{}
		if err := client.Table01ForEnt.Query().
			Where(table01forent.Field05(t.Name())).
			Select(table01forent.FieldField01, table01forent.FieldField02, table01forent.FieldField05).
			Scan(ctx, &result); err != nil {
			t.Fatal(err)
		} else if len(result) != 3 {
			t.Fatal(result)
		}
	}
}

func TestDeleteOne(t *testing.T) {
	ctx := context.Background()

	prepareData(t, ctx)

	for _, client := range clients {
		if table01, err := client.Table01ForEnt.Query().Where(table01forent.Field01(t.Name())).Only(ctx); err != nil {
			t.Fatal(err)
		} else if err := client.Table01ForEnt.DeleteOne(table01).Exec(ctx); err != nil {
			t.Fatal(err)
		}

		if table01s, err := client.Table01ForEnt.Query().Where(table01forent.Field01(t.Name())).All(ctx); err != nil {
			t.Fatal(err)
		} else if len(table01s) != 0 {
			t.Fatal(table01s)
		}
	}
}

func TestDelete(t *testing.T) {
	ctx := context.Background()

	prepareData(t, ctx)

	for _, client := range clients {
		if affected, err := client.Table01ForEnt.Delete().Where(table01forent.Field05(t.Name())).Exec(ctx); err != nil {
			t.Fatal(err)
		} else if affected != 3 {
			t.Fatal(affected)
		}

		if table01s, err := client.Table01ForEnt.Query().Where(table01forent.Field01(t.Name())).All(ctx); err != nil {
			t.Fatal(err)
		} else if len(table01s) != 1 {
			t.Fatal(table01s)
		}
	}
}

func TestTx(t *testing.T) {
	ctx := context.Background()

	prepareData(t, ctx)

	job := func() error { return nil }

	for _, client := range clients {
		if tx, err := client.Tx(ctx); err != nil {
			t.Fatal(err)
		} else if err := job(); err != nil {
			if err := tx.Rollback(); err != nil {
				t.Fatal(err)
			}
		} else {
			if err := tx.Commit(); err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestCommit(t *testing.T) {
	ctx := context.Background()

	prepareData(t, ctx)

	getTable01 := func(client *ent.Client) (*ent.Table01ForEnt, error) {
		return client.Table01ForEnt.Query().Where(table01forent.Field01(t.Name())).Only(ctx)
	}

	for _, client := range clients {
		if tx, err := client.Tx(ctx); err != nil {
			t.Fatal(err)
		} else if _, err := tx.Client().Table01ForEnt.Update().SetField02(2).Where(table01forent.Field01(t.Name())).Save(ctx); err != nil {
			t.Fatal(err)
		} else if table01, err := getTable01(tx.Client()); err != nil {
			t.Fatal(err)
		} else if table01.Field01 != t.Name() || table01.Field02 != 2 {
			t.Fatal(table01)
		} else if table01, err := getTable01(client); err != nil {
			t.Fatal(table01)
		} else if table01.Field01 != t.Name() || table01.Field02 != 1 {
			t.Fatal(table01)
		} else if tx.Commit(); err != nil {
			t.Fatal(err)
		} else if table01, err := getTable01(client); err != nil {
			t.Fatal(table01)
		} else if table01.Field01 != t.Name() || table01.Field02 != 2 {
			t.Fatal(table01)
		}
	}
}

func TestRollback(t *testing.T) {
	ctx := context.Background()

	prepareData(t, ctx)

	getTable01 := func(client *ent.Client) (*ent.Table01ForEnt, error) {
		return client.Table01ForEnt.Query().Where(table01forent.Field01(t.Name())).Only(ctx)
	}

	for _, client := range clients {
		if tx, err := client.Tx(ctx); err != nil {
			t.Fatal(err)
		} else if _, err := tx.Client().Table01ForEnt.Update().SetField02(2).Where(table01forent.Field01(t.Name())).Save(ctx); err != nil {
			t.Fatal(err)
		} else if table01, err := getTable01(tx.Client()); err != nil {
			t.Fatal(err)
		} else if table01.Field01 != t.Name() || table01.Field02 != 2 {
			t.Fatal(table01)
		} else if table01, err := getTable01(client); err != nil {
			t.Fatal(table01)
		} else if table01.Field01 != t.Name() || table01.Field02 != 1 {
			t.Fatal(table01)
		} else if tx.Rollback(); err != nil {
			t.Fatal(err)
		} else if table01, err := getTable01(client); err != nil {
			t.Fatal(table01)
		} else if table01.Field01 != t.Name() || table01.Field02 != 1 {
			t.Fatal(table01)
		}
	}
}

func TestAggregate(t *testing.T) {
	ctx := context.Background()

	prepareData(t, ctx)

	for _, client := range clients {
		if max, err := client.Table01ForEnt.Query().Where(table01forent.Field05(t.Name())).
			Aggregate(ent.Max(table01forent.FieldField02)).Int(ctx); err != nil {
			t.Fatal(err)
		} else if max != 3 {
			t.Fatal(max)
		}

		result := []struct {
			Sum, Min, Max, Count int
		}{}
		if err := client.Table01ForEnt.Query().Where(table01forent.Field05(t.Name())).
			Aggregate(
				ent.Sum(table01forent.FieldField02),
				ent.Min(table01forent.FieldField02),
				ent.Max(table01forent.FieldField02),
				ent.Count(),
			).Scan(ctx, &result); err != nil {
			t.Fatal(err)
		} else if len(result) != 1 {
			t.Fatal(result)
		} else if result[0].Sum != 6 || result[0].Min != 1 || result[0].Max != 3 || result[0].Count != 3 {
			t.Fatal(result)
		}
	}
}

func TestGroupBy(t *testing.T) {
	ctx := context.Background()

	prepareData(t, ctx)

	for _, client := range clients {
		if field01s, err := client.Table01ForEnt.Query().Where(table01forent.Field05(t.Name())).
			GroupBy(table01forent.FieldField05).
			Strings(ctx); err != nil {
			t.Fatal(err)
		} else if len(field01s) != 1 {
			t.Fatal(field01s)
		} else if field01s[0] != t.Name() {
			t.Fatal(field01s)
		}

		result := []struct {
			Field09 string
			Sum     int
		}{}
		if err := client.Table01ForEnt.Query().Where(table01forent.Field05(t.Name())).
			GroupBy(table01forent.FieldField09).
			Aggregate(ent.Sum(table01forent.FieldField02)).
			Scan(ctx, &result); err != nil {
			t.Fatal(err)
		} else if len(result) != 1 {
			t.Fatal(result)
		} else if result[0].Field09 != "Field" || result[0].Sum != 6 {
			t.Fatal(result)
		}
	}
}

func TestLimit(t *testing.T) {
	ctx := context.Background()

	prepareData(t, ctx)

	for _, client := range clients {
		limit := 2
		if table01s, err := client.Table01ForEnt.Query().Where(table01forent.Field05(t.Name())).
			Limit(limit).
			All(ctx); err != nil {
			t.Fatal(err)
		} else if len(table01s) != limit {
			t.Fatal(table01s)
		}
	}
}

func TestOffset(t *testing.T) {
	ctx := context.Background()

	prepareData(t, ctx)

	for _, client := range clients {
		if table01s, err := client.Table01ForEnt.Query().Where(table01forent.Field05(t.Name())).
			Order(table01forent.ByField01()).
			Offset(2).
			All(ctx); err != nil {
			t.Fatal(err)
		} else if len(table01s) != 1 {
			t.Fatal(table01s)
		} else if table01s[0].Field01 != t.Name()+"-3" {
			t.Fatal(table01s)
		}
	}
}

func TestOrder(t *testing.T) {
	ctx := context.Background()

	prepareData(t, ctx)

	for _, client := range clients {
		if table01s, err := client.Table01ForEnt.Query().Where(table01forent.Field05(t.Name())).
			Order(ent.Asc(table01forent.FieldField02)).
			All(ctx); err != nil {
			t.Fatal(err)
		} else if len(table01s) != 3 {
			t.Fatal(table01s)
		}

		if table01s, err := client.Table01ForEnt.Query().Where(table01forent.Field05(t.Name())).
			Order(
				table01forent.ByField03(),
				table01forent.ByField04(ent_sql.OrderDesc()),
			).
			All(ctx); err != nil {
			t.Fatal(err)
		} else if len(table01s) != 3 {
			t.Fatal(table01s)
		}
	}
}
