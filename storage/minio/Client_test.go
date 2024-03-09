package minio_test

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/heaven-chp/common-library-go/file"
	"github.com/heaven-chp/common-library-go/storage/minio"
)

func getClient(t *testing.T) minio.Client {
	client := minio.Client{}

	if err := client.CreateClient("127.0.0.1:9090", "dummy", "dummy", false); err != nil {
		t.Fatal(err)
	}

	return client
}

func removeBucket(t *testing.T, client minio.Client, bucketName string) {
	if exist, err := client.BucketExists(bucketName); err != nil {
		t.Fatal(err)
	} else if exist == false {
		return
	} else if objectInfos, err := client.ListObjects(bucketName, "", true); err != nil {
		t.Fatal(err)
	} else if len(objectInfos) != 0 {
		for _, objectInfo := range objectInfos {
			if err := client.RemoveObject(bucketName, objectInfo.Key, true, true, ""); err != nil {
				t.Fatal(err)
			}
		}
	}

	if err := client.RemoveBucket(bucketName); err != nil {
		t.Fatal(err)
	}
}

func TestCreateClient(t *testing.T) {
	_ = getClient(t)
}

func TestMakeBucket(t *testing.T) {
	if err := (&minio.Client{}).MakeBucket("", "", true); err.Error() != "please call CreateClient first" {
		t.Fatal(err)
	}

	client := getClient(t)
	bucketName := uuid.New().String()
	defer removeBucket(t, client, bucketName)

	if err := client.MakeBucket(bucketName, "", true); err != nil {
		t.Fatal(err)
	} else if exist, err := client.BucketExists(bucketName); err != nil {
		t.Fatal(err)
	} else if exist == false {
		t.Fatal("exist == false")
	}
}

func TestListBuckets(t *testing.T) {
	if _, err := (&minio.Client{}).ListBuckets(); err.Error() != "please call CreateClient first" {
		t.Fatal(err)
	}

	client := getClient(t)
	bucketName := uuid.New().String()
	defer removeBucket(t, client, bucketName)

	if buckets, err := client.ListBuckets(); err != nil {
		t.Fatal(err)
	} else if len(buckets) != 0 {
		t.Fatal("len(buckets) != 0")
	} else if err := client.MakeBucket(bucketName, "", true); err != nil {
		t.Fatal(err)
	} else if buckets, err := client.ListBuckets(); err != nil {
		t.Fatal(err)
	} else if buckets[0].Name != bucketName {
		t.Fatal(bucketName, buckets[0].Name)
	}
}

func TestBucketExists(t *testing.T) {
	if _, err := (&minio.Client{}).BucketExists(""); err.Error() != "please call CreateClient first" {
		t.Fatal(err)
	}

	client := getClient(t)
	bucketName := uuid.New().String()
	defer removeBucket(t, client, bucketName)

	if exist, err := client.BucketExists(bucketName); err != nil {
		t.Fatal(err)
	} else if exist {
		t.Fatal(exist)
	} else if err := client.MakeBucket(bucketName, "", true); err != nil {
		t.Fatal(err)
	} else if exist, err := client.BucketExists(bucketName); err != nil {
		t.Fatal(err)
	} else if exist == false {
		t.Fatal("exist == false")
	}
}

func TestRemoveBucket(t *testing.T) {
	if err := (&minio.Client{}).RemoveBucket(""); err.Error() != "please call CreateClient first" {
		t.Fatal(err)
	}

	client := getClient(t)
	bucketName := uuid.New().String()

	if err := client.MakeBucket(bucketName, "", true); err != nil {
		t.Fatal(err)
	} else if err := client.RemoveBucket(bucketName); err != nil {
		t.Fatal(err)
	}
}

func TestListObjects(t *testing.T) {
	if _, err := (&minio.Client{}).ListObjects("", "", true); err.Error() != "please call CreateClient first" {
		t.Fatal(err)
	}

	client := getClient(t)
	bucketName := uuid.New().String()
	objectName := "test"
	filePath := "./test.txt"
	contentType := "text/plain"
	defer removeBucket(t, client, bucketName)

	if err := client.MakeBucket(bucketName, "", true); err != nil {
		t.Fatal(err)
	} else if err := client.FPutObject(bucketName, objectName, filePath, contentType); err != nil {
		t.Fatal(err)
	} else if objectInfos, err := client.ListObjects(bucketName, "", true); err != nil {
		t.Fatal(err)
	} else if objectInfos[0].Err != nil {
		t.Fatal(objectInfos[0].Err)
	} else if objectInfos[0].Key != objectName {
		t.Fatal("objectInfos[0].Key != objectName")
	}
}

