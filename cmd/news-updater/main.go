package main

import (
	"APIGateway/pkg/rss"
	"APIGateway/pkg/rss/configReader"
	"APIGateway/pkg/storage"
	"APIGateway/pkg/storage/postgres"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const configFile = "config.json"

func main() {
	// Чит, чтобы дождаться инициализации БД
	time.Sleep(5 * time.Second)

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
	st, err := postgres.New(connstr)
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan storage.News, 50)
	defer close(ch)

	runNewsGetting(ch)
	saveNews(st, ch)
}

func runNewsGetting(ch chan storage.News) {
	rssConfig := configReader.Read(configFile)
	for _, url := range rssConfig.Rss {
		go rss.UpdateNews(url, rssConfig.RequestPeriod, ch)
	}
}

func saveNews(st *postgres.DbStorage, ch <-chan storage.News) {
	for news := range ch {
		st.AddNews(news)
	}
}
