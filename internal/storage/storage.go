package storage

// Store is an interface that defines the methods needed to interact with a storage engine such as NATS KV
type Store interface {
	Get(db uint32, key string) ([]byte, error)
	Set(db uint32, key string, value []byte, ttl int64) error
	Delete(db uint32, key string) error
}
