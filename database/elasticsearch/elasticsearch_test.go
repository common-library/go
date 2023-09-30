package elasticsearch_test

import (
	"testing"

	"github.com/heaven-chp/common-library-go/database/elasticsearch"
	v7 "github.com/heaven-chp/common-library-go/database/elasticsearch/v7"
	v8 "github.com/heaven-chp/common-library-go/database/elasticsearch/v8"
)

func TestInterface(t *testing.T) {
	func(elasticsearch.Elasticsearch) {
	}(&v7.Elasticsearch{})

	func(elasticsearch.Elasticsearch) {
	}(&v8.Elasticsearch{})
}
