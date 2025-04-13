package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sudeeya/avito-assignment/internal/service"
	"go.uber.org/zap"
)

func newReceptionsRouter(services *service.Services) *chi.Mux {
	router := chi.NewRouter()

	router.Post("/", createReceptionHandler(services.Reception))

	return router
}

type createReceptionInput struct {
	PVZID uuid.UUID `json:"pvz_id"`
}

func createReceptionHandler(receptionService service.Reception) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input createReceptionInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		reception, err := receptionService.CreateReception(r.Context(), input.PVZID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(reception); err != nil {
			zap.S().Errorf("encoding pvz: %v", err)
		}
	}
}
