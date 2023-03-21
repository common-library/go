// Package elasticsearch provides elasticsearch interface.
//
// used "github.com/elastic/go-elasticsearch/v8".
package elasticsearch

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
	es_v8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/thedevsaddam/gojsonq/v2"
)

// Elasticsearch is object that provides elasticsearch interface.
type Elasticsearch struct {
	addresses []string
	timeout   int

	client *es_v8.Client
}

// Initialize is initialize.
//
// ex) err := elasticsearch.Initialize([]string{"127.0.0.1:9200"}, 60, "", "", "", "", "", []byte(""))
func (this *Elasticsearch) Initialize(addresses []string, timeout int, cloudID, apiKey, username, password, certificateFingerprint string, caCert []byte) error {
	this.addresses = addresses
	this.timeout = timeout

	config := es_v8.Config{
		CloudID:                cloudID,
		APIKey:                 apiKey,
		Username:               username,
		Password:               password,
		CertificateFingerprint: certificateFingerprint,

		Addresses:         this.addresses,
		EnableDebugLogger: true,
		Logger:            &elastictransport.ColorLogger{Output: os.Stdout},
		Transport: &http.Transport{
			ResponseHeaderTimeout: time.Second * time.Duration(timeout),
		},
	}
	if len(caCert) != 0 {
		config.CACert = caCert
	}

	var err error
	this.client, err = es_v8.NewClient(config)
	if err != nil {
		return err
	}

	return nil
}

// Exists is checks if a document exists in the index.
//
// ex) exist, err := elasticsearch.Exists("index", "document_id")
func (this *Elasticsearch) Exists(index, documentID string) (bool, error) {
	if this.client == nil {
		return false, errors.New("please call Initialize first")
	}

	refresh := true

	request := esapi.ExistsRequest{
		Index:      index,
		DocumentID: documentID,
		Refresh:    &refresh,
		Human:      true,
		ErrorTrace: true,
	}

	response, err := request.Do(context.Background(), this.client)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return false, nil
	}

	if response.IsError() {
		return false, this.responseErrorToError(response.Status(), response.Body)
	}

	return true, nil
}

// Index is stores document.
//
// ex) err := elasticsearch.Index("index", "document_id", "{...}")
func (this *Elasticsearch) Index(index, documentID, body string) error {
	if this.client == nil {
		return errors.New("please call Initialize first")
	}

	var builder strings.Builder
	builder.WriteString(body)

	request := esapi.IndexRequest{
		Index:      index,
		DocumentID: documentID,
		Body:       strings.NewReader(builder.String()),
		Refresh:    "true",
		Human:      true,
		ErrorTrace: true,
	}

	response, err := request.Do(context.Background(), this.client)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.IsError() {
		return this.responseErrorToError(response.Status(), response.Body)
	}

	return nil
}

// Delete is deletes document.
//
// ex) err := elasticsearch.Delete("index", "document_id")
func (this *Elasticsearch) Delete(index, documentID string) error {
	if this.client == nil {
		return errors.New("please call Initialize first")
	}

	request := esapi.DeleteRequest{
		Index:      index,
		DocumentID: documentID,
		Refresh:    "true",
		Human:      true,
		ErrorTrace: true,
	}

	response, err := request.Do(context.Background(), this.client)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.IsError() {
		return this.responseErrorToError(response.Status(), response.Body)
	}

	return nil
}

// DeleteByQuery is perform a delete query on index.
//
// ex) err := elasticsearch.DeleteByQuery([]string{"index"},"{...}")
func (this *Elasticsearch) DeleteByQuery(indices []string, body string) error {
	if this.client == nil {
		return errors.New("please call Initialize first")
	}

	var builder strings.Builder
	builder.WriteString(body)

	refresh := true

	request := esapi.DeleteByQueryRequest{
		Index:      indices,
		Body:       strings.NewReader(builder.String()),
		Refresh:    &refresh,
		Human:      true,
		ErrorTrace: true,
	}

	response, err := request.Do(context.Background(), this.client)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.IsError() {
		return this.responseErrorToError(response.Status(), response.Body)
	}

	return nil
}

// IndicesExists is checks if an index exists within indices.
//
// ex) exist, err := elasticsearch.IndicesExists([]string{"index"})
func (this *Elasticsearch) IndicesExists(indices []string) (bool, error) {
	if this.client == nil {
		return false, errors.New("please call Initialize first")
	}

	request := esapi.IndicesExistsRequest{
		Index:      indices,
		Human:      true,
		ErrorTrace: true,
	}

	response, err := request.Do(context.Background(), this.client)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return false, nil
	}

	if response.IsError() {
		return false, this.responseErrorToError(response.Status(), response.Body)
	}

	return true, nil
}

