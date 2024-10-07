package grpc

import (
	"assignment-task-ewallet/wallet/entity"
	pb "assignment-task-ewallet/wallet/proto/transaction_service/v1"
	"assignment-task-ewallet/wallet/service"
	"context"
	"fmt"
	"log"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type TransactionHandler struct {
	pb.UnimplementedTransactionSvcServer
	transactionService service.ITransactionService
}

func NewTransactionHandler(transactionService service.ITransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

func (u *TransactionHandler) CreateTransaction(ctx context.Context, req *pb.CreateTransactionRequest) (*pb.MutationTransResponse, error) {
	log.Println(req)
	createdTransaction, err := u.transactionService.CreateTransaction(ctx, &entity.TransactionRequest{
		SourceId:      int(req.GetSourceId()),
		DestinationId: int(req.GetDestinationId()),
		Type:          string(req.GetType()),
		Amount:        float64(req.GetAmount()),
		Notes:         string(req.GetNotes()),
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pb.MutationTransResponse{
		Message: fmt.Sprintf("Success created Transaction with ID %d", createdTransaction.Id),
	}, nil
}

func (u *TransactionHandler) GetTransactions(ctx context.Context, req *pb.GetTransactionRequest) (*pb.GetTransactionResponse, error) {
	transactions, err := u.transactionService.GetTransactions(ctx, entity.TransactionGetRequest{
		Type:   req.GetType(),
		UserId: int(req.GetUserId()),
		Size:   int(req.GetPageSize()),
		Page:   int(req.GetPage()),
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var transactionProto []*pb.Transaction
	for _, transaction := range transactions.Transaction {
		transactionProto = append(transactionProto, &pb.Transaction{
			Id:        int32(transaction.Id),
			UserId:    int32(transaction.SourceId),
			Amount:    float32(transaction.Amount),
			Type:      string(transaction.Type),
			Category:  string(transaction.Category),
			CreatedAt: timestamppb.New(transaction.CreatedAt),
			UpdatedAt: timestamppb.New(transaction.UpdatedAt),
		})
	}

	return &pb.GetTransactionResponse{
		Transactions: transactionProto,
		Count:        int32(transactions.Count),
	}, nil
}
