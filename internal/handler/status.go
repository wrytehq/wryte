package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

func (h *Handler) StatusPage() http.HandlerFunc {
	tmpl := h.templates.MustRender("status")

	type data struct {
		Error      error
		HealthJSON []byte
		HealthData map[string]string
		DatabaseOK bool
	}

	return func(w http.ResponseWriter, r *http.Request) {
		healthData := h.db.Health()
		health, err := json.Marshal(healthData)

		if err != nil {
			log.Printf("Error marshaling health data: %v", err)
			tmpl.ExecuteTemplate(w, "layout.html", data{
				Error:      err,
				HealthJSON: nil,
				HealthData: nil,
				DatabaseOK: false,
			})
			return
		}

		// Check if database is healthy
		dbOK := healthData["status"] == "up"

		err = tmpl.ExecuteTemplate(w, "layout.html", data{
			Error:      nil,
			HealthJSON: health,
			HealthData: healthData,
			DatabaseOK: dbOK,
		})
		if err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
