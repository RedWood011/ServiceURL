package memory

import (
	"context"
	"sync"
)

type MemoryStorage struct {
	m              sync.Mutex
	cacheShortURL  map[string]map[string]string
	cacheLongURL   map[string]map[string]string
	LongByShortURL map[string]string
}

func NewMemoryStorage() (*MemoryStorage, error) {
	return &MemoryStorage{
		cacheShortURL:  make(map[string]map[string]string),
		cacheLongURL:   make(map[string]map[string]string),
		LongByShortURL: make(map[string]string),
	}, nil
}

func (s *MemoryStorage) SaveDone() error {
	return nil
}

func (s *MemoryStorage) Ping(ctx context.Context) error {
	return nil
}
