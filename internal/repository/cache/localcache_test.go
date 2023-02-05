package cache

import (
	"context"
	"testing"

	"github.com/RedWood011/ServiceURL/internal/apperror"
	"github.com/RedWood011/ServiceURL/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestGetFullURL(t *testing.T) {
	testTable := []struct {
		name     string
		fullURL  string
		shortURL string
		err      error
	}{
		{
			name:     "ExistURL",
			fullURL:  "adsasdasdfasdqwe",
			shortURL: "adr",
			err:      nil,
		},
		{
			name:     "DoesNotExistURL",
			fullURL:  "",
			shortURL: "dasda",
			err:      apperror.ErrNotFound,
		},
	}

	s, _ := NewMemoryStorage()
	err := s.CreateShortURL(context.Background(), []entities.URL{{
		ID:      "adr",
		FullURL: "adsasdasdfasdqwe",
	}})
	assert.NoError(t, err)

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			url, err := s.GetFullURLByID(context.Background(), testCase.shortURL)
			assert.Equal(t, err, testCase.err)
			assert.Equal(t, testCase.fullURL, url)
		})
	}
}
