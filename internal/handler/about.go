package handler

import (
	"log"
	"net/http"
)

func (h *Handler) About() http.HandlerFunc {
	// Template se obtiene una sola vez al crear el handler
	tmpl := h.templates.MustRender("about")

	type data struct {
		Title string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.Execute(w, data{
			Title: "About Wryte",
		})
		if err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
