package main

import (
	"assignment-task-ewallet/user/entity"
	grpc_handler "assignment-task-ewallet/user/handler/grpc"
	pb "assignment-task-ewallet/user/proto/user_service/v1"
	pgsql_gorm "assignment-task-ewallet/user/repository/pgsql_gorm"
	"assignment-task-ewallet/user/service"
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

	err = gormDB.AutoMigrate(entity.User{})
	if err != nil {
		fmt.Println("Failed to migrate database schema:", err)
	} else {
		fmt.Println("Database schema migrated successfully")
	}

	userRepo := pgsql_gorm.NewUserRepository(gormDB)
	userService := service.NewUserService(userRepo)
	userHandler := grpc_handler.NewUserHandler(userService)

	grpcServer := grpc.NewServer()
	pb.RegisterUserSvcServer(grpcServer, userHandler)

	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Running grpc user service in port :50051")
	grpcServer.Serve(lis)
}
