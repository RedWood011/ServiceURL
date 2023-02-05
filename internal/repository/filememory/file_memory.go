package filememory

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type FileStorage struct {
	m             sync.Mutex
	cacheShortURL map[string]string
	filepath      string
}

func NewFileStorage(path string) (*FileStorage, error) {
	cacheShortURL := make(map[string]string)
	file, err := os.Stat(path)

	if errors.Is(err, os.ErrNotExist) || file.Size() == 0 {
		return &FileStorage{
			filepath:      path,
			cacheShortURL: cacheShortURL,
		}, nil
	}

	fp, err := os.Open(path)
	if err != nil {
		return &FileStorage{}, err
	}

	err = json.NewDecoder(fp).Decode(&cacheShortURL)
	if err != nil {
		return &FileStorage{}, err
	}

	err = fp.Close()
	if err != nil {
		return &FileStorage{}, err
	}

	return &FileStorage{
		cacheShortURL: cacheShortURL,
		filepath:      path,
	}, nil
}

func (s *FileStorage) SaveDone() error {
	file, err := os.Create(s.filepath)
	if err != nil {
		return err
	}

	err = json.NewEncoder(file).Encode(s.cacheShortURL)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}
