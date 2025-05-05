package test

import (
	"os"
	"strings"
	"testing"

	"github.com/common-library/go/database/sql"
)

type mysql struct {
}

func (this *mysql) createDatabase() bool {
	return true
}

func (this *mysql) getClient(t *testing.T, databaseName string) (sql.Client, bool) {
	client := sql.Client{}

	if len(os.Getenv("MYSQL_DSN")) == 0 {
		return client, false
	}

	dsn := strings.Replace(os.Getenv("MYSQL_DSN"), "${database}", databaseName, 1)
	if err := client.Open(this.getDriver(), dsn, 10); err != nil {
		t.Fatal(err)
	}

	return client, true
}

func (this *mysql) getDriver() sql.Driver {
	return sql.DriverMySQL
}

func (this *mysql) getPrepare(index uint64) string {
	return "?"
}

func (this *mysql) getCreateTableQuery(tableName string) []string {
	return []string{`CREATE TABLE ` + tableName + `(field01 INT, field02 VARCHAR(255));`}
}
