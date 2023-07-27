package http

import (
	"io/ioutil"
	net_http "net/http"
	"strings"
	"time"
)

type Response struct {
	Header     net_http.Header
	Body       string
	StatusCode int
}

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
