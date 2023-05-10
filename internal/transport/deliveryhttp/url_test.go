package deliveryhttp

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/RedWood011/ServiceURL/internal/transport/deliveryhttp/usermiddleware"
	"github.com/stretchr/testify/require"
)

func TestGetUserURLsJSON(t *testing.T) {
	chiRouter, workerPool, err := initTestServer()
	require.NoError(t, err)
	go func() {
		workerPool.Run(context.Background())
	}()

	getCookieByCreateURLs := func(t *testing.T) *http.Cookie {
		body := []PostBatchShortURLsJSONBody{
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

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", createReqBody(t, body))
		chiRouter.ServeHTTP(w, req)
		require.Equal(t, http.StatusCreated, w.Code)
		response := w.Result()
		defer response.Body.Close()
		cookies := response.Cookies()
		for _, cookie := range cookies {
			if cookie.Name == usermiddleware.CookieName {
				return cookie
			}
		}
		return nil
	}

	type expected struct {
		code        int
		contentType string
	}

	tests := []struct {
		name     string
		cookie   *http.Cookie
		expected expected
	}{
		{
			name:   "Received all user links successfully",
			cookie: getCookieByCreateURLs(t),
			expected: expected{
				code:        http.StatusOK,
				contentType: `application/json`,
			},
		},
		{
			name:   "User does not have shortened links",
			cookie: &http.Cookie{},
			expected: expected{
				code:        http.StatusNoContent,
				contentType: `application/json`,
			},
		},
	}

	time.Sleep(3 * time.Second)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
			req.AddCookie(test.cookie)
			chiRouter.ServeHTTP(w, req)
			require.Equal(t, test.expected.code, w.Code)
			require.Equal(t, test.expected.contentType, w.Header().Get("Content-Type"))
		})
	}
}

func TestGetURLByIDText(t *testing.T) {
	chiRouter, _, err := initTestServer()
	require.NoError(t, err)
	link := "https://www.gmail.com"

	createShortURL := func(t *testing.T, link string) string {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(link))
		chiRouter.ServeHTTP(w, req)
		require.Equal(t, http.StatusCreated, w.Code)
		return w.Body.String()
	}

	type expected struct {
		code     int
		Location string
		body     string
	}
	tests := []struct {
		name     string
		shortURL string
		expected expected
	}{
		{
			name:     "Get original URL correctly",
			shortURL: createShortURL(t, link),
			expected: expected{
				code:     http.StatusTemporaryRedirect,
				Location: "Location",
				body:     link,
			},
		},
		{
			name:     "Get original InternalServerError",
			shortURL: "/12345",
			expected: expected{
				code:     http.StatusInternalServerError,
				Location: "",
				body:     "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, test.shortURL, nil)
			chiRouter.ServeHTTP(w, req)
			require.Equal(t, test.expected.code, w.Code)
			require.Equal(t, test.expected.body, w.Header().Get("Location"))
		})
	}
}

func TestPostBatchURLText(t *testing.T) {
	chiRouter, _, err := initTestServer()
	require.NoError(t, err)

	type expected struct {
		code        int
		contentType string
	}
	tests := []struct {
		name     string
		body     string
		expected expected
	}{
		{
			name: "Post Correct URL",
			body: "https://www.google.com",
			expected: expected{
				code:        http.StatusCreated,
				contentType: "text/plain",
			},
		},
		{
			name: "URL is exist in database",
			body: "https://www.google.com",
			expected: expected{
				code:        http.StatusConflict,
				contentType: "text/plain",
			},
		},
		{
			name: "Empty URL body",
			body: "",
			expected: expected{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	cookie := usermiddleware.CreateValidCookie()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.body))
			req.AddCookie(cookie)
			chiRouter.ServeHTTP(w, req)
			require.Equal(t, test.expected.code, w.Code)
			require.Equal(t, test.expected.contentType, w.Header().Get("Content-Type"))
		})
	}

}

