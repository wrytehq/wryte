package handler

import (
	"log"
	"net/http"
)

func (h *Handler) SignInPage() http.HandlerFunc {
	tmpl := h.templates.MustRender("auth/sign-in")

	return func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.Execute(w, nil)
		if err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
