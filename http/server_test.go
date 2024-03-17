package http_test

import (
	"math/rand/v2"
	net_http "net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/heaven-chp/common-library-go/http"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

var address string

func setUp(server *http.Server) {
	address = ":" + strconv.Itoa(10000+rand.IntN(1000))

	server.AddPathPrefixHandler("/swagger/", httpSwagger.WrapHandler)

	server.AddHandler("/test/{id}", net_http.MethodGet, func(responseWriter net_http.ResponseWriter, request *net_http.Request) {
		responseWriter.WriteHeader(net_http.StatusOK)
		responseWriter.Write([]byte(`{"field_1":1}`))
	})

	middlewareFunction := func(nextHandler net_http.Handler) net_http.Handler {
		return net_http.HandlerFunc(func(responseWriter net_http.ResponseWriter, request *net_http.Request) {
			nextHandler.ServeHTTP(responseWriter, request)
		})
	}

	err := server.Start(address, func(err error) { panic(err) }, middlewareFunction)
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Duration(500) * time.Millisecond)
}

func tearDown(server *http.Server) {
	err := server.Stop(5)
	if err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	server := http.Server{}

	setUp(&server)

	code := m.Run()

	tearDown(&server)

	os.Exit(code)
}

func TestAddHandler(t *testing.T) {
	response, err := http.Request("http://"+address+"/test/id-01", net_http.MethodGet, nil, "", 3, "", "")
	if err != nil {
		t.Error(err)
	} else if response.StatusCode != 200 {
		t.Errorf("invalid status code : (%d)", response.StatusCode)
	} else if response.Body != `{"field_1":1}` {
		t.Errorf("invalid response body : (%s)", response.Body)
	}
}

func TestAddPathPrefixHandler(t *testing.T) {
	response, err := http.Request("http://"+address+"/swagger/index.html", net_http.MethodGet, nil, "", 3, "", "")
	if err != nil {
		t.Error(err)
	} else if response.StatusCode != net_http.StatusOK {
		t.Errorf("invalid status code : (%d)", response.StatusCode)
	}
}

func TestStart1(t *testing.T) {
	listenAndServeFailureFunc := func(err error) {
		if err.Error() != "listen tcp "+address+": bind: address already in use" {
			t.Fatal(err)
		}
	}

	server := http.Server{}
	err := server.Start(address, listenAndServeFailureFunc)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStart2(t *testing.T) {
	response, err := http.Request("http://"+address+"/test/id-01", net_http.MethodGet, nil, "", 3, "", "")
	if err != nil {
		t.Error(err)
	} else if response.StatusCode != net_http.StatusOK {
		t.Errorf("invalid status code : (%d)", response.StatusCode)
	} else if response.Body != `{"field_1":1}` {
		t.Errorf("invalid response body : (%s)", response.Body)
	}
}

func TestStop(t *testing.T) {
	server := http.Server{}

	err := server.Stop(5)
	if err != nil {
		t.Error(err)
	}
}

func TestSetRouter(t *testing.T) {
	router := mux.NewRouter()

	server := http.Server{}

	server.SetRouter(router)
}
