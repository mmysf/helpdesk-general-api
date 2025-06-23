package http_superadmin

import (
	"app/domain"
	"net/http"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/gin-gonic/gin"
)

func (h *routeHandler) handleAuthRoute(prefixPath string) {
	// (optional). add prefix api version
	api := h.Route.Group(prefixPath)

	api.POST("/login", h.Login)

	api.GET("/me", h.Middleware.AuthSuperadmin(), h.GetMe)
}

func (r *routeHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.SuperadminLoginRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	response := r.Usecase.Login(ctx, payload)
	c.JSON(response.Status, response)
}

func (r *routeHandler) GetMe(c *gin.Context) {
	ctx := c.Request.Context()

	response := r.Usecase.GetMe(ctx, c.MustGet("token_data").(domain.JWTClaimSuperadmin))
	c.JSON(response.Status, response)
}
