package deliveryhttp

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/RedWood011/ServiceURL/internal/apperror"
	"github.com/RedWood011/ServiceURL/internal/entities"
	"github.com/RedWood011/ServiceURL/internal/transport/deliveryhttp/usermiddleware"
	"github.com/go-chi/chi/v5"
)

const cookieName usermiddleware.CookieType = "uuid"

type URL struct {
	FullURL string `json:"url"`
	ID      string
}

type PostBatchShortURLJSONBody struct {
	URLs []URL `json:"URLs"`
}

type PostBatchShortURLsJSONBody struct {
	CorrelationID string `json:"correlation_id"`
	FullURL       string `json:"original_url"`
}

type ResponseBatchShortURLsJSONBody struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type GetAllURLsUserID struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func URLsByUserIDFromService(urls []entities.URL) []GetAllURLsUserID {
	res := make([]GetAllURLsUserID, 0, len(urls))

	for _, url := range urls {
		res = append(res, GetAllURLsUserID{
			ShortURL:    url.ShortURL,
			OriginalURL: url.FullURL,
		})
	}
	return res
}

func ResponseBatchShortURLsJSONBodyFromService(urls []entities.URL) []ResponseBatchShortURLsJSONBody {
	res := make([]ResponseBatchShortURLsJSONBody, 0, len(urls))

	for _, url := range urls {
		res = append(res, ResponseBatchShortURLsJSONBody{
			ShortURL:      url.ShortURL,
			CorrelationID: url.CorrelationID,
		})
	}
	return res
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

func toEntity(createURLs []PostBatchShortURLsJSONBody, id string) []entities.URL {
	urls := make([]entities.URL, 0, len(createURLs))
	for _, url := range createURLs {
		urls = append(urls, entities.URL{
			UserID:        id,
			FullURL:       url.FullURL,
			CorrelationID: url.CorrelationID,
		})
	}

	return urls
}

// GetUserURLsJSON - получение списка URL пользователя.
// При успешном запросе - код ответа 200 и списко URL пользователя в
// формате GetURL.
// В случае ошибки получение ссылок из базы данных - код ответа 500.
// В случае отсутствия ссылок у пользователя - код ответа 204.

func (r *Router) GetUserURLsJSON(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	id := ctx.Value(cookieName).(string)
	urls, err := r.service.GetAllURLsByUserID(ctx, id)

	var appErr *apperror.AppError
	if errors.As(err, &appErr) {

		if errors.Is(err, apperror.ErrDataBase) {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(err.Error()))
			return
		}

		if errors.Is(err, apperror.ErrNoContent) {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusNoContent)
			er := json.NewEncoder(writer).Encode(URLsByUserIDFromService(urls))
			if er != nil {
				http.Error(writer, "Unmarshalling error", http.StatusBadRequest)
				return
			}
			return
		}
	}
	writeSuccessful(ctx, writer, http.StatusOK, URLsByUserIDFromService(urls))

}

// GetURLByIDText - получение оригинальной ссылки по укороченному URL.
// Обязательный параметр URL - id.
// Если ссылка верная - код ответа 307 и заголовок "location" с искомой ссылкой.
// Если ссылка не найдена - код ответа 404.
// Если ошибка базы данных - код ответа 500

func (r *Router) GetURLByIDText(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	shortURL := chi.URLParam(request, "id")

	if len(shortURL) == 0 {
		http.Error(writer, "Emplty url", http.StatusBadRequest)
		return
	}

	fullURL, err := r.service.GetURLByID(ctx, shortURL)

	var appErr *apperror.AppError
	if err != nil {
		if errors.As(err, &appErr) {
			if errors.Is(err, apperror.ErrNotFound) {
				writer.WriteHeader(http.StatusNotFound)
				writer.Write([]byte(err.Error()))
				return
			}
			if errors.Is(err, apperror.ErrDataBase) {
				writer.WriteHeader(http.StatusInternalServerError)
				writer.Write([]byte(err.Error()))
				return
			}
		}
		//writeSpecifiedError(ctx, writer, err, "text", createdID)
	}

	writer.Header().Set("Location", fullURL)
	writer.WriteHeader(http.StatusTemporaryRedirect)
}

// PostBatchURLText - создание укороченной ссылки.
// Формат запроса - строка с URL (plain text).
// При успешном создании код ответа 201, а так же в ответе будет укороченная ссылка.
// В случае ошибки в формате запроса - код ответа 400.
// В случае ошибки при записи в базу данных - код ответа 500.

