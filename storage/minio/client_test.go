package minio_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/common-library/go/storage/minio"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	minioContainer testcontainers.Container
	minioEndpoint  string
	accessKey      = "minioadmin"
	secretKey      = "minioadmin"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "minio/minio:RELEASE.2024-12-13T22-19-12Z",
		ExposedPorts: []string{"9000/tcp"},
		Env: map[string]string{
			"MINIO_ROOT_USER":     accessKey,
			"MINIO_ROOT_PASSWORD": secretKey,
		},
		Cmd:        []string{"server", "/data"},
		WaitingFor: wait.ForLog("MinIO Object Storage Server").WithPollInterval(1),
	}

	var err error
	minioContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Failed to start MinIO container: %v", err)
	}

	host, err := minioContainer.Host(ctx)
	if err != nil {
		log.Fatalf("Failed to get MinIO container host: %v", err)
	}

	port, err := minioContainer.MappedPort(ctx, "9000")
	if err != nil {
		log.Fatalf("Failed to get MinIO container port: %v", err)
	}

	minioEndpoint = fmt.Sprintf("%s:%s", host, port.Port())

	code := m.Run()

	if err := minioContainer.Terminate(ctx); err != nil {
		log.Printf("Failed to terminate MinIO container: %v", err)
	}

	os.Exit(code)
}

func setupClient(t *testing.T) *minio.Client {
	client := &minio.Client{}
	err := client.CreateClient(minioEndpoint, accessKey, secretKey, false)
	require.NoError(t, err, "Failed to create MinIO client")
	return client
}

func cleanupBucket(client *minio.Client, bucketName string) {
	objects, _ := client.ListObjects(bucketName, "", true)
	for _, obj := range objects {
		client.RemoveObject(bucketName, obj.Key, false, false, "")
	}
	client.RemoveBucket(bucketName)
}

func TestClient_CreateClient(t *testing.T) {
	client := &minio.Client{}

	err := client.CreateClient(minioEndpoint, accessKey, secretKey, false)
	assert.NoError(t, err)

	_ = client.CreateClient("invalid-endpoint", accessKey, secretKey, false)
}

func TestClient_MakeBucket(t *testing.T) {
	client := setupClient(t)
	bucketName := "test-make-bucket"

	err := client.MakeBucket(bucketName, "us-east-1", false)
	assert.NoError(t, err)

	exists, err := client.BucketExists(bucketName)
	assert.NoError(t, err)
	assert.True(t, exists)

	err = client.MakeBucket(bucketName, "us-east-1", false)
	assert.Error(t, err)

	defer cleanupBucket(client, bucketName)
}

func TestClient_ListBuckets(t *testing.T) {
	client := setupClient(t)
	bucketName := "test-list-buckets"

	initialBuckets, err := client.ListBuckets()
	assert.NoError(t, err)
	initialCount := len(initialBuckets)

	err = client.MakeBucket(bucketName, "us-east-1", false)
	assert.NoError(t, err)

	buckets, err := client.ListBuckets()
	assert.NoError(t, err)
	assert.Len(t, buckets, initialCount+1)

	found := false
	for _, bucket := range buckets {
		if bucket.Name == bucketName {
			found = true
			break
		}
	}
	assert.True(t, found, "Created bucket not found in the list")

	defer cleanupBucket(client, bucketName)
}

func TestClient_BucketExists(t *testing.T) {
	client := setupClient(t)
	bucketName := "test-bucket-exists"

	exists, err := client.BucketExists(bucketName)
	assert.NoError(t, err)
	assert.False(t, exists)

	err = client.MakeBucket(bucketName, "us-east-1", false)
	assert.NoError(t, err)

	exists, err = client.BucketExists(bucketName)
	assert.NoError(t, err)
	assert.True(t, exists)

	defer cleanupBucket(client, bucketName)
}

