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

func TestTextCreateURLOk(t *testing.T) {
	// initial preparations
	r, err := initTestEnv()
	assert.NoError(t, err)
	fullURL := "https://www.google.com/?safe=active&ssui=on"
	uuid := "ff8f5d2a-56c1-4c24-b0e6-44ed243f5d4d"
	URL := createTextShortURL(t, r, fullURL, uuid)

	adr, err := url.Parse(URL)
	assert.NoError(t, err)
	id := strings.ReplaceAll(adr.Path, "/", "")

	actual := getFullURLByID(t, r, id, uuid)
	// get result

	assert.Equal(t, fullURL, actual)
}

func TestJSONCreateSingleURLOk(t *testing.T) {
	// initial preparations
	r, err := initTestEnv()
	assert.NoError(t, err)

	fullURL := "https://www.google.com/?safe=active&ssui=on"
	uuid := "ff8f5d2a-56c1-4c24-b0e6-44ed243f5d4d"

	createdShortURL := createTextShortURL(t, r, fullURL, uuid)

	adr, err := url.Parse(createdShortURL)
	assert.NoError(t, err)

	id := strings.ReplaceAll(adr.Path, "/", "")
	// get result
	actual := getFullURLByID(t, r, id, uuid)

	assert.Equal(t, fullURL, actual)
}

func TestGetAllURLsByUserID(t *testing.T) {
	r, err := initTestEnv()
	assert.NoError(t, err)
	uuid := "ff8f5d2a-56c1-4c24-b0e6-44ed243f5d4d"
	URLs := []string{"https://www.google.com/?safe=active&ssui=on", "www.vk.com"}

	getAllURLsUserID := func(create func(t *testing.T, router *deliveryhttp.Router, fullURL, uuid string) string, t *testing.T,
		router *deliveryhttp.Router, fullURL, uuid string) deliveryhttp.GetAllURLsUserID {
		shortURL := create(t, r, fullURL, uuid)
		return deliveryhttp.GetAllURLsUserID{
			ShortURL:    shortURL,
			OriginalURL: fullURL,
		}
	}

	getWant := func(URLs []string) []deliveryhttp.GetAllURLsUserID {
		expected := make([]deliveryhttp.GetAllURLsUserID, 0, 2)
		expected = append(expected, getAllURLsUserID(createJSONSingleShortURL, t, r, URLs[0], uuid))
		expected = append(expected, getAllURLsUserID(createTextShortURL, t, r, URLs[1], uuid))
		return expected
	}

	testTable := []struct {
		name      string
		want      []deliveryhttp.GetAllURLsUserID
		getUserID string
		code      int
	}{
		{
			name:      "ExistAllURLsUserID",
			want:      getWant(URLs),
			getUserID: uuid,
			code:      http.StatusOK,
		},
		{
			name:      "ErrNoContent",
			want:      []deliveryhttp.GetAllURLsUserID{},
			getUserID: "a144253e-0192-4dbe-96b2-bd80ebe45386",
			code:      http.StatusNoContent,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			URLs, Code := getAllURLsByUserID(t, r, testCase.getUserID)
			assert.Equal(t, URLs, testCase.want)
			assert.Equal(t, Code, testCase.code)
		})
	}

}

func createJSONSingleShortURL(t *testing.T, router *deliveryhttp.Router, url, uuid string) string {
	expected := deliveryhttp.URL{FullURL: url}

	r, w := newReqResp(http.MethodPost, createReqBody(t, expected))

	ctx := r.Context()
	ctx = context.WithValue(ctx, "uuid", uuid)

	r = r.WithContext(ctx)

	router.PostBatchSingleURLJSON(w, r)

	require.Equal(t, http.StatusCreated, w.Code)
	require.Equal(t, w.Header().Get("Content-Type"), "application/json")

	var createdItem deliveryhttp.CreatedItem

	parseRespBody(t, w.Body.Bytes(), &createdItem)
	require.NotEqual(t, createdItem.ID, "")

	return createdItem.ID
}

func createTextShortURL(t *testing.T, router *deliveryhttp.Router, fullURL, uuid string) string {
	// request execution

	r := httptest.NewRequest(http.MethodPost, "/anything", bytes.NewBuffer([]byte(fullURL)))
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	ctx := r.Context()
	ctx = context.WithValue(ctx, "uuid", uuid)

	r = r.WithContext(ctx)
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	router.PostBatchURLText(w, r)

	require.Equal(t, http.StatusCreated, w.Code)
	require.Equal(t, w.Header().Get("Content-Type"), "text/plain")
	// get results

	var createdItem string
	body, err := io.ReadAll(w.Body)
	assert.NoError(t, err)
	createdItem = string(body)
	assert.NoError(t, r.Body.Close())

	return createdItem
}

func getFullURLByID(t *testing.T, router *deliveryhttp.Router, shortURL, uuid string) string {
	r, w := newReqResp(http.MethodGet, nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, "uuid", uuid)
	r = r.WithContext(ctx)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", shortURL)
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	router.GetURLByIDText(w, r)

	require.Equal(t, 307, w.Code)
	fullURL := w.Header().Get("Location")
	return fullURL
}

func getAllURLsByUserID(t *testing.T, router *deliveryhttp.Router, uuid string) ([]deliveryhttp.GetAllURLsUserID, int) {
	r, w := newReqResp(http.MethodGet, nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, "uuid", uuid)
	r = r.WithContext(ctx)

	router.GetUserURLsJSON(w, r)
	require.Equal(t, w.Header().Get("Content-Type"), "application/json")
	var createdItems []deliveryhttp.GetAllURLsUserID
	parseRespBody(t, w.Body.Bytes(), &createdItems)
	return createdItems, w.Code
}
