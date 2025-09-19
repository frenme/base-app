package echo

import (
	"context"
	echopb "shared/pkg/proto/echo"
)

type EchoServer struct {
	echopb.UnimplementedEchoServiceServer
}

func NewEchoServer() *EchoServer { return &EchoServer{} }

func (s *EchoServer) Echo(ctx context.Context, in *echopb.EchoRequest) (*echopb.EchoResponse, error) {
	return &echopb.EchoResponse{Message: in.GetMessage()}, nil
}
