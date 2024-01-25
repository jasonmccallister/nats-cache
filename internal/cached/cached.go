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

// NewServer returns a new server for the cache service.
func NewServer(l *slog.Logger, a auth.Authorizer, s storage.Store) cachev1connect.CacheServiceHandler {
	return &server{
		Logger:     l,
		Authorizer: a,
		Store:      s,
	}
}

func (s *server) Delete(ctx context.Context, req *connect.Request[cachev1.DeleteRequest]) (*connect.Response[cachev1.DeleteResponse], error) {
	t, err := s.Authorizer.Authorize(req.Header().Get("Authorization"))
	if err != nil {
		s.Logger.ErrorContext(ctx, "failed to authorize request", "error", err.Error())
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("failed to authorize request: %w", err))
	}

	start := time.Now()

	internalKey, _, err := keygen.FromToken(*t, req.Msg.GetDatabase(), req.Msg.GetKey())
	if err != nil {
		s.Logger.ErrorContext(ctx, "failed to create internal key", "error", err.Error())
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create key: %w", err))
	}

	// maybe consider removing the db from the delete request and rely on a generic key?
	if err := s.Store.Delete(ctx, internalKey); err != nil {
		s.Logger.ErrorContext(ctx, "failed to delete key", "error", err.Error())
		return nil, err
	}

	s.Logger.DebugContext(ctx, "delete", "key", internalKey, "duration", time.Since(start).String())

	return connect.NewResponse(&cachev1.DeleteResponse{
		Deleted: true,
	}), nil
}

// Exists checks if any of the provided keys exists. It will return the keys (by name) that do exist.
func (s *server) Exists(ctx context.Context, req *connect.Request[cachev1.ExistsRequest]) (*connect.Response[cachev1.ExistsResponse], error) {
	t, err := s.Authorizer.Authorize(req.Header().Get("Authorization"))
	if err != nil {
		s.Logger.ErrorContext(ctx, "failed to authorize request", "error", err.Error())
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("failed to authorize request: %w", err))
	}

	start := time.Now()

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

	s.Logger.DebugContext(ctx, "exists", "keys", found, "duration", time.Since(start).String())

	return connect.NewResponse(&cachev1.ExistsResponse{
		Keys:  found,
		Count: uint32(len(found)),
	}), nil
}

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

		start := time.Now()
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

		s.Logger.DebugContext(ctx, "get", "key", internalKey, "duration", time.Since(start).String())

		if err := stream.Send(&cachev1.GetResponse{
			Key:   key,
			Value: value,
			Ttl:   ttl,
		}); err != nil {
			s.Logger.ErrorContext(ctx, "failed to send response", "error", err.Error())
			return err
		}

		s.Logger.DebugContext(ctx, "get", "key", internalKey, "duration", time.Since(start).String())
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

		start := time.Now()

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
			return connect.NewError(connect.CodeInternal, fmt.Errorf("failed to set key: %w", err))
		}

		s.Logger.DebugContext(ctx, "set", "key", internalKey, "duration", time.Since(start).String())

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
func (s *server) Purge(ctx context.Context, req *connect.Request[cachev1.PurgeRequest]) (*connect.Response[cachev1.PurgeResponse], error) {
	t, err := s.Authorizer.Authorize(req.Header().Get("Authorization"))
	if err != nil {
		s.Logger.ErrorContext(ctx, "failed to authorize request", "error", err.Error())
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("failed to authorize request: %w", err))
	}

	internalKey, _, err := keygen.FromToken(*t, req.Msg.GetDatabase(), req.Msg.GetPrefix())
	if err != nil {
		s.Logger.ErrorContext(ctx, "failed to create internal key", "error", err.Error())
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create key: %w", err))
	}

	if err := s.Store.Purge(ctx, internalKey); err != nil {
		s.Logger.ErrorContext(ctx, "failed to purge keys", "error", err.Error())
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to purge keys: %w", err))
	}

	return connect.NewResponse(&cachev1.PurgeResponse{
		Purged: true,
	}), nil
}
