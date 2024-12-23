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
)

var address string

type handler struct {
}

func (this handler) ServeHTTP(w net_http.ResponseWriter, r *net_http.Request) {
	w.WriteHeader(net_http.StatusOK)
	w.Write([]byte(`{"field_1":1}`))
}

func setUp(server *http.Server) {
	address = ":" + strconv.Itoa(10000+rand.IntN(1000))

	server.RegisterHandler("/test-01/{id}", handler{})
	server.RegisterHandlerFunc("/test-02/{id}", handler{}.ServeHTTP, net_http.MethodGet)

	server.RegisterPathPrefixHandler("/test-03", handler{}, net_http.MethodGet)
	server.RegisterPathPrefixHandlerFunc("/test-04", handler{}.ServeHTTP, net_http.MethodGet, net_http.MethodPost)

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
	if err := server.Stop(10 * time.Second); err != nil {
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

func TestRegisterHandler(t *testing.T) {
	t.Parallel()

	if response, err := http.Request("http://"+address+"/test-01/id-01", net_http.MethodGet, nil, "", 10, "", "", nil); err != nil {
		t.Fatal(err)
	} else if response.StatusCode != 200 {
		t.Fatal("invalid -", response.StatusCode)
	} else if response.Body != `{"field_1":1}` {
		t.Fatal("invalid -", response.Body)
	}
}

func TestRegisterHandlerFunc(t *testing.T) {
	t.Parallel()

	if response, err := http.Request("http://"+address+"/test-02/id-01", net_http.MethodGet, nil, "", 10, "", "", nil); err != nil {
		t.Fatal(err)
	} else if response.StatusCode != 200 {
		t.Fatal("invalid -", response.StatusCode)
	} else if response.Body != `{"field_1":1}` {
		t.Fatal("invalid -", response.Body)
	}
}

func TestRegisterPathPrefixHandler(t *testing.T) {
	t.Parallel()

	if response, err := http.Request("http://"+address+"/test-03", net_http.MethodGet, nil, "", 10, "", "", nil); err != nil {
		t.Fatal(err)
	} else if response.StatusCode != 200 {
		t.Fatal("invalid -", response.StatusCode)
	} else if response.Body != `{"field_1":1}` {
		t.Fatal("invalid -", response.Body)
	}
}

func TestRegisterPathPrefixHandlerFunc(t *testing.T) {
	t.Parallel()

	if response, err := http.Request("http://"+address+"/test-04", net_http.MethodGet, nil, "", 10, "", "", nil); err != nil {
		t.Fatal(err)
	} else if response.StatusCode != 200 {
		t.Fatal("invalid -", response.StatusCode)
	} else if response.Body != `{"field_1":1}` {
		t.Fatal("invalid -", response.Body)
	}
}

func TestStart1(t *testing.T) {
	t.Parallel()

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

	if err := server.Stop(10 * time.Second); err != nil {
		panic(err)
	}
}

func TestStart2(t *testing.T) {
	t.Parallel()

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

	if err := server.Stop(10 * time.Second); err != nil {
		panic(err)
	}
}

func TestStop(t *testing.T) {
	t.Parallel()

	server := http.Server{}

	if err := server.Stop(10 * time.Second); err != nil {
		t.Fatal(err)
	}
}

func TestSetRouter(t *testing.T) {
	t.Parallel()

	router := mux.NewRouter()

	server := http.Server{}

	server.SetRouter(router)
}
