package http_superadmin

import (
	"app/domain"
	"net/http"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/gin-gonic/gin"
)

func (h *routeHandler) handleServerPackageRoute(prefixPath string) {
	api := h.Route.Group(prefixPath, h.Middleware.AuthSuperadmin())
	api.GET("/list", h.ServerPackageList)
	api.POST("/create", h.ServerPackageCreate)
	api.GET("/detail/:id", h.ServerPackageDetail)
	api.PUT("/update/:id", h.ServerPackageUpdate)
	api.PATCH("/update-status/:id", h.ServerPackageUpdateStatus)
	api.DELETE("/delete/:id", h.ServerPackageDelete)
}

func (h *routeHandler) ServerPackageList(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimSuperadmin)
	query := c.Request.URL.Query()

	response := h.Usecase.GetServerPackageList(ctx, claim, query)
	c.JSON(response.Status, response)
}

func (r *routeHandler) ServerPackageCreate(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.ServerPackageRequest{}
	err := c.Bind(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	claim := c.MustGet("token_data").(domain.JWTClaimSuperadmin)

	response := r.Usecase.CreateServerPackage(ctx, claim, payload)
	c.JSON(response.Status, response)
}

func (r *routeHandler) ServerPackageDetail(c *gin.Context) {
	ctx := c.Request.Context()

	productId := c.Param("id")

	response := r.Usecase.GetServerPackageDetail(ctx, productId)
	c.JSON(response.Status, response)
}

func (r *routeHandler) ServerPackageUpdate(c *gin.Context) {
	ctx := c.Request.Context()

	productId := c.Param("id")

	payload := domain.ServerPackageUpdate{}
	err := c.Bind(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	claim := c.MustGet("token_data").(domain.JWTClaimSuperadmin)

	response := r.Usecase.UpdateServerPackage(ctx, claim, productId, payload)
	c.JSON(response.Status, response)
}

func (r *routeHandler) ServerPackageUpdateStatus(c *gin.Context) {
	ctx := c.Request.Context()

	productId := c.Param("id")

	payload := domain.ServerPackageStatusUpdate{}
	err := c.Bind(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	claim := c.MustGet("token_data").(domain.JWTClaimSuperadmin)

	response := r.Usecase.UpdateStatusServerPackage(ctx, claim, productId, payload)
	c.JSON(response.Status, response)
}

func (r *routeHandler) ServerPackageDelete(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimSuperadmin)
	productId := c.Param("id")

	response := r.Usecase.DeleteServerPackage(ctx, claim, productId)
	c.JSON(response.Status, response)
}
