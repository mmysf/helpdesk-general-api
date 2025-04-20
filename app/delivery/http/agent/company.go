package http_agent

import (
	"app/domain"

	"github.com/gin-gonic/gin"
)

func (h *routeHandler) handleCompanyRoute(prefixPath string) {
	// (optional). add prefix api version
	api := h.Route.Group(prefixPath)

	api.GET("/detail", h.Middleware.AuthAgent(), h.CompanyDetail)

}

func (r *routeHandler) CompanyDetail(c *gin.Context) {
	ctx := c.Request.Context()

	response := r.Usecase.GetCompanyDetail(ctx, c.MustGet("token_data").(domain.JWTClaimAgent))
	c.JSON(response.Status, response)
}
