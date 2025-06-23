package http_agent

import (
	"app/app/delivery/http/middleware"
	usecase_agent "app/app/usecase/agent"

	"github.com/gin-gonic/gin"
)

type routeHandler struct {
	Usecase    usecase_agent.AgentUsecase
	Route      *gin.RouterGroup
	Middleware middleware.Middleware
}

func NewAgentRouteHandler(route *gin.RouterGroup, middleware middleware.Middleware, u usecase_agent.AgentUsecase) {
	handler := &routeHandler{
		Usecase:    u,
		Route:      route,
		Middleware: middleware,
	}
	handler.handleAuthRoute("/auth")
	handler.handleTicketRoute("/ticket")
	handler.handleProductRoute("/product")
	handler.handleAttachmentRoute("/attachment")
	handler.handleDashboardRoute("/dashboard")
	handler.handleTicketTimelogsRoute("/timelogs")
	handler.handleCompanyProductRoute("/company-product")
	handler.handleCustomerRoute("/customer")
	handler.handleCompanyRoute("/company")
	handler.handleSettingRoute("/setting")
	handler.handleUserRoute("/user")
	handler.handleConfigRoute("/config")
	handler.handleTicketCategoryRoute("/ticket-category")
}
