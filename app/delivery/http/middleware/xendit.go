package middleware

import (
	"net/http"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/gin-gonic/gin"
)

func (m *appMiddleware) VerifyXenditWebhookToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		callbackToken := c.GetHeader("X-Callback-Token")
		if callbackToken == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				response.Error(http.StatusUnauthorized, "Unauthorized: Callback Token header is required"),
			)
			return
		}

		// verify webhook token
		if callbackToken != m.xenditWebhookToken {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				response.Error(http.StatusUnauthorized, "Unauthorized: Invalid Callback Token"),
			)
			return
		}

		c.Next()
	}
}
