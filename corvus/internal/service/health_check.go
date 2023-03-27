package service

import (
	"context"
	"github.com/andibalo/ramein/corvus/internal/proto"
	"github.com/golang/protobuf/ptypes/empty"
)

type HealthCheckService struct {
	proto.HealthCheckServer
}

func NewHealthCheckService() *HealthCheckService {
	return &HealthCheckService{}
}

func (s *HealthCheckService) HealthCheck(ctx context.Context, void *empty.Empty) (*proto.HealthCheckResponse, error) {

	return &proto.HealthCheckResponse{Status: "Ok"}, nil
}
