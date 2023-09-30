// Package elasticsearch provides elasticsearch interface.
//
// used "github.com/elastic/go-elasticsearch/v7".
package v7

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/thedevsaddam/gojsonq/v2"
)

// Client is object that provides elasticsearch interface.
type Client struct {
	client *elasticsearch.Client
}

// Initialize is initialize.
//
// ex) err := client.Initialize([]string{"127.0.0.1:9200"}, 60, "", "", "", "", "", []byte(""))
func (this *Client) Initialize(addresses []string, timeout uint64, cloudID, apiKey, username, password, certificateFingerprint string, caCert []byte) error {
	config := elasticsearch.Config{
		CloudID:                cloudID,
		APIKey:                 apiKey,
		Username:               username,
		Password:               password,
		CertificateFingerprint: certificateFingerprint,

		Addresses:         addresses,
		EnableDebugLogger: true,
		Logger:            &elastictransport.ColorLogger{Output: os.Stdout},
		Transport: &http.Transport{
			ResponseHeaderTimeout: time.Second * time.Duration(timeout),
		},
	}
	if len(caCert) != 0 {
		config.CACert = caCert
	}

	if client, err := elasticsearch.NewClient(config); err != nil {
		return err
	} else {
		this.client = client
	}

	return nil
}

// Exists is checks if a document exists in the index.
//
// ex) exist, err := client.Exists("index", "document_id")
func (this *Client) Exists(index, documentID string) (bool, error) {
	if this.client == nil {
		return false, errors.New("please call Initialize first")
	}

	refresh := true

	request := esapi.ExistsRequest{
		Index:      index,
		DocumentID: documentID,
		Refresh:    &refresh,
		Human:      true,
		ErrorTrace: true}

	if response, err := request.Do(context.Background(), this.client); err != nil {
		return false, err
	} else {
		defer response.Body.Close()

		if response.StatusCode == http.StatusNotFound {
			return false, nil
		} else if response.IsError() {
			return false, this.responseErrorToError(response.Status(), response.Body)
		}
	}

	return true, nil
}

// Index is stores document.
//
// ex) err := client.Index("index", "document_id", "{...}")
func (this *Client) Index(index, documentID, body string) error {
	if this.client == nil {
		return errors.New("please call Initialize first")
	}

	builder := strings.Builder{}
	if _, err := builder.WriteString(body); err != nil {
		return err
	}

	request := esapi.IndexRequest{
		Index:      index,
		DocumentID: documentID,
		Body:       strings.NewReader(builder.String()),
		Refresh:    "true",
		Human:      true,
		ErrorTrace: true}

	if response, err := request.Do(context.Background(), this.client); err != nil {
		return err
	} else {
		defer response.Body.Close()

		if response.IsError() {
			return this.responseErrorToError(response.Status(), response.Body)
		}
	}

	return nil
}

// Delete is deletes document.
//
// ex) err := client.Delete("index", "document_id")
func (this *Client) Delete(index, documentID string) error {
	if this.client == nil {
		return errors.New("please call Initialize first")
	}

	request := esapi.DeleteRequest{
		Index:      index,
		DocumentID: documentID,
		Refresh:    "true",
		Human:      true,
		ErrorTrace: true}

	if response, err := request.Do(context.Background(), this.client); err != nil {
		return err
	} else {
		defer response.Body.Close()

		if response.IsError() {
			return this.responseErrorToError(response.Status(), response.Body)
		}
	}

	return nil
}

// DeleteByQuery is perform a delete query on index.
//
// ex) err := client.DeleteByQuery([]string{"index"},"{...}")
func (this *Client) DeleteByQuery(indices []string, body string) error {
	if this.client == nil {
		return errors.New("please call Initialize first")
	}

	builder := strings.Builder{}
	if _, err := builder.WriteString(body); err != nil {
		return err
	}

	refresh := true

	request := esapi.DeleteByQueryRequest{
		Index:      indices,
		Body:       strings.NewReader(builder.String()),
		Refresh:    &refresh,
		Human:      true,
		ErrorTrace: true}

	if response, err := request.Do(context.Background(), this.client); err != nil {
		return err
	} else {
		defer response.Body.Close()

		if response.IsError() {
			return this.responseErrorToError(response.Status(), response.Body)
		}
	}

	return nil
}

// IndicesExists is checks if an index exists within indices.
//
// ex) exist, err := client.IndicesExists([]string{"index"})
func (this *Client) IndicesExists(indices []string) (bool, error) {
	if this.client == nil {
		return false, errors.New("please call Initialize first")
	}

	request := esapi.IndicesExistsRequest{
		Index:      indices,
		Human:      true,
		ErrorTrace: true}

	if response, err := request.Do(context.Background(), this.client); err != nil {
		return false, err
	} else {
		defer response.Body.Close()

		if response.StatusCode == http.StatusNotFound {
			return false, nil
		} else if response.IsError() {
			return false, this.responseErrorToError(response.Status(), response.Body)
		}
	}

	return true, nil
}

