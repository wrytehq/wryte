package handler

import (
	"log"
	"net/http"
)

func (h *Handler) Home() http.HandlerFunc {
	tmpl := h.templates.MustRender("home")

	return func(w http.ResponseWriter, r *http.Request) {
		data := map[string]any{
			// Add home page data here
		}

		err := tmpl.ExecuteTemplate(w, "layout.html", data)
		if err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
