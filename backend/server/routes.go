package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

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
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// Set a max timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(app.DefaultTimeout))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("."))
	})

	// API version 1
	r.Mount("/v1", v1Router())

	return r
}

func v1Router() chi.Router {
	r := chi.NewRouter()
	r.Use(apiVersionCtx("v1"))

	// RESTy routes for "articles" resource
	r.Mount("/articles", articleRouter())

	// Mount the admin sub-router, which btw is the same as:
	// r.Route("/admin", func(r chi.Router) { admin routes here })
	r.Mount("/admins", adminRouter())

	// Slow handlers/operations
	r.Group(func(r chi.Router) {
		// Stop processing after 2.5 seconds
		r.Use(middleware.Timeout(2500 * time.Millisecond))

		r.Get("/slow", func(w http.ResponseWriter, r *http.Request) {
			rand.Seed(time.Now().Unix())

			// Processing will take 1-5 seconds
			processTime := time.Duration(rand.Intn(4)+1) * time.Second
			log.Println("Process time:", processTime)

			select {
			case <-r.Context().Done():
				log.Println("Cancelled by chi middleware")
				return
			case <-time.After(processTime):
				// The above channel simulates some hard work
			}

			w.Write([]byte(fmt.Sprintf("Processed in %v seconds\n", processTime)))
		})
	})

	// Throttle very expensive handlers/operations
	r.Group(func(r chi.Router) {
		// Stop processing after 30 seconds
		// We already have a default time in app config
		// r.Use(middleware.Timeout(30 * time.Second))

		// Only one request will be processed at a time
		r.Use(middleware.Throttle(1))

		r.Get("/throttled", func(w http.ResponseWriter, r *http.Request) {
			select {
			case <-r.Context().Done():
				switch r.Context().Err() {
				case context.DeadlineExceeded:
					w.WriteHeader(504)
					w.Write([]byte("Processing too slow\n"))
				default:
					w.Write([]byte("Cancelled\n"))
				}
				return

			case <-time.After(5 * time.Second):
				// The above channel simulates some hard work
			}

			w.Write([]byte("Processed\n"))
		})
	})

	return r
}

// AdminOnly middleware restricts access to just administrators.
func adminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAdmin, ok := r.Context().Value("acl.admin").(bool)
		if !ok || !isAdmin {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func articleRouter() chi.Router {
	r := chi.NewRouter()
	// r.With(paginate).Get("/", ListArticles)
	// r.Post("/", CreateArticle)       // POST /articles
	// r.Get("/search", SearchArticles) // GET /articles/search

	// r.Route("/{articleID}", func(r chi.Router) {
	// 	r.Use(ArticleCtx)            // Load the *Article on the request context
	// 	r.Get("/", GetArticle)       // GET /articles/123
	// 	r.Put("/", UpdateArticle)    // PUT /articles/123
	// 	r.Delete("/", DeleteArticle) // DELETE /articles/123
	// })

	// GET /articles/whats-up
	// r.With(ArticleCtx).Get("/{articleSlug:[a-z-]+}", GetArticle)

	return r
}

// A completely separate router for administrator routes
func adminRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(adminOnly)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("admin: index"))
	})
	r.Get("/accounts", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("admin: list accounts.."))
	})
	r.Get("/users/{userId}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("admin: view user id %v", chi.URLParam(r, "userId"))))
	})
	return r
}

// paginate is a stub, but very possible to implement middleware logic
// to handle the request params for handling a paginated request.
func Paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// just a stub.. some ideas are to look at URL query params for something like
		// the page number, or the limit, and send a query cursor down the chain
		next.ServeHTTP(w, r)
	})
}

func apiVersionCtx(version string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), "api.version", version))
			next.ServeHTTP(w, r)
		})
	}
}
