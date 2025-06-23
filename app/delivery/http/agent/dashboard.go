package http_agent

import (
	"app/domain"

	"github.com/gin-gonic/gin"
)

func (h *routeHandler) handleDashboardRoute(prefixPath string) {
	// (optional). add prefix api version
	api := h.Route.Group(prefixPath)

	api.GET("", h.Middleware.AuthAgent(), h.Dashboard)
	api.GET("/total-ticket", h.Middleware.AuthAgent(), h.TotalTicket)
	api.GET("/total-ticket-now", h.Middleware.AuthAgent(), h.TotalTicketNow)
}

func (r *routeHandler) TotalTicket(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimAgent)

	response := r.Usecase.GetTotalTicket(ctx, claim)
	c.AbortWithStatusJSON(response.Status, response)
}

func (r *routeHandler) TotalTicketNow(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimAgent)

	response := r.Usecase.GetTotalTicketNow(ctx, claim)
	c.AbortWithStatusJSON(response.Status, response)
}

func (r *routeHandler) Dashboard(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimAgent)

	response := r.Usecase.GetDataDashboard(ctx, claim, c.Request.URL.Query())
	c.AbortWithStatusJSON(response.Status, response)
}
