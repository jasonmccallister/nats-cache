package localbucket

import (
	"context"
	"errors"
	"github.com/jasonmccallister/nats-cache/getenv"
	"github.com/nats-io/nats.go/jetstream"
	"os"
)

// OptionsFunc is a function that sets options for the bucket
type OptionsFunc func(*Option)

// WithMaxBytes sets the max bytes for the bucket
func WithMaxBytes(maxBytes int64) OptionsFunc {
	return func(o *Option) {
		o.MaxBytes = maxBytes
	}
}

// WithBucketName sets the name for the bucket
func WithBucketName(bucketName string) OptionsFunc {
	return func(o *Option) {
		o.BucketName = bucketName
	}
}

// WithStreamSourceName sets the stream source name for the bucket
func WithStreamSourceName(streamSourceName string) OptionsFunc {
	return func(o *Option) {
		o.StreamSourceName = streamSourceName
	}
}

// WithStorage sets the storage type for the bucket
func WithStorage(storage jetstream.StorageType) OptionsFunc {
	return func(o *Option) {
		o.Storage = storage
	}
}

func defaultOptions() Option {
	return Option{
		BucketName:       "cache",
		MaxBytes:         1024 * 1024 * 1024,
		StreamSourceName: "cache",
		Storage:          jetstream.FileStorage,
	}
}

type Option struct {
	BucketName       string
	StreamSourceName string
	Storage          jetstream.StorageType
	MaxBytes         int64
}

// CreateFromEnv sets the default options for the bucket and checks the environment for overrides.
func CreateFromEnv(ctx context.Context, js jetstream.JetStream) (jetstream.KeyValue, error) {
	opts := defaultOptions()

	if v, ok := os.LookupEnv("NATS_BUCKET_NAME"); ok {
		opts.BucketName = v
	}

	if _, ok := os.LookupEnv("NATS_BUCKET_MAX_BYTES"); ok {
		opts.MaxBytes = getenv.Int64("NATS_BUCKET_MAX_BYTES", 1024*1024*1024)
	}

	if v, ok := os.LookupEnv("NATS_STREAM_SOURCE_NAME"); ok {
		opts.StreamSourceName = v
	}

	if v, ok := os.LookupEnv("NATS_LOCAL_STORAGE"); ok {
		switch v {
		case "memory":
			opts.Storage = jetstream.MemoryStorage
		default:
			opts.Storage = jetstream.FileStorage
		}
	}

	kv, err := Create(ctx, js)
	if err != nil {
		return nil, err
	}

	return kv, nil
}

// Create creates a new bucket if it does not exist or returns the existing bucket
func Create(ctx context.Context, js jetstream.JetStream, opts ...OptionsFunc) (jetstream.KeyValue, error) {
	o := defaultOptions()
	for _, fn := range opts {
		fn(&o)
	}

	kv, err := js.KeyValue(ctx, o.BucketName)
	if err != nil {
		if errors.Is(err, jetstream.ErrBucketNotFound) {
			kv, err := js.CreateKeyValue(ctx, jetstream.KeyValueConfig{
				Bucket: o.BucketName,
				Mirror: &jetstream.StreamSource{
					Name: o.StreamSourceName,
					External: &jetstream.ExternalStream{
						APIPrefix: "$JS.ngs.API",
					},
				},
				Storage:  o.Storage,
				MaxBytes: o.MaxBytes,
			})
			if err != nil {
				return nil, err
			}

			return kv, nil
		}
	}

	return kv, nil
}
