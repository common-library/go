package gorm_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/common-library/go/testutil"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	mysqlDriver "gorm.io/driver/mysql"
	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbs = map[string]*gorm.DB{}

type Table01ForGorm struct {
	Field01 string  `gorm:"primaryKey"`
	Field02 *string `gorm:"primaryKey;size:128"`

	Field03 uint `gorm:"index"`
	Field04 *uint

	Field05 bool `gorm:"index;default:true"`
	Field06 *bool

	Created time.Time `gorm:"autoCreateTime"`
	Updated time.Time `gorm:"autoUpdateTime"`
}

type Table02ForGorm struct {
	Field01 string `gorm:"primaryKey"`
	Field02 *string

	Table03Field01 string

	Table03ForGorm Table03ForGorm `gorm:"foreignKey:Table03Field01"`
}

type Table03ForGorm struct {
	Field01 string `gorm:"primaryKey;size:128"`
	Field02 *string
}

func (t *Table01ForGorm) BeforeSave(tx *gorm.DB) error {
	return nil
}

func (t *Table01ForGorm) AfterSave(tx *gorm.DB) error {
	return nil
}

func (t *Table01ForGorm) BeforeCreate(tx *gorm.DB) error {
	return nil
}

func (t *Table01ForGorm) AfterCreate(tx *gorm.DB) error {
	return nil
}

func (t *Table01ForGorm) BeforeUpdate(tx *gorm.DB) error {
	return nil
}

func (t *Table01ForGorm) AfterUpdate(tx *gorm.DB) error {
	return nil
}

func (t *Table01ForGorm) BeforeDelete(tx *gorm.DB) error {
	return nil
}

func (t *Table01ForGorm) AfterDelete(tx *gorm.DB) error {
	return nil
}

func prepareData(t *testing.T) {
	for kind, db := range dbs {
		name := t.Name()
		table01s := []*Table01ForGorm{
			{Field01: "b", Field02: &name, Field03: 2},
			{Field01: "a", Field02: &name, Field03: 1},
			{Field01: "c", Field02: &name, Field03: 3},
		}
		if result := db.Create(table01s); result.Error != nil {
			t.Log(name)
			t.Log(kind)
			t.Fatal(result.Error)
		}
	}

	t.Parallel()
}

