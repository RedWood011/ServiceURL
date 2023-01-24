package tests

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/RedWood011/ServiceURL/internal/transport/deliveryhttp"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateURLTextOk(t *testing.T) {
	// initial preparations
	r := initTestEnv()
	fullURL := "https://www.google.com/?safe=active&ssui=on"
	urls := createShortURL(t, fullURL, r)

	adr, err := url.Parse(urls)
	assert.NoError(t, err)
	id := strings.ReplaceAll(adr.Path, "/", "")

	actual := getFullURLByID(t, r, id)
	// get result
	assert.Equal(t, fullURL, actual)
}

func createShortURL(t *testing.T, fullURL string, router *deliveryhttp.Router) string {
	// request execution
	r := httptest.NewRequest(http.MethodPost, "/anything", bytes.NewBuffer([]byte(fullURL)))
	w := httptest.NewRecorder()
	router.PostBatchURLText(w, r)
	require.Equal(t, http.StatusCreated, w.Code)

	// get results
	var createdItem string
	body, err := io.ReadAll(w.Body)
	assert.NoError(t, err)
	createdItem = string(body)
	assert.NoError(t, r.Body.Close())

	return createdItem
}

func getFullURLByID(t *testing.T, router *deliveryhttp.Router, id string) string {
	r, w := newReqResp(http.MethodGet, nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id)
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	router.GetTextURLByID(w, r)

	require.Equal(t, 307, w.Code)
	fullURL := w.Header().Get("Location")
	return fullURL

}
