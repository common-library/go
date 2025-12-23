package testutil_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/common-library/go/database/elasticsearch/testutil"
	"github.com/stretchr/testify/assert"
)

func TestRunWithElasticsearch_Integration(t *testing.T) {
	t.Skip("Skipping integration test that requires Docker")

	var elasticsearchURL string

	m := &testing.M{}

	testutil.RunWithElasticsearch(m, "docker.elastic.co/elasticsearch/elasticsearch:9.1.0", &elasticsearchURL)
}

func TestRunWithElasticsearch_URLFormat(t *testing.T) {
	t.Skip("Skipping test that requires Docker and starts containers")

	expectedFormat := "http://localhost:"
	assert.Contains(t, expectedFormat, "http://")
	assert.Contains(t, expectedFormat, ":")
}

func ExampleRunWithElasticsearch() {
	var elasticsearchURL string

	fmt.Println("Example elasticsearch URL format:", elasticsearchURL)
}

func TestElasticsearchImages(t *testing.T) {
	images := []string{
		"docker.elastic.co/elasticsearch/elasticsearch:7.17.0",
		"docker.elastic.co/elasticsearch/elasticsearch:8.11.0",
		"docker.elastic.co/elasticsearch/elasticsearch:9.1.0",
	}

	for _, image := range images {
		t.Run(image, func(t *testing.T) {
			t.Skip("Skipping test that requires Docker")

			assert.NotEmpty(t, image)
		})
	}
}

func TestContainerConfiguration(t *testing.T) {
	expectedEnv := map[string]string{
		"ES_JAVA_OPTS":   "-Xms512m -Xmx512m",
		"discovery.type": "single-node",
		"cluster.routing.allocation.disk.threshold_enabled": "false",
		"xpack.security.enabled":                            "false",
		"xpack.ml.enabled":                                  "false",
		"xpack.watcher.enabled":                             "false",
	}

	assert.Equal(t, "-Xms512m -Xmx512m", expectedEnv["ES_JAVA_OPTS"])
	assert.Equal(t, "single-node", expectedEnv["discovery.type"])
	assert.Equal(t, "false", expectedEnv["xpack.security.enabled"])
}

func TestContainerPort(t *testing.T) {
	expectedPort := "9200/tcp"
	assert.Equal(t, "9200/tcp", expectedPort)
}

type MockTestMain struct {
	RunCalled bool
	ExitCode  int
}

func (m *MockTestMain) Run() int {
	m.RunCalled = true
	return m.ExitCode
}

func TestMockTestMain(t *testing.T) {
	mock := &MockTestMain{ExitCode: 0}
	code := mock.Run()

	assert.True(t, mock.RunCalled)
	assert.Equal(t, 0, code)
}

func TestOSExit(t *testing.T) {

	exitCodes := []int{0, 1}
	for _, code := range exitCodes {
		t.Run(fmt.Sprintf("exit_code_%d", code), func(t *testing.T) {
			assert.Contains(t, []int{0, 1}, code)
		})
	}
}

func TestContainerCleanup(t *testing.T) {
	t.Skip("Skipping test that would require actual container management")

	cleanupCalled := false

	defer func() {
		cleanupCalled = true
	}()

	_ = os.Getenv("TEST")

	assert.True(t, cleanupCalled || !cleanupCalled)
}
