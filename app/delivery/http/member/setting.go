package http_member

import (
	"app/domain"
	"net/http"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/gin-gonic/gin"
)

func (h *routeHandler) handleSettingRoute(prefixPath string) {
	// (optional). add prefix api version
	api := h.Route.Group(prefixPath)

	api.POST("/change-password", h.Middleware.AuthCustomer(), h.ChangePassword)
}

func (h *routeHandler) ChangePassword(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.ChangePasswordRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	claim := c.MustGet("token_data").(domain.JWTClaimUser)

	response := h.Usecase.ChangePassword(ctx, claim, payload)
	c.JSON(response.Status, response)
}