func TestClient_PutObjectAndGetObject(t *testing.T) {
	client := setupClient(t)
	bucketName := "test-put-get-object"
	objectName := "test-object.txt"
	content := "Hello, MinIO!"

	err := client.MakeBucket(bucketName, "us-east-1", false)
	require.NoError(t, err)

	err = client.PutObject(bucketName, objectName, "text/plain", strings.NewReader(content), int64(len(content)))
	assert.NoError(t, err)

	object, err := client.GetObject(bucketName, objectName)
	assert.NoError(t, err)
	defer object.Close()

	data, err := io.ReadAll(object)
	assert.NoError(t, err)
	assert.Equal(t, content, string(data))

	defer cleanupBucket(client, bucketName)
}

func TestClient_ListObjects(t *testing.T) {
	client := setupClient(t)
	bucketName := "test-list-objects"

	err := client.MakeBucket(bucketName, "us-east-1", false)
	require.NoError(t, err)

	objects := []string{"file1.txt", "file2.txt", "dir/file3.txt"}
	for _, objectName := range objects {
		err = client.PutObject(bucketName, objectName, "text/plain", strings.NewReader("test content"), 12)
		assert.NoError(t, err)
	}

	objectInfos, err := client.ListObjects(bucketName, "", true)
	assert.NoError(t, err)
	assert.Len(t, objectInfos, 3)

	foundObjects := make(map[string]bool)
	for _, info := range objectInfos {
		foundObjects[info.Key] = true
	}

	for _, expectedObject := range objects {
		assert.True(t, foundObjects[expectedObject], "Object %s not found", expectedObject)
	}

	defer cleanupBucket(client, bucketName)
}

func TestClient_StatObject(t *testing.T) {
	client := setupClient(t)
	bucketName := "test-stat-object"
	objectName := "test-stat.txt"
	content := "Test content for stat"

	err := client.MakeBucket(bucketName, "us-east-1", false)
	require.NoError(t, err)

	err = client.PutObject(bucketName, objectName, "text/plain", strings.NewReader(content), int64(len(content)))
	require.NoError(t, err)

	objectInfo, err := client.StatObject(bucketName, objectName)
	assert.NoError(t, err)
	assert.Equal(t, objectName, objectInfo.Key)
	assert.Equal(t, int64(len(content)), objectInfo.Size)
	assert.Equal(t, "text/plain", objectInfo.ContentType)

	defer cleanupBucket(client, bucketName)
}

func TestClient_RemoveObject(t *testing.T) {
	client := setupClient(t)
	bucketName := "test-remove-object"
	objectName := "test-remove.txt"
	content := "Test content to be removed"

	err := client.MakeBucket(bucketName, "us-east-1", false)
	require.NoError(t, err)

	err = client.PutObject(bucketName, objectName, "text/plain", strings.NewReader(content), int64(len(content)))
	require.NoError(t, err)

	_, err = client.StatObject(bucketName, objectName)
	assert.NoError(t, err)

	err = client.RemoveObject(bucketName, objectName, false, false, "")
	assert.NoError(t, err)

	_, err = client.StatObject(bucketName, objectName)
	assert.Error(t, err)

	defer cleanupBucket(client, bucketName)
}

func TestClient_CopyObject(t *testing.T) {
	client := setupClient(t)
	bucketName := "test-copy-object"
	sourceObjectName := "source.txt"
	destObjectName := "destination.txt"
	content := "Test content for copy"

	err := client.MakeBucket(bucketName, "us-east-1", false)
	require.NoError(t, err)

	err = client.PutObject(bucketName, sourceObjectName, "text/plain", strings.NewReader(content), int64(len(content)))
	require.NoError(t, err)

	err = client.CopyObject(bucketName, sourceObjectName, bucketName, destObjectName)
	assert.NoError(t, err)

	destObjectInfo, err := client.StatObject(bucketName, destObjectName)
	assert.NoError(t, err)
	assert.Equal(t, destObjectName, destObjectInfo.Key)
	assert.Equal(t, int64(len(content)), destObjectInfo.Size)

	defer cleanupBucket(client, bucketName)
}

func TestClient_WithoutCreateClient(t *testing.T) {
	client := &minio.Client{}

	err := client.MakeBucket("test", "us-east-1", false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "please call CreateClient first")

	_, err = client.ListBuckets()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "please call CreateClient first")

	_, err = client.BucketExists("test")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "please call CreateClient first")
}
