package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RedWood011/ServiceURL/internal/config"
	"github.com/RedWood011/ServiceURL/internal/repository"
	"github.com/RedWood011/ServiceURL/internal/service"
	"github.com/RedWood011/ServiceURL/internal/transport/deliveryhttp"
)

func initTestEnv() (*deliveryhttp.Router, error) {
	cfg := &config.Config{
		Port:     ":8080",
		Address:  "http://localhost:8080/",
		FilePath: "",
	}
	repo, err := repository.NewRepository(cfg)

	sv := service.New(repo, cfg.Address)
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
