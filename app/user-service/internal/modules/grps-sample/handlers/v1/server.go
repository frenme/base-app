package v1

import (
	"context"
	echopb "shared/pkg/grps/proto/echo"
	appsvc "user/internal/modules/grps-sample"
)

type Server struct {
	echopb.UnimplementedEchoServiceServer
	service *appsvc.Service
}

func NewServer(service *appsvc.Service) *Server {
	return &Server{service: service}
}

func (s *Server) Echo(ctx context.Context, in *echopb.EchoRequest) (*echopb.EchoResponse, error) {
	msg, err := s.service.Echo(ctx, in.GetMessage())
	if err != nil {
		return nil, err
	}
	return &echopb.EchoResponse{Message: msg}, nil
}
