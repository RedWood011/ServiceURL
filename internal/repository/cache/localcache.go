package cache

import (
	"context"
	"sync"

	"github.com/RedWood011/ServiceURL/internal/apperror"
	"github.com/RedWood011/ServiceURL/internal/entities"
)

type URLStorage struct {
	m             sync.Mutex
	cacheURL      map[string]string
	cacheShortURL map[string]string
}

func NewURLStorage() *URLStorage {
	return &URLStorage{
		cacheShortURL: make(map[string]string),
		cacheURL:      make(map[string]string),
	}
}

func (s *URLStorage) CreateShortURL(ctx context.Context, urls []entities.URL) error {
	s.m.Lock()
	defer s.m.Unlock()
	for _, url := range urls {
		_, ok := s.cacheURL[url.FullURL]
		if !ok {
			s.cacheURL[url.FullURL] = url.ID
			s.cacheShortURL[url.ID] = url.FullURL
		}
	}

	return nil
}

func (s *URLStorage) GetFullURLByID(ctx context.Context, id string) (res string, err error) {
	url, ok := s.cacheShortURL[id]
	if ok {
		return url, nil
	}
	return "", apperror.ErrNotFound
}
func (s *URLStorage) GetIDsByURLs(ctx context.Context, urls []string) (map[string]string, error) {
	IDs := make(map[string]string, len(urls))
	for _, url := range urls {
		_, ok := s.cacheURL[url]
		if ok {
			IDs[url] = s.cacheURL[url]
		}
	}
	if len(IDs) > 0 {
		return IDs, nil
	}
	return nil, apperror.ErrNotFound
}
