package v8_test

import (
	"errors"
	"fmt"
	"testing"

	v8 "github.com/common-library/go/database/elasticsearch/v8"
	"github.com/google/uuid"
	"github.com/thedevsaddam/gojsonq/v2"
)

var addresses []string = []string{"http://127.0.0.1:29200"}
var timeout uint64 = 10
var index string = uuid.NewString()
var documentId string = uuid.NewString()
var template string = uuid.NewString()

func initialize(client *v8.Client, addresses []string) error {
	return client.Initialize(addresses, timeout, "", "", "", "", "", []byte(""))
}

func indicesExists(client *v8.Client) error {
	if err := initialize(client, addresses); err != nil {
		return err
	}

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
	if err := initialize(client, addresses); err != nil {
		return err
	}

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
	client := v8.Client{}

	err := client.Initialize(addresses, timeout, "", "", "", "", "", []byte("invalid"))
	if err.Error() != "error creating transport: unable to add CA certificate" {
		t.Error(err)
	}

	if err := client.Initialize(addresses, 0, "", "", "", "", "", []byte("")); err != nil {
		t.Error(err)
	}

	if err := initialize(&client, []string{"invalid_address"}); err != nil {
		t.Error(err)
	}

	if err := initialize(&client, addresses); err != nil {
		t.Error(err)
	}
}

func TestExists(t *testing.T) {
	client := v8.Client{}

	if _, err := client.Exists("", ""); err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	if err := initialize(&client, []string{"invalid_address"}); err != nil {
		t.Error(err)
	}

	if _, err := client.Exists(index, documentId); err.Error() != `unsupported protocol scheme ""` {
		t.Error(err)
	}

	if err := indicesExists(&client); err != nil {
		t.Fatal(err)
	}

	if _, err := client.Exists("", ""); err.Error() != `response error - status : (405 Method Not Allowed)` {
		t.Error(err)
	}

	if _, err := client.Exists("*", ""); err.Error() != `response error - status : (405 Method Not Allowed)` {
		t.Error(err)
	}

	if err := client.Index(index, documentId, `{"field":"value"}`); err != nil {
		t.Error(err)
	}

	if exist, err := client.Exists(index, documentId); err != nil {
		t.Error(err)
	} else if exist == false {
		t.Errorf("invalid exist - exist : (%t)", exist)
	}

	if err := indicesDelete(&client); err != nil {
		t.Error(err)
	}
}

func TestIndex(t *testing.T) {
	client := v8.Client{}

	if err := client.Index(index, documentId, ""); err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	if err := initialize(&client, []string{"invalid_address"}); err != nil {
		t.Error(err)
	}

	if err := client.Index(index, documentId, `{"field":"value"}`); err.Error() != `unsupported protocol scheme ""` {
		t.Error(err)
	}

	if err := indicesExists(&client); err != nil {
		t.Fatal(err)
	}

	if err := client.Index(index, documentId, ""); err.Error() != `response error - status : (400 Bad Request), type : (parse_exception), reason : (request body is required)` {
		t.Error(err)
	}

	if err := client.Index(index, documentId, `{"field":"value"}`); err != nil {
		t.Error(err)
	}

	if exist, err := client.Exists(index, documentId); err != nil {
		t.Error(err)
	} else if exist == false {
		t.Errorf("invalid exist - exist : (%t)", exist)
	}

	if err := indicesDelete(&client); err != nil {
		t.Error(err)
	}
}

func TestDelete(t *testing.T) {
	client := v8.Client{}

	if err := client.Delete(index, documentId); err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	if err := initialize(&client, []string{"invalid_address"}); err != nil {
		t.Error(err)
	}

	if err := client.Delete(index, documentId); err.Error() != `unsupported protocol scheme ""` {
		t.Error(err)
	}

	if err := indicesExists(&client); err != nil {
		t.Fatal(err)
	}

	if err := client.Delete(index, documentId); err.Error() != `response error - status : (404 Not Found), type : (index_not_found_exception), reason : (no such index [`+index+`])` {
		t.Error(err)
	}

	if err := client.Index(index, documentId, `{"field":"value"}`); err != nil {
		t.Error(err)
	}

	if err := client.Delete(index, documentId); err != nil {
		t.Error(err)
	}

	if exist, err := client.Exists(index, documentId); err != nil {
		t.Error(err)
	} else if exist {
		t.Errorf("invalid exist - exist : (%t)", exist)
	}

	if err := indicesDelete(&client); err != nil {
		t.Error(err)
	}
}

