package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ppay/internal/controllers"
	"github.com/ppay/internal/middlewares"
)

func AuthRoutes(router *gin.Engine) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", controllers.Register)
		authGroup.POST("/login", controllers.Login)
		authGroup.POST("/pin", middlewares.ValidateToken(), controllers.VerifPin)
		authGroup.GET("/pin", middlewares.ValidateToken(), controllers.CheckAvailPin)
	}
}