func TestGetObject(t *testing.T) {
	if _, err := (&minio.Client{}).GetObject("", ""); err.Error() != "please call CreateClient first" {
		t.Fatal(err)
	}

	client := getClient(t)
	bucketName := uuid.New().String()
	objectName := "test"
	filePath := "./test.txt"
	contentType := "text/plain"
	defer removeBucket(t, client, bucketName)

	if err := client.MakeBucket(bucketName, "", true); err != nil {
		t.Fatal(err)
	} else if err := client.FPutObject(bucketName, objectName, filePath, contentType); err != nil {
		t.Fatal(err)
	} else if object, err := client.GetObject(bucketName, objectName); err != nil {
		t.Fatal(err)
	} else if objectInfo, err := object.Stat(); err != nil {
		t.Fatal(err)
	} else if objectInfo.Key != objectName {
		t.Fatal("objectInfo.Key != objectName")
	}
}

func TestPutObject(t *testing.T) {
	if err := (&minio.Client{}).PutObject("", "", "", nil, -1); err.Error() != "please call CreateClient first" {
		t.Fatal(err)
	}

	client := getClient(t)
	bucketName := uuid.New().String()
	objectName := "test"
	filePath := "./test.txt"
	contentType := "text/plain"
	defer removeBucket(t, client, bucketName)

	file, err := os.Open(filePath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}

	if err := client.MakeBucket(bucketName, "", true); err != nil {
		t.Fatal(err)
	} else if err := client.PutObject(bucketName, objectName, contentType, file, fileStat.Size()); err != nil {
		t.Fatal(err)
	} else if objectInfo, err := client.StatObject(bucketName, objectName); err != nil {
		t.Fatal(err)
	} else if objectInfo.Key != objectName {
		t.Fatal("objectInfo.Key != objectName")
	}
}

func TestCopyObject(t *testing.T) {
	if err := (&minio.Client{}).CopyObject("", "", "", ""); err.Error() != "please call CreateClient first" {
		t.Fatal(err)
	}

	client := getClient(t)
	filePath := "./test.txt"
	contentType := "text/plain"
	sourceBucketName := uuid.New().String()
	sourceObjectName := uuid.New().String()
	destinationBucketName := uuid.New().String()
	destinationObjectName := uuid.New().String()
	defer removeBucket(t, client, sourceBucketName)
	defer removeBucket(t, client, destinationBucketName)

	if err := client.MakeBucket(sourceBucketName, "", true); err != nil {
		t.Fatal(err)
	} else if err := client.MakeBucket(destinationBucketName, "", true); err != nil {
		t.Fatal(err)
	} else if err := client.FPutObject(sourceBucketName, sourceObjectName, filePath, contentType); err != nil {
		t.Fatal(err)
	} else if err := client.CopyObject(sourceBucketName, sourceObjectName, destinationBucketName, destinationObjectName); err != nil {
		t.Fatal(err)
	} else if object, err := client.GetObject(destinationBucketName, destinationObjectName); err != nil {
		t.Fatal(err)
	} else if objectInfo, err := object.Stat(); err != nil {
		t.Fatal(err)
	} else if objectInfo.Key != destinationObjectName {
		t.Fatal("objectInfo.Key != destinationObjectName")
	}
}

func TestStatObject(t *testing.T) {
	if _, err := (&minio.Client{}).StatObject("", ""); err.Error() != "please call CreateClient first" {
		t.Fatal(err)
	}

	client := getClient(t)
	bucketName := uuid.New().String()
	objectName := "test"
	filePath := "./test.txt"
	contentType := "text/plain"
	defer removeBucket(t, client, bucketName)

	if err := client.MakeBucket(bucketName, "", true); err != nil {
		t.Fatal(err)
	} else if err := client.FPutObject(bucketName, objectName, filePath, contentType); err != nil {
		t.Fatal(err)
	} else if objectInfo, err := client.StatObject(bucketName, objectName); err != nil {
		t.Fatal(err)
	} else if objectInfo.Key != objectName {
		t.Fatal("objectInfo.Key != objectName")
	}
}

