package postgres_gorm

import (
	"assignment-task-ewallet/wallet/entity"
	"assignment-task-ewallet/wallet/service"
	"context"
	"errors"
	"log"

	"gorm.io/gorm"
)

type transactionRepository struct {
	db GormDBIface
}

func NewTransactionRepository(db GormDBIface) service.ITransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) CreateTransaction(ctx context.Context, transaction *entity.Transaction) (entity.Transaction, error) {
	if err := r.db.WithContext(ctx).Create(transaction).Error; err != nil {
		log.Printf("Error creating transaction: %v\n", err)
		return entity.Transaction{}, err
	}
	return *transaction, nil
}

func (r *transactionRepository) GetTransactions(ctx context.Context, param entity.TransactionGetRequest) ([]entity.Transaction, int64, error) {
	var transactions []entity.Transaction
	var totalCount int64

	// Query to get total count
	if err := r.db.WithContext(ctx).Model(&entity.Transaction{}).Where("type = ? AND user_id = ?", param.Type, param.UserId).Count(&totalCount).Error; err != nil {
		log.Printf("Error getting total count of transactions: %v\n", err)
		return nil, 0, err
	}

	//Query to get data using pagination
	offset := (param.Page - 1) * param.Size
	if err := r.db.WithContext(ctx).Where("type = ? AND user_id = ?", param.Type, param.UserId).Limit(param.Size).Offset(offset).Find(&transactions).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transactions, 0, nil
		}
		log.Printf("Error getting all wallets: %v\n", err)
		return nil, 0, err
	}
	return transactions, totalCount, nil
}
