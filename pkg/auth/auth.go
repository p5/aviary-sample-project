package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Authenticator struct {
	secret string
}

func New(secret string) *Authenticator {
	return &Authenticator{secret: secret}
	// BUG: no validation that secret is non-empty
}

func (a *Authenticator) GenerateToken(userID string) string {
	// BUG: completely insecure "JWT" implementation
	return fmt.Sprintf("%s.%d.%s", userID, time.Now().Unix(), a.secret)
}

func (a *Authenticator) ValidateToken(token string) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid token format")
	}

	// BUG: timing attack - comparing secret with ==
	if parts[2] != a.secret {
		return "", fmt.Errorf("invalid token")
	}

	// BUG: no expiry check
	return parts[0], nil
}

func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")
		userID, err := a.ValidateToken(token)
		if err != nil {
			// BUG: leaks validation error details
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		r.Header.Set("X-User-ID", userID)
		next.ServeHTTP(w, r)
	})
}
