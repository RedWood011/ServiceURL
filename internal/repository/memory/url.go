package memory

import (
	"context"

	"github.com/RedWood011/ServiceURL/internal/apperror"
	"github.com/RedWood011/ServiceURL/internal/entities"
)

func (s *MemoryStorage) CreateShortURL(ctx context.Context, url entities.URL) (string, error) {
	s.m.Lock()
	defer s.m.Unlock()
	// проверка, что не существует такого ID
	exitLongURL := s.cacheLongURL[url.UserID]
	if _, ok := exitLongURL[url.FullURL]; ok {
		return exitLongURL[url.FullURL], apperror.ErrConflict
	}

	s.LongByShortURL[url.ShortURL] = url.FullURL

	//записать  новый URL от ID
	createShortByLong := s.cacheShortURL[url.UserID]
	createLongByShort := s.cacheLongURL[url.UserID]

	if createLongByShort != nil {
		createLongByShort[url.FullURL] = url.ShortURL
		s.cacheLongURL[url.UserID] = createLongByShort
	}

	if createShortByLong != nil {
		createShortByLong[url.ShortURL] = url.FullURL
		s.cacheShortURL[url.UserID] = createShortByLong
		return "", nil
	}
	createLongByShort = make(map[string]string, 1)
	createLongByShort[url.FullURL] = url.ShortURL
	s.cacheLongURL[url.UserID] = createLongByShort

	createShortByLong = make(map[string]string, 1)
	createShortByLong[url.ShortURL] = url.FullURL
	s.cacheShortURL[url.UserID] = createShortByLong

	return "", nil
}

func (s *MemoryStorage) GetFullURLByID(ctx context.Context, shortURL string) (res string, err error) {
	s.m.Lock()
	defer s.m.Unlock()

	if fullURL, ok := s.LongByShortURL[shortURL]; ok {
		return fullURL, nil
	}

	return "", apperror.ErrNotFound
}

func (s *MemoryStorage) GetAllURLsByUserID(ctx context.Context, userID string) ([]entities.URL, error) {
	existData, ok := s.cacheShortURL[userID]
	if !ok {
		return nil, apperror.ErrNoContent
	}
	res := make([]entities.URL, 0, len(existData))
	for shortURL, LongURL := range existData {
		res = append(res, entities.URL{
			UserID:   userID,
			ShortURL: shortURL,
			FullURL:  LongURL,
		})
	}
	return res, nil
}

// TODO нормально реализовать, а не заглушка
func (s *MemoryStorage) CreateShortURLs(ctx context.Context, urls []entities.URL) ([]entities.URL, error) {
	return nil, nil
}
