package config

import (
	transaction_service "assignment-task-ewallet/aggregator/proto/transaction_service/v1"
	wallet_service "assignment-task-ewallet/aggregator/proto/wallet_service/v1"
	"log"

	"google.golang.org/grpc"
)

func InitWalletSvc() wallet_service.WalletSvcClient {
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return wallet_service.NewWalletSvcClient(conn)
}

func InitTransactionSvc() transaction_service.TransactionSvcClient {
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return transaction_service.NewTransactionSvcClient(conn)
}