func TestDeleteByQuery(t *testing.T) {
	client := v8.Client{}

	if err := client.DeleteByQuery([]string{index}, ``); err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	if err := initialize(&client, []string{"invalid_address"}); err != nil {
		t.Error(err)
	}

	if err := client.DeleteByQuery([]string{index}, `{"query":{"match":{"field":"value_1"}}}`); err.Error() != `unsupported protocol scheme ""` {
		t.Error(err)
	}

	if err := indicesExists(&client); err != nil {
		t.Fatal(err)
	}

	if err := client.DeleteByQuery([]string{index}, `{}`); err.Error() != `response error - status : (400 Bad Request), type : (action_request_validation_exception), reason : (Validation Failed: 1: query is missing;)` {
		t.Error(err)
	}

	if err := client.Index(index, documentId, `{"field":"value_1"}`); err != nil {
		t.Error(err)
	}

	if err := client.Index(index, documentId+"_temp", `{"field":"value_2"}`); err != nil {
		t.Error(err)
	}

	if err := client.DeleteByQuery([]string{index}, `{"query":{"match":{"field":"value_1"}}}`); err != nil {
		t.Error(err)
	}

	if exist, err := client.Exists(index, documentId); err != nil {
		t.Error(err)
	} else if exist {
		t.Errorf("invalid exist - exist : (%t)", exist)
	}

	if err := indicesDelete(&client); err != nil {
		t.Error(err)
	}
}

func TestIndicesExists(t *testing.T) {
	client := v8.Client{}

	if _, err := client.IndicesExists([]string{index}); err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	if err := initialize(&client, []string{"invalid_address"}); err != nil {
		t.Error(err)
	}

	if _, err := client.IndicesExists([]string{index}); err.Error() != `unsupported protocol scheme ""` {
		t.Error(err)
	}

	if err := indicesExists(&client); err != nil {
		t.Fatal(err)
	}

	if _, err := client.IndicesExists([]string{"<>"}); err.Error() != `response error - status : (400 Bad Request)` {
		t.Error(err)
	}

	if err := client.IndicesCreate(index, ""); err != nil {
		t.Error(err)
	}

	if exist, err := client.IndicesExists([]string{index}); err != nil {
		t.Error(err)
	} else if exist == false {
		t.Fatalf("invalid exist - exist : (%t)", exist)
	}

	if err := indicesDelete(&client); err != nil {
		t.Error(err)
	}
}

func TestIndicesCreate(t *testing.T) {
	client := v8.Client{}

	if err := client.IndicesCreate(index, ""); err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	if err := initialize(&client, []string{"invalid_address"}); err != nil {
		t.Error(err)
	}

	if err := client.IndicesCreate(index, ""); err.Error() != `unsupported protocol scheme ""` {
		t.Error(err)
	}

	if err := indicesExists(&client); err != nil {
		t.Fatal(err)
	}

	if err := client.IndicesCreate(index, "~"); err.Error() != `response error - status : (500 Internal Server Error), type : (not_x_content_exception), reason : (Compressor detection can only be called on some xcontent bytes or compressed xcontent bytes)` {
		t.Error(err)
	}

	if err := client.IndicesCreate(index, ""); err != nil {
		t.Error(err)
	}

	if exist, err := client.IndicesExists([]string{index}); err != nil {
		t.Error(err)
	} else if exist == false {
		t.Fatalf("invalid exist - exist : (%t)", exist)
	}

	if err := indicesDelete(&client); err != nil {
		t.Error(err)
	}
}

func TestIndicesDelete(t *testing.T) {
	client := v8.Client{}

	if err := client.IndicesDelete([]string{""}); err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	if err := initialize(&client, []string{"invalid_address"}); err != nil {
		t.Error(err)
	}

	if err := client.IndicesDelete([]string{""}); err.Error() != `unsupported protocol scheme ""` {
		t.Error(err)
	}

	if err := indicesExists(&client); err != nil {
		t.Fatal(err)
	}

	if err := client.IndicesDelete([]string{""}); err.Error() != `response error - status : (400 Bad Request), type : (action_request_validation_exception), reason : (Validation Failed: 1: index / indices is missing;)` {
		t.Fatal(err)
	}

	if err := client.IndicesCreate(index, ""); err != nil {
		t.Error(err)
	}

	if err := indicesDelete(&client); err != nil {
		t.Error(err)
	}
}

func TestIndicesExistsTemplate(t *testing.T) {
	client := v8.Client{}

	if _, err := client.IndicesExistsTemplate([]string{template}); err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	if err := initialize(&client, []string{"invalid_address"}); err != nil {
		t.Error(err)
	}

	if _, err := client.IndicesExistsTemplate([]string{template}); err.Error() != `unsupported protocol scheme ""` {
		t.Error(err)
	}

	if err := existsTemplate(&client); err != nil {
		t.Fatal(err)
	}

	if _, err := client.IndicesExistsTemplate([]string{""}); err.Error() != `response error - status : (405 Method Not Allowed)` {
		t.Error(err)
	}

	if err := client.IndicesPutTemplate(template, `{"index_patterns": ["*"]}`); err != nil {
		t.Error(err)
	}

	if exist, err := client.IndicesExistsTemplate([]string{template}); err != nil {
		t.Error(err)
	} else if exist == false {
		t.Fatalf("invalid exist - exist : (%t)", exist)
	}

	if err := indicesDeleteTemplate(&client); err != nil {
		t.Error(err)
	}
}

