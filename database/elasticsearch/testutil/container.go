package testutil

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func RunWithElasticsearch(m *testing.M, image string, elasticsearchURL *string) {
	ctx := context.Background()

	// Configure Elasticsearch container request
	req := testcontainers.ContainerRequest{
		Image:        image,
		ExposedPorts: []string{"9200/tcp"},
		Env: map[string]string{
			"ES_JAVA_OPTS":   "-Xms512m -Xmx512m",
			"discovery.type": "single-node",
			"cluster.routing.allocation.disk.threshold_enabled": "false",
			"xpack.security.enabled":                            "false",
			"xpack.ml.enabled":                                  "false",
			"xpack.watcher.enabled":                             "false",
		},
		WaitingFor: wait.ForHTTP("/_cluster/health?wait_for_status=yellow&timeout=1s").
			WithPort(nat.Port("9200/tcp")).
			WithStatusCodeMatcher(func(status int) bool {
				return status == 200
			}).
			WithStartupTimeout(60 * time.Second).     // Reduced from 1 minute to 60 seconds
			WithPollInterval(500 * time.Millisecond), // Reduced from 2 seconds to 500ms
	}

	// Start container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		fmt.Printf("Failed to start Elasticsearch container: %v\n", err)
		os.Exit(1)
	}

	// Get container host and port information
	host, err := container.Host(ctx)
	if err != nil {
		fmt.Printf("Failed to get container host: %v\n", err)
		os.Exit(1)
	}

	natPort, err := container.MappedPort(ctx, nat.Port("9200"))
	if err != nil {
		fmt.Printf("Failed to get container port: %v\n", err)
		os.Exit(1)
	}

	*elasticsearchURL = fmt.Sprintf("http://%s:%s", host, natPort.Port())

	// Run tests
	code := m.Run()

	// Cleanup container after tests
	if err := container.Terminate(ctx); err != nil {
		fmt.Printf("Failed to terminate container: %v\n", err)
	}

	os.Exit(code)
}
