package deliveryhttp

import (
	"github.com/RedWood011/ServiceURL/internal/service"
	"github.com/RedWood011/ServiceURL/internal/transport/deliveryhttp/usermiddleware"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
)

const compressionLevel = 5

type Router struct {
	service service.Translation
}

func NewRout(service service.Translation) *Router {
	return &Router{
		service: service,
	}
}

// NewRouter Создание маршрутизатора
func NewRouter(r chi.Router, serv service.Translation, key string) chi.Router {
	router := &Router{service: serv}
	logger, _ := zap.NewProduction()
	r.Use(usermiddleware.LoggerMiddleware(logger))
	r.Use(middleware.Compress(compressionLevel))
	r.Use(usermiddleware.GzipHeader)
	r.Use(usermiddleware.Cookie(key))

	r.Get("/{id}", router.GetURLByIDText)
	r.Post("/", router.PostBatchURLText)

	r.Get("/api/user/urls", router.GetUserURLsJSON)

	r.Post("/api/shorten", router.PostBatchSingleURLJSON)

	r.Post("/api/shorten/batch", router.PostBatchURLsJSON)

	r.Get("/ping", router.PingDB)

	r.Delete("/api/user/urls", router.DeleteBatchURLs)
	return r
}
