package elasticsearch

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/thedevsaddam/gojsonq/v2"
)

var addresses []string = []string{"http://127.0.0.1:9200"}
var timeout int = 10
var index string = uuid.NewString()
var documentId string = uuid.NewString()
var template string = uuid.NewString()

func initialize(elasticsearch *Elasticsearch, addresses []string) error {
	return elasticsearch.Initialize(addresses, timeout, "", "", "", "", "", []byte(""))
}

func exists(elasticsearch *Elasticsearch) error {
	err := initialize(elasticsearch, addresses)
	if err != nil {
		return err
	}

	exist, err := elasticsearch.IndicesExists([]string{index})
	if err != nil {
		return err
	}
	if exist {
		return errors.New(fmt.Sprintf("invalid exist - exist : (%t)", exist))
	}

	exist, err = elasticsearch.Exists(index, documentId)
	if err != nil {
		return err
	}
	if exist {
		return errors.New(fmt.Sprintf("invalid exist - exist : (%t)", exist))
	}

	return nil
}

func deletes(elasticsearch *Elasticsearch) error {
	err := elasticsearch.IndicesDelete([]string{index})
	if err != nil {
		return err
	}

	exist, err := elasticsearch.IndicesExists([]string{index})
	if err != nil {
		return err
	}
	if exist {
		return errors.New(fmt.Sprintf("invalid exist - exist : (%t)", exist))
	}

	return nil
}

func existsTemplate(elasticsearch *Elasticsearch) error {
	err := initialize(elasticsearch, addresses)
	if err != nil {
		return err
	}

	exist, err := elasticsearch.IndicesExistsTemplate([]string{template})
	if err != nil {
		return err
	}
	if exist {
		return errors.New(fmt.Sprintf("invalid exist - exist : (%t)", exist))
	}

	return nil
}

func deletesTemplate(elasticsearch *Elasticsearch) error {
	err := elasticsearch.IndicesDeleteTemplate(template)
	if err != nil {
		return err
	}

	exist, err := elasticsearch.IndicesExistsTemplate([]string{template})
	if err != nil {
		return err
	}
	if exist {
		return errors.New(fmt.Sprintf("invalid exist - exist : (%t)", exist))
	}

	return nil
}

func TestInitialize(t *testing.T) {
	elasticsearch := Elasticsearch{}

	err := initialize(&elasticsearch, addresses)
	if err != nil {
		t.Error(err)
	}
}

func TestExists(t *testing.T) {
	elasticsearch := Elasticsearch{}

	_, err := elasticsearch.Exists("", "")
	if err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	err = initialize(&elasticsearch, []string{"invalid address"})
	if err != nil {
		t.Error(err)
	}
	_, err = elasticsearch.Exists(index, documentId)
	if err.Error() != `unsupported protocol scheme ""` {
		t.Error(err)
	}

	err = exists(&elasticsearch)
	if err != nil {
		t.Fatal(err)
	}

	_, err = elasticsearch.Exists("", "")
	if err.Error() != `response error - status : (400 Bad Request)` {
		t.Error(err)
	}

	_, err = elasticsearch.Exists("*", "")
	if err.Error() != `response error - status : (405 Method Not Allowed)` {
		t.Error(err)
	}

	err = elasticsearch.Index(index, documentId, `{"field":"value"}`)
	if err != nil {
		t.Error(err)
	}

	exist, err := elasticsearch.Exists(index, documentId)
	if err != nil {
		t.Error(err)
	}
	if exist == false {
		t.Errorf("invalid exist - exist : (%t)", exist)
	}

	err = deletes(&elasticsearch)
	if err != nil {
		t.Error(err)
	}
}

func TestIndex(t *testing.T) {
	elasticsearch := Elasticsearch{}

	err := elasticsearch.Index(index, documentId, "")
	if err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	err = initialize(&elasticsearch, []string{"invalid address"})
	if err != nil {
		t.Error(err)
	}
	err = elasticsearch.Index(index, documentId, `{"field":"value"}`)
	if err.Error() != `unsupported protocol scheme ""` {
		t.Error(err)
	}

	err = exists(&elasticsearch)
	if err != nil {
		t.Fatal(err)
	}

	err = elasticsearch.Index(index, documentId, "")
	if err.Error() != `response error - status : (400 Bad Request), type : (parse_exception), reason : (request body is required)` {
		t.Error(err)
	}

	err = elasticsearch.Index(index, documentId, `{"field":"value"}`)
	if err != nil {
		t.Error(err)
	}

	exist, err := elasticsearch.Exists(index, documentId)
	if err != nil {
		t.Error(err)
	}
	if exist == false {
		t.Errorf("invalid exist - exist : (%t)", exist)
	}

	err = deletes(&elasticsearch)
	if err != nil {
		t.Error(err)
	}
}

