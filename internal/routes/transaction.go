package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ppay/internal/controllers"
	"github.com/ppay/internal/middlewares"
)

func TransactionRoute(router *gin.Engine) {
	transaction := router.Group("/transaction")
	{
		transaction.GET("/history", middlewares.ValidateToken(), controllers.GetTransactionHistory)
		transaction.GET("/expense", middlewares.ValidateToken(), controllers.GetUserExpenses)
		transaction.GET("/income", middlewares.ValidateToken(), controllers.GetUserIncome)
		transaction.GET("/payment", middlewares.ValidateToken(), controllers.GetPaymentMethod)
	}
}
