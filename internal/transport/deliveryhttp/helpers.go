package deliveryhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// APIError JSON ответ ошибки
type APIError struct {
	Message string `json:"message"`
}

// CreatedItem Структура ответа
type CreatedItem struct {
	ID string `json:"result"`
}

func batchCreatedItemFromService(createdID string) CreatedItem {
	return CreatedItem{ID: createdID}
}

func readBody(body io.Reader, receiver interface{}) error {
	return json.NewDecoder(body).Decode(receiver)
}

func writeSuccessful(ctx context.Context, w http.ResponseWriter, code int, payload interface{}) {

	err := respondWithJSON(w, code, payload)
	if err != nil {
		fmt.Println("write response")
	}
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(payload)
}
