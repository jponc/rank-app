package s3manager

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Client interface
type Client interface {
	// Upload uploads content to s3 using concurrent upload
	Upload(input *s3manager.UploadInput) (*s3manager.UploadOutput, error)
}

type client struct {
	s3Manager *s3manager.Uploader
}

// NewClient instantiates an S3 client
func NewClient(awsRegion string) (Client, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})

	if err != nil {
		return nil, fmt.Errorf("cannot create aws session: %v", err)
	}

	s3Manager := s3manager.NewUploader(sess)

	c := &client{
		s3Manager: s3Manager,
	}

	return c, nil
}

func (c *client) Upload(input *s3manager.UploadInput) (*s3manager.UploadOutput, error) {
	return c.s3Manager.Upload(input)
}
