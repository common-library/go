package s3_test

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	aws_s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/common-library/go/aws/s3"
	"github.com/common-library/go/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/localstack"
)

var (
	sharedContainer testcontainers.Container
	sharedEndpoint  string
	setupOnce       sync.Once
	setupErr        error
)

func setupSharedContainer() (string, error) {
	setupOnce.Do(func() {
		ctx := context.Background()

		container, err := localstack.Run(ctx,
			testutil.LocalstackImage,
			testcontainers.WithEnv(map[string]string{
				"SERVICES": "s3",
			}),
		)
		if err != nil {
			setupErr = err
			return
		}

		sharedContainer = container

		mappedPort, err := container.MappedPort(ctx, "4566/tcp")
		if err != nil {
			setupErr = err
			return
		}

		host, err := container.Host(ctx)
		if err != nil {
			setupErr = err
			return
		}

		sharedEndpoint = fmt.Sprintf("http://%s:%s", host, mappedPort.Port())
	})

	return sharedEndpoint, setupErr
}

func createTestClient(t *testing.T) *s3.Client {
	endpoint, err := setupSharedContainer()
	require.NoError(t, err)

	ctx := context.Background()
	client := &s3.Client{}
	err = client.CreateClient(
		ctx,
		"us-east-1",
		"test",
		"test",
		"",
		func(o *aws_s3.Options) {
			o.BaseEndpoint = aws.String(endpoint)
			o.UsePathStyle = true
		},
	)
	require.NoError(t, err)

	return client
}

func TestMain(m *testing.M) {
	code := m.Run()

	if sharedContainer != nil {
		ctx := context.Background()
		_ = sharedContainer.Terminate(ctx)
	}

	os.Exit(code)
}

func TestS3Client_Integration(t *testing.T) {
	client := createTestClient(t)

	t.Run("CreateBucket", func(t *testing.T) {
		t.Parallel()
		testCreateBucket(t, client)
	})

	t.Run("ListBuckets", func(t *testing.T) {
		t.Parallel()
		testListBuckets(t, client)
	})

	t.Run("PutObject", func(t *testing.T) {
		t.Parallel()
		testPutObject(t, client)
	})

	t.Run("GetObject", func(t *testing.T) {
		t.Parallel()
		testGetObject(t, client)
	})

	t.Run("DeleteObject", func(t *testing.T) {
		t.Parallel()
		testDeleteObject(t, client)
	})

	t.Run("DeleteBucket", func(t *testing.T) {
		t.Parallel()
		testDeleteBucket(t, client)
	})
}

func testCreateBucket(t *testing.T, client *s3.Client) {
	bucketName := "test-bucket-" + fmt.Sprintf("%d", time.Now().UnixNano())

	output, err := client.CreateBucket(bucketName, "eu-west-1")

	if err != nil && strings.Contains(strings.ToLower(err.Error()), "locationconstraint") {
		output, err = client.CreateBucket(bucketName, "us-west-2")
	}

	assert.NoError(t, err)
	assert.NotNil(t, output)
	if output != nil {
		assert.NotNil(t, output.Location)
	}
}

func testListBuckets(t *testing.T, client *s3.Client) {
	bucketName := "test-list-bucket-" + fmt.Sprintf("%d", time.Now().UnixNano())
	_, err := client.CreateBucket(bucketName, "eu-west-1")
	require.NoError(t, err)

	defer func() {
		_, _ = client.DeleteBucket(bucketName)
	}()

	output, err := client.ListBuckets()
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.Buckets)

	found := false
	for _, bucket := range output.Buckets {
		if *bucket.Name == bucketName {
			found = true
			break
		}
	}
	assert.True(t, found, "Created bucket should be in the list")
}

func testPutObject(t *testing.T, client *s3.Client) {
	bucketName := "test-put-bucket-" + fmt.Sprintf("%d", time.Now().UnixNano())
	key := "test-key"
	body := "test-content"

	_, err := client.CreateBucket(bucketName, "eu-west-1")
	require.NoError(t, err)

	defer func() {
		_, _ = client.DeleteObject(bucketName, key)
		_, _ = client.DeleteBucket(bucketName)
	}()

	output, err := client.PutObject(bucketName, key, body)
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotNil(t, output.ETag)
}

