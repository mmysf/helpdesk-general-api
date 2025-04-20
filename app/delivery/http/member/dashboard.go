package http_member

import (
	"app/domain"

	"github.com/gin-gonic/gin"
)

func (h *routeHandler) handleDashboardRoute(prefixPath string) {
	// (optional). add prefix api version
	api := h.Route.Group(prefixPath)

	api.GET("", h.Middleware.AuthCustomer(), h.Dashboard)
	api.GET("/total-ticket", h.Middleware.AuthCustomer(), h.TotalTicket)
	api.GET("/total-ticket-now", h.Middleware.AuthCustomer(), h.TotalTicketNow)
	api.GET("/average-duration", h.Middleware.AuthCustomer(), h.AverageDuration)
}

func (r *routeHandler) TotalTicket(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimUser)

	response := r.Usecase.GetTotalTicket(ctx, claim)
	c.AbortWithStatusJSON(response.Status, response)
}

func (r *routeHandler) TotalTicketNow(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimUser)

	response := r.Usecase.GetTotalTicketNow(ctx, claim)
	c.AbortWithStatusJSON(response.Status, response)
}

func (r *routeHandler) Dashboard(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimUser)

	response := r.Usecase.GetDataDashboard(ctx, claim, c.Request.URL.Query())
	c.AbortWithStatusJSON(response.Status, response)
}

func (r *routeHandler) AverageDuration(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimUser)

	response := r.Usecase.GetAverageDurationDashboard(ctx, claim, c.Request.URL.Query())
	c.AbortWithStatusJSON(response.Status, response)
}
