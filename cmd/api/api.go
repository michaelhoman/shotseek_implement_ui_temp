package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/michaelhoman/ShotSeek/cmd/ui"
	"github.com/michaelhoman/ShotSeek/docs" // This is required to run Swagger Docs
	"github.com/michaelhoman/ShotSeek/internal/auth"
	"github.com/michaelhoman/ShotSeek/internal/config"
	"github.com/michaelhoman/ShotSeek/internal/mailer"
	int_middleware "github.com/michaelhoman/ShotSeek/internal/middleware"
	"github.com/michaelhoman/ShotSeek/internal/store"
	"github.com/michaelhoman/ShotSeek/internal/utils"

	// store "github.com/michaelhoman/ShotSeek/internal/store/postgres"
	httpSwagger "github.com/swaggo/http-swagger" // http-swagger middleware
)

type application struct {
	config     config.Config
	store      store.Storage
	mailer     mailer.Client
	jwtService *auth.JWTService
	jwtAuth    *auth.JWTAuth
}

//	type config struct {
//		addr   string
//		db     dbConfig
//		env    string
//		apiURL string
//		mail   mailConfig
//		auth   authConfig
//	}
type authConfig struct {
	basic basicConfig
	token tokenConfig
}
type tokenConfig struct {
	secret string
	exp    time.Duration
	iss    string
	aud    string
}

type basicConfig struct {
	user string
	pass string
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
	authHandler := auth.NewAuthHandler(app.store, app.config, app.jwtService, app.jwtAuth)
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

	// Initialize JWT service
	// jwtService := auth.NewJWTService(app.config.auth.token.secret, app.config.auth.token.exp)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.Addr)
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
			//r.Post("/", app.createUserHandler)
			r.Use(int_middleware.JwtMiddleware(authHandler))
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
			r.Post("/register", authHandler.RegisterUserHandler)
			r.Post("/login", authHandler.LoginHandler)
			r.Post("/logout", authHandler.LogoutHandler)
			//r.Post("/logout", app.logoutHandler)
		})

	})

	ui.RegisterUIRoutes(r) // Call the function and pass its return value

	return r
}

func (app *application) run(mux http.Handler) error {
	//docs
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.ApiURL
	docs.SwaggerInfo.BasePath = "/v1"

	srv := &http.Server{
		Addr:         app.config.Addr,
		Handler:      mux,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	utils.Logger.Info("Server has started at ", "ADDR", app.config.Addr, "ENV", app.config.Env)
	return srv.ListenAndServe()
}
