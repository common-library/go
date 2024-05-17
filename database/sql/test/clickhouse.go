package test

import (
	"os"
	"testing"

	"github.com/common-library/go/database/sql"
)

type clickhouse struct {
}

func (this *clickhouse) createDatabase() bool {
	return true
}

func (this *clickhouse) getClient(t *testing.T, databaseName string) (sql.Client, bool) {
	client := sql.Client{}

	if len(os.Getenv("CLICKHOUSE_DSN")) == 0 {
		return client, false
	}

	dsn := os.Getenv("CLICKHOUSE_DSN")
	if len(databaseName) != 0 {
		dsn += "/" + databaseName
	}

	if err := client.Open(this.getDriver(), dsn, 10); err != nil {
		t.Fatal(err)
	}

	return client, true
}

func (this *clickhouse) getDriver() sql.Driver {
	return sql.DriverClickHouse
}

func (this *clickhouse) getPrepare(index uint64) string {
	return ""
}

func (this *clickhouse) getCreateTableQuery(tableName string) []string {
	return []string{
		`CREATE TABLE ` + tableName + `_
(
    field01 UInt32,
    field02 String
)
ENGINE = MergeTree
PRIMARY KEY (field01)`,
		`CREATE TABLE ` + tableName + `
(
    field01 UInt64,
    field02 String
)
ENGINE = Distributed(default, ` + tableName + `, ` + tableName + `_, rand())`,
	}
}
