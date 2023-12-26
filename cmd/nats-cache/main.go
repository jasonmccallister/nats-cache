package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	addr := flag.Uint("addr", 50051, "address of the server")
	flag.Parse()

	ctx := context.Background()

	logger := log.New(os.Stdout, "nats-cache: ", log.LstdFlags)

	if err := run(ctx, logger, uint32(*addr)); err != nil {
		logger.Fatal(err)
	}
}

func run(ctx context.Context, logger *log.Logger, addr uint32) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", addr))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	opts := &server.Options{
		HTTPPort:  8222,
		JetStream: true,
	}
	ns, err := server.NewServer(opts)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	go ns.Start()

	if !ns.ReadyForConnections(10 * time.Second) {
		return fmt.Errorf("failed to start server")
	}

	if _, err := nats.Connect(ns.ClientURL()); err != nil {
		return fmt.Errorf("failed to connect to local nats: %w", err)
	}

	go func() {
		ns.WaitForShutdown()
	}()

	grpcServer := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())

	reflection.Register(grpcServer)

	logger.Println("starting server on port", addr)

	return grpcServer.Serve(lis)
}
