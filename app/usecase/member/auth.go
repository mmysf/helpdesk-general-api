package usecase_member

import (
	"app/domain"
	"app/domain/model"
	"app/helpers"
	"context"
	"net/http"
	"time"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func (u *appUsecase) Login(ctx context.Context, payload domain.LoginRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	errValidation := make(map[string]string)
	// validating request
	if payload.Email == "" {
		errValidation["email"] = "email field is required"
	}

	if payload.Password == "" {
		errValidation["password"] = "password field is required"
	}

	if payload.AccessKey == "" {
		errValidation["accessKey"] = "accessKey field is required"
	}

	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	// check company
	company, err := u.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{
		"accessKey": payload.AccessKey,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if company == nil {
		return response.Error(http.StatusBadRequest, "company not found")
	}

	// check the db
	user, err := u.mongodbRepo.FetchOneCustomer(ctx, map[string]interface{}{
		"email":     payload.Email,
		"companyID": company.ID.Hex(),
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if user == nil {
		return response.Error(http.StatusBadRequest, "user not found")
	}

	// check companyProduct
	// companyProduct, err := u.mongodbRepo.FetchOneCompanyProduct(ctx, map[string]interface{}{
	// 	"id": user.CompanyProduct.ID,
	// })
	// if err != nil {
	// 	return response.Error(http.StatusInternalServerError, err.Error())
	// }

	// if companyProduct == nil {
	// 	return response.Error(http.StatusBadRequest, "companyProduct not found")
	// }

	// check password
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		return response.Error(http.StatusBadRequest, "Wrong password")
	}

	if !user.IsVerified {
		return response.Error(http.StatusBadRequest, "user email not verified")
	}

	// generate token
	tokenString, err := helpers.GenerateJWTTokenCustomer(domain.JWTClaimUser{

		UserID:    user.ID.Hex(),
		CompanyID: user.Company.ID,
		// CompanyProductID: user.CompanyProduct.ID,
		Role: string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Issuer:    "member",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(helpers.GetJWTTTL()) * time.Minute)),
		},
	})
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error())
	}

	return response.Success(map[string]interface{}{
		"user":  user,
		"token": tokenString,
	})
}

