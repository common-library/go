package v7_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/common-library/go/database/elasticsearch/testutil"
	v7 "github.com/common-library/go/database/elasticsearch/v7"
	imgtestutil "github.com/common-library/go/testutil"
)

var (
	elasticsearchURL string
)

func TestMain(m *testing.M) {
	testutil.RunWithElasticsearch(m, imgtestutil.ElasticsearchV7Image, &elasticsearchURL)
}

func TestClient_Initialize(t *testing.T) {
	client := &v7.Client{}
	err := client.Initialize(
		[]string{elasticsearchURL},
		30*time.Second,
		"", "", "", "", "", nil,
	)
	assert.NoError(t, err)
}

func TestClient_IndicesCreate(t *testing.T) {
	client := testutil.GetTestClient(t, "v7", []string{elasticsearchURL})

	indexName := "test-index-create"

	exists, err := client.IndicesExists([]string{indexName})
	require.NoError(t, err)
	assert.False(t, exists)

	mapping := `{
		"mappings": {
			"properties": {
				"title": {"type": "text"},
				"content": {"type": "text"},
				"timestamp": {"type": "date"}
			}
		}
	}`

	err = client.IndicesCreate(indexName, mapping)
	assert.NoError(t, err)

	exists, err = client.IndicesExists([]string{indexName})
	require.NoError(t, err)
	assert.True(t, exists)

	err = client.IndicesDelete([]string{indexName})
	assert.NoError(t, err)
}

func TestClient_Index_and_Exists(t *testing.T) {
	client := testutil.GetTestClient(t, "v7", []string{elasticsearchURL})

	indexName := "test-index-document"
	documentID := "test-doc-1"

	mapping := `{
		"mappings": {
			"properties": {
				"title": {"type": "text"},
				"content": {"type": "text"}
			}
		}
	}`

	err := client.IndicesCreate(indexName, mapping)
	require.NoError(t, err)

	exists, err := client.Exists(indexName, documentID)
	require.NoError(t, err)
	assert.False(t, exists)

	document := `{
		"title": "Test Document",
		"content": "This is a test document for Elasticsearch testing"
	}`

	err = client.Index(indexName, documentID, document)
	assert.NoError(t, err)

	exists, err = client.Exists(indexName, documentID)
	require.NoError(t, err)
	assert.True(t, exists)

	err = client.IndicesDelete([]string{indexName})
	assert.NoError(t, err)
}

func TestClient_Search(t *testing.T) {
	client := testutil.GetTestClient(t, "v7", []string{elasticsearchURL})

	indexName := "test-index-search"

	mapping := `{
		"mappings": {
			"properties": {
				"title": {"type": "text"},
				"content": {"type": "text"},
				"category": {"type": "keyword"}
			}
		}
	}`

	err := client.IndicesCreate(indexName, mapping)
	require.NoError(t, err)

	documents := []struct {
		id   string
		data string
	}{
		{
			id: "doc1",
			data: `{
				"title": "Elasticsearch Guide",
				"content": "Learn how to use Elasticsearch effectively",
				"category": "tutorial"
			}`,
		},
		{
			id: "doc2",
			data: `{
				"title": "Search Optimization",
				"content": "Tips for optimizing search performance",
				"category": "guide"
			}`,
		},
	}

	for _, doc := range documents {
		err = client.Index(indexName, doc.id, doc.data)
		require.NoError(t, err)
	}

	searchQuery := `{
		"query": {
			"match": {
				"content": "search"
			}
		}
	}`

	result, err := client.Search(indexName, searchQuery)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "hits")

	err = client.IndicesDelete([]string{indexName})
	assert.NoError(t, err)
}

func TestClient_Delete(t *testing.T) {
	client := testutil.GetTestClient(t, "v7", []string{elasticsearchURL})

	indexName := "test-index-delete"
	documentID := "test-doc-delete"

	mapping := `{
		"mappings": {
			"properties": {
				"title": {"type": "text"}
			}
		}
	}`

	err := client.IndicesCreate(indexName, mapping)
	require.NoError(t, err)

	document := `{"title": "Document to be deleted"}`
	err = client.Index(indexName, documentID, document)
	require.NoError(t, err)

	exists, err := client.Exists(indexName, documentID)
	require.NoError(t, err)
	assert.True(t, exists)

	err = client.Delete(indexName, documentID)
	assert.NoError(t, err)

	exists, err = client.Exists(indexName, documentID)
	require.NoError(t, err)
	assert.False(t, exists)

	err = client.IndicesDelete([]string{indexName})
	assert.NoError(t, err)
}

func TestClient_DeleteByQuery(t *testing.T) {
	client := testutil.GetTestClient(t, "v7", []string{elasticsearchURL})

	indexName := "test-index-delete-query"

	mapping := `{
		"mappings": {
			"properties": {
				"title": {"type": "text"},
				"category": {"type": "keyword"}
			}
		}
	}`

	err := client.IndicesCreate(indexName, mapping)
	require.NoError(t, err)

	documents := []struct {
		id   string
		data string
	}{
		{
			id:   "doc1",
			data: `{"title": "Test Document 1", "category": "test"}`,
		},
		{
			id:   "doc2",
			data: `{"title": "Test Document 2", "category": "test"}`,
		},
		{
			id:   "doc3",
			data: `{"title": "Production Document", "category": "prod"}`,
		},
	}

	for _, doc := range documents {
		err = client.Index(indexName, doc.id, doc.data)
		require.NoError(t, err)
	}

	deleteQuery := `{
		"query": {
			"term": {
				"category": "test"
			}
		}
	}`

	err = client.DeleteByQuery([]string{indexName}, deleteQuery)
	assert.NoError(t, err)

	searchQuery := `{
		"query": {
			"match_all": {}
		}
	}`

	result, err := client.Search(indexName, searchQuery)
	assert.NoError(t, err)
	assert.Contains(t, result, "prod")

	err = client.IndicesDelete([]string{indexName})
	assert.NoError(t, err)
}

func TestClient_IndicesTemplate(t *testing.T) {
	client := testutil.GetTestClient(t, "v7", []string{elasticsearchURL})

	templateName := "test-template"

	exists, err := client.IndicesExistsTemplate([]string{templateName})
	require.NoError(t, err)
	assert.False(t, exists)

	template := `{
		"index_patterns": ["test-template-*"],
		"mappings": {
			"properties": {
				"timestamp": {"type": "date"},
				"message": {"type": "text"}
			}
		}
	}`

	err = client.IndicesPutTemplate(templateName, template)
	assert.NoError(t, err)

	exists, err = client.IndicesExistsTemplate([]string{templateName})
	require.NoError(t, err)
	assert.True(t, exists)

	err = client.IndicesDeleteTemplate(templateName)
	assert.NoError(t, err)

	exists, err = client.IndicesExistsTemplate([]string{templateName})
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestClient_IndicesForcemerge(t *testing.T) {
	client := testutil.GetTestClient(t, "v7", []string{elasticsearchURL})

	indexName := "test-index-forcemerge"

	mapping := `{
		"mappings": {
			"properties": {
				"data": {"type": "text"}
			}
		}
	}`

	err := client.IndicesCreate(indexName, mapping)
	require.NoError(t, err)

	for i := 0; i < 2; i++ {
		document := fmt.Sprintf(`{"data": "test document %d"}`, i)
		err = client.Index(indexName, fmt.Sprintf("doc%d", i), document)
		require.NoError(t, err)
	}

	err = client.IndicesForcemerge([]string{indexName})
	assert.NoError(t, err)

	err = client.IndicesDelete([]string{indexName})
	assert.NoError(t, err)
}
