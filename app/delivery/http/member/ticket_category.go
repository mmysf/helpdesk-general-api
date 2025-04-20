package http_member

import (
	"app/domain"

	"github.com/gin-gonic/gin"
)

func (h *routeHandler) handleTicketCategoryRoute(prefixPath string) {
	// (optional). add prefix api version
	api := h.Route.Group(prefixPath)

	api.GET("/list", h.Middleware.AuthCustomer(), h.TicketCategoryList)
}

func (r *routeHandler) TicketCategoryList(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimUser)
	query := c.Request.URL.Query()

	response := r.Usecase.GetTicketCategoriesList(ctx, claim, query)
	c.JSON(response.Status, response)
}
