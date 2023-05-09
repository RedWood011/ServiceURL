package service

import (
	"context"

	"github.com/RedWood011/ServiceURL/internal/entities"
)

const numberElement = 6
const sizeDeleted = 50

// GetURLByID Получить оригинальную ссылку по shortURL
func (s *TranslationServer) GetURLByID(ctx context.Context, shortURL string) (string, error) {
	return s.Repo.GetFullURLByID(ctx, shortURL)
}

// GetAllURLsByUserID Получить все оригинальные ссылки по userID
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

// CreateShortURL Создать короткую ссылку
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

// CreateShortURLs Создать короткие ссылки
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

// DeleteShortURLs Удалить короткие ссылки
func (s *TranslationServer) DeleteShortURLs(ctx context.Context, urls []string, userID string) {
	if ctx.Err() != nil {
		return
	}
	batchDeleted := splitURLs(urls, sizeDeleted)
	for _, urlsDeleted := range batchDeleted {
		s.wp.Add(func(ctx context.Context) error {
			err := s.Repo.DeleteShortURLs(ctx, urlsDeleted, userID)
			return err
		})
	}
}

// splitURLs -.
func splitURLs(urls []string, size int) [][]string {
	res := make([][]string, 0, len(urls)/size+1)
	url := make([]string, 0, size)
	count := 0
	for i := 0; i < len(urls); i++ {
		if count < size {
			url = append(url, urls[i])
			count++
		}

		if count == size {
			res = append(res, url)
			url = nil
			count = 0
		}

		if i == len(urls)-1 && count < size && count != 0 {
			res = append(res, url)
		}
	}
	return res
}
