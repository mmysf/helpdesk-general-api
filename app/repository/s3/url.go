package s3Repo

import (
	s3_model "app/domain/model/s3"
	"context"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sirupsen/logrus"
)

func (r *s3Repo) GetPublicLink(objectKey string) string {
	url := &url.URL{}
	if r.publicURL == nil {
		endpoint := r.client.Options().BaseEndpoint
		url, _ = url.Parse(*endpoint)
	} else {
		url = r.publicURL
	}

	// add path with object key
	url.Path = objectKey

	return url.String()
}

func (r *s3Repo) GetPresignedLink(objectKey string, expires *time.Duration) (uploadData *s3_model.UploadResponse, err error) {
	now := time.Now().UTC()
	var expiredUrlAt *time.Time

	resSigned, err := r.presigner.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		var exp *time.Duration

		if expires != nil {
			exp = expires
		} else {
			sec, _ := strconv.Atoi(os.Getenv("S3_EXPIRES_TIME"))
			if sec == 0 {
				sec = 300 // default 5 minutes
			}

			duration := time.Duration(sec) * time.Second
			exp = &duration
		}

		// Adjust expiration time to now + exp - 15 seconds
		adjustedExp := now.Add(*exp - 15*time.Second)
		expiredUrlAt = &adjustedExp

		opts.Expires = *exp
	})

	if err != nil {
		logrus.Error("GetPresignedLink error: ", err)
		return nil, err
	}

	return &s3_model.UploadResponse{URL: resSigned.URL, Key: objectKey, ExpiredAt: expiredUrlAt}, nil
}
