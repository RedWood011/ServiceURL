package deliveryhttp

import "github.com/go-chi/chi/v5"

func ExamplePostBatchURLText() {
	router, err := initTestEnv()
	if err != nil {
		return
	}

	rtr := chi.NewRouter()
	rtr.Post("/", router.PostBatchURLText)
}

func ExamplePostBatchSingleURLJSON() {
	router, err := initTestEnv()
	if err != nil {
		return
	}

	rtr := chi.NewRouter()
	rtr.Post("/api/shorten", router.PostBatchSingleURLJSON)
}

func ExamplePostBatchURLsJSON() {
	router, err := initTestEnv()
	if err != nil {
		return
	}

	rtr := chi.NewRouter()
	rtr.Post("/api/shorten/batch", router.PostBatchURLsJSON)
}

func ExampleGetURLByIDText() {
	router, err := initTestEnv()
	if err != nil {
		return
	}

	rtr := chi.NewRouter()
	rtr.Post("/{id}", router.GetURLByIDText)
}

func ExampleGetUserURLsJSON() {
	router, err := initTestEnv()
	if err != nil {
		return
	}

	rtr := chi.NewRouter()
	rtr.Get("/api/user/urls", router.GetUserURLsJSON)
}
