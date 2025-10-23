package flash

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
)

type Type string

const (
	Success Type = "success"
	Error   Type = "error"
	Warning Type = "warning"
	Info    Type = "info"
)

type Message struct {
	Type    Type   `json:"type"`
	Content string `json:"content"`
}

const cookieName = "wryte_flash"

// Set adds a flash message that will be displayed on the next page load
func Set(w http.ResponseWriter, flashType Type, message string) error {
	msg := Message{
		Type:    flashType,
		Content: message,
	}

	// Encode message to JSON
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// Encode to base64 to avoid cookie issues with special characters
	encoded := base64.StdEncoding.EncodeToString(jsonData)

	// Set cookie
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    encoded,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   300, // 5 minutes
	}

	http.SetCookie(w, cookie)
	return nil
}

// Get retrieves and removes the flash message from the cookie
func Get(w http.ResponseWriter, r *http.Request) (*Message, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		// No flash message, not an error
		return nil, nil
	}

	// Delete the cookie immediately
	deleteCookie := &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	}
	http.SetCookie(w, deleteCookie)

	// Decode base64
	decoded, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return nil, err
	}

	// Decode JSON
	var msg Message
	if err := json.Unmarshal(decoded, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

// Helper functions for common flash message types

func SetSuccess(w http.ResponseWriter, message string) error {
	return Set(w, Success, message)
}

func SetError(w http.ResponseWriter, message string) error {
	return Set(w, Error, message)
}

func SetWarning(w http.ResponseWriter, message string) error {
	return Set(w, Warning, message)
}

func SetInfo(w http.ResponseWriter, message string) error {
	return Set(w, Info, message)
}
