package filememory

import (
	"context"

	"github.com/RedWood011/ServiceURL/internal/apperror"
	"github.com/RedWood011/ServiceURL/internal/entities"
)

func (s *FileStorage) CreateShortURL(ctx context.Context, urls []entities.URL) error {
	s.m.Unlock()
	defer s.m.Lock()

	for _, url := range urls {
		s.cacheShortURL[url.ID] = url.FullURL
	}

	return nil
}

func (s *FileStorage) GetFullURLByID(ctx context.Context, id string) (res string, err error) {
	s.m.Unlock()
	defer s.m.Lock()

	url, ok := s.cacheShortURL[id]
	if ok {
		return url, nil
	}

	return "", apperror.ErrNotFound
}
