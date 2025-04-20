package http_member

import (
	"app/domain"
	"net/http"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/gin-gonic/gin"
)

func (h *routeHandler) handleProjectRoute(prefixPath string) {
	// (optional). add prefix api version
	api := h.Route.Group(prefixPath)

	api.GET("/list", h.Middleware.AuthCustomer(), h.ProjectList)
	api.GET("/detail/:id", h.Middleware.AuthCustomer(), h.ProjectDetail)
	api.POST("/create", h.Middleware.AuthCustomer(), h.ProjectCreate)
	api.PUT("/update/:id", h.Middleware.AuthCustomer(), h.ProjectUpdate)
	api.DELETE("/delete/:id", h.Middleware.AuthCustomer(), h.ProjectDelete)
}

func (r *routeHandler) ProjectList(c *gin.Context) {
	ctx := c.Request.Context()

	claim := c.MustGet("token_data").(domain.JWTClaimUser)
	query := c.Request.URL.Query()

	response := r.Usecase.GetProjectList(ctx, claim, query)
	c.JSON(response.Status, response)
}

func (r *routeHandler) ProjectDetail(c *gin.Context) {
	ctx := c.Request.Context()

	projectID := c.Param("id")

	response := r.Usecase.GetProjectDetail(ctx, c.MustGet("token_data").(domain.JWTClaimUser), projectID)
	c.JSON(response.Status, response)
}

func (r *routeHandler) ProjectCreate(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.ProjectRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	claim := c.MustGet("token_data").(domain.JWTClaimUser)

	response := r.Usecase.CreateProject(ctx, claim, payload)
	c.JSON(response.Status, response)
}

func (r *routeHandler) ProjectUpdate(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.ProjectRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	claim := c.MustGet("token_data").(domain.JWTClaimUser)
	projectID := c.Param("id")

	response := r.Usecase.UpdateProject(ctx, claim, projectID, payload)
	c.JSON(response.Status, response)
}

func (r *routeHandler) ProjectDelete(c *gin.Context) {
	ctx := c.Request.Context()

	projectID := c.Param("id")

	response := r.Usecase.DeleteProject(ctx, c.MustGet("token_data").(domain.JWTClaimUser), projectID)
	c.JSON(response.Status, response)
}
