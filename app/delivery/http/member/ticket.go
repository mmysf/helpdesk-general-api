package http_member

import (
	"app/domain"
	"net/http"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/gin-gonic/gin"
)

func (h *routeHandler) handleTicketRoute(prefixPath string) {
	// (optional). add prefix api version
	api := h.Route.Group(prefixPath)

	api.GET("/list", h.Middleware.AuthCustomer(), h.TicketList)
	api.GET("/detail/:id", h.Middleware.AuthCustomer(), h.TicketDetail)
	api.POST("/create", h.Middleware.AuthCustomer(), h.TicketCreate)
	api.POST("/close", h.Middleware.AuthCustomer(), h.TicketClose)
	api.POST("/close-by-email", h.TicketCloseByEmail)
	api.POST("/comments/add", h.Middleware.AuthCustomer(), h.TicketCommentCreate)
	api.GET("/comments/list/:idTicket", h.Middleware.AuthCustomer(), h.TicketCommentList)
	api.GET("/comments/detail/:idComment", h.Middleware.AuthCustomer(), h.TicketCommentDetail)
	api.POST("/reopen", h.Middleware.AuthCustomer(), h.TicketReopen)
	api.POST("/cancel", h.Middleware.AuthCustomer(), h.CancelTicket)
}

func (r *routeHandler) TicketList(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimUser)
	query := c.Request.URL.Query()

	response := r.Usecase.GetTicketList(ctx, claim, query)
	c.JSON(response.Status, response)
}

func (r *routeHandler) TicketDetail(c *gin.Context) {
	ctx := c.Request.Context()

	ticketID := c.Param("id")

	response := r.Usecase.GetTicketDetail(ctx, c.MustGet("token_data").(domain.JWTClaimUser), ticketID)
	c.JSON(response.Status, response)
}

func (r *routeHandler) TicketCreate(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.TicketRequest{}
	err := c.Bind(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	claim := c.MustGet("token_data").(domain.JWTClaimUser)

	response := r.Usecase.CreateTicket(ctx, claim, payload)
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
	response := r.Usecase.CloseTicket(ctx, c.MustGet("token_data").(domain.JWTClaimUser), payload)
	c.JSON(response.Status, response)
}

func (r *routeHandler) TicketCloseByEmail(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.CloseTicketbyEmailRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}
	response := r.Usecase.CloseTicketByEmail(ctx, payload)
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

	claim := c.MustGet("token_data").(domain.JWTClaimUser)

	response := r.Usecase.CreateTicketComment(ctx, claim, payload)
	c.JSON(response.Status, response)
}

func (r *routeHandler) TicketCommentList(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimUser)
	ticketId := c.Param("idTicket")
	query := c.Request.URL.Query()

	response := r.Usecase.GetTicketCommentList(ctx, claim, ticketId, query)
	c.JSON(response.Status, response)
}

func (r *routeHandler) TicketCommentDetail(c *gin.Context) {
	ctx := c.Request.Context()

	commentId := c.Param("idComment")

	response := r.Usecase.GetTicketCommentDetail(ctx, c.MustGet("token_data").(domain.JWTClaimUser), commentId)
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

	response := r.Usecase.ReopenTicket(ctx, c.MustGet("token_data").(domain.JWTClaimUser), payload)
	c.JSON(response.Status, response)
}

func (r *routeHandler) CancelTicket(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.CancelTicketRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	response := r.Usecase.CancelTicket(ctx, c.MustGet("token_data").(domain.JWTClaimUser), payload)
	c.JSON(response.Status, response)
}
