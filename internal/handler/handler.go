package handler

import (
	"github.com/wrytehq/wryte/internal/config"
	"github.com/wrytehq/wryte/internal/database"
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
