package cached

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	cachev1 "github.com/jasonmccallister/nats-cache/gen"
	"github.com/jasonmccallister/nats-cache/gen/cachev1connect"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	cachev1connect.UnimplementedCacheServiceHandler
}

// Delete implements cachev1connect.CacheServiceHandler.
func (*server) Delete(context.Context, *connect.BidiStream[cachev1.DeleteRequest, cachev1.DeleteResponse]) error {
	return status.Error(codes.Unimplemented, "not implemented yet")
}

// Get implements cachev1connect.CacheServiceHandler.
func (*server) Get(ctx context.Context, stream *connect.BidiStream[cachev1.GetRequest, cachev1.GetResponse]) error {
	for {
		req, err := stream.Receive()
		if err != nil {
			return err
		}

		fmt.Println("Received:", req.GetKey())
	}

}

// Set implements cachev1connect.CacheServiceHandler.
func (*server) Set(context.Context, *connect.BidiStream[cachev1.SetRequest, cachev1.SetResponse]) error {
	return status.Error(codes.Unimplemented, "not implemented yet")
}

// NewServer returns a new server for the cache service.
func NewServer() cachev1connect.CacheServiceHandler {
	return &server{}
}
