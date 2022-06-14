package main

import (
	"authentication/handlers/auth"
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func (app *Config) routes() http.Handler {
	r := chi.NewRouter()

	// Middleware stack
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
	r.Use(middleware.Timeout(app.DefaultTimeout))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("."))
	})

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		log.Panic("test")
	})

	// API version 1
	r.Mount("/v1", v1Router())

	return r
}

func v1Router() chi.Router {
	r := chi.NewRouter()
	r.Use(setCtx("api.version", "v1"))

	r.Mount("/auth", authRouter())

	return r
}

func authRouter() chi.Router {
	r := chi.NewRouter()

	r.Post("/signin", auth.Signin)

	r.Get("/welcome", auth.Welcome)

	r.Get("/refresh", auth.Refresh)

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
