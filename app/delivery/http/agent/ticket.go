package http_agent

import (
	"app/domain"
	"net/http"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/gin-gonic/gin"
)

func (h *routeHandler) handleTicketRoute(prefixPath string) {
	// (optional). add prefix api version
	api := h.Route.Group(prefixPath)

	api.GET("/list", h.Middleware.AuthAgent(), h.TicketList)
	api.GET("/mine/list", h.Middleware.AuthAgent(), h.MyTicketList)
	api.GET("/detail/:id", h.Middleware.AuthAgent(), h.TicketDetail)
	api.POST("/close", h.Middleware.AuthAgent(), h.TicketClose)
	api.POST("/reopen", h.Middleware.AuthAgent(), h.TicketReopen)
	// api.POST("/logging/start", h.Middleware.AuthAgent(), h.TicketLogStart)
	// api.POST("/logging/stop", h.Middleware.AuthAgent(), h.TicketLogStop)
	api.POST("/logging/pause", h.Middleware.AuthAgent(), h.TicketLogPause)
	api.POST("/logging/resume", h.Middleware.AuthAgent(), h.TicketLogResume)
	api.POST("/comments/add", h.Middleware.AuthAgent(), h.TicketCommentCreate)
	api.GET("/comments/list/:idTicket", h.Middleware.AuthAgent(), h.TicketCommentList)
	api.GET("/comments/detail/:idComment", h.Middleware.AuthAgent(), h.TicketCommentDetail)
	api.PUT("/time-track/update/:idTicket", h.Middleware.AuthAgent(), h.TimeTrack)
	api.GET("/export-csv", h.Middleware.AuthAgent(), h.ExportTicketsToCSV)
	api.POST("/assign-me/:ticket_id", h.Middleware.AuthAgent(), h.AssignMe)
	api.GET("/total-ticket/:id", h.Middleware.AuthAgent(), h.TotalTicketCustomer)
	api.GET("/total-ticket-day/:id", h.Middleware.AuthAgent(), h.TotalTicketCustomerDays)
}

func (r *routeHandler) TicketList(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimAgent)
	query := c.Request.URL.Query()

	response := r.Usecase.GetTicketList(ctx, claim, query)
	c.JSON(response.Status, response)
}

func (r *routeHandler) MyTicketList(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimAgent)
	query := c.Request.URL.Query()

	response := r.Usecase.GetMyTicketList(ctx, claim, query)
	c.JSON(response.Status, response)
}

func (r *routeHandler) TicketDetail(c *gin.Context) {
	ctx := c.Request.Context()

	ticketId := c.Param("id")

	response := r.Usecase.GetTicketDetail(ctx, c.MustGet("token_data").(domain.JWTClaimAgent), ticketId)
	c.JSON(response.Status, response)
}

func (r *routeHandler) TicketClose(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.CloseTicketRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	response := r.Usecase.CloseTicket(ctx, c.MustGet("token_data").(domain.JWTClaimAgent), payload)
	c.JSON(response.Status, response)
}

func (r *routeHandler) TicketReopen(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.ReopenTicketRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	response := r.Usecase.ReopenTicket(ctx, c.MustGet("token_data").(domain.JWTClaimAgent), payload)
	c.JSON(response.Status, response)
}

func (r *routeHandler) TicketLogStart(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.LoggingTicketRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	response := r.Usecase.StartLoggingTicket(ctx, c.MustGet("token_data").(domain.JWTClaimAgent), payload)
	c.JSON(response.Status, response)
}

func (r *routeHandler) TicketLogStop(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.LoggingTicketRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	response := r.Usecase.StopLoggingTicket(ctx, c.MustGet("token_data").(domain.JWTClaimAgent), payload)
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

	response := r.Usecase.PauseLoggingTicket(ctx, c.MustGet("token_data").(domain.JWTClaimAgent), payload)
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

	response := r.Usecase.ResumeLoggingTicket(ctx, c.MustGet("token_data").(domain.JWTClaimAgent), payload)
	c.JSON(response.Status, response)
}

func (r *routeHandler) TicketCommentCreate(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.TicketCommentRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	claim := c.MustGet("token_data").(domain.JWTClaimAgent)

	response := r.Usecase.CreateTicketComment(ctx, claim, payload)
	c.JSON(response.Status, response)
}

func (r *routeHandler) TicketCommentList(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimAgent)
	ticketId := c.Param("idTicket")
	query := c.Request.URL.Query()

	response := r.Usecase.GetTicketCommentList(ctx, claim, ticketId, query)
	c.JSON(response.Status, response)
}

func (r *routeHandler) TicketCommentDetail(c *gin.Context) {
	ctx := c.Request.Context()

	commentId := c.Param("idComment")

	response := r.Usecase.GetTicketCommentDetail(ctx, c.MustGet("token_data").(domain.JWTClaimAgent), commentId)
	c.JSON(response.Status, response)
}

func (r *routeHandler) TimeTrack(c *gin.Context) {
	ctx := c.Request.Context()

	ticketId := c.Param("idTicket")

	payload := domain.TimeTrackRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	response := r.Usecase.EditTimeTrack(ctx, c.MustGet("token_data").(domain.JWTClaimAgent), ticketId, payload)
	c.JSON(response.Status, response)
}

func (r *routeHandler) ExportTicketsToCSV(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimAgent)
	query := c.Request.URL.Query()

	response := r.Usecase.ExportTicketsToCSV(ctx, claim, query, c.Writer)
	if response.Status != http.StatusOK {
		c.JSON(response.Status, response)
	}
}

func (r *routeHandler) AssignMe(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimAgent)
	ticketId := c.Param("ticket_id")

	response := r.Usecase.AssignTicketToMe(ctx, claim, ticketId)
	c.JSON(response.Status, response)
}

func (r *routeHandler) TotalTicketCustomer(c *gin.Context) {
	ctx := c.Request.Context()

	options := map[string]interface{}{
		"id": c.Param("id"),
	}

	response := r.Usecase.GetTotalTicketCustomer(ctx, c.MustGet("token_data").(domain.JWTClaimAgent), options)
	c.JSON(response.Status, response)
}

func (r *routeHandler) TotalTicketCustomerDays(c *gin.Context) {
	ctx := c.Request.Context()

	options := map[string]interface{}{
		"id": c.Param("id"),
	}

	response := r.Usecase.GetDataCustomerTicket(ctx, c.MustGet("token_data").(domain.JWTClaimAgent), options)
	c.JSON(response.Status, response)
}
