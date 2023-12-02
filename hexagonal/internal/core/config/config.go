package config

import (
	"google.golang.org/grpc/keepalive"
)

type HttpServerConfig struct {
	Port uint
}

type GrpcServerConfig struct {
	Port            uint32
	KeepaliveParams keepalive.ServerParameters
	KeepalivePolicy keepalive.EnforcementPolicy
}
