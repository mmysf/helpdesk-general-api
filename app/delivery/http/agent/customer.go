package http_agent

import (
	"app/domain"
	"net/http"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/gin-gonic/gin"
)

func (r *routeHandler) handleCustomerRoute(prefix string) {
	// (optional). add prefix api version
	api := r.Route.Group(prefix)

	api.GET("/list", r.Middleware.AuthAgent(), r.CustomerList)
	api.GET("/detail/:id", r.Middleware.AuthAgent(), r.CustomerDetail)
	api.POST("/create", r.Middleware.AuthAgent(), r.Middleware.Role("admin"), r.CustomerCreate)
	api.PUT("/update/:id", r.Middleware.AuthAgent(), r.Middleware.Role("admin"), r.CustomerUpdate)
	api.DELETE("/delete/:id", r.Middleware.AuthAgent(), r.Middleware.Role("admin"), r.CustomerDelete)
}

func (r *routeHandler) CustomerList(c *gin.Context) {
	ctx := c.Request.Context()

	options := map[string]interface{}{
		"query": c.Request.URL.Query(),
	}

	response := r.Usecase.GetCustomerList(ctx, c.MustGet("token_data").(domain.JWTClaimAgent), options)
	c.JSON(response.Status, response)
}

func (r *routeHandler) CustomerDetail(c *gin.Context) {
	ctx := c.Request.Context()

	options := map[string]interface{}{
		"id": c.Param("id"),
	}

	response := r.Usecase.GetCustomerDetail(ctx, c.MustGet("token_data").(domain.JWTClaimAgent), options)
	c.JSON(response.Status, response)
}

func (r *routeHandler) CustomerCreate(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.CreateUserRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	options := map[string]interface{}{
		"payload": payload,
	}

	response := r.Usecase.CreateCustomer(ctx, c.MustGet("token_data").(domain.JWTClaimAgent), options)
	c.JSON(response.Status, response)
}

func (r *routeHandler) CustomerUpdate(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.CreateUserRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	options := map[string]interface{}{
		"id":      c.Param("id"),
		"payload": payload,
	}

	response := r.Usecase.UpdateCustomer(ctx, c.MustGet("token_data").(domain.JWTClaimAgent), options)
	c.JSON(response.Status, response)
}

func (r *routeHandler) CustomerDelete(c *gin.Context) {
	ctx := c.Request.Context()

	options := map[string]interface{}{
		"id": c.Param("id"),
	}

	response := r.Usecase.DeleteCustomer(ctx, c.MustGet("token_data").(domain.JWTClaimAgent), options)
	c.JSON(response.Status, response)
}
