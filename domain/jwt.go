package domain

import (
	"app/domain/model"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaimUser struct {
	UserID           string                     `json:"userID"`
	CompanyID        string                     `json:"companyID"`
	CompanyProductID string                     `json:"companyProductID"`
	Role             string                     `json:"role"`
	Company          model.CompanyNested        `json:"-"`
	CompanyProduct   model.CompanyProductNested `json:"-"`
	User             model.UserNested           `json:"-"`
	jwt.RegisteredClaims
}

type JWTClaimAgent struct {
	UserID    string              `json:"userID"`
	CompanyID string              `json:"companyID"`
	Role      string              `json:"role"`
	Company   model.CompanyNested `json:"-"`
	User      model.UserNested    `json:"-"`
	jwt.RegisteredClaims
}

type JWTClaimSuperadmin struct {
	UserID string           `json:"userID"`
	User   model.UserNested `json:"-"`
	jwt.RegisteredClaims
}
