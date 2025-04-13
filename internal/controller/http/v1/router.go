package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/sudeeya/avito-assignment/internal/service"
)

func NewRouter(services *service.Services) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/api/v1", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
			w.Write([]byte("API v1 is running"))
		})
		r.Post("/dummyLogin", dummyLoginHandler(services.Auth))

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware(services.Auth))

			r.Mount("/pvz", newPVZRouter(services))
			r.Mount("/receptions", newReceptionsRouter(services))
			r.Mount("/products", newProductsRouter(services))
		})
	})

	return router
}
