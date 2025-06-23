package http_member

import (
	"app/domain"

	"github.com/gin-gonic/gin"
)

func (h *routeHandler) handleHourPackageRoute(prefixPath string) {
	// (optional). add prefix api version
	api := h.Route.Group(prefixPath)

	api.GET("/list", h.Middleware.AuthCustomer(), h.HourPackageList)
	api.GET("/detail/:id", h.Middleware.AuthCustomer(), h.HourPackageDetail)
}

func (r *routeHandler) HourPackageList(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimUser)
	query := c.Request.URL.Query()

	response := r.Usecase.GetHourPackageList(ctx, claim, query)
	c.JSON(response.Status, response)
}

func (r *routeHandler) HourPackageDetail(c *gin.Context) {
	ctx := c.Request.Context()

	packageID := c.Param("id")

	response := r.Usecase.GetHourPackageDetail(ctx, c.MustGet("token_data").(domain.JWTClaimUser), packageID)
	c.JSON(response.Status, response)
}
