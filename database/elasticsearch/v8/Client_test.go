package v8_test

import (
	"os"
	"strings"
	"testing"

	v8 "github.com/common-library/go/database/elasticsearch/v8"
	"github.com/google/uuid"
	"github.com/thedevsaddam/gojsonq/v2"
)

var documentId string = uuid.NewString()

func getClient(t *testing.T) (*v8.Client, bool) {
	t.Parallel()

	client := &v8.Client{}

	if len(os.Getenv("ELASTICSEARCH_ADDRESS_V8")) == 0 {
		return nil, true
	}

	if err := client.Initialize([]string{os.Getenv("ELASTICSEARCH_ADDRESS_V8")}, 10, "", "", "", "", "", []byte("")); err != nil {
		t.Fatal(err)
	}

	return client, false
}

func getIndex(t *testing.T) string {
	return strings.ToLower(t.Name()) + uuid.NewString()
}

func getTemplate(t *testing.T) string {
	return strings.ToLower(t.Name()) + uuid.NewString()
}

func indicesExists(t *testing.T, client *v8.Client, index string) {
	if exist, err := client.IndicesExists([]string{index}); err != nil {
		t.Fatal(err)
	} else if exist {
		t.Fatal(exist)
	}

	if exist, err := client.Exists(index, documentId); err != nil {
		t.Fatal(err)
	} else if exist {
		t.Fatal(exist)
	}

}

func indicesDelete(t *testing.T, client *v8.Client, index string) {
	if err := client.IndicesDelete([]string{index}); err != nil {
		t.Fatal(err)
	}

	if exist, err := client.IndicesExists([]string{index}); err != nil {
		t.Fatal(err)
	} else if exist {
		t.Fatal(exist)
	}
}

func existsTemplate(t *testing.T, client *v8.Client, template string) {
	if exist, err := client.IndicesExistsTemplate([]string{template}); err != nil {
		t.Fatal(err)
	} else if exist {
		t.Fatal(exist)
	}
}

func indicesDeleteTemplate(t *testing.T, client *v8.Client, template string) {
	if err := client.IndicesDeleteTemplate(template); err != nil {
		t.Fatal(err)
	}

	if exist, err := client.IndicesExistsTemplate([]string{template}); err != nil {
		t.Fatal(err)
	} else if exist {
		t.Fatal(exist)
	}
}

func TestInitialize(t *testing.T) {
	_, _ = getClient(t)
}

func TestExists(t *testing.T) {
	index := getIndex(t)

	if _, err := (&v8.Client{}).Exists("", ""); err.Error() != `please call Initialize first` {
		t.Fatal(err)
	}

	client, stop := getClient(t)
	if stop {
		return
	}
	defer indicesDelete(t, client, index)

	indicesExists(t, client, index)

	if _, err := client.Exists("", ""); err.Error() != `response error - status : (405 Method Not Allowed)` {
		t.Fatal(err)
	}

	if _, err := client.Exists("*", ""); err.Error() != `response error - status : (405 Method Not Allowed)` {
		t.Fatal(err)
	}

	if err := client.Index(index, documentId, `{"field":"value"}`); err != nil {
		t.Fatal(err)
	}

	if exist, err := client.Exists(index, documentId); err != nil {
		t.Fatal(err)
	} else if exist == false {
		t.Fatal(exist)
	}
}

func TestIndex(t *testing.T) {
	index := getIndex(t)

	if err := (&v8.Client{}).Index(index, documentId, ""); err.Error() != `please call Initialize first` {
		t.Fatal(err)
	}

	client, stop := getClient(t)
	if stop {
		return
	}
	defer indicesDelete(t, client, index)

	indicesExists(t, client, index)

	if err := client.Index(index, documentId, ""); err.Error() != `response error - status : (400 Bad Request), type : (parse_exception), reason : (request body is required)` {
		t.Fatal(err)
	}

	if err := client.Index(index, documentId, `{"field":"value"}`); err != nil {
		t.Fatal(err)
	}

	if exist, err := client.Exists(index, documentId); err != nil {
		t.Fatal(err)
	} else if exist == false {
		t.Fatal(exist)
	}
}

