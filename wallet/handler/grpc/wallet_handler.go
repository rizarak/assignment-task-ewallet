package grpc

import (
	"assignment-task-ewallet/wallet/entity"
	pb "assignment-task-ewallet/wallet/proto/wallet_service/v1"
	"assignment-task-ewallet/wallet/service"
	"context"
	"fmt"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type WalletHandler struct {
	pb.UnimplementedWalletSvcServer
	walletService service.IWalletService
}

func NewWalletHandler(walletService service.IWalletService) *WalletHandler {
	return &WalletHandler{
		walletService: walletService,
	}
}

func (u *WalletHandler) GetWallets(ctx context.Context, _ *emptypb.Empty) (*pb.GetWalletsResponse, error) {
	wallets, err := u.walletService.GetAllWallets(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var walletsProto []*pb.Wallet
	for _, wallet := range wallets {
		walletsProto = append(walletsProto, &pb.Wallet{
			Id:        int32(wallet.Id),
			UserId:    int32(wallet.UserId),
			Name:      string(wallet.Name),
			Balance:   float32(wallet.Balance),
			CreatedAt: timestamppb.New(wallet.CreatedAt),
			UpdatedAt: timestamppb.New(wallet.UpdatedAt),
		})
	}

	return &pb.GetWalletsResponse{
		Wallets: walletsProto,
	}, nil
}

func (u *WalletHandler) GetWalletByUserId(ctx context.Context, req *pb.GetWalletByUserIdRequest) (*pb.GetWalletByUserIdResponse, error) {
	wallet, err := u.walletService.GetWalletByUserId(ctx, int(req.GetUserId()))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	res := &pb.GetWalletByUserIdResponse{
		Wallet: &pb.Wallet{
			Id:        int32(wallet.Id),
			UserId:    int32(wallet.UserId),
			Name:      string(wallet.Name),
			Balance:   float32(wallet.Balance),
			CreatedAt: timestamppb.New(wallet.CreatedAt),
			UpdatedAt: timestamppb.New(wallet.UpdatedAt),
		},
	}
	return res, nil
}

func (u *WalletHandler) CreateWallet(ctx context.Context, req *pb.CreateWalletRequest) (*pb.MutationResponse, error) {
	createdWallet, err := u.walletService.CreateWallet(ctx, &entity.Wallet{
		UserId:  int(req.GetUserId()),
		Balance: float64(req.GetBalance()),
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pb.MutationResponse{
		Message: fmt.Sprintf("Success created wallet with ID %d", createdWallet.Id),
	}, nil
}
