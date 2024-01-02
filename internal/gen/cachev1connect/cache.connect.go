// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: cache.proto

package cachev1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	gen "github.com/jasonmccallister/nats-cache/internal/gen"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

const (
	// CacheServiceName is the fully-qualified name of the CacheService service.
	CacheServiceName = "cache.v1.CacheService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// CacheServiceExistsProcedure is the fully-qualified name of the CacheService's Exists RPC.
	CacheServiceExistsProcedure = "/cache.v1.CacheService/Exists"
	// CacheServiceGetProcedure is the fully-qualified name of the CacheService's Get RPC.
	CacheServiceGetProcedure = "/cache.v1.CacheService/Get"
	// CacheServiceSetProcedure is the fully-qualified name of the CacheService's Set RPC.
	CacheServiceSetProcedure = "/cache.v1.CacheService/Set"
	// CacheServiceDeleteProcedure is the fully-qualified name of the CacheService's Delete RPC.
	CacheServiceDeleteProcedure = "/cache.v1.CacheService/Delete"
	// CacheServicePurgeProcedure is the fully-qualified name of the CacheService's Purge RPC.
	CacheServicePurgeProcedure = "/cache.v1.CacheService/Purge"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	cacheServiceServiceDescriptor      = gen.File_cache_proto.Services().ByName("CacheService")
	cacheServiceExistsMethodDescriptor = cacheServiceServiceDescriptor.Methods().ByName("Exists")
	cacheServiceGetMethodDescriptor    = cacheServiceServiceDescriptor.Methods().ByName("Get")
	cacheServiceSetMethodDescriptor    = cacheServiceServiceDescriptor.Methods().ByName("Set")
	cacheServiceDeleteMethodDescriptor = cacheServiceServiceDescriptor.Methods().ByName("Delete")
	cacheServicePurgeMethodDescriptor  = cacheServiceServiceDescriptor.Methods().ByName("Purge")
)

// CacheServiceClient is a client for the cache.v1.CacheService service.
type CacheServiceClient interface {
	Exists(context.Context, *connect.Request[gen.ExistsRequest]) (*connect.Response[gen.ExistsResponse], error)
	// Get is responsible for retrieving a value from the cache. If the value is not found, the value will return as nil.
	Get(context.Context) *connect.BidiStreamForClient[gen.GetRequest, gen.GetResponse]
	// Set is responsible for setting a value in the cache.
	Set(context.Context) *connect.BidiStreamForClient[gen.SetRequest, gen.SetResponse]
	// Delete is responsible for deleting a value from the cache.
	Delete(context.Context) *connect.BidiStreamForClient[gen.DeleteRequest, gen.DeleteResponse]
	// Purge is responsible for purging all values from the cache.
	Purge(context.Context) *connect.BidiStreamForClient[gen.PurgeRequest, gen.PurgeResponse]
}

