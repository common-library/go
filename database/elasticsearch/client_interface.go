// Package elasticsearch provides a unified interface for Elasticsearch clients across multiple versions.
//
// This package supports Elasticsearch versions 7, 8, and 9 through a common interface,
// allowing applications to switch between versions without changing their code.
//
// Features:
//   - Multi-version support (v7, v8, v9)
//   - Document operations (index, exists, delete)
//   - Index management (create, delete, exists)
//   - Template management
//   - Search operations
//   - Force merge support
//
// Example:
//
//	import v7 "github.com/common-library/go/database/elasticsearch/v7"
//
//	client := &v7.Client{}
//	client.Initialize([]string{"localhost:9200"}, 10*time.Second, "", "", "", "", "", nil)
//	client.Index("myindex", "doc1", `{"field":"value"}`)
package elasticsearch

import "time"

type ClientInterface interface {
	Initialize(addresses []string, timeout time.Duration, cloudID, apiKey, username, password, certificateFingerprint string, caCert []byte) error

	Exists(index, documentID string) (bool, error)

	Index(index, documentID, body string) error

	Delete(index, documentID string) error
	DeleteByQuery(indices []string, body string) error

	IndicesExists(indices []string) (bool, error)
	IndicesCreate(index, body string) error
	IndicesDelete(indices []string) error

	IndicesExistsTemplate(name []string) (bool, error)
	IndicesPutTemplate(name, body string) error
	IndicesDeleteTemplate(name string) error

	IndicesForcemerge(indices []string) error

	Search(index, body string) (string, error)
}