func TestRemoveObject(t *testing.T) {
	if err := (&minio.Client{}).RemoveObject("", "", true, true, ""); err.Error() != "please call CreateClient first" {
		t.Fatal(err)
	}

	client := getClient(t)
	bucketName := uuid.New().String()
	objectName := "test"
	filePath := "./test.txt"
	contentType := "text/plain"
	defer removeBucket(t, client, bucketName)

	if err := client.MakeBucket(bucketName, "", true); err != nil {
		t.Fatal(err)
	} else if err := client.FPutObject(bucketName, objectName, filePath, contentType); err != nil {
		t.Fatal(err)
	} else if err := client.RemoveObject(bucketName, objectName, false, true, ""); err != nil {
		t.Fatal(err)
	} else if objectInfos, err := client.ListObjects(bucketName, "", true); err != nil {
		t.Fatal(err)
	} else if len(objectInfos) != 0 {
		t.Fatal("len(objectInfos) != 0")
	}
}

func TestRemoveObjects(t *testing.T) {
	if err := (&minio.Client{}).RemoveObjects("", nil, true); err[0].Err.Error() != "please call CreateClient first" {
		t.Fatal(err[0].Err)
	}

	client := getClient(t)
	bucketName := uuid.New().String()
	objectName := "test"
	filePath := "./test.txt"
	contentType := "text/plain"
	defer removeBucket(t, client, bucketName)

	if err := client.MakeBucket(bucketName, "", true); err != nil {
		t.Fatal(err)
	} else if err := client.FPutObject(bucketName, objectName, filePath, contentType); err != nil {
		t.Fatal(err)
	} else if objectInfos, err := client.ListObjects(bucketName, "", true); err != nil {
		t.Fatal(err)
	} else if err := client.RemoveObjects(bucketName, objectInfos, true); len(err) != 0 {
		t.Fatal(err)
	} else if objectInfos, err := client.ListObjects(bucketName, "", true); err != nil {
		t.Fatal(err)
	} else if len(objectInfos) != 0 {
		t.Fatal("len(objectInfos) != 0")
	}
}

func TestFPutObject(t *testing.T) {
	if err := (&minio.Client{}).FPutObject("", "", "", ""); err.Error() != "please call CreateClient first" {
		t.Fatal(err)
	}

	client := getClient(t)
	bucketName := uuid.New().String()
	objectName := "test"
	filePath := "./test.txt"
	contentType := "text/plain"
	defer removeBucket(t, client, bucketName)

	if err := client.MakeBucket(bucketName, "", true); err != nil {
		t.Fatal(err)
	} else if err := client.FPutObject(bucketName, objectName, filePath, contentType); err != nil {
		t.Fatal(err)
	} else if objectInfo, err := client.StatObject(bucketName, objectName); err != nil {
		t.Fatal(err)
	} else if objectInfo.Key != objectName {
		t.Fatal("objectInfo.Key != objectName")
	}
}

func TestFGetObject(t *testing.T) {
	if err := (&minio.Client{}).FGetObject("", "", ""); err.Error() != "please call CreateClient first" {
		t.Fatal(err)
	}

	client := getClient(t)
	bucketName := uuid.New().String()
	objectName := "test"
	filePath := "./test.txt"
	contentType := "text/plain"
	getFilePath := "./" + uuid.New().String()

	defer removeBucket(t, client, bucketName)
	defer file.Remove(getFilePath)

	if err := client.MakeBucket(bucketName, "", true); err != nil {
		t.Fatal(err)
	} else if err := client.FPutObject(bucketName, objectName, filePath, contentType); err != nil {
		t.Fatal(err)
	} else if err := client.FGetObject(bucketName, objectName, getFilePath); err != nil {
		t.Fatal(err)
	} else if data, err := file.Read(getFilePath); err != nil {
		t.Fatal(err)
	} else if answer, err := file.Read(filePath); err != nil {
		t.Fatal(err)
	} else if data != answer {
		t.Fatal("invalid -", data)
	}
}
