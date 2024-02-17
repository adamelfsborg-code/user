package server

import (
	"net/http"

	"github.com/adamelfsborg-code/food/user/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (a *Server) loadRoutes() {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/api/v1/users", a.loadUserRoutes)
	a.router = router
}

func (a *Server) loadUserRoutes(router chi.Router) {
	userHandler := &handler.UserHandler{
		Data: a.data,
	}

	router.Post("/register", userHandler.Register)
	router.Post("/login", userHandler.Login)

	router.Group(func(r chi.Router) {
		r.Use(CustomAuthMiddleware())
		r.Get("/ping", userHandler.Ping)
	})
}
