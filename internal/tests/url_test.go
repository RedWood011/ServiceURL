package tests

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/RedWood011/ServiceURL/internal/transport/deliveryhttp"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateShortUrlOk(t *testing.T) {
	// initial preparations
	r := initTestEnv()
	fullUrl := "https://www.google.com/?safe=active&ssui=on"
	urls := createShortUrl(t, fullUrl, r)

	adr, err := url.Parse(urls)
	assert.NoError(t, err)
	id := strings.ReplaceAll(adr.Path, "/", "")

	actual := getFullUrlByID(t, r, id)
	// get result
	assert.Equal(t, fullUrl, actual)
}

func createShortUrl(t *testing.T, fullUrl string, router *deliveryhttp.Router) string {
	expected := deliveryhttp.PostBatchShortUrlJSONBody{
		Urls: []deliveryhttp.Url{
			{
				FullUrl: fullUrl,
				ID:      "",
			},
		},
	}

	// request execution
	r, w := newReqResp(http.MethodPost, createReqBody(t, expected))
	router.PostShortUrl(w, r)
	require.Equal(t, http.StatusCreated, w.Code)

	// get results
	var createdItem []deliveryhttp.CreatedItem

	parseRespBody(t, w.Body.Bytes(), &createdItem)

	return createdItem[0].Id
}

func getFullUrlByID(t *testing.T, router *deliveryhttp.Router, id string) string { //[]deliveryhttp.Url
	ids := make([]string, 0, 0)
	ids = append(ids, id)
	r, w := newReqResp(http.MethodGet, nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id)
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	router.GetFullUrlByID(w, r)

	require.Equal(t, 307, w.Code)
	fullUrl := w.Header().Get("Location")
	return fullUrl

}
