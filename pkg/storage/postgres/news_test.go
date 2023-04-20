package postgres

import (
	"APIGateway/pkg/storage"
	"errors"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"testing"
)

func TestPostgres_New(t *testing.T) {
	_, err := connect()
	if err != nil {
		log.Fatal(err)
	}
}

func TestPostgres_AddNews(t *testing.T) {
	news := storage.NewNews(
		"guid",
		"title",
		"content",
		111,
		"url.ru",
	)

	st, err := connect()
	if err != nil {
		t.Fatal(err)
	}
	err = st.AddNews(news)
	if err != nil {
		log.Fatal(err)
	}
}

func TestPostgres_News(t *testing.T) {
	st, err := connect()
	if err != nil {
		t.Fatal(err)
	}

	for i := 1; i < 10; i++ {
		news := storage.NewNews(
			"guid"+strconv.Itoa(i),
			"title",
			"content",
			111,
			"url.ru",
		)
		err = st.AddNews(news)
		if err != nil {
			t.Fatal(err)
		}
	}
	news, err := st.News(10)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", news)
}

func connect() (*DbStorage, error) {
	err := godotenv.Load("./../../../.env")
	if err != nil {
		log.Fatal(err)
	}

	user := os.Getenv("POSTGRES_USER")
	pwd := os.Getenv("POSTGRES_PASSWORD")
	dbService := "localhost"
	dbPort := os.Getenv("POSTGRES_PORT")
	dbName := os.Getenv("POSTGRES_DB")

	if user == "" || pwd == "" || dbService == "" || dbPort == "" || dbName == "" {
		return nil, errors.New("Empty environment variables")
	}

	connstr := "postgres://" + user + ":" + pwd +
		"@" + dbService + ":" + dbPort + "/" + dbName
	st, err := New(connstr)

	return st, err
}
