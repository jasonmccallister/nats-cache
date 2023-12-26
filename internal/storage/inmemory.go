package storage

import "sync"

type inMemory struct {
	mu sync.RWMutex
	db map[string][]byte
}

// Delete implements Store.
func (s *inMemory) Delete(db uint32, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.db[key]; !ok {
		return nil
	}

	delete(s.db, key)

	return nil
}

// Get implements Store.
func (s *inMemory) Get(db uint32, key string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, ok := s.db[key]; !ok {
		return nil, nil
	}

	return s.db[key], nil
}

// Set implements Store.
func (s *inMemory) Set(db uint32, key string, value []byte, ttl int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.db[key] = value

	return nil
}

// NewInMemory returns a new in memory storage engine
func NewInMemory() Store {
	return &inMemory{
		mu: sync.RWMutex{},
		db: make(map[string][]byte),
	}
}
