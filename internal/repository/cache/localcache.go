package cache

import (
	"sync"
)

type MemoryStorage struct {
	m             sync.Mutex
	cacheShortURL map[string]string
}

func NewMemoryStorage() (*MemoryStorage, error) {
	return &MemoryStorage{
		cacheShortURL: make(map[string]string),
	}, nil
}

func (s *MemoryStorage) SaveDone() error {
	return nil
}
