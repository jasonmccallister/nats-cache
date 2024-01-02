package storage

import (
	"context"
	"time"
)

// Item is a struct that holds the value and ttl of a key.
// Since NATS does not natively support a TTL per key we need to store it in the value.
// See this issue for more details https://github.com/nats-io/nats-server/issues/3251
type Item struct {
	Value []byte `json:"value"`
	TTL   int64  `json:"ttl"`
}

func (i Item) IsExpired() bool {
	if i.TTL == 0 {
		return false
	}

	return i.TTL < time.Now().Unix()
}

// Store is an interface that defines the methods needed to interact with a storage engine such as NATS KV
type Store interface {
	Delete(ctx context.Context, key string) error
	Get(ctx context.Context, key string) ([]byte, int64, error)
	Purge(ctx context.Context, prefix string) error
	Set(ctx context.Context, key string, value []byte, ttl int64) error
}
