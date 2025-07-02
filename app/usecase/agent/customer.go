package usecase_agent

import (
	"app/domain"
	"app/domain/model"
	"app/helpers"
	"context"
	"net/http"
	"net/url"
	"time"

	yurekahelpers "github.com/Yureka-Teknologi-Cipta/yureka/helpers"
	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func (u *agentUsecase) GetCustomerList(ctx context.Context, claim domain.JWTClaimAgent, options map[string]interface{}) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	paramQuery := options["query"].(url.Values)
	page, limit, offset := yurekahelpers.GetLimitOffset(paramQuery)

	fetchOptions := map[string]interface{}{
		"limit":     limit,
		"offset":    offset,
		"companyID": claim.CompanyID,
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

	if paramQuery.Get("companyProductID") != "" {
		fetchOptions["companyProductID"] = paramQuery.Get("companyProductID")
	}

	// count first
	totalDocuments := u.mongodbRepo.CountCustomer(ctx, fetchOptions)

	if totalDocuments == 0 {
		return response.Success(response.List{
			List:  []interface{}{},
			Page:  page,
			Limit: limit,
			Total: totalDocuments,
		})
	}

	// check customer list
	cur, err := u.mongodbRepo.FetchCustomerList(ctx, fetchOptions)

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
		row := model.Customer{}
		formattedRow := helpers.CustomerBalanceFormat(&row)
		err := cur.Decode(&row)
		if err != nil {
			logrus.Error("Customer Decode ", err)
			return response.Success(response.List{
				List:  []interface{}{},
				Page:  page,
				Limit: limit,
				Total: totalDocuments,
			})
		}

		list = append(list, formattedRow)
	}

	return response.Success(response.List{
		List:  list,
		Page:  page,
		Limit: limit,
		Total: totalDocuments,
	})
}

func (u *agentUsecase) GetCustomerDetail(ctx context.Context, claim domain.JWTClaimAgent, options map[string]interface{}) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	customerID := options["id"]

	// check customer
	customer, err := u.mongodbRepo.FetchOneCustomer(ctx, map[string]interface{}{
		"id": customerID,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if customer == nil {
		return response.Error(http.StatusBadRequest, "customer not found")
	}

	return response.Success(helpers.CustomerBalanceFormat(customer))
}

func (u *agentUsecase) CreateCustomer(ctx context.Context, claim domain.JWTClaimAgent, options map[string]interface{}) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// payload
	payload := options["payload"].(domain.CreateUserRequest)

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

	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	// check customer
	customer, err := u.mongodbRepo.FetchOneCustomer(ctx, map[string]interface{}{
		"email": payload.Email,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if customer != nil {
		return response.Error(http.StatusBadRequest, "email already in use")
	}

	// get company
	company, err := u.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{
		"id": claim.CompanyID,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if company == nil {
		return response.Error(http.StatusBadRequest, "company not found")
	}

	// // check product
	// product, err := u.mongodbRepo.FetchOneCompanyProduct(ctx, map[string]interface{}{
	// 	"id":        payload.CompanyProductId,
	// 	"companyID": claim.CompanyID,
	// })

	// if err != nil {
	// 	return response.Error(http.StatusInternalServerError, err.Error())
	// }

	// if product == nil {
	// 	return response.Error(http.StatusBadRequest, "product not found")
	// }

	password := helpers.RandomString(5)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	t := time.Now()

	isNeedBalance := false
	subscription := &model.Subscription{}
	if claim.Company.Type == "B2C" {
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

	// create customer
	newCustomer := &model.Customer{
		ID:       primitive.NewObjectID(),
		Name:     payload.Name,
		Email:    payload.Email,
		Password: string(hashedPassword),
		// CompanyProduct: model.CompanyProductNested{
		// 	ID:    product.ID.Hex(),
		// 	Name:  product.Name,
		// 	Image: product.Logo.URL,
		// 	Code:  product.Code,
		// },
		Company:       claim.Company,
		IsNeedBalance: isNeedBalance,
		Subscription:  subscription,
		JobTitle:      "customer",
		Role:          "customer",
		IsVerified:    true,
		VerifiedAt:    &t,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err = u.mongodbRepo.CreateCustomer(ctx, newCustomer)

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	// get from config
	config := u._CacheConfig(ctx)

	go func() {
		_sendEmailCustomerCrendential(config, *newCustomer, password, *company)

		u.mongodbRepo.IncrementOneCompany(ctx, claim.Company.ID, map[string]int64{
			"customerTotal": 1,
		})
	}()

	return response.Success(newCustomer)
}

func (u *agentUsecase) UpdateCustomer(ctx context.Context, claim domain.JWTClaimAgent, options map[string]interface{}) response.Base {
	payload := options["payload"].(domain.CreateUserRequest)

	errValidation := make(map[string]string)

	// validating request
	if payload.Name == "" {
		errValidation["name"] = "name field is required"
	}
	// if payload.Email == "" {
	// 	errValidation["email"] = "email field is required"
	// } else if !helpers.IsValidEmail(payload.Email) {
	// 	errValidation["email"] = "email field is invalid"
	// }

	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	// check customer
	customer, err := u.mongodbRepo.FetchOneCustomer(ctx, map[string]interface{}{
		"id": options["id"],
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if customer == nil {
		return response.Error(http.StatusBadRequest, "customer not found")
	}

	// check product
	// product, err := u.mongodbRepo.FetchOneCompanyProduct(ctx, map[string]interface{}{
	// 	"id":        payload.CompanyProductId,
	// 	"companyID": claim.CompanyID,
	// })

	// if err != nil {
	// 	return response.Error(http.StatusInternalServerError, err.Error())
	// }

	// if product == nil {
	// 	return response.Error(http.StatusBadRequest, "product not found")
	// }

	// check customer
	existingCustomer, err := u.mongodbRepo.FetchOneCustomer(ctx, map[string]interface{}{
		"email": payload.Email,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if existingCustomer != nil && existingCustomer.ID.Hex() != customer.ID.Hex() {
		return response.Error(http.StatusBadRequest, "email already in use ")
	}

	// update customer
	customer.Name = payload.Name
	customer.Email = payload.Email
	// customer.CompanyProduct = model.CompanyProductNested{
	// 	ID:    product.ID.Hex(),
	// 	Name:  product.Name,
	// 	Image: product.Logo.URL,
	// 	Code:  product.Code,
	// }
	customer.UpdatedAt = time.Now()

	if err := u.mongodbRepo.UpdateOneCustomer(
		ctx,
		map[string]interface{}{"id": customer.ID},
		map[string]interface{}{
			"name":     customer.Name,
			"jobTitle": customer.JobTitle,
			// "company":   customer.CompanyProduct,
			"email":     customer.Email,
			"role":      customer.Role,
			"updatedAt": customer.UpdatedAt,
		}); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.Success(customer)
}

func (u *agentUsecase) DeleteCustomer(ctx context.Context, claim domain.JWTClaimAgent, options map[string]interface{}) response.Base {
	// check customer
	customer, err := u.mongodbRepo.FetchOneCustomer(ctx, map[string]interface{}{
		"id": options["id"],
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if customer == nil {
		return response.Error(http.StatusBadRequest, "customer not found")
	}

	t := time.Now()

	// update customer
	customer.UpdatedAt = t
	customer.DeletedAt = &t

	if err := u.mongodbRepo.UpdateOneCustomer(
		ctx,
		map[string]interface{}{"id": customer.ID},
		map[string]interface{}{
			"updatedAt": customer.UpdatedAt,
			"deletedAt": customer.DeletedAt,
		}); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.Success("Customer deleted successfully")
}
