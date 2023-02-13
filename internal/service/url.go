package service

import (
	"context"
	"fmt"
	"net/url"

	"github.com/RedWood011/ServiceURL/internal/entities"
)

const numberElement = 6

func (s *TranslationServer) GetURLByID(ctx context.Context, id string) (string, error) {
	return s.Repo.GetFullURLByID(ctx, id)
}

func (s *TranslationServer) CreateShortURL(ctx context.Context, urls []entities.URL) (IDs []string, err error) {
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

		adr.Path = urls[i].ID
		IDs = append(IDs, adr.String())
	}

	return IDs, s.Repo.CreateShortURL(ctx, urls)
}
