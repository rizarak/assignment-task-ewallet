package service

import (
	"assignment-task-ewallet/wallet/entity"
	"context"
	"fmt"
	"log"
)

type ITransactionService interface {
	CreateTransaction(ctx context.Context, transaction *entity.TransactionRequest) (entity.Transaction, error)
	GetTransactions(ctx context.Context, transaction entity.TransactionGetRequest) (entity.TransactionResponse, error)
}

type ITransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction *entity.Transaction) (entity.Transaction, error)
	GetTransactions(ctx context.Context, transaction entity.TransactionGetRequest) ([]entity.Transaction, int64, error)
}

type transactionService struct {
	transactionRepo ITransactionRepository
	walletRepo      IWalletRepository
}

func NewTransactionService(transactionRepo ITransactionRepository, walletRepo IWalletRepository) ITransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
		walletRepo:      walletRepo,
	}
}

func (s *transactionService) CreateTransaction(ctx context.Context, transaction *entity.TransactionRequest) (entity.Transaction, error) {
	if transaction.Category == "topup" {
		//create transaction
		userId := transaction.SourceId
		dataTransaction := entity.Transaction{
			SourceId:      userId,
			DestinationId: userId,
			Amount:        transaction.Amount,
			Type:          "credit",
			Category:      transaction.Category,
			Notes:         transaction.Notes,
		}
		_, err := s.transactionRepo.CreateTransaction(ctx, &dataTransaction)
		if err != nil {
			return dataTransaction, fmt.Errorf("Failed to create transaction: %v", err)
		}
		//update wallet
		getWallet, err := s.walletRepo.GetWalletByUserId(ctx, userId)
		if err != nil {
			return dataTransaction, fmt.Errorf("Failed to get the wallet: %v", err)
		}
		balance := getWallet.Balance + dataTransaction.Amount
		dataWallet := entity.Wallet{
			Id:      getWallet.Id,
			Balance: balance,
		}
		updatedWallet, err := s.walletRepo.UpdateWallet(ctx, getWallet.Id, dataWallet)
		log.Println(updatedWallet)
		if err != nil {
			return dataTransaction, fmt.Errorf("Failed to update the wallet: %v", err)
		}
		return dataTransaction, nil
	} else if transaction.Category == "transfer" {
		//transfer from
		dataTransaction := entity.Transaction{
			SourceId:      transaction.SourceId,
			DestinationId: transaction.DestinationId,
			Amount:        transaction.Amount,
			Type:          "debt",
			Category:      transaction.Category,
		}
		_, err := s.transactionRepo.CreateTransaction(ctx, &dataTransaction)
		if err != nil {
			return dataTransaction, fmt.Errorf("Failed to create transaction: %v", err)
		}

		getWalletFrom, err := s.walletRepo.GetWalletByUserId(ctx, transaction.SourceId)
		if err != nil {
			return dataTransaction, fmt.Errorf("Failed to get the source wallet: %v", err)
		}
		balance := getWalletFrom.Balance - dataTransaction.Amount
		dataWallet := entity.Wallet{
			Id:      getWalletFrom.Id,
			Balance: balance,
		}
		updatedWallet, err := s.walletRepo.UpdateWallet(ctx, getWalletFrom.Id, dataWallet)
		log.Println(updatedWallet)
		if err != nil {
			return dataTransaction, fmt.Errorf("Failed to update the source wallet balance: %v", err)
		}

		// transfer to
		dtTransaction := entity.Transaction{
			SourceId:      transaction.SourceId,
			DestinationId: transaction.DestinationId,
			Amount:        transaction.Amount,
			Type:          "credit",
			Category:      transaction.Category,
		}
		createdTrans, err := s.transactionRepo.CreateTransaction(ctx, &dtTransaction)
		log.Println(createdTrans)
		if err != nil {
			return dtTransaction, fmt.Errorf("Failed to create transaction: %v", err)
		}

		getWalletTo, err := s.walletRepo.GetWalletByUserId(ctx, transaction.DestinationId)
		if err != nil {
			return dtTransaction, fmt.Errorf("Failed to get the destination wallet: %v", err)
		}
		balanceTo := getWalletTo.Balance + dataTransaction.Amount
		dtWallet := entity.Wallet{
			Id:      getWalletTo.Id,
			Balance: balanceTo,
		}
		updatedWalletTo, err := s.walletRepo.UpdateWallet(ctx, getWalletTo.Id, dtWallet)
		log.Println(updatedWalletTo)
		if err != nil {
			return dtTransaction, fmt.Errorf("Failed to update the destination wallet balance: %v", err)
		}
		return dtTransaction, nil
	} else {
		return entity.Transaction{}, fmt.Errorf("Unknown Transaction Category")
	}
}

func (s *transactionService) GetTransactions(ctx context.Context, param entity.TransactionGetRequest) (entity.TransactionResponse, error) {
	transactions, totalCount, err := s.transactionRepo.GetTransactions(ctx, param)

	if err != nil {
		return entity.TransactionResponse{}, fmt.Errorf("Failed to get transactions: %v", err)
	}

	return entity.TransactionResponse{
		Transaction: transactions,
		Count:       totalCount,
	}, nil
}
