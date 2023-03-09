package service

import (
	"context"

	"github.com/RedWood011/ServiceURL/internal/entities"
)

const numberElement = 6

func (s *TranslationServer) GetURLByID(ctx context.Context, shortURL string) (string, error) {
	return s.Repo.GetFullURLByID(ctx, shortURL)
}

func (s *TranslationServer) GetAllURLsByUserID(ctx context.Context, userID string) ([]entities.URL, error) {
	urls, err := s.Repo.GetAllURLsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(urls); i++ {
		urls[i].ShortURL = s.address + urls[i].ShortURL
	}

	return urls, nil
}

func (s *TranslationServer) CreateShortURL(ctx context.Context, url entities.URL) (ID string, err error) {
	url.GenerateRandomString(numberElement)

	shortURL, err := s.Repo.CreateShortURL(ctx, url)
	if err != nil {
		return s.address + shortURL, err
	}

	shortURL = s.address + url.ShortURL
	return shortURL, nil
}
func (s *TranslationServer) PingDB(ctx context.Context) error {
	return s.Repo.Ping(ctx)
}

func (s *TranslationServer) CreateShortURLs(ctx context.Context, urls []entities.URL) (URLs []entities.URL, err error) {
	for i := 0; i < len(urls); i++ {
		urls[i].GenerateRandomString(numberElement)
	}

	createURLs, err := s.Repo.CreateShortURLs(ctx, urls)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(createURLs); i++ {
		createURLs[i].ShortURL = s.address + urls[i].ShortURL
	}

	return createURLs, err
}
