package cache

import (
	"context"
	"testing"

	"github.com/RedWood011/ServiceURL/internal/apperror"
	"github.com/RedWood011/ServiceURL/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestGetFullUrl(t *testing.T) {
	testTable := []struct {
		name     string
		fullUrl  string
		shortUrl string
		err      error
	}{
		{
			name:     "ExistUrl",
			fullUrl:  "adsasdasdfasdqwe",
			shortUrl: "adr",
			err:      nil,
		},
		{
			name:     "DoesNotExistUrl",
			fullUrl:  "",
			shortUrl: "dasda",
			err:      apperror.ErrNotFound,
		},
	}

	s := NewUrlStorage()
	err := s.CreateShortUrl(context.Background(), []entities.Url{{
		ID:      "adr",
		FullUrl: "adsasdasdfasdqwe",
	}})
	assert.NoError(t, err)

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			url, err := s.GetFullUrlByID(context.Background(), testCase.shortUrl)
			assert.Equal(t, err, testCase.err)
			assert.Equal(t, testCase.fullUrl, url)
		})
	}
}
