package main

import (
	"log"

	"github.com/michaelhoman/ShotSeek/internal/env"
	"github.com/michaelhoman/ShotSeek/internal/store/postgres"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
	}
	store := postgres.NewPostgresStorage(nil)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
