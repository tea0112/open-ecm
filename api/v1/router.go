package v1api

import (
	"database/sql"
	"open-ecm/users"

	"github.com/go-chi/chi/v5"
)

func NewUsersRouter(db *sql.DB) *chi.Mux {
	v1ApiRouter := chi.NewRouter()

	v1ApiRouter.Route("/", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			userController := users.NewController(db)
			r.Get("/{userId}", userController.HandleGetUser)
			r.Get("/", userController.HandleGetUsers)
			r.Post("/", userController.HandleSaveUser)
			r.Delete("/{userId}", userController.HandleDeleteUser)
		})
	})

	return v1ApiRouter
}
