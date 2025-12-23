// Package minio provides MinIO object storage client.
//
// This package wraps the MinIO Go SDK with simplified methods for common
// object storage operations including bucket management, object upload/download,
// and metadata operations.
//
// # Features
//
//   - Bucket creation and management
//   - Object upload and download
//   - Object listing and search
//   - File-based operations (FPutObject, FGetObject)
//   - Bulk object removal
//
// # Basic Example
//
//	client := &minio.Client{}
//	err := client.CreateClient("localhost:9000", "accessKey", "secretKey", false)
//	err = client.MakeBucket("mybucket", "us-east-1", false)
//	err = client.FPutObject("mybucket", "file.txt", "/path/to/file.txt", "text/plain")
package minio

import (
	"context"
	"errors"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Client provides MinIO object storage operations.
type Client struct {
	client *minio.Client
}

// CreateClient initializes the MinIO client connection.
//
// This method creates a new MinIO client with the provided credentials and
// connection settings. It must be called before any other operations.
//
// # Parameters
//
//   - endpoint: MinIO server address (e.g., "localhost:9000", "s3.amazonaws.com")
//   - accessKeyID: Access key for authentication
//   - secretAccessKey: Secret key for authentication
//   - secure: Use HTTPS if true, HTTP if false
//
// # Returns
//
//   - error: Error if client creation fails, nil on success
//
// # Examples
//
// Local MinIO server:
//
//	client := &minio.Client{}
//	err := client.CreateClient("localhost:9000", "minioadmin", "minioadmin", false)
//
// AWS S3:
//
//	err := client.CreateClient("s3.amazonaws.com", "ACCESS_KEY", "SECRET_KEY", true)
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

// MakeBucket creates a new storage bucket.
//
// # Parameters
//
//   - bucketName: Name of the bucket (must be globally unique, lowercase)
//   - region: AWS region (e.g., "us-east-1", "" for MinIO)
//   - objectLocking: Enable object versioning and locking
//
// # Returns
//
//   - error: Error if bucket creation fails, nil on success
//
// # Examples
//
//	err := client.MakeBucket("mybucket", "us-east-1", false)
func (c *Client) MakeBucket(bucketName, region string, objectLocking bool) error {
	if c.client == nil {
		return errors.New("please call CreateClient first")
	}

	options := minio.MakeBucketOptions{Region: region, ObjectLocking: objectLocking}

	return c.client.MakeBucket(context.Background(), bucketName, options)
}

// ListBuckets returns all buckets owned by the user.
//
// # Returns
//
//   - []minio.BucketInfo: Slice of bucket information
//   - error: Error if listing fails, nil on success
//
// # Examples
//
//	buckets, err := client.ListBuckets()
//	for _, bucket := range buckets {
//	    fmt.Printf("%s (created: %v)\n", bucket.Name, bucket.CreationDate)
//	}
func (c *Client) ListBuckets() ([]minio.BucketInfo, error) {
	if c.client == nil {
		return nil, errors.New("please call CreateClient first")
	}

	return c.client.ListBuckets(context.Background())
}

// BucketExists checks if a bucket exists and is accessible.
//
// # Parameters
//
//   - bucketName: Name of the bucket to check
//
// # Returns
//
//   - bool: true if bucket exists, false otherwise
//   - error: Error if check fails, nil on success
//
// # Examples
//
//	exists, err := client.BucketExists("mybucket")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	if exists {
//	    fmt.Println("Bucket exists")
//	}
func (c *Client) BucketExists(bucketName string) (bool, error) {
	if c.client == nil {
		return false, errors.New("please call CreateClient first")
	}

	return c.client.BucketExists(context.Background(), bucketName)
}

// RemoveBucket deletes a bucket.
//
// The bucket must be empty before it can be removed.
//
// # Parameters
//
//   - bucketName: Name of the bucket to remove
//
// # Returns
//
//   - error: Error if removal fails, nil on success
//
// # Examples
//
//	err := client.RemoveBucket("mybucket")
func (c *Client) RemoveBucket(bucketName string) error {
	if c.client == nil {
		return errors.New("please call CreateClient first")
	}

	return c.client.RemoveBucket(context.Background(), bucketName)
}

// ListObjects lists objects in a bucket.
//
// # Parameters
//
//   - bucketName: Name of the bucket
//   - prefix: Filter objects by prefix (e.g., "photos/2024/")
//   - recursive: List all objects recursively if true, only top-level if false
//
// # Returns
//
//   - []minio.ObjectInfo: Slice of object information
//   - error: Error if listing fails, nil on success
//
// # Examples
//
// List all objects:
//
//	objects, err := client.ListObjects("mybucket", "", true)
//
// List objects with prefix:
//
//	objects, err := client.ListObjects("photos", "2024/", true)
//	for _, obj := range objects {
//	    fmt.Printf("%s - %d bytes\n", obj.Key, obj.Size)
//	}
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

// GetObject retrieves an object from a bucket.
//
// Returns a reader for streaming the object content. The caller must close
// the returned object when done.
//
// # Parameters
//
//   - bucketName: Name of the bucket
//   - objectName: Name/key of the object
//
// # Returns
//
//   - *minio.Object: Object reader (must be closed by caller)
//   - error: Error if retrieval fails, nil on success
//
// # Examples
//
//	object, err := client.GetObject("mybucket", "document.pdf")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer object.Close()
//
//	data, err := io.ReadAll(object)
func (c *Client) GetObject(bucketName, objectName string) (*minio.Object, error) {
	if c.client == nil {
		return nil, errors.New("please call CreateClient first")
	}

	return c.client.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
}

// PutObject uploads an object to a bucket from a reader.
//
// # Parameters
//
//   - bucketName: Name of the bucket
//   - objectName: Name/key for the object
//   - contentType: MIME type (e.g., "image/jpeg", "application/pdf")
//   - reader: Data reader (e.g., file, buffer)
//   - objectSize: Size of the object in bytes (-1 for unknown size)
//
// # Returns
//
//   - error: Error if upload fails, nil on success
//
// # Examples
//
//	file, _ := os.Open("/path/to/file.pdf")
//	defer file.Close()
//
//	stat, _ := file.Stat()
//	err := client.PutObject("documents", "file.pdf", "application/pdf", file, stat.Size())
func (c *Client) PutObject(bucketName, objectName, contentType string, reader io.Reader, objectSize int64) error {
	if c.client == nil {
		return errors.New("please call CreateClient first")
	}

	_, err := c.client.PutObject(context.Background(), bucketName, objectName, reader, objectSize, minio.PutObjectOptions{ContentType: contentType})

	return err
}

// CopyObject copies an object from source to destination.
//
// # Parameters
//
//   - sourceBucketName: Source bucket name
//   - sourceObjectName: Source object name/key
//   - destinationBucketName: Destination bucket name
//   - destinationObjectName: Destination object name/key
//
// # Returns
//
//   - error: Error if copy fails, nil on success
//
// # Examples
//
// Copy within same bucket:
//
//	err := client.CopyObject("photos", "original.jpg", "photos", "backup.jpg")
//
// Copy to different bucket:
//
//	err := client.CopyObject("source-bucket", "file.txt", "backup-bucket", "file.txt")
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

// StatObject retrieves metadata and information about an object.
//
// # Parameters
//
//   - bucketName: Name of the bucket
//   - objectName: Name/key of the object
//
// # Returns
//
//   - minio.ObjectInfo: Object metadata (size, modified time, etag, etc.)
//   - error: Error if stat fails, nil on success
//
// # Examples
//
//	info, err := client.StatObject("mybucket", "document.pdf")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Size: %d bytes\n", info.Size)
//	fmt.Printf("Modified: %v\n", info.LastModified)
//	fmt.Printf("Content-Type: %s\n", info.ContentType)
func (c *Client) StatObject(bucketName, objectName string) (minio.ObjectInfo, error) {
	if c.client == nil {
		return minio.ObjectInfo{}, errors.New("please call CreateClient first")
	}

	return c.client.StatObject(context.Background(), bucketName, objectName, minio.StatObjectOptions{})
}

// RemoveObject deletes an object from a bucket.
//
// # Parameters
//
//   - bucketName: Name of the bucket
//   - objectName: Name/key of the object to delete
//   - forceDelete: Force delete even if object is locked
//   - governanceBypass: Bypass governance retention
//   - versionID: Specific version to delete (empty for latest)
//
// # Returns
//
//   - error: Error if deletion fails, nil on success
//
// # Examples
//
// Simple delete:
//
//	err := client.RemoveObject("mybucket", "old-file.txt", false, false, "")
//
// Delete specific version:
//
//	err := client.RemoveObject("mybucket", "file.txt", false, false, "version-id")
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

// RemoveObjects deletes multiple objects from a bucket in a single operation.
//
// This is more efficient than calling RemoveObject multiple times.
//
// # Parameters
//
//   - bucketName: Name of the bucket
//   - objectInfos: Slice of objects to delete
//   - governanceBypass: Bypass governance retention for all objects
//
// # Returns
//
//   - []minio.RemoveObjectError: Errors for failed deletions (empty if all succeeded)
//
// # Examples
//
//	objects, _ := client.ListObjects("mybucket", "temp/", true)
//	errors := client.RemoveObjects("mybucket", objects, false)
//	if len(errors) > 0 {
//	    for _, err := range errors {
//	        log.Printf("Failed to delete %s: %v", err.ObjectName, err.Err)
//	    }
//	}
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

// FPutObject uploads a file to a bucket.
//
// This is a convenience method for uploading files directly without
// manually opening and reading them.
//
// # Parameters
//
//   - bucketName: Name of the bucket
//   - objectName: Name/key for the uploaded object
//   - filePath: Local file path to upload
//   - contentType: MIME type (e.g., "image/jpeg", "application/pdf")
//
// # Returns
//
//   - error: Error if upload fails, nil on success
//
// # Examples
//
//	err := client.FPutObject("photos", "vacation.jpg", "/home/user/vacation.jpg", "image/jpeg")
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

// FGetObject downloads an object to a file.
//
// This is a convenience method for downloading objects directly to disk
// without manually handling readers and writers.
//
// # Parameters
//
//   - bucketName: Name of the bucket
//   - objectName: Name/key of the object to download
//   - filePath: Local file path to save the downloaded object
//
// # Returns
//
//   - error: Error if download fails, nil on success
//
// # Examples
//
//	err := client.FGetObject("photos", "vacation.jpg", "/tmp/vacation.jpg")
func (c *Client) FGetObject(bucketName, objectName, filePath string) error {
	if c.client == nil {
		return errors.New("please call CreateClient first")
	}

	return c.client.FGetObject(context.Background(), bucketName, objectName, filePath, minio.GetObjectOptions{})
}
