package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	"connectrpc.com/otelconnect"
	"github.com/jasonmccallister/nats-cache/internal/auth"
	"github.com/jasonmccallister/nats-cache/internal/cached"
	"github.com/jasonmccallister/nats-cache/internal/gen/cachev1connect"
	"github.com/jasonmccallister/nats-cache/internal/storage"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	addr := flag.Uint("addr", 50051, "address of the server")
	publicKey := flag.String("public-key", "", "The public key to use for authorizing requests")
	flag.Parse()

	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// if the public key is not set, exit
	if *publicKey == "" {
		logger.ErrorContext(ctx, "public key is required")
		os.Exit(1)
	}

	if err := run(ctx, logger, *addr, *publicKey); err != nil {
		logger.ErrorContext(ctx, fmt.Errorf("failed to run server: %w", err).Error())
		os.Exit(2)
	}
}

func run(ctx context.Context, logger *slog.Logger, addr uint, publicKey string) error {
	mux := http.NewServeMux()
	reflector := grpcreflect.NewStaticReflector(cachev1connect.CacheServiceName)
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	// Add tracing middleware.
	var handlerOpts []connect.HandlerOption
	handlerOpts = append(handlerOpts, connect.WithInterceptors(
		otelconnect.NewInterceptor(),
	))

	// Add authorization middleware.
	authorizer := auth.NewPasetoPublicKey(publicKey)

	logger.InfoContext(ctx, "authorizing requests with public key")

	store := storage.NewInMemory()
	// register the cache service
	mux.Handle(cachev1connect.NewCacheServiceHandler(cached.NewServer(authorizer, store), handlerOpts...))

	opts := &server.Options{
		JetStream: true,
		HTTPPort:  8222,
	}

	ns, err := server.NewServer(opts)
	if err != nil {
		return fmt.Errorf("failed to create nats server: %w", err)
	}

	ns.Start()

	if !ns.ReadyForConnections(10 * time.Second) {
		logger.ErrorContext(ctx, "nats server did not start")
		os.Exit(1)
	}

	nc, err := nats.Connect(ns.ClientURL())
	if err != nil {
		return fmt.Errorf("failed to connect to nats server: %w", err)
	}
	defer nc.Close()

	// register the health service
	checker := grpchealth.NewStaticChecker(cachev1connect.CacheServiceName)
	mux.Handle(grpchealth.NewHandler(checker))

	logger.InfoContext(ctx, "starting server", "port", addr)

	// Use h2c so we can serve HTTP/2 without TLS.
	return http.ListenAndServe(fmt.Sprintf(":%d", addr), h2c.NewHandler(mux, &http2.Server{}))
}
