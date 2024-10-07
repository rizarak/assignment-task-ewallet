package main

import (
	"assignment-task-ewallet/aggregator/config"
	"assignment-task-ewallet/aggregator/handler"
	"assignment-task-ewallet/aggregator/router"
	"assignment-task-ewallet/aggregator/service"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Inisialisasi router Gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Initialize service user
	userClient := config.InitUserSvc()
	if userClient == nil {
		log.Fatal("Failed to initialize User Service Client")
	}

	// Initialize wallet service
	walletClient := config.InitWalletSvc()
	if walletClient == nil {
		log.Fatal("Failed to initialize Wallet Service Client")
	}

	// Initialize wallet service
	transactionClient := config.InitTransactionSvc()
	if transactionClient == nil {
		log.Fatal("Failed to initialize Transaction Service Client")
	}

	aggregationService := service.NewAggregatorService(userClient, walletClient, transactionClient)
	aggregatorHandler := handler.NewAggregatorHandler(*aggregationService)

	router.SetupRouter(r, aggregatorHandler)

	// Jalankan server pada port 8080
	log.Println("Running main service on port 8080")
	r.Run(":8080")
}
