package http_member

import (
	"app/domain"
	"net/http"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/gin-gonic/gin"
)

func (r *routeHandler) handleUserRoute(prefix string) {
	// (optional). add prefix api version
	api := r.Route.Group(prefix)

	api.GET("/list", r.Middleware.AuthCustomer(), r.Middleware.Role("admin"), r.UserList)
	api.GET("/detail/:id", r.Middleware.AuthCustomer(), r.Middleware.Role("admin"), r.UserDetail)
	api.POST("/create", r.Middleware.AuthCustomer(), r.Middleware.Role("admin"), r.UserCreate)
	api.PUT("/update/:id", r.Middleware.AuthCustomer(), r.Middleware.Role("admin"), r.UserUpdate)
	api.DELETE("/delete/:id", r.Middleware.AuthCustomer(), r.Middleware.Role("admin"), r.UserDelete)
}

func (r *routeHandler) UserList(c *gin.Context) {
	ctx := c.Request.Context()

	options := map[string]interface{}{
		"query": c.Request.URL.Query(),
	}

	response := r.Usecase.GetUserList(ctx, c.MustGet("token_data").(domain.JWTClaimUser), options)
	c.JSON(response.Status, response)
}

func (r *routeHandler) UserDetail(c *gin.Context) {
	ctx := c.Request.Context()

	options := map[string]interface{}{
		"id": c.Param("id"),
	}

	response := r.Usecase.GetUserDetail(ctx, c.MustGet("token_data").(domain.JWTClaimUser), options)
	c.JSON(response.Status, response)
}

func (r *routeHandler) UserCreate(c *gin.Context) {
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

	response := r.Usecase.CreateUser(ctx, c.MustGet("token_data").(domain.JWTClaimUser), options)
	c.JSON(response.Status, response)
}

func (r *routeHandler) UserUpdate(c *gin.Context) {
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

	response := r.Usecase.UpdateUser(ctx, c.MustGet("token_data").(domain.JWTClaimUser), options)
	c.JSON(response.Status, response)
}

func (r *routeHandler) UserDelete(c *gin.Context) {
	ctx := c.Request.Context()

	options := map[string]interface{}{
		"id": c.Param("id"),
	}

	response := r.Usecase.DeleteUser(ctx, c.MustGet("token_data").(domain.JWTClaimUser), options)
	c.JSON(response.Status, response)
}
