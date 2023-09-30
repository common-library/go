package v8_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	v8 "github.com/heaven-chp/common-library-go/database/elasticsearch/v8"
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

func exists(client *v8.Client) error {
	err := initialize(client, addresses)
	if err != nil {
		return err
	}

	exist, err := client.IndicesExists([]string{index})
	if err != nil {
		return err
	}
	if exist {
		return errors.New(fmt.Sprintf("invalid exist - exist : (%t)", exist))
	}

	exist, err = client.Exists(index, documentId)
	if err != nil {
		return err
	}
	if exist {
		return errors.New(fmt.Sprintf("invalid exist - exist : (%t)", exist))
	}

	return nil
}

func deletes(client *v8.Client) error {
	err := client.IndicesDelete([]string{index})
	if err != nil {
		return err
	}

	exist, err := client.IndicesExists([]string{index})
	if err != nil {
		return err
	}
	if exist {
		return errors.New(fmt.Sprintf("invalid exist - exist : (%t)", exist))
	}

	return nil
}

func existsTemplate(client *v8.Client) error {
	err := initialize(client, addresses)
	if err != nil {
		return err
	}

	exist, err := client.IndicesExistsTemplate([]string{template})
	if err != nil {
		return err
	}
	if exist {
		return errors.New(fmt.Sprintf("invalid exist - exist : (%t)", exist))
	}

	return nil
}

