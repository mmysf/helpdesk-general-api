package http_agent

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
	api.POST("/request-password-reset", h.RequestPasswordReset)
	api.POST("/password-reset", h.ResetPassword)

	api.GET("/me", h.Middleware.AuthAgent(), h.GetMe)
}

func (r *routeHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.LoginRequest{}
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

	claim := c.MustGet("token_data").(domain.JWTClaimAgent)

	response := r.Usecase.GetMe(ctx, claim)
	c.JSON(response.Status, response)
}

func (r *routeHandler) RequestPasswordReset(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.EmailPasswordResetRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	response := r.Usecase.SendEmailPasswordReset(ctx, payload)
	c.JSON(response.Status, response)
}

func (r *routeHandler) ResetPassword(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.PasswordResetRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	response := r.Usecase.PasswordReset(ctx, payload)
	c.JSON(response.Status, response)
}
