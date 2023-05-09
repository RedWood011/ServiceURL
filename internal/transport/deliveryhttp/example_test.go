package deliveryhttp

import "github.com/go-chi/chi/v5"

func ExampleRouter_PostBatchURLText() {
	router, err := initTestEnv()
	if err != nil {
		return
	}

	rtr := chi.NewRouter()
	rtr.Post("/", router.PostBatchURLText)
}

func ExampleRouter_PostBatchSingleURLJSON() {
	router, err := initTestEnv()
	if err != nil {
		return
	}

	rtr := chi.NewRouter()
	rtr.Post("/api/shorten", router.PostBatchSingleURLJSON)
}

func ExampleRouter_PostBatchURLsJSON() {
	router, err := initTestEnv()
	if err != nil {
		return
	}

	rtr := chi.NewRouter()
	rtr.Post("/api/shorten/batch", router.PostBatchURLsJSON)
}

func ExampleRouter_GetURLByIDText() {
	router, err := initTestEnv()
	if err != nil {
		return
	}

	rtr := chi.NewRouter()
	rtr.Post("/{id}", router.GetURLByIDText)
}

func ExampleRouter_GetUserURLsJSON() {
	router, err := initTestEnv()
	if err != nil {
		return
	}

	rtr := chi.NewRouter()
	rtr.Get("/api/user/urls", router.GetUserURLsJSON)
}
