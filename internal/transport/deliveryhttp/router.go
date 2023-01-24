package deliveryhttp

import (
	"github.com/RedWood011/ServiceURL/internal/service"

	"github.com/go-chi/chi/v5"
	_ "github.com/go-chi/chi/v5/middleware"
)

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
	r.Get("/{id}", router.GetTextURLByID)
	r.Post("/url", router.PostBatchURLJSON)
	r.Post("/", router.PostBatchURLText)

	return r
}
