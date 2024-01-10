package main

import (
	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	"connectrpc.com/otelconnect"
	"context"
	"fmt"
	"github.com/jasonmccallister/nats-cache/internal/auth"
	"github.com/jasonmccallister/nats-cache/internal/cached"
	"github.com/jasonmccallister/nats-cache/internal/embeddednats"
	"github.com/jasonmccallister/nats-cache/internal/gen/cache/v1/cachev1connect"
	"github.com/jasonmccallister/nats-cache/internal/localbucket"
	"github.com/jasonmccallister/nats-cache/internal/storage"
	"github.com/jasonmccallister/nats-cache/logs"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	ctx := context.Background()

	logger := logs.NewFromEnvironment()

	// create the nats server and start it
	ns, creds, err := embeddednats.NewServerFromEnvironment()
	if err != nil {
		logger.ErrorContext(ctx, fmt.Errorf("failed to create nats server: %w", err).Error())
		os.Exit(1)
	}
	// remove the creds file when we exit if it exists
	if creds != "" {
		logger.InfoContext(ctx, "store creds file", "path", creds)
		defer os.Remove(creds)
	}

	logger.DebugContext(ctx, "starting nats server", "url", ns.ClientURL())
	logger.DebugContext(ctx, "nats server leaf nodes", "leaf-nodes", ns.NumLeafNodes())

	go ns.Start()

	if !ns.ReadyForConnections(10 * time.Second) {
		logger.ErrorContext(ctx, "nats server failed to start")
		os.Exit(1)
	}

	nc, err := nats.Connect(ns.ClientURL(), nats.Name(os.Getenv("HOSTNAME")))
	if err != nil {
		logger.ErrorContext(ctx, fmt.Errorf("failed to connect to nats: %w", err).Error())
		os.Exit(1)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		logger.ErrorContext(ctx, fmt.Errorf("failed to create jetstream: %w", err).Error())
		os.Exit(1)
	}

	kv, err := localbucket.CreateFromEnv(ctx, js)
	if err != nil {
		logger.ErrorContext(ctx, fmt.Errorf("failed to create local bucket: %w", err).Error())
		os.Exit(1)
	}

	authorizer, err := auth.NewFromEnvironment()
	if err != nil {
		logger.ErrorContext(ctx, fmt.Errorf("failed to create authorizer: %w", err).Error())
		os.Exit(1)
	}

	if err := run(ctx, logger, authorizer, kv); err != nil {
		logger.ErrorContext(ctx, fmt.Errorf("failed to run server: %w", err).Error())
		os.Exit(2)
	}
}

func run(ctx context.Context, logger *slog.Logger, authorizer auth.Authorizer, kv jetstream.KeyValue) error {
	store := storage.NewNATSKeyValue(kv, logger)

	var opts []connect.HandlerOption
	opts = append(opts, connect.WithInterceptors(otelconnect.NewInterceptor()))
	server := cached.NewServer(logger, authorizer, store)

	// create the services
	cachePath, cacheHandler := cachev1connect.NewCacheServiceHandler(server, opts...)
	healthCheck, healthCheckHandler := grpchealth.NewHandler(grpchealth.NewStaticChecker(cachev1connect.CacheServiceName))
	reflectPath, reflectHandler := grpcreflect.NewHandlerV1Alpha(
		grpcreflect.NewStaticReflector(cachev1connect.CacheServiceName, grpchealth.HealthV1ServiceName),
	)

	mux := http.NewServeMux()

	// register the reflection service
	logger.InfoContext(ctx, "registering reflection service", "service", reflectPath, "path", reflectPath)
	mux.Handle(reflectPath, reflectHandler)

	// register the cache service
	logger.InfoContext(ctx, "registering cache service", "service", cachev1connect.CacheServiceName, "path", cachePath)
	mux.Handle(cachePath, cacheHandler)

	// register the health service
	logger.InfoContext(ctx, "registering health check", "service", cachev1connect.CacheServiceName, "path", healthCheck)
	mux.Handle(healthCheck, healthCheckHandler)

	port := 50051
	if v, ok := os.LookupEnv("APP_PORT"); ok {
		port, _ = strconv.Atoi(v)
	}

	logger.InfoContext(ctx, "starting server", "port", port)

	// Use h2c so we can serve HTTP/2 without TLS.
	return http.ListenAndServe(fmt.Sprintf(":%d", port), h2c.NewHandler(mux, &http2.Server{}))
}
