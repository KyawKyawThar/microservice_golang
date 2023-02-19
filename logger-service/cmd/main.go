package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"logger-serviice/data"
	"net/http"
	"time"
)

var (
	webPort  = "80"
	rpcPort  = "5001"
	gRpcPort = "50001"
	mongoURL = "mongodb://mongo:27017"
)

var client *mongo.Client

type Config struct {
	Model data.Models
}

//mongodb://admin:password@localhost:27018

func main() {

	log.Printf("Starting Logger service on port :%s\n", webPort)

	mongoClient, err := connectToMongo()

	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	//create context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Panic(err)
		}
	}()

	app := Config{
		Model: data.New(client),
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

func connectToMongo() (*mongo.Client, error) {

	clientOptions := options.Client().ApplyURI(mongoURL)

	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	//connect mongo
	con, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Println("Error in connection", err)
		return con, err
	}
	log.Println("connect to Mongo")
	return con, nil
}