func TestDelete(t *testing.T) {
	elasticsearch := Elasticsearch{}

	err := elasticsearch.Delete(index, documentId)
	if err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	err = initialize(&elasticsearch, []string{"invalid address"})
	if err != nil {
		t.Error(err)
	}
	err = elasticsearch.Delete(index, documentId)
	if err.Error() != `unsupported protocol scheme ""` {
		t.Error(err)
	}

	err = exists(&elasticsearch)
	if err != nil {
		t.Fatal(err)
	}

	err = elasticsearch.Delete(index, documentId)
	if err.Error() != `response error - status : (404 Not Found), type : (index_not_found_exception), reason : (no such index [`+index+`])` {
		t.Error(err)
	}

	err = elasticsearch.Index(index, documentId, `{"field":"value"}`)
	if err != nil {
		t.Error(err)
	}

	err = elasticsearch.Delete(index, documentId)
	if err != nil {
		t.Error(err)
	}

	exist, err := elasticsearch.Exists(index, documentId)
	if err != nil {
		t.Error(err)
	}
	if exist {
		t.Errorf("invalid exist - exist : (%t)", exist)
	}

	err = deletes(&elasticsearch)
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteByQuery(t *testing.T) {
	elasticsearch := Elasticsearch{}

	err := elasticsearch.DeleteByQuery([]string{index}, ``)
	if err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	err = initialize(&elasticsearch, []string{"invalid address"})
	if err != nil {
		t.Error(err)
	}
	err = elasticsearch.DeleteByQuery([]string{index}, `{"query":{"match":{"field":"value_1"}}}`)
	if err.Error() != `unsupported protocol scheme ""` {
		t.Error(err)
	}

	err = exists(&elasticsearch)
	if err != nil {
		t.Fatal(err)
	}

	err = elasticsearch.DeleteByQuery([]string{index}, `{}`)
	if err.Error() != `response error - status : (400 Bad Request), type : (action_request_validation_exception), reason : (Validation Failed: 1: query is missing;)` {
		t.Error(err)
	}

	err = elasticsearch.Index(index, documentId, `{"field":"value_1"}`)
	if err != nil {
		t.Error(err)
	}

	err = elasticsearch.Index(index, documentId+"_temp", `{"field":"value_2"}`)
	if err != nil {
		t.Error(err)
	}

	err = elasticsearch.DeleteByQuery([]string{index}, `{"query":{"match":{"field":"value_1"}}}`)
	if err != nil {
		t.Error(err)
	}

	exist, err := elasticsearch.Exists(index, documentId)
	if err != nil {
		t.Error(err)
	}
	if exist {
		t.Errorf("invalid exist - exist : (%t)", exist)
	}

	err = deletes(&elasticsearch)
	if err != nil {
		t.Error(err)
	}
}

func TestIndicesExists(t *testing.T) {
	elasticsearch := Elasticsearch{}

	_, err := elasticsearch.IndicesExists([]string{index})
	if err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	err = initialize(&elasticsearch, []string{"invalid address"})
	if err != nil {
		t.Error(err)
	}
	_, err = elasticsearch.IndicesExists([]string{index})
	if err.Error() != `unsupported protocol scheme ""` {
		t.Error(err)
	}

	err = exists(&elasticsearch)
	if err != nil {
		t.Fatal(err)
	}

	_, err = elasticsearch.IndicesExists([]string{"<>"})
	if err.Error() != `response error - status : (400 Bad Request)` {
		t.Error(err)
	}

	err = elasticsearch.IndicesCreate(index, "")
	if err != nil {
		t.Error(err)
	}

	exist, err := elasticsearch.IndicesExists([]string{index})
	if err != nil {
		t.Error(err)
	}
	if exist == false {
		t.Fatalf("invalid exist - exist : (%t)", exist)
	}

	err = deletes(&elasticsearch)
	if err != nil {
		t.Error(err)
	}
}

func TestIndicesCreate(t *testing.T) {
	elasticsearch := Elasticsearch{}

	err := elasticsearch.IndicesCreate(index, "")
	if err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	err = initialize(&elasticsearch, []string{"invalid address"})
	if err != nil {
		t.Error(err)
	}
	err = elasticsearch.IndicesCreate(index, "")
	if err.Error() != `unsupported protocol scheme ""` {
		t.Error(err)
	}

	err = exists(&elasticsearch)
	if err != nil {
		t.Fatal(err)
	}

	err = elasticsearch.IndicesCreate(index, "~")
	if err.Error() != `response error - status : (500 Internal Server Error), type : (not_x_content_exception), reason : (Compressor detection can only be called on some xcontent bytes or compressed xcontent bytes)` {
		t.Error(err)
	}

	err = elasticsearch.IndicesCreate(index, "")
	if err != nil {
		t.Error(err)
	}

	exist, err := elasticsearch.IndicesExists([]string{index})
	if err != nil {
		t.Error(err)
	}
	if exist == false {
		t.Fatalf("invalid exist - exist : (%t)", exist)
	}

	err = deletes(&elasticsearch)
	if err != nil {
		t.Error(err)
	}
}

func TestIndicesDelete(t *testing.T) {
	elasticsearch := Elasticsearch{}

	err := elasticsearch.IndicesDelete([]string{""})
	if err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	err = initialize(&elasticsearch, []string{"invalid address"})
	if err != nil {
		t.Error(err)
	}
	err = elasticsearch.IndicesDelete([]string{""})
	if err.Error() != `unsupported protocol scheme ""` {
		t.Error(err)
	}

	err = exists(&elasticsearch)
	if err != nil {
		t.Fatal(err)
	}

	err = elasticsearch.IndicesDelete([]string{""})
	if err.Error() != `response error - status : (400 Bad Request), type : (action_request_validation_exception), reason : (Validation Failed: 1: index / indices is missing;)` {
		t.Fatal(err)
	}

	err = elasticsearch.IndicesCreate(index, "")
	if err != nil {
		t.Error(err)
	}

	err = deletes(&elasticsearch)
	if err != nil {
		t.Error(err)
	}
}

func TestIndicesExistsTemplate(t *testing.T) {
	elasticsearch := Elasticsearch{}

	_, err := elasticsearch.IndicesExistsTemplate([]string{template})
	if err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	err = initialize(&elasticsearch, []string{"invalid address"})
	if err != nil {
		t.Error(err)
	}
	_, err = elasticsearch.IndicesExistsTemplate([]string{template})
	if err.Error() != `unsupported protocol scheme ""` {
		t.Error(err)
	}

	err = existsTemplate(&elasticsearch)
	if err != nil {
		t.Fatal(err)
	}

	_, err = elasticsearch.IndicesExistsTemplate([]string{""})
	if err.Error() != `response error - status : (405 Method Not Allowed)` {
		t.Error(err)
	}

	err = elasticsearch.IndicesPutTemplate(template, `{"index_patterns": ["*"]}`)
	if err != nil {
		t.Error(err)
	}

	exist, err := elasticsearch.IndicesExistsTemplate([]string{template})
	if err != nil {
		t.Error(err)
	}
	if exist == false {
		t.Fatalf("invalid exist - exist : (%t)", exist)
	}

	err = deletesTemplate(&elasticsearch)
	if err != nil {
		t.Error(err)
	}
}

func TestIndicesPutTemplate(t *testing.T) {
	elasticsearch := Elasticsearch{}

	err := elasticsearch.IndicesPutTemplate(template, ``)
	if err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	err = initialize(&elasticsearch, []string{"invalid address"})
	if err != nil {
		t.Error(err)
	}
	err = elasticsearch.IndicesPutTemplate(template, `{"index_patterns": ["*"]}`)
	if err.Error() != `unsupported protocol scheme ""` {
		t.Error(err)
	}

	err = existsTemplate(&elasticsearch)
	if err != nil {
		t.Fatal(err)
	}

	err = elasticsearch.IndicesPutTemplate(template, "")
	if err.Error() != `response error - status : (400 Bad Request), type : (parse_exception), reason : (request body is required)` {
		t.Error(err)
	}

	err = elasticsearch.IndicesPutTemplate(template, `{"index_patterns": ["*"]}`)
	if err != nil {
		t.Error(err)
	}

	exist, err := elasticsearch.IndicesExistsTemplate([]string{template})
	if err != nil {
		t.Error(err)
	}
	if exist == false {
		t.Fatalf("invalid exist - exist : (%t)", exist)
	}

	err = deletesTemplate(&elasticsearch)
	if err != nil {
		t.Error(err)
	}
}

func TestIndicesDeleteTemplate(t *testing.T) {
	elasticsearch := Elasticsearch{}

	err := elasticsearch.IndicesDeleteTemplate(template)
	if err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	err = initialize(&elasticsearch, []string{"invalid address"})
	if err != nil {
		t.Error(err)
	}
	err = elasticsearch.IndicesDeleteTemplate(template)
	if err.Error() != `unsupported protocol scheme ""` {
		t.Error(err)
	}

	err = existsTemplate(&elasticsearch)
	if err != nil {
		t.Fatal(err)
	}

	err = elasticsearch.IndicesDeleteTemplate(template)
	if err.Error() != `response error - status : (404 Not Found), type : (index_template_missing_exception), reason : (index_template [`+template+`] missing)` {
		t.Error(err)
	}

	err = elasticsearch.IndicesPutTemplate(template, `{"index_patterns": ["*"]}`)
	if err != nil {
		t.Error(err)
	}

	exist, err := elasticsearch.IndicesExistsTemplate([]string{template})
	if err != nil {
		t.Error(err)
	}
	if exist == false {
		t.Fatalf("invalid exist - exist : (%t)", exist)
	}

	err = deletesTemplate(&elasticsearch)
	if err != nil {
		t.Error(err)
	}
}

func TestSearch(t *testing.T) {
	elasticsearch := Elasticsearch{}

	_, err := elasticsearch.Search(index, ``)
	if err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	err = initialize(&elasticsearch, []string{"invalid address"})
	if err != nil {
		t.Error(err)
	}
	_, err = elasticsearch.Search(index, `{"query":{"match_all":{}}}`)
	if err.Error() != `unsupported protocol scheme ""` {
		t.Error(err)
	}

	err = exists(&elasticsearch)
	if err != nil {
		t.Fatal(err)
	}

	_, err = elasticsearch.Search(index, `{}`)
	if err.Error() != `response error - status : (404 Not Found), type : (index_not_found_exception), reason : (no such index [`+index+`])` {
		t.Error(err)
	}

	err = elasticsearch.Index(index, documentId, `{"field":"value"}`)
	if err != nil {
		t.Error(err)
	}

	result, err := elasticsearch.Search(index, `{"query":{"match_all":{}}}`)
	if err != nil {
		t.Error(err)
	}
	if gojsonq.New().FromString(result).Find("hits.total.value").(float64) != 1 {
		t.Errorf("invalid result - result : (\n%s)", result)
	}

	err = deletes(&elasticsearch)
	if err != nil {
		t.Error(err)
	}
}

func TestIndicesForcemerge(t *testing.T) {
	elasticsearch := Elasticsearch{}

	err := elasticsearch.IndicesForcemerge([]string{index})
	if err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	err = initialize(&elasticsearch, []string{"invalid address"})
	if err != nil {
		t.Error(err)
	}
	err = elasticsearch.IndicesForcemerge([]string{index})
	if err.Error() != `unsupported protocol scheme ""` {
		t.Error(err)
	}

	err = exists(&elasticsearch)
	if err != nil {
		t.Fatal(err)
	}

	err = elasticsearch.IndicesForcemerge([]string{index})
	if err.Error() != `response error - status : (404 Not Found), type : (index_not_found_exception), reason : (no such index [`+index+`])` {
		t.Error(err)
	}

	err = elasticsearch.Index(index, documentId, `{"field":"value"}`)
	if err != nil {
		t.Error(err)
	}

	err = elasticsearch.IndicesForcemerge([]string{index})
	if err != nil {
		t.Error(err)
	}

	err = deletes(&elasticsearch)
	if err != nil {
		t.Error(err)
	}
}
