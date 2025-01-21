package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ppay/lib"
)

func ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		response := lib.NewResponse(c)

		head := c.GetHeader("Authorization")
		if head == "" {
			response.Unauthorized("Unauthorized", nil)
			return
		}
		token := strings.Split(head, " ")[1]

		if token == "" {
			response.Unauthorized("Unauthorized", nil)
		}

		claims, err := lib.VerifyToken(token)
		if err != nil {
			response.Unauthorized("Unauthorized", nil)
		}
		c.Set("UserId", claims.UserId)
		c.Next()
	}
}
