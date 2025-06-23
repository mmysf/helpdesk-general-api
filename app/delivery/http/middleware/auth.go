package middleware

import (
	"app/domain"
	"app/domain/model"
	"app/helpers"
	"errors"
	"net/http"
	"strings"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (m *appMiddleware) AuthCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		hAuth := c.GetHeader("Authorization")
		if hAuth == "" {
			response := response.Error(http.StatusUnauthorized, "Unauthorized: Header authorization is required")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		splitToken := strings.Split(hAuth, "Bearer ")
		if len(splitToken) != 2 {
			response := response.Error(http.StatusUnauthorized, "Unauthorized: Token is invalid")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// get token without 'Bearer '
		tokenString := splitToken[1]

		// validating token
		token, err := jwt.ParseWithClaims(tokenString, &domain.JWTClaimUser{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.secretKeyCustomer), nil
		})
		if err != nil {
			response := response.Error(http.StatusUnauthorized, err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// check validity token
		if !token.Valid {
			if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
				c.AbortWithStatusJSON(
					http.StatusUnauthorized,
					response.Error(http.StatusUnauthorized, "Unauthorized: Token signature invalid"),
				)
				return
			}

			if errors.Is(err, jwt.ErrTokenExpired) {
				c.AbortWithStatusJSON(
					http.StatusUnauthorized,
					response.Error(http.StatusUnauthorized, "Unauthorized: Token expired"),
				)
				return
			}
		}

		claims, tokenOK := token.Claims.(*domain.JWTClaimUser)
		if !tokenOK {
			response := response.Error(http.StatusUnauthorized, "Unauthorized: Token data not valid")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		//check customer
		customer, err := m.mongo.FetchOneCustomer(c, map[string]interface{}{
			"id": claims.UserID,
		})
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				response.Error(http.StatusInternalServerError, err.Error()),
			)
			return
		}
		if customer == nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				response.Error(http.StatusUnauthorized, "Unauthorized: User not found"),
			)
			return
		}
		claims.User = model.UserNested{
			ID:    customer.ID.Hex(),
			Name:  customer.Name,
			Email: customer.Email,
		}

		//check company
		company, err := m.mongo.FetchOneCompany(c, map[string]interface{}{
			"id": claims.CompanyID,
		})
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				response.Error(http.StatusInternalServerError, err.Error()),
			)
			return
		}
		if company == nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				response.Error(http.StatusUnauthorized, "Unauthorized: Company not found"),
			)
			return
		}
		claims.Company = model.CompanyNested{
			ID:    company.ID.Hex(),
			Name:  company.Name,
			Image: company.Logo.URL,
			Type:  company.Type,
		}

		//check companyProduct
		companyProduct, err := m.mongo.FetchOneCompanyProduct(c, map[string]interface{}{
			"id": claims.CompanyProductID,
		})
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				response.Error(http.StatusInternalServerError, err.Error()),
			)
			return
		}
		if companyProduct == nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				response.Error(http.StatusUnauthorized, "Unauthorized: Company product not found"),
			)
			return
		}
		claims.CompanyProduct = model.CompanyProductNested{
			ID:    companyProduct.ID.Hex(),
			Name:  companyProduct.Name,
			Image: companyProduct.Logo.URL,
			Code:  companyProduct.Code,
		}

		c.Set("token_data", *claims)
		c.Next()
	}
}

