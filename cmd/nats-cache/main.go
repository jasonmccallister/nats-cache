package main

import (
	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	"connectrpc.com/otelconnect"
	"context"
	"flag"
	"fmt"
	"github.com/jasonmccallister/nats-cache/getenv"
	"github.com/jasonmccallister/nats-cache/internal/auth"
	"github.com/jasonmccallister/nats-cache/internal/cached"
	"github.com/jasonmccallister/nats-cache/internal/embeddednats"
	"github.com/jasonmccallister/nats-cache/internal/gen/cache/v1/cachev1connect"
	"github.com/jasonmccallister/nats-cache/internal/localbucket"
	"github.com/jasonmccallister/nats-cache/internal/storage"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	port := flag.Int("port", getenv.Int("APP_PORT", 50051), "address of the server")
	publicKey := flag.String("public-key", getenv.String("AUTH_PUBLIC_KEY", ""), "The public key to use for authorizing requests")
	logLevel := flag.String("log-level", getenv.String("LOG_LEVEL", "error"), "The log level to use")
	logFormat := flag.String("log-format", getenv.String("LOG_FORMAT", "json"), "The log format to use")
	logOutput := flag.String("log-output", getenv.String("LOG_OUTPUT", "stdout"), "The log output to use")
	natsNKEY := flag.String("nats-nkey", getenv.String("NATS_NKEY", ""), "The NATS nkey to use for NGS")
	natsJWT := flag.String("nats-jwt", getenv.String("NATS_JWT", ""), "The NATS jwt to use for NGS")
	natsHttpPort := flag.Int("nats-http-port", getenv.Int("NATS_HTTP_PORT", 0), "The NATS http port to use for the embedded server")
	natsPort := flag.Int("nats-port", getenv.Int("NATS_PORT", 0), "The NATS port to use for the embedded server")
	kvBucket := flag.String("nats-kv-bucket-name", getenv.String("NATS_KV_BUCKET_NAME", "cache"), "The NATS KV bucket name to use for the embedded server")
	streamSourceName := flag.String("nats-stream-source-name", getenv.String("NATS_STREAM_SOURCE_NAME", "cache"), "The NATS stream source name to use for the embedded server (will be prefixed with KV_")
	natsLocalStorage := flag.String("nats-local-storage", getenv.String("NATS_LOCAL_STORAGE", ""), "The NATS local storage to use for the embedded server")
	natsBucketMaxBytes := flag.Int64("nats-bucket-max-bytes", getenv.Int64("NATS_BUCKET_MAX_BYTES", 1024*1024*1024), "The NATS bucket max bytes to use for the embedded server")

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

	// create the nats server and start it
	ns, creds, err := embeddednats.NewServer(*natsNKEY, *natsJWT, *natsPort, *natsHttpPort)
	if err != nil {
		logger.ErrorContext(ctx, fmt.Errorf("failed to create nats server: %w", err).Error())
		os.Exit(1)
	}
	// remove the creds file when we exit
	defer os.Remove(creds)

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

	storageType := jetstream.FileStorage
	if *natsLocalStorage == "memory" {
		storageType = jetstream.MemoryStorage
	}

	var bucketOpts []localbucket.OptionsFunc
	bucketOpts = append(bucketOpts, localbucket.WithBucketName(*kvBucket))
	bucketOpts = append(bucketOpts, localbucket.WithStreamSourceName(*streamSourceName))
	bucketOpts = append(bucketOpts, localbucket.WithStorage(storageType))
	bucketOpts = append(bucketOpts, localbucket.WithMaxBytes(*natsBucketMaxBytes))

	kv, err := localbucket.Create(ctx, js, bucketOpts...)
	if err != nil {
		logger.ErrorContext(ctx, fmt.Errorf("failed to create local bucket: %w", err).Error())
		os.Exit(1)
	}

	if err := run(ctx, logger, *port, *publicKey, kv); err != nil {
		logger.ErrorContext(ctx, fmt.Errorf("failed to run server: %w", err).Error())
		os.Exit(2)
	}
}

func run(ctx context.Context, logger *slog.Logger, addr int, publicKey string, kv jetstream.KeyValue) error {
	store := storage.NewNATSKeyValue(kv, logger)
	authorizer := auth.NewPasetoPublicKey(publicKey)

	var opts []connect.HandlerOption
	opts = append(opts, connect.WithInterceptors(otelconnect.NewInterceptor()))

	// create the services
	cachePath, cacheHandler := cachev1connect.NewCacheServiceHandler(cached.NewServer(logger, authorizer, store), opts...)
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
