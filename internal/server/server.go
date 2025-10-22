package server

import (
	"log"
	"net/http"

	"github.com/wrytehq/wryte/internal/config"
	"github.com/wrytehq/wryte/internal/handler"
	"github.com/wrytehq/wryte/internal/templates"
)

type Server struct {
	config *config.Config
}

func New() *http.Server {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Printf("Server starting on %s (env: %s)", cfg.Addr(), cfg.Server.Env)

	tmpl, err := templates.New()
	if err != nil {
		log.Fatalf("Failed to initialize templates: %v", err)
	}
	log.Printf("Templates loaded: %v", tmpl.List())

	h := handler.New(tmpl)

	newServer := &Server{
		config: cfg,
	}

	s := &http.Server{
		Addr:    cfg.Addr(),
		Handler: newServer.Routes(h),
	}

	return s

}
