package usecase_superadmin

import (
	"app/domain"
	"app/helpers"
	"context"
	"net/http"
	"time"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (u *superadminUsecase) Login(ctx context.Context, payload domain.SuperadminLoginRequest) response.Base {
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
	user, err := u.mongodbRepo.FetchOneSuperadmin(ctx, map[string]interface{}{
		"email": payload.Email,
	})
	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}
	if user == nil {
		return response.Error(http.StatusBadRequest, "user not found")
	}

	// check password
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		return response.Error(http.StatusBadRequest, "Wrong password")
	}

	// generate token
	tokenString, err := helpers.GenerateJWTTokenSuperadmin(domain.JWTClaimSuperadmin{
		UserID: user.ID.Hex(),
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Issuer:    "superadmin",
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

func (u *superadminUsecase) GetMe(ctx context.Context, claim domain.JWTClaimSuperadmin) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	userID := claim.UserID

	// check the db
	user, err := u.mongodbRepo.FetchOneSuperadmin(ctx, map[string]interface{}{
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
