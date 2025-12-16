package testutil_test

import (
	"testing"

	"github.com/common-library/go/database/elasticsearch/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTestClient_UnsupportedVersion(t *testing.T) {
	t.Skip("Cannot test t.Fatalf behavior without mock testing.T")
}

func TestGetTestClient_EmptyAddresses(t *testing.T) {

	t.Skip("Skipping test that requires actual Elasticsearch connection")
}

func TestGetTestClient_V7(t *testing.T) {
	t.Skip("Skipping test that requires actual Elasticsearch connection")

	addresses := []string{"http://localhost:9200"}
	client := testutil.GetTestClient(t, "v7", addresses)

	require.NotNil(t, client)
	assert.NotNil(t, client)
}

func TestGetTestClient_V8(t *testing.T) {
	t.Skip("Skipping test that requires actual Elasticsearch connection")

	addresses := []string{"http://localhost:9200"}
	client := testutil.GetTestClient(t, "v8", addresses)

	require.NotNil(t, client)
	assert.NotNil(t, client)
}

func TestGetTestClient_V9(t *testing.T) {
	t.Skip("Skipping test that requires actual Elasticsearch connection")

	addresses := []string{"http://localhost:9200"}
	client := testutil.GetTestClient(t, "v9", addresses)

	require.NotNil(t, client)
	assert.NotNil(t, client)
}

func TestGetTestClient_MultipleAddresses(t *testing.T) {
	t.Skip("Skipping test that requires actual Elasticsearch connection")

	addresses := []string{
		"http://localhost:9200",
		"http://localhost:9201",
		"http://localhost:9202",
	}
	client := testutil.GetTestClient(t, "v9", addresses)

	require.NotNil(t, client)
	assert.NotNil(t, client)
}

func TestGetTestClient_InvalidAddress(t *testing.T) {
	t.Skip("Skipping test that requires actual Elasticsearch connection")

	addresses := []string{"invalid-address"}
	client := testutil.GetTestClient(t, "v9", addresses)

	require.NotNil(t, client)
}

func BenchmarkGetTestClient_V9(b *testing.B) {
	b.Skip("Skipping benchmark that requires actual Elasticsearch connection")

	addresses := []string{"http://localhost:9200"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = testutil.GetTestClient(&testing.T{}, "v9", addresses)
	}
}
