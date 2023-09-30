package elasticsearch

type ClientInterface interface {
	Initialize(addresses []string, timeout uint64, cloudID, apiKey, username, password, certificateFingerprint string, caCert []byte) error

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
