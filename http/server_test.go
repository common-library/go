package http_test

import (
	"math/rand/v2"
	net_http "net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/common-library/go/http"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

var address string

func setUp(server *http.Server) {
	address = ":" + strconv.Itoa(10000+rand.IntN(1000))

	server.RegisterPathPrefixHandler("/swagger/", httpSwagger.WrapHandler)

	server.RegisterHandlerFunc("/test/{id}", net_http.MethodGet, func(responseWriter net_http.ResponseWriter, request *net_http.Request) {
		responseWriter.WriteHeader(net_http.StatusOK)
		responseWriter.Write([]byte(`{"field_1":1}`))
	})

	middlewareFunction := func(nextHandler net_http.Handler) net_http.Handler {
		return net_http.HandlerFunc(func(responseWriter net_http.ResponseWriter, request *net_http.Request) {
			nextHandler.ServeHTTP(responseWriter, request)
		})
	}

	if err := server.Start(address, func(err error) { panic(err) }, middlewareFunction); err != nil {
		panic(err)
	}
	time.Sleep(200 * time.Millisecond)
}

func tearDown(server *http.Server) {
	if err := server.Stop(10); err != nil {
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

func TestRegisterHandlerFunc(t *testing.T) {
	if response, err := http.Request("http://"+address+"/test/id-01", net_http.MethodGet, nil, "", 10, "", "", nil); err != nil {
		t.Fatal(err)
	} else if response.StatusCode != 200 {
		t.Fatal("invalid -", response.StatusCode)
	} else if response.Body != `{"field_1":1}` {
		t.Fatal("invalid -", response.Body)
	}
}

func TestRegisterPathPrefixHandler(t *testing.T) {
	if response, err := http.Request("http://"+address+"/swagger/index.html", net_http.MethodGet, nil, "", 10, "", "", nil); err != nil {
		t.Fatal(err)
	} else if response.StatusCode != net_http.StatusOK {
		t.Fatal("invalid -", response.StatusCode)
	}
}

func TestStart1(t *testing.T) {
	listenAndServeFailureFunc := func(err error) {
		if err.Error() != "listen tcp "+address+": bind: address already in use" {
			t.Fatal(err)
		}
	}

	server := http.Server{}
	if err := server.Start(address, listenAndServeFailureFunc); err != nil {
		t.Fatal(err)
	}
	time.Sleep(200 * time.Millisecond)

	if err := server.Stop(10); err != nil {
		panic(err)
	}
}

func TestStart2(t *testing.T) {
	listenAndServeFailureFunc := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}

	server := http.Server{}
	addressTemp := ":" + strconv.Itoa(11000+rand.IntN(1000))
	if err := server.Start(addressTemp, listenAndServeFailureFunc); err != nil {
		t.Fatal(err)
	}
	time.Sleep(200 * time.Millisecond)

	if err := server.Stop(10); err != nil {
		panic(err)
	}
}

func TestStop(t *testing.T) {
	server := http.Server{}

	if err := server.Stop(10); err != nil {
		t.Fatal(err)
	}
}

func TestSetRouter(t *testing.T) {
	router := mux.NewRouter()

	server := http.Server{}

	server.SetRouter(router)
}
