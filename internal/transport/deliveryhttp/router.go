package deliveryhttp

import (
	"github.com/RedWood011/ServiceURL/internal/service"
	"github.com/RedWood011/ServiceURL/internal/transport/deliveryhttp/usermiddleware"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/chi/v5"
	_ "github.com/go-chi/chi/v5/middleware"
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

func NewRouter(r chi.Router, serv service.Translation) chi.Router {
	router := &Router{service: serv}
	r.Use(middleware.Compress(compressionLevel))
	r.Use(usermiddleware.GzipHeader)
	r.Get("/{id}", router.GetTextURLByID)
	r.Post("/url", router.PostBatchURLsJSON)
	r.Post("/api/shorten", router.PostBatchSingleURLJSON)
	r.Post("/", router.PostBatchURLText)

	return r
}
