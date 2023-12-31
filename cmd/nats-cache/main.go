package main

import (
	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	"connectrpc.com/otelconnect"
	"context"
	"flag"
	"fmt"
	"github.com/jasonmccallister/nats-cache/internal/auth"
	"github.com/jasonmccallister/nats-cache/internal/cached"
	"github.com/jasonmccallister/nats-cache/internal/gen/cache/v1/cachev1connect"
	"github.com/jasonmccallister/nats-cache/internal/storage"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	addr := flag.Uint("addr", 50051, "address of the server")
	publicKey := flag.String("public-key", os.Getenv("NC_PUBLIC_KEY"), "The public key to use for authorizing requests")
	logLevel := flag.String("log-level", "info", "The log level to use")
	logFormat := flag.String("log-format", "text", "The log format to use")
	logOutput := flag.String("log-output", "stdout", "The log output to use")
	natsNKEY := flag.String("nats-nkey", os.Getenv("NATS_NKEY"), "The NATS nkey to use for NGS")
	natsJWT := flag.String("nats-jwt", os.Getenv("NATS_JWT"), "The NATS jwt to use for NGS")
	flag.Parse()

	ctx := context.Background()

	var logLevelOption slog.Level
	switch *logLevel {
	case "debug":
		logLevelOption = slog.LevelDebug
	case "warn":
		logLevelOption = slog.LevelWarn
	case "error":
		logLevelOption = slog.LevelError
	default:
		logLevelOption = slog.LevelInfo
	}

	var output *os.File
	switch *logOutput {
	case "stderr":
		output = os.Stderr
	default:
		output = os.Stdout
	}

	var logHandler slog.Handler
	switch *logFormat {
	case "json":
		logHandler = slog.NewJSONHandler(output, &slog.HandlerOptions{
			Level: logLevelOption,
		})
	default:
		logHandler = slog.NewTextHandler(output, &slog.HandlerOptions{
			Level: logLevelOption,
		})
	}

	logger := slog.New(logHandler)

	// if the public key is not set, exit
	if *publicKey == "" {
		logger.ErrorContext(ctx, "public key is required")
		os.Exit(1)
	}

	// if the nats nkey is not set, exit
	if *natsNKEY == "" {
		logger.ErrorContext(ctx, "nats-nkey is required")
		os.Exit(1)
	}

	// if the nats jwt is not set, exit
	if *natsJWT == "" {
		logger.ErrorContext(ctx, "nats-jwt is required")
		os.Exit(1)
	}

	if err := run(ctx, logger, *addr, *publicKey, *natsJWT, *natsNKEY); err != nil {
		logger.ErrorContext(ctx, fmt.Errorf("failed to run server: %w", err).Error())
		os.Exit(2)
	}
}

func run(ctx context.Context, logger *slog.Logger, addr uint, publicKey, jwt, nkey string) error {
	nc, err := nats.Connect("tls://connect.ngs.global", nats.UserJWTAndSeed(jwt, nkey), nats.Name(os.Getenv("HOSTNAME")))
	if err != nil {
		return fmt.Errorf("failed to connect to nats server: %w", err)
	}
	defer nc.Close()

	js, err := jetstream.New(nc)
	if err != nil {
		return fmt.Errorf("failed to create jetstream kv: %w", err)
	}

	// TODO(jasonmccallister) make this configurable
	kv, err := js.KeyValue(ctx, "cache")
	if err != nil {
		return fmt.Errorf("failed to create jetstream kv: %w", err)
	}

	store := storage.NewNATSKeyValue(kv, logger)
	authorizer := auth.NewPasetoPublicKey(publicKey)

	var opts []connect.HandlerOption
	opts = append(opts, connect.WithInterceptors(otelconnect.NewInterceptor()))

	// create all of the services
	cachePath, cacheHandler := cachev1connect.NewCacheServiceHandler(cached.NewServer(authorizer, store), opts...)
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

	logger.InfoContext(ctx, "starting server", "port", addr)

	// Use h2c so we can serve HTTP/2 without TLS.
	return http.ListenAndServe(fmt.Sprintf(":%d", addr), h2c.NewHandler(mux, &http2.Server{}))
}
