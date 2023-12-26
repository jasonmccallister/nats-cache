package cached

import (
	"context"

	"connectrpc.com/connect"
	cachev1 "github.com/jasonmccallister/nats-cache/gen"
	"github.com/jasonmccallister/nats-cache/gen/cachev1connect"
	"github.com/jasonmccallister/nats-cache/internal/storage"
)

type server struct {
	cachev1connect.UnimplementedCacheServiceHandler
	Store storage.Store
}

// Delete implements cachev1connect.CacheServiceHandler.
func (s *server) Delete(ctx context.Context, stream *connect.BidiStream[cachev1.DeleteRequest, cachev1.DeleteResponse]) error {
	for {
		req, err := stream.Receive()
		if err != nil {
			return err
		}

		if err := s.Store.Delete(req.GetDatabase(), req.GetKey()); err != nil {
			return err
		}

		if err := stream.Send(&cachev1.DeleteResponse{
			Database: req.GetDatabase(),
			Key:      req.GetKey(),
		}); err != nil {
			return err
		}
	}
}

// Get implements cachev1connect.CacheServiceHandler.
func (s *server) Get(ctx context.Context, stream *connect.BidiStream[cachev1.GetRequest, cachev1.GetResponse]) error {
	for {
		req, err := stream.Receive()
		if err != nil {
			return err
		}

		value, err := s.Store.Get(req.GetDatabase(), req.GetKey())
		if err != nil {
			stream.Send(&cachev1.GetResponse{
				Database: req.GetDatabase(),
				Key:      req.GetKey(),
				Value:    nil,
			})

			continue
		}

		if err := stream.Send(&cachev1.GetResponse{
			Database: req.GetDatabase(),
			Key:      req.GetKey(),
			Value:    value,
		}); err != nil {
			return err
		}
	}
}

// Set implements cachev1connect.CacheServiceHandler.
func (s *server) Set(ctx context.Context, stream *connect.BidiStream[cachev1.SetRequest, cachev1.SetResponse]) error {
	for {
		req, err := stream.Receive()
		if err != nil {
			return err
		}

		if err := s.Store.Set(req.GetDatabase(), req.GetKey(), req.GetValue(), 0); err != nil {
			return err
		}

		if err := stream.Send(&cachev1.SetResponse{
			Database: req.GetDatabase(),
			Key:      req.GetKey(),
			Value:    req.GetValue(),
		}); err != nil {
			return err
		}
	}
}

// NewServer returns a new server for the cache service.
func NewServer(store storage.Store) cachev1connect.CacheServiceHandler {
	return &server{
		Store: store,
	}
}
