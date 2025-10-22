package handler

import (
	"log"
	"net/http"
)

func (h *Handler) Home() http.HandlerFunc {
	// Obtener el template una sola vez al crear el handler
	tmpl := h.templates.MustRender("home")

	type data struct {
		Title string
		User  string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Ejecutar el template y manejar el error correctamente
		err := tmpl.Execute(w, data{
			Title: "Wryte",
			User:  "Jhon Doe",
		})
		if err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
