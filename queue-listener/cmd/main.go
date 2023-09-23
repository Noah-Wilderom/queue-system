package main

import (
	"database/sql"
	"github.com/Noah-Wilderom/queue-system/queue-listener/data"
	"log"
	"os"
	"sync"
	"time"
)

const (
	gRPC = 5002
)

var (
	counts int64
)

type Config struct {
}

func main() {
	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres")
	}

	data.SetConnection(conn)

	app := Config{}

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		app.gRPCListen()
	}()

	go func() {
		defer wg.Done()
		app.ListenQueue()
	}()

	wg.Wait()
}

func (app *Config) ListenQueue() {
	//
}

func openDBConnection(dsn string) (*sql.DB, error) {
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
		connection, err := openDBConnection(dsn)
		if err != nil {
			log.Println()
			log.Println("Postgres not yet ready...", err)
			counts++
		} else {
			log.Println("Connected to Postgres")
			return connection
		}

		if counts > 10 {
			log.Println(err)

			return nil
		}

		log.Println("Backing off for two seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}
