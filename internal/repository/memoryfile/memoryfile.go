package memoryfile

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type FileMap struct {
	m              sync.Mutex
	LongByShortURL map[string]string
	cacheShortURL  map[string]map[string]string
	cacheLongURL   map[string]map[string]string
	filepath       string
}

type Params struct {
	UserID        string
	ShortURL      string
	FullURL       string
	CorrelationID string
}

type write struct {
	ShortURL string `json:"short_url"`
	LongURL  string `json:"long_url"`
	UserID   string `json:"userID"`
}

func NewFileMap(path string) (*FileMap, error) {
	if path == "" {
		return &FileMap{
			cacheShortURL:  make(map[string]map[string]string),
			cacheLongURL:   make(map[string]map[string]string),
			LongByShortURL: make(map[string]string),
		}, nil
	}

	longByShortURL := make(map[string]string)
	file, err := os.Stat(path)

	if errors.Is(err, os.ErrNotExist) || file.Size() == 0 {
		return &FileMap{
			filepath:       path,
			LongByShortURL: longByShortURL,
			cacheShortURL:  make(map[string]map[string]string),
			cacheLongURL:   make(map[string]map[string]string),
		}, nil
	}

	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	var writer []write
	err = json.NewDecoder(fp).Decode(&writer)
	if err != nil {
		return nil, err
	}
	cacheShortURL := make(map[string]map[string]string)
	cacheLongURL := make(map[string]map[string]string)
	var longByShort, shortByLong map[string]string
	for i := 0; i < len(writer); i++ {
		if cacheShortURL[writer[i].UserID] == nil && cacheLongURL[writer[i].UserID] == nil {
			longByShort = make(map[string]string, 1)
			shortByLong = make(map[string]string, 1)
		} else {
			longByShort = cacheShortURL[writer[i].UserID]
			shortByLong = cacheLongURL[writer[i].UserID]
		}

		longByShortURL[writer[i].ShortURL] = writer[i].LongURL

		longByShort[writer[i].ShortURL] = writer[i].LongURL
		shortByLong[writer[i].LongURL] = writer[i].ShortURL
		cacheShortURL[writer[i].UserID] = longByShort
		cacheLongURL[writer[i].UserID] = shortByLong

	}

	err = fp.Close()
	if err != nil {
		return nil, err
	}

	return &FileMap{
		LongByShortURL: longByShortURL,
		filepath:       path,
		cacheShortURL:  cacheShortURL,
		cacheLongURL:   cacheLongURL,
	}, nil
}

func (f *FileMap) SaveDone() error {
	if f.filepath == "" {
		return nil
	}

	file, err := os.Create(f.filepath)
	if err != nil {
		return err
	}
	writer := make([]write, 0, len(f.LongByShortURL))
	for userID, cacheShortByLong := range f.cacheShortURL {
		for shortURL, longURL := range cacheShortByLong {
			writer = append(writer, write{
				UserID:   userID,
				ShortURL: shortURL,
				LongURL:  longURL,
			})
		}
	}
	err = json.NewEncoder(file).Encode(writer)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

func (f *FileMap) Ping(ctx context.Context) error {
	return nil
}
