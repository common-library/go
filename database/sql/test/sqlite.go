package test

import (
	"testing"

	"github.com/common-library/go/database/sql"
)

type sqlite struct {
}

func (this *sqlite) createDatabase() bool {
	return false
}

func (this *sqlite) getClient(t *testing.T, databaseName string) (sql.Client, bool) {
	client := sql.Client{}

	if err := client.Open(this.getDriver(), databaseName, 10); err != nil {
		t.Fatal(err)
	}

	return client, true
}

func (this *sqlite) getDriver() sql.Driver {
	return sql.DriverSQLite
}

func (this *sqlite) getPrepare(index uint64) string {
	return "?"
}

func (this *sqlite) getCreateTableQuery(tableName string) []string {
	return []string{`CREATE TABLE ` + tableName + `(field01 INT, field02 VARCHAR(255));`}
}
