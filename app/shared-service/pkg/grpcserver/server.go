package grpcserver

import (
	"context"
	"fmt"
	"net"
	"time"

	"shared/pkg/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func Start(ctx context.Context, log *logger.Logger, listenAddress string) {
	go func() {
		lis, err := net.Listen("tcp", listenAddress)
		if err != nil {
			if log != nil {
				log.WithError(err).Error("grpc: failed to listen")
			} else {
				fmt.Println("grpc: failed to listen:", err)
			}
			return
		}

		srv := grpc.NewServer()

		hs := health.NewServer()
		healthpb.RegisterHealthServer(srv, hs)

		reflection.Register(srv)

		if log != nil {
			log.WithField("addr", listenAddress).Info("grpc: server starting")
		} else {
			fmt.Println("grpc: server starting on", listenAddress)
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
