package api

import (
	"APIGateway/pkg/storage"
	"APIGateway/pkg/storage/postgres"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/urfave/negroni"

	"github.com/gorilla/mux"
)

const REQUEST_ID = "request_id"

// Программный интерфейс сервера GoNews
type API struct {
	dbNews        storage.NewsInterface
	dbComments    storage.CommentsInterface
	dbRequestLogs storage.RequestLogsInterface
	router        *mux.Router
}

// Конструктор объекта API
func New(db *postgres.DbStorage) *API {
	api := API{
		dbNews:        db,
		dbComments:    db,
		dbRequestLogs: db,
		router:        mux.NewRouter(),
	}

	api.router.Use(api.RequestIdMiddleware)

	api.newsEndpoints()
	api.commentsEndpoints()
	api.censorEndpoints()

	api.router.Use(api.ResponseHeadersMiddleware)
	api.router.Use(api.RequestLoggingMiddleware)

	return &api
}

// Получение маршрутизатора запросов.
// Требуется для передачи маршрутизатора веб-серверу.
func (api *API) Router() *mux.Router {
	return api.router
}

// RequestIdMiddleware
func (api *API) RequestIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := r.URL.Query().Get(REQUEST_ID)
		if requestId == "" {
			requestId := generateRequestId()
			q := r.URL.Query()
			q.Add(REQUEST_ID, requestId)
			r.URL.RawQuery = q.Encode()
		}
		next.ServeHTTP(w, r)
	})
}

// RequestIdMiddleware
func (api *API) RequestLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := negroni.NewResponseWriter(w)
		next.ServeHTTP(lrw, r)
		ip := r.RemoteAddr
		timestamp := uint64(time.Now().Unix())
		requestId := r.URL.Query().Get(REQUEST_ID)
		statusCode := strconv.Itoa(lrw.Status())
		rl := storage.NewRequestLog(ip, timestamp, statusCode, requestId)
		api.dbRequestLogs.AddLog(rl)
	})
}

// ResponseHeadersMiddleware
func (api *API) ResponseHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func generateRequestId() string {
	min := 100000
	max := 999999
	rand.Seed(time.Now().UnixNano())
	r := min + rand.Intn(max-min)
	return strconv.Itoa(r)
}