// IndicesCreate is create an index.
//
// ex) err := elasticsearch.IndicesCreate("index", "{...}")
func (this *Elasticsearch) IndicesCreate(index, body string) error {
	if this.client == nil {
		return errors.New("please call Initialize first")
	}

	var builder strings.Builder
	builder.WriteString(body)

	request := esapi.IndicesCreateRequest{
		Index:      index,
		Body:       strings.NewReader(builder.String()),
		Human:      true,
		ErrorTrace: true,
	}

	response, err := request.Do(context.Background(), this.client)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.IsError() {
		return this.responseErrorToError(response.Status(), response.Body)
	}

	return nil
}

// IndicesDelete is delete an index.
//
// ex) err := elasticsearch.IndicesDelete([]string{"Index"})
func (this *Elasticsearch) IndicesDelete(indices []string) error {
	if this.client == nil {
		return errors.New("please call Initialize first")
	}

	request := esapi.IndicesDeleteRequest{
		Index:      indices,
		Human:      true,
		ErrorTrace: true,
	}

	response, err := request.Do(context.Background(), this.client)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.IsError() {
		return this.responseErrorToError(response.Status(), response.Body)
	}

	return nil
}

// IndicesExistsTemplate is checks if a template exists.
//
// ex) exist, err := elasticsearch.IndicesExistsTemplate([]string{"template"})
func (this *Elasticsearch) IndicesExistsTemplate(name []string) (bool, error) {
	if this.client == nil {
		return false, errors.New("please call Initialize first")
	}

	request := esapi.IndicesExistsTemplateRequest{
		Name:       name,
		Human:      true,
		ErrorTrace: true,
	}

	response, err := request.Do(context.Background(), this.client)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return false, nil
	}

	if response.IsError() {
		return false, this.responseErrorToError(response.Status(), response.Body)
	}

	return true, nil
}

// IndicesPutTemplate is stores templates.
//
// ex) err := elasticsearch.IndicesPutTemplate("template", "{...}")
func (this *Elasticsearch) IndicesPutTemplate(name, body string) error {
	if this.client == nil {
		return errors.New("please call Initialize first")
	}

	var builder strings.Builder
	builder.WriteString(body)

	request := esapi.IndicesPutTemplateRequest{
		Name:       name,
		Body:       strings.NewReader(builder.String()),
		Human:      true,
		ErrorTrace: true,
	}

	response, err := request.Do(context.Background(), this.client)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.IsError() {
		return this.responseErrorToError(response.Status(), response.Body)
	}

	return nil
}

// IndicesDeleteTemplate is delete an template.
//
// ex) err := elasticsearch.IndicesDeleteTemplate("template")
func (this *Elasticsearch) IndicesDeleteTemplate(name string) error {
	if this.client == nil {
		return errors.New("please call Initialize first")
	}

	request := esapi.IndicesDeleteTemplateRequest{
		Name:       name,
		Human:      true,
		ErrorTrace: true,
	}

	response, err := request.Do(context.Background(), this.client)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.IsError() {
		return this.responseErrorToError(response.Status(), response.Body)
	}

	return nil
}

// Search is search
//
// ex) result, err := elasticsearch.Search("index", "{...}")
func (this *Elasticsearch) Search(index, body string) (string, error) {
	if this.client == nil {
		return "", errors.New("please call Initialize first")
	}

	var builder strings.Builder
	builder.WriteString(body)

	response, err := this.client.Search(
		this.client.Search.WithContext(context.Background()),
		this.client.Search.WithIndex(index),
		this.client.Search.WithBody(strings.NewReader(builder.String())),
		this.client.Search.WithTrackTotalHits(true),
		this.client.Search.WithPretty(),
	)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.IsError() {
		return "", this.responseErrorToError(response.Status(), response.Body)
	}

	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

// IndicesForcemerge is perform a force merge on index.
//
// ex) err := elasticsearch.IndicesForcemerge([]string{"index"})
func (this *Elasticsearch) IndicesForcemerge(indices []string) error {
	if this.client == nil {
		return errors.New("please call Initialize first")
	}

	onlyExpungeDeletes := true

	request := esapi.IndicesForcemergeRequest{
		Index:              indices,
		OnlyExpungeDeletes: &onlyExpungeDeletes,
		Human:              true,
		ErrorTrace:         true,
	}

	response, err := request.Do(context.Background(), this.client)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.IsError() {
		return this.responseErrorToError(response.Status(), response.Body)
	}

	return nil
}

func (this *Elasticsearch) responseErrorToError(status string, reader io.Reader) error {
	buffer := new(bytes.Buffer)
	buffer.ReadFrom(reader)

	if len(buffer.String()) == 0 {
		return errors.New(fmt.Sprintf("response error - status : (%s)", status))
	}

	return errors.New(fmt.Sprintf("response error - status : (%s), type : (%s), reason : (%s)",
		status,
		gojsonq.New().FromString(buffer.String()).Find("error.type").(string),
		gojsonq.New().FromString(buffer.String()).Find("error.reason").(string)))
}
