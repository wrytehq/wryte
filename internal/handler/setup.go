package handler

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/wrytehq/wryte/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) SetupPage() func(http.ResponseWriter, *http.Request) {
	tmpl := h.templates.MustRender("auth/setup")

	return func(w http.ResponseWriter, r *http.Request) {
		data := map[string]any{
			"Form":   &validator.SetupForm{},
			"Errors": &validator.ValidationErrors{},
		}
		err := tmpl.ExecuteTemplate(w, "layout.html", data)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handler) SetupForm() http.HandlerFunc {
	v := validator.New()
	tmpl := h.templates.MustRender("auth/setup")

	return func(w http.ResponseWriter, r *http.Request) {
		var form validator.SetupForm

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
				"Errors": validationErrs,
				"Form":   &form,
			}
			if err := tmpl.ExecuteTemplate(w, "setup_form", data); err != nil {
				log.Printf("Error rendering template: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}

		// Generate password hash
		hash, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error generating password hash: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// Start transaction
		tx, err := h.db.GetDB().BeginTx(r.Context(), nil)
		if err != nil {
			log.Printf("Error starting transaction: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		// Insert user and get the ID
		var userID string
		query := `INSERT INTO users (username, email, password_hash, created_at, updated_at)
		          VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id`
		err = tx.QueryRowContext(r.Context(), query, form.Name, form.Email, hash).Scan(&userID)
		if err != nil {
			// Check for duplicate email or username
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				// Unique constraint violation - 23505 is PostgreSQL's unique_violation code
				formErrors := &validator.ValidationErrors{}

				// Check which constraint was violated
				switch pgErr.ConstraintName {
				case "users_email_key":
					formErrors.AddError("email", "This email is already registered")
					log.Printf("Duplicate email attempt: %s", form.Email)
				case "users_username_key":
					formErrors.AddError("name", "This username is already taken")
					log.Printf("Duplicate username attempt: %s", form.Name)
				}

				// Re-render form with error
				data := map[string]any{
					"Errors": formErrors,
					"Form":   &form,
				}
				tmpl.ExecuteTemplate(w, "setup_form", data)
				return
			}

			log.Printf("Error creating user: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		token := uuid.NewString()
		expiresAt := time.Now().Add(time.Hour * 24 * 7)

		sessionQuery := `INSERT INTO sessions (user_id, token, expires_at, created_at)
		                 VALUES ($1, $2, $3, NOW())`
		_, err = tx.ExecContext(r.Context(), sessionQuery, userID, token, expiresAt)
		if err != nil {
			log.Printf("Error creating session: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

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

		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusOK)
	}
}
