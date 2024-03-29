package usermiddleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// CookieType Тип куки
type CookieType string

// Настройки cookie
const (
	CookieName     = "uuid"
	timeSecondLive = 900
)

// Cookie Проверка куки
func Cookie(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var nameCookie CookieType = CookieName

			cookie, err := r.Cookie(CookieName)
			if err != nil {
				uid := setUUIDCookie(w, uuid.NewString(), key)
				ctx := r.Context()

				ctx = context.WithValue(ctx, nameCookie, uid)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			}

			parts := strings.Split(cookie.Value, ":")
			if len(parts) != 2 {
				// если в куки нет обоих параметров, то генерируем новый uid
				uid := setUUIDCookie(w, uuid.NewString(), key)
				ctx := r.Context()
				ctx = context.WithValue(ctx, nameCookie, uid)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			}
			uid, hash := parts[0], parts[1]
			if checkHash(uid, hash, key) {
				ctx := r.Context()
				ctx = context.WithValue(ctx, nameCookie, uid)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return

			}
			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}

func calculateHash(uid, key string) string {
	hash := hmac.New(sha256.New, []byte(key))
	hash.Write([]byte(uid))
	return hex.EncodeToString(hash.Sum(nil))
}

func checkHash(uid, hash, key string) bool {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(uid))
	sign, err := hex.DecodeString(hash)
	if err != nil {
		return false
	}
	return hmac.Equal(sign, h.Sum(nil))
}

func setUUIDCookie(w http.ResponseWriter, uid string, key string) string {
	uuid := fmt.Sprintf("%s:%s", uid, calculateHash(uid, key))

	http.SetCookie(w, &http.Cookie{
		Name:   CookieName,
		Value:  uuid,
		MaxAge: timeSecondLive,
	})
	return uid
}

// CreateValidCookie TODO вподумать над нормальным механизмом передачи ключа
// CreateValidCookie создает куку
func CreateValidCookie() *http.Cookie {
	id := uuid.NewString()
	uuid := fmt.Sprintf("%s:%s", id, calculateHash(id, "7cdb395a-e63e-445f-b2c4-90a400438ee4"))
	return &http.Cookie{
		Name:   CookieName,
		Value:  uuid,
		MaxAge: timeSecondLive,
	}
}
