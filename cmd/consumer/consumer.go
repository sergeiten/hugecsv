package main

import (
	"os"
	"strconv"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/sergeiten/hugecsv"
	"github.com/sergeiten/hugecsv/consumer"
)

func main() {
	cfg := config()
	consumer, err := consumer.New(cfg)
	hugecsv.LogFatal(err, "failed to create consumer instance")

	consumer.Serve()
}

func config() *consumer.Config {
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		hugecsv.LogFatal(err, "failed to convert database port")
	}

	return &consumer.Config{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME"),
		Port:     port,
	}
}
