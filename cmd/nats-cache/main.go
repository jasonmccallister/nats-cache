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
	"github.com/nats-io/nats.go/jetstream"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	addr := flag.Uint("addr", 50051, "address of the server")
	publicKey := flag.String("public-key", "", "The public key to use for authorizing requests")
	logLevel := flag.String("log-level", "info", "The log level to use")
	logFormat := flag.String("log-format", "text", "The log format to use")
	flag.Parse()

	ctx := context.Background()

	var logLevelOption slog.Level
	switch *logLevel {
	case "debug":
		logLevelOption = slog.LevelDebug
	case "info":
		logLevelOption = slog.LevelInfo
	case "warn":
		logLevelOption = slog.LevelWarn
	case "error":
		logLevelOption = slog.LevelError
	default:
		logLevelOption = slog.LevelInfo
	}

	var logHandler slog.Handler
	switch *logFormat {
	case "json":
		logHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevelOption,
		})
	default:
		logHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevelOption,
		})
	}

	logger := slog.New(logHandler)

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
	ns, err := server.NewServer(&server.Options{
		JetStream: true,
		HTTPPort:  8222, // disable http for production
	})
	if err != nil {
		return fmt.Errorf("failed to create nats server: %w", err)
	}
	defer ns.Shutdown()

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

	js, err := jetstream.New(nc)
	if err != nil {
		return fmt.Errorf("failed to create jetstream kv: %w", err)
	}

	kv, err := js.CreateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket:   "cache",
		MaxBytes: 1024 * 1024,
	})
	if err != nil {
		return fmt.Errorf("failed to create jetstream kv: %w", err)
	}

	store := storage.NewNATSKeyValue(kv)
	authorizer := auth.NewPasetoPublicKey(publicKey)

	var opts []connect.HandlerOption
	opts = append(opts, connect.WithInterceptors(otelconnect.NewInterceptor()))

	mux := http.NewServeMux()

	// register the reflection service
	reflectPath, reflectHandler := grpcreflect.NewHandlerV1Alpha(grpcreflect.NewStaticReflector(cachev1connect.CacheServiceName))
	logger.InfoContext(ctx, "registering reflection service", "service", cachev1connect.CacheServiceName, "path", reflectPath)
	mux.Handle(reflectPath, reflectHandler)

	// register the cache service
	cachePath, cacheHandler := cachev1connect.NewCacheServiceHandler(cached.NewServer(authorizer, store), opts...)
	logger.InfoContext(ctx, "registering cache service", "service", cachev1connect.CacheServiceName, "path", cachePath)
	mux.Handle(cachePath, cacheHandler)

	// register the health service
	healthCheck, healthCheckHandler := grpchealth.NewHandler(grpchealth.NewStaticChecker(cachev1connect.CacheServiceName))
	logger.InfoContext(ctx, "registering health check", "service", cachev1connect.CacheServiceName, "path", healthCheck)
	mux.Handle(healthCheck, healthCheckHandler)

	logger.InfoContext(ctx, "starting server", "port", addr)

	// Use h2c so we can serve HTTP/2 without TLS.
	return http.ListenAndServe(fmt.Sprintf(":%d", addr), h2c.NewHandler(mux, &http2.Server{}))
}
