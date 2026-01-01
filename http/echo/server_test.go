package echo_test

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/common-library/go/http/echo"
	echo_lib "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func TestRegisterHandler(t *testing.T) {
	t.Parallel()

	server := &echo.Server{}
	defer server.Stop(5 * time.Second)

	server.RegisterHandler(echo_lib.GET, "/test", func(c echo_lib.Context) error {
		return c.String(http.StatusOK, "test response")
	})

	if err := server.Start(":18080", nil); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:18080/test")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "test response" {
		t.Fatalf("expected 'test response', got '%s'", body)
	}
}

func TestRegisterHandlerMethods(t *testing.T) {
	t.Parallel()

	server := &echo.Server{}
	defer server.Stop(5 * time.Second)

	// Test GET
	server.RegisterHandler(echo_lib.GET, "/get", func(c echo_lib.Context) error {
		return c.String(http.StatusOK, "GET response")
	})

	// Test POST
	server.RegisterHandler(echo_lib.POST, "/post", func(c echo_lib.Context) error {
		return c.String(http.StatusCreated, "POST response")
	})

	// Test PUT
	server.RegisterHandler(echo_lib.PUT, "/put", func(c echo_lib.Context) error {
		return c.String(http.StatusOK, "PUT response")
	})

	// Test DELETE
	server.RegisterHandler(echo_lib.DELETE, "/delete", func(c echo_lib.Context) error {
		return c.NoContent(http.StatusNoContent)
	})

	// Test PATCH
	server.RegisterHandler(echo_lib.PATCH, "/patch", func(c echo_lib.Context) error {
		return c.String(http.StatusOK, "PATCH response")
	})

	if err := server.Start(":18081", nil); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	// Test GET request
	resp, err := http.Get("http://localhost:18081/get")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET: expected 200, got %d", resp.StatusCode)
	}

	// Test POST request
	resp, err = http.Post("http://localhost:18081/post", "text/plain", strings.NewReader("data"))
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("POST: expected 201, got %d", resp.StatusCode)
	}
}

func TestUseMiddleware(t *testing.T) {
	t.Parallel()

	server := &echo.Server{}
	defer server.Stop(5 * time.Second)

	// Add global middleware
	server.Use(middleware.Recover())

	server.RegisterHandler(echo_lib.GET, "/panic", func(c echo_lib.Context) error {
		panic("test panic")
	})

	if err := server.Start(":18082", nil); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	// Should recover from panic
	resp, err := http.Get("http://localhost:18082/panic")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Recover middleware returns 500
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", resp.StatusCode)
	}
}

func TestGroup(t *testing.T) {
	t.Parallel()

	server := &echo.Server{}
	defer server.Stop(5 * time.Second)

	// Create API group
	api := server.Group("/api")
	api.GET("/users", func(c echo_lib.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "users"})
	})

	v1 := api.Group("/v1")
	v1.GET("/info", func(c echo_lib.Context) error {
		return c.String(http.StatusOK, "v1 info")
	})

	if err := server.Start(":18083", nil); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	// Test /api/users
	resp, err := http.Get("http://localhost:18083/api/users")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	// Test /api/v1/info
	resp, err = http.Get("http://localhost:18083/api/v1/info")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestStartStop(t *testing.T) {
	t.Parallel()

	server := &echo.Server{}

	server.RegisterHandler(echo_lib.GET, "/test", func(c echo_lib.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	// Server should not be running initially
	if server.IsRunning() {
		t.Fatal("server should not be running initially")
	}

	// Start server
	if err := server.Start(":18084", nil); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	// Server should be running
	if !server.IsRunning() {
		t.Fatal("server should be running")
	}

	// Test server is responding
	resp, err := http.Get("http://localhost:18084/test")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	// Stop server
	if err := server.Stop(5 * time.Second); err != nil {
		t.Fatal(err)
	}

	// Server should not be running
	time.Sleep(100 * time.Millisecond)
	if server.IsRunning() {
		t.Fatal("server should not be running after stop")
	}

	// Server should not respond
	_, err = http.Get("http://localhost:18084/test")
	if err == nil {
		t.Fatal("expected connection error after server stop")
	}
}

func TestAlreadyStarted(t *testing.T) {
	t.Parallel()

	server := &echo.Server{}
	defer server.Stop(5 * time.Second)

	if err := server.Start(":18085", nil); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	// Try to start again
	err := server.Start(":18086", nil)
	if err == nil {
		t.Fatal("expected error when starting already running server")
	}
	if err.Error() != "server already started" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetEcho(t *testing.T) {
	t.Parallel()

	server := &echo.Server{}
	defer server.Stop(5 * time.Second)

	// Get Echo instance
	e := server.GetEcho()
	if e == nil {
		t.Fatal("Echo instance should not be nil")
	}

	// Configure Echo directly
	e.HideBanner = true
	e.Debug = false

	// Should still work with our wrapper methods
	server.RegisterHandler(echo_lib.GET, "/direct", func(c echo_lib.Context) error {
		return c.String(http.StatusOK, "direct config works")
	})

	if err := server.Start(":18087", nil); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:18087/direct")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRouteMiddleware(t *testing.T) {
	t.Parallel()

	server := &echo.Server{}
	defer server.Stop(5 * time.Second)

	// Middleware that adds a header
	addHeader := func(next echo_lib.HandlerFunc) echo_lib.HandlerFunc {
		return func(c echo_lib.Context) error {
			c.Response().Header().Set("X-Custom", "test")
			return next(c)
		}
	}

	server.RegisterHandler(echo_lib.GET, "/with-middleware", func(c echo_lib.Context) error {
		return c.String(http.StatusOK, "ok")
	}, addHeader)

	if err := server.Start(":18088", nil); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:18088/with-middleware")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.Header.Get("X-Custom") != "test" {
		t.Fatal("middleware header not set")
	}
}

func TestWrapHandler(t *testing.T) {
	t.Parallel()

	server := &echo.Server{}
	defer server.Stop(5 * time.Second)

	// Standard http.HandlerFunc
	stdHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Handler-Type", "standard")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from standard handler"))
	})

	server.RegisterHandler(echo_lib.GET, "/standard", echo.WrapHandler(stdHandler))

	if err := server.Start(":18092", nil); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:18092/standard")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "Hello from standard handler" {
		t.Fatalf("expected 'Hello from standard handler', got '%s'", body)
	}

	if resp.Header.Get("X-Handler-Type") != "standard" {
		t.Fatal("custom header not set")
	}
}

