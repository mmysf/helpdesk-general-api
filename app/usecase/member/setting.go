package usecase_member

import (
	"app/domain"
	"context"
	"net/http"
	"time"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"golang.org/x/crypto/bcrypt"
)

func (u *appUsecase) ChangePassword(ctx context.Context, claim domain.JWTClaimUser, payload domain.ChangePasswordRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// jwt claim
	UserID := claim.UserID

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

	// check customer
	customer, err := u.mongodbRepo.FetchOneCustomer(ctx, map[string]interface{}{
		"id": UserID,
	})

	if err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	if customer == nil {
		return response.Error(http.StatusBadRequest, "customer not found")
	}

	// check old password
	if err = bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(payload.OldPassword)); err != nil {
		return response.Error(http.StatusBadRequest, "Wrong old password")
	}

	// hash new password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(payload.NewPassword), bcrypt.DefaultCost)

	// update customer password
	customer.Password = string(hashedPassword)

	// save
	if err := u.mongodbRepo.UpdateOneCustomer(ctx, map[string]interface{}{
		"id": customer.ID,
	}, map[string]interface{}{
		"password":  customer.Password,
		"updatedAt": time.Now(),
	}); err != nil {
		return response.Error(http.StatusInternalServerError, err.Error())
	}

	return response.Success(customer)
}
