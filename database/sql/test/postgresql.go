package test

import (
	"os"
	"strconv"
	"testing"

	"github.com/common-library/go/database/sql"
)

type postgresql struct {
	prepare int
}

func (this *postgresql) createDatabase() bool {
	return true
}

func (this *postgresql) getClient(t *testing.T, databaseName string) (sql.Client, bool) {
	client := sql.Client{}

	if len(os.Getenv("POSTGRESQL_DSN")) == 0 {
		return client, false
	}

	dsn := os.Getenv("POSTGRESQL_DSN") + " dbname=" + databaseName
	if err := client.Open(this.getDriver(), dsn, 10); err != nil {
		t.Fatal(err)
	}

	return client, true
}

func (this *postgresql) getDriver() sql.Driver {
	return sql.DriverPostgreSQL
}

func (this *postgresql) getPrepare(index uint64) string {
	return "$" + strconv.FormatUint(index, 10)
}

func (this *postgresql) getCreateTableQuery(tableName string) []string {
	return []string{`CREATE TABLE ` + tableName + `(field01 INT, field02 VARCHAR(255));`}
}
