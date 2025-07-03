package usecase_superadmin

import (
	"app/domain"
	"app/domain/model"
	"app/helpers"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	yurekahelpers "github.com/Yureka-Teknologi-Cipta/yureka/helpers"
	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func (u *superadminUsecase) GetCompanyList(ctx context.Context, claim domain.JWTClaimSuperadmin, options map[string]interface{}) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	paramQuery := options["query"].(url.Values)
	page, limit, offset := yurekahelpers.GetLimitOffset(paramQuery)

	fetchOptions := map[string]interface{}{
		"limit":  limit,
		"offset": offset,
	}

	// filtering
	if paramQuery.Get("sort") != "" {
		fetchOptions["sort"] = paramQuery.Get("sort")
	}

	if paramQuery.Get("dir") != "" {
		fetchOptions["dir"] = paramQuery.Get("dir")
	}

	if paramQuery.Get("q") != "" {
		fetchOptions["q"] = paramQuery.Get("q")
	}

	if paramQuery.Get("search") != "" {
		fetchOptions["search"] = paramQuery.Get("search")
	}

	// count first
	totalDocuments := u.mongodbRepo.CountCompany(ctx, fetchOptions)

	if totalDocuments == 0 {
		return response.Success(response.List{
			List:  []interface{}{},
			Page:  page,
			Limit: limit,
			Total: totalDocuments,
		})
	}

	// check company list
	cur, err := u.mongodbRepo.FetchCompanyList(ctx, fetchOptions)

	if err != nil {
		return response.Success(response.List{
			List:  []interface{}{},
			Page:  page,
			Limit: limit,
			Total: totalDocuments,
		})
	}

	defer cur.Close(ctx)

	list := make([]interface{}, 0)
	for cur.Next(ctx) {
		row := model.Company{}
		err := cur.Decode(&row)
		if err != nil {
			logrus.Error("Company Decode ", err)
			return response.Success(response.List{
				List:  []interface{}{},
				Page:  page,
				Limit: limit,
				Total: totalDocuments,
			})
		}

		list = append(list, row)
	}

	return response.Success(response.List{
		List:  list,
		Page:  page,
		Limit: limit,
		Total: totalDocuments,
	})
}

func (u *superadminUsecase) GetCompanyDetail(ctx context.Context, claim domain.JWTClaimSuperadmin, options map[string]interface{}) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// get id
	companyID := options["id"]

	// check company
	company, err := u.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{
		"id": companyID,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if company == nil {
		return response.Error(http.StatusBadRequest, "Company not found")
	}

	return response.Success(company)
}

