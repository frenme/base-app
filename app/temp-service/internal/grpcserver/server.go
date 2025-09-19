package grpcserver

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"shared/pkg/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func getPort() string {
	port := os.Getenv("TEMP_SERVICE_GRPC_PORT")
	if port == "" {
		port = "50051"
	}
	return port
}

// StartGRPCServer запускает gRPC-сервер неблокирующе (в горутине)
func StartGRPCServer(ctx context.Context, log *logger.Logger) {
	go func() {
		addr := ":" + getPort()
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			if log != nil {
				log.WithError(err).Error("grpc: failed to listen")
			} else {
				fmt.Println("grpc: failed to listen:", err)
			}
			return
		}

		srv := grpc.NewServer()

		// health-checks
		hs := health.NewServer()
		healthpb.RegisterHealthServer(srv, hs)

		// reflection для дебага
		reflection.Register(srv)

		if log != nil {
			log.WithField("addr", addr).Info("grpc: server starting")
		} else {
			fmt.Println("grpc: server starting on", addr)
		}

		// graceful shutdown
		go func() {
			<-ctx.Done()
			hs.Shutdown()
			stopped := make(chan struct{})
			go func() {
				srv.GracefulStop()
				close(stopped)
			}()
			select {
			case <-stopped:
			case <-time.After(5 * time.Second):
				srv.Stop()
			}
		}()

		if err := srv.Serve(lis); err != nil {
			if log != nil {
				log.WithError(err).Error("grpc: server stopped")
			} else {
				fmt.Println("grpc: server stopped:", err)
			}
		}
	}()
}


