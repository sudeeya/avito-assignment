package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sudeeya/avito-assignment/internal/service"
	"go.uber.org/zap"
)

func newPVZRouter(services *service.Services) *chi.Mux {
	router := chi.NewRouter()

	router.Get("/", getPVZPaginationHandler(services.PVZ))
	router.Post("/", createPVZHandler(services.PVZ))
	router.Post("/{pvzID}/close_last_reception", closeLastReceptionHandler(services.Reception))
	router.Post("/{pvzID}/delete_last_product", deleteLastProductHandler(services.Product))

	return router
}

type createPVZInput struct {
	ID               uuid.UUID `json:"id"`
	RegistrationDate time.Time `json:"registration_date"`
	City             string    `json:"city"`
}

func createPVZHandler(pvzService service.PVZ) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input createPVZInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		pvz, err := pvzService.CreatePVZ(r.Context(), input.City)
		if errors.Is(err, service.ErrUnsupportedCity) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(pvz); err != nil {
			zap.S().Errorf("encoding pvz: %v", err)
		}
	}
}

func getPVZPaginationHandler(pvzService service.PVZ) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()

		fmt.Println(params.Get("startDate"))

		start, err := time.Parse(time.DateOnly, params.Get("startDate"))
		if err != nil {
			http.Error(w, "invalid startDate", http.StatusBadRequest)
			return
		}

		end, err := time.Parse(time.DateOnly, params.Get("endDate"))
		if err != nil {
			http.Error(w, "invalid endDate", http.StatusBadRequest)
			return
		}

		limit, err := strconv.Atoi(params.Get("limit"))
		if err != nil {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}

		offset, err := strconv.Atoi(params.Get("page"))
		if err != nil {
			http.Error(w, "invalid page", http.StatusBadRequest)
			return
		}

		pvzs, err := pvzService.GetPVZPagination(r.Context(), start, end, limit, offset)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(pvzs); err != nil {
			zap.S().Errorf("encoding pvzs: %v", err)
		}
	}
}

func closeLastReceptionHandler(receptionService service.Reception) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pvzID, err := uuid.Parse(chi.URLParam(r, "pvzID"))
		if err != nil {
			http.Error(w, "invalid UUID", http.StatusBadRequest)
			return
		}

		reception, err := receptionService.CloseLastReception(r.Context(), pvzID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(reception); err != nil {
			zap.S().Errorf("encoding reception: %v", err)
		}
	}
}

func deleteLastProductHandler(productService service.Product) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pvzID, err := uuid.Parse(chi.URLParam(r, "pvzID"))
		if err != nil {
			http.Error(w, "invalid UUID", http.StatusBadRequest)
			return
		}

		err = productService.DeleteLastProduct(r.Context(), pvzID)
		if errors.Is(err, service.ErrReceptionIsEmpty) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
