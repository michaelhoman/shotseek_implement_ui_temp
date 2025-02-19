package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/michaelhoman/ShotSeek/docs" // This is required to run Swagger Docs
	"github.com/michaelhoman/ShotSeek/internal/store"

	// store "github.com/michaelhoman/ShotSeek/internal/store/postgres"
	httpSwagger "github.com/swaggo/http-swagger" // http-swagger middleware
)

type application struct {
	config config
	store  store.Storage
	logger *zap.SugaredLogger
}

type config struct {
	addr   string
	db     dbConfig
	env    string
	apiURL string
	mail   mailConfig
}

type mailConfig struct {
	exp time.Duration
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL(docsURL), //The url pointing to API definition
		))
		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostsHandler)

			// Comments
			r.Route("/comments/{commentID}", func(r chi.Router) {
				r.Use(app.commentsContextMiddleware)
				r.Get("/", app.getCommentHandler)
				r.Patch("/", app.updateCommentHandler)
				r.Delete("/", app.DeleteByCommentIDHandler)
			})

			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware)
				r.Post("/comments", app.createCommentHandler)
				r.Get("/", app.getPostHandler)
				r.Get("/", app.getPostHandler)
				r.Patch("/", app.updatePostHandler)
				r.Delete("/", app.deletePostHandler)

			})
		})

		r.Route("/users", func(r chi.Router) {
			// r.Post("/", app.createUserHandler)
			r.Put("/activate/{token}", app.activateUserHandler)
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.usersContextMiddleware)
				r.Get("/", app.getUserHandler)
				r.Patch("/", app.updateUserHandler)
				r.Delete("/", app.deleteUserHandler)
			})
		})

		//public
		r.Route("/authentication", func(r chi.Router) {
			r.Post("/register", app.registerUserHandler)
			//r.Post("/login", app.loginHandler)
			//r.Post("/logout", app.logoutHandler)
		})

	})

	return r
}

func (app *application) run(mux http.Handler) error {
	//docs
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.apiURL
	docs.SwaggerInfo.BasePath = "/v1"

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	app.logger.Info("Server has started at ", "ADDR", app.config.addr, "ENV", app.config.env)
	return srv.ListenAndServe()
}
