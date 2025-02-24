package http_member

import (
	"app/app/delivery/http/middleware"
	usecase_member "app/app/usecase/member"

	"github.com/gin-gonic/gin"
)

type routeHandler struct {
	Usecase    usecase_member.AppUsecase
	Route      *gin.RouterGroup
	Middleware middleware.Middleware
}

func NewRouteHandler(route *gin.RouterGroup, middleware middleware.Middleware, u usecase_member.AppUsecase) {
	handler := &routeHandler{
		Usecase:    u,
		Route:      route,
		Middleware: middleware,
	}

	handler.handleAuthRoute("/auth")
}
