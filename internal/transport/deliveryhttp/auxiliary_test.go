package deliveryhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/RedWood011/ServiceURL/internal/config"
	"github.com/RedWood011/ServiceURL/internal/repository"
	"github.com/RedWood011/ServiceURL/internal/repository/memoryfile"
	"github.com/RedWood011/ServiceURL/internal/repository/postgres"
	"github.com/RedWood011/ServiceURL/internal/service"
	"github.com/RedWood011/ServiceURL/internal/workers"
	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slog"
)

func initTestEnv() (*Router, error) {
	cfg := &config.Config{
		Address:  "http://localhost:8080/",
		FilePath: "",
	}
	var db *postgres.Repository
	var dbFile *memoryfile.FileMap
	var err error
	var worker workers.WorkerPool

	if cfg.DatabaseDSN != "" {
		db, err = postgres.NewDatabase(context.Background(), cfg.DatabaseDSN, cfg.CountRepetitionBD)
	} else {
		dbFile, err = memoryfile.NewFileMap(cfg.FilePath)
	}
	repo := repository.NewRepository(cfg.DatabaseDSN, db, dbFile)

	logger := slog.New(slog.NewTextHandler(os.Stderr))
	sv := service.New(repo, logger, &worker, cfg.Address)

	router := NewRout(sv)
	return router, err
}

func createReqBody(t *testing.T, raw interface{}) *bytes.Buffer {
	body, err := json.Marshal(raw)
	if err != nil {
		t.Fatalf("marshaling request body with err %v", err)
	}
	return bytes.NewBuffer(body)
}
func cleanup(r *postgres.Repository) error {
	_, err := r.DB.Exec(context.Background(), "TRUNCATE TABLE urls")
	if err != nil {
		return err
	}
	return nil
}

func initTestServer() (chi.Router, *workers.WorkerPool, error) {
	var (
		db     *postgres.Repository
		dbFile *memoryfile.FileMap
		err    error
	)
	cfg := &config.Config{
		Address:  "http://localhost:8080",
		FilePath: "",
		KeyHash:  "7cdb395a-e63e-445f-b2c4-90a400438ee4",
		//DatabaseDSN:       "postgres://qwerty:qwerty@localhost:5438/postgres?sslmode=disable",
		DatabaseDSN:       "",
		CountRepetitionBD: "5",
		NumWorkers:        5,
		SizeBufWorker:     100,
	}

	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stderr))
	if cfg.DatabaseDSN != "" {
		db, err = postgres.NewDatabase(ctx, cfg.DatabaseDSN, cfg.CountRepetitionBD)
	} else {
		dbFile, err = memoryfile.NewFileMap(cfg.FilePath)
	}
	if err != nil {
		return nil, nil, err
	}
	repo := repository.NewRepository(cfg.DatabaseDSN, db, dbFile)
	workerPool := workers.New(cfg.NumWorkers, cfg.SizeBufWorker)

	serv := service.New(repo, logger, workerPool, cfg.Address)
	chiRouter := NewRouter(chi.NewRouter(), serv, cfg.KeyHash)
	if cfg.DatabaseDSN != "" {
		cleanup(db)
	}
	return chiRouter, workerPool, nil
}

func parseRespBody(t *testing.T, body []byte, result interface{}) {
	if len(body) == 0 {
		t.Fatal("no body in response")
	}
	err := json.Unmarshal(body, result)
	if err != nil {
		t.Fatal(err, " on resp body parsing")
	}
}
