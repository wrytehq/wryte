package server

import (
	"io/fs"
	"net/http"

	"github.com/wrytehq/wryte/internal/handler"
	"github.com/wrytehq/wryte/internal/middleware"
	"github.com/wrytehq/wryte/web"
)

func (s *Server) Routes(h *handler.Handler) http.Handler {
	r := http.NewServeMux()

	assetsFS, err := fs.Sub(web.Files, "assets")
	if err != nil {
		panic(err)
	}
	r.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(assetsFS))))

	r.HandleFunc("/", h.Home())
	r.HandleFunc("GET /sign-in", h.SignInPage())

	if s.config.IsSelfHosted() && !s.config.IsCloud() {
		r.HandleFunc("GET /setup", h.SetupPage())
		r.HandleFunc("POST /setup", h.SetupForm())
	}

	return middleware.Chain(
		r,
		middleware.Recovery,
		middleware.Logger,
		middleware.Cors,
	)
}
