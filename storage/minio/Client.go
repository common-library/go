package minio

import (
	"context"
	"errors"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	client *minio.Client
}

func (this *Client) CreateClient(endpoint, accessKeyID, secretAccessKey string, secure bool) error {
	if client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: secure,
	}); err != nil {
		return err
	} else {
		this.client = client
		return nil
	}
}

func (this *Client) MakeBucket(bucketName, region string, objectLocking bool) error {
	if this.client == nil {
		return errors.New("please call CreateClient first")
	}

	options := minio.MakeBucketOptions{Region: region, ObjectLocking: objectLocking}

	return this.client.MakeBucket(context.Background(), bucketName, options)
}

func (this *Client) ListBuckets() ([]minio.BucketInfo, error) {
	if this.client == nil {
		return nil, errors.New("please call CreateClient first")
	}

	return this.client.ListBuckets(context.Background())
}

func (this *Client) BucketExists(bucketName string) (bool, error) {
	if this.client == nil {
		return false, errors.New("please call CreateClient first")
	}

	return this.client.BucketExists(context.Background(), bucketName)
}

func (this *Client) RemoveBucket(bucketName string) error {
	if this.client == nil {
		return errors.New("please call CreateClient first")
	}

	return this.client.RemoveBucket(context.Background(), bucketName)
}

func (this *Client) ListObjects(bucketName, prefix string, recursive bool) ([]minio.ObjectInfo, error) {
	if this.client == nil {
		return nil, errors.New("please call CreateClient first")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	options := minio.ListObjectsOptions{Prefix: prefix, Recursive: recursive}

	objectInfosChan := this.client.ListObjects(ctx, bucketName, options)

	objectInfos := []minio.ObjectInfo{}
	for objectInfo := range objectInfosChan {
		objectInfos = append(objectInfos, objectInfo)
	}

	return objectInfos, nil
}

func (this *Client) GetObject(bucketName, objectName string) (*minio.Object, error) {
	if this.client == nil {
		return nil, errors.New("please call CreateClient first")
	}

	return this.client.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
}

func (this *Client) PutObject(bucketName, objectName, contentType string, reader io.Reader, objectSize int64) error {
	if this.client == nil {
		return errors.New("please call CreateClient first")
	}

	_, err := this.client.PutObject(context.Background(), bucketName, objectName, reader, objectSize, minio.PutObjectOptions{ContentType: contentType})

	return err
}

func (this *Client) CopyObject(sourceBucketName, sourceObjectName, destinationBucketName, destinationObjectName string) error {
	if this.client == nil {
		return errors.New("please call CreateClient first")
	}

	sourceOptions := minio.CopySrcOptions{
		Bucket: sourceBucketName,
		Object: sourceObjectName,
	}

	destinationOptions := minio.CopyDestOptions{
		Bucket: destinationBucketName,
		Object: destinationObjectName,
	}

	_, err := this.client.CopyObject(context.Background(), destinationOptions, sourceOptions)

	return err
}

func (this *Client) StatObject(bucketName, objectName string) (minio.ObjectInfo, error) {
	if this.client == nil {
		return minio.ObjectInfo{}, errors.New("please call CreateClient first")
	}

	return this.client.StatObject(context.Background(), bucketName, objectName, minio.StatObjectOptions{})

}

func (this *Client) RemoveObject(bucketName, objectName string, forceDelete bool, governanceBypass bool, versionID string) error {
	if this.client == nil {
		return errors.New("please call CreateClient first")
	}

	options := minio.RemoveObjectOptions{
		ForceDelete:      forceDelete,
		GovernanceBypass: governanceBypass,
		VersionID:        versionID,
	}

	return this.client.RemoveObject(context.Background(), bucketName, objectName, options)
}

func (this *Client) RemoveObjects(bucketName string, objectInfos []minio.ObjectInfo, governanceBypass bool) []minio.RemoveObjectError {
	if this.client == nil {
		return []minio.RemoveObjectError{minio.RemoveObjectError{Err: errors.New("please call CreateClient first")}}
	}

	objectInfoChan := make(chan minio.ObjectInfo)

	go func() {
		defer close(objectInfoChan)

		for _, objectInfo := range objectInfos {
			objectInfoChan <- objectInfo
		}
	}()

	options := minio.RemoveObjectsOptions{GovernanceBypass: governanceBypass}

	removeObjectError := []minio.RemoveObjectError{}
	for err := range this.client.RemoveObjects(context.Background(), bucketName, objectInfoChan, options) {
		removeObjectError = append(removeObjectError, err)
	}

	return removeObjectError
}

func (this *Client) FPutObject(bucketName, objectName, filePath, contentType string) error {
	if this.client == nil {
		return errors.New("please call CreateClient first")
	}

	options := minio.PutObjectOptions{
		ContentType: contentType,
	}

	_, err := this.client.FPutObject(context.Background(), bucketName, objectName, filePath, options)

	return err
}

func (this *Client) FGetObject(bucketName, objectName, filePath string) error {
	if this.client == nil {
		return errors.New("please call CreateClient first")
	}

	return this.client.FGetObject(context.Background(), bucketName, objectName, filePath, minio.GetObjectOptions{})
}
