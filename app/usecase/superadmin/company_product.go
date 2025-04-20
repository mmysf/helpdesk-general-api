package usecase_superadmin

import (
	"app/domain"
	"app/domain/model"
	"app/helpers"
	"context"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	yurekahelpers "github.com/Yureka-Teknologi-Cipta/yureka/helpers"
	"github.com/Yureka-Teknologi-Cipta/yureka/response"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func (u *superadminUsecase) GetCompanyProductList(ctx context.Context, claim domain.JWTClaimSuperadmin, options map[string]interface{}) response.Base {
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

	// count first
	totalDocuments := u.mongodbRepo.CountCompanyProduct(ctx, fetchOptions)

	if totalDocuments == 0 {
		return response.Success(response.List{
			List:  []interface{}{},
			Page:  page,
			Limit: limit,
			Total: totalDocuments,
		})
	}

	// check company list
	cur, err := u.mongodbRepo.FetchCompanyProductList(ctx, fetchOptions)

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
		row := model.CompanyProduct{}
		err := cur.Decode(&row)
		if err != nil {
			logrus.Error("Company Product Decode ", err)
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

func (u *superadminUsecase) GetCompanyProductDetail(ctx context.Context, claim domain.JWTClaimSuperadmin, options map[string]interface{}) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// get id
	companyID := options["id"]

	// check company
	company, err := u.mongodbRepo.FetchOneCompanyProduct(ctx, map[string]interface{}{
		"id": companyID,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if company == nil {
		return response.Error(http.StatusBadRequest, "CompanyProduct not found")
	}

	return response.Success(company)
}

func (u *superadminUsecase) CreateCompanyProduct(ctx context.Context, claim domain.JWTClaimSuperadmin, options map[string]interface{}) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// payload
	payload := options["payload"].(domain.CreateCompanyProductRequest)

	errValidation := make(map[string]string)

	// validating request
	if payload.Name == "" {
		errValidation["name"] = "name field is required"
	}

	if payload.Code == "" {
		errValidation["code"] = "code field is required"
	} else {
		if payload.Code == payload.Name {
			errValidation["code"] = "code cannot be same with name"
		}
		// trim space
		payload.Code = strings.TrimSpace(payload.Code)
		payload.Code = strings.ReplaceAll(payload.Code, " ", "")
		if !helpers.IsValidAlphanumeric(payload.Code) {
			errValidation["code"] = "code only accept aplhanumeric"
		}
	}

	if payload.Email == "" {
		errValidation["email"] = "email field is required"
	} else if !helpers.IsValidEmail(payload.Email) {
		errValidation["email"] = "email field is invalid"
	}

	if payload.LogoAttachId == "" {
		errValidation["logoAttachId"] = "logoAttachId field is required"
	}

	if payload.CompanyId == "" {
		errValidation["companyId"] = "companyId field is required"
	}

	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	// check exist product by code payload
	existProduct, err := u.mongodbRepo.FetchOneCompanyProduct(ctx, map[string]interface{}{
		"code": strings.ToUpper(payload.Code),
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if existProduct != nil {
		return response.Error(http.StatusBadRequest, "code already used by another company product")
	}

	// check company
	company, err := u.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{
		"id": payload.CompanyId,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if company == nil {
		return response.Error(http.StatusBadRequest, "company not found")
	}

	// check logo
	logo, err := u.mongodbRepo.FetchOneMedia(ctx, map[string]interface{}{
		"id":       payload.LogoAttachId,
		"category": model.BrandLogo,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if logo == nil {
		return response.Error(http.StatusBadRequest, "Brand logo media not found")
	}

	// check customer
	customer, err := u.mongodbRepo.FetchOneCustomer(ctx, map[string]interface{}{
		"email": payload.Email,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if customer != nil {
		return response.Error(http.StatusBadRequest, "email already in use for brand "+customer.CompanyProduct.Name)
	}

	// create product
	product := model.CompanyProduct{
		ID: primitive.NewObjectID(),
		Company: model.CompanyNested{
			ID:    company.ID.Hex(),
			Name:  company.Name,
			Image: company.Logo.URL,
			Type:  company.Type,
		},
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
		Name:      payload.Name,
		Code:      strings.ToUpper(payload.Code),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = u.mongodbRepo.CreateCompanyProduct(ctx, &product)

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	password := helpers.RandomString(5)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	t := time.Now()

	isNeedBalance := false
	if company.Type == "B2C" {
		isNeedBalance = true
	}

	// create customer
	newUser := &model.Customer{
		ID:       primitive.NewObjectID(),
		Name:     payload.Name,
		Email:    payload.Email,
		Password: string(hashedPassword),
		CompanyProduct: model.CompanyProductNested{
			ID:    product.ID.Hex(),
			Name:  product.Name,
			Image: product.Logo.URL,
			Code:  product.Code,
		},
		Company: model.CompanyNested{
			ID:    company.ID.Hex(),
			Name:  company.Name,
			Image: company.Logo.URL,
			Type:  company.Type,
		},
		IsNeedBalance: isNeedBalance,
		Subscription:  nil,
		JobTitle:      "Admin",
		Role:          model.AdminRole,
		IsVerified:    true,
		VerifiedAt:    &t,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err = u.mongodbRepo.CreateCustomer(ctx, newUser)

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	// get from config
	config := u._CacheConfig(ctx)

	go func() {
		_sendEmailCustomerCrendential(config, *newUser, password, *company)

		u.mongodbRepo.UpdateManyMediaPartial(ctx, []primitive.ObjectID{logo.ID}, map[string]interface{}{
			"isUsed": true,
		})

		u.mongodbRepo.IncrementOneCompany(ctx, company.ID.Hex(), map[string]int64{
			"productTotal": 1,
		})
	}()

	return response.Success(product)
}

func _sendEmailCustomerCrendential(config model.Config, customer model.Customer, password string, company model.Company) {
	loginLink := helpers.StringReplacer(config.LoginLink, map[string]string{
		"base_url_frontend": company.Settings.Domain.FullUrl,
	})
	// send email
	mail := helpers.NewSMTPMailer(&company)
	mail.To([]string{customer.Email})
	mail.Subject(config.Email.Template.DefaultUser.Title)
	mail.Body(helpers.StringReplacer(config.Email.Template.DefaultUser.Body, map[string]string{
		"title":      config.Email.Template.DefaultUser.Title,
		"email":      customer.Email,
		"password":   password,
		"login_link": loginLink,
	}))

	// send
	if err := mail.Send(); err != nil {
		logrus.Errorf("Send Email to %s error %v", customer.Email, err)
	}
}

func (u *superadminUsecase) UpdateCompanyProduct(ctx context.Context, claim domain.JWTClaimSuperadmin, options map[string]interface{}) response.Base {
	// payload
	payload := options["payload"].(domain.CreateCompanyProductRequest)

	errValidation := make(map[string]string)

	// validating request
	if payload.Name == "" {
		errValidation["name"] = "name field is required"
	}

	if payload.Code == "" {
		errValidation["code"] = "code field is required"
	} else {
		if payload.Code == payload.Name {
			errValidation["code"] = "code cannot be same with name"
		}
		// trim space
		payload.Code = strings.TrimSpace(payload.Code)
		payload.Code = strings.ReplaceAll(payload.Code, " ", "")
		if !helpers.IsValidAlphanumeric(payload.Code) {
			errValidation["code"] = "code only accept aplhanumeric"
		}
	}

	if payload.LogoAttachId == "" {
		errValidation["logoAttachId"] = "logoAttachId field is required"
	}

	if payload.CompanyId == "" {
		errValidation["companyId"] = "companyId field is required"
	}

	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	// check company
	company, err := u.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{
		"id": payload.CompanyId,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if company == nil {
		return response.Error(http.StatusBadRequest, "company not found")
	}

	// check company product
	companyProduct, err := u.mongodbRepo.FetchOneCompanyProduct(ctx, map[string]interface{}{
		"id": options["id"],
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if companyProduct == nil {
		return response.Error(http.StatusBadRequest, "company product not found")
	}

	// check the product company
	company, err = u.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{
		"id": companyProduct.Company.ID,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if company == nil {
		return response.Error(http.StatusBadRequest, "product company not found")
	}

	if company.Type == "B2C" {
		return response.Error(http.StatusBadRequest, "company type must be B2B")
	}

	// check exist product by code payload
	existProduct, err := u.mongodbRepo.FetchOneCompanyProduct(ctx, map[string]interface{}{
		"code": strings.ToUpper(payload.Code),
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if existProduct != nil && existProduct.ID.Hex() != companyProduct.ID.Hex() {
		return response.Error(http.StatusBadRequest, "code already used by another company product")
	}

	// check logo
	logo, err := u.mongodbRepo.FetchOneMedia(ctx, map[string]interface{}{
		"id":       payload.LogoAttachId,
		"category": model.BrandLogo,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if logo == nil {
		return response.Error(http.StatusBadRequest, "Brand logo media not found")
	}

	// update product
	companyProduct.Logo = model.MediaFK{
		ID:          logo.ID.Hex(),
		Name:        logo.Name,
		Size:        logo.Size,
		URL:         logo.URL,
		Type:        logo.Type,
		Category:    logo.Category,
		ProviderKey: logo.ProviderKey,
		IsPrivate:   logo.IsPrivate,
	}
	companyProduct.Name = payload.Name
	companyProduct.Code = strings.ToUpper(payload.Code)
	companyProduct.UpdatedAt = time.Now()

	err = u.mongodbRepo.UpdateCompanyProduct(ctx, companyProduct)

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	// update company product nested
	go u.mongodbRepo.UpdatePartialCustomer(ctx, map[string]any{
		"companyProductID": companyProduct.ID.Hex(),
	}, map[string]any{
		"companyProduct.name":  payload.Name,
		"companyProduct.image": companyProduct.Logo.URL,
		"companyProduct.code":  strings.ToUpper(payload.Code),
	})

	return response.Success(companyProduct)
}

func (u *superadminUsecase) DeleteCompanyProduct(ctx context.Context, claim domain.JWTClaimSuperadmin, options map[string]interface{}) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// check company
	companyProduct, err := u.mongodbRepo.FetchOneCompanyProduct(ctx, map[string]interface{}{
		"id": options["id"],
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if companyProduct == nil {
		return response.Error(http.StatusBadRequest, "company not found")
	}

	err = u.mongodbRepo.DeleteCompanyProduct(ctx, companyProduct)

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	u.mongodbRepo.IncrementOneCompany(ctx, companyProduct.Company.ID, map[string]int64{
		"productTotal": -1,
	})

	return response.Success("company product deleted")
}

func (u *superadminUsecase) UploadCompanyProductLogo(ctx context.Context, claim domain.JWTClaimSuperadmin, payload domain.UploadAttachment, request *http.Request) response.Base {
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
	objName = "brand-logos/" + category + "/" + strconv.Itoa(year) + "/" + strconv.Itoa(int(month)) + "/" + helpers.GenerateCleanName(uploadedFile.Filename)

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
		Category:     model.BrandLogo,
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
