package main

import (
	"log"

	"go.uber.org/zap"
	googleGrpc "google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	proto "user-service/api"
	grpcCtrl "user-service/internal/controller/grpc"
	"user-service/internal/core/config"
	"user-service/internal/core/server"
	"user-service/internal/core/server/grpc"
	"user-service/internal/core/service"
	infraConf "user-service/internal/infra/config"
	"user-service/internal/infra/repository"
)

func main() {
	// Initialize logger
	logger, _ := zap.NewProduction()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	// Initialize the database connection
	db, err := repository.NewDB(
		infraConf.DatabaseConfig{
			Driver:                  "mysql",
			Url:                     "user:password@tcp(127.0.0.1:3306)/your_database_name?charset=utf8mb4&parseTime=true&loc=UTC&tls=false&readTimeout=3s&writeTimeout=3s&timeout=3s&clientFoundRows=true",
			ConnMaxLifetimeInMinute: 3,
			MaxOpenConns:            10,
			MaxIdleConns:            1,
		},
	)
	if err != nil {
		log.Fatalf("failed to new database err=%s\n", err.Error())
	}

	// Create the UserRepository
	userRepo := repository.NewUserRepository(db)

	// Create the UserService
	userService := service.NewUserService(userRepo)

	// Create the UserController
	userController := grpcCtrl.NewUserController(userService)

	// Create the gRPC server
	grpcServer, err := grpc.NewGrpcServer(
		config.GrpcServerConfig{
			Port: 9090,
			KeepaliveParams: keepalive.ServerParameters{
				MaxConnectionIdle:     100,
				MaxConnectionAge:      7200,
				MaxConnectionAgeGrace: 60,
				Time:                  10,
				Timeout:               3,
			},
			KeepalivePolicy: keepalive.EnforcementPolicy{
				MinTime:             10,
				PermitWithoutStream: true,
			},
		},
	)
	if err != nil {
		log.Fatalf("failed to new grpc server err=%s\n", err.Error())
	}

	// Start the gRPC server
	go grpcServer.Start(
		func(server *googleGrpc.Server) {
			proto.RegisterUserServiceServer(server, userController)
		},
	)

	// Add shutdown hook to trigger closer resources of service
	server.AddShutdownHook(grpcServer, db)
}
