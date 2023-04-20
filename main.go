package main

import (
	"APIGateway/pkg/api"
	"APIGateway/pkg/storage/postgres"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	user := os.Getenv("POSTGRES_USER")
	pwd := os.Getenv("POSTGRES_PASSWORD")
	dbService := os.Getenv("POSTGRES_DB_SERVICE")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbName := os.Getenv("POSTGRES_DB")

	if user == "" || pwd == "" || dbService == "" || dbPort == "" || dbName == "" {
		os.Exit(1)
	}

	connstr := "postgres://" + user + ":" + pwd + "@" + dbService + ":" + dbPort + "/" + dbName
	// Чит, чтобы дождаться инициализации БД
	time.Sleep(5 * time.Second)
	db, err := postgres.New(connstr)
	if err != nil {
		log.Fatal(err)
	}
	api := api.New(db)

	err = http.ListenAndServe(":8080", api.Router())
	if err != nil {
		log.Fatal(err)
	}
}
