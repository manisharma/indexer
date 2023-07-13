package handler

import (
	"encoding/json"
	"fmt"
	"indexer/pkg/store"
	"log"
	"net/http"
)

type HTTP struct {
	repo store.Repository
}

func New(repo store.Repository) *HTTP {
	return &HTTP{repo}
}

func (h *HTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {
		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(http.StatusText(http.StatusNotFound)))
			return
		case "/":
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}
	epochs, err := h.repo.Get(r.Context())
	if err != nil {
		message := fmt.Sprintf("repo.Get() failed, err: %v", err.Error())
		log.Print(message)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(message))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	if len(epochs) > 0 {
		encoder.Encode(epochs)
	} else {
		encoder.Encode(map[string]string{"message": "no blocks yet"})
	}
}
