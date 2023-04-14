package props

import (
	"github.com/SumoLogic-Labs/props/pkg/aws/s3"
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/magiconair/properties"
)

type S3 struct {
	*s3.Client
	bucket string
	key    string
}

func (s S3) Poll(ctx context.Context) (*properties.Properties, error) {
	r, err := s.Client.Download(ctx, s.bucket, s.key)
	if err != nil {
		return nil, fmt.Errorf("unable to poll s3 bucket %s for %s: %w", s.bucket, s.key, err)
	}
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("unable to read props: %w", err)
	}
	l := &properties.Loader{
		DisableExpansion: true,
		Encoding:         properties.UTF8,
	}
	return l.LoadBytes(b)
}

func NewS3Source(cfg aws.Config, bucket, key string) *S3 {
	return &S3{
		Client: s3.New(cfg),
		bucket: bucket,
		key:    key,
	}
}
