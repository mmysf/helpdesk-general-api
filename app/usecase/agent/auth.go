package usecase_agent

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
	"golang.org/x/crypto/bcrypt"
)

func (u *agentUsecase) Login(ctx context.Context, payload domain.LoginRequest) response.Base {
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

	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	// check the db
	user, err := u.mongodbRepo.FetchOneAgent(ctx, map[string]interface{}{
		"email": payload.Email,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if user == nil {
		return response.Error(http.StatusBadRequest, "user not found")
	}

	// check company
	company, err := u.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{
		"id": user.Company.ID,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if company == nil {
		return response.Error(http.StatusBadRequest, "company not found")
	}

	// check password
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		return response.Error(http.StatusBadRequest, "Wrong password")
	}

	// generate token
	tokenString, err := helpers.GenerateJWTTokenAgent(domain.JWTClaimAgent{
		UserID:    user.ID.Hex(),
		CompanyID: user.Company.ID,
		Role:      string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Issuer:    "agent",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(helpers.GetJWTTTL()) * time.Minute)),
		},
	})
	if err != nil {
		return response.Error(http.StatusBadRequest, err.Error())
	}

	u.mongodbRepo.UpdateOneAgent(
		ctx,
		map[string]interface{}{"id": user.ID},
		map[string]interface{}{
			"updatedAt":      time.Now(),
			"lastActivityAt": time.Now(),
		})

	return response.Success(map[string]interface{}{
		"user":  user,
		"token": tokenString,
	})
}

func (u *agentUsecase) SendEmailPasswordReset(ctx context.Context, payload domain.EmailPasswordResetRequest) response.Base {
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
	agent, err := u.mongodbRepo.FetchOneAgent(ctx, map[string]interface{}{
		"email": payload.Email,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if agent == nil {
		return response.Error(http.StatusBadRequest, "user not found")
	}

	company, err := u.mongodbRepo.FetchOneCompany(ctx, map[string]interface{}{
		"id": agent.Company.ID,
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
	agent.PasswordResetToken = helpers.RandomString(64)
	agent.UpdatedAt = time.Now()

	// save
	u.mongodbRepo.UpdateAgent(ctx, agent)

	go _sendEmailPasswordReset(config, *agent, company)

	return response.Success("The password reset email has been sent.")

}

func _sendEmailPasswordReset(config model.Config, agent model.Agent, company *model.Company) {
	passwordResetLink := helpers.StringReplacer(config.ResetPasswordLink, map[string]string{
		"base_url_frontend":  config.AgentLink,
		"passwordResetToken": agent.PasswordResetToken,
	})
	// send email
	mail := helpers.NewSMTPMailer(company)
	mail.To([]string{agent.Email})
	mail.Subject(config.Email.Template.ResetPassword.Title)
	mail.Body(helpers.StringReplacer(config.Email.Template.ResetPassword.Body, map[string]string{
		"title":               config.Email.Template.ResetPassword.Title,
		"reset_password_link": passwordResetLink,
	}))

	// send
	if err := mail.Send(); err != nil {
		logrus.Errorf("Send Email to %s error %v", agent.Email, err)
	}
}

func (u *agentUsecase) PasswordReset(ctx context.Context, payload domain.PasswordResetRequest) response.Base {
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
	agent, err := u.mongodbRepo.FetchOneAgent(ctx, map[string]interface{}{
		"passwordResetToken": payload.Token,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if agent == nil {
		return response.Error(http.StatusBadRequest, "password reset token not valid")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

	// update user
	agent.PasswordResetToken = ""
	agent.Password = string(hashedPassword)
	agent.UpdatedAt = time.Now()

	// save
	u.mongodbRepo.UpdateAgent(ctx, agent)

	return response.Success(agent)
}

func (u *agentUsecase) GetMe(ctx context.Context, claim domain.JWTClaimAgent) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	userID := claim.UserID

	// check the db
	user, err := u.mongodbRepo.FetchOneAgent(ctx, map[string]interface{}{
		"id": userID,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if user == nil {
		return response.Error(http.StatusBadRequest, "user not found")
	}

	return response.Success(user)
}
