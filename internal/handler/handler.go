package handler

import (
	"net/http"

	"github.com/wrytehq/wryte/internal/config"
	"github.com/wrytehq/wryte/internal/database"
	"github.com/wrytehq/wryte/internal/flash"
	"github.com/wrytehq/wryte/internal/middleware"
	"github.com/wrytehq/wryte/internal/templates"
)

type Handler struct {
	templates *templates.Manager
	db        database.Service
	config    *config.Config
}

func New(tmpl *templates.Manager, db database.Service, cfg *config.Config) *Handler {
	return &Handler{
		templates: tmpl,
		db:        db,
		config:    cfg,
	}
}

func (h *Handler) Authenticated(next http.Handler) http.Handler {
	return middleware.Authenticated(h.db)(next)
}

func (h *Handler) Guest(next http.Handler) http.Handler {
	return middleware.Guest(h.db)(next)
}

func (h *Handler) SelfHosted(next http.Handler) http.Handler {
	return middleware.SelfHosted(h.db)(next)
}

func (h *Handler) GetFlashMessage(w http.ResponseWriter, r *http.Request) *flash.Message {
	msg, _ := flash.Get(w, r)
	return msg
}
