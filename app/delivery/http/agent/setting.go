package http_agent

import (
	"app/domain"
	"net/http"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/gin-gonic/gin"
)

func (h *routeHandler) handleSettingRoute(prefixPath string) {
	// (optional). add prefix api version
	api := h.Route.Group(prefixPath)

	api.POST("/change-password", h.Middleware.AuthAgent(), h.ChangePassword)
	api.POST("/change-domain", h.Middleware.AuthAgent(), h.Middleware.Role("admin"), h.ChangeDomain)
	api.POST("/update-profile", h.Middleware.AuthAgent(), h.UpdateProfile)
	api.POST("/change-color", h.Middleware.AuthAgent(), h.Middleware.Role("admin"), h.ChangeColor)
	api.POST("/upload-profile-picture", h.Middleware.AuthAgent(), h.UploadAgentProfilePicture)
}

func (h *routeHandler) ChangePassword(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.ChangePasswordRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	claim := c.MustGet("token_data").(domain.JWTClaimAgent)

	response := h.Usecase.ChangePassword(ctx, claim, payload)
	c.JSON(response.Status, response)
}

func (h *routeHandler) ChangeDomain(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.ChangeDomainRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	claim := c.MustGet("token_data").(domain.JWTClaimAgent)

	response := h.Usecase.ChangeDomain(ctx, claim, payload)
	c.JSON(response.Status, response)
}

func (h *routeHandler) UpdateProfile(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.UpdateProfileRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	claim := c.MustGet("token_data").(domain.JWTClaimAgent)

	response := h.Usecase.UpdateProfile(ctx, claim, payload)
	c.JSON(response.Status, response)
}

func (h *routeHandler) ChangeColor(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.ChangeColorMode{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	claim := c.MustGet("token_data").(domain.JWTClaimAgent)

	response := h.Usecase.ChangeColor(ctx, claim, payload)
	c.JSON(response.Status, response)
}

func (r *routeHandler) UploadAgentProfilePicture(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimAgent)

	payload := domain.UploadAttachment{}
	c.Bind(&payload)

	response := r.Usecase.UploadAgentProfilePicture(ctx, claim, payload, c.Request)
	c.AbortWithStatusJSON(response.Status, response)
}
