package usecase_member

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

func (u *appUsecase) GetUserList(ctx context.Context, claim domain.JWTClaimUser, options map[string]interface{}) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	paramQuery := options["query"].(url.Values)
	page, limit, offset := yurekahelpers.GetLimitOffset(paramQuery)

	fetchOptions := map[string]interface{}{
		"limit":     limit,
		"offset":    offset,
		"companyID": claim.CompanyID,
		// "companyProductID": claim.CompanyProductID,
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
			logrus.Error("User Decode ", err)
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

func (u *appUsecase) GetUserDetail(ctx context.Context, claim domain.JWTClaimUser, options map[string]interface{}) response.Base {
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
		return response.Error(http.StatusUnauthorized, "customer not found")
	}

	return response.Success(helpers.CustomerBalanceFormat(customer))
}

func (u *appUsecase) CreateUser(ctx context.Context, claim domain.JWTClaimUser, options map[string]interface{}) response.Base {
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

	if payload.JobTitle == "" {
		errValidation["jobTitle"] = "jobTitle field is required"
	}

	if payload.Role == "" {
		errValidation["role"] = "role field is required"
	} else if !helpers.InArrayString(payload.Role, []string{"admin", "customer"}) {
		errValidation["role"] = "role field must be admin or customer"
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
		"id": claim.Company.ID,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if company == nil {
		return response.Error(http.StatusUnauthorized, "company not found")
	}

	password := helpers.RandomString(5)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	t := time.Now()

	isNeedBalance := false
	if claim.Company.Type == "B2C" {
		isNeedBalance = true
	}

	// create customer
	newUser := &model.Customer{
		ID:       primitive.NewObjectID(),
		Name:     payload.Name,
		Email:    payload.Email,
		Password: string(hashedPassword),
		// CompanyProduct: claim.CompanyProduct,
		Company:       claim.Company,
		IsNeedBalance: isNeedBalance,
		Subscription:  nil,
		JobTitle:      payload.JobTitle,
		Role:          model.UserRole(payload.Role),
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
	}()

	return response.Success(newUser)
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

func (u *appUsecase) UpdateUser(ctx context.Context, claim domain.JWTClaimUser, options map[string]interface{}) response.Base {
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
	if payload.JobTitle == "" {
		errValidation["jobTitle"] = "jobTitle field is required"
	}
	if payload.Role == "" {
		errValidation["role"] = "role field is required"
	} else if !helpers.InArrayString(payload.Role, []string{"admin", "customer"}) {
		errValidation["role"] = "role field must be admin or customer"
	}

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
		return response.Error(http.StatusUnauthorized, "customer not found")
	}

	// check customer
	existingCustomer, err := u.mongodbRepo.FetchOneCustomer(ctx, map[string]interface{}{
		"email": payload.Email,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if existingCustomer != nil && existingCustomer.Email != customer.Email {
		return response.Error(http.StatusBadRequest, "email already in use")
	}

	// update customer
	customer.Name = payload.Name
	customer.Email = payload.Email
	customer.JobTitle = payload.JobTitle
	customer.Role = model.UserRole(payload.Role)
	customer.UpdatedAt = time.Now()

	if err := u.mongodbRepo.UpdateOneCustomer(
		ctx,
		map[string]interface{}{"id": customer.ID},
		map[string]interface{}{
			"name":      customer.Name,
			"jobTitle":  customer.JobTitle,
			"email":     customer.Email,
			"role":      customer.Role,
			"updatedAt": customer.UpdatedAt,
		}); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.Success(customer)
}

func (u *appUsecase) DeleteUser(ctx context.Context, claim domain.JWTClaimUser, options map[string]interface{}) response.Base {
	// check customer
	customer, err := u.mongodbRepo.FetchOneCustomer(ctx, map[string]interface{}{
		"id": options["id"],
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if customer == nil {
		return response.Error(http.StatusUnauthorized, "customer not found")
	}

	// validate cant delete self
	if customer.ID.Hex() == claim.UserID {
		return response.Error(http.StatusBadRequest, "cannot delete your own account")
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

	return response.Success("User deleted successfully")
}
