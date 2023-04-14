package s3

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	client *s3.Client
}

func (u Client) Upload(ctx context.Context, bucket, filename string, r io.Reader) error {
	_, err := u.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   r,
	})
	return err
}

func (u Client) Download(ctx context.Context, bucket, filename string) (io.ReadCloser, error) {
	res, err := u.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to get object: %w", err)
	}
	return res.Body, nil
}

func New(cfg aws.Config) *Client {
	return &Client{s3.NewFromConfig(cfg)}
}