// IndicesCreate is create an index.
//
// ex) err := client.IndicesCreate("index", "{...}")
func (this *Client) IndicesCreate(index, body string) error {
	if this.client == nil {
		return errors.New("please call Initialize first")
	}

	builder := strings.Builder{}
	if _, err := builder.WriteString(body); err != nil {
		return err
	}

	request := esapi.IndicesCreateRequest{
		Index:      index,
		Body:       strings.NewReader(builder.String()),
		Human:      true,
		ErrorTrace: true}

	if response, err := request.Do(context.Background(), this.client); err != nil {
		return err
	} else {
		defer response.Body.Close()

		if response.IsError() {
			return this.responseErrorToError(response.Status(), response.Body)
		}
	}

	return nil
}

// IndicesDelete is delete an index.
//
// ex) err := client.IndicesDelete([]string{"Index"})
func (this *Client) IndicesDelete(indices []string) error {
	if this.client == nil {
		return errors.New("please call Initialize first")
	}

	request := esapi.IndicesDeleteRequest{
		Index:      indices,
		Human:      true,
		ErrorTrace: true}

	if response, err := request.Do(context.Background(), this.client); err != nil {
		return err
	} else {
		defer response.Body.Close()

		if response.IsError() {
			return this.responseErrorToError(response.Status(), response.Body)
		}
	}

	return nil
}

// IndicesExistsTemplate is checks if a template exists.
//
// ex) exist, err := client.IndicesExistsTemplate([]string{"template"})
func (this *Client) IndicesExistsTemplate(name []string) (bool, error) {
	if this.client == nil {
		return false, errors.New("please call Initialize first")
	}

	request := esapi.IndicesExistsTemplateRequest{
		Name:       name,
		Human:      true,
		ErrorTrace: true}

	if response, err := request.Do(context.Background(), this.client); err != nil {
		return false, err
	} else {
		defer response.Body.Close()

		if response.StatusCode == http.StatusNotFound {
			return false, nil
		} else if response.IsError() {
			return false, this.responseErrorToError(response.Status(), response.Body)
		}
	}

	return true, nil
}

// IndicesPutTemplate is stores templates.
//
// ex) err := client.IndicesPutTemplate("template", "{...}")
func (this *Client) IndicesPutTemplate(name, body string) error {
	if this.client == nil {
		return errors.New("please call Initialize first")
	}

	builder := strings.Builder{}
	if _, err := builder.WriteString(body); err != nil {
		return err
	}

	request := esapi.IndicesPutTemplateRequest{
		Name:       name,
		Body:       strings.NewReader(builder.String()),
		Human:      true,
		ErrorTrace: true}

	if response, err := request.Do(context.Background(), this.client); err != nil {
		return err
	} else {
		defer response.Body.Close()

		if response.IsError() {
			return this.responseErrorToError(response.Status(), response.Body)
		}
	}

	return nil
}

// IndicesDeleteTemplate is delete an template.
//
// ex) err := client.IndicesDeleteTemplate("template")
func (this *Client) IndicesDeleteTemplate(name string) error {
	if this.client == nil {
		return errors.New("please call Initialize first")
	}

	request := esapi.IndicesDeleteTemplateRequest{
		Name:       name,
		Human:      true,
		ErrorTrace: true}

	if response, err := request.Do(context.Background(), this.client); err != nil {
		return err
	} else {
		defer response.Body.Close()

		if response.IsError() {
			return this.responseErrorToError(response.Status(), response.Body)
		}
	}

	return nil
}

// IndicesForcemerge is perform a force merge on index.
//
// ex) err := client.IndicesForcemerge([]string{"index"})
func (this *Client) IndicesForcemerge(indices []string) error {
	if this.client == nil {
		return errors.New("please call Initialize first")
	}

	onlyExpungeDeletes := true

	request := esapi.IndicesForcemergeRequest{
		Index:              indices,
		OnlyExpungeDeletes: &onlyExpungeDeletes,
		Human:              true,
		ErrorTrace:         true}

	if response, err := request.Do(context.Background(), this.client); err != nil {
		return err
	} else {
		defer response.Body.Close()

		if response.IsError() {
			return this.responseErrorToError(response.Status(), response.Body)
		}
	}

	return nil
}

// Search is search
//
// ex) result, err := client.Search("index", "{...}")
func (this *Client) Search(index, body string) (string, error) {
	if this.client == nil {
		return "", errors.New("please call Initialize first")
	}

	builder := strings.Builder{}
	if _, err := builder.WriteString(body); err != nil {
		return "", err
	}

	if response, err := this.client.Search(
		this.client.Search.WithContext(context.Background()),
		this.client.Search.WithIndex(index),
		this.client.Search.WithBody(strings.NewReader(builder.String())),
		this.client.Search.WithTrackTotalHits(true),
		this.client.Search.WithPretty()); err != nil {
		return "", err
	} else {
		defer response.Body.Close()

		if response.IsError() {
			return "", this.responseErrorToError(response.Status(), response.Body)
		}

		result, err := ioutil.ReadAll(response.Body)
		return string(result), err
	}
}

func (this *Client) responseErrorToError(status string, reader io.Reader) error {
	buffer := new(bytes.Buffer)
	buffer.ReadFrom(reader)

	if len(buffer.String()) == 0 {
		return errors.New(fmt.Sprintf("response error - status : (%s)", status))
	}

	return errors.New(
		fmt.Sprintf("response error - status : (%s), type : (%s), reason : (%s)",
			status,
			gojsonq.New().FromString(buffer.String()).Find("error.type").(string),
			gojsonq.New().FromString(buffer.String()).Find("error.reason").(string)))
}
