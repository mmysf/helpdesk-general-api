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
	"golang.org/x/crypto/bcrypt"
)

func (u *agentUsecase) ChangePassword(ctx context.Context, claim domain.JWTClaimAgent, payload domain.ChangePasswordRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	errValidation := make(map[string]string)
	// validating request
	if payload.OldPassword == "" {
		errValidation["oldPassword"] = "old password field is required"
	}
	if payload.NewPassword == "" {
		errValidation["newPassword"] = "new password field is required"
	}
	if payload.NewPassword == payload.OldPassword {
		errValidation["newPassword"] = "new password must be different from old password"
	}
	if len(payload.NewPassword) < 8 {
		errValidation["newPassword"] = "new password must be at least 8 characters"
	}
	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	// check agent
	agent, err := u.mongodbRepo.FetchOneAgent(ctx, map[string]interface{}{
		"id": claim.UserID,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if agent == nil {
		return response.Error(http.StatusUnauthorized, "user not found")
	}

	// check password
	if err = bcrypt.CompareHashAndPassword([]byte(agent.Password), []byte(payload.OldPassword)); err != nil {
		return response.Error(http.StatusBadRequest, "Wrong password")
	}

	// hash password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(payload.NewPassword), bcrypt.DefaultCost)

	// update password
	agent.Password = string(hashedPassword)
	err = u.mongodbRepo.UpdateAgent(ctx, agent)
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error())
	}

	return response.Success(agent)
}

func (u *agentUsecase) ChangeDomain(ctx context.Context, claim domain.JWTClaimAgent, payload domain.ChangeDomainRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	errValidation := make(map[string]string)
	// validating request
	if payload.IsCustom == nil {
		errValidation["isCustom"] = "isCustom field only accept true or false"
	}

	if payload.IsCustom != nil && !*payload.IsCustom {
		if payload.Subdomain == "" {
			errValidation["subdomain"] = "subdomain field is required"
		} else if !helpers.IsValidSubdomain(payload.Subdomain) {
			errValidation["subdomain"] = "subdomain only accept alphanumeric and must not start or end with dash"
		}
	}

	if payload.IsCustom != nil && *payload.IsCustom {
		if payload.FullUrl == "" {
			errValidation["fullUrl"] = "fullUrl field is required"
		} else if !helpers.IsValidURLWithoutProtocol(payload.FullUrl) {
			errValidation["fullUrl"] = "fullUrl is invalid"
		}
	}

	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	// check company
	company, err := u.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{
		"id": claim.CompanyID,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if company == nil {
		return response.Error(http.StatusBadRequest, "company not found")
	}

	if company.Type == "B2C" {
		return response.Error(http.StatusBadRequest, "company type must be B2B")
	}

	// check domain
	if !*payload.IsCustom {
		subdomain, err := u.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{
			"subdomain": payload.Subdomain,
		})
		if err != nil {
			return response.Error(http.StatusInternalServerError, err.Error())
		}
		if subdomain != nil && subdomain.ID.Hex() != claim.CompanyID {
			return response.Error(http.StatusBadRequest, "subdomain is already taken")
		}

		// get from config
		config := u._CacheConfig(ctx)

		// check blacklist subdomain
		for _, blacklisted := range config.BlacklistSubdomain {
			if payload.Subdomain == blacklisted {
				return response.Error(http.StatusBadRequest, "subdomain is prohibited")
			}
		}

		payload.FullUrl = payload.Subdomain + "." + config.MainDomain
	} else {
		payload.Subdomain = ""
		fullUrl, err := u.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{
			"fullUrl": payload.FullUrl,
		})
		if err != nil {
			return response.Error(http.StatusInternalServerError, err.Error())
		}
		if fullUrl != nil && fullUrl.ID.Hex() != claim.CompanyID {
			return response.Error(http.StatusBadRequest, "fullUrl is already taken")
		}
	}

	company.Settings.Domain.IsCustom = *payload.IsCustom
	company.Settings.Domain.Subdomain = payload.Subdomain
	company.Settings.Domain.FullUrl = payload.FullUrl

	if err = u.mongodbRepo.UpdatePartialCompany(
		ctx,
		map[string]interface{}{
			"id": claim.CompanyID,
		},
		map[string]interface{}{
			"settings.domain.isCustom":  company.Settings.Domain.IsCustom,
			"settings.domain.subdomain": company.Settings.Domain.Subdomain,
			"settings.domain.fullUrl":   company.Settings.Domain.FullUrl,
		}); err != nil {
		return response.Error(http.StatusBadRequest, err.Error())
	}

	return response.Success(company)
}

