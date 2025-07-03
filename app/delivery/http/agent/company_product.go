package http_agent

func (h *routeHandler) handleCompanyProductRoute(prefixPath string) {
	// (optional). add prefix api version
	// api := h.Route.Group(prefixPath)

	// api.GET("/list", h.Middleware.AuthAgent(), h.CompanyProductList)
	// api.GET("/detail/:id", h.Middleware.AuthAgent(), h.CompanyProductDetail)
	// api.POST("/upload-logo", h.Middleware.AuthAgent(), h.UploadCompanyProductLogo)
	// api.POST("/create", h.Middleware.AuthAgent(), h.Middleware.Role("admin"), h.CompanyProductCreate)
	// api.PUT("/update/:id", h.Middleware.AuthAgent(), h.Middleware.Role("admin"), h.CompanyProductUpdate)
	// api.DELETE("/delete/:id", h.Middleware.AuthAgent(), h.Middleware.Role("admin"), h.CompanyProductDelete)
}

// func (r *routeHandler) CompanyProductList(c *gin.Context) {
// 	ctx := c.Request.Context()

// 	options := map[string]interface{}{
// 		"query": c.Request.URL.Query(),
// 	}

// 	response := r.Usecase.GetCompanyProductList(ctx, c.MustGet("token_data").(domain.JWTClaimAgent), options)
// 	c.JSON(response.Status, response)
// }

// func (r *routeHandler) CompanyProductDetail(c *gin.Context) {
// 	ctx := c.Request.Context()

// 	options := map[string]interface{}{
// 		"id": c.Param("id"),
// 	}

// 	response := r.Usecase.GetCompanyProductDetail(ctx, c.MustGet("token_data").(domain.JWTClaimAgent), options)
// 	c.JSON(response.Status, response)
// }

// func (r *routeHandler) CompanyProductCreate(c *gin.Context) {
// 	ctx := c.Request.Context()

// 	payload := domain.CreateCompanyProductRequest{}
// 	err := c.ShouldBindJSON(&payload)
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
// 		return
// 	}

// 	options := map[string]interface{}{
// 		"payload": payload,
// 	}

// 	response := r.Usecase.CreateCompanyProduct(ctx, c.MustGet("token_data").(domain.JWTClaimAgent), options)
// 	c.JSON(response.Status, response)
// }

// func (r *routeHandler) CompanyProductUpdate(c *gin.Context) {
// 	ctx := c.Request.Context()

// 	payload := domain.CreateCompanyProductRequest{}
// 	err := c.ShouldBindJSON(&payload)
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
// 		return
// 	}

// 	options := map[string]interface{}{
// 		"id":      c.Param("id"),
// 		"payload": payload,
// 	}

// 	response := r.Usecase.UpdateCompanyProduct(ctx, c.MustGet("token_data").(domain.JWTClaimAgent), options)
// 	c.JSON(response.Status, response)
// }

// func (r *routeHandler) CompanyProductDelete(c *gin.Context) {
// 	ctx := c.Request.Context()

// 	options := map[string]interface{}{
// 		"id": c.Param("id"),
// 	}

// 	response := r.Usecase.DeleteCompanyProduct(ctx, c.MustGet("token_data").(domain.JWTClaimAgent), options)
// 	c.JSON(response.Status, response)
// }

// func (r *routeHandler) UploadCompanyProductLogo(c *gin.Context) {
// 	ctx := c.Request.Context()

// 	claim := c.MustGet("token_data").(domain.JWTClaimAgent)

// 	payload := domain.UploadAttachment{}
// 	c.Bind(&payload)

// 	response := r.Usecase.UploadCompanyProductLogo(ctx, claim, payload, c.Request)
// 	c.AbortWithStatusJSON(response.Status, response)
// }
