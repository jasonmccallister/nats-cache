package storage

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/nats-io/nats.go/jetstream"
)

type natsKeyValue struct {
	bucket jetstream.KeyValue
}

func (n *natsKeyValue) Get(ctx context.Context, key string) ([]byte, error) {
	v, err := n.bucket.Get(ctx, key)
	if err != nil {
		if errors.Is(err, jetstream.ErrKeyNotFound) {
			return nil, nil
		}
	}

	var i Item
	if err := json.Unmarshal(v.Value(), &i); err != nil {
		return nil, err
	}

	// check if the item has expired
	if i.IsExpired() {
		return nil, nil
	}

	return i.Value, nil
}

func (n *natsKeyValue) Set(ctx context.Context, key string, value []byte, ttl int64) error {
	i := Item{
		Value: value,
		TTL:   ttl, // this should already be in unix time if its more than 0
	}

	b, err := json.Marshal(i)
	if err != nil {
		return err
	}

	if _, err := n.bucket.Put(ctx, key, b); err != nil {
		return err
	}

	return nil
}

func (n *natsKeyValue) Purge(ctx context.Context, prefix string) error {
	keys, err := n.bucket.Keys(ctx)
	if err != nil {
		return err
	}

	for _, key := range keys {
		if prefix != "" {
			if strings.HasPrefix(key, prefix) {
				if err := n.bucket.Delete(ctx, key); err != nil {
					return err
				}
			}
		} else {
			if err := n.bucket.Delete(ctx, key); err != nil {
				return err
			}
		}
	}

	return nil
}

func (n *natsKeyValue) Delete(ctx context.Context, key string) error {
	return n.bucket.Delete(ctx, key)
}

// NewNATSKeyValue returns a new instance of a natsKeyValue.
func NewNATSKeyValue(bucket jetstream.KeyValue) Store {
	return &natsKeyValue{
		bucket: bucket,
	}
}
