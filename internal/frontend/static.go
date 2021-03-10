package frontend

import (
	"io/fs"
	"net/http"

	"github.com/gorilla/mux"
)

type Static struct {
	fs fs.FS
}

func NewStatic(sfs fs.FS) *Static {
	return &Static{fs: sfs}
}

func (s *Static) RegisterHandlers(r *mux.Router) {
	r.PathPrefix("/").Handler(http.FileServer(http.FS(s.fs)))
}
