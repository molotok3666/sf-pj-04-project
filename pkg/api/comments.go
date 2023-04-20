package api

import (
	"APIGateway/pkg/storage"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Регистрация обработчиков API.
func (api *API) commentsEndpoints() {
	api.router.HandleFunc("/news/{newsId}/comments/", api.commentsHandler).Methods(http.MethodGet, http.MethodOptions)
	api.router.HandleFunc("/news/{newsId}/comments/", api.addCommentHandler).Methods(http.MethodPost, http.MethodOptions)
}

// Получение комментариев новости
func (api *API) commentsHandler(w http.ResponseWriter, r *http.Request) {
	newsId, err := strconv.ParseUint(mux.Vars(r)["newsId"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	comments, err := api.dbComments.Comments(newsId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bytes, err := json.Marshal(comments)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
	w.Header().Set("Content-Type", "application/json")
}

// Добавление комментария
func (api *API) addCommentHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var c storage.Comment
	err = json.Unmarshal(b, &c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !isValidComment(c) {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	newsId, err := strconv.ParseUint(mux.Vars(r)["newsId"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c.NewsId = uint64(newsId)
	id, err := api.dbComments.AddComment(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
	w.Header().Set("Content-Type", "application/json")
}

func isValidComment(c storage.Comment) bool {
	requestURL := "http://app:8080/censor/"
	jsonBody, err := json.Marshal(c)
	if err != nil {
		return false
	}

	requestBody := bytes.NewReader(jsonBody)
	req, err := http.NewRequest(http.MethodPost, requestURL, requestBody)
	if err != nil {
		return false
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}

	if res.StatusCode != http.StatusOK {
		return false
	}

	return true
}
