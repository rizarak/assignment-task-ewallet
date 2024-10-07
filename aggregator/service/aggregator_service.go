package service

import (
	"assignment-task-ewallet/aggregator/entity"
	transaction_service "assignment-task-ewallet/aggregator/proto/transaction_service/v1"
	user_service "assignment-task-ewallet/aggregator/proto/user_service/v1"
	wallet_service "assignment-task-ewallet/aggregator/proto/wallet_service/v1"
	"context"
	"log"
	"strconv"
	"time"
)

type IAggregatorService interface {
	GetUser(ctx context.Context, id int) (entity.UserResponse, error)
	CreateUser(ctx context.Context, request entity.UserCreateRequest) (entity.UserResponse, error)
	TopupTransaction(ctx context.Context, request entity.TransactionRequest) (entity.TransactionResponse, error)
	TransferTransaction(ctx context.Context, request entity.TransactionRequest) (entity.TransactionResponse, error)
	GetTransactions(ctx context.Context, param entity.TransactionGetRequest) (entity.TransactionGetResponseWithPagination, error)
}

type AggregatorService struct {
	userService        user_service.UserSvcClient
	walletService      wallet_service.WalletSvcClient
	transactionService transaction_service.TransactionSvcClient
}

func NewAggregatorService(userService user_service.UserSvcClient, walletService wallet_service.WalletSvcClient, transactionService transaction_service.TransactionSvcClient) *AggregatorService {
	return &AggregatorService{
		userService:        userService,
		walletService:      walletService,
		transactionService: transactionService,
	}
}

func (svc *AggregatorService) TopupTransaction(ctx context.Context, request entity.TopUpRequest) (entity.TransactionResponse, error) {
	transactionResp, err := svc.transactionService.CreateTransaction(ctx, &transaction_service.CreateTransactionRequest{
		SourceId:      int32(request.Id),
		DestinationId: int32(request.Id),
		Type:          "credit",
		Category:      "topup",
		Amount:        float32(request.Amount),
		Notes:         request.Notes,
	})

	if err != nil {
		return entity.TransactionResponse{}, err
	}

	return entity.TransactionResponse{Message: transactionResp.Message}, nil
}

func (svc *AggregatorService) TransferTransaction(ctx context.Context, request entity.TransferRequest) (entity.TransactionResponse, error) {
	//Get Balance User sender
	walletResp, err := svc.walletService.GetWalletByUserId(ctx, &wallet_service.GetWalletByUserIdRequest{UserId: int32(request.SourceId)})
	if err != nil {
		return entity.TransactionResponse{Message: "Wallet Not Fount"}, err
	}

	if walletResp.Wallet.Balance < float32(request.Amount) {
		return entity.TransactionResponse{Message: "The balance is not enough. Your balance: " + strconv.FormatFloat(float64(walletResp.Wallet.Balance), 'f', 2, 64)}, err
	}

	//create transaction
	transactionResp, err := svc.transactionService.CreateTransaction(ctx, &transaction_service.CreateTransactionRequest{
		SourceId:      int32(request.SourceId),
		DestinationId: int32(request.DestinationId),
		Type:          "debt",
		Category:      "transfer",
		Amount:        float32(request.Amount),
		Notes:         request.Notes,
	})

	if err != nil {
		return entity.TransactionResponse{}, err
	}

	return entity.TransactionResponse{Message: transactionResp.Message}, nil
}

func (svc *AggregatorService) GetTransactions(ctx context.Context, request entity.TransactionGetRequest) (entity.TransactionGetResponseWithPagination, error) {
	transactionResp, err := svc.transactionService.GetTransactions(ctx, &transaction_service.GetTransactionRequest{
		Type:     request.Type,
		Category: request.Category,
		UserId:   int32(request.UserId),
		PageSize: int32(request.Size),
		Page:     int32(request.Page),
	})

	if err != nil {
		return entity.TransactionGetResponseWithPagination{}, err
	}
	var transactionsWithUser []entity.TransactionGetResponse
	for _, tx := range transactionResp.Transactions {
		_, err := svc.userService.GetUserById(ctx, &user_service.GetUserByIdRequest{Id: tx.UserId})
		if err != nil {
			return entity.TransactionGetResponseWithPagination{}, err
		}

		var createdAt, updatedAt time.Time
		if tx.CreatedAt != nil {
			createdAt = tx.CreatedAt.AsTime()
		}
		if tx.UpdatedAt != nil {
			updatedAt = tx.UpdatedAt.AsTime()
		}

		transactionsWithUser = append(transactionsWithUser, entity.TransactionGetResponse{
			Id:        int(tx.Id),
			UserId:    int(tx.UserId),
			Amount:    float64(tx.Amount),
			Type:      tx.Type,
			Category:  tx.Category,
			Notes:     tx.Notes,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
	}

	return entity.TransactionGetResponseWithPagination{
		Data: transactionsWithUser,
		Pagination: entity.Pagination{
			TotalData: int(transactionResp.Count),
			TotalPage: (int(transactionResp.Count) + request.Size - 1) / request.Size,
			PageSize:  request.Size,
			Page:      request.Page,
		},
	}, nil
}

func (svc *AggregatorService) GetUser(ctx context.Context, id int) (entity.UserResponse, error) {
	userId := int32(id)
	userResp, err := svc.userService.GetUserById(ctx, &user_service.GetUserByIdRequest{Id: userId})
	if err != nil {
		return entity.UserResponse{}, err
	}
	log.Println(userId)
	walletResp, err := svc.walletService.GetWalletByUserId(ctx, &wallet_service.GetWalletByUserIdRequest{UserId: userId})
	if err != nil {
		return entity.UserResponse{}, err
	}

	userWallet := entity.Wallet{
		Id:        walletResp.Wallet.Id,
		UserId:    walletResp.Wallet.UserId,
		Name:      walletResp.Wallet.Name,
		Balance:   walletResp.Wallet.Balance,
		CreatedAt: walletResp.Wallet.CreatedAt.AsTime(),
		UpdatedAt: walletResp.Wallet.UpdatedAt.AsTime(),
	}

	user := entity.UserResponse{
		Id:        userResp.User.Id,
		Name:      userResp.User.Name,
		Email:     userResp.User.Email,
		Password:  userResp.User.Password,
		CreatedAt: userResp.User.CreatedAt.AsTime(),
		UpdatedAt: userResp.User.UpdatedAt.AsTime(),
		Wallet:    userWallet,
	}
	return user, nil
}

func (svc *AggregatorService) CreateUser(ctx context.Context, request entity.UserCreateRequest) (entity.UserResponse, error) {
	//create user
	createdUser, err := svc.userService.CreateUser(ctx, &user_service.CreateUserRequest{
		Name:  request.Name,
		Email: request.Email,
	})

	if err != nil {
		return entity.UserResponse{}, err
	}

	//create wallet for balance
	_, err = svc.walletService.CreateWallet(ctx, &wallet_service.CreateWalletRequest{
		UserId:  createdUser.Id,
		Balance: 0,
	})
	if err != nil {
		return entity.UserResponse{}, err
	}

	return entity.UserResponse{
		Id: createdUser.Id,
	}, nil
}
