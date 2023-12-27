package cached

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	cachev1 "github.com/jasonmccallister/nats-cache/gen"
	"github.com/jasonmccallister/nats-cache/gen/cachev1connect"
	"github.com/jasonmccallister/nats-cache/internal/auth"
	"github.com/jasonmccallister/nats-cache/internal/keygen"
	"github.com/jasonmccallister/nats-cache/internal/storage"
)

type server struct {
	Authorizer auth.Authorizer
	Store      storage.Store

	cachev1connect.UnimplementedCacheServiceHandler
}

// Delete implements cachev1connect.CacheServiceHandler.
func (s *server) Delete(ctx context.Context, stream *connect.BidiStream[cachev1.DeleteRequest, cachev1.DeleteResponse]) error {
	t, err := s.Authorizer.Authorize(stream.RequestHeader().Get("Authorization"))
	if err != nil {
		return connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("failed to authorize request: %w", err))
	}

	for {
		req, err := stream.Receive()
		if err != nil {
			return err
		}

		internalKey, key, err := keygen.FromToken(*t, req.GetDatabase(), req.GetKey())
		if err != nil {
			return connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create key: %w", err))
		}

		// maybe consider removing the db from the delete request and rely on a generic key?
		if err := s.Store.Delete(internalKey); err != nil {
			return err
		}

		if err := stream.Send(&cachev1.DeleteResponse{
			Database: req.GetDatabase(),
			Key:      key,
		}); err != nil {
			return err
		}
	}
}

// Get implements cachev1connect.CacheServiceHandler.
func (s *server) Get(ctx context.Context, stream *connect.BidiStream[cachev1.GetRequest, cachev1.GetResponse]) error {
	t, err := s.Authorizer.Authorize(stream.RequestHeader().Get("Authorization"))
	if err != nil {
		return connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("failed to authorize request: %w", err))
	}

	for {
		req, err := stream.Receive()
		if err != nil {
			return err
		}

		internalKey, _, err := keygen.FromToken(*t, req.GetDatabase(), req.GetKey())
		if err != nil {
			return connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create key: %w", err))
		}

		value, err := s.Store.Get(internalKey)
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
	t, err := s.Authorizer.Authorize(stream.RequestHeader().Get("Authorization"))
	if err != nil {
		return connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("failed to authorize request: %w", err))
	}

	for {
		req, err := stream.Receive()
		if err != nil {
			return err
		}

		internalKey, _, err := keygen.FromToken(*t, req.GetDatabase(), req.GetKey())
		if err != nil {
			return connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create key: %w", err))
		}

		// TODO(jasonmccallister) did the user provide a value for the ttl?

		if err := s.Store.Set(internalKey, req.GetValue(), 0); err != nil {
			return err
		}

		if err := stream.Send(&cachev1.SetResponse{
			Database: req.GetDatabase(),
			Key:      req.GetKey(),
			Value:    req.GetValue(),
			Ttl:      0,
		}); err != nil {
			return err
		}
	}
}

// NewServer returns a new server for the cache service.
func NewServer(a auth.Authorizer, s storage.Store) cachev1connect.CacheServiceHandler {
	return &server{
		Authorizer: a,
		Store:      s,
	}
}
