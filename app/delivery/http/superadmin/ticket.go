package http_superadmin

import (
	"app/domain"
	"net/http"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/gin-gonic/gin"
)

func (h *routeHandler) handleTicketRoute(prefixPath string) {
	api := h.Route.Group(prefixPath)

	api.GET("/total-ticket", h.Middleware.AuthSuperadmin(), h.TotalTicket)
	api.GET("/list", h.Middleware.AuthSuperadmin(), h.TicketList)
	api.GET("/detail/:id", h.TicketDetail)
	api.POST("/assign-agent/:id", h.Middleware.AuthSuperadmin(), h.AssignAgent)
	api.GET("/total-ticket-day/:id", h.Middleware.AuthSuperadmin(), h.TotalTicketClientDays)
	api.GET("/average-duration/:id", h.Middleware.AuthSuperadmin(), h.AverageTicketClient)
	api.POST("/logging/pause", h.Middleware.AuthSuperadmin(), h.TicketLogPause)
	api.POST("/logging/resume", h.Middleware.AuthSuperadmin(), h.TicketLogResume)
}

func (r *routeHandler) TotalTicket(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimSuperadmin)

	response := r.Usecase.GetTotalTicket(ctx, claim)
	c.AbortWithStatusJSON(response.Status, response)
}

func (h *routeHandler) TicketList(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimSuperadmin)
	query := c.Request.URL.Query()

	response := h.Usecase.GetTicketList(ctx, claim, query)
	c.JSON(response.Status, response)
}

func (h *routeHandler) TicketDetail(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	response := h.Usecase.GetTicketDetail(ctx, id)
	c.JSON(response.Status, response)
}

func (h *routeHandler) AssignAgent(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")

	claim := c.MustGet("token_data").(domain.JWTClaimSuperadmin)
	payload := domain.AssignAgentRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	response := h.Usecase.AssignAgent(ctx, claim, id, payload)
	c.JSON(response.Status, response)
}

func (r *routeHandler) TotalTicketClientDays(c *gin.Context) {
	ctx := c.Request.Context()

	options := map[string]interface{}{
		"id": c.Param("id"),
	}

	response := r.Usecase.GetDataClientTicket(ctx, c.MustGet("token_data").(domain.JWTClaimSuperadmin), options)
	c.JSON(response.Status, response)
}

func (r *routeHandler) AverageTicketClient(c *gin.Context) {
	ctx := c.Request.Context()

	options := map[string]interface{}{
		"id": c.Param("id"),
	}

	response := r.Usecase.GetAverageDurationClient(ctx, c.MustGet("token_data").(domain.JWTClaimSuperadmin), options)
	c.JSON(response.Status, response)
}

func (r *routeHandler) TicketLogPause(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.LoggingTicketRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	response := r.Usecase.PauseLoggingTicket(ctx, c.MustGet("token_data").(domain.JWTClaimSuperadmin), payload)
	c.JSON(response.Status, response)
}

func (r *routeHandler) TicketLogResume(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.LoggingTicketRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	response := r.Usecase.ResumeLoggingTicket(ctx, c.MustGet("token_data").(domain.JWTClaimSuperadmin), payload)
	c.JSON(response.Status, response)
}
