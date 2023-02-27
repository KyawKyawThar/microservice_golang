package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var webPort = "80"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	log.Printf("Starting broker service on port :%s\n", webPort)

	rabbitConn, err := connect()
	app := Config{

		Rabbit: rabbitConn,
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()

	if err != nil {
		log.Panicln(err)
	}
}

func connect() (*amqp.Connection, error) {
	var count int64
	var backoff = 1 * time.Second
	var connection *amqp.Connection

	//don't continue until rabbitmq is readu

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")

		if err != nil {
			log.Println("Rabbit mq is not ready yet...")
			count++
		} else {
			connection = c
			log.Println("Connect to rabbit....")
			break

		}

		if count > 5 {
			fmt.Println(err)
			return nil, err
		}

		backoff = time.Duration(math.Pow(float64(count), 2)) * time.Second

		log.Println("backing off.....")
		time.Sleep(backoff)

	}

	return connection, nil
}
