package main

import (
	"authentication/data"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/jackc/pgx/v4"
)

const webPort = "80"

var counts int64

type Config struct {
	DB *sql.DB
	Models data.Models
}

func main() {

	log.Println("Starting authentication service...")

	// conect to db
	conn := connectToDB()
	if conn == nil {
		log.Panic("Could not connect to database")
	}

	// setup config
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}


	srv := &http.Server{
		Addr:    ":" + webPort,
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Error connecting to database:", err)
			counts++
		} else {
			log.Println("Connected to database")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Retrying in 5 seconds...")
		// wait for 5 seconds before retrying
		time.Sleep(5 * time.Second)
		log.Println("Retrying connection to database...")
		continue
	}
}