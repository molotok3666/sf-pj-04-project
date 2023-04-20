package api

import (
	"APIGateway/pkg/storage"
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Pagination struct {
	Page       int
	PagesTotal int
	PerPage    int
}

type ResponseWithPagination struct {
	News       []storage.News
	Pagination Pagination
}

type ResponseDetailNewsWithComments struct {
	News     storage.News
	Comments []storage.Comment
}

func NewPagination(page int, count int, perPage int) Pagination {
	var pagesTotal int
	if count != 0 {
		pagesTotal = int(math.Ceil(float64(count / perPage)))
	}

	return Pagination{
		page,
		pagesTotal,
		perPage,
	}
}

func NewResponseWithPagination(n []storage.News, p Pagination) ResponseWithPagination {
	return ResponseWithPagination{n, p}
}

func NewResponseDetailNewsWithComments(n storage.News, c []storage.Comment) ResponseDetailNewsWithComments {
	return ResponseDetailNewsWithComments{n, c}

}

// Регистрация обработчиков API.
func (api *API) newsEndpoints() {
	api.router.HandleFunc("/news/{id}", api.newsDetailHandler).Methods(http.MethodGet, http.MethodOptions)
	api.router.HandleFunc("/news/", api.newsHandler).Methods(http.MethodGet, http.MethodOptions)
}

// Получение новостей из БД.
func (api *API) newsHandler(w http.ResponseWriter, r *http.Request) {
	pageParam := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		page = 1
	}

	s := r.URL.Query().Get("s")
	news, totalPages, err := api.dbNews.News(page, s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(NewResponseWithPagination(news, NewPagination(page, totalPages, storage.NEWS_PAGE_LIMIT)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
	w.Header().Set("Content-Type", "application/json")
}

// Получение новости из БД.
func (api *API) newsDetailHandler(w http.ResponseWriter, r *http.Request) {
	id32, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if id32 <= 0 {
		http.Error(w, "ID must be positive", http.StatusBadRequest)
		return
	}

	id64 := uint64(id32)

	n, err := api.dbNews.NewsDetail(id64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c, err := api.dbComments.Comments(id64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(NewResponseDetailNewsWithComments(n, c))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
	w.Header().Set("Content-Type", "application/json")
}
