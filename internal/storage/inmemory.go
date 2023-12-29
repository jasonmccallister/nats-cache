package storage

import (
	"context"
	"encoding/json"
	"fmt"
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
func (s *inMemory) Get(ctx context.Context, key string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, ok := s.db[key]; !ok {
		return nil, nil
	}

	var i Item
	if err := json.Unmarshal(s.db[key], &i); err != nil {
		return nil, err
	}

	fmt.Println(i)

	// check if the item has expired
	if i.IsExpired() {
		return nil, nil
	}

	return i.Value, nil
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

// NewInMemory returns a new in memory storage engine
func NewInMemory() Store {
	return &inMemory{
		mu: sync.RWMutex{},
		db: make(map[string][]byte),
	}
}
