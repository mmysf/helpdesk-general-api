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

func NewCustomerHandler(route *gin.RouterGroup, middleware middleware.Middleware, u usecase_member.AppUsecase) {
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
	handler.handleOrderRoute("/order")
	handler.handleHourPackageRoute("/package/hour")
	handler.handleCustomerSubscriptionRoute("/customer-subscription")
	handler.handleUserRoute("/user")
	handler.handleProjectRoute("/project")
	handler.handleCompanyRoute("/company")
	handler.handleSettingRoute("/setting")
	handler.handleTicketCategoryRoute("/ticket-category")
	handler.handleServerPackageRoute("/package/server")
	handler.handleConfigRoute("/config")
	handler.handleNotificationRoute("/notification")

}