func TestDelete(t *testing.T) {
	index := getIndex(t)

	if err := (&v8.Client{}).Delete(index, documentId); err.Error() != `please call Initialize first` {
		t.Fatal(err)
	}

	client, stop := getClient(t)
	if stop {
		return
	}
	defer indicesDelete(t, client, index)

	indicesExists(t, client, index)

	if err := client.Delete(index, documentId); err.Error() != `response error - status : (404 Not Found), type : (index_not_found_exception), reason : (no such index [`+index+`])` {
		t.Fatal(err)
	}

	if err := client.Index(index, documentId, `{"field":"value"}`); err != nil {
		t.Fatal(err)
	}

	if err := client.Delete(index, documentId); err != nil {
		t.Fatal(err)
	}

	if exist, err := client.Exists(index, documentId); err != nil {
		t.Fatal(err)
	} else if exist {
		t.Fatal(exist)
	}
}

func TestDeleteByQuery(t *testing.T) {
	index := getIndex(t)

	if err := (&v8.Client{}).DeleteByQuery([]string{index}, ``); err.Error() != `please call Initialize first` {
		t.Fatal(err)
	}

	client, stop := getClient(t)
	if stop {
		return
	}
	defer indicesDelete(t, client, index)

	indicesExists(t, client, index)

	if err := client.DeleteByQuery([]string{index}, `{}`); err.Error() != `response error - status : (400 Bad Request), type : (action_request_validation_exception), reason : (Validation Failed: 1: query is missing;)` {
		t.Fatal(err)
	}

	if err := client.Index(index, documentId, `{"field":"value_1"}`); err != nil {
		t.Fatal(err)
	}

	if err := client.Index(index, documentId+"_temp", `{"field":"value_2"}`); err != nil {
		t.Fatal(err)
	}

	if err := client.DeleteByQuery([]string{index}, `{"query":{"match":{"field":"value_1"}}}`); err != nil {
		t.Fatal(err)
	}

	if exist, err := client.Exists(index, documentId); err != nil {
		t.Fatal(err)
	} else if exist {
		t.Fatal(exist)
	}
}

func TestIndicesExists(t *testing.T) {
	index := getIndex(t)

	if _, err := (&v8.Client{}).IndicesExists([]string{index}); err.Error() != `please call Initialize first` {
		t.Fatal(err)
	}

	client, stop := getClient(t)
	if stop {
		return
	}
	defer indicesDelete(t, client, index)

	indicesExists(t, client, index)

	if _, err := client.IndicesExists([]string{"<>"}); err.Error() != `response error - status : (400 Bad Request)` {
		t.Fatal(err)
	}

	if err := client.IndicesCreate(index, ""); err != nil {
		t.Fatal(err)
	}

	if exist, err := client.IndicesExists([]string{index}); err != nil {
		t.Fatal(err)
	} else if exist == false {
		t.Fatal(exist)
	}
}

func TestIndicesCreate(t *testing.T) {
	index := getIndex(t)

	if err := (&v8.Client{}).IndicesCreate(index, ""); err.Error() != `please call Initialize first` {
		t.Fatal(err)
	}

	client, stop := getClient(t)
	if stop {
		return
	}
	defer indicesDelete(t, client, index)

	indicesExists(t, client, index)

	if err := client.IndicesCreate(index, "~"); err.Error() != `response error - status : (500 Internal Server Error), type : (not_x_content_exception), reason : (Compressor detection can only be called on some xcontent bytes or compressed xcontent bytes)` {
		t.Fatal(err)
	}

	if err := client.IndicesCreate(index, ""); err != nil {
		t.Fatal(err)
	}

	if exist, err := client.IndicesExists([]string{index}); err != nil {
		t.Fatal(err)
	} else if exist == false {
		t.Fatal(exist)
	}
}

func TestIndicesDelete(t *testing.T) {
	index := getIndex(t)

	if err := (&v8.Client{}).IndicesDelete([]string{""}); err.Error() != `please call Initialize first` {
		t.Fatal(err)
	}

	client, stop := getClient(t)
	if stop {
		return
	}
	defer indicesDelete(t, client, index)

	indicesExists(t, client, index)

	if err := client.IndicesDelete([]string{""}); err.Error() != `response error - status : (400 Bad Request), type : (action_request_validation_exception), reason : (Validation Failed: 1: index / indices is missing;)` {
		t.Fatal(err)
	}

	if err := client.IndicesCreate(index, ""); err != nil {
		t.Fatal(err)
	}
}

