package helpers

import (
	"app/domain/model"
	"errors"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	s3repo "app/app/repository/s3"

	"github.com/gosimple/slug"
)

func DetectFileExistsByURL(url string) error {
	if url == "" {
		return errors.New("url cannot be empty")
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("page / media not found")
	}
	return nil
}

func DetectContentType(url string) (string, error) {
	if url == "" {
		return "", errors.New("url cannot be empty")
	}
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	return resp.Header.Get("Content-Type"), nil
}

func GetCategoryByContentType(contentType string, plural bool) string {
	if InArrayString(contentType, []string{"application/pdf", "application/msword", "application/vnd.openxmlformats-officedocument.wordprocessingml.document", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}) {
		if plural {
			return "documents"
		}
		return "document"
	}
	if InArrayString(contentType, []string{"image/jpeg", "image/png", "image/gif", "image/svg+xml", "image/webp"}) {
		if plural {
			return "images"
		}
		return "image"
	}
	if InArrayString(contentType, []string{"video/mp4", "video/mpeg", "video/ogg", "video/webm"}) {
		if plural {
			return "videos"
		}
		return "video"
	}
	if InArrayString(contentType, []string{"audio/mpeg", "audio/ogg", "audio/wav", "audio/webm"}) {
		if plural {
			return "audios"
		}
		return "audio"
	}

	return "uncategorized"
}

func GenerateCleanName(originalName string) string {
	nano := time.Now().UnixNano()
	extName := filepath.Ext(originalName)
	baseName := strings.TrimSuffix(originalName, extName)
	slugBaseName := slug.Make(baseName)
	name := strconv.Itoa(int(nano)) + "-" + slugBaseName + extName
	return name
}

func GeneratePublicURL(objectName string) string {
	return objectName
}

func GenerateLinkAttachment(attachments *[]model.AttachmentFK, s3 s3repo.S3Repo) (isUpdate bool, err error) {
	if attachments == nil {
		return
	}
	for i, attachment := range *attachments {
		if attachment.IsPrivate {
			// errUrl := DetectFileExistsByURL(attachment.URL)

			// check distance expiredAt with now
			// if attachment.ExpiredUrlAt == nil || time.Until(*attachment.ExpiredUrlAt) <= 15 || errUrl != nil {
			isUpdate = true

			url, err := s3.GetPresignedLink(attachment.ProviderKey, nil)
			if err != nil {
				return false, err
			}

			// check file really uploaded public by google storage
			err = DetectFileExistsByURL(url.URL)
			if err != nil {
				return false, err
			}

			(*attachments)[i].URL = url.URL
			// (*attachments)[i].ExpiredUrlAt = url.ExpiredAt
			// }
		}
	}

	return
}
