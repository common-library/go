// Package s3 provides Amazon S3 client implementations.
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

// CreateClient creates S3 client with service-specific options.
// This is the recommended way to create S3 client with custom endpoint.
//
// Example:
//
//	err := client.CreateClient(context.TODO(), "us-east-1", "access-key", "secret-key", "",
//	    func(o *s3.Options) {
//	        o.BaseEndpoint = aws.String("http://localhost:9090")
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

// CreateBucket creates a bucket.
//
// Example: _, err := client.CreateBucket(bucketName, "us-west-2")
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

// ListBuckets retrieves all buckets.
//
// Example: output, err := client.ListBuckets()
func (c *Client) ListBuckets() (*aws_s3.ListBucketsOutput, error) {
	return c.client.ListBuckets(c.ctx, &aws_s3.ListBucketsInput{})
}

// DeleteBucket deletes a bucket.
//
// Example: _, err := client.DeleteBucket(bucketName)
func (c *Client) DeleteBucket(name string) (*aws_s3.DeleteBucketOutput, error) {
	return c.client.DeleteBucket(c.ctx, &aws_s3.DeleteBucketInput{Bucket: aws.String(name)})
}

// PutObject uploads an object to the specified bucket.
//
// Example: _, err := client.PutObject(bucketName, key, data)
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
// Example: output, err := client.GetObject(bucketName, key)
func (c *Client) GetObject(bucketName string, key string) (*aws_s3.GetObjectOutput, error) {
	return c.client.GetObject(
		c.ctx,
		&aws_s3.GetObjectInput{Bucket: aws.String(bucketName), Key: aws.String(key)})
}

// DeleteObject deletes an object from the specified bucket.
//
// Example: _, err := client.DeleteObject(bucketName, key)
func (c *Client) DeleteObject(bucketName string, key string) (*aws_s3.DeleteObjectOutput, error) {
	return c.client.DeleteObject(
		c.ctx,
		&aws_s3.DeleteObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key)})
}
