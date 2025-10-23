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

	// Public assets - no auth required
	assetsFS, err := fs.Sub(web.Files, "assets")
	if err != nil {
		panic(err)
	}
	r.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(assetsFS))))

	{
		mux := http.NewServeMux()

		// Public routes - status page (only for self-hosted)
		if s.config.IsSelfHosted() && !s.config.IsCloud() {
			mux.HandleFunc("GET /status", h.StatusPage())
		}

		// Guest routes - login
		{
			loginMux := http.NewServeMux()

			loginMux.HandleFunc("GET /", h.LoginPage())
			loginMux.HandleFunc("POST /", h.LoginForm())

			mux.Handle("/login", h.Guest(loginMux))
		}

		// Guest routes - setup (only for self-hosted)
		if s.config.IsSelfHosted() && !s.config.IsCloud() {
			guestMux := http.NewServeMux()
			guestMux.HandleFunc("GET /", h.SetupPage())
			guestMux.HandleFunc("POST /", h.SetupForm())

			mux.Handle("/setup", h.Guest(guestMux))
		}

		// Guest routes - register (only for cloud)
		if !s.config.IsSelfHosted() && s.config.IsCloud() {
			cloudMux := http.NewServeMux()

			cloudMux.HandleFunc("GET /", h.RegisterPage())
			cloudMux.HandleFunc("POST /", h.RegisterForm())

			mux.Handle("/register", h.Guest(cloudMux))
		}

		// Authenticated routes
		{
			authenticatedMux := http.NewServeMux()
			authenticatedMux.HandleFunc("GET /{$}", h.Home())
			authenticatedMux.HandleFunc("GET /logout", h.Logout())
			authenticatedMux.HandleFunc("GET /documents/{documentId}", h.ViewDocument())

			mux.Handle("/", h.Authenticated(authenticatedMux))
		}

		r.Handle("/", mux)
	}

	return middleware.Chain(
		r,
		middleware.Recovery,
		middleware.Logger,
		middleware.Cors,
	)
}
