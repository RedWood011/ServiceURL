package deliveryhttp

import (
	"net/http"

	"github.com/RedWood011/ServiceURL/internal/entities"
	"github.com/go-chi/chi/v5"
)

type Url struct {
	FullUrl string `json:"fullUrl"`
	ID      string
}

type PostBatchShortUrlJSONBody struct {
	Urls []Url `json:"urls"`
}

func (u PostBatchShortUrlJSONBody) toEntity() []entities.Url {
	urls := make([]entities.Url, 0, len(u.Urls))
	for _, url := range u.Urls {
		urls = append(urls, entities.Url{
			FullUrl: url.FullUrl,
		})
	}
	return urls
}

func (r *Router) GetFullUrlByID(writer http.ResponseWriter, request *http.Request) {
	urlID := chi.URLParam(request, "id")
	ctx := request.Context()
	if len(urlID) == 0 {
		http.Error(writer, "Emplty urlID", http.StatusBadRequest)
		return
	}
	fullUrl, err := r.service.GetUrlByID(ctx, urlID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
	}
	writer.Header().Set("Location", fullUrl)
	writer.WriteHeader(http.StatusTemporaryRedirect)
}

func (r *Router) PostShortUrl(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	var urlsFull PostBatchShortUrlJSONBody
	err := readBody(request.Body, &urlsFull)
	if err != nil {
		writeProcessBodyError(ctx, writer, err)
	}
	urls := urlsFull.toEntity()
	if len(urls) == 0 {
		writeSpecifiedError(ctx, writer, err)
	}
	createdIDs, err := r.service.CreateShortUrl(ctx, urls)
	if err != nil {
		writeSpecifiedError(ctx, writer, err)
	}
	//writeSuccessful(ctx, writer, http.StatusCreated, batchCreatedItemsFromService(createdIDs))
	writeSuccessful(ctx, writer, http.StatusCreated, createdIDs)
}
