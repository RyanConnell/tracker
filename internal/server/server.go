package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Server struct {
	s *http.Server
	r *mux.Router

	log *zap.Logger
}

// Component allows to extend the functionality of the server by adding
// more routes.
type Component interface {
	RegisterHandlers(r *mux.Router)
}

// NewServer creates a new server
func NewServer(components map[string]Component, opts ...Option) *Server {
	s := &Server{
		s: &http.Server{},
		r: mux.NewRouter(),

		log: zap.NewNop(),
	}

	for _, opt := range opts {
		opt(s)
	}

	for prefix, component := range components {
		component.RegisterHandlers(s.r.PathPrefix(prefix).Subrouter())
	}

	return s
}

// Run the server
func (s *Server) Run(port int) error {
	s.s.Handler = LoggingMiddleware(s.log, s.r)

	s.s.Addr = fmt.Sprintf(":%d", port)
	return s.s.ListenAndServe()
}

// Option allows to extend the functionality of the server
type Option func(s *Server)

func Logger(log *zap.Logger) Option {
	return func(s *Server) {
		s.log = log
	}
}

// ServeError returns back a error to be served.
func ServeError(err error, w http.ResponseWriter) {
	fmt.Fprint(w, err)
}