func (u *agentUsecase) UpdateProfile(ctx context.Context, claim domain.JWTClaimAgent, payload domain.UpdateProfileRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	errValidation := make(map[string]string)
	// validating request
	if payload.Name == "" {
		errValidation["name"] = "name field is required"
	}

	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	// check agent
	agent, err := u.mongodbRepo.FetchOneAgent(ctx, map[string]interface{}{
		"id": claim.UserID,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if agent == nil {
		return response.Error(http.StatusBadRequest, "user not found")
	}

	// check media
	if payload.AttachId != "" {
		media, err := u.mongodbRepo.FetchOneMedia(ctx, map[string]interface{}{
			"id":       payload.AttachId,
			"category": model.Avatar,
		})
		if err != nil {
			return response.Error(http.StatusInternalServerError, err.Error())
		}
		if media == nil {
			return response.Error(http.StatusBadRequest, "avatar media not found")
		}

		agent.ProfilePicture = model.MediaFK{
			ID:          media.ID.Hex(),
			Name:        media.Name,
			Size:        media.Size,
			URL:         media.URL,
			Type:        media.Type,
			Category:    media.Category,
			IsPrivate:   media.IsPrivate,
			ProviderKey: media.ProviderKey,
		}
	}

	// update profile
	agent.Name = payload.Name
	agent.Bio = payload.Bio
	agent.Contact = payload.Contact

	if err = u.mongodbRepo.UpdateAgent(ctx, agent); err != nil {
		return response.Error(http.StatusBadRequest, err.Error())
	}

	return response.Success(agent)
}

func (u *agentUsecase) ChangeColor(ctx context.Context, claim domain.JWTClaimAgent, payload domain.ChangeColorMode) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	config := u._CacheConfig(ctx)

	errValidation := make(map[string]string)

	// validating request
	if payload.Dark.Primary != "" {
		if !helpers.IsValidHexColor(payload.Dark.Primary) {
			errValidation["dark.primary"] = "colorMode dark primary field is invalid"
		}
	} else {
		payload.Dark.Primary = config.DefaultColor.Dark.Primary
	}

	if payload.Light.Primary != "" {
		if !helpers.IsValidHexColor(payload.Light.Primary) {
			errValidation["light.primary"] = "colorMode light primary field is invalid"
		}
	} else {
		payload.Light.Primary = config.DefaultColor.Light.Primary
	}

	if payload.Dark.Secondary != "" {
		if !helpers.IsValidHexColor(payload.Dark.Secondary) {
			errValidation["dark.secondary"] = "colorMode dark secondary field is invalid"
		}
	} else {
		payload.Dark.Secondary = config.DefaultColor.Dark.Secondary
	}

	if payload.Light.Secondary != "" {
		if !helpers.IsValidHexColor(payload.Light.Secondary) {
			errValidation["light.secondary"] = "colorMode light secondary field is invalid"
		}
	} else {
		payload.Light.Secondary = config.DefaultColor.Light.Secondary
	}

	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	company, err := u.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{"id": claim.CompanyID})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if company == nil {
		return response.Error(http.StatusBadRequest, "company not found")
	}

	company.Settings.ColorMode.Light.Primary = payload.Light.Primary
	company.Settings.ColorMode.Light.Secondary = payload.Light.Secondary
	company.Settings.ColorMode.Dark.Primary = payload.Dark.Primary
	company.Settings.ColorMode.Dark.Secondary = payload.Dark.Secondary

	if err = u.mongodbRepo.UpdatePartialCompany(
		ctx,
		map[string]interface{}{"id": claim.CompanyID},
		map[string]interface{}{
			"settings.colorMode.light.primary":   company.Settings.ColorMode.Light.Primary,
			"settings.colorMode.light.secondary": company.Settings.ColorMode.Light.Secondary,
			"settings.colorMode.dark.primary":    company.Settings.ColorMode.Dark.Primary,
			"settings.colorMode.dark.secondary":  company.Settings.ColorMode.Dark.Secondary,
		}); err != nil {
		return response.Error(http.StatusBadRequest, err.Error())
	}

	return response.Success(company)
}

func (u *agentUsecase) UploadAgentProfilePicture(ctx context.Context, claim domain.JWTClaimAgent, payload domain.UploadAttachment, request *http.Request) response.Base {
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
	if !helpers.InArrayString(typeDocument, domain.AllowedImgMimeTypes) {
		validation["file"] = "field file is not valid type"
	}

	fileSize = uploadedFile.Size
	maxFileSize := int64(1 * 1024 * 1024) // 1 MB in bytes

	if fileSize > maxFileSize {
		validation["file"] = "file size exceeds the maximum limit of 1 MB"
	}
	defer file.Close()

	if len(validation) > 0 {
		return response.ErrorValidation(validation, "error validation")
	}

	category := helpers.GetCategoryByContentType(typeDocument, true)
	year, month, _ := time.Now().Date()
	objName = "avatars/" + category + "/" + strconv.Itoa(year) + "/" + strconv.Itoa(int(month)) + "/" + helpers.GenerateCleanName(uploadedFile.Filename)

	uploadData, err := u.s3Repo.UploadFilePublic(objName, file, typeDocument)

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	typeDocument = helpers.GetCategoryByContentType(typeDocument, false)

	now := time.Now().UTC()
	media := model.Media{
		ID:           primitive.NewObjectID(),
		Name:         payload.Title,
		Provider:     "s3",
		ProviderKey:  objName,
		Type:         typeDocument,
		Category:     model.Avatar,
		Size:         fileSize,
		URL:          helpers.GeneratePublicURL(uploadData.URL),
		ExpiredUrlAt: nil,
		IsPrivate:    false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// insert into db
	err = u.mongodbRepo.CreateMedia(ctx, &media)
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	return response.Success(media)
}
