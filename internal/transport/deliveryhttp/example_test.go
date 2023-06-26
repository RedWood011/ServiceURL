package deliveryhttp

import "github.com/go-chi/chi/v5"

func ExampleRouter_PostOneURL() {
	router, err := initTestEnv()
	if err != nil {
		return
	}

	rtr := chi.NewRouter()
	rtr.Post("/", router.PostOneURL)
}

func ExampleRouter_PostBatchSingleURL() {
	router, err := initTestEnv()
	if err != nil {
		return
	}

	rtr := chi.NewRouter()
	rtr.Post("/api/shorten", router.PostBatchSingleURL)
}

func ExampleRouter_PostBatchURLs() {
	router, err := initTestEnv()
	if err != nil {
		return
	}

	rtr := chi.NewRouter()
	rtr.Post("/api/shorten/batch", router.PostBatchURLs)
}

func ExampleRouter_GetURLByID() {
	router, err := initTestEnv()
	if err != nil {
		return
	}

	rtr := chi.NewRouter()
	rtr.Post("/{id}", router.GetURLByID)
}

func ExampleRouter_GetUserURLs() {
	router, err := initTestEnv()
	if err != nil {
		return
	}

	rtr := chi.NewRouter()
	rtr.Get("/api/user/urls", router.GetUserURLs)
}
