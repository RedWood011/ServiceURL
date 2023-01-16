package deliveryhttp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/RedWood011/ServiceURL/internal/apperror"
)

type APIError struct {
	Message string `json:"message"`
}
type CreatedItem struct {
	ID string `json:"id"`
}

func batchCreatedItemsFromService(createdIDs []string) []CreatedItem {
	batchItems := make([]CreatedItem, 0, len(createdIDs))
	for _, value := range createdIDs {
		batchItems = append(batchItems, CreatedItem{ID: value})
	}
	return batchItems
}

func readBody(body io.Reader, receiver interface{}) error {
	return json.NewDecoder(body).Decode(receiver)
}

func writeProcessBodyError(ctx context.Context, w http.ResponseWriter, err error) {
	writeError(ctx, w, http.StatusBadRequest, fmt.Sprintf("wrong request body: %v", err.Error()))
}

func writeSuccessful(ctx context.Context, w http.ResponseWriter, code int, payload interface{}) {

	err := respondWithJSON(w, code, payload)
	if err != nil {
		fmt.Println("write response")
	}
}

func writeSpecifiedError(ctx context.Context, w http.ResponseWriter, err error) {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		if errors.Is(err, apperror.ErrNotFound) {
			writeErrStatus(ctx, w, http.StatusNotFound)
			return
		}
		writeErrStatus(ctx, w, http.StatusBadRequest)
		return
	}
	//writeErrStatus(ctx, w, http.StatusTeapot)
	writeErrStatus(ctx, w, http.StatusBadRequest)
}
func writeErrStatus(ctx context.Context, w http.ResponseWriter, status int) {

	err := respondWithJSON(w, status, APIError{Message: http.StatusText(status)})
	if err != nil {
		fmt.Println("write response")
	}

}

func writeError(ctx context.Context, w http.ResponseWriter, status int, message string) {

	err := respondWithJSON(w, status, APIError{Message: message})
	if err != nil {
		fmt.Println("write response")
	}
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(payload)
}