func TestIndicesExistsTemplate(t *testing.T) {
	template := getTemplate(t)

	if _, err := (&v8.Client{}).IndicesExistsTemplate([]string{template}); err.Error() != `please call Initialize first` {
		t.Fatal(err)
	}

	client, stop := getClient(t)
	if stop {
		return
	}
	defer indicesDeleteTemplate(t, client, template)

	existsTemplate(t, client, template)

	if _, err := client.IndicesExistsTemplate([]string{""}); err.Error() != `response error - status : (405 Method Not Allowed)` {
		t.Fatal(err)
	}

	if err := client.IndicesPutTemplate(template, `{"index_patterns": ["*"]}`); err != nil {
		t.Fatal(err)
	}

	if exist, err := client.IndicesExistsTemplate([]string{template}); err != nil {
		t.Fatal(err)
	} else if exist == false {
		t.Fatal(exist)
	}
}

func TestIndicesPutTemplate(t *testing.T) {
	template := getTemplate(t)

	if err := (&v8.Client{}).IndicesPutTemplate(template, ``); err.Error() != `please call Initialize first` {
		t.Fatal(err)
	}

	client, stop := getClient(t)
	if stop {
		return
	}
	defer indicesDeleteTemplate(t, client, template)

	existsTemplate(t, client, template)

	if err := client.IndicesPutTemplate(template, ""); err.Error() != `response error - status : (400 Bad Request), type : (parse_exception), reason : (request body is required)` {
		t.Fatal(err)
	}

	if err := client.IndicesPutTemplate(template, `{"index_patterns": ["*"]}`); err != nil {
		t.Fatal(err)
	}

	if exist, err := client.IndicesExistsTemplate([]string{template}); err != nil {
		t.Fatal(err)
	} else if exist == false {
		t.Fatal(exist)
	}
}

func TestIndicesDeleteTemplate(t *testing.T) {
	template := getTemplate(t)

	if err := (&v8.Client{}).IndicesDeleteTemplate(template); err.Error() != `please call Initialize first` {
		t.Fatal(err)
	}

	client, stop := getClient(t)
	if stop {
		return
	}
	defer indicesDeleteTemplate(t, client, template)

	existsTemplate(t, client, template)

	if err := client.IndicesDeleteTemplate(template); err.Error() != `response error - status : (404 Not Found), type : (index_template_missing_exception), reason : (index_template [`+template+`] missing)` {
		t.Fatal(err)
	}

	if err := client.IndicesPutTemplate(template, `{"index_patterns": ["*"]}`); err != nil {
		t.Fatal(err)
	}

	if exist, err := client.IndicesExistsTemplate([]string{template}); err != nil {
		t.Fatal(err)
	} else if exist == false {
		t.Fatal(exist)
	}
}

func TestSearch(t *testing.T) {
	index := getIndex(t)

	if _, err := (&v8.Client{}).Search(index, ``); err.Error() != `please call Initialize first` {
		t.Fatal(err)
	}

	client, stop := getClient(t)
	if stop {
		return
	}
	defer indicesDelete(t, client, index)

	indicesExists(t, client, index)

	if _, err := client.Search(index, `{}`); err.Error() != `response error - status : (404 Not Found), type : (index_not_found_exception), reason : (no such index [`+index+`])` {
		t.Fatal(err)
	}

	if err := client.Index(index, documentId, `{"field":"value"}`); err != nil {
		t.Fatal(err)
	}

	if result, err := client.Search(index, `{"query":{"match_all":{}}}`); err != nil {
		t.Fatal(err)
	} else if gojsonq.New().FromString(result).Find("hits.total.value").(float64) != 1 {
		t.Fatal(result)
	}
}

func TestIndicesForcemerge(t *testing.T) {
	index := getIndex(t)

	if err := (&v8.Client{}).IndicesForcemerge([]string{index}); err.Error() != `please call Initialize first` {
		t.Fatal(err)
	}

	client, stop := getClient(t)
	if stop {
		return
	}
	defer indicesDelete(t, client, index)

	indicesExists(t, client, index)

	if err := client.IndicesForcemerge([]string{index}); err.Error() != `response error - status : (404 Not Found), type : (index_not_found_exception), reason : (no such index [`+index+`])` {
		t.Fatal(err)
	}

	if err := client.Index(index, documentId, `{"field":"value"}`); err != nil {
		t.Fatal(err)
	}

	if err := client.IndicesForcemerge([]string{index}); err != nil {
		t.Fatal(err)
	}
}
