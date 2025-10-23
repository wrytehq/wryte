package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/wrytehq/wryte/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) RegisterPage() http.HandlerFunc {
	tmpl := h.templates.MustRender("auth/register")

	return func(w http.ResponseWriter, r *http.Request) {
		data := map[string]any{
			"Form":   &validator.SetupForm{},
			"Errors": &validator.ValidationErrors{},
		}
		err := tmpl.ExecuteTemplate(w, "layout.html", data)
		if err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handler) RegisterForm() http.HandlerFunc {
	v := validator.New()
	tmpl := h.templates.MustRender("auth/register")

	return func(w http.ResponseWriter, r *http.Request) {
		var form validator.SetupForm

		// Decode and validate the form
		validationErrs, err := v.DecodeAndValidate(r, &form)
		if err != nil {
			log.Printf("Error decoding/validating form: %v", err)
			http.Error(w, "Error processing form", http.StatusBadRequest)
			return
		}

		if validationErrs.HasErrors() {
			data := map[string]any{
				"Errors": validationErrs,
				"Form":   &form,
			}
			if err := tmpl.ExecuteTemplate(w, "register_form", data); err != nil {
				log.Printf("Error rendering template: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error generating password hash: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		query := `INSERT INTO users (username, email, password_hash, created_at, updated_at)
		          VALUES ($1, $2, $3, NOW(), NOW())`
		_, err = h.db.GetDB().ExecContext(r.Context(), query, form.Name, form.Email, hash)
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
					formErrors.AddError("name", "This name is already taken")
					log.Printf("Duplicate username attempt: %s", form.Name)
				}

				// Re-render form with error
				data := map[string]any{
					"Errors": formErrors,
					"Form":   &form,
				}
				tmpl.ExecuteTemplate(w, "register_form", data)
				return
			}

			log.Printf("Error creating user: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("HX-Redirect", "/login")
		w.WriteHeader(http.StatusOK)
	}
}
