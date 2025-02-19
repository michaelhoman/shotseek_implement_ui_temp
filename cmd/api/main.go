package main

import (
	"time"

	"github.com/michaelhoman/ShotSeek/internal/env"
	"github.com/michaelhoman/ShotSeek/internal/postgres_db"
	"go.uber.org/zap"

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
	cfg := config{
		addr:   env.GetString("ADDR", ":8080"),
		apiURL: env.GetString("EXTERNAL_URL", "localhost:8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5432/shotseek?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			exp: time.Hour * 1, // 1 hour
		},
	}

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Database
	db, err := postgres_db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("Database connection pool established")

	store := store.NewPostgresStorage(db)

	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
	}

	mux := app.mount()

	logger.Fatal(app.run(mux))
}
