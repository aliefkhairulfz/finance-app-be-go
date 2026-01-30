package middlewares

import (
	"context"
	"fmt"
	"net/http"
)

type CsrfToken struct {
	Token string
}

func (c *CsrfToken) GetToken() string {
	return c.Token
}

func (c *CsrfToken) SetToken(token string) {
	c.Token = token
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		csrfToken := r.Header.Get("X-CSRF-Token")
		if csrfToken == "" {
			fmt.Println("csrf token is missing")
			http.Error(w, "CSRF token is missing", http.StatusUnauthorized)
			return
		}

		cookie, err := r.Cookie("session")
		if err != nil {
			fmt.Println("token session is missing")
			http.Error(w, "Session cookie is missing", http.StatusUnauthorized)
			return
		}

		token := cookie.Value
		if token == "" {
			http.Error(w, "Token is missing", http.StatusUnauthorized)
			return
		}

		fmt.Println("csrf-token is: ", csrfToken)
		fmt.Println("session-token is: ", cookie)
		ctx := context.WithValue(r.Context(), "token", token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
