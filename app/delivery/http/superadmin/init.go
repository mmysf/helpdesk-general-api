package http_superadmin

import (
	"app/app/delivery/http/middleware"
	usecase_superadmin "app/app/usecase/superadmin"

	"github.com/gin-gonic/gin"
)

type routeHandler struct {
	Usecase    usecase_superadmin.SuperadminUsecase
	Route      *gin.RouterGroup
	Middleware middleware.Middleware
}

func NewSuperadminRouteHandler(route *gin.RouterGroup, middleware middleware.Middleware, u usecase_superadmin.SuperadminUsecase) {

	handler := &routeHandler{
		Usecase:    u,
		Route:      route,
		Middleware: middleware,
	}
	handler.handleAuthRoute("/auth")
	handler.handleDashboardRoute("/dashboard")
	handler.handleTicketRoute("/ticket")
	handler.handleTicketCommentRoute("/ticket-comment")
	handler.handleOrderRoute("/order")
	handler.handleCustomerRoute("/customer")
	handler.handleHourPackageRoute("/package/hour")
	handler.handleAgentRoute("/agent")
	handler.handleCompanyRoute("/company")
	handler.handleCompanyProductRoute("/company-product")
	handler.handleConfigRoute("/config")
	handler.handleServerPackageRoute("/package/server")
}
