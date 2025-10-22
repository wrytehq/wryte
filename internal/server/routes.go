package server

import (
	"net/http"

	"github.com/wrytehq/wryte/internal/handler"
	"github.com/wrytehq/wryte/internal/middleware"
)

func (s *Server) Routes(h *handler.Handler) http.Handler {
	r := http.NewServeMux()

	r.HandleFunc("/", h.Home())
	r.HandleFunc("/about", h.About())

	return middleware.Chain(
		r,
		middleware.Recovery,
		middleware.Logger,
		middleware.Cors,
	)
}
