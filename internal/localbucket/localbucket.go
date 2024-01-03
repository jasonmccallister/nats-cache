package localbucket

import (
	"context"
	"errors"
	"github.com/nats-io/nats.go/jetstream"
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
