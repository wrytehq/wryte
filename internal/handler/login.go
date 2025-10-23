package handler

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/wrytehq/wryte/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) LoginPage() http.HandlerFunc {
	tmpl := h.templates.MustRender("auth/login")

	return func(w http.ResponseWriter, r *http.Request) {
		data := map[string]any{
			"Form":         &validator.LoginForm{},
			"Errors":       &validator.ValidationErrors{},
			"IsSelfHosted": h.config.IsSelfHosted(),
			"Flash":        h.GetFlashMessage(w, r),
		}
		err := tmpl.ExecuteTemplate(w, "layout.html", data)
		if err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handler) LoginForm() http.HandlerFunc {
	v := validator.New()
	tmpl := h.templates.MustRender("auth/login")

	return func(w http.ResponseWriter, r *http.Request) {
		var form validator.LoginForm

		// Decode and validate the form
		validationErrs, err := v.DecodeAndValidate(r, &form)
		if err != nil {
			log.Printf("Error decoding/validating form: %v", err)
			http.Error(w, "Error processing form", http.StatusBadRequest)
			return
		}

		// If there are validation errors, render the form with errors
		if validationErrs.HasErrors() {
			data := map[string]any{
				"Errors":       validationErrs,
				"Form":         &form,
				"IsSelfHosted": h.config.IsSelfHosted(),
			}
			if err := tmpl.ExecuteTemplate(w, "login_form", data); err != nil {
				log.Printf("Error rendering template: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}

		// Query user by email
		var userID, passwordHash string
		query := `SELECT id, password_hash FROM users WHERE email = $1`
		err = h.db.GetDB().QueryRowContext(r.Context(), query, form.Email).Scan(&userID, &passwordHash)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// User not found - return generic error for security
				validationErrs.AddError("email", "Invalid credentials")
				data := map[string]any{
					"Errors":       validationErrs,
					"Form":         &form,
					"IsSelfHosted": h.config.IsSelfHosted(),
				}
				tmpl.ExecuteTemplate(w, "login_form", data)
				return
			}
			log.Printf("Error querying user: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// Verify password
		err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(form.Password))
		if err != nil {
			// Password doesn't match - return generic error for security
			validationErrs.AddError("email", "Invalid credentials")
			data := map[string]any{
				"Errors":       validationErrs,
				"Form":         &form,
				"IsSelfHosted": h.config.IsSelfHosted(),
			}
			tmpl.ExecuteTemplate(w, "login_form", data)
			return
		}

		// Start transaction for session creation
		tx, err := h.db.GetDB().BeginTx(r.Context(), nil)
		if err != nil {
			log.Printf("Error starting transaction: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		// Create session
		token := uuid.NewString()
		expiresAt := time.Now().Add(time.Hour * 24 * 7) // 7 days

		sessionQuery := `INSERT INTO sessions (user_id, token, expires_at, created_at)
		                 VALUES ($1, $2, $3, NOW())`
		_, err = tx.ExecContext(r.Context(), sessionQuery, userID, token, expiresAt)
		if err != nil {
			log.Printf("Error creating session: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// Commit transaction
		if err := tx.Commit(); err != nil {
			log.Printf("Error committing transaction: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// Set session cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "wryte_session",
			Value:    token,
			Expires:  expiresAt,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})

		// Redirect to home
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusOK)
	}
}
