package s3_model

import "time"

type UploadResponse struct {
	Key         string     `json:"key"`
	ContentType string     `json:"content_type"`
	URL         string     `json:"url"`
	ExpiredAt   *time.Time `json:"expiredAt"`
}
