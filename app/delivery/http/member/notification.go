package http_member

import (
	"app/domain"

	"github.com/gin-gonic/gin"
)

func (h *routeHandler) handleNotificationRoute(prefixPath string) {
	// (optional). add prefix api version
	api := h.Route.Group(prefixPath)

	api.GET("/list", h.Middleware.AuthCustomer(), h.NotificationList)
	api.GET("/detail/:id", h.Middleware.AuthCustomer(), h.NotificationDetail)
	api.POST("/read-all", h.Middleware.AuthCustomer(), h.NotificationReadAll)
	api.GET("/count", h.Middleware.AuthCustomer(), h.NotificationCount)
}

func (h *routeHandler) NotificationList(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimUser)
	response := h.Usecase.GetNotificationList(ctx, claim, c.Request.URL.Query())
	c.JSON(response.Status, response)
}

func (h *routeHandler) NotificationDetail(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimUser)
	id := c.Param("id")
	response := h.Usecase.GetNotificationDetail(ctx, claim, id)
	c.JSON(response.Status, response)
}

func (h *routeHandler) NotificationReadAll(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimUser)

	response := h.Usecase.ReadAllNotification(ctx, claim)
	c.JSON(response.Status, response)
}

func (h *routeHandler) NotificationCount(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimUser)
	response := h.Usecase.GetNotificationCount(ctx, claim)
	c.JSON(response.Status, response)
}
