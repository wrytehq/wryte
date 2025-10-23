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

		// Guest routes - login
		{
			loginMux := http.NewServeMux()

			loginMux.HandleFunc("GET /", h.LoginPage())
			loginMux.HandleFunc("POST /", h.LoginForm())

			mux.Handle("/login", h.Guest(loginMux))
		}

		// Guest routes - setup (only for self-hosted)
		if s.config.IsSelfHosted() && !s.config.IsCloud() {
			setupMux := http.NewServeMux()

			setupMux.HandleFunc("GET /", h.SetupPage())
			setupMux.HandleFunc("POST /", h.SetupForm())

			mux.Handle("/setup", h.Guest(setupMux))
		}

		// Guest routes - register (only for cloud)
		if !s.config.IsSelfHosted() && s.config.IsCloud() {
			registerMux := http.NewServeMux()

			registerMux.HandleFunc("GET /", h.RegisterPage())
			registerMux.HandleFunc("POST /", h.RegisterForm())

			mux.Handle("/register", h.Guest(registerMux))
		}

		// Authenticated routes
		{
			authenticatedMux := http.NewServeMux()

			authenticatedMux.HandleFunc("GET /{$}", h.Home())
			authenticatedMux.HandleFunc("GET /logout", h.Logout())

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
