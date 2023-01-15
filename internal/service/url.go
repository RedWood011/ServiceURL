package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/RedWood011/ServiceURL/internal/apperror"
	"github.com/RedWood011/ServiceURL/internal/entities"
)

func (s *TranslationServer) GetUrlByID(ctx context.Context, id string) (string, error) {
	return s.repo.GetFullUrlByID(ctx, id)
}

func (s *TranslationServer) GetIDsByUrls(ctx context.Context, urls []string) (map[string]string, error) {
	return s.repo.GetIDsByUrls(ctx, urls)
}

func (s *TranslationServer) CreateShortUrl(ctx context.Context, urls []entities.Url) (IDs []string, err error) {
	idByUrl := make(map[string]string, len(urls))
	urlByID := make(map[string]string, len(urls))
	createUrls := make([]entities.Url, len(urls))
	fullUrls := make([]string, 0, len(urls))
	IDs = make([]string, 0, len(urls))

	for i := 0; i < len(urls); i++ {
		urls[i].GenerateRandomString(6)
		if urls[i].ID == "" {
			return nil, fmt.Errorf("CreateShortUrl.EmptyShortUrl %w", err)
		}

		idByUrl[urls[i].FullUrl] = urls[i].ID
		urlByID[urls[i].ID] = urls[i].FullUrl
		fullUrls = append(fullUrls, urls[i].FullUrl)
		IDs = append(IDs, urls[i].ID)
	}

	existIDs, err := s.GetIDsByUrls(ctx, fullUrls)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return IDs, s.repo.CreateShortUrl(ctx, urls)
		}
		return nil, err
	}

	IDs = nil
	for _, fullUrl := range fullUrls {
		_, ok := existIDs[fullUrl]
		if ok {
			idByUrl[fullUrl] = existIDs[fullUrl]
			IDs = append(IDs, idByUrl[fullUrl])
		} else {
			createUrls = append(createUrls, entities.Url{
				ID:      idByUrl[fullUrl],
				FullUrl: fullUrl,
			})
			IDs = append(IDs, idByUrl[fullUrl])
		}
	}

	return IDs, s.repo.CreateShortUrl(ctx, createUrls)

}
