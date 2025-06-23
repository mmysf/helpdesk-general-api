package usecase_superadmin

import (
	"app/domain"
	"app/domain/model"
	"app/helpers"
	"bufio"
	"context"
	"encoding/csv"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	yurekahelpers "github.com/Yureka-Teknologi-Cipta/yureka/helpers"
	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func (u *superadminUsecase) GetCustomers(ctx context.Context, claim domain.JWTClaimSuperadmin, query url.Values) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	page, limit, offset := yurekahelpers.GetLimitOffset(query)

	fetchOptions := map[string]interface{}{
		"limit":  limit,
		"offset": offset,
	}

	// filtering
	if query.Get("sort") != "" {
		fetchOptions["sort"] = query.Get("sort")
	}

	if query.Get("dir") != "" {
		fetchOptions["dir"] = query.Get("dir")
	}

	if query.Get("q") != "" {
		fetchOptions["q"] = query.Get("q")
	}

	if query.Get("companyID") != "" {
		fetchOptions["companyID"] = query.Get("companyID")
	}

	if query.Get("companyProductID") != "" {
		fetchOptions["companyProductID"] = query.Get("companyProductID")
	}

	if query.Get("type") != "" {
		fetchOptions["type"] = query.Get("type")
	}

	// count first
	totalCustomers := u.mongodbRepo.CountCustomer(ctx, fetchOptions)
	if totalCustomers == 0 {
		return response.Success(domain.ResponseList{
			List: response.List{
				List:  []interface{}{},
				Page:  page,
				Limit: limit,
				Total: totalCustomers,
			},
			TotalPage: helpers.GetTotalPage(totalCustomers, limit),
		})
	}

	// check customer list
	customers, err := u.mongodbRepo.FetchCustomerList(ctx, fetchOptions)
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	defer customers.Close(ctx)

	list := make([]interface{}, 0)
	for customers.Next(ctx) {
		row := model.Customer{}
		formattedRow := helpers.CustomerBalanceFormat(&row)
		err := customers.Decode(&row)
		if err != nil {
			logrus.Error("Topup Decode ", err)
			return response.Error(http.StatusInternalServerError, err.Error())
		}
		list = append(list, formattedRow)
	}

	return response.Success(domain.ResponseList{
		List: response.List{
			List:  list,
			Page:  page,
			Limit: limit,
			Total: totalCustomers,
		},
		TotalPage: helpers.GetTotalPage(totalCustomers, limit),
	})
}

