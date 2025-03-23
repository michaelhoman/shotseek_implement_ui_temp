package main

import (
	"fmt"
	"log"

	"github.com/michaelhoman/ShotSeek/internal/auth"
	"github.com/michaelhoman/ShotSeek/internal/config"
	"github.com/michaelhoman/ShotSeek/internal/postgres_db"
	"github.com/michaelhoman/ShotSeek/internal/utils"

	// postgres_store "github.com/michaelhoman/ShotSeek/internal/store/postgres"
	"github.com/michaelhoman/ShotSeek/internal/store"
)

const version = "0.0.1"

//	@title			ShotSeek API
//	@description	This is the API for ShotSeek Cinematographer Finder
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	https://www.hintproductions.com
//	@contact.email	homanstudio@proton.me

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath					/v1
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description

func main() {
	// cfg := config{
	// 	addr:   env.GetString("ADDR", ":8080"),
	// 	apiURL: env.GetString("EXTERNAL_URL", "localhost:8080"),
	// 	db: dbConfig{
	// 		addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5432/shotseek?sslmode=disable"),
	// 		maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
	// 		maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
	// 		maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
	// 	},
	// 	env: env.GetString("ENV", "development"),
	// 	mail: mailConfig{
	// 		exp: time.Hour * 1, // 1 hour
	// 	},
	// 	auth: authConfig{
	// 		basic: basicConfig{
	// 			user: env.GetString("AUTH_BASIC_USER", "admin"),
	// 			pass: env.GetString("AUTH_BASIC_PASS", "admin"),
	// 		},
	// 		token: tokenConfig{
	// 			secret: env.GetString("JWT_SIGNING_KEY", "example"),
	// 			exp:    time.Minute * 6, // 6 minutes
	// 			iss:    "shotseek-auth-service",
	// 			aud:    "shotseek-api",
	// 		},
	// 	},
	// }

	// Load configuration
	cfg := config.Load()

	// Initialize logger
	utils.InitLogger()
	defer utils.CleanupLogger()

	// Access logger
	logger := utils.Logger
	// Database
	db, err := postgres_db.New(
		cfg.Db.Addr,
		cfg.Db.MaxOpenConns,
		cfg.Db.MaxIdleConns,
		cfg.Db.MaxIdleTime,
	)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("Database connection pool established")

	store := store.NewPostgresStorage(db)

	jwtService := auth.NewJWTService(cfg.Auth.Token.Secret, cfg.Auth.Token.Exp)

	// authHandler := auth.NewAuthHandler(store, cfg, jwtService)
	// Initialize JWTAuth with the ECDSA keys
	jwtAuth, err := auth.NewJWTAuth()
	if err != nil {
		log.Fatalf("Error initializing JWTAuth: %v", err)
	}

	// Example usage of the jwtAuth instance
	fmt.Println("JWT Auth initialized:", jwtAuth)

	app := &application{
		config:     cfg,
		store:      store,
		jwtService: jwtService,
		jwtAuth:    jwtAuth,
		auth:       auth.NewAuthHandler(store, cfg, jwtService, jwtAuth),
	}

	mux := app.mount()

	logger.Fatal(app.run(mux))
}
