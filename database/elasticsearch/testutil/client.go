package testutil

import (
	"testing"
	"time"

	"github.com/common-library/go/database/elasticsearch"
	v7 "github.com/common-library/go/database/elasticsearch/v7"
	v8 "github.com/common-library/go/database/elasticsearch/v8"
	v9 "github.com/common-library/go/database/elasticsearch/v9"
	"github.com/stretchr/testify/require"
)

func GetTestClient(t *testing.T, version string, addresses []string) elasticsearch.ClientInterface {
	var client elasticsearch.ClientInterface
	var err error

	timeout := 10 * time.Second
	switch version {
	case "v7":
		client = &v7.Client{}
		err = client.Initialize(addresses, timeout, "", "", "", "", "", nil)
	case "v8":
		client = &v8.Client{}
		err = client.Initialize(addresses, timeout, "", "", "", "", "", nil)
	case "v9":
		client = &v9.Client{}
		err = client.Initialize(addresses, timeout, "", "", "", "", "", nil)
	default:
		t.Fatalf("Unsupported Elasticsearch version: %s", version)
	}

	require.NoError(t, err)
	return client
}
