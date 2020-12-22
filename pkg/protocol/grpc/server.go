package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	// "getitqec.com/server/auth/pkg/handlers"
	pb "getitqec.com/server/file/pkg/api/v1"
	"getitqec.com/server/file/pkg/logger"
	"getitqec.com/server/file/pkg/protocol/grpc/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// RunServer ...
func RunServer(ctx context.Context, service pb.FileServiceServer, port string, certFilePath string, keyFilePath string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	// gRPC server statup options
	opts := []grpc.ServerOption{}
	// cert
	if certFilePath != "" && keyFilePath != "" {
		creds, err := credentials.NewServerTLSFromFile(certFilePath, keyFilePath)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		// opts = []grpc.ServerOption{grpc.Creds(creds)}
		opts = append(opts, grpc.Creds(creds))
	}

	// add middleware
	uLog, sLog := middleware.LogInterceptor(logger.Log)
	opts = middleware.AddInterceptor([]grpc.UnaryServerInterceptor{
		middleware.UnaryTagInterceptor(),
		uLog,
		// middleware.UnaryAuthInterceptor(),
	}, []grpc.StreamServerInterceptor{
		middleware.StreamTagInterceptor(),
		sLog,
		// middleware.StreamAuthInterceptor(),
	}, opts)
	// auth middleware
	fmt.Printf("opts len : %d\n", len(opts))
	fmt.Printf("opts[1] : %v\n", opts[1])
	// opts = append(opts, grpc.UnaryInterceptor(common.AuthUnaryInterceptor))

	// register handlers
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterFileServiceServer(grpcServer, service)

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			logger.Log.Warn("shutting down gRPC server...")

			grpcServer.GracefulStop()

			<-ctx.Done()
		}
	}()

	// start gRPC server
	logger.Log.Info("starting gRPC server...")
	return grpcServer.Serve(listen)
}
