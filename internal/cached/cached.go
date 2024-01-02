package cached

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"
	"github.com/jasonmccallister/nats-cache/internal/auth"
	cachev1 "github.com/jasonmccallister/nats-cache/internal/gen/cache/v1"
	"github.com/jasonmccallister/nats-cache/internal/gen/cache/v1/cachev1connect"
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

		internalKey, _, err := keygen.FromToken(*t, req.GetDatabase(), req.GetKey())
		if err != nil {
			return connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create key: %w", err))
		}

		// maybe consider removing the db from the delete request and rely on a generic key?
		if err := s.Store.Delete(ctx, internalKey); err != nil {
			return err
		}

		if err := stream.Send(&cachev1.DeleteResponse{
			Deleted: true,
		}); err != nil {
			return err
		}
	}
}

// Exists checks if any of the provided keys exists. It will return the keys (by name) that do exist.
func (s *server) Exists(ctx context.Context, req *connect.Request[cachev1.ExistsRequest]) (*connect.Response[cachev1.ExistsResponse], error) {
	t, err := s.Authorizer.Authorize(req.Header().Get("Authorization"))
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("failed to authorize request: %w", err))
	}

	var found []string
	for _, k := range req.Msg.GetKeys() {
		internalKey, _, err := keygen.FromToken(*t, req.Msg.GetDatabase(), k)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create key: %w", err))
		}

		val, err := s.Store.Get(ctx, internalKey)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get key: %w", err))
		}

		if val != nil {
			found = append(found, k)
		}
	}

	return connect.NewResponse(&cachev1.ExistsResponse{
		Keys:  found,
		Count: uint32(len(found)),
	}), nil
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

		internalKey, key, err := keygen.FromToken(*t, req.GetDatabase(), req.GetKey())
		if err != nil {
			return connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create key: %w", err))
		}

		value, err := s.Store.Get(ctx, internalKey)
		if err != nil {
			err := stream.Send(&cachev1.GetResponse{
				Key:   key,
				Value: nil,
			})
			if err != nil {
				return err
			}

			continue
		}

		if err := stream.Send(&cachev1.GetResponse{
			Key:   key,
			Value: value,
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

		// did the user provide a value for the ttl?
		var ttl int64
		if req.GetTtl() > 0 {
			ttl = time.Now().Add(time.Duration(req.GetTtl()) * time.Second).Unix()
		}

		if err := s.Store.Set(ctx, internalKey, req.GetValue(), ttl); err != nil {
			return err
		}

		if err := stream.Send(&cachev1.SetResponse{
			Value: req.GetValue(),
			Ttl:   uint32(ttl),
		}); err != nil {
			return err
		}
	}
}

// Purge implements cachev1connect.CacheServiceHandler.
func (s *server) Purge(ctx context.Context, stream *connect.BidiStream[cachev1.PurgeRequest, cachev1.PurgeResponse]) error {
	t, err := s.Authorizer.Authorize(stream.RequestHeader().Get("Authorization"))
	if err != nil {
		return connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("failed to authorize request: %w", err))
	}

	for {
		req, err := stream.Receive()
		if err != nil {
			return err
		}

		internalKey, _, err := keygen.FromToken(*t, req.GetDatabase(), req.GetPrefix())
		if err != nil {
			return connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create key: %w", err))
		}

		if err := s.Store.Purge(ctx, internalKey); err != nil {
			return err
		}

		if err := stream.Send(&cachev1.PurgeResponse{
			Purged: true,
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
