package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/RedWood011/ServiceURL/internal/config"
	"github.com/RedWood011/ServiceURL/internal/repository"
	"github.com/RedWood011/ServiceURL/internal/repository/memoryfile"
	"github.com/RedWood011/ServiceURL/internal/repository/postgres"
	"github.com/RedWood011/ServiceURL/internal/service"
	"github.com/RedWood011/ServiceURL/internal/transport/deliveryhttp"
	"golang.org/x/exp/slog"
)

func initTestEnv() (*deliveryhttp.Router, error) {
	cfg := &config.Config{
		Port:     ":8080",
		Address:  "http://localhost:8080/",
		FilePath: "",
	}
	var db *postgres.Repository
	var dbFile *memoryfile.FileMap
	var err error

	if cfg.DatabaseDSN != "" {
		db, err = postgres.NewDatabase(context.Background(), cfg.DatabaseDSN, cfg.CountRepetitionBD)
	} else {
		dbFile, err = memoryfile.NewFileMap(cfg.FilePath)
	}
	repo := repository.NewRepository(cfg.DatabaseDSN, db, dbFile)

	logger := slog.New(slog.NewTextHandler(os.Stderr))
	sv := service.New(repo, logger, cfg.Address)

	router := deliveryhttp.NewRout(sv)
	return router, err
}

func newReqResp(method string, body io.Reader) (*http.Request, *httptest.ResponseRecorder) {
	return httptest.NewRequest(method, "/anything", body), httptest.NewRecorder()
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

func createReqBody(t *testing.T, raw interface{}) *bytes.Buffer {
	body, err := json.Marshal(raw)
	if err != nil {
		t.Fatalf("marshaling request body with err %v", err)
	}
	return bytes.NewBuffer(body)
}
