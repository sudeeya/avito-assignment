package v1

import (
	"net/http"
	"strings"

	"github.com/sudeeya/avito-assignment/internal/service"
)

const (
	_authorizationHeader = "Authorization"
	_bearer              = "Bearer "
)

func authMiddleware(authService service.Auth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		h := func(w http.ResponseWriter, r *http.Request) {
			bearerString := r.Header.Get(_authorizationHeader)

			if bearerString == "" {
				http.Error(w, "empty bearer string", http.StatusForbidden)
				return
			}

			tokenString := strings.TrimPrefix(bearerString, _bearer)

			err := authService.VerifyToken(r.Context(), tokenString)
			if err != nil {
				http.Error(w, "wrong token", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(h)
	}
}
