package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"listener-service/event"
	"log"
	"math"
	"os"
	"time"
)

func main() {

	//try to connect to rabbitmq

	rabbitConn, err := connect()

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	defer rabbitConn.Close()

	//Start listening for service
	fmt.Println("Listening for and consuming RabbitMQ MESSAGE....")

	//create consumer
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		return
	}

	//watch the queue and consume events
	err = consumer.Listener([]string{"log.INFO", "log.WARNING", "log.ERROR"})

	if err != nil {
		log.Println(err)
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
