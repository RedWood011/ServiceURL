package usermiddleware

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGzipHeader(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello world")
	})

	req, err := http.NewRequest("GET", "/", nil)

	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	gzipMiddleware := GzipHeader(handler)

	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	_, err = gz.Write([]byte("hello world"))
	assert.NoError(t, err)
	err = gz.Close()
	assert.NoError(t, err)

	req.Header.Set("Content-Encoding", "gzip")
	req.Body = ioutil.NopCloser(bytes.NewReader(b.Bytes()))

	gzipMiddleware.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK)

	expected := "hello world\n"
	assert.Equal(t, rr.Body.String(), expected)
}
