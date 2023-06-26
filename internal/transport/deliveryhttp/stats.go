package deliveryhttp

import (
	"net"
	"net/http"
)

// ResponseStats Статистика
type ResponseStats struct {
	CountURL  int `json:"urls"`
	CountUser int `json:"users"`
}

// GetStats Получить статистику по количеству пользователей и сокращенных сылок
func (rout *Router) GetStats(trustedSubnet string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if trustedSubnet == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		address := net.ParseIP(r.Header.Get("X-Real-IP"))
		_, subnet, err := net.ParseCIDR(trustedSubnet)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		if !subnet.Contains(address) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		ctx := r.Context()

		stats, err := rout.service.GetAllStats(ctx)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		writeSuccessful(ctx, w, http.StatusOK, ResponseStats{CountUser: stats.CountUser, CountURL: stats.CountURL})
	})
}
