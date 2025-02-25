package middlewares

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/palashbhasme/ecommerce_microservices/inventory_service/internals/domain/models"
	"github.com/palashbhasme/ecommerce_microservices/inventory_service/utils"
)

const unauthorizedErrorMessage = "unauthorized"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": unauthorizedErrorMessage})
			c.Abort()
			return
		}

		secret := os.Getenv("JWT_SECRET")
		claims, err := utils.ParseToken(cookie, secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": unauthorizedErrorMessage})
			c.Abort()
			return
		}

		c.Set("user", claims)
		c.Next()
	}

}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user claims from context
		claims, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		// Fix: Type assert properly
		userClaims, ok := claims.(*models.Claims)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}

		// Fix: Ensure Role is properly compared
		if userClaims.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}

		c.Next()
	}
}
