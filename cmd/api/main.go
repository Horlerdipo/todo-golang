package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/horlerdipo/todo-golang/env"
	"github.com/horlerdipo/todo-golang/internal/app"
	"github.com/horlerdipo/todo-golang/internal/database"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	_ "modernc.org/sqlite"
	"net/http"
	"time"
)

func main() {
	env.LoadEnv(".env")
	db, err := gorm.Open(sqlite.Dialector{
		DriverName: "sqlite", // <-- must match the imported driver
		DSN:        "test.db",
	}, &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Migrate models
	err = db.AutoMigrate(
		&database.User{},
		&database.TokenBlacklist{},
		&database.Todo{},
		&database.Checklist{},
	)
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	//r.Use(middleware.SupressNotFound(r))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CleanPath)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Heartbeat("/ping"))
	//r.Use(middleware.Timeout(60 * time.Second))
	rateLimitMins := env.FetchInt("HTTP_RATE_LIMIT_MINUTES", 60)
	if rateLimitMins > 0 {
		r.Use(httprate.LimitByIP(rateLimitMins, time.Minute))
	}

	appContainer := app.NewAppContainer(db)
	appContainer.RegisterRoutes(r)
	appContainer.RegisterListeners()

	port := env.FetchString("PORT", ":8000")
	log.Println("🚀🚀🚀 Starting server on port " + port)
	err = http.ListenAndServe(port, r)

	if err != nil {
		panic(err.Error())
	}
}