func TestIndicesPutTemplate(t *testing.T) {
	client := v8.Client{}

	if err := client.IndicesPutTemplate(template, ``); err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	if err := initialize(&client, []string{"invalid_address"}); err != nil {
		t.Error(err)
	}

	if err := client.IndicesPutTemplate(template, `{"index_patterns": ["*"]}`); err.Error() != `unsupported protocol scheme ""` {
		t.Error(err)
	}

	if err := existsTemplate(&client); err != nil {
		t.Fatal(err)
	}

	if err := client.IndicesPutTemplate(template, ""); err.Error() != `response error - status : (400 Bad Request), type : (parse_exception), reason : (request body is required)` {
		t.Error(err)
	}

	if err := client.IndicesPutTemplate(template, `{"index_patterns": ["*"]}`); err != nil {
		t.Error(err)
	}

	if exist, err := client.IndicesExistsTemplate([]string{template}); err != nil {
		t.Error(err)
	} else if exist == false {
		t.Fatalf("invalid exist - exist : (%t)", exist)
	}

	if err := indicesDeleteTemplate(&client); err != nil {
		t.Error(err)
	}
}

func TestIndicesDeleteTemplate(t *testing.T) {
	client := v8.Client{}

	if err := client.IndicesDeleteTemplate(template); err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	if err := initialize(&client, []string{"invalid_address"}); err != nil {
		t.Error(err)
	}

	if err := client.IndicesDeleteTemplate(template); err.Error() != `unsupported protocol scheme ""` {
		t.Error(err)
	}

	if err := existsTemplate(&client); err != nil {
		t.Fatal(err)
	}

	if err := client.IndicesDeleteTemplate(template); err.Error() != `response error - status : (404 Not Found), type : (index_template_missing_exception), reason : (index_template [`+template+`] missing)` {
		t.Error(err)
	}

	if err := client.IndicesPutTemplate(template, `{"index_patterns": ["*"]}`); err != nil {
		t.Error(err)
	}

	if exist, err := client.IndicesExistsTemplate([]string{template}); err != nil {
		t.Error(err)
	} else if exist == false {
		t.Fatalf("invalid exist - exist : (%t)", exist)
	}

	if err := indicesDeleteTemplate(&client); err != nil {
		t.Error(err)
	}
}

func TestSearch(t *testing.T) {
	client := v8.Client{}

	if _, err := client.Search(index, ``); err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	if err := initialize(&client, []string{"invalid_address"}); err != nil {
		t.Error(err)
	}

	if _, err := client.Search(index, `{"query":{"match_all":{}}}`); err.Error() != `unsupported protocol scheme ""` {
		t.Error(err)
	}

	if err := indicesExists(&client); err != nil {
		t.Fatal(err)
	}

	if _, err := client.Search(index, `{}`); err.Error() != `response error - status : (404 Not Found), type : (index_not_found_exception), reason : (no such index [`+index+`])` {
		t.Error(err)
	}

	if err := client.Index(index, documentId, `{"field":"value"}`); err != nil {
		t.Error(err)
	}

	if result, err := client.Search(index, `{"query":{"match_all":{}}}`); err != nil {
		t.Error(err)
	} else if gojsonq.New().FromString(result).Find("hits.total.value").(float64) != 1 {
		t.Errorf("invalid result - result : (\n%s)", result)
	}

	if err := indicesDelete(&client); err != nil {
		t.Error(err)
	}
}

func TestIndicesForcemerge(t *testing.T) {
	client := v8.Client{}

	if err := client.IndicesForcemerge([]string{index}); err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	if err := initialize(&client, []string{"invalid_address"}); err != nil {
		t.Error(err)
	}

	if err := client.IndicesForcemerge([]string{index}); err.Error() != `unsupported protocol scheme ""` {
		t.Error(err)
	}

	if err := indicesExists(&client); err != nil {
		t.Fatal(err)
	}

	if err := client.IndicesForcemerge([]string{index}); err.Error() != `response error - status : (404 Not Found), type : (index_not_found_exception), reason : (no such index [`+index+`])` {
		t.Error(err)
	}

	if err := client.Index(index, documentId, `{"field":"value"}`); err != nil {
		t.Error(err)
	}

	if err := client.IndicesForcemerge([]string{index}); err != nil {
		t.Error(err)
	}

	if err := indicesDelete(&client); err != nil {
		t.Error(err)
	}
}
