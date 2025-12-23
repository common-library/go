package v9_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/common-library/go/database/elasticsearch/testutil"
	imgtestutil "github.com/common-library/go/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	elasticsearchURL string
)

func TestMain(m *testing.M) {
	testutil.RunWithElasticsearch(m, imgtestutil.ElasticsearchV9Image, &elasticsearchURL)
}

func generateUniqueIndexName(prefix string) string {
	return fmt.Sprintf("%s-%d-%d", prefix, time.Now().UnixNano(), rand.Int31())
}

func TestClient_IndicesCreate(t *testing.T) {
	t.Parallel()
	client := testutil.GetTestClient(t, "v9", []string{elasticsearchURL})

	indexName := generateUniqueIndexName("test-index-create")
	indexBody := `{
		"settings": {
			"number_of_shards": 1,
			"number_of_replicas": 0
		},
		"mappings": {
			"properties": {
				"title": {
					"type": "text"
				},
				"content": {
					"type": "text"
				}
			}
		}
	}`

	err := client.IndicesCreate(indexName, indexBody)
	assert.NoError(t, err, "Failed to create index")

	exists, err := client.IndicesExists([]string{indexName})
	assert.NoError(t, err, "Failed to check if index exists")
	assert.True(t, exists, "Index should exist after creation")

	err = client.IndicesDelete([]string{indexName})
	assert.NoError(t, err, "Failed to delete index")
}

func TestClient_Index(t *testing.T) {
	t.Parallel()
	client := testutil.GetTestClient(t, "v9", []string{elasticsearchURL})

	indexName := generateUniqueIndexName("test-index-document")
	documentID := "test-doc-1"
	documentBody := `{
		"title": "Test Document",
		"content": "This is a test document for indexing",
		"timestamp": "2025-01-09T10:00:00Z"
	}`

	indexMapping := `{
		"settings": {
			"number_of_shards": 1,
			"number_of_replicas": 0
		}
	}`
	err := client.IndicesCreate(indexName, indexMapping)
	require.NoError(t, err, "Failed to create index for document test")

	err = client.Index(indexName, documentID, documentBody)
	assert.NoError(t, err, "Failed to index document")

	exists, err := client.Exists(indexName, documentID)
	assert.NoError(t, err, "Failed to check if document exists")
	assert.True(t, exists, "Document should exist after indexing")

	err = client.IndicesDelete([]string{indexName})
	assert.NoError(t, err, "Failed to delete index")
}

func TestClient_Search(t *testing.T) {
	t.Parallel()
	client := testutil.GetTestClient(t, "v9", []string{elasticsearchURL})

	indexName := generateUniqueIndexName("test-index-search")

	indexMapping := `{
		"settings": {
			"number_of_shards": 1,
			"number_of_replicas": 0
		},
		"mappings": {
			"properties": {
				"title": {
					"type": "text"
				},
				"content": {
					"type": "text"
				}
			}
		}
	}`
	err := client.IndicesCreate(indexName, indexMapping)
	require.NoError(t, err, "Failed to create index for search test")

	documents := []struct {
		id   string
		body string
	}{
		{
			id: "doc1",
			body: `{
				"title": "First Document",
				"content": "This is the first test document"
			}`,
		},
		{
			id: "doc2",
			body: `{
				"title": "Second Document",
				"content": "This is the second test document"
			}`,
		},
	}

	for _, doc := range documents {
		err := client.Index(indexName, doc.id, doc.body)
		require.NoError(t, err, "Failed to index document %s", doc.id)
	}

	searchQuery := `{
		"query": {
			"match": {
				"title": "Document"
			}
		}
	}`

	result, err := client.Search(indexName, searchQuery)
	assert.NoError(t, err, "Failed to search documents")
	assert.NotEmpty(t, result, "Search result should not be empty")
	assert.Contains(t, result, "hits", "Search result should contain hits")

	err = client.IndicesDelete([]string{indexName})
	assert.NoError(t, err, "Failed to delete index")
}

func TestClient_Delete(t *testing.T) {
	t.Parallel()
	client := testutil.GetTestClient(t, "v9", []string{elasticsearchURL})

	indexName := generateUniqueIndexName("test-index-delete")
	documentID := "test-doc-delete"
	documentBody := `{
		"title": "Document to Delete",
		"content": "This document will be deleted"
	}`

	indexMapping := `{
		"settings": {
			"number_of_shards": 1,
			"number_of_replicas": 0
		}
	}`
	err := client.IndicesCreate(indexName, indexMapping)
	require.NoError(t, err, "Failed to create index for delete test")

	err = client.Index(indexName, documentID, documentBody)
	require.NoError(t, err, "Failed to index document for delete test")

	exists, err := client.Exists(indexName, documentID)
	require.NoError(t, err, "Failed to check if document exists before delete")
	require.True(t, exists, "Document should exist before deletion")

	err = client.Delete(indexName, documentID)
	assert.NoError(t, err, "Failed to delete document")

	exists, err = client.Exists(indexName, documentID)
	assert.NoError(t, err, "Failed to check if document exists after delete")
	assert.False(t, exists, "Document should not exist after deletion")

	err = client.IndicesDelete([]string{indexName})
	assert.NoError(t, err, "Failed to delete index")
}

