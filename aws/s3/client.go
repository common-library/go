// Package s3 provides utilities for working with AWS S3 (Simple Storage Service).
//
// This package wraps the AWS SDK v2 S3 client to provide simplified
// functions for bucket and object management operations.
//
// Features:
//   - Bucket operations (create, list, delete)
//   - Object operations (put, get, delete)
//   - Custom endpoint support (for S3-compatible services)
//   - Path-style and virtual-hosted-style access
//
// Example usage:
//
//	var client s3.Client
//	err := client.CreateClient(ctx, "us-east-1", "key", "secret", "")
//	_, err = client.PutObject("bucket", "key", "data")
package s3

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	aws_s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// Client is a struct that provides client related methods.
type Client struct {
	ctx context.Context

	client *aws_s3.Client
}

// CreateClient creates an S3 client with service-specific options.
//
// This is the recommended way to create an S3 client with custom endpoint
// and other service-specific configurations (e.g., for MinIO or LocalStack).
//
// Parameters:
//   - ctx: context for the client lifecycle
//   - region: AWS region (e.g., "us-east-1")
//   - accessKey: AWS access key ID
//   - secretAccessKey: AWS secret access key
//   - sessionToken: optional session token (empty string if not using temporary credentials)
//   - s3OptionsFuncs: optional service-specific configuration functions
//
// Example:
//
//	err := client.CreateClient(context.TODO(), "us-east-1", "access-key", "secret-key", "",
//	    func(o *s3.Options) {
//	        o.BaseEndpoint = aws.String("http://localhost:9000")
//	        o.UsePathStyle = true
//	    })
func (c *Client) CreateClient(ctx context.Context, region, accessKey, secretAccessKey, sessionToken string, s3OptionsFuncs ...func(*aws_s3.Options)) error {
	c.ctx = ctx

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretAccessKey, sessionToken)))
	if err != nil {
		return err
	}

	c.client = aws_s3.NewFromConfig(cfg, s3OptionsFuncs...)

	return nil
}

// CreateBucket creates a new S3 bucket.
//
// Parameters:
//   - name: bucket name (must be globally unique)
//   - region: AWS region for the bucket (empty string for default region)
//
// Returns the CreateBucketOutput and any error encountered.
//
// Example:
//
//	_, err := client.CreateBucket("my-bucket", "us-west-2")
func (c *Client) CreateBucket(name, region string) (*aws_s3.CreateBucketOutput, error) {
	input := &aws_s3.CreateBucketInput{
		Bucket: aws.String(name),
	}

	if len(region) != 0 {
		input.CreateBucketConfiguration = &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(region),
		}
	}

	return c.client.CreateBucket(c.ctx, input)
}

// ListBuckets retrieves all buckets owned by the authenticated sender.
//
// Returns the ListBucketsOutput containing all buckets and any error encountered.
//
// Example:
//
//	output, err := client.ListBuckets()
func (c *Client) ListBuckets() (*aws_s3.ListBucketsOutput, error) {
	return c.client.ListBuckets(c.ctx, &aws_s3.ListBucketsInput{})
}

// DeleteBucket deletes an S3 bucket.
//
// The bucket must be empty before it can be deleted.
//
// Parameters:
//   - name: bucket name to delete
//
// Returns the DeleteBucketOutput and any error encountered.
//
// Example:
//
//	_, err := client.DeleteBucket("my-bucket")
func (c *Client) DeleteBucket(name string) (*aws_s3.DeleteBucketOutput, error) {
	return c.client.DeleteBucket(c.ctx, &aws_s3.DeleteBucketInput{Bucket: aws.String(name)})
}

// PutObject uploads an object to the specified bucket.
//
// Parameters:
//   - bucketName: name of the bucket
//   - key: object key (path) in the bucket
//   - body: object content as string
//
// Returns the PutObjectOutput and any error encountered.
//
// Example:
//
//	_, err := client.PutObject("my-bucket", "docs/file.txt", "Hello, World!")
func (c *Client) PutObject(bucketName, key, body string) (*aws_s3.PutObjectOutput, error) {
	return c.client.PutObject(
		c.ctx,
		&aws_s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
			Body:   strings.NewReader(body)})
}

// GetObject retrieves an object from the specified bucket.
//
// Parameters:
//   - bucketName: name of the bucket
//   - key: object key (path) in the bucket
//
// Returns the GetObjectOutput containing the object data and any error encountered.
//
// Example:
//
//	output, err := client.GetObject("my-bucket", "docs/file.txt")
func (c *Client) GetObject(bucketName string, key string) (*aws_s3.GetObjectOutput, error) {
	return c.client.GetObject(
		c.ctx,
		&aws_s3.GetObjectInput{Bucket: aws.String(bucketName), Key: aws.String(key)})
}

// DeleteObject deletes an object from the specified bucket.
//
// Parameters:
//   - bucketName: name of the bucket
//   - key: object key (path) to delete
//
// Returns the DeleteObjectOutput and any error encountered.
//
// Example:
//
//	_, err := client.DeleteObject("my-bucket", "docs/file.txt")
func (c *Client) DeleteObject(bucketName string, key string) (*aws_s3.DeleteObjectOutput, error) {
	return c.client.DeleteObject(
		c.ctx,
		&aws_s3.DeleteObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key)})
}
