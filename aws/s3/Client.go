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

type Client struct {
	ctx context.Context

	client *aws_s3.Client
}

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

func (this *Client) CreateBucket(name, region string) (*aws_s3.CreateBucketOutput, error) {
	return this.client.CreateBucket(
		this.ctx,
		&aws_s3.CreateBucketInput{
			Bucket: aws.String(name),
			CreateBucketConfiguration: &types.CreateBucketConfiguration{
				LocationConstraint: types.BucketLocationConstraint(region),
			}})
}

func (this *Client) ListBuckets() (*aws_s3.ListBucketsOutput, error) {
	return this.client.ListBuckets(this.ctx, &aws_s3.ListBucketsInput{})
}

func (this *Client) DeleteBucket(name string) (*aws_s3.DeleteBucketOutput, error) {
	return this.client.DeleteBucket(this.ctx, &aws_s3.DeleteBucketInput{Bucket: aws.String(name)})
}

func (this *Client) PutObject(bucketName, key, body string) (*aws_s3.PutObjectOutput, error) {
	return this.client.PutObject(
		this.ctx,
		&aws_s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
			Body:   strings.NewReader(body)})
}

func (this *Client) GetObject(bucketName string, key string) (*aws_s3.GetObjectOutput, error) {
	return this.client.GetObject(
		this.ctx,
		&aws_s3.GetObjectInput{Bucket: aws.String(bucketName), Key: aws.String(key)})
}

func (this *Client) DeleteObject(bucketName string, key string) (*aws_s3.DeleteObjectOutput, error) {
	return this.client.DeleteObject(
		this.ctx,
		&aws_s3.DeleteObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key)})
}
