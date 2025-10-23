package handler

import (
	"log"
	"net/http"
)

func (h *Handler) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the session cookie
		cookie, err := r.Cookie("wryte_session")
		if err == nil {
			// Delete session from database
			query := `DELETE FROM sessions WHERE token = $1`
			_, err = h.db.GetDB().ExecContext(r.Context(), query, cookie.Value)
			if err != nil {
				log.Printf("Error deleting session: %v", err)
			}
		}

		// Clear the session cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "wryte_session",
			Value:    "",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})

		// Redirect to login
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
