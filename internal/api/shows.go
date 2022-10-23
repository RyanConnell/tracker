package server

import (
	"net/http"
)

type ShowHandler struct {}

func (h *ShowHandler) List(w http.ResponseWriter, r *http.Request) {
	serveError(w, 400, "unimplemented")
}

func (h *ShowHandler) Info(w http.ResponseWriter, r *http.Request) {
	serveError(w, 400, "unimplemented")

}

func (h *ShowHandler) Schedule(w http.ResponseWriter, r *http.Request) {
	serveError(w, 400, "unimplemented")
}
