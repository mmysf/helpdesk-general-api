package http_superadmin

import (
	"app/domain"
	"net/http"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/gin-gonic/gin"
)

func (h *routeHandler) handleOrderRoute(prefixPath string) {
	// (optional). add prefix api version
	api := h.Route.Group(prefixPath)

	api.GET("/list", h.Middleware.AuthSuperadmin(), h.OrderList)
	api.GET("/detail/:id", h.Middleware.AuthSuperadmin(), h.OrderDetail)
	api.PATCH("/update-manual-payment/:id", h.Middleware.AuthSuperadmin(), h.UpdateManualPayment)
	api.POST("/upload-attachment", h.Middleware.AuthSuperadmin(), h.UploadAttachmentOrder)
}

func (h *routeHandler) OrderList(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimSuperadmin)
	query := c.Request.URL.Query()

	response := h.Usecase.GetOrderList(ctx, claim, query)
	c.JSON(response.Status, response)
}

func (h *routeHandler) OrderDetail(c *gin.Context) {
	ctx := c.Request.Context()

	orderId := c.Param("id")

	response := h.Usecase.GetOrderDetail(ctx, orderId)
	c.JSON(response.Status, response)
}

func (h *routeHandler) UpdateManualPayment(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimSuperadmin)
	orderId := c.Param("id")
	payload := domain.UpdateManualPaymentRequest{}

	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	response := h.Usecase.UpdateManualPayment(ctx, claim, orderId, payload)
	c.JSON(response.Status, response)
}

func (r *routeHandler) UploadAttachmentOrder(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimSuperadmin)

	payload := domain.UploadAttachment{}
	c.Bind(&payload)

	response := r.Usecase.UploadAttachmentOrder(ctx, claim, payload, c.Request)
	c.AbortWithStatusJSON(response.Status, response)
}
