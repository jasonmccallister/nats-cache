package localbucket

import (
	"context"
	"errors"
	"github.com/nats-io/nats.go/jetstream"
)

// Create creates a new bucket if it does not exist or returns the existing bucket
func Create(ctx context.Context, js jetstream.JetStream, bucketName, streamSourceName string) (jetstream.KeyValue, error) {
	kv, err := js.KeyValue(ctx, bucketName)
	if err != nil {
		if errors.Is(err, jetstream.ErrBucketNotFound) {
			kv, err := js.CreateKeyValue(ctx, jetstream.KeyValueConfig{
				Bucket: bucketName,
				Mirror: &jetstream.StreamSource{
					Name: streamSourceName,
					External: &jetstream.ExternalStream{
						APIPrefix: "$JS.ngs.API",
					},
				},
				MaxBytes: 1024 * 1024 * 1024,
			})
			if err != nil {
				return nil, err
			}

			return kv, nil
		}
	}

	return kv, nil
}
