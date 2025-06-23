package http_superadmin

import (
	"app/domain"

	"github.com/gin-gonic/gin"
)

func (h *routeHandler) handleDashboardRoute(prefixPath string) {
	// (optional). add prefix api version
	api := h.Route.Group(prefixPath)

	api.GET("/data", h.Middleware.AuthSuperadmin(), h.Dashboard)
	api.GET("/hour-packages", h.Middleware.AuthSuperadmin(), h.HourPackageList)
}

func (r *routeHandler) Dashboard(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimSuperadmin)

	response := r.Usecase.GetDataDashboard(ctx, claim)
	c.AbortWithStatusJSON(response.Status, response)
}

func (h *routeHandler) HourPackageList(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimSuperadmin)
	query := c.Request.URL.Query()

	response := h.Usecase.GetHourPackagesDashboard(ctx, claim, query)
	c.JSON(response.Status, response)
}
