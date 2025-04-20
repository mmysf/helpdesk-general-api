package http_agent

import (
	"github.com/gin-gonic/gin"
)

func (r *routeHandler) handleConfigRoute(prefixPath string) {
	api := r.Route.Group(prefixPath)

	api.GET("", r.GetConfig)
}

func (r *routeHandler) GetConfig(c *gin.Context) {
	ctx := c.Request.Context()

	response := r.Usecase.GetConfig(ctx)
	c.JSON(response.Status, response)
}
