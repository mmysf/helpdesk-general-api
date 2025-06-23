package http_member

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
	api.POST("/register", h.Register)
	api.POST("/verify", h.VerifyRegistration)
	api.POST("/request-password-reset", h.RequestPasswordReset)
	api.POST("/password-reset", h.ResetPassword)
	api.POST("/register-b2b", h.CreateCustomer)

	api.GET("/me", h.Middleware.AuthCustomer(), h.GetMe)
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

func (r *routeHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.RegisterRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	response := r.Usecase.Register(ctx, payload)
	c.JSON(response.Status, response)
}

func (r *routeHandler) GetMe(c *gin.Context) {
	ctx := c.Request.Context()

	tokenData := c.MustGet("token_data")

	response := r.Usecase.GetMe(ctx, tokenData.(domain.JWTClaimUser))
	c.JSON(response.Status, response)
}

func (r *routeHandler) VerifyRegistration(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.VerifyRegisterRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	response := r.Usecase.VerifyRegistration(ctx, payload)
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

func (r *routeHandler) CreateCustomer(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.RegisterRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	response := r.Usecase.RegisterB2B(ctx, payload)
	c.JSON(response.Status, response)
}
