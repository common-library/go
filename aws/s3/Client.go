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

// CreateClient is create client.
// ex)
//
//	err := client.CreateClient(context.TODO(), "dummy", "dummy", "dummy", "dummy",
//	    config.WithEndpointResolver(aws.EndpointResolverFunc(
//	        func(service, region string) (aws.Endpoint, error) {
//	            return aws.Endpoint{URL: fmt.Sprintf("http://127.0.0.1:9090"), HostnameImmutable: true}, nil
//	        })),
//	)
func (this *Client) CreateClient(ctx context.Context, region, accessKey, secretAccessKey, sessionToken string, loadOptionFunctions ...func(*config.LoadOptions) error) error {
	this.ctx = ctx

	loadOptionFunctions = append(loadOptionFunctions, config.WithRegion(region))
	loadOptionFunctions = append(loadOptionFunctions, config.WithCredentialsProvider(
		credentials.NewStaticCredentialsProvider(accessKey, secretAccessKey, sessionToken)))

	cfg, err := config.LoadDefaultConfig(this.ctx, loadOptionFunctions...)
	if err != nil {
		return err
	}

	this.client = aws_s3.NewFromConfig(cfg)

	return nil
}

// CreateBucket is create a bucket.
//
// ex) _, err := client.CreateBucket(bucketName, "dummy")
func (this *Client) CreateBucket(name, region string) (*aws_s3.CreateBucketOutput, error) {
	return this.client.CreateBucket(
		this.ctx,
		&aws_s3.CreateBucketInput{
			Bucket: aws.String(name),
			CreateBucketConfiguration: &types.CreateBucketConfiguration{
				LocationConstraint: types.BucketLocationConstraint(region),
			}})
}

// ListBuckets is get buckets.
//
// ex) output, err := client.ListBuckets()
func (this *Client) ListBuckets() (*aws_s3.ListBucketsOutput, error) {
	return this.client.ListBuckets(this.ctx, &aws_s3.ListBucketsInput{})
}

// DeleteBucket is delete a bucket.
//
// ex) _, err := client.DeleteBucket(bucketName)
func (this *Client) DeleteBucket(name string) (*aws_s3.DeleteBucketOutput, error) {
	return this.client.DeleteBucket(this.ctx, &aws_s3.DeleteBucketInput{Bucket: aws.String(name)})
}

// PutObject is put an object.
//
// ex) _, err := client.PutObject(bucketName, key, data)
func (this *Client) PutObject(bucketName, key, body string) (*aws_s3.PutObjectOutput, error) {
	return this.client.PutObject(
		this.ctx,
		&aws_s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
			Body:   strings.NewReader(body)})
}

// GetObject is get an object.
//
// ex) output, err := client.GetObject(bucketName, key)
func (this *Client) GetObject(bucketName string, key string) (*aws_s3.GetObjectOutput, error) {
	return this.client.GetObject(
		this.ctx,
		&aws_s3.GetObjectInput{Bucket: aws.String(bucketName), Key: aws.String(key)})
}

// DeleteObject is delete an object.
//
// ex) _, err := client.DeleteObject(bucketName, key)
func (this *Client) DeleteObject(bucketName string, key string) (*aws_s3.DeleteObjectOutput, error) {
	return this.client.DeleteObject(
		this.ctx,
		&aws_s3.DeleteObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key)})
}
