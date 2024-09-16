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
	t.Parallel()

	_, err := http.Request("invalid_url", net_http.MethodGet, nil, "", 1, "", "", nil)
	if err.Error() != `Get "invalid_url": unsupported protocol scheme ""` {
		t.Fatal(err)
	}
}

func TestRequest2(t *testing.T) {
	t.Parallel()

	address := ":" + strconv.Itoa(10000+rand.IntN(10000))

	server := http.Server{}

	server.RegisterHandlerFunc("/test/{id}", net_http.MethodGet, func(responseWriter net_http.ResponseWriter, request *net_http.Request) {
		responseWriter.WriteHeader(net_http.StatusOK)
		responseWriter.Write([]byte(`{"field_1":1}`))
	})

	middlewareFunction := func(nextHandler net_http.Handler) net_http.Handler {
		return net_http.HandlerFunc(func(responseWriter net_http.ResponseWriter, request *net_http.Request) {
			nextHandler.ServeHTTP(responseWriter, request)
		})
	}

	listenAndServeFailureFunc := func(err error) { t.Fatal(err) }
	if err := server.Start(address, listenAndServeFailureFunc, middlewareFunction); err != nil {
		t.Fatal(err)
	}
	time.Sleep(200 * time.Millisecond)

	if response, err := http.Request("http://"+address+"/test/id-01", net_http.MethodGet, map[string][]string{"header-1": {"value-1"}}, "", 10, "", "", nil); err != nil {
		t.Fatal(err)
	} else if response.StatusCode != 200 {
		t.Fatal(response.StatusCode)
	} else if response.Body != `{"field_1":1}` {
		t.Fatal(response.Body)
	}

	if err := server.Stop(10); err != nil {
		t.Fatal(err)
	}
}