func (u *superadminUsecase) CreateCompany(ctx context.Context, claim domain.JWTClaimSuperadmin, options map[string]interface{}) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	config := u._CacheConfig(ctx)

	// payload
	payload := options["payload"].(domain.CreateCompanyRequest)

	errValidation := make(map[string]string)

	// validating request
	if payload.Name == "" {
		errValidation["name"] = "name field is required"
	}

	if payload.Email == "" {
		errValidation["email"] = "email field is required"
	} else if !helpers.IsValidEmail(payload.Email) {
		errValidation["email"] = "email field is invalid"
	}

	if payload.Domain.IsCustom == nil {
		errValidation["domain.isCustom"] = "isCustom field only accept true or false"
	}

	if payload.Domain.IsCustom != nil && !*payload.Domain.IsCustom {
		if payload.Domain.Subdomain == "" {
			errValidation["domain.subdomain"] = "subdomain field is required"
		} else if !helpers.IsValidSubdomain(payload.Domain.Subdomain) {
			errValidation["domain.subdomain"] = "subdomain only accept alphanumeric and must not start or end with dash"
		}
	}

	if payload.Domain.IsCustom != nil && *payload.Domain.IsCustom {
		if payload.Domain.FullUrl == "" {
			errValidation["domain.fullUrl"] = "fullUrl field is required"
		} else if !helpers.IsValidURLWithoutProtocol(payload.Domain.FullUrl) {
			errValidation["domain.fullUrl"] = "fullUrl is invalid"
		}
	}

	if payload.LogoAttachId == "" {
		errValidation["logoAttachId"] = "logoAttachId field is required"
	}

	if payload.ColorMode.Dark.Primary != "" {
		if !helpers.IsValidHexColor(payload.ColorMode.Dark.Primary) {
			errValidation["colorMode.dark.primary"] = "colorMode dark primary field is invalid"
		}
	} else {
		payload.ColorMode.Dark.Primary = config.DefaultColor.Dark.Primary
	}

	if payload.ColorMode.Light.Primary != "" {
		if !helpers.IsValidHexColor(payload.ColorMode.Light.Primary) {
			errValidation["colorMode.light.primary"] = "colorMode light primary field is invalid"
		}
	} else {
		payload.ColorMode.Light.Primary = config.DefaultColor.Light.Primary
	}

	if payload.ColorMode.Dark.Secondary != "" {
		if !helpers.IsValidHexColor(payload.ColorMode.Dark.Secondary) {
			errValidation["colorMode.dark.secondary"] = "colorMode dark secondary field is invalid"
		}
	} else {
		payload.ColorMode.Dark.Secondary = config.DefaultColor.Dark.Secondary
	}

	if payload.ColorMode.Light.Secondary != "" {
		if !helpers.IsValidHexColor(payload.ColorMode.Light.Secondary) {
			errValidation["colorMode.light.secondary"] = "colorMode light secondary field is invalid"
		}
	} else {
		payload.ColorMode.Light.Secondary = config.DefaultColor.Light.Secondary
	}

	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	// check domain
	if !*payload.Domain.IsCustom {
		subdomain, err := u.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{
			"subdomain": payload.Domain.Subdomain,
		})
		if err != nil {
			return response.Error(http.StatusInternalServerError, err.Error())
		}
		if subdomain != nil {
			return response.Error(http.StatusBadRequest, "subdomain is already taken")
		}

		// get from config
		config := u._CacheConfig(ctx)

		// check blacklist subdomain
		for _, blacklisted := range config.BlacklistSubdomain {
			if payload.Domain.Subdomain == blacklisted {
				return response.Error(http.StatusBadRequest, "subdomain is prohibited")
			}
		}

		payload.Domain.FullUrl = payload.Domain.Subdomain + "." + config.MainDomain
	} else {
		fullUrl, err := u.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{
			"fullUrl": payload.Domain.FullUrl,
		})
		if err != nil {
			return response.Error(http.StatusInternalServerError, err.Error())
		}
		if fullUrl != nil {
			return response.Error(http.StatusBadRequest, "fullUrl is already taken")
		}
	}

	// check email
	email, err := u.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{
		"email": payload.Email,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if email != nil {
		return response.Error(http.StatusBadRequest, "email already in use")
	}

	// check logo
	logo, err := u.mongodbRepo.FetchOneMedia(ctx, map[string]interface{}{
		"id":       payload.LogoAttachId,
		"category": model.CompanyLogo,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if logo == nil {
		return response.Error(http.StatusBadRequest, "Company logo media not found")
	}

	// check agent
	agent, err := u.mongodbRepo.FetchOneAgent(ctx, map[string]interface{}{
		"email": payload.Email,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if agent != nil {
		return response.Error(http.StatusBadRequest, "email already in use for agent "+agent.Company.Name)
	}

	// generate code
	code, err := u._generateCode(payload.Name)

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	// create company
	company := model.Company{
		ID:        primitive.NewObjectID(),
		AccessKey: uuid.Must(uuid.NewRandom()).String(),
		Name:      payload.Name,
		Bio:       "",
		Type:      "B2B",
		Code:      code,
		Logo: model.MediaFK{
			ID:          logo.ID.Hex(),
			Name:        logo.Name,
			Size:        logo.Size,
			URL:         logo.URL,
			Type:        logo.Type,
			Category:    logo.Category,
			ProviderKey: logo.ProviderKey,
			IsPrivate:   logo.IsPrivate,
		},
		Settings: model.CompanySeting{
			Code:  uuid.Must(uuid.NewRandom()).String(),
			Email: payload.Email,
			ColorMode: model.ColorMode{
				Light: model.Color{
					Primary:   payload.ColorMode.Light.Primary,
					Secondary: payload.ColorMode.Light.Secondary,
				},
				Dark: model.Color{
					Primary:   payload.ColorMode.Dark.Primary,
					Secondary: payload.ColorMode.Dark.Secondary,
				},
			},
			Domain: model.CompanyDomain{
				IsCustom:  *payload.Domain.IsCustom,
				Subdomain: payload.Domain.Subdomain,
				FullUrl:   payload.Domain.FullUrl,
			},
			SMTP: model.SMTP{
				FromAddress: config.Email.SenderEmail,
				FromName:    config.Email.SenderName,
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = u.mongodbRepo.CreateCompany(ctx, &company)

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	password := helpers.RandomString(5)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// create agent
	newAgent := &model.Agent{
		ID:       primitive.NewObjectID(),
		Name:     payload.Name,
		Email:    payload.Email,
		JobTitle: "Admin",
		Password: string(hashedPassword),
		Role:     model.AdminRole,
		Company: model.CompanyNested{
			ID:    company.ID.Hex(),
			Name:  company.Name,
			Image: company.Logo.URL,
			Type:  company.Type,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = u.mongodbRepo.CreateAgent(ctx, newAgent)

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	go _sendEmailAgentCrendential(config, *newAgent, password, &company)

	go u.mongodbRepo.UpdateManyMediaPartial(ctx, []primitive.ObjectID{logo.ID}, map[string]interface{}{
		"isUsed": true,
	})

	return response.Success(company)
}

func (u *superadminUsecase) _generateCode(name string) (string, error) {
	words := strings.Fields(name)
	baseCode := ""

	// Take first letters of words
	for _, word := range words {
		if len(baseCode) < 3 && len(word) > 0 {
			baseCode += strings.ToUpper(string(word[0]))
		}
	}

	if len(baseCode) < 3 && len(words) > 0 {
		firstWord := strings.ToUpper(words[0])
		for i := 1; len(baseCode) < 3 && i < len(firstWord); i++ {
			baseCode += string(firstWord[i])
		}
	}

	for i := -1; i < 10000; i++ {
		code := baseCode
		if i >= 0 {
			code = baseCode + letterSuffix(i)
		}

		company, err := u.mongodbRepo.FetchOneCompany(context.Background(), map[string]interface{}{
			"code": code,
		})
		if err != nil {
			return "", err
		}
		if company == nil {
			return code, nil // unique!
		}
	}

	return "", fmt.Errorf("unable to generate unique code for: %s", name)
}

func letterSuffix(n int) string {
	result := ""
	for n >= 0 {
		result = string('A'+(n%26)) + result
		n = n/26 - 1
		if n < 0 {
			break
		}
	}
	return result
}

func _sendEmailAgentCrendential(config model.Config, agent model.Agent, password string, company *model.Company) {
	loginLink := helpers.StringReplacer(config.LoginLink, map[string]string{
		"base_url_frontend": config.AgentLink,
	})
	// send email
	mail := helpers.NewSMTPMailer(company)
	mail.To([]string{agent.Email})
	mail.Subject(config.Email.Template.DefaultUser.Title)
	mail.Body(helpers.StringReplacer(config.Email.Template.DefaultUser.Body, map[string]string{
		"title":      config.Email.Template.DefaultUser.Title,
		"fe_link":    config.AgentLink,
		"email":      agent.Email,
		"password":   password,
		"login_link": loginLink,
	}))

	// send
	if err := mail.Send(); err != nil {
		logrus.Errorf("Send Email to %s error %v", agent.Email, err)
	}
}

func (u *superadminUsecase) UpdateCompany(ctx context.Context, claim domain.JWTClaimSuperadmin, options map[string]interface{}) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	config := u._CacheConfig(ctx)

	// payload
	payload := options["payload"].(domain.CreateCompanyRequest)

	errValidation := make(map[string]string)

	// validating request
	if payload.Name == "" {
		errValidation["name"] = "name field is required"
	}

	if payload.Email == "" {
		errValidation["email"] = "email field is required"
	} else if !helpers.IsValidEmail(payload.Email) {
		errValidation["email"] = "email field is invalid"
	}

	if payload.Domain.IsCustom == nil {
		errValidation["isCustom"] = "isCustom field only accept true or false"
	}

	if payload.Domain.IsCustom != nil && !*payload.Domain.IsCustom {
		if payload.Domain.Subdomain == "" {
			errValidation["subdomain"] = "subdomain field is required"
		} else if !helpers.IsValidSubdomain(payload.Domain.Subdomain) {
			errValidation["subdomain"] = "subdomain only accept alphanumeric and must not start or end with dash"
		}
	}

	if payload.Domain.IsCustom != nil && *payload.Domain.IsCustom {
		if payload.Domain.FullUrl == "" {
			errValidation["fullUrl"] = "fullUrl field is required"
		} else if !helpers.IsValidURLWithoutProtocol(payload.Domain.FullUrl) {
			errValidation["fullUrl"] = "fullUrl is invalid"
		}
	}

	if payload.LogoAttachId == "" {
		errValidation["logoAttachId"] = "logoAttachId field is required"
	}

	if payload.ColorMode.Dark.Primary != "" {
		if !helpers.IsValidHexColor(payload.ColorMode.Dark.Primary) {
			errValidation["colorMode.dark.primary"] = "colorMode dark primary field is invalid"
		}
	} else {
		payload.ColorMode.Dark.Primary = config.DefaultColor.Dark.Primary
	}

	if payload.ColorMode.Light.Primary != "" {
		if !helpers.IsValidHexColor(payload.ColorMode.Light.Primary) {
			errValidation["colorMode.light.primary"] = "colorMode light primary field is invalid"
		}
	} else {
		payload.ColorMode.Light.Primary = config.DefaultColor.Light.Primary
	}

	if payload.ColorMode.Dark.Secondary != "" {
		if !helpers.IsValidHexColor(payload.ColorMode.Dark.Secondary) {
			errValidation["colorMode.dark.secondary"] = "colorMode dark secondary field is invalid"
		}
	} else {
		payload.ColorMode.Dark.Secondary = config.DefaultColor.Dark.Secondary
	}

	if payload.ColorMode.Light.Secondary != "" {
		if !helpers.IsValidHexColor(payload.ColorMode.Light.Secondary) {
			errValidation["colorMode.light.secondary"] = "colorMode light secondary field is invalid"
		}
	} else {
		payload.ColorMode.Light.Secondary = config.DefaultColor.Light.Secondary
	}

	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	// check company
	company, err := u.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{
		"id": options["id"],
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
	if !*payload.Domain.IsCustom {
		subdomain, err := u.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{
			"subdomain": payload.Domain.Subdomain,
		})
		if err != nil {
			return response.Error(http.StatusInternalServerError, err.Error())
		}
		if subdomain != nil && subdomain.ID.Hex() != company.ID.Hex() {
			return response.Error(http.StatusBadRequest, "subdomain is already taken")
		}

		// get from config
		config := u._CacheConfig(ctx)

		// check blacklist subdomain
		for _, blacklisted := range config.BlacklistSubdomain {
			if payload.Domain.Subdomain == blacklisted {
				return response.Error(http.StatusBadRequest, "subdomain is prohibited")
			}
		}

		payload.Domain.FullUrl = payload.Domain.Subdomain + "." + config.MainDomain
	} else {
		fullUrl, err := u.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{
			"fullUrl": payload.Domain.FullUrl,
		})
		if err != nil {
			return response.Error(http.StatusInternalServerError, err.Error())
		}
		if fullUrl != nil && fullUrl.ID.Hex() != company.ID.Hex() {
			return response.Error(http.StatusBadRequest, "fullUrl is already taken")
		}
	}

	// check email
	email, err := u.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{
		"email": payload.Email,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if email != nil && email.ID.Hex() != company.ID.Hex() {
		return response.Error(http.StatusBadRequest, "email already in use")
	}

	// check logoa
	logo, err := u.mongodbRepo.FetchOneMedia(ctx, map[string]interface{}{
		"id":       payload.LogoAttachId,
		"category": model.CompanyLogo,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if logo == nil {
		return response.Error(http.StatusBadRequest, "Company logo media not found")
	}

	// uopdate company
	company.Name = payload.Name
	company.Logo = model.MediaFK{
		ID:          logo.ID.Hex(),
		Name:        logo.Name,
		Size:        logo.Size,
		URL:         logo.URL,
		Type:        logo.Type,
		Category:    logo.Category,
		ProviderKey: logo.ProviderKey,
		IsPrivate:   logo.IsPrivate,
	}
	company.Settings.Email = payload.Email
	company.Settings.ColorMode = model.ColorMode{
		Light: model.Color{
			Primary:   payload.ColorMode.Light.Primary,
			Secondary: payload.ColorMode.Light.Secondary,
		},
		Dark: model.Color{
			Primary:   payload.ColorMode.Dark.Primary,
			Secondary: payload.ColorMode.Dark.Secondary,
		},
	}
	company.Settings.Domain = model.CompanyDomain{
		IsCustom:  *payload.Domain.IsCustom,
		Subdomain: payload.Domain.Subdomain,
		FullUrl:   payload.Domain.FullUrl,
	}
	company.UpdatedAt = time.Now()

	err = u.mongodbRepo.UpdateCompany(ctx, company)

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	// update company nested
	// go u.mongodbRepo.UpdatePartialCompanyProduct(ctx, map[string]any{
	// 	"companyID": company.ID.Hex(),
	// }, map[string]any{
	// 	"company.name":  payload.Name,
	// 	"company.image": company.Logo.URL,
	// })
	go u.mongodbRepo.UpdatePartialAgent(ctx, map[string]any{
		"companyID": company.ID.Hex(),
	}, map[string]any{
		"company.name":  payload.Name,
		"company.image": company.Logo.URL,
	})
	go u.mongodbRepo.UpdatePartialCustomer(ctx, map[string]any{
		"companyID": company.ID.Hex(),
	}, map[string]any{
		"company.name":  payload.Name,
		"company.image": company.Logo.URL,
	})

	return response.Success(company)
}

func (u *superadminUsecase) DeleteCompany(ctx context.Context, claim domain.JWTClaimSuperadmin, options map[string]interface{}) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// check company
	subdomain, err := u.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{
		"id": options["id"],
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if subdomain == nil {
		return response.Error(http.StatusBadRequest, "company not found")
	}

	err = u.mongodbRepo.DeleteCompany(ctx, subdomain)

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.Success("company deleted")
}

func (u *superadminUsecase) UploadCompanyLogo(ctx context.Context, claim domain.JWTClaimSuperadmin, payload domain.UploadAttachment, request *http.Request) response.Base {
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
	objName = "company-logos/" + category + "/" + strconv.Itoa(year) + "/" + strconv.Itoa(int(month)) + "/" + helpers.GenerateCleanName(uploadedFile.Filename)

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
		Category:     model.CompanyLogo,
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
