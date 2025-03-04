package fileupload

import (
	"context"
	c "go-starter/config"
	"go-starter/internal/domain/ports"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

// S3Adapter is an adapter for the ports.FileUploadAdapter interface.
type S3Adapter struct {
	uploader   *manager.Uploader
	errTracker ports.ErrTrackerAdapter
	cfg        *c.FileUpload
}

// NewS3Adapter creates a new S3Adapter instance.
func NewS3Adapter(fileUploadCfg *c.FileUpload, errTracker ports.ErrTrackerAdapter) (*S3Adapter, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(fileUploadCfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			fileUploadCfg.AccessKey,
			fileUploadCfg.SecretKey,
			"",
		)),
	)

	if err != nil {
		errTracker.CaptureException(err)
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.HTTPClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	})
	uploader := manager.NewUploader(client)

	return &S3Adapter{uploader: uploader, errTracker: errTracker, cfg: fileUploadCfg}, nil
}

// Upload uploads a file to the S3 bucket.
// Returns the URL of the uploaded file or an error if the upload fails.
func (s *S3Adapter) Upload(ctx context.Context, key string, body io.Reader) (string, error) {
	result, err := s.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.cfg.Bucket),
		Key:    aws.String(key),
		Body:   body,
	})

	if err != nil {
		s.errTracker.CaptureException(err)
		return "", err
	}

	return result.Location, nil
}
