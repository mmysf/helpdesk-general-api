package http_superadmin

import (
	"app/domain"
	"net/http"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/gin-gonic/gin"
)

func (h *routeHandler) handleCompanyRoute(prefixPath string) {
	// (optional). add prefix api version
	api := h.Route.Group(prefixPath)

	api.GET("/list", h.Middleware.AuthSuperadmin(), h.CompanyList)
	api.GET("/detail/:id", h.Middleware.AuthSuperadmin(), h.CompanyDetail)
	api.POST("/upload-logo", h.Middleware.AuthSuperadmin(), h.UploadCompanyLogo)
	api.POST("/create", h.Middleware.AuthSuperadmin(), h.CompanyCreate)
	api.PUT("/update/:id", h.Middleware.AuthSuperadmin(), h.CompanyUpdate)
	api.DELETE("/delete/:id", h.Middleware.AuthSuperadmin(), h.CompanyDelete)
}

func (r *routeHandler) CompanyList(c *gin.Context) {
	ctx := c.Request.Context()

	options := map[string]interface{}{
		"query": c.Request.URL.Query(),
	}

	response := r.Usecase.GetCompanyList(ctx, c.MustGet("token_data").(domain.JWTClaimSuperadmin), options)
	c.JSON(response.Status, response)
}

func (r *routeHandler) CompanyDetail(c *gin.Context) {
	ctx := c.Request.Context()

	options := map[string]interface{}{
		"id": c.Param("id"),
	}

	response := r.Usecase.GetCompanyDetail(ctx, c.MustGet("token_data").(domain.JWTClaimSuperadmin), options)
	c.JSON(response.Status, response)
}

func (r *routeHandler) CompanyCreate(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.CreateCompanyRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	options := map[string]interface{}{
		"payload": payload,
	}

	response := r.Usecase.CreateCompany(ctx, c.MustGet("token_data").(domain.JWTClaimSuperadmin), options)
	c.JSON(response.Status, response)
}

func (r *routeHandler) CompanyUpdate(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.CreateCompanyRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	options := map[string]interface{}{
		"id":      c.Param("id"),
		"payload": payload,
	}

	response := r.Usecase.UpdateCompany(ctx, c.MustGet("token_data").(domain.JWTClaimSuperadmin), options)
	c.JSON(response.Status, response)
}

func (r *routeHandler) CompanyDelete(c *gin.Context) {
	ctx := c.Request.Context()

	options := map[string]interface{}{
		"id": c.Param("id"),
	}

	response := r.Usecase.DeleteCompany(ctx, c.MustGet("token_data").(domain.JWTClaimSuperadmin), options)
	c.JSON(response.Status, response)
}

func (r *routeHandler) UploadCompanyLogo(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimSuperadmin)

	payload := domain.UploadAttachment{}
	c.Bind(&payload)

	response := r.Usecase.UploadCompanyLogo(ctx, claim, payload, c.Request)
	c.AbortWithStatusJSON(response.Status, response)
}
