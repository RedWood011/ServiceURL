package deliveryhttp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func BenchmarkPostBatchURLs(b *testing.B) {
	b.Run("", func(b *testing.B) {
		chiRouter, _, _ := initTestServer()
		for i := 0; i < b.N; i++ {
			data := []PostBatchShortURLsJSONBody{
				{
					CorrelationID: "e6ae8f2c-8596-4ca2-81d4-17daa467039f",
					FullURL:       "https://www.yandex.ru"},
				{
					CorrelationID: "d424040b-9b16-44b5-be0f-e78968674e9d",
					FullURL:       "https://www.ya.сom",
				},
				{
					CorrelationID: "78022ed0-badc-4e2d-8e5d-8daa7467826e",
					FullURL:       "https://www.jira.com",
				},
			}
			body, _ := json.Marshal(data)
			res := bytes.NewBuffer(body)
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", res)
			chiRouter.ServeHTTP(w, req)
		}
	})
}

func BenchmarkPostBatchSingleURL(b *testing.B) {
	b.Run("", func(b *testing.B) {
		chiRouter, _, _ := initTestServer()
		for i := 0; i < b.N; i++ {

			link := "https://www.gmail.com"
			existURL := URL{FullURL: link}
			body, _ := json.Marshal(existURL)
			res := bytes.NewBuffer(body)
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/shorten", res)
			chiRouter.ServeHTTP(w, req)
		}
	})
}

func BenchmarkPostOneURL(b *testing.B) {
	b.Run("", func(b *testing.B) {
		chiRouter, _, _ := initTestServer()
		for i := 0; i < b.N; i++ {
			link := "https://www.gmail.com"
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(link))
			chiRouter.ServeHTTP(w, req)
		}
	})
}
