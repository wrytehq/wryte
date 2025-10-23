package handler

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/wrytehq/wryte/internal/middleware"
)

type Document struct {
	ID        string
	Title     string
	Content   string
	CreatedAt string
	UpdatedAt string
	UserID    string
}

func (h *Handler) ViewDocument() http.HandlerFunc {
	tmpl := h.templates.MustRender("document")

	return func(w http.ResponseWriter, r *http.Request) {
		// Get document ID from URL path
		documentID := r.PathValue("documentId")
		if documentID == "" {
			http.Error(w, "Document ID is required", http.StatusBadRequest)
			return
		}

		// Get authenticated user ID from context
		userID, ok := middleware.GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Query document from database
		var doc Document
		query := `SELECT id, title, content, created_at, updated_at, user_id
		          FROM documents
		          WHERE id = $1`
		err := h.db.GetDB().QueryRowContext(r.Context(), query, documentID).Scan(
			&doc.ID,
			&doc.Title,
			&doc.Content,
			&doc.CreatedAt,
			&doc.UpdatedAt,
			&doc.UserID,
		)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "Document not found", http.StatusNotFound)
				return
			}
			log.Printf("Error querying document: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Check if user has access to this document
		// For now, simple check: document must belong to the user
		// TODO: Add sharing/permissions logic
		if doc.UserID != userID {
			http.Error(w, "Forbidden - You don't have access to this document", http.StatusForbidden)
			return
		}

		// Render template
		data := map[string]any{
			"Document": doc,
		}

		err = tmpl.ExecuteTemplate(w, "layout.html", data)
		if err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
