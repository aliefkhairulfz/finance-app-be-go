package utils

import "net/http"

func SetCookie(w http.ResponseWriter, token string) {
	cookie := http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
}

func SetCsrfCookie(w http.ResponseWriter, token string) {
	cookie := http.Cookie{
		Name:     "csrf",
		Value:    token,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
}
