package deliveryhttp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetStats(t *testing.T) {
	chiRouter, _, err := initTestServer()
	require.NoError(t, err)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/internal/stats", nil)
	req.Header.Set("X-Real-IP", "127.0.0.1")
	chiRouter.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

}