var containers = []testcontainers.Container{}

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

		var db *gorm.DB
		maxRetries := 10
		for i := 0; i < maxRetries; i++ {
			db, err = gorm.Open(mysqlDriver.Open(mysqlDSN), &gorm.Config{})
			if err == nil {
				sqlDB, sqlErr := db.DB()
				if sqlErr == nil {
					if pingErr := sqlDB.Ping(); pingErr == nil {
						break
					}
				}
			}
			if i == maxRetries-1 {
				panic(fmt.Errorf("failed to connect to MySQL after %d retries: %w", maxRetries, err))
			}
			time.Sleep(500 * time.Millisecond)
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

		maxRetries = 10
		for i := 0; i < maxRetries; i++ {
			db, err = gorm.Open(postgresDriver.Open(postgresDSN), &gorm.Config{})
			if err == nil {
				sqlDB, sqlErr := db.DB()
				if sqlErr == nil {
					if pingErr := sqlDB.Ping(); pingErr == nil {
						break
					}
				}
			}
			if i == maxRetries-1 {
				panic(fmt.Errorf("failed to connect to PostgreSQL after %d retries: %w", maxRetries, err))
			}
			time.Sleep(500 * time.Millisecond)
		}
		dbs["postgres"] = db

		tables := []any{
			&Table01ForGorm{},
			&Table02ForGorm{},
			&Table03ForGorm{},
		}
		for _, db := range dbs {
			if err := db.AutoMigrate(tables...); err != nil {
				panic(err)
			}
		}
	}

	teardown := func() {
		tables := []string{
			"table01_for_gorms",
			"table02_for_gorms",
			"table03_for_gorms",
		}
		for _, db := range dbs {
			for _, table := range tables {
				if result := db.Exec("DROP TABLE IF EXISTS " + table); result.Error != nil {
					panic(result.Error)
				}
			}
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

func TestCreate(t *testing.T) {
	for _, db := range dbs {
		name := t.Name()
		var table01 Table01ForGorm
		var table01_map map[string]any
		var table01s []*Table01ForGorm
		var table01s_map []map[string]any

		table01 = Table01ForGorm{Field01: "a", Field02: &name, Field03: 1}
		if result := db.Create(&table01); result.Error != nil {
			t.Fatal(result.Error)
		}

		table01s = []*Table01ForGorm{{Field01: "b", Field02: &name, Field03: 2}, {Field01: "c", Field02: &name, Field03: 3}}
		if result := db.Create(table01s); result.Error != nil {
			t.Fatal(result.Error)
		}

		table01s = []*Table01ForGorm{{Field01: "d", Field02: &name, Field03: 4}, {Field01: "e", Field02: &name, Field03: 5}}
		if result := db.CreateInBatches(table01s, 10); result.Error != nil {
			t.Fatal(result.Error)
		}

		table01_map = map[string]any{"field01": "f", "field02": &name, "field03": 6}
		if result := db.Model(&Table01ForGorm{}).Create(&table01_map); result.Error != nil {
			t.Fatal(result.Error)
		}

		table01s_map = []map[string]any{{"field01": "g", "field02": &name, "field03": 7}, {"field01": "h", "field02": &name, "field03": 8}}
		if result := db.Model(&Table01ForGorm{}).Create(table01s_map); result.Error != nil {
			t.Fatal(result.Error)
		}
	}
}
func TestCreateSelect(t *testing.T) {
	for _, db := range dbs {
		name := t.Name()
		table01 := Table01ForGorm{Field01: "y", Field02: &name, Field03: 25}
		table01_02 := Table01ForGorm{}
		if result := db.Select("field01", "field02").Create(&table01); result.Error != nil {
			t.Fatal(result.Error)
		} else if result := db.First(&table01_02, "field01 = ? AND field02 = ?", "y", t.Name()); result.Error != nil {
			t.Fatal(result.Error)
		} else if table01_02.Field01 != "y" || table01_02.Field03 != 0 {
			t.Fatal(table01)
		}
	}
}

func TestCreateOmit(t *testing.T) {
	for _, db := range dbs {
		name := t.Name()
		table01 := Table01ForGorm{Field01: "z", Field02: &name, Field03: 26}
		table01_02 := Table01ForGorm{}
		if result := db.Omit("field03").Create(&table01); result.Error != nil {
			t.Fatal(result.Error)
		} else if result := db.First(&table01_02, "field01 = ? AND field02 = ?", "z", t.Name()); result.Error != nil {
			t.Fatal(result.Error)
		} else if table01_02.Field01 != "z" || table01_02.Field03 != 0 {
			t.Fatal(table01)
		}
	}
}

func TestFirst(t *testing.T) {
	prepareData(t)

	for _, db := range dbs {
		table01 := Table01ForGorm{}
		table01_map := map[string]any{}

		table01 = Table01ForGorm{}
		if result := db.First(&table01, "field02 = ?", t.Name()); result.Error != nil {
			t.Fatal(result.Error)
		} else if table01.Field01 != "a" || table01.Field03 != 1 {
			t.Fatal(table01)
		}

		table01_map = map[string]any{}
		if result := db.Model(&Table01ForGorm{}).First(&table01_map, "field02 = ?", t.Name()); result.Error != nil {
			t.Fatal(result.Error)
		} else if table01_map["field01"].(string) != "a" || table01_map["field03"].(uint) != 1 {
			t.Fatal(table01_map)
		}

		table01 = Table01ForGorm{}
		if result := db.First(&table01, "field01 = ? AND field02 = ?", "a", t.Name()); result.Error != nil {
			t.Fatal(result.Error)
		} else if table01.Field01 != "a" || table01.Field03 != 1 {
			t.Fatal(table01)
		}

		table01 = Table01ForGorm{Field01: "b"}
		if result := db.First(&table01, "field02 = ?", t.Name()); result.Error != nil {
			t.Fatal(result.Error)
		} else if table01.Field01 != "b" || table01.Field03 != 2 {
			t.Fatal(table01)
		}
	}
}

func TestFirstWhere(t *testing.T) {
	prepareData(t)

	for _, db := range dbs {
		table01 := Table01ForGorm{}

		table01 = Table01ForGorm{}
		if result := db.First(&table01, "field01 = ? AND field02 = ?", "a", t.Name()); result.Error != nil {
			t.Fatal(result.Error)
		} else if table01.Field01 != "a" || table01.Field03 != 1 {
			t.Fatal(table01)
		}

		table01 = Table01ForGorm{}
		if result := db.First(&table01, map[string]any{"field01": "a", "field02": t.Name()}); result.Error != nil {
			t.Fatal(result.Error)
		} else if table01.Field01 != "a" || table01.Field03 != 1 {
			t.Fatal(table01)
		}

		table01 = Table01ForGorm{}
		if result := db.Where("field01 = ?", "a").First(&table01, "field02 = ?", t.Name()); result.Error != nil {
			t.Fatal(result.Error)
		} else if table01.Field01 != "a" || table01.Field03 != 1 {
			t.Fatal(table01)
		}

		table01 = Table01ForGorm{}
		if result := db.Where(map[string]any{"field01": "a"}).First(&table01, "field02 = ?", t.Name()); result.Error != nil {
			t.Fatal(result.Error)
		} else if table01.Field01 != "a" || table01.Field03 != 1 {
			t.Fatal(table01)
		}
	}
}

func TestLast(t *testing.T) {
	prepareData(t)

	for _, db := range dbs {
		table01 := Table01ForGorm{}
		if result := db.Last(&table01, "field02 = ?", t.Name()); result.Error != nil {
			t.Fatal(result.Error)
		} else if table01.Field01 != "c" || table01.Field03 != 3 {
			t.Fatal(table01)
		}
	}
}

func TestTake(t *testing.T) {
	prepareData(t)

	for _, db := range dbs {
		table01 := Table01ForGorm{}
		if result := db.Take(&table01, "field02 = ?", t.Name()); result.Error != nil {
			t.Fatal(result.Error)
		}

		table01_map := map[string]any{}
		if result := db.Table("table01_for_gorms").Take(&table01_map, "field02 = ?", t.Name()); result.Error != nil {
			t.Fatal(result.Error)
		}
	}
}

func TestFind(t *testing.T) {
	prepareData(t)

	for _, db := range dbs {
		table01s := []*Table01ForGorm{}
		if result := db.Find(&table01s, "field02 = ?", t.Name()); result.Error != nil {
			t.Fatal(result.Error)
		} else if result.RowsAffected != 3 {
			t.Fatal(result.RowsAffected)
		}
	}
}

func TestFindWhere(t *testing.T) {
	prepareData(t)

	for _, db := range dbs {
		table01s := []*Table01ForGorm{}

		table01s = []*Table01ForGorm{}
		if result := db.Find(&table01s, "field01 <> ? AND field02 = ?", "a", t.Name()); result.Error != nil {
			t.Fatal(result.Error)
		} else if result.RowsAffected != 2 {
			t.Fatal(result.RowsAffected)
		}

		table01s = []*Table01ForGorm{}
		if result := db.Find(&table01s, "field01 IN ? AND field02 = ?", []string{"a", "b"}, t.Name()); result.Error != nil {
			t.Fatal(result.Error)
		} else if result.RowsAffected != 2 {
			t.Fatal(result.RowsAffected)
		}

		table01s = []*Table01ForGorm{}
		if result := db.Where("field01 <> ?", "a").Find(&table01s, "field02 = ?", t.Name()); result.Error != nil {
			t.Fatal(result.Error)
		} else if result.RowsAffected != 2 {
			t.Fatal(result.RowsAffected)
		}

		table01s = []*Table01ForGorm{}
		if result := db.Where("field01 IN ?", []string{"a", "b"}).Find(&table01s, "field02 = ?", t.Name()); result.Error != nil {
			t.Fatal(result.Error)
		} else if result.RowsAffected != 2 {
			t.Fatal(result.RowsAffected)
		}
	}
}

func TestAnd(t *testing.T) {
	prepareData(t)

	for _, db := range dbs {
		table01s := []*Table01ForGorm{}

		table01s = []*Table01ForGorm{}
		if result := db.Where("field01 = ?", "a").Where("field03 = ?", 1).Find(&table01s, "field02 = ?", t.Name()); result.Error != nil {
			t.Fatal(result.Error)
		} else if result.RowsAffected != 1 {
			t.Fatal(result.RowsAffected)
		}

		table01s = []*Table01ForGorm{}
		if result := db.Where("field01 = ? AND field02 = ? AND field03 = ?", "a", t.Name(), 1).Find(&table01s); result.Error != nil {
			t.Fatal(result.Error)
		} else if result.RowsAffected != 1 {
			t.Fatal(result.RowsAffected)
		}
	}
}

func TestOr(t *testing.T) {
	prepareData(t)

	for _, db := range dbs {
		table01s := []*Table01ForGorm{}
		if result := db.Where("field01 = ? AND field02 = ?", "a", t.Name()).Or("field01 = ? AND field02 = ?", "b", t.Name()).Find(&table01s); result.Error != nil {
			t.Fatal(result.Error)
		} else if result.RowsAffected != 2 {
			t.Fatal(result.RowsAffected)
		}
	}
}

func TestNot(t *testing.T) {
	prepareData(t)

	for _, db := range dbs {
		table01s := []*Table01ForGorm{}
		if result := db.Not("field01 = ?", "a").Find(&table01s, "field02 = ?", t.Name()); result.Error != nil {
			t.Fatal(result.Error)
		} else if result.RowsAffected != 2 {
			t.Fatal(result.RowsAffected)
		}
	}
}

func TestSelect(t *testing.T) {
	prepareData(t)

	for _, db := range dbs {
		table01s := []*Table01ForGorm{}
		if result := db.Select("field01", "field02").Find(&table01s, "field02 = ?", t.Name()); result.Error != nil {
			t.Fatal(result.Error)
		} else {
			for _, table01 := range table01s {
				if table01.Field03 != 0 {
					t.Fatal(table01)
				}
			}
		}
	}
}

func TestOrder(t *testing.T) {
	prepareData(t)

	for _, db := range dbs {
		table01 := Table01ForGorm{}
		if result := db.Order("field01 desc").First(&table01, "field02 = ?", t.Name()); result.Error != nil {
			t.Fatal(result.Error)
		} else if table01.Field01 != "c" || table01.Field03 != 3 {
			t.Fatal(table01)
		}
	}
}

func TestLimit(t *testing.T) {
	prepareData(t)

	for _, db := range dbs {
		table01s := []*Table01ForGorm{}
		if result := db.Limit(2).Find(&table01s, "field02 = ?", t.Name()); result.Error != nil {
			t.Fatal(result.Error)
		} else if result.RowsAffected != 2 {
			t.Fatal(result.RowsAffected)
		} else if result := db.Limit(3).Find(&table01s, "field02 = ?", t.Name()); result.Error != nil {
			t.Fatal(result.Error)
		} else if result.RowsAffected != 3 {
			t.Fatal(result.RowsAffected)
		}
	}
}

func TestOffset(t *testing.T) {
	prepareData(t)

	for _, db := range dbs {
		table01s := []*Table01ForGorm{}
		if result := db.Limit(10).Offset(2).Find(&table01s, "field02 = ?", t.Name()); result.Error != nil {
			t.Fatal(result.Error)
		} else if result.RowsAffected != 1 {
			t.Fatal(result.RowsAffected)
		} else if result := db.Limit(3).Offset(3).Find(&table01s, "field02 = ?", t.Name()); result.Error != nil {
			t.Fatal(result.Error)
		} else if result.RowsAffected != 0 {
			t.Fatal(result.RowsAffected)
		}
	}
}

func TestGroup(t *testing.T) {
	prepareData(t)

	for _, db := range dbs {
		table01s := []*Table01ForGorm{}
		if result := db.Group("field01, field02").Find(&table01s, "field02 = ?", t.Name()); result.Error != nil {
			t.Fatal(result.Error)
		} else if result.RowsAffected != 3 {
			t.Fatal(result.RowsAffected)
		}
	}
}

func TestHaving(t *testing.T) {
	prepareData(t)

	for _, db := range dbs {
		table01s := []*Table01ForGorm{}
		if result := db.Group("field01, field02").Having("field01 = ?", "a").Find(&table01s, "field02 = ?", t.Name()); result.Error != nil {
			t.Fatal(result.Error)
		} else if result.RowsAffected != 1 {
			t.Fatal(result.RowsAffected)
		}
	}
}

func TestDistinct(t *testing.T) {
	prepareData(t)

	for _, db := range dbs {
		table01s := []*Table01ForGorm{}
		if result := db.Distinct("field02").Find(&table01s, "field02 = ?", t.Name()); result.Error != nil {
			t.Fatal(result.Error)
		} else if result.RowsAffected != 1 {
			t.Fatal(result.RowsAffected)
		}
	}
}

func TestJoins(t *testing.T) {
}

func TestScan(t *testing.T) {
	prepareData(t)

	for _, db := range dbs {
		table01 := Table01ForGorm{}
		if result := db.Model(&Table01ForGorm{}).Where("field01 = ?", "a").Where("field02 = ?", t.Name()).Scan(&table01); result.Error != nil {
			t.Fatal(result.Error)
		} else if table01.Field01 != "a" || table01.Field03 != 1 {
			t.Fatal(table01)
		}
	}
}

func TestSave(t *testing.T) {
	prepareData(t)

	for _, db := range dbs {
		name := t.Name()
		table01 := Table01ForGorm{}
		table01_temp := Table01ForGorm{}

		table01 = Table01ForGorm{Field01: "z", Field02: &name}
		table01_temp = Table01ForGorm{Field01: "z", Field02: &name}
		if result := db.Save(&table01); result.Error != nil {
			t.Fatal(result.Error)
		} else if result := db.First(&table01_temp); result.Error == gorm.ErrRecordNotFound {
			t.Fatal(result.Error)
		} else if table01_temp.Field01 != "z" || *table01_temp.Field02 != name || table01_temp.Field03 != 0 {
			t.Fatal(table01_temp)
		}

		table01 = Table01ForGorm{Field01: "z", Field02: &name, Field03: 26, Created: table01.Created}
		table01_temp = Table01ForGorm{Field01: "z", Field02: &name}
		if result := db.Save(&table01); result.Error != nil {
			t.Fatal(result.Error)
		} else if result := db.First(&table01_temp); result.Error == gorm.ErrRecordNotFound {
			t.Fatal(result.Error)
		} else if table01_temp.Field01 != "z" || *table01_temp.Field02 != name || table01_temp.Field03 != 26 {
			t.Fatal(table01_temp)
		}

		table01 = Table01ForGorm{Field01: "z", Field02: &name, Field03: 0}
		table01_temp = Table01ForGorm{Field01: "z", Field02: &name}
		if result := db.Select("field01", "field02").Save(&table01); result.Error != nil {
			t.Fatal(result.Error)
		} else if result := db.First(&table01_temp); result.Error == gorm.ErrRecordNotFound {
			t.Fatal(result.Error)
		} else if table01_temp.Field01 != "z" || *table01_temp.Field02 != name || table01_temp.Field03 != 26 {
			t.Fatal(table01_temp)
		}
	}
}

func TestUpdate(t *testing.T) {
	prepareData(t)

	for _, db := range dbs {
		name := t.Name()
		table01 := Table01ForGorm{}

		table01 = Table01ForGorm{Field01: "z", Field02: &name}
		if result := db.Model(&table01).Update("field03", 26); result.Error != nil {
			t.Fatal(result.Error)
		} else if result := db.First(&table01); result.Error != gorm.ErrRecordNotFound {
			t.Fatal(result.Error)
		}

		table01 = Table01ForGorm{}
		if result := db.Model(&Table01ForGorm{}).Where("field01 = ?", "a").Where("field02 = ?", t.Name()).Update("field03", 11); result.Error != nil {
			t.Fatal(result.Error)
		} else if result := db.First(&table01, "field01 = ? AND field02 = ?", "a", name); result.Error != nil {
			t.Fatal(result.Error)
		} else if table01.Field01 != "a" || table01.Field03 != 11 {
			t.Fatal(table01)
		}
	}
}

func TestUpdates(t *testing.T) {
	prepareData(t)

	for _, db := range dbs {
		name := t.Name()
		table01 := Table01ForGorm{}

		table01 = Table01ForGorm{Field01: "a", Field02: &name}
		if result := db.Model(&table01).Updates(Table01ForGorm{Field03: 11}); result.Error != nil {
			t.Fatal(result.Error)
		} else if result.RowsAffected != 1 {
			t.Fatal(result.RowsAffected)
		} else if table01.Field01 != "a" || table01.Field03 != 11 {
			t.Fatal(table01)
		} else if result := db.Model(&Table01ForGorm{Field02: &name}).Updates(Table01ForGorm{Field03: 11}); result.Error != nil {
			t.Fatal(result.Error)
		} else if result.RowsAffected != 3 {
			t.Fatal(result.RowsAffected)
		}

		table01 = Table01ForGorm{Field01: "a", Field02: &name}
		if result := db.Model(&table01).Updates(map[string]any{"field03": 111}); result.Error != nil {
			t.Fatal(result.Error)
		} else if result.RowsAffected != 1 {
			t.Fatal(result.RowsAffected)
		} else if table01.Field01 != "a" || table01.Field03 != 111 {
			t.Fatal(table01)
		}
	}
}

func TestUpdatesSelect(t *testing.T) {
	prepareData(t)

	for _, db := range dbs {
		name := t.Name()
		table01 := Table01ForGorm{}
		table01_temp := Table01ForGorm{}

		table01 = Table01ForGorm{Field01: "a", Field02: &name}
		table01_temp = Table01ForGorm{}
		if result := db.Model(&table01).Select("field03").Updates(Table01ForGorm{Field03: 11, Field05: false}); result.Error != nil {
			t.Fatal(result.Error)
		} else if result := db.First(&table01_temp, "field01 = ? AND field02 = ?", "a", t.Name()); result.Error != nil {
			t.Fatal(result.Error)
		} else if table01_temp.Field01 != "a" || table01_temp.Field03 != 11 || table01_temp.Field05 != true {
			t.Fatal(table01_temp)
		}
	}
}

func TestUpdatesOmit(t *testing.T) {
	prepareData(t)

	for _, db := range dbs {
		name := t.Name()
		table01 := Table01ForGorm{}
		table01_temp := Table01ForGorm{}

		table01 = Table01ForGorm{Field01: "a", Field02: &name}
		table01_temp = Table01ForGorm{}
		if result := db.Model(&table01).Select("field05").Omit("field03").Updates(Table01ForGorm{Field03: 11, Field05: false}); result.Error != nil {
			t.Fatal(result.Error)
		} else if result := db.First(&table01_temp, "field01 = ? AND field02 = ?", "a", t.Name()); result.Error != nil {
			t.Fatal(result.Error)
		} else if table01_temp.Field01 != "a" || table01_temp.Field03 != 1 || table01_temp.Field05 != false {
			t.Fatal(table01_temp)
		}
	}
}

func TestUpdateColumn(t *testing.T) {
	prepareData(t)

	for _, db := range dbs {
		name := t.Name()
		table01 := Table01ForGorm{Field01: "a", Field02: &name}
		table01_temp := Table01ForGorm{}
		if result := db.Model(&table01).UpdateColumn("field03", 11); result.Error != nil {
			t.Fatal(result.Error)
		} else if result := db.First(&table01_temp, "field01 = ? AND field02 = ?", "a", t.Name()); result.Error != nil {
			t.Fatal(result.Error)
		} else if table01_temp.Field01 != "a" || table01_temp.Field03 != 11 {
			t.Fatal(table01_temp)
		}
	}
}

func TestUpdateColumns(t *testing.T) {
	prepareData(t)

	for _, db := range dbs {
		name := t.Name()
		table01 := Table01ForGorm{Field01: "a", Field02: &name}
		table01_temp := Table01ForGorm{}
		if result := db.Model(&table01).Select("field03", "field05").UpdateColumns(Table01ForGorm{Field03: 11, Field05: false}); result.Error != nil {
			t.Fatal(result.Error)
		} else if result := db.First(&table01_temp, "field01 = ? AND field02 = ?", "a", t.Name()); result.Error != nil {
			t.Fatal(result.Error)
		} else if table01_temp.Field01 != "a" || table01_temp.Field03 != 11 || table01_temp.Field05 != false {
			t.Fatal(table01_temp)
		}
	}
}

func TestUpdateExpr(t *testing.T) {
	prepareData(t)

	for _, db := range dbs {
		name := t.Name()
		table01 := Table01ForGorm{Field01: "a", Field02: &name}
		table01_temp := Table01ForGorm{}
		if result := db.Model(&table01).Update("field03", gorm.Expr("field03 * ? + ?", 3, 10)); result.Error != nil {
			t.Fatal(result.Error)
		} else if result := db.First(&table01_temp, "field01 = ? AND field02 = ?", "a", t.Name()); result.Error != nil {
			t.Fatal(result.Error)
		} else if table01_temp.Field01 != "a" || table01_temp.Field03 != 13 {
			t.Fatal(table01_temp)
		}
	}
}

func TestUpdateSubQuery(t *testing.T) {
	prepareData(t)

	for kind, db := range dbs {
		if kind == "mysql" {
			continue
		}

		name := t.Name()
		table01 := Table01ForGorm{Field01: "a", Field02: &name}
		table01_temp := Table01ForGorm{}
		if result := db.Model(&table01).Update("field03", db.Model(&Table01ForGorm{}).Select("field03").Where("field01 = ?", "b").Where("field02 = ?", name)); result.Error != nil {
			t.Fatal(result.Error)
		} else if result := db.First(&table01_temp, "field01 = ? AND field02 = ?", "a", name); result.Error != nil {
			t.Fatal(result.Error)
		} else if table01_temp.Field01 != "a" || table01_temp.Field03 != 2 {
			t.Fatal(table01_temp)
		}
	}
}

func TestDelete(t *testing.T) {
	prepareData(t)

	for _, db := range dbs {
		name := t.Name()
		table01 := Table01ForGorm{}

		table01 = Table01ForGorm{Field01: "a", Field02: &name}
		if result := db.Delete(&table01); result.Error != nil {
			t.Fatal(result.Error)
		}

		table01 = Table01ForGorm{}
		if result := db.Where("field01 = ?", "b").Where("field02 = ?", name).Delete(&table01); result.Error != nil {
			t.Fatal(result.Error)
		}

		table01s := []*Table01ForGorm{}
		if result := db.Find(&table01s, "field02 = ?", name); result.Error != nil {
			t.Fatal(result.Error)
		} else if result.RowsAffected != 1 {
			t.Fatal(result.RowsAffected)
		}
	}
}

func TestRaw(t *testing.T) {
	prepareData(t)

	for _, db := range dbs {
		table01 := Table01ForGorm{}

		table01 = Table01ForGorm{}
		if result := db.Raw("SELECT * FROM table01_for_gorms WHERE field01 = ? AND field02 = ?", "a", t.Name()).Scan(&table01); result.Error != nil {
			t.Fatal(result.Error)
		} else if table01.Field01 != "a" || table01.Field03 != 1 {
			t.Fatal(table01)
		}

		table01 = Table01ForGorm{}
		if result := db.Raw("SELECT * FROM table01_for_gorms WHERE field01 = ? AND field02 = ?", "a", t.Name()).First(&table01); result.Error != nil {
			t.Fatal(result.Error)
		} else if table01.Field01 != "a" || table01.Field03 != 1 {
			t.Fatal(table01)
		}
	}
}

func TestRows(t *testing.T) {
	prepareData(t)

	for _, db := range dbs {
		if rows, err := db.Model(&Table01ForGorm{}).Where("field02 = ?", t.Name()).Rows(); err != nil {
			t.Fatal(err)
		} else {
			rows.Close()
		}

		if rows, err := db.Raw("SELECT * FROM table01_for_gorms WHERE field02 = ?", t.Name()).Rows(); err != nil {
			t.Fatal(err)
		} else {
			rows.Close()
		}
	}
}

func TestExec(t *testing.T) {
	prepareData(t)

	for _, db := range dbs {
		table01 := Table01ForGorm{}
		if result := db.Exec("UPDATE table01_for_gorms SET field03 = ? WHERE field01 = ? AND field02 = ?", 11, "a", t.Name()); result.Error != nil {
			t.Fatal(result.Error)
		} else if result := db.First(&table01, "field02 = ?", t.Name()); result.Error != nil {
			t.Fatal(result.Error)
		} else if table01.Field01 != "a" || table01.Field03 != 11 {
			t.Fatal(table01)
		}
	}
}

func TestNamedArgument(t *testing.T) {
	prepareData(t)

	for _, db := range dbs {
		table01 := Table01ForGorm{}

		table01 = Table01ForGorm{}
		if result := db.Where("field01 = @f1 AND field02 = @f2", sql.Named("f1", "a"), sql.Named("f2", t.Name())).Find(&table01); result.Error != nil {
			t.Fatal(result.Error)
		} else if table01.Field01 != "a" || table01.Field03 != 1 {
			t.Fatal(table01)
		}

		table01 = Table01ForGorm{}
		if result := db.Where("field01 = @f1 AND field02 = @f2", map[string]any{"f1": "a", "f2": t.Name()}).Find(&table01); result.Error != nil {
			t.Fatal(result.Error)
		} else if table01.Field01 != "a" || table01.Field03 != 1 {
			t.Fatal(table01)
		}

		table01 = Table01ForGorm{}
		if result := db.Raw("SELECT * FROM table01_for_gorms WHERE field01 = @f1 AND field02 = @f2", sql.Named("f1", "a"), sql.Named("f2", t.Name())).Scan(&table01); result.Error != nil {
			t.Fatal(result.Error)
		} else if table01.Field01 != "a" || table01.Field03 != 1 {
			t.Fatal(table01)
		}

		table01 = Table01ForGorm{}
		if result := db.Exec("UPDATE table01_for_gorms SET field03 = @f3 WHERE field01 = @f1 AND field02 = @f2", map[string]any{"f1": "a", "f2": t.Name(), "f3": 11}); result.Error != nil {
			t.Fatal(result.Error)
		} else if result := db.First(&table01, "field02 = ?", t.Name()); result.Error != nil {
			t.Fatal(result.Error)
		} else if table01.Field01 != "a" || table01.Field03 != 11 {
			t.Fatal(table01)
		}
	}
}

func TestDryRun(t *testing.T) {
	prepareData(t)

	for _, db := range dbs {
		table01 := Table01ForGorm{}
		statement := db.Session(&gorm.Session{DryRun: true}).First(&table01).Statement
		if statement.SQL.String() != `SELECT * FROM "table01_for_gorms" ORDER BY "table01_for_gorms"."field01" LIMIT $1` &&
			statement.SQL.String() != "SELECT * FROM `table01_for_gorms` ORDER BY `table01_for_gorms`.`field01` LIMIT ?" {
			t.Fatal(statement.SQL.String())
		} else if table01.Field01 != "" || table01.Field03 != 0 {
			t.Fatal(table01)
		}
	}
}

func TestToSQL(t *testing.T) {
	prepareData(t)

	for _, db := range dbs {
		table01 := Table01ForGorm{}
		sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
			return tx.First(&table01, "field02 = ?", "a")
		})
		if sql != `SELECT * FROM "table01_for_gorms" WHERE field02 = 'a' ORDER BY "table01_for_gorms"."field01" LIMIT 1` &&
			sql != "SELECT * FROM `table01_for_gorms` WHERE field02 = 'a' ORDER BY `table01_for_gorms`.`field01` LIMIT 1" {
			t.Fatal(sql)
		} else if table01.Field01 != "" || table01.Field03 != 0 {
			t.Fatal(table01)
		}
	}
}

func TestConnection(t *testing.T) {
	prepareData(t)

	for _, db := range dbs {
		err := db.Connection(func(tx *gorm.DB) error {
			table01 := Table01ForGorm{}

			for i := 0; i < 10; i++ {
				if result := tx.First(&table01, "field02 = ?", t.Name()); result.Error != nil {
					return result.Error
				} else if table01.Field01 != "a" || table01.Field03 != 1 {
					return fmt.Errorf("%#v", table01)
				}
			}

			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}
