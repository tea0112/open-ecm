package api

import (
	"database/sql"
	"fmt"
	"net/http"
	v1api "open-ecm/api/v1"
	"open-ecm/apps"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func NewChiMainRouter(db *sql.DB, logger *zap.Logger) *chi.Mux {
	r := chi.NewRouter()

	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(apps.CtxWithLogger(r.Context(), logger))
			h.ServeHTTP(w, r)
		})
	})
	r.Use(middleware.Recoverer)
	r.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info(fmt.Sprintf("%s - %s%s", r.Method, r.Host, r.URL))
			handler.ServeHTTP(w, r)
		})
	})

	r.Mount("/api/v1", v1api.NewUsersRouter(db))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("This URL doesn't belong to the OpenECM"))
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
	})

	return r
}
