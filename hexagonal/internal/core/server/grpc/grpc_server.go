package grpc

import (
	"io"
	"net"
	"strconv"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"user-service/internal/core/config"
)

type GRPCServer interface {
	Start(serviceRegister func(server *grpc.Server))
	io.Closer
}

type gRPCServer struct {
	grpcServer *grpc.Server
	config     config.GrpcServerConfig
}

func NewGrpcServer(config config.GrpcServerConfig) (GRPCServer, error) {
	options, err := buildOptions(config)
	if err != nil {
		return nil, err
	}

	server := grpc.NewServer(options...)

	return &gRPCServer{
		config:     config,
		grpcServer: server,
	}, err
}

func buildOptions(config config.GrpcServerConfig) ([]grpc.ServerOption, error) {
	return []grpc.ServerOption{
		grpc.KeepaliveParams(buildKeepaliveParams(config.KeepaliveParams)),
		grpc.KeepaliveEnforcementPolicy(buildKeepalivePolicy(config.KeepalivePolicy)),
	}, nil
}

func buildKeepalivePolicy(config keepalive.EnforcementPolicy) keepalive.EnforcementPolicy {
	return keepalive.EnforcementPolicy{
		MinTime:             config.MinTime * time.Second,
		PermitWithoutStream: config.PermitWithoutStream,
	}
}

func buildKeepaliveParams(config keepalive.ServerParameters) keepalive.ServerParameters {
	return keepalive.ServerParameters{
		MaxConnectionIdle:     config.MaxConnectionIdle * time.Second,
		MaxConnectionAge:      config.MaxConnectionAge * time.Second,
		MaxConnectionAgeGrace: config.MaxConnectionAgeGrace * time.Second,
		Time:                  config.Time * time.Second,
		Timeout:               config.Timeout * time.Second,
	}
}

func (g gRPCServer) Start(serviceRegister func(server *grpc.Server)) {
	grpcListener, err := net.Listen("tcp", ":"+strconv.Itoa(int(g.config.Port)))
	if err != nil {
		zap.L().Fatal("failed to start grpc server", zap.Any("err", err))
	}

	serviceRegister(g.grpcServer)

	zap.L().Info("start grpc server success ", zap.Any("endpoint", grpcListener.Addr()))
	if err := g.grpcServer.Serve(grpcListener); err != nil {
		zap.L().Fatal("failed to grpc server serve", zap.Any("err", err))
	}
}

func (g gRPCServer) Close() error {
	zap.L().Info("close gRPC server")
	g.grpcServer.GracefulStop()
	return nil
}
