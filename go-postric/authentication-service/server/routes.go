package main

import (
	"authentication/handlers"
	"authentication/util"
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func (app *Config) routes() http.Handler {
	r := chi.NewRouter()

	// Middleware stack
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Heartbeat("/health"))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// Set a max timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(30 * time.Second))

	r.Get("/", util.RouteHandler(handlers.Home))

	r.Get("/ping", util.RouteHandler(handlers.Ping))

	// API version 1
	r.Mount("/v1", v1Router())

	return r
}

func v1Router() chi.Router {
	r := chi.NewRouter()
	r.Use(setCtx("api.version", "v1"))

	r.Post("/signin", util.RouteHandler(handlers.Signin))

	r.Group(func(r chi.Router) {
		r.Use(handlers.AuthGuard())

		r.Get("/welcome", util.RouteHandler(handlers.Welcome))

		r.Get("/refresh", util.RouteHandler(handlers.Refresh))
	})

	return r
}

func setCtx(key, val any) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), key, val))
			next.ServeHTTP(w, r)
		})
	}
}
