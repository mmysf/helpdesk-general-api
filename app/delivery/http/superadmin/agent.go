package http_superadmin

import (
	"app/domain"
	"net/http"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/gin-gonic/gin"
)

func (h *routeHandler) handleAgentRoute(prefixPath string) {
	api := h.Route.Group(prefixPath)

	api.GET("/list", h.Middleware.AuthSuperadmin(), h.AgentList)
	api.POST("/create", h.Middleware.AuthSuperadmin(), h.AgentCreate)
	api.GET("/detail/:id", h.Middleware.AuthSuperadmin(), h.AgentDetail)
	api.PATCH("/update/:id", h.Middleware.AuthSuperadmin(), h.AgentUpdate)
	api.DELETE("/delete/:id", h.Middleware.AuthSuperadmin(), h.AgentDelete)
	api.PATCH("/reset-password/:id", h.Middleware.AuthSuperadmin(), h.AgentResetPassword)
}

func (h *routeHandler) AgentList(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimSuperadmin)
	query := c.Request.URL.Query()

	response := h.Usecase.GetAgents(ctx, claim, query)
	c.JSON(response.Status, response)
}

func (h *routeHandler) AgentCreate(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.AccountRequest{}
	err := c.Bind(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	claim := c.MustGet("token_data").(domain.JWTClaimSuperadmin)

	response := h.Usecase.CreateAgent(ctx, claim, payload)
	c.JSON(response.Status, response)
}

func (h *routeHandler) AgentDetail(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimSuperadmin)

	userId := c.Param("id")

	response := h.Usecase.GetAgentDetail(ctx, claim, userId)
	c.JSON(response.Status, response)
}

func (h *routeHandler) AgentUpdate(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimSuperadmin)
	userId := c.Param("id")

	payload := domain.AccountRequest{}
	err := c.Bind(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	response := h.Usecase.UpdateAgent(ctx, claim, userId, payload)
	c.JSON(response.Status, response)
}

func (h *routeHandler) AgentDelete(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimSuperadmin)
	userId := c.Param("id")

	response := h.Usecase.DeleteAgent(ctx, claim, userId)
	c.JSON(response.Status, response)
}

func (h *routeHandler) AgentResetPassword(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimSuperadmin)
	userId := c.Param("id")

	response := h.Usecase.ResetPasswordAgent(ctx, claim, userId)
	c.JSON(response.Status, response)
}
