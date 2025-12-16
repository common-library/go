// Package http provides http client and server implementations.
package http

import (
	"io"
	"net/http"
	"strings"
	"time"
)

// Response is response information.
type Response struct {
	Header     http.Header
	Body       string
	StatusCode int
}

// Request is request.
//
// ex) response, err := http.Request("http://127.0.0.1:10000/test/id-01", http.MethodGet, map[string][]string{"header-1": {"value-1"}}, "", 10, "", "", nil)
func Request(url, method string, header map[string][]string, body string, timeout time.Duration, username, password string, transport *http.Transport) (Response, error) {
	if request, err := getRequest(url, method, header, body, username, password); err != nil {
		return Response{}, err
	} else {
		return getResponse(request, timeout, transport)
	}
}

func getRequest(url, method string, header map[string][]string, body string, username, password string) (*http.Request, error) {
	if request, err := http.NewRequest(method, url, strings.NewReader(body)); err != nil {
		return nil, err
	} else {
		request.SetBasicAuth(username, password)

		for key, array := range header {
			for _, value := range array {
				request.Header.Add(key, value)
			}
		}

		return request, nil
	}
}

func getResponse(request *http.Request, timeout time.Duration, transport *http.Transport) (Response, error) {
	if transport == nil {
		transport = &http.Transport{}
	}

	client := http.Client{Transport: transport, Timeout: timeout * time.Second}

	if response, err := client.Do(request); err != nil {
		return Response{}, err
	} else {
		defer response.Body.Close()

		if responseBody, err := io.ReadAll(response.Body); err != nil {
			return Response{}, err
		} else {
			return Response{Header: response.Header, Body: string(responseBody), StatusCode: response.StatusCode}, nil
		}
	}
}
