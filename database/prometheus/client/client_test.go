package client_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/common-library/go/database/prometheus/client"
	"github.com/common-library/go/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	prometheusContainer testcontainers.Container
	prometheusEndpoint  string
	containerOnce       sync.Once
	cleanupOnce         sync.Once
)

func setupPrometheusContainer() error {
	var err error
	containerOnce.Do(func() {
		ctx := context.Background()
		req := testcontainers.ContainerRequest{
			Image:        testutil.PrometheusImage,
			ExposedPorts: []string{"9090/tcp"},
			WaitingFor:   wait.ForListeningPort("9090/tcp"),
			Cmd:          []string{"--config.file=/etc/prometheus/prometheus.yml", "--storage.tsdb.path=/prometheus", "--web.console.libraries=/etc/prometheus/console_libraries", "--web.console.templates=/etc/prometheus/consoles", "--web.enable-lifecycle"},
		}

		prometheusContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})

		if err == nil {
			prometheusEndpoint, err = prometheusContainer.Endpoint(ctx, "")
			if err == nil {
				err = waitForPrometheusReady(prometheusEndpoint, 10*time.Second)
			}
		}
	})
	return err
}

func teardownPrometheusContainer() {
	cleanupOnce.Do(func() {
		if prometheusContainer != nil {
			_ = prometheusContainer.Terminate(context.Background())
		}
	})
}

func TestMain(m *testing.M) {
	if err := setupPrometheusContainer(); err != nil {
		fmt.Printf("Failed to setup Prometheus container: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()

	teardownPrometheusContainer()

	os.Exit(code)
}

func waitForPrometheusReady(endpoint string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get("http://" + endpoint + "/-/ready")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return nil
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("prometheus not ready within %v", timeout)
}

func TestNewClient(t *testing.T) {
	c, err := client.NewClient("http://localhost:9090")
	assert.NoError(t, err)
	assert.NotNil(t, c)

	c, err = client.NewClient("invalid-address")
	assert.NoError(t, err)
	assert.NotNil(t, c)
}

func TestNewClientWithBasicAuth(t *testing.T) {
	c, err := client.NewClientWithBasicAuth("http://localhost:9090", "user", "pass")
	assert.NoError(t, err)
	assert.NotNil(t, c)

	c, err = client.NewClientWithBasicAuth("invalid-address", "user", "pass")
	assert.NoError(t, err)
	assert.NotNil(t, c)
}

func TestNewClientWithBearerToken(t *testing.T) {
	c, err := client.NewClientWithBearerToken("http://localhost:9090", "test-token")
	assert.NoError(t, err)
	assert.NotNil(t, c)

	c, err = client.NewClientWithBearerToken("invalid-address", "test-token")
	assert.NoError(t, err)
	assert.NotNil(t, c)
}

func TestClientQuery(t *testing.T) {
	c, err := client.NewClient("http://" + prometheusEndpoint)
	require.NoError(t, err)

	value, warnings, err := c.Query("up", time.Now(), 10*time.Second)

	assert.NoError(t, err)
	_ = warnings
	_ = value
}

func TestClientQueryRange(t *testing.T) {
	c, err := client.NewClient("http://" + prometheusEndpoint)
	require.NoError(t, err)

	now := time.Now()
	r := client.Range{
		Start: now.Add(-5 * time.Minute),
		End:   now,
		Step:  time.Minute,
	}

	value, warnings, err := c.QueryRange("up", r, 10*time.Second)

	assert.NoError(t, err)
	_ = warnings
	_ = value
}

func TestClientWithInvalidQueries(t *testing.T) {
	c, err := client.NewClient("http://localhost:9090")
	require.NoError(t, err)

	_, _, err = c.Query("invalid_query{", time.Now(), 10*time.Second)
	assert.Error(t, err)

	now := time.Now()
	r := client.Range{
		Start: now.Add(-5 * time.Minute),
		End:   now,
		Step:  time.Minute,
	}
	_, _, err = c.QueryRange("invalid_query{", r, 10*time.Second)
	assert.Error(t, err)
}
