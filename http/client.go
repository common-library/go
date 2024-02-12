// Package http provides http client and server implementations.
package http

import (
	"io/ioutil"
	net_http "net/http"
	"strings"
	"time"
)

// Response is response information.
type Response struct {
	Header     net_http.Header
	Body       string
	StatusCode int
}

// Request is request.
//
// ex) response, err := http.Request("http://127.0.0.1:10000/test/id-01", net_http.MethodGet, map[string][]string{"header-1": {"value-1"}}, "", 3, "", "")
func Request(url, method string, header map[string][]string, body string, timeout int, username, password string) (Response, error) {
	client := &net_http.Client{
		Transport: &net_http.Transport{},
		Timeout:   time.Duration(timeout) * time.Second,
	}

	request, err := net_http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return Response{}, err
	}

	request.SetBasicAuth(username, password)

	for key, array := range header {
		for _, value := range array {
			request.Header.Add(key, value)
		}
	}

	response, err := client.Do(request)
	if err != nil {
		return Response{}, err
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return Response{}, err
	}

	return Response{Header: response.Header, Body: string(responseBody), StatusCode: response.StatusCode}, nil
}
