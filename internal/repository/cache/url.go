package cache

import (
	"context"

	"github.com/RedWood011/ServiceURL/internal/apperror"
	"github.com/RedWood011/ServiceURL/internal/entities"
)

func (s *MemoryStorage) CreateShortURL(ctx context.Context, urls []entities.URL) error {
	s.m.Lock()
	defer s.m.Unlock()

	for _, url := range urls {
		s.cacheShortURL[url.ID] = url.FullURL
	}

	return nil
}

func (s *MemoryStorage) GetFullURLByID(ctx context.Context, id string) (res string, err error) {
	s.m.Lock()
	defer s.m.Unlock()

	url, ok := s.cacheShortURL[id]
	if ok {
		return url, nil
	}

	return "", apperror.ErrNotFound
}
