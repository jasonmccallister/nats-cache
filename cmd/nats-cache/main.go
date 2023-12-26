package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jasonmccallister/nats-cache/gen/cachev1connect"
	"github.com/jasonmccallister/nats-cache/internal/cached"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
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
	mux := http.NewServeMux()

	mux.Handle(cachev1connect.NewCacheServiceHandler(cached.NewServer()))

	fmt.Println("... Listening on", addr)

	return http.ListenAndServe(
		fmt.Sprintf(":%d", addr),
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
