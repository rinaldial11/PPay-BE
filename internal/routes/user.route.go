package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ppay/internal/controllers"
	"github.com/ppay/internal/middlewares"
)

func UserRoutes(router *gin.Engine) {
	userGroup := router.Group("/users")
	{
		userGroup.GET("", controllers.GetUsers)
		userGroup.GET("/:id", middlewares.ValidateToken(), controllers.GetUserByID)
		userGroup.GET("/balance/:id", middlewares.ValidateToken(), controllers.GetBalance)
		userGroup.POST("", controllers.CreateUser)
		userGroup.PATCH("/:id", middlewares.ValidateToken(), controllers.UpdateUser)
		userGroup.DELETE("/:id", middlewares.ValidateToken(), controllers.DeleteUser)
		userGroup.GET("/transfer/:id", middlewares.ValidateToken(), controllers.GetTargetUserById)
	}
}
