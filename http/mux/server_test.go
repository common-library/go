package mux_test

import (
	"fmt"
	"io"
	net_http "net/http"
	"testing"
	"time"

	"github.com/common-library/go/http/mux"
	gorilla_mux "github.com/gorilla/mux"
)

func TestRegisterHandler(t *testing.T) {
	t.Parallel()

	server := &mux.Server{}
	handlerCalled := false

	handler := net_http.HandlerFunc(func(w net_http.ResponseWriter, r *net_http.Request) {
		handlerCalled = true
		w.WriteHeader(net_http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	server.RegisterHandlerAny("/test/{id}", handler)
	server.Start(":20080", func(err error) {
		t.Errorf("Server error: %v", err)
	})
	defer server.Stop(5 * time.Second)

	time.Sleep(100 * time.Millisecond)

	resp, err := net_http.Get("http://localhost:20080/test/123")
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != net_http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != `{"status":"ok"}` {
		t.Errorf("Expected body %q, got %q", `{"status":"ok"}`, string(body))
	}

	if !handlerCalled {
		t.Error("Handler was not called")
	}
}

func TestRegisterHandlerFunc(t *testing.T) {
	t.Parallel()

	server := &mux.Server{}

	server.RegisterHandlerFunc(net_http.MethodGet, "/test/{id}", func(w net_http.ResponseWriter, r *net_http.Request) {
		vars := gorilla_mux.Vars(r)
		id := vars["id"]
		w.WriteHeader(net_http.StatusOK)
		fmt.Fprintf(w, `{"id":"%s"}`, id)
	})

	server.Start(":20081", func(err error) {
		t.Errorf("Server error: %v", err)
	})
	defer server.Stop(5 * time.Second)

	time.Sleep(100 * time.Millisecond)

	resp, err := net_http.Get("http://localhost:20081/test/456")
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	expected := `{"id":"456"}`
	if string(body) != expected {
		t.Errorf("Expected body %q, got %q", expected, string(body))
	}
}

func TestRegisterHandlerFuncMethods(t *testing.T) {
	t.Parallel()

	server := &mux.Server{}

	// Register GET handler
	server.RegisterHandlerFunc(net_http.MethodGet, "/data", func(w net_http.ResponseWriter, r *net_http.Request) {
		w.WriteHeader(net_http.StatusOK)
		w.Write([]byte(`{"method":"` + r.Method + `"}`))
	})

	// Register POST handler
	server.RegisterHandlerFunc(net_http.MethodPost, "/data", func(w net_http.ResponseWriter, r *net_http.Request) {
		w.WriteHeader(net_http.StatusOK)
		w.Write([]byte(`{"method":"` + r.Method + `"}`))
	})

	server.Start(":20082", func(err error) {
		t.Errorf("Server error: %v", err)
	})
	defer server.Stop(5 * time.Second)

	time.Sleep(100 * time.Millisecond)

	// Test GET
	resp, _ := net_http.Get("http://localhost:20082/data")
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if string(body) != `{"method":"GET"}` {
		t.Errorf("GET failed: %s", string(body))
	}

	// Test POST
	resp, _ = net_http.Post("http://localhost:20082/data", "text/plain", nil)
	body, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	if string(body) != `{"method":"POST"}` {
		t.Errorf("POST failed: %s", string(body))
	}

	// Test PUT (should fail - method not allowed)
	req, _ := net_http.NewRequest(net_http.MethodPut, "http://localhost:20082/data", nil)
	resp, _ = net_http.DefaultClient.Do(req)
	resp.Body.Close()
	if resp.StatusCode != net_http.StatusMethodNotAllowed {
		t.Errorf("PUT should fail with 405, got %d", resp.StatusCode)
	}
}

func TestRegisterPathPrefixHandler(t *testing.T) {
	t.Parallel()

	server := &mux.Server{}

	handler := net_http.HandlerFunc(func(w net_http.ResponseWriter, r *net_http.Request) {
		w.WriteHeader(net_http.StatusOK)
		w.Write([]byte(`{"prefix":"matched"}`))
	})

	server.RegisterPathPrefixHandler(net_http.MethodGet, "/api/", handler)
	server.Start(":20083", func(err error) {
		t.Errorf("Server error: %v", err)
	})
	defer server.Stop(5 * time.Second)

	time.Sleep(100 * time.Millisecond)

	// Should match /api/users
	resp, err := net_http.Get("http://localhost:20083/api/users")
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != net_http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Should also match /api/posts
	resp2, _ := net_http.Get("http://localhost:20083/api/posts")
	defer resp2.Body.Close()
	if resp2.StatusCode != net_http.StatusOK {
		t.Errorf("Prefix should match /api/posts")
	}
}

func TestRegisterPathPrefixHandlerFunc(t *testing.T) {
	t.Parallel()

	server := &mux.Server{}

	server.RegisterPathPrefixHandlerFuncAny("/static/", func(w net_http.ResponseWriter, r *net_http.Request) {
		w.WriteHeader(net_http.StatusOK)
		w.Write([]byte(`static content`))
	})

	server.Start(":20084", func(err error) {
		t.Errorf("Server error: %v", err)
	})
	defer server.Stop(5 * time.Second)

	time.Sleep(100 * time.Millisecond)

	resp, err := net_http.Get("http://localhost:20084/static/file.js")
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if string(body) != `static content` {
		t.Errorf("Expected 'static content', got %q", string(body))
	}
}

func TestUseMiddleware(t *testing.T) {
	t.Parallel()

	server := &mux.Server{}

	middlewareCalled := false
	middleware := func(next net_http.Handler) net_http.Handler {
		return net_http.HandlerFunc(func(w net_http.ResponseWriter, r *net_http.Request) {
			middlewareCalled = true
			w.Header().Set("X-Custom-Header", "middleware")
			next.ServeHTTP(w, r)
		})
	}

	server.Use(middleware)
	server.RegisterHandlerFuncAny("/test", func(w net_http.ResponseWriter, r *net_http.Request) {
		w.WriteHeader(net_http.StatusOK)
		w.Write([]byte("ok"))
	})

	server.Start(":20085", func(err error) {
		t.Errorf("Server error: %v", err)
	})
	defer server.Stop(5 * time.Second)

	time.Sleep(100 * time.Millisecond)

	resp, err := net_http.Get("http://localhost:20085/test")
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if !middlewareCalled {
		t.Error("Middleware was not called")
	}

	if resp.Header.Get("X-Custom-Header") != "middleware" {
		t.Error("Middleware did not set header")
	}
}

func TestStartStop(t *testing.T) {
	t.Parallel()

	server := &mux.Server{}

	server.RegisterHandlerFuncAny("/", func(w net_http.ResponseWriter, r *net_http.Request) {
		w.WriteHeader(net_http.StatusOK)
	})

	err := server.Start(":20086", func(err error) {
		t.Errorf("Server error: %v", err)
	})
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	if !server.IsRunning() {
		t.Error("Server should be running")
	}

	err = server.Stop(5 * time.Second)
	if err != nil {
		t.Errorf("Stop failed: %v", err)
	}

	if server.IsRunning() {
		t.Error("Server should not be running after Stop")
	}
}

func TestAlreadyStarted(t *testing.T) {
	t.Parallel()

	server := &mux.Server{}

	server.RegisterHandlerFuncAny("/", func(w net_http.ResponseWriter, r *net_http.Request) {
		w.WriteHeader(net_http.StatusOK)
	})

	err := server.Start(":20087", nil)
	if err != nil {
		t.Fatalf("First start failed: %v", err)
	}
	defer server.Stop(5 * time.Second)

	time.Sleep(100 * time.Millisecond)

	// Try to start again
	err = server.Start(":20088", nil)
	if err == nil {
		t.Error("Second start should fail with 'server already started' error")
	}
	if err != nil && err.Error() != "server already started" {
		t.Errorf("Expected 'server already started' error, got: %v", err)
	}
}

func TestGetRouter(t *testing.T) {
	t.Parallel()

	server := &mux.Server{}

	router := server.GetRouter()
	if router == nil {
		t.Fatal("GetRouter returned nil")
	}

	// Configure router directly
	router.StrictSlash(true)
	router.HandleFunc("/custom", func(w net_http.ResponseWriter, r *net_http.Request) {
		w.WriteHeader(net_http.StatusOK)
		w.Write([]byte("custom route"))
	})

	server.Start(":20089", func(err error) {
		t.Errorf("Server error: %v", err)
	})
	defer server.Stop(5 * time.Second)

	time.Sleep(100 * time.Millisecond)

	resp, err := net_http.Get("http://localhost:20089/custom")
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "custom route" {
		t.Errorf("Expected 'custom route', got %q", string(body))
	}
}

func TestIsRunning(t *testing.T) {
	t.Parallel()

	server := &mux.Server{}

	if server.IsRunning() {
		t.Error("New server should not be running")
	}

	server.RegisterHandlerFuncAny("/", func(w net_http.ResponseWriter, r *net_http.Request) {
		w.WriteHeader(net_http.StatusOK)
	})

	server.Start(":20090", nil)
	time.Sleep(100 * time.Millisecond)

	if !server.IsRunning() {
		t.Error("Server should be running after Start")
	}

	server.Stop(5 * time.Second)

	if server.IsRunning() {
		t.Error("Server should not be running after Stop")
	}
}
