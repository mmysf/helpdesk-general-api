package usecase_agent

import (
	"app/domain"
	"app/domain/model"
	"app/helpers"
	"context"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (u *agentUsecase) UploadAttachment(ctx context.Context, claim domain.JWTClaimAgent, payload domain.UploadAttachment, request *http.Request) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// declare variables
	var err error
	validation := make(map[string]string)
	var typeDocument, objName string
	var file multipart.File
	var uploadedFile *multipart.FileHeader
	var fileSize int64

	// validation
	file, uploadedFile, err = request.FormFile("file")
	if err != nil {
		validation["file"] = "file field is required"
	}

	typeDocument = uploadedFile.Header.Get("Content-Type")
	if !helpers.InArrayString(typeDocument, domain.AllowedMimeTypes) {
		validation["file"] = "field file is not valid type"
	}

	fileSize = uploadedFile.Size
	maxFileSize := int64(10 * 1024 * 1024) // 10 MB in bytes

	if fileSize > maxFileSize {
		validation["file"] = "file size exceeds the maximum limit of 10 MB"
	}
	defer file.Close()

	if len(validation) > 0 {
		return response.ErrorValidation(validation, "error validation")
	}

	category := helpers.GetCategoryByContentType(typeDocument, true)
	year, month, _ := time.Now().Date()
	objName = "attachments/" + category + "/" + strconv.Itoa(year) + "/" + strconv.Itoa(int(month)) + "/" + helpers.GenerateCleanName(uploadedFile.Filename)

	uploadData, err := u.s3Repo.UploadFilePrivate(objName, file, typeDocument, nil)
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	typeDocument = helpers.GetCategoryByContentType(typeDocument, false)

	// check file really uploaded public by google storage
	err = helpers.DetectFileExistsByURL(uploadData.URL)
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	now := time.Now().UTC()
	attachment := model.Attachment{
		ID:           primitive.NewObjectID(),
		Company:      claim.Company,
		IsPrivate:    true,
		Name:         payload.Title,
		Provider:     "s3",
		ProviderKey:  objName,
		Type:         typeDocument,
		Size:         fileSize,
		URL:          uploadData.URL,
		ExpiredUrlAt: uploadData.ExpiredAt,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// insert into db
	err = u.mongodbRepo.CreateAttachment(ctx, &attachment)
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	attachment.URL = uploadData.URL
	return response.Success(attachment)
}

func (u *agentUsecase) GetAttachmentDetail(ctx context.Context, claim domain.JWTClaimAgent, attachmentId string) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// check attachment
	attachment, err := u.mongodbRepo.FetchOneAttachment(ctx, map[string]interface{}{
		"id": attachmentId,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if attachment == nil {
		return response.Error(http.StatusBadRequest, "attachment not found")
	}

	updateAttachmentPartial := map[string]interface{}{
		"url":          "",
		"expiredUrlAt": nil,
	}
	isUpdate := false

	if attachment.IsPrivate {
		// check if url expired or not
		if attachment.ExpiredUrlAt == nil || time.Now().UTC().After(*attachment.ExpiredUrlAt) || attachment.URL == "" {
			url, err := u.s3Repo.GetPresignedLink(attachment.ProviderKey, nil)
			if err != nil {
				return response.Error(http.StatusBadRequest, err.Error())
			}

			// // check file really uploaded public by google storage
			// err = helpers.DetectFileExistsByURL(url.URL)
			// if err != nil {
			// 	return response.Error(http.StatusBadRequest, err.Error())
			// }

			attachment.URL = url.URL
			attachment.ExpiredUrlAt = url.ExpiredAt

			updateAttachmentPartial["url"] = url.URL
			updateAttachmentPartial["expiredUrlAt"] = url.ExpiredAt

			isUpdate = true
		}
	} else {
		// check if attachment.URL is can access
		err = helpers.DetectFileExistsByURL(attachment.URL)
		if err != nil {
			url := helpers.GeneratePublicURL(attachment.ProviderKey)
			if url == "" {
				return response.Error(http.StatusBadRequest, "attachment not found")
			}

			attachment.URL = url
			attachment.ExpiredUrlAt = nil

			updateAttachmentPartial["url"] = url

			isUpdate = true
		}

	}

	if isUpdate {
		go (func() {
			u.mongodbRepo.UpdateAttachmentPartial(ctx, attachment.ID, updateAttachmentPartial)
		})()
	}

	return response.Success(attachment)
}
