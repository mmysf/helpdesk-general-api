package http_webhook

import (
	"app/domain"
	"net/http"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/gin-gonic/gin"
)

func (h *routeHandler) handleWebhookRoute(prefixPath string) {
	// (optional). add prefix api version
	api := h.Route.Group(prefixPath)

	api.POST("/xendit", h.Middleware.VerifyXenditWebhookToken(), h.WebhookXendit)
}

func (h *routeHandler) WebhookXendit(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.SnapWebhookRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	options := map[string]interface{}{
		"payload": payload,
		"query":   c.Request.URL.Query(),
	}

	response := h.Usecase.HandleWebhook(ctx, options)
	c.JSON(response.Status, response)
}
