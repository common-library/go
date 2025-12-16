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

	// Elasticsearch 컨테이너 요청 구성
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
			WithStartupTimeout(60 * time.Second).     // 1분 → 60초로 단축
			WithPollInterval(500 * time.Millisecond), // 2초 → 500ms로 단축
	}

	// 컨테이너 시작
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		fmt.Printf("Failed to start Elasticsearch container: %v\n", err)
		os.Exit(1)
	}

	// 컨테이너의 호스트와 포트 정보 가져오기
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

	// 테스트 실행
	code := m.Run()

	// 테스트 완료 후 컨테이너 정리
	if err := container.Terminate(ctx); err != nil {
		fmt.Printf("Failed to terminate container: %v\n", err)
	}

	os.Exit(code)
}
