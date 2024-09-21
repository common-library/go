package beego_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type dataBaseInfo struct {
	aliasName  string
	driverName string
	dataSource string
}

var dataBaseInfos = []dataBaseInfo{}

type Table01ForBeego struct {
	Field01 string `orm:"pk;size(512)"`
	Field02 int
	Field03 bool
}

func getOrmers(t *testing.T) []orm.Ormer {
	if t != nil {
		t.Parallel()
	}

	ormers := []orm.Ormer{}

	for _, dataBaseInfo := range dataBaseInfos {
		ormers = append(ormers, orm.NewOrmUsingDB(dataBaseInfo.aliasName))
	}

	return ormers
}

func TestMain(m *testing.M) {
	setup := func() {
		models := []any{
			new(Table01ForBeego),
		}
		orm.RegisterModel(models...)

		if len(os.Getenv("MYSQL_DSN")) != 0 {
			dataBaseInfos = append(dataBaseInfos, dataBaseInfo{aliasName: "mysql", driverName: "mysql", dataSource: os.Getenv("MYSQL_DSN") + "mysql"})
		}
		if len(os.Getenv("POSTGRESQL_DSN")) != 0 {
			dataBaseInfos = append(dataBaseInfos, dataBaseInfo{aliasName: "postgres", driverName: "postgres", dataSource: os.Getenv("POSTGRESQL_DSN") + " dbname=postgres"})
		}

		if len(dataBaseInfos) != 0 {
			dataBaseInfos[0].aliasName = "default"
		}

		for _, dataBaseInfo := range dataBaseInfos {
			if err := orm.RegisterDataBase(dataBaseInfo.aliasName, dataBaseInfo.driverName, dataBaseInfo.dataSource); err != nil {
				panic(err)
			} else if err := orm.RunSyncdb(dataBaseInfo.aliasName, true, true); err != nil {
				panic(err)
			}
		}
	}

	teardown := func() {
		tables := []string{
			"table01_for_beego",
		}

		for _, ormer := range getOrmers(nil) {
			for _, table := range tables {
				rawSeter := ormer.Raw("DROP TABLE IF EXISTS " + table)
				if _, err := rawSeter.Exec(); err != nil {
					panic(err)
				}
			}
		}
	}

	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestRead(t *testing.T) {
	table01 := Table01ForBeego{}

	for _, ormer := range getOrmers(t) {
		table01 = Table01ForBeego{Field01: t.Name(), Field02: 1, Field03: true}
		if _, err := ormer.Insert(&table01); err != nil && err != orm.ErrLastInsertIdUnavailable {
			t.Fatal(err)
		}

		table01 = Table01ForBeego{Field01: t.Name()}
		if err := ormer.Read(&table01); err != nil {
			t.Fatal(err)
		} else if table01.Field01 != t.Name() || table01.Field02 != 1 || table01.Field03 != true {
			t.Fatal(table01)
		}

		table01 = Table01ForBeego{Field01: t.Name(), Field02: 1}
		if err := ormer.Read(&table01, "Field01"); err != nil {
			t.Fatal(err)
		} else if table01.Field02 != 1 {
			t.Fatal(table01)
		}
	}
}

func TestReadForUpdate(t *testing.T) {
	table01 := Table01ForBeego{}

	for _, ormer := range getOrmers(t) {
		table01 = Table01ForBeego{Field01: t.Name(), Field02: 1, Field03: true}
		if _, err := ormer.Insert(&table01); err != nil && err != orm.ErrLastInsertIdUnavailable {
			t.Fatal(err)
		}

		table01 = Table01ForBeego{Field01: t.Name()}
		if err := ormer.ReadForUpdate(&table01); err != nil {
			t.Fatal(err)
		} else if table01.Field01 != t.Name() || table01.Field02 != 1 || table01.Field03 != true {
			t.Fatal(table01)
		}
	}
}

func TestReadOrCreate(t *testing.T) {
	table01 := Table01ForBeego{}

	for _, ormer := range getOrmers(t) {
		t.Log(ormer.Driver().Name())

		table01 = Table01ForBeego{Field01: t.Name(), Field02: 1, Field03: true}
		if created, _, err := ormer.ReadOrCreate(&table01, "field01"); err != nil && err != orm.ErrLastInsertIdUnavailable {
			t.Fatal(err)
		} else if ormer.Driver().Type() == orm.DRMySQL && created == false {
			t.Fatal(created)
		} else if ormer.Driver().Type() == orm.DRPostgres && created {
			t.Fatal(created)
		}

		table01 = Table01ForBeego{Field01: t.Name()}
		if err := ormer.Read(&table01, "field01"); err != nil && err != orm.ErrLastInsertIdUnavailable {
			t.Fatal(err)
		} else if table01.Field01 != t.Name() || table01.Field02 != 1 || table01.Field03 != true {
			t.Fatal(table01)
		}
	}
}

func TestLoadRelated(t *testing.T) {
}

func TestQueryM2M(t *testing.T) {
}

func TestQueryTable(t *testing.T) {
}

func TestInsert(t *testing.T) {
	TestRead(t)
}

func TestInsertOrUpdate(t *testing.T) {
	table01 := Table01ForBeego{}

	for _, ormer := range getOrmers(t) {
		table01 = Table01ForBeego{Field01: t.Name(), Field02: 1, Field03: true}
		if _, err := ormer.InsertOrUpdate(&table01, "field01"); err != nil && err != orm.ErrLastInsertIdUnavailable {
			t.Fatal(err)
		}

		table01 = Table01ForBeego{Field01: t.Name(), Field02: 2, Field03: true}
		if _, err := ormer.InsertOrUpdate(&table01, "field01"); err != nil && err != orm.ErrLastInsertIdUnavailable {
			t.Fatal(err)
		}

		table01 = Table01ForBeego{Field01: t.Name()}
		if err := ormer.Read(&table01); err != nil {
			t.Fatal(err)
		} else if table01.Field01 != t.Name() || table01.Field02 != 2 || table01.Field03 != true {
			t.Fatal(table01)
		}
	}
}

func TestInsertMulti(t *testing.T) {
	table01s := []Table01ForBeego{
		{Field01: t.Name() + "-1", Field02: 1, Field03: true},
		{Field01: t.Name() + "-2", Field02: 1, Field03: true},
		{Field01: t.Name() + "-3", Field02: 1, Field03: true},
	}

	for _, ormer := range getOrmers(t) {
		if successNums, err := ormer.InsertMulti(1000, table01s); err != nil {
			t.Fatal(err)
		} else if successNums != int64(len(table01s)) {
			t.Fatal(successNums)
		}
	}
}

func TestUpdate(t *testing.T) {
	table01 := Table01ForBeego{}

	for _, ormer := range getOrmers(t) {
		table01 = Table01ForBeego{Field01: t.Name(), Field02: 1, Field03: true}
		if _, err := ormer.Insert(&table01); err != nil && err != orm.ErrLastInsertIdUnavailable {
			t.Fatal(err)
		}

		table01 = Table01ForBeego{Field01: t.Name(), Field02: 2, Field03: true}
		if updateNums, err := ormer.Update(&table01); err != nil {
			t.Fatal(err)
		} else if updateNums != 1 {
			t.Fatal(updateNums)
		}
		table01 = Table01ForBeego{Field01: t.Name()}
		if err := ormer.Read(&table01); err != nil {
			t.Fatal(err)
		} else if table01.Field01 != t.Name() || table01.Field02 != 2 || table01.Field03 != true {
			t.Fatal(table01)
		}

		table01 = Table01ForBeego{Field01: t.Name(), Field02: 3, Field03: false}
		if updateNums, err := ormer.Update(&table01, "field02"); err != nil {
			t.Fatal(err)
		} else if updateNums != 1 {
			t.Fatal(updateNums)
		}
		table01 = Table01ForBeego{Field01: t.Name()}
		if err := ormer.Read(&table01); err != nil {
			t.Fatal(err)
		} else if table01.Field01 != t.Name() || table01.Field02 != 3 || table01.Field03 != true {
			t.Fatal(table01)
		}
	}
}

func TestDelete(t *testing.T) {
	table01 := Table01ForBeego{}

	for _, ormer := range getOrmers(t) {
		table01 = Table01ForBeego{Field01: t.Name(), Field02: 1, Field03: true}
		if _, err := ormer.Insert(&table01); err != nil && err != orm.ErrLastInsertIdUnavailable {
			t.Fatal(err)
		}

		table01 = Table01ForBeego{Field01: t.Name()}
		if deleteNums, err := ormer.Delete(&table01); err != nil {
			t.Fatal(err)
		} else if deleteNums != 1 {
			t.Fatal(err)
		}

		table01 = Table01ForBeego{Field01: t.Name()}
		if err := ormer.Read(&table01); err != orm.ErrNoRows {
			t.Fatal(err)
		}
	}
}

func TestRaw(t *testing.T) {
	table01 := Table01ForBeego{}

	for _, ormer := range getOrmers(t) {
		table01 = Table01ForBeego{Field01: t.Name(), Field02: 1, Field03: true}
		if _, err := ormer.Insert(&table01); err != nil && err != orm.ErrLastInsertIdUnavailable {
			t.Fatal(err)
		}

		rawSeter := ormer.Raw("SELECT field03 FROM table01_for_beego WHERE field01 = ? AND field02 = ?", t.Name(), 1)

		field03 := false
		if err := rawSeter.QueryRow(&field03); err != nil {
			t.Fatal(err)
		} else if field03 == false {
			t.Fatal(field03)
		}
	}
}

func TestDoTx(t *testing.T) {
	table01 := Table01ForBeego{}

	for _, ormer := range getOrmers(t) {
		const errorString = "error"
		if err := ormer.DoTx(func(ctx context.Context, txOrmer orm.TxOrmer) error {
			table01 := Table01ForBeego{Field01: t.Name(), Field02: 1, Field03: true}
			if _, err := txOrmer.Insert(&table01); err != nil && err != orm.ErrLastInsertIdUnavailable {
				t.Fatal(err)
			}

			return errors.New(errorString)
		}); err.Error() != errorString {
			t.Fatal(err)
		}
		table01 = Table01ForBeego{Field01: t.Name()}
		if err := ormer.Read(&table01); err != orm.ErrNoRows {
			t.Fatal(err)
		}

		if err := ormer.DoTx(func(ctx context.Context, txOrmer orm.TxOrmer) error {
			table01 := Table01ForBeego{Field01: t.Name(), Field02: 1, Field03: true}
			if _, err := txOrmer.Insert(&table01); err != nil && err != orm.ErrLastInsertIdUnavailable {
				t.Fatal(err)
			}

			return nil
		}); err != nil {
			t.Fatal(err)
		}
		table01 = Table01ForBeego{Field01: t.Name()}
		if err := ormer.Read(&table01); err != nil {
			t.Fatal(err)
		} else if table01.Field01 != t.Name() || table01.Field02 != 1 || table01.Field03 != true {
			t.Fatal(table01)
		}
	}
}

func TestBegin(t *testing.T) {
	for _, ormer := range getOrmers(t) {
		if txOrmer, err := ormer.Begin(); err != nil {
			t.Fatal(err)
		} else if err := func() error { return nil }; err != nil {
			if err := txOrmer.Rollback(); err != nil {
				t.Fatal(err)
			}
		} else {
			if err := txOrmer.Commit(); err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestCommit(t *testing.T) {
	table01 := Table01ForBeego{}

	for _, ormer := range getOrmers(t) {
		txOrmer, err := ormer.Begin()
		if err != nil {
			t.Fatal(err)
		}

		table01 = Table01ForBeego{Field01: t.Name(), Field02: 1, Field03: true}
		if _, err := txOrmer.Insert(&table01); err != nil && err != orm.ErrLastInsertIdUnavailable {
			t.Fatal(err)
		}

		table01 = Table01ForBeego{Field01: t.Name()}
		if err := ormer.Read(&table01); err != orm.ErrNoRows {
			t.Fatal(err)
		}

		if err := txOrmer.Commit(); err != nil {
			t.Fatal(err)
		}

		table01 = Table01ForBeego{Field01: t.Name()}
		if err := ormer.Read(&table01); err != nil {
			t.Fatal(err)
		} else if table01.Field01 != t.Name() || table01.Field02 != 1 || table01.Field03 != true {
			t.Fatal(table01)
		}
	}
}

func TestRollback(t *testing.T) {
	table01 := Table01ForBeego{}

	for _, ormer := range getOrmers(t) {
		txOrmer, err := ormer.Begin()
		if err != nil {
			t.Fatal(err)
		}

		table01 = Table01ForBeego{Field01: t.Name(), Field02: 1, Field03: true}
		if _, err := txOrmer.Insert(&table01); err != nil && err != orm.ErrLastInsertIdUnavailable {
			t.Fatal(err)
		}

		table01 = Table01ForBeego{Field01: t.Name()}
		if err := ormer.Read(&table01); err != orm.ErrNoRows {
			t.Fatal(err)
		}

		if err := txOrmer.Rollback(); err != nil {
			t.Fatal(err)
		}

		table01 = Table01ForBeego{Field01: t.Name()}
		if err := ormer.Read(&table01); err != orm.ErrNoRows {
			t.Fatal(err)
		}
	}
}

func TestRollbackUnlessCommit(t *testing.T) {
	test := func(ormer orm.Ormer, commit bool) {
		txOrmer, err := ormer.Begin()
		if err != nil {
			t.Fatal(err)
		}

		defer func() {
			if err := txOrmer.RollbackUnlessCommit(); err != nil {
				t.Fatal(err)
			}
		}()

		table01 := Table01ForBeego{Field01: t.Name(), Field02: 1, Field03: true}
		if _, err := txOrmer.Insert(&table01); err != nil && err != orm.ErrLastInsertIdUnavailable {
			t.Fatal(err)
		}

		if commit {
			if err := txOrmer.Commit(); err != nil {
				t.Fatal(err)
			}
		}
	}

	table01 := Table01ForBeego{}
	for _, ormer := range getOrmers(t) {
		test(ormer, false)
		table01 = Table01ForBeego{Field01: t.Name()}
		if err := ormer.Read(&table01); err != orm.ErrNoRows {
			t.Fatal(err)
		}

		test(ormer, true)
		table01 = Table01ForBeego{Field01: t.Name()}
		if err := ormer.Read(&table01); err != nil {
			t.Fatal(err)
		} else if table01.Field01 != t.Name() || table01.Field02 != 1 || table01.Field03 != true {
			t.Fatal(table01)
		}
	}
}

func TestDriver(t *testing.T) {
	for _, ormer := range getOrmers(t) {
		t.Log(ormer.Driver().Name(), ormer.Driver().Type())
	}
}