func testGetObject(t *testing.T, client *s3.Client) {
	bucketName := "test-get-bucket-" + fmt.Sprintf("%d", time.Now().UnixNano())
	key := "test-key"
	body := "test-content"

	_, err := client.CreateBucket(bucketName, "eu-west-1")
	require.NoError(t, err)

	defer func() {
		_, _ = client.DeleteObject(bucketName, key)
		_, _ = client.DeleteBucket(bucketName)
	}()

	_, err = client.PutObject(bucketName, key, body)
	require.NoError(t, err)

	output, err := client.GetObject(bucketName, key)
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotNil(t, output.Body)

	defer output.Body.Close()
	content, err := io.ReadAll(output.Body)
	assert.NoError(t, err)
	assert.Equal(t, body, string(content))
}

func testDeleteObject(t *testing.T, client *s3.Client) {
	bucketName := "test-delete-object-bucket-" + fmt.Sprintf("%d", time.Now().UnixNano())
	key := "test-key"
	body := "test-content"

	_, err := client.CreateBucket(bucketName, "eu-west-1")
	require.NoError(t, err)

	defer func() {
		_, _ = client.DeleteBucket(bucketName)
	}()

	_, err = client.PutObject(bucketName, key, body)
	require.NoError(t, err)

	output, err := client.DeleteObject(bucketName, key)
	assert.NoError(t, err)
	assert.NotNil(t, output)

	_, err = client.GetObject(bucketName, key)
	assert.Error(t, err, "Getting deleted object should return an error")
}

func testDeleteBucket(t *testing.T, client *s3.Client) {
	bucketName := "test-delete-bucket-" + fmt.Sprintf("%d", time.Now().UnixNano())

	_, err := client.CreateBucket(bucketName, "eu-west-1")
	require.NoError(t, err)

	output, err := client.DeleteBucket(bucketName)
	assert.NoError(t, err)
	assert.NotNil(t, output)

	listOutput, err := client.ListBuckets()
	assert.NoError(t, err)

	found := false
	for _, bucket := range listOutput.Buckets {
		if *bucket.Name == bucketName {
			found = true
			break
		}
	}
	assert.False(t, found, "Deleted bucket should not be in the list")
}

func TestS3Client_ErrorCases(t *testing.T) {
	client := createTestClient(t)

	t.Run("GetObject_NonExistentBucket", func(t *testing.T) {
		t.Parallel()
		_, err := client.GetObject("non-existent-bucket-"+fmt.Sprintf("%d", time.Now().UnixNano()), "test-key")
		assert.Error(t, err)
		errStr := strings.ToLower(err.Error())
		assert.True(t,
			strings.Contains(errStr, "nosuchbucket") ||
				strings.Contains(errStr, "notfound") ||
				strings.Contains(errStr, "not found") ||
				strings.Contains(errStr, "bucket") && strings.Contains(errStr, "exist"),
			"Error should indicate bucket doesn't exist, got: %s", err.Error())
	})

	t.Run("GetObject_NonExistentKey", func(t *testing.T) {
		t.Parallel()
		bucketName := "test-error-bucket-" + fmt.Sprintf("%d", time.Now().UnixNano())

		_, err := client.CreateBucket(bucketName, "eu-west-1")
		require.NoError(t, err)

		defer func() {
			_, _ = client.DeleteBucket(bucketName)
		}()

		_, err = client.GetObject(bucketName, "non-existent-key")
		assert.Error(t, err)
		errStr := strings.ToLower(err.Error())
		assert.True(t,
			strings.Contains(errStr, "nosuchkey") ||
				strings.Contains(errStr, "notfound") ||
				strings.Contains(errStr, "not found") ||
				strings.Contains(errStr, "key") && strings.Contains(errStr, "exist"),
			"Error should indicate key doesn't exist, got: %s", err.Error())
	})

	t.Run("DeleteBucket_NonEmpty", func(t *testing.T) {
		t.Parallel()
		bucketName := "test-non-empty-bucket-" + fmt.Sprintf("%d", time.Now().UnixNano())

		_, err := client.CreateBucket(bucketName, "eu-west-1")
		require.NoError(t, err)

		_, err = client.PutObject(bucketName, "test-key", "test-content")
		require.NoError(t, err)

		defer func() {
			_, _ = client.DeleteObject(bucketName, "test-key")
			_, _ = client.DeleteBucket(bucketName)
		}()

		_, err = client.DeleteBucket(bucketName)
		assert.Error(t, err)
		errStr := strings.ToLower(err.Error())
		assert.True(t,
			strings.Contains(errStr, "bucketnotempty") ||
				strings.Contains(errStr, "bucket not empty") ||
				strings.Contains(errStr, "empty"),
			"Error should indicate bucket is not empty, got: %s", err.Error())
	})
}
