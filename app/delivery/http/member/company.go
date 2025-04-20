package http_member

import (
	"github.com/gin-gonic/gin"
)

func (h *routeHandler) handleCompanyRoute(prefixPath string) {
	// (optional). add prefix api version
	api := h.Route.Group(prefixPath)

	api.GET("/detail-by-domain/:domain", h.CompanyDetailByDomain)
}

func (r *routeHandler) CompanyDetailByDomain(c *gin.Context) {
	ctx := c.Request.Context()

	options := map[string]interface{}{
		"domain": c.Param("domain"),
	}

	response := r.Usecase.GetCompanyDetailByDomain(ctx, options)
	c.JSON(response.Status, response)
}
