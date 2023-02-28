package usermiddleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMiddlewareGetCookie(t *testing.T) {
	key := "50526924-9432-4c7c-a32d-324643a7b927"
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello world")
	})
	rr := httptest.NewRecorder()

	cookieMiddleware := Cookie(key)
	cookie := cookieMiddleware(handler)
	cookie.ServeHTTP(rr, req)
	authorizationFirst := rr.Header().Get("Set-Cookie")

	rq, err := http.NewRequest("GET", "/", nil)
	r := httptest.NewRecorder()
	assert.NoError(t, err)
	//устанавливаем куки
	rq.Header.Set("Cookie", authorizationFirst)
	//делаем новый запрос с установленными cookie
	cookieMiddleware = Cookie(key)
	cookie = cookieMiddleware(handler)
	cookie.ServeHTTP(r, rq)
	//если cookie не валидна пришел бы новый Set-Cookie
	assert.Equal(t, r.Header().Get("Set-Cookie"), "")

}
