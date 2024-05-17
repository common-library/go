package v8_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	v8 "github.com/common-library/go/database/elasticsearch/v8"
	"github.com/google/uuid"
	"github.com/thedevsaddam/gojsonq/v2"
)

var index string = uuid.NewString()
var documentId string = uuid.NewString()
var template string = uuid.NewString()

func getClient(t *testing.T) (v8.Client, bool) {
	client := v8.Client{}

	if len(os.Getenv("ELASTICSEARCH_ADDRESS_V8")) == 0 {
		return client, false
	}

	if err := client.Initialize([]string{os.Getenv("ELASTICSEARCH_ADDRESS_V8")}, 10, "", "", "", "", "", []byte("")); err != nil {
		t.Fatal(err)
	}

	return client, true
}

func indicesExists(client *v8.Client) error {
	if exist, err := client.IndicesExists([]string{index}); err != nil {
		return err
	} else if exist {
		return errors.New(fmt.Sprintf("invalid exist - exist : (%t)", exist))
	}

	if exist, err := client.Exists(index, documentId); err != nil {
		return err
	} else if exist {
		return errors.New(fmt.Sprintf("invalid exist - exist : (%t)", exist))
	}

	return nil
}

func indicesDelete(client *v8.Client) error {
	if err := client.IndicesDelete([]string{index}); err != nil {
		return err
	}

	if exist, err := client.IndicesExists([]string{index}); err != nil {
		return err
	} else if exist {
		return errors.New(fmt.Sprintf("invalid exist - exist : (%t)", exist))
	}

	return nil
}

func existsTemplate(client *v8.Client) error {
	if exist, err := client.IndicesExistsTemplate([]string{template}); err != nil {
		return err
	} else if exist {
		return errors.New(fmt.Sprintf("invalid exist - exist : (%t)", exist))
	}

	return nil
}

func indicesDeleteTemplate(client *v8.Client) error {
	if err := client.IndicesDeleteTemplate(template); err != nil {
		return err
	}

	if exist, err := client.IndicesExistsTemplate([]string{template}); err != nil {
		return err
	} else if exist {
		return errors.New(fmt.Sprintf("invalid exist - exist : (%t)", exist))
	}

	return nil
}

func TestInitialize(t *testing.T) {
	_, _ = getClient(t)
}

func TestExists(t *testing.T) {
	client := v8.Client{}
	if _, err := client.Exists("", ""); err.Error() != `please call Initialize first` {
		t.Fatal(err)
	}

	client, ok := getClient(t)
	if ok == false {
		return
	}

	if err := indicesExists(&client); err != nil {
		t.Fatal(err)
	}

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
		t.Fatal("invalid exist")
	}

	if err := indicesDelete(&client); err != nil {
		t.Fatal(err)
	}
}

func TestIndex(t *testing.T) {
	client := v8.Client{}
	if err := client.Index(index, documentId, ""); err.Error() != `please call Initialize first` {
		t.Fatal(err)
	}

	client, ok := getClient(t)
	if ok == false {
		return
	}

	if err := indicesExists(&client); err != nil {
		t.Fatal(err)
	}

	if err := client.Index(index, documentId, ""); err.Error() != `response error - status : (400 Bad Request), type : (parse_exception), reason : (request body is required)` {
		t.Fatal(err)
	}

	if err := client.Index(index, documentId, `{"field":"value"}`); err != nil {
		t.Fatal(err)
	}

	if exist, err := client.Exists(index, documentId); err != nil {
		t.Fatal(err)
	} else if exist == false {
		t.Fatal("invalid exist")
	}

	if err := indicesDelete(&client); err != nil {
		t.Fatal(err)
	}
}

func TestDelete(t *testing.T) {
	client := v8.Client{}
	if err := client.Delete(index, documentId); err.Error() != `please call Initialize first` {
		t.Fatal(err)
	}

	client, ok := getClient(t)
	if ok == false {
		return
	}

	if err := indicesExists(&client); err != nil {
		t.Fatal(err)
	}

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
		t.Fatal("invalid exist")
	}

	if err := indicesDelete(&client); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteByQuery(t *testing.T) {
	client := v8.Client{}
	if err := client.DeleteByQuery([]string{index}, ``); err.Error() != `please call Initialize first` {
		t.Fatal(err)
	}

	client, ok := getClient(t)
	if ok == false {
		return
	}

	if err := indicesExists(&client); err != nil {
		t.Fatal(err)
	}

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
		t.Fatal("invalid exist")
	}

	if err := indicesDelete(&client); err != nil {
		t.Fatal(err)
	}
}

func TestIndicesExists(t *testing.T) {
	client := v8.Client{}
	if _, err := client.IndicesExists([]string{index}); err.Error() != `please call Initialize first` {
		t.Fatal(err)
	}

	client, ok := getClient(t)
	if ok == false {
		return
	}

	if err := indicesExists(&client); err != nil {
		t.Fatal(err)
	}

	if _, err := client.IndicesExists([]string{"<>"}); err.Error() != `response error - status : (400 Bad Request)` {
		t.Fatal(err)
	}

	if err := client.IndicesCreate(index, ""); err != nil {
		t.Fatal(err)
	}

	if exist, err := client.IndicesExists([]string{index}); err != nil {
		t.Fatal(err)
	} else if exist == false {
		t.Fatal("invalid exist")
	}

	if err := indicesDelete(&client); err != nil {
		t.Fatal(err)
	}
}

