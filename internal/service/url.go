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

	adr, err := url.Parse(s.address)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(urls); i++ {
		urls[i].GenerateRandomString(numberElement)
		if urls[i].ID == "" {
			return nil, fmt.Errorf("CreateShortUrl.EmptyShortUrl %w", err)
		}

		idByUrl[urls[i].FullUrl] = urls[i].ID
		urlByID[urls[i].ID] = urls[i].FullUrl
		fullUrls = append(fullUrls, urls[i].FullUrl)

		adr.Path = urls[i].ID
		IDs = append(IDs, adr.String())
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
			adr.Path = idByUrl[fullUrl]
			IDs = append(IDs, adr.String())
		} else {
			createUrls = append(createUrls, entities.Url{
				ID:      idByUrl[fullUrl],
				FullUrl: fullUrl,
			})
			adr.Path = idByUrl[fullUrl]
			IDs = append(IDs, adr.String())
		}
	}
	return IDs, s.repo.CreateShortUrl(ctx, createUrls)
}
