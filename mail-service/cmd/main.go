package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

var webPort = "80"

type Config struct {
	Mailer Mail
}

func main() {
	log.Printf("Starting broker service on port :%s\n", webPort)

	app := Config{Mailer: createMail()}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()

	if err != nil {
		log.Panicln(err)
	}
}
func createMail() Mail {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))

	m := Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromAddress: os.Getenv("FROM_ADDRESS"),
		FromName:    os.Getenv("FROM_NAME"),
	}

	fmt.Printf("create Mail %v\n", m)
	return m
}
