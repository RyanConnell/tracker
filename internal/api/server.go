package server

import (
	"net/http"
)

type HTTPHandler struct {
	showHandler *ShowHandler
}

func (h *HTTPHandler) registerHandlers() {
	r := mux.NewRouter()

	// Register API handlers.
	r.HandleFunc("/v1/shows", h.showHandler.List)
	r.HandleFunc("/v1/shows/{id:[0-9]+}", h.showHandler.Info)
	r.HandleFunc("/v1/schedule/{start:[0-9-]+}/{end:[0-9-]+}", h.showHandler.Schedule)

	http.Handle("/", r)
}

type Error struct {
	message string `json:"error"`
}

func serveError(w http.ResponseWriter, statusCode int, msg string) {
	w.WriteHeader(statusCode)
	w.Write(json.Encode(map[string]string{"error": msg}))
}
