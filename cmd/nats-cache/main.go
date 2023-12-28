package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	"connectrpc.com/otelconnect"
	"github.com/jasonmccallister/nats-cache/gen/cachev1connect"
	"github.com/jasonmccallister/nats-cache/internal/auth"
	"github.com/jasonmccallister/nats-cache/internal/cached"
	"github.com/jasonmccallister/nats-cache/internal/storage"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	addr := flag.Uint("addr", 50051, "address of the server")
	publicKey := flag.String("public-key", "", "The public key to use for authorizing requests")
	flag.Parse()

	ctx := context.Background()

	//logger := log.New(os.Stdout, "nats-cache: ", log.LstdFlags)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// if the public key is not set, exit
	if *publicKey == "" {
		logger.Error("public key is required")
		os.Exit(1)
	}

	if err := run(ctx, logger, *addr, *publicKey); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func run(ctx context.Context, logger *slog.Logger, addr uint, publicKey string) error {
	mux := http.NewServeMux()

	reflector := grpcreflect.NewStaticReflector(
		cachev1connect.CacheServiceName,
	)
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	store := storage.NewInMemory()

	var opts []connect.HandlerOption

	// Add tracing middleware.
	opts = append(opts, connect.WithInterceptors(
		otelconnect.NewInterceptor(),
	))

	// Add authorization middleware.
	authorizer := auth.NewPasetoPublicKey(publicKey)

	logger.Info("authorizing requests with public key")

	// register the cache service
	mux.Handle(cachev1connect.NewCacheServiceHandler(
		cached.NewServer(authorizer, store),
		opts...,
	))

	// register the health service
	checker := grpchealth.NewStaticChecker(cachev1connect.CacheServiceName)
	mux.Handle(grpchealth.NewHandler(checker))

	logger.Info("starting server", "port", addr)

	// Use h2c so we can serve HTTP/2 without TLS.
	return http.ListenAndServe(fmt.Sprintf(":%d", addr), h2c.NewHandler(mux, &http2.Server{}))
}
