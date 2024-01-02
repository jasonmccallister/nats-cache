package storage

import (
	"context"
	"encoding/json"
	"sync"
)

type inMemory struct {
	mu sync.RWMutex
	db map[string][]byte
}

// Delete implements Store.
func (s *inMemory) Delete(ctx context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.db[key]; !ok {
		return nil
	}

	delete(s.db, key)

	return nil
}

// Get implements Store.
func (s *inMemory) Get(ctx context.Context, key string) ([]byte, int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, ok := s.db[key]; !ok {
		return nil, 0, nil
	}

	var i Item
	if err := json.Unmarshal(s.db[key], &i); err != nil {
		return nil, 0, err
	}

	// check if the item has expired
	if i.IsExpired() {
		return nil, 0, nil
	}

	return i.Value, i.TTL, nil
}

// Set implements Store.
func (s *inMemory) Set(ctx context.Context, key string, value []byte, ttl int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	i := Item{
		Value: value,
		TTL:   ttl,
	}

	b, err := json.Marshal(i)
	if err != nil {
		return err
	}

	s.db[key] = b

	return nil
}

func (s *inMemory) Purge(ctx context.Context, prefix string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for k := range s.db {
		if prefix != "" && k == prefix {
			delete(s.db, k)
		} else if prefix == "" {
			delete(s.db, k)
		}
	}

	return nil
}

// NewInMemory returns a new in memory storage engine
func NewInMemory() Store {
	return &inMemory{
		mu: sync.RWMutex{},
		db: make(map[string][]byte),
	}
}
