package server

import (
	"net/http"

	"github.com/aqin97/cache/cache"
)

type Server struct {
	cache.Cache
}

func (s *Server) Listen() {
	http.Handle("/cache/", s.cacheHandler())
	http.Handle("/status", s.statusHandler())
	http.ListenAndServe(":8080", nil)
}

func New(c cache.Cache) *Server {
	return &Server{c}
}
