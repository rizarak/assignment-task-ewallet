package config

import (
	user_service "assignment-task-ewallet/aggregator/proto/user_service/v1"
	"log"

	"google.golang.org/grpc"
)

func InitUserSvc() user_service.UserSvcClient {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return user_service.NewUserSvcClient(conn)
}