func (r *Router) PostBatchURLText(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	ID := ctx.Value(cookieName).(string)

	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, "Something wrong with request", http.StatusBadRequest)
		return
	}
	defer func() {
		if err := request.Body.Close(); err != nil {
			log.Println(err)
		}
	}()

	if len(body) == 0 {
		http.Error(writer, "Request body is empty", http.StatusBadRequest)
		return
	}

	createdID, err := r.service.CreateShortURL(ctx, entities.URL{
		UserID:  ID,
		FullURL: string(body),
	})
	var appErr *apperror.AppError
	if err != nil {
		if errors.As(err, &appErr) {
			if errors.Is(err, apperror.ErrConflict) {
				writer.Header().Set("Content-Type", "text/plain")
				writer.WriteHeader(http.StatusConflict)
				writer.Write([]byte(createdID))
			}

			if errors.Is(err, apperror.ErrDataBase) {
				writer.WriteHeader(http.StatusInternalServerError)
				writer.Write([]byte(err.Error()))
			}
		}
		return
	}

	//TODO Вынести в helpers
	writer.WriteHeader(http.StatusCreated)
	writer.Header().Set("Content-Type", "text/plain")
	_, err = writer.Write([]byte(createdID))

	if err != nil {
		log.Println(err)
	}
}

// PostBatchSingleURLJSON - создание укороченной ссылки.
// Формат запроса PostURL.
// При успешном создании код ответа 201, а так же в ответе будет укороченная ссылка
// в result.
// В случае ошибки в формате запроса - код ответа 400.
// В случае, если такая ссылка уже имеется - код ответа 409.
// В случае ошибки при записи в базу данных - код ответа 500.

func (r *Router) PostBatchSingleURLJSON(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	ID := ctx.Value(cookieName).(string)

	var url URL

	err := readBody(request.Body, &url)
	if err != nil {
		http.Error(writer, "Something wrong with request", http.StatusBadRequest)
		return
	}

	if len(url.FullURL) == 0 {
		http.Error(writer, "Request body is empty", http.StatusBadRequest)
		return
	}

	createdID, err := r.service.CreateShortURL(ctx, entities.URL{
		UserID:  ID,
		FullURL: url.FullURL,
	})

	var appErr *apperror.AppError
	if err != nil {
		if errors.As(err, &appErr) {

			if errors.Is(err, apperror.ErrDataBase) {
				writer.WriteHeader(http.StatusInternalServerError)
				writer.Write([]byte(err.Error()))
				return
			}

			if errors.Is(err, apperror.ErrConflict) {
				writer.Header().Set("Content-Type", "application/json")
				writer.WriteHeader(http.StatusConflict)
				er := json.NewEncoder(writer).Encode(batchCreatedItemFromService(createdID))
				if er != nil {
					http.Error(writer, "Unmarshalling error", http.StatusBadRequest)
					return
				}

				return
			}
		}
	}

	writeSuccessful(ctx, writer, http.StatusCreated, batchCreatedItemFromService(createdID))
}

func (r *Router) PostBatchURLsJSON(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	ID := ctx.Value(cookieName).(string)

	var batch []PostBatchShortURLsJSONBody

	err := readBody(request.Body, &batch)
	if err != nil {
		http.Error(writer, "Something wrong with request", http.StatusBadRequest)
	}

	urls := toEntity(batch, ID)
	if len(urls) == 0 {
		http.Error(writer, "Request body is empty", http.StatusBadRequest)
	}

	var appErr *apperror.AppError
	createdIDs, err := r.service.CreateShortURLs(ctx, urls)
	if err != nil {
		if errors.As(err, &appErr) {
			if errors.Is(err, apperror.ErrDataBase) {
				writer.WriteHeader(http.StatusInternalServerError)
				writer.Write([]byte(err.Error()))
				return
			}
		}
	}
	writeSuccessful(ctx, writer, http.StatusCreated, ResponseBatchShortURLsJSONBodyFromService(createdIDs))
}

func (r *Router) PingDB(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	err := r.service.PingDB(ctx)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}
	writer.WriteHeader(http.StatusOK)
}

func (r *Router) DeleteBatchURLs(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	ID := ctx.Value(cookieName).(string)

	var shortURls []string
	if err := json.NewDecoder(request.Body).Decode(&shortURls); err != nil {
		http.Error(writer, "Something wrong with request", http.StatusBadRequest)
		return
	}
	var URLs []string
	for _, short := range shortURls {
		shortURL, _ := url.Parse(short)
		URLs = append(URLs, strings.TrimLeft(shortURL.Path, "/"))
	}
	r.service.DeleteShortURLs(ctx, URLs, ID)

	writer.WriteHeader(http.StatusAccepted)
}
