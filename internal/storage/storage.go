package storage

// Item is a struct that holds the value and ttl of a key.
// Since NATS does not natively support a TTL per key we need to store it in the value.
// See this issue for more details https://github.com/nats-io/nats-server/issues/3251
type Item struct {
	Value []byte
	TTL   int64
}

// Store is an interface that defines the methods needed to interact with a storage engine such as NATS KV
type Store interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte, ttl int64) error
	Delete(key string) error
}
