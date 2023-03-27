package corvus

import (
	"fmt"
	"github.com/andibalo/ramein/corvus/internal/config"
	"github.com/andibalo/ramein/corvus/internal/proto"
	"github.com/andibalo/ramein/corvus/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type GRPCServer struct {
	listener   net.Listener
	grpcServer *grpc.Server
}

func NewGRPCServer(cfg config.Config) *GRPCServer {

	lis, err := net.Listen("tcp", cfg.AppAddress())
	if err != nil {
		cfg.Logger().Error(fmt.Sprintf("failed to listen at port %v", cfg.AppAddress()), zap.Error(err))
		panic(fmt.Sprintf("failed to listen at port %v", cfg.AppAddress()))
	}

	s := grpc.NewServer()

	healthCheckService := service.NewHealthCheckService()

	proto.RegisterHealthCheckServer(s, healthCheckService)

	return &GRPCServer{
		listener:   lis,
		grpcServer: s,
	}
}

func (s *GRPCServer) Start() error {
	if err := s.grpcServer.Serve(s.listener); err != nil {
		return err
	}

	return nil
}
