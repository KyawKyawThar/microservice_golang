package main

import (
	"fmt"
	"log"
	"net/http"
)

var webPort = "80"

type Config struct {
}

func main() {
	log.Printf("Starting broker service on port :%s\n", webPort)

	app := Config{}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()

	if err != nil {
		log.Panicln(err)
	}
}
