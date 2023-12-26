package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"
	"github.com/jasonmccallister/nats-cache/gen/cachev1connect"
	"github.com/jasonmccallister/nats-cache/internal/cached"
	"github.com/jasonmccallister/nats-cache/internal/storage"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	addr := flag.Uint("addr", 50051, "address of the server")
	flag.Parse()

	ctx := context.Background()

	logger := log.New(os.Stdout, "nats-cache: ", log.LstdFlags)

	if err := run(ctx, logger, *addr); err != nil {
		logger.Fatal(err)
	}
}

func run(ctx context.Context, logger *log.Logger, addr uint) error {
	mux := http.NewServeMux()

	store := storage.NewInMemory()

	var opts []connect.HandlerOption

	// Add tracing middleware.
	opts = append(opts, connect.WithInterceptors(
		otelconnect.NewInterceptor(),
	))

	mux.Handle(cachev1connect.NewCacheServiceHandler(
		cached.NewServer(store),
		opts...,
	))

	// TODO(jasonmccallister) register reflection service on gRPC server using connectrpc.com/grpcreflect
	// TODO(jasonmccallister) register health service on gRPC server using connectrpc.com/grpchealth

	logger.Println("starting server on port", addr)

	// Use h2c so we can serve HTTP/2 without TLS.
	return http.ListenAndServe(fmt.Sprintf(":%d", addr), h2c.NewHandler(mux, &http2.Server{}))
}
