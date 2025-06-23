package http_agent

import (
	"app/domain"

	"github.com/gin-gonic/gin"
)

func (h *routeHandler) handleProductRoute(prefixPath string) {
	// (optional). add prefix api version
	api := h.Route.Group(prefixPath)

	api.GET("/list", h.Middleware.AuthAgent(), h.ProductList)
}

func (r *routeHandler) ProductList(c *gin.Context) {
	ctx := c.Request.Context()

	response := r.Usecase.GetProductList(ctx, c.MustGet("token_data").(domain.JWTClaimAgent), c.Request.URL.Query())
	c.JSON(response.Status, response)
}
