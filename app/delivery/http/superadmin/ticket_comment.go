package http_superadmin

import (
	"app/domain"
	"net/http"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/gin-gonic/gin"
)

func (h *routeHandler) handleTicketCommentRoute(prefixPath string) {
	api := h.Route.Group(prefixPath)

	api.POST("/add", h.Middleware.AuthSuperadmin(), h.TicketCommentCreate)
	api.GET("/list/:idTicket", h.Middleware.AuthSuperadmin(), h.TicketCommentList)
	api.GET("/detail/:idComment", h.Middleware.AuthSuperadmin(), h.TicketCommentDetail)
}

func (r *routeHandler) TicketCommentCreate(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.SuperadminTicketCommentRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	claim := c.MustGet("token_data").(domain.JWTClaimSuperadmin)

	response := r.Usecase.CreateTicketComment(ctx, claim, payload)
	c.JSON(response.Status, response)
}

func (r *routeHandler) TicketCommentList(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimSuperadmin)
	ticketId := c.Param("idTicket")
	query := c.Request.URL.Query()

	response := r.Usecase.GetTicketCommentList(ctx, claim, ticketId, query)
	c.JSON(response.Status, response)
}

func (r *routeHandler) TicketCommentDetail(c *gin.Context) {
	ctx := c.Request.Context()

	commentId := c.Param("idComment")
	claim := c.MustGet("token_data").(domain.JWTClaimSuperadmin)

	response := r.Usecase.GetTicketCommentDetail(ctx, claim, commentId)
	c.JSON(response.Status, response)
}
