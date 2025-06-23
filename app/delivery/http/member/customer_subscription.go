package http_member

import (
	"app/domain"

	"github.com/gin-gonic/gin"
)

func (h *routeHandler) handleCustomerSubscriptionRoute(prefixPath string) {
	// (optional). add prefix api version
	api := h.Route.Group(prefixPath)

	api.GET("/list", h.Middleware.AuthCustomer(), h.CustomerSubscriptionList)
	api.GET("/detail/:id", h.Middleware.AuthCustomer(), h.CustomerSubscriptionDetail)
}

func (h *routeHandler) CustomerSubscriptionList(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimUser)
	query := c.Request.URL.Query()

	response := h.Usecase.GetCustomerSubscriptionList(ctx, claim, query)
	c.JSON(response.Status, response)
}

func (h *routeHandler) CustomerSubscriptionDetail(c *gin.Context) {
	ctx := c.Request.Context()

	customerSubscriptionID := c.Param("id")

	response := h.Usecase.GetCustomerSubscriptionDetail(ctx, c.MustGet("token_data").(domain.JWTClaimUser), customerSubscriptionID)
	c.JSON(response.Status, response)
}
