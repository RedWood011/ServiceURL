package memoryfile

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
			err:      apperror.ErrDataBase,
		},
	}

	s, _ := NewFileMap("")
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

	s, _ := NewFileMap("")

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

func TestCreateURLs(t *testing.T) {

	testTable := []struct {
		name     string
		urls     []entities.URL
		expected []entities.URL

		err error
	}{
		{
			name: "CreateURLs",
			urls: []entities.URL{{
				UserID:   "5555",
				ShortURL: "1a2a3a4a5a",
				FullURL:  "aaaaaaaaaaa",
			},
				{
					UserID:   "5555",
					ShortURL: "1b2b3b4b5b",
					FullURL:  "bbbbbbbbbbb",
				},
				{
					UserID:   "6666",
					ShortURL: "1b2b3b4b5b",
					FullURL:  "cccccccccc",
				},
				{
					UserID:   "6666",
					ShortURL: "1d2d3d4d5d",
					FullURL:  "dddddddddd",
				}},
			expected: []entities.URL{{
				UserID:  "5555",
				FullURL: "aaaaaaaaaaa",
			},
				{
					UserID:  "5555",
					FullURL: "bbbbbbbbbbb",
				},
				{
					UserID:  "6666",
					FullURL: "cccccccccc",
				},
				{
					UserID:  "6666",
					FullURL: "dddddddddd",
				}},
			err: nil,
		},
		{
			name: "ErrConflict",
			urls: []entities.URL{{
				UserID:   "5555",
				ShortURL: "1a2a3a4a5a",
				FullURL:  "aaaaaaaaaaa",
			},
				{
					UserID:   "5555",
					ShortURL: "1b2b3b4b5b",
					FullURL:  "bbbbbbbbbbb",
				},
				{
					UserID:   "5555",
					ShortURL: "1b2b3b4b5b",
					FullURL:  "bbbbbbbbbbb",
				},
				{
					UserID:   "6666",
					ShortURL: "1b2b3b4b5b",
					FullURL:  "bbbbbbbbbbb",
				}},
			expected: nil,
			err:      apperror.ErrConflict,
		},
	}

	for _, testCase := range testTable {
		s, _ := NewFileMap("")
		t.Run(testCase.name, func(t *testing.T) {
			res, err := s.CreateShortURLs(context.Background(), testCase.urls)
			assert.Equal(t, res, testCase.expected)
			assert.Equal(t, err, err)

		})
	}
}
