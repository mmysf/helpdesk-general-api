package http_webhook

import (
	"app/app/delivery/http/middleware"
	usecase_webhook "app/app/usecase/webhook"

	"github.com/gin-gonic/gin"
)

type routeHandler struct {
	Usecase    usecase_webhook.WebhookUsecase
	Route      *gin.RouterGroup
	Middleware middleware.Middleware
}

func NewWebhookRouteHandler(route *gin.RouterGroup, middleware middleware.Middleware, u usecase_webhook.WebhookUsecase) {
	handler := &routeHandler{
		Usecase:    u,
		Route:      route,
		Middleware: middleware,
	}
	handler.handleWebhookRoute("/webhook")
}