func (u *superadminUsecase) CreateCustomer(ctx context.Context, claim domain.JWTClaimSuperadmin, payload domain.AccountRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

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
	if payload.CompanyID == "" {
		errValidation["companyID"] = "companyID field is required"
	}
	if payload.CompanyProductID == "" {
		errValidation["companyProductID"] = "companyProductID field is required"
	}
	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(helpers.RandomString(5)), bcrypt.DefaultCost)

	// check the user
	existingUser, err := u.mongodbRepo.FetchOneCustomer(ctx, map[string]interface{}{
		"email": payload.Email,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if existingUser != nil {
		return response.Error(http.StatusBadRequest, "email already taken")
	}

	// check company
	company, err := u.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{
		"id": payload.CompanyID,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if company == nil {
		return response.Error(http.StatusBadRequest, "company not found")
	}

	// check company product
	companyProduct, err := u.mongodbRepo.FetchOneCompanyProduct(ctx, map[string]interface{}{
		"id": payload.CompanyProductID,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if companyProduct == nil {
		return response.Error(http.StatusBadRequest, "company product not found")
	}

	// get from config
	config := u._CacheConfig(ctx)

	newCompanyNested := model.CompanyNested{
		ID:    company.ID.Hex(),
		Name:  company.Name,
		Image: company.Logo.URL,
		Type:  company.Type,
	}

	newCompanyProductNested := model.CompanyProductNested{
		ID:    companyProduct.ID.Hex(),
		Name:  companyProduct.Name,
		Image: companyProduct.Logo.URL,
		Code:  companyProduct.Code,
	}

	isNeedBalance := false
	subscription := &model.Subscription{}
	if company.Type == "B2C" {
		isNeedBalance = true
		now := time.Now()

		// subscription expiry
		startAt := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		expiredAt := startAt.AddDate(0, helpers.GetSubscriptionDuration(), 0)

		// set starter subscription
		subscription = &model.Subscription{
			Status:      model.Active,
			HourPackage: nil,
			Balance:     nil,
			StartAt:     startAt,
			EndAt:       expiredAt,
		}
	}

	newUser := model.Customer{
		ID:             primitive.NewObjectID(),
		Company:        newCompanyNested,
		CompanyProduct: newCompanyProductNested,
		Name:           payload.Name,
		Email:          payload.Email,
		Password:       string(hashedPassword),
		Role:           model.AdminRole,
		IsVerified:     true,
		IsNeedBalance:  isNeedBalance,
		Subscription:   subscription,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = u.mongodbRepo.CreateCustomer(ctx, &newUser)
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	go _sendEmailCustomerCrendential(config, newUser, string(hashedPassword), *company)

	return response.Success(newUser)

}

func (u *superadminUsecase) GetCustomerDetail(ctx context.Context, claim domain.JWTClaimSuperadmin, customerId string) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	customer, err := u.mongodbRepo.FetchOneCustomer(ctx, map[string]interface{}{
		"id": customerId,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if customer == nil {
		return response.Error(http.StatusBadRequest, "customer not found")
	}

	return response.Success(helpers.CustomerBalanceFormat(customer))
}

func (u *superadminUsecase) UpdateCustomer(ctx context.Context, claim domain.JWTClaimSuperadmin, customerId string, payload domain.AccountRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// check customer
	customer, err := u.mongodbRepo.FetchOneCustomer(ctx, map[string]interface{}{
		"id": customerId,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if customer == nil {
		return response.Error(http.StatusBadRequest, "customer not found")
	}

	customerEmail, err := u.mongodbRepo.FetchOneCustomer(ctx, map[string]interface{}{
		"email": payload.Email,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if customerEmail != nil && customerEmail.ID.Hex() != customer.ID.Hex() {
		return response.Error(http.StatusBadRequest, "email already taken")
	}

	if payload.Name != "" {
		customer.Name = payload.Name
	}
	if payload.Email != "" {
		if !helpers.IsValidEmail(payload.Email) {
			return response.Error(http.StatusBadRequest, "email field is invalid")
		}
		customer.Email = payload.Email
	}

	customer.UpdatedAt = time.Now()

	if err := u.mongodbRepo.UpdateOneCustomer(
		ctx,
		map[string]interface{}{"id": customer.ID},
		map[string]interface{}{
			"name":      customer.Name,
			"email":     customer.Email,
			"updatedAt": customer.UpdatedAt,
		}); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.Success(customer)
}

func (u *superadminUsecase) DeleteCustomer(ctx context.Context, claim domain.JWTClaimSuperadmin, customerId string) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	customer, err := u.mongodbRepo.FetchOneCustomer(ctx, map[string]interface{}{
		"id": customerId,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if customer == nil {
		return response.Error(http.StatusBadRequest, "customer not found")
	}

	now := time.Now()
	customer.UpdatedAt = now
	customer.DeletedAt = &now

	if err := u.mongodbRepo.UpdateOneCustomer(
		ctx,
		map[string]interface{}{"id": customer.ID},
		map[string]interface{}{
			"updatedAt": customer.UpdatedAt,
			"deletedAt": customer.DeletedAt,
		}); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.Success(customer)
}

func (u *superadminUsecase) ResetPasswordCustomer(ctx context.Context, claim domain.JWTClaimSuperadmin, customerId string) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	customer, err := u.mongodbRepo.FetchOneCustomer(ctx, map[string]interface{}{
		"id": customerId,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if customer == nil {
		return response.Error(http.StatusBadRequest, "customer not found")
	}

	defaultPassword := helpers.RandomString(5)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)

	now := time.Now()
	customer.UpdatedAt = now
	customer.Password = string(hashedPassword)

	if err := u.mongodbRepo.UpdateOneCustomer(
		ctx,
		map[string]interface{}{"id": customer.ID},
		map[string]interface{}{
			"updatedAt": customer.UpdatedAt,
			"password":  customer.Password,
		}); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.Success(customer)
}

func (u *superadminUsecase) ImportCustomer(ctx context.Context, claim domain.JWTClaimSuperadmin, request *http.Request) response.Base {
	_, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	var err error
	validation := make(map[string]string)
	var file multipart.File
	var uploadedFile *multipart.FileHeader

	file, uploadedFile, err = request.FormFile("file")

	if err != nil {
		validation["file"] = "file field is required"
		// return response.Error(http.StatusInternalServerError, err.Error())
	}
	if file == nil || uploadedFile == nil {
		validation["file"] = "file field is required"
	} else {
		typeDocument := uploadedFile.Header.Get("Content-Type")
		if !helpers.InArrayString(typeDocument, []string{"text/csv"}) {
			validation["file"] = "field file is not valid type"
		}
	}

	if len(validation) > 0 {
		return response.ErrorValidation(validation, "error validation")
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))

	companyMap := map[string]model.Company{}
	companyProductMap := map[string]model.CompanyProduct{}
	customers := []*model.Customer{}
	count := 0

	for {
		// Read one row at a time
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" { // End of file
				break
			}
		}

		if record[0] == "No" {
			continue
		}

		// check the user
		existingUser, err := u.mongodbRepo.FetchOneCustomer(ctx, map[string]interface{}{
			"email": record[2],
		})
		if err != nil {
			return response.Error(http.StatusInternalServerError, err.Error())
		}
		if existingUser != nil {
			return response.Error(http.StatusBadRequest, existingUser.Email+" email already taken")
		}

		if _, exists := companyMap[record[5]]; !exists {
			// check company
			company, err := u.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{
				"id": record[5],
			})
			if err != nil {
				return response.Error(http.StatusInternalServerError, err.Error())
			}
			if company == nil {
				return response.Error(http.StatusBadRequest, "company not found")
			}
			companyMap[record[5]] = *company
		}

		if _, exists := companyProductMap[record[6]]; !exists {
			// check company product
			companyProduct, err := u.mongodbRepo.FetchOneCompanyProduct(ctx, map[string]interface{}{
				"id": record[6],
			})
			if err != nil {
				return response.Error(http.StatusInternalServerError, err.Error())
			}
			if companyProduct == nil {
				return response.Error(http.StatusBadRequest, "company product not found")
			}
			companyProductMap[record[6]] = *companyProduct
		}

		isNeedBalance := false
		if companyMap[record[5]].Type == "B2C" {
			isNeedBalance = true
		}

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(record[4]), bcrypt.DefaultCost)
		t := time.Now()

		customers = append(customers, &model.Customer{
			ID: primitive.NewObjectID(),
			Company: model.CompanyNested{
				ID:    companyMap[record[5]].ID.Hex(),
				Name:  companyMap[record[5]].Name,
				Image: companyMap[record[5]].Logo.URL,
				Type:  companyMap[record[5]].Type,
			},
			CompanyProduct: model.CompanyProductNested{
				ID:    companyProductMap[record[6]].ID.Hex(),
				Name:  companyProductMap[record[6]].Name,
				Image: companyProductMap[record[6]].Logo.URL,
				Code:  companyProductMap[record[6]].Code,
			},
			Name:          record[1],
			Email:         record[2],
			JobTitle:      record[3],
			IsVerified:    true,
			Password:      string(hashedPassword),
			IsNeedBalance: isNeedBalance,
			Role:          model.AdminRole,
			CreatedAt:     t,
			UpdatedAt:     t,
			VerifiedAt:    &t,
		})

		count++
		if count%10 == 0 {
			err = u.mongodbRepo.CreateManyCustomer(context.Background(), customers)
			if err != nil {
				return response.Error(http.StatusInternalServerError, err.Error())
			}
			helpers.Dump("ini insert 1000")
			customers = make([]*model.Customer, 0)
			break
		}
	}

	if len(customers) > 0 {
		err = u.mongodbRepo.CreateManyCustomer(context.Background(), customers)
		if err != nil {
			return response.Error(http.StatusInternalServerError, err.Error())
		}
		helpers.Dump("ini insert di akhir")
	}

	return response.Success(nil)
}
