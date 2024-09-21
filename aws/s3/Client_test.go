package s3_test

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	aws_s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/common-library/go/aws/s3"
)

func getClient(t *testing.T) (s3.Client, bool) {
	t.Parallel()

	client := s3.Client{}

	if len(os.Getenv("S3_URL")) == 0 {
		return client, false
	}

	if err := client.CreateClient(context.TODO(), "dummy", "dummy", "dummy", "dummy",
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{URL: os.Getenv("S3_URL"), HostnameImmutable: true}, nil
			})),
	); err != nil {
		t.Fatal(err)
	}

	return client, true
}

func TestCreateClient(t *testing.T) {
	_, _ = getClient(t)
}

func TestCreateBucket(t *testing.T) {
	client, ok := getClient(t)
	if ok == false {
		return
	}

	bucketName := strings.ToLower(t.Name())
	if _, err := client.CreateBucket(bucketName, "dummy"); err != nil {
		t.Fatal(err)
	} else if _, err := client.DeleteBucket(bucketName); err != nil {
		t.Fatal(err)
	}
}

func TestListBuckets(t *testing.T) {
	bucketName := strings.ToLower(t.Name())

	client, ok := getClient(t)
	if ok == false {
		return
	}

	exist := func(output *aws_s3.ListBucketsOutput, bucketName string) bool {
		for _, bucket := range output.Buckets {
			if *bucket.Name == bucketName {
				return true
			}
		}

		return false
	}

	if output, err := client.ListBuckets(); err != nil {
		t.Fatal(err)
	} else if exist(output, bucketName) {
		t.Fatal(output)
	}

	if _, err := client.CreateBucket(bucketName, "dummy"); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if _, err := client.DeleteBucket(bucketName); err != nil {
			t.Fatal(err)
		}
	}()

	if output, err := client.ListBuckets(); err != nil {
		t.Fatal(err)
	} else if exist(output, bucketName) == false {
		t.Fatal(output)
	}
}

func TestDeleteBucket(t *testing.T) {
	TestCreateBucket(t)
}

func TestPutObject(t *testing.T) {
	const key = "key"
	const data = "data"
	bucketName := strings.ToLower(t.Name())

	client, ok := getClient(t)
	if ok == false {
		return
	}

	if _, err := client.CreateBucket(bucketName, "dummy"); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if _, err := client.DeleteBucket(bucketName); err != nil {
			t.Fatal(err)
		}
	}()

	if _, err := client.PutObject(bucketName, key, data); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if _, err := client.DeleteObject(bucketName, key); err != nil {
			t.Fatal(err)
		}
	}()

	if output, err := client.GetObject(bucketName, key); err != nil {
		t.Fatal(err)
	} else {
		defer output.Body.Close()

		if body, err := io.ReadAll(output.Body); err != nil {
			t.Fatal(err)
		} else if string(body) != data {
			t.Fatal(string(body), ",", data)
		}
	}

}

func TestGetObject(t *testing.T) {
	TestPutObject(t)
}

func TestDeleteObject(t *testing.T) {
	TestPutObject(t)
}
