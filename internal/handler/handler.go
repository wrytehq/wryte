package handler

import "github.com/wrytehq/wryte/internal/templates"

type Handler struct {
	templates *templates.Manager
}

func New(tmpl *templates.Manager) *Handler {
	return &Handler{
		templates: tmpl,
	}
}
