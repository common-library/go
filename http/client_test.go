package http_test

import (
	"fmt"
	net_http "net/http"
	"testing"
	"time"

	"github.com/common-library/go/http"
	"github.com/common-library/go/http/mux"
	gorilla_mux "github.com/gorilla/mux"
)

func TestRequestInvalidURL(t *testing.T) {
	t.Parallel()

	_, err := http.Request("invalid_url", net_http.MethodGet, nil, "", 1*time.Second, "", "", nil)
	if err == nil {
		t.Fatal("Expected error for invalid URL")
	}
	if err.Error() != `Get "invalid_url": unsupported protocol scheme ""` {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestRequestGET(t *testing.T) {
	t.Parallel()

	server := &mux.Server{}

	server.RegisterHandlerFunc(net_http.MethodGet, "/test/{id}", func(w net_http.ResponseWriter, r *net_http.Request) {
		vars := gorilla_mux.Vars(r)
		id := vars["id"]
		w.WriteHeader(net_http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"id":"%s","method":"GET"}`, id)))
	})

	if err := server.Start(":30080", func(err error) { t.Fatalf("Server error: %v", err) }); err != nil {
		t.Fatal(err)
	}
	defer server.Stop(5 * time.Second)

	time.Sleep(100 * time.Millisecond)

	response, err := http.Request("http://localhost:30080/test/id-01", net_http.MethodGet, nil, "", 10*time.Second, "", "", nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if response.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", response.StatusCode)
	}

	expected := `{"id":"id-01","method":"GET"}`
	if response.Body != expected {
		t.Errorf("Expected body %q, got %q", expected, response.Body)
	}
}

func TestRequestPOST(t *testing.T) {
	t.Parallel()

	server := &mux.Server{}

	server.RegisterHandlerFunc(net_http.MethodPost, "/users", func(w net_http.ResponseWriter, r *net_http.Request) {
		if r.Method != net_http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		w.WriteHeader(net_http.StatusCreated)
		w.Write([]byte(`{"status":"created"}`))
	})

	if err := server.Start(":30081", func(err error) { t.Fatalf("Server error: %v", err) }); err != nil {
		t.Fatal(err)
	}
	defer server.Stop(5 * time.Second)

	time.Sleep(100 * time.Millisecond)

	body := `{"name":"Alice","email":"alice@example.com"}`
	response, err := http.Request("http://localhost:30081/users", net_http.MethodPost, nil, body, 10*time.Second, "", "", nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if response.StatusCode != 201 {
		t.Errorf("Expected status 201, got %d", response.StatusCode)
	}

	expected := `{"status":"created"}`
	if response.Body != expected {
		t.Errorf("Expected body %q, got %q", expected, response.Body)
	}
}

func TestRequestWithHeaders(t *testing.T) {
	t.Parallel()

	server := &mux.Server{}

	server.RegisterHandlerFunc(net_http.MethodGet, "/headers", func(w net_http.ResponseWriter, r *net_http.Request) {
		contentType := r.Header.Get("Content-Type")
		apiKey := r.Header.Get("X-API-Key")
		customHeader := r.Header.Get("X-Custom-Header")

		w.WriteHeader(net_http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"content-type":"%s","api-key":"%s","custom":"%s"}`,
			contentType, apiKey, customHeader)))
	})

	if err := server.Start(":30082", func(err error) { t.Fatalf("Server error: %v", err) }); err != nil {
		t.Fatal(err)
	}
	defer server.Stop(5 * time.Second)

	time.Sleep(100 * time.Millisecond)

	headers := map[string][]string{
		"Content-Type":    {"application/json"},
		"X-API-Key":       {"secret123"},
		"X-Custom-Header": {"value1"},
	}

	response, err := http.Request("http://localhost:30082/headers", net_http.MethodGet, headers, "", 10*time.Second, "", "", nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	expected := `{"content-type":"application/json","api-key":"secret123","custom":"value1"}`
	if response.Body != expected {
		t.Errorf("Expected body %q, got %q", expected, response.Body)
	}
}