// NewCacheServiceClient constructs a client for the cache.v1.CacheService service. By default, it
// uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewCacheServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) CacheServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &cacheServiceClient{
		exists: connect.NewClient[gen.ExistsRequest, gen.ExistsResponse](
			httpClient,
			baseURL+CacheServiceExistsProcedure,
			connect.WithSchema(cacheServiceExistsMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		get: connect.NewClient[gen.GetRequest, gen.GetResponse](
			httpClient,
			baseURL+CacheServiceGetProcedure,
			connect.WithSchema(cacheServiceGetMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		set: connect.NewClient[gen.SetRequest, gen.SetResponse](
			httpClient,
			baseURL+CacheServiceSetProcedure,
			connect.WithSchema(cacheServiceSetMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		delete: connect.NewClient[gen.DeleteRequest, gen.DeleteResponse](
			httpClient,
			baseURL+CacheServiceDeleteProcedure,
			connect.WithSchema(cacheServiceDeleteMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		purge: connect.NewClient[gen.PurgeRequest, gen.PurgeResponse](
			httpClient,
			baseURL+CacheServicePurgeProcedure,
			connect.WithSchema(cacheServicePurgeMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// cacheServiceClient implements CacheServiceClient.
type cacheServiceClient struct {
	exists *connect.Client[gen.ExistsRequest, gen.ExistsResponse]
	get    *connect.Client[gen.GetRequest, gen.GetResponse]
	set    *connect.Client[gen.SetRequest, gen.SetResponse]
	delete *connect.Client[gen.DeleteRequest, gen.DeleteResponse]
	purge  *connect.Client[gen.PurgeRequest, gen.PurgeResponse]
}

// Exists calls cache.v1.CacheService.Exists.
func (c *cacheServiceClient) Exists(ctx context.Context, req *connect.Request[gen.ExistsRequest]) (*connect.Response[gen.ExistsResponse], error) {
	return c.exists.CallUnary(ctx, req)
}

// Get calls cache.v1.CacheService.Get.
func (c *cacheServiceClient) Get(ctx context.Context) *connect.BidiStreamForClient[gen.GetRequest, gen.GetResponse] {
	return c.get.CallBidiStream(ctx)
}

// Set calls cache.v1.CacheService.Set.
func (c *cacheServiceClient) Set(ctx context.Context) *connect.BidiStreamForClient[gen.SetRequest, gen.SetResponse] {
	return c.set.CallBidiStream(ctx)
}

// Delete calls cache.v1.CacheService.Delete.
func (c *cacheServiceClient) Delete(ctx context.Context) *connect.BidiStreamForClient[gen.DeleteRequest, gen.DeleteResponse] {
	return c.delete.CallBidiStream(ctx)
}

// Purge calls cache.v1.CacheService.Purge.
func (c *cacheServiceClient) Purge(ctx context.Context) *connect.BidiStreamForClient[gen.PurgeRequest, gen.PurgeResponse] {
	return c.purge.CallBidiStream(ctx)
}

// CacheServiceHandler is an implementation of the cache.v1.CacheService service.
type CacheServiceHandler interface {
	Exists(context.Context, *connect.Request[gen.ExistsRequest]) (*connect.Response[gen.ExistsResponse], error)
	// Get is responsible for retrieving a value from the cache. If the value is not found, the value will return as nil.
	Get(context.Context, *connect.BidiStream[gen.GetRequest, gen.GetResponse]) error
	// Set is responsible for setting a value in the cache.
	Set(context.Context, *connect.BidiStream[gen.SetRequest, gen.SetResponse]) error
	// Delete is responsible for deleting a value from the cache.
	Delete(context.Context, *connect.BidiStream[gen.DeleteRequest, gen.DeleteResponse]) error
	// Purge is responsible for purging all values from the cache.
	Purge(context.Context, *connect.BidiStream[gen.PurgeRequest, gen.PurgeResponse]) error
}

// NewCacheServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewCacheServiceHandler(svc CacheServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	cacheServiceExistsHandler := connect.NewUnaryHandler(
		CacheServiceExistsProcedure,
		svc.Exists,
		connect.WithSchema(cacheServiceExistsMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	cacheServiceGetHandler := connect.NewBidiStreamHandler(
		CacheServiceGetProcedure,
		svc.Get,
		connect.WithSchema(cacheServiceGetMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	cacheServiceSetHandler := connect.NewBidiStreamHandler(
		CacheServiceSetProcedure,
		svc.Set,
		connect.WithSchema(cacheServiceSetMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	cacheServiceDeleteHandler := connect.NewBidiStreamHandler(
		CacheServiceDeleteProcedure,
		svc.Delete,
		connect.WithSchema(cacheServiceDeleteMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	cacheServicePurgeHandler := connect.NewBidiStreamHandler(
		CacheServicePurgeProcedure,
		svc.Purge,
		connect.WithSchema(cacheServicePurgeMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/cache.v1.CacheService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case CacheServiceExistsProcedure:
			cacheServiceExistsHandler.ServeHTTP(w, r)
		case CacheServiceGetProcedure:
			cacheServiceGetHandler.ServeHTTP(w, r)
		case CacheServiceSetProcedure:
			cacheServiceSetHandler.ServeHTTP(w, r)
		case CacheServiceDeleteProcedure:
			cacheServiceDeleteHandler.ServeHTTP(w, r)
		case CacheServicePurgeProcedure:
			cacheServicePurgeHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedCacheServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedCacheServiceHandler struct{}

func (UnimplementedCacheServiceHandler) Exists(context.Context, *connect.Request[gen.ExistsRequest]) (*connect.Response[gen.ExistsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cache.v1.CacheService.Exists is not implemented"))
}

func (UnimplementedCacheServiceHandler) Get(context.Context, *connect.BidiStream[gen.GetRequest, gen.GetResponse]) error {
	return connect.NewError(connect.CodeUnimplemented, errors.New("cache.v1.CacheService.Get is not implemented"))
}

func (UnimplementedCacheServiceHandler) Set(context.Context, *connect.BidiStream[gen.SetRequest, gen.SetResponse]) error {
	return connect.NewError(connect.CodeUnimplemented, errors.New("cache.v1.CacheService.Set is not implemented"))
}

func (UnimplementedCacheServiceHandler) Delete(context.Context, *connect.BidiStream[gen.DeleteRequest, gen.DeleteResponse]) error {
	return connect.NewError(connect.CodeUnimplemented, errors.New("cache.v1.CacheService.Delete is not implemented"))
}

func (UnimplementedCacheServiceHandler) Purge(context.Context, *connect.BidiStream[gen.PurgeRequest, gen.PurgeResponse]) error {
	return connect.NewError(connect.CodeUnimplemented, errors.New("cache.v1.CacheService.Purge is not implemented"))
}