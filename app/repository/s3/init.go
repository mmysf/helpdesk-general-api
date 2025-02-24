package s3repo

import (
	s3_model "app/domain/model/s3"
	"context"
	"io"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sirupsen/logrus"
)

type s3Repo struct {
	bucketName string
	publicURL  *url.URL
	client     *s3.Client
	presigner  *s3.PresignClient
}

func NewS3Repo() S3Repo {
	cfg, _ := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(
			aws.NewCredentialsCache(
				credentials.NewStaticCredentialsProvider(
					os.Getenv("S3_ACCESS_KEY"),
					os.Getenv("S3_SECRET_KEY"),
					"",
				),
			),
		),
		config.WithRegion(os.Getenv("S3_REGION")),
	)

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(os.Getenv("S3_ENDPOINT"))
	})

	publicURL, err := url.Parse(os.Getenv("S3_PUBLIC_URL"))
	if err != nil {
		logrus.Info("s3: without public url", err)
	}

	return &s3Repo{
		bucketName: os.Getenv("S3_BUCKET_NAME"),
		publicURL:  publicURL,
		client:     client,
		presigner:  s3.NewPresignClient(client),
	}
}

type S3Repo interface {
	GetPresignedLink(objectKey string, expires *time.Duration) string
	GetPublicLink(objectKey string) string
	UploadFilePublic(objectKey string, body io.Reader, contentType string) (uploadData *s3_model.UploadResponse, err error)
	UploadFilePrivate(objectKey string, body io.Reader, contentType string, expires *time.Duration) (uploadData *s3_model.UploadResponse, err error)
}