func (m *appMiddleware) AuthAgent() gin.HandlerFunc {
	return func(c *gin.Context) {
		hAuth := c.GetHeader("Authorization")
		if hAuth == "" {
			response := response.Error(http.StatusUnauthorized, "Unauthorized: Header authorization is required")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		splitToken := strings.Split(hAuth, "Bearer ")
		if len(splitToken) != 2 {
			response := response.Error(http.StatusUnauthorized, "Unauthorized: Token is invalid")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// get token without 'Bearer '
		tokenString := splitToken[1]

		// validating token
		token, err := jwt.ParseWithClaims(tokenString, &domain.JWTClaimAgent{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.secretKeyAgent), nil
		})
		if err != nil {
			response := response.Error(http.StatusUnauthorized, err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// check validity token
		if !token.Valid {
			if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
				c.AbortWithStatusJSON(
					http.StatusUnauthorized,
					response.Error(http.StatusUnauthorized, "Unauthorized: Token signature invalid"),
				)
				return
			}

			if errors.Is(err, jwt.ErrTokenExpired) {
				c.AbortWithStatusJSON(
					http.StatusUnauthorized,
					response.Error(http.StatusUnauthorized, "Unauthorized: Token expired"),
				)
				return
			}
			return
		}

		claims, tokenOK := token.Claims.(*domain.JWTClaimAgent)
		if !tokenOK {
			response := response.Error(http.StatusUnauthorized, "Unauthorized: Token data not valid")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		//check agent
		agent, err := m.mongo.FetchOneAgent(c, map[string]interface{}{
			"id": claims.UserID,
		})
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				response.Error(http.StatusInternalServerError, err.Error()),
			)
			return
		}
		if agent == nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				response.Error(http.StatusUnauthorized, "Unauthorized: User not found"),
			)
			return
		}
		claims.User = model.UserNested{
			ID:    agent.ID.Hex(),
			Name:  agent.Name,
			Email: agent.Email,
		}

		//check company
		company, err := m.mongo.FetchOneCompany(c, map[string]interface{}{
			"id": claims.CompanyID,
		})
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				response.Error(http.StatusInternalServerError, err.Error()),
			)
			return
		}
		if company == nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				response.Error(http.StatusUnauthorized, "Unauthorized: Company not found"),
			)
			return
		}
		claims.Company = model.CompanyNested{
			ID:    company.ID.Hex(),
			Name:  company.Name,
			Image: company.Logo.URL,
			Type:  company.Type,
		}

		c.Set("token_data", *claims)
		c.Next()
	}
}

func (m *appMiddleware) AuthSuperadmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		hAuth := c.GetHeader("Authorization")
		if hAuth == "" {
			response := response.Error(http.StatusUnauthorized, "Unauthorized: Header authorization is required")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		splitToken := strings.Split(hAuth, "Bearer ")
		if len(splitToken) != 2 {
			response := response.Error(http.StatusUnauthorized, "Unauthorized: Token is invalid")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// get token without 'Bearer '
		tokenString := splitToken[1]

		// validating token
		token, err := jwt.ParseWithClaims(tokenString, &domain.JWTClaimSuperadmin{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.secretKeySuperadmin), nil
		})
		if err != nil {
			response := response.Error(http.StatusUnauthorized, err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// check validity token
		if !token.Valid {
			if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
				c.AbortWithStatusJSON(
					http.StatusUnauthorized,
					response.Error(http.StatusUnauthorized, "Unauthorized: Token signature invalid"),
				)
				return
			}

			if errors.Is(err, jwt.ErrTokenExpired) {
				c.AbortWithStatusJSON(
					http.StatusUnauthorized,
					response.Error(http.StatusUnauthorized, "Unauthorized: Token expired"),
				)
				return
			}
		}

		claims, tokenOK := token.Claims.(*domain.JWTClaimSuperadmin)
		if !tokenOK {
			response := response.Error(http.StatusUnauthorized, "Unauthorized: Token data not valid")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		//check admin
		admin, err := m.mongo.FetchOneSuperadmin(c, map[string]interface{}{
			"id": claims.UserID,
		})
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				response.Error(http.StatusInternalServerError, err.Error()),
			)
			return
		}
		if admin == nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				response.Error(http.StatusUnauthorized, "Unauthorized: User not found"),
			)
			return
		}
		claims.User = model.UserNested{
			ID:    admin.ID.Hex(),
			Name:  admin.Name,
			Email: admin.Email,
		}

		c.Set("token_data", *claims)
		c.Next()
	}
}

func (m *appMiddleware) Role(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenData := c.MustGet("token_data")

		var userRole string

		switch tokenData.(type) {
		case domain.JWTClaimAgent:
			userRole = tokenData.(domain.JWTClaimAgent).Role
		case domain.JWTClaimUser:
			userRole = tokenData.(domain.JWTClaimUser).Role
		default:
			response := response.Error(http.StatusForbidden, "Forbidden: You don't have permission to access this resource")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		if !helpers.InArrayString(userRole, allowedRoles) {
			response := response.Error(http.StatusForbidden, "Forbidden: You don't have permission to access this resource")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Next()
	}
}
