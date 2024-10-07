package main

import (
	"assignment-task-ewallet/wallet/entity"
	grpcHandler "assignment-task-ewallet/wallet/handler/grpc"
	pbTransaction "assignment-task-ewallet/wallet/proto/transaction_service/v1"
	pbWallet "assignment-task-ewallet/wallet/proto/wallet_service/v1"
	postgres_gorm "assignment-task-ewallet/wallet/repository/pgsql_gorm"
	"assignment-task-ewallet/wallet/service"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "postgresql://postgres:P4ssw0rd@192.168.26.50:5432/latihangolang"
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		log.Fatalln(err)
	}

	err = gormDB.AutoMigrate(entity.Wallet{}, entity.Transaction{})
	if err != nil {
		fmt.Println("Failed to migrate database schema:", err)
	} else {
		fmt.Println("Database schema migrated successfully")
	}

	walletRepo := postgres_gorm.NewWalletRepository(gormDB)
	walletService := service.NewWalletService(walletRepo)
	walletHandler := grpcHandler.NewWalletHandler(walletService)

	transactionRepo := postgres_gorm.NewTransactionRepository(gormDB)
	transactionService := service.NewTransactionService(transactionRepo, walletRepo)
	transactionHandler := grpcHandler.NewTransactionHandler(transactionService)

	grpcServer := grpc.NewServer()
	pbWallet.RegisterWalletSvcServer(grpcServer, walletHandler)
	pbTransaction.RegisterTransactionSvcServer(grpcServer, transactionHandler)
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("Running grpc wallet service in port :50052")
	grpcServer.Serve(lis)
}
