package http_member

import (
	"app/domain"

	"github.com/gin-gonic/gin"
)

func (h *routeHandler) handleServerPackageRoute(prefixPath string) {
	api := h.Route.Group(prefixPath, h.Middleware.AuthCustomer())
	api.GET("/list", h.ServerPackageList)
	api.GET("/detail/:id", h.ServerPackageDetail)

}

func (h *routeHandler) ServerPackageList(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimUser)
	query := c.Request.URL.Query()

	response := h.Usecase.GetServerPackageList(ctx, claim, query)
	c.JSON(response.Status, response)
}

func (r *routeHandler) ServerPackageDetail(c *gin.Context) {
	ctx := c.Request.Context()

	productId := c.Param("id")

	response := r.Usecase.GetServerPackageDetail(ctx, productId)
	c.JSON(response.Status, response)
}
