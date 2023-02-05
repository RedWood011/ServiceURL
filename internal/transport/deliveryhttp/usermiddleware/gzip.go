package usermiddleware

import (
	"compress/gzip"
	"net/http"
	"strings"
)

func GzipHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Type"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		//TODO стоит в дальшнейшем подумать, как не принимать большой объем данных
		body, err := gzip.NewReader(r.Body)
		defer r.Body.Close()

		if err != nil {
			http.Error(w, "failed to gzip", http.StatusInternalServerError)
			return
		}

		r.Body = body
		next.ServeHTTP(w, r)
	})
}
