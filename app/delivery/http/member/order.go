package http_member

import (
	"app/domain"
	"net/http"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/gin-gonic/gin"
)

func (h *routeHandler) handleOrderRoute(prefixPath string) {
	// (optional). add prefix api version
	api := h.Route.Group(prefixPath)

	api.GET("/list", h.Middleware.AuthCustomer(), h.OrderList)
	api.GET("/detail/:id", h.Middleware.AuthCustomer(), h.OrderDetail)
	api.POST("/hour/create", h.Middleware.AuthCustomer(), h.OrderHourCreate)
	api.POST("/server/create", h.Middleware.AuthCustomer(), h.OrderServerCreate)
	api.POST("/confirm", h.Middleware.AuthCustomer(), h.ConfirmOrder)
	api.POST("/upload-attachment", h.Middleware.AuthCustomer(), h.UploadAttachmentOrder)
}

func (h *routeHandler) OrderList(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimUser)
	query := c.Request.URL.Query()

	response := h.Usecase.GetOrderList(ctx, claim, query)
	c.JSON(response.Status, response)
}

func (h *routeHandler) OrderDetail(c *gin.Context) {
	ctx := c.Request.Context()

	orderID := c.Param("id")

	response := h.Usecase.GetOrderDetail(ctx, c.MustGet("token_data").(domain.JWTClaimUser), orderID)
	c.JSON(response.Status, response)
}

func (h *routeHandler) OrderHourCreate(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.OrderRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	claim := c.MustGet("token_data").(domain.JWTClaimUser)

	response := h.Usecase.CreateHourOrder(ctx, claim, payload)
	c.JSON(response.Status, response)
}

func (h *routeHandler) OrderServerCreate(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.OrderRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	claim := c.MustGet("token_data").(domain.JWTClaimUser)

	response := h.Usecase.CreateServerOrder(ctx, claim, payload)
	c.JSON(response.Status, response)
}

func (h *routeHandler) ConfirmOrder(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.ConfrimOrderRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	claim := c.MustGet("token_data").(domain.JWTClaimUser)

	response := h.Usecase.ConfirmOrder(ctx, claim, payload)
	c.JSON(response.Status, response)
}

func (r *routeHandler) UploadAttachmentOrder(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimUser)

	payload := domain.UploadAttachment{}
	c.Bind(&payload)

	response := r.Usecase.UploadAttachmentOrder(ctx, claim, payload, c.Request)
	c.AbortWithStatusJSON(response.Status, response)
}
