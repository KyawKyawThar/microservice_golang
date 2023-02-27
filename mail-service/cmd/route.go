package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"log"
	"net/http"
)

func (app *Config) routes() http.Handler {

	log.Println("mail-service route hits")
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},

		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	//Heartbeat endpoint middleware useful to setting up a path like /ping
	//that load balancers or uptime testing external services can make a request
	//before hitting any routes. It's also convenient to place this above
	//ACL middlewares as well.

	mux.Use(middleware.Heartbeat("/ping"))
	mux.Post("/send", app.SendMail)
	return mux
}