func deletesTemplate(client *v8.Client) error {
	err := client.IndicesDeleteTemplate(template)
	if err != nil {
		return err
	}

	exist, err := client.IndicesExistsTemplate([]string{template})
	if err != nil {
		return err
	}
	if exist {
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

	err = client.Initialize(addresses, 0, "", "", "", "", "", []byte(""))
	if err != nil {
		t.Error(err)
	}

	err = initialize(&client, []string{"invalid_address"})
	if err != nil {
		t.Error(err)
	}

	err = initialize(&client, addresses)
	if err != nil {
		t.Error(err)
	}
}

func TestExists(t *testing.T) {
	client := v8.Client{}

	_, err := client.Exists("", "")
	if err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	err = initialize(&client, []string{"invalid_address"})
	if err != nil {
		t.Error(err)
	}
	_, err = client.Exists(index, documentId)
	if err.Error() != `unsupported protocol scheme ""` {
		t.Error(err)
	}

	err = exists(&client)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Exists("", "")
	if err.Error() != `response error - status : (405 Method Not Allowed)` {
		t.Error(err)
	}

	_, err = client.Exists("*", "")
	if err.Error() != `response error - status : (405 Method Not Allowed)` {
		t.Error(err)
	}

	err = client.Index(index, documentId, `{"field":"value"}`)
	if err != nil {
		t.Error(err)
	}

	exist, err := client.Exists(index, documentId)
	if err != nil {
		t.Error(err)
	}
	if exist == false {
		t.Errorf("invalid exist - exist : (%t)", exist)
	}

	err = deletes(&client)
	if err != nil {
		t.Error(err)
	}
}

func TestIndex(t *testing.T) {
	client := v8.Client{}

	err := client.Index(index, documentId, "")
	if err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	err = initialize(&client, []string{"invalid_address"})
	if err != nil {
		t.Error(err)
	}
	err = client.Index(index, documentId, `{"field":"value"}`)
	if err.Error() != `unsupported protocol scheme ""` {
		t.Error(err)
	}

	err = exists(&client)
	if err != nil {
		t.Fatal(err)
	}

	err = client.Index(index, documentId, "")
	if err.Error() != `response error - status : (400 Bad Request), type : (parse_exception), reason : (request body is required)` {
		t.Error(err)
	}

	err = client.Index(index, documentId, `{"field":"value"}`)
	if err != nil {
		t.Error(err)
	}

	exist, err := client.Exists(index, documentId)
	if err != nil {
		t.Error(err)
	}
	if exist == false {
		t.Errorf("invalid exist - exist : (%t)", exist)
	}

	err = deletes(&client)
	if err != nil {
		t.Error(err)
	}
}

func TestDelete(t *testing.T) {
	client := v8.Client{}

	err := client.Delete(index, documentId)
	if err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	err = initialize(&client, []string{"invalid_address"})
	if err != nil {
		t.Error(err)
	}
	err = client.Delete(index, documentId)
	if err.Error() != `unsupported protocol scheme ""` {
		t.Error(err)
	}

	err = exists(&client)
	if err != nil {
		t.Fatal(err)
	}

	err = client.Delete(index, documentId)
	if err.Error() != `response error - status : (404 Not Found), type : (index_not_found_exception), reason : (no such index [`+index+`])` {
		t.Error(err)
	}

	err = client.Index(index, documentId, `{"field":"value"}`)
	if err != nil {
		t.Error(err)
	}

	err = client.Delete(index, documentId)
	if err != nil {
		t.Error(err)
	}

	exist, err := client.Exists(index, documentId)
	if err != nil {
		t.Error(err)
	}
	if exist {
		t.Errorf("invalid exist - exist : (%t)", exist)
	}

	err = deletes(&client)
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteByQuery(t *testing.T) {
	client := v8.Client{}

	err := client.DeleteByQuery([]string{index}, ``)
	if err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	err = initialize(&client, []string{"invalid_address"})
	if err != nil {
		t.Error(err)
	}
	err = client.DeleteByQuery([]string{index}, `{"query":{"match":{"field":"value_1"}}}`)
	if err.Error() != `unsupported protocol scheme ""` {
		t.Error(err)
	}

	err = exists(&client)
	if err != nil {
		t.Fatal(err)
	}

	err = client.DeleteByQuery([]string{index}, `{}`)
	if err.Error() != `response error - status : (400 Bad Request), type : (action_request_validation_exception), reason : (Validation Failed: 1: query is missing;)` {
		t.Error(err)
	}

	err = client.Index(index, documentId, `{"field":"value_1"}`)
	if err != nil {
		t.Error(err)
	}

	err = client.Index(index, documentId+"_temp", `{"field":"value_2"}`)
	if err != nil {
		t.Error(err)
	}

	err = client.DeleteByQuery([]string{index}, `{"query":{"match":{"field":"value_1"}}}`)
	if err != nil {
		t.Error(err)
	}

	exist, err := client.Exists(index, documentId)
	if err != nil {
		t.Error(err)
	}
	if exist {
		t.Errorf("invalid exist - exist : (%t)", exist)
	}

	err = deletes(&client)
	if err != nil {
		t.Error(err)
	}
}

func TestIndicesExists(t *testing.T) {
	client := v8.Client{}

	_, err := client.IndicesExists([]string{index})
	if err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	err = initialize(&client, []string{"invalid_address"})
	if err != nil {
		t.Error(err)
	}
	_, err = client.IndicesExists([]string{index})
	if err.Error() != `unsupported protocol scheme ""` {
		t.Error(err)
	}

	err = exists(&client)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.IndicesExists([]string{"<>"})
	if err.Error() != `response error - status : (400 Bad Request)` {
		t.Error(err)
	}

	err = client.IndicesCreate(index, "")
	if err != nil {
		t.Error(err)
	}

	exist, err := client.IndicesExists([]string{index})
	if err != nil {
		t.Error(err)
	}
	if exist == false {
		t.Fatalf("invalid exist - exist : (%t)", exist)
	}

	err = deletes(&client)
	if err != nil {
		t.Error(err)
	}
}

func TestIndicesCreate(t *testing.T) {
	client := v8.Client{}

	err := client.IndicesCreate(index, "")
	if err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	err = initialize(&client, []string{"invalid_address"})
	if err != nil {
		t.Error(err)
	}
	err = client.IndicesCreate(index, "")
	if err.Error() != `unsupported protocol scheme ""` {
		t.Error(err)
	}

	err = exists(&client)
	if err != nil {
		t.Fatal(err)
	}

	err = client.IndicesCreate(index, "~")
	if err.Error() != `response error - status : (500 Internal Server Error), type : (not_x_content_exception), reason : (Compressor detection can only be called on some xcontent bytes or compressed xcontent bytes)` {
		t.Error(err)
	}

	err = client.IndicesCreate(index, "")
	if err != nil {
		t.Error(err)
	}

	exist, err := client.IndicesExists([]string{index})
	if err != nil {
		t.Error(err)
	}
	if exist == false {
		t.Fatalf("invalid exist - exist : (%t)", exist)
	}

	err = deletes(&client)
	if err != nil {
		t.Error(err)
	}
}

func TestIndicesDelete(t *testing.T) {
	client := v8.Client{}

	err := client.IndicesDelete([]string{""})
	if err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	err = initialize(&client, []string{"invalid_address"})
	if err != nil {
		t.Error(err)
	}
	err = client.IndicesDelete([]string{""})
	if err.Error() != `unsupported protocol scheme ""` {
		t.Error(err)
	}

	err = exists(&client)
	if err != nil {
		t.Fatal(err)
	}

	err = client.IndicesDelete([]string{""})
	if err.Error() != `response error - status : (400 Bad Request), type : (action_request_validation_exception), reason : (Validation Failed: 1: index / indices is missing;)` {
		t.Fatal(err)
	}

	err = client.IndicesCreate(index, "")
	if err != nil {
		t.Error(err)
	}

	err = deletes(&client)
	if err != nil {
		t.Error(err)
	}
}

func TestIndicesExistsTemplate(t *testing.T) {
	client := v8.Client{}

	_, err := client.IndicesExistsTemplate([]string{template})
	if err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	err = initialize(&client, []string{"invalid_address"})
	if err != nil {
		t.Error(err)
	}
	_, err = client.IndicesExistsTemplate([]string{template})
	if err.Error() != `unsupported protocol scheme ""` {
		t.Error(err)
	}

	err = existsTemplate(&client)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.IndicesExistsTemplate([]string{""})
	if err.Error() != `response error - status : (405 Method Not Allowed)` {
		t.Error(err)
	}

	err = client.IndicesPutTemplate(template, `{"index_patterns": ["*"]}`)
	if err != nil {
		t.Error(err)
	}

	exist, err := client.IndicesExistsTemplate([]string{template})
	if err != nil {
		t.Error(err)
	}
	if exist == false {
		t.Fatalf("invalid exist - exist : (%t)", exist)
	}

	err = deletesTemplate(&client)
	if err != nil {
		t.Error(err)
	}
}

func TestIndicesPutTemplate(t *testing.T) {
	client := v8.Client{}

	err := client.IndicesPutTemplate(template, ``)
	if err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	err = initialize(&client, []string{"invalid_address"})
	if err != nil {
		t.Error(err)
	}
	err = client.IndicesPutTemplate(template, `{"index_patterns": ["*"]}`)
	if err.Error() != `unsupported protocol scheme ""` {
		t.Error(err)
	}

	err = existsTemplate(&client)
	if err != nil {
		t.Fatal(err)
	}

	err = client.IndicesPutTemplate(template, "")
	if err.Error() != `response error - status : (400 Bad Request), type : (parse_exception), reason : (request body is required)` {
		t.Error(err)
	}

	err = client.IndicesPutTemplate(template, `{"index_patterns": ["*"]}`)
	if err != nil {
		t.Error(err)
	}

	exist, err := client.IndicesExistsTemplate([]string{template})
	if err != nil {
		t.Error(err)
	}
	if exist == false {
		t.Fatalf("invalid exist - exist : (%t)", exist)
	}

	err = deletesTemplate(&client)
	if err != nil {
		t.Error(err)
	}
}

func TestIndicesDeleteTemplate(t *testing.T) {
	client := v8.Client{}

	err := client.IndicesDeleteTemplate(template)
	if err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	err = initialize(&client, []string{"invalid_address"})
	if err != nil {
		t.Error(err)
	}
	err = client.IndicesDeleteTemplate(template)
	if err.Error() != `unsupported protocol scheme ""` {
		t.Error(err)
	}

	err = existsTemplate(&client)
	if err != nil {
		t.Fatal(err)
	}

	err = client.IndicesDeleteTemplate(template)
	if err.Error() != `response error - status : (404 Not Found), type : (index_template_missing_exception), reason : (index_template [`+template+`] missing)` {
		t.Error(err)
	}

	err = client.IndicesPutTemplate(template, `{"index_patterns": ["*"]}`)
	if err != nil {
		t.Error(err)
	}

	exist, err := client.IndicesExistsTemplate([]string{template})
	if err != nil {
		t.Error(err)
	}
	if exist == false {
		t.Fatalf("invalid exist - exist : (%t)", exist)
	}

	err = deletesTemplate(&client)
	if err != nil {
		t.Error(err)
	}
}

func TestSearch(t *testing.T) {
	client := v8.Client{}

	_, err := client.Search(index, ``)
	if err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	err = initialize(&client, []string{"invalid_address"})
	if err != nil {
		t.Error(err)
	}
	_, err = client.Search(index, `{"query":{"match_all":{}}}`)
	if err.Error() != `unsupported protocol scheme ""` {
		t.Error(err)
	}

	err = exists(&client)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Search(index, `{}`)
	if err.Error() != `response error - status : (404 Not Found), type : (index_not_found_exception), reason : (no such index [`+index+`])` {
		t.Error(err)
	}

	err = client.Index(index, documentId, `{"field":"value"}`)
	if err != nil {
		t.Error(err)
	}

	result, err := client.Search(index, `{"query":{"match_all":{}}}`)
	if err != nil {
		t.Error(err)
	}
	if gojsonq.New().FromString(result).Find("hits.total.value").(float64) != 1 {
		t.Errorf("invalid result - result : (\n%s)", result)
	}

	err = deletes(&client)
	if err != nil {
		t.Error(err)
	}
}

func TestIndicesForcemerge(t *testing.T) {
	client := v8.Client{}

	err := client.IndicesForcemerge([]string{index})
	if err.Error() != `please call Initialize first` {
		t.Error(err)
	}

	err = initialize(&client, []string{"invalid_address"})
	if err != nil {
		t.Error(err)
	}
	err = client.IndicesForcemerge([]string{index})
	if err.Error() != `unsupported protocol scheme ""` {
		t.Error(err)
	}

	err = exists(&client)
	if err != nil {
		t.Fatal(err)
	}

	err = client.IndicesForcemerge([]string{index})
	if err.Error() != `response error - status : (404 Not Found), type : (index_not_found_exception), reason : (no such index [`+index+`])` {
		t.Error(err)
	}

	err = client.Index(index, documentId, `{"field":"value"}`)
	if err != nil {
		t.Error(err)
	}

	err = client.IndicesForcemerge([]string{index})
	if err != nil {
		t.Error(err)
	}

	err = deletes(&client)
	if err != nil {
		t.Error(err)
	}
}
