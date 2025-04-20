package http_agent

import (
	"app/domain"

	"github.com/gin-gonic/gin"
)

func (h *routeHandler) handleTicketTimelogsRoute(prefixPath string) {
	// (optional). add prefix api version
	api := h.Route.Group(prefixPath)

	api.GET("/list", h.Middleware.AuthAgent(), h.List)
}

func (r *routeHandler) List(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimAgent)

	response := r.Usecase.GetTicketTimeLogsList(ctx, claim, c.Request.URL.Query())
	c.AbortWithStatusJSON(response.Status, response)
}

