package storage

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"

	"github.com/nats-io/nats.go/jetstream"
)

type natsKeyValue struct {
	bucket jetstream.KeyValue
	logger *slog.Logger
}

func (n *natsKeyValue) Get(ctx context.Context, key string) ([]byte, int64, error) {
	v, err := n.bucket.Get(ctx, key)
	if err != nil {
		if errors.Is(err, jetstream.ErrKeyNotFound) {
			n.logger.InfoContext(ctx, "key not found", "key", key)

			return nil, 0, nil
		}

		n.logger.ErrorContext(ctx, "failed to get key", "key", key, "error", err.Error())

		return nil, 0, err
	}

	var i Item
	if err := json.Unmarshal(v.Value(), &i); err != nil {
		n.logger.ErrorContext(ctx, "failed to unmarshal item", "key", key, "error", err.Error())

		return nil, 0, err
	}

	// check if the item has expired
	if i.IsExpired() {
		n.logger.InfoContext(ctx, "key expired", "key", key, "ttl", i.TTL)

		defer n.purgeKey(ctx, key)

		return nil, 0, nil
	}

	n.logger.InfoContext(ctx, "got key", "key", key, "ttl", i.TTL)

	return i.Value, i.TTL, nil
}

func (n *natsKeyValue) purgeKey(ctx context.Context, keys ...string) error {
	for _, k := range keys {
		if err := n.bucket.Purge(ctx, k); err != nil {
			n.logger.ErrorContext(ctx, "failed to purge key", "key", k, "error", err.Error())
			continue
		}

		n.logger.InfoContext(ctx, "purged key", "key", k)
	}

	return nil
}

func (n *natsKeyValue) Set(ctx context.Context, key string, value []byte, ttl int64) error {
	i := Item{
		Value: value,
		TTL:   ttl, // this should already be in unix time if its more than 0
	}

	b, err := json.Marshal(i)
	if err != nil {
		n.logger.ErrorContext(ctx, "failed to marshal item", "key", key, "error", err.Error())

		return err
	}

	if _, err := n.bucket.Put(ctx, key, b); err != nil {
		n.logger.ErrorContext(ctx, "failed to set key", "key", key, "error", err.Error())
		return err
	}

	n.logger.InfoContext(ctx, "set key", "key", key, "ttl", ttl)

	return nil
}

func (n *natsKeyValue) Purge(ctx context.Context, prefix string) error {
	keys, err := n.bucket.Keys(ctx)
	if err != nil {
		n.logger.ErrorContext(ctx, "failed to get keys", "error", err.Error())

		return err
	}

	for _, key := range keys {
		if prefix != "" {
			if strings.HasPrefix(key, prefix) {
				if err := n.bucket.Delete(ctx, key); err != nil {
					n.logger.ErrorContext(ctx, "failed to delete key", "key", key, "error", err.Error())

					return err
				}
			}
		} else {
			if err := n.bucket.Delete(ctx, key); err != nil {
				n.logger.ErrorContext(ctx, "failed to delete key", "key", key, "error", err.Error())

				return err
			}
		}
	}

	return nil
}

func (n *natsKeyValue) Delete(ctx context.Context, key string) error {
	if err := n.bucket.Delete(ctx, key); err != nil {
		n.logger.ErrorContext(ctx, "failed to delete key", "key", key, "error", err.Error())

		return err
	}

	n.logger.InfoContext(ctx, "deleted key", "key", key)

	return nil
}

// NewNATSKeyValue returns a new instance of a natsKeyValue.
func NewNATSKeyValue(bucket jetstream.KeyValue, logger *slog.Logger) Store {
	return &natsKeyValue{
		bucket: bucket,
		logger: logger,
	}
}
