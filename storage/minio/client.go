// Package minio provides MinIO client implementations.
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

func (c *Client) CreateClient(endpoint, accessKeyID, secretAccessKey string, secure bool) error {
	if client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: secure,
	}); err != nil {
		return err
	} else {
		c.client = client
		return nil
	}
}

func (c *Client) MakeBucket(bucketName, region string, objectLocking bool) error {
	if c.client == nil {
		return errors.New("please call CreateClient first")
	}

	options := minio.MakeBucketOptions{Region: region, ObjectLocking: objectLocking}

	return c.client.MakeBucket(context.Background(), bucketName, options)
}

func (c *Client) ListBuckets() ([]minio.BucketInfo, error) {
	if c.client == nil {
		return nil, errors.New("please call CreateClient first")
	}

	return c.client.ListBuckets(context.Background())
}

func (c *Client) BucketExists(bucketName string) (bool, error) {
	if c.client == nil {
		return false, errors.New("please call CreateClient first")
	}

	return c.client.BucketExists(context.Background(), bucketName)
}

func (c *Client) RemoveBucket(bucketName string) error {
	if c.client == nil {
		return errors.New("please call CreateClient first")
	}

	return c.client.RemoveBucket(context.Background(), bucketName)
}

func (c *Client) ListObjects(bucketName, prefix string, recursive bool) ([]minio.ObjectInfo, error) {
	if c.client == nil {
		return nil, errors.New("please call CreateClient first")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	options := minio.ListObjectsOptions{Prefix: prefix, Recursive: recursive}

	objectInfosChan := c.client.ListObjects(ctx, bucketName, options)

	objectInfos := []minio.ObjectInfo{}
	for objectInfo := range objectInfosChan {
		objectInfos = append(objectInfos, objectInfo)
	}

	return objectInfos, nil
}

func (c *Client) GetObject(bucketName, objectName string) (*minio.Object, error) {
	if c.client == nil {
		return nil, errors.New("please call CreateClient first")
	}

	return c.client.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
}

func (c *Client) PutObject(bucketName, objectName, contentType string, reader io.Reader, objectSize int64) error {
	if c.client == nil {
		return errors.New("please call CreateClient first")
	}

	_, err := c.client.PutObject(context.Background(), bucketName, objectName, reader, objectSize, minio.PutObjectOptions{ContentType: contentType})

	return err
}

func (c *Client) CopyObject(sourceBucketName, sourceObjectName, destinationBucketName, destinationObjectName string) error {
	if c.client == nil {
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

	_, err := c.client.CopyObject(context.Background(), destinationOptions, sourceOptions)

	return err
}

func (c *Client) StatObject(bucketName, objectName string) (minio.ObjectInfo, error) {
	if c.client == nil {
		return minio.ObjectInfo{}, errors.New("please call CreateClient first")
	}

	return c.client.StatObject(context.Background(), bucketName, objectName, minio.StatObjectOptions{})

}

func (c *Client) RemoveObject(bucketName, objectName string, forceDelete bool, governanceBypass bool, versionID string) error {
	if c.client == nil {
		return errors.New("please call CreateClient first")
	}

	options := minio.RemoveObjectOptions{
		ForceDelete:      forceDelete,
		GovernanceBypass: governanceBypass,
		VersionID:        versionID,
	}

	return c.client.RemoveObject(context.Background(), bucketName, objectName, options)
}

func (c *Client) RemoveObjects(bucketName string, objectInfos []minio.ObjectInfo, governanceBypass bool) []minio.RemoveObjectError {
	if c.client == nil {
		return []minio.RemoveObjectError{{Err: errors.New("please call CreateClient first")}}
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
	for err := range c.client.RemoveObjects(context.Background(), bucketName, objectInfoChan, options) {
		removeObjectError = append(removeObjectError, err)
	}

	return removeObjectError
}

func (c *Client) FPutObject(bucketName, objectName, filePath, contentType string) error {
	if c.client == nil {
		return errors.New("please call CreateClient first")
	}

	options := minio.PutObjectOptions{
		ContentType: contentType,
	}

	_, err := c.client.FPutObject(context.Background(), bucketName, objectName, filePath, options)

	return err
}

func (c *Client) FGetObject(bucketName, objectName, filePath string) error {
	if c.client == nil {
		return errors.New("please call CreateClient first")
	}

	return c.client.FGetObject(context.Background(), bucketName, objectName, filePath, minio.GetObjectOptions{})
}
