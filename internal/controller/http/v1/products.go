package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/sudeeya/avito-assignment/internal/service"
)

func newProductsRouter(services *service.Services) *chi.Mux {
	router := chi.NewRouter()

	router.Post("/", addProductHandler(services.Product))

	return router
}

type addProductInput struct {
	Type  string    `json:"type"`
	PVZID uuid.UUID `json:"pvz_id"`
}

func addProductHandler(productService service.Product) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input addProductInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		product, err := productService.AddProduct(r.Context(), input.PVZID, input.Type)
		if errors.Is(err, service.ErrUnsupportedProductType) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(product); err != nil {
			zap.S().Errorf("encoding pvz: %v", err)
		}
	}
}
