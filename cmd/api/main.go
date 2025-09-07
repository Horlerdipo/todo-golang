package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/horlerdipo/todo-golang/env"
	"github.com/horlerdipo/todo-golang/internal/app"
	"github.com/horlerdipo/todo-golang/internal/users"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	_ "modernc.org/sqlite"
	"net/http"
	"time"
)

func main() {

	db, err := gorm.Open(sqlite.Dialector{
		DriverName: "sqlite", // <-- must match the imported driver
		DSN:        "test.db",
	}, &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Migrate models
	err = db.AutoMigrate(&users.User{})
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{\"message\": \"pong\"}"))
	})

	appContainer := app.NewAppContainer(db)
	appContainer.AuthContainer.RegisterRoutes(r)

	port := env.FetchString("PORT", ":8000")
	log.Println("🚀🚀🚀 Starting server on port " + port)

	err = http.ListenAndServe(port, r)

	if err != nil {
		panic(err.Error())
	}
}
