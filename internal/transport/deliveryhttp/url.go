package deliveryhttp

import (
	"io"
	"log"
	"net/http"

	"github.com/RedWood011/ServiceURL/internal/entities"
	"github.com/go-chi/chi/v5"
)

type URL struct {
	FullURL string `json:"fullURL"`
	ID      string
}

type PostBatchShortURLJSONBody struct {
	URLs []URL `json:"URLs"`
}

func (u PostBatchShortURLJSONBody) toEntity() []entities.URL {
	urls := make([]entities.URL, 0, len(u.URLs))
	for _, url := range u.URLs {
		urls = append(urls, entities.URL{
			FullURL: url.FullURL,
		})
	}
	return urls
}

func (r *Router) GetTextURLByID(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	urlID := chi.URLParam(request, "id")
	if len(urlID) == 0 {
		http.Error(writer, "Emplty urlID", http.StatusBadRequest)
		return
	}
	fullURL, err := r.service.GetURLByID(ctx, urlID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}
	writer.Header().Set("Location", fullURL)
	writer.WriteHeader(http.StatusTemporaryRedirect)
}

func (r *Router) PostBatchURLJSON(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	var urlsFull PostBatchShortURLJSONBody

	err := readBody(request.Body, &urlsFull)
	if err != nil {
		writeProcessBodyError(ctx, writer, err)
	}

	urls := urlsFull.toEntity()
	if len(urls) == 0 {
		writeSpecifiedError(ctx, writer, err)
	}
	createdIDs, err := r.service.CreateShortURL(ctx, urls)
	if err != nil {
		writeSpecifiedError(ctx, writer, err)
	}
	writeSuccessful(ctx, writer, http.StatusCreated, batchCreatedItemsFromService(createdIDs))
}

func (router *Router) PostBatchURLText(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		//http.Error(w, "Wrong with request", http.StatusBadRequest)
		writeSpecifiedError(ctx, w, err)
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Println(err)
		}
	}()

	if len(body) == 0 {
		writeProcessBodyError(ctx, w, err)
		return
	}

	createdIDs, err := router.service.CreateShortURL(ctx, []entities.URL{{
		ID:      "",
		FullURL: string(body),
	}})
	if err != nil {
		writeSpecifiedError(ctx, w, err)
	}
	//TODO Вынести в helpers
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(createdIDs[0]))

	if err != nil {
		log.Println(err)
	}
}
