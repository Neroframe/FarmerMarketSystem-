package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"time"
)

func SetCSRFToken(w http.ResponseWriter) (string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}
	token := base64.StdEncoding.EncodeToString(tokenBytes)

	cookie := &http.Cookie{
		Name:     "csrf_token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true, // Prevents JavaScript access
		Secure:   true,
		SameSite: http.SameSiteNoneMode, // Allows cross-site cookie
		// SameSite: http.SameSiteStrictMode,   // Prevents CSRF attacks
	}
	http.SetCookie(w, cookie)
	return token, nil
}

// Compares the CSRF tokens (form with cookie)
func ValidateCSRFToken(r *http.Request) error {
	formToken := r.FormValue("csrf_token")
	if formToken == "" {
		return errors.New("csrf token not provided")
	}

	cookie, err := r.Cookie("csrf_token")
	if err != nil {
		return errors.New("csrf token cookie not found")
	}

	if formToken != cookie.Value {
		return errors.New("invalid csrf token")
	}

	return nil
}
