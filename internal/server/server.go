package server

import (
	"database/sql"
	"log"
	"time"

	"github.com/AsapDolly/EM_test/internal"
	"github.com/AsapDolly/EM_test/internal/handlers"
	"github.com/AsapDolly/EM_test/internal/storage"
	"github.com/go-chi/chi/v5"
)

// GetServer возвращает Chi сервер со всеми хэндлерами.
func GetServer(dbConnection *sql.DB) (internal.Config, *chi.Mux) {

	cfg := internal.GetConfig()

	h := handlers.Handler{}

	if cfg.DBAddress != "" {
		var err error

		for i := 1; i <= 5; i++ {
			dbConnection, err = sql.Open("postgres", cfg.DBAddress)
			if err == nil {
				break
			}
			time.Sleep(30 * time.Second)
		}

		if err != nil {
			log.Fatalf("unable to connect to database %v\n", cfg.DBAddress)
		}

		h.Storage = storage.GetNewConnection(dbConnection, cfg.DBAddress, "file://migrations/postgres")
	}

	r := chi.NewRouter()

	r.Route("/api/v1/person", func(r chi.Router) {

		r.Route("/delete", func(r chi.Router) {
			r.Delete("/{id}", h.DeleteData)
		})

		r.Route("/", func(r chi.Router) {
			r.Get("/", h.GetData)
			r.Post("/", h.SendData)
			r.Put("/", h.UpdateData)
		})

	})

	return cfg, r
}