func TestIndicesCreate(t *testing.T) {
	client := v8.Client{}
	if err := client.IndicesCreate(index, ""); err.Error() != `please call Initialize first` {
		t.Fatal(err)
	}

	client, ok := getClient(t)
	if ok == false {
		return
	}

	if err := indicesExists(&client); err != nil {
		t.Fatal(err)
	}

	if err := client.IndicesCreate(index, "~"); err.Error() != `response error - status : (500 Internal Server Error), type : (not_x_content_exception), reason : (Compressor detection can only be called on some xcontent bytes or compressed xcontent bytes)` {
		t.Fatal(err)
	}

	if err := client.IndicesCreate(index, ""); err != nil {
		t.Fatal(err)
	}

	if exist, err := client.IndicesExists([]string{index}); err != nil {
		t.Fatal(err)
	} else if exist == false {
		t.Fatal("invalid exist")
	}

	if err := indicesDelete(&client); err != nil {
		t.Fatal(err)
	}
}

func TestIndicesDelete(t *testing.T) {
	client := v8.Client{}
	if err := client.IndicesDelete([]string{""}); err.Error() != `please call Initialize first` {
		t.Fatal(err)
	}

	client, ok := getClient(t)
	if ok == false {
		return
	}

	if err := indicesExists(&client); err != nil {
		t.Fatal(err)
	}

	if err := client.IndicesDelete([]string{""}); err.Error() != `response error - status : (400 Bad Request), type : (action_request_validation_exception), reason : (Validation Failed: 1: index / indices is missing;)` {
		t.Fatal(err)
	}

	if err := client.IndicesCreate(index, ""); err != nil {
		t.Fatal(err)
	}

	if err := indicesDelete(&client); err != nil {
		t.Fatal(err)
	}
}

func TestIndicesExistsTemplate(t *testing.T) {
	client := v8.Client{}
	if _, err := client.IndicesExistsTemplate([]string{template}); err.Error() != `please call Initialize first` {
		t.Fatal(err)
	}

	client, ok := getClient(t)
	if ok == false {
		return
	}

	if err := existsTemplate(&client); err != nil {
		t.Fatal(err)
	}

	if _, err := client.IndicesExistsTemplate([]string{""}); err.Error() != `response error - status : (405 Method Not Allowed)` {
		t.Fatal(err)
	}

	if err := client.IndicesPutTemplate(template, `{"index_patterns": ["*"]}`); err != nil {
		t.Fatal(err)
	}

	if exist, err := client.IndicesExistsTemplate([]string{template}); err != nil {
		t.Fatal(err)
	} else if exist == false {
		t.Fatal("invalid exist")
	}

	if err := indicesDeleteTemplate(&client); err != nil {
		t.Fatal(err)
	}
}

func TestIndicesPutTemplate(t *testing.T) {
	client := v8.Client{}
	if err := client.IndicesPutTemplate(template, ``); err.Error() != `please call Initialize first` {
		t.Fatal(err)
	}

	client, ok := getClient(t)
	if ok == false {
		return
	}

	if err := existsTemplate(&client); err != nil {
		t.Fatal(err)
	}

	if err := client.IndicesPutTemplate(template, ""); err.Error() != `response error - status : (400 Bad Request), type : (parse_exception), reason : (request body is required)` {
		t.Fatal(err)
	}

	if err := client.IndicesPutTemplate(template, `{"index_patterns": ["*"]}`); err != nil {
		t.Fatal(err)
	}

	if exist, err := client.IndicesExistsTemplate([]string{template}); err != nil {
		t.Fatal(err)
	} else if exist == false {
		t.Fatal("invalid exist")
	}

	if err := indicesDeleteTemplate(&client); err != nil {
		t.Fatal(err)
	}
}

func TestIndicesDeleteTemplate(t *testing.T) {
	client := v8.Client{}
	if err := client.IndicesDeleteTemplate(template); err.Error() != `please call Initialize first` {
		t.Fatal(err)
	}

	client, ok := getClient(t)
	if ok == false {
		return
	}

	if err := existsTemplate(&client); err != nil {
		t.Fatal(err)
	}

	if err := client.IndicesDeleteTemplate(template); err.Error() != `response error - status : (404 Not Found), type : (index_template_missing_exception), reason : (index_template [`+template+`] missing)` {
		t.Fatal(err)
	}

	if err := client.IndicesPutTemplate(template, `{"index_patterns": ["*"]}`); err != nil {
		t.Fatal(err)
	}

	if exist, err := client.IndicesExistsTemplate([]string{template}); err != nil {
		t.Fatal(err)
	} else if exist == false {
		t.Fatal("invalid exist")
	}

	if err := indicesDeleteTemplate(&client); err != nil {
		t.Fatal(err)
	}
}

func TestSearch(t *testing.T) {
	client := v8.Client{}
	if _, err := client.Search(index, ``); err.Error() != `please call Initialize first` {
		t.Fatal(err)
	}

	client, ok := getClient(t)
	if ok == false {
		return
	}

	if err := indicesExists(&client); err != nil {
		t.Fatal(err)
	}

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

	if err := indicesDelete(&client); err != nil {
		t.Fatal(err)
	}
}

func TestIndicesForcemerge(t *testing.T) {
	client := v8.Client{}
	if err := client.IndicesForcemerge([]string{index}); err.Error() != `please call Initialize first` {
		t.Fatal(err)
	}

	client, ok := getClient(t)
	if ok == false {
		return
	}

	if err := indicesExists(&client); err != nil {
		t.Fatal(err)
	}

	if err := client.IndicesForcemerge([]string{index}); err.Error() != `response error - status : (404 Not Found), type : (index_not_found_exception), reason : (no such index [`+index+`])` {
		t.Fatal(err)
	}

	if err := client.Index(index, documentId, `{"field":"value"}`); err != nil {
		t.Fatal(err)
	}

	if err := client.IndicesForcemerge([]string{index}); err != nil {
		t.Fatal(err)
	}

	if err := indicesDelete(&client); err != nil {
		t.Fatal(err)
	}
}
