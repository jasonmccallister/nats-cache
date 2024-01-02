package cached

import (
	"context"
	"fmt"
	"log/slog"
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
	Logger     *slog.Logger

	cachev1connect.UnimplementedCacheServiceHandler
}

// Delete implements cachev1connect.CacheServiceHandler.
func (s *server) Delete(ctx context.Context, stream *connect.BidiStream[cachev1.DeleteRequest, cachev1.DeleteResponse]) error {
	t, err := s.Authorizer.Authorize(stream.RequestHeader().Get("Authorization"))
	if err != nil {
		s.Logger.ErrorContext(ctx, "failed to authorize request", "error", err.Error())
		return connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("failed to authorize request: %w", err))
	}

	for {
		req, err := stream.Receive()
		if err != nil {
			s.Logger.ErrorContext(ctx, "failed to receive request", "error", err.Error())
			return err
		}

		internalKey, _, err := keygen.FromToken(*t, req.GetDatabase(), req.GetKey())
		if err != nil {
			s.Logger.ErrorContext(ctx, "failed to create internal key", "error", err.Error())
			return connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create key: %w", err))
		}

		// maybe consider removing the db from the delete request and rely on a generic key?
		if err := s.Store.Delete(ctx, internalKey); err != nil {
			s.Logger.ErrorContext(ctx, "failed to delete key", "error", err.Error())
			return err
		}

		if err := stream.Send(&cachev1.DeleteResponse{
			Deleted: true,
		}); err != nil {
			s.Logger.ErrorContext(ctx, "failed to send response", "error", err.Error())
			return err
		}
	}
}

// Exists checks if any of the provided keys exists. It will return the keys (by name) that do exist.
func (s *server) Exists(ctx context.Context, req *connect.Request[cachev1.ExistsRequest]) (*connect.Response[cachev1.ExistsResponse], error) {
	t, err := s.Authorizer.Authorize(req.Header().Get("Authorization"))
	if err != nil {
		s.Logger.ErrorContext(ctx, "failed to authorize request", "error", err.Error())
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("failed to authorize request: %w", err))
	}

	var found []string
	for _, k := range req.Msg.GetKeys() {
		internalKey, _, err := keygen.FromToken(*t, req.Msg.GetDatabase(), k)
		if err != nil {
			s.Logger.ErrorContext(ctx, "failed to create internal key", "error", err.Error())
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create key: %w", err))
		}

		val, _, err := s.Store.Get(ctx, internalKey)
		if err != nil {
			s.Logger.ErrorContext(ctx, "failed to get key", "error", err.Error())
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get key: %w", err))
		}

		if val != nil {
			found = append(found, k)
		}
	}

	s.Logger.InfoContext(ctx, "found keys", "keys", found)

	return connect.NewResponse(&cachev1.ExistsResponse{
		Keys:  found,
		Count: uint32(len(found)),
	}), nil
}

// Get implements cachev1connect.CacheServiceHandler.
func (s *server) Get(ctx context.Context, stream *connect.BidiStream[cachev1.GetRequest, cachev1.GetResponse]) error {
	t, err := s.Authorizer.Authorize(stream.RequestHeader().Get("Authorization"))
	if err != nil {
		s.Logger.ErrorContext(ctx, "failed to authorize request", "error", err.Error())
		return connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("failed to authorize request: %w", err))
	}

	for {
		req, err := stream.Receive()
		if err != nil {
			s.Logger.ErrorContext(ctx, "failed to receive request", "error", err.Error())
			return err
		}

		internalKey, key, err := keygen.FromToken(*t, req.GetDatabase(), req.GetKey())
		if err != nil {
			s.Logger.ErrorContext(ctx, "failed to create internal key", "error", err.Error())
			return connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create key: %w", err))
		}

		value, ttl, err := s.Store.Get(ctx, internalKey)
		if err != nil {
			s.Logger.ErrorContext(ctx, "failed to get key", "error", err.Error())
			err := stream.Send(&cachev1.GetResponse{
				Key:   key,
				Value: nil,
				Ttl:   ttl,
			})
			if err != nil {
				return err
			}

			continue
		}

		if err := stream.Send(&cachev1.GetResponse{
			Key:   key,
			Value: value,
			Ttl:   ttl,
		}); err != nil {
			s.Logger.ErrorContext(ctx, "failed to send response", "error", err.Error())
			return err
		}
	}
}

// Set implements cachev1connect.CacheServiceHandler.
func (s *server) Set(ctx context.Context, stream *connect.BidiStream[cachev1.SetRequest, cachev1.SetResponse]) error {
	t, err := s.Authorizer.Authorize(stream.RequestHeader().Get("Authorization"))
	if err != nil {
		s.Logger.ErrorContext(ctx, "failed to authorize request", "error", err.Error())
		return connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("failed to authorize request: %w", err))
	}

	for {
		req, err := stream.Receive()
		if err != nil {
			s.Logger.ErrorContext(ctx, "failed to receive request", "error", err.Error())
			return err
		}

		internalKey, key, err := keygen.FromToken(*t, req.GetDatabase(), req.GetKey())
		if err != nil {
			s.Logger.ErrorContext(ctx, "failed to create internal key", "error", err.Error())
			return connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create key: %w", err))
		}

		// did the user provide a value for the ttl?
		var ttl int64
		if req.GetTtl() > 0 {
			ttl = time.Now().Add(time.Duration(req.GetTtl()) * time.Second).Unix()
		}

		if err := s.Store.Set(ctx, internalKey, req.GetValue(), ttl); err != nil {
			s.Logger.ErrorContext(ctx, "failed to set key", "error", err.Error())
			return err
		}

		if err := stream.Send(&cachev1.SetResponse{
			Key:   key,
			Value: req.GetValue(),
			Ttl:   ttl,
		}); err != nil {
			s.Logger.ErrorContext(ctx, "failed to send response", "error", err.Error())
			return err
		}
	}
}

// Purge implements cachev1connect.CacheServiceHandler.
func (s *server) Purge(ctx context.Context, stream *connect.BidiStream[cachev1.PurgeRequest, cachev1.PurgeResponse]) error {
	t, err := s.Authorizer.Authorize(stream.RequestHeader().Get("Authorization"))
	if err != nil {
		s.Logger.ErrorContext(ctx, "failed to authorize request", "error", err.Error())
		return connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("failed to authorize request: %w", err))
	}

	for {
		req, err := stream.Receive()
		if err != nil {
			s.Logger.ErrorContext(ctx, "failed to receive request", "error", err.Error())
			return err
		}

		internalKey, _, err := keygen.FromToken(*t, req.GetDatabase(), req.GetPrefix())
		if err != nil {
			s.Logger.ErrorContext(ctx, "failed to create internal key", "error", err.Error())
			return connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create key: %w", err))
		}

		if err := s.Store.Purge(ctx, internalKey); err != nil {
			s.Logger.ErrorContext(ctx, "failed to purge keys", "error", err.Error())
			return err
		}

		if err := stream.Send(&cachev1.PurgeResponse{
			Purged: true,
		}); err != nil {
			s.Logger.ErrorContext(ctx, "failed to send response", "error", err.Error())
			return err
		}
	}
}

// NewServer returns a new server for the cache service.
func NewServer(l *slog.Logger, a auth.Authorizer, s storage.Store) cachev1connect.CacheServiceHandler {
	return &server{
		Logger:     l,
		Authorizer: a,
		Store:      s,
	}
}