func TestPostBatchSingleURLJSON(t *testing.T) {
	chiRouter, _, err := initTestServer()
	require.NoError(t, err)
	link := "https://www.gmail.com"
	existURL := URL{FullURL: link}
	type expected struct {
		code        int
		contentType string
	}
	tests := []struct {
		name     string
		body     URL
		expected expected
	}{
		{
			name: "Post Correct URL",
			body: URL{
				FullURL: "https://www.yandex.ru",
			},
			expected: expected{
				code:        http.StatusCreated,
				contentType: `application/json`,
			},
		},
		{
			name: "URL is exist in database",
			body: URL{
				FullURL: link,
			},
			expected: expected{
				code:        http.StatusConflict,
				contentType: `application/json`,
			},
		},
	}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", createReqBody(t, existURL))
	cookie := usermiddleware.CreateValidCookie()
	req.AddCookie(cookie)
	chiRouter.ServeHTTP(w, req)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w = httptest.NewRecorder()
			req = httptest.NewRequest(http.MethodPost, "/api/shorten", createReqBody(t, test.body))
			req.AddCookie(cookie)
			chiRouter.ServeHTTP(w, req)
			require.Equal(t, test.expected.code, w.Code)
			require.Equal(t, test.expected.contentType, w.Header().Get("Content-Type"))
		})
	}
}

func TestPostBatchURLsJSON(t *testing.T) {
	chiRouter, _, err := initTestServer()
	require.NoError(t, err)

	type expected struct {
		code        int
		contentType string
	}
	tests := []struct {
		name     string
		body     []PostBatchShortURLsJSONBody
		expected expected
	}{
		{
			name: "Post batch many URL success",
			body: []PostBatchShortURLsJSONBody{
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
			},
			expected: expected{
				code:        http.StatusCreated,
				contentType: `application/json`,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", createReqBody(t, test.body))
			chiRouter.ServeHTTP(w, req)
			require.Equal(t, test.expected.code, w.Code)
			require.Equal(t, test.expected.contentType, w.Header().Get("Content-Type"))
		})
	}
}

func TestDeleteBatchURLs(t *testing.T) {
	chiRouter, workerPool, err := initTestServer()
	require.NoError(t, err)
	go func() {
		workerPool.Run(context.Background())
	}()

	getCreatedShortURls := func(t *testing.T) (string, []string, *http.Cookie) {
		body := []PostBatchShortURLsJSONBody{
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

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", createReqBody(t, body))
		chiRouter.ServeHTTP(w, req)
		require.Equal(t, http.StatusCreated, w.Code)
		response := w.Result()
		defer response.Body.Close()
		var cook *http.Cookie
		cookies := response.Cookies()
		for _, cookie := range cookies {
			if cookie.Name == usermiddleware.CookieName {
				cook = cookie
			}
		}

		var object []ResponseBatchShortURLsJSONBody
		parseRespBody(t, w.Body.Bytes(), &object)
		var shortURLs string
		shorts := make([]string, 0, len(object))
		for _, v := range object {
			if shortURLs == "" {
				shortURLs += fmt.Sprintf("[\"%s\"", v.ShortURL)
				shorts = append(shorts, v.ShortURL)
			} else {
				shortURLs += fmt.Sprintf(",\"%s\"", v.ShortURL)
				shorts = append(shorts, v.ShortURL)
			}
		}
		shortURLs += "]"

		return shortURLs, shorts, cook
	}

	shortURLs, shorts, cookie := getCreatedShortURls(t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", strings.NewReader(shortURLs))
	req.AddCookie(cookie)
	chiRouter.ServeHTTP(w, req)
	require.Equal(t, http.StatusAccepted, w.Code)

	time.Sleep(3 * time.Second)

	for _, short := range shorts {
		req = httptest.NewRequest(http.MethodGet, short, nil)
		w = httptest.NewRecorder()
		chiRouter.ServeHTTP(w, req)
		require.Equal(t, http.StatusGone, w.Code)
	}

}

func TestPingDB(t *testing.T) {
	chiRouter, _, err := initTestServer()
	require.NoError(t, err)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	chiRouter.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

}
