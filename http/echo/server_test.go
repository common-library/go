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