func TestRequestWithMultipleHeaderValues(t *testing.T) {
	t.Parallel()

	server := &mux.Server{}

	server.RegisterHandlerFunc(net_http.MethodGet, "/multi-headers", func(w net_http.ResponseWriter, r *net_http.Request) {
		values := r.Header["X-Multi"]
		w.WriteHeader(net_http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"count":%d,"values":["%s","%s"]}`,
			len(values), values[0], values[1])))
	})

	if err := server.Start(":30083", func(err error) { t.Fatalf("Server error: %v", err) }); err != nil {
		t.Fatal(err)
	}
	defer server.Stop(5 * time.Second)

	time.Sleep(100 * time.Millisecond)

	headers := map[string][]string{
		"X-Multi": {"value1", "value2"},
	}

	response, err := http.Request("http://localhost:30083/multi-headers", net_http.MethodGet, headers, "", 10*time.Second, "", "", nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	expected := `{"count":2,"values":["value1","value2"]}`
	if response.Body != expected {
		t.Errorf("Expected body %q, got %q", expected, response.Body)
	}
}

func TestRequestWithBasicAuth(t *testing.T) {
	t.Parallel()

	server := &mux.Server{}

	server.RegisterHandlerFunc(net_http.MethodGet, "/auth", func(w net_http.ResponseWriter, r *net_http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			w.WriteHeader(net_http.StatusUnauthorized)
			w.Write([]byte(`{"error":"no auth"}`))
			return
		}
		w.WriteHeader(net_http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password)))
	})

	if err := server.Start(":30084", func(err error) { t.Fatalf("Server error: %v", err) }); err != nil {
		t.Fatal(err)
	}
	defer server.Stop(5 * time.Second)

	time.Sleep(100 * time.Millisecond)

	response, err := http.Request("http://localhost:30084/auth", net_http.MethodGet, nil, "", 10*time.Second, "admin", "password123", nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if response.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", response.StatusCode)
	}

	expected := `{"username":"admin","password":"password123"}`
	if response.Body != expected {
		t.Errorf("Expected body %q, got %q", expected, response.Body)
	}
}

func TestRequestWithoutBasicAuth(t *testing.T) {
	t.Parallel()

	server := &mux.Server{}

	server.RegisterHandlerFunc(net_http.MethodGet, "/no-auth", func(w net_http.ResponseWriter, r *net_http.Request) {
		_, _, ok := r.BasicAuth()
		if ok {
			t.Error("Should not have basic auth")
		}
		w.WriteHeader(net_http.StatusOK)
		w.Write([]byte(`{"auth":"none"}`))
	})

	if err := server.Start(":30085", func(err error) { t.Fatalf("Server error: %v", err) }); err != nil {
		t.Fatal(err)
	}
	defer server.Stop(5 * time.Second)

	time.Sleep(100 * time.Millisecond)

	response, err := http.Request("http://localhost:30085/no-auth", net_http.MethodGet, nil, "", 10*time.Second, "", "", nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	expected := `{"auth":"none"}`
	if response.Body != expected {
		t.Errorf("Expected body %q, got %q", expected, response.Body)
	}
}

func TestRequestPUT(t *testing.T) {
	t.Parallel()

	server := &mux.Server{}

	server.RegisterHandlerFunc(net_http.MethodPut, "/users/{id}", func(w net_http.ResponseWriter, r *net_http.Request) {
		vars := gorilla_mux.Vars(r)
		id := vars["id"]
		w.WriteHeader(net_http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"id":"%s","method":"PUT"}`, id)))
	})

	if err := server.Start(":30086", func(err error) { t.Fatalf("Server error: %v", err) }); err != nil {
		t.Fatal(err)
	}
	defer server.Stop(5 * time.Second)

	time.Sleep(100 * time.Millisecond)

	body := `{"name":"Updated Name"}`
	response, err := http.Request("http://localhost:30086/users/123", net_http.MethodPut, nil, body, 10*time.Second, "", "", nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	expected := `{"id":"123","method":"PUT"}`
	if response.Body != expected {
		t.Errorf("Expected body %q, got %q", expected, response.Body)
	}
}

func TestRequestDELETE(t *testing.T) {
	t.Parallel()

	server := &mux.Server{}

	server.RegisterHandlerFunc(net_http.MethodDelete, "/users/{id}", func(w net_http.ResponseWriter, r *net_http.Request) {
		vars := gorilla_mux.Vars(r)
		id := vars["id"]
		w.WriteHeader(net_http.StatusNoContent)
		w.Write([]byte(fmt.Sprintf(`{"id":"%s","deleted":true}`, id)))
	})

	if err := server.Start(":30087", func(err error) { t.Fatalf("Server error: %v", err) }); err != nil {
		t.Fatal(err)
	}
	defer server.Stop(5 * time.Second)

	time.Sleep(100 * time.Millisecond)

	response, err := http.Request("http://localhost:30087/users/456", net_http.MethodDelete, nil, "", 10*time.Second, "", "", nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if response.StatusCode != 204 {
		t.Errorf("Expected status 204, got %d", response.StatusCode)
	}
}

func TestRequestResponseHeaders(t *testing.T) {
	t.Parallel()

	server := &mux.Server{}

	server.RegisterHandlerFunc(net_http.MethodGet, "/response-headers", func(w net_http.ResponseWriter, r *net_http.Request) {
		w.Header().Set("X-Custom-Response", "custom-value")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(net_http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	if err := server.Start(":30088", func(err error) { t.Fatalf("Server error: %v", err) }); err != nil {
		t.Fatal(err)
	}
	defer server.Stop(5 * time.Second)

	time.Sleep(100 * time.Millisecond)

	response, err := http.Request("http://localhost:30088/response-headers", net_http.MethodGet, nil, "", 10*time.Second, "", "", nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if response.Header.Get("X-Custom-Response") != "custom-value" {
		t.Errorf("Expected header 'X-Custom-Response: custom-value', got %q",
			response.Header.Get("X-Custom-Response"))
	}

	if response.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got %q",
			response.Header.Get("Content-Type"))
	}
}

func TestRequestWithMiddleware(t *testing.T) {
	t.Parallel()

	server := &mux.Server{}

	middlewareCalled := false
	middleware := func(next net_http.Handler) net_http.Handler {
		return net_http.HandlerFunc(func(w net_http.ResponseWriter, r *net_http.Request) {
			middlewareCalled = true
			w.Header().Set("X-Middleware", "applied")
			next.ServeHTTP(w, r)
		})
	}

	server.Use(middleware)
	server.RegisterHandlerFunc(net_http.MethodGet, "/test", func(w net_http.ResponseWriter, r *net_http.Request) {
		w.WriteHeader(net_http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	if err := server.Start(":30089", func(err error) { t.Fatalf("Server error: %v", err) }); err != nil {
		t.Fatal(err)
	}
	defer server.Stop(5 * time.Second)

	time.Sleep(100 * time.Millisecond)

	response, err := http.Request("http://localhost:30089/test", net_http.MethodGet, nil, "", 10*time.Second, "", "", nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if !middlewareCalled {
		t.Error("Middleware was not called")
	}

	if response.Header.Get("X-Middleware") != "applied" {
		t.Error("Middleware did not set header")
	}
}
