package router

import (
	"assignment-task-ewallet/aggregator/handler"
	"assignment-task-ewallet/aggregator/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine, aggregator handler.IAggregatorHandler) {
	userEndpoint := r.Group("/user")
	userEndpoint.Use(middleware.AuthMiddleware())
	userEndpoint.GET("/:userId", aggregator.GetUser)
	userEndpoint.POST("/", aggregator.CreateUser)

	transactionEndpoint := r.Group("/transaction")
	transactionEndpoint.Use(middleware.AuthMiddleware())
	transactionEndpoint.GET("/", aggregator.GetTransactions)
	transactionEndpoint.POST("/topup", aggregator.TopupTransaction)
	transactionEndpoint.POST("/transfer", aggregator.TransferTransaction)
}
