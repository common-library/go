package gin_test

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/common-library/go/http/gin"
	gin_lib "github.com/gin-gonic/gin"
)

func init() {
	gin_lib.SetMode(gin_lib.TestMode)
}

func TestRegisterHandler(t *testing.T) {
	t.Parallel()

	server := &gin.Server{}
	defer server.Stop(5 * time.Second)

	server.RegisterHandler(http.MethodGet, "/test", func(c *gin_lib.Context) {
		c.String(http.StatusOK, "test response")
	})

	if err := server.Start(":19080", nil); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:19080/test")
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

	server := &gin.Server{}
	defer server.Stop(5 * time.Second)

	// Test GET
	server.RegisterHandler(http.MethodGet, "/get", func(c *gin_lib.Context) {
		c.String(http.StatusOK, "GET response")
	})

	// Test POST
	server.RegisterHandler(http.MethodPost, "/post", func(c *gin_lib.Context) {
		c.String(http.StatusCreated, "POST response")
	})

	// Test PUT
	server.RegisterHandler(http.MethodPut, "/put", func(c *gin_lib.Context) {
		c.String(http.StatusOK, "PUT response")
	})

	// Test DELETE
	server.RegisterHandler(http.MethodDelete, "/delete", func(c *gin_lib.Context) {
		c.Status(http.StatusNoContent)
	})

	if err := server.Start(":19081", nil); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	// Test GET request
	resp, err := http.Get("http://localhost:19081/get")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET: expected 200, got %d", resp.StatusCode)
	}

	// Test POST request
	resp, err = http.Post("http://localhost:19081/post", "text/plain", strings.NewReader("data"))
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("POST: expected 201, got %d", resp.StatusCode)
	}
}

func TestRegisterHandlerAny(t *testing.T) {
	t.Parallel()

	server := &gin.Server{}
	defer server.Stop(5 * time.Second)

	server.RegisterHandlerAny("/any", func(c *gin_lib.Context) {
		c.String(http.StatusOK, c.Request.Method)
	})

	if err := server.Start(":19082", nil); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	// Test GET
	resp, err := http.Get("http://localhost:19082/any")
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if string(body) != "GET" {
		t.Fatalf("expected 'GET', got '%s'", body)
	}

	// Test POST
	resp, err = http.Post("http://localhost:19082/any", "text/plain", nil)
	if err != nil {
		t.Fatal(err)
	}
	body, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	if string(body) != "POST" {
		t.Fatalf("expected 'POST', got '%s'", body)
	}
}

func TestUseMiddleware(t *testing.T) {
	t.Parallel()

	server := &gin.Server{}
	defer server.Stop(5 * time.Second)

	// Add global middleware
	server.Use(gin_lib.Recovery())

	server.RegisterHandler(http.MethodGet, "/panic", func(c *gin_lib.Context) {
		panic("test panic")
	})

	if err := server.Start(":19083", nil); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	// Should recover from panic
	resp, err := http.Get("http://localhost:19083/panic")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Recovery middleware returns 500
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", resp.StatusCode)
	}
}

func TestGroup(t *testing.T) {
	t.Parallel()

	server := &gin.Server{}
	defer server.Stop(5 * time.Second)

	// Create API group
	api := server.Group("/api")
	api.GET("/users", func(c *gin_lib.Context) {
		c.JSON(http.StatusOK, gin_lib.H{"message": "users"})
	})

	v1 := api.Group("/v1")
	v1.GET("/info", func(c *gin_lib.Context) {
		c.String(http.StatusOK, "v1 info")
	})

	if err := server.Start(":19084", nil); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	// Test /api/users
	resp, err := http.Get("http://localhost:19084/api/users")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	// Test /api/v1/info
	resp, err = http.Get("http://localhost:19084/api/v1/info")
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

	server := &gin.Server{}

	server.RegisterHandler(http.MethodGet, "/test", func(c *gin_lib.Context) {
		c.String(http.StatusOK, "ok")
	})

	// Server should not be running initially
	if server.IsRunning() {
		t.Fatal("server should not be running initially")
	}

	// Start server
	if err := server.Start(":19085", nil); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	// Server should be running
	if !server.IsRunning() {
		t.Fatal("server should be running")
	}

	// Test server is responding
	resp, err := http.Get("http://localhost:19085/test")
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
	_, err = http.Get("http://localhost:19085/test")
	if err == nil {
		t.Fatal("expected connection error after server stop")
	}
}

func TestAlreadyStarted(t *testing.T) {
	t.Parallel()

	server := &gin.Server{}
	defer server.Stop(5 * time.Second)

	if err := server.Start(":19086", nil); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	// Try to start again
	err := server.Start(":19087", nil)
	if err == nil {
		t.Fatal("expected error when starting already running server")
	}
	if err.Error() != "server already started" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetEngine(t *testing.T) {
	t.Parallel()

	server := &gin.Server{}
	defer server.Stop(5 * time.Second)

	// Get Engine instance
	engine := server.GetEngine()
	if engine == nil {
		t.Fatal("Engine instance should not be nil")
	}

	// Configure Gin directly
	engine.SetTrustedProxies([]string{"127.0.0.1"})

	// Should still work with our wrapper methods
	server.RegisterHandler(http.MethodGet, "/direct", func(c *gin_lib.Context) {
		c.String(http.StatusOK, "direct config works")
	})

	if err := server.Start(":19088", nil); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:19088/direct")
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

	server := &gin.Server{}
	defer server.Stop(5 * time.Second)

	// Middleware that adds a header
	addHeader := func(c *gin_lib.Context) {
		c.Header("X-Custom", "test")
		c.Next()
	}

	server.RegisterHandler(http.MethodGet, "/with-middleware", addHeader, func(c *gin_lib.Context) {
		c.String(http.StatusOK, "ok")
	})

	if err := server.Start(":19089", nil); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:19089/with-middleware")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.Header.Get("X-Custom") != "test" {
		t.Fatal("middleware header not set")
	}
}
