package http_superadmin

import (
	"app/domain"
	"net/http"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/gin-gonic/gin"
)

func (h *routeHandler) handleCustomerRoute(prefixPath string) {
	api := h.Route.Group(prefixPath)
	api.GET("/list", h.Middleware.AuthSuperadmin(), h.CustomerList)
	api.POST("/create", h.Middleware.AuthSuperadmin(), h.CustomerCreate)
	api.GET("/detail/:id", h.Middleware.AuthSuperadmin(), h.CustomerDetail)
	api.PATCH("/update/:id", h.Middleware.AuthSuperadmin(), h.CustomerUpdate)
	api.DELETE("/delete/:id", h.Middleware.AuthSuperadmin(), h.CustomerDelete)
	api.PATCH("/reset-password/:id", h.Middleware.AuthSuperadmin(), h.CustomerResetPassword)
	api.POST("/import", h.Middleware.AuthSuperadmin(), h.CustomerImport)
}

func (h *routeHandler) CustomerList(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimSuperadmin)
	query := c.Request.URL.Query()

	response := h.Usecase.GetCustomers(ctx, claim, query)
	c.JSON(response.Status, response)
}

func (h *routeHandler) CustomerCreate(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.AccountRequest{}
	err := c.Bind(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	claim := c.MustGet("token_data").(domain.JWTClaimSuperadmin)

	response := h.Usecase.CreateCustomer(ctx, claim, payload)
	c.JSON(response.Status, response)
}

func (h *routeHandler) CustomerDetail(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimSuperadmin)

	userId := c.Param("id")

	response := h.Usecase.GetCustomerDetail(ctx, claim, userId)
	c.JSON(response.Status, response)
}

func (h *routeHandler) CustomerUpdate(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimSuperadmin)
	userId := c.Param("id")

	payload := domain.AccountRequest{}

	err := c.Bind(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	response := h.Usecase.UpdateCustomer(ctx, claim, userId, payload)
	c.JSON(response.Status, response)
}

func (h *routeHandler) CustomerDelete(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimSuperadmin)
	userId := c.Param("id")

	response := h.Usecase.DeleteCustomer(ctx, claim, userId)
	c.JSON(response.Status, response)
}

func (h *routeHandler) CustomerResetPassword(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimSuperadmin)
	userId := c.Param("id")

	response := h.Usecase.ResetPasswordCustomer(ctx, claim, userId)
	c.JSON(response.Status, response)
}

func (h *routeHandler) CustomerImport(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimSuperadmin)

	response := h.Usecase.ImportCustomer(ctx, claim, c.Request)
	c.JSON(response.Status, response)
}