func TestWrapHandlerWithPathParams(t *testing.T) {
	t.Parallel()

	server := &echo.Server{}
	defer server.Stop(5 * time.Second)

	// Standard handler that reads path from request
	stdHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Path: " + r.URL.Path))
	})

	server.RegisterHandler(echo_lib.GET, "/users/:id/posts/:postId", echo.WrapHandler(stdHandler))

	if err := server.Start(":18093", nil); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:18093/users/123/posts/456")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "Path: /users/123/posts/456" {
		t.Fatalf("expected 'Path: /users/123/posts/456', got '%s'", body)
	}
}

func TestWrapHandlerFileServer(t *testing.T) {
	t.Parallel()

	server := &echo.Server{}
	defer server.Stop(5 * time.Second)

	// Create a simple in-memory file server simulation
	stdHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate file serving
		if strings.HasSuffix(r.URL.Path, ".txt") {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("This is a text file"))
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("File not found"))
		}
	})

	server.RegisterHandlerAny("/static/*", echo.WrapHandler(stdHandler))

	if err := server.Start(":18094", nil); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	// Test existing file
	resp, err := http.Get("http://localhost:18094/static/test.txt")
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if string(body) != "This is a text file" {
		t.Fatalf("unexpected body: %s", body)
	}
	if resp.Header.Get("Content-Type") != "text/plain" {
		t.Fatalf("expected Content-Type 'text/plain', got '%s'", resp.Header.Get("Content-Type"))
	}

	// Test non-existing file
	resp, err = http.Get("http://localhost:18094/static/missing.jpg")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", resp.StatusCode)
	}
}

func TestWrapHandlerFunc(t *testing.T) {
	t.Parallel()

	server := &echo.Server{}
	defer server.Stop(5 * time.Second)

	// Standard http.HandlerFunc
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Custom-Header", "test-value")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Response from HandlerFunc"))
	}

	server.RegisterHandler(echo_lib.GET, "/handlerfunc", echo.WrapHandlerFunc(handler))

	if err := server.Start(":18095", nil); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:18095/handlerfunc")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "Response from HandlerFunc" {
		t.Fatalf("expected 'Response from HandlerFunc', got '%s'", body)
	}

	if resp.Header.Get("X-Custom-Header") != "test-value" {
		t.Fatal("custom header not set")
	}
}

func TestWrapHandlerFuncVsWrapHandler(t *testing.T) {
	t.Parallel()

	server := &echo.Server{}
	defer server.Stop(5 * time.Second)

	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
	}

	// Both should work identically for http.HandlerFunc
	server.RegisterHandler(echo_lib.GET, "/wrap-handler", echo.WrapHandler(http.HandlerFunc(handlerFunc)))
	server.RegisterHandler(echo_lib.GET, "/wrap-handlerfunc", echo.WrapHandlerFunc(handlerFunc))

	if err := server.Start(":18096", nil); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	// Test WrapHandler
	resp1, err := http.Get("http://localhost:18096/wrap-handler")
	if err != nil {
		t.Fatal(err)
	}
	body1, _ := io.ReadAll(resp1.Body)
	resp1.Body.Close()

	// Test WrapHandlerFunc
	resp2, err := http.Get("http://localhost:18096/wrap-handlerfunc")
	if err != nil {
		t.Fatal(err)
	}
	body2, _ := io.ReadAll(resp2.Body)
	resp2.Body.Close()

	// Both should return the same result
	if string(body1) != string(body2) {
		t.Fatalf("WrapHandler and WrapHandlerFunc should produce same result")
	}
	if string(body1) != "test" {
		t.Fatalf("unexpected response: %s", body1)
	}
}

