package frontend

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Redirect to a certain place
type Redirect struct {
	code int
	to   string
}

func NewRedirect(code int, to string) *Redirect {
	return &Redirect{
		code: code,
		to:   to,
	}
}

func (r *Redirect) RegisterHandlers(router *mux.Router) {
	router.Handle("/", http.RedirectHandler(r.to, r.code))
}
