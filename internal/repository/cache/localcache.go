package cache

import (
	"context"
	"sync"

	"github.com/RedWood011/ServiceURL/internal/apperror"
	"github.com/RedWood011/ServiceURL/internal/entities"
)

type UrlStorage struct {
	m             sync.Mutex
	cacheUrl      map[string]string
	cacheShortUrl map[string]string
}

func NewUrlStorage() *UrlStorage {
	return &UrlStorage{
		cacheShortUrl: make(map[string]string),
		cacheUrl:      make(map[string]string),
	}
}

func (s *UrlStorage) CreateShortUrl(ctx context.Context, urls []entities.Url) error {
	s.m.Lock()
	defer s.m.Unlock()
	for _, url := range urls {
		_, ok := s.cacheUrl[url.FullUrl]
		if !ok {
			s.cacheUrl[url.FullUrl] = url.ID
			s.cacheShortUrl[url.ID] = url.FullUrl
		}
	}

	return nil
}

func (s *UrlStorage) GetFullUrlByID(ctx context.Context, id string) (res string, err error) {
	url, ok := s.cacheShortUrl[id]
	if ok {
		return url, nil
	}
	return "", apperror.ErrNotFound
}
func (s *UrlStorage) GetIDsByUrls(ctx context.Context, urls []string) (map[string]string, error) {
	IDs := make(map[string]string, len(urls))
	for _, url := range urls {
		_, ok := s.cacheUrl[url]
		if ok {
			IDs[url] = s.cacheUrl[url]
		}
	}
	if len(IDs) > 0 {
		return IDs, nil
	}
	return nil, apperror.ErrNotFound
}