func TestWrapHandlerFuncWithQueryParams(t *testing.T) {
	t.Parallel()

	server := &echo.Server{}
	defer server.Stop(5 * time.Second)

	handler := func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		age := r.URL.Query().Get("age")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("name=" + name + "&age=" + age))
	}

	server.RegisterHandler(echo_lib.GET, "/query", echo.WrapHandlerFunc(handler))

	if err := server.Start(":18097", nil); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:18097/query?name=john&age=30")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "name=john&age=30" {
		t.Fatalf("unexpected response: %s", body)
	}
}

func TestRegisterHandlerAny(t *testing.T) {
	t.Parallel()

	server := &echo.Server{}
	defer server.Stop(5 * time.Second)

	// Register handler that accepts any HTTP method
	server.RegisterHandlerAny("/any", func(c echo_lib.Context) error {
		return c.String(http.StatusOK, "method: "+c.Request().Method)
	})

	if err := server.Start(":18089", nil); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	baseURL := "http://localhost:18089/any"

	// Test GET
	resp, err := http.Get(baseURL)
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET: expected 200, got %d", resp.StatusCode)
	}
	if string(body) != "method: GET" {
		t.Fatalf("GET: expected 'method: GET', got '%s'", body)
	}

	// Test POST
	resp, err = http.Post(baseURL, "text/plain", strings.NewReader("data"))
	if err != nil {
		t.Fatal(err)
	}
	body, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("POST: expected 200, got %d", resp.StatusCode)
	}
	if string(body) != "method: POST" {
		t.Fatalf("POST: expected 'method: POST', got '%s'", body)
	}

	// Test PUT
	req, _ := http.NewRequest(http.MethodPut, baseURL, strings.NewReader("data"))
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	body, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("PUT: expected 200, got %d", resp.StatusCode)
	}
	if string(body) != "method: PUT" {
		t.Fatalf("PUT: expected 'method: PUT', got '%s'", body)
	}

	// Test DELETE
	req, _ = http.NewRequest(http.MethodDelete, baseURL, nil)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	body, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("DELETE: expected 200, got %d", resp.StatusCode)
	}
	if string(body) != "method: DELETE" {
		t.Fatalf("DELETE: expected 'method: DELETE', got '%s'", body)
	}

	// Test PATCH
	req, _ = http.NewRequest(http.MethodPatch, baseURL, strings.NewReader("data"))
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	body, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("PATCH: expected 200, got %d", resp.StatusCode)
	}
	if string(body) != "method: PATCH" {
		t.Fatalf("PATCH: expected 'method: PATCH', got '%s'", body)
	}
}

func TestRegisterHandlerAnyWithMiddleware(t *testing.T) {
	t.Parallel()

	server := &echo.Server{}
	defer server.Stop(5 * time.Second)

	// Middleware that adds a custom header
	addHeader := func(next echo_lib.HandlerFunc) echo_lib.HandlerFunc {
		return func(c echo_lib.Context) error {
			c.Response().Header().Set("X-Method", c.Request().Method)
			return next(c)
		}
	}

	// Register handler with middleware
	server.RegisterHandlerAny("/any-with-middleware", func(c echo_lib.Context) error {
		return c.String(http.StatusOK, "ok")
	}, addHeader)

	if err := server.Start(":18090", nil); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	baseURL := "http://localhost:18090/any-with-middleware"

	// Test GET with middleware
	resp, err := http.Get(baseURL)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if resp.Header.Get("X-Method") != "GET" {
		t.Fatalf("expected header 'X-Method: GET', got '%s'", resp.Header.Get("X-Method"))
	}

	// Test POST with middleware
	resp, err = http.Post(baseURL, "text/plain", strings.NewReader("data"))
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if resp.Header.Get("X-Method") != "POST" {
		t.Fatalf("expected header 'X-Method: POST', got '%s'", resp.Header.Get("X-Method"))
	}
}

func TestMultipleStops(t *testing.T) {
	t.Parallel()

	server := &echo.Server{}

	server.RegisterHandler(echo_lib.GET, "/test", func(c echo_lib.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	if err := server.Start(":18091", nil); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	// First stop
	if err := server.Stop(5 * time.Second); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	// Second stop should not error
	if err := server.Stop(5 * time.Second); err != nil {
		t.Fatalf("multiple stops should not error: %v", err)
	}
}

func TestStopBeforeStart(t *testing.T) {
	t.Parallel()

	server := &echo.Server{}

	// Stop before starting should not error
	if err := server.Stop(5 * time.Second); err != nil {
		t.Fatalf("stop before start should not error: %v", err)
	}
}
