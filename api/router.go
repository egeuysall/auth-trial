package api

import (
	"time"

	appmid "github.com/egeuysall/auth-trial/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(
		middleware.Recoverer,
		middleware.RealIP,
		middleware.Timeout(10*time.Second),
		middleware.NoCache,
		middleware.Compress(5),
		httprate.LimitByIP(3, 3*time.Second),
		appmid.SetContentType(),
		appmid.Cors(),
	)

	r.Get("/", HandleRoot)
	r.Get("/ping", CheckPing)

	r.Route("/v1", func(r chi.Router) {
		r.Post("/login", Login)
		r.Post("/users", CreateUser)
	})

	return r
}
