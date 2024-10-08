package handler

import (
	_ "github.com/DSorbon/effective-mobile-task/docs"
	"github.com/DSorbon/effective-mobile-task/internal/service"
	mw "github.com/DSorbon/effective-mobile-task/internal/transport/http/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	songService service.Song
}

func NewHandler(songService service.Song) *Handler {
	return &Handler{
		songService: songService,
	}
}

func (h *Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(mw.LogMiddleware)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
	))

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/songs", func(r chi.Router) {
			r.Get("/", h.List)
			r.Get("/{id}", h.Get)
			r.Post("/", h.Create)
			r.Patch("/{id}", h.Update)
			r.Delete("/{id}", h.Delete)
		})
	})

	return r
}
