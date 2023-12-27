package storage

// Store is an interface that defines the methods needed to interact with a storage engine such as NATS KV
type Store interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte, ttl int64) error
	Delete(key string) error
}
