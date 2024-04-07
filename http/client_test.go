package http_test

import (
	"math/rand/v2"
	net_http "net/http"
	"strconv"
	"testing"
	"time"

	"github.com/common-library/go/http"
)

func TestRequest1(t *testing.T) {
	_, err := http.Request("invalid_url", net_http.MethodGet, nil, "", 1, "", "")
	if err.Error() != `Get "invalid_url": unsupported protocol scheme ""` {
		t.Error(err)
	}
}

func TestRequest2(t *testing.T) {
	address := ":" + strconv.Itoa(10000+rand.IntN(10000))

	server := http.Server{}

	server.AddHandler("/test/{id}", net_http.MethodGet, func(responseWriter net_http.ResponseWriter, request *net_http.Request) {
		responseWriter.WriteHeader(net_http.StatusOK)
		responseWriter.Write([]byte(`{"field_1":1}`))
	})

	middlewareFunction := func(nextHandler net_http.Handler) net_http.Handler {
		return net_http.HandlerFunc(func(responseWriter net_http.ResponseWriter, request *net_http.Request) {
			nextHandler.ServeHTTP(responseWriter, request)
		})
	}

	err := server.Start(address, func(err error) { t.Error(err) }, middlewareFunction)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Duration(200) * time.Millisecond)

	{
		response, err := http.Request("http://"+address+"/test/id-01", net_http.MethodGet, map[string][]string{"header-1": {"value-1"}}, "", 3, "", "")
		if err != nil {
			t.Error(err)
		} else if response.StatusCode != 200 {
			t.Errorf("invalid status code : (%d)", response.StatusCode)
		} else if response.Body != `{"field_1":1}` {
			t.Errorf("invalid response body : (%s)", response.Body)
		}
	}

	err = server.Stop(5)
	if err != nil {
		t.Error(err)
	}
}
