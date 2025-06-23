package http_superadmin

import (
	"app/domain"
	"net/http"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/gin-gonic/gin"
)

func (h *routeHandler) handleHourPackageRoute(prefixPath string) {
	api := h.Route.Group(prefixPath, h.Middleware.AuthSuperadmin())
	api.GET("/list", h.List)
	api.POST("/create", h.Create)
	api.GET("/detail/:id", h.Detail)
	api.PUT("/update/:id", h.Update)
	api.PATCH("/update-status/:id", h.UpdateStatus)
	api.DELETE("/delete/:id", h.Delete)
}

func (h *routeHandler) List(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimSuperadmin)
	query := c.Request.URL.Query()

	response := h.Usecase.GetHourPackages(ctx, claim, query)
	c.JSON(response.Status, response)
}

func (r *routeHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.HourPackageRequest{}
	err := c.Bind(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	claim := c.MustGet("token_data").(domain.JWTClaimSuperadmin)

	response := r.Usecase.CreateHourPackage(ctx, claim, payload)
	c.JSON(response.Status, response)
}

func (r *routeHandler) Detail(c *gin.Context) {
	ctx := c.Request.Context()

	productId := c.Param("id")

	response := r.Usecase.GetHourPackageDetail(ctx, productId)
	c.JSON(response.Status, response)
}

func (r *routeHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()

	productId := c.Param("id")

	payload := domain.HourPackageUpdate{}
	err := c.Bind(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	claim := c.MustGet("token_data").(domain.JWTClaimSuperadmin)

	response := r.Usecase.UpdateHourPackage(ctx, claim, productId, payload)
	c.JSON(response.Status, response)
}

func (r *routeHandler) UpdateStatus(c *gin.Context) {
	ctx := c.Request.Context()

	productId := c.Param("id")

	payload := domain.HourPackageStatusUpdate{}
	err := c.Bind(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	claim := c.MustGet("token_data").(domain.JWTClaimSuperadmin)

	response := r.Usecase.UpdateStatusHourPackage(ctx, claim, productId, payload)
	c.JSON(response.Status, response)
}

func (r *routeHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimSuperadmin)
	productId := c.Param("id")

	response := r.Usecase.DeleteHourPackage(ctx, claim, productId)
	c.JSON(response.Status, response)
}
