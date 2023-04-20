package api

import (
	"APIGateway/pkg/storage"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

var BAD_WORDS = []string{
	"qwerty", "йцукен", "zxvbnm",
}

// Регистрация обработчиков API.
func (api *API) censorEndpoints() {
	api.router.HandleFunc("/censor/", api.checkCommentHandler).Methods(http.MethodPost, http.MethodOptions)
}

func (api *API) checkCommentHandler(w http.ResponseWriter, r *http.Request) {
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

	re := regexp.MustCompile(strings.Join(BAD_WORDS, "|"))
	hasBadWord := re.MatchString(c.Content)

	if hasBadWord {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	w.Write([]byte(""))
}
