package routers

import (
	"net/http"
	"xeed/apps/cp-api/internal/http/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func InitRouter(userHandler *handlers.UserHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/users/register", userHandler.Register)
		r.Post("/auth/login", userHandler.Login)
	})

	return r
}