func (u *appUsecase) Register(ctx context.Context, payload domain.RegisterRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	errValidation := make(map[string]string)
	// validating request
	if payload.Name == "" {
		errValidation["name"] = "name field is required"
	}

	if payload.Email == "" {
		errValidation["email"] = "email field is required"
	}

	if payload.Password == "" {
		errValidation["password"] = "password field is required"
	}

	if payload.AccessKey == "" {
		errValidation["accessKey"] = "accessKey field is required"

	}

	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	// check company
	company, err := u.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{
		"accessKey": payload.AccessKey,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if company == nil {
		return response.Error(http.StatusBadRequest, "company not found")
	}

	if company.Type != "B2C" {
		return response.Error(http.StatusBadRequest, "company type must be B2C")
	}

	// check company product
	// companyProduct, err := u.mongodbRepo.FetchOneCompanyProduct(ctx, map[string]interface{}{
	// 	"id": os.Getenv("B2C_COMPANY_PRODUCT_ID"),
	// })

	// if err != nil {
	// 	return response.Error(http.StatusInternalServerError, err.Error())
	// }

	// if companyProduct == nil {
	// 	return response.Error(http.StatusBadRequest, "company product not found")
	// }

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

	newCompanyNested := model.CompanyNested{
		ID:    company.ID.Hex(),
		Name:  company.Name,
		Image: company.Logo.URL,
		Type:  company.Type,
	}

	// newCompanyProduct := model.CompanyProductNested{
	// 	ID:    companyProduct.ID.Hex(),
	// 	Name:  companyProduct.Name,
	// 	Image: companyProduct.Logo.URL,
	// 	Code:  companyProduct.Code,
	// }

	// check the user
	existingUser, err := u.mongodbRepo.FetchOneCustomer(ctx, map[string]interface{}{
		"email":     payload.Email,
		"companyID": company.ID.Hex(),
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	// get from config
	config := u._CacheConfig(ctx)

	if existingUser != nil {
		// resend verify email
		if config.Registration.VerifyEmail && !existingUser.IsVerified {

			// new token
			existingUser.Token = helpers.RandomString(64)
			existingUser.UpdatedAt = time.Now()

			// save
			u.mongodbRepo.UpdateOneCustomer(context.Background(), map[string]interface{}{
				"id": existingUser.ID,
			}, map[string]interface{}{
				"token":     existingUser.Token,
				"updatedAt": existingUser.UpdatedAt,
			})

			go _sendEmailRegistration(config, *existingUser, company)

			return response.Success(existingUser)
		}

		return response.Error(http.StatusBadRequest, "email already taken")
	}

	defaultVerified := true
	defaultToken := ""

	// need verify email
	if config.Registration.VerifyEmail {
		defaultVerified = false
		defaultToken = helpers.RandomString(64)
	}

	now := time.Now()

	// subscription expiry
	startAt := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	expiredAt := startAt.AddDate(0, helpers.GetSubscriptionDuration(), 0)

	// set starter subscription
	subscription := model.Subscription{
		Status:      model.Active,
		HourPackage: nil,
		Balance:     nil,
		StartAt:     startAt,
		EndAt:       expiredAt,
	}

	newUser := model.Customer{
		ID:      primitive.NewObjectID(),
		Name:    payload.Name,
		Email:   payload.Email,
		Company: newCompanyNested,
		// CompanyProduct: newCompanyProduct,
		Password:      string(hashedPassword),
		Role:          model.AdminRole,
		IsNeedBalance: true,
		Subscription:  &subscription,
		Token:         defaultToken,
		IsVerified:    defaultVerified,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	err = u.mongodbRepo.CreateCustomer(ctx, &newUser)

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	// need verify email
	if config.Registration.VerifyEmail {
		go _sendEmailRegistration(config, newUser, company)
	}

	return response.Success(newUser)
}

func (u *appUsecase) VerifyRegistration(ctx context.Context, payload domain.VerifyRegisterRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	errValidation := make(map[string]string)
	// validating request
	if payload.Token == "" {
		errValidation["token"] = "token field is required"
	}

	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	// check the db
	existingUser, err := u.mongodbRepo.FetchOneCustomer(ctx, map[string]interface{}{
		"token": payload.Token,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if existingUser == nil {
		return response.Error(http.StatusBadRequest, "verification token not valid")
	}

	// update verified
	existingUser.IsVerified = true
	existingUser.Token = ""
	existingUser.UpdatedAt = time.Now()
	existingUser.VerifiedAt = &existingUser.UpdatedAt

	// save
	u.mongodbRepo.UpdateOneCustomer(context.Background(), map[string]interface{}{
		"id": existingUser.ID,
	}, map[string]interface{}{
		"token":      existingUser.Token,
		"isVerified": existingUser.IsVerified,
		"updatedAt":  existingUser.UpdatedAt,
		"verifiedAt": existingUser.VerifiedAt,
	})

	return response.Success(existingUser)
}

func _sendEmailRegistration(config model.Config, user model.Customer, company *model.Company) {
	verificationLink := helpers.StringReplacer(config.Registration.VerificationLink, map[string]string{
		"base_url_frontend": company.Settings.Domain.FullUrl,
		"token":             user.Token,
	})
	// send email
	mail := helpers.NewSMTPMailer(company)
	mail.To([]string{user.Email})
	mail.Subject(config.Email.Template.Register.Title)
	mail.Body(helpers.StringReplacer(config.Email.Template.Register.Body, map[string]string{
		"title":             config.Email.Template.Register.Title,
		"verification_link": verificationLink,
	}))

	// send
	if err := mail.Send(); err != nil {
		logrus.Errorf("Send Email to %s error %v", user.Email, err)
	}
}

func (u *appUsecase) SendEmailPasswordReset(ctx context.Context, payload domain.EmailPasswordResetRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	errValidation := make(map[string]string)
	// validating request
	if payload.Email == "" {
		errValidation["email"] = "email field is required"
	}

	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	// check the db
	user, err := u.mongodbRepo.FetchOneCustomer(ctx, map[string]interface{}{
		"email": payload.Email,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if user == nil {
		return response.Error(http.StatusBadRequest, "user not found")
	}

	// get company
	company, err := u.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{
		"id": user.Company.ID,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if company == nil {
		return response.Error(http.StatusBadRequest, "company not found")
	}

	// get from config
	config := u._CacheConfig(ctx)

	// new token
	user.PasswordResetToken = helpers.RandomString(64)
	user.UpdatedAt = time.Now()

	// save
	u.mongodbRepo.UpdateOneCustomer(context.Background(), map[string]interface{}{
		"id": user.ID,
	}, map[string]interface{}{
		"passwordResetToken": user.PasswordResetToken,
		"updatedAt":          user.UpdatedAt,
	})

	go _sendEmailPasswordReset(config, *user, *company)

	return response.Success("The password reset email has been sent.")

}

func _sendEmailPasswordReset(config model.Config, user model.Customer, company model.Company) {
	passwordResetLink := helpers.StringReplacer(config.ResetPasswordLink, map[string]string{
		"base_url_frontend":  company.Settings.Domain.FullUrl,
		"passwordResetToken": user.PasswordResetToken,
	})
	// send email
	mail := helpers.NewSMTPMailer(&company)
	mail.To([]string{user.Email})
	mail.Subject(config.Email.Template.ResetPassword.Title)
	mail.Body(helpers.StringReplacer(config.Email.Template.ResetPassword.Body, map[string]string{
		"title":               config.Email.Template.ResetPassword.Title,
		"reset_password_link": passwordResetLink,
	}))

	// send
	if err := mail.Send(); err != nil {
		logrus.Errorf("Send Email to %s error %v", user.Email, err)
	}
}

func (u *appUsecase) PasswordReset(ctx context.Context, payload domain.PasswordResetRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	errValidation := make(map[string]string)
	// validating request
	if payload.Token == "" {
		errValidation["Token"] = "token field is required"
	}
	if payload.Password == "" {
		errValidation["password"] = "password field is required"
	}
	if len(payload.Password) < 8 {
		errValidation["password"] = "password must be at least 8 characters"
	}
	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	// check the db
	user, err := u.mongodbRepo.FetchOneCustomer(ctx, map[string]interface{}{
		"passwordResetToken": payload.Token,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if user == nil {
		return response.Error(http.StatusBadRequest, "password reset token not valid")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

	// update user
	user.PasswordResetToken = ""
	user.Password = string(hashedPassword)
	user.UpdatedAt = time.Now()

	// save
	u.mongodbRepo.UpdateOneCustomer(context.Background(), map[string]interface{}{
		"id": user.ID,
	}, map[string]interface{}{
		"passwordResetToken": user.PasswordResetToken,
		"password":           user.Password,
		"updatedAt":          user.UpdatedAt,
	})

	return response.Success(user)
}

func (u *appUsecase) GetMe(ctx context.Context, claim domain.JWTClaimUser) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	userID := claim.UserID

	// check the db
	user, err := u.mongodbRepo.FetchOneCustomer(ctx, map[string]interface{}{
		"id": userID,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if user == nil {
		return response.Error(http.StatusBadRequest, "user not found")
	}

	return response.Success(helpers.CustomerBalanceFormat(user))
}

func (u *appUsecase) RegisterB2B(ctx context.Context, payload domain.RegisterRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	errValidation := make(map[string]string)
	// validating request
	if payload.Name == "" {
		errValidation["name"] = "name field is required"
	}

	if payload.Email == "" {
		errValidation["email"] = "email field is required"
	}

	if payload.Password == "" {
		errValidation["password"] = "password field is required"
	}

	// if payload.CompanyProductId == "" {
	// 	errValidation["companyProductId"] = "companyProductId field is required"
	// }

	if payload.AccessKey == "" {
		errValidation["accessKey"] = "accessKey field is required"
	}

	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	// check company
	company, err := u.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{
		"accessKey": payload.AccessKey,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if company == nil {
		return response.Error(http.StatusBadRequest, "company not found")
	}

	//check company product
	// companyProduct, err := u.mongodbRepo.FetchOneCompanyProduct(ctx, map[string]interface{}{
	// 	"id": payload.CompanyProductId,
	// })

	// if err != nil {
	// 	return response.Error(http.StatusInternalServerError, err.Error())
	// }

	// if companyProduct == nil {
	// 	return response.Error(http.StatusBadRequest, "company product not found")
	// }

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

	newCompanyNested := model.CompanyNested{
		ID:    company.ID.Hex(),
		Name:  company.Name,
		Image: company.Logo.ID,
		Type:  company.Type,
	}

	// newCompanyProduct := model.CompanyProductNested{
	// 	ID:    companyProduct.ID.Hex(),
	// 	Name:  companyProduct.Name,
	// 	Code:  companyProduct.Code,
	// 	Image: companyProduct.Logo.ID,
	// }

	// check the user
	existingUser, err := u.mongodbRepo.FetchOneCustomer(ctx, map[string]interface{}{
		"email":     payload.Email,
		"companyID": company.ID.Hex(),
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	// get from config
	config := u._CacheConfig(ctx)

	if existingUser != nil {
		// resend verify email
		if config.Registration.VerifyEmail && !existingUser.IsVerified {

			// new token
			existingUser.Token = helpers.RandomString(64)
			existingUser.UpdatedAt = time.Now()

			// save
			u.mongodbRepo.UpdateOneCustomer(context.Background(), map[string]interface{}{
				"id": existingUser.ID,
			}, map[string]interface{}{
				"token":     existingUser.Token,
				"updatedAt": existingUser.UpdatedAt,
			})

			go _sendEmailRegistration(config, *existingUser, company)

			return response.Success(existingUser)
		}

		return response.Error(http.StatusBadRequest, "email already taken")
	}

	defaultVerified := true
	defaultToken := ""

	// need verify email
	if config.Registration.VerifyEmail {
		defaultVerified = false
		defaultToken = helpers.RandomString(64)
	}

	newUser := model.Customer{

		ID:      primitive.NewObjectID(),
		Name:    payload.Name,
		Email:   payload.Email,
		Company: newCompanyNested,
		// CompanyProduct: newCompanyProduct,
		Password:   string(hashedPassword),
		Role:       model.CustomerRole,
		Token:      defaultToken,
		IsVerified: defaultVerified,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err = u.mongodbRepo.CreateCustomer(ctx, &newUser)

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	// need verify email
	if config.Registration.VerifyEmail {
		go _sendEmailRegistration(config, newUser, company)
	}

	return response.Success(newUser)
}
