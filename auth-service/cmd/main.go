package main

import (
	"auth-service/data"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"net/http"
	"os"
	"time"
)

var webPort = "80"

var count int64

type Config struct {
	DB    *sql.DB
	Model data.Models
}

func main() {

	log.Printf("Starting authentication service: %s", webPort)

	//Connect DB

	conn := ConnectDB()

	if conn == nil {
		log.Println("Can't connect to Postgres....")
	}

	app := Config{
		DB:    conn,
		Model: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()

	if err != nil {
		log.Panicln(err)
	}
}

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	log.Println("Successfully ping DB")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func ConnectDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		conn, err := OpenDB(dsn)

		if err != nil {
			log.Println("Postgres not yet ready ...")
			count++
		} else {
			log.Println("Connected to Postgres!")
			return conn
		}

		if count > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds....")
		time.Sleep(time.Second * 2)
		continue
	}

}
