package test

import (
	"os"
	"testing"

	"github.com/common-library/go/database/sql"
)

type dynamodb struct {
}

func (this *dynamodb) createDatabase() bool {
	return false
}

func (this *dynamodb) getClient(t *testing.T, databaseName string) (sql.Client, bool) {
	client := sql.Client{}

	if len(os.Getenv("DYNAMODB_URL")) == 0 {
		return client, false
	}

	dsn := "Region=dummy;AkId=dummy;SecretKey=dummy;Endpoint=" + os.Getenv("DYNAMODB_URL")

	if err := client.Open(this.getDriver(), dsn, 10); err != nil {
		t.Fatal(err)
	}

	return client, true
}

func (this *dynamodb) getDriver() sql.Driver {
	return sql.DriverAmazonDynamoDB
}

func (this *dynamodb) getPrepare(index uint64) string {
	return "?"
}

func (this *dynamodb) getCreateTableQuery(tableName string) []string {
	return []string{`CREATE TABLE ` + tableName + ` WITH pk=field01:number`}
}
