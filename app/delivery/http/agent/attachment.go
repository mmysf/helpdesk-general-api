package http_agent

import (
	"app/domain"

	"github.com/gin-gonic/gin"
)

func (h *routeHandler) handleAttachmentRoute(prefixPath string) {
	// (optional). add prefix api version
	api := h.Route.Group(prefixPath)

	api.POST("/upload", h.Middleware.AuthAgent(), h.UploadAttachment)
	api.GET("/detail/:id", h.Middleware.AuthAgent(), h.AttachmentDetail)
}

func (r *routeHandler) UploadAttachment(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimAgent)

	payload := domain.UploadAttachment{}
	c.Bind(&payload)

	response := r.Usecase.UploadAttachment(ctx, claim, payload, c.Request)
	c.AbortWithStatusJSON(response.Status, response)
}

func (r *routeHandler) AttachmentDetail(c *gin.Context) {
	ctx := c.Request.Context()

	attachmentId := c.Param("id")

	response := r.Usecase.GetAttachmentDetail(ctx, c.MustGet("token_data").(domain.JWTClaimAgent), attachmentId)
	c.JSON(response.Status, response)
}
