package middleware

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/wrytehq/wryte/internal/database"
)

type contextKey string

const UserIDKey contextKey = "userID"

type SessionInfo struct {
	UserID string
	Token  string
}

func GetSession(r *http.Request, db database.Service) (*SessionInfo, error) {
	cookie, err := r.Cookie("wryte_session")
	if err != nil {
		return nil, err
	}

	var userID string
	var expiresAt time.Time
	query := `SELECT user_id, expires_at FROM sessions WHERE token = $1`
	err = db.GetDB().QueryRowContext(r.Context(), query, cookie.Value).Scan(&userID, &expiresAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invalid session")
		}
		return nil, err
	}

	if time.Now().After(expiresAt) {
		return nil, errors.New("session expired")
	}

	return &SessionInfo{
		UserID: userID,
		Token:  cookie.Value,
	}, nil
}

func Authenticated(db database.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := GetSession(r, db)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, session.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Guest(db database.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := GetSession(r, db)
			if err == nil {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func GetUserID(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(UserIDKey).(string)
	return userID, ok
}
