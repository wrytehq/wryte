package handler

import (
	"github.com/wrytehq/wryte/internal/database"
	"github.com/wrytehq/wryte/internal/templates"
)

type Handler struct {
	templates *templates.Manager
	db        database.Service
}

func New(tmpl *templates.Manager, db database.Service) *Handler {
	return &Handler{
		templates: tmpl,
		db:        db,
	}
}
