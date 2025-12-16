package elasticsearch_test

import (
	"testing"

	"github.com/common-library/go/database/elasticsearch"
	v7 "github.com/common-library/go/database/elasticsearch/v7"
	v8 "github.com/common-library/go/database/elasticsearch/v8"
	v9 "github.com/common-library/go/database/elasticsearch/v9"
)

func TestClientInterface(t *testing.T) {
	t.Parallel()

	func(elasticsearch.ClientInterface) {}(&v7.Client{})
	func(elasticsearch.ClientInterface) {}(&v8.Client{})
	func(elasticsearch.ClientInterface) {}(&v9.Client{})
}