func TestClient_DeleteByQuery(t *testing.T) {
	t.Parallel()
	client := testutil.GetTestClient(t, "v9", []string{elasticsearchURL})

	indexName := generateUniqueIndexName("test-index-delete-by-query")

	indexMapping := `{
		"settings": {
			"number_of_shards": 1,
			"number_of_replicas": 0
		},
		"mappings": {
			"properties": {
				"status": {
					"type": "keyword"
				},
				"title": {
					"type": "text"
				}
			}
		}
	}`
	err := client.IndicesCreate(indexName, indexMapping)
	require.NoError(t, err, "Failed to create index for delete by query test")

	documents := []struct {
		id   string
		body string
	}{
		{
			id: "doc1",
			body: `{
				"title": "Active Document",
				"status": "active"
			}`,
		},
		{
			id: "doc2",
			body: `{
				"title": "Inactive Document",
				"status": "inactive"
			}`,
		},
		{
			id: "doc3",
			body: `{
				"title": "Another Inactive Document",
				"status": "inactive"
			}`,
		},
	}

	for _, doc := range documents {
		err := client.Index(indexName, doc.id, doc.body)
		require.NoError(t, err, "Failed to index document %s", doc.id)
	}

	deleteQuery := `{
		"query": {
			"term": {
				"status": "inactive"
			}
		}
	}`

	err = client.DeleteByQuery([]string{indexName}, deleteQuery)
	assert.NoError(t, err, "Failed to delete documents by query")

	searchQuery := `{
		"query": {
			"match_all": {}
		}
	}`

	result, err := client.Search(indexName, searchQuery)
	assert.NoError(t, err, "Failed to search documents after delete by query")
	assert.Contains(t, result, "active", "Only active documents should remain")
	assert.NotContains(t, result, "inactive", "Inactive documents should be deleted")

	err = client.IndicesDelete([]string{indexName})
	assert.NoError(t, err, "Failed to delete index")
}

func TestClient_Template(t *testing.T) {
	t.Parallel()
	client := testutil.GetTestClient(t, "v9", []string{elasticsearchURL})

	templateName := generateUniqueIndexName("test-template")
	templateBody := `{
		"index_patterns": ["test-template-*"],
		"settings": {
			"number_of_shards": 1,
			"number_of_replicas": 0
		},
		"mappings": {
			"properties": {
				"timestamp": {
					"type": "date"
				},
				"message": {
					"type": "text"
				}
			}
		}
	}`

	err := client.IndicesPutTemplate(templateName, templateBody)
	assert.NoError(t, err, "Failed to create template")

	exists, err := client.IndicesExistsTemplate([]string{templateName})
	assert.NoError(t, err, "Failed to check if template exists")
	assert.True(t, exists, "Template should exist after creation")

	err = client.IndicesDeleteTemplate(templateName)
	assert.NoError(t, err, "Failed to delete template")

	exists, err = client.IndicesExistsTemplate([]string{templateName})
	assert.NoError(t, err, "Failed to check if template exists after delete")
	assert.False(t, exists, "Template should not exist after deletion")
}

func TestClient_Forcemerge(t *testing.T) {
	t.Parallel()
	client := testutil.GetTestClient(t, "v9", []string{elasticsearchURL})

	indexName := generateUniqueIndexName("test-index-forcemerge")

	indexMapping := `{
		"settings": {
			"number_of_shards": 1,
			"number_of_replicas": 0
		}
	}`
	err := client.IndicesCreate(indexName, indexMapping)
	require.NoError(t, err, "Failed to create index for forcemerge test")

	for i := 0; i < 2; i++ {
		documentBody := fmt.Sprintf(`{
			"id": %d,
			"message": "Test document %d"
		}`, i, i)

		err := client.Index(indexName, fmt.Sprintf("doc-%d", i), documentBody)
		require.NoError(t, err, "Failed to index document %d", i)
	}

	err = client.IndicesForcemerge([]string{indexName})
	assert.NoError(t, err, "Failed to perform forcemerge")

	err = client.IndicesDelete([]string{indexName})
	assert.NoError(t, err, "Failed to delete index")
}
