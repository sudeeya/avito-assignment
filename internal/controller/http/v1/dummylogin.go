package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sudeeya/avito-assignment/internal/service"
)

// User roles.
const (
	_moderator = "moderator"
)

func dummyLoginHandler(authService service.Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user struct {
			Role string `json:"role"`
		}
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if user.Role != _moderator {
			http.Error(w, "invalid role", http.StatusBadRequest)
			return
		}

		token, err := authService.IssueToken(r.Context())
		if err != nil {
			fmt.Println(err)
			http.Error(w, "issuing token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(token))
	}
}
