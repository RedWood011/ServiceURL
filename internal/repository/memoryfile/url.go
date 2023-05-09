package memoryfile

import (
	"context"

	"github.com/RedWood011/ServiceURL/internal/apperror"
	"github.com/RedWood011/ServiceURL/internal/entities"
)

func (f *FileMap) CreateShortURL(_ context.Context, url entities.URL) (string, error) {
	f.m.Lock()
	defer f.m.Unlock()
	// проверка, что не существует такого ID
	exitLongURL := f.cacheLongURL[url.UserID]
	if _, ok := exitLongURL[url.FullURL]; ok {
		return exitLongURL[url.FullURL], apperror.ErrConflict
	}

	f.LongByShortURL[url.ShortURL] = url.FullURL

	//записать  новый URL от ID
	createShortByLong := f.cacheShortURL[url.UserID]
	createLongByShort := f.cacheLongURL[url.UserID]

	if createLongByShort != nil {
		createLongByShort[url.FullURL] = url.ShortURL
		f.cacheLongURL[url.UserID] = createLongByShort
	}

	if createShortByLong != nil {
		createShortByLong[url.ShortURL] = url.FullURL
		f.cacheShortURL[url.UserID] = createShortByLong
		return "", nil
	}
	createLongByShort = make(map[string]string, 1)
	createLongByShort[url.FullURL] = url.ShortURL
	f.cacheLongURL[url.UserID] = createLongByShort

	createShortByLong = make(map[string]string, 1)
	createShortByLong[url.ShortURL] = url.FullURL
	f.cacheShortURL[url.UserID] = createShortByLong

	return "", nil
}

func (f *FileMap) GetFullURLByID(_ context.Context, shortURL string) (res string, err error) {
	f.m.Lock()
	defer f.m.Unlock()

	if f.shortURLByIsDeleted[shortURL] {
		return "", apperror.ErrGone
	}

	if fullURL, ok := f.LongByShortURL[shortURL]; ok {
		return fullURL, nil
	}

	return "", apperror.ErrDataBase
}

func (f *FileMap) GetAllURLsByUserID(_ context.Context, userID string) ([]entities.URL, error) {
	existData, ok := f.cacheShortURL[userID]
	if !ok {
		return nil, apperror.ErrNoContent
	}
	res := make([]entities.URL, 0, len(existData))
	for shortURL, LongURL := range existData {
		res = append(res, entities.URL{
			UserID:   userID,
			ShortURL: shortURL,
			FullURL:  LongURL,
		})
	}
	return res, nil
}

func (f *FileMap) CreateShortURLs(ctx context.Context, urls []entities.URL) ([]entities.URL, error) {
	result := make([]entities.URL, 0, len(urls))
	for _, url := range urls {
		_, err := f.CreateShortURL(ctx, url)
		if err != nil {
			f.rollback(result)
			return nil, err
		}
		result = append(result, entities.URL{
			CorrelationID: url.CorrelationID,
			FullURL:       url.FullURL,
			UserID:        url.UserID,
		})

	}
	return result, nil
}
func (f *FileMap) rollback(urls []entities.URL) {
	if len(urls) == 0 {
		return
	}
	for _, url := range urls {
		exitShortURL := f.cacheShortURL[url.UserID]
		exitLongURL := f.cacheLongURL[url.UserID]
		shortURL := exitLongURL[url.FullURL]
		existLongByShort := f.LongByShortURL

		delete(exitLongURL, url.FullURL)
		delete(exitShortURL, shortURL)
		delete(existLongByShort, shortURL)
		f.cacheShortURL[url.UserID] = exitShortURL
		f.cacheLongURL[url.UserID] = exitLongURL
		f.LongByShortURL = existLongByShort
	}
}

func (f *FileMap) DeleteShortURLs(_ context.Context, urls []string, _ string) error {
	f.m.Lock()
	defer f.m.Unlock()
	for _, short := range urls {
		f.shortURLByIsDeleted[short] = true
	}

	return nil
}
