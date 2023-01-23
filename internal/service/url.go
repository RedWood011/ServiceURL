package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/RedWood011/ServiceURL/internal/apperror"
	"github.com/RedWood011/ServiceURL/internal/entities"
)

const numberElement = 6

func (s *TranslationServer) GetURLByID(ctx context.Context, id string) (string, error) {
	return s.repo.GetFullURLByID(ctx, id)
}

func (s *TranslationServer) GetIDsByURLs(ctx context.Context, urls []string) (map[string]string, error) {
	return s.repo.GetIDsByURLs(ctx, urls)
}

func (s *TranslationServer) CreateShortURL(ctx context.Context, urls []entities.URL) (IDs []string, err error) {
	idByURL := make(map[string]string, len(urls))
	urlByID := make(map[string]string, len(urls))
	createURLs := make([]entities.URL, 0, len(urls))
	fullURLs := make([]string, 0, len(urls))
	IDs = make([]string, 0, len(urls))

	adr, err := url.Parse(s.address)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(urls); i++ {
		urls[i].GenerateRandomString(numberElement)
		if urls[i].ID == "" {
			return nil, fmt.Errorf("CreateShortURL.EmptyShortURL %w", err)
		}

		idByURL[urls[i].FullURL] = urls[i].ID
		urlByID[urls[i].ID] = urls[i].FullURL
		fullURLs = append(fullURLs, urls[i].FullURL)

		adr.Path = urls[i].ID
		IDs = append(IDs, adr.String())
	}

	existIDs, err := s.GetIDsByURLs(ctx, fullURLs)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return IDs, s.repo.CreateShortURL(ctx, urls)
		}
		return nil, err
	}

	IDs = nil
	for _, fullURL := range fullURLs {
		_, ok := existIDs[fullURL]
		if ok {
			idByURL[fullURL] = existIDs[fullURL]
			adr.Path = idByURL[fullURL]
			IDs = append(IDs, adr.String())
		} else {
			createURLs = append(createURLs, entities.URL{
				ID:      idByURL[fullURL],
				FullURL: fullURL,
			})

			adr.Path = idByURL[fullURL]
			IDs = append(IDs, adr.String())
		}
	}

	return IDs, s.repo.CreateShortURL(ctx, createURLs)
}
