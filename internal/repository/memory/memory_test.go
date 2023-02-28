package memory

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
	_, err := s.CreateShortURL(context.Background(), entities.URL{
		ShortURL: "adr",
		FullURL:  "adsasdasdfasdqwe",
	})
	assert.NoError(t, err)

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			url, err := s.GetFullURLByID(context.Background(), testCase.shortURL)
			assert.Equal(t, err, testCase.err)
			assert.Equal(t, testCase.fullURL, url)
		})
	}
}

func TestGetAllURLsByUserID(t *testing.T) {

	testTable := []struct {
		name      string
		want      []entities.URL
		getUserID string
		err       error
	}{
		{
			name: "ExistAllURLsUserID",
			want: []entities.URL{{
				UserID:   "5555",
				ShortURL: "1a2a3a4a5a",
				FullURL:  "aaaaaaaaaaa",
			},
				{
					UserID:   "5555",
					ShortURL: "1b2b3b4b5b",
					FullURL:  "bbbbbbbbbbb",
				}},
			err:       nil,
			getUserID: "5555",
		},
		{
			name:      "ErrNoContent",
			want:      nil,
			getUserID: "9999",
			err:       apperror.ErrNoContent,
		},
	}

	s, _ := NewMemoryStorage()

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			if len(testCase.want) > 0 {
				for i := 0; i < len(testCase.want); i++ {
					_, err := s.CreateShortURL(context.Background(), entities.URL{
						UserID:   testCase.want[i].UserID,
						FullURL:  testCase.want[i].FullURL,
						ShortURL: testCase.want[i].ShortURL,
					})
					assert.NoError(t, err)
				}
			}
			urls, err := s.GetAllURLsByUserID(context.Background(), testCase.getUserID)
			assert.Equal(t, urls, testCase.want)
			assert.Equal(t, err, testCase.err)

		})
	}
}
