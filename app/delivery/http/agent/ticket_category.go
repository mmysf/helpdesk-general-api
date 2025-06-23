package http_agent

import (
	"app/domain"
	"net/http"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/gin-gonic/gin"
)

func (h *routeHandler) handleTicketCategoryRoute(prefixPath string) {
	// (optional). add prefix api version
	api := h.Route.Group(prefixPath)

	api.GET("/list", h.Middleware.AuthAgent(), h.TicketCategoryList)
	api.GET("/detail/:id", h.Middleware.AuthAgent(), h.TicketCategoryDetail)
	api.POST("/create", h.Middleware.AuthAgent(), h.Middleware.Role("admin"), h.TicketCategoryCreate)
	api.PUT("/update/:id", h.Middleware.AuthAgent(), h.Middleware.Role("admin"), h.TicketCategoryUpdate)
	api.DELETE("/delete/:id", h.Middleware.AuthAgent(), h.Middleware.Role("admin"), h.TicketCategoryDelete)
}

func (r *routeHandler) TicketCategoryList(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimAgent)
	query := c.Request.URL.Query()

	response := r.Usecase.GetTicketCategoriesList(ctx, claim, query)
	c.JSON(response.Status, response)
}

func (r *routeHandler) TicketCategoryDetail(c *gin.Context) {
	ctx := c.Request.Context()

	ticketCategoryID := c.Param("id")

	response := r.Usecase.GetTicketCategoryDetail(ctx, c.MustGet("token_data").(domain.JWTClaimAgent), ticketCategoryID)
	c.JSON(response.Status, response)
}

func (r *routeHandler) TicketCategoryCreate(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.TicketCategoryRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	claim := c.MustGet("token_data").(domain.JWTClaimAgent)

	response := r.Usecase.CreateTicketCategory(ctx, claim, payload)
	c.JSON(response.Status, response)
}

func (r *routeHandler) TicketCategoryUpdate(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.TicketCategoryRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	claim := c.MustGet("token_data").(domain.JWTClaimAgent)
	ticketCategoryID := c.Param("id")

	response := r.Usecase.UpdateTicketCategory(ctx, claim, ticketCategoryID, payload)
	c.JSON(response.Status, response)
}

func (r *routeHandler) TicketCategoryDelete(c *gin.Context) {
	ctx := c.Request.Context()

	ticketCategoryID := c.Param("id")

	response := r.Usecase.DeleteTicketCategory(ctx, c.MustGet("token_data").(domain.JWTClaimAgent), ticketCategoryID)
	c.JSON(response.Status, response)
}
